package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

// PDEFeeWithdrawalRequest is a request to withdraw trading fee rewards from the pDEX.
// The user needs to sign this request to make sure he/she is authorized to withdraw the rewards.
type PDEFeeWithdrawalRequest struct {
	WithdrawerAddressStr  string
	WithdrawalToken1IDStr string
	WithdrawalToken2IDStr string
	WithdrawalFeeAmt      uint64
	MetadataBaseWithSignature
}

// Hash overrides MetadataBase.Hash().
func (pc PDEFeeWithdrawalRequest) Hash() *common.Hash {
	record := pc.MetadataBase.Hash().String()
	record += pc.WithdrawerAddressStr
	record += pc.WithdrawalToken1IDStr
	record += pc.WithdrawalToken2IDStr
	record += strconv.FormatUint(pc.WithdrawalFeeAmt, 10)
	if pc.Sig != nil && len(pc.Sig) != 0 {
		record += string(pc.Sig)
	}
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// HashWithoutSig overrides MetadataBase.HashWithoutSig().
func (pc PDEFeeWithdrawalRequest) HashWithoutSig() *common.Hash {
	record := pc.MetadataBase.Hash().String()
	record += pc.WithdrawerAddressStr
	record += pc.WithdrawalToken1IDStr
	record += pc.WithdrawalToken2IDStr
	record += strconv.FormatUint(pc.WithdrawalFeeAmt, 10)
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (pc *PDEFeeWithdrawalRequest) CalculateSize() uint64 {
	return calculateSize(pc)
}

func NewPDEFeeWithdrawalRequest(
	withdrawerAddressStr string,
	withdrawalToken1IDStr string,
	withdrawalToken2IDStr string,
	withdrawalFeeAmt uint64,
	metaType int,
) (*PDEFeeWithdrawalRequest, error) {
	metadataBase := NewMetadataBaseWithSignature(metaType)
	pdeFeeWithdrawalRequest := &PDEFeeWithdrawalRequest{
		WithdrawerAddressStr:  withdrawerAddressStr,
		WithdrawalToken1IDStr: withdrawalToken1IDStr,
		WithdrawalToken2IDStr: withdrawalToken2IDStr,
		WithdrawalFeeAmt:      withdrawalFeeAmt,
	}
	pdeFeeWithdrawalRequest.MetadataBaseWithSignature = *metadataBase
	return pdeFeeWithdrawalRequest, nil
}
