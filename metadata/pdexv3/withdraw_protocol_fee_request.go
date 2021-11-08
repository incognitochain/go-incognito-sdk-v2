package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type WithdrawalProtocolFeeRequest struct {
	metadataCommon.MetadataBaseWithSignature
	PoolPairID string `json:"PoolPairID"`
}

type WithdrawalProtocolFeeContent struct {
	PoolPairID string      `json:"PoolPairID"`
	Address    string      `json:"Address"`
	TokenID    common.Hash `json:"TokenID"`
	Amount     uint64      `json:"Amount"`
	IsLastInst bool        `json:"IsLastInst"`
	TxReqID    common.Hash `json:"TxReqID"`
	ShardID    byte        `json:"ShardID"`
}

type WithdrawalProtocolFeeStatus struct {
	Status int                    `json:"Status"`
	Amount map[common.Hash]uint64 `json:"Amount"`
}

func NewPdexv3WithdrawalProtocolFeeRequest(
	metaType int,
	pairID string,
) (*WithdrawalProtocolFeeRequest, error) {
	metadataBase := metadataCommon.NewMetadataBaseWithSignature(metaType)

	return &WithdrawalProtocolFeeRequest{
		MetadataBaseWithSignature: *metadataBase,
		PoolPairID:                pairID,
	}, nil
}

func (withdrawal WithdrawalProtocolFeeRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(withdrawal)
	if withdrawal.Sig != nil && len(withdrawal.Sig) != 0 {
		rawBytes = append(rawBytes, withdrawal.Sig...)
	}

	// final hash
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (withdrawal WithdrawalProtocolFeeRequest) HashWithoutSig() *common.Hash {
	rawBytes, _ := json.Marshal(struct {
		Type       int    `json:"Type"`
		PoolPairID string `json:"PoolPairID"`
	}{
		Type:       metadataCommon.Pdexv3WithdrawProtocolFeeRequestMeta,
		PoolPairID: withdrawal.PoolPairID,
	})

	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (withdrawal *WithdrawalProtocolFeeRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(withdrawal)
}
