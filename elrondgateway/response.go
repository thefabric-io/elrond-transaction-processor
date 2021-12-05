package elrondgateway

type GetShardsResponse struct {
	Data struct {
		Config struct {
			ErdChainId                   string `json:"erd_chain_id"`
			ErdDenomination              int    `json:"erd_denomination"`
			ErdGasPerDataByte            int    `json:"erd_gas_per_data_byte"`
			ErdGasPriceModifier          string `json:"erd_gas_price_modifier"`
			ErdMaxGasPerTransaction      int    `json:"erd_max_gas_per_transaction"`
			ErdMetaConsensusGroupSize    int    `json:"erd_meta_consensus_group_size"`
			ErdMinGasLimit               int    `json:"erd_min_gas_limit"`
			ErdMinGasPrice               int    `json:"erd_min_gas_price"`
			ErdMinTransactionVersion     int    `json:"erd_min_transaction_version"`
			ErdNumMetachainNodes         int    `json:"erd_num_metachain_nodes"`
			ErdNumNodesInShard           int    `json:"erd_num_nodes_in_shard"`
			ErdNumShardsWithoutMeta      int    `json:"erd_num_shards_without_meta"`
			ErdRewardsTopUpGradientPoint string `json:"erd_rewards_top_up_gradient_point"`
			ErdRoundDuration             int    `json:"erd_round_duration"`
			ErdRoundsPerEpoch            int    `json:"erd_rounds_per_epoch"`
			ErdShardConsensusGroupSize   int    `json:"erd_shard_consensus_group_size"`
			ErdStartTime                 int    `json:"erd_start_time"`
			ErdTopUpFactor               string `json:"erd_top_up_factor"`
		} `json:"config"`
	} `json:"data"`
	Code  string `json:"code"`
	Error string `json:"error"`
}

type GetCurrentNonceForShardResponse struct {
	Data struct {
		Status struct {
			ErdCurrentRound               int `json:"erd_current_round"`
			ErdEpochNumber                int `json:"erd_epoch_number"`
			ErdHighestFinalNonce          int `json:"erd_highest_final_nonce"`
			ErdNonce                      int `json:"erd_nonce"`
			ErdNonceAtEpochStart          int `json:"erd_nonce_at_epoch_start"`
			ErdNoncesPassedInCurrentEpoch int `json:"erd_nonces_passed_in_current_epoch"`
			ErdRoundAtEpochStart          int `json:"erd_round_at_epoch_start"`
			ErdRoundsPassedInCurrentEpoch int `json:"erd_rounds_passed_in_current_epoch"`
			ErdRoundsPerEpoch             int `json:"erd_rounds_per_epoch"`
		} `json:"status"`
	} `json:"data"`
	Code  string `json:"code"`
	Error string `json:"error"`
}

type GetShardTransactionsResponse struct {
	Data struct {
		Block struct {
			Nonce         int    `json:"nonce"`
			Round         int    `json:"round"`
			Hash          string `json:"hash"`
			PrevBlockHash string `json:"prevBlockHash"`
			Epoch         int    `json:"epoch"`
			Shard         int    `json:"shard"`
			NumTxs        int    `json:"numTxs"`
			MiniBlocks    []struct {
				Hash             string `json:"hash"`
				Type             string `json:"type"`
				SourceShard      int64  `json:"sourceShard"`
				DestinationShard int    `json:"destinationShard"`
				Transactions     []struct {
					Type                    string `json:"type"`
					Hash                    string `json:"hash"`
					Nonce                   int    `json:"nonce"`
					Value                   string `json:"value"`
					Receiver                string `json:"receiver"`
					Sender                  string `json:"sender"`
					GasPrice                int    `json:"gasPrice"`
					GasLimit                int    `json:"gasLimit,omitempty"`
					Data                    string `json:"data"`
					PreviousTransactionHash string `json:"previousTransactionHash,omitempty"`
					OriginalTransactionHash string `json:"originalTransactionHash,omitempty"`
					SourceShard             int    `json:"sourceShard"`
					DestinationShard        int    `json:"destinationShard"`
					MiniblockType           string `json:"miniblockType"`
					MiniblockHash           string `json:"miniblockHash"`
					Status                  string `json:"status"`
					Signature               string `json:"signature,omitempty"`
				} `json:"transactions"`
			} `json:"miniBlocks"`
			Timestamp       int    `json:"timestamp"`
			AccumulatedFees string `json:"accumulatedFees"`
			DeveloperFees   string `json:"developerFees"`
			Status          string `json:"status"`
		} `json:"block"`
	} `json:"data"`
	Code  string `json:"code"`
	Error string `json:"error"`
}
