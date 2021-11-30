package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
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
	amount := uint64(4300000)

	txHash, err := client.CreateAndSendPdexv3StakingTransaction(privateKey, tokenIDStr, nftIDStr, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Staking TX submitted %v\n", txHash)

	time.Sleep(100 * time.Second)
	status, err := client.CheckDEXStakingStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status: %v\n", string(jsb))
}
