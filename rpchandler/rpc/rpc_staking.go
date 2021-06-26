package rpc

// GetListRewardAmount returns the current reward amounts on the network.
func (server *RPCServer) GetListRewardAmount() ([]byte, error) {
	return server.SendQuery(listRewardAmount, nil)
}

// GetRewardAmount gets the reward amounts of a user.
func (server *RPCServer) GetRewardAmount(paymentAddress string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, paymentAddress)

	return server.SendQuery(getRewardAmount, params)
}

// GetMiningInfo retrieves the mining status of a remote (validator) node.
//
// This RPC should call to the (staked) node, instead of a full-node.
func (server *RPCServer) GetMiningInfo() ([]byte, error) {
	return server.SendQuery(getMiningInfo, nil)
}

// GetSyncStats retrieves the sync statistics of a remote (validator) node.
//
// This RPC should call to the (staked) node, instead of a full-node.
func (server *RPCServer) GetSyncStats() ([]byte, error) {
	return server.SendQuery(getSyncStats, nil)
}
