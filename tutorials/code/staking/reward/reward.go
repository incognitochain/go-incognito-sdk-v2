package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	rewardAddress := incclient.PrivateKeyToPaymentAddress(privateKey, -1)


	txHash, err := client.CreateAndSendWithDrawRewardTransaction(privateKey, rewardAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
