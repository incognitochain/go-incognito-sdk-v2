package schnorr

import (
	"errors"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
)

// SchnorrPublicKey represents a public key used in the Schnorr signature scheme.
//
// PK = G^SK + H^R.
type SchnorrPublicKey struct {
	publicKey *crypto.Point
	g, h      *crypto.Point
}

// GetPublicKey returns the real public key of a SchnorrPublicKey.
func (publicKey SchnorrPublicKey) GetPublicKey() *crypto.Point {
	return publicKey.publicKey
}

// Set sets v as the real public key of a SchnorrPublicKey.
func (publicKey *SchnorrPublicKey) Set(v *crypto.Point) {
	pubKey := v.GetKey()
	pedRandom := crypto.PedCom.G[crypto.PedersenRandomnessIndex].GetKey()
	pedPrivate := crypto.PedCom.G[crypto.PedersenPrivateKeyIndex].GetKey()

	publicKey.publicKey, _ = new(crypto.Point).SetKey(&pubKey)
	publicKey.g, _ = new(crypto.Point).SetKey(&pedPrivate)
	publicKey.h, _ = new(crypto.Point).SetKey(&pedRandom)
}

// SchnorrPrivateKey represents a private key used to sign messages in the Schnorr signature scheme.
type SchnorrPrivateKey struct {
	privateKey *crypto.Scalar
	randomness *crypto.Scalar
	publicKey  *SchnorrPublicKey
}

// GetPublicKey returns the corresponding public key of a SchnorrPrivateKey.
func (privateKey SchnorrPrivateKey) GetPublicKey() *SchnorrPublicKey {
	return privateKey.publicKey
}

// Set creats a new SchnorrPrivateKey.
func (privateKey *SchnorrPrivateKey) Set(sk *crypto.Scalar, r *crypto.Scalar) {
	pedRandom := crypto.PedCom.G[crypto.PedersenRandomnessIndex].GetKey()
	pedPrivate := crypto.PedCom.G[crypto.PedersenPrivateKeyIndex].GetKey()

	privateKey.privateKey = sk
	privateKey.randomness = r
	privateKey.publicKey = new(SchnorrPublicKey)
	privateKey.publicKey.g, _ = new(crypto.Point).SetKey(&pedPrivate)
	privateKey.publicKey.h, _ = new(crypto.Point).SetKey(&pedRandom)
	privateKey.publicKey.publicKey = new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], sk)
	privateKey.publicKey.publicKey.Add(privateKey.publicKey.publicKey, new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenRandomnessIndex], r))
}

// SchnSignature represents a Schnorr signature. The Schnorr signature is used to sign a transaction of version 1,
// or to sign metadata in a transaction of version 2.
type SchnSignature struct {
	e, z1, z2 *crypto.Scalar
}

// Bytes returns the byte-representation of a SchnSignature.
func (sig SchnSignature) Bytes() []byte {
	bytes := append(sig.e.ToBytesS(), sig.z1.ToBytesS()...)
	// Z2 is nil when has no privacy
	if sig.z2 != nil {
		bytes = append(bytes, sig.z2.ToBytesS()...)
	}
	return bytes
}

// SetBytes returns a SchnSignature given its byte-representation.
func (sig *SchnSignature) SetBytes(bytes []byte) error {
	if len(bytes) != 2*crypto.Ed25519KeySize && len(bytes) != 3*crypto.Ed25519KeySize {
		return utils.NewPrivacyErr(utils.InvalidInputToSetBytesErr, nil)
	}
	sig.e = new(crypto.Scalar).FromBytesS(bytes[0:crypto.Ed25519KeySize])
	sig.z1 = new(crypto.Scalar).FromBytesS(bytes[crypto.Ed25519KeySize : 2*crypto.Ed25519KeySize])
	if len(bytes) == 3*crypto.Ed25519KeySize {
		sig.z2 = new(crypto.Scalar).FromBytesS(bytes[2*crypto.Ed25519KeySize:])
	} else {
		sig.z2 = nil
	}

	return nil
}

// Sign returns a valid signature for the given data.
func (privateKey SchnorrPrivateKey) Sign(data []byte) (*SchnSignature, error) {
	if len(data) != common.HashSize {
		return nil, utils.NewPrivacyErr(utils.UnexpectedErr, errors.New("hash length must be 32 bytes"))
	}

	signature := new(SchnSignature)

	// has privacy
	if !privateKey.randomness.IsZero() {
		// generates random numbers s1, s2 in [0, Curve.Params().N - 1]

		s1 := crypto.RandomScalar()
		s2 := crypto.RandomScalar()

		// t = s1*G + s2*H
		t := new(crypto.Point).ScalarMult(privateKey.publicKey.g, s1)
		t.Add(t, new(crypto.Point).ScalarMult(privateKey.publicKey.h, s2))

		// E is the hash of elliptic point t and data need to be signed
		msg := append(t.ToBytesS(), data...)

		signature.e = crypto.HashToScalar(msg)

		signature.z1 = new(crypto.Scalar).Mul(privateKey.privateKey, signature.e)
		signature.z1 = new(crypto.Scalar).Sub(s1, signature.z1)

		signature.z2 = new(crypto.Scalar).Mul(privateKey.randomness, signature.e)
		signature.z2 = new(crypto.Scalar).Sub(s2, signature.z2)

		return signature, nil
	}

	// generates random numbers s, k2 in [0, Curve.Params().N - 1]
	s := crypto.RandomScalar()

	// t = s*G
	t := new(crypto.Point).ScalarMult(privateKey.publicKey.g, s)

	// E is the hash of elliptic point t and data need to be signed
	msg := append(t.ToBytesS(), data...)
	signature.e = crypto.HashToScalar(msg)

	// Z1 = s - e*sk
	signature.z1 = new(crypto.Scalar).Mul(privateKey.privateKey, signature.e)
	signature.z1 = new(crypto.Scalar).Sub(s, signature.z1)

	signature.z2 = nil

	return signature, nil
}
