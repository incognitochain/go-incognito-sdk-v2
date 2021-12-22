package pdexv3

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type ParamsModifyingRequest struct {
	metadataCommon.MetadataBaseWithSignature
	jsonresult.Pdexv3Params `json:"Pdexv3Params"`
}

type ParamsModifyingContent struct {
	Content  jsonresult.Pdexv3Params `json:"Content"`
	ErrorMsg string                  `json:"ErrorMsg"`
	TxReqID  common.Hash             `json:"TxReqID"`
	ShardID  byte                    `json:"ShardID"`
}

type ParamsModifyingRequestStatus struct {
	Status                  int    `json:"Status"`
	ErrorMsg                string `json:"ErrorMsg"`
	jsonresult.Pdexv3Params `json:"Pdexv3Params"`
}

func (paramsModifying ParamsModifyingRequest) Hash() *common.Hash {
	record := paramsModifying.MetadataBaseWithSignature.Hash().String()
	if paramsModifying.Sig != nil && len(paramsModifying.Sig) != 0 {
		record += string(paramsModifying.Sig)
	}
	contentBytes, _ := json.Marshal(paramsModifying.Pdexv3Params)
	hashParams := common.HashH(contentBytes)
	record += hashParams.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (paramsModifying ParamsModifyingRequest) HashWithoutSig() *common.Hash {
	record := paramsModifying.MetadataBaseWithSignature.Hash().String()
	contentBytes, _ := json.Marshal(paramsModifying.Pdexv3Params)
	hashParams := common.HashH(contentBytes)
	record += hashParams.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (paramsModifying *ParamsModifyingRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(paramsModifying)
}
