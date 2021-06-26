package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	tokenID := "0000000000000000000000000000000000000000000000000000000000000100"
	tokenReceivers := []string{"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H"}
	tokenAmounts := []uint64{10000000}
	hasTokenFee := false
	tokenFee := uint64(0)
	txVersion := int8(1)

	txTokenParam := incclient.NewTxTokenParam(tokenID, 1, tokenReceivers, tokenAmounts, hasTokenFee, tokenFee, nil)
	txParam := incclient.NewTxParam(privateKey, nil, nil, 0, txTokenParam, nil, nil)

	encodedTx, txHash, err := client.CreateRawTokenTransaction(txParam, txVersion)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create and send tx token successfull, txhash: %v\n", txHash)
}
