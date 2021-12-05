package processor

import (
	"log"
	"time"
)

const (
	minPruningElapsedTime = 600 * time.Second
)

func NewCrossShardDictionary() CrossShardDictionary {
	return CrossShardDictionary{}
}

type CrossShardDictionary map[string]*CrossShardTransaction

func (m CrossShardDictionary) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func (m CrossShardDictionary) Set(id string, tx *CrossShardTransaction) {
	m[id] = tx
}

func (m CrossShardDictionary) FindTransaction(h string) *CrossShardTransaction {
	crossShardItem, found := m[h]
	if !found {
		return nil
	}

	return crossShardItem
}

func (m CrossShardDictionary) Delete(h string) {
	delete(m, h)
}

func (m CrossShardDictionary) PruneTransactions() {
	for h, crossShardTransaction := range m {
		elapsed := time.Now().Sub(crossShardTransaction.created)
		if elapsed > minPruningElapsedTime {
			log.Printf("pruning transaction with hash %s since its elapsed time is %.2f seconds", h, elapsed.Seconds())
			delete(m, h)
		}
	}
}
