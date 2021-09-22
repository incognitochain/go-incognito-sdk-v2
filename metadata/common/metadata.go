package common

import (
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
)

// Metadata is an interface describing all common methods of a metadata.
// A Metadata is a special piece of information enclosed with a transaction to indicate additional purpose of
// the transaction.
type Metadata interface {
	// GetType returns the type of a Metadata
	GetType() int

	// Sign signs the metadata with the provided private key.
	Sign(*key.PrivateKey, Transaction) error

	// Hash calculates the hash of a metadata.
	Hash() *common.Hash

	// HashWithoutSig calculates the hash of a metadata without including its sig.
	HashWithoutSig() *common.Hash

	// CalculateSize returns the size of a metadata in bytes.
	CalculateSize() uint64
}

// Transaction is an interface describing all common methods of a transaction.
type Transaction interface {
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
