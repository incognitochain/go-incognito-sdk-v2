package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	//Trade PRV to tokens
	tokenToSell := common.PRVIDStr
	tokenToBuy := "0000000000000000000000000000000000000000000000000000000000000100"
	sellAmount := uint64(500000000)
	expectedAmount, err := client.CheckXPrice(tokenToSell, tokenToBuy, sellAmount)
	if err != nil {
		log.Fatal(err)
	}
	tradingFee := uint64(10)

	txHash, err := client.CreateAndSendPDETradeTransaction(privateKey, tokenToSell, tokenToBuy, sellAmount, expectedAmount, tradingFee)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
