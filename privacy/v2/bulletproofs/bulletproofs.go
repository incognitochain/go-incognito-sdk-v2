package bulletproofs

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
	"math"
)

type Witness struct {
	values []uint64
	rands  []*crypto.Scalar
}

type RangeProof struct {
	cmsValue          []*crypto.Point
	a                 *crypto.Point
	s                 *crypto.Point
	t1                *crypto.Point
	t2                *crypto.Point
	tauX              *crypto.Scalar
	tHat              *crypto.Scalar
	mu                *crypto.Scalar
	innerProductProof *InnerProductProof
}

type bulletproofParams struct {
	g  []*crypto.Point
	h  []*crypto.Point
	u  *crypto.Point
	cs *crypto.Point
}

var AggParam = newBulletproofParams(utils.MaxOutputCoin)

func (proof RangeProof) ValidateSanity() bool {
	for i := 0; i < len(proof.cmsValue); i++ {
		if !proof.cmsValue[i].PointValid() {
			return false
		}
	}
	if !proof.a.PointValid() || !proof.s.PointValid() || !proof.t1.PointValid() || !proof.t2.PointValid() {
		return false
	}
	if !proof.tauX.ScalarValid() || !proof.tHat.ScalarValid() || !proof.mu.ScalarValid() {
		return false
	}

	return proof.innerProductProof.ValidateSanity()
}

func (proof *RangeProof) Init() {
	proof.a = new(crypto.Point).Identity()
	proof.s = new(crypto.Point).Identity()
	proof.t1 = new(crypto.Point).Identity()
	proof.t2 = new(crypto.Point).Identity()
	proof.tauX = new(crypto.Scalar)
	proof.tHat = new(crypto.Scalar)
	proof.mu = new(crypto.Scalar)
	proof.innerProductProof = new(InnerProductProof).Init()
}

func (proof RangeProof) IsNil() bool {
	if proof.a == nil {
		return true
	}
	if proof.s == nil {
		return true
	}
	if proof.t1 == nil {
		return true
	}
	if proof.t2 == nil {
		return true
	}
	if proof.tauX == nil {
		return true
	}
	if proof.tHat == nil {
		return true
	}
	if proof.mu == nil {
		return true
	}
	return proof.innerProductProof == nil
}

func (proof RangeProof) Bytes() []byte {
	var res []byte

	if proof.IsNil() {
		return []byte{}
	}

	res = append(res, byte(len(proof.cmsValue)))
	for i := 0; i < len(proof.cmsValue); i++ {
		res = append(res, proof.cmsValue[i].ToBytesS()...)
	}

	res = append(res, proof.a.ToBytesS()...)
	res = append(res, proof.s.ToBytesS()...)
	res = append(res, proof.t1.ToBytesS()...)
	res = append(res, proof.t2.ToBytesS()...)

	res = append(res, proof.tauX.ToBytesS()...)
	res = append(res, proof.tHat.ToBytesS()...)
	res = append(res, proof.mu.ToBytesS()...)
	res = append(res, proof.innerProductProof.Bytes()...)

	return res
}

func (proof RangeProof) GetCommitments() []*crypto.Point { return proof.cmsValue }

func (proof *RangeProof) SetCommitments(cmsValue []*crypto.Point) {
	proof.cmsValue = cmsValue
}

func (proof *RangeProof) SetBytes(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	lenValues := int(bytes[0])
	offset := 1
	var err error

	proof.cmsValue = make([]*crypto.Point, lenValues)
	for i := 0; i < lenValues; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("unmarshalling from bytes failed")
		}
		proof.cmsValue[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset += crypto.Ed25519KeySize
	}

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling from bytes failed")
	}
	proof.a, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling from bytes failed")
	}
	proof.s, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling from bytes failed")
	}
	proof.t1, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling from bytes failed")
	}
	proof.t2, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling from bytes failed")
	}
	proof.tauX = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling from bytes failed")
	}
	proof.tHat = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("unmarshalling from bytes failed")
	}
	proof.mu = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	if offset >= len(bytes) {
		return fmt.Errorf("unmarshalling from bytes failed")
	}

	proof.innerProductProof = new(InnerProductProof)
	err = proof.innerProductProof.SetBytes(bytes[offset:])

	return err
}

func (wit *Witness) Set(values []uint64, rands []*crypto.Scalar) {
	numValue := len(values)
	wit.values = make([]uint64, numValue)
	wit.rands = make([]*crypto.Scalar, numValue)

	for i := range values {
		wit.values[i] = values[i]
		wit.rands[i] = new(crypto.Scalar).Set(rands[i])
	}
}

func (wit Witness) Prove() (*RangeProof, error) {
	proof := new(RangeProof)
	numValue := len(wit.values)
	if numValue > utils.MaxOutputCoin {
		return nil, fmt.Errorf("must less than MaxOutputCoin")
	}
	numValuePad := roundUpPowTwo(numValue)
	maxExp := utils.MaxExp
	N := maxExp * numValuePad

	aggParam := setAggregateParams(N)

	values := make([]uint64, numValuePad)
	rands := make([]*crypto.Scalar, numValuePad)
	for i := range wit.values {
		values[i] = wit.values[i]
		rands[i] = new(crypto.Scalar).Set(wit.rands[i])
	}
	for i := numValue; i < numValuePad; i++ {
		values[i] = uint64(0)
		rands[i] = new(crypto.Scalar).FromUint64(0)
	}

	proof.cmsValue = make([]*crypto.Point, numValue)
	for i := 0; i < numValue; i++ {
		proof.cmsValue[i] = crypto.PedCom.CommitAtIndex(new(crypto.Scalar).FromUint64(values[i]), rands[i], crypto.PedersenValueIndex)
	}
	// Convert values to binary array
	aL := make([]*crypto.Scalar, N)
	aR := make([]*crypto.Scalar, N)
	sL := make([]*crypto.Scalar, N)
	sR := make([]*crypto.Scalar, N)

	for i, value := range values {
		tmp := ConvertUint64ToBinary(value, maxExp)
		for j := 0; j < maxExp; j++ {
			aL[i*maxExp+j] = tmp[j]
			aR[i*maxExp+j] = new(crypto.Scalar).Sub(tmp[j], new(crypto.Scalar).FromUint64(1))
			sL[i*maxExp+j] = crypto.RandomScalar()
			sR[i*maxExp+j] = crypto.RandomScalar()
		}
	}
	// LINE 40-50
	// Commitment to aL, aR: A = h^alpha * G^aL * H^aR
	// Commitment to sL, sR : S = h^rho * G^sL * H^sR
	var alpha, rho *crypto.Scalar
	if A, err := encodeVectors(aL, aR, aggParam.g, aggParam.h); err != nil {
		return nil, err
	} else if S, err := encodeVectors(sL, sR, aggParam.g, aggParam.h); err != nil {
		return nil, err
	} else {
		alpha = crypto.RandomScalar()
		rho = crypto.RandomScalar()
		A.Add(A, new(crypto.Point).ScalarMult(crypto.HBase, alpha))
		S.Add(S, new(crypto.Point).ScalarMult(crypto.HBase, rho))
		proof.a = A
		proof.s = S
	}
	// challenge y, z
	y := generateChallenge(aggParam.cs.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	z := generateChallenge(y.ToBytesS(), []*crypto.Point{proof.a, proof.s})

	// LINE 51-54
	twoNumber := new(crypto.Scalar).FromUint64(2)
	twoVectorN := powerVector(twoNumber, maxExp)

	// HPrime = H^(y^(1-i)
	HPrime := computeHPrime(y, N, aggParam.h)

	// l(X) = (aL -z*1^n) + sL*X; r(X) = y^n Hadamard (aR +z*1^n + sR*X) + z^2 * 2^n
	yVector := powerVector(y, N)
	hadProduct, err := hadamardProduct(yVector, vectorAddScalar(aR, z))
	if err != nil {
		return nil, err
	}
	vectorSum := make([]*crypto.Scalar, N)
	zTmp := new(crypto.Scalar).Set(z)
	for j := 0; j < numValuePad; j++ {
		zTmp.Mul(zTmp, z)
		for i := 0; i < maxExp; i++ {
			vectorSum[j*maxExp+i] = new(crypto.Scalar).Mul(twoVectorN[i], zTmp)
		}
	}
	zNeg := new(crypto.Scalar).Sub(new(crypto.Scalar).FromUint64(0), z)
	l0 := vectorAddScalar(aL, zNeg)
	l1 := sL
	var r0, r1 []*crypto.Scalar
	if r0, err = vectorAdd(hadProduct, vectorSum); err != nil {
		return nil, err
	} else {
		if r1, err = hadamardProduct(yVector, sR); err != nil {
			return nil, err
		}
	}

	// t(X) = <l(X), r(X)> = t0 + t1*X + t2*X^2
	// t1 = <l1, ro> + <l0, r1>, t2 = <l1, r1>
	var t1, t2 *crypto.Scalar
	if ip3, err := innerProduct(l1, r0); err != nil {
		return nil, err
	} else if ip4, err := innerProduct(l0, r1); err != nil {
		return nil, err
	} else {
		t1 = new(crypto.Scalar).Add(ip3, ip4)
		if t2, err = innerProduct(l1, r1); err != nil {
			return nil, err
		}
	}

	// commitment to t1, t2
	tau1 := crypto.RandomScalar()
	tau2 := crypto.RandomScalar()
	proof.t1 = crypto.PedCom.CommitAtIndex(t1, tau1, crypto.PedersenValueIndex)
	proof.t2 = crypto.PedCom.CommitAtIndex(t2, tau2, crypto.PedersenValueIndex)

	x := generateChallenge(z.ToBytesS(), []*crypto.Point{proof.t1, proof.t2})
	xSquare := new(crypto.Scalar).Mul(x, x)

	// lVector = aL - z*1^n + sL*x
	// rVector = y^n Hadamard (aR +z*1^n + sR*x) + z^2*2^n
	// tHat = <lVector, rVector>
	lVector, err := vectorAdd(vectorAddScalar(aL, zNeg), vectorMulScalar(sL, x))
	if err != nil {
		return nil, err
	}
	tmpVector, err := vectorAdd(vectorAddScalar(aR, z), vectorMulScalar(sR, x))
	if err != nil {
		return nil, err
	}
	rVector, err := hadamardProduct(yVector, tmpVector)
	if err != nil {
		return nil, err
	}
	rVector, err = vectorAdd(rVector, vectorSum)
	if err != nil {
		return nil, err
	}
	proof.tHat, err = innerProduct(lVector, rVector)
	if err != nil {
		return nil, err
	}

	// blinding value for tHat: tauX = tau2*x^2 + tau1*x + z^2*rand
	proof.tauX = new(crypto.Scalar).Mul(tau2, xSquare)
	proof.tauX.Add(proof.tauX, new(crypto.Scalar).Mul(tau1, x))
	zTmp = new(crypto.Scalar).Set(z)
	tmpBN := new(crypto.Scalar)
	for j := 0; j < numValuePad; j++ {
		zTmp.Mul(zTmp, z)
		proof.tauX.Add(proof.tauX, tmpBN.Mul(zTmp, rands[j]))
	}

	// alpha, rho blind A, S
	// mu = alpha + rho*x
	proof.mu = new(crypto.Scalar).Add(alpha, new(crypto.Scalar).Mul(rho, x))

	// instead of sending left vector and right vector, we use inner sum argument to reduce proof size from 2*n to 2(log2(n)) + 2
	innerProductWit := new(InnerProductWitness)
	innerProductWit.a = lVector
	innerProductWit.b = rVector
	innerProductWit.p, err = encodeVectors(lVector, rVector, aggParam.g, HPrime)
	if err != nil {
		return nil, err
	}
	uPrime := new(crypto.Point).ScalarMult(aggParam.u, crypto.HashToScalar(x.ToBytesS()))
	innerProductWit.p = innerProductWit.p.Add(innerProductWit.p, new(crypto.Point).ScalarMult(uPrime, proof.tHat))

	proof.innerProductProof, err = innerProductWit.Prove(aggParam.g, HPrime, uPrime, x.ToBytesS())
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func (proof RangeProof) Verify() (bool, error) {
	numValue := len(proof.cmsValue)
	if numValue > utils.MaxOutputCoin {
		return false, fmt.Errorf("must less than MaxOutputNumber")
	}
	numValuePad := roundUpPowTwo(numValue)
	maxExp := utils.MaxExp
	N := numValuePad * maxExp
	twoVectorN := powerVector(new(crypto.Scalar).FromUint64(2), maxExp)
	aggParam := setAggregateParams(N)

	cmsValue := proof.cmsValue
	for i := numValue; i < numValuePad; i++ {
		cmsValue = append(cmsValue, new(crypto.Point).Identity())
	}

	// recalculate challenge y, z
	y := generateChallenge(aggParam.cs.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	z := generateChallenge(y.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	zSquare := new(crypto.Scalar).Mul(z, z)
	zNeg := new(crypto.Scalar).Sub(new(crypto.Scalar).FromUint64(0), z)

	x := generateChallenge(z.ToBytesS(), []*crypto.Point{proof.t1, proof.t2})
	xSquare := new(crypto.Scalar).Mul(x, x)

	// HPrime = H^(y^(1-i)
	HPrime := computeHPrime(y, N, aggParam.h)

	// g^tHat * h^tauX = V^(z^2) * g^delta(y,z) * T1^x * T2^(x^2)
	yVector := powerVector(y, N)
	deltaYZ, err := computeDeltaYZ(z, zSquare, yVector, N)
	if err != nil {
		return false, err
	}

	LHS := crypto.PedCom.CommitAtIndex(proof.tHat, proof.tauX, crypto.PedersenValueIndex)
	RHS := new(crypto.Point).ScalarMult(proof.t2, xSquare)
	RHS.Add(RHS, new(crypto.Point).AddPedersen(deltaYZ, crypto.PedCom.G[crypto.PedersenValueIndex], x, proof.t1))

	expVector := vectorMulScalar(powerVector(z, numValuePad), zSquare)
	RHS.Add(RHS, new(crypto.Point).MultiScalarMult(expVector, cmsValue))

	if !crypto.IsPointEqual(LHS, RHS) {
		return false, fmt.Errorf("verify aggregated range proof statement 1 failed")
	}

	// verify eq (66)
	uPrime := new(crypto.Point).ScalarMult(aggParam.u, crypto.HashToScalar(x.ToBytesS()))

	vectorSum := make([]*crypto.Scalar, N)
	zTmp := new(crypto.Scalar).Set(z)
	for j := 0; j < numValuePad; j++ {
		zTmp.Mul(zTmp, z)
		for i := 0; i < maxExp; i++ {
			vectorSum[j*maxExp+i] = new(crypto.Scalar).Mul(twoVectorN[i], zTmp)
			vectorSum[j*maxExp+i].Add(vectorSum[j*maxExp+i], new(crypto.Scalar).Mul(z, yVector[j*maxExp+i]))
		}
	}
	tmpHPrime := new(crypto.Point).MultiScalarMult(vectorSum, HPrime)
	tmpG := new(crypto.Point).Set(aggParam.g[0])
	for i := 1; i < N; i++ {
		tmpG.Add(tmpG, aggParam.g[i])
	}
	ASx := new(crypto.Point).Add(proof.a, new(crypto.Point).ScalarMult(proof.s, x))
	P := new(crypto.Point).Add(new(crypto.Point).ScalarMult(tmpG, zNeg), tmpHPrime)
	P.Add(P, ASx)
	P.Add(P, new(crypto.Point).ScalarMult(uPrime, proof.tHat))
	PPrime := new(crypto.Point).Add(proof.innerProductProof.p, new(crypto.Point).ScalarMult(crypto.HBase, proof.mu))

	if !crypto.IsPointEqual(P, PPrime) {
		return false, fmt.Errorf("verify aggregated range proof statement 2-1 failed")
	}

	// verify eq (68)
	innerProductArgValid := proof.innerProductProof.Verify(aggParam.g, HPrime, uPrime, x.ToBytesS())
	if !innerProductArgValid {
		return false, fmt.Errorf("verify aggregated range proof statement 2 failed")
	}

	return true, nil
}

func (proof RangeProof) VerifyFaster() (bool, error) {
	numValue := len(proof.cmsValue)
	if numValue > utils.MaxOutputCoin {
		return false, fmt.Errorf("must less than MaxOutputNumber")
	}
	numValuePad := roundUpPowTwo(numValue)
	maxExp := utils.MaxExp
	N := maxExp * numValuePad
	aggParam := setAggregateParams(N)
	twoVectorN := powerVector(new(crypto.Scalar).FromUint64(2), maxExp)

	cmsValue := proof.cmsValue
	for i := numValue; i < numValuePad; i++ {
		cmsValue = append(cmsValue, new(crypto.Point).Identity())
	}

	// recalculate challenge y, z
	y := generateChallenge(aggParam.cs.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	z := generateChallenge(y.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	zSquare := new(crypto.Scalar).Mul(z, z)
	zNeg := new(crypto.Scalar).Sub(new(crypto.Scalar).FromUint64(0), z)

	x := generateChallenge(z.ToBytesS(), []*crypto.Point{proof.t1, proof.t2})
	xSquare := new(crypto.Scalar).Mul(x, x)

	// g^tHat * h^tauX = V^(z^2) * g^delta(y,z) * T1^x * T2^(x^2)
	yVector := powerVector(y, N)
	deltaYZ, err := computeDeltaYZ(z, zSquare, yVector, N)
	if err != nil {
		return false, err
	}
	// HPrime = H^(y^(1-i)
	HPrime := computeHPrime(y, N, aggParam.h)
	uPrime := new(crypto.Point).ScalarMult(aggParam.u, crypto.HashToScalar(x.ToBytesS()))

	// Verify eq (65)
	LHS := crypto.PedCom.CommitAtIndex(proof.tHat, proof.tauX, crypto.PedersenValueIndex)
	RHS := new(crypto.Point).ScalarMult(proof.t2, xSquare)
	RHS.Add(RHS, new(crypto.Point).AddPedersen(deltaYZ, crypto.PedCom.G[crypto.PedersenValueIndex], x, proof.t1))
	expVector := vectorMulScalar(powerVector(z, numValuePad), zSquare)
	RHS.Add(RHS, new(crypto.Point).MultiScalarMult(expVector, cmsValue))
	if !crypto.IsPointEqual(LHS, RHS) {
		return false, fmt.Errorf("verify aggregated range proof statement 1 failed")
	}

	// Verify eq (66)
	vectorSum := make([]*crypto.Scalar, N)
	zTmp := new(crypto.Scalar).Set(z)
	for j := 0; j < numValuePad; j++ {
		zTmp.Mul(zTmp, z)
		for i := 0; i < maxExp; i++ {
			vectorSum[j*maxExp+i] = new(crypto.Scalar).Mul(twoVectorN[i], zTmp)
			vectorSum[j*maxExp+i].Add(vectorSum[j*maxExp+i], new(crypto.Scalar).Mul(z, yVector[j*maxExp+i]))
		}
	}
	tmpHPrime := new(crypto.Point).MultiScalarMult(vectorSum, HPrime)
	tmpG := new(crypto.Point).Set(aggParam.g[0])
	for i := 1; i < N; i++ {
		tmpG.Add(tmpG, aggParam.g[i])
	}
	ASx := new(crypto.Point).Add(proof.a, new(crypto.Point).ScalarMult(proof.s, x))
	P := new(crypto.Point).Add(new(crypto.Point).ScalarMult(tmpG, zNeg), tmpHPrime)
	P.Add(P, ASx)
	P.Add(P, new(crypto.Point).ScalarMult(uPrime, proof.tHat))
	PPrime := new(crypto.Point).Add(proof.innerProductProof.p, new(crypto.Point).ScalarMult(crypto.HBase, proof.mu))

	if !crypto.IsPointEqual(P, PPrime) {
		return false, fmt.Errorf("verify aggregated range proof statement 2-1 failed")
	}

	// Verify eq (68)
	hashCache := x.ToBytesS()
	L := proof.innerProductProof.l
	R := proof.innerProductProof.r
	s := make([]*crypto.Scalar, N)
	sInverse := make([]*crypto.Scalar, N)
	logN := int(math.Log2(float64(N)))
	vSquareList := make([]*crypto.Scalar, logN)
	vInverseSquareList := make([]*crypto.Scalar, logN)

	for i := 0; i < N; i++ {
		s[i] = new(crypto.Scalar).Set(proof.innerProductProof.a)
		sInverse[i] = new(crypto.Scalar).Set(proof.innerProductProof.b)
	}

	for i := range L {
		v := generateChallenge(hashCache, []*crypto.Point{L[i], R[i]})
		hashCache = v.ToBytesS()
		vInverse := new(crypto.Scalar).Invert(v)
		vSquareList[i] = new(crypto.Scalar).Mul(v, v)
		vInverseSquareList[i] = new(crypto.Scalar).Mul(vInverse, vInverse)

		for j := 0; j < N; j++ {
			if j&int(math.Pow(2, float64(logN-i-1))) != 0 {
				s[j] = new(crypto.Scalar).Mul(s[j], v)
				sInverse[j] = new(crypto.Scalar).Mul(sInverse[j], vInverse)
			} else {
				s[j] = new(crypto.Scalar).Mul(s[j], vInverse)
				sInverse[j] = new(crypto.Scalar).Mul(sInverse[j], v)
			}
		}
	}

	c := new(crypto.Scalar).Mul(proof.innerProductProof.a, proof.innerProductProof.b)
	tmp1 := new(crypto.Point).MultiScalarMult(s, aggParam.g)
	tmp2 := new(crypto.Point).MultiScalarMult(sInverse, HPrime)
	rightHS := new(crypto.Point).Add(tmp1, tmp2)
	rightHS.Add(rightHS, new(crypto.Point).ScalarMult(uPrime, c))

	tmp3 := new(crypto.Point).MultiScalarMult(vSquareList, L)
	tmp4 := new(crypto.Point).MultiScalarMult(vInverseSquareList, R)
	leftHS := new(crypto.Point).Add(tmp3, tmp4)
	leftHS.Add(leftHS, proof.innerProductProof.p)

	res := crypto.IsPointEqual(rightHS, leftHS)
	if !res {
		return false, fmt.Errorf("verify aggregated range proof statement 2 failed")
	}

	return true, nil
}

func VerifyBatch(proofs []*RangeProof) (bool, error, int) {
	maxExp := utils.MaxExp
	baseG := crypto.PedCom.G[crypto.PedersenValueIndex]
	baseH := crypto.PedCom.G[crypto.PedersenRandomnessIndex]

	sumTHat := new(crypto.Scalar).FromUint64(0)
	sumTauX := new(crypto.Scalar).FromUint64(0)
	xAlphaList := make([]*crypto.Scalar, 0)
	xBetaList := make([]*crypto.Scalar, 0)
	xSquareList := make([]*crypto.Scalar, 0)
	zSquareList := make([]*crypto.Scalar, 0)

	t1List := make([]*crypto.Point, 0)
	txList := make([]*crypto.Point, 0)
	vList := make([]*crypto.Point, 0)

	muSum := new(crypto.Scalar).FromUint64(0)
	tempSum := new(crypto.Scalar).FromUint64(0) // sum of ab - tHat

	sList := make([]*crypto.Point, 0)
	aList := make([]*crypto.Point, 0)
	betaList := make([]*crypto.Scalar, 0)
	LRList := make([]*crypto.Point, 0)
	lVectorList := make([]*crypto.Scalar, 0)
	rVectorList := make([]*crypto.Scalar, 0)
	gVectorList := make([]*crypto.Point, 0)
	hVectorList := make([]*crypto.Point, 0)

	twoNumber := new(crypto.Scalar).FromUint64(2)
	twoVectorN := powerVector(twoNumber, maxExp)

	for k, proof := range proofs {
		numValue := len(proof.cmsValue)
		if numValue > utils.MaxOutputCoin {
			return false, fmt.Errorf("must less than MaxOutputNumber"), k
		}
		numValuePad := roundUpPowTwo(numValue)
		N := maxExp * numValuePad
		aggParam := setAggregateParams(N)

		cmsValue := proof.cmsValue
		for i := numValue; i < numValuePad; i++ {
			identity := new(crypto.Point).Identity()
			cmsValue = append(cmsValue, identity)
		}

		// recalculate challenge y, z, x
		y := generateChallenge(aggParam.cs.ToBytesS(), []*crypto.Point{proof.a, proof.s})
		z := generateChallenge(y.ToBytesS(), []*crypto.Point{proof.a, proof.s})
		x := generateChallenge(z.ToBytesS(), []*crypto.Point{proof.t1, proof.t2})
		zSquare := new(crypto.Scalar).Mul(z, z)
		xSquare := new(crypto.Scalar).Mul(x, x)

		// Random alpha and beta for batch equations check
		alpha := crypto.RandomScalar()
		beta := crypto.RandomScalar()
		betaList = append(betaList, beta)

		// Compute first equation check
		yVector := powerVector(y, N)
		deltaYZ, err := computeDeltaYZ(z, zSquare, yVector, N)
		if err != nil {
			return false, err, k
		}
		sumTHat.Add(sumTHat, new(crypto.Scalar).Mul(alpha, new(crypto.Scalar).Sub(proof.tHat, deltaYZ)))
		sumTauX.Add(sumTauX, new(crypto.Scalar).Mul(alpha, proof.tauX))

		xAlphaList = append(xAlphaList, new(crypto.Scalar).Mul(x, alpha))
		xBetaList = append(xBetaList, new(crypto.Scalar).Mul(x, beta))
		xSquareList = append(xSquareList, new(crypto.Scalar).Mul(xSquare, alpha))
		tmp := vectorMulScalar(powerVector(z, numValuePad), new(crypto.Scalar).Mul(zSquare, alpha))
		zSquareList = append(zSquareList, tmp...)

		vList = append(vList, cmsValue...)
		t1List = append(t1List, proof.t1)
		txList = append(txList, proof.t2)

		// Verify the second argument
		hashCache := x.ToBytesS()
		L := proof.innerProductProof.l
		R := proof.innerProductProof.r
		s := make([]*crypto.Scalar, N)
		sInverse := make([]*crypto.Scalar, N)
		logN := int(math.Log2(float64(N)))
		vSquareList := make([]*crypto.Scalar, logN)
		vInverseSquareList := make([]*crypto.Scalar, logN)

		for i := 0; i < N; i++ {
			s[i] = new(crypto.Scalar).Set(proof.innerProductProof.a)
			sInverse[i] = new(crypto.Scalar).Set(proof.innerProductProof.b)
		}

		for i := range L {
			v := generateChallenge(hashCache, []*crypto.Point{L[i], R[i]})
			hashCache = v.ToBytesS()
			vInverse := new(crypto.Scalar).Invert(v)
			vSquareList[i] = new(crypto.Scalar).Mul(v, v)
			vInverseSquareList[i] = new(crypto.Scalar).Mul(vInverse, vInverse)

			for j := 0; j < N; j++ {
				if j&int(math.Pow(2, float64(logN-i-1))) != 0 {
					s[j] = new(crypto.Scalar).Mul(s[j], v)
					sInverse[j] = new(crypto.Scalar).Mul(sInverse[j], vInverse)
				} else {
					s[j] = new(crypto.Scalar).Mul(s[j], vInverse)
					sInverse[j] = new(crypto.Scalar).Mul(sInverse[j], v)
				}
			}
		}

		lVector := make([]*crypto.Scalar, N)
		rVector := make([]*crypto.Scalar, N)

		vectorSum := make([]*crypto.Scalar, N)
		zTmp := new(crypto.Scalar).Set(z)
		for j := 0; j < numValuePad; j++ {
			zTmp.Mul(zTmp, z)
			for i := 0; i < maxExp; i++ {
				vectorSum[j*maxExp+i] = new(crypto.Scalar).Mul(twoVectorN[i], zTmp)
			}
		}
		yInverse := new(crypto.Scalar).Invert(y)
		yTmp := new(crypto.Scalar).Set(y)
		for j := 0; j < N; j++ {
			yTmp.Mul(yTmp, yInverse)
			lVector[j] = new(crypto.Scalar).Add(s[j], z)
			rVector[j] = new(crypto.Scalar).Sub(sInverse[j], vectorSum[j])
			rVector[j].Mul(rVector[j], yTmp)
			rVector[j].Sub(rVector[j], z)

			lVector[j].Mul(lVector[j], beta)
			rVector[j].Mul(rVector[j], beta)
		}

		lVectorList = append(lVectorList, lVector...)
		rVectorList = append(rVectorList, rVector...)

		tmp1 := new(crypto.Point).MultiScalarMult(vSquareList, L)
		tmp2 := new(crypto.Point).MultiScalarMult(vInverseSquareList, R)
		LRList = append(LRList, new(crypto.Point).Add(tmp1, tmp2))

		gVectorList = append(gVectorList, aggParam.g...)
		hVectorList = append(hVectorList, aggParam.h...)

		muSum.Add(muSum, new(crypto.Scalar).Mul(proof.mu, beta))
		ab := new(crypto.Scalar).Mul(proof.innerProductProof.a, proof.innerProductProof.b)
		tmpDiff := new(crypto.Scalar).Sub(ab, proof.tHat)       // ab - tHat
		tmpDiff.Mul(tmpDiff, crypto.HashToScalar(x.ToBytesS())) // (ab - tHat) * Hash(x)
		tempSum.Add(tempSum, new(crypto.Scalar).Mul(tmpDiff, beta))
		aList = append(aList, proof.a)
		sList = append(sList, proof.s)
	}

	tmp1 := new(crypto.Point).MultiScalarMult(lVectorList, gVectorList)
	tmp2 := new(crypto.Point).MultiScalarMult(rVectorList, hVectorList)
	tmp3 := new(crypto.Point).ScalarMult(AggParam.u, tempSum)
	tmp4 := new(crypto.Point).ScalarMult(baseH, muSum)
	LHSPrime := new(crypto.Point).Add(tmp1, tmp2)
	LHSPrime.Add(LHSPrime, tmp3)
	LHSPrime.Add(LHSPrime, tmp4)

	LHS := new(crypto.Point).AddPedersen(sumTHat, baseG, sumTauX, baseH)
	LHSPrime.Add(LHSPrime, LHS)

	tmp5 := new(crypto.Point).MultiScalarMult(betaList, aList)
	tmp6 := new(crypto.Point).MultiScalarMult(xBetaList, sList)
	RHSPrime := new(crypto.Point).Add(tmp5, tmp6)
	RHSPrime.Add(RHSPrime, new(crypto.Point).MultiScalarMult(betaList, LRList))

	part1 := new(crypto.Point).MultiScalarMult(xAlphaList, t1List)
	part2 := new(crypto.Point).MultiScalarMult(xSquareList, txList)
	RHS := new(crypto.Point).Add(part1, part2)
	RHS.Add(RHS, new(crypto.Point).MultiScalarMult(zSquareList, vList))
	RHSPrime.Add(RHSPrime, RHS)

	if !crypto.IsPointEqual(LHSPrime, RHSPrime) {
		return false, fmt.Errorf("batch verify aggregated range proof failed"), -1
	}
	return true, nil, -1
}

// EstimateMultiRangeProofSize estimates the size of a multi-range proof.
func EstimateMultiRangeProofSize(nOutput int) uint64 {
	return uint64((nOutput+2*int(math.Log2(float64(utils.MaxExp*roundUpPowTwo(nOutput))))+5)*crypto.Ed25519KeySize + 5*crypto.Ed25519KeySize + 2)
}
