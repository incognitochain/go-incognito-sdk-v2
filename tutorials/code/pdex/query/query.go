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

	pdeState, err := client.GetPdexState(0, nil)
	if err != nil {
		log.Fatal(err)
	}

	common.PrintJson(pdeState, "Full Pdex State")

	allPairs, err := client.GetAllPdexPoolPairs(0)
	if err != nil {
		log.Fatal(err)
	}

	common.PrintJson(allPairs, "Pool Pairs")

	tokenID1 := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48"
	tokenID2 := "fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50"
	pair, err := client.GetPdexPoolPair(0, tokenID1, tokenID2)
	if err != nil {
		log.Fatal(err)
	}
	common.PrintJson(pair, "Found Pool(s) for tokenIDs")

	tokenToSell := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48"
	pairID := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48-fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50-579050a8274e029a567debf87a17725baff195ed155d7f01d2bd62d8d77fdc3d"
	sellAmount := uint64(10000)
	expectedAmount, err := client.CheckPrice(pairID, tokenToSell, sellAmount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Expected amount: %v\n", expectedAmount)
}
