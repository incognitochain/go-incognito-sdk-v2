package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
)

// Interface for all types of metadata in tx
type Metadata interface {
	GetType() int
	Sign(*key.PrivateKey, Transaction) error
	Hash() *common.Hash
	HashWithoutSig() *common.Hash
	CalculateSize() uint64
}

// Interface for all type of transaction
type Transaction interface {
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
	GetMetadata() Metadata
	SetMetadata(Metadata)

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
