package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// PDEWithdrawalResponse is the response for a PDEWithdrawalRequest.
type PDEWithdrawalResponse struct {
	MetadataBase
	RequestedTxID common.Hash
	TokenIDStr    string
	SharedRandom  []byte `json:"SharedRandom,omitempty"`
}

// Hash overrides MetadataBase.Hash().
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

// CalculateSize overrides MetadataBase.CalculateSize().
func (iRes *PDEWithdrawalResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
