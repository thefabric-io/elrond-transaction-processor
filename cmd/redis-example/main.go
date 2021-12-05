package main

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thefabric-io/elrond-transaction-processor/elrondgateway"
	"github.com/thefabric-io/elrond-transaction-processor/processor"
)

func onTransactionReceivedFunc(shard processor.Shard, nonce processor.Nonce, transactions []*processor.Transaction, blockHash string) {
	/*
		Commented due to bad readability of progress bar when enabled
		fmt.Printf("%d transaction(s) received from %s with nonce %d\n", len(transactions), shard.Name(), nonce)
	*/

	/*
		TODO:
			Filter transactions to be processed (e.g. by contract(s))
			Send shard, nonce and transactions to a message broker (e.g. Kafka)
			Consume message and persist for later querying (e.g. relational/document database, elasticsearch)
	*/
}

func main() {
	_ = godotenv.Load(".env")

	/*
		Make sure to update the state of your actual state storage to avoid reading unnecessary and older blocks (nonces)
		when switching to MainNetGatewayURL
	*/
	elrondGateway := elrondgateway.NewClient(elrondgateway.TestNetGatewayURL)

	/*
		Make sure you set an environment variable to reach your storage server.
	*/
	redisServerURI := os.Getenv("REDIS_SERVER_URI")
	if redisServerURI == "" {
		panic(errors.New("could not fetch env variable REDIS_SERVER_URI"))
	}

	/*
		!Disclaimer: In this example, stateStorage (as Redis server), do not persist the cross shard transactions dictionary.
		Please make sure your processor state repository persists the cross shard transactions dictionary if
		IncludeCrossShardStartedTransactions is true.
	*/
	stateStorage := NewProcessorStateStorage(redisServerURI)

	opts := processor.Options{}
	proc, err := processor.NewProcessor(
		opts.DataSource(elrondGateway),
		opts.StateStorage(stateStorage),
		opts.OnTransactionsReceived(onTransactionReceivedFunc),
		opts.IncludeCrossShardStartedTransactions(true),
		opts.PastTransactionBufferPerShard(2),
		opts.WaitForFinalizedCrossShardSmartContractResults(false),
		opts.Verbose(),
		opts.DisplayProgressBar(),
	)
	if err != nil {
		panic(err)
	}

	if err = proc.Start(); err != nil {
		log.Println(err)
	}
}
