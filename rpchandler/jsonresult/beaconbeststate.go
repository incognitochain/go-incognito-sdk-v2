package jsonresult

import "github.com/incognitochain/go-incognito-sdk-v2/common"

type MissingSignature struct {
	Total   uint
	Missing uint
}

type Penalty struct {
	MinPercent   uint
	Time         int64
	ForceUnstake bool
}

type GetBeaconBestState struct {
	BestBlockHash                          common.Hash                 `json:"BestBlockHash"`         // The hash of the block.
	PreviousBestBlockHash                  common.Hash                 `json:"PreviousBestBlockHash"` // The hash of the block.
	BestShardHash                          map[byte]common.Hash        `json:"BestShardHash"`
	BestShardHeight                        map[byte]uint64             `json:"BestShardHeight"`
	Epoch                                  uint64                      `json:"Epoch"`
	BeaconHeight                           uint64                      `json:"BeaconHeight"`
	BeaconProposerIndex                    int                         `json:"BeaconProposerIndex"`
	BeaconCommittee                        []string                    `json:"BeaconCommittee"`
	BeaconPendingValidator                 []string                    `json:"BeaconPendingValidator"`
	CandidateShardWaitingForCurrentRandom  []string                    `json:"CandidateShardWaitingForCurrentRandom"` // snapshot shard candidate list, waiting to be shuffled in this current epoch
	CandidateBeaconWaitingForCurrentRandom []string                    `json:"CandidateBeaconWaitingForCurrentRandom"`
	CandidateShardWaitingForNextRandom     []string                    `json:"CandidateShardWaitingForNextRandom"` // shard candidate list, waiting to be shuffled in next epoch
	CandidateBeaconWaitingForNextRandom    []string                    `json:"CandidateBeaconWaitingForNextRandom"`
	RewardReceiver                         map[string]string           `json:"RewardReceiver"`        // key: incognito public key of committee, value: payment address reward receiver
	ShardCommittee                         map[byte][]string           `json:"ShardCommittee"`        // current committee and validator of all shard
	ShardPendingValidator                  map[byte][]string           `json:"ShardPendingValidator"` // pending candidate waiting for swap to get in committee of all shard
	AutoStaking                            map[string]bool             `json:"AutoStaking"`
	StakingTx                              map[string]common.Hash      `json:"StakingTx"`
	CurrentRandomNumber                    int64                       `json:"CurrentRandomNumber"`
	CurrentRandomTimeStamp                 int64                       `json:"CurrentRandomTimeStamp"` // random timestamp for this epoch
	IsGetRandomNumber                      bool                        `json:"IsGetRandomNumber"`
	MaxBeaconCommitteeSize                 int                         `json:"MaxBeaconCommitteeSize"`
	MinBeaconCommitteeSize                 int                         `json:"MinBeaconCommitteeSize"`
	MaxShardCommitteeSize                  int                         `json:"MaxShardCommitteeSize"`
	MinShardCommitteeSize                  int                         `json:"MinShardCommitteeSize"`
	ActiveShards                           int                         `json:"ActiveShards"`
	LastCrossShardState                    map[byte]map[byte]uint64    `json:"LastCrossShardState"`
	ShardHandle                            map[byte]bool               `json:"ShardHandle"` // lock sync.RWMutex
	CommitteeEngineVersion                 uint                        `json:"CommitteeEngineVersion"`
	NumberOfMissingSignature               map[string]MissingSignature `json:"MissingSignature"`        // lock sync.RWMutex
	MissingSignaturePenalty                map[string]Penalty          `json:"MissingSignaturePenalty"` // lock sync.RWMutex
}
