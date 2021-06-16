package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	firstToken := common.PRVIDStr
	secondToken := "0000000000000000000000000000000000000000000000000000000000000100"
	addr := incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	sharedAmount, err := client.GetShareAmount(0, firstToken, secondToken, addr) // get our current shared amount

	txHash, err := client.CreateAndSendPDEWithdrawalTransaction(privateKey, firstToken, secondToken, sharedAmount, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdrawal transaction %v\n", txHash)
}
