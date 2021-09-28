package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

// PDEContribution is a request to contribute to a pool of the pDEX.
type PDEContribution struct {
	PDEContributionPairID string
	ContributorAddressStr string
	ContributedAmount     uint64
	TokenIDStr            string
	MetadataBase
}

// NewPDEContribution creates a new PDEContribution.
func NewPDEContribution(
	pdeContributionPairID string,
	contributorAddressStr string,
	contributedAmount uint64,
	tokenIDStr string,
	metaType int,
) (*PDEContribution, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	pdeContribution := &PDEContribution{
		PDEContributionPairID: pdeContributionPairID,
		ContributorAddressStr: contributorAddressStr,
		ContributedAmount:     contributedAmount,
		TokenIDStr:            tokenIDStr,
	}
	pdeContribution.MetadataBase = metadataBase
	return pdeContribution, nil
}

// Hash overrides MetadataBase.Hash().
func (pc PDEContribution) Hash() *common.Hash {
	record := pc.MetadataBase.Hash().String()
	record += pc.PDEContributionPairID
	record += pc.ContributorAddressStr
	record += pc.TokenIDStr
	record += strconv.FormatUint(pc.ContributedAmount, 10)
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (pc *PDEContribution) CalculateSize() uint64 {
	return calculateSize(pc)
}
