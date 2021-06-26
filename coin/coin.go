package coin

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// Coin represents a input or output coin of a transaction.
type Coin interface {
	GetVersion() uint8
	GetCommitment() *crypto.Point
	GetInfo() []byte
	GetPublicKey() *crypto.Point
	GetKeyImage() *crypto.Point
	GetValue() uint64
	GetRandomness() *crypto.Scalar
	GetShardID() (uint8, error)
	GetSNDerivator() *crypto.Scalar
	GetCoinDetailEncrypted() []byte
	IsEncrypted() bool
	GetTxRandom() *TxRandom
	GetSharedRandom() *crypto.Scalar
	GetSharedConcealRandom() *crypto.Scalar
	GetAssetTag() *crypto.Point
	SetValue(v uint64)
	Decrypt(*key.KeySet) (PlainCoin, error)

	Bytes() []byte
	SetBytes([]byte) error

	CheckCoinValid(key.PaymentAddress, []byte, uint64) bool
	DoesCoinBelongToKeySet(keySet *key.KeySet) (bool, *crypto.Point)
}

// PlainCoin represents an un-encrypted coin of a transaction.
type PlainCoin interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error

	GetVersion() uint8
	GetCommitment() *crypto.Point
	GetInfo() []byte
	GetPublicKey() *crypto.Point
	GetValue() uint64
	GetKeyImage() *crypto.Point
	GetRandomness() *crypto.Scalar
	GetShardID() (uint8, error)
	GetSNDerivator() *crypto.Scalar
	GetCoinDetailEncrypted() []byte
	IsEncrypted() bool
	GetTxRandom() *TxRandom
	GetSharedRandom() *crypto.Scalar
	GetSharedConcealRandom() *crypto.Scalar
	GetAssetTag() *crypto.Point

	SetKeyImage(*crypto.Point)
	SetPublicKey(*crypto.Point)
	SetCommitment(*crypto.Point)
	SetInfo([]byte)
	SetValue(uint64)
	SetRandomness(*crypto.Scalar)

	ParseKeyImageWithPrivateKey(key.PrivateKey) (*crypto.Point, error)
	ParsePrivateKeyOfCoin(key.PrivateKey) (*crypto.Scalar, error)

	ConcealOutputCoin(additionalData interface{}) error

	Bytes() []byte
	SetBytes([]byte) error
}

// NewPlainCoinFromByte parse a new PlainCoin from its bytes.
//
// The first byte should determine the coin version.
func NewPlainCoinFromByte(b []byte) (PlainCoin, error) {
	version := byte(CoinVersion2)
	if len(b) >= 1 {
		version = b[0]
	}
	var c PlainCoin
	if version == CoinVersion2 {
		c = new(CoinV2)
	} else {
		c = new(PlainCoinV1)
	}
	err := c.SetBytes(b)
	return c, err
}

// NewCoinFromByte creates a new Coin from its bytes.
//
// The first byte should determine the coin version or json marshal "34".
func NewCoinFromByte(b []byte) (Coin, error) {
	coinV1 := new(CoinV1)
	coinV2 := new(CoinV2)
	if errV2 := json.Unmarshal(b, coinV2); errV2 != nil {
		if errV1 := json.Unmarshal(b, coinV1); errV1 != nil {
			version := b[0]
			if version == CoinVersion2 {
				err := coinV2.SetBytes(b)
				return coinV2, err
			} else {
				err := coinV1.SetBytes(b)
				return coinV1, err
			}
		} else {
			return coinV1, nil
		}
	} else {
		return coinV2, nil
	}
}

// ParseCoinsFromBytes parses a list of raw bytes into a list of corresponding Coin objects.
func ParseCoinsFromBytes(data []json.RawMessage) ([]Coin, error) {
	coinList := make([]Coin, len(data))
	for i := 0; i < len(data); i++ {
		if coin, err := NewCoinFromByte(data[i]); err != nil {
			return nil, err
		} else {
			coinList[i] = coin
		}
	}
	return coinList, nil
}
