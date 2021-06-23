package tx_generic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"sort"
)

// Tx is a alias for metadata.Transaction
type Tx = metadata.Transaction

// TxTokenBase represents a token transaction field. It is used in both TxTokenVer1 and TxTokenVer2.
// A TxTokenBase consists of a PRV sub-transaction for paying the transaction fee, and the token data.
type TxTokenBase struct {
	Tx
	TxTokenData TxTokenData `json:"TxTokenPrivacyData"`
	cachedHash  *common.Hash
}

// TxTokenParams consists of parameters used to create a new token transaction.
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
	KvArgs          map[string]interface{}
}

// TokenParam represents the parameters of a token transaction.
type TokenParam struct {
	PropertyID     string             `json:"TokenID"`
	PropertyName   string             `json:"TokenName"`
	PropertySymbol string             `json:"TokenSymbol"`
	Amount         uint64             `json:"TokenAmount"`
	TokenTxType    int                `json:"TokenTxType"`
	Receiver       []*key.PaymentInfo `json:"TokenReceiver"`
	TokenInput     []coin.PlainCoin   `json:"TokenInput"`
	Mintable       bool               `json:"TokenMintable"`
	Fee            uint64             `json:"TokenFee"`
	KvArgs         map[string]interface{}
}

// NewTokenParam creates a new TokenParam based on the given inputs.
func NewTokenParam(propertyID, propertyName, propertySymbol string,
	amount uint64,
	tokenTxType int,
	receivers []*key.PaymentInfo,
	tokenInput []coin.PlainCoin,
	mintable bool,
	fee uint64,
	kvArgs map[string]interface{}) *TokenParam {

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
		KvArgs:         kvArgs,
	}
	return params
}

// NewTxTokenParams creates a new TxTokenParams based on the given inputs.
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
	kvArgs map[string]interface{}) *TxTokenParams {
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
		KvArgs:          kvArgs,
	}
	return params
}

// GetTxBase returns the PRV sub-transaction of a TxTokenBase.
func (txToken TxTokenBase) GetTxBase() metadata.Transaction { return txToken.Tx }

// GetTxNormal returns the token sub-transaction of a TxTokenBase.
func (txToken TxTokenBase) GetTxNormal() metadata.Transaction { return txToken.TxTokenData.TxNormal }

// GetTxTokenData returns token data of a TxTokenBase.
func (txToken TxTokenBase) GetTxTokenData() TxTokenData { return txToken.TxTokenData }

// GetTxMintData returns the minting data of a TxTokenBase.
func (txToken TxTokenBase) GetTxMintData() (bool, coin.Coin, *common.Hash, error) {
	tokenID := txToken.TxTokenData.GetPropertyID()
	return GetTxMintData(txToken.TxTokenData.TxNormal, &tokenID)
}

// GetTxBurnData returns the burning data (token only) of a TxTokenBase.
func (txToken TxTokenBase) GetTxBurnData() (bool, coin.Coin, *common.Hash, error) {
	tokenID := txToken.TxTokenData.GetPropertyID()
	isBurn, burnCoin, _, err := txToken.TxTokenData.TxNormal.GetTxBurnData()
	return isBurn, burnCoin, &tokenID, err
}

// GetTxFullBurnData returns the full burning data (both PRV and token) of a TxTokenBase.
func (txToken TxTokenBase) GetTxFullBurnData() (bool, coin.Coin, coin.Coin, *common.Hash, error) {
	isBurnToken, burnToken, burnedTokenID, errToken := txToken.GetTxBurnData()
	isBurnPrv, burnPrv, _, errPrv := txToken.GetTxBase().GetTxBurnData()

	if errToken != nil && errPrv != nil {
		return false, nil, nil, nil, fmt.Errorf("%v and %v", errPrv, errToken)
	}

	return isBurnPrv || isBurnToken, burnPrv, burnToken, burnedTokenID, nil
}

// GetSigPubKey returns the sigPubKey of a TxTokenBase.
func (txToken TxTokenBase) GetSigPubKey() []byte {
	return txToken.TxTokenData.TxNormal.GetSigPubKey()
}

// GetTxFeeToken returns to transaction fee paid in the token of a TxTokenBase.
func (txToken TxTokenBase) GetTxFeeToken() uint64 {
	return txToken.TxTokenData.TxNormal.GetTxFee()
}

// GetTokenID returns the tokenID of a TxTokenBase.
func (txToken TxTokenBase) GetTokenID() *common.Hash {
	return &txToken.TxTokenData.PropertyID
}

// GetTransferData returns the transferred data (receivers and amounts) of a TxTokenBase.
// The result does not the transferred data of PRV.
func (txToken TxTokenBase) GetTransferData() (bool, []byte, uint64, *common.Hash) {
	pubKeys, amounts := txToken.TxTokenData.TxNormal.GetReceivers()
	if len(pubKeys) == 0 {
		return false, nil, 0, &txToken.TxTokenData.PropertyID
	}
	if len(pubKeys) > 1 {
		return false, nil, 0, &txToken.TxTokenData.PropertyID
	}
	return true, pubKeys[0], amounts[0], &txToken.TxTokenData.PropertyID
}

// GetTxFee returns the transaction fee paid in PRV of a TxTokenBase.
func (txToken TxTokenBase) GetTxFee() uint64 {
	return txToken.Tx.GetTxFee()
}

// IsSalaryTx checks if a TxTokenBase is a salary transaction.
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

// SetTxBase sets v as the TxBase of a TxTokenBase.
func (txToken *TxTokenBase) SetTxBase(v metadata.Transaction) error {
	txToken.Tx = v
	return nil
}

// SetTxNormal sets v as the TxNormal of a TxTokenBase.
func (txToken *TxTokenBase) SetTxNormal(v metadata.Transaction) error {
	txToken.TxTokenData.TxNormal = v
	return nil
}

// SetTxTokenData sets v as the token data of a TxTokenBase.
func (txToken *TxTokenBase) SetTxTokenData(v TxTokenData) error {
	txToken.TxTokenData = v
	return nil
}

// CheckAuthorizedSender checks if the sender of a TxTokenBase is authorized w.r.t to a public key.
func (txToken TxTokenBase) CheckAuthorizedSender(publicKey []byte) (bool, error) {
	sigPubKey := txToken.TxTokenData.TxNormal.GetSigPubKey()
	if bytes.Equal(sigPubKey, publicKey) {
		return true, nil
	} else {
		return false, nil
	}
}

// MarshalJSON does the JSON-marshalling operation for a TxBase.
func (txToken TxTokenBase) MarshalJSON() ([]byte, error) {
	type TemporaryTxToken struct {
		TxBase
		TxTokenData TxTokenData `json:"TxTokenPrivacyData"`
	}
	tempTx := TemporaryTxToken{}
	tempTx.TxTokenData = txToken.GetTxTokenData()
	tx := txToken.GetTxBase()
	if tx == nil {
		return nil, fmt.Errorf("cannot unmarshal transaction: txfee cannot be nil")
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

// String returns the string-representation of a TxTokenBase.
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

// Hash calculates the hash of a TxTokenBase.
func (txToken *TxTokenBase) Hash() *common.Hash {
	if txToken.cachedHash != nil {
		return txToken.cachedHash
	}
	// final hash
	hash := common.HashH([]byte(txToken.String()))
	//txToken.cachedHash = &hash
	return &hash
}

// HashWithoutMetadataSig calculates the hash of a TxTokenBase with out adding the signature of its metadata.
func (txToken *TxTokenBase) HashWithoutMetadataSig() *common.Hash {
	return nil
}

// CalculateTxValue calculates total output values (not including the coins which are sent back to the sender).
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

// ListSerialNumbersHashH returns the hash list of all serial numbers in a TxTokenBase.
func (txToken TxTokenBase) ListSerialNumbersHashH() []common.Hash {
	tx := txToken.Tx
	result := make([]common.Hash, 0)
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

// ValidateType checks if the type of a TxTokenBase is valid.
func (txToken TxTokenBase) ValidateType() bool {
	return txToken.Tx.GetType() == common.TxCustomTokenPrivacyType
}
