package processor

type NonceByShard map[Shard]Nonce

func (nn NonceByShard) PutNonce(shard Shard, nonce Nonce) {
	nn[shard] = nonce
}

func (nn NonceByShard) Nonce(shard Shard) (Nonce, bool) {
	nonce, found := nn[shard]

	return nonce, found
}
