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
	// Status represents the status of the trade: 0 - trade not found or trade refunded; 1 - trade accepted.
	Status     int    `json:"Status"`

	// BuyAmount is the receiving amount of the trade (in case of failure, it equals to 0).
	BuyAmount  uint64 `json:"BuyAmount"`

	// TokenToBuy is the buying tokenId.
	TokenToBuy string `json:"TokenToBuy"`
}

// CheckTradeStatus retrieves the status of a trading transaction.
func (server *RPCServer) CheckTradeStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(pdexv3GetTradeStatus, params)
}

// CheckNFTMintingStatus retrieves the status of an NFT minting transaction.
func (server *RPCServer) CheckNFTMintingStatus(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(getPdexv3MintNftStatus, params)
}

// GetPdexState retrieves the pDEX state at the given beacon height.
func (server *RPCServer) GetPdexState(beaconHeight uint64, filter map[string]interface{}) ([]byte, error) {
	mapParams := make(map[string]interface{})
	mapParams["BeaconHeight"] = beaconHeight
	if filter == nil {
		filter = make(map[string]interface{})
	}
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
