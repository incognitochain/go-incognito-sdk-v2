package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type StakingRequest struct {
	metadataCommon.MetadataBase

	// tokenID is the token we wish to stake. This token must be in the list of allowed staking tokens.
	tokenID string

	// otaReceiver is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
	otaReceiver string

	// nftID is the ID of the NFT associated with the staking request.
	nftID string

	// tokenAmount is the staking amount.
	tokenAmount uint64
}

func NewStakingRequestWithValue(
	tokenID, nftID, otaReceiver string, tokenAmount uint64,
) *StakingRequest {
	return &StakingRequest{
		MetadataBase: metadataCommon.MetadataBase{
			Type: metadataCommon.Pdexv3StakingRequestMeta,
		},
		tokenID:     tokenID,
		nftID:       nftID,
		tokenAmount: tokenAmount,
		otaReceiver: otaReceiver,
	}
}

func (request *StakingRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&request)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (request *StakingRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(request)
}

func (request *StakingRequest) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		OtaReceiver string `json:"OtaReceiver"`
		TokenID     string `json:"TokenID"`
		NftID       string `json:"NftID"`
		TokenAmount uint64 `json:"TokenAmount"`
		metadataCommon.MetadataBase
	}{
		OtaReceiver:  request.otaReceiver,
		TokenID:      request.tokenID,
		NftID:        request.nftID,
		TokenAmount:  request.tokenAmount,
		MetadataBase: request.MetadataBase,
	})
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (request *StakingRequest) UnmarshalJSON(data []byte) error {
	temp := struct {
		OtaReceiver string `json:"OtaReceiver"`
		TokenID     string `json:"TokenID"`
		NftID       string `json:"NftID"`
		TokenAmount uint64 `json:"TokenAmount"`
		metadataCommon.MetadataBase
	}{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	request.otaReceiver = temp.OtaReceiver
	request.tokenID = temp.TokenID
	request.nftID = temp.NftID
	request.tokenAmount = temp.TokenAmount
	request.MetadataBase = temp.MetadataBase
	return nil
}

func (request *StakingRequest) OtaReceiver() string {
	return request.otaReceiver
}

func (request *StakingRequest) TokenID() string {
	return request.tokenID
}

func (request *StakingRequest) TokenAmount() uint64 {
	return request.tokenAmount
}

func (request *StakingRequest) NftID() string {
	return request.nftID
}
