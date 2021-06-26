package metadata

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

// StopAutoStakingMetadata is a request to stop staking for a staked node.
type StopAutoStakingMetadata struct {
	MetadataBaseWithSignature
	CommitteePublicKey string
}

// NewStopAutoStakingMetadata creates a new StopAutoStakingMetadata.
func NewStopAutoStakingMetadata(stopStakingType int, committeePublicKey string) (*StopAutoStakingMetadata, error) {
	if stopStakingType != StopAutoStakingMeta {
		return nil, fmt.Errorf("invalid stop staking type")
	}
	metadataBase := NewMetadataBaseWithSignature(stopStakingType)
	return &StopAutoStakingMetadata{
		MetadataBaseWithSignature: *metadataBase,
		CommitteePublicKey:        committeePublicKey,
	}, nil
}

// Hash overrides MetadataBase.Hash().
func (req *StopAutoStakingMetadata) Hash() *common.Hash {
	record := strconv.Itoa(req.Type)
	data := []byte(record)
	data = append(data, req.Sig...)
	hash := common.HashH(data)
	return &hash
}

// HashWithoutSig overrides MetadataBase.HashWithoutSig().
func (req *StopAutoStakingMetadata) HashWithoutSig() *common.Hash {
	return req.MetadataBase.Hash()
}

// ShouldSignMetaData returns true
func (*StopAutoStakingMetadata) ShouldSignMetaData() bool { return true }

// GetType overrides MetadataBase.GetType().
func (req StopAutoStakingMetadata) GetType() int {
	return req.Type
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (req *StopAutoStakingMetadata) CalculateSize() uint64 {
	return calculateSize(req)
}
