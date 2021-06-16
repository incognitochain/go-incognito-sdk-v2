package metadata

import "github.com/incognitochain/go-incognito-sdk-v2/common"

type InitTokenResponse struct {
	MetadataBase
	RequestedTxID common.Hash
}

func (iRes InitTokenResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += iRes.MetadataBase.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *InitTokenResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
