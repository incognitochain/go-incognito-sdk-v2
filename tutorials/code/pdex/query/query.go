package main

import (
	"encoding/json"
	"fmt"
	"log"

	// "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	pdeState, err := client.GetPDEState(0)
	if err != nil {
		log.Fatal(err)
	}

	jsb, _ := json.Marshal(pdeState)
	fmt.Printf("pdeState: \n%s\n", string(jsb))

	// allPairs, err := client.GetAllPDEPoolPairs(0)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("allPairs: \n%v\n", allPairs)

	// tokenID1 := common.PRVIDStr
	// tokenID2 := "0000000000000000000000000000000000000000000000000000000000000100"
	// pair, err := client.GetPDEPoolPair(0, tokenID1, tokenID2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("pair: %v\n", pair)

	// tokenToSell := common.PRVIDStr
	// tokenToBuy := "0000000000000000000000000000000000000000000000000000000000000100"
	// sellAmount := uint64(1000000000)
	// expectedAmount, err := client.CheckXPrice(tokenToSell, tokenToBuy, sellAmount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Expected amount: %v\n", expectedAmount)
}
