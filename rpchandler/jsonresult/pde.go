package jsonresult

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"sort"
)

var (
	WaitingPDEContributionPrefix = []byte("waitingpdecontribution-")
	PDEPoolPrefix                = []byte("pdepool-")
	PDESharePrefix               = []byte("pdeshare-")
	PDETradingFeePrefix          = []byte("pdetradingfee-")
	PDETradeFeePrefix            = []byte("pdetradefee-")
	PDEContributionStatusPrefix  = []byte("pdecontributionstatus-")
	PDETradeStatusPrefix         = []byte("pdetradestatus-")
	PDEWithdrawalStatusPrefix    = []byte("pdewithdrawalstatus-")
	PDEFeeWithdrawalStatusPrefix = []byte("pdefeewithdrawalstatus-")
)

// PoolInfo represents a pDEX pool of two tokenIDs.
type PoolInfo struct {
	Token1IDStr     string
	Token1PoolValue uint64
	Token2IDStr     string
	Token2PoolValue uint64
}

// CurrentPDEState describes the state of the pDEX at a specific beacon height.
type CurrentPDEState struct {
	WaitingPDEContributions map[string]*PDEContribution `json:"WaitingPDEContributions"`
	PDEPoolPairs            map[string]*PoolInfo        `json:"PDEPoolPairs"`
	PDEShares               map[string]uint64           `json:"PDEShares"`
	PDETradingFees          map[string]uint64           `json:"PDETradingFees"`
	BeaconTimeStamp         int64                       `json:"BeaconTimeStamp"`
}

// PDEContribution describes a contribution on the pDEX.
type PDEContribution struct {
	ContributorAddressStr string
	TokenIDStr            string
	Amount                uint64
	TxReqID               common.Hash
}

// BuildPDEPoolForPairKey builds a pDEX pool-key for a pair of tokenIDs at the given beacon height.
func BuildPDEPoolForPairKey(
	beaconHeight uint64,
	token1IDStr string,
	token2IDStr string,
) []byte {
	beaconHeightBytes := []byte(fmt.Sprintf("%d-", beaconHeight))
	pdePoolForPairByBCHeightPrefix := append(PDEPoolPrefix, beaconHeightBytes...)
	tokenIDStrings := []string{token1IDStr, token2IDStr}
	sort.Strings(tokenIDStrings)
	return append(pdePoolForPairByBCHeightPrefix, []byte(tokenIDStrings[0]+"-"+tokenIDStrings[1])...)
}
