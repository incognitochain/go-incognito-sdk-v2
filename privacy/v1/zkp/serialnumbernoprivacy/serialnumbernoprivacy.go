package serialnumbernoprivacy

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/utils"
)

// SerialNumberNoPrivacyStatement represents the statement of an SNNoPrivacyProof.
type SerialNumberNoPrivacyStatement struct {
	sn     *crypto.Point
	pubKey *crypto.Point
	snd    *crypto.Scalar
}

// SNNoPrivacyWitness represents the witness of an SNNoPrivacyProof.
type SNNoPrivacyWitness struct {
	stmt SerialNumberNoPrivacyStatement
	seed *crypto.Scalar
}

// SNNoPrivacyProof represents a zero-knowledge proof for proving the the serial number is correctly generated from the
// secret key and the SND with the following formula: SN = (sk + snd)^-1 * G[0].
//
// It is only used in non-private transactions of version 1, and conversion transactions.
type SNNoPrivacyProof struct {
	// general info
	stmt SerialNumberNoPrivacyStatement

	tSeed   *crypto.Point
	tOutput *crypto.Point

	zSeed *crypto.Scalar
}

// GetPubKey returns the public key in the statement of an SNNoPrivacyProof.
func (proof SNNoPrivacyProof) GetPubKey() *crypto.Point {
	return proof.stmt.pubKey
}

// GetSN returns the serial number in the statement of an SNNoPrivacyProof.
func (proof SNNoPrivacyProof) GetSN() *crypto.Point {
	return proof.stmt.sn
}

// GetSND returns the serial number derivator in the statement of an SNNoPrivacyProof.
func (proof SNNoPrivacyProof) GetSND() *crypto.Scalar {
	return proof.stmt.snd
}

// Init creates an empty ValidateSanity.
func (proof *SNNoPrivacyProof) Init() *SNNoPrivacyProof {
	proof.stmt.sn = new(crypto.Point)
	proof.stmt.pubKey = new(crypto.Point)
	proof.stmt.snd = new(crypto.Scalar)

	proof.tSeed = new(crypto.Point)
	proof.tOutput = new(crypto.Point)

	proof.zSeed = new(crypto.Scalar)

	return proof
}

// Set sets data to an SNNoPrivacyWitness.
func (wit *SNNoPrivacyWitness) Set(
	output *crypto.Point,
	vKey *crypto.Point,
	input *crypto.Scalar,
	seed *crypto.Scalar) {

	wit.stmt.sn = output
	wit.stmt.pubKey = vKey
	wit.stmt.snd = input

	wit.seed = seed
}

// Set sets data to an SNNoPrivacyProof.
func (proof *SNNoPrivacyProof) Set(
	output *crypto.Point,
	vKey *crypto.Point,
	input *crypto.Scalar,
	tSeed *crypto.Point,
	tOutput *crypto.Point,
	zSeed *crypto.Scalar) {

	proof.stmt.sn = output
	proof.stmt.pubKey = vKey
	proof.stmt.snd = input

	proof.tSeed = tSeed
	proof.tOutput = tOutput

	proof.zSeed = zSeed
}

func (proof SNNoPrivacyProof) Bytes() []byte {
	// if proof is nil, return an empty array
	if proof.isNil() {
		return []byte{}
	}

	var bytes []byte
	bytes = append(bytes, proof.stmt.sn.ToBytesS()...)
	bytes = append(bytes, proof.stmt.pubKey.ToBytesS()...)
	bytes = append(bytes, proof.stmt.snd.ToBytesS()...)

	bytes = append(bytes, proof.tSeed.ToBytesS()...)
	bytes = append(bytes, proof.tOutput.ToBytesS()...)

	bytes = append(bytes, proof.zSeed.ToBytesS()...)

	return bytes
}

func (proof *SNNoPrivacyProof) SetBytes(bytes []byte) error {
	if len(bytes) < crypto.Ed25519KeySize*6 {
		return fmt.Errorf("not enough bytes to unmarshal Serial Number No Privacy Proof")
	}

	offset := 0
	var err error
	proof.stmt.sn, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	proof.stmt.pubKey, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	proof.stmt.snd.FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	proof.tSeed, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	proof.tOutput, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	proof.zSeed = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])

	return nil
}

func (wit SNNoPrivacyWitness) Prove(mess []byte) (*SNNoPrivacyProof, error) {
	// randomness
	eSK := crypto.RandomScalar()

	// calculate tSeed = g_SK^eSK
	tSK := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], eSK)

	// calculate tOutput = sn^eSK
	tE := new(crypto.Point).ScalarMult(wit.stmt.sn, eSK)

	x := new(crypto.Scalar)
	if mess == nil {
		// calculate x = hash(tSeed || tInput || tSND2 || tOutput)
		x = utils.GenerateChallenge([][]byte{wit.stmt.sn.ToBytesS(), wit.stmt.pubKey.ToBytesS(), tSK.ToBytesS(), tE.ToBytesS()})
	} else {
		x.FromBytesS(mess)
	}

	// Calculate zSeed = SK * x + eSK
	zSK := new(crypto.Scalar).Mul(wit.seed, x)
	zSK.Add(zSK, eSK)

	proof := new(SNNoPrivacyProof).Init()
	proof.Set(wit.stmt.sn, wit.stmt.pubKey, wit.stmt.snd, tSK, tE, zSK)
	return proof, nil
}

func (proof SNNoPrivacyProof) isNil() bool {
	if proof.stmt.sn == nil {
		return true
	}
	if proof.stmt.pubKey == nil {
		return true
	}
	if proof.stmt.snd == nil {
		return true
	}
	if proof.tSeed == nil {
		return true
	}
	if proof.tOutput == nil {
		return true
	}
	if proof.zSeed == nil {
		return true
	}
	return false
}
