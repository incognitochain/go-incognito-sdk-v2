package bridge

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type IssuingEVMResponse struct {
	metadataCommon.MetadataBase
	RequestedTxID   common.Hash `json:"RequestedTxID"`
	UniqTx          []byte      `json:"UniqETHTx"`
	ExternalTokenID []byte      `json:"ExternalTokenID"`
	SharedRandom    []byte      `json:"SharedRandom,omitempty"`
}

type IssuingEVMResAction struct {
	Meta       *IssuingEVMResponse `json:"meta"`
	IncTokenID *common.Hash        `json:"incTokenID"`
}

func NewIssuingEVMResponse(
	requestedTxID common.Hash,
	uniqTx []byte,
	externalTokenID []byte,
	metaType int,
) *IssuingEVMResponse {
	metadataBase := metadataCommon.MetadataBase{
		Type: metaType,
	}
	return &IssuingEVMResponse{
		RequestedTxID:   requestedTxID,
		UniqTx:          uniqTx,
		ExternalTokenID: externalTokenID,
		MetadataBase:    metadataBase,
	}
}

func (iRes IssuingEVMResponse) Hash() *common.Hash {
	record := iRes.RequestedTxID.String()
	record += string(iRes.UniqTx)
	record += string(iRes.ExternalTokenID)
	record += iRes.MetadataBase.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *IssuingEVMResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(iRes)
}

func (iRes *IssuingEVMResponse) SetSharedRandom(r []byte) {
	iRes.SharedRandom = r
}
