package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// WithdrawOrderStatus containns the info tracked by feature statedb, which is then displayed in RPC status queries.
// For refunded `add order` requests, all fields except Status are ignored
type WithdrawOrderStatus struct {
	Status         int         `json:"Status"`
	TokenID        common.Hash `json:"TokenID"`
	WithdrawAmount uint64      `json:"Amount"`
}

type WithdrawOrderResponse struct {
	Status      int         `json:"Status"`
	RequestTxID common.Hash `json:"RequestTxID"`
	metadataCommon.MetadataBase
}

type AcceptedWithdrawOrder struct {
	PoolPairID string           `json:"PoolPairID"`
	OrderID    string           `json:"OrderID"`
	TokenID    common.Hash      `json:"TokenID"`
	Receiver   coin.OTAReceiver `json:"Receiver"`
	Amount     uint64           `json:"Amount"`
}

func (md AcceptedWithdrawOrder) GetType() int {
	return metadataCommon.Pdexv3WithdrawOrderRequestMeta
}

func (md AcceptedWithdrawOrder) GetStatus() int {
	return WithdrawOrderAcceptedStatus
}

type RejectedWithdrawOrder struct {
	PoolPairID string `json:"PoolPairID"`
	OrderID    string `json:"OrderID"`
}

func (md RejectedWithdrawOrder) GetType() int {
	return metadataCommon.Pdexv3WithdrawOrderRequestMeta
}

func (md RejectedWithdrawOrder) GetStatus() int {
	return WithdrawOrderRejectedStatus
}

func (res WithdrawOrderResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(res)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (res *WithdrawOrderResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(res)
}
