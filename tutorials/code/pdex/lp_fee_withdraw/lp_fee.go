package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	withdrawTokenIDs := make([]string, 0) // leave it empty if you want to withdraw all fees in the pool. Otherwise, specify which token.
	pairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"

	txHash, err := client.CreateAndSendPdexv3WithdrawLPFeeTransaction(privateKey, pairID, nftIDStr, withdrawTokenIDs...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdraw Liquidity-Provider Fee submitted in TX %v\n", txHash)
}
