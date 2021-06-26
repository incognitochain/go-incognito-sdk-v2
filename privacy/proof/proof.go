package proof

import (
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/proof/range_proof"
)

// Paymentproof
type Proof interface {
	GetVersion() uint8
	Init()
	GetInputCoins() []coin.PlainCoin
	GetOutputCoins() []coin.Coin
	GetRangeProof() range_proof.RangeProof

	SetInputCoins([]coin.PlainCoin) error
	SetOutputCoins([]coin.Coin) error

	Bytes() []byte
	SetBytes(proofBytes []byte) error

	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error

	IsPrivacy() bool
}
