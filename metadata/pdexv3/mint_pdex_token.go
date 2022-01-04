package pdexv3

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type MintPDEXGenesisResponse struct {
	metadataCommon.MetadataBase
	MintingPaymentAddress string `json:"MintingPaymentAddress"`
	MintingAmount         uint64 `json:"MintingAmount"`
	SharedRandom          []byte `json:"SharedRandom"`
}

type MintBlockRewardContent struct {
	PoolPairID string      `json:"PoolPairID"`
	Amount     uint64      `json:"Amount"`
	TokenID    common.Hash `json:"TokenID"`
}

type MintPDEXGenesisContent struct {
	MintingPaymentAddress string `json:"MintingPaymentAddress"`
	MintingAmount         uint64 `json:"MintingAmount"`
	ShardID               byte   `json:"ShardID"`
}

func (mintResponse MintPDEXGenesisResponse) Hash() *common.Hash {
	record := mintResponse.MetadataBase.Hash().String()
	record += mintResponse.MintingPaymentAddress
	record += strconv.FormatUint(mintResponse.MintingAmount, 10)

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (mintResponse *MintPDEXGenesisResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(mintResponse)
}

func (mintResponse *MintPDEXGenesisResponse) SetSharedRandom(r []byte) {
	mintResponse.SharedRandom = r
}
