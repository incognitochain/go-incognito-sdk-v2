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
	// set withdrawn amount to 0 to withdraw all remaining balance
	withdrawAmount := uint64(0)
	// specify which token(s) in this pool to withdraw
	withdrawTokenIDs := []string{"fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50", "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48"}
	pairID := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48-fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50-579050a8274e029a567debf87a17725baff195ed155d7f01d2bd62d8d77fdc3d"
	nftIDStr := "941c5e6879c5f690d151b227e30bfee72e4cdbdd5709bc8ae22aa1c46b41a7df"
	orderID := "b7ef57b7c2837934279036f70f57199bd02f17d4ade7626508434911b36056c9"

	txHash, err := client.CreateAndSendPdexv3WithdrawOrderTransaction(privateKey, pairID, orderID, withdrawTokenIDs, nftIDStr, withdrawAmount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdraw Order %s...\nSubmitted in TX %v\n", orderID, txHash)
}
