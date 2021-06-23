package rpc

// GetBlockchainInfo returns the current state of the Incognito network.
func (server *RPCServer) GetBlockchainInfo() ([]byte, error) {
	return server.SendQuery(getBlockChainInfo, nil)
}

// GetBestBlock returns the best block numbers (for beacon and shard chains).
func (server *RPCServer) GetBestBlock() ([]byte, error) {
	return server.SendQuery(getBestBlock, nil)
}

// GetBestBlockHash returns the current best block hashes.
func (server *RPCServer) GetBestBlockHash() ([]byte, error) {
	return server.SendQuery(getBestBlockHash, nil)
}

// RetrieveBlock returns the detail of a block given its hash.
func (server *RPCServer) RetrieveBlock(blockHash string, verbosity string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, blockHash)
	params = append(params, verbosity)

	return server.SendQuery(retrieveBlock, params)
}

// GetShardBestState returns the best state of a shard chain.
func (server *RPCServer) GetShardBestState(shardID byte) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, shardID)

	return server.SendQuery(getShardBestState, params)
}

// GetBeaconBestState returns the best state of the beacon chain.
func (server *RPCServer) GetBeaconBestState() ([]byte, error) {
	return server.SendQuery(getBeaconBestState, nil)
}

// GetRawMemPool returns a list of transactions currently in the pool.
func (server *RPCServer) GetRawMemPool() ([]byte, error) {
	return server.SendQuery(getRawMempool, nil)
}
