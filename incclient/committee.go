package incclient

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

type ShardCommitteeState struct {
	Root       string   `json:"root"`
	ShardID    uint64   `json:"shardID"`
	Committee  []string `json:"committee"`
	Substitute []string `json:"substitute"`
}

// GetCommitteeStateByShard retrieves the committee state of the shardID at the provided root hash.
func (client *IncClient) GetCommitteeStateByShard(shardID int, shardRootHash string) (*ShardCommitteeState, error) {
	responseInBytes, err := client.rpcServer.GetCommitteeStateByShardID(shardID, shardRootHash)
	if err != nil {
		return nil, err
	}

	response, err := rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return nil, err
	}

	var res ShardCommitteeState
	err = json.Unmarshal(response.Result, &res)

	return &res, err
}
