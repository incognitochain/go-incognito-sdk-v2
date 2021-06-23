package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// IssuingEVMResponse is the response for a IssuingEVMRequest.
type IssuingEVMResponse struct {
	MetadataBase
	RequestedTxID   common.Hash
	UniqETHTx       []byte
	ExternalTokenID []byte
	SharedRandom    []byte `json:"SharedRandom,omitempty"`
}

// Hash overrides MetadataBase.Hash().
func (iRes IssuingEVMResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += string(iRes.UniqETHTx)
	record += string(iRes.ExternalTokenID)
	record += iRes.MetadataBase.Hash().String()
	if iRes.SharedRandom != nil && len(iRes.SharedRandom) > 0 {
		record += string(iRes.SharedRandom)
	}
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (iRes *IssuingEVMResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
