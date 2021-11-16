package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// WithdrawOrderRequest
type WithdrawOrderRequest struct {
	// PoolPairID is the ID of the target pool from which the user wants to withdraw his order.
	PoolPairID string                           `json:"PoolPairID"`

	// OrderID is the ID of the added order.
	OrderID    string                           `json:"OrderID"`

	// Amount is the amount in which we want to withdraw (0 for all).
	Amount     uint64                           `json:"Amount"`

	// Receiver is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
	Receiver   map[common.Hash]coin.OTAReceiver `json:"Receiver"`

	// NftID is the ID of the NFT associated with the order.
	NftID      common.Hash                      `json:"NftID"`

	metadataCommon.MetadataBase
}

func NewWithdrawOrderRequest(
	pairID, orderID string,
	amount uint64,
	recv map[common.Hash]coin.OTAReceiver,
	nftID common.Hash,
	metaType int,
) (*WithdrawOrderRequest, error) {
	r := &WithdrawOrderRequest{
		PoolPairID: pairID,
		OrderID:    orderID,
		Amount:     amount,
		Receiver:   recv,
		NftID:      nftID,
		MetadataBase: metadataCommon.MetadataBase{
			Type: metaType,
		},
	}
	return r, nil
}

func (req WithdrawOrderRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(req)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (req *WithdrawOrderRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(req)
}
