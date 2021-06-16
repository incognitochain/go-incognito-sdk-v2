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
	receiverList := []string{"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H"}
	amountList := []uint64{10000000}
	txVersion := int8(1) //txVersion should be -1, 1, or 2
	txFee := uint64(10)

	txParam := incclient.NewTxParam(privateKey, receiverList, amountList, txFee, nil, nil, nil)

	encodedTx, txHash, err := client.CreateRawTransaction(txParam, txVersion)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create and send tx successfully, txHash: %v\n", txHash)
}