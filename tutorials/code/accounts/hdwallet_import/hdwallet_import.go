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
	mnemonic := "search trophy awake proud sponsor toe lumber toilet sugar smoke soup joke"

	wallets, err := client.ImportAccount(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	for i, w := range wallets {
		privateKey, err := w.GetPrivateKey()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("idx: %v, privateKey: %v\n", i, privateKey)
	}
}
