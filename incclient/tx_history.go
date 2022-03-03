package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"sort"
	"strings"
	"time"
)

// TxHistoryInterface implements necessary methods for a history action.
type TxHistoryInterface interface {
	GetLockTime() int64
	GetAmount() uint64
	String() string
	Summarize() string
	GetTxHash() string
}

// TxIn is an in-coming transaction.
// A transaction is considered to be a TxIn if
//	- it receives at least 1 output coin; and
//	- the input coins are does not belong to the receivers.
// In case a user A sends some coins to a user B and A receives a sent-back output coin, this is not considered to
// be a TxIn.
type TxIn struct {
	Version  int8
	LockTime int64
	TxHash   string
	Amount   uint64
	TokenID  string
	Metadata metadata.Metadata
	OutCoins map[string]uint64
	Note     string
}

// GetLockTime returns the lock-time.
func (txIn TxIn) GetLockTime() int64 {
	return txIn.LockTime
}

// GetAmount returns the amount.
func (txIn TxIn) GetAmount() uint64 {
	return txIn.Amount
}

// String returns the string-representation.
func (txIn TxIn) String() string {
	resBytes, err := json.Marshal(txIn)
	if err != nil {
		return ""
	}

	lockTimeStr := time.Unix(txIn.GetLockTime(), 0).Format(common.DateOutputFormat)
	res := fmt.Sprintf("[TxIn] Timestamp: %v, Detail: %v", lockTimeStr, string(resBytes))
	return res
}

// Summarize prints out a summary of a TxIn.
func (txIn TxIn) Summarize() string {
	lockTimeStr := time.Unix(txIn.GetLockTime(), 0).Format(common.DateOutputFormat)
	res := fmt.Sprintf("[TxIn] Timestamp: %v, TxHash: %v, TokenID: %v, Amount: %v", lockTimeStr, txIn.TxHash, txIn.TokenID, txIn.Amount)
	if txIn.Note != "" {
		res = res + fmt.Sprintf(", Note: %v", txIn.Note)
	}
	return res
}

// GetTxHash returns the txHash.
func (txIn TxIn) GetTxHash() string {
	return txIn.TxHash
}

// TxOut is an out-going transaction.
// A transaction is considered to be a TxOut if it spends input coins.
type TxOut struct {
	Version    int8
	LockTime   int64
	TxHash     string
	Amount     uint64
	TokenID    string
	SpentCoins map[string]uint64
	Receivers  []string
	PRVFee     uint64
	TokenFee   uint64
	Metadata   metadata.Metadata
	Note       string
}

// GetLockTime returns the lock-time.
func (txOut TxOut) GetLockTime() int64 {
	return txOut.LockTime
}

// GetAmount returns the amount.
func (txOut TxOut) GetAmount() uint64 {
	return txOut.Amount
}

// GetLockTime returns the lock-time.
func (txOut TxOut) String() string {
	resBytes, err := json.Marshal(txOut)
	if err != nil {
		return ""
	}

	lockTimeStr := time.Unix(txOut.GetLockTime(), 0).Format(common.DateOutputFormat)
	res := fmt.Sprintf("[TxOut] Timestamp: %v, Detail: %v", lockTimeStr, string(resBytes))
	return res
}

// Summarize prints out a summary of a TxOut.
func (txOut TxOut) Summarize() string {
	lockTimeStr := time.Unix(txOut.GetLockTime(), 0).Format(common.DateOutputFormat)
	res := fmt.Sprintf("[TxOut] Timestamp: %v, TxHash: %v, TokenID: %v, Amount: %v", lockTimeStr, txOut.TxHash, txOut.TokenID, txOut.Amount)
	if txOut.Note != "" {
		res = res + fmt.Sprintf(", Note: %v", txOut.Note)
	}
	return res
}

// GetTxHash returns the txHash.
func (txOut TxOut) GetTxHash() string {
	return txOut.TxHash
}

// TxHistory consists of a list of TxIn's and a list of TxOut's.
type TxHistory struct {
	TxInList  []TxIn
	TxOutList []TxOut
}

// GetListTxsInV1 returns a list of all in-coming tokenIDStr transactions (V1) to a private key.
func (client *IncClient) GetListTxsInV1(privateKey string, tokenIDStr string) ([]TxIn, error) {
	kWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot deserialize private key %v: %v", privateKey, err)
	}

	addrStr := PrivateKeyToPaymentAddress(privateKey, -1)
	if addrStr == "" {
		return nil, fmt.Errorf("cannot get payment address")
	}
	listDecryptedCoins, err := client.GetListDecryptedOutCoin(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, err
	}

	mapCmt := makeMapCMToPlainCoin(listDecryptedCoins)

	txList, err := client.GetTransactionsByReceiver(addrStr)
	if err != nil {
		return nil, err
	}

	res := make([]TxIn, 0)
	for txHash, tx := range txList {
		if isOut, err := isTxOut(tx, tokenIDStr, listDecryptedCoins); err != nil {
			return nil, err
		} else if isOut {
			continue
		}

		outCoins, err := getTxOutputCoinsByKeySet(tx, tokenIDStr, &kWallet.KeySet)
		if err != nil {
			return nil, err
		}

		amount := uint64(0)
		for cmtStr := range outCoins {
			if outCoin, ok := mapCmt[cmtStr]; ok {
				amount += outCoin.GetValue()
				continue
			}
		}
		if amount > 0 {
			newTxIn := TxIn{
				Version:  tx.GetVersion(),
				LockTime: tx.GetLockTime(),
				TxHash:   txHash,
				TokenID:  tx.GetTokenID().String(),
				Metadata: tx.GetMetadata(),
			}
			newTxIn.Amount = amount
			res = append(res, newTxIn)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})

	return res, nil
}

// GetListTxsInV2 returns a list of all in-coming tokenIDStr transactions (V2) to a private key.
func (client *IncClient) GetListTxsInV2(privateKey string, tokenIDStr string) ([]TxIn, error) {
	res := make([]TxIn, 0)
	if client.version != 2 {
		return res, nil
	}
	kWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot deserialize private key %v: %v", privateKey, err)
	}

	listDecryptedCoins, err := client.GetListDecryptedOutCoin(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, err
	}

	publicKeys := make([]string, 0)
	for _, outCoin := range listDecryptedCoins {
		if outCoin.GetVersion() == 2 {
			publicKeyStr := base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), common.ZeroByte)
			publicKeys = append(publicKeys, publicKeyStr)
		}
	}

	if len(publicKeys) == 0 {
		return nil, nil
	}

	mapCmt := makeMapCMToPlainCoin(listDecryptedCoins)

	txMap, err := client.GetTransactionsByPublicKeys(publicKeys)
	if err != nil {
		return nil, err
	}

	mapRes := make(map[string]TxIn)
	for _, tmpTxMap := range txMap {
		for txHash, tx := range tmpTxMap {
			if _, ok := mapRes[txHash]; ok {
				continue
			}
			if isOut, err := isTxOut(tx, tokenIDStr, listDecryptedCoins); err != nil {
				return nil, err
			} else if isOut {
				continue
			}

			outCoins, err := getTxOutputCoinsByKeySet(tx, tokenIDStr, &kWallet.KeySet)
			if err != nil {
				return nil, err
			}

			pubKeys := make(map[string]uint64)
			amount := uint64(0)
			for cmtStr := range outCoins {
				if outCoin, ok := mapCmt[cmtStr]; ok {
					amount += outCoin.GetValue()
					pubKeys[base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), 0)] = outCoin.GetValue()
					continue
				}
			}
			if amount > 0 {
				txIn := TxIn{
					Version:  tx.GetVersion(),
					OutCoins: pubKeys,
					LockTime: tx.GetLockTime(),
					TxHash:   txHash,
					TokenID:  tx.GetTokenID().String(),
					Metadata: tx.GetMetadata(),
					Amount:   amount,
					Note:     txMetadataNote[tx.GetMetadataType()],
				}
				mapRes[txHash] = txIn
			}
		}
	}

	for _, txIn := range mapRes {
		res = append(res, txIn)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})

	return res, nil
}

// GetListTxsIn returns a list of all in-coming tokenIDStr transactions to a private key.
func (client *IncClient) GetListTxsIn(privateKey string, tokenIDStr string) ([]TxIn, error) {
	res := make([]TxIn, 0)

	txInsV1, err := client.GetListTxsInV1(privateKey, tokenIDStr)
	if err != nil {
		return nil, err
	}
	res = append(res, txInsV1...)

	txInsV2, err := client.GetListTxsInV2(privateKey, tokenIDStr)
	if err != nil {
		return nil, err
	}
	res = append(res, txInsV2...)

	// sort the results based on lock-time
	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})

	return res, nil
}

// GetListTxsOutV1 returns a list of all out-going tokenIDStr transactions (V1) of a private key.
func (client *IncClient) GetListTxsOutV1(privateKey string, tokenIDStr string) ([]TxOut, error) {
	kWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot deserialize private key %v: %v", privateKey, err)
	}

	addr := kWallet.KeySet.PaymentAddress
	shardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])

	listSpentCoins, _, err := client.GetSpentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, err
	}

	// Create a map from serial numbers to coins
	mapSpentCoins := make(map[string]coin.PlainCoin)
	snList := make([]string, 0)
	for _, spentCoin := range listSpentCoins {
		if spentCoin.GetVersion() != 2 {
			snStr := base58.Base58Check{}.Encode(spentCoin.GetKeyImage().ToBytesS(), common.ZeroByte)
			mapSpentCoins[snStr] = spentCoin
			snList = append(snList, snStr)
		}
	}

	if len(snList) == 0 {
		return nil, nil
	}

	// Retrieve the list of transactions which spent these coins
	mapSpentTxs, err := client.GetTxHashBySerialNumbers(snList, tokenIDStr, shardID)
	if err != nil {
		if strings.Contains(err.Error(), "Method not found") {
			return nil, fmt.Errorf("method not supported by the remote node configurations")
		}
		return nil, err
	}

	mapRes := make(map[string]TxOut)
	for _, txHash := range mapSpentTxs {
		// check if the txHash has been processed
		if _, ok := mapRes[txHash]; ok {
			continue
		}

		tx, err := client.GetTx(txHash)
		if err != nil {
			return nil, err
		}

		//get transaction fee
		fee, isPRVFee := getTxFeeBy(tx)

		//calculate transaction's amount
		inputAmount, spentCoins, err := getTxInputAmount(tx, tokenIDStr, mapSpentCoins)
		if err != nil {
			return nil, err
		}
		outputAmount, err := getTxOutputAmountByKeySet(tx, tokenIDStr, &kWallet.KeySet)
		if err != nil {
			return nil, err
		}
		amount := inputAmount - outputAmount
		if isPRVFee && tokenIDStr == common.PRVIDStr {
			amount -= fee
		}
		if !isPRVFee && tokenIDStr != common.PRVIDStr {
			amount -= fee
		}

		//get list of receivers' public keys
		receivers, err := getTxReceivers(tx, tokenIDStr)
		if err != nil {
			return nil, err
		}

		if amount > 0 || tokenIDStr == common.PRVIDStr {
			note := txMetadataNote[tx.GetMetadataType()]
			if tokenIDStr == common.PRVIDStr && amount == 0 {
				note += " (Tx Fee)"
			}
			newTxOut := TxOut{
				Version:    tx.GetVersion(),
				LockTime:   tx.GetLockTime(),
				TxHash:     txHash,
				TokenID:    tx.GetTokenID().String(),
				SpentCoins: spentCoins,
				Receivers:  receivers,
				Amount:     amount,
				Metadata:   tx.GetMetadata(),
				PRVFee:     fee,
				Note:       note,
			}
			if !isPRVFee {
				newTxOut.PRVFee = 0
				newTxOut.TokenFee = fee
			}

			mapRes[txHash] = newTxOut
		}
	}

	res := make([]TxOut, 0)
	for _, txOut := range mapRes {
		res = append(res, txOut)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})

	return res, nil
}

// GetListTxsOutV2 returns a list of all out-going tokenIDStr transactions (V2) of a private key.
func (client *IncClient) GetListTxsOutV2(privateKey string, tokenIDStr string) ([]TxOut, error) {
	kWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot deserialize private key %v: %v", privateKey, err)
	}

	addr := kWallet.KeySet.PaymentAddress
	shardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])

	listSpentCoins, _, err := client.GetSpentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, err
	}

	// Create a map from serial numbers to coins
	mapSpentCoins := make(map[string]coin.PlainCoin)
	snList := make([]string, 0)
	for _, spentCoin := range listSpentCoins {
		if spentCoin.GetVersion() == 2 {
			snStr := base58.Base58Check{}.Encode(spentCoin.GetKeyImage().ToBytesS(), common.ZeroByte)
			mapSpentCoins[snStr] = spentCoin
			snList = append(snList, snStr)
		}
	}

	if len(snList) == 0 {
		return nil, nil
	}

	// Retrieve the list of transactions which spent these coins
	mapSpentTxs, err := client.GetTxHashBySerialNumbers(snList, tokenIDStr, shardID)
	if err != nil {
		if strings.Contains(err.Error(), "Method not found") {
			return nil, fmt.Errorf("method not supported by the remote node configurations")
		}
		return nil, err
	}

	mapRes := make(map[string]TxOut)
	for _, txHash := range mapSpentTxs {
		// check if the txHash has been processed
		if _, ok := mapRes[txHash]; ok {
			continue
		}

		tx, err := client.GetTx(txHash)
		if err != nil {
			return nil, err
		}

		//get transaction fee
		fee, isPRVFee := getTxFeeBy(tx)

		//calculate transaction's amount
		inputAmount, spentCoins, err := getTxInputAmount(tx, tokenIDStr, mapSpentCoins)
		if err != nil {
			return nil, err
		}
		outputAmount, err := getTxOutputAmountByKeySet(tx, tokenIDStr, &kWallet.KeySet)
		if err != nil {
			return nil, err
		}
		amount := inputAmount - outputAmount
		if isPRVFee && tokenIDStr == common.PRVIDStr {
			amount -= fee
		}
		if !isPRVFee && tokenIDStr != common.PRVIDStr {
			amount -= fee
		}

		//get list of receivers' public keys
		receivers, err := getTxReceivers(tx, tokenIDStr)
		if err != nil {
			return nil, err
		}

		if amount > 0 || tokenIDStr == common.PRVIDStr {
			note := txMetadataNote[tx.GetMetadataType()]
			if tokenIDStr == common.PRVIDStr && amount == 0 {
				note += " (Tx Fee)"
			}
			newTxOut := TxOut{
				Version:    tx.GetVersion(),
				LockTime:   tx.GetLockTime(),
				TxHash:     txHash,
				TokenID:    tx.GetTokenID().String(),
				SpentCoins: spentCoins,
				Receivers:  receivers,
				Amount:     amount,
				Metadata:   tx.GetMetadata(),
				PRVFee:     fee,
				Note:       note,
			}
			if !isPRVFee {
				newTxOut.PRVFee = 0
				newTxOut.TokenFee = fee
			}

			mapRes[txHash] = newTxOut
		}
	}

	res := make([]TxOut, 0)
	for _, txOut := range mapRes {
		res = append(res, txOut)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})

	return res, nil
}

// GetListTxsOut returns a list of all out-coming tokenIDStr transactions of a private key.
func (client *IncClient) GetListTxsOut(privateKey string, tokenIDStr string) ([]TxOut, error) {
	res := make([]TxOut, 0)

	txOutsV1, err := client.GetListTxsOutV1(privateKey, tokenIDStr)
	if err != nil {
		return nil, err
	}
	res = append(res, txOutsV1...)

	txOutsV2, err := client.GetListTxsOutV2(privateKey, tokenIDStr)
	if err != nil {
		return nil, err
	}
	res = append(res, txOutsV2...)

	// sort the results based on lock-time
	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})

	return res, nil
}

// GetTxHistoryV1 retrieves the history (V1) of a private key w.r.t a tokenID.
func (client *IncClient) GetTxHistoryV1(privateKey string, tokenIDStr string) (*TxHistory, error) {
	var err error
	res := TxHistory{}
	res.TxInList, err = client.GetListTxsInV1(privateKey, tokenIDStr)
	if err != nil {
		return nil, err
	}

	res.TxOutList, err = client.GetListTxsOutV1(privateKey, tokenIDStr)
	if err != nil && !strings.Contains(err.Error(), "method not supported") {
		return nil, err
	}

	return &res, nil
}

// GetTxHistoryV2 retrieves the history (V2) of a private key w.r.t a tokenID.
func (client *IncClient) GetTxHistoryV2(privateKey string, tokenIDStr string) (*TxHistory, error) {
	var err error
	res := TxHistory{}
	res.TxInList, err = client.GetListTxsInV2(privateKey, tokenIDStr)
	if err != nil {
		return nil, err
	}

	res.TxOutList, err = client.GetListTxsOutV2(privateKey, tokenIDStr)
	if err != nil && !strings.Contains(err.Error(), "method not supported") {
		return nil, err
	}

	return &res, nil
}
