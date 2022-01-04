package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type UnstakingResponse struct {
	metadataCommon.MetadataBase
	status  string
	txReqID string
}

func (response *UnstakingResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&response)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (response *UnstakingResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(response)
}

func (response *UnstakingResponse) MarshalJSON() ([]byte, error) {
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

func (response *UnstakingResponse) UnmarshalJSON(data []byte) error {
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

func (response *UnstakingResponse) TxReqID() string {
	return response.txReqID
}

func (response *UnstakingResponse) Status() string {
	return response.status
}
