package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

type ReturnStakingMetadata struct {
	MetadataBase
	TxID          string
	StakerAddress key.PaymentAddress
	SharedRandom []byte `json:"SharedRandom,omitempty"`
}

func (sbsRes ReturnStakingMetadata) Hash() *common.Hash {
	record := sbsRes.StakerAddress.String()
	record += sbsRes.TxID
	if sbsRes.SharedRandom != nil && len(sbsRes.SharedRandom) > 0 {
		record += string(sbsRes.SharedRandom)
	}
	// final hash
	record += sbsRes.MetadataBase.Hash().String()
	hash := common.HashH([]byte(record))
	return &hash
}

func (sbsRes *ReturnStakingMetadata) SetSharedRandom(r []byte) {
	sbsRes.SharedRandom = r
}

