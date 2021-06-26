package privacy

import (
	"errors"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/conversion"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/proof"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/proof/range_proof"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/hybridencryption"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/schnorr"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp"
	bulletProofsV1 "github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/bulletproofs"
	v2 "github.com/incognitochain/go-incognito-sdk-v2/privacy/v2"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v2/bulletproofs"
)

type PrivacyError = utils.PrivacyError

var ErrCodeMessage = utils.ErrCodeMessage

// Public Constants
const (
	CStringBurnAddress    = "burningaddress"
	Ed25519KeySize        = crypto.Ed25519KeySize
	CStringBulletProof    = crypto.CStringBulletProof
	CommitmentRingSize    = utils.CommitmentRingSize
	CommitmentRingSizeExp = utils.CommitmentRingSizeExp

	PedersenSndIndex        = crypto.PedersenSndIndex
	PedersenValueIndex      = crypto.PedersenValueIndex
	PedersenShardIDIndex    = crypto.PedersenShardIDIndex
	PedersenPrivateKeyIndex = crypto.PedersenPrivateKeyIndex
	PedersenRandomnessIndex = crypto.PedersenRandomnessIndex

	RingSize          = utils.RingSize
	MaxTriesOta       = coin.MaxTriesOTA
	TxRandomGroupSize = coin.TxRandomGroupSize
)

var PedCom = crypto.PedCom

const (
	MaxSizeInfoCoin = coin.MaxSizeInfoCoin // byte
)

// Export as package privacy for other packages easily use it

type HybridCipherText = hybridencryption.HybridCipherText

type SchnSignature = schnorr.SchnSignature
type SchnorrPublicKey = schnorr.SchnorrPublicKey
type SchnorrPrivateKey = schnorr.SchnorrPrivateKey

type Proof = proof.Proof
type ProofV1 = zkp.ProofV1
type PaymentWitnessParam = zkp.PaymentWitnessParam
type PaymentWitness = zkp.PaymentWitness
type ProofV2 = v2.ProofV2
type ProofForConversion = conversion.ConversionProofVer1ToVer2
type RangeProof = range_proof.RangeProof

// RangeProofV1 represents a RangeProof of version 1 used in transactions v1.
type RangeProofV1 = bulletProofsV1.RangeProof

// RangeProofV2 represents a RangeProof of version 2 used in transactions v2.
type RangeProofV2 = bulletproofs.RangeProof

func NewProofWithVersion(version int8) Proof {
	var result Proof
	if version == 1 {
		result = &zkp.ProofV1{}
	} else {
		result = &v2.ProofV2{}
	}
	return result
}

func ArrayScalarToBytes(arr *[]*crypto.Scalar) ([]byte, error) {
	scalarArr := *arr

	n := len(scalarArr)
	if n > 255 {
		return nil, errors.New("ArrayScalarToBytes: length of scalar array is too big")
	}
	b := make([]byte, 1)
	b[0] = byte(n)

	for _, sc := range scalarArr {
		b = append(b, sc.ToBytesS()...)
	}
	return b, nil
}

func ArrayScalarFromBytes(b []byte) (*[]*crypto.Scalar, error) {
	if len(b) == 0 {
		return nil, errors.New("ArrayScalarFromBytes error: length of byte is 0")
	}
	n := int(b[0])
	if n*Ed25519KeySize+1 != len(b) {
		return nil, errors.New("ArrayScalarFromBytes error: length of byte is not correct")
	}
	scalarArr := make([]*crypto.Scalar, n)
	offset := 1
	for i := 0; i < n; i += 1 {
		curByte := b[offset : offset+Ed25519KeySize]
		scalarArr[i] = new(crypto.Scalar).FromBytesS(curByte)
		offset += Ed25519KeySize
	}
	return &scalarArr, nil
}

func ProveV2(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, sharedSecrets []*crypto.Point, hasPrivacy bool, paymentInfo []*key.PaymentInfo) (*ProofV2, error) {
	return v2.Prove(inputCoins, outputCoins, sharedSecrets, hasPrivacy, paymentInfo)
}
