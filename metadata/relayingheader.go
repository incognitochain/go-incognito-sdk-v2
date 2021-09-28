package metadata

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// RelayingHeader is a relay request from other public blockchains.
type RelayingHeader struct {
	MetadataBase
	IncAddressStr string
	Header        string
	BlockHeight   uint64
}

// NewRelayingHeader creates a new RelayingHeader.
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
		IncAddressStr: incognitoAddrStr,
		Header:        header,
		BlockHeight:   blockHeight,
	}
	relayingHeader.MetadataBase = metadataBase
	return relayingHeader, nil
}

// Hash overrides MetadataBase.Hash().
func (rh RelayingHeader) Hash() *common.Hash {
	record := rh.MetadataBase.Hash().String()
	record += rh.IncAddressStr
	record += rh.Header
	record += strconv.Itoa(int(rh.BlockHeight))

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (rh *RelayingHeader) CalculateSize() uint64 {
	return calculateSize(rh)
}
