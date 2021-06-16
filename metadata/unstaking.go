package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

//UnStakingMetadata : unstaking metadata
type UnStakingMetadata struct {
	MetadataBaseWithSignature
	CommitteePublicKey string
}

func (meta *UnStakingMetadata) Hash() *common.Hash {
	record := strconv.Itoa(meta.Type)
	data := []byte(record)
	hash := common.HashH(data)
	return &hash
}

func (meta *UnStakingMetadata) HashWithoutSig() *common.Hash {
	return meta.MetadataBaseWithSignature.Hash()
}

func (*UnStakingMetadata) ShouldSignMetaData() bool { return true }

//GetType :
func (unStakingMetadata UnStakingMetadata) GetType() int {
	return unStakingMetadata.Type
}

//CalculateSize :
func (unStakingMetadata *UnStakingMetadata) CalculateSize() uint64 {
	return calculateSize(unStakingMetadata)
}
