package zkp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/proof/range_proof"
	privacyUtils "github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/bulletproofs"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/oneoutofmany"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/serialnumbernoprivacy"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/serialnumberprivacy"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/utils"
	"math/big"
)

const FixedRandomnessString = "fixedrandomness"

// FixedRandomnessShardID is fixed randomness for shardID commitment from param.BCHeightBreakPointFixRandShardCM
// is result from HashToScalar([]byte(privacy.FixedRandomnessString))
var FixedRandomnessShardID = new(crypto.Scalar).FromBytesS([]byte{0x60, 0xa2, 0xab, 0x35, 0x26, 0x9, 0x97, 0x7c, 0x6b, 0xe1, 0xba, 0xec, 0xbf, 0x64, 0x27, 0x2, 0x6a, 0x9c, 0xe8, 0x10, 0x9e, 0x93, 0x4a, 0x0, 0x47, 0x83, 0x15, 0x48, 0x63, 0xeb, 0xda, 0x6})

// ProofV1 represents a payment proof for a transaction of version 1.
// A ProofV1 consists of the following
//	- oneOfManyProof: used to prove the existence of real input coins within a set of input coins
//	(used in private transactions only).
// 	- serialNumberProof: a sigma protocol for proving that the serial numbers are derived from the real input coins.
//	It is used to avoid double-spending (used in private transactions only).
//	- serialNumberNoPrivacyProof: same as serialNumberProof but used in non-private transaction.
//	- rangeProofWitness: a proof proving each output coin's value lies in a specific range (i.e, [0, 2^64-1]) without
//	revealing the output coin's value.
type ProofV1 struct {
	// for input coins
	oneOfManyProof    []*oneoutofmany.OneOutOfManyProof
	serialNumberProof []*serialnumberprivacy.SNPrivacyProof
	// it is exits when tx has no privacy
	serialNumberNoPrivacyProof []*serialnumbernoprivacy.SNNoPrivacyProof

	// for output coins
	// for proving each value and sum of them are less than a threshold value
	rangeProof *bulletproofs.RangeProof

	inputCoins  []coin.PlainCoin
	outputCoins []*coin.CoinV1

	commitmentOutputValue   []*crypto.Point
	commitmentOutputSND     []*crypto.Point
	commitmentOutputShardID []*crypto.Point

	commitmentInputSecretKey *crypto.Point
	commitmentInputValue     []*crypto.Point
	commitmentInputSND       []*crypto.Point
	commitmentInputShardID   *crypto.Point

	commitmentIndices []uint64
}

func (proof *ProofV1) GetVersion() uint8 { return 1 }

// GET/SET function
func (proof ProofV1) GetOneOfManyProof() []*oneoutofmany.OneOutOfManyProof {
	return proof.oneOfManyProof
}
func (proof ProofV1) GetSerialNumberProof() []*serialnumberprivacy.SNPrivacyProof {
	return proof.serialNumberProof
}
func (proof ProofV1) GetSerialNumberNoPrivacyProof() []*serialnumbernoprivacy.SNNoPrivacyProof {
	return proof.serialNumberNoPrivacyProof
}
func (proof ProofV1) GetRangeProof() range_proof.RangeProof {
	return proof.rangeProof
}
func (proof ProofV1) GetCommitmentOutputValue() []*crypto.Point {
	return proof.commitmentOutputValue
}
func (proof ProofV1) GetCommitmentOutputSND() []*crypto.Point {
	return proof.commitmentOutputSND
}
func (proof ProofV1) GetCommitmentOutputShardID() []*crypto.Point {
	return proof.commitmentOutputShardID
}
func (proof ProofV1) GetCommitmentInputSecretKey() *crypto.Point {
	return proof.commitmentInputSecretKey
}
func (proof ProofV1) GetCommitmentInputValue() []*crypto.Point {
	return proof.commitmentInputValue
}
func (proof ProofV1) GetCommitmentInputSND() []*crypto.Point { return proof.commitmentInputSND }
func (proof ProofV1) GetCommitmentInputShardID() *crypto.Point {
	return proof.commitmentInputShardID
}
func (proof ProofV1) GetCommitmentIndices() []uint64  { return proof.commitmentIndices }
func (proof ProofV1) GetInputCoins() []coin.PlainCoin { return proof.inputCoins }
func (proof ProofV1) GetOutputCoins() []coin.Coin {
	res := make([]coin.Coin, len(proof.outputCoins))
	for i := 0; i < len(proof.outputCoins); i += 1 {
		res[i] = proof.outputCoins[i]
	}
	return res
}

func (proof *ProofV1) SetCommitmentShardID(commitmentShardID *crypto.Point){proof.commitmentInputShardID = commitmentShardID}
func (proof *ProofV1) SetCommitmentInputSND(commitmentInputSND []*crypto.Point){proof.commitmentInputSND = commitmentInputSND}
func (proof *ProofV1) SetAggregatedRangeProof(aggregatedRangeProof *aggregatedrange.AggregatedRangeProof) {proof.rangeProof = aggregatedRangeProof}
func (proof *ProofV1) SetSerialNumberProof(serialNumberProof []*serialnumberprivacy.SNPrivacyProof) {proof.serialNumberProof = serialNumberProof}
func (proof *ProofV1) SetOneOfManyProof(oneOfManyProof []*oneoutofmany.OneOutOfManyProof) {proof.oneOfManyProof = oneOfManyProof}
func (proof *ProofV1) SetSerialNumberNoPrivacyProof(serialNumberNoPrivacyProof []*serialnumbernoprivacy.SNNoPrivacyProof) {proof.serialNumberNoPrivacyProof = serialNumberNoPrivacyProof}
func (proof *ProofV1) SetCommitmentInputValue(commitmentInputValue []*crypto.Point) {proof.commitmentInputValue = commitmentInputValue}

// SetCommitmentInputSND sets v as the inputSND commitments of a ProofV1.
func (proof *ProofV1) SetCommitmentInputSND(v []*crypto.Point) {
	proof.commitmentInputSND = v
}

// SetAggregatedRangeProof sets v as the range proof of a ProofV1.
func (proof *ProofV1) SetAggregatedRangeProof(v *bulletproofs.RangeProof) {
	proof.rangeProof = v
}

// SetSerialNumberProof sets v as the serial number proofs of a ProofV1.
func (proof *ProofV1) SetSerialNumberProof(v []*serialnumberprivacy.SNPrivacyProof) {
	proof.serialNumberProof = v
}

// SetOneOfManyProof sets v as the one-of-many proofs of a ProofV1.
func (proof *ProofV1) SetOneOfManyProof(v []*oneoutofmany.OneOutOfManyProof) {
	proof.oneOfManyProof = v
}

// SetSerialNumberNoPrivacyProof sets v as the serial number no privacy proofs of a ProofV1.
func (proof *ProofV1) SetSerialNumberNoPrivacyProof(v []*serialnumbernoprivacy.SNNoPrivacyProof) {
	proof.serialNumberNoPrivacyProof = v
}

// SetCommitmentInputValue sets v as the commitments to input values of a ProofV1.
func (proof *ProofV1) SetCommitmentInputValue(v []*crypto.Point) {
	proof.commitmentInputValue = v
}

// SetInputCoins sets v as the input coins of a ProofV1.
func (proof *ProofV1) SetInputCoins(v []coin.PlainCoin) error {
	var err error
	proof.inputCoins = make([]coin.PlainCoin, len(v))
	for i := 0; i < len(v); i += 1 {
		b := v[i].Bytes()
		if proof.inputCoins[i], err = coin.NewPlainCoinFromByte(b); err != nil {
			return err
		}
	}
	return nil
}

// SetOutputCoins's input should be all coinv1
func (proof *ProofV1) SetOutputCoins(v []coin.Coin) error {
	var err error
	proof.outputCoins = make([]*coin.CoinV1, len(v))
	for i := 0; i < len(v); i += 1 {
		b := v[i].Bytes()
		proof.outputCoins[i] = new(coin.CoinV1)
		if err = proof.outputCoins[i].SetBytes(b); err != nil {
			return err
		}
	}
	return nil
}

// End GET/SET function

// Init
func (proof *ProofV1) Init() {
	rangeProof := &bulletproofs.RangeProof{}
	rangeProof.Init()
	proof.oneOfManyProof = []*oneoutofmany.OneOutOfManyProof{}
	proof.serialNumberProof = []*serialnumberprivacy.SNPrivacyProof{}
	proof.rangeProof = rangeProof
	proof.inputCoins = []coin.PlainCoin{}
	proof.outputCoins = []*coin.CoinV1{}

	proof.commitmentOutputValue = []*crypto.Point{}
	proof.commitmentOutputSND = []*crypto.Point{}
	proof.commitmentOutputShardID = []*crypto.Point{}

	proof.commitmentInputSecretKey = new(crypto.Point)
	proof.commitmentInputValue = []*crypto.Point{}
	proof.commitmentInputSND = []*crypto.Point{}
	proof.commitmentInputShardID = new(crypto.Point)
}

func (proof ProofV1) MarshalJSON() ([]byte, error) {
	data := proof.Bytes()
	//temp := base58.Base58Check{}.Encode(data, common.ZeroByte)
	temp := base64.StdEncoding.EncodeToString(data)
	return json.Marshal(temp)
}

func (proof *ProofV1) UnmarshalJSON(data []byte) error {
	dataStr := common.EmptyString
	errJson := json.Unmarshal(data, &dataStr)
	if errJson != nil {
		return errJson
	}

	//temp, _, err := base58.Base58Check{}.Decode(dataStr)
	temp, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return err
	}

	errSetBytes := proof.SetBytes(temp)
	if errSetBytes != nil {
		return errSetBytes
	}

	return nil
}

func (proof ProofV1) Bytes() []byte {
	var bytes []byte
	hasPrivacy := len(proof.oneOfManyProof) > 0

	// OneOfManyProofSize
	bytes = append(bytes, byte(len(proof.oneOfManyProof)))
	for i := 0; i < len(proof.oneOfManyProof); i++ {
		oneOfManyProof := proof.oneOfManyProof[i].Bytes()
		bytes = append(bytes, common.IntToBytes(utils.OneOfManyProofSize)...)
		bytes = append(bytes, oneOfManyProof...)
	}

	// SerialNumberProofSize
	bytes = append(bytes, byte(len(proof.serialNumberProof)))
	for i := 0; i < len(proof.serialNumberProof); i++ {
		serialNumberProof := proof.serialNumberProof[i].Bytes()
		bytes = append(bytes, common.IntToBytes(utils.SnPrivacyProofSize)...)
		bytes = append(bytes, serialNumberProof...)
	}

	// SNNoPrivacyProofSize
	bytes = append(bytes, byte(len(proof.serialNumberNoPrivacyProof)))
	for i := 0; i < len(proof.serialNumberNoPrivacyProof); i++ {
		snNoPrivacyProof := proof.serialNumberNoPrivacyProof[i].Bytes()
		bytes = append(bytes, byte(utils.SnNoPrivacyProofSize))
		bytes = append(bytes, snNoPrivacyProof...)
	}

	//ComOutputMultiRangeProofSize
	if hasPrivacy {
		comOutputMultiRangeProof := proof.rangeProof.Bytes()
		bytes = append(bytes, common.IntToBytes(len(comOutputMultiRangeProof))...)
		bytes = append(bytes, comOutputMultiRangeProof...)
	} else {
		bytes = append(bytes, []byte{0, 0}...)
	}

	// InputCoins
	bytes = append(bytes, byte(len(proof.inputCoins)))
	for i := 0; i < len(proof.inputCoins); i++ {
		inputCoins := proof.inputCoins[i].Bytes()
		bytes = append(bytes, byte(len(inputCoins)))
		bytes = append(bytes, inputCoins...)
	}

	// OutputCoins
	bytes = append(bytes, byte(len(proof.outputCoins)))
	for i := 0; i < len(proof.outputCoins); i++ {
		outputCoins := proof.outputCoins[i].Bytes()
		lenOutputCoins := len(outputCoins)
		lenOutputCoinsBytes := []byte{}
		if lenOutputCoins < 256 {
			lenOutputCoinsBytes = []byte{byte(lenOutputCoins)}
		} else {
			lenOutputCoinsBytes = common.IntToBytes(lenOutputCoins)
		}

		bytes = append(bytes, lenOutputCoinsBytes...)
		bytes = append(bytes, outputCoins...)
	}

	// ComOutputValue
	bytes = append(bytes, byte(len(proof.commitmentOutputValue)))
	for i := 0; i < len(proof.commitmentOutputValue); i++ {
		comOutputValue := proof.commitmentOutputValue[i].ToBytesS()
		bytes = append(bytes, byte(crypto.Ed25519KeySize))
		bytes = append(bytes, comOutputValue...)
	}

	// ComOutputSND
	bytes = append(bytes, byte(len(proof.commitmentOutputSND)))
	for i := 0; i < len(proof.commitmentOutputSND); i++ {
		comOutputSND := proof.commitmentOutputSND[i].ToBytesS()
		bytes = append(bytes, byte(crypto.Ed25519KeySize))
		bytes = append(bytes, comOutputSND...)
	}

	// ComOutputShardID
	bytes = append(bytes, byte(len(proof.commitmentOutputShardID)))
	for i := 0; i < len(proof.commitmentOutputShardID); i++ {
		comOutputShardID := proof.commitmentOutputShardID[i].ToBytesS()
		bytes = append(bytes, byte(crypto.Ed25519KeySize))
		bytes = append(bytes, comOutputShardID...)
	}

	//ComInputSK 				*crypto.Point
	if proof.commitmentInputSecretKey != nil {
		comInputSK := proof.commitmentInputSecretKey.ToBytesS()
		bytes = append(bytes, byte(crypto.Ed25519KeySize))
		bytes = append(bytes, comInputSK...)
	} else {
		bytes = append(bytes, byte(0))
	}

	//ComInputValue 		[]*crypto.Point
	bytes = append(bytes, byte(len(proof.commitmentInputValue)))
	for i := 0; i < len(proof.commitmentInputValue); i++ {
		comInputValue := proof.commitmentInputValue[i].ToBytesS()
		bytes = append(bytes, byte(crypto.Ed25519KeySize))
		bytes = append(bytes, comInputValue...)
	}

	//ComInputSND 			[]*privacy.Point
	bytes = append(bytes, byte(len(proof.commitmentInputSND)))
	for i := 0; i < len(proof.commitmentInputSND); i++ {
		comInputSND := proof.commitmentInputSND[i].ToBytesS()
		bytes = append(bytes, byte(crypto.Ed25519KeySize))
		bytes = append(bytes, comInputSND...)
	}

	//ComInputShardID 	*privacy.Point
	if proof.commitmentInputShardID != nil {
		comInputShardID := proof.commitmentInputShardID.ToBytesS()
		bytes = append(bytes, byte(crypto.Ed25519KeySize))
		bytes = append(bytes, comInputShardID...)
	} else {
		bytes = append(bytes, byte(0))
	}

	// convert commitment index to bytes array
	for i := 0; i < len(proof.commitmentIndices); i++ {
		bytes = append(bytes, common.AddPaddingBigInt(big.NewInt(int64(proof.commitmentIndices[i])), common.Uint64Size)...)
	}
	//fmt.Printf("BYTES ------------------ %v\n", bytes)
	//fmt.Printf("LEN BYTES ------------------ %v\n", len(bytes))

	return bytes
}

func (proof *ProofV1) SetBytes(proofbytes []byte) error {
	if len(proofbytes) == 0 {
		return fmt.Errorf("length of proof is zero")
	}
	var err error
	offset := 0

	// Set OneOfManyProofSize
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range one out of many proof")
	}
	lenOneOfManyProofArray := int(proofbytes[offset])
	offset += 1
	proof.oneOfManyProof = make([]*oneoutofmany.OneOutOfManyProof, lenOneOfManyProofArray)
	for i := 0; i < lenOneOfManyProofArray; i++ {
		if offset+2 > len(proofbytes) {
			return fmt.Errorf("out of range one out of many proof")
		}
		lenOneOfManyProof := common.BytesToInt(proofbytes[offset : offset+2])
		offset += 2
		proof.oneOfManyProof[i] = new(oneoutofmany.OneOutOfManyProof).Init()

		if offset+lenOneOfManyProof > len(proofbytes) {
			return fmt.Errorf("out of range one out of many proof")
		}
		err := proof.oneOfManyProof[i].SetBytes(proofbytes[offset : offset+lenOneOfManyProof])
		if err != nil {
			return err
		}
		offset += lenOneOfManyProof
	}

	// Set serialNumberProofSize
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range serial number proof")
	}
	lenSerialNumberProofArray := int(proofbytes[offset])
	offset += 1
	proof.serialNumberProof = make([]*serialnumberprivacy.SNPrivacyProof, lenSerialNumberProofArray)
	for i := 0; i < lenSerialNumberProofArray; i++ {
		if offset+2 > len(proofbytes) {
			return fmt.Errorf("out of range serial number proof")
		}
		lenSerialNumberProof := common.BytesToInt(proofbytes[offset : offset+2])
		offset += 2
		proof.serialNumberProof[i] = new(serialnumberprivacy.SNPrivacyProof).Init()

		if offset+lenSerialNumberProof > len(proofbytes) {
			return fmt.Errorf("out of range serial number proof")
		}
		err := proof.serialNumberProof[i].SetBytes(proofbytes[offset : offset+lenSerialNumberProof])
		if err != nil {
			return err
		}
		offset += lenSerialNumberProof
	}

	// Set SNNoPrivacyProofSize
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range serial number no privacy proof")
	}
	lenSNNoPrivacyProofArray := int(proofbytes[offset])
	offset += 1
	proof.serialNumberNoPrivacyProof = make([]*serialnumbernoprivacy.SNNoPrivacyProof, lenSNNoPrivacyProofArray)
	for i := 0; i < lenSNNoPrivacyProofArray; i++ {
		if offset >= len(proofbytes) {
			return fmt.Errorf("out of range serial number no privacy proof")
		}
		lenSNNoPrivacyProof := int(proofbytes[offset])
		offset += 1

		proof.serialNumberNoPrivacyProof[i] = new(serialnumbernoprivacy.SNNoPrivacyProof).Init()
		if offset+lenSNNoPrivacyProof > len(proofbytes) {
			return fmt.Errorf("out of range serial number no privacy proof")
		}
		err := proof.serialNumberNoPrivacyProof[i].SetBytes(proofbytes[offset : offset+lenSNNoPrivacyProof])
		if err != nil {
			return err
		}
		offset += lenSNNoPrivacyProof
	}

	//ComOutputMultiRangeProofSize *rangeProofWitness
	if offset+2 > len(proofBytes) {
		return fmt.Errorf("out of range aggregated range proof")
	}
	lenComOutputMultiRangeProof := common.BytesToInt(proofbytes[offset : offset+2])
	offset += 2
	if lenComOutputMultiRangeProof > 0 {
		rangeProof := &bulletproofs.RangeProof{}
		rangeProof.Init()
		proof.rangeProof = rangeProof
		if offset+lenComOutputMultiRangeProof > len(proofBytes) {
			return fmt.Errorf("out of range aggregated range proof")
		}
		err := proof.rangeProof.SetBytes(proofbytes[offset : offset+lenComOutputMultiRangeProof])
		if err != nil {
			return err
		}
		offset += lenComOutputMultiRangeProof
	}

	//InputCoins  []*coin.PlainCoinV1
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range input coins")
	}
	lenInputCoinsArray := int(proofbytes[offset])
	offset += 1
	proof.inputCoins = make([]coin.PlainCoin, lenInputCoinsArray)
	for i := 0; i < lenInputCoinsArray; i++ {
		if offset >= len(proofbytes) {
			return fmt.Errorf("out of range input coins")
		}
		lenInputCoin := int(proofbytes[offset])
		offset += 1

		if offset+lenInputCoin > len(proofbytes) {
			return fmt.Errorf("out of range input coins")
		}
		coinBytes := proofbytes[offset : offset+lenInputCoin]
		proof.inputCoins[i], err = coin.NewPlainCoinFromByte(coinBytes)
		if err != nil {
			return fmt.Errorf("set byte to inputCoin got error")
		}
		offset += lenInputCoin
	}

	//OutputCoins []*privacy.OutputCoin
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range output coins")
	}
	lenOutputCoinsArray := int(proofbytes[offset])
	offset += 1
	proof.outputCoins = make([]*coin.CoinV1, lenOutputCoinsArray)
	for i := 0; i < lenOutputCoinsArray; i++ {
		proof.outputCoins[i] = new(coin.CoinV1)
		// try get 1-byte for len
		if offset >= len(proofbytes) {
			return fmt.Errorf("out of range output coins")
		}
		lenOutputCoin := int(proofbytes[offset])
		offset += 1

		if offset+lenOutputCoin > len(proofbytes) {
			return fmt.Errorf("out of range output coins")
		}
		err := proof.outputCoins[i].SetBytes(proofbytes[offset : offset+lenOutputCoin])
		if err != nil {
			// 1-byte is wrong
			// try get 2-byte for len
			if offset+1 > len(proofbytes) {
				return fmt.Errorf("out of range output coins")
			}
			lenOutputCoin = common.BytesToInt(proofbytes[offset-1 : offset+1])
			offset += 1

			if offset+lenOutputCoin > len(proofbytes) {
				return fmt.Errorf("out of range output coins")
			}
			err1 := proof.outputCoins[i].SetBytes(proofbytes[offset : offset+lenOutputCoin])
			if err1 != nil {
				return err
			}
		}
		offset += lenOutputCoin
	}
	//ComOutputValue   []*privacy.Point
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range commitment output coins value")
	}
	lenComOutputValueArray := int(proofbytes[offset])
	offset += 1
	proof.commitmentOutputValue = make([]*crypto.Point, lenComOutputValueArray)
	for i := 0; i < lenComOutputValueArray; i++ {
		if offset >= len(proofbytes) {
			return fmt.Errorf("out of range commitment output coins value")
		}
		lenComOutputValue := int(proofbytes[offset])
		offset += 1

		if offset+lenComOutputValue > len(proofbytes) {
			return fmt.Errorf("out of range commitment output coins value")
		}
		proof.commitmentOutputValue[i], err = new(crypto.Point).FromBytesS(proofbytes[offset : offset+lenComOutputValue])
		if err != nil {
			return err
		}
		offset += lenComOutputValue
	}
	//ComOutputSND     []*crypto.Point
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range commitment output coins snd")
	}
	lenComOutputSNDArray := int(proofbytes[offset])
	offset += 1
	proof.commitmentOutputSND = make([]*crypto.Point, lenComOutputSNDArray)
	for i := 0; i < lenComOutputSNDArray; i++ {
		if offset >= len(proofbytes) {
			return fmt.Errorf("out of range commitment output coins snd")
		}
		lenComOutputSND := int(proofbytes[offset])
		offset += 1

		if offset+lenComOutputSND > len(proofbytes) {
			return fmt.Errorf("out of range commitment output coins snd")
		}
		proof.commitmentOutputSND[i], err = new(crypto.Point).FromBytesS(proofbytes[offset : offset+lenComOutputSND])

		if err != nil {
			return err
		}
		offset += lenComOutputSND
	}

	// commitmentOutputShardID
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range commitment output coins shardid")
	}
	lenComOutputShardIdArray := int(proofbytes[offset])
	offset += 1
	proof.commitmentOutputShardID = make([]*crypto.Point, lenComOutputShardIdArray)
	for i := 0; i < lenComOutputShardIdArray; i++ {
		if offset >= len(proofbytes) {
			return fmt.Errorf("out of range commitment output coins shardid")
		}
		lenComOutputShardId := int(proofbytes[offset])
		offset += 1

		if offset+lenComOutputShardId > len(proofbytes) {
			return fmt.Errorf("out of range commitment output coins shardid")
		}
		proof.commitmentOutputShardID[i], err = new(crypto.Point).FromBytesS(proofbytes[offset : offset+lenComOutputShardId])

		if err != nil {
			return err
		}
		offset += lenComOutputShardId
	}

	//ComInputSK 				*crypto.Point
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range commitment input coins private key")
	}
	lenComInputSK := int(proofbytes[offset])
	offset += 1
	if lenComInputSK > 0 {
		if offset+lenComInputSK > len(proofbytes) {
			return fmt.Errorf("out of range commitment input coins private key")
		}
		proof.commitmentInputSecretKey, err = new(crypto.Point).FromBytesS(proofbytes[offset : offset+lenComInputSK])

		if err != nil {
			return err
		}
		offset += lenComInputSK
	}
	//ComInputValue 		[]*crypto.Point
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range commitment input coins value")
	}
	lenComInputValueArr := int(proofbytes[offset])
	offset += 1
	proof.commitmentInputValue = make([]*crypto.Point, lenComInputValueArr)
	for i := 0; i < lenComInputValueArr; i++ {
		if offset >= len(proofbytes) {
			return fmt.Errorf("out of range commitment input coins value")
		}
		lenComInputValue := int(proofbytes[offset])
		offset += 1

		if offset+lenComInputValue > len(proofbytes) {
			return fmt.Errorf("out of range commitment input coins value")
		}
		proof.commitmentInputValue[i], err = new(crypto.Point).FromBytesS(proofbytes[offset : offset+lenComInputValue])

		if err != nil {
			return err
		}
		offset += lenComInputValue
	}
	//ComInputSND 			[]*crypto.Point
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range commitment input coins snd")
	}
	lenComInputSNDArr := int(proofbytes[offset])
	offset += 1
	proof.commitmentInputSND = make([]*crypto.Point, lenComInputSNDArr)
	for i := 0; i < lenComInputSNDArr; i++ {
		if offset >= len(proofbytes) {
			return fmt.Errorf("out of range commitment input coins snd")
		}
		lenComInputSND := int(proofbytes[offset])
		offset += 1

		if offset+lenComInputSND > len(proofbytes) {
			return fmt.Errorf("out of range commitment input coins snd")
		}
		proof.commitmentInputSND[i], err = new(crypto.Point).FromBytesS(proofbytes[offset : offset+lenComInputSND])

		if err != nil {
			return err
		}
		offset += lenComInputSND
	}
	//ComInputShardID 	*crypto.Point
	if offset >= len(proofbytes) {
		return fmt.Errorf("out of range commitment input coins shardid")
	}
	lenComInputShardID := int(proofbytes[offset])
	offset += 1
	if lenComInputShardID > 0 {
		if offset+lenComInputShardID > len(proofbytes) {
			return fmt.Errorf("out of range commitment input coins shardid")
		}
		proof.commitmentInputShardID, err = new(crypto.Point).FromBytesS(proofbytes[offset : offset+lenComInputShardID])

		if err != nil {
			return err
		}
		offset += lenComInputShardID
	}

	// get commitments list
	proof.commitmentIndices = make([]uint64, len(proof.oneOfManyProof)*privacy_util.CommitmentRingSize)
	for i := 0; i < len(proof.oneOfManyProof)*privacy_util.CommitmentRingSize; i++ {
		if offset+common.Uint64Size > len(proofbytes) {
			return fmt.Errorf("out of range commitment indices")
		}
		proof.commitmentIndices[i] = new(big.Int).SetBytes(proofbytes[offset : offset+common.Uint64Size]).Uint64()
		offset = offset + common.Uint64Size
	}

	//fmt.Printf("SETBYTES ------------------ %v\n", proof.Bytes())

	return nil
}

func (proof *ProofV1) IsPrivacy() bool {
	if proof == nil || len(proof.GetOneOfManyProof()) == 0 {
		return false
	}
	return true
}

func isBadScalar(sc *crypto.Scalar) bool {
	if sc == nil || !sc.ScalarValid() {
		return true
	}
	return false
}

func isBadPoint(point *crypto.Point) bool {
	if point == nil || !point.PointValid() {
		return true
	}
	return false
}
