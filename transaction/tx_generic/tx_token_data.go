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

// TransactionToken describes all functions of a token transaction.
type TransactionToken interface {
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

	CheckTxVersion(int8) bool
	IsSalaryTx() bool
	IsPrivacy() bool

	Init(interface{}) error
}

// TxTokenData represents all data of a token transaction.
type TxTokenData struct {
	TxNormal       metadata.Transaction // TxNormal is the PRV transaction, it will never be token transaction
	PropertyID     common.Hash          // = hash of TxCustomTokenPrivacy data
	PropertyName   string
	PropertySymbol string

	Type     int    // action type
	Mintable bool   // default false
	Amount   uint64 // init amount
}

// GetPropertyID returns the tokenID of a TxTokenData.
func (txData TxTokenData) GetPropertyID() common.Hash { return txData.PropertyID }

// GetPropertyName returns the token name of a TxTokenData.
func (txData TxTokenData) GetPropertyName() string { return txData.PropertyName }

// GetPropertySymbol returns the token symbol of a TxTokenData.
func (txData TxTokenData) GetPropertySymbol() string { return txData.PropertySymbol }

// GetType returns the type of a TxTokenData.
func (txData TxTokenData) GetType() int { return txData.Type }

// IsMintable checks if a TxTokenData is mintable.
func (txData TxTokenData) IsMintable() bool { return txData.Mintable }

// GetAmount returns the amount of a TxTokenData.
func (txData TxTokenData) GetAmount() uint64 { return txData.Amount }

// SetPropertyID sets v as the tokenID of a TxTokenData.
func (txData *TxTokenData) SetPropertyID(v common.Hash) { txData.PropertyID = v }

// SetPropertyName sets v as the token name of a TxTokenData.
func (txData *TxTokenData) SetPropertyName(v string) { txData.PropertyName = v }

// SetPropertySymbol sets v as the token symbol of a TxTokenData.
func (txData *TxTokenData) SetPropertySymbol(v string) {
	txData.PropertySymbol = v
}

// SetType sets v as the type of a TxTokenData.
func (txData *TxTokenData) SetType(v int) { txData.Type = v }

// SetMintable sets v as the mintable flag of a TxTokenData.
func (txData *TxTokenData) SetMintable(v bool) { txData.Mintable = v }

// SetAmount sets v as the amount of a TxTokenData.
func (txData *TxTokenData) SetAmount(v uint64) { txData.Amount = v }

// String returns the string-representation of a TxTokenData.
func (txData TxTokenData) String() string {
	record := txData.PropertyName
	record += txData.PropertySymbol
	record += fmt.Sprintf("%d", txData.Amount)
	if txData.TxNormal.GetProof() != nil {
		inputCoins := txData.TxNormal.GetProof().GetInputCoins()
		outputCoins := txData.TxNormal.GetProof().GetOutputCoins()
		for _, out := range outputCoins {
			publicKeyBytes := make([]byte, 0)
			if out.GetPublicKey() != nil {
				publicKeyBytes = out.GetPublicKey().ToBytesS()
			}
			record += string(publicKeyBytes)
			record += strconv.FormatUint(out.GetValue(), 10)
		}
		for _, in := range inputCoins {
			publicKeyBytes := make([]byte, 0)
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

// Hash calculates the hash of a TxTokenData.
func (txData TxTokenData) Hash() (*common.Hash, error) {
	point := crypto.HashToPoint([]byte(txData.String()))
	hash := new(common.Hash)
	err := hash.SetBytes(point.ToBytesS())
	if err != nil {
		return nil, err
	}
	return hash, nil
}
