package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// IssuingResponse is the response for a IssuingRequest.
type IssuingResponse struct {
	MetadataBase
	RequestedTxID common.Hash
	SharedRandom  []byte `json:"SharedRandom,omitempty"`
}

// Hash overrides MetadataBase.Hash().
func (iRes IssuingResponse) Hash() *common.Hash {
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
func (iRes *IssuingResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
