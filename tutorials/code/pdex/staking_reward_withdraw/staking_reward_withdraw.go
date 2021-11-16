package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	tokenIDStr := common.PRVIDStr
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	// specify which token(s) in this pool to withdraw, leave it empty if withdrawing all tokens.
	withdrawTokenIDs := make([]string, 0)

	// check the current rewards
	res, err := client.GetEstimatedDEXStakingReward(0, tokenIDStr, nftIDStr)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	// withdraw the rewards
	txHash, err := client.CreateAndSendPdexv3WithdrawStakeRewardTransaction(privateKey, tokenIDStr, nftIDStr, withdrawTokenIDs...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Staking TX submitted %v\n", txHash)

	// check the withdrawing status
	time.Sleep(100 * time.Second)
	status, err := client.CheckDEXStakingRewardWithdrawalStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err = json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status: %v\n", string(jsb))
}
