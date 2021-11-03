package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	// Trade between some tokens
	tokenToSell := "0000000000000000000000000000000000000000000000000000000000000004"
	tokenToBuy := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48"
	sellAmount := uint64(5000)
	expectedAmount := uint64(1)
	tradePath := []string{"0000000000000000000000000000000000000000000000000000000000000004-0e2ceb130ec236ecb14630a63e9de3ed4bd29ee739c2253b7fc95a0857a53e93-129a559ebe38bdd71c186a91747fe3f62fa03737446bd86eedecfbd8af5c21da"}
	// expectedAmount, err := client.CheckXPrice(tokenToSell, tokenToBuy, sellAmount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	tradingFee := uint64(5000000000)
	feeInPRV := false

encodedTx, txHash, err := client.CreatePdexv3Trade(privateKey, tradePath, tokenToSell, tokenToBuy, sellAmount, expectedAmount, tradingFee, feeInPRV)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction %s\nHash %v\n", encodedTx, txHash)
}

