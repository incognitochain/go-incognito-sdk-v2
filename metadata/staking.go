package metadata

import (
	"errors"
)

// StakingMetadata is a request to stake a node to beacon a validator.
// The user has to burn 1750 PRV to stake a node.
//	- FunderPaymentAddress: the address of the user burning 1750 PRV.
//	- CommitteePublicKey: the public key that is used in the consensus protocol.
//	- RewardReceiverPaymentAddress: the address to which staking rewards will be paid.
//	- AutoReStaking: the indicator of whether to stay staked after being swapped out of a committee.
type StakingMetadata struct {
	MetadataBase
	FunderPaymentAddress         string
	RewardReceiverPaymentAddress string
	StakingAmountShard           uint64
	AutoReStaking                bool
	CommitteePublicKey           string
}

// NewStakingMetadata creates a new StakingMetadata.
func NewStakingMetadata(
	stakingType int,
	funderPaymentAddress string,
	rewardReceiverPaymentAddress string,
	stakingAmountShard uint64,
	committeePublicKey string,
	autoReStaking bool,
) (
	*StakingMetadata,
	error,
) {
	if stakingType != ShardStakingMeta && stakingType != BeaconStakingMeta {
		return nil, errors.New("invalid staking type")
	}
	metadataBase := NewMetadataBase(stakingType)
	return &StakingMetadata{
		MetadataBase:                 *metadataBase,
		FunderPaymentAddress:         funderPaymentAddress,
		RewardReceiverPaymentAddress: rewardReceiverPaymentAddress,
		StakingAmountShard:           stakingAmountShard,
		CommitteePublicKey:           committeePublicKey,
		AutoReStaking:                autoReStaking,
	}, nil
}

// GetType overrides MetadataBase.GetType().
func (stakingMetadata StakingMetadata) GetType() int {
	return stakingMetadata.Type
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (stakingMetadata *StakingMetadata) CalculateSize() uint64 {
	return calculateSize(stakingMetadata)
}
