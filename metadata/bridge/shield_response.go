package bridge

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type ShieldResponseData struct {
	ExternalTokenID []byte      `json:"ExternalTokenID"`
	UniqTx          []byte      `json:"UniqTx"`
	IncTokenID      common.Hash `json:"IncTokenID"`
}

type ShieldResponse struct {
	metadataCommon.MetadataBase
	RequestedTxID common.Hash          `json:"RequestedTxID"`
	ShieldAmount  uint64               `json:"ShieldAmount"`
	Reward        uint64               `json:"Reward"`
	Data          []ShieldResponseData `json:"Data"`
	SharedRandom  []byte               `json:"SharedRandom,omitempty"`
}

func NewShieldResponse(metaType int) *ShieldResponse {
	return &ShieldResponse{
		MetadataBase: metadataCommon.MetadataBase{
			Type: metaType,
		},
	}
}

func NewShieldResponseWithValue(
	metaType int, shieldAmount, reward uint64, data []ShieldResponseData, requestedTxID common.Hash, sharedRandom []byte,
) *ShieldResponse {
	return &ShieldResponse{
		RequestedTxID: requestedTxID,
		ShieldAmount:  shieldAmount,
		Reward:        reward,
		Data:          data,
		SharedRandom:  sharedRandom,
		MetadataBase: metadataCommon.MetadataBase{
			Type: metaType,
		},
	}
}

func (response *ShieldResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&response)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (response *ShieldResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(response)
}

func (response *ShieldResponse) SetSharedRandom(r []byte) {
	response.SharedRandom = r
}
