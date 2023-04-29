package bulletproofs

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
)

type bulletproofParams struct {
	g  []*crypto.Point
	h  []*crypto.Point
	u  *crypto.Point
	cs *crypto.Point
}

// Witness represents a Bulletproofs witness.
type Witness struct {
	values []uint64
	rands  []*crypto.Scalar
}

// RangeProof represents a Bulletproofs proof.
type RangeProof struct {
	version           uint8
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

var aggParam = newBulletproofParams(utils.MaxOutputCoin)

// Init creates an empty RangeProof.
func (proof *RangeProof) Init() {
	proof.a = new(crypto.Point).Identity()
	proof.s = new(crypto.Point).Identity()
	proof.t1 = new(crypto.Point).Identity()
	proof.t2 = new(crypto.Point).Identity()
	proof.tauX = new(crypto.Scalar)
	proof.tHat = new(crypto.Scalar)
	proof.mu = new(crypto.Scalar)
	proof.innerProductProof = new(InnerProductProof).Init()
	proof.version = 2
}

func (proof RangeProof) GetVersion() uint8 {
	return proof.version
}

// IsNil checks if a RangeProof is empty.
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

// Bytes returns the byte-representation of a RangeProof.
func (proof RangeProof) Bytes() []byte {
	var res []byte

	if proof.version >= 2 {
		res = append(res, byte(0))
		res = append(res, byte(proof.version))
	}

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

// GetCommitments returns the commitments of a RangeProof.
func (proof RangeProof) GetCommitments() []*crypto.Point { return proof.cmsValue }

// SetCommitments sets v as the commitments of a RangeProof.
func (proof *RangeProof) SetCommitments(v []*crypto.Point) {
	proof.cmsValue = v
}

// SetBytes sets byte-representation data to a RangeProof.
func (proof *RangeProof) SetBytes(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	if bytes[0] == 0 {
		// parse versions
		proof.version = uint8(bytes[1])
		bytes = bytes[2:]
	} else {
		proof.version = 1
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

// Set sets parameters for a Witness.
func (wit *Witness) Set(values []uint64, rands []*crypto.Scalar) {
	numValue := len(values)
	wit.values = make([]uint64, numValue)
	wit.rands = make([]*crypto.Scalar, numValue)

	for i := range values {
		wit.values[i] = values[i]
		wit.rands[i] = new(crypto.Scalar).Set(rands[i])
	}
}

// Prove returns the RangeProof for a Witness.
func (wit Witness) Prove() (*RangeProof, error) {
	proof := new(RangeProof)
	proof.Init()
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
	initChal := aggParam.cs.ToBytesS()
	for i := 0; i < numValue; i++ {
		proof.cmsValue[i] = crypto.PedCom.CommitAtIndex(new(crypto.Scalar).FromUint64(values[i]), rands[i], crypto.PedersenValueIndex)
		if proof.version >= 2 {
			initChal = append(initChal, proof.cmsValue[i].ToBytesS()...)
		}
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
	// y := generateChallenge(aggParam.cs.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	y := generateChallenge(initChal, []*crypto.Point{proof.a, proof.s})
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
