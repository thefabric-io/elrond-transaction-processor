package processor

import "time"

func NewCrossShardTransaction(t *Transaction) *CrossShardTransaction {
	return &CrossShardTransaction{transaction: *t}
}

type CrossShardTransaction struct {
	transaction Transaction
	counter     int
	created     time.Time
}

func (t *CrossShardTransaction) CounterIsZero() bool {
	return t.counter == 0
}

func (t *CrossShardTransaction) IncrementCounter() {
	t.counter++
}

func (t *CrossShardTransaction) DecrementCounter() {
	t.counter--
}
