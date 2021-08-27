package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

type PortalUnshieldResponse struct {
	MetadataBase
	RequestStatus  string
	ReqTxID        common.Hash
	OTAPubKeyStr   string
	TxRandomStr    string
	UnshieldAmount uint64
	IncTokenID     string
}

func NewPortalV4UnshieldResponse(
	requestStatus string,
	reqTxID common.Hash,
	requesterAddressStr string,
	txRandomStr string,
	amount uint64,
	tokenID string,
	metaType int,
) *PortalUnshieldResponse {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	return &PortalUnshieldResponse{
		RequestStatus:  requestStatus,
		ReqTxID:        reqTxID,
		MetadataBase:   metadataBase,
		OTAPubKeyStr:   requesterAddressStr,
		TxRandomStr: txRandomStr,
		UnshieldAmount: amount,
		IncTokenID:     tokenID,
	}
}

func (iRes PortalUnshieldResponse) Hash() *common.Hash {
	record := iRes.MetadataBase.Hash().String()
	record += iRes.RequestStatus
	record += iRes.ReqTxID.String()
	record += iRes.OTAPubKeyStr
	record += iRes.TxRandomStr
	record += strconv.FormatUint(iRes.UnshieldAmount, 10)
	record += iRes.IncTokenID
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *PortalUnshieldResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}
