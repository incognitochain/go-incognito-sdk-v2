package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := ""
	tokenIDStr := common.PRVIDStr
	nftIDStr := "941c5e6879c5f690d151b227e30bfee72e4cdbdd5709bc8ae22aa1c46b41a7df"
	amount := uint64(4300)

	txHash, err := client.CreateAndSendPdexv3StakingTransaction(privateKey, tokenIDStr, nftIDStr, amount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Staking TX submitted %v\n", txHash)
}
