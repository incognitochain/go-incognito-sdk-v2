package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type MintNftResponse struct {
	nftID       string
	otaReceiver string
	metadataCommon.MetadataBase
}

func (response *MintNftResponse) Hash() *common.Hash {
	record := response.MetadataBase.Hash().String()
	record += response.nftID
	record += response.otaReceiver
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (response *MintNftResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(response)
}

func (response *MintNftResponse) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		NftID       string `json:"NftID"`
		OtaReceiver string `json:"OtaReceiver"`
		metadataCommon.MetadataBase
	}{
		NftID:        response.nftID,
		OtaReceiver:  response.otaReceiver,
		MetadataBase: response.MetadataBase,
	})
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (response *MintNftResponse) UnmarshalJSON(data []byte) error {
	temp := struct {
		NftID       string `json:"NftID"`
		OtaReceiver string `json:"OtaReceiver"`
		metadataCommon.MetadataBase
	}{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	response.otaReceiver = temp.OtaReceiver
	response.nftID = temp.NftID
	response.MetadataBase = temp.MetadataBase
	return nil
}

func (response *MintNftResponse) OtaReceiver() string {
	return response.otaReceiver
}

func (response *MintNftResponse) NftID() string {
	return response.nftID
}

type MintNftData struct {
	NftID       common.Hash `json:"NftID"`
	OtaReceiver string      `json:"OtaReceiver"`
	ShardID     byte        `json:"ShardID"`
}
