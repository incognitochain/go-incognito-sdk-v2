package rpc

// GetPortalShieldingRequestStatus retrieves the status of a port shielding request.
func (server *RPCServer) GetPortalShieldingRequestStatus(shieldID string) ([]byte, error) {
	method := getPortalShieldingRequestStatus
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["ReqTxID"] = shieldID
	params = append(params, mapParams)
	return server.SendQuery(method, params)
}

// GenerateShieldingMultiSigAddress calls the remote node to generate the depositing address for a payment address w.r.t to a tokenID.
func (server *RPCServer) GenerateShieldingMultiSigAddress(paymentAddress, tokenID string) ([]byte, error) {
	method := generatePortalShieldMultisigAddress
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["IncAddressStr"] = paymentAddress
	mapParams["TokenID"] = tokenID
	params = append(params, mapParams)
	return server.SendQuery(method, params)
}

// GetPortalUnShieldingRequestStatus retrieves the status of a portal un-shielding request.
func (server *RPCServer) GetPortalUnShieldingRequestStatus(unShieldID string) ([]byte, error) {
	method := getPortalUnShieldingRequestStatus
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["UnshieldID"] = unShieldID
	params = append(params, mapParams)
	return server.SendQuery(method, params)
}
