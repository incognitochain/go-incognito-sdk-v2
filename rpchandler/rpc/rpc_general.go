package rpc

// GetActiveShards retrieves the active shards of the network.
func (server *RPCServer) GetActiveShards() ([]byte, error) {
	return server.SendQuery(getActiveShards, nil)
}

// ConvertPaymentAddress calls the full-node to convert a payment address into the oldest version.
func (server *RPCServer) ConvertPaymentAddress(addr string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, addr)

	return server.SendQuery(convertPaymentAddress, params)
}
