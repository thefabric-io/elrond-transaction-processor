package processor

import (
	"encoding/base64"
	"log"
)

type Transaction struct {
	value                   string
	data                    string
	hash                    string
	sender                  string
	receiver                string
	status                  string
	sourceShard             Shard
	destinationShard        Shard
	nonce                   Nonce
	previousTransactionHash string
	originalTransactionHash string
	gasPrice                int
	gasLimit                int
}

func (t *Transaction) Sender() string {
	return t.sender
}

func (t *Transaction) HasOriginalTransactionHash() bool {
	return len(t.originalTransactionHash) != 0
}

func (t *Transaction) IsFromShard(s Shard) bool {
	return t.sourceShard == s
}

func (t *Transaction) IsDestinationTo(s Shard) bool {
	return t.destinationShard == s
}

func (t *Transaction) IsPendingAndOutgoingFromShard(s Shard) bool {
	return t.HasOriginalTransactionHash() && t.IsFromShard(s) && !t.IsDestinationTo(s)
}

func (t *Transaction) IsPendingAndIncomingToShard(s Shard) bool {
	return t.HasOriginalTransactionHash() && !t.IsFromShard(s) && t.IsDestinationTo(s)
}

func (t *Transaction) B64DataDecoded() (string, error) {
	if len(t.data) != 0 {
		data, err := base64.StdEncoding.DecodeString(t.data)
		if err != nil {
			log.Printf("could not decode data from transaction %s in transaction list\n", t.hash)

			return "", err
		}

		return string(data), nil
	}

	return "", nil
}

func (t *Transaction) DataEquals(d string) bool {
	data, err := t.B64DataDecoded()
	if err != nil {
		return false
	}

	return data == d
}

type Transactions []*Transaction

func (tt Transactions) FindByHash(h string) *Transaction {
	for i := range tt {
		if tt[i].hash == h {
			return tt[i]
		}
	}

	return nil
}

func (tt Transactions) IsEmpty() bool {
	return len(tt) == 0
}
