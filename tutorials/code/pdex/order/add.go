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
	privateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	tokenToSell := "fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50"
	sellAmount := uint64(500)
	expectedAmount := uint64(600)
	pairID := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48-fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50-579050a8274e029a567debf87a17725baff195ed155d7f01d2bd62d8d77fdc3d"
	nftIDStr := "941c5e6879c5f690d151b227e30bfee72e4cdbdd5709bc8ae22aa1c46b41a7df"

	txHash, err := client.CreateAndSendPdexv3AddOrderTransaction(privateKey, pairID, tokenToSell, nftIDStr, sellAmount, expectedAmount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Add New Order: submitted in TX %v\n", txHash)
}
