package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// PDETradeResponse is the response for a PDETradeRequest.
type PDETradeResponse struct {
	MetadataBase
	TradeStatus   string
	RequestedTxID common.Hash
}

// Hash overrides MetadataBase.Hash().
func (iRes PDETradeResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += iRes.TradeStatus
	record += iRes.MetadataBase.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (iRes *PDETradeResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
