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
	tokenIDStr := "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
	evmTxHash := "0xb31d963b3f183d60532ca60d534e0113ca56070af795fde450dd456945a7be42"

	// specify which EVM network we are interacting with. evmNetworkID could be one of the following:
	// 	- rpc.ETHNetworkID
	//	- rpc.BSCNetworkID
	//	- rpc.PLGNetworkID
	evmNetworkID := rpc.ETHNetworkID

	evmProof, depositAmount, err := ic.GetEVMDepositProof(evmTxHash, evmNetworkID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deposited amount: %v\n", depositAmount)

	txHashStr, err := ic.CreateAndSendIssuingEVMRequestTransaction(privateKey, tokenIDStr, *evmProof, evmNetworkID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TxHash: %v\n", txHashStr)

	time.Sleep(10 * time.Second)

	fmt.Printf("check shielding status\n")
	for {
		status, err := ic.CheckShieldStatus(txHashStr)
		if err != nil {
			log.Fatal(err)
		}
		if status == 1 || status == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Printf("shielding status: %v\n", status)
		break
	}
}
