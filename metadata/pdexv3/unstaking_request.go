package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type UnstakingRequest struct {
	metadataCommon.MetadataBase
	stakingPoolID   string
	otaReceivers    map[string]string
	nftID           string
	unstakingAmount uint64
}

func NewUnstakingRequestWithValue(
	stakingPoolID, nftID string,
	otaReceivers map[string]string,
	unstakingAmount uint64,
) *UnstakingRequest {
	return &UnstakingRequest{
		MetadataBase: metadataCommon.MetadataBase{
			Type: metadataCommon.Pdexv3UnstakingRequestMeta,
		},
		stakingPoolID:   stakingPoolID,
		nftID:           nftID,
		otaReceivers:    otaReceivers,
		unstakingAmount: unstakingAmount,
	}
}

func (request *UnstakingRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&request)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (request *UnstakingRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(request)
}

func (request *UnstakingRequest) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		StakingPoolID   string            `json:"StakingPoolID"`
		NftID           string            `json:"NftID"`
		OtaReceivers    map[string]string `json:"OtaReceivers"`
		UnstakingAmount uint64            `json:"UnstakingAmount"`
		metadataCommon.MetadataBase
	}{
		StakingPoolID:   request.stakingPoolID,
		NftID:           request.nftID,
		OtaReceivers:    request.otaReceivers,
		UnstakingAmount: request.unstakingAmount,
		MetadataBase:    request.MetadataBase,
	})
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (request *UnstakingRequest) UnmarshalJSON(data []byte) error {
	temp := struct {
		StakingPoolID   string            `json:"StakingPoolID"`
		NftID           string            `json:"NftID"`
		OtaReceivers    map[string]string `json:"OtaReceivers"`
		UnstakingAmount uint64            `json:"UnstakingAmount"`
		metadataCommon.MetadataBase
	}{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	request.stakingPoolID = temp.StakingPoolID
	request.nftID = temp.NftID
	request.otaReceivers = temp.OtaReceivers
	request.unstakingAmount = temp.UnstakingAmount
	request.MetadataBase = temp.MetadataBase
	return nil
}

func (request *UnstakingRequest) StakingPoolID() string {
	return request.stakingPoolID
}

func (request *UnstakingRequest) OtaReceivers() map[string]string {
	return request.otaReceivers
}

func (request *UnstakingRequest) UnstakingAmount() uint64 {
	return request.unstakingAmount
}

func (request *UnstakingRequest) NftID() string {
	return request.nftID
}
