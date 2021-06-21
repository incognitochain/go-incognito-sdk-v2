package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func (client *IncClient) CreatePortalShieldTransaction(
	privateKey, tokenID, paymentAddr, shieldingProof string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalShieldingMetadata, err := metadata.NewPortalShieldingRequest(
		metadata.PortalV4ShieldingRequestMeta,
		tokenID,
		paymentAddr,
		shieldingProof,
	)
	if err != nil {
		return nil, "", err
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, portalShieldingMetadata, nil)
	if len(inputCoins) > 0 {
		return client.CreateRawTransactionWithInputCoins(txParam, inputCoins, coinIndices)
	}
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalShieldTransaction(
	privateKey, tokenID, paymentAddr, shieldingProof string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreatePortalShieldTransaction(privateKey, tokenID, paymentAddr, shieldingProof, inputCoins, coinIndices)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalReplaceByFeeTransaction(
	privateKey, tokenID, batchID string, fee uint, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalRBFMetadata, err := metadata.NewPortalReplacementFeeRequest(
		metadata.PortalV4FeeReplacementRequestMeta,
		tokenID,
		batchID,
		fee,
	)
	if err != nil {
		return nil, "", err
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, portalRBFMetadata, nil)
	if len(inputCoins) > 0 {
		return client.CreateRawTransactionWithInputCoins(txParam, inputCoins, coinIndices)
	}
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalReplaceByFeeTransaction(
	privateKey, tokenID, batchID string, fee uint, inputCoins []coin.PlainCoin, coinIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreatePortalReplaceByFeeTransaction(privateKey, tokenID, batchID, fee, inputCoins, coinIndices)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalSubmitConfirmationTransaction(
	privateKey, tokenID, unShieldProof, batchID string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalSubmitConfirmationMetadata, err := metadata.NewPortalSubmitConfirmedTxRequest(
		metadata.PortalV4SubmitConfirmedTxMeta,
		unShieldProof,
		tokenID,
		batchID,
	)
	if err != nil {
		return nil, "", err
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, portalSubmitConfirmationMetadata, nil)
	if len(inputCoins) > 0 {
		return client.CreateRawTransactionWithInputCoins(txParam, inputCoins, coinIndices)
	}
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalSubmitConfirmationTransaction(
	privateKey, tokenID, unShieldProof, batchID string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreatePortalSubmitConfirmationTransaction(privateKey, tokenID, unShieldProof, batchID, inputCoins, coinIndices)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalConvertVaultTransaction(
	privateKey, tokenID, paymentAddr, convertingProof string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalConvertVaultMetadata, err := metadata.NewPortalConvertVaultRequest(
		metadata.PortalV4ConvertVaultRequestMeta,
		tokenID,
		convertingProof,
		paymentAddr,
	)
	if err != nil {
		return nil, "", err
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, portalConvertVaultMetadata, nil)
	if len(inputCoins) > 0 {
		return client.CreateRawTransactionWithInputCoins(txParam, inputCoins, coinIndices)
	}
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalConvertVaultTransaction(
	privateKey, tokenID, paymentAddr, convertingProof string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreatePortalConvertVaultTransaction(privateKey, tokenID, paymentAddr, convertingProof, inputCoins, coinIndices)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalUnShieldTransaction(
	privateKey, tokenID, remoteAddr string, unShieldingAmount uint64, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	addr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	pubKeyStr, txRandomStr, err := GenerateOTAFromPaymentAddress(addr)
	if err != nil {
		return nil, "", err
	}

	portalUnShieldingMetadata, err := metadata.NewPortalUnshieldRequest(
		metadata.PortalV4UnShieldingRequestMeta,
		pubKeyStr,
		txRandomStr,
		tokenID,
		remoteAddr,
		unShieldingAmount,
	)
	if err != nil {
		return nil, "", err
	}

	tokenParam := NewTxTokenParam(tokenID, 1, []string{common.BurningAddress2}, []uint64{unShieldingAmount}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, portalUnShieldingMetadata, nil)
	if len(inputCoins) > 0 {
		return client.CreateRawTransactionWithInputCoins(txParam, inputCoins, coinIndices)
	}
	return client.CreateRawTokenTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalUnShieldTransaction(
	privateKey, tokenID, remoteAddr string, unShieldingAmount uint64, inputCoins []coin.PlainCoin, coinIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreatePortalUnShieldTransaction(privateKey, tokenID, remoteAddr, unShieldingAmount, inputCoins, coinIndices)
	if err != nil {
		return "", err
	}
	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalRelayHeaderTransaction(
	privateKey, paymentAddr, header string, blockHeight uint64, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalRelayHeaderMetadata, err := metadata.NewRelayingHeader(
		metadata.RelayingBTCHeaderMeta,
		paymentAddr,
		header,
		blockHeight,
	)
	if err != nil {
		return nil, "", err
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, portalRelayHeaderMetadata, nil)
	if len(inputCoins) > 0 {
		return client.CreateRawTransactionWithInputCoins(txParam, inputCoins, coinIndices)
	}
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalRelayHeaderTransaction(
	privateKey, paymentAddr, header string, blockHeight uint64, inputCoins []coin.PlainCoin, coinIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreatePortalRelayHeaderTransaction(privateKey, paymentAddr, header, blockHeight, inputCoins, coinIndices)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}
