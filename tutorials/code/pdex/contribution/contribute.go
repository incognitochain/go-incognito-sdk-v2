package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	addr := incclient.PrivateKeyToPaymentAddress(privateKey, -1)

	pairID := "newPairID"
	firstToken := common.PRVIDStr
	secondToken := "0000000000000000000000000000000000000000000000000000000000000100"
	firstAmount := uint64(1000000000)
	secondAmount := uint64(1000000000)

	expectedSecondAmount, err := client.CheckPrice(firstToken, secondToken, firstAmount)
	if err == nil {
		secondAmount = expectedSecondAmount
	} else {
		log.Println("pool has bot been initialized")
	}

	firstTx, err := client.CreateAndSendPDEContributeTransaction(privateKey, pairID, firstToken, firstAmount, 2)
	if err != nil {
		log.Fatal(err)
	}

	secondTx, err := client.CreateAndSendPDEContributeTransaction(privateKey, pairID, secondToken, secondAmount, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v - %v\n", firstTx, secondTx)

	fmt.Println("checking share updated...")
	time.Sleep(60 * time.Second)

	// retrieve contributed share
	myShare, err := client.GetShareAmount(0, firstToken, secondToken, addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("myShare: %v\n", myShare)
}
