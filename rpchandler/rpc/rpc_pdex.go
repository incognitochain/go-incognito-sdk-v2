package rpc

// ConvertedPrice represents a price conversion between two tokenIDs.
type ConvertedPrice struct {
	FromTokenIDStr string
	ToTokenIDStr   string
	Amount         uint64
	Price          uint64
}

// CheckTradeStatus retrieves the status of a trading transaction.
func (server *RPCServer) CheckTradeStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(pdexv3GetTradeStatus, params)
}

// CheckDEXLiquidityContributionStatus retrieves the status of a liquidity-contributing transaction.
func (server *RPCServer) CheckDEXLiquidityContributionStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(getPdexv3ContributionStatus, params)
}

// CheckDEXLiquidityWithdrawalStatus retrieves the status of a liquidity-withdrawal transaction.
func (server *RPCServer) CheckDEXLiquidityWithdrawalStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(getPdexv3WithdrawLiquidityStatus, params)
}

// CheckAddOrderStatus retrieves the status of an order-adding transaction.
func (server *RPCServer) CheckAddOrderStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(pdexv3GetAddOrderStatus, params)
}

// CheckOrderWithdrawalStatus retrieves the status of an order-canceling transaction.
func (server *RPCServer) CheckOrderWithdrawalStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(pdexv3GetWithdrawOrderStatus, params)
}

// CheckNFTMintingStatus retrieves the status of an NFT minting transaction.
func (server *RPCServer) CheckNFTMintingStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(getPdexv3MintNftStatus, params)
}

// CheckDEXStakingStatus retrieves the status of a pDEX staking transaction.
func (server *RPCServer) CheckDEXStakingStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(pdexv3GetStakingStatus, params)
}

// CheckDEXUnStakingStatus retrieves the status of a pDEX un-staking transaction.
func (server *RPCServer) CheckDEXUnStakingStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(pdexv3GetUnstakingStatus, params)
}

// CheckDEXStakingReward retrieves the estimated amount of staking reward for a nftID.
func (server *RPCServer) CheckDEXStakingReward(beaconHeight uint64, stakingPoolID, nftID string) ([]byte, error) {
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["BeaconHeight"] = beaconHeight
	mapParams["NftID"] = nftID
	mapParams["StakingPoolID"] = stakingPoolID
	params = append(params, mapParams)

	return server.SendQuery(getPdexv3EstimatedStakingReward, params)
}

// CheckDEXStakingRewardWithdrawalStatus retrieves the status of a pDEX staking-reward withdrawal transaction.
func (server *RPCServer) CheckDEXStakingRewardWithdrawalStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["ReqTxID"] = txHash
	params = append(params, mapParams)

	return server.SendQuery(getPdexv3WithdrawalStakingRewardStatus, params)
}

// CheckDEXLPFeeWithdrawalStatus retrieves the status of a pDEX LP fee withdrawal transaction.
func (server *RPCServer) CheckDEXLPFeeWithdrawalStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["ReqTxID"] = txHash
	params = append(params, mapParams)

	return server.SendQuery(getPdexv3WithdrawalLPFeeStatus, params)
}

// CheckDEXProtocolFeeWithdrawalStatus retrieves the status of a pDEX protocol fee withdrawal transaction.
func (server *RPCServer) CheckDEXProtocolFeeWithdrawalStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["ReqTxID"] = txHash
	params = append(params, mapParams)

	return server.SendQuery(getPdexv3WithdrawalProtocolFeeStatus, params)
}

// CheckDEXLPValue retrieves the estimated LP value in a pool pairID for a given nftID.
func (server *RPCServer) CheckDEXLPValue(beaconHeight uint64, pairID, nftID string) ([]byte, error) {
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["BeaconHeight"] = beaconHeight
	mapParams["NftID"] = nftID
	mapParams["PoolPairID"] = pairID
	params = append(params, mapParams)

	return server.SendQuery(getPdexv3EstimatedLPValue, params)
}

// GetPdexState retrieves the pDEX state at the given beacon height.
func (server *RPCServer) GetPdexState(beaconHeight uint64, filters ...map[string]interface{}) ([]byte, error) {
	filter := make(map[string]interface{})
	if len(filters) > 0 {
		filter = filters[0]
	}
	mapParams := make(map[string]interface{})
	mapParams["BeaconHeight"] = beaconHeight
	mapParams["Filter"] = filter

	params := make([]interface{}, 0)
	params = append(params, mapParams)

	return server.SendQuery(getPdexv3State, params)
}

// ConvertPdexPrice gets the pDEX to check the price between to tokens.
func (server *RPCServer) ConvertPdexPrice(tokenToSell, tokenToBuy string, amount uint64) ([]byte, error) {
	mapParam := make(map[string]interface{})
	mapParam["FromTokenIDStr"] = tokenToSell
	mapParam["ToTokenIDStr"] = tokenToBuy
	mapParam["Amount"] = amount

	params := make([]interface{}, 0)
	params = append(params, mapParam)

	return server.SendQuery(convertPDEPrices, params)
}
