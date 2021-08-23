package rpc

// GetBurnProof retrieves the burning proof of a transaction.
func (server *RPCServer) GetBurnProof(txHash string, isBSC ...bool) ([]byte, error) {
	method := getBurnProof
	if len(isBSC) > 0 && isBSC[0] {
		method = getBSCBurnProof
	}
	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(method, params)
}

// GetBurnProofForSC retrieves the burning proof of a transaction for depositing to smart contracts.
func (server *RPCServer) GetBurnProofForSC(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(getBurnProofForDepositToSC, params)
}

// CheckShieldStatus checks the status of a decentralized shielding transaction.
func (server *RPCServer) CheckShieldStatus(txHash string) ([]byte, error) {
	tmpParams := make(map[string]interface{})
	tmpParams["TxReqID"] = txHash

	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	return server.SendQuery(getBridgeReqWithStatus, params)
}

// GetAllBridgeTokens retrieves the list of bridge tokens in the network.
func (server *RPCServer) GetAllBridgeTokens() ([]byte, error) {
	return server.SendQuery(getAllBridgeTokens, nil)
}
