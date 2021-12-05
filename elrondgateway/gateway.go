package elrondgateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/thefabric-io/elrond-transaction-processor/processor"
)

const (
	MainNetGatewayURL string = "https://gateway.elrond.com"
	TestNetGatewayURL string = "https://testnet-gateway.elrond.com"
	CodeSuccessful    string = "successful"
)

func NewClient(url string) *Client {
	return &Client{url: url}
}

type Client struct {
	url string
}

func (e *Client) GetShards() ([]processor.Shard, error) {
	b, err := e.get("network/config")
	if err != nil {
		return nil, err
	}

	response := GetShardsResponse{}
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	if response.Code != CodeSuccessful {
		return nil, errors.New(fmt.Sprintf("%s: %s", response.Code, response.Error))
	}

	shardCount := response.Data.Config.ErdNumShardsWithoutMeta

	result := make([]processor.Shard, shardCount)
	for i := 0; i < shardCount; i++ {
		result[i] = processor.Shard(i)
	}

	result = append(result, processor.ShardMetachain)

	return result, nil
}

func (e *Client) GetCurrentNoncesForShards(shards []processor.Shard) (processor.NonceByShard, error) {
	var err error

	result := make(processor.NonceByShard, len(shards))
	for _, shard := range shards {
		result[shard], err = e.GetCurrentNonceForShard(shard)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (e *Client) GetCurrentNonceForShard(shard processor.Shard) (processor.Nonce, error) {
	b, err := e.get(fmt.Sprintf("network/status/%d", shard))
	if err != nil {
		return 0, err
	}

	response := GetCurrentNonceForShardResponse{}
	if err := json.Unmarshal(b, &response); err != nil {
		return 0, err
	}

	if response.Code != CodeSuccessful {
		return 0, errors.New(fmt.Sprintf("%s: %s", response.Code, response.Error))
	}

	return processor.Nonce(response.Data.Status.ErdNonce), nil
}

func (e *Client) GetShardTransactions(shard processor.Shard, nonce processor.Nonce) (string, []*processor.Transaction, error) {
	b, err := e.get(fmt.Sprintf("block/%d/by-nonce/%d?withTxs=true", shard, nonce))
	if err != nil {
		return "", nil, err
	}

	response := GetShardTransactionsResponse{}
	if err := json.Unmarshal(b, &response); err != nil {
		return "", nil, err
	}

	if response.Code != CodeSuccessful {
		return "", nil, errors.New(fmt.Sprintf("%s: %s", response.Code, response.Error))
	}

	if len(response.Data.Block.Hash) == 0 {
		return "", nil, errors.New(fmt.Sprintf("Block for shard %d and nonce %d is undefined or block not available\n", shard, nonce))
	}

	if len(response.Data.Block.MiniBlocks) == 0 {
		return response.Data.Block.Hash, []*processor.Transaction{}, nil
	}

	results := make([]*processor.Transaction, 0)

	txB := processor.NewTransactionBuilder()

	for _, mb := range response.Data.Block.MiniBlocks {
		for _, mbTx := range mb.Transactions {
			tx := txB.NewTransaction().
				Value(mbTx.Value).
				Data(mbTx.Data).
				Hash(mbTx.Hash).
				Sender(mbTx.Sender).
				Receiver(mbTx.Receiver).
				Status(mbTx.Status).
				SourceShard(processor.Shard(mbTx.SourceShard)).
				DestinationShard(processor.Shard(mbTx.DestinationShard)).
				Nonce(processor.Nonce(mbTx.Nonce)).
				PreviousTransactionHash(mbTx.PreviousTransactionHash).
				OriginalTransactionHash(mbTx.OriginalTransactionHash).
				GasPrice(mbTx.GasPrice).
				GasLimit(mbTx.GasLimit).
				Build()
			results = append(results, tx)
		}
	}

	return response.Data.Block.Hash, results, nil

}

func (e *Client) get(path string) ([]byte, error) {
	if e.url == "" {
		e.url = MainNetGatewayURL
	}

	fullUrl := fmt.Sprintf("%s/%s", e.url, path)
	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
