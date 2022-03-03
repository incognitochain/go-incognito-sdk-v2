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
	poolPairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
	tokenToSell := common.PRVIDStr
	tokenToBuy := "00000000000000000000000000000000000000000000000000000000000115d7"
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	sellAmount := uint64(100000)
	minAcceptableAmount := uint64(100000)

	txHash, err := client.CreateAndSendPdexv3AddOrderTransaction(privateKey,
		poolPairID, tokenToSell, tokenToBuy, nftIDStr,
		sellAmount, minAcceptableAmount,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("txHash: %v\n", txHash)

	time.Sleep(100 * time.Second)
	status, err := client.CheckOrderAddingStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status: %v\n", string(jsb))
}
