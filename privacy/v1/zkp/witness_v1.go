package zkp

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/bulletproofs"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/oneoutofmany"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/serialnumbernoprivacy"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/serialnumberprivacy"
)

// PaymentWitness contains all of witnesses to create a ProofV1.
type PaymentWitness struct {
	privateKey          *crypto.Scalar
	inputCoins          []coin.PlainCoin
	outputCoins         []*coin.CoinV1
	commitmentIndices   []uint64
	myCommitmentIndices []uint64

	oneOfManyWitness             []*oneoutofmany.OneOutOfManyWitness
	serialNumberWitness          []*serialnumberprivacy.SNPrivacyWitness
	serialNumberNoPrivacyWitness []*serialnumbernoprivacy.SNNoPrivacyWitness

	rangeProofWitness *bulletproofs.Witness

	comOutputValue                 []*crypto.Point
	comOutputSerialNumberDerivator []*crypto.Point
	comOutputShardID               []*crypto.Point

	comInputSecretKey             *crypto.Point
	comInputValue                 []*crypto.Point
	comInputSerialNumberDerivator []*crypto.Point
	comInputShardID               *crypto.Point

	randSecretKey *crypto.Scalar
}

// GetRandSecretKey returns the random factor in the commitment of the secret key.
func (w PaymentWitness) GetRandSecretKey() *crypto.Scalar {
	return w.randSecretKey
}

// PaymentWitnessParam consists of parameters for initializing a PaymentWitness.
type PaymentWitnessParam struct {
	HasPrivacy              bool
	PrivateKey              *crypto.Scalar
	InputCoins              []coin.PlainCoin
	OutputCoins             []*coin.CoinV1
	PublicKeyLastByteSender byte
	Commitments             []*crypto.Point
	CommitmentIndices       []uint64
	MyCommitmentIndices     []uint64
	Fee                     uint64
}

// Init creates a PaymentWitness from the given PaymentWitnessParam.
//
// If hasPrivacy = false, the returned PaymentWitness only consists of the spending key, input coins, output coins,
// and the non-private serial number proofs.
//
// Otherwise, it includes all attributes in the PaymentWitness.
func (w *PaymentWitness) Init(PaymentWitnessParam PaymentWitnessParam) error {
	_ = PaymentWitnessParam.Fee
	w.privateKey = PaymentWitnessParam.PrivateKey
	w.inputCoins = PaymentWitnessParam.InputCoins
	w.outputCoins = PaymentWitnessParam.OutputCoins
	w.commitmentIndices = PaymentWitnessParam.CommitmentIndices
	w.myCommitmentIndices = PaymentWitnessParam.MyCommitmentIndices

	randInputSK := crypto.RandomScalar()
	w.randSecretKey = new(crypto.Scalar).Set(randInputSK)

	if !PaymentWitnessParam.HasPrivacy {
		for _, outCoin := range w.outputCoins {
			outCoin.CoinDetails.SetRandomness(crypto.RandomScalar())
			err := outCoin.CoinDetails.CommitAll()
			if err != nil {
				return err
			}
		}
		lenInputs := len(w.inputCoins)
		if lenInputs > 0 {
			w.serialNumberNoPrivacyWitness = make([]*serialnumbernoprivacy.SNNoPrivacyWitness, lenInputs)
			for i := 0; i < len(w.inputCoins); i++ {
				/***** Build witness for proving that serial number is derived from the committed derivator *****/
				if w.serialNumberNoPrivacyWitness[i] == nil {
					w.serialNumberNoPrivacyWitness[i] = new(serialnumbernoprivacy.SNNoPrivacyWitness)
				}
				w.serialNumberNoPrivacyWitness[i].Set(w.inputCoins[i].GetKeyImage(), w.inputCoins[i].GetPublicKey(),
					w.inputCoins[i].GetSNDerivator(), w.privateKey)
			}
		}
		return nil
	}

	numInputCoin := len(w.inputCoins)
	numOutputCoin := len(w.outputCoins)

	cmInputSK := crypto.PedCom.CommitAtIndex(w.privateKey, randInputSK, crypto.PedersenPrivateKeyIndex)
	w.comInputSecretKey = new(crypto.Point).Set(cmInputSK)

	randInputShardID := FixedRandomnessShardID
	senderShardID := common.GetShardIDFromLastByte(PaymentWitnessParam.PublicKeyLastByteSender)
	w.comInputShardID = crypto.PedCom.CommitAtIndex(new(crypto.Scalar).FromUint64(uint64(senderShardID)), randInputShardID, crypto.PedersenShardIDIndex)

	w.comInputValue = make([]*crypto.Point, numInputCoin)
	w.comInputSerialNumberDerivator = make([]*crypto.Point, numInputCoin)

	randInputValue := make([]*crypto.Scalar, numInputCoin)
	randInputSND := make([]*crypto.Scalar, numInputCoin)

	// cmInputValueAll is sum of all input coins' value commitments
	cmInputValueAll := new(crypto.Point).Identity()
	randInputValueAll := new(crypto.Scalar).FromUint64(0)

	// Summing all commitments of each input coin into one commitment and proving the knowledge of its Openings
	cmInputSum := make([]*crypto.Point, numInputCoin)
	randInputSum := make([]*crypto.Scalar, numInputCoin)
	// randInputSumAll is sum of all randomness of coin commitments
	randInputSumAll := new(crypto.Scalar).FromUint64(0)

	w.oneOfManyWitness = make([]*oneoutofmany.OneOutOfManyWitness, numInputCoin)
	w.serialNumberWitness = make([]*serialnumberprivacy.SNPrivacyWitness, numInputCoin)

	commitmentTemps := make([][]*crypto.Point, numInputCoin)
	randInputIsZero := make([]*crypto.Scalar, numInputCoin)

	preIndex := 0
	commitments := PaymentWitnessParam.Commitments
	for i, inputCoin := range w.inputCoins {
		// tx only has fee, no output, Rand_Value_Input = 0
		if numOutputCoin == 0 {
			randInputValue[i] = new(crypto.Scalar).FromUint64(0)
		} else {
			randInputValue[i] = crypto.RandomScalar()
		}
		// commit each component of coin commitment
		randInputSND[i] = crypto.RandomScalar()

		w.comInputValue[i] = crypto.PedCom.CommitAtIndex(new(crypto.Scalar).FromUint64(inputCoin.GetValue()), randInputValue[i], crypto.PedersenValueIndex)
		w.comInputSerialNumberDerivator[i] = crypto.PedCom.CommitAtIndex(inputCoin.GetSNDerivator(), randInputSND[i], crypto.PedersenSndIndex)

		cmInputValueAll.Add(cmInputValueAll, w.comInputValue[i])
		randInputValueAll.Add(randInputValueAll, randInputValue[i])

		/***** Build witness for proving one-out-of-N commitments is a commitment to the coins being spent *****/
		cmInputSum[i] = new(crypto.Point).Add(cmInputSK, w.comInputValue[i])
		cmInputSum[i].Add(cmInputSum[i], w.comInputSerialNumberDerivator[i])
		cmInputSum[i].Add(cmInputSum[i], w.comInputShardID)

		randInputSum[i] = new(crypto.Scalar).Set(randInputSK)
		randInputSum[i].Add(randInputSum[i], randInputValue[i])
		randInputSum[i].Add(randInputSum[i], randInputSND[i])
		randInputSum[i].Add(randInputSum[i], randInputShardID)

		randInputSumAll.Add(randInputSumAll, randInputSum[i])

		// commitmentTemps is a list of commitments for protocol one-out-of-N
		commitmentTemps[i] = make([]*crypto.Point, utils.CommitmentRingSize)

		randInputIsZero[i] = new(crypto.Scalar).FromUint64(0)
		randInputIsZero[i].Sub(inputCoin.GetRandomness(), randInputSum[i])

		for j := 0; j < utils.CommitmentRingSize; j++ {
			commitmentTemps[i][j] = new(crypto.Point).Sub(commitments[preIndex+j], cmInputSum[i])
		}

		if w.oneOfManyWitness[i] == nil {
			w.oneOfManyWitness[i] = new(oneoutofmany.OneOutOfManyWitness)
		}
		indexIsZero := w.myCommitmentIndices[i] % utils.CommitmentRingSize

		w.oneOfManyWitness[i].Set(commitmentTemps[i], randInputIsZero[i], indexIsZero)
		preIndex = utils.CommitmentRingSize * (i + 1)
		// ---------------------------------------------------

		/***** Build witness for proving that serial number is derived from the committed derivator *****/
		if w.serialNumberWitness[i] == nil {
			w.serialNumberWitness[i] = new(serialnumberprivacy.SNPrivacyWitness)
		}
		stmt := new(serialnumberprivacy.SerialNumberPrivacyStatement)
		stmt.Set(inputCoin.GetKeyImage(), cmInputSK, w.comInputSerialNumberDerivator[i])
		w.serialNumberWitness[i].Set(stmt, w.privateKey, randInputSK, inputCoin.GetSNDerivator(), randInputSND[i])
		// ---------------------------------------------------
	}

	cmOutputSND := make([]*crypto.Point, numOutputCoin)
	cmOutputSum := make([]*crypto.Point, numOutputCoin)
	cmOutputValue := make([]*crypto.Point, numOutputCoin)
	randOutputSum := make([]*crypto.Scalar, numOutputCoin)
	randOutputSND := make([]*crypto.Scalar, numOutputCoin)
	randOutputValue := make([]*crypto.Scalar, numOutputCoin)

	// cmOutputValueAll is sum of all value coin commitments
	cmOutputSumAll := new(crypto.Point).Identity()
	cmOutputValueAll := new(crypto.Point).Identity()
	randOutputValueAll := new(crypto.Scalar).FromUint64(0)

	cmOutputShardID := make([]*crypto.Point, numOutputCoin)
	randOutputShardID := make([]*crypto.Scalar, numOutputCoin)

	outputCoins := w.outputCoins
	for i, outputCoin := range outputCoins {
		if i == len(outputCoins)-1 {
			randOutputValue[i] = new(crypto.Scalar).Sub(randInputValueAll, randOutputValueAll)
		} else {
			randOutputValue[i] = crypto.RandomScalar()
		}

		randOutputSND[i] = crypto.RandomScalar()
		randOutputShardID[i] = crypto.RandomScalar()

		cmOutputValue[i] = crypto.PedCom.CommitAtIndex(new(crypto.Scalar).FromUint64(outputCoin.CoinDetails.GetValue()), randOutputValue[i], crypto.PedersenValueIndex)
		cmOutputSND[i] = crypto.PedCom.CommitAtIndex(outputCoin.CoinDetails.GetSNDerivator(), randOutputSND[i], crypto.PedersenSndIndex)

		receiverShardID, err := outputCoins[i].GetShardID()
		if err != nil {
			return fmt.Errorf("cannot parse shardID of outputCoins")
		}

		cmOutputShardID[i] = crypto.PedCom.CommitAtIndex(new(crypto.Scalar).FromUint64(uint64(receiverShardID)), randOutputShardID[i], crypto.PedersenShardIDIndex)

		randOutputSum[i] = new(crypto.Scalar).FromUint64(0)
		randOutputSum[i].Add(randOutputValue[i], randOutputSND[i])
		randOutputSum[i].Add(randOutputSum[i], randOutputShardID[i])

		cmOutputSum[i] = new(crypto.Point).Identity()
		cmOutputSum[i].Add(cmOutputValue[i], cmOutputSND[i])
		cmOutputSum[i].Add(cmOutputSum[i], outputCoins[i].GetPublicKey())
		cmOutputSum[i].Add(cmOutputSum[i], cmOutputShardID[i])

		cmOutputValueAll.Add(cmOutputValueAll, cmOutputValue[i])
		randOutputValueAll.Add(randOutputValueAll, randOutputValue[i])

		// calculate final commitment for output coins
		outputCoins[i].CoinDetails.SetCommitment(cmOutputSum[i])
		outputCoins[i].CoinDetails.SetRandomness(randOutputSum[i])

		cmOutputSumAll.Add(cmOutputSumAll, cmOutputSum[i])
	}

	// For Multi Range Protocol
	// proving each output value is less than vMax
	// proving sum of output values is less than vMax
	outputValue := make([]uint64, numOutputCoin)
	for i := 0; i < numOutputCoin; i++ {
		if outputCoins[i].CoinDetails.GetValue() >= 0 {
			outputValue[i] = outputCoins[i].CoinDetails.GetValue()
		} else {
			return fmt.Errorf("output coin's value is less than 0")
		}
	}
	if w.rangeProofWitness == nil {
		w.rangeProofWitness = new(bulletproofs.Witness)
	}
	w.rangeProofWitness.Set(outputValue, randOutputValue)
	// ---------------------------------------------------

	// save partial commitments (value, input, shardID)
	w.comOutputValue = cmOutputValue
	w.comOutputSerialNumberDerivator = cmOutputSND
	w.comOutputShardID = cmOutputShardID

	return nil
}

// Prove creates big ProofV1.
func (w *PaymentWitness) Prove(hasPrivacy bool, paymentInfo []*key.PaymentInfo) (*ProofV1, error) {
	proof := new(ProofV1)
	proof.Init()

	proof.inputCoins = w.inputCoins
	proof.outputCoins = w.outputCoins
	proof.commitmentOutputValue = w.comOutputValue
	proof.commitmentOutputSND = w.comOutputSerialNumberDerivator
	proof.commitmentOutputShardID = w.comOutputShardID

	proof.commitmentInputSecretKey = w.comInputSecretKey
	proof.commitmentInputValue = w.comInputValue
	proof.commitmentInputSND = w.comInputSerialNumberDerivator
	proof.commitmentInputShardID = w.comInputShardID
	proof.commitmentIndices = w.commitmentIndices

	// if hasPrivacy == false, don't need to create the zero knowledge proof
	// proving user has spending key corresponding with public key in input coins
	// is proved by signing with spending key
	if !hasPrivacy {
		// Proving that serial number is derived from the committed derivator
		for i := 0; i < len(w.inputCoins); i++ {
			snNoPrivacyProof, err := w.serialNumberNoPrivacyWitness[i].Prove(nil)
			if err != nil {
				return nil, err
			}
			proof.serialNumberNoPrivacyProof = append(proof.serialNumberNoPrivacyProof, snNoPrivacyProof)
		}
		for i := 0; i < len(proof.outputCoins); i++ {
			proof.outputCoins[i].CoinDetails.SetKeyImage(nil)
		}
		return proof, nil
	}

	// if hasPrivacy == true
	numInputCoins := len(w.oneOfManyWitness)

	for i := 0; i < numInputCoins; i++ {
		// Proving one-out-of-N commitments is a commitment to the coins being spent
		oneOfManyProof, err := w.oneOfManyWitness[i].Prove()
		if err != nil {
			return nil, err
		}
		proof.oneOfManyProof = append(proof.oneOfManyProof, oneOfManyProof)

		// Proving that serial number is derived from the committed derivator
		serialNumberProof, err := w.serialNumberWitness[i].Prove(nil)
		if err != nil {
			return nil, err
		}
		proof.serialNumberProof = append(proof.serialNumberProof, serialNumberProof)
	}
	var err error

	// Proving that each output values and sum of them does not exceed v_max
	proof.rangeProof, err = w.rangeProofWitness.Prove()
	if err != nil {
		return nil, err
	}

	if len(proof.inputCoins) == 0 {
		proof.commitmentIndices = nil
		proof.commitmentInputSecretKey = nil
		proof.commitmentInputShardID = nil
		proof.commitmentInputSND = nil
		proof.commitmentInputValue = nil
	}

	if len(proof.outputCoins) == 0 {
		proof.commitmentOutputValue = nil
		proof.commitmentOutputSND = nil
		proof.commitmentOutputShardID = nil
	}

	// After Prove, we should hide all information in coin details.
	// encrypt coin details (Randomness)
	// hide information of output coins except coin commitments, public key, snDerivator
	for i := 0; i < len(proof.outputCoins); i++ {
		errEncrypt := proof.outputCoins[i].Encrypt(paymentInfo[i].PaymentAddress.Tk)
		if errEncrypt != nil {
			return nil, errEncrypt
		}
		proof.outputCoins[i].CoinDetails.SetKeyImage(nil)
		proof.outputCoins[i].CoinDetails.SetValue(0)
		proof.outputCoins[i].CoinDetails.SetRandomness(nil)
	}

	for i := 0; i < len(proof.GetInputCoins()); i++ {
		err = proof.inputCoins[i].ConcealOutputCoin(nil)
		if err != nil {
			return nil, err
		}
	}
	return proof, nil
}
