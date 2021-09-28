package conversion

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/proof/range_proof"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/serialnumbernoprivacy"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/utils"
)

const (
	ProofVersion = 255
)

// ConversionProof represents a payment proof used in conversion transactions.
// For a conversion proof, its version will be counted down from 255 -> 0
// It should contain inputCoins of v1 and outputCoins of v2 because it convert v1 to v2.
type ConversionProof struct {
	Version                    uint8
	inputCoins                 []*coin.PlainCoinV1
	outputCoins                []*coin.CoinV2
	serialNumberNoPrivacyProof []*serialnumbernoprivacy.SNNoPrivacyProof
}

// MarshalJSON returns the JSON-marshalled data of a ConversionProof.
func (proof ConversionProof) MarshalJSON() ([]byte, error) {
	data := proof.Bytes()
	//temp := base58.Base58Check{}.Encode(data, common.ZeroByte)
	temp := base64.StdEncoding.EncodeToString(data)
	return json.Marshal(temp)
}

// UnmarshalJSON un-marshals raw-byte data into a ConversionProof.
func (proof *ConversionProof) UnmarshalJSON(data []byte) error {
	dataStr := common.EmptyString
	errJson := json.Unmarshal(data, &dataStr)
	if errJson != nil {
		return errJson
	}
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

// Init creates an empty ConversionProof.
func (proof ConversionProof) Init() {
	proof.Version = ProofVersion
	proof.inputCoins = []*coin.PlainCoinV1{}
	proof.outputCoins = []*coin.CoinV2{}
	proof.serialNumberNoPrivacyProof = []*serialnumbernoprivacy.SNNoPrivacyProof{}
}

// GetVersion returns the version of a ConversionProof, which is 255.
func (proof ConversionProof) GetVersion() uint8 { return ProofVersion }

// SetVersion sets the version of a ConversionProof to 255.
func (proof *ConversionProof) SetVersion(uint8) { proof.Version = ProofVersion }

// GetInputCoins returns the input coins of a ConversionProof.
func (proof ConversionProof) GetInputCoins() []coin.PlainCoin {
	res := make([]coin.PlainCoin, len(proof.inputCoins))
	for i := 0; i < len(proof.inputCoins); i += 1 {
		res[i] = proof.inputCoins[i]
	}
	return res
}

// GetOutputCoins returns the output coins of a ConversionProof.
func (proof ConversionProof) GetOutputCoins() []coin.Coin {
	res := make([]coin.Coin, len(proof.outputCoins))
	for i := 0; i < len(proof.outputCoins); i += 1 {
		res[i] = proof.outputCoins[i]
	}
	return res
}

// SetInputCoins sets v as the input coins of a ConversionProof.
// All input coins must must be of version 1, otherwise, it would crash.
func (proof *ConversionProof) SetInputCoins(v []coin.PlainCoin) error {
	proof.inputCoins = make([]*coin.PlainCoinV1, len(v))
	for i := 0; i < len(v); i += 1 {
		c, ok := v[i].(*coin.PlainCoinV1)
		if !ok {
			return fmt.Errorf("input coins should all be PlainCoinV1")
		}
		proof.inputCoins[i] = c
	}
	return nil
}

// SetOutputCoins sets v as the output coins of a ConversionProof.
// All output coins must be of version 2, otherwise, it would crash.
func (proof *ConversionProof) SetOutputCoins(v []coin.Coin) error {
	proof.outputCoins = make([]*coin.CoinV2, len(v))
	for i := 0; i < len(v); i += 1 {
		c, ok := v[i].(*coin.CoinV2)
		if !ok {
			return fmt.Errorf("output coins should all be CoinV2")
		}
		proof.outputCoins[i] = c
	}
	return nil
}

// GetRangeProof returns the range proof of a ConversionProof.
// A ConversionProof does not have range proof, everything is non-private.
func (proof ConversionProof) GetRangeProof() range_proof.RangeProof {
	return nil
}

// Bytes returns the byte-representation of a ConversionProof.
func (proof ConversionProof) Bytes() []byte {
	proofBytes := []byte{ProofVersion}

	// InputCoins
	proofBytes = append(proofBytes, byte(len(proof.inputCoins)))
	for i := 0; i < len(proof.inputCoins); i++ {
		inputCoins := proof.inputCoins[i].Bytes()
		proofBytes = append(proofBytes, byte(len(inputCoins)))
		proofBytes = append(proofBytes, inputCoins...)
	}

	// OutputCoins
	proofBytes = append(proofBytes, byte(len(proof.outputCoins)))
	for i := 0; i < len(proof.outputCoins); i++ {
		outputCoins := proof.outputCoins[i].Bytes()
		lenOutputCoins := len(outputCoins)
		lenOutputCoinsBytes := []byte{}
		if lenOutputCoins < 256 {
			lenOutputCoinsBytes = []byte{byte(lenOutputCoins)}
		} else {
			lenOutputCoinsBytes = common.IntToBytes(lenOutputCoins)
		}

		proofBytes = append(proofBytes, lenOutputCoinsBytes...)
		proofBytes = append(proofBytes, outputCoins...)
	}

	// SNNoPrivacyProofSize
	proofBytes = append(proofBytes, byte(len(proof.serialNumberNoPrivacyProof)))
	for i := 0; i < len(proof.serialNumberNoPrivacyProof); i++ {
		snNoPrivacyProof := proof.serialNumberNoPrivacyProof[i].Bytes()
		proofBytes = append(proofBytes, byte(utils.SnNoPrivacyProofSize))
		proofBytes = append(proofBytes, snNoPrivacyProof...)
	}

	return proofBytes
}

// SetBytes sets byte-representation data into a ConversionProof.
func (proof *ConversionProof) SetBytes(proofBytes []byte) error {
	if len(proofBytes) == 0 {
		return fmt.Errorf("proof bytes = 0")
	}
	if proofBytes[0] != proof.GetVersion() {
		return fmt.Errorf("proof bytes version is not correct")
	}
	if proof == nil {
		proof = new(ConversionProof)
	}
	proof.SetVersion(ProofVersion)
	offset := 1

	//InputCoins  []*coin.PlainCoinV1
	if offset >= len(proofBytes) {
		return fmt.Errorf("out of range input coins")
	}
	lenInputCoinsArray := int(proofBytes[offset])
	offset += 1
	proof.inputCoins = make([]*coin.PlainCoinV1, lenInputCoinsArray)
	for i := 0; i < lenInputCoinsArray; i++ {
		if offset >= len(proofBytes) {
			return fmt.Errorf("out of range input coins")
		}
		lenInputCoin := int(proofBytes[offset])
		offset += 1

		if offset+lenInputCoin > len(proofBytes) {
			return fmt.Errorf("out of range input coins")
		}
		coinBytes := proofBytes[offset : offset+lenInputCoin]
		if pc, err := coin.NewPlainCoinFromByte(coinBytes); err != nil {
			return err
		} else {
			var ok bool
			if proof.inputCoins[i], ok = pc.(*coin.PlainCoinV1); !ok {
				err := fmt.Errorf("cannot assert type of PlainCoin to PlainCoinV1")
				return err
			}
		}
		offset += lenInputCoin
	}

	//OutputCoins  []*coin.CoinV2
	if offset >= len(proofBytes) {
		return fmt.Errorf("out of range output coins")
	}
	lenOutputCoinsArray := int(proofBytes[offset])
	offset += 1
	proof.outputCoins = make([]*coin.CoinV2, lenOutputCoinsArray)
	for i := 0; i < lenOutputCoinsArray; i++ {
		proof.outputCoins[i] = new(coin.CoinV2)
		// try get 1-byte for len
		if offset >= len(proofBytes) {
			return fmt.Errorf("out of range output coins")
		}
		lenOutputCoin := int(proofBytes[offset])
		offset += 1

		if offset+lenOutputCoin > len(proofBytes) {
			return fmt.Errorf("out of range output coins")
		}
		err := proof.outputCoins[i].SetBytes(proofBytes[offset : offset+lenOutputCoin])
		if err != nil {
			// 1-byte is wrong
			// try get 2-byte for len
			if offset+1 > len(proofBytes) {
				return fmt.Errorf("out of range output coins")
			}
			lenOutputCoin = common.BytesToInt(proofBytes[offset-1 : offset+1])
			offset += 1

			if offset+lenOutputCoin > len(proofBytes) {
				return fmt.Errorf("out of range output coins")
			}
			err1 := proof.outputCoins[i].SetBytes(proofBytes[offset : offset+lenOutputCoin])
			if err1 != nil {
				return err1
			}
		}
		offset += lenOutputCoin

	}

	// SNNoPrivacyProof
	// Set SNNoPrivacyProofSize
	if offset >= len(proofBytes) {
		return fmt.Errorf("out of range serial number no privacy proof")
	}
	lenSNNoPrivacyProofArray := int(proofBytes[offset])
	offset += 1
	proof.serialNumberNoPrivacyProof = make([]*serialnumbernoprivacy.SNNoPrivacyProof, lenSNNoPrivacyProofArray)
	for i := 0; i < lenSNNoPrivacyProofArray; i++ {
		if offset >= len(proofBytes) {
			return fmt.Errorf("out of range serial number no privacy proof")
		}
		lenSNNoPrivacyProof := int(proofBytes[offset])
		offset += 1

		proof.serialNumberNoPrivacyProof[i] = new(serialnumbernoprivacy.SNNoPrivacyProof).Init()
		if offset+lenSNNoPrivacyProof > len(proofBytes) {
			return fmt.Errorf("out of range serial number no privacy proof")
		}
		err := proof.serialNumberNoPrivacyProof[i].SetBytes(proofBytes[offset : offset+lenSNNoPrivacyProof])
		if err != nil {
			return err
		}
		offset += lenSNNoPrivacyProof
	}
	return nil
}

// IsPrivacy returns false.
func (proof *ConversionProof) IsPrivacy() bool {
	return false
}

// ProveConversion returns a ConversionProof given a list of input coins, a list of output coins,
// and a list of serial number witnesses.
func ProveConversion(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, snWitness []*serialnumbernoprivacy.SNNoPrivacyWitness) (*ConversionProof, error) {
	var err error
	proof := new(ConversionProof)
	proof.SetVersion(ProofVersion)
	if err = proof.SetInputCoins(inputCoins); err != nil {
		return nil, err
	}
	outputCoinsV2 := make([]coin.Coin, len(outputCoins))
	for i := 0; i < len(outputCoins); i += 1 {
		outputCoinsV2[i] = outputCoins[i]
	}
	if err = proof.SetOutputCoins(outputCoinsV2); err != nil {
		return nil, err
	}

	// Proving that serial number is derived from the committed derivator
	for i := 0; i < len(inputCoins); i++ {
		snNoPrivacyProof, err := snWitness[i].Prove(nil)
		if err != nil {
			return nil, err
		}
		proof.serialNumberNoPrivacyProof = append(proof.serialNumberNoPrivacyProof, snNoPrivacyProof)
	}
	// Hide the keyimage :D
	for i := 0; i < len(proof.outputCoins); i++ {
		proof.outputCoins[i].SetKeyImage(nil)
	}
	return proof, nil
}
