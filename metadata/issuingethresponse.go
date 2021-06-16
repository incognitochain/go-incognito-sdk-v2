package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type IssuingETHResponse struct {
	MetadataBase
	RequestedTxID   common.Hash
	UniqETHTx       []byte
	ExternalTokenID []byte
	SharedRandom       []byte `json:"SharedRandom,omitempty"`
}

func (iRes IssuingETHResponse) Hash() *common.Hash {
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

func (iRes *IssuingETHResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}