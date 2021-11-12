package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"math/big"
	"testing"
)

func calculatePoolAmount(pool *jsonresult.PoolInfo, totalShare uint64, shareAmount uint64) (uint64, uint64) {
	shareBig := new(big.Int).SetUint64(shareAmount)
	totalShareBig := new(big.Int).SetUint64(totalShare)

	value1 := new(big.Int).SetUint64(pool.Token1PoolValue)
	value1 = value1.Mul(value1, shareBig)
	value1 = value1.Div(value1, totalShareBig)

	value2 := new(big.Int).SetUint64(pool.Token2PoolValue)
	value2 = value2.Mul(value2, shareBig)
	value2 = value2.Div(value2, totalShareBig)

	return value1.Uint64(), value2.Uint64()
}

func TestIncClient_CreateAndSendPdexv3WithdrawStakeRewardTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	stakingPoolIDStr := common.PRVIDStr
	nftIDStr := "eb1ec0987a37829831c8d947ef2c48f8ab6ada4b02d99e82039ca5977570bd0c"

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	txHash, err := ic.CreateAndSendPdexv3WithdrawStakeRewardTransaction(privateKey, stakingPoolIDStr, nftIDStr)
	if err != nil {
		panic(err)
	}
	Logger.Printf("TxHash: %v\n", txHash)
}
