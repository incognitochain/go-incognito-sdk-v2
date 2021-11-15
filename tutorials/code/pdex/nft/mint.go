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
	// burn some PRV to get your NFTID to use in pdex operations
	privateKey := ""

	txHash, err := client.CreateAndSendPdexv3UserMintNFTransaction(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Mint-NFT submitted in TX %v\n", txHash)
}
