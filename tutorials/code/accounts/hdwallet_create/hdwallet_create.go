package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func main() {
	// create a master wallet
	masterWallet, mnemonic, err := wallet.NewMasterKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("mnemonic: `%v`\n", mnemonic)

	// derive the first wallet from the master (i.e., the `Anon` account on the Incognito wallet)
	childIdx := uint32(1)
	firstWallet, err := masterWallet.DeriveChild(childIdx)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := firstWallet.GetPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("firstPrivateKey: %v\n", privateKey)

}
