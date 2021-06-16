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
	tokenName := "INC"
	tokenSymbol := "INC"
	tokenCap := uint64(1000000000000)
	tokenVersion := 1

	encodedTx, txHash, err := client.CreateTokenInitTransaction(privateKey, tokenName, tokenSymbol, tokenCap, tokenVersion)
	if err != nil {
		log.Fatal(err)
	}

	if tokenVersion == 1 {
		err = client.SendRawTokenTx(encodedTx)
	} else {
		err = client.SendRawTx(encodedTx)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Initialize a token successfully, txHash: %v\n", txHash)
}
