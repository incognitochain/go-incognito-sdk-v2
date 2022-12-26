package bridge

import (
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type IssuingReshieldResponse struct {
	metadataCommon.MetadataBase
	RequestedTxID   common.Hash `json:"RequestedTxID"`
	UniqTx          []byte      `json:"UniqETHTx"`
	ExternalTokenID []byte      `json:"ExternalTokenID"`
}

type AcceptedReshieldRequest struct {
	UnifiedTokenID *common.Hash              `json:"UnifiedTokenID"`
	Receiver       coin.OTAReceiver          `json:"Receiver"`
	TxReqID        common.Hash               `json:"TxReqID"`
	ReshieldData   AcceptedShieldRequestData `json:"ReshieldData"`
}

func NewIssuingReshieldResponse(
	requestedTxID common.Hash,
	uniqTx []byte,
	externalTokenID []byte,
	metaType int,
) *IssuingReshieldResponse {
	metadataBase := metadataCommon.MetadataBase{
		Type: metaType,
	}
	return &IssuingReshieldResponse{
		RequestedTxID:   requestedTxID,
		UniqTx:          uniqTx,
		ExternalTokenID: externalTokenID,
		MetadataBase:    metadataBase,
	}
}

func (iRes IssuingReshieldResponse) ValidateMetadataByItself() bool {
	// The validation just need to check at tx level, so returning true here
	return true
}

func (iRes IssuingReshieldResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += string(iRes.UniqTx)
	record += string(iRes.ExternalTokenID)
	record += iRes.MetadataBase.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *IssuingReshieldResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(iRes)
}
