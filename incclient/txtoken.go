package incclient

import (
	"encoding/json"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver1"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// CreateRawTokenTransaction creates a token transaction with the provided version.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTokenTransaction(txParam *TxParam, version int8) ([]byte, string, error) {
	if version == -1 { //Try either one of the version, if possible
		encodedTx, txHash, err := client.CreateRawTokenTransactionVer1(txParam)
		if err != nil {
			encodedTx, txHash, err1 := client.CreateRawTokenTransactionVer2(txParam)
			if err1 != nil {
				return nil, "", fmt.Errorf("cannot create raw token transaction for either version: %v, %v", err, err1)
			}

			return encodedTx, txHash, nil
		}

		return encodedTx, txHash, nil
	} else if version == 2 {
		return client.CreateRawTokenTransactionVer2(txParam)
	} else {
		return client.CreateRawTokenTransactionVer1(txParam)
	}
}

// CreateRawTokenTransactionVer1 creates a token transaction version 1.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTokenTransactionVer1(txParam *TxParam) ([]byte, string, error) {
	privateKey := txParam.senderPrivateKey

	tokenIDStr := txParam.txTokenParam.tokenID
	_, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	lastByteSender := senderWallet.KeySet.PaymentAddress.Pk[len(senderWallet.KeySet.PaymentAddress.Pk)-1]
	shardID := common.GetShardIDFromLastByte(lastByteSender)

	hasTokenFee := txParam.txTokenParam.hasTokenFee

	//Calculate the total transacted amount
	totalAmount := uint64(0)
	for _, amount := range txParam.txTokenParam.amountList {
		totalAmount += amount
	}

	//Create list of payment infos
	var tokenReceivers []*key.PaymentInfo
	if txParam.txTokenParam.tokenType == utils.CustomTokenInit {
		uniqueReceiver := key.PaymentInfo{PaymentAddress: senderWallet.KeySet.PaymentAddress, Amount: totalAmount, Message: []byte{}}
		tokenReceivers = []*key.PaymentInfo{&uniqueReceiver}
	} else {
		tokenReceivers, err = createPaymentInfos(txParam.txTokenParam.receiverList, txParam.txTokenParam.amountList)
		if err != nil {
			return nil, "", err
		}
	}

	hasPrivacyToken := true
	hasPrivacyPRV := true
	if txParam.md != nil {
		hasPrivacyToken = false
		hasPrivacyPRV = false
	} else if txParam.txTokenParam.tokenType == utils.CustomTokenInit {
		hasPrivacyPRV = true
		hasPrivacyToken = false
	}

	//Init PRV fee param when not paying fee by token
	var coinsPRVToSpend []coin.PlainCoin
	var kvArgsPRV map[string]interface{}
	prvFee := txParam.fee
	if prvFee == 0 {
		prvFee = DefaultPRVFee
	}

	totalPRVAmount := prvFee
	for _, amount := range txParam.amountList {
		totalPRVAmount += amount
	}

	tokenFee := uint64(0)
	if !hasTokenFee {
		//if totalPRVAmount >= prvFee {
		//	totalPRVAmount -= prvFee
		//}
		coinsPRVToSpend, kvArgsPRV, err = client.initParams(privateKey, common.PRVIDStr, totalPRVAmount, hasPrivacyPRV, 1)
		if err != nil {
			return nil, "", err
		}
	} else {
		//set prv fee to 0
		prvFee = 0

		//calculate the token amount to pay transaction fee
		tokenFee, err = client.GetTokenFee(shardID, tokenIDStr)
		if err != nil {
			return nil, "", err
		}
		totalAmount += tokenFee
	}
	//End init PRV fee param

	//Init token param
	var coinsTokenToSpend []coin.PlainCoin
	var kvArgsToken map[string]interface{}
	if txParam.txTokenParam.tokenType != utils.CustomTokenInit {
		coinsTokenToSpend, kvArgsToken, err = client.initParams(privateKey, tokenIDStr, totalAmount, true, 1)
		if err != nil {
			return nil, "", err
		}
	}
	//End init token param

	//Create token param for transactions
	tokenParam := tx_generic.NewTokenParam(tokenIDStr, "", "",
		totalAmount, txParam.txTokenParam.tokenType, tokenReceivers, coinsTokenToSpend, false, tokenFee, kvArgsToken)

	prvReceivers := make([]*key.PaymentInfo, 0)
	if len(txParam.receiverList) > 0 {
		prvReceivers, err = createPaymentInfos(txParam.receiverList, txParam.amountList)
		if err != nil {
			return nil, "", err
		}
	}

	txTokenParam := tx_generic.NewTxTokenParams(&senderWallet.KeySet.PrivateKey, prvReceivers, coinsPRVToSpend, prvFee,
		tokenParam, txParam.md, hasPrivacyPRV, hasPrivacyToken, shardID, nil, kvArgsPRV)

	tx := new(tx_ver1.TxToken)
	err = tx.Init(txTokenParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txtokenver1 error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txtokenver1: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateRawTokenTransactionVer2 creates a token transaction version 2.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateRawTokenTransactionVer2(txParam *TxParam) ([]byte, string, error) {
	privateKey := txParam.senderPrivateKey

	tokenIDStr := txParam.txTokenParam.tokenID
	_, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	lastByteSender := senderWallet.KeySet.PaymentAddress.Pk[len(senderWallet.KeySet.PaymentAddress.Pk)-1]
	shardID := common.GetShardIDFromLastByte(lastByteSender)

	//Calculate the total transacted amount
	totalAmount := uint64(0)
	for _, amount := range txParam.txTokenParam.amountList {
		totalAmount += amount
	}

	//Create list of payment infos
	tokenReceivers, err := createPaymentInfos(txParam.txTokenParam.receiverList, txParam.txTokenParam.amountList)
	if err != nil {
		return nil, "", err
	}

	prvFee := txParam.fee
	if prvFee == 0 {
		prvFee = DefaultPRVFee
	}

	totalPRVAmount := prvFee
	for _, amount := range txParam.amountList {
		totalPRVAmount += amount
	}

	//Init PRV fee param
	coinsToSpendPRV, kvArgsPRV, err := client.initParams(privateKey, common.PRVIDStr, totalPRVAmount, true, 2)
	if err != nil {
		return nil, "", err
	}
	//End init PRV fee param

	//Init token param
	coinsTokenToSpend, kvArgsToken, err := client.initParams(privateKey, tokenIDStr, totalAmount, true, 2)
	if err != nil {
		return nil, "", err
	}
	//End init token param

	//Create token param for transactions
	tokenParam := tx_generic.NewTokenParam(tokenIDStr, "", "",
		totalAmount, txParam.txTokenParam.tokenType, tokenReceivers, coinsTokenToSpend, false, 0, kvArgsToken)

	prvReceivers := make([]*key.PaymentInfo, 0)
	if len(txParam.receiverList) > 0 {
		prvReceivers, err = createPaymentInfos(txParam.receiverList, txParam.amountList)
		if err != nil {
			return nil, "", err
		}
	}
	txTokenParam := tx_generic.NewTxTokenParams(&senderWallet.KeySet.PrivateKey, prvReceivers, coinsToSpendPRV, prvFee,
		tokenParam, txParam.md, true, true, shardID, nil, kvArgsPRV)

	tx := new(tx_ver2.TxToken)
	err = tx.Init(txTokenParam)
	if err != nil {
		return nil, "", fmt.Errorf("init txtokenver2 error: %v", err)
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot marshal txtokenver2: %v", err)
	}

	base58CheckData := base58.Base58Check{}.Encode(txBytes, common.ZeroByte)

	return []byte(base58CheckData), tx.Hash().String(), nil
}

// CreateAndSendRawTokenTransaction creates a token transaction with the provided version, and submits it to the Incognito network.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendRawTokenTransaction(privateKey string, addrList []string, amountList []uint64, tokenID string, version int8, md metadata.Metadata) (string, error) {
	tokenParams := NewTxTokenParam(tokenID, 1, addrList, amountList, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, tokenParams, md, nil)
	encodedTx, txHash, err := client.CreateRawTokenTransaction(txParam, version)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateTokenInitTransaction creates a token init transaction with the provided version.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateTokenInitTransaction(privateKey, tokenName, tokenSymbol string, amount uint64, version int) ([]byte, string, error) {
	if version == -1 { //Try either one of the version, if possible
		encodedTx, txHash, err := client.CreateTokenInitTransactionV1(privateKey, tokenName, tokenSymbol, amount)
		if err != nil {
			encodedTx, txHash, err1 := client.CreateTokenInitTransactionV2(privateKey, tokenName, tokenSymbol, amount)
			if err1 != nil {
				return nil, "", fmt.Errorf("cannot create raw token init transaction for either version: %v, %v", err, err1)
			}

			return encodedTx, txHash, nil
		}

		return encodedTx, txHash, nil
	} else if version == 2 {
		return client.CreateTokenInitTransactionV2(privateKey, tokenName, tokenSymbol, amount)
	} else {
		return client.CreateTokenInitTransactionV1(privateKey, tokenName, tokenSymbol, amount)
	}
}

// CreateTokenInitTransactionV1 inits a new token version 1. In version 1, users are free to choose the tokenID of their own
// as long as it does not collide with any existing token in the Incognito network.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateTokenInitTransactionV1(privateKey, _, _ string, amount uint64) ([]byte, string, error) {
	addr := PrivateKeyToPaymentAddress(privateKey, 0)
	tokenParam := NewTxTokenParam("", 0, []string{addr}, []uint64{amount}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, nil, nil)

	return client.CreateRawTokenTransactionVer1(txParam)

}

// CreateTokenInitTransactionV2 inits a new token version 2. In version 2, to ensure that no collision happens between
// the initialized tokenID and the existing ones. The new tokenID will be generated by the shard committee based on
// information given in the requesting transaction.
//
// Note: that this transaction is a PRV transaction.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateTokenInitTransactionV2(privateKey, tokenName, tokenSymbol string, amount uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	addr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	pubKeyStr, txRandomStr, err := GenerateOTAFromPaymentAddress(addr)
	if err != nil {
		return nil, "", err
	}

	if tokenName == "" {
		tokenName = "INC_" + common.RandChars(5)
	}
	if tokenSymbol == "" {
		tokenSymbol = "INC_" + common.RandChars(5)
	}

	tokenInitReq, err := metadata.NewInitTokenRequest(pubKeyStr, txRandomStr, amount, tokenName, tokenSymbol, metadata.InitTokenRequestMeta)
	if err != nil {
		return nil, "", err
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, tokenInitReq, nil)

	return client.CreateRawTransaction(txParam, 2)
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

	//We need to use PRV coinV2 to payment (it's a must)
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

// SendRawTokenTx sends submits a raw token transaction to the Incognito blockchain.
func (client *IncClient) SendRawTokenTx(encodedTx []byte) error {
	responseInBytes, err := client.rpcServer.SendRawTokenTx(string(encodedTx))
	if err != nil {
		return nil
	}

	_, err = rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return err
	}

	return nil
}
