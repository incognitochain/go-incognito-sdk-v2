package main

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := ""
	// addr := incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	pairID := ""
	pairHash := "PH1"
	firstToken := "3609431c4404eb5fd91607f5afcb427afe02c9cf2ff64bf0970880eb56c03b48"
	secondToken := "fd0febf5a30be293a3e241aeb860ce843f49415ac5914e4e96b428e195af9d50"
	firstAmount := uint64(10000)
	secondAmount := uint64(10000)
	nftIDStr := "941c5e6879c5f690d151b227e30bfee72e4cdbdd5709bc8ae22aa1c46b41a7df"
	amplifier := uint64(30000)

	// expectedSecondAmount, err := client.CheckPrice(firstToken, secondToken, firstAmount)
	// if err == nil {
	// 	secondAmount = expectedSecondAmount
	// } else {
	// 	log.Println("pool has not been initialized")
	// }

	firstTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, pairID, pairHash, firstToken, nftIDStr, firstAmount, amplifier)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)
	secondTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, pairID, pairHash, secondToken, nftIDStr, secondAmount, amplifier)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v - %v\n", firstTx, secondTx)

	fmt.Println("checking share updated...")
	time.Sleep(60 * time.Second)

	// retrieve contributed share
	// myShare, err := client.GetShareAmount(0, firstToken, secondToken, addr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("myShare: %v\n", myShare)
}
