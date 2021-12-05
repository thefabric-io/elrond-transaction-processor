package processor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrDataSourceIsUndefined         = errors.New("ledger gatewayAPI for transaction processor is not defined")
	ErrStateStorageIsUndefined       = errors.New("state storage is undefined")
	ErrPastTransactionMustBePositive = errors.New("past transaction buffer must be positive")
	ErrLastNonceToProcessNotFound    = errors.New("last nonce to process is not found")
	ErrLastProcessedNonceNotFound    = errors.New("last processed nonce is not found")
)

const (
	networkResetNonceThreshold   Nonce = 10000
	defaultPastTransactionBuffer int   = 10
)

type OnTransactionReceivedFunc func(shard Shard, nonce Nonce, transactions []*Transaction, blockHash string)

var defaultTransactionProcessor = Processor{
	pastBlocksBuffer: defaultPastTransactionBuffer,
	waitForFinalizedCrossShardSmartContractResults: false,
	notifyEmptyBlocks:                    true,
	includeCrossShardStartedTransactions: false,
	onTransactionsReceivedFunc:           nil,
	verbose:                              false,
	internalState: &State{
		crossShardDictionary:        NewCrossShardDictionary(),
		lastProcessedNoncesInternal: NonceByShard{},
	},
}

func NewProcessor(opts ...Option) (*Processor, error) {
	p := defaultTransactionProcessor

	for _, opt := range opts {
		opt(&p)
	}

	if p.pastBlocksBuffer < 0 {
		p.pastBlocksBuffer = 1
	}

	if p.displayProgressBar && p.verbose {
		log.Printf("for readability, verbose cannot be displayed while progress bar is on, falling back to not displaying verbose")
		p.verbose = false
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return &p, nil
}

type Processor struct {
	dataSource                                     DataSource
	stateStorage                                   StateStorage
	startDate                                      time.Time
	shards                                         Shards
	onTransactionsReceivedFunc                     OnTransactionReceivedFunc
	pastBlocksBuffer                               int
	waitForFinalizedCrossShardSmartContractResults bool
	notifyEmptyBlocks                              bool
	includeCrossShardStartedTransactions           bool
	internalState                                  *State
	progressBar                                    *progressbar.ProgressBar
	displayProgressBar                             bool
	verbose                                        bool
}

func (p *Processor) Validate() error {
	if p.dataSource == nil {
		return ErrDataSourceIsUndefined
	}

	if p.stateStorage == nil {
		return ErrStateStorageIsUndefined
	}

	if p.pastBlocksBuffer < 0 {
		return ErrPastTransactionMustBePositive
	}

	return nil
}

func (p *Processor) logIfVerbose(s string) {
	if p.verbose {
		log.Printf(s)
	}
}

func (p *Processor) Start() (err error) {
	p.shards, err = p.dataSource.GetShards()
	if err != nil {
		panic(err)
	}

	p.internalState, err = p.stateStorage.FetchLastState(p.shards)
	if err != nil {
		return fmt.Errorf("could not fetch last state of processor: %w", err)
	}

	p.internalState.AddBufferToLastProcessNonces(p.pastBlocksBuffer)

	if p.internalState.toNonces == nil {
		p.internalState.toNonces, err = p.dataSource.GetCurrentNoncesForShards(p.shards)
		if err != nil {
			return fmt.Errorf("could not fetch current nonces for shards: %w", err)
		}
	}

	if err = p.Validate(); err != nil {
		return err
	}

	p.internalState.PruneCrossShardDictionary()

	p.startDate = time.Now()

	log.Printf("Targeted shards: %s\n\n", p.shards)

	if err = p.initProgressBar(p.internalState.NumberOfRemainingNonces()); err != nil {
		return err
	}

	defer p.end()

	var reachedTip bool

	for run := true; run; run = !reachedTip {
		reachedTip = true

		for _, shard := range p.shards {
			shardName := shard.Name()

			lastNonceToProcess, found := p.internalState.LastNonceToProcessInShard(shard)
			if !found {
				return ErrLastNonceToProcessNotFound
			}

			lastProcessedNonce, found := p.internalState.LastProcessedNonceInShard(shard)
			if !found {
				return ErrLastProcessedNonceNotFound
			}

			if lastProcessedNonce.Equals(lastNonceToProcess) {
				p.logIfVerbose(fmt.Sprintf("Nonce %d in %s has already been processed", lastNonceToProcess, shardName))

				continue
			}

			/*
				Handle the situation where the last nonce to process is reset.
				(e.g. devnet/testnet resets where the nonces start again from zero)
			*/
			if p.currentNonceIsReset(lastProcessedNonce, lastNonceToProcess) {
				p.logIfVerbose(fmt.Sprintf("Detected network reset. Setting last processed nonce to %d for %s\n", lastNonceToProcess, shardName))

				lastProcessedNonce = lastNonceToProcess.Decrement()
			}

			if lastProcessedNonce.IsGreaterThan(lastNonceToProcess) {
				p.logIfVerbose(fmt.Sprintf("The last processed nonce is superior to the current nonce"))

				continue
			}

			reachedTip = false

			nonce := lastProcessedNonce.Increment()

			if err := p.processValidTransactions(shard, nonce); err != nil {
				return err
			}

			p.logIfVerbose(fmt.Sprintf("Setting last processed nonce for %s to %d\n\n", shard.Name(), nonce))
			p.internalState.lastProcessedNoncesInternal.PutNonce(shard, nonce)

			p.incrementProgressBar()
		}
	}

	return nil
}

func (p *Processor) processValidTransactions(shard Shard, nonce Nonce) error {
	p.logIfVerbose(fmt.Sprintf("Begin transaction processing for nonce %d in %s\n", nonce, shard.Name()))

	blockHash, transactions, err := p.dataSource.GetShardTransactions(shard, nonce)
	if err != nil {
		log.Println(err)

		return err
	}

	validTransactions := make(Transactions, 0)

	if p.waitForFinalizedCrossShardSmartContractResults {
		finalizedTransactions := p.finalizedCrossShardScrTransactions(shard, transactions)
		for _, tx := range finalizedTransactions {
			validTransactions = append(validTransactions, tx)
		}
	}

	for _, tx := range transactions {
		// we only care about transactions that are finalized in the given shard
		if !tx.IsDestinationTo(shard) && !p.includeCrossShardStartedTransactions {
			p.logIfVerbose(fmt.Sprintf("\t| Transaction with hash %s is not on the destination shard, skipping...", tx.hash))

			continue
		}

		// we skip transactions that are cross shard and still pending for smart-contract results
		if cst := p.internalState.FindCrossShardTransactionByHash(tx.hash); cst != nil {
			p.logIfVerbose(fmt.Sprintf("\t| Transaction with hash %s is still awaiting cross shard SCRs, skipping...", tx.hash))

			continue
		}

		validTransactions = append(validTransactions, tx)
	}

	if !validTransactions.IsEmpty() || p.notifyEmptyBlocks {
		p.logIfVerbose(fmt.Sprintf("\t| Sending %d valid transaction(s) to event consumer...\n", len(validTransactions)))

		p.onTransactionsReceivedFunc(shard, nonce, validTransactions, blockHash)
	}

	return nil
}

func (p *Processor) finalizedCrossShardScrTransactions(shard Shard, transactions Transactions) []*Transaction {
	finalizedTransactions := make(Transactions, 0)

	/*
		Step 1: we add (incrementing counter) pending transactions in the dictionary from current shard to another one
	*/
	for _, tx := range transactions {
		if tx.IsPendingAndOutgoingFromShard(shard) {
			crossShardTransaction := p.internalState.FindCrossShardTransactionByHash(tx.originalTransactionHash)
			if crossShardTransaction == nil {
				originalTx := transactions.FindByHash(tx.originalTransactionHash)
				if originalTx == nil {
					p.logIfVerbose(fmt.Sprintf("\t| Could not identify transaction with hash %s in transaction list\n", tx.originalTransactionHash))

					continue
				}

				p.logIfVerbose(fmt.Sprintf("\t| Creating dictionary for original tx hash %s\n", tx.originalTransactionHash))

				crossShardTransaction = NewCrossShardTransaction(originalTx)
				p.internalState.SetCrossShardTransactionByHash(originalTx.hash, crossShardTransaction)
			}

			if tx.DataEquals("@6f6b") {
				p.logIfVerbose(fmt.Sprintf("\t| Not incrementing counter for cross-shard SCR, original tx hash %s, tx hash %s since the data is @ok (@6f6b)\n", tx.originalTransactionHash, tx.hash))

				continue
			}

			crossShardTransaction.IncrementCounter()

			p.logIfVerbose(fmt.Sprintf("\t| Detected new cross-shard SCR for original tx hash %s, tx hash %s, counter = %d\n", tx.originalTransactionHash, tx.hash, crossShardTransaction.counter))

			p.internalState.crossShardDictionary.Set(tx.originalTransactionHash, crossShardTransaction)
		}
	}

	/*
		Step 2: We remove (decrementing counter) pending transactions in the dictionary from another shard to current shard
	*/
	for _, tx := range transactions {
		if tx.IsPendingAndIncomingToShard(shard) {
			cst := p.internalState.FindCrossShardTransactionByHash(tx.originalTransactionHash)
			if cst == nil {
				p.logIfVerbose(fmt.Sprintf("\t| No counter available for cross-shard SCR, original tx hash %s, tx hash %s", tx.originalTransactionHash, tx.hash))

				continue
			}

			// Ignore the transaction if data is '@ok'
			if tx.DataEquals("@6f6b") {
				p.logIfVerbose(fmt.Sprintf("\t| Not decrementing counter for cross-shard SCR, original tx hash %s, tx hash %s since the data is @ok (@6f6b)\n", tx.originalTransactionHash, tx.hash))

				continue
			}

			cst.DecrementCounter()

			p.logIfVerbose(fmt.Sprintf("\t  Finalized cross-shard SCR for original tx hash %s, tx hash %s, counter = %d\n", tx.originalTransactionHash, tx.hash, cst.counter))

			p.internalState.crossShardDictionary.Set(tx.originalTransactionHash, cst)
		}
	}

	/*
		Step 3. If the counter reaches zero, we remove the value from the cross shard dictionary
	*/
	for hash, crossShardTransaction := range p.internalState.crossShardDictionary {
		if crossShardTransaction.CounterIsZero() {
			p.logIfVerbose(fmt.Sprintf("\t| Completed cross-shard transaction for original tx hash %s", hash))

			// We only add the transaction to the finalized transactions if it isn't already in the list of completed transactions
			tx := transactions.FindByHash(hash)
			if tx != nil {
				finalizedTransactions = append(finalizedTransactions, tx)
			}

			p.internalState.DeleteCrossShardTransaction(hash)
		}
	}

	return finalizedTransactions
}

func (p *Processor) end() {
	p.stateStorage.PersistLastState(p.shards, p.internalState)

	if p.progressBar != nil {
		_ = p.progressBar.Finish()
	}
}

func (p *Processor) currentNonceIsReset(lastProcessedNonce, lastNonceToProcess Nonce) bool {
	return lastProcessedNonce > lastNonceToProcess+networkResetNonceThreshold
}

func (p *Processor) initProgressBar(max int) error {
	if p.displayProgressBar {
		p.progressBar = progressbar.NewOptions64(
			int64(max),
			progressbar.OptionSetDescription("Processing Blocks"),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSetWidth(10),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
			progressbar.OptionSetItsString("Blocks"),
			progressbar.OptionOnCompletion(func() {
				fmt.Fprint(os.Stderr, "\n\n")
			}),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionFullWidth(),
		)

		return p.progressBar.RenderBlank()
	}

	return nil
}

func (p *Processor) incrementProgressBar() {
	if p.progressBar != nil {
		p.progressBar.Add(1)
	}
}
