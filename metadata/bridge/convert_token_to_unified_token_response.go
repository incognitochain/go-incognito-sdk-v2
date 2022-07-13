package bridge

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type ConvertTokenToUnifiedTokenResponse struct {
	metadataCommon.MetadataBase
	ConvertAmount uint64      `json:"ConvertAmount"`
	Reward        uint64      `json:"Reward"`
	Status        string      `json:"Status"`
	TxReqID       common.Hash `json:"TxReqID"`
}

func NewConvertTokenToUnifiedTokenResponse() *ConvertTokenToUnifiedTokenResponse {
	return &ConvertTokenToUnifiedTokenResponse{
		MetadataBase: metadataCommon.MetadataBase{
			Type: metadataCommon.BridgeAggConvertTokenToUnifiedTokenResponseMeta,
		},
	}
}

func NewBridgeAggConvertTokenToUnifiedTokenResponseWithValue(
	status string, txReqID common.Hash, convertAmount uint64, reward uint64,
) *ConvertTokenToUnifiedTokenResponse {
	return &ConvertTokenToUnifiedTokenResponse{
		MetadataBase: metadataCommon.MetadataBase{
			Type: metadataCommon.BridgeAggConvertTokenToUnifiedTokenResponseMeta,
		},
		ConvertAmount: convertAmount,
		Reward:        reward,
		Status:        status,
		TxReqID:       txReqID,
	}
}

func (response *ConvertTokenToUnifiedTokenResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&response)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (response *ConvertTokenToUnifiedTokenResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(response)
}
