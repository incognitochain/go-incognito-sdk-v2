package incclient

import (
	"testing"
	"time"
)

func TestIncClient_CheckTradeStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	txHash := "e4c13e368eb4da34ebcd04aaf9da9a401d5f55df752f3d1c650331a19f69a53a"
	status, err := ic.CheckTradeStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_CheckNFTMintingStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	encodedTx, txHash, err := ic.CreatePdexv3MintNFT(privateKey)
	if err != nil {
		panic(err)
	}
	err = ic.SendRawTx(encodedTx)
	if err != nil {
		panic(err)
	}
	Logger.Printf("TxHash: %v\n", txHash)

	time.Sleep(100 * time.Second)

	status, ID, err := ic.CheckNFTMintingStatus(txHash)
	if err != nil {
		panic(err)
	}
	Logger.Printf("status: %v, NftID: %v\n", status, ID)
}
