package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type PDEFeeWithdrawalResponse struct {
	MetadataBase
	RequestedTxID common.Hash
	SharedRandom  []byte `json:"SharedRandom,omitempty"`
}

func (iRes PDEFeeWithdrawalResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += iRes.MetadataBase.Hash().String()
	if iRes.SharedRandom != nil && len(iRes.SharedRandom) > 0 {
		record += string(iRes.SharedRandom)
	}
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *PDEFeeWithdrawalResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}

func (iRes *PDEFeeWithdrawalResponse) SetSharedRandom(r []byte) {
	iRes.SharedRandom = r
}
