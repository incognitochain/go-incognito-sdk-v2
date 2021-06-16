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
	privateSeed :=  incclient.PrivateKeyToMiningKey(privateKey) //NOTE: the private seed (a.k.a the mining key) can be randomly generated and not be dependent on the private key
	candidateAddress := incclient.PrivateKeyToPaymentAddress(privateKey, -1)

	txHash, err := client.CreateAndSendUnStakingTransaction(privateKey, privateSeed, candidateAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}