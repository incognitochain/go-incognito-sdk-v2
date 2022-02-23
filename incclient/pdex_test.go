package incclient

import (
	"math/big"
	"testing"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataPdexv3 "github.com/incognitochain/go-incognito-sdk-v2/metadata/pdexv3"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
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

func TestIncClient_CreateAndPdexv3ModifyParamsTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	newParams := metadataPdexv3.Pdexv3Params{
		DefaultFeeRateBPS:               30,
		FeeRateBPS:                      map[string]uint{},
		PRVDiscountPercent:              25,
		TradingProtocolFeePercent:       0,
		TradingStakingPoolRewardPercent: 10,
		PDEXRewardPoolPairsShare: map[string]uint{
			"0000000000000000000000000000000000000000000000000000000000000b7c-00000000000000000000000000000000000000000000000000000000000115d7-a1747b936740aac9a34d9841bd57f6b969ef2b75402bab3907da070ed7f6d343": 4000,
		},
		StakingPoolsShare:                 map[string]uint{},
		StakingRewardTokens:               []common.Hash{},
		MintNftRequireAmount:              1000000000,
		MaxOrdersPerNft:                   10,
		AutoWithdrawOrderLimitAmount:      100,
		MinPRVReserveTradingRate:          1000000000000,
		DefaultOrderTradingRewardRatioBPS: 2500,
		OrderTradingRewardRatioBPS: map[string]uint{
			"0000000000000000000000000000000000000000000000000000000000000b7c-00000000000000000000000000000000000000000000000000000000000115d7-a1747b936740aac9a34d9841bd57f6b969ef2b75402bab3907da070ed7f6d343": 4000,
			"000000000000000000000000000000000000000000000000000000000000e776-00000000000000000000000000000000000000000000000000000000000115d7-82035dd88dba066741855783d9f37a5fb5cc285f008cac0b8ffbf43338bb687c": 3500,
		},
		OrderLiquidityMiningBPS: map[string]uint{
			"0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115dc-235ea1cc5122c3d35b4c8b466d20c8308f2f3548d234adcf2031836f544e3831": 2000,
			"0000000000000000000000000000000000000000000000000000000000000b7c-00000000000000000000000000000000000000000000000000000000000115d7-a1747b936740aac9a34d9841bd57f6b969ef2b75402bab3907da070ed7f6d343": 2000,
			"00000000000000000000000000000000000000000000000000000000000115d7-00000000000000000000000000000000000000000000000000000000000115dc-6398ddb36dec2c75f8df62e9ea4d5e11aad274c61e6dab1371d5975d915f2aeb": 2000,
		},
		DAOContributingPercent:    50,
		MiningRewardPendingBlocks: 50,
	}
	txHash, err := ic.CreateAndSendPdexv3ModifyParamsTransaction(privateKey, newParams)
	if err != nil {
		panic(err)
	}
	Logger.Printf("TxHash: %v\n", txHash)
}
