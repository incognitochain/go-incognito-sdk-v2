package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type WithdrawalStakingRewardResponse struct {
	metadataCommon.MetadataBase
	ReqTxID common.Hash `json:"ReqTxID"`
}

type DistributeStakingRewardContent struct {
	StakingPoolID string                 `json:"StakingPoolID"`
	Rewards       map[common.Hash]uint64 `json:"Rewards"`
}

func (withdrawalResponse WithdrawalStakingRewardResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(withdrawalResponse)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (withdrawalResponse *WithdrawalStakingRewardResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(withdrawalResponse)
}
