package rpc

// GetCommitteeState retrieves the committee state at the given beacon height and beacon root hash.
// This RPC is mainly used for debugging purposes.
func (server *RPCServer) GetCommitteeState(beaconHeight uint64, beaconRootHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, beaconHeight)
	params = append(params, beaconRootHash)

	return server.SendQuery(getCommitteeState, params)
}

// GetCommitteeStateByShardID retrieves the committee state of a shard given a root hash.
// This RPC is mainly used for debugging purposes.
func (server *RPCServer) GetCommitteeStateByShardID(shardID int, shardRootHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, shardID)
	params = append(params, shardRootHash)

	return server.SendQuery(getCommitteeStateByShard, params)
}
