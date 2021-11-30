package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// AddOrderStatus containns the info tracked by feature statedb, which is then displayed in RPC status queries.
// For refunded `add order` requests, all fields except Status are ignored
type AddOrderStatus struct {
	Status  int    `json:"Status"`
	OrderID string `json:"OrderID"`
}

// AddOrderResponse is the metadata inside response tx for `add order` (applicable for refunded case only)
type AddOrderResponse struct {
	Status      int         `json:"Status"`
	RequestTxID common.Hash `json:"RequestTxID"`
	metadataCommon.MetadataBase
}

// AcceptedAddOrder is added as Content for produced beacon instruction after to handling an order successfully
type AcceptedAddOrder struct {
	PoolPairID     string      `json:"PoolPairID"`
	OrderID        string      `json:"OrderID"`
	NftID          common.Hash `json:"NftID"`
	Token0Rate     uint64      `json:"Token0Rate"`
	Token1Rate     uint64      `json:"Token1Rate"`
	Token0Balance  uint64      `json:"Token0Balance"`
	Token1Balance  uint64      `json:"Token1Balance"`
	TradeDirection byte        `json:"TradeDirection"`
}

func (md AcceptedAddOrder) GetType() int {
	return metadataCommon.Pdexv3AddOrderRequestMeta
}

func (md AcceptedAddOrder) GetStatus() int {
	return OrderAcceptedStatus
}

// RefundedAddOrder is added as Content for produced beacon instruction after failure to handle an order
type RefundedAddOrder struct {
	Receiver coin.OTAReceiver `json:"Receiver"`
	TokenID  common.Hash      `json:"TokenID"`
	Amount   uint64           `json:"Amount"`
}

func (md RefundedAddOrder) GetType() int {
	return metadataCommon.Pdexv3AddOrderRequestMeta
}

func (md RefundedAddOrder) GetStatus() int {
	return OrderRefundedStatus
}

func (res AddOrderResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(res)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (res *AddOrderResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(res)
}
