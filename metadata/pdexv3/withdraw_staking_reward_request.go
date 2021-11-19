package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type WithdrawalStakingRewardRequest struct {
	metadataCommon.MetadataBase

	// StakingPoolID
	StakingPoolID string `json:"StakingPoolID"`

	// NftID is theID of the NFT associated with the staking request.
	NftID common.Hash `json:"NftID"`

	// Receivers is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
	Receivers map[common.Hash]coin.OTAReceiver `json:"Receivers"`
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
