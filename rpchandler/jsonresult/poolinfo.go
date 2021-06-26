package jsonresult

type PoolInfo struct {
	Info map[int][]BlockInfo `json:"Info"`
}

type SyncInfo struct {
	IsSync      bool
	LastInsert  string
	BlockHeight uint64
	BlockTime   string
	BlockHash   string
}

type SyncStats struct {
	Beacon SyncInfo
	Shard  map[int]*SyncInfo
}

type BlockInfo struct {
	Height  uint64 `json:"BlockHeight"`
	Hash    string `json:"BlockHash"`
	PreHash string `json:"PreHash"`
}
