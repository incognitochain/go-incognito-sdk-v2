package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	firstToken := common.PRVIDStr
	secondToken := "00000000000000000000000000000000000000000000000000000000000115d7"
	pairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	sharedAmount := uint64(5000)

	txHash, err := client.CreateAndSendPdexv3WithdrawLiquidityTransaction(privateKey, pairID, firstToken, secondToken, nftIDStr, sharedAmount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdrawal transaction %v\n", txHash)
}
