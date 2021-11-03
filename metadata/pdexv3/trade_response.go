package pdexv3

import (
	"encoding/json"
	"math/big"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// TradeStatus containns the info tracked by feature statedb, which is then displayed in RPC status queries.
// For refunded trade, all fields except Status are ignored
type TradeStatus struct {
	Status     int         `json:"Status"`
	BuyAmount  uint64      `json:"BuyAmount"`
	TokenToBuy common.Hash `json:"TokenToBuy"`
}

// TradeResponse is the metadata inside response tx for trade
type TradeResponse struct {
	Status      int         `json:"Status"`
	RequestTxID common.Hash `json:"RequestTxID"`
	metadataCommon.MetadataBase
}

// AcceptedTrade is added as Content for produced beacon Instructions after handling a trade successfully
type AcceptedTrade struct {
	Receiver     coin.OTAReceiver         `json:"Receiver"`
	Amount       uint64                   `json:"Amount"`
	TradePath    []string                 `json:"TradePath"`
	TokenToBuy   common.Hash              `json:"TokenToBuy"`
	PairChanges  [][2]*big.Int            `json:"PairChanges"`
	RewardEarned []map[common.Hash]uint64 `json:"RewardEarned"`
	OrderChanges []map[string][2]*big.Int `json:"OrderChanges"`
}

func (md AcceptedTrade) GetType() int {
	return metadataCommon.Pdexv3TradeRequestMeta
}

func (md AcceptedTrade) GetStatus() int {
	return TradeAcceptedStatus
}

// RefundedTrade is added as Content for produced beacon instruction after failure to handle a trade
type RefundedTrade struct {
	Receiver coin.OTAReceiver `json:"Receiver"`
	TokenID  common.Hash      `json:"TokenToSell"`
	Amount   uint64           `json:"Amount"`
}

func (md RefundedTrade) GetType() int {
	return metadataCommon.Pdexv3TradeRequestMeta
}

func (md RefundedTrade) GetStatus() int {
	return TradeRefundedStatus
}

func (res TradeResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(res)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (res *TradeResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(res)
}
