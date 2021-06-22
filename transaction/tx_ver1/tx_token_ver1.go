package tx_ver1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"math"
)

type TxToken struct {
	tx_generic.TxTokenBase
}

func (txToken *TxToken) Init(paramsInterface interface{}) error {
	params, ok := paramsInterface.(*tx_generic.TxTokenParams)
	if !ok {
		return errors.New("Cannot init TxTokenBase because params is not correct")
	}
	// init data for tx PRV for fee
	txPrivacyParams := tx_generic.NewTxPrivacyInitParams(
		params.SenderKey,
		params.PaymentInfo,
		params.InputCoin,
		params.FeeNativeCoin,
		params.HasPrivacyCoin,
		nil,
		params.MetaData,
		params.Info,
		params.Kvargs)
	txToken.Tx = new(Tx)
	if err := txToken.Tx.Init(txPrivacyParams); err != nil {
		return err
	}
	// override TxCustomTokenPrivacyType type
	txToken.Tx.SetType(common.TxCustomTokenPrivacyType)

	// check action type and create privacy custom toke data
	var handled = false
	// Add token data component
	txToken.TxTokenData.SetType(params.TokenParams.TokenTxType)
	txToken.TxTokenData.SetPropertyName(params.TokenParams.PropertyName)
	txToken.TxTokenData.SetPropertySymbol(params.TokenParams.PropertySymbol)

	switch params.TokenParams.TokenTxType {
	case utils.CustomTokenInit: {
		// case init a new privacy custom token
		handled = true
		txToken.TxTokenData.SetAmount(params.TokenParams.Amount)

		temp := new(Tx)
		temp.SetVersion(utils.TxVersion1Number)
		temp.Type = common.TxNormalType
		temp.Proof = new(privacy.ProofV1)
		tempOutputCoin := make([]*coin.CoinV1, 1)
		tempOutputCoin[0] = new(coin.CoinV1)
		tempOutputCoin[0].CoinDetails = new(coin.PlainCoinV1)
		tempOutputCoin[0].CoinDetails.SetValue(params.TokenParams.Amount)
		PK, err := new(crypto.Point).FromBytesS(params.TokenParams.Receiver[0].PaymentAddress.Pk)
		if err != nil {
			return err
		}
		tempOutputCoin[0].CoinDetails.SetPublicKey(PK)
		tempOutputCoin[0].CoinDetails.SetRandomness(crypto.RandomScalar())

		// set info coin for output coin
		if len(params.TokenParams.Receiver[0].Message) > 0 {
			if len(params.TokenParams.Receiver[0].Message) > coin.MaxSizeInfoCoin {
				return fmt.Errorf("len of message (%v) too large", len(params.TokenParams.Receiver[0].Message))
			}
			tempOutputCoin[0].CoinDetails.SetInfo(params.TokenParams.Receiver[0].Message)
		}
		tempOutputCoin[0].CoinDetails.SetSNDerivator(crypto.RandomScalar())
		err = tempOutputCoin[0].CoinDetails.CommitAll()
		if err != nil {
			return err
		}
		outputCoinsAsGeneric := make([]coin.Coin, len(tempOutputCoin))
		for i := 0; i < len(tempOutputCoin); i += 1 {
			outputCoinsAsGeneric[i] = tempOutputCoin[i]
		}
		temp.Proof.SetOutputCoins(outputCoinsAsGeneric)

		// get last byte
		lastBytes := params.TokenParams.Receiver[0].PaymentAddress.Pk[len(params.TokenParams.Receiver[0].PaymentAddress.Pk)-1]
		temp.PubKeyLastByteSender = common.GetShardIDFromLastByte(lastBytes)

		// signOnMessage Tx
		temp.SigPubKey = params.TokenParams.Receiver[0].PaymentAddress.Pk
		temp.SetPrivateKey(*params.SenderKey)
		err = temp.sign()
		if err != nil {
			return err
		}
		txToken.TxTokenData.TxNormal = temp

		hashInitToken, err := txToken.TxTokenData.Hash()
		if err != nil {
			return err
		}

		if params.TokenParams.Mintable {
			propertyID, err := common.Hash{}.NewHashFromStr(params.TokenParams.PropertyID)
			if err != nil {
				return err
			}
			txToken.TxTokenData.PropertyID = *propertyID
			txToken.TxTokenData.Mintable = true
		} else {
			//NOTICE: @merman update PropertyID calculated from hash of tokendata and shardID
			newHashInitToken := common.HashH(append(hashInitToken.GetBytes(), params.ShardID))
			txToken.TxTokenData.PropertyID = newHashInitToken
		}
	}
	case utils.CustomTokenTransfer: {
		handled = true
		// make a transfering for privacy custom token
		// fee always 0 and reuse function of normal tx for custom token ID
		propertyID, _ := common.Hash{}.NewHashFromStr(params.TokenParams.PropertyID)

		txToken.TxTokenData.SetPropertyID(*propertyID)
		txToken.TxTokenData.SetMintable(params.TokenParams.Mintable)

		txToken.TxTokenData.TxNormal = new(Tx)
		err := txToken.TxTokenData.TxNormal.Init(tx_generic.NewTxPrivacyInitParams(params.SenderKey,
			params.TokenParams.Receiver,
			params.TokenParams.TokenInput,
			params.TokenParams.Fee,
			params.HasPrivacyToken,
			propertyID,
			nil,
			nil,
			params.TokenParams.Kvargs))
		if err != nil {
			fmt.Printf("Init PRV fee transaction returns an error: %v\n", err)
			return err
		}
	}
	}
	if !handled {
		return errors.New("can't handle this TokenTxType")
	}
	return nil
}

func (txToken TxToken) GetTxActualSize() uint64 {
	normalTxSize := txToken.Tx.GetTxActualSize()
	tokenDataSize := uint64(0)
	tokenDataSize += txToken.TxTokenData.TxNormal.GetTxActualSize()
	tokenDataSize += uint64(len(txToken.TxTokenData.PropertyName))
	tokenDataSize += uint64(len(txToken.TxTokenData.PropertySymbol))
	tokenDataSize += uint64(len(txToken.TxTokenData.PropertyID))
	tokenDataSize += 4 // for TxPrivacyTokenDataVersion1.Type
	tokenDataSize += 8 // for TxPrivacyTokenDataVersion1.Amount
	meta := txToken.GetMetadata()
	if meta != nil {
		tokenDataSize += meta.CalculateSize()
	}
	return normalTxSize + uint64(math.Ceil(float64(tokenDataSize)/1024))
}

func (txToken *TxToken) UnmarshalJSON(data []byte) error {
	var err error
	txToken.Tx = &Tx{}
	if err = json.Unmarshal(data, txToken.Tx); err != nil {
		return err
	}
	temp := &struct {
		TxTokenData tx_generic.TxTokenData `json:"TxTokenPrivacyData"`
	}{}
	temp.TxTokenData.TxNormal = &Tx{}
	err = json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	txToken.TxTokenData = temp.TxTokenData
	if txToken.Tx.GetMetadata() != nil && txToken.Tx.GetMetadata().GetType() == 81 {
		if txToken.TxTokenData.Amount == 37772966455153490 {
			txToken.TxTokenData.Amount = 37772966455153487
		}
	}
	return nil
}

func (txToken TxToken) ListOTAHashH() []common.Hash {
	return []common.Hash{}
}