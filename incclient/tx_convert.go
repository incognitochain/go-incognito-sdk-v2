package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// CreateRawConversionTransaction creates a PRV transaction that converts PRV coins version 1 to version 2.
// This type of transactions is non-private by default.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawConversionTransaction(privateKey string) ([]byte, string, error) {
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	//Get list of UTXOs
	utxoList, _, err := client.GetUnspentOutputCoins(privateKey, common.PRVIDStr, 0)
	if err != nil {
		return nil, "", err
	}

	//Get list of coinV1 to convert.
	coinV1List, _, _, err := divideCoins(utxoList, nil, true)
	if err != nil {
		return nil, "", fmt.Errorf("cannot divide coin: %v", err)
	}

	if len(coinV1List) == 0 {
		return nil, "", fmt.Errorf("no CoinV1 left to be converted")
	}

	//Calculating the total amount being converted.
	totalAmount := uint64(0)
	for _, utxo := range coinV1List {
		totalAmount += utxo.GetValue()
	}
	if totalAmount < DefaultPRVFee {
		fmt.Printf("Total amount (%v) is less than txFee (%v).\n", totalAmount, DefaultPRVFee)
		return nil, "", fmt.Errorf("Total amount (%v) is less than txFee (%v).\n", totalAmount, DefaultPRVFee)
	}
	totalAmount -= DefaultPRVFee

	uniquePayment := key.PaymentInfo{PaymentAddress: senderWallet.KeySet.PaymentAddress, Amount: totalAmount, Message: []byte{}}

	//Create tx conversion params
	txParam := tx_ver2.NewTxConvertVer1ToVer2InitParams(&(senderWallet.KeySet.PrivateKey), []*key.PaymentInfo{&uniquePayment}, coinV1List,
		DefaultPRVFee, nil, nil, nil, nil)

	tx := new(tx_ver2.Tx)
	err = tx_ver2.InitConversion(tx, txParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txconvert error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txconvert: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateRawTokenConversionTransaction creates a token transaction that converts token UTXOs version 1 to version 2.
// This type of transactions is non-private by default.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTokenConversionTransaction(privateKey, tokenIDStr string) ([]byte, string, error) {
	if tokenIDStr == common.PRVIDStr {
		return nil, "", fmt.Errorf("try conversion transaction")
	}

	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid token ID: %v", tokenID)
	}

	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	//We need to use PRV coinV2 to pay fee (it's a must)
	prvFee := DefaultPRVFee
	coinsToSpendPRV, kvArgsPRV, err := client.initParams(privateKey, common.PRVIDStr, prvFee, true, 2)
	if err != nil {
		return nil, "", err
	}

	//Get list of UTXOs
	utxoListToken, _, err := client.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, "", err
	}

	//We only need to convert token version 1
	coinV1ListToken, _, _, err := divideCoins(utxoListToken, nil, true)
	if err != nil {
		return nil, "", fmt.Errorf("cannot divide coin: %v", err)
	}

	//Calculate the total token amount to be converted
	totalAmount := uint64(0)
	for _, utxo := range coinV1ListToken {
		totalAmount += utxo.GetValue()
	}

	//Create unique receiver for token
	uniquePayment := key.PaymentInfo{PaymentAddress: senderWallet.KeySet.PaymentAddress, Amount: totalAmount, Message: []byte{}}

	txTokenParam := tx_ver2.NewTxTokenConvertVer1ToVer2InitParams(&(senderWallet.KeySet.PrivateKey), coinsToSpendPRV, []*key.PaymentInfo{}, coinV1ListToken,
		[]*key.PaymentInfo{&uniquePayment}, prvFee, tokenID,
		nil, nil, kvArgsPRV)

	tx := new(tx_ver2.TxToken)
	err = tx_ver2.InitTokenConversion(tx, txTokenParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txtokenconversion error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txtokenconversion: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateAndSendRawConversionTransaction creates a PRV transaction that converts PRV coins version 1 to version 2 and broadcasts it to the network.
// This type of transactions is non-private by default.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendRawConversionTransaction(privateKey string, tokenID string) (string, error) {
	var txHash string
	var err error
	var encodedTx []byte

	if tokenID == common.PRVIDStr {
		encodedTx, txHash, err = client.CreateRawConversionTransaction(privateKey)
		if err != nil {
			return "", err
		}

		err = client.SendRawTx(encodedTx)
		if err != nil {
			return "", err
		}
	} else {
		encodedTx, txHash, err = client.CreateRawTokenConversionTransaction(privateKey, tokenID)
		if err != nil {
			return "", err
		}

		err = client.SendRawTokenTx(encodedTx)
		if err != nil {
			return "", err
		}
	}

	return txHash, nil
}

// CreateConversionTransactionWithInputCoins convert a list of PRV UTXOs V1 into PRV UTXOs v2.
// Parameters:
//	- privateKey: the private key of the user.
//	- inputCoins: a list of decrypted, unspent PRV output coins (with the same version).
//
// This function uses the DefaultPRVFee to pay the transaction fee.
//
// NOTE: this servers PRV transactions only.
func (client *IncClient) CreateConversionTransactionWithInputCoins(privateKey string, coinV1List []coin.PlainCoin) ([]byte, string, error) {
	var txHash string

	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, txHash, fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	// check version of coins
	version, err := getVersionFromInputCoins(coinV1List)
	if err != nil {
		return nil, txHash, err
	}
	if version != 1 {
		return nil, txHash, fmt.Errorf("input coins must be of version 1")
	}

	// check number of input coins
	if len(coinV1List) > MaxInputSize {
		return nil, txHash, fmt.Errorf("support at most %v input coins, got %v", MaxInputSize, len(coinV1List))
	}
	if len(coinV1List) == 0 {
		return nil, txHash, fmt.Errorf("no CoinV1 to be converted")
	}

	//Calculating the total amount being converted.
	totalAmount := uint64(0)
	for _, utxo := range coinV1List {
		totalAmount += utxo.GetValue()
	}
	if totalAmount < DefaultPRVFee {
		fmt.Printf("Total amount (%v) is less than txFee (%v).\n", totalAmount, DefaultPRVFee)
		return nil, txHash, fmt.Errorf("Total amount (%v) is less than txFee (%v).\n", totalAmount, DefaultPRVFee)
	}
	totalAmount -= DefaultPRVFee

	uniquePayment := key.PaymentInfo{PaymentAddress: senderWallet.KeySet.PaymentAddress, Amount: totalAmount, Message: []byte{}}

	//Create tx conversion params
	txParam := tx_ver2.NewTxConvertVer1ToVer2InitParams(&(senderWallet.KeySet.PrivateKey), []*key.PaymentInfo{&uniquePayment}, coinV1List,
		DefaultPRVFee, nil, nil, nil, nil)

	tx := new(tx_ver2.Tx)
	err = tx_ver2.InitConversion(tx, txParam)
	if err != nil {
		return nil, txHash, fmt.Errorf("init txconvert error: %v", err)
	}
	txHash = tx.Hash().String()

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, txHash, fmt.Errorf("cannot marshal txconvert: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), txHash, nil
}

// CreateTokenConversionTransactionWithInputCoins convert a list of token UTXOs V1 into PRV UTXOs v2.
// Parameters:
//	- privateKey: the private key of the user.
//	- tokenIDStr: the id of the asset being converted.
//	- tokenInCoins: a list of decrypted, unspent token output coins v1.
//	- prvInCoins: a list of decrypted, unspent PRV output coins v2 for paying the transaction fee.
//	- prvIndices: a list of corresponding indices for the prv input coins.
//
// This function uses the DefaultPRVFee to pay the transaction fee.
//
// NOTE: this servers PRV transactions only.
func (client *IncClient) CreateTokenConversionTransactionWithInputCoins(privateKey,
	tokenIDStr string,
	tokenInCoins []coin.PlainCoin,
	prvInCoins []coin.PlainCoin,
	prvIndices []uint64,
) ([]byte, string, error) {
	var txHash string

	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, txHash, fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}
	shardID := GetShardIDFromPrivateKey(privateKey)

	// check number of token input coins
	if len(tokenInCoins) > MaxInputSize {
		return nil, txHash, fmt.Errorf("support at most %v token input coins, got %v", MaxInputSize, len(tokenInCoins))
	}
	if len(tokenInCoins) == 0 {
		return nil, txHash, fmt.Errorf("no token CoinV1 to be converted")
	}

	// check version of token input coins
	version, err := getVersionFromInputCoins(tokenInCoins)
	if err != nil {
		return nil, txHash, err
	}
	if version != 1 {
		return nil, txHash, fmt.Errorf("token input coins must be of version 1")
	}

	// check number of token input coins
	if len(prvInCoins) > MaxInputSize {
		return nil, txHash, fmt.Errorf("support at most %v PRV input coins, got %v", MaxInputSize, len(prvInCoins))
	}
	if len(prvInCoins) == 0 {
		return nil, txHash, fmt.Errorf("no PRV CoinV2 to pay fee")
	}

	// check version of PRV input coins
	version, err = getVersionFromInputCoins(prvInCoins)
	if err != nil {
		return nil, txHash, err
	}
	if version != 2 {
		return nil, txHash, fmt.Errorf("PRV input coins must be of version 2")
	}
	if len(prvIndices) != len(prvInCoins) {
		return nil, txHash, fmt.Errorf("need %v PRV indices, got %v", len(prvInCoins), len(prvIndices))
	}

	// check and parse tokenID
	if tokenIDStr == common.PRVIDStr {
		return nil, txHash, fmt.Errorf("try conversion transaction")
	}
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, txHash, fmt.Errorf("invalid token ID: %v", tokenID)
	}

	//We need to use PRV coinV2 to pay fee (it's a must)
	prvFee := DefaultPRVFee
	totalPRVAmount := uint64(0)
	for _, utxo := range prvInCoins {
		totalPRVAmount += utxo.GetValue()
	}
	if totalPRVAmount < prvFee {
		return nil, txHash, fmt.Errorf("not enough PRV to pay fee, need %v, got %v", DefaultPRVFee, totalPRVAmount)
	}

	//Calculate the total token amount to be converted
	totalAmount := uint64(0)
	for _, utxo := range tokenInCoins {
		totalAmount += utxo.GetValue()
	}

	// init PRV parameters
	kvArgs, err := client.getRandomCommitmentV2(shardID, common.PRVIDStr, len(prvInCoins)*(privacy.RingSize-1))
	if err != nil {
		return nil, txHash, err
	}
	kvArgs[utils.MyIndices] = prvIndices

	//Create unique receiver for token
	uniquePayment := key.PaymentInfo{PaymentAddress: senderWallet.KeySet.PaymentAddress, Amount: totalAmount, Message: []byte{}}

	txTokenParam := tx_ver2.NewTxTokenConvertVer1ToVer2InitParams(&(senderWallet.KeySet.PrivateKey),
		prvInCoins, []*key.PaymentInfo{}, tokenInCoins,
		[]*key.PaymentInfo{&uniquePayment}, prvFee, tokenID,
		nil, nil, kvArgs)

	tx := new(tx_ver2.TxToken)
	err = tx_ver2.InitTokenConversion(tx, txTokenParam)
	if err != nil {
		return nil, txHash, fmt.Errorf("init txtokenconversion error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, txHash, fmt.Errorf("cannot marshal txtokenconversion: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}
