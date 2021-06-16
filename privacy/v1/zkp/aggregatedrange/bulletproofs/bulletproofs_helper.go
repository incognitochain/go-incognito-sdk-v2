package bulletproofs

import (
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
	"github.com/pkg/errors"
)

// ConvertIntToBinary represents a integer number in binary
func ConvertUint64ToBinary(number uint64, n int) []*crypto.Scalar {
	if number == 0 {
		res := make([]*crypto.Scalar, n)
		for i := 0; i < n; i++ {
			res[i] = new(crypto.Scalar).FromUint64(0)
		}
		return res
	}

	binary := make([]*crypto.Scalar, n)

	for i := 0; i < n; i++ {
		binary[i] = new(crypto.Scalar).FromUint64(number % 2)
		number = number / 2
	}
	return binary
}

func computeHPrime(y *crypto.Scalar, N int, H []*crypto.Point) []*crypto.Point {
	yInverse := new(crypto.Scalar).Invert(y)
	HPrime := make([]*crypto.Point, N)
	expyInverse := new(crypto.Scalar).FromUint64(1)
	for i := 0; i < N; i++ {
		HPrime[i] = new(crypto.Point).ScalarMult(H[i], expyInverse)
		expyInverse.Mul(expyInverse, yInverse)
	}
	return HPrime
}

func computeDeltaYZ(z, zSquare *crypto.Scalar, yVector []*crypto.Scalar, N int) (*crypto.Scalar, error) {
	oneNumber := new(crypto.Scalar).FromUint64(1)
	twoNumber := new(crypto.Scalar).FromUint64(2)
	oneVectorN := powerVector(oneNumber, utils.MaxExp)
	twoVectorN := powerVector(twoNumber, utils.MaxExp)
	oneVector := powerVector(oneNumber, N)

	deltaYZ := new(crypto.Scalar).Sub(z, zSquare)
	// ip1 = <1^(n*m), y^(n*m)>
	var ip1, ip2 *crypto.Scalar
	var err error
	if ip1, err = innerProduct(oneVector, yVector); err != nil {
		return nil, err
	} else if ip2, err = innerProduct(oneVectorN, twoVectorN); err != nil {
		return nil, err
	} else {
		deltaYZ.Mul(deltaYZ, ip1)
		sum := new(crypto.Scalar).FromUint64(0)
		zTmp := new(crypto.Scalar).Set(zSquare)
		for j := 0; j < int(N/utils.MaxExp); j++ {
			zTmp.Mul(zTmp, z)
			sum.Add(sum, zTmp)
		}
		sum.Mul(sum, ip2)
		deltaYZ.Sub(deltaYZ, sum)
	}
	return deltaYZ, nil
}

func innerProduct(a []*crypto.Scalar, b []*crypto.Scalar) (*crypto.Scalar, error) {
	if len(a) != len(b) {
		return nil, errors.New("Incompatible sizes of a and b")
	}
	result := new(crypto.Scalar).FromUint64(uint64(0))
	for i := range a {
		//res = a[i]*b[i] + res % l
		result.MulAdd(a[i], b[i], result)
	}
	return result, nil
}

func vectorAdd(a []*crypto.Scalar, b []*crypto.Scalar) ([]*crypto.Scalar, error) {
	if len(a) != len(b) {
		return nil, errors.New("Incompatible sizes of a and b")
	}
	result := make([]*crypto.Scalar, len(a))
	for i := range a {
		result[i] = new(crypto.Scalar).Add(a[i], b[i])
	}
	return result, nil
}

func setAggregateParams(N int) *bulletproofParams {
	aggParam := new(bulletproofParams)
	aggParam.g = AggParam.g[0:N]
	aggParam.h = AggParam.h[0:N]
	aggParam.u = AggParam.u
	aggParam.cs = AggParam.cs
	return aggParam
}

func roundUpPowTwo(v int) int {
	if v == 0 {
		return 1
	} else {
		v--
		v |= v >> 1
		v |= v >> 2
		v |= v >> 4
		v |= v >> 8
		v |= v >> 16
		v++
		return v
	}
}

func hadamardProduct(a []*crypto.Scalar, b []*crypto.Scalar) ([]*crypto.Scalar, error) {
	if len(a) != len(b) {
		return nil, errors.New("Invalid input")
	}
	result := make([]*crypto.Scalar, len(a))
	for i := 0; i < len(result); i++ {
		result[i] = new(crypto.Scalar).Mul(a[i], b[i])
	}
	return result, nil
}

// powerVector calculates base^n
func powerVector(base *crypto.Scalar, n int) []*crypto.Scalar {
	result := make([]*crypto.Scalar, n)
	result[0] = new(crypto.Scalar).FromUint64(1)
	if n > 1 {
		result[1] = new(crypto.Scalar).Set(base)
		for i := 2; i < n; i++ {
			result[i] = new(crypto.Scalar).Mul(result[i-1], base)
		}
	}
	return result
}

// vectorAddScalar adds a vector to a big int, returns big int array
func vectorAddScalar(v []*crypto.Scalar, s *crypto.Scalar) []*crypto.Scalar {
	result := make([]*crypto.Scalar, len(v))
	for i := range v {
		result[i] = new(crypto.Scalar).Add(v[i], s)
	}
	return result
}

// vectorMulScalar mul a vector to a big int, returns a vector
func vectorMulScalar(v []*crypto.Scalar, s *crypto.Scalar) []*crypto.Scalar {
	result := make([]*crypto.Scalar, len(v))
	for i := range v {
		result[i] = new(crypto.Scalar).Mul(v[i], s)
	}
	return result
}

// CommitAll commits a list of PCM_CAPACITY value(s)
func encodeVectors(l []*crypto.Scalar, r []*crypto.Scalar, g []*crypto.Point, h []*crypto.Point) (*crypto.Point, error) {
	if len(l) != len(r) || len(g) != len(l) || len(h) != len(g) {
		return nil, errors.New("Invalid input")
	}
	tmp1 := new(crypto.Point).MultiScalarMult(l, g)
	tmp2 := new(crypto.Point).MultiScalarMult(r, h)
	res := new(crypto.Point).Add(tmp1, tmp2)
	return res, nil
}

// bulletproofParams includes all generator for aggregated range proof
func newBulletproofParams(m int) *bulletproofParams {
	maxExp := utils.MaxExp
	numCommitValue := utils.NumBase
	maxOutputCoin := utils.MaxOutputCoin
	capacity := maxExp * m // fixed value
	param := new(bulletproofParams)
	param.g = make([]*crypto.Point, capacity)
	param.h = make([]*crypto.Point, capacity)
	csByte := []byte{}

	for i := 0; i < capacity; i++ {
		param.g[i] = crypto.HashToPointFromIndex(int64(numCommitValue+i), crypto.CStringBulletProof)
		param.h[i] = crypto.HashToPointFromIndex(int64(numCommitValue+i+maxOutputCoin*maxExp), crypto.CStringBulletProof)
		csByte = append(csByte, param.g[i].ToBytesS()...)
		csByte = append(csByte, param.h[i].ToBytesS()...)
	}

	param.u = new(crypto.Point)
	param.u = crypto.HashToPointFromIndex(int64(numCommitValue+2*maxOutputCoin*maxExp), crypto.CStringBulletProof)
	csByte = append(csByte, param.u.ToBytesS()...)

	param.cs = crypto.HashToPoint(csByte)
	return param
}

func generateChallenge(hashCache []byte, values []*crypto.Point) *crypto.Scalar {
	bytes := []byte{}
	bytes = append(bytes, hashCache...)
	for i := 0; i < len(values); i++ {
		bytes = append(bytes, values[i].ToBytesS()...)
	}
	hash := crypto.HashToScalar(bytes)
	return hash
}
