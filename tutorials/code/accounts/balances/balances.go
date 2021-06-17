package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	incClient, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	balancePRV, err := incClient.GetBalance(privateKey, common.PRVIDStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("balancePRV: %v\n", balancePRV)

	tokenID := "0000000000000000000000000000000000000000000000000000000000000100"
	balanceToken, err := incClient.GetBalance(privateKey, tokenID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("balanceToken: %v\n", balanceToken)
}
