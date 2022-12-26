package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"log"
)

func main() {
	client, err := incclient.NewTestNetClientWithCache()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	receiverAddr := "12sm5BLDevJsUbevkJd7eaU9zAiuNixEeUsNnYddnsqEYTrcFMfE7aSS2J6mK3GeHbdT7LMm4VcRETaJCRzzU8xKKa1Tn2t9XcGiqWSDpG7jewQkDeRDY3czMHVEgwWGfUWMvkd2pWr1QpMw1i4s"
	w, err := wallet.Base58CheckDeserialize(receiverAddr)
	if err != nil {
		log.Fatal(err)
	}
	otaReceiver := new(coin.OTAReceiver)
	err = otaReceiver.FromAddress(w.KeySet.PaymentAddress)
	if err != nil {
		log.Fatal(err)
	}
	otaReceiverStr := otaReceiver.String(true) // set `isConcealable = true` to enable receiving within confidential transactions
	fmt.Printf("otaReceiver: %v\n", otaReceiverStr)

	txParam := incclient.NewTxParam(privateKey, []string{otaReceiverStr}, []uint64{100000}, 0, nil, nil, nil)
	encodedTx, txHash, err := client.CreateRawTransaction(txParam, -1)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create and send tx successfully, txHash: %v\n", txHash)
}
