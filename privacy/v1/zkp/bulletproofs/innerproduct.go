package bulletproofs

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
)

// InnerProductWitness represents a witness for an inner-product proof, described in the Bulletproofs paper.
type InnerProductWitness struct {
	a []*crypto.Scalar
	b []*crypto.Scalar
	p *crypto.Point
}

// InnerProductProof represents an inner-product proof. It is used as a sub-proof for a RangeProof.
type InnerProductProof struct {
	l []*crypto.Point
	r []*crypto.Point
	a *crypto.Scalar
	b *crypto.Scalar
	p *crypto.Point
}

// Init creates an empty InnerProductProof.
func (proof *InnerProductProof) Init() *InnerProductProof {
	proof.l = []*crypto.Point{}
	proof.r = []*crypto.Point{}
	proof.a = new(crypto.Scalar)
	proof.b = new(crypto.Scalar)
	proof.p = new(crypto.Point).Identity()

	return proof
}

// Bytes returns the byte-representation of an InnerProductProof.
func (proof InnerProductProof) Bytes() []byte {
	var res []byte

	res = append(res, byte(len(proof.l)))
	for _, l := range proof.l {
		res = append(res, l.ToBytesS()...)
	}

	for _, r := range proof.r {
		res = append(res, r.ToBytesS()...)
	}

	res = append(res, proof.a.ToBytesS()...)
	res = append(res, proof.b.ToBytesS()...)
	res = append(res, proof.p.ToBytesS()...)

	return res
}

// SetBytes sets byte-representation data to an InnerProductProof.
func (proof *InnerProductProof) SetBytes(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	lenLArray := int(bytes[0])
	offset := 1
	var err error

	proof.l = make([]*crypto.Point, lenLArray)
	for i := 0; i < lenLArray; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("unmarshalling failed")
		}
		proof.l[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset += crypto.Ed25519KeySize
	}

	proof.r = make([]*crypto.Point, lenLArray)
	for i := 0; i < lenLArray; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("unmarshalling failed")
		}
		proof.r[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset += crypto.Ed25519KeySize
	}

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling failed")
	}
	proof.a = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling failed")
	}
	proof.b = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling failed")
	}
	proof.p, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}

	return nil
}

// Prove returns an InnerProductProof for an InnerProductWitness.
func (wit InnerProductWitness) Prove(GParam []*crypto.Point, HParam []*crypto.Point, uParam *crypto.Point, hashCache []byte) (*InnerProductProof, error) {
	if len(wit.a) != len(wit.b) {
		return nil, fmt.Errorf("invalid inputs")
	}

	N := len(wit.a)

	a := make([]*crypto.Scalar, N)
	b := make([]*crypto.Scalar, N)

	for i := range wit.a {
		a[i] = new(crypto.Scalar).Set(wit.a[i])
		b[i] = new(crypto.Scalar).Set(wit.b[i])
	}

	p := new(crypto.Point).Set(wit.p)
	G := make([]*crypto.Point, N)
	H := make([]*crypto.Point, N)
	for i := range G {
		G[i] = new(crypto.Point).Set(GParam[i])
		H[i] = new(crypto.Point).Set(HParam[i])
	}

	proof := new(InnerProductProof)
	proof.l = make([]*crypto.Point, 0)
	proof.r = make([]*crypto.Point, 0)
	proof.p = new(crypto.Point).Set(wit.p)

	for N > 1 {
		nPrime := N / 2

		cL, err := innerProduct(a[:nPrime], b[nPrime:])
		if err != nil {
			return nil, err
		}
		cR, err := innerProduct(a[nPrime:], b[:nPrime])
		if err != nil {
			return nil, err
		}

		L, err := encodeVectors(a[:nPrime], b[nPrime:], G[nPrime:], H[:nPrime])
		if err != nil {
			return nil, err
		}
		L.Add(L, new(crypto.Point).ScalarMult(uParam, cL))
		proof.l = append(proof.l, L)

		R, err := encodeVectors(a[nPrime:], b[:nPrime], G[:nPrime], H[nPrime:])
		if err != nil {
			return nil, err
		}
		R.Add(R, new(crypto.Point).ScalarMult(uParam, cR))
		proof.r = append(proof.r, R)

		x := generateChallenge(hashCache, []*crypto.Point{L, R})
		hashCache = new(crypto.Scalar).Set(x).ToBytesS()

		xInverse := new(crypto.Scalar).Invert(x)
		xSquare := new(crypto.Scalar).Mul(x, x)
		xSquareInverse := new(crypto.Scalar).Mul(xInverse, xInverse)

		// calculate GPrime, HPrime, PPrime for the next loop
		GPrime := make([]*crypto.Point, nPrime)
		HPrime := make([]*crypto.Point, nPrime)

		for i := range GPrime {
			GPrime[i] = new(crypto.Point).AddPedersen(xInverse, G[i], x, G[i+nPrime])
			HPrime[i] = new(crypto.Point).AddPedersen(x, H[i], xInverse, H[i+nPrime])
		}

		// x^2 * l + P + xInverse^2 * r
		PPrime := new(crypto.Point).AddPedersen(xSquare, L, xSquareInverse, R)
		PPrime.Add(PPrime, p)

		// calculate aPrime, bPrime
		aPrime := make([]*crypto.Scalar, nPrime)
		bPrime := make([]*crypto.Scalar, nPrime)

		for i := range aPrime {
			aPrime[i] = new(crypto.Scalar).Mul(a[i], x)
			aPrime[i] = new(crypto.Scalar).MulAdd(a[i+nPrime], xInverse, aPrime[i])

			bPrime[i] = new(crypto.Scalar).Mul(b[i], xInverse)
			bPrime[i] = new(crypto.Scalar).MulAdd(b[i+nPrime], x, bPrime[i])
		}

		a = aPrime
		b = bPrime
		p.Set(PPrime)
		G = GPrime
		H = HPrime
		N = nPrime
	}

	proof.a = new(crypto.Scalar).Set(a[0])
	proof.b = new(crypto.Scalar).Set(b[0])

	return proof, nil
}
