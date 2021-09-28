package jsonresult

// BestBlockResult represents the best block detail of each shard chain and the beacon chain.
type BestBlockResult struct {
	BestBlocks map[int]BestBlockItem `json:"BestBlocks"`
}

// BestBlockItem describes the information of a best block.
type BestBlockItem struct {
	Height              uint64 `json:"Height"`
	Hash                string `json:"Hash"`
	TotalTxs            uint64 `json:"TotalTxs"`
	BlockProducer       string `json:"BlockProducer"`
	ValidationData      string `json:"ValidationData"`
	Epoch               uint64 `json:"Epoch"`
	Time                int64  `json:"Time"`
	RemainingBlockEpoch uint64 `json:"RemainingBlockEpoch"`
	EpochBlock          uint64 `json:"EpochBlock"`
}

// BlockHashResult represents the best block hash of each shard chain and the beacon chain.
type BlockHashResult struct {
	BestBlockHashes map[int]string `json:"BestBlockHashes"`
}
