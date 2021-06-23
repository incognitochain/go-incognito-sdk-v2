package tx_ver1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"math/big"
)

// Tx represents a PRV transaction of version 1. It is a embedded TxBase with some overridden functions.
// A transaction version is mainly composed of
//	- Schnorr signature: to sign the transaction.
//	- One-of-many proof: to anonymize the true sender.
//	- BulletProofs: a range proof used to prove that a value lies within an interval without revealing it.
//	- SerialNumberProof: to prove the correctness of serial numbers.
//	- SerialNumberProof (no privacy): to prove the correctness of serial numbers used in non-private transactions.
type Tx struct {
	tx_generic.TxBase
}

// GetReceiverData returns a list of output coins (not including sent-back coins) of a Tx.
func (tx *Tx) GetReceiverData() ([]coin.Coin, error) {
	pubKeys := make([]*crypto.Point, 0)
	amounts := make([]uint64, 0)

	if tx.Proof != nil && len(tx.Proof.GetOutputCoins()) > 0 {
		for _, c := range tx.Proof.GetOutputCoins() {
			coinPubKey := c.GetPublicKey()
			added := false
			for i, k := range pubKeys {
				if bytes.Equal(coinPubKey.ToBytesS(), k.ToBytesS()) {
					added = true
					amounts[i] += c.GetValue()
					break
				}
			}
			if !added {
				pubKeys = append(pubKeys, coinPubKey)
				amounts = append(amounts, c.GetValue())
			}
		}
	}
	coins := make([]coin.Coin, 0)
	for i := 0; i < len(pubKeys); i++ {
		c := new(coin.CoinV1).Init()
		c.CoinDetails.SetPublicKey(pubKeys[i])
		c.CoinDetails.SetValue(amounts[i])
		coins = append(coins, c)
	}
	return coins, nil
}

// GetTxMintData returns the minting data of a Tx.
func (tx Tx) GetTxMintData() (bool, coin.Coin, *common.Hash, error) {
	return tx_generic.GetTxMintData(&tx, &common.PRVCoinID)
}

// GetTxBurnData returns the burning data (token only) of a Tx.
func (tx Tx) GetTxBurnData() (bool, coin.Coin, *common.Hash, error) {
	return tx_generic.GetTxBurnData(&tx)
}

// GetTxFullBurnData is the same as GetTxBurnData.
func (tx Tx) GetTxFullBurnData() (bool, coin.Coin, coin.Coin, *common.Hash, error) {
	isBurn, burnedCoin, burnedTokenID, err := tx.GetTxBurnData()
	return isBurn, burnedCoin, nil, burnedTokenID, err
}

// CheckAuthorizedSender checks if the sender of a Tx is authorized w.r.t to a public key.
func (tx *Tx) CheckAuthorizedSender(publicKey []byte) (bool, error) {
	sigPubKey := tx.GetSigPubKey()
	if bytes.Equal(sigPubKey, publicKey) {
		return true, nil
	} else {
		return false, nil
	}
}

// Init creates a PRV transaction version 1 from the given parameter.
// The input parameter should be a *tx_generic.TxPrivacyInitParams.
func (tx *Tx) Init(txParams interface{}) error {
	params, ok := txParams.(*tx_generic.TxPrivacyInitParams)
	if !ok {
		return fmt.Errorf("cannot parse the input as a TxPrivacyInitParams")
	}

	if err := tx_generic.ValidateTxParams(params); err != nil {
		return err
	}

	// Init tx and params (tx and params will be changed)
	if err := tx.InitializeTxAndParams(params); err != nil {
		return err
	}
	tx.SetVersion(utils.TxVersion1Number)

	if check, err := tx.IsNonPrivacyNonInput(params); check {
		return err
	}

	if err := tx.prove(params); err != nil {
		return err
	}
	return nil
}

// Sign re-signs a Tx using the given private key.
func (tx *Tx) Sign(sigPrivateKey []byte) error {
	if sigPrivateKey != nil {
		tx.SetPrivateKey(sigPrivateKey)
	}
	return tx.sign()
}

func (tx *Tx) prove(params *tx_generic.TxPrivacyInitParams) error {
	// PrepareTransaction paymentWitness params
	paymentWitnessParamPtr, err := tx.initPaymentWitnessParam(params)
	if err != nil {
		return err
	}
	return tx.proveAndSignCore(params, paymentWitnessParamPtr)
}

func (tx *Tx) sign() error {
	//Check input transaction
	if tx.Sig != nil {
		return fmt.Errorf("input transaction must be an unsigned one")
	}

	sk := new(crypto.Scalar).FromBytesS(tx.GetPrivateKey()[:common.BigIntSize])
	r := new(crypto.Scalar).FromBytesS(tx.GetPrivateKey()[common.BigIntSize:])
	sigKey := new(privacy.SchnorrPrivateKey)
	sigKey.Set(sk, r)

	// save public key for verification signature tx
	tx.SigPubKey = sigKey.GetPublicKey().GetPublicKey().ToBytesS()

	// signing
	signature, err := sigKey.Sign(tx.Hash()[:])
	if err != nil {
		return err
	}

	// convert signature to byte array
	tx.Sig = signature.Bytes()

	return nil
}

func (tx *Tx) initPaymentWitnessParam(params *tx_generic.TxPrivacyInitParams) (*privacy.PaymentWitnessParam, error) {
	var commitmentIndices []uint64
	var inputCoinCommitmentIndices []uint64
	var commitments []*crypto.Point

	if params.HasPrivacy && len(params.InputCoins) > 0 {
		//Get list of decoy indices.
		tmp, ok := params.KvArgs[utils.CommitmentIndices]
		if !ok {
			return nil, fmt.Errorf("decoy commitment indices not found: %v", params.KvArgs)
		}

		commitmentIndices, ok = tmp.([]uint64)
		if !ok {
			return nil, fmt.Errorf("cannot parse commitment indices: %v", tmp)
		}

		//Get list of decoy commitments.
		tmp, ok = params.KvArgs[utils.Commitments]
		if !ok {
			return nil, fmt.Errorf("decoy commitment list not found: %v", params.KvArgs)
		}

		commitments, ok = tmp.([]*crypto.Point)
		if !ok {
			return nil, fmt.Errorf("cannot parse sender commitment indices: %v", tmp)
		}

		//Get list of input coin indices
		tmp, ok = params.KvArgs[utils.MyIndices]
		if !ok {
			return nil, fmt.Errorf("inputCoin commitment indices not found: %v", params.KvArgs)
		}

		inputCoinCommitmentIndices, ok = tmp.([]uint64)
		if !ok {
			return nil, fmt.Errorf("cannot parse inputCoin commitment indices: %v", tmp)
		}
	}

	outputCoins, err := generateOutputCoinV1s(params.PaymentInfo)
	if err != nil {
		return nil, err
	}

	// prepare witness for proving
	paymentWitnessParam := privacy.PaymentWitnessParam{
		HasPrivacy:              params.HasPrivacy,
		PrivateKey:              new(crypto.Scalar).FromBytesS(*params.SenderSK),
		InputCoins:              params.InputCoins,
		OutputCoins:             outputCoins,
		PublicKeyLastByteSender: common.GetShardIDFromLastByte(tx.PubKeyLastByteSender),
		Commitments:             commitments,
		CommitmentIndices:       commitmentIndices,
		MyCommitmentIndices:     inputCoinCommitmentIndices,
		Fee:                     params.Fee,
	}
	return &paymentWitnessParam, nil
}

func (tx *Tx) proveAndSignCore(params *tx_generic.TxPrivacyInitParams, paymentWitnessParamPtr *privacy.PaymentWitnessParam) error {
	paymentWitnessParam := *paymentWitnessParamPtr
	witness := new(privacy.PaymentWitness)
	err := witness.Init(paymentWitnessParam)
	if err != nil {
		jsonParam, _ := json.MarshalIndent(paymentWitnessParam, common.EmptyString, "  ")
		return fmt.Errorf("witnessParam init error. Param %v, error %v", string(jsonParam), err)
	}

	paymentProof, err := witness.Prove(params.HasPrivacy, params.PaymentInfo)
	if err != nil {
		jsonParam, _ := json.MarshalIndent(paymentWitnessParam, common.EmptyString, "  ")
		return fmt.Errorf("witnessParam prove error. Param %v, error %v", string(jsonParam), err)
	}
	tx.Proof = paymentProof

	// set private key for signing tx
	if params.HasPrivacy {
		randSK := witness.GetRandSecretKey()
		tx.SetPrivateKey(append(*params.SenderSK, randSK.ToBytesS()...))
	} else {
		tx.SetPrivateKey([]byte{})
		randSK := big.NewInt(0)
		tx.SetPrivateKey(append(*params.SenderSK, randSK.Bytes()...))
	}

	// sign tx
	signErr := tx.sign()
	if signErr != nil {
		return fmt.Errorf("tx sign error %v", err)
	}
	return nil
}

func generateOutputCoinV1s(paymentInfo []*key.PaymentInfo) ([]*coin.CoinV1, error) {
	outputCoins := make([]*coin.CoinV1, len(paymentInfo))
	for i, pInfo := range paymentInfo {
		outputCoins[i] = new(coin.CoinV1)
		outputCoins[i].CoinDetails = new(coin.PlainCoinV1)
		outputCoins[i].CoinDetails.SetValue(pInfo.Amount)
		if len(pInfo.Message) > 0 {
			if len(pInfo.Message) > coin.MaxSizeInfoCoin {
				return nil, fmt.Errorf("length of message (%v) too large", len(pInfo.Message))
			}
		}
		outputCoins[i].CoinDetails.SetInfo(pInfo.Message)

		PK, err := new(crypto.Point).FromBytesS(pInfo.PaymentAddress.Pk)
		if err != nil {
			return nil, fmt.Errorf("can not decompress public key from %v: %v", pInfo.PaymentAddress, err)
		}
		outputCoins[i].CoinDetails.SetPublicKey(PK)
		outputCoins[i].CoinDetails.SetSNDerivator(crypto.RandomScalar())
	}
	return outputCoins, nil
}
