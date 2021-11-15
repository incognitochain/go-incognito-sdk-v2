package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	// burn some PRV to get your NFTID to use in pdex operations
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	txHash, err := client.CreateAndSendPdexv3UserMintNFTransaction(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Mint-NFT submitted in TX %v\n", txHash)

	// check the minting status
	time.Sleep(100 * time.Second)
	status, err := client.CheckNFTMintingStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status: %v\n", string(jsb))
}
