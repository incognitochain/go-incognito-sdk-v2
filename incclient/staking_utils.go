package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
)

// GetRewardAmount returns the current reward for a base58-encoded payment address.
// The returned results is a mapping from a tokenID to the corresponding reward amount.
//
// An example result (marshalled):
// 	{
//		"0000000000000000000000000000000000000000000000000000000000000004": 9918708109355,
//		"002ffd86f6b6d0342ebb641e7b89748ba44075db1765173b7d4e77289fbf28fd": 48000,
//		"0fa3e49c7d01a3df067c55293705844ae7d41befd3dfc2f231ab763e9c7daa04": 5,
//		"42f4bee6e1c14f94697fb35b0b0bd7e08da1b3ab8a0311563a6793175e31e93b": 70714,
//		"4946b16a08a9d4afbdf416edf52ef15073db0fc4a63e78eb9de80f94f6c0852a": 2,
//		"5c562893dc38c3c2899143ec32cf67051912fc5b6cf8a8c8c7f8d7397fa64418": 29032,
//		"880ea0787f6c1555e59e3958a595086b7802fc7a38276bcd80d4525606557fbc": 4,
//		"8ba3466c61cbcdd895be8ccbdcc74e7f56a764d6bf390a9abdc8bfe1322e67d6": 4,
//		"961179e5a1c6b354e3544cb7e3c74d1cd1625e59d1138fdafca7b0f9c0c9eaad": 54545,
//		"96d4ee94024abb55c0f000978f73dee078682f94b93a9fa67afdcdc11b79e4ef": 189000,
//		"a37469618aa6e768e6d511db6414fcfe8668b914651976b9509a01ce9e855f58": 94500,
//		"f6c3b18679aff8d307b08d4724697bb8dca123a536b863cbe55dc59c110f5c10": 27272,
//		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854": 9
//	}
func (client *IncClient) GetRewardAmount(paymentAddress string) (map[common.Hash]uint64, error) {
	responseInBytes, err := client.rpcServer.GetRewardAmount(paymentAddress)
	if err != nil {
		return nil, err
	}

	var res map[common.Hash]uint64
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	// the full-node return "PRV":amount so we have to parse it into "0000000000000000000000000000000000000000000000000000000000000004":amount
	tmpPRVHash, err := common.Hash{}.NewHashFromStr("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		return nil, err
	}
	res[common.PRVCoinID] = res[*tmpPRVHash]
	delete(res, *tmpPRVHash)

	return res, err
}

// ListReward returns the staking rewards on the blockchain.
// The returned results is a mapping from a public key to a tokenID-reward mapping.
//
// An example result (marshalled):
//	"117XEowF5Y4eYs6mTQPkZxaK3H9GSwGHdcuv6emK7NVV698XSq": {
//		"0000000000000000000000000000000000000000000000000000000000000004": 9792062854956,
//		"002ffd86f6b6d0342ebb641e7b89748ba44075db1765173b7d4e77289fbf28fd": 48000,
//		"15aeb4c4ea24a50695a0cb425b711a49cc9b7ab2e56af8184a570db8c3e34ff8": 122916,
//		"2d04e28959cf3767734d9a7adbe639f8818d32c4531e467108c07e2254a6e4eb": 242314,
//		"42f4bee6e1c14f94697fb35b0b0bd7e08da1b3ab8a0311563a6793175e31e93b": 24827,
//		"4946b16a08a9d4afbdf416edf52ef15073db0fc4a63e78eb9de80f94f6c0852a": 3,
//		"880ea0787f6c1555e59e3958a595086b7802fc7a38276bcd80d4525606557fbc": 0,
//		"961179e5a1c6b354e3544cb7e3c74d1cd1625e59d1138fdafca7b0f9c0c9eaad": 45000,
//		"9abfda385c6700656778da12c21b36698bdf9fff250d94314d53c6069c5c45a4": 34758,
//		"9fca0a0947f4393994145ef50eecd2da2aa15da2483b310c2c0650301c59b17d": 0,
//		"f6c3b18679aff8d307b08d4724697bb8dca123a536b863cbe55dc59c110f5c10": 88306,
//		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854": 0
//	}
func (client *IncClient) ListReward() (map[string]map[common.Hash]uint64, error) {
	responseInBytes, err := client.rpcServer.GetListRewardAmount()
	if err != nil {
		return nil, err
	}

	var res map[string]map[common.Hash]uint64
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, err
}

// GetMiningInfo returns the mining information of a node.
//
// Create an IncClient instance pointing to your node and call this function to gather the node's mining information.
func (client *IncClient) GetMiningInfo() (*jsonresult.MiningInfoResult, error) {
	responseInBytes, err := client.rpcServer.GetMiningInfo()
	if err != nil {
		return nil, err
	}

	var res jsonresult.MiningInfoResult
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetSyncStats returns the statistics of data-synchronizing status.
//
// Create an IncClient instance pointing to your node and call this function to get the statistics.
func (client *IncClient) GetSyncStats() (*jsonresult.SyncStats, error) {
	responseInBytes, err := client.rpcServer.GetSyncStats()
	if err != nil {
		return nil, err
	}

	var res jsonresult.SyncStats
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
