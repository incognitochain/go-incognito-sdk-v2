package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"log"
)

func main() {
	ic, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	tokenIDStr := "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
	remoteAddr := "b446151522b8f1c9d27cacedce93f398a016f84337c1b79fc54c8436af5f7900"
	burnedAmount := uint64(50000000)

	// specify which EVM network we are interacting with. evmNetworkID could be one of the following:
	// 	- rpc.ETHNetworkID
	//	- rpc.BSCNetworkID
	//	- rpc.PLGNetworkID
	evmNetworkID := rpc.ETHNetworkID

	txHash, err := ic.CreateAndSendBurningRequestTransaction(privateKey, remoteAddr, tokenIDStr, burnedAmount, evmNetworkID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TxHash: %v\n", txHash)
}
