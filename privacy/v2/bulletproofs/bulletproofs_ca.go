package bulletproofs

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
)

var cACommitmentScheme = copyPedersenCommitmentScheme(crypto.PedCom)

func copyPedersenCommitmentScheme(sch crypto.PedersenCommitment) crypto.PedersenCommitment {
	var result crypto.PedersenCommitment
	var generators []*crypto.Point
	for _, gen := range sch.G {
		generators = append(generators, new(crypto.Point).Set(gen))
	}
	result.G = generators
	return result
}

// GetFirstAssetTag returns the first asset tag of a list out CoinV2's.
func GetFirstAssetTag(coins []*coin.CoinV2) (*crypto.Point, error) {
	if coins == nil || len(coins) == 0 {
		return nil, fmt.Errorf("cannot get asset tag from empty input")
	}
	result := coins[0].GetAssetTag()
	if result == nil {
		return nil, fmt.Errorf("the coin does not have an asset tag")
	}
	return result, nil
}

// ProveUsingBase returns a RangeProof for a Witness created using the base point v.
func (wit Witness) ProveUsingBase(b *crypto.Point) (*RangeProof, error) {
	cACommitmentScheme.G[crypto.PedersenValueIndex] = b
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
		A.Add(A, new(crypto.Point).ScalarMult(cACommitmentScheme.G[crypto.PedersenRandomnessIndex], alpha))
		S.Add(S, new(crypto.Point).ScalarMult(cACommitmentScheme.G[crypto.PedersenRandomnessIndex], rho))
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
	proof.t1 = cACommitmentScheme.CommitAtIndex(t1, tau1, crypto.PedersenValueIndex)
	proof.t2 = cACommitmentScheme.CommitAtIndex(t2, tau2, crypto.PedersenValueIndex)

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

// TransformWitnessToCAWitness transforms a regular Witness into a new Witness used in confidential assets,
// given the asset tag blinders.
func TransformWitnessToCAWitness(wit *Witness, assetTagBlinders []*crypto.Scalar) (*Witness, error) {
	if len(assetTagBlinders) != len(wit.values) || len(assetTagBlinders) != len(wit.rands) {
		return nil, fmt.Errorf("cannot transform witness: parameter lengths mismatch")
	}
	newRands := make([]*crypto.Scalar, len(wit.values))

	for i := range wit.values {
		temp := new(crypto.Scalar).Sub(assetTagBlinders[i], assetTagBlinders[0])
		temp.Mul(temp, new(crypto.Scalar).FromUint64(wit.values[i]))
		temp.Add(temp, wit.rands[i])
		newRands[i] = temp
	}
	result := new(Witness)
	result.Set(wit.values, newRands)
	return result, nil
}
