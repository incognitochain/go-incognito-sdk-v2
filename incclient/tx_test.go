package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"testing"
)

func TestIncClient_CreateRawTransaction(t *testing.T) {
	ic, err := NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	paymentAddress := PrivateKeyToPaymentAddress("112t8rnzyZWHhboZMZYMmeMGj1nDuVNkXB3FzwpPbhnNbWcSrbytAeYjDdNLfLSJhauvzYLWM2DQkWW2hJ14BGvmFfH1iDFAxgc4ywU6qMqW", -1)

	receiverList := []string{paymentAddress}
	amountList := []uint64{1000000}

	txHash, err := ic.CreateAndSendRawTransaction(privateKey, receiverList, amountList, -1, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("TxHash: %v\n", txHash)
}

func TestIncClient_GetTx(t *testing.T) {
	ic, err := NewDevNetClient()
	if err != nil {
		panic(err)
	}

	var txHash string
	var tx metadata.Transaction
	var jsb []byte

	////TxNormal
	txHash = "5012d9c28f42e597e93a4695c5de16b3f44bb0acf8101ab4e6ebf6ec777b5101"

	tx, err = ic.GetTx(txHash)
	if err != nil {
		panic(err)
	}

	jsb, err = json.Marshal(tx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tx Normal: %v\n", string(jsb))

	//TxToken
	txHash = "b1129d473c2bd81646d7d348cdeb15a77066ae4fa378a510dd63973a583de8fb"

	tx, err = ic.GetTx(txHash)
	if err != nil {
		panic(err)
	}

	jsb, err = json.Marshal(tx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tx Token: %v\n", string(jsb))
}
