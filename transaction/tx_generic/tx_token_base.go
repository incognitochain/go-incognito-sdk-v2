package tx_generic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"sort"
)

type Tx = metadata.Transaction

type TxTokenBase struct {
	Tx
	TxTokenData TxTokenData `json:"TxTokenPrivacyData"`
	cachedHash  *common.Hash
}

type TxTokenParams struct {
	SenderKey       *key.PrivateKey
	PaymentInfo     []*key.PaymentInfo
	InputCoin       []coin.PlainCoin
	FeeNativeCoin   uint64
	TokenParams     *TokenParam
	MetaData        metadata.Metadata
	HasPrivacyCoin  bool
	HasPrivacyToken bool
	ShardID         byte
	Info            []byte
	Kvargs          map[string]interface{}
}

// CustomTokenParamTx - use for rpc request json body
type TokenParam struct {
	PropertyID     string                 `json:"TokenID"`
	PropertyName   string                 `json:"TokenName"`
	PropertySymbol string                 `json:"TokenSymbol"`
	Amount         uint64                 `json:"TokenAmount"`
	TokenTxType    int                    `json:"TokenTxType"`
	Receiver       []*key.PaymentInfo `json:"TokenReceiver"`
	TokenInput     []coin.PlainCoin    `json:"TokenInput"`
	Mintable       bool                   `json:"TokenMintable"`
	Fee            uint64                 `json:"TokenFee"`
	Kvargs         map[string]interface{}
}

func NewTokenParam(propertyID, propertyName, propertySymbol string,
	amount uint64,
	tokenTxType int,
	receivers []*key.PaymentInfo,
	tokenInput []coin.PlainCoin,
	mintable bool,
	fee uint64,
	kvargs map[string]interface{}) *TokenParam {

	params := &TokenParam{
		PropertyID:     propertyID,
		PropertyName:   propertyName,
		PropertySymbol: propertySymbol,
		Amount:         amount,
		TokenTxType:    tokenTxType,
		Receiver:       receivers,
		TokenInput:     tokenInput,
		Mintable:       mintable,
		Fee:            fee,
		Kvargs:         kvargs,
	}
	return params
}

func NewTxTokenParams(senderKey *key.PrivateKey,
	paymentInfo []*key.PaymentInfo,
	inputCoin []coin.PlainCoin,
	feeNativeCoin uint64,
	tokenParams *TokenParam,
	metaData metadata.Metadata,
	hasPrivacyCoin bool,
	hasPrivacyToken bool,
	shardID byte,
	info []byte,
	kargs map[string]interface{}) *TxTokenParams {
	params := &TxTokenParams{
		ShardID:         shardID,
		PaymentInfo:     paymentInfo,
		MetaData:        metaData,
		FeeNativeCoin:   feeNativeCoin,
		HasPrivacyCoin:  hasPrivacyCoin,
		HasPrivacyToken: hasPrivacyToken,
		InputCoin:       inputCoin,
		SenderKey:       senderKey,
		TokenParams:     tokenParams,
		Info:            info,
		Kvargs:          kargs,
	}
	return params
}

// ========== Get/Set FUNCTION ============

func (txToken TxTokenBase) GetTxBase() metadata.Transaction { return txToken.Tx }
func (txToken *TxTokenBase) SetTxBase(tx metadata.Transaction) error {
	txToken.Tx = tx
	return nil
}
func (txToken TxTokenBase) GetTxNormal() metadata.Transaction { return txToken.TxTokenData.TxNormal }
func (txToken *TxTokenBase) SetTxNormal(tx metadata.Transaction) error {
	txToken.TxTokenData.TxNormal = tx
	return nil
}
func (txToken TxTokenBase) GetTxTokenData() TxTokenData { return txToken.TxTokenData }
func (txToken *TxTokenBase) SetTxTokenData(data TxTokenData) error {
	txToken.TxTokenData = data
	return nil
}

func (txToken TxTokenBase) GetTxMintData() (bool, coin.Coin, *common.Hash, error) {
	tokenID := txToken.TxTokenData.GetPropertyID()
	return GetTxMintData(txToken.TxTokenData.TxNormal, &tokenID)
}

func (txToken TxTokenBase) GetTxBurnData() (bool, coin.Coin, *common.Hash, error) {
	tokenID := txToken.TxTokenData.GetPropertyID()
	isBurn, burnCoin, _, err := txToken.TxTokenData.TxNormal.GetTxBurnData()
	return isBurn, burnCoin, &tokenID, err
}

func (txToken TxTokenBase) GetTxFullBurnData() (bool, coin.Coin, coin.Coin, *common.Hash, error) {
	isBurnToken, burnToken, burnedTokenID, errToken := txToken.GetTxBurnData()
	isBurnPrv, burnPrv, _, errPrv := txToken.GetTxBase().GetTxBurnData()

	if errToken != nil && errPrv != nil {
		return false, nil, nil, nil, fmt.Errorf("%v and %v", errPrv, errToken)
	}

	return isBurnPrv || isBurnToken, burnPrv, burnToken, burnedTokenID, nil
}

// ========== CHECK FUNCTION ===========

func (txToken TxTokenBase) CheckAuthorizedSender(publicKey []byte) (bool, error) {
	sigPubKey := txToken.TxTokenData.TxNormal.GetSigPubKey()
	if bytes.Equal(sigPubKey, publicKey) {
		return true, nil
	} else {
		return false, nil
	}
}

func (txToken TxTokenBase) IsSalaryTx() bool {
	if txToken.GetType() != common.TxRewardType {
		return false
	}
	if txToken.GetProof() != nil {
		return false
	}
	if len(txToken.TxTokenData.TxNormal.GetProof().GetInputCoins()) > 0 {
		return false
	}
	return true
}

// ==========  PARSING JSON FUNCTIONS ==========

func (txToken TxTokenBase) MarshalJSON() ([]byte, error) {
	type TemporaryTxToken struct {
		TxBase
		TxTokenData TxTokenData `json:"TxTokenPrivacyData"`
	}
	tempTx := TemporaryTxToken{}
	tempTx.TxTokenData = txToken.GetTxTokenData()
	tx := txToken.GetTxBase()
	if tx == nil {
		return nil, errors.New("Cannot unmarshal transaction: txfee cannot be nil")
	}
	tempTx.TxBase.SetVersion(tx.GetVersion())
	tempTx.TxBase.SetType(tx.GetType())
	tempTx.TxBase.SetLockTime(tx.GetLockTime())
	tempTx.TxBase.SetTxFee(tx.GetTxFee())
	tempTx.TxBase.SetInfo(tx.GetInfo())
	tempTx.TxBase.SetSigPubKey(tx.GetSigPubKey())
	tempTx.TxBase.SetSig(tx.GetSig())
	tempTx.TxBase.SetProof(tx.GetProof())
	tempTx.TxBase.SetGetSenderAddrLastByte(tx.GetSenderAddrLastByte())
	tempTx.TxBase.SetMetadata(tx.GetMetadata())
	tempTx.TxBase.SetGetSenderAddrLastByte(tx.GetSenderAddrLastByte())

	return json.Marshal(tempTx)
}

func (txToken TxTokenBase) String() string {
	// get hash of tx
	record := txToken.Tx.Hash().String()
	// add more hash of tx custom token data privacy
	tokenPrivacyDataHash, _ := txToken.TxTokenData.Hash()
	record += tokenPrivacyDataHash.String()

	meta := txToken.GetMetadata()
	if meta != nil {
		record += string(meta.Hash()[:])
	}
	return record
}

func (txToken TxTokenBase) JSONString() string {
	data, err := json.MarshalIndent(txToken, "", "\t")
	if err != nil {
		return ""
	}
	return string(data)
}

// =================== FUNCTIONS THAT GET STUFF ===================

// Hash returns the hash of all fields of the transaction
func (txToken *TxTokenBase) Hash() *common.Hash {
	if txToken.cachedHash != nil {
		return txToken.cachedHash
	}
	// final hash
	hash := common.HashH([]byte(txToken.String()))
	//txToken.cachedHash = &hash
	return &hash
}

func (txToken *TxTokenBase) HashWithoutMetadataSig() *common.Hash {
	// hashing to sign metadata is version-specific
	return nil
}

// Get SigPubKey of ptoken
func (txToken TxTokenBase) GetSigPubKey() []byte {
	return txToken.TxTokenData.TxNormal.GetSigPubKey()
}

// GetTxFeeToken - return Token Fee use to pay for privacy token Tx
func (txToken TxTokenBase) GetTxFeeToken() uint64 {
	return txToken.TxTokenData.TxNormal.GetTxFee()
}

func (txToken TxTokenBase) GetTokenID() *common.Hash {
	return &txToken.TxTokenData.PropertyID
}

func (txToken TxTokenBase) GetTransferData() (bool, []byte, uint64, *common.Hash) {
	pubkeys, amounts := txToken.TxTokenData.TxNormal.GetReceivers()
	if len(pubkeys) == 0 {
		return false, nil, 0, &txToken.TxTokenData.PropertyID
	}
	if len(pubkeys) > 1 {
		return false, nil, 0, &txToken.TxTokenData.PropertyID
	}
	return true, pubkeys[0], amounts[0], &txToken.TxTokenData.PropertyID
}

// CalculateBurnAmount - get tx value for pToken
func (txToken TxTokenBase) CalculateTxValue() uint64 {
	proof := txToken.TxTokenData.TxNormal.GetProof()
	if proof == nil {
		return 0
	}
	if proof.GetOutputCoins() == nil || len(proof.GetOutputCoins()) == 0 {
		return 0
	}
	if proof.GetInputCoins() == nil || len(proof.GetInputCoins()) == 0 { // coinbase tx
		txValue := uint64(0)
		for _, outCoin := range proof.GetOutputCoins() {
			txValue += outCoin.GetValue()
		}
		return txValue
	}

	if txToken.TxTokenData.TxNormal.IsPrivacy() {
		return 0
	}

	senderPKBytes := proof.GetInputCoins()[0].GetPublicKey().ToBytesS()
	txValue := uint64(0)
	for _, outCoin := range proof.GetOutputCoins() {
		outPKBytes := outCoin.GetPublicKey().ToBytesS()
		if bytes.Equal(senderPKBytes, outPKBytes) {
			continue
		}
		txValue += outCoin.GetValue()
	}
	return txValue
}

func (txToken TxTokenBase) ListSerialNumbersHashH() []common.Hash {
	tx := txToken.Tx
	result := []common.Hash{}
	if tx.GetProof() != nil {
		for _, d := range tx.GetProof().GetInputCoins() {
			hash := common.HashH(d.GetKeyImage().ToBytesS())
			result = append(result, hash)
		}
	}
	customTokenPrivacy := txToken.TxTokenData
	if customTokenPrivacy.TxNormal.GetProof() != nil {
		for _, d := range customTokenPrivacy.TxNormal.GetProof().GetInputCoins() {
			hash := common.HashH(d.GetKeyImage().ToBytesS())
			result = append(result, hash)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})
	return result
}

// GetTxFee - return fee PRV of Tx which contain privacy token Tx
func (txToken TxTokenBase) GetTxFee() uint64 {
	return txToken.Tx.GetTxFee()
}

// ================== NORMAL INIT FUNCTIONS ===================

// =================== FUNCTION THAT CHECK STUFFS  ===================

// ========== VALIDATE FUNCTIONS ===========

func (txToken TxTokenBase) ValidateType() bool {
	return txToken.Tx.GetType() == common.TxCustomTokenPrivacyType
}
