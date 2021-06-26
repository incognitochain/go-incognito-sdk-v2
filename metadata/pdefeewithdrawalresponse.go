package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// PDEFeeWithdrawalResponse is the response for a PDEFeeWithdrawalRequest.
type PDEFeeWithdrawalResponse struct {
	MetadataBase
	RequestedTxID common.Hash
	SharedRandom  []byte `json:"SharedRandom,omitempty"`
}

// Hash overrides MetadataBase.Hash().
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

// CalculateSize overrides MetadataBase.CalculateSize().
func (iRes *PDEFeeWithdrawalResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}

// SetSharedRandom sets v as the shared random of a PDEFeeWithdrawalResponse.
func (iRes *PDEFeeWithdrawalResponse) SetSharedRandom(v []byte) {
	iRes.SharedRandom = v
}
