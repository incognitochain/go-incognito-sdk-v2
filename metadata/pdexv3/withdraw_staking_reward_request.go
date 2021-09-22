package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type WithdrawalStakingRewardRequest struct {
	metadataCommon.MetadataBase
	StakingPoolID string                           `json:"StakingPoolID"`
	NftID         common.Hash                      `json:"NftID"`
	Receivers     map[common.Hash]coin.OTAReceiver `json:"Receivers"`
}

type WithdrawalStakingRewardContent struct {
	StakingPoolID string       `json:"StakingPoolID"`
	NftID         common.Hash  `json:"NftID"`
	TokenID       common.Hash  `json:"TokenID"`
	Receiver      ReceiverInfo `json:"Receiver"`
	IsLastInst    bool         `json:"IsLastInst"`
	TxReqID       common.Hash  `json:"TxReqID"`
	ShardID       byte         `json:"ShardID"`
}

type WithdrawalStakingRewardStatus struct {
	Status    int                          `json:"Status"`
	Receivers map[common.Hash]ReceiverInfo `json:"Receivers"`
}

func NewPdexv3WithdrawalStakingRewardRequest(
	metaType int,
	stakingToken string,
	nftID common.Hash,
	receivers map[common.Hash]coin.OTAReceiver,
) (*WithdrawalStakingRewardRequest, error) {
	metadataBase := metadataCommon.NewMetadataBase(metaType)

	return &WithdrawalStakingRewardRequest{
		MetadataBase:  *metadataBase,
		StakingPoolID: stakingToken,
		NftID:         nftID,
		Receivers:     receivers,
	}, nil
}

func (withdrawal WithdrawalStakingRewardRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(withdrawal)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (withdrawal *WithdrawalStakingRewardRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(withdrawal)
}
