package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type UserMintNftRequest struct {
	metadataCommon.MetadataBase
	otaReceiver string
	amount      uint64
}

func NewUserMintNftRequestWithValue(otaReceiver string, amount uint64) *UserMintNftRequest {
	metadataBase := metadataCommon.MetadataBase{
		Type: metadataCommon.Pdexv3UserMintNftRequestMeta,
	}
	return &UserMintNftRequest{
		otaReceiver:  otaReceiver,
		amount:       amount,
		MetadataBase: metadataBase,
	}
}

func (request *UserMintNftRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&request)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (request *UserMintNftRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(request)
}

func (request *UserMintNftRequest) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		OtaReceiver string `json:"OtaReceiver"`
		Amount      uint64 `json:"Amount"`
		metadataCommon.MetadataBase
	}{
		Amount:       request.amount,
		OtaReceiver:  request.otaReceiver,
		MetadataBase: request.MetadataBase,
	})
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (request *UserMintNftRequest) UnmarshalJSON(data []byte) error {
	temp := struct {
		OtaReceiver string `json:"OtaReceiver"`
		Amount      uint64 `json:"Amount"`
		metadataCommon.MetadataBase
	}{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	request.amount = temp.Amount
	request.otaReceiver = temp.OtaReceiver
	request.MetadataBase = temp.MetadataBase
	return nil
}

func (request *UserMintNftRequest) OtaReceiver() string {
	return request.otaReceiver
}

func (request *UserMintNftRequest) Amount() uint64 {
	return request.amount
}
