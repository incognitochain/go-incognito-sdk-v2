package bridge

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type ModifyBridgeAggParamReq struct {
	metadataCommon.MetadataBaseWithSignature
	PercentFeeWithDec uint64 `json:"PercentFeeWithDec"`
}

type ModifyBridgeAggParamContentInst struct {
	PercentFeeWithDec uint64      `json:"PercentFeeWithDec"`
	TxReqID           common.Hash `json:"TxReqID"`
}

func NewModifyBridgeAggParamReq() *ModifyBridgeAggParamReq {
	return &ModifyBridgeAggParamReq{}
}

func NewModifyBridgeAggParamReqWithValue(percentFeeWithDec uint64) *ModifyBridgeAggParamReq {
	metadataBase := metadataCommon.NewMetadataBaseWithSignature(metadataCommon.BridgeAggModifyParamMeta)
	request := &ModifyBridgeAggParamReq{}
	request.MetadataBaseWithSignature = *metadataBase
	request.PercentFeeWithDec = percentFeeWithDec
	return request
}

func (request *ModifyBridgeAggParamReq) Hash() *common.Hash {
	record := request.MetadataBaseWithSignature.Hash().String()
	if request.Sig != nil && len(request.Sig) != 0 {
		record += string(request.Sig)
	}
	contentBytes, _ := json.Marshal(request)
	hashParams := common.HashH(contentBytes)
	record += hashParams.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (request *ModifyBridgeAggParamReq) HashWithoutSig() *common.Hash {
	record := request.MetadataBaseWithSignature.Hash().String()
	contentBytes, _ := json.Marshal(request.PercentFeeWithDec)
	hashParams := common.HashH(contentBytes)
	record += hashParams.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (request *ModifyBridgeAggParamReq) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(request)
}
