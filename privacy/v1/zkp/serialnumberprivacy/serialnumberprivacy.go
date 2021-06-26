package serialnumberprivacy

import (
	"errors"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/utils"
)

type SerialNumberPrivacyStatement struct {
	sn       *crypto.Point // serial number
	comSK    *crypto.Point // commitment to private key
	comInput *crypto.Point // commitment to input of the pseudo-random function
}

type SNPrivacyWitness struct {
	stmt *SerialNumberPrivacyStatement // statement to be proved

	sk     *crypto.Scalar // private key
	rSK    *crypto.Scalar // blinding factor in the commitment to private key
	input  *crypto.Scalar // input of pseudo-random function
	rInput *crypto.Scalar // blinding factor in the commitment to input
}

type SNPrivacyProof struct {
	stmt *SerialNumberPrivacyStatement // statement to be proved

	tSK    *crypto.Point // random commitment related to private key
	tInput *crypto.Point // random commitment related to input
	tSN    *crypto.Point // random commitment related to serial number

	zSK     *crypto.Scalar // first challenge-dependent information to open the commitment to private key
	zRSK    *crypto.Scalar // second challenge-dependent information to open the commitment to private key
	zInput  *crypto.Scalar // first challenge-dependent information to open the commitment to input
	zRInput *crypto.Scalar // second challenge-dependent information to open the commitment to input
}

// ValidateSanity validates sanity of proof
func (proof SNPrivacyProof) ValidateSanity() bool {
	if !proof.stmt.sn.PointValid() {
		return false
	}
	if !proof.stmt.comSK.PointValid() {
		return false
	}
	if !proof.stmt.comInput.PointValid() {
		return false
	}
	if !proof.tSK.PointValid() {
		return false
	}
	if !proof.tInput.PointValid() {
		return false
	}
	if !proof.tSN.PointValid() {
		return false
	}
	if !proof.zSK.ScalarValid() {
		return false
	}
	if !proof.zRSK.ScalarValid() {
		return false
	}
	if !proof.zInput.ScalarValid() {
		return false
	}
	if !proof.zRInput.ScalarValid() {
		return false
	}
	return true
}

func (proof SNPrivacyProof) isNil() bool {
	if proof.stmt.sn == nil {
		return true
	}
	if proof.stmt.comSK == nil {
		return true
	}
	if proof.stmt.comInput == nil {
		return true
	}
	if proof.tSK == nil {
		return true
	}
	if proof.tInput == nil {
		return true
	}
	if proof.tSN == nil {
		return true
	}
	if proof.zSK == nil {
		return true
	}
	if proof.zRSK == nil {
		return true
	}
	if proof.zInput == nil {
		return true
	}
	return proof.zRInput == nil
}

// Init inits Proof
func (proof *SNPrivacyProof) Init() *SNPrivacyProof {
	proof.stmt = new(SerialNumberPrivacyStatement)

	proof.tSK = new(crypto.Point)
	proof.tInput = new(crypto.Point)
	proof.tSN = new(crypto.Point)

	proof.zSK = new(crypto.Scalar)
	proof.zRSK = new(crypto.Scalar)
	proof.zInput = new(crypto.Scalar)
	proof.zRInput = new(crypto.Scalar)

	return proof
}

func (proof SNPrivacyProof) GetComSK() *crypto.Point {
	return proof.stmt.comSK
}

func (proof SNPrivacyProof) GetSN() *crypto.Point {
	return proof.stmt.sn
}

func (proof SNPrivacyProof) GetComInput() *crypto.Point {
	return proof.stmt.comInput
}



// Set sets Statement
func (stmt *SerialNumberPrivacyStatement) Set(
	SN *crypto.Point,
	comSK *crypto.Point,
	comInput *crypto.Point) {
	stmt.sn = SN
	stmt.comSK = comSK
	stmt.comInput = comInput
}

// Set sets Witness
func (wit *SNPrivacyWitness) Set(
	stmt *SerialNumberPrivacyStatement,
	SK *crypto.Scalar,
	rSK *crypto.Scalar,
	input *crypto.Scalar,
	rInput *crypto.Scalar) {

	wit.stmt = stmt
	wit.sk = SK
	wit.rSK = rSK
	wit.input = input
	wit.rInput = rInput
}

// Set sets Proof
func (proof *SNPrivacyProof) Set(
	stmt *SerialNumberPrivacyStatement,
	tSK *crypto.Point,
	tInput *crypto.Point,
	tSN *crypto.Point,
	zSK *crypto.Scalar,
	zRSK *crypto.Scalar,
	zInput *crypto.Scalar,
	zRInput *crypto.Scalar) {
	proof.stmt = stmt
	proof.tSK = tSK
	proof.tInput = tInput
	proof.tSN = tSN

	proof.zSK = zSK
	proof.zRSK = zRSK
	proof.zInput = zInput
	proof.zRInput = zRInput
}

func (proof SNPrivacyProof) Bytes() []byte {
	// if proof is nil, return an empty array
	if proof.isNil() {
		return []byte{}
	}

	var bytes []byte
	bytes = append(bytes, proof.stmt.sn.ToBytesS()...)
	bytes = append(bytes, proof.stmt.comSK.ToBytesS()...)
	bytes = append(bytes, proof.stmt.comInput.ToBytesS()...)

	bytes = append(bytes, proof.tSK.ToBytesS()...)
	bytes = append(bytes, proof.tInput.ToBytesS()...)
	bytes = append(bytes, proof.tSN.ToBytesS()...)

	bytes = append(bytes, proof.zSK.ToBytesS()...)
	bytes = append(bytes, proof.zRSK.ToBytesS()...)
	bytes = append(bytes, proof.zInput.ToBytesS()...)
	bytes = append(bytes, proof.zRInput.ToBytesS()...)

	return bytes
}

func (proof *SNPrivacyProof) SetBytes(bytes []byte) error {
	if len(bytes) == 0 {
		return errors.New("Bytes array is empty")
	}
	if len(bytes) < 9*crypto.Ed25519KeySize {
		return fmt.Errorf("not enough bytes to unmarshal Serial Number Proof")
	}

	offset := 0
	var err error

	proof.stmt.sn = new(crypto.Point)
	proof.stmt.sn, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += crypto.Ed25519KeySize

	proof.stmt.comSK = new(crypto.Point)
	proof.stmt.comSK, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}

	offset += crypto.Ed25519KeySize
	proof.stmt.comInput = new(crypto.Point)
	proof.stmt.comInput, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}

	offset += crypto.Ed25519KeySize
	proof.tSK = new(crypto.Point)
	proof.tSK, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}

	offset += crypto.Ed25519KeySize
	proof.tInput = new(crypto.Point)
	proof.tInput, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}

	offset += crypto.Ed25519KeySize
	proof.tSN = new(crypto.Point)
	proof.tSN, err = new(crypto.Point).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])
	if err != nil {
		return err
	}

	offset += crypto.Ed25519KeySize
	proof.zSK = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])

	offset += crypto.Ed25519KeySize
	proof.zRSK = new(crypto.Scalar).FromBytesS(bytes[offset : offset+crypto.Ed25519KeySize])

	offset += crypto.Ed25519KeySize
	proof.zInput = new(crypto.Scalar).FromBytesS(bytes[offset : offset+common.BigIntSize])

	offset += crypto.Ed25519KeySize
	proof.zRInput = new(crypto.Scalar).FromBytesS(bytes[offset : offset+common.BigIntSize])

	return nil
}

func (wit SNPrivacyWitness) Prove(mess []byte) (*SNPrivacyProof, error) {

	eSK := crypto.RandomScalar()
	eSND := crypto.RandomScalar()
	dSK := crypto.RandomScalar()
	dSND := crypto.RandomScalar()

	// calculate tSeed = g_SK^eSK * h^dSK
	tSeed := crypto.PedCom.CommitAtIndex(eSK, dSK, crypto.PedersenPrivateKeyIndex)

	// calculate tSND = g_SND^eSND * h^dSND
	tInput := crypto.PedCom.CommitAtIndex(eSND, dSND, crypto.PedersenSndIndex)

	// calculate tSND = g_SK^eSND * h^dSND2
	tOutput := new(crypto.Point).ScalarMult(wit.stmt.sn, new(crypto.Scalar).Add(eSK, eSND))

	// calculate x = hash(tSeed || tInput || tSND2 || tOutput)
	x := new(crypto.Scalar)
	if mess == nil {
		x = utils.GenerateChallenge([][]byte{
			wit.stmt.sn.ToBytesS(),
			wit.stmt.comSK.ToBytesS(),
			tSeed.ToBytesS(),
			tInput.ToBytesS(),
			tOutput.ToBytesS()})
	} else {
		x.FromBytesS(mess)
	}

	// Calculate zSeed = sk * x + eSK
	zSeed := new(crypto.Scalar).Mul(wit.sk, x)
	zSeed.Add(zSeed, eSK)
	//zSeed.Mod(zSeed, crypto.Curve.Params().N)

	// Calculate zRSeed = rSK * x + dSK
	zRSeed := new(crypto.Scalar).Mul(wit.rSK, x)
	zRSeed.Add(zRSeed, dSK)
	//zRSeed.Mod(zRSeed, crypto.Curve.Params().N)

	// Calculate zInput = input * x + eSND
	zInput := new(crypto.Scalar).Mul(wit.input, x)
	zInput.Add(zInput, eSND)
	//zInput.Mod(zInput, crypto.Curve.Params().N)

	// Calculate zRInput = rInput * x + dSND
	zRInput := new(crypto.Scalar).Mul(wit.rInput, x)
	zRInput.Add(zRInput, dSND)
	//zRInput.Mod(zRInput, crypto.Curve.Params().N)

	proof := new(SNPrivacyProof).Init()
	proof.Set(wit.stmt, tSeed, tInput, tOutput, zSeed, zRSeed, zInput, zRInput)
	return proof, nil
}

func (proof SNPrivacyProof) Verify(mess []byte) (bool, error) {
	// re-calculate x = hash(tSeed || tInput || tSND2 || tOutput)
	x := new(crypto.Scalar)
	if mess == nil {
		x = utils.GenerateChallenge([][]byte{
			proof.stmt.sn.ToBytesS(),
			proof.stmt.comSK.ToBytesS(),
			proof.tSK.ToBytesS(),
			proof.tInput.ToBytesS(),
			proof.tSN.ToBytesS()})
	} else {
		x.FromBytesS(mess)
	}

	// Check gSND^zInput * h^zRInput = input^x * tInput
	leftPoint1 := crypto.PedCom.CommitAtIndex(proof.zInput, proof.zRInput, crypto.PedersenSndIndex)

	rightPoint1 := new(crypto.Point).ScalarMult(proof.stmt.comInput, x)
	rightPoint1.Add(rightPoint1, proof.tInput)

	if !crypto.IsPointEqual(leftPoint1, rightPoint1) {
		//Logger.Log.Errorf("verify serial number privacy proof statement 1 failed")
		return false, errors.New("verify serial number privacy proof statement 1 failed")
	}

	// Check gSK^zSeed * h^zRSeed = vKey^x * tSeed
	leftPoint2 := crypto.PedCom.CommitAtIndex(proof.zSK, proof.zRSK, crypto.PedersenPrivateKeyIndex)

	rightPoint2 := new(crypto.Point).ScalarMult(proof.stmt.comSK, x)
	rightPoint2.Add(rightPoint2, proof.tSK)

	if !crypto.IsPointEqual(leftPoint2, rightPoint2) {
		return false, errors.New("verify serial number privacy proof statement 2 failed")
	}

	// Check sn^(zSeed + zInput) = gSK^x * tOutput
	leftPoint3 := new(crypto.Point).ScalarMult(proof.stmt.sn, new(crypto.Scalar).Add(proof.zSK, proof.zInput))

	rightPoint3 := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], x)
	rightPoint3.Add(rightPoint3, proof.tSN)

	if !crypto.IsPointEqual(leftPoint3, rightPoint3) {
		//privacy.Logger.Log.Errorf("verify serial number privacy proof statement 3 failed")
		return false, errors.New("verify serial number privacy proof statement 3 failed")
	}

	return true, nil
}

func (proof SNPrivacyProof) VerifyOld(mess []byte) (bool, error) {
	// re-calculate x = hash(tSeed || tInput || tSND2 || tOutput)
	x := new(crypto.Scalar)
	if mess == nil {
		x = utils.GenerateChallenge([][]byte{
			proof.tSK.ToBytesS(),
			proof.tInput.ToBytesS(),
			proof.tSN.ToBytesS()})
	} else {
		x.FromBytesS(mess)
	}

	// Check gSND^zInput * h^zRInput = input^x * tInput
	leftPoint1 := crypto.PedCom.CommitAtIndex(proof.zInput, proof.zRInput, crypto.PedersenSndIndex)

	rightPoint1 := new(crypto.Point).ScalarMult(proof.stmt.comInput, x)
	rightPoint1.Add(rightPoint1, proof.tInput)

	if !crypto.IsPointEqual(leftPoint1, rightPoint1) {
		//Logger.Log.Errorf("verify serial number privacy proof statement 1 failed")
		return false, errors.New("verifyOld serial number privacy proof statement 1 failed")
	}

	// Check gSK^zSeed * h^zRSeed = vKey^x * tSeed
	leftPoint2 := crypto.PedCom.CommitAtIndex(proof.zSK, proof.zRSK, crypto.PedersenPrivateKeyIndex)

	rightPoint2 := new(crypto.Point).ScalarMult(proof.stmt.comSK, x)
	rightPoint2.Add(rightPoint2, proof.tSK)

	if !crypto.IsPointEqual(leftPoint2, rightPoint2) {
		return false, errors.New("verifyOld serial number privacy proof statement 2 failed")
	}

	// Check sn^(zSeed + zInput) = gSK^x * tOutput
	leftPoint3 := new(crypto.Point).ScalarMult(proof.stmt.sn, new(crypto.Scalar).Add(proof.zSK, proof.zInput))

	rightPoint3 := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], x)
	rightPoint3.Add(rightPoint3, proof.tSN)

	if !crypto.IsPointEqual(leftPoint3, rightPoint3) {
		//privacy.Logger.Log.Errorf("verify serial number privacy proof statement 3 failed")
		return false, errors.New("verifyOld serial number privacy proof statement 3 failed")
	}

	return true, nil
}


func Copy(proof SNPrivacyProof) *SNPrivacyProof{
	tmpProof := new(SNPrivacyProof)
	tmpProof.tInput = new(crypto.Point).Set(proof.tInput)
	tmpProof.tSK = new(crypto.Point).Set(proof.tSK)
	tmpProof.tSN = new(crypto.Point).Set(proof.tSN)
	tmpProof.zInput = new(crypto.Scalar).Set(proof.zInput)
	tmpProof.zRInput = new(crypto.Scalar).Set(proof.zRInput)
	tmpProof.zSK = new(crypto.Scalar).Set(proof.zSK)
	tmpProof.zRSK = new(crypto.Scalar).Set(proof.zRSK)

	sn := new(crypto.Point).Set(proof.stmt.sn)
	comSK := new(crypto.Point).Set(proof.stmt.comSK)
	comInput := new(crypto.Point).Set(proof.stmt.comInput)
	tmpProof.stmt = new(SerialNumberPrivacyStatement)
	tmpProof.stmt.Set(sn, comSK, comInput)

	return tmpProof
}
