package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type PDETradeResponse struct {
	MetadataBase
	TradeStatus   string
	RequestedTxID common.Hash
}

func (iRes PDETradeResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += iRes.TradeStatus
	record += iRes.MetadataBase.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *PDETradeResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
