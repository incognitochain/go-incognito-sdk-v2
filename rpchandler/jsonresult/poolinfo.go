package jsonresult

// NodePoolInfo describes the information of a node.
type NodePoolInfo struct {
	Info map[int][]BlockInfo `json:"Info"`
}

// SyncInfo describes the sync information of a node.
type SyncInfo struct {
	IsSync      bool
	LastInsert  string
	BlockHeight uint64
	BlockTime   string
	BlockHash   string
}

// SyncStats describes the sync statistics of a node.
type SyncStats struct {
	Beacon SyncInfo
	Shard  map[int]*SyncInfo
}

// BlockInfo consists of simplified information of a block.
type BlockInfo struct {
	Height  uint64 `json:"BlockHeight"`
	Hash    string `json:"BlockHash"`
	PreHash string `json:"PreHash"`
}
