package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// AddOrderRequest
type AddOrderRequest struct {
	TokenToSell         common.Hash                      `json:"TokenToSell"`
	PoolPairID          string                           `json:"PoolPairID"`
	SellAmount          uint64                           `json:"SellAmount"`
	MinAcceptableAmount uint64                           `json:"MinAcceptableAmount"`
	Receiver            map[common.Hash]coin.OTAReceiver `json:"Receiver"`
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
