package incclient

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// CreatePDETradeTransaction creates a trading transaction with the provided version.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePDETradeTransaction(privateKey, tokenIDToSell, tokenIDToBuy string, amount, expectedBuy, tradingFee uint64, version int8) ([]byte, string, error) {
	if version == 2 {
		return client.CreatePDETradeTransactionVer2(privateKey, tokenIDToSell, tokenIDToBuy, amount, expectedBuy, tradingFee)
	} else if version == 1 {
		return client.CreatePDETradeTransactionVer1(privateKey, tokenIDToSell, tokenIDToBuy, amount, expectedBuy, tradingFee)
	} else { //Try either one of the version, if possible
		encodedTx, txHash, err := client.CreatePDETradeTransactionVer1(privateKey, tokenIDToSell, tokenIDToBuy, amount, expectedBuy, tradingFee)
		if err != nil {
			encodedTx, txHash, err1 := client.CreatePDETradeTransactionVer2(privateKey, tokenIDToSell, tokenIDToBuy, amount, expectedBuy, tradingFee)
			if err1 != nil {
				return nil, "", fmt.Errorf("cannot create raw pdetradetransaction for either version: %v, %v", err, err1)
			}
			return encodedTx, txHash, nil
		}
		return encodedTx, txHash, nil
	}
}

// CreatePDETradeTransactionVer1 creates a trading transaction version 1.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePDETradeTransactionVer1(privateKey, tokenIDToSell, tokenIDToBuy string, amount, expectedBuy, tradingFee uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	minAccept := expectedBuy
	//uncomment this code if you want to get the best price
	if minAccept == 0 {
		minAccept, err = client.CheckXPrice(tokenIDToSell, tokenIDToBuy, amount)
		if err != nil {
			return nil, "", err
		}
	}

	addr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	addr, err = wallet.GetPaymentAddressV1(addr, false)
	if err != nil {
		return nil, "", err
	}

	var pdeTradeMetadata *metadata.PDETradeRequest
	if tokenIDToSell == common.PRVIDStr || tokenIDToBuy == common.PRVIDStr {
		pdeTradeMetadata, err = metadata.NewPDETradeRequest(tokenIDToBuy, tokenIDToSell, amount, minAccept, tradingFee,
			addr, "", metadata.PDETradeRequestMeta)
	} else {
		pdeTradeMetadata, err = metadata.NewPDETradeRequest(tokenIDToBuy, tokenIDToSell, amount, minAccept, tradingFee,
			addr, "", metadata.PDECrossPoolTradeRequestMeta)
	}
	if err != nil {
		return nil, "", fmt.Errorf("cannot init trade request for %v to %v with amount %v: %v", tokenIDToSell, tokenIDToBuy, amount, err)
	}

	if tokenIDToSell == common.PRVIDStr {
		txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{amount + tradingFee}, 0, nil, pdeTradeMetadata, nil)
		return client.CreateRawTransaction(txParam, 1)
	} else {
		var tokenParam *TxTokenParam
		var txParam *TxParam
		if tokenIDToBuy == common.PRVIDStr {
			tokenParam = NewTxTokenParam(tokenIDToSell, 1, []string{common.BurningAddress2}, []uint64{amount + tradingFee}, false, 0, nil)
			txParam = NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, pdeTradeMetadata, nil)
		} else {
			tokenParam = NewTxTokenParam(tokenIDToSell, 1, []string{common.BurningAddress2}, []uint64{amount}, false, 0, nil)
			if tradingFee > 0 {
				txParam = NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{tradingFee}, 0, tokenParam, pdeTradeMetadata, nil)
			} else {
				txParam = NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, pdeTradeMetadata, nil)
			}
		}
		return client.CreateRawTokenTransaction(txParam, 1)
	}
}

// CreatePDETradeTransactionVer2 creates a trading transaction version 2.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePDETradeTransactionVer2(privateKey, tokenIDToSell, tokenIDToBuy string, amount uint64, expectedBuy, tradingFee uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	minAccept := expectedBuy
	////uncomment this code if you want to get the best price
	//minAccept, err = CheckPrice(tokenIDToSell, tokenIDToBuy, amount)
	//if err != nil {
	//	return nil, "", err
	//}
	addr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	pubKeyStr, txRandomStr, err := GenerateOTAFromPaymentAddress(addr)
	if err != nil {
		return nil, "", err
	}

	subPubKeyStr, subTxRandomStr, err := GenerateOTAFromPaymentAddress(addr)
	if err != nil {
		return nil, "", err
	}

	pdeTradeMetadata, err := metadata.NewPDECrossPoolTradeRequest(tokenIDToBuy, tokenIDToSell, amount, minAccept, tradingFee,
		pubKeyStr, txRandomStr, subPubKeyStr, subTxRandomStr, metadata.PDECrossPoolTradeRequestMeta)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init trade request for %v to %v with amount %v: %v", tokenIDToSell, tokenIDToBuy, amount, err)
	}

	if tokenIDToSell == common.PRVIDStr {
		txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{amount + tradingFee}, 0, nil, pdeTradeMetadata, nil)
		return client.CreateRawTransaction(txParam, 2)
	} else {
		tokenParam := NewTxTokenParam(tokenIDToSell, 1, []string{common.BurningAddress2}, []uint64{amount}, false, 0, nil)
		var txParam *TxParam
		if tradingFee > 0 {
			txParam = NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{tradingFee}, 0, tokenParam, pdeTradeMetadata, nil)
		} else {
			txParam = NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, pdeTradeMetadata, nil)
		}

		return client.CreateRawTokenTransaction(txParam, 2)
	}

}

// CreateAndSendPDETradeTransaction creates a trading transaction with the provided version, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPDETradeTransaction(privateKey, tokenIDToSell, tokenIDToBuy string, amount, expectedBuy, tradingFee uint64) (string, error) {
	encodedTx, txHash, err := client.CreatePDETradeTransaction(privateKey, tokenIDToSell, tokenIDToBuy, amount, expectedBuy, tradingFee, -1)
	if err != nil {
		return "", err
	}

	if tokenIDToSell == common.PRVIDStr {
		err = client.SendRawTx(encodedTx)
		if err != nil {
			return "", err
		}
	} else {
		err = client.SendRawTokenTx(encodedTx)
		if err != nil {
			return "", err
		}
	}

	return txHash, nil
}

// CreatePDEContributeTransaction creates a contributing transaction which contributes an amount of tokenID to the pDEX.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePDEContributeTransaction(privateKey, pairID, tokenID string, amount uint64, version int8) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	addr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	if version == 1 {
		addr, err = wallet.GetPaymentAddressV1(addr, false)
		if err != nil {
			return nil, "", err
		}
	}

	md, err := metadata.NewPDEContribution(pairID, addr, amount, tokenID, metadata.PDEPRVRequiredContributionRequestMeta)
	if err != nil {
		return nil, "", err
	}

	if tokenID == common.PRVIDStr {
		txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{amount}, 0, nil, md, nil)
		return client.CreateRawTransaction(txParam, version)
	} else {
		tokenParam := NewTxTokenParam(tokenID, 1, []string{common.BurningAddress2}, []uint64{amount}, false, 0, nil)
		txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)

		return client.CreateRawTokenTransaction(txParam, version)
	}
}

// CreateAndSendPDEContributeTransaction creates a contributing transaction which contributes an amount of tokenID to the pDEX, and then submits it to the Incognito network.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPDEContributeTransaction(privateKey, pairID, tokenID string, amount uint64, version int8) (string, error) {
	encodedTx, txHash, err := client.CreatePDEContributeTransaction(privateKey, pairID, tokenID, amount, version)
	if err != nil {
		return "", err
	}

	if tokenID == common.PRVIDStr {
		err = client.SendRawTx(encodedTx)
		if err != nil {
			return "", err
		}
	} else {
		err = client.SendRawTokenTx(encodedTx)
		if err != nil {
			return "", err
		}
	}

	return txHash, nil
}

// CreatePDEWithdrawalTransaction creates a withdrawing transaction which withdraws a pair of tokenIDs from the pDEX.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePDEWithdrawalTransaction(privateKey, tokenID1, tokenID2 string, sharedAmount uint64, version int8) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	addr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	if version == 1 {
		addr, err = wallet.GetPaymentAddressV1(addr, false)
		if err != nil {
			return nil, "", err
		}
	}

	md, err := metadata.NewPDEWithdrawalRequest(addr, tokenID2, tokenID1, sharedAmount, metadata.PDEWithdrawalRequestMeta)

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, md, nil)

	return client.CreateRawTransaction(txParam, version)
}

// CreateAndSendPDEWithdrawalTransaction creates a withdrawing transaction which withdraws a pair of tokenIDs from the pDEX, and submits it to the Incognito network.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPDEWithdrawalTransaction(privateKey, tokenID1, tokenID2 string, sharedAmount uint64, version int8) (string, error) {
	encodedTx, txHash, err := client.CreatePDEWithdrawalTransaction(privateKey, tokenID1, tokenID2, sharedAmount, version)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreatePDEFeeWithdrawalTransaction creates a withdrawing transaction which withdraws pDEX LP fee.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePDEFeeWithdrawalTransaction(privateKey, tokenID1, tokenID2 string, version int8) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	addr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	if version == 1 {
		addr, err = wallet.GetPaymentAddressV1(addr, false)
		if err != nil {
			return nil, "", err
		}
	}
	withdrawnAmount, err := client.GetLPFeeAmount(0, tokenID1, tokenID2, addr)
	if err != nil {
		return nil, "", err
	}
	if withdrawnAmount == 0 {
		return nil, "", fmt.Errorf("no trading fee to collect")
	}

	md, err := metadata.NewPDEFeeWithdrawalRequest(addr, tokenID2, tokenID1, withdrawnAmount, metadata.PDEFeeWithdrawalRequestMeta)

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, md, nil)

	return client.CreateRawTransaction(txParam, version)
}

// CreateAndSendPDEFeeWithdrawalTransaction creates a withdrawing transaction which withdraws pDEX LP fees, and submits it to the Incognito network.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPDEFeeWithdrawalTransaction(privateKey, tokenID1, tokenID2 string, version int8) (string, error) {
	encodedTx, txHash, err := client.CreatePDEFeeWithdrawalTransaction(privateKey, tokenID1, tokenID2, version)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
