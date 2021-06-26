package metadata

import "github.com/incognitochain/go-incognito-sdk-v2/common"

// InitTokenResponse is the response for a InitTokenRequest.
type InitTokenResponse struct {
	MetadataBase
	RequestedTxID common.Hash
}

// Hash overrides MetadataBase.Hash().
func (iRes InitTokenResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += iRes.MetadataBase.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (iRes *InitTokenResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
