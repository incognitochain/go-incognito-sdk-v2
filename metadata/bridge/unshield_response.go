package bridge

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type UnshieldResponse struct {
	metadataCommon.MetadataBase
	Status        string      `json:"Status"`
	RequestedTxID common.Hash `json:"RequestedTxID"`
}

func NewUnshieldResponse() *UnshieldResponse {
	return &UnshieldResponse{
		MetadataBase: metadataCommon.MetadataBase{
			Type: metadataCommon.BurningUnifiedTokenResponseMeta,
		},
	}
}

func NewUnshieldResponseWithValue(
	status string, requestedTxID common.Hash,
) *UnshieldResponse {
	return &UnshieldResponse{
		MetadataBase: metadataCommon.MetadataBase{
			Type: metadataCommon.BurningUnifiedTokenResponseMeta,
		},
		Status:        status,
		RequestedTxID: requestedTxID,
	}
}

func (response *UnshieldResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&response)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (response *UnshieldResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(response)
}
