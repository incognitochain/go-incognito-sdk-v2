package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestIncClient_CreateRawTokenTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rnzyZWHhboZMZYMmeMGj1nDuVNkXB3FzwpPbhnNbWcSrbytAeYjDdNLfLSJhauvzYLWM2DQkWW2hJ14BGvmFfH1iDFAxgc4ywU6qMqW"
	paymentAddress := PrivateKeyToPaymentAddress("112t8rnzyZWHhboZMZYMmeMGj1nDuVNkXB3FzwpPbhnNbWcSrbytAeYjDdNLfLSJhauvzYLWM2DQkWW2hJ14BGvmFfH1iDFAxgc4ywU6qMqW", -1)
	tokenID := "974ff9005a6769b4159b5b3e718f12cd8218673797870cb95d76784addf65066"

	receiverList := []string{paymentAddress}
	amountList := []uint64{1000000}

	txHash, err := ic.CreateAndSendRawTokenTransaction(privateKey, receiverList, amountList, tokenID, 2, nil)
	if err != nil {
		panic(err)
	}

	log.Printf("TxHash: %v\n", txHash)
}

func TestIncClient_CreateTokenInitTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rnzyZWHhboZMZYMmeMGj1nDuVNkXB3FzwpPbhnNbWcSrbytAeYjDdNLfLSJhauvzYLWM2DQkWW2hJ14BGvmFfH1iDFAxgc4ywU6qMqW"
	shardID := GetShardIDFromPrivateKey(privateKey)

	initAmount := common.RandUint64() % uint64(1000000000000000)
	tokenName := "INC_" + common.RandChars(5)
	log.Printf("tokenName: %v\n", tokenName)
	tokenSymbol := "INC_" + common.RandChars(5)
	log.Printf("tokenSymbol: %v\n", tokenSymbol)

	// create a token init transaction
	encodedTx, txHash, err := ic.CreateTokenInitTransaction(privateKey, tokenName, tokenSymbol,
		initAmount, 2)
	if err != nil {
		panic(err)
	}

	// broadcast the transaction to the network
	err = ic.SendRawTx(encodedTx)
	if err != nil {
		panic(err)
	}

	log.Printf("TxHash: %v\n", txHash)

	tokenID := common.HashH([]byte(txHash + strconv.FormatUint(uint64(shardID), 10)))
	log.Printf("generated tokenID: %v\n", tokenID.String())

	err = waitingCheckTxInBlock(txHash)
	if err != nil {
		panic(err)
	}

	for i := 0; i < MaxAttempts; i++ {
		log.Printf("Attempt %v\n", i)
		balance, err := ic.GetBalance(privateKey, tokenID.String())
		if err != nil {
			panic(err)
		}

		if balance == initAmount {
			log.Printf("balance updated correctly!!\n")
			break
		}
		time.Sleep(10 * time.Second)
	}
}
