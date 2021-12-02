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
	poolPairID := ""                // for pool-initializing, leave it empty. Otherwise, input the poolPairID of the existing pool
	pairHash := "JUSTARANDOMSTRING" // a string to match the two transactions of the contribution
	firstToken := common.PRVIDStr
	secondToken := "00000000000000000000000000000000000000000000000000000000000115d7"
	firstAmount := uint64(3000)
	secondAmount := uint64(3000)
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	amplifier := uint64(15000)

	firstTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, poolPairID, pairHash, firstToken, nftIDStr, firstAmount, amplifier)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("firstTx: %v\n", firstTx)

	// wait for the first transaction to be confirmed, so the nftID has been re-minted to proceed.
	//time.Sleep(60 * time.Second)
	secondTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, poolPairID, pairHash, secondToken, nftIDStr, secondAmount, amplifier)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("secondTx: %v\n", secondTx)

	// check the minting status
	time.Sleep(100 * time.Second)
	status, err := client.CheckDEXLiquidityContributionStatus(firstTx)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("firstTxstatus: %v\n", string(jsb))

	status, err = client.CheckDEXLiquidityContributionStatus(secondTx)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err = json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("secondTxstatus: %v\n", string(jsb))
}
