package bulletproofs

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"math"
)

type InnerProductWitness struct {
	a []*crypto.Scalar
	b []*crypto.Scalar
	p *crypto.Point
}

type InnerProductProof struct {
	l []*crypto.Point
	r []*crypto.Point
	a *crypto.Scalar
	b *crypto.Scalar
	p *crypto.Point
}

func (proof *InnerProductProof) Init() *InnerProductProof {
	proof.l = []*crypto.Point{}
	proof.r = []*crypto.Point{}
	proof.a = new(crypto.Scalar)
	proof.b = new(crypto.Scalar)
	proof.p = new(crypto.Point).Identity()

	return proof
}

func (proof InnerProductProof) ValidateSanity() bool {
	if len(proof.l) != len(proof.r) {
		return false
	}
	for i := 0; i < len(proof.l); i++ {
		if !proof.l[i].PointValid() || !proof.r[i].PointValid() {
			return false
		}
	}
	if !proof.a.ScalarValid() || !proof.b.ScalarValid() {
		return false
	}
	return proof.p.PointValid()
}

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
func (proof InnerProductProof) Verify(GParam []*crypto.Point, HParam []*crypto.Point, uParam *crypto.Point, hashCache []byte) bool {
	//var aggParam = newBulletproofParams(1)
	p := new(crypto.Point)
	p.Set(proof.p)

	n := len(GParam)
	G := make([]*crypto.Point, n)
	H := make([]*crypto.Point, n)
	for i := range G {
		G[i] = new(crypto.Point).Set(GParam[i])
		H[i] = new(crypto.Point).Set(HParam[i])
	}

	for i := range proof.l {
		nPrime := n / 2
		x := generateChallenge(hashCache, []*crypto.Point{proof.l[i], proof.r[i]})
		hashCache = new(crypto.Scalar).Set(x).ToBytesS()
		xInverse := new(crypto.Scalar).Invert(x)
		xSquare := new(crypto.Scalar).Mul(x, x)
		xSquareInverse := new(crypto.Scalar).Mul(xInverse, xInverse)

		// calculate GPrime, HPrime, PPrime for the next loop
		GPrime := make([]*crypto.Point, nPrime)
		HPrime := make([]*crypto.Point, nPrime)

		for j := 0; j < len(GPrime); j++ {
			GPrime[j] = new(crypto.Point).AddPedersen(xInverse, G[j], x, G[j+nPrime])
			HPrime[j] = new(crypto.Point).AddPedersen(x, H[j], xInverse, H[j+nPrime])
		}
		// calculate x^2 * l + P + xInverse^2 * r
		PPrime := new(crypto.Point).AddPedersen(xSquare, proof.l[i], xSquareInverse, proof.r[i])
		PPrime.Add(PPrime, p)

		p = PPrime
		G = GPrime
		H = HPrime
		n = nPrime
	}

	c := new(crypto.Scalar).Mul(proof.a, proof.b)
	rightPoint := new(crypto.Point).AddPedersen(proof.a, G[0], proof.b, H[0])
	rightPoint.Add(rightPoint, new(crypto.Point).ScalarMult(uParam, c))
	res := crypto.IsPointEqual(rightPoint, p)
	if !res {
		fmt.Println("Inner product argument failed:")
		fmt.Printf("p: %v\n", p)
		fmt.Printf("RightPoint: %v\n", rightPoint)
	}

	return res
}

func (proof InnerProductProof) VerifyFaster(GParam []*crypto.Point, HParam []*crypto.Point, uParam *crypto.Point, hashCache []byte) bool {
	//var aggParam = newBulletproofParams(1)
	p := new(crypto.Point)
	p.Set(proof.p)
	n := len(GParam)
	G := make([]*crypto.Point, n)
	H := make([]*crypto.Point, n)
	s := make([]*crypto.Scalar, n)
	sInverse := make([]*crypto.Scalar, n)

	for i := range G {
		G[i] = new(crypto.Point).Set(GParam[i])
		H[i] = new(crypto.Point).Set(HParam[i])
		s[i] = new(crypto.Scalar).FromUint64(1)
		sInverse[i] = new(crypto.Scalar).FromUint64(1)
	}
	logN := int(math.Log2(float64(n)))
	xList := make([]*crypto.Scalar, logN)
	xInverseList := make([]*crypto.Scalar, logN)
	xSquareList := make([]*crypto.Scalar, logN)
	xInverseSquareList := make([]*crypto.Scalar, logN)

	//a*s ; b*s^-1

	for i := range proof.l {
		// calculate challenge x = hash(hash(G || H || u || p) || x || l || r)
		xList[i] = generateChallenge(hashCache, []*crypto.Point{proof.l[i], proof.r[i]})
		hashCache = new(crypto.Scalar).Set(xList[i]).ToBytesS()

		xInverseList[i] = new(crypto.Scalar).Invert(xList[i])
		xSquareList[i] = new(crypto.Scalar).Mul(xList[i], xList[i])
		xInverseSquareList[i] = new(crypto.Scalar).Mul(xInverseList[i], xInverseList[i])

		//Update s, s^-1
		for j := 0; j < n; j++ {
			if j&int(math.Pow(2, float64(logN-i-1))) != 0 {
				s[j] = new(crypto.Scalar).Mul(s[j], xList[i])
				sInverse[j] = new(crypto.Scalar).Mul(sInverse[j], xInverseList[i])
			} else {
				s[j] = new(crypto.Scalar).Mul(s[j], xInverseList[i])
				sInverse[j] = new(crypto.Scalar).Mul(sInverse[j], xList[i])
			}
		}
	}

	// Compute (g^s)^a (h^-s)^b u^(ab) = p l^(x^2) r^(-x^2)
	c := new(crypto.Scalar).Mul(proof.a, proof.b)
	rightHSPart1 := new(crypto.Point).MultiScalarMult(s, G)
	rightHSPart1.ScalarMult(rightHSPart1, proof.a)
	rightHSPart2 := new(crypto.Point).MultiScalarMult(sInverse, H)
	rightHSPart2.ScalarMult(rightHSPart2, proof.b)

	rightHS := new(crypto.Point).Add(rightHSPart1, rightHSPart2)
	rightHS.Add(rightHS, new(crypto.Point).ScalarMult(uParam, c))

	leftHSPart1 := new(crypto.Point).MultiScalarMult(xSquareList, proof.l)
	leftHSPart2 := new(crypto.Point).MultiScalarMult(xInverseSquareList, proof.r)

	leftHS := new(crypto.Point).Add(leftHSPart1, leftHSPart2)
	leftHS.Add(leftHS, proof.p)

	res := crypto.IsPointEqual(rightHS, leftHS)
	if !res {
		fmt.Println("Inner product argument failed:")
		fmt.Printf("LHS: %v\n", leftHS)
		fmt.Printf("RHS: %v\n", rightHS)
	}

	return res
}
