package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type AddLiquidityRequest struct {
	// poolPairID is the ID of the target pool in which the LP wants to add liquidity to.
	// If this is the first contribution (i.e, pool-initialization), the poolPairID must be left empty.
	poolPairID string

	// pairHash is a string for matching the two contributing transactions. It can be anything as long as it is the same in
	// both contributing transaction.
	pairHash string

	// otaReceiver is a one-time address for receiving back the token in case of being refunded.
	otaReceiver string

	// tokenID is the ID of the contributing token.
	tokenID string

	// nftID is the ID of the NFT associated with this contribution. This value must be the same in both contributing transactions.
	nftID string

	// tokenAmount is the contributing amount of this token.
	tokenAmount uint64

	// amplifier is the amplifier of the pool. In the case of contributing to an existing pool, this value must match that of the existing pool.
	// The detail of this param can be found in the Kyber's White-paper (https://files.kyber.network/DMM-Feb21.pdf).
	amplifier uint

	metadataCommon.MetadataBase
}

func NewAddLiquidityRequestWithValue(
	poolPairID, pairHash, otaReceiver, tokenID, nftID string,
	tokenAmount uint64, amplifier uint,
) *AddLiquidityRequest {
	metadataBase := metadataCommon.MetadataBase{
		Type: metadataCommon.Pdexv3AddLiquidityRequestMeta,
	}
	return &AddLiquidityRequest{
		poolPairID:   poolPairID,
		pairHash:     pairHash,
		otaReceiver:  otaReceiver,
		tokenID:      tokenID,
		nftID:        nftID,
		tokenAmount:  tokenAmount,
		amplifier:    amplifier,
		MetadataBase: metadataBase,
	}
}

func (request *AddLiquidityRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(&request)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (request *AddLiquidityRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(request)
}

func (request *AddLiquidityRequest) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		PoolPairID  string `json:"PoolPairID"` // only "" for the first contribution of pool
		PairHash    string `json:"PairHash"`
		OtaReceiver string `json:"OtaReceiver"` // receive pToken
		TokenID     string `json:"TokenID"`
		NftID       string `json:"NftID"`
		TokenAmount uint64 `json:"TokenAmount"`
		Amplifier   uint   `json:"Amplifier"` // only set for the first contribution
		metadataCommon.MetadataBase
	}{
		PoolPairID:   request.poolPairID,
		PairHash:     request.pairHash,
		OtaReceiver:  request.otaReceiver,
		TokenID:      request.tokenID,
		NftID:        request.nftID,
		TokenAmount:  request.tokenAmount,
		Amplifier:    request.amplifier,
		MetadataBase: request.MetadataBase,
	})
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (request *AddLiquidityRequest) UnmarshalJSON(data []byte) error {
	temp := struct {
		PoolPairID  string `json:"PoolPairID"` // only "" for the first contribution of pool
		PairHash    string `json:"PairHash"`
		OtaReceiver string `json:"OtaReceiver"` // receive pToken
		TokenID     string `json:"TokenID"`
		NftID       string `json:"NftID"`
		TokenAmount uint64 `json:"TokenAmount"`
		Amplifier   uint   `json:"Amplifier"` // only set for the first contribution
		metadataCommon.MetadataBase
	}{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	request.poolPairID = temp.PoolPairID
	request.pairHash = temp.PairHash
	request.otaReceiver = temp.OtaReceiver
	request.tokenID = temp.TokenID
	request.nftID = temp.NftID
	request.tokenAmount = temp.TokenAmount
	request.amplifier = temp.Amplifier
	request.MetadataBase = temp.MetadataBase
	return nil
}

func (request *AddLiquidityRequest) PoolPairID() string {
	return request.poolPairID
}

func (request *AddLiquidityRequest) PairHash() string {
	return request.pairHash
}

func (request *AddLiquidityRequest) OtaReceiver() string {
	return request.otaReceiver
}

func (request *AddLiquidityRequest) TokenID() string {
	return request.tokenID
}

func (request *AddLiquidityRequest) TokenAmount() uint64 {
	return request.tokenAmount
}

func (request *AddLiquidityRequest) Amplifier() uint {
	return request.amplifier
}

func (request *AddLiquidityRequest) NftID() string {
	return request.nftID
}
