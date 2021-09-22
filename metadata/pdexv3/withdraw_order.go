package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// WithdrawOrderRequest
type WithdrawOrderRequest struct {
	PoolPairID string                           `json:"PoolPairID"`
	OrderID    string                           `json:"OrderID"`
	Amount     uint64                           `json:"Amount"`
	Receiver   map[common.Hash]coin.OTAReceiver `json:"Receiver"`
	NftID      common.Hash                      `json:"NftID"`
	metadataCommon.MetadataBase
}

func (req WithdrawOrderRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(req)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (req *WithdrawOrderRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(req)
}
