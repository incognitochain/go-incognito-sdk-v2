package jsonresult

import "github.com/incognitochain/go-incognito-sdk-v2/common"

// ShardBestState describes the best state result of a shard.
type ShardBestState struct {
	BestBlockHash          common.Hash       `json:"BestBlockHash"` // hash of block.
	BestBeaconHash         common.Hash       `json:"BestBeaconHash"`
	BeaconHeight           uint64            `json:"BeaconHeight"`
	ShardID                byte              `json:"ShardID"`
	Epoch                  uint64            `json:"Epoch"`
	ShardHeight            uint64            `json:"ShardHeight"`
	MaxShardCommitteeSize  int               `json:"MaxShardCommitteeSize"`
	MinShardCommitteeSize  int               `json:"MinShardCommitteeSize"`
	ShardProposerIdx       int               `json:"ShardProposerIdx"`
	ShardCommittee         []string          `json:"ShardCommittee"`
	ShardPendingValidator  []string          `json:"ShardPendingValidator"`
	BestCrossShard         map[byte]uint64   `json:"BestCrossShard"` // Best cross shard block by heigh
	StakingTx              map[string]string `json:"StakingTx"`
	NumTxns                uint64            `json:"NumTxns"`                // The number of txns in the block.
	TotalTxns              uint64            `json:"TotalTxns"`              // The total number of txns in the chain.
	TotalTxnsExcludeSalary uint64            `json:"TotalTxnsExcludeSalary"` // for testing and benchmark
	ActiveShards           int               `json:"ActiveShards"`
	MetricBlockHeight      uint64            `json:"MetricBlockHeight"`
	CommitteeFromBlock     common.Hash       `json:"CommitteeFromBlock"`
	CommitteeEngineVersion uint              `json:"CommitteeEngineVersion"`
}
