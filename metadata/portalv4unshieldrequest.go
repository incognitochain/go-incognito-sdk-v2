package metadata

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type PortalUnshieldRequest struct {
	MetadataBase
	OTAPubKeyStr   string // OTA
	TxRandomStr    string
	RemoteAddress  string
	TokenID        string
	UnshieldAmount uint64
}

func NewPortalUnshieldRequest(metaType int, otaPubKeyStr, txRandomStr string, tokenID, remoteAddress string, burnAmount uint64) (*PortalUnshieldRequest, error) {
	portalUnshieldReq := &PortalUnshieldRequest{
		OTAPubKeyStr:   otaPubKeyStr,
		TxRandomStr:    txRandomStr,
		UnshieldAmount: burnAmount,
		RemoteAddress:  remoteAddress,
		TokenID:        tokenID,
	}

	portalUnshieldReq.MetadataBase = MetadataBase{
		Type: metaType,
	}

	return portalUnshieldReq, nil
}

func (uReq PortalUnshieldRequest) Hash() *common.Hash {
	record := uReq.MetadataBase.Hash().String()
	record += uReq.OTAPubKeyStr
	record += uReq.TxRandomStr
	record += uReq.RemoteAddress
	record += strconv.FormatUint(uReq.UnshieldAmount, 10)
	record += uReq.TokenID

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (uReq *PortalUnshieldRequest) CalculateSize() uint64 {
	return calculateSize(uReq)
}
