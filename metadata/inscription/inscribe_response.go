package inscription

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type InscribeResponse struct {
	metadataCommon.MetadataBase
	status  string
	txReqID string
}

func NewInscribeResponseWithValue(status, txReqID string) *InscribeResponse {
	return &InscribeResponse{
		MetadataBase: metadataCommon.MetadataBase{
			Type: metadataCommon.InscribeResponseMeta,
		},
		status:  status,
		txReqID: txReqID,
	}
}

func (response *InscribeResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&response)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (response *InscribeResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(response)
}

func (response *InscribeResponse) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		Status  string `json:"Status"`
		TxReqID string `json:"TxReqID"`
		metadataCommon.MetadataBase
	}{
		Status:       response.status,
		TxReqID:      response.txReqID,
		MetadataBase: response.MetadataBase,
	})
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (response *InscribeResponse) UnmarshalJSON(data []byte) error {
	temp := struct {
		Status  string `json:"Status"`
		TxReqID string `json:"TxReqID"`
		metadataCommon.MetadataBase
	}{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	response.txReqID = temp.TxReqID
	response.status = temp.Status
	response.MetadataBase = temp.MetadataBase
	return nil
}

func (response *InscribeResponse) TxReqID() string {
	return response.txReqID
}

func (response *InscribeResponse) Status() string {
	return response.status
}

type MintNftData struct {
	NftID       common.Hash `json:"NftID"`
	OtaReceiver string      `json:"OtaReceiver"`
	ShardID     byte        `json:"ShardID"`
}
