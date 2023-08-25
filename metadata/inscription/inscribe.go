package inscription

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type InscribeRequest struct {
	Data     string           `json:"Data"`
	Receiver coin.OTAReceiver `json:"Receiver"`
	metadataCommon.MetadataBase
}

func (iReq InscribeRequest) GetType() int {
	return iReq.Type
}

// NewInscribeRequest creates a new InscribeRequest.
func NewInscribeRequest(
	data string,
	receiver coin.OTAReceiver,
	metaType int,
) (*InscribeRequest, error) {
	metadataBase := metadataCommon.MetadataBase{
		Type: metaType,
	}

	inscribeReq := &InscribeRequest{
		Data:     data,
		Receiver: receiver,
	}
	inscribeReq.MetadataBase = metadataBase
	return inscribeReq, nil
}

func (iReq InscribeRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(iReq)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (iReq *InscribeRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(iReq)
}
