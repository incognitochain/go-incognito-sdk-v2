package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

type WithDrawRewardResponse struct {
	MetadataBase
	TxRequest       *common.Hash
	TokenID         common.Hash
	RewardPublicKey []byte
	SharedRandom    []byte
	Version         int
}

func (withDrawRewardResponse WithDrawRewardResponse) Hash() *common.Hash {
	if withDrawRewardResponse.Version == common.SALARY_VER_FIX_HASH {
		if withDrawRewardResponse.TxRequest == nil {
			return &common.Hash{}
		}
		bArr := append(withDrawRewardResponse.TxRequest.GetBytes(), withDrawRewardResponse.TokenID.GetBytes()...)
		version := strconv.Itoa(withDrawRewardResponse.Version)
		if len(withDrawRewardResponse.SharedRandom) != 0 {
			bArr = append(bArr, withDrawRewardResponse.SharedRandom...)
		}
		if len(withDrawRewardResponse.RewardPublicKey) != 0 {
			bArr = append(bArr, withDrawRewardResponse.RewardPublicKey...)
		}

		bArr = append(bArr, []byte(version)...)
		txResHash := common.HashH(bArr)
		return &txResHash
	} else {
		return withDrawRewardResponse.TxRequest
	}
}