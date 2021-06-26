package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

// UnStakingMetadata is a request to un-stake a running validator.
// The node will un-stake faster when using this metadata compared to the StopAutoStakingMetadata.
type UnStakingMetadata struct {
	MetadataBaseWithSignature
	CommitteePublicKey string
}

// Hash overrides MetadataBase.Hash().
func (req *UnStakingMetadata) Hash() *common.Hash {
	record := strconv.Itoa(req.Type)
	data := []byte(record)
	hash := common.HashH(data)
	return &hash
}

// HashWithoutSig overrides MetadataBase.HashWithoutSig().
func (req *UnStakingMetadata) HashWithoutSig() *common.Hash {
	return req.MetadataBaseWithSignature.Hash()
}

// ShouldSignMetaData returns true.
func (*UnStakingMetadata) ShouldSignMetaData() bool { return true }

// GetType overrides MetadataBase.GetType().
func (req UnStakingMetadata) GetType() int {
	return req.Type
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (req *UnStakingMetadata) CalculateSize() uint64 {
	return calculateSize(req)
}
