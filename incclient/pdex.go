package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	// "github.com/incognitochain/go-incognito-sdk-v2/metadata"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
	metadataPdexv3 "github.com/incognitochain/go-incognito-sdk-v2/metadata/pdexv3"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func GenerateOTAReceivers(
	tokens []common.Hash, addr key.PaymentAddress,
) (map[common.Hash]coin.OTAReceiver, error) {
	result := make(map[common.Hash]coin.OTAReceiver)
	var err error
	for _, tokenID := range tokens {
		temp := coin.OTAReceiver{}
		err = temp.FromAddress(addr)
		if err != nil {
			return nil, err
		}
		result[tokenID] = temp
	}
	return result, nil
}

func toStringKeys(inputMap map[common.Hash]coin.OTAReceiver) map[string]string {
	result := make(map[string]string)
	for k, v := range inputMap {
		s, _ := v.String()
		result[k.String()] = s
	}
	return result
}

// CreatePdexv3TradeVer2 creates a trading transaction version 2.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3Trade(privateKey string, tradePath []string, tokenIDToSellStr,
	tokenIDToBuyStr string, amount uint64, expectedBuy, tradingFee uint64, feeInPRV bool,
) ([]byte, string, error) {
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

	tokenSell, err := common.Hash{}.NewHashFromStr(tokenIDToSellStr)
	if err != nil {
		return nil, "", err
	}
	tokenBuy, err := common.Hash{}.NewHashFromStr(tokenIDToBuyStr)
	if err != nil {
		return nil, "", err
	}

	// construct trade metadata
	md, _ := metadataPdexv3.NewTradeRequest(
		tradePath, *tokenSell, amount,
		minAccept, tradingFee, nil,
		metadataCommon.Pdexv3TradeRequestMeta,
	)
	// create one-time receivers for response TX
	isPRV := md.TokenToSell == common.PRVCoinID
	tokenList := []common.Hash{md.TokenToSell, *tokenBuy}
	// add a receiver for PRV if necessary
	if feeInPRV && !isPRV && *tokenBuy != common.PRVCoinID {
		tokenList = append(tokenList, common.PRVCoinID)
	}
	md.Receiver, err = GenerateOTAReceivers(
		tokenList, senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}

	if isPRV {
		txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{amount + tradingFee}, 0, nil, md, nil)
		return client.CreateRawTransaction(txParam, 2)
	} else {
		var txParam *TxParam
		if feeInPRV {
			tokenParam := NewTxTokenParam(tokenIDToSellStr, 1, []string{common.BurningAddress2}, []uint64{amount}, false, 0, nil)
			txParam = NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{tradingFee}, 0, tokenParam, md, nil)
		} else {
			tokenParam := NewTxTokenParam(tokenIDToSellStr, 1, []string{common.BurningAddress2}, []uint64{amount + tradingFee}, false, 0, nil)
			txParam = NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)
		}
		return client.CreateRawTokenTransaction(txParam, 2)
	}

}

// CreateAndSendPdexv3TradeTransaction creates a trading transaction with the provided version, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3TradeTransaction(privateKey string, tradePath []string, tokenIDToSellStr, tokenIDToBuyStr string, amount uint64,
	expectedBuy, tradingFee uint64, feeInPRV bool,
) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3Trade(privateKey, tradePath, tokenIDToSellStr, tokenIDToBuyStr, amount, expectedBuy, tradingFee, feeInPRV)
	if err != nil {
		return "", err
	}

	if tokenIDToSellStr == common.PRVIDStr {
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

// CreatePdexv3Contribute creates a contributing transaction which contributes an amount of tokenID to the pDEX.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3Contribute(privateKey, pairID, pairHash, tokenIDStr, nftIDStr string, amount uint64, amplifier uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	tokenID, err := common.Hash{}.NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}
	isPRV := *tokenID == common.PRVCoinID

	// construct metadata for contribution
	temp := coin.OTAReceiver{}
	err = temp.FromAddress(senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}
	otaReceiverStr, _ := temp.String()
	md := metadataPdexv3.NewAddLiquidityRequestWithValue(
		pairID, pairHash, otaReceiverStr,
		tokenIDStr, nftIDStr, amount, uint(amplifier),
	)

	if isPRV {
		txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{amount}, 0, nil, md, nil)
		return client.CreateRawTransaction(txParam, 2)
	} else {
		tokenParam := NewTxTokenParam(tokenIDStr, 1, []string{common.BurningAddress2}, []uint64{amount}, false, 0, nil)
		txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)

		return client.CreateRawTokenTransaction(txParam, 2)
	}
}

// CreateAndSendPdexv3ContributeTransaction creates a contributing transaction which contributes an amount of tokenID to the pDEX, and then submits it to the Incognito network.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3ContributeTransaction(privateKey, pairID, pairHash, tokenIDStr, nftIDStr string, amount uint64, amplifier uint64) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3Contribute(privateKey, pairID, pairHash, tokenIDStr, nftIDStr, amount, amplifier)
	if err != nil {
		return "", err
	}

	if tokenIDStr == common.PRVIDStr {
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

// CreatePdexv3Withdraw creates a withdrawing transaction which withdraws a pair of tokenIDs from the pDEX.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3WithdrawLiquidity(privateKey, pairID, token0IDStr, token1IDStr, nftIDStr string, shareAmount uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	token0ID, err := common.Hash{}.NewHashFromStr(token0IDStr)
	if err != nil {
		return nil, "", err
	}
	token1ID, err := common.Hash{}.NewHashFromStr(token1IDStr)
	if err != nil {
		return nil, "", err
	}
	nftID, err := common.Hash{}.NewHashFromStr(nftIDStr)
	if err != nil {
		return nil, "", err
	}
	
	tokenList := []common.Hash{*token0ID, *token1ID, *nftID}
	otaReceivers, err := GenerateOTAReceivers(tokenList, senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}
	md := metadataPdexv3.NewWithdrawLiquidityRequestWithValue(
		pairID, nftIDStr,
		toStringKeys(otaReceivers), shareAmount,
	)

	tokenParam := NewTxTokenParam(nftIDStr, 1, []string{common.BurningAddress2}, []uint64{1}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)

	return client.CreateRawTransaction(txParam, 2)
}

// CreateAndSendPdexv3WithdrawalTransaction creates a withdrawing transaction which withdraws a pair of tokenIDs from the pDEX, and submits it to the Incognito network.
// Version = -1 indicates that whichever version is accepted.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3WithdrawalTransaction(privateKey, pairID, token0IDStr, token1IDStr, nftIDStr string, shareAmount uint64) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3WithdrawLiquidity(privateKey, pairID, token0IDStr, token1IDStr, nftIDStr, shareAmount)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
