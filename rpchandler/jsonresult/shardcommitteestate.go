package jsonresult

// ShardCommitteeState describes a committee state of a shard.
type ShardCommitteeState struct {
	Root       string   `json:"root"`
	ShardID    uint64   `json:"shardID"`
	Committee  []string `json:"committee"`
	Substitute []string `json:"substitute"`
}
