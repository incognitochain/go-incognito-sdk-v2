package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type PDECrossPoolTradeResponse struct {
	MetadataBase
	TradeStatus   string
	RequestedTxID common.Hash
}

func (iRes PDECrossPoolTradeResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += iRes.TradeStatus
	record += iRes.MetadataBase.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *PDECrossPoolTradeResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
