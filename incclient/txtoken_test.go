package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"log"
	"math"
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
	Logger.IsEnable = false

	privateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"

	receiverPrivateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"
	paymentAddress := PrivateKeyToPaymentAddress(receiverPrivateKey, -1)

	var tokenIDStr string
	found := false
	allBalances, _ := ic.GetAllBalancesV2(privateKey)
	for tokenID, balance := range allBalances {
		if tokenID == common.PRVIDStr {
			continue
		}
		if balance > 0 {
			tokenIDStr = tokenID
			found = true
			break
		}
	}

	if !found {
		log.Printf("Token balance is zero. Mint a new token.\n")
		mintingAmount := common.RandUint64() % uint64(10000000000000000000)
		encodedTx, _, err := ic.CreateTokenInitTransaction(privateKey, "TEST-TOKEN", "TTT", mintingAmount, 2)
		if err != nil {
			panic(err)
		}
		err = ic.SendRawTx(encodedTx)
		if err != nil {
			panic(err)
		}
		time.Sleep(20 * time.Second)

		for {
			allBalances, _ := ic.GetAllBalancesV2(privateKey)
			found := false
			for tokenID, balance := range allBalances {
				if balance == mintingAmount {
					tokenIDStr = tokenID
					found = true
					break
				}
			}
			if found {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}

	log.Printf("TESTING WITH TOKENID %v\n\n", tokenIDStr)

	for i := 0; i < numTests; i++ {
		version := 2
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

		sendingAmount := common.RandUint64() % (oldSenderBalance / 10000)
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
		err = waitingCheckBalanceUpdated(privateKey, tokenIDStr, oldSenderBalance, expectedSenderBalance, uint8(version))
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(receiverPrivateKey, tokenIDStr, oldReceiverBalance, expectedReceiverBalance, uint8(version))
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

func TestIncClient_CreateRawTokenTransactionWithInputCoinsV1WithTokenFee(t *testing.T) {
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
		log.Printf("TEST %v\n", i)
		oldSenderBalance, err := getBalanceByVersion(privateKey, tokenIDStr, 1)
		if err != nil {
			panic(err)
		}
		log.Printf("oldSenderBalance: %v\n", oldSenderBalance)

		oldReceiverBalance, err := getBalanceByVersion(receiverPrivateKey, tokenIDStr, 1)
		if err != nil {
			panic(err)
		}
		log.Printf("oldReceiverBalance: %v\n", oldReceiverBalance)

		utxoList, listIndices, err := ic.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
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

		txFee := 40 + (common.RandUint64()%10)*10
		totalAmount := uint64(0)
		for _, c := range coinsToSpend {
			totalAmount += c.GetValue()
		}
		if totalAmount <= txFee {
			panic("not enough coins to spend")
		}

		// choose the sending amount
		sendingAmount := common.RandUint64() % (totalAmount - txFee)
		log.Printf("SendingAmount: %v, txTokenFee: %v\n", sendingAmount, txFee)

		txTokenParam := NewTxTokenParam(tokenIDStr, 1, []string{paymentAddress}, []uint64{sendingAmount}, true, txFee, nil)
		txParam := NewTxParam(privateKey, nil, nil, 0, txTokenParam, nil, nil)
		encodedTx, txHash, err := ic.CreateRawTokenTransactionWithInputCoins(txParam, coinsToSpend, nil, nil, nil)
		if err != nil {
			panic(err)
		}
		err = ic.SendRawTokenTx(encodedTx)
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
		tokenTx, ok := tx.(tx_generic.TransactionToken)
		if !ok {
			panic("not a token transaction")
		}
		_, err = compareInputCoins(coinsToSpend, tokenTx.GetTxNormal().GetProof().GetInputCoins())
		if err != nil {
			panic(err)
		}
		log.Printf("Checked input coins SUCCEEDED\n")

		// checking updated balance
		expectedReceiverBalance := oldReceiverBalance + sendingAmount
		expectedSenderBalance := oldSenderBalance - sendingAmount - txFee
		if privateKey == receiverPrivateKey {
			expectedReceiverBalance = oldReceiverBalance - txFee
			expectedSenderBalance = oldSenderBalance - txFee
		}
		err = waitingCheckBalanceUpdated(privateKey, tokenIDStr, oldSenderBalance, expectedSenderBalance, 1)
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(receiverPrivateKey, tokenIDStr, oldReceiverBalance, expectedReceiverBalance, 1)
		if err != nil {
			panic(err)
		}

		log.Printf("FINISHED TEST %v\n\n", i)
	}

}

func TestIncClient_CreateRawTokenTransactionWithInputCoinsV1WithPRVFee(t *testing.T) {
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
		log.Printf("TEST %v\n", i)
		oldSenderBalance, err := getBalanceByVersion(privateKey, tokenIDStr, 1)
		if err != nil {
			panic(err)
		}
		log.Printf("oldSenderBalance: %v\n", oldSenderBalance)

		oldReceiverBalance, err := getBalanceByVersion(receiverPrivateKey, tokenIDStr, 1)
		if err != nil {
			panic(err)
		}
		log.Printf("oldReceiverBalance: %v\n", oldReceiverBalance)

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
		prvUTXOList, _, err := ic.GetUnspentOutputCoins(privateKey, common.PRVIDStr, 0)
		if err != nil {
			panic(err)
		}
		prvCoinV1s, _, _, err := divideCoins(prvUTXOList, nil, true)
		if err != nil {
			panic(err)
		}
		if len(prvCoinV1s) == 0 {
			panic("no PRV UTXO v1 to spend")
		}

		log.Printf("TESTING WITH %v INPUT COINs V1, %v PRV COINs V1\n", len(coinV1s), len(prvCoinV1s))

		// choose random UTXOs to spend
		coinsToSpend := coinV1s
		if len(coinV1s) > 1 {
			r := 1 + common.RandInt()%(int(math.Min(float64(len(coinV1s)-1), MaxInputSize)))
			coinsToSpend, _ = chooseRandomCoins(coinV1s, nil, r)
		}
		log.Printf("#coinsToSpend: %v\n", len(coinsToSpend))
		totalAmount := uint64(0)
		for _, c := range coinsToSpend {
			totalAmount += c.GetValue()
		}

		// choose random PRV UTXOs to pay fee
		prvCoinsToSpend := prvCoinV1s
		if len(prvCoinV1s) > 1 {
			r := 1 + common.RandInt()%(int(math.Min(float64(len(prvCoinV1s)-1), MaxInputSize)))
			prvCoinsToSpend, _ = chooseRandomCoins(prvCoinV1s, nil, r)
		}
		log.Printf("#prvCoinsToSpend: %v\n", len(prvCoinsToSpend))
		txFee := 40 + (common.RandUint64()%10)*10
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

		txTokenParam := NewTxTokenParam(tokenIDStr, 1, []string{paymentAddress}, []uint64{sendingAmount}, false, 0, nil)
		txParam := NewTxParam(privateKey, nil, nil, txFee, txTokenParam, nil, nil)
		encodedTx, txHash, err := ic.CreateRawTokenTransactionWithInputCoins(txParam, coinsToSpend, nil, prvCoinsToSpend, nil)
		if err != nil {
			panic(err)
		}
		err = ic.SendRawTokenTx(encodedTx)
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

		// checking updated balance
		expectedReceiverBalance := oldReceiverBalance + sendingAmount
		expectedSenderBalance := oldSenderBalance - sendingAmount
		if privateKey == receiverPrivateKey {
			expectedReceiverBalance = oldReceiverBalance
			expectedSenderBalance = oldSenderBalance
		}
		err = waitingCheckBalanceUpdated(privateKey, tokenIDStr, oldSenderBalance, expectedSenderBalance, 1)
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(receiverPrivateKey, tokenIDStr, oldReceiverBalance, expectedReceiverBalance, 1)
		if err != nil {
			panic(err)
		}

		log.Printf("FINISHED TEST %v\n\n", i)
	}

}

func TestIncClient_CreateRawTokenTransactionWithInputCoinsV2(t *testing.T) {
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
		log.Printf("TEST %v\n", i)
		oldSenderBalance, err := getBalanceByVersion(privateKey, tokenIDStr, 2)
		if err != nil {
			panic(err)
		}
		log.Printf("oldSenderBalance: %v\n", oldSenderBalance)

		oldReceiverBalance, err := getBalanceByVersion(receiverPrivateKey, tokenIDStr, 2)
		if err != nil {
			panic(err)
		}
		log.Printf("oldReceiverBalance: %v\n", oldReceiverBalance)

		// choose token coins to spend
		utxoList, idxList, err := ic.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
		if err != nil {
			panic(err)
		}
		_, coinV2s, idxV2List, err := divideCoins(utxoList, idxList, true)
		if err != nil {
			panic(err)
		}
		if len(coinV2s) == 0 {
			panic("no UTXO v2 to spend")
		}

		// choose PRV coins to spend
		prvUTXOList, prvIdxList, err := ic.GetUnspentOutputCoins(privateKey, common.PRVIDStr, 0)
		if err != nil {
			panic(err)
		}
		_, prvCoinV2s, prvIdxV2List, err := divideCoins(prvUTXOList, prvIdxList, true)
		if err != nil {
			panic(err)
		}
		if len(prvCoinV2s) == 0 {
			panic("no PRV UTXO v2 to spend")
		}

		log.Printf("TESTING WITH %v INPUT COINs V2, %v PRV COINs V2\n", len(coinV2s), len(prvCoinV2s))

		// choose random UTXOs to spend
		coinsToSpend := coinV2s
		indices := idxV2List
		if len(coinV2s) > 1 {
			r := 1 + common.RandInt()%(int(math.Min(float64(len(coinV2s)-1), MaxInputSize)))
			coinsToSpend, indices = chooseRandomCoins(coinV2s, idxV2List, r)
		}
		log.Printf("#coinsToSpend: %v\n", len(coinsToSpend))
		totalAmount := uint64(0)
		for _, c := range coinsToSpend {
			totalAmount += c.GetValue()
		}

		// choose random PRV UTXOs to pay fee
		prvCoinsToSpend := prvCoinV2s
		prvIndices := prvIdxV2List
		if len(prvCoinV2s) > 1 {
			r := 1 + common.RandInt()%(int(math.Min(float64(len(prvCoinV2s)-1), MaxInputSize)))
			prvCoinsToSpend, prvIndices = chooseRandomCoins(prvCoinV2s, prvIdxV2List, r)
		}
		log.Printf("#prvCoinsToSpend: %v\n", len(prvCoinsToSpend))
		txFee := 40 + (common.RandUint64()%10)*10
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

		txTokenParam := NewTxTokenParam(tokenIDStr, 1, []string{paymentAddress}, []uint64{sendingAmount}, false, 0, nil)
		txParam := NewTxParam(privateKey, nil, nil, txFee, txTokenParam, nil, nil)
		encodedTx, txHash, err := ic.CreateRawTokenTransactionWithInputCoins(txParam, coinsToSpend, indices, prvCoinsToSpend, prvIndices)
		if err != nil {
			panic(err)
		}
		err = ic.SendRawTokenTx(encodedTx)
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
		if tx.GetVersion() != 2 {
			panic(fmt.Sprintf("tx version must be 2, got %v", tx.GetVersion()))
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

		// checking updated balance
		expectedReceiverBalance := oldReceiverBalance + sendingAmount
		expectedSenderBalance := oldSenderBalance - sendingAmount
		if privateKey == receiverPrivateKey {
			expectedReceiverBalance = oldReceiverBalance
			expectedSenderBalance = oldSenderBalance
		}
		err = waitingCheckBalanceUpdated(privateKey, tokenIDStr, oldSenderBalance, expectedSenderBalance, 2)
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(receiverPrivateKey, tokenIDStr, oldReceiverBalance, expectedReceiverBalance, 2)
		if err != nil {
			panic(err)
		}

		log.Printf("FINISHED TEST %v\n\n", i)
	}

}
