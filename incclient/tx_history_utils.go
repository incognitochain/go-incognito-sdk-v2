package incclient

import (
	"encoding/csv"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"os"
)

const (
	DefaultTxIn      = "txIn_history.csv"
	DefaultTxOut     = "txOut_history.csv"
	DefaultTxHistory = "txHistory.csv"
)

// getListKeyImagesFromTx returns the list of key images of a transaction based on the provided tokenID.
func getListKeyImagesFromTx(tx metadata.Transaction, tokenIDStr string) (map[string]string, error) {
	var proof privacy.Proof
	res := make(map[string]string)

	switch tx.GetType() {
	case common.TxNormalType, common.TxRewardType, common.TxReturnStakingType, common.TxConversionType:
		if tokenIDStr != common.PRVIDStr {
			return nil, nil
		}
		proof = tx.GetProof()

	case common.TxCustomTokenPrivacyType, common.TxTokenConversionType:
		tmpTx, ok := tx.(tx_generic.TransactionToken)
		if !ok {
			return nil, fmt.Errorf("cannot parse the transaction as a transaction token")
		}

		if tokenIDStr == common.PRVIDStr {
			proof = tmpTx.GetTxBase().GetProof()
		} else if tokenIDStr == tx.GetTokenID().String() {
			proof = tmpTx.GetTxNormal().GetProof()
		} else {
			return nil, nil
		}
	}

	if proof != nil {
		for _, inCoin := range proof.GetInputCoins() {
			keyImageStr := base58.Base58Check{}.Encode(inCoin.GetKeyImage().ToBytesS(), common.ZeroByte)
			res[keyImageStr] = keyImageStr
		}
	}

	return res, nil
}

// getTxOutputCoinsByKeySet returns the list of output coins of a transaction sent to a key-set based on the provided tokenID.
// It returns a map from the commitment (base58-encoded) of an output coin to the output coin itself.
func getTxOutputCoinsByKeySet(tx metadata.Transaction, tokenIDStr string, keySet *key.KeySet) (map[string]coin.Coin, error) {
	var proof privacy.Proof
	res := make(map[string]coin.Coin)

	switch tx.GetType() {
	case common.TxNormalType, common.TxRewardType, common.TxReturnStakingType, common.TxConversionType:
		if tokenIDStr != common.PRVIDStr {
			return nil, nil
		}
		proof = tx.GetProof()

	case common.TxCustomTokenPrivacyType, common.TxTokenConversionType:
		tmpTx, ok := tx.(tx_generic.TransactionToken)
		if !ok {
			return nil, fmt.Errorf("cannot parse the transaction as a transaction token")
		}

		if tokenIDStr == common.PRVIDStr {
			proof = tmpTx.GetTxBase().GetProof()
		} else if tokenIDStr == tmpTx.GetTokenID().String() || tmpTx.GetVersion() == 2 {
			proof = tmpTx.GetTxNormal().GetProof()
		} else {
			return nil, nil
		}
	}

	if proof != nil {
		for _, outCoin := range proof.GetOutputCoins() {
			isOwned, _ := outCoin.DoesCoinBelongToKeySet(keySet)
			if isOwned {
				cmtStr := base58.Base58Check{}.Encode(outCoin.GetCommitment().ToBytesS(), common.ZeroByte)
				res[cmtStr] = outCoin
			}
		}
	}

	return res, nil
}

// getTxReceivers returns a list of base58-encoded public keys of a transaction w.r.t a tokenID.
func getTxReceivers(tx metadata.Transaction, tokenIDStr string) ([]string, error) {
	var proof privacy.Proof
	res := make([]string, 0)

	switch tx.GetType() {
	case common.TxNormalType, common.TxRewardType, common.TxReturnStakingType, common.TxConversionType:
		if tokenIDStr != common.PRVIDStr {
			return nil, nil
		}
		proof = tx.GetProof()

	case common.TxCustomTokenPrivacyType, common.TxTokenConversionType:
		tmpTx, ok := tx.(tx_generic.TransactionToken)
		if !ok {
			return nil, fmt.Errorf("cannot parse the transaction as a transaction token")
		}

		if tokenIDStr == common.PRVIDStr {
			proof = tmpTx.GetTxBase().GetProof()
		} else if tokenIDStr == tmpTx.GetTokenID().String() {
			proof = tmpTx.GetTxNormal().GetProof()
		} else {
			return nil, nil
		}
	}

	if proof != nil {
		for _, outCoin := range proof.GetOutputCoins() {
			publicKeyStr := base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), common.ZeroByte)
			res = append(res, publicKeyStr)
		}
	}

	return res, nil
}

// getTxFeeBy returns the transaction fee, a boolean indicating if the fee is paid by PRV.
func getTxFeeBy(tx metadata.Transaction) (uint64, bool) {
	prvFee := tx.GetTxFee()
	tokenFee := tx.GetTxFeeToken()

	if prvFee > 0 {
		return prvFee, true
	} else {
		return tokenFee, false
	}
}

// getTxInputAmount returns the total input amount and the amount of each input coins of a transaction.
func getTxInputAmount(tx metadata.Transaction, tokenIDStr string, listTXOs map[string]coin.PlainCoin) (uint64, map[string]uint64, error) {
	listKeyImages, err := getListKeyImagesFromTx(tx, tokenIDStr)
	if err != nil {
		return 0, nil, err
	}

	amount := uint64(0)
	txInputs := make(map[string]uint64)
	for keyImageStr := range listKeyImages {
		if outCoin, ok := listTXOs[keyImageStr]; ok {
			amount += outCoin.GetValue()
			txInputs[keyImageStr] = outCoin.GetValue()
		}
	}

	return amount, txInputs, nil
}

// getTxOutputAmountByKeySet returns the total output amount of a transaction sent to a key-set based on the provided tokenID.
func getTxOutputAmountByKeySet(tx metadata.Transaction, tokenIDStr string, keySet *key.KeySet) (uint64, error) {
	outCoins, err := getTxOutputCoinsByKeySet(tx, tokenIDStr, keySet)
	if err != nil {
		return 0, err
	}

	res := uint64(0)
	for _, outCoin := range outCoins {
		decryptedCoin, err := outCoin.Decrypt(keySet)
		if err != nil {
			return 0, fmt.Errorf("cannot decrypt outputCoin")
		}
		res += decryptedCoin.GetValue()
	}

	return res, nil
}

// makeMapCMToPlainCoin returns a map from commitments to a plain-coins
func makeMapCMToPlainCoin(outCoins map[string]coin.PlainCoin) map[string]coin.PlainCoin {
	res := make(map[string]coin.PlainCoin)
	for _, outCoin := range outCoins {
		snStr := base58.Base58Check{}.Encode(outCoin.GetCommitment().ToBytesS(), common.ZeroByte)
		res[snStr] = outCoin
	}

	return res
}

// makeMapPublicKeyToPlainCoinV2 returns a map from public keys to a plain-coins (for CoinV2 only).
func makeMapPublicKeyToPlainCoinV2(outCoins map[string]coin.PlainCoin) map[string]coin.PlainCoin {
	res := make(map[string]coin.PlainCoin)
	for _, outCoin := range outCoins {
		if outCoin.GetVersion() != 2 {
			continue
		}
		publicKeyStr := base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), common.ZeroByte)
		res[publicKeyStr] = outCoin
	}

	return res
}

// isTxOut checks if a transaction is an out-going transaction w.r.t a list out TXOs or not.
func isTxOut(tx metadata.Transaction, tokenIDStr string, listTXOs map[string]coin.PlainCoin) (bool, error) {
	listKeyImages, err := getListKeyImagesFromTx(tx, tokenIDStr)
	if err != nil {
		return false, err
	}

	for keyImageStr := range listKeyImages {
		if _, ok := listTXOs[keyImageStr]; ok {
			return true, nil
		}
	}

	return false, nil
}

// SaveTxHistory saves a TxHistory in a csv file.
func SaveTxHistory(txHistory *TxHistory, filePath string) error {
	if len(filePath) == 0 {
		filePath = DefaultTxHistory
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	totalIn := uint64(0)
	for _, txIn := range txHistory.TxInList {
		totalIn += txIn.Amount
		err = w.Write([]string{txIn.GetTxHash(), txIn.String()})
		if err != nil {
			return fmt.Errorf("write txHash %v error: %v", txIn.GetTxHash(), err)
		}
	}
	err = w.Write([]string{"totalIn", fmt.Sprintf("%v", totalIn)})

	err = w.Write([]string{"-----", "-----"})
	if err != nil {
		return fmt.Errorf("cannot write csv file")
	}

	totalOut := uint64(0)
	for _, txOut := range txHistory.TxOutList {
		totalOut += txOut.Amount
		err = w.Write([]string{txOut.GetTxHash(), txOut.String()})
		if err != nil {
			return fmt.Errorf("write txHash %v error: %v", txOut.GetTxHash(), err)
		}
	}
	err = w.Write([]string{"totalOut", fmt.Sprintf("%v", totalOut)})

	return nil
}
