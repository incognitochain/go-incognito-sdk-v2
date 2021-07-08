package main

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewMainNetClient()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "YOUR_PRIVATE_KEY" // input your private key
	tokenIDStr := common.PRVIDStr
	version := int8(1)
	numThreads := 20

	txList, err := client.Consolidate(privateKey, tokenIDStr, version, numThreads)
	if err != nil {
		log.Printf("txList: %v\n", txList)
		log.Fatal(err)
	}
	log.Printf("txList: %v\n", txList)
}
