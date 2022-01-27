package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"math/big"
	"sort"
	"strings"

	// "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	// "github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

const (
	WaitingContributions  = "WaitingContributions"
	PoolPairs             = "PoolPairs"
	PoolPair              = "PoolPair"
	PoolPairShares        = "PoolPairShares"
	PoolPairOrders        = "PoolPairOrders"
	Params                = "Params"
	StakingPools          = "StakingPools"
	StakingPool           = "StakingPool"
	NftIDs                = "NftIDs"
	All                   = "All"
	SimpleVerbosity       = 1
	IntermediateVerbosity = 2
	FullVerbosity         = 3
)

// Share represents a pDEX contribution share.
type Share struct {
	TokenID1Str string
	TokenID2Str string
	ShareAmount uint64
}

// GetPdexState retrieves the state of pDEX at the provided beacon height.
// If the beacon height is set to 0, it returns the latest pDEX state.
func (client *IncClient) GetPdexState(beaconHeight uint64) (*jsonresult.CurrentPdexState, error) {
	responseInBytes, err := client.rpcServer.GetPdexState(beaconHeight)
	if err != nil {
		return nil, err
	}

	var pdeState jsonresult.CurrentPdexState
	err = rpchandler.ParseResponse(responseInBytes, &pdeState)
	if err != nil {
		return nil, err
	}

	return &pdeState, nil
}

// GetPdexParams retrieves the params of pDEX at the provided beacon height.
// If the beacon height is set to 0, it returns the latest pDEX param.
func (client *IncClient) GetPdexParams(beaconHeight uint64) (*jsonresult.Pdexv3Params, error) {
	filter := make(map[string]interface{})
	filter["Key"] = "Params"
	filter["Verbosity"] = 3
	responseInBytes, err := client.rpcServer.GetPdexState(beaconHeight, filter)
	if err != nil {
		return nil, err
	}

	type Res struct {
		Params map[string]interface{}
	}
	var tmpRes *Res
	err = rpchandler.ParseResponse(responseInBytes, &tmpRes)
	if err != nil {
		return nil, err
	}

	jsb, _ := json.Marshal(tmpRes.Params)
	var param jsonresult.Pdexv3Params
	err = json.Unmarshal(jsb, &param)
	if err != nil {
		return nil, err
	}

	return &param, nil
}

// GetAllPdexPoolPairs retrieves all pools in pDEX at the provided beacon height.
// If the beacon height is set to 0, it returns the latest pDEX pool pairs.
func (client *IncClient) GetAllPdexPoolPairs(beaconHeight uint64) (map[string]*jsonresult.Pdexv3PoolPairState, error) {
	filter := make(map[string]interface{})
	filter["Key"] = PoolPairs
	filter["Verbosity"] = FullVerbosity
	filter["ID"] = ""

	responseInBytes, err := client.rpcServer.GetPdexState(beaconHeight, filter)
	if err != nil {
		return nil, err
	}

	var res jsonresult.CurrentPdexState
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res.PoolPairs, nil
}

// GetPdexPoolPair retrieves the pDEX pool information for pair tokenID1-tokenID2 at the provided beacon height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetPdexPoolPair(beaconHeight uint64, tokenID1, tokenID2 string) (map[string]*jsonresult.Pdexv3PoolPairState, error) {
	allPoolPairs, err := client.GetAllPdexPoolPairs(beaconHeight)
	if err != nil {
		return nil, err
	}

	results := make(map[string]*jsonresult.Pdexv3PoolPairState)
	// search for the pool by concatenating tokenIDs. Try both orderings
	prefix1 := fmt.Sprintf("%s-%s", tokenID1, tokenID2)
	prefix2 := fmt.Sprintf("%s-%s", tokenID2, tokenID1)
	for k, v := range allPoolPairs {
		if strings.HasPrefix(k, prefix1) || strings.HasPrefix(k, prefix2) {
			results[k] = v
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("cannot found pool pair for tokenID %v and %v", tokenID1, tokenID2)
	}
	return results, nil
}

// GetPoolPairStateByID returns the pool pair state of a given poolID at the provided beacon height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetPoolPairStateByID(beaconHeight uint64, poolID string) (*jsonresult.Pdexv3PoolPairState, error) {
	filter := make(map[string]interface{})
	filter["Key"] = PoolPair
	filter["Verbosity"] = FullVerbosity
	filter["ID"] = poolID

	responseInBytes, err := client.rpcServer.GetPdexState(beaconHeight, filter)
	if err != nil {
		return nil, err
	}

	var res jsonresult.CurrentPdexState
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res.PoolPairs[poolID], nil
}

// GetPoolShareAmount returns the share amount of a pDEX nftID with-in a given poolID.
func (client *IncClient) GetPoolShareAmount(poolID, nftID string) (uint64, error) {
	pool, err := client.GetPoolPairStateByID(0, poolID)
	if err != nil {
		return 0, err
	}

	share, ok := pool.Shares[nftID]
	if !ok {
		return 0, fmt.Errorf("share of nftID %v not found for poolID %v", nftID, poolID)
	}

	return share.Amount, nil
}

func calculateBuyAmount(amountIn uint64, virtualReserveIn *big.Int, virtualReserveOut *big.Int) (uint64, error) {
	if amountIn <= 0 {
		return 0, fmt.Errorf("invalid input amount %d", amountIn)
	}
	amount := big.NewInt(0).SetUint64(amountIn)
	num := big.NewInt(0).Mul(amount, virtualReserveOut)
	den := big.NewInt(0).Add(amount, virtualReserveIn)
	result := num.Div(num, den)
	if !result.IsUint64() {
		return 0, fmt.Errorf("buy amount out %s of uint64 range", result.String())
	}
	return result.Uint64(), nil
}

// CheckPrice gets the remote server to check price for trading things.
func (client *IncClient) CheckPrice(pairID, tokenToSell string, sellAmount uint64) (uint64, error) {
	pairs, err := client.GetAllPdexPoolPairs(0)
	if err != nil {
		return 0, err
	}
	pair, exists := pairs[pairID]
	if !exists {
		return 0, fmt.Errorf("No pool found for ID %s", pairID)
	}

	var virtualAmtSell, virtualAmtBuy *big.Int
	switch tokenToSell {
	case pair.State.Token0ID.String():
		virtualAmtSell = big.NewInt(0).Set(pair.State.Token0VirtualAmount)
		virtualAmtBuy = big.NewInt(0).Set(pair.State.Token1VirtualAmount)
	case pair.State.Token1ID.String():
		virtualAmtSell = big.NewInt(0).Set(pair.State.Token1VirtualAmount)
		virtualAmtBuy = big.NewInt(0).Set(pair.State.Token0VirtualAmount)
	default:
		return 0, fmt.Errorf("No tokenID %s in pool %s", tokenToSell, pairID)
	}

	buyAmount, err := calculateBuyAmount(sellAmount, virtualAmtSell, virtualAmtBuy)
	if err != nil {
		return 0, err
	}
	return buyAmount, nil
}

// CheckNFTMintingStatus retrieves the status of a (pDEX) NFT minting transaction.
func (client *IncClient) CheckNFTMintingStatus(txHash string) (*jsonresult.MintNFTStatus, error) {
	responseInBytes, err := client.rpcServer.CheckNFTMintingStatus(txHash)
	if err != nil {
		return nil, err
	}
	type TmpResult struct {
		ID     string `json:"NftID"`
		Status int    `json:"Status"`
	}
	var res jsonresult.MintNFTStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckTradeStatus checks the status of a trading transaction.
func (client *IncClient) CheckTradeStatus(txHash string) (*jsonresult.DEXTradeStatus, error) {
	responseInBytes, err := client.rpcServer.CheckTradeStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXTradeStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckDEXLiquidityContributionStatus checks the status of a liquidity-contributing transaction.
func (client *IncClient) CheckDEXLiquidityContributionStatus(txHash string) (*jsonresult.DEXAddLiquidityStatus, error) {
	responseInBytes, err := client.rpcServer.CheckDEXLiquidityContributionStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXAddLiquidityStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckDEXLiquidityWithdrawalStatus checks the status of a liquidity-withdrawal transaction.
func (client *IncClient) CheckDEXLiquidityWithdrawalStatus(txHash string) (*jsonresult.DEXWithdrawLiquidityStatus, error) {
	responseInBytes, err := client.rpcServer.CheckDEXLiquidityWithdrawalStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXWithdrawLiquidityStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckOrderAddingStatus checks the status of an order-book adding transaction.
func (client *IncClient) CheckOrderAddingStatus(txHash string) (*jsonresult.AddOrderStatus, error) {
	responseInBytes, err := client.rpcServer.CheckAddOrderStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.AddOrderStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckOrderWithdrawalStatus checks the status of an order-book withdrawing transaction.
// There are at most two statuses for an order (in case of partially-filled). If `tokenIDs` is not provided, the function
// will automatically return one of them.
//
// NOTE: only the first value of `tokenIDs` is used.
func (client *IncClient) CheckOrderWithdrawalStatus(txHash string, tokenIDs ...string) (*jsonresult.WithdrawOrderStatus, error) {
	responseInBytes, err := client.rpcServer.CheckOrderWithdrawalStatus(txHash, tokenIDs...)
	if err != nil {
		return nil, err
	}

	var res jsonresult.WithdrawOrderStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckDEXStakingStatus checks the status of a pDEX staking transaction.
func (client *IncClient) CheckDEXStakingStatus(txHash string) (*jsonresult.DEXStakeStatus, error) {
	responseInBytes, err := client.rpcServer.CheckDEXStakingStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXStakeStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckDEXUnStakingStatus checks the status of a pDEX un-staking transaction.
func (client *IncClient) CheckDEXUnStakingStatus(txHash string) (*jsonresult.DEXUnStakeStatus, error) {
	responseInBytes, err := client.rpcServer.CheckDEXUnStakingStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXUnStakeStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckDEXStakingRewardWithdrawalStatus retrieves the status of a pDEX staking-reward withdrawal transaction.
func (client *IncClient) CheckDEXStakingRewardWithdrawalStatus(txHash string) (*jsonresult.DEXWithdrawStakingRewardStatus, error) {
	responseInBytes, err := client.rpcServer.CheckDEXStakingRewardWithdrawalStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXWithdrawStakingRewardStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckDEXLPFeeWithdrawalStatus retrieves the status of a pDEX LP fee withdrawal transaction.
func (client *IncClient) CheckDEXLPFeeWithdrawalStatus(txHash string) (*jsonresult.DEXWithdrawLPFeeStatus, error) {
	responseInBytes, err := client.rpcServer.CheckDEXLPFeeWithdrawalStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXWithdrawLPFeeStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CheckDEXProtocolFeeWithdrawalStatus retrieves the status of a pDEX protocol fee withdrawal transaction.
func (client *IncClient) CheckDEXProtocolFeeWithdrawalStatus(txHash string) (*jsonresult.DEXWithdrawProtocolFeeStatus, error) {
	responseInBytes, err := client.rpcServer.CheckDEXProtocolFeeWithdrawalStatus(txHash)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXWithdrawProtocolFeeStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetEstimatedDEXStakingReward returns the estimated pDEX staking rewards for an nftID with the given staking pool at a specific beacon height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetEstimatedDEXStakingReward(beaconHeight uint64, stakingPoolID, nftID string) (map[string]uint64, error) {
	responseInBytes, err := client.rpcServer.CheckDEXStakingReward(beaconHeight, stakingPoolID, nftID)
	if err != nil {
		return nil, err
	}

	var res map[string]uint64
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetEstimatedLPValue returns the estimated LP value in a pool pairID for a given nftID at a specific beacon height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetEstimatedLPValue(beaconHeight uint64, pairID, nftIDStr string) (*jsonresult.DEXLPValue, error) {
	responseInBytes, err := client.rpcServer.CheckDEXLPValue(beaconHeight, pairID, nftIDStr)
	if err != nil {
		return nil, err
	}

	var res jsonresult.DEXLPValue
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetListNftIDs returns the all pDEX minted nftIDs information till the given beacon block height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetListNftIDs(beaconHeight uint64) (map[string]uint64, error) {
	filter := make(map[string]interface{})
	filter["Key"] = NftIDs
	filter["Verbosity"] = SimpleVerbosity
	filter["ID"] = ""

	responseInBytes, err := client.rpcServer.GetPdexState(beaconHeight, filter)
	if err != nil {
		return nil, err
	}
	type NftResults struct {
		NftIDs map[string]uint64 `json:"NftIDs"`
	}
	var res NftResults
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res.NftIDs, nil
}

// GetDexParams returns the pDEX parameters at given beacon block height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetDexParams(beaconHeight uint64) (*jsonresult.Pdexv3Params, error) {
	filter := make(map[string]interface{})
	filter["Key"] = Params
	filter["Verbosity"] = SimpleVerbosity
	filter["ID"] = ""

	responseInBytes, err := client.rpcServer.GetPdexState(beaconHeight, filter)
	if err != nil {
		return nil, err
	}

	var res jsonresult.CurrentPdexState
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res.Params, nil
}

// GetListStakingPoolShares returns the list of all tokens (and their share) allowed to stake at given beacon block height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetListStakingPoolShares(beaconHeight uint64) (map[string]uint, error) {
	pdeState, err := client.GetPdexState(beaconHeight)
	if err != nil {
		return nil, err
	}

	return pdeState.Params.StakingPoolsShare, nil
}

// GetListStakingRewardTokens returns the list of all available staking reward tokens at given beacon block height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetListStakingRewardTokens(beaconHeight uint64) ([]common.Hash, error) {
	pdeState, err := client.GetPdexState(beaconHeight)
	if err != nil {
		return nil, err
	}

	return pdeState.Params.StakingRewardTokens, nil
}

// GetOrderByID returns the detail of an order given its id. The data is subject to the given beacon block height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetOrderByID(beaconHeight uint64, orderID string) (*jsonresult.Pdexv3Order, error) {
	dexState, err := client.GetPdexState(beaconHeight)
	if err != nil {
		return nil, err
	}

	for id, pair := range dexState.PoolPairs {
		for _, order := range pair.Orderbook.Orders {
			if order.Id == orderID {
				order.PoolID = id
				return order, nil
			}
		}
	}

	return nil, fmt.Errorf("order not found")
}

// GetMinPRVRequiredToMintNFT returns the minimum PRV amount required to mint an NFT.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetMinPRVRequiredToMintNFT(beaconHeight uint64) uint64 {
	params, err := client.GetPdexParams(beaconHeight)
	if err != nil {
		return defaultNftRequiredAmount
	}

	return params.MintNftRequireAmount
}

// BuildDEXShareKey constructs a key for retrieving contributed shares in pDEX.
func BuildDEXShareKey(beaconHeight uint64, token1ID string, token2ID string, contributorAddress string) ([]byte, error) {
	pdeSharePrefix := []byte("pdeshare-")
	prefix := append(pdeSharePrefix, []byte(fmt.Sprintf("%d-", beaconHeight))...)
	tokenIDs := []string{token1ID, token2ID}
	sort.Strings(tokenIDs)

	var keyAddr string
	var err error
	if len(contributorAddress) == 0 {
		keyAddr = contributorAddress
	} else {
		//Always parse the contributor address into the oldest version for compatibility
		keyAddr, err = wallet.GetPaymentAddressV1(contributorAddress, false)
		if err != nil {
			return nil, err
		}
	}
	return append(prefix, []byte(tokenIDs[0]+"-"+tokenIDs[1]+"-"+keyAddr)...), nil
}

// BuildDEXPoolKey constructs a key for a pool in pDEX.
func BuildDEXPoolKey(token1ID string, token2ID string) string {
	tokenIDs := []string{token1ID, token2ID}
	sort.Strings(tokenIDs)

	return fmt.Sprintf("%v-%v", tokenIDs[0], tokenIDs[1])
}
