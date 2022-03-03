package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"log"
	"time"
)

func main() {
	ic, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	remoteAddr := "b446151522b8f1c9d27cacedce93f398a016f84337c1b79fc54c8436af5f7900"
	burnedAmount := uint64(50000000)

	// specify which EVM network we are interacting with. evmNetworkID could be one of the following:
	// 	- rpc.ETHNetworkID
	//	- rpc.BSCNetworkID
	evmNetworkID := rpc.ETHNetworkID

	// burn PRV
	txHash, err := ic.CreateAndSendBurningPRVPeggingRequestTransaction(privateKey, remoteAddr, burnedAmount, evmNetworkID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TxHash: %v\n", txHash)

	// wait for the above tx to reach the beacon chain
	time.Sleep(50 * time.Second)

	// retrieve the burn proof
	prvBurnProof, err := ic.GetBurnPRVPeggingProof(txHash, evmNetworkID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prvBurnProof)
}
