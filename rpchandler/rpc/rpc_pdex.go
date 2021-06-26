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
	mapParam := make(map[string]interface{})
	mapParam["TxRequestIDStr"] = txHash

	params := make([]interface{}, 0)
	params = append(params, mapParam)

	return server.SendQuery(getPDETradeStatus, params)
}

// GetPDEState retrieves the pDEX state at the given beacon height.
func (server *RPCServer) GetPDEState(beaconHeight uint64) ([]byte, error) {
	mapParams := make(map[string]interface{})
	mapParams["BeaconHeight"] = beaconHeight

	params := make([]interface{}, 0)
	params = append(params, mapParams)

	return server.SendQuery(getPDEState, params)
}

// ConvertPDEPrice gets the pDEX to check the price between to tokens.
func (server *RPCServer) ConvertPDEPrice(tokenToSell, tokenToBuy string, amount uint64) ([]byte, error) {
	mapParam := make(map[string]interface{})
	mapParam["FromTokenIDStr"] = tokenToSell
	mapParam["ToTokenIDStr"] = tokenToBuy
	mapParam["Amount"] = amount

	params := make([]interface{}, 0)
	params = append(params, mapParam)

	return server.SendQuery(convertPDEPrices, params)
}
