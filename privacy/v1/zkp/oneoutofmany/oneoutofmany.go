package oneoutofmany

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	privacyUtils "github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/oneoutofmany/polynomial"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/utils"
	"math/big"
)

// This protocol proves in zero-knowledge that one-out-of-N commitments contains 0

// OneOutOfManyStatement represents the statement of an OneOutOfManyProof.
type OneOutOfManyStatement struct {
	Commitments []*crypto.Point
}

// OneOutOfManyWitness represents a witness to an OneOutOfManyProof.
type OneOutOfManyWitness struct {
	stmt        *OneOutOfManyStatement
	rand        *crypto.Scalar
	indexIsZero uint64
}

// OneOutOfManyProof represents an one-of-many proof proving that among a list of commitments, there exists one
// commitment that commits to 0. It is used to hide the real input coin among a list of input coins.
type OneOutOfManyProof struct {
	Statement      *OneOutOfManyStatement
	cl, ca, cb, cd []*crypto.Point
	f, za, zb      []*crypto.Scalar
	zd             *crypto.Scalar
}

// Init creates an empty OneOutOfManyProof.
func (proof *OneOutOfManyProof) Init() *OneOutOfManyProof {
	proof.zd = new(crypto.Scalar)
	proof.Statement = new(OneOutOfManyStatement)

	return proof
}

// Set sets data to an OneOutOfManyStatement.
func (stmt *OneOutOfManyStatement) Set(v []*crypto.Point) {
	stmt.Commitments = v
}

// Set sets data to an OneOutOfManyWitness.
func (wit *OneOutOfManyWitness) Set(commitments []*crypto.Point, rand *crypto.Scalar, indexIsZero uint64) {
	wit.stmt = new(OneOutOfManyStatement)
	wit.stmt.Set(commitments)

	wit.indexIsZero = indexIsZero
	wit.rand = rand
}

// Set sets data to an OneOutOfManyProof.
func (proof *OneOutOfManyProof) Set(
	commitments []*crypto.Point,
	cl, ca, cb, cd []*crypto.Point,
	f, za, zb []*crypto.Scalar,
	zd *crypto.Scalar) {

	proof.Statement = new(OneOutOfManyStatement)
	proof.Statement.Set(commitments)

	proof.cl, proof.ca, proof.cb, proof.cd = cl, ca, cb, cd
	proof.f, proof.za, proof.zb = f, za, zb
	proof.zd = zd
}

// Bytes returns the byte-representation of an OneOutOfManyProof.
func (proof OneOutOfManyProof) Bytes() []byte {
	// if proof is nil, return an empty array
	if proof.isNil() {
		return []byte{}
	}

	// N = 2^n
	n := privacyUtils.CommitmentRingSizeExp

	var bytes []byte

	// convert array cl to bytes array
	for i := 0; i < n; i++ {
		bytes = append(bytes, proof.cl[i].ToBytesS()...)
	}
	// convert array ca to bytes array
	for i := 0; i < n; i++ {
		//fmt.Printf("proof.ca[i]: %v\n", proof.ca[i])
		//fmt.Printf("proof.ca[i]: %v\n", proof.ca[i].Compress())
		bytes = append(bytes, proof.ca[i].ToBytesS()...)
	}

	// convert array cb to bytes array
	for i := 0; i < n; i++ {
		bytes = append(bytes, proof.cb[i].ToBytesS()...)
	}

	// convert array cd to bytes array
	for i := 0; i < n; i++ {
		bytes = append(bytes, proof.cd[i].ToBytesS()...)
	}

	// convert array f to bytes array
	for i := 0; i < n; i++ {
		bytes = append(bytes, proof.f[i].ToBytesS()...)
	}

	// convert array za to bytes array
	for i := 0; i < n; i++ {
		bytes = append(bytes, proof.za[i].ToBytesS()...)
	}

	// convert array zb to bytes array
	for i := 0; i < n; i++ {
		bytes = append(bytes, proof.zb[i].ToBytesS()...)
	}

	// convert array zd to bytes array
	bytes = append(bytes, proof.zd.ToBytesS()...)

	return bytes
}

// SetBytes set raw-data to an OneOutOfManyProof.
func (proof *OneOutOfManyProof) SetBytes(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	n := privacyUtils.CommitmentRingSizeExp

	offset := 0
	var err error

	// get cl array
	proof.cl = make([]*crypto.Point, n)
	for i := 0; i < n; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("one-of-many proof byte unmarshalling failed")
		}
		proof.cl[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset = offset + crypto.Ed25519KeySize
	}

	// get ca array
	proof.ca = make([]*crypto.Point, n)
	for i := 0; i < n; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("one-of-many proof byte unmarshalling failed")
		}
		proof.ca[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset = offset + crypto.Ed25519KeySize
	}

	// get cb array
	proof.cb = make([]*crypto.Point, n)
	for i := 0; i < n; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("one-of-many proof byte unmarshalling failed")
		}
		proof.cb[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset = offset + crypto.Ed25519KeySize
	}

	// get cd array
	proof.cd = make([]*crypto.Point, n)
	for i := 0; i < n; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("one-of-many proof byte unmarshalling failed")
		}
		proof.cd[i], err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset = offset + crypto.Ed25519KeySize
	}

	// get f array
	proof.f = make([]*crypto.Scalar, n)
	for i := 0; i < n; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("one-of-many proof byte unmarshalling failed")
		}
		proof.f[i] = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		offset = offset + crypto.Ed25519KeySize
	}

	// get za array
	proof.za = make([]*crypto.Scalar, n)
	for i := 0; i < n; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("one-of-many proof byte unmarshalling failed")
		}
		proof.za[i] = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		offset = offset + crypto.Ed25519KeySize
	}

	// get zb array
	proof.zb = make([]*crypto.Scalar, n)
	for i := 0; i < n; i++ {
		if offset+crypto.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("one-of-many proof byte unmarshalling failed")
		}
		proof.zb[i] = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
		offset = offset + crypto.Ed25519KeySize
	}

	// get zd
	if offset+crypto.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("one-of-many proof byte unmarshalling failed")
	}
	proof.zd = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])

	return nil
}

// Prove produces an OneOutOfManyProof from an OneOutOfManyWitness.
func (wit OneOutOfManyWitness) Prove() (*OneOutOfManyProof, error) {
	// Check the number of Commitment list's elements
	N := len(wit.stmt.Commitments)
	if N != privacyUtils.CommitmentRingSize {
		return nil, fmt.Errorf("the number of Commitment list's elements must be equal to CMRingSize")
	}

	n := privacyUtils.CommitmentRingSizeExp

	// Check indexIsZero
	if wit.indexIsZero > uint64(N) {
		return nil, fmt.Errorf("index is zero must be Index in list of commitments")
	}

	// represent indexIsZero in binary
	indexIsZeroBinary := privacyUtils.ConvertIntToBinary(int(wit.indexIsZero), n)

	//
	r := make([]*crypto.Scalar, n)
	a := make([]*crypto.Scalar, n)
	s := make([]*crypto.Scalar, n)
	t := make([]*crypto.Scalar, n)
	u := make([]*crypto.Scalar, n)

	cl := make([]*crypto.Point, n)
	ca := make([]*crypto.Point, n)
	cb := make([]*crypto.Point, n)
	cd := make([]*crypto.Point, n)

	for j := 0; j < n; j++ {
		// Generate random numbers
		r[j] = crypto.RandomScalar()
		a[j] = crypto.RandomScalar()
		s[j] = crypto.RandomScalar()
		t[j] = crypto.RandomScalar()
		u[j] = crypto.RandomScalar()

		// convert indexIsZeroBinary[j] to crypto.Scalar
		indexInt := new(crypto.Scalar).FromUint64(uint64(indexIsZeroBinary[j]))

		// Calculate cl, ca, cb, cd
		// cl = Com(l, r)
		cl[j] = crypto.PedCom.CommitAtIndex(indexInt, r[j], crypto.PedersenPrivateKeyIndex)

		// ca = Com(a, s)
		ca[j] = crypto.PedCom.CommitAtIndex(a[j], s[j], crypto.PedersenPrivateKeyIndex)

		// cb = Com(la, t)
		la := new(crypto.Scalar).Mul(indexInt, a[j])
		//la.Mod(la, crypto.Curve.Params().N)
		cb[j] = crypto.PedCom.CommitAtIndex(la, t[j], crypto.PedersenPrivateKeyIndex)
	}

	// Calculate: cd_k = ci^pi,k
	for k := 0; k < n; k++ {
		// Calculate pi,k which is coefficient of x^k in polynomial pi(x)
		cd[k] = new(crypto.Point).Identity()

		for i := 0; i < N; i++ {
			iBinary := privacyUtils.ConvertIntToBinary(i, n)
			pik := getCoefficient(iBinary, k, n, a, indexIsZeroBinary)
			cd[k].Add(cd[k], new(crypto.Point).ScalarMult(wit.stmt.Commitments[i], pik))
		}

		cd[k].Add(cd[k], crypto.PedCom.CommitAtIndex(new(crypto.Scalar).FromUint64(0), u[k], crypto.PedersenPrivateKeyIndex))
	}

	// Calculate x
	commitmentsInBytes := make([][]byte, 0)
	for _, commitment := range wit.stmt.Commitments {
		commitmentsInBytes = append(commitmentsInBytes, commitment.ToBytesS())
	}
	x := utils.GenerateChallenge(commitmentsInBytes)
	for j := 0; j < n; j++ {
		x = utils.GenerateChallenge([][]byte{
			x.ToBytesS(),
			cl[j].ToBytesS(),
			ca[j].ToBytesS(),
			cb[j].ToBytesS(),
			cd[j].ToBytesS(),
		})
	}

	// Calculate za, zb zd
	za := make([]*crypto.Scalar, n)
	zb := make([]*crypto.Scalar, n)

	f := make([]*crypto.Scalar, n)

	for j := 0; j < n; j++ {
		// f = lx + a
		f[j] = new(crypto.Scalar).Mul(new(crypto.Scalar).FromUint64(uint64(indexIsZeroBinary[j])), x)
		f[j].Add(f[j], a[j])

		// za = s + rx
		za[j] = new(crypto.Scalar).Mul(r[j], x)
		za[j].Add(za[j], s[j])

		// zb = r(x - f) + t
		zb[j] = new(crypto.Scalar).Sub(x, f[j])
		zb[j].Mul(zb[j], r[j])
		zb[j].Add(zb[j], t[j])
	}

	// zd = rand * x^n - sum_{k=0}^{n-1} u[k] * x^k
	xi := new(crypto.Scalar).FromUint64(1)
	sum := new(crypto.Scalar).FromUint64(0)
	for k := 0; k < n; k++ {
		tmp := new(crypto.Scalar).Mul(xi, u[k])
		sum.Add(sum, tmp)
		xi.Mul(xi, x)
	}
	zd := new(crypto.Scalar).Mul(xi, wit.rand)
	zd.Sub(zd, sum)

	proof := new(OneOutOfManyProof).Init()
	proof.Set(wit.stmt.Commitments, cl, ca, cb, cd, f, za, zb, zd)

	return proof, nil
}

func (proof OneOutOfManyProof) isNil() bool {
	if proof.cl == nil {
		return true
	}
	if proof.ca == nil {
		return true
	}
	if proof.cb == nil {
		return true
	}
	if proof.cd == nil {
		return true
	}
	if proof.f == nil {
		return true
	}
	if proof.za == nil {
		return true
	}
	if proof.zb == nil {
		return true
	}
	return proof.zd == nil
}

// getCoefficient gets the coefficients of x^k in the polynomial p_i(x).
func getCoefficient(iBinary []byte, k int, n int, scLs []*crypto.Scalar, l []byte) *crypto.Scalar {

	a := make([]*big.Int, len(scLs))
	for i := 0; i < len(scLs); i++ {
		a[i] = privacyUtils.ScalarToBigInt(scLs[i])
	}

	//AP2
	curveOrder := polynomial.LInt
	res := polynomial.Poly{big.NewInt(1)}
	var fji polynomial.Poly
	for j := n - 1; j >= 0; j-- {
		fj := polynomial.Poly{a[j], big.NewInt(int64(l[j]))}
		if iBinary[j] == 0 {
			fji = polynomial.Poly{big.NewInt(0), big.NewInt(1)}.Sub(fj, curveOrder)
		} else {
			fji = fj
		}
		res = res.Mul(fji, curveOrder)
	}

	var sc2 *crypto.Scalar
	if res.GetDegree() < k {
		sc2 = new(crypto.Scalar).FromUint64(0)
	} else {
		sc2 = privacyUtils.BigIntToScalar(res[k])
	}
	return sc2
}
