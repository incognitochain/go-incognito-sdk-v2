package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"time"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	// Trade between some tokens
	tokenToSell := "00000000000000000000000000000000000000000000000000000000000115d7"
	tokenToBuy := common.PRVIDStr
	sellAmount := uint64(10000)
	expectedAmount := uint64(7000000)
	tradePath := []string{"00000000000000000000000000000000000000000000000000000000000115d7-00000000000000000000000000000000000000000000000000000000000115dc-aeb37b2be73b62b6b5b95086e47687767950e66772e14db6daeef01e40344dd5", "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115dc-03696365b2ff79bb9ef35bf43a74e655ffadae0fa139b8016148d7a036716c5c"}
	tradingFee := uint64(50)
	feeInPRV := false

	txHash, err := client.CreateAndSendPdexv3TradeTransaction(privateKey, tradePath, tokenToSell, tokenToBuy, sellAmount, expectedAmount, tradingFee, feeInPRV)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("txHash: %v\n", txHash)

	time.Sleep(100 * time.Second)
	status, err := client.CheckTradeStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	common.PrintJson(status, "TradeStatus")
}
