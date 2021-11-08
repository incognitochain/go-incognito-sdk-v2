package incclient

import (
	"testing"
	"time"
)

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
