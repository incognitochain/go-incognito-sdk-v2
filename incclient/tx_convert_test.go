package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"log"
	"math"
	"testing"
)

func TestIncClient_CreateConversionTransactionWithInputCoins(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	privateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	for i := 0; i < numTests; i++ {
		log.Printf("TEST %v\n", i)
		oldBalanceV1, err := getBalanceByVersion(privateKey, common.PRVIDStr, 1)
		if err != nil {
			panic(err)
		}
		log.Printf("oldBalanceV1: %v\n", oldBalanceV1)

		oldBalanceV2, err := getBalanceByVersion(privateKey, common.PRVIDStr, 2)
		if err != nil {
			panic(err)
		}
		log.Printf("oldBalanceV2: %v\n", oldBalanceV2)

		utxoList, listIndices, err := ic.GetUnspentOutputCoins(privateKey, common.PRVIDStr, 0)
		if err != nil {
			panic(err)
		}

		coinV1s, _, _, err := divideCoins(utxoList, listIndices, true)
		if err != nil {
			panic(err)
		}

		if len(coinV1s) == 0 {
			panic("no UTXO v1 to spend")
		}

		log.Printf("TESTING WITH %v INPUT COINs V1\n", len(coinV1s))

		// choose random UTXOs to spend
		coinsToSpend := coinV1s
		if len(coinV1s) > 1 {
			r := 1 + common.RandInt()%(int(math.Min(float64(len(coinV1s)-1), MaxInputSize)))
			coinsToSpend, _ = chooseRandomCoins(coinV1s, nil, r)
		}

		log.Printf("#coinsToSpend: %v\n", len(coinsToSpend))

		txFee := ic.cfg.DefaultPRVFee
		totalAmount := uint64(0)
		for _, c := range coinsToSpend {
			totalAmount += c.GetValue()
		}
		if totalAmount <= txFee {
			panic("not enough coins to spend")
		}

		// choose the sending amount
		log.Printf("totalConvertingAmount: %v, txFee: %v\n", totalAmount, txFee)

		encodedTx, txHash, err := ic.CreateConversionTransactionWithInputCoins(privateKey, coinsToSpend)
		if err != nil {
			panic(err)
		}
		err = ic.SendRawTx(encodedTx)
		if err != nil {
			panic(err)
		}
		log.Printf("TxHash created: %v\n", txHash)

		// checking if tx is in blocks
		log.Printf("Checking status of tx %v...\n", txHash)
		err = waitingCheckTxInBlock(txHash)
		if err != nil {
			panic(err)
		}

		// checking if tx has spent the coinsToSpend
		tx, err := ic.GetTx(txHash)
		if err != nil {
			panic(err)
		}
		_, err = compareInputCoins(coinsToSpend, tx.GetProof().GetInputCoins())
		if err != nil {
			panic(err)
		}
		log.Printf("Checked input coins SUCCEEDED\n")

		// checking updated balance
		expectedBalanceV1 := oldBalanceV1 - totalAmount
		expectedBalanceV2 := oldBalanceV2 + totalAmount - ic.cfg.DefaultPRVFee
		err = waitingCheckBalanceUpdated(privateKey, common.PRVIDStr, oldBalanceV1, expectedBalanceV1, 1)
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(privateKey, common.PRVIDStr, oldBalanceV2, expectedBalanceV2, 2)
		if err != nil {
			panic(err)
		}

		log.Printf("FINISHED TEST %v\n\n", i)
	}

}

func TestIncClient_CreateTokenConversionTransactionWithInputCoins(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	privateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	tokenIDStr := "f3e586e281d275ea2059e35ae434d0431947d2b49466b6d2479808378268f822"
	for i := 0; i < numTests; i++ {
		log.Printf("TEST %v\n", i)
		oldBalanceV1, err := getBalanceByVersion(privateKey, tokenIDStr, 1)
		if err != nil {
			panic(err)
		}
		log.Printf("oldBalanceV1: %v\n", oldBalanceV1)

		oldBalanceV2, err := getBalanceByVersion(privateKey, tokenIDStr, 2)
		if err != nil {
			panic(err)
		}
		log.Printf("oldBalanceV2: %v\n", oldBalanceV2)

		// choose token coins to spend
		utxoList, _, err := ic.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
		if err != nil {
			panic(err)
		}
		coinV1s, _, _, err := divideCoins(utxoList, nil, true)
		if err != nil {
			panic(err)
		}
		if len(coinV1s) == 0 {
			panic("no UTXO v1 to spend")
		}

		// choose PRV coins to spend
		prvUTXOList, prvIdxList, err := ic.GetUnspentOutputCoins(privateKey, common.PRVIDStr, 0)
		if err != nil {
			panic(err)
		}
		_, prvCoinV2s, prvIdxV2s, err := divideCoins(prvUTXOList, prvIdxList, true)
		if err != nil {
			panic(err)
		}
		if len(prvCoinV2s) == 0 {
			panic("no PRV UTXO v2 to spend")
		}

		log.Printf("TESTING WITH %v TOKEN INPUT COINs V1, %v PRV COINs V2\n", len(coinV1s), len(prvCoinV2s))

		// choose random UTXOs to spend
		coinsToSpend := coinV1s
		if len(coinV1s) > 1 {
			r := 1 + common.RandInt()%(int(math.Min(float64(len(coinV1s)-1), MaxInputSize+MaxInputSize/3)))
			coinsToSpend, _ = chooseRandomCoins(coinV1s, nil, r)
		}
		log.Printf("#coinsToSpend: %v\n", len(coinsToSpend))
		totalAmount := uint64(0)
		for _, c := range coinsToSpend {
			totalAmount += c.GetValue()
		}

		// choose random PRV UTXOs to pay fee
		prvCoinsToSpend := prvCoinV2s
		prvIdxToSpend := prvIdxV2s
		if len(prvCoinV2s) > 1 {
			r := 1 + common.RandInt()%(int(math.Min(float64(len(prvCoinV2s)-1), MaxInputSize+MaxInputSize/3)))
			prvCoinsToSpend, prvIdxToSpend = chooseRandomCoins(prvCoinV2s, prvIdxV2s, r)
		}
		log.Printf("#prvCoinsToSpend: %v\n", len(prvCoinsToSpend))

		txFee := ic.cfg.DefaultPRVFee
		totalPRVAmount := uint64(0)
		for _, c := range prvCoinsToSpend {
			totalPRVAmount += c.GetValue()
		}
		if totalPRVAmount <= txFee {
			panic(fmt.Sprintf("not enough PRV coins to spend, want %v, have %v", txFee, totalPRVAmount))
		}

		// choose the sending amount
		sendingAmount := common.RandUint64() % totalAmount
		log.Printf("SendingAmount: %v, txPRVFee: %v\n", sendingAmount, txFee)

		encodedTx, txHash, err := ic.CreateTokenConversionTransactionWithInputCoins(privateKey,
			tokenIDStr, coinsToSpend, prvCoinsToSpend, prvIdxToSpend)
		if err != nil {
			if len(coinsToSpend) > MaxInputSize || len(prvCoinsToSpend) > MaxInputSize {
				log.Printf("Should rejected SUCCEEDED\n")
				log.Printf("FINISHED TEST %v\n\n", i)
				continue
			}
			panic(err)
		}
		log.Printf("TxHash created: %v\n", txHash)
		err = ic.SendRawTokenTx(encodedTx)
		if err != nil {
			panic(err)
		}

		// checking if tx is in blocks
		log.Printf("Checking status of tx %v...\n", txHash)
		err = waitingCheckTxInBlock(txHash)
		if err != nil {
			panic(err)
		}

		// checking if tx has spent the coinsToSpend
		tx, err := ic.GetTx(txHash)
		if err != nil {
			panic(err)
		}
		tokenTx, ok := tx.(tx_generic.TransactionToken)
		if !ok {
			panic("not a token transaction")
		}
		_, err = compareInputCoins(coinsToSpend, tokenTx.GetTxNormal().GetProof().GetInputCoins())
		if err != nil {
			panic(err)
		}
		log.Printf("Checked token input coins SUCCEEDED\n")
		_, err = compareInputCoins(prvCoinsToSpend, tokenTx.GetTxBase().GetProof().GetInputCoins())
		if err != nil {
			panic(err)
		}
		log.Printf("Checked PRV input coins SUCCEEDED\n")

		// checking updated balances
		expectedBalanceV1 := oldBalanceV1 - totalAmount
		expectedBalanceV2 := oldBalanceV2 + totalAmount
		err = waitingCheckBalanceUpdated(privateKey, tokenIDStr, oldBalanceV1, expectedBalanceV1, 1)
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(privateKey, tokenIDStr, oldBalanceV2, expectedBalanceV2, 2)
		if err != nil {
			panic(err)
		}

		log.Printf("FINISHED TEST %v\n\n", i)
	}

}
