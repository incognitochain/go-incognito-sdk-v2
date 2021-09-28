package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// PDECrossPoolTradeResponse is the response for a PDECrossPoolTradeRequest.
type PDECrossPoolTradeResponse struct {
	MetadataBase
	TradeStatus   string
	RequestedTxID common.Hash
}

// Hash overrides MetadataBase.Hash().
func (iRes PDECrossPoolTradeResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += iRes.TradeStatus
	record += iRes.MetadataBase.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (iRes *PDECrossPoolTradeResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
