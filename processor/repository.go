package processor

type StateStorage interface {
	PersistLastState(shards Shards, state *State) error
	FetchLastState(shards []Shard) (*State, error)
}

type DataSource interface {
	GetShards() ([]Shard, error)
	GetCurrentNonceForShard(shard Shard) (Nonce, error)
	GetCurrentNoncesForShards([]Shard) (NonceByShard, error)
	GetShardTransactions(shard Shard, nonce Nonce) (string, []*Transaction, error)
}
