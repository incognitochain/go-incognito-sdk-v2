package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"strings"

	// "github.com/incognitochain/go-incognito-sdk-v2/metadata"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
	metadataPdexv3 "github.com/incognitochain/go-incognito-sdk-v2/metadata/pdexv3"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

const MintNftRequiredAmount = 1000000000

// CreatePdexv3MintNFT creates a transaction minting a new pDEX NFT for the given private key.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3MintNFT(privateKey string) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}
	otaReceiver := coin.OTAReceiver{}
	err = otaReceiver.FromAddress(senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}
	otaReceiveStr, err := otaReceiver.String()
	md := metadataPdexv3.NewUserMintNftRequestWithValue(otaReceiveStr, MintNftRequiredAmount)

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{MintNftRequiredAmount}, 0, nil, md, nil)

	return client.CreateRawTransaction(txParam, 2)
}

// CreateAndSendPdexv3UserMintNFTransaction creates a transaction minting a new pDEX NFT for the given private key, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3UserMintNFTransaction(privateKey string) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3MintNFT(privateKey)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreatePdexv3Trade creates a trading transaction.
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

// CreateAndSendPdexv3TradeTransaction creates a trading transaction (version 2 only), and submits it to the Incognito network.
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

// CreatePdexv3AddOrder creates a transaction that adds a new order in pdex v3.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3AddOrder(privateKey, pairID string,
	tokenIDToSellStr, tokenIDToBuyStr, nftIDStr string,
	sellAmount, expectedBuy uint64,
) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	tokenSell, err := common.Hash{}.NewHashFromStr(tokenIDToSellStr)
	if err != nil {
		return nil, "", err
	}
	tokenBuy, err := common.Hash{}.NewHashFromStr(tokenIDToBuyStr)
	if err != nil {
		return nil, "", err
	}
	nftID, err := common.Hash{}.NewHashFromStr(nftIDStr)
	if err != nil {
		return nil, "", err
	}

	// construct order metadata
	md, _ := metadataPdexv3.NewAddOrderRequest(
		*tokenSell, pairID, sellAmount,
		expectedBuy, nil, *nftID,
		metadataCommon.Pdexv3AddOrderRequestMeta,
	)
	// create one-time receivers for response TX
	isPRV := md.TokenToSell == common.PRVCoinID
	tokenList := []common.Hash{md.TokenToSell, *tokenBuy}
	md.Receiver, err = GenerateOTAReceivers(
		tokenList, senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}

	if isPRV {
		txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{sellAmount}, 0, nil, md, nil)
		return client.CreateRawTransaction(txParam, 2)
	} else {
		tokenParam := NewTxTokenParam(tokenIDToSellStr, 1, []string{common.BurningAddress2}, []uint64{sellAmount}, false, 0, nil)
		txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)
		return client.CreateRawTokenTransaction(txParam, 2)
	}
}

// CreateAndSendPdexv3AddOrderTransaction creates an order transaction, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3AddOrderTransaction(privateKey, pairID string,
	tokenIDToSellStr, tokenIDToBuyStr, nftIDStr string,
	sellAmount, expectedBuy uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3AddOrder(privateKey, pairID, tokenIDToSellStr, tokenIDToBuyStr, nftIDStr, sellAmount, expectedBuy)
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

// CreatePdexv3WithdrawOrder creates a transaction that withdraws all outstanding balance from an order in pdex v3.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3WithdrawOrder(privateKey, pairID, orderID string,
	nftIDStr string, amount uint64, withdrawTokenIDs ...string) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	nftID, err := common.Hash{}.NewHashFromStr(nftIDStr)
	if err != nil {
		return nil, "", err
	}

	tokenList := []common.Hash{*nftID}
	if len(withdrawTokenIDs) == 0 {
		withdrawTokenIDs, err = getTokenIDsFromPairID(pairID)
		if err != nil {
			return nil, "", err
		}
	}
	for _, v := range withdrawTokenIDs {
		temp, err := common.Hash{}.NewHashFromStr(v)
		if err != nil {
			return nil, "", err
		}
		tokenList = append(tokenList, *temp)
	}

	otaReceivers, err := GenerateOTAReceivers(tokenList, senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}
	md, _ := metadataPdexv3.NewWithdrawOrderRequest(pairID, orderID, amount,
		otaReceivers, *nftID, metadataCommon.Pdexv3WithdrawOrderRequestMeta)

	tokenParam := NewTxTokenParam(nftIDStr, 1, []string{common.BurningAddress2}, []uint64{1}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)

	return client.CreateRawTokenTransaction(txParam, 2)
}

// CreateAndSendPdexv3WithdrawOrderTransaction creates an order-withdrawing transaction, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3WithdrawOrderTransaction(privateKey, pairID, orderID string,
	nftIDStr string, amount uint64, withdrawTokenIDs ...string) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3WithdrawOrder(privateKey, pairID, orderID, nftIDStr, amount, withdrawTokenIDs...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreatePdexv3Contribute creates a contributing transaction which contributes an amount of tokenID to the pDEX.
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

// CreatePdexv3WithdrawLiquidity creates a withdrawing transaction which withdraws a pair of tokenIDs from pdex v3.
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

	return client.CreateRawTokenTransaction(txParam, 2)
}

// CreateAndSendPdexv3WithdrawLiquidityTransaction creates a withdrawing transaction which withdraws a pair of tokenIDs from the pDEX, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3WithdrawLiquidityTransaction(privateKey, pairID, token0IDStr, token1IDStr, nftIDStr string, shareAmount uint64) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3WithdrawLiquidity(privateKey, pairID, token0IDStr, token1IDStr, nftIDStr, shareAmount)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreatePdexv3WithdrawLPFee creates a transaction that withdraws all LP fee rewards earned by a liquidity provider in one pool in pdex v3.
// If `withdrawTokenIDs` are not specified, it will get the tokenIDs from the given pairID.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3WithdrawLPFee(privateKey, pairID string,
	nftIDStr string, withdrawTokenIDs ...string) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	nftID, err := common.Hash{}.NewHashFromStr(nftIDStr)
	if err != nil {
		return nil, "", err
	}

	tokenList := []common.Hash{*nftID, common.PRVCoinID, common.PDEXCoinID}
	if len(withdrawTokenIDs) == 0 {
		withdrawTokenIDs, err = getTokenIDsFromPairID(pairID)
		if err != nil {
			return nil, "", err
		}
	}
	for _, v := range withdrawTokenIDs {
		temp, err := common.Hash{}.NewHashFromStr(v)
		if err != nil {
			return nil, "", err
		}
		tokenList = append(tokenList, *temp)
	}

	otaReceivers, err := GenerateOTAReceivers(tokenList, senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}
	md, _ := metadataPdexv3.NewPdexv3WithdrawalLPFeeRequest(
		metadataCommon.Pdexv3WithdrawLPFeeRequestMeta,
		pairID,
		*nftID,
		otaReceivers,
	)

	tokenParam := NewTxTokenParam(nftIDStr, 1, []string{common.BurningAddress2}, []uint64{1}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)

	return client.CreateRawTokenTransaction(txParam, 2)
}

// CreateAndSendPdexv3WithdrawLPFeeTransaction creates a transaction that withdraws a liquidity provider's reward in one pool, and submits it to the Incognito network.
// If `withdrawTokenIDs` are not specified, it will get the tokenIDs from the given pairID.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3WithdrawLPFeeTransaction(privateKey, pairID string,
	nftIDStr string, withdrawTokenIDs ...string) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3WithdrawLPFee(privateKey, pairID, nftIDStr, withdrawTokenIDs...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreatePdexv3WithdrawProtocolFee creates a transaction that withdraws all protocol fee rewards earned by a liquidity provider in one pool in pdex v3.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3WithdrawProtocolFee(privateKey, pairID string) ([]byte, string, error) {
	md, _ := metadataPdexv3.NewPdexv3WithdrawalProtocolFeeRequest(
		metadataCommon.Pdexv3WithdrawProtocolFeeRequestMeta,
		pairID,
	)

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, md, nil)

	return client.CreateRawTransaction(txParam, 2)
}

// CreateAndSendPdexv3WithdrawProtocolFeeTransaction creates a protocol-fee-withdrawing transaction, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3WithdrawProtocolFeeTransaction(privateKey, pairID string) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3WithdrawProtocolFee(privateKey, pairID)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreatePdexv3Staking creates a staking transaction in pdex v3.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3Staking(privateKey, tokenIDStr, nftIDStr string, amount uint64) ([]byte, string, error) {
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
	md := metadataPdexv3.NewStakingRequestWithValue(tokenIDStr, nftIDStr, otaReceiverStr, amount)

	if isPRV {
		txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{amount}, 0, nil, md, nil)
		return client.CreateRawTransaction(txParam, 2)
	} else {
		tokenParam := NewTxTokenParam(tokenIDStr, 1, []string{common.BurningAddress2}, []uint64{amount}, false, 0, nil)
		txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)

		return client.CreateRawTokenTransaction(txParam, 2)
	}
}

// CreateAndSendPdexv3StakingTransaction creates a staking transaction in pdex v3, and then submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3StakingTransaction(privateKey, tokenIDStr, nftIDStr string, amount uint64) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3Staking(privateKey, tokenIDStr, nftIDStr, amount)
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

// CreatePdexv3Unstaking creates an unstaking transaction in pdex v3.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3Unstaking(privateKey, tokenIDStr, nftIDStr string,
	amount uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	nftID, err := common.Hash{}.NewHashFromStr(nftIDStr)
	if err != nil {
		return nil, "", err
	}
	tokenID, err := common.Hash{}.NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	tokenList := []common.Hash{*nftID, *tokenID}
	otaReceivers, err := GenerateOTAReceivers(tokenList, senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}
	md := metadataPdexv3.NewUnstakingRequestWithValue(
		tokenIDStr, nftIDStr, toStringKeys(otaReceivers), amount,
	)

	tokenParam := NewTxTokenParam(nftIDStr, 1, []string{common.BurningAddress2}, []uint64{1}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)

	return client.CreateRawTokenTransaction(txParam, 2)
}

// CreateAndSendPdexv3UnstakingTransaction creates an unstaking transaction, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3UnstakingTransaction(privateKey, tokenIDStr string,
	nftIDStr string, amount uint64) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3Unstaking(privateKey, tokenIDStr, nftIDStr, amount)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreatePdexv3WithdrawStakeRewardTransaction creates a transaction that withdraws all rewards (from trading fees) earned by staking in one pool in pdex v3.
// If `withdrawTokenIDs` are not specified, it will get all available staking tokens.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePdexv3WithdrawStakeRewardTransaction(
	privateKey, stakingPoolIDStr, nftIDStr string, withdrawTokenIDs ...string) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	nftID, err := common.Hash{}.NewHashFromStr(nftIDStr)
	if err != nil {
		return nil, "", err
	}

	stakingPoolID, err := common.Hash{}.NewHashFromStr(stakingPoolIDStr)
	if err != nil {
		return nil, "", err
	}

	tokenList := []common.Hash{*nftID, common.PRVCoinID, *stakingPoolID}
	if len(withdrawTokenIDs) == 0 {
		tmpTokenIDs, err := client.GetListStakingRewardTokens(0)
		if err != nil {
			return nil, "", err
		}
		tokenList = append(tokenList, tmpTokenIDs...)
	} else {
		for _, v := range withdrawTokenIDs {
			temp, err := common.Hash{}.NewHashFromStr(v)
			if err != nil {
				return nil, "", err
			}
			tokenList = append(tokenList, *temp)
		}
	}

	otaReceivers, err := GenerateOTAReceivers(tokenList, senderWallet.KeySet.PaymentAddress)
	if err != nil {
		return nil, "", err
	}
	md, _ := metadataPdexv3.NewPdexv3WithdrawalStakingRewardRequest(
		metadataCommon.Pdexv3WithdrawStakingRewardRequestMeta,
		stakingPoolIDStr,
		*nftID,
		otaReceivers,
	)

	tokenParam := NewTxTokenParam(nftIDStr, 1, []string{common.BurningAddress2}, []uint64{1}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, md, nil)

	return client.CreateRawTokenTransaction(txParam, 2)
}

// CreateAndSendPdexv3WithdrawStakeRewardTransaction creates a transaction that withdraws all rewards (from trading fees) earned by staking in one pool in pdex v3.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPdexv3WithdrawStakeRewardTransaction(
	privateKey, stakingPoolIDStr, nftIDStr string, withdrawTokenIDs ...string,
) (string, error) {
	encodedTx, txHash, err := client.CreatePdexv3WithdrawStakeRewardTransaction(privateKey, stakingPoolIDStr, nftIDStr, withdrawTokenIDs...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

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

func getTokenIDsFromPairID(pairID string) ([]string, error) {
	res := strings.Split(pairID, "-")
	if len(res) != 3 {
		return nil, fmt.Errorf("invalid pairID %v", pairID)
	}
	res = res[:2]
	for _, tokenIDStr := range res {
		_, err := common.Hash{}.NewHashFromStr(tokenIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid tokenID %v in pairID %v", tokenIDStr, pairID)
		}
	}

	return res, nil
}
