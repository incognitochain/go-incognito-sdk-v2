package mlsag

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
)


func (ml *Mlsag) SignConfidentialAsset(message []byte) (*MlsagSig, error) {
	if len(message) != common.HashSize {
		return nil, fmt.Errorf("cannot mlsag sign the message because its length is not 32, maybe it has not been hashed")
	}
	var message32byte [32]byte
	copy(message32byte[:], message)

	alpha, r := ml.createRandomChallenges()            // step 2 in paper
	c, err := ml.calculateCCA(message32byte, alpha, r) // step 3 and 4 in paper

	if err != nil {
		return nil, err
	}
	return &MlsagSig{
		c[0], ml.keyImages, r,
	}, nil
}

func calculateNextCCA(digest [common.HashSize]byte, r []*crypto.Scalar, c *crypto.Scalar, K []*crypto.Point, keyImages []*crypto.Point) (*crypto.Scalar, error) {
	if len(r) != len(K) || len(r) != len(keyImages) {
		return nil, fmt.Errorf("error in MLSAG: Calculating next C must have length of r be the same with length of ring R and same with length of keyImages")
	}
	var b []byte
	b = append(b, digest[:]...)

	// Below is the mathematics within the Monero paper:
	// If you are reviewing my code, please refer to paper
	// rG: r*G
	// cK: c*R
	// rG_cK: rG + cK
	//
	// HK: H_p(K_i)
	// rHK: r_i*H_p(K_i)
	// cKI: c*R~ (KI as keyImage)
	// rHK_cKI: rHK + cKI

	// Process columns before the last
	for i := 0; i < len(K)-2; i += 1 {
		rG := new(crypto.Point).ScalarMultBase(r[i])
		cK := new(crypto.Point).ScalarMult(K[i], c)
		rG_cK := new(crypto.Point).Add(rG, cK)

		HK := crypto.HashToPoint(K[i].ToBytesS())
		rHK := new(crypto.Point).ScalarMult(HK, r[i])
		cKI := new(crypto.Point).ScalarMult(keyImages[i], c)
		rHK_cKI := new(crypto.Point).Add(rHK, cKI)

		b = append(b, rG_cK.ToBytesS()...)
		b = append(b, rHK_cKI.ToBytesS()...)
	}

	// Process last column
	rG := new(crypto.Point).ScalarMult(
		crypto.PedCom.G[crypto.PedersenRandomnessIndex],
		r[len(K)-2],
	)
	cK := new(crypto.Point).ScalarMult(K[len(K)-2], c)
	rG_cK := new(crypto.Point).Add(rG, cK)
	b = append(b, rG_cK.ToBytesS()...)

	rG = new(crypto.Point).ScalarMult(
		crypto.PedCom.G[crypto.PedersenRandomnessIndex],
		r[len(K)-1],
	)
	cK = new(crypto.Point).ScalarMult(K[len(K)-1], c)
	rG_cK = new(crypto.Point).Add(rG, cK)
	b = append(b, rG_cK.ToBytesS()...)

	return crypto.HashToScalar(b), nil
}

func calculateFirstCCA(digest [common.HashSize]byte, alpha []*crypto.Scalar, K []*crypto.Point) (*crypto.Scalar, error) {
	if len(alpha) != len(K) {
		return nil, fmt.Errorf("error in MLSAG: Calculating first C must have length of alpha be the same with length of ring R")
	}
	var b []byte
	b = append(b, digest[:]...)

	// Process columns before the last
	for i := 0; i < len(K)-2; i += 1 {
		alphaG := new(crypto.Point).ScalarMultBase(alpha[i])

		H := crypto.HashToPoint(K[i].ToBytesS())
		alphaH := new(crypto.Point).ScalarMult(H, alpha[i])

		b = append(b, alphaG.ToBytesS()...)
		b = append(b, alphaH.ToBytesS()...)
	}

	// Process last column
	alphaG := new(crypto.Point).ScalarMult(
		// TODO : which g here ?
		crypto.PedCom.G[crypto.PedersenRandomnessIndex],
		alpha[len(K)-2],
	)
	b = append(b, alphaG.ToBytesS()...)
	alphaG = new(crypto.Point).ScalarMult(
		crypto.PedCom.G[crypto.PedersenRandomnessIndex],
		alpha[len(K)-1],
	)
	b = append(b, alphaG.ToBytesS()...)

	return crypto.HashToScalar(b), nil
}

func (ml *Mlsag) calculateCCA(message [common.HashSize]byte, alpha []*crypto.Scalar, r [][]*crypto.Scalar) ([]*crypto.Scalar, error) {
	m := len(ml.privateKeys)
	n := len(ml.R.keys)

	c := make([]*crypto.Scalar, n)
	firstC, err := calculateFirstCCA(
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
		nextC, err := calculateNextCCA(
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