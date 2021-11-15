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
	// addr := incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	pairID := "MyContributionID"
	pairHash := "PH1"
	firstToken := common.PRVIDStr
	secondToken := "fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50"
	firstAmount := uint64(10000)
	secondAmount := uint64(10000)
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	amplifier := uint64(30000)

	firstTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, pairID, pairHash, firstToken, nftIDStr, firstAmount, amplifier)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(60 * time.Second)	// expectedSecondAmount, err := client.CheckPrice(firstToken, secondToken, firstAmount)
	secondTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, pairID, pairHash, secondToken, nftIDStr, secondAmount, amplifier)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v - %v\n", firstTx, secondTx)

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
