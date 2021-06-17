package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func (client *IncClient) CreatePortalShieldTransaction(privateKey, tokenID, paymentAddr, shieldingProof string) ([]byte, string, error) {
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
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalShieldTransaction(privateKey, tokenID, paymentAddr, shieldingProof string) (string, error) {
	encodedTx, txHash, err := client.CreatePortalShieldTransaction(privateKey, tokenID, paymentAddr, shieldingProof)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalReplaceByFeeTransaction(privateKey, tokenID, batchID string, fee uint) ([]byte, string, error) {
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
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalReplaceByFeeTransaction(privateKey, tokenID, batchID string, fee uint) (string, error) {
	encodedTx, txHash, err := client.CreatePortalReplaceByFeeTransaction(privateKey, tokenID, batchID, fee)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalSubmitConfirmationTransaction(privateKey, tokenID, unshieldProof, batchID string) ([]byte, string, error) {
	portalSubmitConfirmationMetadata, err := metadata.NewPortalSubmitConfirmedTxRequest(
		metadata.PortalV4SubmitConfirmedTxMeta,
		unshieldProof,
		tokenID,
		batchID,
	)
	if err != nil {
		return nil, "", err
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, portalSubmitConfirmationMetadata, nil)
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalSubmitConfirmationTransaction(privateKey, tokenID, unshieldProof, batchID string) (string, error) {
	encodedTx, txHash, err := client.CreatePortalSubmitConfirmationTransaction(privateKey, tokenID, unshieldProof, batchID)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalConvertVaultTransaction(privateKey, tokenID, paymentAddr, convertingProof string) ([]byte, string, error) {
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
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalConvertVaultTransaction(privateKey, tokenID, paymentAddr, convertingProof string) (string, error) {
	encodedTx, txHash, err := client.CreatePortalSubmitConfirmationTransaction(privateKey, tokenID, paymentAddr, convertingProof)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalUnshieldTransaction(privateKey, tokenID, remoteAddr string, unshieldingAmount uint64) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	addr := senderWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	pubKeyStr, txRandomStr, err := GenerateOTAFromPaymentAddress(addr)
	if err != nil {
		return nil, "", err
	}

	portalUnshieldingMetadata, err := metadata.NewPortalUnshieldRequest(
		metadata.PortalV4UnshieldingRequestMeta,
		pubKeyStr,
		txRandomStr,
		tokenID,
		remoteAddr,
		unshieldingAmount,
	)
	if err != nil {
		return nil, "", err
	}

	tokenParam := NewTxTokenParam(tokenID, 1, []string{common.BurningAddress2}, []uint64{unshieldingAmount}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, tokenParam, portalUnshieldingMetadata, nil)
	return client.CreateRawTokenTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalUnshieldTransaction(privateKey, tokenID, remoteAddr string, unshieldingAmount uint64) (string, error) {
	encodedTx, txHash, err := client.CreatePortalUnshieldTransaction(privateKey, tokenID, remoteAddr, unshieldingAmount)
	if err != nil {
		return "", err
	}
	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (client *IncClient) CreatePortalRelayHeaderTransaction(privateKey, paymentAddr, header string, blockHeight uint64) ([]byte, string, error) {
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
	return client.CreateRawTransaction(txParam, 2)
}

func (client *IncClient) CreateAndSendPortalRelayHeaderTransaction(privateKey, paymentAddr, header string, blockHeight uint64) (string, error) {
	encodedTx, txHash, err := client.CreatePortalRelayHeaderTransaction(privateKey, paymentAddr, header, blockHeight)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}
