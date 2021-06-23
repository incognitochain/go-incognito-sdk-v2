package privacy

import (
	"errors"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/conversion"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/proof"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/proof/range_proof"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/hybridencryption"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/schnorr"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/aggregatedrange"
	v2 "github.com/incognitochain/go-incognito-sdk-v2/privacy/v2"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v2/bulletproofs"
)

const (
	RingSize = utils.RingSize
)

// PedCom represents the parameters for the Pedersen commitment scheme.
var PedCom = crypto.PedCom

// HybridCipherText represents a ciphertext in the hybrid encryption scheme.
type HybridCipherText = hybridencryption.HybridCipherText

// SchnorrSignature represents a Schnorr signature.
type SchnorrSignature = schnorr.SchnSignature

// SchnorrPublicKey is a public key used in the Schnorr signature scheme.
type SchnorrPublicKey = schnorr.SchnorrPublicKey

// SchnorrPrivateKey is a private key used in the Schnorr signature scheme.
type SchnorrPrivateKey = schnorr.SchnorrPrivateKey

// Proof represents a payment proof.
type Proof = proof.Proof

// ProofV1 represents a Proof of version 1 used in transactions v1.
type ProofV1 = zkp.ProofV1

// ProofV2 represents a Proof of version 2 used in transactions v2.
type ProofV2 = v2.ProofV2

// ProofForConversion represents a Proof used in conversion transactions. (e.g. to convert UTXOs v1 into UTXOs v2).
type ProofForConversion = conversion.ConversionProof

// PaymentWitness is a witness used to construct a ProofV1.
type PaymentWitness = zkp.PaymentWitness

// PaymentWitnessParam is used to initialize a PaymentWitness.
type PaymentWitnessParam = zkp.PaymentWitnessParam

// RangeProof represents a range proof.
type RangeProof = range_proof.RangeProof

// RangeProofV1 represents a RangeProof of version 1 used in transactions v1.
type RangeProofV1 = aggregatedrange.AggregatedRangeProof

// RangeProofV2 represents a RangeProof of version 2 used in transactions v2.
type RangeProofV2 = bulletproofs.RangeProof

// ArrayScalarToBytes parses a slice of scalars into a flatten slice of bytes.
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

// ProveV2 returns a ProofV2 based on the given input coins, output coins, shared secrets, etc.
// It is usually used in constructing a transaction of version 2.
func ProveV2(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, sharedSecrets []*crypto.Point, hasPrivacy bool, paymentInfo []*key.PaymentInfo) (*ProofV2, error) {
	return v2.Prove(inputCoins, outputCoins, sharedSecrets, hasPrivacy, paymentInfo)
}
