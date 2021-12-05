package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/thefabric-io/elrond-transaction-processor/processor"
)

func NewProcessorStateStorage(uri string) processor.StateStorage {
	opt, _ := redis.ParseURL(uri)

	s := ProcessorStateStore{
		redis.NewClient(opt)}

	return &s
}

type ProcessorStateStore struct {
	client *redis.Client
}

func (p *ProcessorStateStore) FetchLastState(shards []processor.Shard) (*processor.State, error) {
	var lpn = processor.NonceByShard{}

	for _, shard := range shards {
		cmd := p.client.Get(context.Background(), fmt.Sprintf("%d", shard))
		cmd.Val()
		nonce, err := strconv.Atoi(cmd.Val())
		if err != nil {
			return nil, err
		}

		lpn[shard] = processor.Nonce(nonce)
	}

	csDictionary := processor.NewCrossShardDictionary()

	state := processor.NewState(csDictionary, lpn, nil)
	log.Printf("fetched last processed nonces: %v", state.LastProcessedNonces())

	return state, nil
}

func (p *ProcessorStateStore) PersistLastState(shards processor.Shards, state *processor.State) error {
	for _, shard := range shards {
		nonce, _ := state.LastProcessedNonceInShard(shard)
		if err := p.client.Set(context.Background(), fmt.Sprintf("%d", shard), int(nonce), 0).Err(); err != nil {
			return err
		}

	}
	log.Printf("persisted last processed nonces: %v", state.LastProcessedNonces())

	return nil
}
