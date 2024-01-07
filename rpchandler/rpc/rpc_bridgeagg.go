package rpc

// // GetBridgeAggState get bridge aggregator state
// func (server *RPCServer) GetBridgeAggState(beaconHeight uint64) ([]byte, error) {
// 	tmpParams := make(map[string]interface{})
// 	tmpParams["BeaconHeight"] = beaconHeight

// 	params := make([]interface{}, 0)
// 	params = append(params, tmpParams)
// 	return server.SendQuery(bridgeaggState, params)
// }

// // CheckUnshieldUnifiedStatus checks the status of a decentralized unshielding transaction.
// func (server *RPCServer) CheckUnshieldUnifiedStatus(txHash string) ([]byte, error) {

// 	params := make([]interface{}, 0)
// 	params = append(params, txHash)
// 	return server.SendQuery(bridgeaggStatusUnshield, params)
// }

// // CheckShieldUnifiedStatus checks the status of a decentralized shielding transaction.
// func (server *RPCServer) CheckShieldUnifiedStatus(txHash string) ([]byte, error) {

// 	params := make([]interface{}, 0)
// 	params = append(params, txHash)
// 	return server.SendQuery(bridgeaggStatusShield, params)
// }

// // CheckConvertStatuspUnifiedStatus checks the status of a decentralized convert transaction.
// func (server *RPCServer) CheckConvertStatuspUnifiedStatus(txHash string) ([]byte, error) {

// 	params := make([]interface{}, 0)
// 	params = append(params, txHash)
// 	return server.SendQuery(bridgeaggStatusConvert, params)
// }

// func (server *RPCServer) GetBridgeAggEstimateFeeByExpectedAmount(pUnifiedTokenID, tokenID string, expectedAmount uint64) ([]byte, error) {
// 	tmpParams := make(map[string]interface{})
// 	tmpParams["UnifiedTokenID"] = pUnifiedTokenID
// 	tmpParams["TokenID"] = tokenID
// 	tmpParams["ExpectedAmount"] = expectedAmount

// 	params := make([]interface{}, 0)
// 	params = append(params, tmpParams)
// 	return server.SendQuery(bridgeaggEstimateFeeByExpectedAmount, params)
// }

// func (server *RPCServer) GetBridgeAggEstimateFeeByBurntAmount(pUnifiedTokenID, tokenID string, burnAmount uint64) ([]byte, error) {
// 	tmpParams := make(map[string]interface{})
// 	tmpParams["UnifiedTokenID"] = pUnifiedTokenID
// 	tmpParams["TokenID"] = tokenID
// 	tmpParams["BurntAmount"] = burnAmount

// 	params := make([]interface{}, 0)
// 	params = append(params, tmpParams)
// 	return server.SendQuery(bridgeaggEstimateFeeByBurntAmount, params)
// }

// func (server *RPCServer) GetBridgeAggEstimateReward(pUnifiedTokenID, tokenID string, amount uint64) ([]byte, error) {
// 	tmpParams := make(map[string]interface{})
// 	tmpParams["UnifiedTokenID"] = pUnifiedTokenID
// 	tmpParams["TokenID"] = tokenID
// 	tmpParams["Amount"] = amount

// 	params := make([]interface{}, 0)
// 	params = append(params, tmpParams)
// 	return server.SendQuery(bridgeaggEstimateReward, params)
// }

// func (server *RPCServer) GetBridgeAggGetBurnProof(txReqID string, dataIndex ...int) ([]byte, error) {
// 	tmpParams := make(map[string]interface{})
// 	tmpParams["TxReqID"] = txReqID
// 	if len(dataIndex) > 0 {
// 		tmpParams["DataIndex"] = dataIndex[0]
// 	}
// 	params := make([]interface{}, 0)
// 	params = append(params, tmpParams)
// 	return server.SendQuery(bridgeaggGetBurnProof, params)
// }
