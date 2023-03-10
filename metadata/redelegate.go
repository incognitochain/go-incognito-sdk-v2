package metadata

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type ReDelegateMetadata struct {
	MetadataBaseWithSignature
	CommitteePublicKey string
	NewDelegate        string
	DelegateUID        string
}

func (meta *ReDelegateMetadata) Hash() *common.Hash {
	record := strconv.Itoa(meta.Type)
	data := []byte(record)
	data = append(data, meta.Sig...)
	data = append(data, []byte(meta.CommitteePublicKey)...)
	data = append(data, []byte(meta.NewDelegate)...)
	data = append(data, []byte(meta.DelegateUID)...)
	hash := common.HashH(data)
	return &hash
}

func (meta *ReDelegateMetadata) HashWithoutSig() *common.Hash {
	record := strconv.Itoa(meta.Type)
	data := []byte(record)
	data = append(data, []byte(meta.CommitteePublicKey)...)
	data = append(data, []byte(meta.NewDelegate)...)
	data = append(data, []byte(meta.DelegateUID)...)
	hash := common.HashH(data)
	return &hash
}

func NewReDelegateMetadata(committeePublicKey, newDelegate string, newDelegateUID string) (*ReDelegateMetadata, error) {
	metadataBase := NewMetadataBaseWithSignature(ReDelegateMeta)
	return &ReDelegateMetadata{
		MetadataBaseWithSignature: *metadataBase,
		CommitteePublicKey:        committeePublicKey,
		NewDelegate:               newDelegate,
		DelegateUID:               newDelegateUID,
	}, nil
}

/*
 */

func (redelegateMetadata ReDelegateMetadata) GetType() int {
	return redelegateMetadata.Type
}

func (redelegateMetadata *ReDelegateMetadata) CalculateSize() uint64 {
	return calculateSize(redelegateMetadata)
}

// ShouldSignMetaData returns true
func (*ReDelegateMetadata) ShouldSignMetaData() bool { return true }
