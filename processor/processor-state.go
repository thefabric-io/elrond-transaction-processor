package processor

import (
	"log"
	"time"
)

func NewState(dictionary CrossShardDictionary, fromNonces, toNonces NonceByShard) *State {
	s := State{
		crossShardDictionary:        dictionary,
		lastProcessedNoncesInternal: fromNonces,
		toNonces:                    toNonces,
	}

	return &s
}

type State struct {
	crossShardDictionary        CrossShardDictionary
	lastProcessedNoncesInternal NonceByShard
	toNonces                    NonceByShard
}

func (s *State) LastProcessedNonces() NonceByShard {
	return s.lastProcessedNoncesInternal
}

func (s *State) LastProcessedNonceInShard(shard Shard) (Nonce, bool) {
	n, f := s.lastProcessedNoncesInternal[shard]

	return n, f
}

func (s *State) LastNonceToProcessInShard(shard Shard) (Nonce, bool) {
	n, f := s.toNonces[shard]

	return n, f
}

func (s *State) AddBufferToLastProcessNonces(buffer int) {
	for shard, nonce := range s.lastProcessedNoncesInternal {
		s.lastProcessedNoncesInternal[shard] = nonce.Subtract(Nonce(buffer + 1))
	}
}

func (s *State) FindCrossShardTransactionByHash(h string) *CrossShardTransaction {
	return s.crossShardDictionary.FindTransaction(h)
}

func (s *State) SetCrossShardTransactionByHash(h string, t *CrossShardTransaction) {
	s.crossShardDictionary.Set(h, t)
}

func (s *State) DeleteCrossShardTransaction(h string) {
	s.crossShardDictionary.Delete(h)
}

func (s *State) PruneCrossShardDictionary() {
	for h, crossShardTransaction := range s.crossShardDictionary {
		elapsed := time.Now().Sub(crossShardTransaction.created)
		if elapsed > 600*time.Second {
			log.Printf("pruning transaction with hash %s since its elapsed time is %.2f seconds", h, elapsed.Seconds())
			s.crossShardDictionary.Delete(h)
		}
	}
}

func (s *State) NumberOfRemainingNonces() int {
	var fromNonce, toNonce Nonce
	for _, nonce := range s.lastProcessedNoncesInternal {
		fromNonce += nonce
	}

	for _, nonce := range s.toNonces {
		toNonce += nonce
	}

	return int(toNonce.Subtract(fromNonce))
}
