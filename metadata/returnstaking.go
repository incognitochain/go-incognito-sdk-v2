package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// ReturnStakingMetadata is a Metadata for returning the staking amount after a un-staking request.
type ReturnStakingMetadata struct {
	MetadataBase
	TxID          string
	StakerAddress key.PaymentAddress
	SharedRandom  []byte `json:"SharedRandom,omitempty"`
}

// Hash overrides MetadataBase.Hash().
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

// SetSharedRandom sets v as the shared random of a ReturnStakingMetadata.
func (sbsRes *ReturnStakingMetadata) SetSharedRandom(v []byte) {
	sbsRes.SharedRandom = v
}
