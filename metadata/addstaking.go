package metadata

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
