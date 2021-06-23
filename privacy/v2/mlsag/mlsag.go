package mlsag

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	C25519 "github.com/incognitochain/go-incognito-sdk-v2/crypto/curve25519"
)

var curveOrder = new(crypto.Scalar).SetKeyUnsafe(&C25519.L)

// Ring represents a ring of public keys used in the MLSAG signature scheme.
type Ring struct {
	keys [][]*crypto.Point
}

// GetKeys returns the keys of a Ring.
func (ring Ring) GetKeys() [][]*crypto.Point {
	return ring.keys
}

// NewRing creates a new Ring from the given list of public keys.
func NewRing(keys [][]*crypto.Point) *Ring {
	return &Ring{keys}
}

// ToBytes returns the byte-representation of a Ring.
func (ring Ring) ToBytes() ([]byte, error) {
	k := ring.keys
	if len(k) == 0 {
		return nil, fmt.Errorf("RingToBytes: Ring is empty")
	}

	// Make sure that the ring size is a rectangle row*column
	for i := 1; i < len(k); i += 1 {
		if len(k[i]) != len(k[0]) {
			return nil, fmt.Errorf("RingToBytes: Ring is not a proper rectangle row*column")
		}
	}
	n := len(k)
	m := len(k[0])
	if n > 255 || m > 255 {
		return nil, fmt.Errorf("RingToBytes: Ring size is too large")
	}
	b := make([]byte, 3)
	b[0] = SigPrefix
	b[1] = byte(n)
	b[2] = byte(m)

	for i := 0; i < n; i += 1 {
		for j := 0; j < m; j += 1 {
			b = append(b, k[i][j].ToBytesS()...)
		}
	}

	return b, nil
}

// FromBytes sets a byte-representation b to a Ring.
func (ring *Ring) FromBytes(b []byte) (*Ring, error) {
	if len(b) < 3 {
		return nil, fmt.Errorf("RingFromBytes: byte length is too short")
	}
	if b[0] != SigPrefix {
		return nil, fmt.Errorf("RingFromBytes: byte[0] is not SigPrefix")
	}
	n := int(b[1])
	m := int(b[2])

	if len(b) != crypto.Ed25519KeySize*n*m+3 {
		return nil, fmt.Errorf("RingFromBytes: byte length is not correct")
	}
	offset := 3
	key := make([][]*crypto.Point, 0)
	for i := 0; i < n; i += 1 {
		curRow := make([]*crypto.Point, m)
		for j := 0; j < m; j += 1 {
			currentByte := b[offset : offset+crypto.Ed25519KeySize]
			offset += crypto.Ed25519KeySize
			currentPoint, err := new(crypto.Point).FromBytesS(currentByte)
			if err != nil {
				return nil, fmt.Errorf("RingFromBytes: byte contains incorrect point")
			}
			curRow[j] = currentPoint
		}
		key = append(key, curRow)
	}
	ring.keys = key
	return ring, nil
}

// NewRandomRing creates a random ring with dimension: (numFake; len(privateKeys)) where we generate fake public keys inside.
func NewRandomRing(privateKeys []*crypto.Scalar, numFake, pi int) (K *Ring) {
	m := len(privateKeys)

	K = new(Ring)
	K.keys = make([][]*crypto.Point, numFake)
	for i := 0; i < numFake; i += 1 {
		if i != pi {
			K.keys[i] = createFakePublicKeyArray(m)
		} else {
			K.keys[pi] = make([]*crypto.Point, m)
			for j := 0; j < m; j += 1 {
				K.keys[i][j] = parsePublicKey(privateKeys[j], j == m-1)
			}
		}
	}
	return
}

// Mlsag represents the parameters of a MLSAG private key.
type Mlsag struct {
	R           *Ring
	pi          int
	keyImages   []*crypto.Point
	privateKeys []*crypto.Scalar
}

// NewMlsag creates a new Mlsag from the given private keys, a Ring and the index of the true private keys.
func NewMlsag(privateKeys []*crypto.Scalar, R *Ring, pi int) *Mlsag {
	return &Mlsag{
		R,
		pi,
		parseKeyImages(privateKeys),
		privateKeys,
	}
}

// Sign returns a signature for the given message.
func (ml *Mlsag) Sign(message []byte) (*Sig, error) {
	if len(message) != common.HashSize {
		return nil, fmt.Errorf("cannot mlsag sign the message because its length is not 32, maybe it has not been hashed")
	}
	message32byte := [32]byte{}
	copy(message32byte[:], message)

	alpha, r := ml.createRandomChallenges()          // step 2 in paper
	c, err := ml.calculateC(message32byte, alpha, r) // step 3 and 4 in paper

	if err != nil {
		return nil, err
	}
	return &Sig{
		c[0], ml.keyImages, r,
	}, nil
}

// parsePublicKey parses public key from private key.
func parsePublicKey(privateKey *crypto.Scalar, isLast bool) *crypto.Point {
	// isLast will commit to random base G
	if isLast {
		return new(crypto.Point).ScalarMult(
			crypto.PedCom.G[crypto.PedersenRandomnessIndex],
			privateKey,
		)
	}
	return new(crypto.Point).ScalarMultBase(privateKey)
}

func parseKeyImages(privateKeys []*crypto.Scalar) []*crypto.Point {
	m := len(privateKeys)

	result := make([]*crypto.Point, m)
	for i := 0; i < m; i += 1 {
		publicKey := parsePublicKey(privateKeys[i], i == m-1)
		hashPoint := crypto.HashToPoint(publicKey.ToBytesS())
		result[i] = new(crypto.Point).ScalarMult(hashPoint, privateKeys[i])
	}
	return result
}

func (ml *Mlsag) createRandomChallenges() (alpha []*crypto.Scalar, r [][]*crypto.Scalar) {
	m := len(ml.privateKeys)
	n := len(ml.R.keys)

	alpha = make([]*crypto.Scalar, m)
	for i := 0; i < m; i += 1 {
		alpha[i] = crypto.RandomScalar()
	}
	r = make([][]*crypto.Scalar, n)
	for i := 0; i < n; i += 1 {
		r[i] = make([]*crypto.Scalar, m)
		if i == ml.pi {
			continue
		}
		for j := 0; j < m; j += 1 {
			r[i][j] = crypto.RandomScalar()
		}
	}
	return
}

func calculateFirstC(digest [common.HashSize]byte, alpha []*crypto.Scalar, K []*crypto.Point) (*crypto.Scalar, error) {
	if len(alpha) != len(K) {
		return nil, fmt.Errorf("error in MLSAG: Calculating first C must have length of alpha be the same with length of ring R")
	}
	var b []byte
	b = append(b, digest[:]...)

	// Process columns before the last
	for i := 0; i < len(K)-1; i += 1 {
		alphaG := new(crypto.Point).ScalarMultBase(alpha[i])

		H := crypto.HashToPoint(K[i].ToBytesS())
		alphaH := new(crypto.Point).ScalarMult(H, alpha[i])

		b = append(b, alphaG.ToBytesS()...)
		b = append(b, alphaH.ToBytesS()...)
	}

	// Process last column
	alphaG := new(crypto.Point).ScalarMult(
		crypto.PedCom.G[crypto.PedersenRandomnessIndex],
		alpha[len(K)-1],
	)
	b = append(b, alphaG.ToBytesS()...)

	return crypto.HashToScalar(b), nil
}

func calculateNextC(digest [common.HashSize]byte, r []*crypto.Scalar, c *crypto.Scalar, K []*crypto.Point, keyImages []*crypto.Point) (*crypto.Scalar, error) {
	if len(r) != len(K) || len(r) != len(keyImages) {
		return nil, fmt.Errorf("error in MLSAG: Calculating next C must have length of r be the same with length of ring R and same with length of keyImages")
	}
	var b []byte
	b = append(b, digest[:]...)

	// Below is the mathematics within the Monero paper:
	// If you are reviewing my code, please refer to paper
	// rG: r*G
	// cK: c*R
	// sum1: rG + cK
	//
	// HK: H_p(K_i)
	// rHK: r_i*H_p(K_i)
	// cKI: c*R~ (KI as keyImage)
	// sum2: rHK + cKI

	// Process columns before the last
	for i := 0; i < len(K)-1; i += 1 {
		rG := new(crypto.Point).ScalarMultBase(r[i])
		cK := new(crypto.Point).ScalarMult(K[i], c)
		sum1 := new(crypto.Point).Add(rG, cK) // rG + cK

		HK := crypto.HashToPoint(K[i].ToBytesS())
		rHK := new(crypto.Point).ScalarMult(HK, r[i])
		cKI := new(crypto.Point).ScalarMult(keyImages[i], c)
		sum2 := new(crypto.Point).Add(rHK, cKI) // rHK + cKI

		b = append(b, sum1.ToBytesS()...)
		b = append(b, sum2.ToBytesS()...)
	}

	// Process last column
	rG := new(crypto.Point).ScalarMult(
		crypto.PedCom.G[crypto.PedersenRandomnessIndex],
		r[len(K)-1],
	)
	cK := new(crypto.Point).ScalarMult(K[len(K)-1], c)
	sum := new(crypto.Point).Add(rG, cK) // rG + cK
	b = append(b, sum.ToBytesS()...)

	return crypto.HashToScalar(b), nil
}

func (ml *Mlsag) calculateC(message [common.HashSize]byte, alpha []*crypto.Scalar, r [][]*crypto.Scalar) ([]*crypto.Scalar, error) {
	m := len(ml.privateKeys)
	n := len(ml.R.keys)

	c := make([]*crypto.Scalar, n)
	firstC, err := calculateFirstC(
		message,
		alpha,
		ml.R.keys[ml.pi],
	)
	if err != nil {
		return nil, err
	}

	var i int = (ml.pi + 1) % n
	c[i] = firstC
	for next := (i + 1) % n; i != ml.pi; {
		nextC, err := calculateNextC(
			message,
			r[i], c[i],
			(*ml.R).keys[i],
			ml.keyImages,
		)
		if err != nil {
			return nil, err
		}
		c[next] = nextC
		i = next
		next = (next + 1) % n
	}

	for i := 0; i < m; i += 1 {
		ck := new(crypto.Scalar).Mul(c[ml.pi], ml.privateKeys[i])
		r[ml.pi][i] = new(crypto.Scalar).Sub(alpha[i], ck)
	}

	return c, nil
}

func createFakePublicKeyArray(length int) []*crypto.Point {
	K := make([]*crypto.Point, length)
	for i := 0; i < length; i += 1 {
		K[i] = crypto.RandomPoint()
	}
	return K
}
