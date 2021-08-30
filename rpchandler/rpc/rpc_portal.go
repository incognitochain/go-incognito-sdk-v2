package rpc

// GetPortalUnShieldingRequestStatus retrieves the status of a port un-shielding request.
func (server *RPCServer) GetPortalUnShieldingRequestStatus(unShieldID string) ([]byte, error) {
	method := getPortalUnShieldingRequestStatus
	params := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["UnshieldID"] = unShieldID
	params = append(params, mapParams)
	return server.SendQuery(method, params)
}
