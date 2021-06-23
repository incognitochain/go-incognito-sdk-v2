package metadata

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// RelayingHeader - relaying header chain
// metadata - create normal tx with this metadata
type RelayingHeader struct {
	MetadataBase
	IncogAddressStr string
	Header          string
	BlockHeight     uint64
}

func NewRelayingHeader(
	metaType int,
	incognitoAddrStr string,
	header string,
	blockHeight uint64,
) (*RelayingHeader, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	relayingHeader := &RelayingHeader{
		IncogAddressStr: incognitoAddrStr,
		Header:          header,
		BlockHeight:     blockHeight,
	}
	relayingHeader.MetadataBase = metadataBase
	return relayingHeader, nil
}

func (rh RelayingHeader) Hash() *common.Hash {
	record := rh.MetadataBase.Hash().String()
	record += rh.IncogAddressStr
	record += rh.Header
	record += strconv.Itoa(int(rh.BlockHeight))

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (rh *RelayingHeader) CalculateSize() uint64 {
	return calculateSize(rh)
}
