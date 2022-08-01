package bridge

import (
	"encoding/json"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type BurnForCallRequest struct {
	BurnTokenID common.Hash              `json:"BurnTokenID"`
	Data        []BurnForCallRequestData `json:"Data"`
	metadataCommon.MetadataBase
}

type BurnForCallRequestData struct {
	BurningAmount       uint64           `json:"BurningAmount"`
	ExternalNetworkID   uint8            `json:"ExternalNetworkID"`
	IncTokenID          common.Hash      `json:"IncTokenID"`
	ExternalCalldata    string           `json:"ExternalCalldata"`
	ExternalCallAddress string           `json:"ExternalCallAddress"`
	ReceiveToken        string           `json:"ReceiveToken"`
	RedepositReceiver   coin.OTAReceiver `json:"RedepositReceiver"`
	WithdrawAddress     string           `json:"WithdrawAddress"`
}

type RejectedBurnForCallRequest struct {
	BurnTokenID common.Hash      `json:"BurnTokenID"`
	Amount      uint64           `json:"Amount"`
	Receiver    coin.OTAReceiver `json:"Receiver"`
}

func (bReq BurnForCallRequest) TotalBurningAmount() (uint64, error) {
	var totalBurningAmount uint64 = 0
	for _, d := range bReq.Data {
		totalBurningAmount += d.BurningAmount
		if totalBurningAmount < d.BurningAmount {
			return 0, fmt.Errorf("out of range uint64")
		}
	}
	return totalBurningAmount, nil
}

func (bReq BurnForCallRequest) ValidateMetadataByItself() bool {
	return bReq.Type == metadataCommon.BurnForCallRequestMeta
}

func (bReq BurnForCallRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(bReq)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (bReq *BurnForCallRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(bReq)
}

func (bReq *BurnForCallRequest) GetOTADeclarations() []metadataCommon.OTADeclaration {
	var result []metadataCommon.OTADeclaration
	for _, d := range bReq.Data {
		result = append(result, metadataCommon.OTADeclaration{
			PublicKey: d.RedepositReceiver.PublicKey.ToBytes(), TokenID: common.ConfidentialAssetID,
		})
	}
	return result
}

type BurnForCallResponse struct {
	UnshieldResponse
}

func NewBurnForCallResponseWithValue(
	status string, requestedTxID common.Hash,
) *BurnForCallResponse {
	return &BurnForCallResponse{
		UnshieldResponse{
			MetadataBase: metadataCommon.MetadataBase{
				Type: metadataCommon.BurnForCallResponseMeta,
			},
			Status:        status,
			RequestedTxID: requestedTxID,
		}}
}

func (response *BurnForCallResponse) ValidateMetadataByItself() bool {
	return response.Type == metadataCommon.BurnForCallResponseMeta
}
