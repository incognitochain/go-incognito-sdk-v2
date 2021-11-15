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
	pairID := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48-fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50-579050a8274e029a567debf87a17725baff195ed155d7f01d2bd62d8d77fdc3d"

	txHash, err := client.CreateAndSendPdexv3WithdrawProtocolFeeTransaction(privateKey, pairID)
	if err != nil {
		log.Fatal(err)
	}

	addr := incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	fmt.Printf("Withdraw Protocol Fee to address %s...\nSubmitted in TX %v\n", addr, txHash)
}
