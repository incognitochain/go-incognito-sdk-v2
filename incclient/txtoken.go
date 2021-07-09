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
	if txParam.txTokenParam == nil {
		return nil, "", fmt.Errorf("TxTokenParam must not be nil")
	}
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
	if txParam.txTokenParam == nil {
		return nil, "", fmt.Errorf("TxTokenParam must not be nil")
	}

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
		coinsPRVToSpend, kvArgsPRV, err = client.initParamsV1(txParam, common.PRVIDStr, totalPRVAmount, hasPrivacyPRV)
		if err != nil {
			return nil, "", err
		}
	} else {
		//set prv fee to 0
		prvFee = 0

		//calculate the token amount to pay transaction fee
		tokenFee = txParam.txTokenParam.tokenFee
		if tokenFee == 0 {
			tokenFee, err = client.GetTokenFee(shardID, tokenIDStr)
			if err != nil {
				return nil, "", err
			}
		}
		totalAmount += tokenFee
	}
	//End init PRV fee param

	//Init token param
	var coinsTokenToSpend []coin.PlainCoin
	var kvArgsToken map[string]interface{}
	if txParam.txTokenParam.tokenType != utils.CustomTokenInit {
		coinsTokenToSpend, kvArgsToken, err = client.initParamsV1(txParam, tokenIDStr, totalAmount, true)
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
	if txParam.txTokenParam == nil {
		return nil, "", fmt.Errorf("TxTokenParam must not be nil")
	}

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
	coinsToSpendPRV, kvArgsPRV, err := client.initParamsV2(txParam, common.PRVIDStr, totalPRVAmount)
	if err != nil {
		Logger.Printf("init PRVParamsV2 error: %v\n", err)
		return nil, "", err
	}
	//End init PRV fee param

	//Init token param
	coinsTokenToSpend, kvArgsToken, err := client.initParamsV2(txParam, tokenIDStr, totalAmount)
	if err != nil {
		Logger.Printf("init TokenParamsV2 error: %v\n", err)
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

// CreateRawTokenTransactionWithInputCoins creates a raw token transaction from the provided input coins.
// Parameters:
//	- txParam: a regular TxParam.
//	- tokenInCoins: a list of decrypted, unspent token output coins (with the same version).
//	- tokenIndices: a list of corresponding indices for the token input coins. This value must not be `nil` if the caller is
//	creating a transaction v2.
//	- prvInCoins: a list of decrypted, unspent PRV output coins for paying the transaction fee (if have).
//	- prvIndices: a list of corresponding indices for the prv input coins. This value must not be `nil` if the caller is
//	creating a transaction v2.
//
// For transaction with metadata, callers must make sure other values of `param` are valid.
//
// NOTE: this servers PRV transactions only.
func (client *IncClient) CreateRawTokenTransactionWithInputCoins(txParam *TxParam,
	tokenInCoins []coin.PlainCoin,
	tokenIndices []uint64,
	prvInCoins []coin.PlainCoin,
	prvIndices []uint64,
) ([]byte, string, error) {
	var txHash string
	if txParam.txTokenParam == nil {
		return nil, txHash, fmt.Errorf("this function supports token transaction only")
	}

	// check version of coins
	version, err := getVersionFromInputCoins(tokenInCoins)
	if err != nil {
		return nil, txHash, err
	}
	if version == 2 && tokenIndices == nil {
		return nil, txHash, fmt.Errorf("tokenIndices must not be nil")
	}
	if version == 2 && prvInCoins == nil {
		return nil, txHash, fmt.Errorf("must have PRV input coins to pay the transaction fee")
	}

	isPRVFee := prvInCoins != nil
	if isPRVFee {
		// check version of the PRV coins
		prvVersion, err := getVersionFromInputCoins(prvInCoins)
		if err != nil {
			return nil, txHash, err
		}
		if prvVersion != version {
			return nil, txHash, fmt.Errorf("expect PRV version to be %v, got %v", version, prvVersion)
		}

		if prvVersion == 2 && prvIndices == nil {
			return nil, txHash, fmt.Errorf("prvIndices must not be nil")
		}
	}

	// check number of input coins
	if len(tokenInCoins) > MaxInputSize {
		return nil, txHash, fmt.Errorf("support at most %v token input coins, got %v", MaxInputSize, len(tokenInCoins))
	}
	if len(prvInCoins) > MaxInputSize {
		return nil, txHash, fmt.Errorf("support at most %v PRV input coins, got %v", MaxInputSize, len(prvInCoins))
	}

	prvCp := coinParams{
		coinList: prvInCoins,
		idxList:  prvIndices,
	}
	tokenCp := coinParams{
		coinList: tokenInCoins,
		idxList:  tokenIndices,
	}
	if version == 1 && prvInCoins == nil {
		txParam.txTokenParam.hasTokenFee = true
	}

	txParam.kArgs = make(map[string]interface{})
	txParam.kArgs[prvInCoinKey] = prvCp
	txParam.kArgs[tokenInCoinKey] = tokenCp

	return client.CreateRawTokenTransaction(txParam, int8(version))
}

// SendRawTokenTx sends submits a raw token transaction to the Incognito blockchain.
func (client *IncClient) SendRawTokenTx(encodedTx []byte) error {
	responseInBytes, err := client.rpcServer.SendRawTokenTx(string(encodedTx))
	if err != nil {
		return nil
	}

	err = rpchandler.ParseResponse(responseInBytes, nil)
	if err != nil {
		return err
	}

	return nil
}
