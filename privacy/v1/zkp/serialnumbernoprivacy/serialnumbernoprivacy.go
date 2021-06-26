package serialnumbernoprivacy

import (
	"errors"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/utils"
)

type SerialNumberNoPrivacyStatement struct {
	output *crypto.Point
	vKey   *crypto.Point
	input  *crypto.Scalar
}

// SNNoPrivacyWitness is a protocol for Zero-knowledge Proof of Knowledge of one out of many commitments containing 0
// include Witness: CommitedValue, r []byte
type SNNoPrivacyWitness struct {
	stmt SerialNumberNoPrivacyStatement
	seed *crypto.Scalar
}

// serialNumberNNoPrivacyProof contains Proof's value
type SNNoPrivacyProof struct {
	// general info
	stmt SerialNumberNoPrivacyStatement

	tSeed   *crypto.Point
	tOutput *crypto.Point

	zSeed *crypto.Scalar
}

func (proof SNNoPrivacyProof) GetVKey() *crypto.Point {
	return proof.stmt.vKey
}

func (proof SNNoPrivacyProof) GetOutput() *crypto.Point {
	return proof.stmt.output
}

func (proof SNNoPrivacyProof) GetInput() *crypto.Scalar {
	return proof.stmt.input
}

func (proof SNNoPrivacyProof) ValidateSanity() bool {
	if !proof.stmt.output.PointValid() {
		return false
	}
	if !proof.stmt.vKey.PointValid() {
		return false
	}
	if !proof.stmt.input.ScalarValid() {
		return false
	}

	if !proof.tSeed.PointValid() {
		return false
	}
	if !proof.tOutput.PointValid() {
		return false
	}
	return proof.zSeed.ScalarValid()
}

func (pro SNNoPrivacyProof) isNil() bool {
	if pro.stmt.output == nil {
		return true
	}
	if pro.stmt.vKey == nil {
		return true
	}
	if pro.stmt.input == nil {
		return true
	}
	if pro.tSeed == nil {
		return true
	}
	if pro.tOutput == nil {
		return true
	}
	if pro.zSeed == nil {
		return true
	}
	return false
}

func (pro *SNNoPrivacyProof) Init() *SNNoPrivacyProof {
	pro.stmt.output = new(crypto.Point)
	pro.stmt.vKey = new(crypto.Point)
	pro.stmt.input = new(crypto.Scalar)

	pro.tSeed = new(crypto.Point)
	pro.tOutput = new(crypto.Point)

	pro.zSeed = new(crypto.Scalar)

	return pro
}

// Set sets Witness
func (wit *SNNoPrivacyWitness) Set(
	output *crypto.Point,
	vKey *crypto.Point,
	input *crypto.Scalar,
	seed *crypto.Scalar) {

	if wit == nil {
		wit = new(SNNoPrivacyWitness)
	}

	wit.stmt.output = output
	wit.stmt.vKey = vKey
	wit.stmt.input = input

	wit.seed = seed
}

// Set sets Proof
func (pro *SNNoPrivacyProof) Set(
	output *crypto.Point,
	vKey *crypto.Point,
	input *crypto.Scalar,
	tSeed *crypto.Point,
	tOutput *crypto.Point,
	zSeed *crypto.Scalar) {

	if pro == nil {
		pro = new(SNNoPrivacyProof)
	}

	pro.stmt.output = output
	pro.stmt.vKey = vKey
	pro.stmt.input = input

	pro.tSeed = tSeed
	pro.tOutput = tOutput

	pro.zSeed = zSeed
}

func (pro SNNoPrivacyProof) Bytes() []byte {
	// if proof is nil, return an empty array
	if pro.isNil() {
		return []byte{}
	}

	var bytes []byte
	bytes = append(bytes, pro.stmt.output.ToBytesS()...)
	bytes = append(bytes, pro.stmt.vKey.ToBytesS()...)
	bytes = append(bytes, pro.stmt.input.ToBytesS()...)

	bytes = append(bytes, pro.tSeed.ToBytesS()...)
	bytes = append(bytes, pro.tOutput.ToBytesS()...)

	bytes = append(bytes, pro.zSeed.ToBytesS()...)

	return bytes
}

func (pro *SNNoPrivacyProof) SetBytes(bytes []byte) error {
	// if len(bytes) == 0 {
	// 	return errors.New("Bytes array is empty")
	// }
	if len(bytes) < crypto.Ed25519KeySize*6 {
		return errors.New("Not enough bytes to unmarshal Serial Number No Privacy Proof")
	}

	offset := 0
	var err error
	pro.stmt.output, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	pro.stmt.vKey, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	pro.stmt.input.FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	offset += crypto.Ed25519KeySize

	pro.tSeed, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	pro.tOutput, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	pro.zSeed = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])

	return nil
}

func (wit SNNoPrivacyWitness) Prove(mess []byte) (*SNNoPrivacyProof, error) {
	// randomness
	eSK := crypto.RandomScalar()

	// calculate tSeed = g_SK^eSK
	tSK := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], eSK)

	// calculate tOutput = sn^eSK
	tE := new(crypto.Point).ScalarMult(wit.stmt.output, eSK)

	x := new(crypto.Scalar)
	if mess == nil {
		// calculate x = hash(tSeed || tInput || tSND2 || tOutput)
		// recheck frombytes is valid scalar
		x = utils.GenerateChallenge([][]byte{wit.stmt.output.ToBytesS(), wit.stmt.vKey.ToBytesS(), tSK.ToBytesS(), tE.ToBytesS()})
	} else {
		x.FromBytesS(mess)
	}

	// Calculate zSeed = SK * x + eSK
	zSK := new(crypto.Scalar).Mul(wit.seed, x)
	zSK.Add(zSK, eSK)

	proof := new(SNNoPrivacyProof).Init()
	proof.Set(wit.stmt.output, wit.stmt.vKey, wit.stmt.input, tSK, tE, zSK)
	return proof, nil
}

func (pro SNNoPrivacyProof) Verify(mess []byte) (bool, error) {
	// re-calculate x = hash(tSeed || tOutput)
	x := new(crypto.Scalar)
	if mess == nil {
		// calculate x = hash(tSeed || tInput || tSND2 || tOutput)
		x = utils.GenerateChallenge([][]byte{pro.stmt.output.ToBytesS(), pro.stmt.vKey.ToBytesS(), pro.tSeed.ToBytesS(), pro.tOutput.ToBytesS()})
	} else {
		x.FromBytesS(mess)
	}

	// Check gSK^zSeed = vKey^x * tSeed
	leftPoint1 := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], pro.zSeed)

	rightPoint1 := new(crypto.Point).ScalarMult(pro.stmt.vKey, x)
	rightPoint1 = rightPoint1.Add(rightPoint1, pro.tSeed)

	if !crypto.IsPointEqual(leftPoint1, rightPoint1) {
		return false, errors.New("verify serial number no privacy proof statement 1 failed")
	}

	// Check sn^(zSeed + x*input) = gSK^x * tOutput
	tmp := new(crypto.Scalar).Add(pro.zSeed, new(crypto.Scalar).Mul(x, pro.stmt.input))
	leftPoint2 := new(crypto.Point).ScalarMult(pro.stmt.output, tmp)

	rightPoint2 := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], x)
	rightPoint2 = rightPoint2.Add(rightPoint2, pro.tOutput)

	if !crypto.IsPointEqual(leftPoint2, rightPoint2) {
		return false, errors.New("verify serial number no privacy proof statement 2 failed")
	}

	return true, nil
}

func (pro SNNoPrivacyProof) VerifyOld(mess []byte) (bool, error) {
	// re-calculate x = hash(tSeed || tOutput)
	x := new(crypto.Scalar)
	if mess == nil {
		// calculate x = hash(tSeed || tInput || tSND2 || tOutput)
		x = utils.GenerateChallenge([][]byte{pro.tSeed.ToBytesS(), pro.tOutput.ToBytesS()})
	} else {
		x.FromBytesS(mess)
	}

	// Check gSK^zSeed = vKey^x * tSeed
	leftPoint1 := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], pro.zSeed)

	rightPoint1 := new(crypto.Point).ScalarMult(pro.stmt.vKey, x)
	rightPoint1 = rightPoint1.Add(rightPoint1, pro.tSeed)

	if !crypto.IsPointEqual(leftPoint1, rightPoint1) {
		return false, errors.New("verifyOld serial number no privacy proof statement 1 failed")
	}

	// Check sn^(zSeed + x*input) = gSK^x * tOutput
	tmp := new(crypto.Scalar).Add(pro.zSeed, new(crypto.Scalar).Mul(x, pro.stmt.input))
	leftPoint2 := new(crypto.Point).ScalarMult(pro.stmt.output, tmp)

	rightPoint2 := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], x)
	rightPoint2 = rightPoint2.Add(rightPoint2, pro.tOutput)

	if !crypto.IsPointEqual(leftPoint2, rightPoint2) {
		return false, errors.New("verifyOld serial number no privacy proof statement 2 failed")
	}

	return true, nil
}
