package incclient

import (
	"encoding/csv"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"os"
)

const (
	DefaultTxIn      = "txIn_history.csv"
	DefaultTxOut     = "txOut_history.csv"
	DefaultTxHistory = "txHistory.csv"
)

var (
	committeePrefix = "[Committee]"
	bridgePrefix    = "[Bridge]"
	portalPrefix    = "[Portal]"
	pDEXV2Prefix    = "[pDEX v2]"
	pDEXV3Prefix    = "[pDEX v3]"
)

var txMetadataNote = map[int]string{
	metadata.InvalidMeta: "",

	metadata.InitTokenRequestMeta:  fmt.Sprintf("Init Token Request"),
	metadata.InitTokenResponseMeta: fmt.Sprintf("Init Token Response"),

	// Committee
	metadata.ShardStakingMeta:           fmt.Sprintf("%v Staking", committeePrefix),
	metadata.BeaconStakingMeta:          fmt.Sprintf("%v Staking", committeePrefix),
	metadata.UnStakingMeta:              fmt.Sprintf("%v Un-staking", committeePrefix),
	metadata.StopAutoStakingMeta:        fmt.Sprintf("%v Stop-staking", committeePrefix),
	metadata.WithDrawRewardRequestMeta:  fmt.Sprintf("%v Withdraw Reward Request", committeePrefix),
	metadata.WithDrawRewardResponseMeta: fmt.Sprintf("%v Withdraw Reward Response", committeePrefix),

	// Bridge
	metadata.IssuingRequestMeta:                   fmt.Sprintf("%v Shield Request", bridgePrefix),
	metadata.IssuingResponseMeta:                  fmt.Sprintf("%v Shield Response", bridgePrefix),
	metadata.IssuingETHRequestMeta:                fmt.Sprintf("%v Shield Request", bridgePrefix),
	metadata.IssuingETHResponseMeta:               fmt.Sprintf("%v Shield Response", bridgePrefix),
	metadata.IssuingBSCRequestMeta:                fmt.Sprintf("%v Shield Request", bridgePrefix),
	metadata.IssuingBSCResponseMeta:               fmt.Sprintf("%v Shield Response", bridgePrefix),
	metadata.IssuingPLGRequestMeta:                fmt.Sprintf("%v Shield Request", bridgePrefix),
	metadata.IssuingPLGResponseMeta:               fmt.Sprintf("%v Shield Response", bridgePrefix),
	metadata.IssuingPRVERC20RequestMeta:           fmt.Sprintf("%v Shield Request", bridgePrefix),
	metadata.IssuingPRVERC20ResponseMeta:          fmt.Sprintf("%v Shield Response", bridgePrefix),
	metadata.IssuingPRVBEP20RequestMeta:           fmt.Sprintf("%v Shield Request", bridgePrefix),
	metadata.IssuingPRVBEP20ResponseMeta:          fmt.Sprintf("%v Shield Response", bridgePrefix),
	metadata.BurningRequestMeta:                   fmt.Sprintf("%v Unshield Request", bridgePrefix),
	metadata.BurningRequestMetaV2:                 fmt.Sprintf("%v Unshield Request", bridgePrefix),
	metadata.ContractingRequestMeta:               fmt.Sprintf("%v Unshield Request", bridgePrefix),
	metadata.BurningPBSCRequestMeta:               fmt.Sprintf("%v Unshield Request", bridgePrefix),
	metadata.BurningPLGRequestMeta:                fmt.Sprintf("%v Unshield Request", bridgePrefix),
	metadata.BurningForDepositToSCRequestMeta:     fmt.Sprintf("%v Unshield Request", bridgePrefix),
	metadata.BurningForDepositToSCRequestMetaV2:   fmt.Sprintf("%v Unshield Request", bridgePrefix),
	metadata.BurningPBSCForDepositToSCRequestMeta: fmt.Sprintf("%v Unshield Request", bridgePrefix),
	metadata.BurningPLGForDepositToSCRequestMeta:  fmt.Sprintf("%v Unshield Request", bridgePrefix),

	// Portal
	metadata.PortalV4ShieldingRequestMeta:      fmt.Sprintf("%v Shield Request", portalPrefix),
	metadata.PortalV4ShieldingResponseMeta:     fmt.Sprintf("%v Shield Response", portalPrefix),
	metadata.PortalV4UnshieldingRequestMeta:    fmt.Sprintf("%v Unshield Request", portalPrefix),
	metadata.PortalV4UnshieldingResponseMeta:   fmt.Sprintf("%v Unshield Response", portalPrefix),
	metadata.PortalV4ConvertVaultRequestMeta:   fmt.Sprintf("%v Convert Vault", portalPrefix),
	metadata.PortalV4SubmitConfirmedTxMeta:     fmt.Sprintf("%v Submit Confirmed Tx", portalPrefix),
	metadata.PortalV4UnshieldBatchingMeta:      fmt.Sprintf("%v Batch Unshield", portalPrefix),
	metadata.PortalV4FeeReplacementRequestMeta: fmt.Sprintf("%v Fee Replacement Request", portalPrefix),

	// pDEX v2
	metadata.PDETradeRequestMeta:                   fmt.Sprintf("%v Trade Request", pDEXV2Prefix),
	metadata.PDETradeResponseMeta:                  fmt.Sprintf("%v Trade Response", pDEXV2Prefix),
	metadata.PDECrossPoolTradeRequestMeta:          fmt.Sprintf("%v Trade Request", pDEXV2Prefix),
	metadata.PDECrossPoolTradeResponseMeta:         fmt.Sprintf("%v Trade Response", pDEXV2Prefix),
	metadata.PDEContributionMeta:                   fmt.Sprintf("%v Contribution Request", pDEXV2Prefix),
	metadata.PDEPRVRequiredContributionRequestMeta: fmt.Sprintf("%v Contribution Request", pDEXV2Prefix),
	metadata.PDEContributionResponseMeta:           fmt.Sprintf("%v Contribution Response", pDEXV2Prefix),
	metadata.PDEWithdrawalRequestMeta:              fmt.Sprintf("%v Withdrawal Request", pDEXV2Prefix),
	metadata.PDEWithdrawalResponseMeta:             fmt.Sprintf("%v Withdrawal Response", pDEXV2Prefix),
	metadata.PDEFeeWithdrawalRequestMeta:           fmt.Sprintf("%v Fee Withdrawal Request", pDEXV2Prefix),
	metadata.PDEFeeWithdrawalResponseMeta:          fmt.Sprintf("%v Fee Withdrawal Response", pDEXV2Prefix),

	// pDEX v3
	metadataCommon.Pdexv3ModifyParamsMeta:                  fmt.Sprintf("%v Modify Params", pDEXV3Prefix),
	metadataCommon.Pdexv3AddLiquidityRequestMeta:           fmt.Sprintf("%v Contribution Request", pDEXV3Prefix),
	metadataCommon.Pdexv3AddLiquidityResponseMeta:          fmt.Sprintf("%v Contribution Response", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawLiquidityRequestMeta:      fmt.Sprintf("%v Withdrawal Request", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawLiquidityResponseMeta:     fmt.Sprintf("%v Withdrawal Response", pDEXV3Prefix),
	metadataCommon.Pdexv3TradeRequestMeta:                  fmt.Sprintf("%v Trade Request", pDEXV3Prefix),
	metadataCommon.Pdexv3TradeResponseMeta:                 fmt.Sprintf("%v Trade Response", pDEXV3Prefix),
	metadataCommon.Pdexv3AddOrderRequestMeta:               fmt.Sprintf("%v Add Order Request", pDEXV3Prefix),
	metadataCommon.Pdexv3AddOrderResponseMeta:              fmt.Sprintf("%v Add Order Response", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawOrderRequestMeta:          fmt.Sprintf("%v Remove Order Request", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawOrderResponseMeta:         fmt.Sprintf("%v Remove Order Response", pDEXV3Prefix),
	metadataCommon.Pdexv3UserMintNftRequestMeta:            fmt.Sprintf("%v Mint NFT Request", pDEXV3Prefix),
	metadataCommon.Pdexv3UserMintNftResponseMeta:           fmt.Sprintf("%v Mint NFT Response", pDEXV3Prefix),
	metadataCommon.Pdexv3MintNftRequestMeta:                fmt.Sprintf("%v Mint NFT Request", pDEXV3Prefix),
	metadataCommon.Pdexv3MintNftResponseMeta:               fmt.Sprintf("%v Mint NFT Response", pDEXV3Prefix),
	metadataCommon.Pdexv3StakingRequestMeta:                fmt.Sprintf("%v Staking Request", pDEXV3Prefix),
	metadataCommon.Pdexv3StakingResponseMeta:               fmt.Sprintf("%v Staking Response", pDEXV3Prefix),
	metadataCommon.Pdexv3UnstakingRequestMeta:              fmt.Sprintf("%v Staking Request", pDEXV3Prefix),
	metadataCommon.Pdexv3UnstakingResponseMeta:             fmt.Sprintf("%v Staking Response", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawLPFeeRequestMeta:          fmt.Sprintf("%v Withdraw LP Fee Request", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawLPFeeResponseMeta:         fmt.Sprintf("%v Withdraw LP Fee Response", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawProtocolFeeRequestMeta:    fmt.Sprintf("%v Withdraw Protocol Fee Request", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawProtocolFeeResponseMeta:   fmt.Sprintf("%v Withdraw Protocol Fee Response", pDEXV3Prefix),
	metadataCommon.Pdexv3MintBlockRewardMeta:               fmt.Sprintf("%v Block Reward", pDEXV3Prefix),
	metadataCommon.Pdexv3DistributeStakingRewardMeta:       fmt.Sprintf("%v Distribute Staking Reward", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawStakingRewardRequestMeta:  fmt.Sprintf("%v Withdraw Staking Reward Request", pDEXV3Prefix),
	metadataCommon.Pdexv3WithdrawStakingRewardResponseMeta: fmt.Sprintf("%v Withdraw Staking Reward Response", pDEXV3Prefix),
}

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

	_ = f.Truncate(0)

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

	Logger.Printf("Finished storing history to file %v\n", filePath)
	return nil
}
