package jsonresult

type MiningInfoResult struct {
	ShardHeight         uint64 `json:"ShardHeight"`
	BeaconHeight        uint64 `json:"BeaconHeight"`
	CurrentShardBlockTx int    `json:"CurrentShardBlockTx"`
	PoolSize            int    `json:"PoolSize"`
	Chain               string `json:"Chain"`
	ShardID             int    `json:"ShardID"`
	Layer               string `json:"Layer"`
	Role                string `json:"Role"`
	MiningPublickey     string `json:"MiningPublickey"`
	IsEnableMining      bool   `json:"IsEnableMining"`
}
