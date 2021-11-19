package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// AddOrderRequest represents a request for adding an order book to a pool.
type AddOrderRequest struct {
	// TokenToSell is the ID of the selling token.
	TokenToSell         common.Hash                      `json:"TokenToSell"`

	// PoolPairID is the ID of the pool pair where the order belongs to. In Incognito, an order book is subject to a specific pool.
	PoolPairID          string                           `json:"PoolPairID"`

	// SellAmount is the amount of the `TokenToSell` the user wished to sell.
	SellAmount          uint64                           `json:"SellAmount"`

	// MinAcceptableAmount is the minimum amount of the buying token the user wished to receive.
	MinAcceptableAmount uint64                           `json:"MinAcceptableAmount"`

	// Receiver is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
	Receiver            map[common.Hash]coin.OTAReceiver `json:"Receiver"`

	// is the ID of the NFT associated with order.
	NftID               common.Hash                      `json:"NftID"`

	metadataCommon.MetadataBase
}

func NewAddOrderRequest(
	tokenToSell common.Hash,
	pairID string,
	sellAmount uint64,
	minAcceptableAmount uint64,
	recv map[common.Hash]coin.OTAReceiver,
	nftID common.Hash,
	metaType int,
) (*AddOrderRequest, error) {
	r := &AddOrderRequest{
		TokenToSell:         tokenToSell,
		PoolPairID:          pairID,
		SellAmount:          sellAmount,
		MinAcceptableAmount: minAcceptableAmount,
		Receiver:            recv,
		NftID:               nftID,
		MetadataBase: metadataCommon.MetadataBase{
			Type: metaType,
		},
	}
	return r, nil
}

func (req AddOrderRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(req)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (req *AddOrderRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(req)
}
