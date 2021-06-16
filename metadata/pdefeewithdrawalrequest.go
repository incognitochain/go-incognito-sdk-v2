package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

// PDEFeeWithdrawalRequest - privacy dex withdrawal request
type PDEFeeWithdrawalRequest struct {
	WithdrawerAddressStr  string
	WithdrawalToken1IDStr string
	WithdrawalToken2IDStr string
	WithdrawalFeeAmt      uint64
	MetadataBaseWithSignature
}

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

func (pc *PDEFeeWithdrawalRequest) CalculateSize() uint64 {
	return calculateSize(pc)
}
