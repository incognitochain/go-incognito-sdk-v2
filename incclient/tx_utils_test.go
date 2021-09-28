package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"testing"
)

func TestIncClient_GetTransactionsByReceiver(t *testing.T) {
	ic, err := NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rnzyZWHhboZMZYMmeMGj1nDuVNkXB3FzwpPbhnNbWcSrbytAeYjDdNLfLSJhauvzYLWM2DQkWW2hJ14BGvmFfH1iDFAxgc4ywU6qMqW"
	paymentAddress := PrivateKeyToPaymentAddress(privateKey, -1)

	txsReceived, err := ic.GetTransactionsByReceiver(paymentAddress)
	if err != nil {
		panic(err)
	}

	for txHash, tx := range txsReceived {
		jsb, err := json.Marshal(tx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("TxHash: %v, %v\n", txHash, tx.Hash().String())
		fmt.Printf("TxDetail: %v\n\n", string(jsb))
	}
}

func TestIncClient_GetTxHashByPublicKeys(t *testing.T) {
	ic, err := NewLocalClient("")
	if err != nil {
		panic(err)
	}

	publicKeys := []string{
		"1nzmSA2cuYMX5i8UPdeUeWiCduKrevY6TRq5iCbRjUCkYzvCu3",
		"1Yo3VXGdHuBbPzDcGRNa7aYcM39N5GCMBogQdF9Agm7FG5U1LW",
		"12Vz3da29u7oX2GCQzuKZbWESRqaHmKBvjVnUccwtieygGT3N3i",
		"12CeZQq3XFUR7KLhHyWSVkEJUkYb7JigNrch2Fdm2cPezBNkGgJ",
		"12G9MyAF8eenfor27NHJp9jrZ7zvMa3YDzFrJ8BPCkxBLZqgzgY",
	}

	txsByPubKeys, err := ic.GetTxHashByPublicKeys(publicKeys)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(txsByPubKeys, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("res: %v\n", string(jsb))
}

func TestIncClient_GetTransactionsByPublicKeys(t *testing.T) {
	ic, err := NewLocalClient("")
	if err != nil {
		panic(err)
	}

	publicKeys := []string{
		"1WP5E7xE8RkNRYZszm9uLadMypDHYdbZY9kFCbA5tUY97qAMHf",
		"12TH58GsfzFxboRSSZEXqN1Kz4BRt1ouzHTpfrtizGpfmx9Ynmh",
		"1NEiRYzhWYSX9xZvzYmD3ryvUA3RGhDnFGeLsGei8cScrZRCxx",
		"19X63aX6S8RJNZNqPQr9oSUwHACQLroQXC19z2oBUA8m95Sjr2",
		"12e5UFiiiXa4AwsiLEkoXRiFipN46yhZAbiVhS4y3BnSQzdqinS",
		"1qS2zeALEX7SZndPPyc7cFcTCxvXnK1ot89Q6jALnYA3d8BLBJ",
		"12kmiLFDSxaezVQe32Ze7D9TFm9TLs1Rx5cgh4ix5Eh4hZrvtXB",
	}

	txRecv, err := ic.GetTransactionsByPublicKeys(publicKeys)
	if err != nil {
		panic(err)
	}

	for pkStr, txMap := range txRecv {
		fmt.Printf("publicKey: %v\n", pkStr)
		for txHash, tx := range txMap {
			jsb, err := json.Marshal(tx)
			if err != nil {
				panic(err)
			}
			fmt.Printf("TxHash: %v, %v\n", txHash, tx.Hash().String())
			fmt.Printf("TxDetail: %v\n", string(jsb))
		}
		fmt.Printf("End publicKey %v\n\n", pkStr)

	}
}

func TestIncClient_GetTxHashBySerialNumbers(t *testing.T) {
	ic, err := NewLocalClient("")
	if err != nil {
		panic(err)
	}

	serialNumbers := []string{
		"12jFxrGaDfLdoxucbszHeHRyTM8Z3CdzxwAdEhwhuPcMh4MoBzT",
		"126jy6yW5NSfzpYEGUtgmuzps3ey8DRmKbJbAm77KHEEnzUW6oP",
	}

	outTxs, err := ic.GetTxHashBySerialNumbers(serialNumbers, common.PRVIDStr, 255)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(outTxs, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("res: %v\n", string(jsb))
}

func TestIncClient_GetTxs(t *testing.T) {
	var err error
	ic, err = NewIncClient("https://beta-fullnode.incognito.org/fullnode", "", 1)
	if err != nil {
		panic(err)
	}

	txHashList := []string{
		"83834123d04d2d2ceec8970a627f8557ee7737f2e037094bd3ace07e35d160dc",
		"86fd368c7550c61620493f220d60d73e00278d2f129b13f33867ba246654c37c",
	}

	txs, err := ic.GetTxs(txHashList)
	if err != nil {
		panic(err)
	}

	for txHash, tx := range txs {
		txBytes, err := json.Marshal(tx)
		if err != nil {
			panic(err)
		}

		fmt.Printf("txHash %v: %v\n", txHash, string(txBytes))

	}
}

func TestIncClient_GetReceivingInfo(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	receiverOtaKey := "14y77fjaSSxc7LNcy1rf5scTkucn2CndPSRE8NTxkxB7hJfs5s6cAsK5jh1sJWWKGJYKoConDPJ1sLUqh59oTq6DVHd1Rc7S9sCtJwi"
	receiverReadonlyKey := "13hVXWMfgetD6ci9LwgJFBdXqBaapBo7HEDQsPoNg13smvjTn1rKFKpJzSvX96MzDR524Ng7m9RJ5ZYBipVxvkFnJxwHbEqEYW7MRai"
	txHash := "8afd4009134a6d30e46b1b2fc6322d93e84f242001e36b13a54446d4e337ae93"

	received, receivingInfo, err := ic.GetReceivingInfo(txHash, receiverOtaKey, receiverReadonlyKey)
	if err != nil {
		panic(err)
	}

	Logger.Printf("received: %v\n", received)
	Logger.Printf("receivingInfo: %v\n", receivingInfo)

}
