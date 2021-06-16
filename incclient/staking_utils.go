package incclient

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

// GetRewardAmountByPublicKey returns the current reward for a public key.
func (client *IncClient) GetRewardAmountByPublicKey(publicKey string) (uint64, error) {
	responseInBytes, err := client.rpcServer.GetRewardAmountByPublicKey(publicKey)
	if err != nil {
		return 0, err
	}

	response, err := rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return 0, err
	}

	var res uint64
	err = json.Unmarshal(response.Result, &res)

	return res, err
}

// ListReward returns the staking rewards on the blockchain.
func (client *IncClient) ListReward() (map[string]map[common.Hash]uint64, error) {
	responseInBytes, err := client.rpcServer.GetListRewardAmount()
	if err != nil {
		return nil, err
	}

	response, err := rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return nil, err
	}

	var res map[string]map[common.Hash]uint64
	err = json.Unmarshal(response.Result, &res)
	return res, err
}
