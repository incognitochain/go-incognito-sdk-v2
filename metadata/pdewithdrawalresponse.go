package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type PDEWithdrawalResponse struct {
	MetadataBase
	RequestedTxID common.Hash
	TokenIDStr    string
	SharedRandom       []byte `json:"SharedRandom,omitempty"`
}

func (iRes PDEWithdrawalResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += iRes.TokenIDStr
	record += iRes.MetadataBase.Hash().String()
	if iRes.SharedRandom != nil && len(iRes.SharedRandom) > 0 {
		record += string(iRes.SharedRandom)
	}
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *PDEWithdrawalResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}