package incclient

import (
	"encoding/json"
	"testing"
	"time"
)

func TestIncClient_GetPoolPairStateByID(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	poolID := "0000000000000000000000000000000000000000000000000000000000000004-0000000000000000000000000000000000000000000000000000000000000006-56e4e9d710a01dfe865e6d5047fabd6bb98b646465863c2726ebc56538983b5d"
	poolState, err := ic.GetPoolPairStateByID(0, poolID)
	if err != nil {
		panic(err)
	}
	jsb, _ := json.MarshalIndent(poolState, "", "\t")
	Logger.Printf("state: %v\n", string(jsb))
}

func TestIncClient_GetPoolShareAmount(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	poolID := "0000000000000000000000000000000000000000000000000000000000000004-0000000000000000000000000000000000000000000000000000000000000006-56e4e9d710a01dfe865e6d5047fabd6bb98b646465863c2726ebc56538983b5d"
	nftID := "d150bd389f7f881a271e1617aba13dbc6c0dde7b8d184f0cbd637e93aa83c69f"
	share, err := ic.GetPoolShareAmount(poolID, nftID)
	if err != nil {
		panic(err)
	}
	Logger.Printf("share: %v\n", share)
}

func TestIncClient_CheckTradeStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "e4c13e368eb4da34ebcd04aaf9da9a401d5f55df752f3d1c650331a19f69a53a"
	status, err := ic.CheckTradeStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckAddLiquidityStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "928ca5fef7f8274f025c5184240f3b5b13f310e2a11dd553b8fa656901e0827f"
	status, err := ic.CheckDEXLiquidityWithdrawalStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckLiquidityWithdrawalStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "13fd37a90bb838d3402c37fc4b11c3715ef847bfcc397d4fff3a04b351e12388"
	status, err := ic.CheckDEXLiquidityWithdrawalStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckOrderAddedStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "91c1571bd0debf386f3c99c475b7d71394c531d6640d8cafc35515d7e2b0d568"
	status, err := ic.CheckOrderAddingStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckOrderWithdrawalStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "9e10c30df9a042290060561e0367ec21f6fce04c6a51c8f9276605a97643424a"
	status, err := ic.CheckOrderWithdrawalStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckDexStakingStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "c34e765399681aba33189498c37262eb1d4bb2568e0e60378a0428cbaa97f205"
	status, err := ic.CheckDEXStakingStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckDexUnStakingStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "6e2e6bb4a671c9d991cc48d211d57770a15ae90b1324b7da8552efcc57292df8"
	status, err := ic.CheckDEXUnStakingStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckDEXStakingRewardWithdrawalStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "c657f36d7cf5058dc9df694e58ec5cf02509e611830968cb58f353480cf4baa6"
	status, err := ic.CheckDEXStakingRewardWithdrawalStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckDEXLPFeeWithdrawalStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "7201105e9746aeb9479092a9f47b81445c2bb343d13b8519ed29416c246968e4"
	status, err := ic.CheckDEXLPFeeWithdrawalStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckNFTMintingStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	encodedTx, txHash, err := ic.CreatePdexv3MintNFT(privateKey)
	if err != nil {
		panic(err)
	}
	err = ic.SendRawTx(encodedTx)
	if err != nil {
		panic(err)
	}
	Logger.Printf("TxHash: %v\n", txHash)

	time.Sleep(100 * time.Second)

	status, err := ic.CheckNFTMintingStatus(txHash)
	if err != nil {
		panic(err)
	}
	jsb, _ := json.MarshalIndent(status, "", "\t")
	Logger.Printf("status: %v\n", string(jsb))
}

func TestIncClient_GetEstimatedDEXStakingReward(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	stakingPoolID := "0000000000000000000000000000000000000000000000000000000000000004"
	nftID := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	beaconHeight := uint64(0)

	res, err := ic.GetEstimatedDEXStakingReward(beaconHeight, stakingPoolID, nftID)
	if err != nil {
		panic(err)
	}
	jsb, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", string(jsb))
}

func TestIncClient_GetEstimatedLPValue(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	poolPairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-768d95970886ef2eea617d2a0ef4aea80ca394db4f8f1ff1fd724c82955cd228"
	nftID := "eb1ec0987a37829831c8d947ef2c48f8ab6ada4b02d99e82039ca5977570bd0c"
	beaconHeight := uint64(0)

	res, err := ic.GetEstimatedLPValue(beaconHeight, poolPairID, nftID)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", res)
}

func TestIncClient_GetListNftIDs(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	nftList, err := ic.GetListNftIDs(0)
	if err != nil {
		panic(err)
	}
	Logger.Println(nftList)
}

func TestIncClient_GetListStakingPoolShares(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	res, err := ic.GetListStakingPoolShares(0)
	if err != nil {
		panic(err)
	}
	jsb, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))
}

func TestIncClient_GetListStakingRewardTokens(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	res, err := ic.GetListStakingRewardTokens(0)
	if err != nil {
		panic(err)
	}
	jsb, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))
}

func TestIncClient_GetOrderByID(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	orderID := "e16109f5b96bd86fd4a1148b084d0b1caeb6dce2d5904d5ff3261e996585fceb"
	res, err := ic.GetOrderByID(0, orderID)
	if err != nil {
		panic(err)
	}
	jsb, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))
}
