package rpc

// ConvertedPrice represents a price conversion between two tokenIDs.
type ConvertedPrice struct {
	FromTokenIDStr string
	ToTokenIDStr   string
	Amount         uint64
	Price          uint64
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

// CheckDEXStakingRewardWithdrawalStatus retrieves the status of a pDEX staking-reward withdrawal transaction.
func (server *RPCServer) CheckDEXStakingRewardWithdrawalStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["ReqTxID"] = txHash
	params = append(params, mapParams)

	return server.SendQuery(getPdexv3WithdrawalStakingRewardStatus, params)
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
