package metadata

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type ShieldRequest struct {
	Data           []ShieldRequestData `json:"Data"`
	UnifiedTokenID common.Hash         `json:"UnifiedTokenID"`
	MetadataBase
}

type ShieldRequestData struct {
	Proof      []byte      `json:"Proof"`
	NetworkID  uint8       `json:"NetworkID"`
	IncTokenID common.Hash `json:"IncTokenID"`
}

type AcceptedInstShieldRequest struct {
	Receiver       key.PaymentAddress          `json:"Receiver"`
	UnifiedTokenID common.Hash                 `json:"UnifiedTokenID"`
	TxReqID        common.Hash                 `json:"TxReqID"`
	Data           []AcceptedShieldRequestData `json:"Data"`
}

type AcceptedShieldRequestData struct {
	ShieldAmount    uint64      `json:"ShieldAmount"`
	Reward          uint64      `json:"Reward"`
	UniqTx          []byte      `json:"UniqTx"`
	ExternalTokenID []byte      `json:"ExternalTokenID"`
	NetworkID       uint8       `json:"NetworkID"`
	IncTokenID      common.Hash `json:"IncTokenID"`
}

func NewShieldRequest() *ShieldRequest {
	return &ShieldRequest{
		MetadataBase: MetadataBase{
			Type: metadataCommon.IssuingUnifiedTokenRequestMeta,
		},
	}
}

func NewShieldRequestWithValue(
	data []ShieldRequestData, unifiedTokenID common.Hash,
) *ShieldRequest {
	return &ShieldRequest{
		Data:           data,
		UnifiedTokenID: unifiedTokenID,
		MetadataBase: MetadataBase{
			Type: metadataCommon.IssuingUnifiedTokenRequestMeta,
		},
	}
}
func (request *ShieldRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&request)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (request *ShieldRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(request)
}

func UnmarshalActionDataForShieldEVMReq(data []byte) (*types.Receipt, error) {
	txReceipt := types.Receipt{}
	err := json.Unmarshal(data, &txReceipt)
	return &txReceipt, err
}

func MarshalActionDataForShieldEVMReq(txReceipt *types.Receipt) ([]byte, error) {
	return json.Marshal(txReceipt)
}
