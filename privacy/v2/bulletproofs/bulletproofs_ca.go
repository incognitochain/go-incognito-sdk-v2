package bulletproofs

import (
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
	"github.com/pkg/errors"
	"math"
)

var CACommitmentScheme crypto.PedersenCommitment = CopyPedersenCommitmentScheme(crypto.PedCom)

func CopyPedersenCommitmentScheme(sch crypto.PedersenCommitment) crypto.PedersenCommitment{
	var result crypto.PedersenCommitment
	var generators []*crypto.Point
	for _, gen := range sch.G{
		generators = append(generators, new(crypto.Point).Set(gen))
	}
	result.G = generators
	return result
}

func GetFirstAssetTag(coins []*coin.CoinV2) (*crypto.Point,error){
	if coins==nil || len(coins)==0{
		return nil, errors.New("Cannot get asset tag from empty input")
	}
	result := coins[0].GetAssetTag()
	if result==nil{
		return nil, errors.New("The coin does not have an asset tag")
	}
	return result, nil
}

func (wit Witness) ProveUsingBase(anAssetTag *crypto.Point) (*RangeProof, error) {
	CACommitmentScheme.G[crypto.PedersenValueIndex] = anAssetTag
	proof := new(RangeProof)
	numValue := len(wit.values)
	if numValue > utils.MaxOutputCoin {
		return nil, errors.New("Must less than MaxOutputCoin")
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
		A.Add(A, new(crypto.Point).ScalarMult(CACommitmentScheme.G[crypto.PedersenRandomnessIndex], alpha))
		S.Add(S, new(crypto.Point).ScalarMult(CACommitmentScheme.G[crypto.PedersenRandomnessIndex], rho))
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

	// l(X) = (aL -z*1^n) + sL*X; r(X) = y^n hada (aR +z*1^n + sR*X) + z^2 * 2^n
	yVector := powerVector(y, N)
	hadaProduct, err := hadamardProduct(yVector, vectorAddScalar(aR, z))
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
	if r0, err = vectorAdd(hadaProduct, vectorSum); err != nil {
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
	proof.t1 = CACommitmentScheme.CommitAtIndex(t1, tau1, crypto.PedersenValueIndex)
	proof.t2 = CACommitmentScheme.CommitAtIndex(t2, tau2, crypto.PedersenValueIndex)

	x := generateChallenge(z.ToBytesS(), []*crypto.Point{proof.t1, proof.t2})
	xSquare := new(crypto.Scalar).Mul(x, x)

	// lVector = aL - z*1^n + sL*x
	// rVector = y^n hada (aR +z*1^n + sR*x) + z^2*2^n
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

func (proof RangeProof) VerifyUsingBase(anAssetTag *crypto.Point) (bool, error) {
	CACommitmentScheme.G[crypto.PedersenValueIndex] = anAssetTag
	numValue := len(proof.cmsValue)
	if numValue > utils.MaxOutputCoin {
		return false, errors.New("Must less than MaxOutputNumber")
	}
	numValuePad := roundUpPowTwo(numValue)
	maxExp := utils.MaxExp
	N := numValuePad * maxExp
	aggParam := setAggregateParams(N)

	cmsValue := proof.cmsValue
	for i := numValue; i < numValuePad; i++ {
		cmsValue = append(cmsValue, new(crypto.Point).Identity())
	}

	// recalculate challenge y, z
	y := generateChallenge(aggParam.cs.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	z := generateChallenge(y.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	zSquare := new(crypto.Scalar).Mul(z, z)

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

	LHS := CACommitmentScheme.CommitAtIndex(proof.tHat, proof.tauX, crypto.PedersenValueIndex)
	RHS := new(crypto.Point).ScalarMult(proof.t2, xSquare)
	RHS.Add(RHS, new(crypto.Point).AddPedersen(deltaYZ, CACommitmentScheme.G[crypto.PedersenValueIndex], x, proof.t1))

	expVector := vectorMulScalar(powerVector(z, numValuePad), zSquare)
	RHS.Add(RHS, new(crypto.Point).MultiScalarMult(expVector, cmsValue))

	if !crypto.IsPointEqual(LHS, RHS) {
		return false, errors.New("verify aggregated range proof statement 1 failed")
	}
	uPrime := new(crypto.Point).ScalarMult(aggParam.u, crypto.HashToScalar(x.ToBytesS()))
	innerProductArgValid := proof.innerProductProof.Verify(aggParam.g, HPrime, uPrime, x.ToBytesS())
	if !innerProductArgValid {
		return false, errors.New("verify aggregated range proof statement 2 failed")
	}

	return true, nil
}

func (proof RangeProof) VerifyFasterUsingBase(anAssetTag *crypto.Point) (bool, error) {
	CACommitmentScheme.G[crypto.PedersenValueIndex] = anAssetTag
	numValue := len(proof.cmsValue)
	if numValue > utils.MaxOutputCoin {
		return false, errors.New("Must less than MaxOutputNumber")
	}
	numValuePad := roundUpPowTwo(numValue)
	maxExp := utils.MaxExp
	N := maxExp * numValuePad
	aggParam := setAggregateParams(N)

	cmsValue := proof.cmsValue
	for i := numValue; i < numValuePad; i++ {
		cmsValue = append(cmsValue, new(crypto.Point).Identity())
	}

	// recalculate challenge y, z
	y := generateChallenge(aggParam.cs.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	z := generateChallenge(y.ToBytesS(), []*crypto.Point{proof.a, proof.s})
	zSquare := new(crypto.Scalar).Mul(z, z)

	x := generateChallenge(z.ToBytesS(), []*crypto.Point{proof.t1, proof.t2})
	xSquare := new(crypto.Scalar).Mul(x, x)

	// g^tHat * h^tauX = V^(z^2) * g^delta(y,z) * T1^x * T2^(x^2)
	yVector := powerVector(y, N)
	deltaYZ, err := computeDeltaYZ(z, zSquare, yVector, N)
	if err != nil {
		return false, err
	}

	// Verify the first argument
	LHS := CACommitmentScheme.CommitAtIndex(proof.tHat, proof.tauX, crypto.PedersenValueIndex)
	RHS := new(crypto.Point).ScalarMult(proof.t2, xSquare)
	RHS.Add(RHS, new(crypto.Point).AddPedersen(deltaYZ, CACommitmentScheme.G[crypto.PedersenValueIndex], x, proof.t1))
	expVector := vectorMulScalar(powerVector(z, numValuePad), zSquare)
	RHS.Add(RHS, new(crypto.Point).MultiScalarMult(expVector, cmsValue))
	if !crypto.IsPointEqual(LHS, RHS) {
		return false, errors.New("verify aggregated range proof statement 1 failed")
	}

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
	// HPrime = H^(y^(1-i)
	HPrime := computeHPrime(y, N, aggParam.h)
	uPrime := new(crypto.Point).ScalarMult(aggParam.u, crypto.HashToScalar(x.ToBytesS()))
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
		return false, errors.New("verify aggregated range proof statement 2 failed")
	}

	return true, nil
}

// func VerifyBatchUsingBase(proofs []*RangeProof) (bool, error, int) {
// 	maxExp := utils.MaxExp
// 	baseG := CACommitmentScheme.G[crypto.PedersenValueIndex]
// 	baseH := CACommitmentScheme.G[crypto.PedersenRandomnessIndex]

// 	sum_tHat := new(crypto.Scalar).FromUint64(0)
// 	sum_tauX := new(crypto.Scalar).FromUint64(0)
// 	list_x_alpha := make([]*crypto.Scalar, 0)
// 	list_x_beta := make([]*crypto.Scalar, 0)
// 	list_xSquare := make([]*crypto.Scalar, 0)
// 	list_zSquare := make([]*crypto.Scalar, 0)

// 	list_t1 := make([]*crypto.Point, 0)
// 	list_t2 := make([]*crypto.Point, 0)
// 	list_V := make([]*crypto.Point, 0)

// 	sum_mu := new(crypto.Scalar).FromUint64(0)
// 	sum_absubthat := new(crypto.Scalar).FromUint64(0)

// 	list_S := make([]*crypto.Point, 0)
// 	list_A := make([]*crypto.Point, 0)
// 	list_beta := make([]*crypto.Scalar, 0)
// 	list_LR := make([]*crypto.Point, 0)
// 	list_lVector := make([]*crypto.Scalar, 0)
// 	list_rVector := make([]*crypto.Scalar, 0)
// 	list_gVector := make([]*crypto.Point, 0)
// 	list_hVector := make([]*crypto.Point, 0)

// 	twoNumber := new(crypto.Scalar).FromUint64(2)
// 	twoVectorN := powerVector(twoNumber, maxExp)

// 	for k, proof := range proofs {
// 		numValue := len(proof.cmsValue)
// 		if numValue > utils.MaxOutputCoin {
// 			return false, errors.New("Must less than MaxOutputNumber"), k
// 		}
// 		numValuePad := roundUpPowTwo(numValue)
// 		N := maxExp * numValuePad
// 		aggParam := setAggregateParams(N)

// 		cmsValue := proof.cmsValue
// 		for i := numValue; i < numValuePad; i++ {
// 			identity := new(crypto.Point).Identity()
// 			cmsValue = append(cmsValue, identity)
// 		}

// 		// recalculate challenge y, z, x
// 		y := generateChallenge(aggParam.cs.ToBytesS(), []*crypto.Point{proof.a, proof.s})
// 		z := generateChallenge(y.ToBytesS(), []*crypto.Point{proof.a, proof.s})
// 		x := generateChallenge(z.ToBytesS(), []*crypto.Point{proof.t1, proof.t2})
// 		zSquare := new(crypto.Scalar).Mul(z, z)
// 		xSquare := new(crypto.Scalar).Mul(x, x)

// 		// Random alpha and beta for batch equations check
// 		alpha := crypto.RandomScalar()
// 		beta := crypto.RandomScalar()
// 		list_beta = append(list_beta, beta)

// 		// Compute first equation check
// 		yVector := powerVector(y, N)
// 		deltaYZ, err := computeDeltaYZ(z, zSquare, yVector, N)
// 		if err != nil {
// 			return false, err, k
// 		}
// 		sum_tHat.Add(sum_tHat, new(crypto.Scalar).Mul(alpha, new(crypto.Scalar).Sub(proof.tHat, deltaYZ)))
// 		sum_tauX.Add(sum_tauX, new(crypto.Scalar).Mul(alpha, proof.tauX))

// 		list_x_alpha = append(list_x_alpha, new(crypto.Scalar).Mul(x, alpha))
// 		list_x_beta = append(list_x_beta, new(crypto.Scalar).Mul(x, beta))
// 		list_xSquare = append(list_xSquare, new(crypto.Scalar).Mul(xSquare, alpha))
// 		tmp := vectorMulScalar(powerVector(z, numValuePad), new(crypto.Scalar).Mul(zSquare, alpha))
// 		list_zSquare = append(list_zSquare, tmp...)

// 		list_V = append(list_V, cmsValue...)
// 		list_t1 = append(list_t1, proof.t1)
// 		list_t2 = append(list_t2, proof.t2)

// 		// Verify the second argument
// 		hashCache := x.ToBytesS()
// 		L := proof.innerProductProof.l
// 		R := proof.innerProductProof.r
// 		s := make([]*crypto.Scalar, N)
// 		sInverse := make([]*crypto.Scalar, N)
// 		logN := int(math.Log2(float64(N)))
// 		vSquareList := make([]*crypto.Scalar, logN)
// 		vInverseSquareList := make([]*crypto.Scalar, logN)

// 		for i := 0; i < N; i++ {
// 			s[i] = new(crypto.Scalar).Set(proof.innerProductProof.a)
// 			sInverse[i] = new(crypto.Scalar).Set(proof.innerProductProof.b)
// 		}

// 		for i := range L {
// 			v := generateChallenge(hashCache, []*crypto.Point{L[i], R[i]})
// 			hashCache = v.ToBytesS()
// 			vInverse := new(crypto.Scalar).Invert(v)
// 			vSquareList[i] = new(crypto.Scalar).Mul(v, v)
// 			vInverseSquareList[i] = new(crypto.Scalar).Mul(vInverse, vInverse)

// 			for j := 0; j < N; j++ {
// 				if j&int(math.Pow(2, float64(logN-i-1))) != 0 {
// 					s[j] = new(crypto.Scalar).Mul(s[j], v)
// 					sInverse[j] = new(crypto.Scalar).Mul(sInverse[j], vInverse)
// 				} else {
// 					s[j] = new(crypto.Scalar).Mul(s[j], vInverse)
// 					sInverse[j] = new(crypto.Scalar).Mul(sInverse[j], v)
// 				}
// 			}
// 		}

// 		lVector := make([]*crypto.Scalar, N)
// 		rVector := make([]*crypto.Scalar, N)

// 		vectorSum := make([]*crypto.Scalar, N)
// 		zTmp := new(crypto.Scalar).Set(z)
// 		for j := 0; j < numValuePad; j++ {
// 			zTmp.Mul(zTmp, z)
// 			for i := 0; i < maxExp; i++ {
// 				vectorSum[j*maxExp+i] = new(crypto.Scalar).Mul(twoVectorN[i], zTmp)
// 			}
// 		}
// 		yInverse := new(crypto.Scalar).Invert(y)
// 		yTmp := new(crypto.Scalar).Set(y)
// 		for j := 0; j < N; j++ {
// 			yTmp.Mul(yTmp, yInverse)
// 			lVector[j] = new(crypto.Scalar).Add(s[j], z)
// 			rVector[j] = new(crypto.Scalar).Sub(sInverse[j], vectorSum[j])
// 			rVector[j].Mul(rVector[j], yTmp)
// 			rVector[j].Sub(rVector[j], z)

// 			lVector[j].Mul(lVector[j], beta)
// 			rVector[j].Mul(rVector[j], beta)
// 		}

// 		list_lVector = append(list_lVector, lVector...)
// 		list_rVector = append(list_rVector, rVector...)

// 		tmp1 := new(crypto.Point).MultiScalarMult(vSquareList, L)
// 		tmp2 := new(crypto.Point).MultiScalarMult(vInverseSquareList, R)
// 		list_LR = append(list_LR, new(crypto.Point).Add(tmp1, tmp2))

// 		list_gVector = append(list_gVector, aggParam.g...)
// 		list_hVector = append(list_hVector, aggParam.h...)

// 		sum_mu.Add(sum_mu, new(crypto.Scalar).Mul(proof.mu, beta))
// 		ab := new(crypto.Scalar).Mul(proof.innerProductProof.a, proof.innerProductProof.b)
// 		absubthat := new(crypto.Scalar).Sub(ab, proof.tHat)
// 		absubthat.Mul(absubthat, crypto.HashToScalar(x.ToBytesS()))
// 		sum_absubthat.Add(sum_absubthat, new(crypto.Scalar).Mul(absubthat, beta))
// 		list_A = append(list_A, proof.a)
// 		list_S = append(list_S, proof.s)
// 	}

// 	tmp1 := new(crypto.Point).MultiScalarMult(list_lVector, list_gVector)
// 	tmp2 := new(crypto.Point).MultiScalarMult(list_rVector, list_hVector)
// 	tmp3 := new(crypto.Point).ScalarMult(AggParam.u, sum_absubthat)
// 	tmp4 := new(crypto.Point).ScalarMult(baseH, sum_mu)
// 	LHSPrime := new(crypto.Point).Add(tmp1, tmp2)
// 	LHSPrime.Add(LHSPrime, tmp3)
// 	LHSPrime.Add(LHSPrime, tmp4)

// 	LHS := new(crypto.Point).AddPedersen(sum_tHat, baseG, sum_tauX, baseH)
// 	LHSPrime.Add(LHSPrime, LHS)

// 	tmp5 := new(crypto.Point).MultiScalarMult(list_beta, list_A)
// 	tmp6 := new(crypto.Point).MultiScalarMult(list_x_beta, list_S)
// 	RHSPrime := new(crypto.Point).Add(tmp5, tmp6)
// 	RHSPrime.Add(RHSPrime, new(crypto.Point).MultiScalarMult(list_beta, list_LR))

// 	part1 := new(crypto.Point).MultiScalarMult(list_x_alpha, list_t1)
// 	part2 := new(crypto.Point).MultiScalarMult(list_xSquare, list_t2)
// 	RHS := new(crypto.Point).Add(part1, part2)
// 	RHS.Add(RHS, new(crypto.Point).MultiScalarMult(list_zSquare, list_V))
// 	RHSPrime.Add(RHSPrime, RHS)
// 	//fmt.Println("Batch Verification ", LHSPrime)
// 	//fmt.Println("Batch Verification ", RHSPrime)

// 	if !crypto.IsPointEqual(LHSPrime, RHSPrime) {
// 		Logger.Log.Errorf("batch verify aggregated range proof failed")
// 		return false, errors.New("batch verify aggregated range proof failed"), -1
// 	}
// 	return true, nil, -1
// }

func TransformWitnessToCAWitness(wit *Witness, assetTagBlinders []*crypto.Scalar) (*Witness,error){
	if len(assetTagBlinders)!=len(wit.values) || len(assetTagBlinders)!=len(wit.rands){
		return nil, errors.New("Cannot transform witness. Parameter lengths mismatch")
	}
	newRands := make([]*crypto.Scalar, len(wit.values))

	for i,_ := range wit.values{
		temp := new(crypto.Scalar).Sub(assetTagBlinders[i], assetTagBlinders[0])
		temp.Mul(temp, new(crypto.Scalar).FromUint64(wit.values[i]))
		temp.Add(temp, wit.rands[i])
		newRands[i] = temp
	}
	result := new(Witness)
	result.Set(wit.values, newRands)
	return result, nil
}