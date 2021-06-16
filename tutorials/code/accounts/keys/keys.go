package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func main() {
	keyWallet, err := wallet.NewMasterKeyFromSeed(common.RandBytes(32))
	if err != nil {
		log.Fatal(err)
	}

	privateKey := keyWallet.Base58CheckSerialize(wallet.PrivateKeyType)
	accInfo, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AccountInfo: %v\n", accInfo.String())
}
