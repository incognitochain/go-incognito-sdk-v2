package tx_generic

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"strconv"
)

// Interface for Transaction Transfer Token
type TransactionToken interface {
	// GET/SET FUNCTION
	GetVersion() int8
	SetVersion(int8)
	GetMetadataType() int
	GetType() string
	SetType(string)
	GetLockTime() int64
	SetLockTime(int64)
	GetSenderAddrLastByte() byte
	SetGetSenderAddrLastByte(byte)
	GetTxFee() uint64
	SetTxFee(uint64)
	GetTxFeeToken() uint64
	GetInfo() []byte
	SetInfo([]byte)
	GetSigPubKey() []byte
	SetSigPubKey([]byte)
	GetSig() []byte
	SetSig([]byte)
	GetProof() privacy.Proof
	SetProof(privacy.Proof)
	GetTokenID() *common.Hash
	GetMetadata() metadata.Metadata
	SetMetadata(metadata.Metadata)

	GetTxTokenData() TxTokenData
	SetTxTokenData(TxTokenData) error
	GetTxBase() metadata.Transaction
	SetTxBase(metadata.Transaction) error
	GetTxNormal() metadata.Transaction
	SetTxNormal(metadata.Transaction) error

	// =================== FUNCTIONS THAT GET STUFF AND REQUIRE SOME CODING ===================
	GetTxActualSize() uint64
	GetReceivers() ([][]byte, []uint64)
	GetTransferData() (bool, []byte, uint64, *common.Hash)
	GetReceiverData() ([]coin.Coin, error)
	GetTxMintData() (bool, coin.Coin, *common.Hash, error)
	GetTxBurnData() (bool, coin.Coin, *common.Hash, error)
	GetTxFullBurnData() (bool, coin.Coin, coin.Coin, *common.Hash, error)
	ListOTAHashH() []common.Hash
	ListSerialNumbersHashH() []common.Hash
	String() string
	Hash() *common.Hash
	HashWithoutMetadataSig() *common.Hash
	CalculateTxValue() uint64

	// =================== FUNCTION THAT CHECK STUFFS  ===================
	CheckTxVersion(int8) bool
	IsSalaryTx() bool
	IsPrivacy() bool

	// Init Transaction, the input should be params such as: TxPrivacyInitParams
	Init(interface{}) error
}

type TxTokenData struct {
	// TxNormal is the normal transaction, it will never be token transaction
	TxNormal       metadata.Transaction
	PropertyID     common.Hash // = hash of TxCustomTokenprivacy data
	PropertyName   string
	PropertySymbol string

	Type     int    // action type
	Mintable bool   // default false
	Amount   uint64 // init amount
}

func (txData TxTokenData) GetPropertyID() common.Hash { return txData.PropertyID }
func (txData TxTokenData) GetPropertyName() string    { return txData.PropertyName }
func (txData TxTokenData) GetPropertySymbol() string  { return txData.PropertySymbol }
func (txData TxTokenData) GetType() int               { return txData.Type }
func (txData TxTokenData) IsMintable() bool           { return txData.Mintable }
func (txData TxTokenData) GetAmount() uint64          { return txData.Amount }

func (txData *TxTokenData) SetPropertyID(propID common.Hash)    { txData.PropertyID = propID }
func (txData *TxTokenData) SetPropertyName(propertyName string) { txData.PropertyName = propertyName }
func (txData *TxTokenData) SetPropertySymbol(propertySymbol string) {
	txData.PropertySymbol = propertySymbol
}
func (txData *TxTokenData) SetType(t int)             { txData.Type = t }
func (txData *TxTokenData) SetMintable(mintable bool) { txData.Mintable = mintable }
func (txData *TxTokenData) SetAmount(amount uint64)   { txData.Amount = amount }

func (txData TxTokenData) String() string {
	record := txData.PropertyName
	record += txData.PropertySymbol
	record += fmt.Sprintf("%d", txData.Amount)
	if txData.TxNormal.GetProof() != nil {
		inputCoins := txData.TxNormal.GetProof().GetInputCoins()
		outputCoins := txData.TxNormal.GetProof().GetOutputCoins()
		for _, out := range outputCoins {
			publicKeyBytes := []byte{}
			if out.GetPublicKey() != nil {
				publicKeyBytes = out.GetPublicKey().ToBytesS()
			}
			record += string(publicKeyBytes)
			record += strconv.FormatUint(out.GetValue(), 10)
		}
		for _, in := range inputCoins {
			publicKeyBytes := []byte{}
			if in.GetPublicKey() != nil {
				publicKeyBytes = in.GetPublicKey().ToBytesS()
			}
			record += string(publicKeyBytes)
			if in.GetValue() > 0 {
				record += strconv.FormatUint(in.GetValue(), 10)
			}
		}
	}
	return record
}

func (txData TxTokenData) Hash() (*common.Hash, error) {
	point := crypto.HashToPoint([]byte(txData.String()))
	hash := new(common.Hash)
	err := hash.SetBytes(point.ToBytesS())
	if err != nil {
		return nil, err
	}
	return hash, nil
}
