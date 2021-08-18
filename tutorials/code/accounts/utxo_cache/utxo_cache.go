package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	// create a new client
	client, err := incclient.NewTestNet1ClientWithCache()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	tokenID := "0000000000000000000000000000000000000000000000000000000000000004"

	// create a new OutCoinKey
	outCoinKey, err := incclient.NewOutCoinKeyFromPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	outCoinKey.SetReadonlyKey("") // call this if you do not want the remote full-node to decrypt your coin

	outputCoins, idxList, err := client.GetOutputCoins(outCoinKey, tokenID, 0)
	for i, outCoin := range outputCoins {
		fmt.Printf("idx: %v, version: %v, isEncrypted: %v, publicKey: %v\n", idxList[i].Uint64(),
			outCoin.GetVersion(), outCoin.IsEncrypted(), base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), 0x00))
	}
}
