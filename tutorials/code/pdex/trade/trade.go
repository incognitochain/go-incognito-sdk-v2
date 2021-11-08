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
	privateKey := ""
	// Trade between some tokens
	tokenToSell := "fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50"
	tokenToBuy := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48"
	sellAmount := uint64(500)
	expectedAmount := uint64(470)
	tradePath := []string{"3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48-fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50-be93d713532275875bbe5d9411f7e1e2634355b8aeb1039b4b83e3468839c1c4"}
	// expectedAmount, err := client.CheckXPrice(tokenToSell, tokenToBuy, sellAmount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	tradingFee := uint64(50)
	feeInPRV := false

	txHash, err := client.CreateAndSendPdexv3TradeTransaction(privateKey, tradePath, tokenToSell, tokenToBuy, sellAmount, expectedAmount, tradingFee, feeInPRV)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
