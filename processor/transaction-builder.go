package processor

func NewTransactionBuilder() *TransactionBuilder {
	return &TransactionBuilder{shardTransaction: &Transaction{}}
}

type TransactionBuilder struct {
	shardTransaction *Transaction
}

func (b *TransactionBuilder) NewTransaction() *TransactionBuilder {
	b.shardTransaction = &Transaction{}

	return b
}

func (b *TransactionBuilder) Value(v string) *TransactionBuilder {
	b.shardTransaction.value = v

	return b
}

func (b *TransactionBuilder) Data(d string) *TransactionBuilder {
	b.shardTransaction.data = d

	return b
}

func (b *TransactionBuilder) Hash(h string) *TransactionBuilder {
	b.shardTransaction.hash = h

	return b
}

func (b *TransactionBuilder) Sender(s string) *TransactionBuilder {
	b.shardTransaction.sender = s

	return b
}

func (b *TransactionBuilder) Receiver(r string) *TransactionBuilder {
	b.shardTransaction.receiver = r

	return b
}

func (b *TransactionBuilder) Status(s string) *TransactionBuilder {
	b.shardTransaction.status = s

	return b
}

func (b *TransactionBuilder) SourceShard(s Shard) *TransactionBuilder {
	b.shardTransaction.sourceShard = s

	return b
}

func (b *TransactionBuilder) DestinationShard(s Shard) *TransactionBuilder {
	b.shardTransaction.destinationShard = s

	return b
}

func (b *TransactionBuilder) Nonce(n Nonce) *TransactionBuilder {
	b.shardTransaction.nonce = n

	return b
}

func (b *TransactionBuilder) PreviousTransactionHash(h string) *TransactionBuilder {
	b.shardTransaction.previousTransactionHash = h
	return b
}

func (b *TransactionBuilder) OriginalTransactionHash(h string) *TransactionBuilder {
	b.shardTransaction.originalTransactionHash = h
	return b
}

func (b *TransactionBuilder) GasPrice(g int) *TransactionBuilder {
	b.shardTransaction.gasPrice = g
	return b
}

func (b *TransactionBuilder) GasLimit(g int) *TransactionBuilder {
	b.shardTransaction.gasLimit = g
	return b
}

func (b *TransactionBuilder) Build() *Transaction {
	return b.shardTransaction
}
