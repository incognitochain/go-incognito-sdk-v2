package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type WithdrawalLPFeeRequest struct {
	metadataCommon.MetadataBase

	// PoolPairID is the ID of the target pool pair.
	PoolPairID string                           `json:"PoolPairID"`

	// NftID is the ID of the NFT which he used to make contribution.
	NftID      common.Hash                      `json:"NftID"`

	// is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
	Receivers  map[common.Hash]coin.OTAReceiver `json:"Receivers"`
}

type WithdrawalLPFeeContent struct {
	PoolPairID string       `json:"PoolPairID"`
	NftID      common.Hash  `json:"NftID"`
	TokenID    common.Hash  `json:"TokenID"`
	Receiver   ReceiverInfo `json:"Receiver"`
	IsLastInst bool         `json:"IsLastInst"`
	TxReqID    common.Hash  `json:"TxReqID"`
	ShardID    byte         `json:"ShardID"`
}

type WithdrawalLPFeeStatus struct {
	Status    int                          `json:"Status"`
	Receivers map[common.Hash]ReceiverInfo `json:"Receivers"`
}

func NewPdexv3WithdrawalLPFeeRequest(
	metaType int,
	pairID string,
	nftID common.Hash,
	receivers map[common.Hash]coin.OTAReceiver,
) (*WithdrawalLPFeeRequest, error) {
	metadataBase := metadataCommon.NewMetadataBase(metaType)

	return &WithdrawalLPFeeRequest{
		MetadataBase: *metadataBase,
		PoolPairID:   pairID,
		NftID:        nftID,
		Receivers:    receivers,
	}, nil
}

func (withdrawal WithdrawalLPFeeRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(withdrawal)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (withdrawal *WithdrawalLPFeeRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(withdrawal)
}
