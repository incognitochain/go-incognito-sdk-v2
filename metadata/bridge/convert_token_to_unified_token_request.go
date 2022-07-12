package bridge

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type ConvertTokenToUnifiedTokenRequest struct {
	TokenID        common.Hash      `json:"TokenID"`
	UnifiedTokenID common.Hash      `json:"UnifiedTokenID"`
	Amount         uint64           `json:"Amount"`
	Receiver       coin.OTAReceiver `json:"Receiver"`
	metadataCommon.MetadataBase
}

type RejectedConvertTokenToUnifiedToken struct {
	TokenID  common.Hash      `json:"TokenID"`
	Amount   uint64           `json:"Amount"`
	Receiver coin.OTAReceiver `json:"Receiver"`
}

type AcceptedConvertTokenToUnifiedToken struct {
	UnifiedTokenID        common.Hash      `json:"UnifiedTokenID"`
	TokenID               common.Hash      `json:"TokenID"`
	Receiver              coin.OTAReceiver `json:"Receiver"`
	ConvertPUnifiedAmount uint64           `json:"ConvertPUnifiedAmount"`
	ConvertPTokenAmount   uint64           `json:"ConvertPTokenAmount"`
	Reward                uint64           `json:"Reward"`
	TxReqID               common.Hash      `json:"TxReqID"`
}

func NewConvertTokenToUnifiedTokenRequest() *ConvertTokenToUnifiedTokenRequest {
	return &ConvertTokenToUnifiedTokenRequest{}
}

func NewConvertTokenToUnifiedTokenRequestWithValue(
	tokenID, unifiedTokenID common.Hash, amount uint64, receiver coin.OTAReceiver,
) *ConvertTokenToUnifiedTokenRequest {
	metadataBase := metadataCommon.MetadataBase{
		Type: metadataCommon.BridgeAggConvertTokenToUnifiedTokenRequestMeta,
	}
	return &ConvertTokenToUnifiedTokenRequest{
		UnifiedTokenID: unifiedTokenID,
		TokenID:        tokenID,
		Amount:         amount,
		Receiver:       receiver,
		MetadataBase:   metadataBase,
	}
}

func (request *ConvertTokenToUnifiedTokenRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&request)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (request *ConvertTokenToUnifiedTokenRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(request)
}

func (request *ConvertTokenToUnifiedTokenRequest) GetOTADeclarations() []metadataCommon.OTADeclaration {
	var result []metadataCommon.OTADeclaration
	result = append(result, metadataCommon.OTADeclaration{
		PublicKey: request.Receiver.PublicKey.ToBytes(), TokenID: common.ConfidentialAssetID,
	})
	return result
}
