package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	// init client
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	tokenID := "00000000000000000000000000000000000000000000000000000000000000ff" // tokenID wished to convert
	txHash, err := client.CreateAndSendRawConversionTransaction(privateKey, tokenID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Initialize a token successfully, txHash: %v\n", txHash)
}
