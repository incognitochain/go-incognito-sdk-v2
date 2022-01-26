package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// DepositParams consists of parameters for creating a shielding transaction.
// A DepositParams is valid if at least one of the following conditions hold:
//	- Signature is not empty
//		- Receiver and DepositPubKey must not be empty
//	- Signature is empty
//		- If Receiver is empty, it will be generated from the sender's privateKey
//		- If DepositPrivateKey is empty, it will be derived from the DepositKeyIndex
//		- DepositPubKey is derived from DepositPrivateKey.
type DepositParams struct {
	// TokenID is the shielding asset ID.
	TokenID string

	// ShieldProof is a merkel proof for the shielding request.
	ShieldProof string

	// DepositPrivateKey is a base58-encoded deposit privateKey used to sign the request.
	// If set empty, it will be derived from the DepositKeyIndex.
	DepositPrivateKey string

	// DepositPubKey is a base58-encoded deposit publicKey. If Signature is not provided, DepositPubKey will be derived from the DepositPrivateKey.
	DepositPubKey string

	// DepositKeyIndex is the index of the OTDepositKey.
	DepositKeyIndex uint64

	// Receiver is a base58-encoded OTAReceiver. If set empty, it will be generated from the sender's privateKey.
	Receiver string

	// Signature is a valid signature signed by the owner of the shielding asset.
	// If Signature is not empty, DepositPubKey and Receiver must not be empty.
	Signature string
}

// IsValid checks if a DepositParams is valid.
func (dp DepositParams) IsValid() (bool, error) {
	var err error

	_, err = common.Hash{}.NewHashFromStr(dp.TokenID)
	if err != nil || dp.TokenID == "" {
		return false, fmt.Errorf("invalid tokenID %v", dp.TokenID)
	}

	if dp.Signature != "" {
		_, _, err = base58.Base58Check{}.Decode(dp.Signature)
		if err != nil {
			return false, fmt.Errorf("invalid signature")
		}
		if dp.DepositPubKey == "" || dp.Receiver == "" {
			return false, fmt.Errorf("must have both `DepositPubKey` and `Receiver`")
		}
	} else {
		if dp.DepositPrivateKey != "" {
			_, _, err = base58.Base58Check{}.Decode(dp.DepositPrivateKey)
			if err != nil {
				return false, fmt.Errorf("invalid DepositPrivateKey")
			}
		}
	}

	if dp.DepositPubKey != "" {
		_, _, err = base58.Base58Check{}.Decode(dp.DepositPubKey)
		if err != nil {
			return false, fmt.Errorf("invalid DepositPubKey")
		}
	}

	if dp.Receiver != "" {
		otaReceiver := new(coin.OTAReceiver)
		err = otaReceiver.FromString(dp.Receiver)
		if err != nil {
			return false, fmt.Errorf("invalid receiver: %v", err)
		}
	}

	return true, nil
}

// CreatePortalShieldTransaction creates a Portal V4 shielding transaction.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
//
// Deprecated: use CreatePortalShieldTransactionWithDepositKey instead.
func (client *IncClient) CreatePortalShieldTransaction(
	privateKey, tokenID, paymentAddr, shieldingProof string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalShieldingMetadata, err := metadata.NewPortalShieldingRequest(
		metadataCommon.PortalV4ShieldingRequestMeta,
		tokenID,
		paymentAddr,
		shieldingProof,
		"",
		nil,
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

// CreateAndSendPortalShieldTransaction creates a Portal V4 shielding transaction,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
//
// Deprecated: use CreateAndSendPortalShieldTransactionWithDepositKey instead.
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

// CreatePortalShieldTransactionWithDepositKey creates a Portal V4 shielding transaction using one-time depositing key.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePortalShieldTransactionWithDepositKey(
	privateKey string, depositParams DepositParams, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	w, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	_, err = depositParams.IsValid()
	if err != nil {
		return nil, "", err
	}

	tokenID, _ := common.Hash{}.NewHashFromStr(depositParams.TokenID)
	receiver, depositPubKey := depositParams.Receiver, depositParams.DepositPubKey
	var sig []byte
	if depositParams.Signature != "" {
		sig, _, _ = base58.Base58Check{}.Decode(depositParams.Signature)
	} else {
		if receiver == "" {
			otaReceivers, err := GenerateOTAReceivers([]common.Hash{*tokenID}, w.KeySet.PaymentAddress)
			if err != nil {
				return nil, "", err
			}
			receiver = otaReceivers[*tokenID].String()
		}
		otaReceiver := new(coin.OTAReceiver)
		_ = otaReceiver.FromString(receiver)

		var depositPrivateKey *crypto.Scalar
		if depositParams.DepositPrivateKey != "" {
			tmp, _, _ := base58.Base58Check{}.Decode(depositParams.DepositPrivateKey)
			depositPrivateKey = new(crypto.Scalar).FromBytesS(tmp)
		} else {
			depositKey, err := client.GenerateDepositKeyFromPrivateKey(privateKey, depositParams.TokenID, depositParams.DepositKeyIndex)
			if err != nil {
				return nil, "", fmt.Errorf("generate depositKey error: %v", err)
			}
			depositPrivateKey = new(crypto.Scalar).FromBytesS(depositKey.PrivateKey)
			depositPubKeyBytes := new(crypto.Point).ScalarMultBase(depositPrivateKey).ToBytesS()
			depositPubKey = base58.Base58Check{}.NewEncode(depositPubKeyBytes, 0)
		}
		schnorrPrivateKey := new(privacy.SchnorrPrivateKey)
		schnorrPrivateKey.Set(depositPrivateKey, crypto.RandomScalar())
		tmpSig, err := schnorrPrivateKey.Sign(common.HashB(otaReceiver.Bytes()))
		if err != nil {
			return nil, "", fmt.Errorf("sign metadata error: %v", err)
		}
		sig = tmpSig.Bytes()
	}

	portalShieldingMetadata, err := metadata.NewPortalShieldingRequest(
		metadataCommon.PortalV4ShieldingRequestMeta,
		depositParams.TokenID,
		receiver,
		depositParams.ShieldProof,
		depositPubKey,
		sig,
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

// CreateAndSendPortalShieldTransactionWithDepositKey creates a Portal V4 shielding transaction,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendPortalShieldTransactionWithDepositKey(
	privateKey string, depositParam DepositParams, inputCoins []coin.PlainCoin, coinIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreatePortalShieldTransactionWithDepositKey(privateKey,
		depositParam, inputCoins, coinIndices)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

// CreatePortalReplaceByFeeTransaction creates a Portal V4 replace-by-fee transaction.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePortalReplaceByFeeTransaction(
	privateKey, tokenID, batchID string, fee uint, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalRBFMetadata, err := metadata.NewPortalReplacementFeeRequest(
		metadataCommon.PortalV4FeeReplacementRequestMeta,
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

// CreateAndSendPortalReplaceByFeeTransaction creates a Portal V4 replace-by-fee transaction,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
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

// CreatePortalSubmitConfirmationTransaction creates a Portal V4 confirmation submission transaction.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePortalSubmitConfirmationTransaction(
	privateKey, tokenID, unShieldProof, batchID string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalSubmitConfirmationMetadata, err := metadata.NewPortalSubmitConfirmedTxRequest(
		metadataCommon.PortalV4SubmitConfirmedTxMeta,
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

// CreateAndSendPortalSubmitConfirmationTransaction creates a Portal V4 confirmation submission transaction,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
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

// CreatePortalConvertVaultTransaction creates a Portal V4 vault conversion transaction.
// This transaction SHOULD only be created only one time when migrating centralized bridge to portal v4.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreatePortalConvertVaultTransaction(
	privateKey, tokenID, paymentAddr, convertingProof string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	portalConvertVaultMetadata, err := metadata.NewPortalConvertVaultRequest(
		metadataCommon.PortalV4ConvertVaultRequestMeta,
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

// CreateAndSendPortalConvertVaultTransaction creates a Portal V4 vault conversion transaction,
// and submits it to the Incognito network. This transaction SHOULD only be created only one time
// when migrating centralized bridge to portal v4.
//
// It returns the transaction's hash, and an error (if any).
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

// CreatePortalUnShieldTransaction creates a Portal V4 un-shielding transaction.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
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
		metadataCommon.PortalV4UnshieldingRequestMeta,
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

// CreateAndSendPortalUnShieldTransaction creates a Portal V4 un-shielding transaction,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
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

// CreatePortalRelayHeaderTransaction creates block header-relaying transaction used in the Portal V4 protocol.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
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

// CreateAndSendPortalRelayHeaderTransaction creates block header-relaying transaction used in the Portal V4 protocol,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
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
