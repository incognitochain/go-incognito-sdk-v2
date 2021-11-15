package jsonresult

import (
	"math/big"

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

// DEXTradeStatus represents the status of a pDEX v3 trade.
type DEXTradeStatus struct {
	// Status represents the status of the trade, and should be understood as follows:
	// 	- 0: the trade request is refunded;
	//	- 1: the trade request is accepted.
	Status int `json:"Status"`

	// BuyAmount is the receiving amount of the trade (in case of failure, it equals to 0).
	BuyAmount uint64 `json:"BuyAmount"`

	// TokenToBuy is the buying tokenId.
	TokenToBuy string `json:"TokenToBuy"`
}

// DEXAddLiquidityStatus represents the status of a pDEX v3 liquidity contribution.
type DEXAddLiquidityStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 1: the contribution is in the waiting pool;
	//	- 2: the contribution is fully accepted;
	//	- 3: the contribution failed and is refunded;
	//	- 4: the contribution is partially accepted.
	Status int `json:"Status"`

	// Token0ID is the ID of the first token.
	Token0ID string `json:"Token0ID"`

	// Token0ContributedAmount is the contributed amount of the first tokenID.
	Token0ContributedAmount uint64 `json:"Token0ContributedAmount"`

	// Token0ReturnedAmount is the returned amount (in case of over-amount) of the first tokenID.
	Token0ReturnedAmount uint64 `json:"Token0ReturnedAmount"`

	// Token1ID is the ID of the second token.
	Token1ID string `json:"Token1ID"`

	// Token1ContributedAmount is the contributed amount of the second tokenID.
	Token1ContributedAmount uint64 `json:"Token1ContributedAmount"`

	// Token1ReturnedAmount is the returned amount (in case of over-amount) of the second tokenID.
	Token1ReturnedAmount uint64 `json:"Token1ReturnedAmount"`

	// PoolPairID is the pool pair ID of the contribution.
	PoolPairID string `json:"PoolPairID"`
}

// DEXWithdrawLiquidityStatus represents the status of a pDEX v3 liquidity withdrawal.
type DEXWithdrawLiquidityStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 1: the request is accepted;
	//	- 2: the request is rejected.
	Status int `json:"Status"`

	// Token0ID is the ID of the first token.
	Token0ID string `json:"Token0ID"`

	// Token0Amount is the withdrawn amount of the first tokenID.
	Token0Amount uint64 `json:"Token0Amount"`

	// Token1ID is the ID of the second token.
	Token1ID string `json:"Token1ID"`

	// Token1Amount is the withdrawn amount of the second tokenID.
	Token1Amount uint64 `json:"Token1Amount"`
}

// MintNFTStatus represents the status of a pDEX nft minting transaction.
type MintNFTStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 1: the request is accepted;
	//	- 2: the request is rejected.
	Status int `json:"Status"`

	// BurntAmount is the amount of PRV that was burned to mint this NFT.
	BurntAmount uint64 `json:"BurntAmount"`

	// NftID is the ID of the minted NFT.
	NftID string `json:"NftID"`
}

// AddOrderStatus represents the status of a pDEX OB-adding transaction.
type AddOrderStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 0: the request is rejected;
	//	- 1: the request is accepted.
	Status int `json:"Status"`

	// OrderID is the ID of the requesting order.
	OrderID string `json:"OrderID"`
}

// WithdrawOrderStatus represents the status of a pDEX OB-withdrawing transaction.
type WithdrawOrderStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 0: the request is rejected;
	//	- 1: the request is accepted.
	Status int `json:"Status"`

	// TokenID is the ID of the withdrawn token.
	TokenID string `json:"TokenID"`

	// Amount is the withdrawn amount.
	Amount uint64 `json:"Amount"`
}

// DEXStakeStatus represents the status of a pDEX staking transaction.
type DEXStakeStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 0: the request is rejected;
	//	- 1: the request is accepted.
	Status int `json:"Status"`

	// NftID is the ID of the NFT associated with the action.
	NftID string `json:"NftID"`

	// StakingPoolID is the ID of the pool.
	StakingPoolID string `json:"StakingPoolID"`

	// Liquidity is the staked amount.
	Liquidity uint64 `json:"Liquidity"`
}

// DEXUnStakeStatus represents the status of a pDEX un-staking transaction.
type DEXUnStakeStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 0: the request is rejected;
	//	- 1: the request is accepted.
	Status int `json:"Status"`

	// NftID is the ID of the NFT associated with the action.
	NftID string `json:"NftID"`

	// StakingPoolID is the ID of the pool.
	StakingPoolID string `json:"StakingPoolID"`

	// Liquidity is the un-staked amount.
	Liquidity uint64 `json:"Liquidity"`
}

// DEXWithdrawStakingRewardStatus represents the status of a pDEX staking reward withdrawal transaction.
type DEXWithdrawStakingRewardStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 0: the request is rejected;
	//	- 1: the request is accepted.
	Status int `json:"Status"`

	// Receivers is the receiving information.
	Receivers map[string]struct {
		Address string `json:"Address"`
		Amount  uint64 `json:"Amount"`
	} `json:"Receivers"`
}

// DEXWithdrawLPFeeStatus represents the status of a pDEX LP fee withdrawal transaction.
type DEXWithdrawLPFeeStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 0: the request is rejected;
	//	- 1: the request is accepted.
	Status int `json:"Status"`

	// Receivers is the receiving information.
	Receivers map[string]struct {
		Address string `json:"Address"`
		Amount  uint64 `json:"Amount"`
	} `json:"Receivers"`
}

// DEXWithdrawProtocolFeeStatus represents the status of a pDEX protocol fee withdrawal transaction.
type DEXWithdrawProtocolFeeStatus struct {
	// Status represents the status of the transaction, and should be understood as follows:
	//	- 0: the request is rejected;
	//	- 1: the request is accepted.
	Status int `json:"Status"`

	// Receivers is the receiving information.
	Receivers map[string]struct {
		Address string `json:"Address"`
		Amount  uint64 `json:"Amount"`
	} `json:"Receivers"`
}

// DEXLPValue represents the LP value of an LP.
type DEXLPValue struct {
	// PoolValue represents the contributed liquidity in the pool of the LP.
	PoolValue map[string]uint64 `json:"PoolValue"`

	// TradingFee is the trading fee distributed to the LP.
	TradingFee map[string]uint64 `json:"TradingFee"`
}
