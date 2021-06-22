package tx_ver1

import (
	"bytes"
	"encoding/json"
	"errors"
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

type Tx struct {
	tx_generic.TxBase
}

//GETTER FUNCTIONS
func (tx *Tx) GetReceiverData() ([]coin.Coin, error) {
	pubkeys := make([]*crypto.Point, 0)
	amounts := []uint64{}

	if tx.Proof != nil && len(tx.Proof.GetOutputCoins()) > 0 {
		for _, coin := range tx.Proof.GetOutputCoins() {
			coinPubKey := coin.GetPublicKey()
			added := false
			for i, key := range pubkeys {
				if bytes.Equal(coinPubKey.ToBytesS(), key.ToBytesS()) {
					added = true
					amounts[i] += coin.GetValue()
					break
				}
			}
			if !added {
				pubkeys = append(pubkeys, coinPubKey)
				amounts = append(amounts, coin.GetValue())
			}
		}
	}
	coins := make([]coin.Coin, 0)
	for i := 0; i < len(pubkeys); i++ {
		coin := new(coin.CoinV1).Init()
		coin.CoinDetails.SetPublicKey(pubkeys[i])
		coin.CoinDetails.SetValue(amounts[i])
		coins = append(coins, coin)
	}
	return coins, nil
}

func (tx Tx) GetTxMintData() (bool, coin.Coin, *common.Hash, error) { return tx_generic.GetTxMintData(&tx, &common.PRVCoinID) }

func (tx Tx) GetTxBurnData() (bool, coin.Coin, *common.Hash, error) { return tx_generic.GetTxBurnData(&tx) }

func (tx Tx) GetTxFullBurnData() (bool, coin.Coin, coin.Coin, *common.Hash, error) {
	isBurn, burnedCoin, burnedTokenID, err := tx.GetTxBurnData()
	return isBurn, burnedCoin, nil, burnedTokenID, err
}
//END GETTER FUNCTIONS

//INIT FUNCTIONS
func (tx *Tx) Init(paramsInterface interface{}) error {
	params, ok := paramsInterface.(*tx_generic.TxPrivacyInitParams)
	if !ok {
		return errors.New("params of tx Init is not TxPrivacyInitParam")
	}

	if err := tx_generic.ValidateTxParams(params); err != nil {
		return err
	}

	// Init tx and params (tx and params will be changed)
	if err := tx.InitializeTxAndParams(params); err != nil {
		return err
	}
	tx.SetVersion(utils.TxVersion1Number)

	// Check if this tx is nonPrivacyNonInput
	// Case 1: tx ptoken transfer with ptoken fee
	// Case 2: tx Reward
	if check, err := tx.IsNonPrivacyNonInput(params); check {
		return err
	}

	if err := tx.prove(params); err != nil {
		return err
	}
	return nil
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
		return errors.New("input transaction must be an unsigned one")
	}

	/****** using Schnorr signature *******/
	// sign with sigPrivKey
	// prepare private key for Schnorr
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

func (tx *Tx) Sign(sigPrivakey []byte) error {//For testing-purpose only, remove when deploy
	if sigPrivakey != nil{
		tx.SetPrivateKey(sigPrivakey)
	}
	return tx.sign()
}
//END INIT FUNCTIONS

//HELPER FUNCTIONS
func GenerateOutputCoinV1s(paymentInfo []*key.PaymentInfo) ([]*coin.CoinV1, error) {
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
			return nil, fmt.Errorf("can not decompress public key from %v. Error: %v.", pInfo.PaymentAddress, err)
		}
		outputCoins[i].CoinDetails.SetPublicKey(PK)
		outputCoins[i].CoinDetails.SetSNDerivator(crypto.RandomScalar())
	}
	return outputCoins, nil
}

func (tx *Tx) initPaymentWitnessParam(params *tx_generic.TxPrivacyInitParams) (*privacy.PaymentWitnessParam, error) {
	var commitmentIndices []uint64
	var inputCoinCommitmentIndices []uint64
	var commitments []*crypto.Point

	if params.HasPrivacy && len(params.InputCoins) > 0{
		//Get list of decoy indices.
		tmp, ok := params.Kvargs[utils.CommitmentIndices]
		if !ok {
			return nil, fmt.Errorf("decoy commitment indices not found: %v", params.Kvargs)
		}

		commitmentIndices, ok = tmp.([]uint64)
		if !ok {
			return nil,  fmt.Errorf("cannot parse commitment indices: %v", tmp)
		}

		//Get list of decoy commitments.
		tmp, ok = params.Kvargs[utils.Commitments]
		if !ok {
			return nil,  fmt.Errorf("decoy commitment list not found: %v", params.Kvargs)
		}

		commitments, ok = tmp.([]*crypto.Point)
		if !ok {
			return nil,   fmt.Errorf("cannot parse sender commitment indices: %v", tmp)
		}

		//Get list of inputcoin indices
		tmp, ok = params.Kvargs[utils.MyIndices]
		if !ok {
			return nil,  fmt.Errorf("inputCoin commitment indices not found: %v", params.Kvargs)
		}

		inputCoinCommitmentIndices, ok = tmp.([]uint64)
		if !ok {
			return nil,  fmt.Errorf("cannot parse inputCoin commitment indices: %v", tmp)
		}
	}

	outputCoins, err := GenerateOutputCoinV1s(params.PaymentInfo)
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

func (tx *Tx) CheckAuthorizedSender(publicKey []byte) (bool, error) {
	sigPubKey := tx.GetSigPubKey()
	if bytes.Equal(sigPubKey, publicKey) {
		return true, nil
	} else {
		return false, nil
	}
}
//END HELPER FUNCTIONS