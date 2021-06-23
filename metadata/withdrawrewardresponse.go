package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

// WithDrawRewardResponse is the response for a WithDrawRewardRequest.
type WithDrawRewardResponse struct {
	MetadataBase
	TxRequest       *common.Hash
	TokenID         common.Hash
	RewardPublicKey []byte
	SharedRandom    []byte
	Version         int
}

// Hash overrides MetadataBase.Hash().
func (resp WithDrawRewardResponse) Hash() *common.Hash {
	if resp.Version == common.SalaryVerFixHash {
		if resp.TxRequest == nil {
			return &common.Hash{}
		}
		bArr := append(resp.TxRequest.GetBytes(), resp.TokenID.GetBytes()...)
		version := strconv.Itoa(resp.Version)
		if len(resp.SharedRandom) != 0 {
			bArr = append(bArr, resp.SharedRandom...)
		}
		if len(resp.RewardPublicKey) != 0 {
			bArr = append(bArr, resp.RewardPublicKey...)
		}

		bArr = append(bArr, []byte(version)...)
		txResHash := common.HashH(bArr)
		return &txResHash
	} else {
		return resp.TxRequest
	}
}
