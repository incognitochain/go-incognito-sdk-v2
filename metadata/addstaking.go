package metadata

import (
	"fmt"
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type AddStakingMetadata struct {
	MetadataBaseWithSignature
	CommitteePublicKey string
	AddStakingAmount   uint64
}

func NewAddStakingMetadata(committeePublicKey string, addStakingAmount uint64) (*AddStakingMetadata, error) {
	metadataBase := NewMetadataBaseWithSignature(AddStakingMeta)
	return &AddStakingMetadata{
		MetadataBaseWithSignature: *metadataBase,
		CommitteePublicKey:        committeePublicKey,
		AddStakingAmount:          addStakingAmount,
	}, nil
}

func (addStakingMetadata AddStakingMetadata) GetType() int {
	return addStakingMetadata.Type
}

func (addStakingMetadata *AddStakingMetadata) CalculateSize() uint64 {
	return calculateSize(addStakingMetadata)
}

func (meta *AddStakingMetadata) Hash() *common.Hash {
	record := strconv.Itoa(meta.Type)
	data := []byte(record)
	data = append(data, meta.Sig...)
	data = append(data, []byte(meta.CommitteePublicKey)...)
	data = append(data, []byte(fmt.Sprintf("%v", meta.AddStakingAmount))...)
	hash := common.HashH(data)
	return &hash
}

func (meta *AddStakingMetadata) HashWithoutSig() *common.Hash {
	record := strconv.Itoa(meta.Type)
	data := []byte(record)
	data = append(data, []byte(meta.CommitteePublicKey)...)
	data = append(data, []byte(fmt.Sprintf("%v", meta.AddStakingAmount))...)
	hash := common.HashH(data)
	return &hash
}

// ShouldSignMetaData returns true
func (*AddStakingMetadata) ShouldSignMetaData() bool { return true }
