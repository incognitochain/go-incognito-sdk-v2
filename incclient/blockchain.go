package incclient

import (
	"encoding/json"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
)

// GetActiveShard returns the number of active shards on the Incognito network.
func (client *IncClient) GetActiveShard() (int, error) {
	responseInBytes, err := client.rpcServer.GetActiveShards()
	if err != nil {
		return 0, err
	}

	var activeShards int
	err = rpchandler.ParseResponse(responseInBytes, &activeShards)
	if err != nil {
		return 0, err
	}

	return activeShards, err
}

// GetBestBlock returns the best blocks of the beacon chain and each shard.
func (client *IncClient) GetBestBlock() (map[int]uint64, error) {
	responseInBytes, err := client.rpcServer.GetBestBlock()
	if err != nil {
		return nil, err
	}

	var bestBlocksResult jsonresult.BestBlockResult
	err = rpchandler.ParseResponse(responseInBytes, &bestBlocksResult)
	if err != nil {
		return nil, err
	}

	res := make(map[int]uint64)

	for key, value := range bestBlocksResult.BestBlocks {
		res[key] = value.Height
	}

	return res, nil
}

// GetListToken returns all tokens currently on the Incognito network.
func (client *IncClient) GetListToken() (map[string]CustomToken, error) {
	responseInBytes, err := client.rpcServer.ListPrivacyCustomTokenByRPC()
	if err != nil {
		return nil, err
	}
	var res rpc.ListCustomToken
	err = json.Unmarshal(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	tokenCount := 0
	listTokens := make(map[string]CustomToken)
	for _, token := range res.Result.ListCustomToken {
		tmp := CustomToken{
			tokenID:   token.ID,
			tokenName: token.Name,
			amount:    uint64(token.Amount),
		}
		if len(tmp.tokenName) == 0 {
			tmp.tokenName = fmt.Sprintf("%d", tokenCount)
		}

		listTokens[token.ID] = tmp
		tokenCount++
	}

	return listTokens, nil
}

// GetListTokenIDs returns all token IDs currently on the Incognito network.
func (client *IncClient) GetListTokenIDs() ([]string, error) {
	responseInBytes, err := client.rpcServer.ListPrivacyCustomTokenIDsByRPC()
	if err != nil {
		return nil, err
	}
	var res []string
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetRawMemPool returns a list of transaction hashes currently in the pool.
func (client *IncClient) GetRawMemPool() ([]string, error) {
	responseInBytes, err := client.rpcServer.GetRawMemPool()
	if err != nil {
		return nil, err
	}

	var txHashes map[string][]string
	err = rpchandler.ParseResponse(responseInBytes, &txHashes)
	if err != nil {
		return nil, err
	}

	txList, ok := txHashes["TxHashes"]
	if !ok {
		return nil, fmt.Errorf("TxHashes not found in %v", txHashes)
	}

	return txList, nil

}

// GetCommitteeStateByShard retrieves the committee state of the shardID at the provided root hash, usually used for debugging purposes.
func (client *IncClient) GetCommitteeStateByShard(shardID int, shardRootHash string) (*jsonresult.ShardCommitteeState, error) {
	responseInBytes, err := client.rpcServer.GetCommitteeStateByShardID(shardID, shardRootHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.ShardCommitteeState
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, err
}

// GetShardBestState returns the latest state of a shard chain.
func (client *IncClient) GetShardBestState(shardID int) (*jsonresult.ShardBestState, error) {
	if shardID < 0 || shardID >= common.MaxShardNumber {
		return nil, fmt.Errorf("shardID out of range")
	}

	responseInBytes, err := client.rpcServer.GetShardBestState(byte(shardID))
	if err != nil {
		return nil, err
	}

	var res jsonresult.ShardBestState
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetBeaconBestState returns the latest state of the beacon chain.
func (client *IncClient) GetBeaconBestState(shardID int) (*jsonresult.BeaconBestState, error) {
	if shardID < 0 || shardID >= common.MaxShardNumber {
		return nil, fmt.Errorf("shardID out of range")
	}

	responseInBytes, err := client.rpcServer.GetBeaconBestState()
	if err != nil {
		return nil, err
	}

	var res jsonresult.BeaconBestState
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetBeaconBestState returns the latest state of the beacon chain.
func (client *IncClient) GetBeaconStaker(beaconHeight uint64, beaconStakerPublicKey string) (interface{}, error) {
	return client.rpcServer.GetBeaconStaker(beaconHeight, beaconStakerPublicKey)
}

func (client *IncClient) GetShardStaker(beaconHeight uint64, shardStakerPublicKey string) (interface{}, error) {
	return client.rpcServer.GetShardStaker(beaconHeight, shardStakerPublicKey)
}

func (client *IncClient) GetBeaconCommitteeState(beaconHeight uint64) (interface{}, error) {
	return client.rpcServer.GetBeaconCommitteeState(beaconHeight)
}
