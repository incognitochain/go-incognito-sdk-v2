package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	// set withdrawn amount to 0 to withdraw all remaining balance
	withdrawAmount := uint64(0)
	poolPairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	orderID := "4d033bad4ae9ef2104feda1712e2b7b7ef215b25a4e58103e6f5a29bb63fd387"
	// specify which token(s) in this pool to withdraw, leave it empty if withdrawing all tokens.
	withdrawTokenIDs := make([]string, 0)

	txHash, err := client.CreateAndSendPdexv3WithdrawOrderTransaction(privateKey, poolPairID, orderID, nftIDStr, withdrawAmount, withdrawTokenIDs...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdraw Order %s...\nSubmitted in TX %v\n", orderID, txHash)

	time.Sleep(100 * time.Second)
	status, err := client.CheckOrderWithdrawalStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status: %v\n", string(jsb))
}
