package aggregatedrange

import (
	"errors"
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
		proof.l[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset += crypto.Ed25519KeySize
	}

	proof.r = make([]*crypto.Point, lenLArray)
	for i := 0; i < lenLArray; i++ {
		proof.r[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset += crypto.Ed25519KeySize
	}

	proof.a = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	proof.b = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	proof.p, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}

	return nil
}

func (wit InnerProductWitness) Prove(aggParam *bulletproofParams) (*InnerProductProof, error) {
	if len(wit.a) != len(wit.b) {
		return nil, errors.New("invalid inputs")
	}

	n := len(wit.a)

	a := make([]*crypto.Scalar, n)
	b := make([]*crypto.Scalar, n)

	for i := range wit.a {
		a[i] = new(crypto.Scalar).Set(wit.a[i])
		b[i] = new(crypto.Scalar).Set(wit.b[i])
	}

	p := new(crypto.Point).Set(wit.p)
	G := make([]*crypto.Point, n)
	H := make([]*crypto.Point, n)
	for i := range G {
		G[i] = new(crypto.Point).Set(aggParam.g[i])
		H[i] = new(crypto.Point).Set(aggParam.h[i])
	}

	proof := new(InnerProductProof)
	proof.l = make([]*crypto.Point, 0)
	proof.r = make([]*crypto.Point, 0)
	proof.p = new(crypto.Point).Set(wit.p)

	for n > 1 {
		nPrime := n / 2

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
		L.Add(L, new(crypto.Point).ScalarMult(aggParam.u, cL))
		proof.l = append(proof.l, L)

		R, err := encodeVectors(a[nPrime:], b[:nPrime], G[:nPrime], H[nPrime:])
		if err != nil {
			return nil, err
		}
		R.Add(R, new(crypto.Point).ScalarMult(aggParam.u, cR))
		proof.r = append(proof.r, R)

		// calculate challenge x = hash(G || H || u || x || l || r)
		x := generateChallenge([][]byte{aggParam.cs, p.ToBytesS(), L.ToBytesS(), R.ToBytesS()})
		//x := generateChallengeOld(aggParam, [][]byte{p.ToBytesS(), L.ToBytesS(), R.ToBytesS()})
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
		n = nPrime
	}

	proof.a = new(crypto.Scalar).Set(a[0])
	proof.b = new(crypto.Scalar).Set(b[0])

	return proof, nil
}

func (proof InnerProductProof) Verify(aggParam *bulletproofParams) bool {
	//var aggParam = newBulletproofParams(1)
	p := new(crypto.Point)
	p.Set(proof.p)
	n := len(aggParam.g)
	G := make([]*crypto.Point, n)
	H := make([]*crypto.Point, n)
	s := make([]*crypto.Scalar, n)
	sInverse := make([]*crypto.Scalar, n)

	for i := range G {
		G[i] = new(crypto.Point).Set(aggParam.g[i])
		H[i] = new(crypto.Point).Set(aggParam.h[i])
		s[i] = new(crypto.Scalar).FromUint64(1)
		sInverse[i] = new(crypto.Scalar).FromUint64(1)
	}
	logN := int(math.Log2(float64(n)))
	xList := make([]*crypto.Scalar, logN)
	xInverseList := make([]*crypto.Scalar, logN)
	xSquareList := make([]*crypto.Scalar, logN)
	xInverseSquare_List := make([]*crypto.Scalar, logN)

	//a*s ; b*s^-1
	for i := range proof.l {
		// calculate challenge x = hash(hash(G || H || u || p) || x || l || r)
		xList[i] = generateChallenge([][]byte{aggParam.cs, p.ToBytesS(), proof.l[i].ToBytesS(), proof.r[i].ToBytesS()})
		xInverseList[i] = new(crypto.Scalar).Invert(xList[i])
		xSquareList[i] = new(crypto.Scalar).Mul(xList[i], xList[i])
		xInverseSquare_List[i] = new(crypto.Scalar).Mul(xInverseList[i], xInverseList[i])

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
		PPrime := new(crypto.Point).AddPedersen(xSquareList[i], proof.l[i], xInverseSquare_List[i], proof.r[i])
		PPrime.Add(PPrime, p)
		p = PPrime
	}

	// Compute (g^s)^a (h^-s)^b u^(ab) = p l^(x^2) r^(-x^2)
	c := new(crypto.Scalar).Mul(proof.a, proof.b)
	rightHSPart1 := new(crypto.Point).MultiScalarMult(s, G)
	rightHSPart1.ScalarMult(rightHSPart1, proof.a)
	rightHSPart2 := new(crypto.Point).MultiScalarMult(sInverse, H)
	rightHSPart2.ScalarMult(rightHSPart2, proof.b)
	rightHS := new(crypto.Point).Add(rightHSPart1, rightHSPart2)
	rightHS.Add(rightHS, new(crypto.Point).ScalarMult(aggParam.u, c))

	leftHSPart1 := new(crypto.Point).MultiScalarMult(xSquareList, proof.l)
	leftHSPart2 := new(crypto.Point).MultiScalarMult(xInverseSquare_List, proof.r)
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

func VerifyBatchingInnerProductProofs(proofs []*InnerProductProof, csList [][]byte) bool {
	batchSize := len(proofs)
	// Generate list of random value
	sum_abAlpha := new(crypto.Scalar).FromUint64(0)
	pList := make([]*crypto.Point, 0)
	alphaList := make([]*crypto.Scalar, 0)
	LList := make([]*crypto.Point, 0)
	nXSquareList := make([]*crypto.Scalar, 0)
	RList := make([]*crypto.Point, 0)
	nXInverseSquareList := make([]*crypto.Scalar, 0)

	maxN := 0
	asAlphaList := make([]*crypto.Scalar, len(AggParam.g))
	bsInverseAlphaList := make([]*crypto.Scalar, len(AggParam.g))
	for k := 0; k < len(AggParam.g); k++ {
		asAlphaList[k] = new(crypto.Scalar).FromUint64(0)
		bsInverseAlphaList[k] = new(crypto.Scalar).FromUint64(0)
	}
	for i := 0; i < batchSize; i++ {
		alpha := crypto.RandomScalar()
		abAlpha := new(crypto.Scalar).Mul(proofs[i].a, proofs[i].b)
		abAlpha.Mul(abAlpha, alpha)
		sum_abAlpha.Add(sum_abAlpha, abAlpha)

		//prod_PAlpha.Add(prod_PAlpha, new(crypto.Point).ScalarMult(proofs[i].p,alpha))
		pList = append(pList, proofs[i].p)
		alphaList = append(alphaList, alpha)

		n := int(math.Pow(2, float64(len(proofs[i].l))))
		if maxN < n {
			maxN = n
		}
		logN := int(math.Log2(float64(n)))
		s := make([]*crypto.Scalar, n)
		sInverse := make([]*crypto.Scalar, n)
		xList := make([]*crypto.Scalar, logN)
		xInverseList := make([]*crypto.Scalar, logN)
		xSquareList := make([]*crypto.Scalar, logN)
		xSquareAlphaList := make([]*crypto.Scalar, logN)
		xInverseSquareList := make([]*crypto.Scalar, logN)
		xInverseSquareAlphaList := make([]*crypto.Scalar, logN)

		for k := 0; k < n; k++ {
			s[k] = new(crypto.Scalar).Mul(alpha, proofs[i].a)
			sInverse[k] = new(crypto.Scalar).Mul(alpha, proofs[i].b)
		}

		p := new(crypto.Point).Set(proofs[i].p)
		for j := 0; j < len(proofs[i].l); j++ {
			// calculate challenge x = hash(hash(G || H || u || p) || x || l || r)
			xList[j] = generateChallenge([][]byte{csList[i], p.ToBytesS(), proofs[i].l[j].ToBytesS(), proofs[i].r[j].ToBytesS()})
			xInverseList[j] = new(crypto.Scalar).Invert(xList[j])
			xSquareList[j] = new(crypto.Scalar).Mul(xList[j], xList[j])
			xSquareAlphaList[j] = new(crypto.Scalar).Mul(xSquareList[j], alpha)
			xInverseSquareList[j] = new(crypto.Scalar).Mul(xInverseList[j], xInverseList[j])
			xInverseSquareAlphaList[j] = new(crypto.Scalar).Mul(xInverseSquareList[j], alpha)

			pPrime := new(crypto.Point).AddPedersen(xSquareList[j], proofs[i].l[j], xInverseSquareList[j], proofs[i].r[j])
			pPrime.Add(pPrime, p)
			p = pPrime

			//Update s, s^-1
			for k := 0; k < n; k++ {
				if k&int(math.Pow(2, float64(logN-j-1))) != 0 {
					s[k] = new(crypto.Scalar).Mul(s[k], xList[j])
					sInverse[k] = new(crypto.Scalar).Mul(sInverse[k], xInverseList[j])
				} else {
					s[k] = new(crypto.Scalar).Mul(s[k], xInverseList[j])
					sInverse[k] = new(crypto.Scalar).Mul(sInverse[k], xList[j])
				}
			}
		}
		for k := 0; k < n; k++ {
			asAlphaList[k].Add(asAlphaList[k], s[k])
			bsInverseAlphaList[k].Add(bsInverseAlphaList[k], sInverse[k])
		}

		LList = append(LList, proofs[i].l...)
		nXSquareList = append(nXSquareList, xSquareAlphaList...)
		RList = append(RList, proofs[i].r...)
		nXInverseSquareList = append(nXInverseSquareList, xInverseSquareAlphaList...)
	}

	gAlphaAS := new(crypto.Point).MultiScalarMult(asAlphaList[0:maxN], AggParam.g[0:maxN])
	hAlphaBSInverse := new(crypto.Point).MultiScalarMult(bsInverseAlphaList[0:maxN], AggParam.h[0:maxN])
	LHS := new(crypto.Point).Add(gAlphaAS, hAlphaBSInverse)
	LHS.Add(LHS, new(crypto.Point).ScalarMult(AggParam.u, sum_abAlpha))
	//fmt.Println("LHS:", LHS )

	prod_PAlpha := new(crypto.Point).MultiScalarMult(alphaList, pList)
	prod_LX := new(crypto.Point).MultiScalarMult(nXSquareList, LList)
	prod_RX := new(crypto.Point).MultiScalarMult(nXInverseSquareList, RList)

	RHS := new(crypto.Point).Add(prod_LX, prod_RX)
	RHS.Add(RHS, prod_PAlpha)
	//fmt.Println("RHS:", RHS)

	res := crypto.IsPointEqual(RHS, LHS)
	if !res {
		fmt.Println("Inner product argument failed:")
		fmt.Printf("LHS: %v\n", LHS)
		fmt.Printf("RHS: %v\n", RHS)
	}

	return res
}
