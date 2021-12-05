package processor

type Option func(*Processor)

type Options struct{}

func (oo *Options) DataSource(s DataSource) Option {
	return func(p *Processor) {
		p.dataSource = s
	}
}
func (oo *Options) StateStorage(s StateStorage) Option {
	return func(p *Processor) {
		p.stateStorage = s
	}
}

func (oo *Options) WaitForFinalizedCrossShardSmartContractResults(b bool) Option {
	return func(p *Processor) {
		p.waitForFinalizedCrossShardSmartContractResults = b
	}
}

func (oo *Options) NotifyEmptyBlocks(b bool) Option {
	return func(p *Processor) {
		p.notifyEmptyBlocks = b
	}
}

func (oo *Options) IncludeCrossShardStartedTransactions(b bool) Option {
	return func(p *Processor) {
		p.includeCrossShardStartedTransactions = b
	}
}

func (oo *Options) OnTransactionsReceived(f OnTransactionReceivedFunc) Option {
	return func(p *Processor) {
		p.onTransactionsReceivedFunc = f
	}
}

func (oo *Options) PastTransactionBufferPerShard(d int) Option {
	return func(p *Processor) {
		p.pastBlocksBuffer = d
	}
}

func (oo *Options) Verbose() Option {
	return func(p *Processor) {
		p.verbose = true
	}
}

func (oo *Options) DisplayProgressBar() Option {
	return func(p *Processor) {
		p.displayProgressBar = true
	}
}
