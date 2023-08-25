package tx_generic

import (
	"fmt"
	"math"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// GetTxMintData returns the minting data of a transaction w.r.t a given tokenID.
func GetTxMintData(tx metadata.Transaction, tokenID *common.Hash) (bool, coin.Coin, *common.Hash, error) {
	outputCoins, err := tx.GetReceiverData()
	if err != nil {
		return false, nil, nil, err
	}
	if len(outputCoins) != 1 {
		return false, nil, nil, fmt.Errorf("tx mint has more than one receiver")
	}
	if inputCoins := tx.GetProof().GetInputCoins(); len(inputCoins) > 0 {
		return false, nil, nil, fmt.Errorf("this is not a tx mint")
	}
	return true, outputCoins[0], tokenID, nil
}

// GetTxBurnData returns the burning data of a transaction.
func GetTxBurnData(tx metadata.Transaction) (bool, coin.Coin, *common.Hash, error) {
	outputCoins, err := tx.GetReceiverData()
	if err != nil {
		return false, nil, nil, err
	}

	for _, c := range outputCoins {
		if wallet.IsPublicKeyBurningAddress(c.GetPublicKey().ToBytesS()) {
			return true, c, &common.PRVCoinID, nil
		}
	}
	return false, nil, nil, nil
}

// CalculateSumOutputsWithFee returns a sum of given output coins' commitments plus the commitment of the given fee.
func CalculateSumOutputsWithFee(outputCoins []coin.Coin, fee uint64) *crypto.Point {
	sumOutputsWithFee := new(crypto.Point).Identity()
	for i := 0; i < len(outputCoins); i += 1 {
		sumOutputsWithFee.Add(sumOutputsWithFee, outputCoins[i].GetCommitment())
	}
	feeCommitment := new(crypto.Point).ScalarMult(
		crypto.PedCom.G[crypto.PedersenValueIndex],
		new(crypto.Scalar).FromUint64(fee),
	)
	sumOutputsWithFee.Add(sumOutputsWithFee, feeCommitment)
	return sumOutputsWithFee
}

// ValidateTxParams checks sanity of a TxPrivacyInitParams.
func ValidateTxParams(params *TxPrivacyInitParams) error {
	if len(params.InputCoins) > 255 {
		return fmt.Errorf("number of inputs (%v) is too large", len(params.InputCoins))
	}
	if len(params.PaymentInfo) > 254 {
		return fmt.Errorf("number of outputs (%v) is too large", len(params.PaymentInfo))
	}
	if params.TokenID == nil {
		// using default PRV
		params.TokenID = &common.Hash{}
		err := params.TokenID.SetBytes(common.PRVCoinID[:])
		if err != nil {
			return fmt.Errorf("cannot setbytes tokenID %v: %v", params.TokenID.String(), err)
		}
	}
	return nil
}

// SignNoPrivacy signs a message in a non-private manner using the Schnorr signature scheme.
func SignNoPrivacy(privateKey *key.PrivateKey, hashedMessage []byte) (signatureBytes []byte, sigPubKey []byte, err error) {
	sk := new(crypto.Scalar).FromBytesS(*privateKey)
	r := new(crypto.Scalar).FromUint64(0)
	sigKey := new(privacy.SchnorrPrivateKey)
	sigKey.Set(sk, r)
	signature, err := sigKey.Sign(hashedMessage)
	if err != nil {
		return nil, nil, err
	}

	signatureBytes = signature.Bytes()
	sigPubKey = sigKey.GetPublicKey().GetPublicKey().ToBytesS()
	return signatureBytes, sigPubKey, nil
}

type EstimateTxSizeParam struct {
	version                  int
	numInputCoins            int
	numPayments              int
	hasPrivacy               bool
	metadata                 metadata.Metadata
	privacyCustomTokenParams *TokenParam
	limitFee                 uint64
}

func NewEstimateTxSizeParam(version, numInputCoins, numPayments int,
	hasPrivacy bool, metadata metadata.Metadata,
	privacyCustomTokenParams *TokenParam,
	limitFee uint64) *EstimateTxSizeParam {
	estimateTxSizeParam := &EstimateTxSizeParam{
		version:                  version,
		numInputCoins:            numInputCoins,
		numPayments:              numPayments,
		hasPrivacy:               hasPrivacy,
		limitFee:                 limitFee,
		metadata:                 metadata,
		privacyCustomTokenParams: privacyCustomTokenParams,
	}
	return estimateTxSizeParam
}

func toB64Len(numOfBytes uint64) uint64 {
	l := (numOfBytes*4 + 2) / 3
	l = ((l + 3) / 4) * 4
	return l
}

func EstimateProofSizeV2(numIn, numOut uint64) uint64 {
	coinSizeBound := uint64(257) + (privacy.Ed25519KeySize+1)*7 + coin.TxRandomGroupSize + 1
	ipProofLRLen := uint64(math.Log2(float64(numOut))) + 1
	aggProofSizeBound := uint64(4) + 1 + privacy.Ed25519KeySize*uint64(7+numOut) + 1 + uint64(2*ipProofLRLen+3)*privacy.Ed25519KeySize
	// add 10 for rounding
	result := uint64(1) + (coinSizeBound+1)*uint64(numIn+numOut) + 2 + aggProofSizeBound + 10
	return toB64Len(result)
}

func EstimateTxSizeV2(estimateTxSizeParam *EstimateTxSizeParam) uint64 {
	jsonKeysSizeBound := uint64(20*10 + 2)
	sizeVersion := uint64(1)      // int8
	sizeType := uint64(5)         // string, max : 5
	sizeLockTime := uint64(8) * 3 // int64
	sizeFee := uint64(8) * 3      // uint64
	sizeInfo := toB64Len(uint64(512))

	numIn := uint64(estimateTxSizeParam.numInputCoins)
	numOut := uint64(estimateTxSizeParam.numPayments)

	sizeSigPubKey := uint64(numIn)*privacy.RingSize*9 + 2
	sizeSigPubKey = toB64Len(sizeSigPubKey)
	sizeSig := uint64(1) + numIn + (numIn+2)*privacy.RingSize
	sizeSig = sizeSig*33 + 3

	sizeProof := EstimateProofSizeV2(numIn, numOut)

	sizePubKeyLastByte := uint64(1) * 3
	sizeMetadata := uint64(0)
	if estimateTxSizeParam.metadata != nil {
		sizeMetadata += estimateTxSizeParam.metadata.CalculateSize()
	}

	sizeTx := jsonKeysSizeBound + sizeVersion + sizeType + sizeLockTime + sizeFee + sizeInfo + sizeSigPubKey + sizeSig + sizeProof + sizePubKeyLastByte + sizeMetadata
	if estimateTxSizeParam.privacyCustomTokenParams != nil {
		tokenKeysSizeBound := uint64(20*8 + 2)
		tokenSize := toB64Len(uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertyID)))
		tokenSize += uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertySymbol))
		tokenSize += uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertyName))
		tokenSize += 2
		numIn = uint64(len(estimateTxSizeParam.privacyCustomTokenParams.TokenInput))
		numOut = uint64(len(estimateTxSizeParam.privacyCustomTokenParams.Receiver))

		// shadow variable names
		sizeSigPubKey := uint64(numIn)*privacy.RingSize*9 + 2
		sizeSigPubKey = toB64Len(sizeSigPubKey)
		sizeSig := uint64(1) + numIn + (numIn+2)*privacy.RingSize
		sizeSig = sizeSig*33 + 3

		sizeProof := EstimateProofSizeV2(numIn, numOut)
		tokenSize += tokenKeysSizeBound + sizeSigPubKey + sizeSig + sizeProof
		sizeTx += tokenSize
	}
	return sizeTx
}

// // EstimateTxSize returns the estimated size of the tx in kilobyte
// func EstimateTxSize(estimateTxSizeParam *EstimateTxSizeParam) uint64 {
// 	if estimateTxSizeParam.version == 2 {
// 		return uint64(math.Ceil(float64(EstimateTxSizeV2(estimateTxSizeParam)) / 1024))
// 	}
// 	sizeVersion := uint64(1)  // int8
// 	sizeType := uint64(5)     // string, max : 5
// 	sizeLockTime := uint64(8) // int64
// 	sizeFee := uint64(8)      // uint64
// 	sizeInfo := uint64(512)

// 	sizeSigPubKey := uint64(common.SigPubKeySize)
// 	sizeSig := uint64(common.SigNoPrivacySize)
// 	if estimateTxSizeParam.hasPrivacy {
// 		sizeSig = uint64(common.SigPrivacySize)
// 	}

// 	sizeProof := uint64(0)
// 	if estimateTxSizeParam.numInputCoins != 0 || estimateTxSizeParam.numPayments != 0 {
// 		sizeProof = zku.EstimateProofSize(estimateTxSizeParam.numInputCoins, estimateTxSizeParam.numPayments, estimateTxSizeParam.hasPrivacy)
// 	} else if estimateTxSizeParam.limitFee > 0 {
// 		sizeProof = zku.EstimateProofSize(1, 1, estimateTxSizeParam.hasPrivacy)
// 	}

// 	sizePubKeyLastByte := uint64(1)

// 	sizeMetadata := uint64(0)
// 	if estimateTxSizeParam.metadata != nil {
// 		sizeMetadata += estimateTxSizeParam.metadata.CalculateSize()
// 	}

// 	sizeTx := sizeVersion + sizeType + sizeLockTime + sizeFee + sizeInfo + sizeSigPubKey + sizeSig + sizeProof + sizePubKeyLastByte + sizeMetadata

// 	// size of privacy custom token  data
// 	if estimateTxSizeParam.privacyCustomTokenParams != nil {
// 		customTokenDataSize := uint64(0)

// 		customTokenDataSize += uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertyID))
// 		customTokenDataSize += uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertySymbol))
// 		customTokenDataSize += uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertyName))

// 		customTokenDataSize += 8 // for amount
// 		customTokenDataSize += 4 // for TokenTxType

// 		customTokenDataSize += uint64(1) // int8 version
// 		customTokenDataSize += uint64(5) // string, max : 5 type
// 		customTokenDataSize += uint64(8) // int64 locktime
// 		customTokenDataSize += uint64(8) // uint64 fee

// 		customTokenDataSize += uint64(64) // info

// 		customTokenDataSize += uint64(common.SigPubKeySize)  // sig pubkey
// 		customTokenDataSize += uint64(common.SigPrivacySize) // sig

// 		// Proof
// 		customTokenDataSize += zku.EstimateProofSize(len(estimateTxSizeParam.privacyCustomTokenParams.TokenInput), len(estimateTxSizeParam.privacyCustomTokenParams.Receiver), true)

// 		customTokenDataSize += uint64(1) // PubKeyLastByte

// 		sizeTx += customTokenDataSize
// 	}

// 	return uint64(math.Ceil(float64(sizeTx) / 1024))
// }
