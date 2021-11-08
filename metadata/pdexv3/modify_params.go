package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type Pdexv3Params struct {
	DefaultFeeRateBPS               uint            `json:"DefaultFeeRateBPS"`
	FeeRateBPS                      map[string]uint `json:"FeeRateBPS"`
	PRVDiscountPercent              uint            `json:"PRVDiscountPercent"`
	TradingProtocolFeePercent       uint            `json:"TradingProtocolFeePercent"`
	TradingStakingPoolRewardPercent uint            `json:"TradingStakingPoolRewardPercent"`
	PDEXRewardPoolPairsShare        map[string]uint `json:"PDEXRewardPoolPairsShare"`
	StakingPoolsShare               map[string]uint `json:"StakingPoolsShare"`
	StakingRewardTokens             []common.Hash   `json:"StakingRewardTokens"`
	MintNftRequireAmount            uint64          `json:"MintNftRequireAmount"`
	MaxOrdersPerNft                 uint            `json:"MaxOrdersPerNft"`
}

type ParamsModifyingRequest struct {
	metadataCommon.MetadataBaseWithSignature
	Pdexv3Params `json:"Pdexv3Params"`
}

type ParamsModifyingContent struct {
	Content  Pdexv3Params `json:"Content"`
	ErrorMsg string       `json:"ErrorMsg"`
	TxReqID  common.Hash  `json:"TxReqID"`
	ShardID  byte         `json:"ShardID"`
}

type ParamsModifyingRequestStatus struct {
	Status       int    `json:"Status"`
	ErrorMsg     string `json:"ErrorMsg"`
	Pdexv3Params `json:"Pdexv3Params"`
}

func (paramsModifying ParamsModifyingRequest) Hash() *common.Hash {
	record := paramsModifying.MetadataBaseWithSignature.Hash().String()
	if paramsModifying.Sig != nil && len(paramsModifying.Sig) != 0 {
		record += string(paramsModifying.Sig)
	}
	contentBytes, _ := json.Marshal(paramsModifying.Pdexv3Params)
	hashParams := common.HashH(contentBytes)
	record += hashParams.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (paramsModifying ParamsModifyingRequest) HashWithoutSig() *common.Hash {
	record := paramsModifying.MetadataBaseWithSignature.Hash().String()
	contentBytes, _ := json.Marshal(paramsModifying.Pdexv3Params)
	hashParams := common.HashH(contentBytes)
	record += hashParams.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (paramsModifying *ParamsModifyingRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(paramsModifying)
}
