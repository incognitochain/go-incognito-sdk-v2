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
	privateKey := ""
	// Trade between some tokens
	withdrawTokenIDs := []string{"fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50", "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48"}
	pairID := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48-fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50-be93d713532275875bbe5d9411f7e1e2634355b8aeb1039b4b83e3468839c1c4"
	nftIDStr := "941c5e6879c5f690d151b227e30bfee72e4cdbdd5709bc8ae22aa1c46b41a7df"

	txHash, err := client.CreateAndSendPdexv3WithdrawLPFeeTransaction(privateKey, pairID, withdrawTokenIDs, nftIDStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdraw Liquidity-Provider Fee submitted in TX %v\n", txHash)
}
