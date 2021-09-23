package jsonresult

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
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
type CurrentPdexState struct {
	WaitingContributions        map[string]Pdexv3Contribution
	DeletedWaitingContributions map[string]Pdexv3Contribution
	PoolPairs                   map[string]*Pdexv3PoolPairState
	Params                      *Pdexv3Params
	StakingPoolStates           map[string]*Pdexv3StakingPoolState // tokenID -> StakingPoolState
	NftIDs                      map[string]uint64
}

type Pdexv3Contribution struct {
	PoolPairID  string
	OtaReceiver string
	TokenID     common.Hash
	Amount      uint64
	Amplifier   uint
	TxReqID     common.Hash
	NftID       common.Hash
	ShardID     byte
}

type Pdexv3PoolPairState struct {
	State           Pdexv3PoolPair
	Shares          map[string]*Pdexv3Share
	Orderbook       Pdexv3Orderbook
	LpFeesPerShare  map[common.Hash]*big.Int
	ProtocolFees    map[common.Hash]uint64
	StakingPoolFees map[common.Hash]uint64
}

type Pdexv3PoolPair struct {
	ShareAmount         uint64
	Token0ID            common.Hash
	Token1ID            common.Hash
	Token0RealAmount    uint64
	Token1RealAmount    uint64
	Token0VirtualAmount *big.Int
	Token1VirtualAmount *big.Int
	Amplifier           uint
}

type Pdexv3Share struct {
	Amount             uint64
	TradingFees        map[common.Hash]uint64
	LastLPFeesPerShare map[common.Hash]*big.Int
}

type Pdexv3Orderbook struct {
	Orders []*Pdexv3Order `json:"orders"`
}

type Pdexv3Order struct {
	Id             string
	NftID          common.Hash
	Token0Rate     uint64
	Token1Rate     uint64
	Token0Balance  uint64
	Token1Balance  uint64
	TradeDirection byte
	Fee            uint64
}

type Pdexv3Params struct {
	DefaultFeeRateBPS               uint            // the default value if fee rate is not specific in FeeRateBPS (default 0.3% ~ 30 BPS)
	FeeRateBPS                      map[string]uint // map: pool ID -> fee rate (0.1% ~ 10 BPS)
	PRVDiscountPercent              uint            // percent of fee that will be discounted if using PRV as the trading token fee (default: 25%)
	TradingProtocolFeePercent       uint            // percent of fees that is rewarded for the core team (default: 0%)
	TradingStakingPoolRewardPercent uint            // percent of fees that is distributed for staking pools (PRV, PDEX, ..., default: 10%)
	PDEXRewardPoolPairsShare        map[string]uint // map: pool pair ID -> PDEX reward share weight
	StakingPoolsShare               map[string]uint // map: staking tokenID -> pool staking share weight
	StakingRewardTokens             []common.Hash   // list of staking reward tokens
	MintNftRequireAmount            uint64          // amount prv for depositing to pdex
	MaxOrdersPerNft                 uint            // max orders per nft
}

type Pdexv3StakingPoolState struct {
	Liquidity       uint64
	Stakers         map[string]*Pdexv3Staker // nft -> amount staking
	RewardsPerShare map[common.Hash]*big.Int
}


type Pdexv3Staker struct {
	Liquidity           uint64
	Rewards             map[common.Hash]uint64
	LastRewardsPerShare map[common.Hash]*big.Int
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
