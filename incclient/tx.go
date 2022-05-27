package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver1"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// CreateRawTransaction creates a PRV transaction with the provided version.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTransaction(param *TxParam, version int8) ([]byte, string, error) {
	if param.txTokenParam != nil {
		return nil, "", fmt.Errorf("method supports PRV transaction only")
	}
	if version == -1 { //Try either one of the version, if possible
		encodedTx, txHash, err := client.CreateRawTransactionVer1(param)
		if err != nil {
			encodedTx, txHash, err1 := client.CreateRawTransactionVer2(param)
			if err1 != nil {
				return nil, "", fmt.Errorf("cannot create raw transaction for either version: %v, %v", err, err1)
			}

			return encodedTx, txHash, nil
		}

		return encodedTx, txHash, nil
	} else if version == 2 {
		return client.CreateRawTransactionVer2(param)
	} else if version == 1 {
		return client.CreateRawTransactionVer1(param)
	}

	return nil, "", fmt.Errorf("transaction version is invalid")
}

// CreateRawTransactionVer1 creates a PRV transaction version 1.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTransactionVer1(param *TxParam) ([]byte, string, error) {
	privateKey := param.senderPrivateKey
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	//Create list of payment infos
	paymentInfos, err := createPaymentInfos(param.receiverList, param.amountList)
	if err != nil {
		return nil, "", err
	}

	//Calculate the total transacted amount
	if param.fee == 0 {
		param.fee = client.cfg.DefaultPRVFee
	}
	totalAmount := param.fee
	for _, amount := range param.amountList {
		totalAmount += amount
	}

	hasPrivacy := true
	if param.md != nil {
		hasPrivacy = false
	}

	coinsToSpend, kvArgs, err := client.initParamsV1(param, common.PRVIDStr, totalAmount, hasPrivacy)
	if err != nil {
		return nil, "", err
	}

	txInitParam := tx_generic.NewTxPrivacyInitParams(&(senderWallet.KeySet.PrivateKey), paymentInfos, coinsToSpend, param.fee, hasPrivacy, &common.PRVCoinID, param.md, nil, kvArgs)

	tx := new(tx_ver1.Tx)
	err = tx.Init(txInitParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txver1 error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txver1: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateRawTransactionVer2 creates a PRV transaction version 2.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTransactionVer2(param *TxParam) ([]byte, string, error) {
	privateKey := param.senderPrivateKey
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	//Create list of payment infos
	paymentInfos, err := createPaymentInfos(param.receiverList, param.amountList)
	if err != nil {
		return nil, "", err
	}

	txFee := param.fee
	if param.fee == 0 {
		txFee = client.cfg.DefaultPRVFee
	}

	//Calculate the total transacted amount
	totalAmount := txFee
	for _, amount := range param.amountList {
		totalAmount += amount
	}

	hasPrivacy := true
	if param.md != nil {
		hasPrivacy = false
	}

	coinsToSpend, kArgs, err := client.initParamsV2(param, common.PRVIDStr, totalAmount)
	if err != nil {
		return nil, "", err
	}

	txParam := tx_generic.NewTxPrivacyInitParams(&(senderWallet.KeySet.PrivateKey), paymentInfos, coinsToSpend, txFee, hasPrivacy, &common.PRVCoinID, param.md, nil, kArgs)

	tx := new(tx_ver2.Tx)
	err = tx.Init(txParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txver2 error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txver2: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateAndSendRawTransaction creates a PRV transaction with the provided version, and submits it to the Incognito network.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendRawTransaction(privateKey string, addrList []string, amountList []uint64, version int8, md metadata.Metadata) (string, error) {
	txParam := NewTxParam(privateKey, addrList, amountList, 0, nil, md, nil)
	encodedTx, txHash, err := client.CreateRawTransaction(txParam, version)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateRawTransactionWithInputCoins creates a raw PRV transaction from the provided input coins.
// Parameters:
//	- param: a regular TxParam.
//	- inputCoins: a list of decrypted, unspent PRV output coins (with the same version).
//	- coinIndices: a list of corresponding indices for the input coins. This value must not be `nil` if the caller is
//	creating a transaction v2.
//
// For transaction with metadata, callers must make sure other values of `param` are valid.
//
// NOTE: this servers PRV transactions only.
func (client *IncClient) CreateRawTransactionWithInputCoins(param *TxParam, inputCoins []coin.PlainCoin, coinIndices []uint64) ([]byte, string, error) {
	var txHash string
	if param.txTokenParam != nil {
		return nil, txHash, fmt.Errorf("this function supports PRV transaction only")
	}

	// check version of coins
	version, err := getVersionFromInputCoins(inputCoins)
	if err != nil {
		return nil, txHash, err
	}
	if version == 2 && coinIndices == nil {
		return nil, txHash, fmt.Errorf("coinIndices must not be nil")
	}

	// check number of input coins
	if len(inputCoins) > MaxInputSize {
		return nil, txHash, fmt.Errorf("support at most %v input coins, got %v", MaxInputSize, len(inputCoins))
	}

	cp := coinParams{
		coinList: inputCoins,
		idxList:  coinIndices,
	}
	param.kArgs = make(map[string]interface{})
	param.kArgs[prvInCoinKey] = cp

	return client.CreateRawTransaction(param, int8(version))
}

// SendRawTx sends submits a raw PRV transaction to the Incognito blockchain.
func (client *IncClient) SendRawTx(encodedTx []byte) error {
	responseInBytes, err := client.rpcServer.SendRawTx(string(encodedTx))
	if err != nil {
		return nil
	}

	err = rpchandler.ParseResponse(responseInBytes, nil)
	if err != nil {
		return err
	}

	return nil
}
