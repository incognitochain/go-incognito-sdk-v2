package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"log"
	"math"
	"testing"
)

// getBalanceByVersion is for testing purposes ONLY.
func getBalanceByVersion(privateKey, tokenIDStr string, version uint8) (uint64, error) {
	if ic == nil {
		return 0, fmt.Errorf("client not initialized")
	}

	unSpentCoins, _, err := ic.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return 0, err
	}

	balance := uint64(0)
	if version == 255 {
		for _, c := range unSpentCoins {
			balance += c.GetValue()
		}
	} else {
		for _, c := range unSpentCoins {
			if c.GetVersion() == version {
				balance += c.GetValue()
			}
		}
	}

	return balance, nil
}

func TestIncClient_CreateRawTransaction(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"

	receiverPrivateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	paymentAddress := PrivateKeyToPaymentAddress(receiverPrivateKey, -1)

	for i := 0; i < numTests; i++ {
		version := 1 + common.RandInt()%2
		log.Printf("TEST %v, VERSION %v\n", i, version)
		oldSenderBalance, err := getBalanceByVersion(privateKey, common.PRVIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		log.Printf("oldSenderBalance: %v\n", oldSenderBalance)

		oldReceiverBalance, err := getBalanceByVersion(receiverPrivateKey, common.PRVIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		log.Printf("oldReceiverBalance: %v\n", oldReceiverBalance)

		sendingAmount := common.RandUint64() % (oldSenderBalance / 100)
		receiverList := []string{paymentAddress}
		amountList := []uint64{sendingAmount}
		log.Printf("sendingAmount: %v\n", sendingAmount)

		txHash, err := ic.CreateAndSendRawTransaction(privateKey, receiverList, amountList, int8(version), nil)
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
		expectedSenderBalance := oldSenderBalance - sendingAmount - DefaultPRVFee
		if privateKey == receiverPrivateKey {
			expectedReceiverBalance = oldReceiverBalance - DefaultPRVFee
			expectedSenderBalance = oldSenderBalance - DefaultPRVFee
		}
		err = waitingCheckBalanceUpdated(receiverPrivateKey, common.PRVIDStr, oldReceiverBalance, expectedReceiverBalance, uint8(version))
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(privateKey, common.PRVIDStr, oldSenderBalance, expectedSenderBalance, uint8(version))
		if err != nil {
			panic(err)
		}
		log.Printf("FINISHED TEST %v\n\n", i)
	}

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

func TestIncClient_CreateRawTransactionWithInputCoinsV1(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"

	receiverPrivateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	paymentAddress := PrivateKeyToPaymentAddress(receiverPrivateKey, -1)

	for i := 0; i < numTests; i++ {
		log.Printf("TEST %v\n", i)
		oldSenderBalance, err := getBalanceByVersion(privateKey, common.PRVIDStr, 1)
		if err != nil {
			panic(err)
		}
		log.Printf("oldSenderBalance: %v\n", oldSenderBalance)

		oldReceiverBalance, err := getBalanceByVersion(receiverPrivateKey, common.PRVIDStr, 1)
		if err != nil {
			panic(err)
		}
		log.Printf("oldReceiverBalance: %v\n", oldReceiverBalance)

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
		log.Printf("SendingAmount: %v, txFee: %v\n", sendingAmount, txFee)

		txParam := NewTxParam(privateKey, []string{paymentAddress}, []uint64{sendingAmount}, txFee, nil, nil, nil)
		encodedTx, txHash, err := ic.CreateRawTransactionWithInputCoins(txParam, coinsToSpend, nil)
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
		expectedReceiverBalance := oldReceiverBalance + sendingAmount
		expectedSenderBalance := oldSenderBalance - sendingAmount - txFee
		if privateKey == receiverPrivateKey {
			expectedReceiverBalance = oldReceiverBalance - txFee
			expectedSenderBalance = oldSenderBalance - txFee
		}
		err = waitingCheckBalanceUpdated(receiverPrivateKey, common.PRVIDStr, oldReceiverBalance, expectedReceiverBalance, 1)
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(privateKey, common.PRVIDStr, oldSenderBalance, expectedSenderBalance, 1)
		if err != nil {
			panic(err)
		}

		log.Printf("FINISHED TEST %v\n\n", i)
	}

}

func TestIncClient_CreateRawTransactionWithInputCoinsV2(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"

	receiverPrivateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	paymentAddress := PrivateKeyToPaymentAddress(receiverPrivateKey, -1)

	for i := 0; i < numTests; i++ {
		log.Printf("TEST %v\n", i)

		oldSenderBalance, err := getBalanceByVersion(privateKey, common.PRVIDStr, 2)
		if err != nil {
			panic(err)
		}
		log.Printf("oldSenderBalance: %v\n", oldSenderBalance)

		oldReceiverBalance, err := getBalanceByVersion(receiverPrivateKey, common.PRVIDStr, 2)
		if err != nil {
			panic(err)
		}
		log.Printf("oldReceiverBalance: %v\n", oldReceiverBalance)

		utxoList, listIndices, err := ic.GetUnspentOutputCoins(privateKey, common.PRVIDStr, 0)
		if err != nil {
			panic(err)
		}

		_, coinV2s, idxListV2, err := divideCoins(utxoList, listIndices, true)
		if err != nil {
			panic(err)
		}

		if len(coinV2s) == 0 {
			panic("no UTXO v2 to spend")
		}

		log.Printf("TESTING WITH %v INPUT COINs V2\n", len(coinV2s))

		// choose random UTXOs to spend
		coinsToSpend := coinV2s
		idxToSpend := idxListV2
		if len(coinV2s) > 1 {
			r := 1 + common.RandInt()%(int(math.Min(float64(len(coinV2s)-1), MaxInputSize)))
			coinsToSpend, idxToSpend = chooseRandomCoins(coinV2s, idxListV2, r)
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
		log.Printf("SendingAmount: %v, txFee: %v\n", sendingAmount, txFee)

		txParam := NewTxParam(privateKey, []string{paymentAddress}, []uint64{sendingAmount}, txFee, nil, nil, nil)
		encodedTx, txHash, err := ic.CreateRawTransactionWithInputCoins(txParam, coinsToSpend, idxToSpend)
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
		expectedReceiverBalance := oldReceiverBalance + sendingAmount
		expectedSenderBalance := oldSenderBalance - sendingAmount - txFee
		if privateKey == receiverPrivateKey {
			expectedReceiverBalance = oldReceiverBalance - txFee
			expectedSenderBalance = oldSenderBalance - txFee
		}
		err = waitingCheckBalanceUpdated(receiverPrivateKey, common.PRVIDStr, oldReceiverBalance, expectedReceiverBalance, 2)
		if err != nil {
			panic(err)
		}
		err = waitingCheckBalanceUpdated(privateKey, common.PRVIDStr, oldSenderBalance, expectedSenderBalance, 2)
		if err != nil {
			panic(err)
		}

		log.Printf("FINISHED TEST %v\n\n", i)
	}
}

func chooseRandomCoins(inputCoins []coin.PlainCoin, indices []uint64, numCoins int) ([]coin.PlainCoin, []uint64) {
	coinRes := make([]coin.PlainCoin, 0)
	var idxRes []uint64
	if indices != nil {
		idxRes = make([]uint64, 0)
	}

	usedIdx := make(map[int]bool)
	for len(coinRes) < numCoins {
		i := common.RandInt() % len(inputCoins)
		if usedIdx[i] {
			continue
		}

		coinRes = append(coinRes, inputCoins[i])
		if indices != nil {
			idxRes = append(idxRes, indices[i])
		}
		usedIdx[i] = true
	}

	return coinRes, idxRes
}

func compareInputCoins(inCoinList1 []coin.PlainCoin, inCoinList2 []coin.PlainCoin) (bool, error) {
	if len(inCoinList1) != len(inCoinList2) {
		return false, fmt.Errorf("lengths mismatch %v != %v", len(inCoinList1), len(inCoinList2))
	}
	mapCoinList1 := make(map[string]bool)
	for _, inCoin := range inCoinList1 {
		mapCoinList1[inCoin.GetKeyImage().String()] = true
	}

	for i, inCoin := range inCoinList2 {
		if _, ok := mapCoinList1[inCoin.GetKeyImage().String()]; !ok {
			return false, fmt.Errorf("compare input coins failed at index %v", i)
		}
	}

	return true, nil
}
