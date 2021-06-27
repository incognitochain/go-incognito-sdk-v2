package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestIncClient_CreateRawTokenTransaction(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"

	receiverPrivateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	paymentAddress := PrivateKeyToPaymentAddress(receiverPrivateKey, -1)
	tokenIDStr := "f3e586e281d275ea2059e35ae434d0431947d2b49466b6d2479808378268f822"

	for i := 0; i < numTests; i++ {
		version := 1 + common.RandInt()%2
		log.Printf("TEST %v, VERSION %v\n", i, version)

		oldSenderBalance, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		log.Printf("oldSenderBalance: %v\n", oldSenderBalance)

		oldReceiverBalance, err := getBalanceByVersion(receiverPrivateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		log.Printf("oldReceiverBalance: %v\n", oldReceiverBalance)

		sendingAmount := common.RandUint64() % (oldSenderBalance / 100)
		receiverList := []string{paymentAddress}
		amountList := []uint64{sendingAmount}
		log.Printf("sendingAmount: %v\n", sendingAmount)

		txHash, err := ic.CreateAndSendRawTokenTransaction(privateKey, receiverList, amountList, tokenIDStr, int8(version), nil)
		if err != nil {
			panic(err)
		}

		fmt.Printf("TxHash: %v\n", txHash)

		// checking if tx is in blocks
		log.Printf("Checking status of tx %v...\n", txHash)
		err = waitingCheckTxInBlock(txHash)
		if err != nil {
			panic(err)
		}

		// checking updated balance
		log.Printf("Checking balance of tx %v...\n", receiverPrivateKey)
		expectedReceiverBalance := oldReceiverBalance + sendingAmount
		expectedSenderBalance := oldSenderBalance - sendingAmount
		err = waitingCheckBalanceUpdated(receiverPrivateKey, tokenIDStr, oldReceiverBalance, expectedReceiverBalance, uint8(version))
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(privateKey, tokenIDStr, oldSenderBalance, expectedSenderBalance, uint8(version))
		if err != nil {
			panic(err)
		}
		log.Printf("FINISHED TEST %v\n\n", i)
	}

}

func TestIncClient_CreateTokenInitTransaction(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"
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

	for i := 0; i < maxAttempts; i++ {
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
