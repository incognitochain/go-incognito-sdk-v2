package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/tendermint/tendermint/libs/math"
	"log"
	"testing"
)

func TestIncClient_CreateRawTransaction(t *testing.T) {
	ic, err := NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rnzyZWHhboZMZYMmeMGj1nDuVNkXB3FzwpPbhnNbWcSrbytAeYjDdNLfLSJhauvzYLWM2DQkWW2hJ14BGvmFfH1iDFAxgc4ywU6qMqW"
	paymentAddress := PrivateKeyToPaymentAddress("112t8rnzyZWHhboZMZYMmeMGj1nDuVNkXB3FzwpPbhnNbWcSrbytAeYjDdNLfLSJhauvzYLWM2DQkWW2hJ14BGvmFfH1iDFAxgc4ywU6qMqW", -1)

	receiverList := []string{paymentAddress}
	amountList := []uint64{1000000}

	txHash, err := ic.CreateAndSendRawTransaction(privateKey, receiverList, amountList, 2, nil)
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

func TestIncClient_CreateRawTransactionWithInputCoinsV1(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"

	receiverPrivateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	paymentAddress := PrivateKeyToPaymentAddress(receiverPrivateKey, -1)
	oldReceiverBalance, err := ic.GetBalance(receiverPrivateKey, common.PRVIDStr)
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
		r := 1 + common.RandInt()%(math.MinInt(len(coinV1s)-1, MaxInputSize))
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

	// checking updated balance
	log.Printf("Checking balance of tx %v...\n", receiverPrivateKey)
	expectedReceiverBalance := oldReceiverBalance + sendingAmount
	if privateKey == receiverPrivateKey {
		expectedReceiverBalance = oldReceiverBalance - txFee
	}
	err = waitingCheckBalanceUpdated(receiverPrivateKey, common.PRVIDStr, oldReceiverBalance, expectedReceiverBalance)
	if err != nil {
		panic(err)
	}

	log.Printf("FINISHED V1\n\n")

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
	oldReceiverBalance, err := ic.GetBalance(receiverPrivateKey, common.PRVIDStr)
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
		r := 1 + common.RandInt()%(math.MinInt(len(coinV2s)-1, MaxInputSize))
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

	// checking updated balance
	log.Printf("Checking balance of tx %v...\n", receiverPrivateKey)
	expectedReceiverBalance := oldReceiverBalance + sendingAmount
	if privateKey == receiverPrivateKey {
		expectedReceiverBalance = oldReceiverBalance - txFee
	}
	err = waitingCheckBalanceUpdated(receiverPrivateKey, common.PRVIDStr, oldReceiverBalance, expectedReceiverBalance)
	if err != nil {
		panic(err)
	}

	log.Printf("FINISHED V2\n\n")

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
