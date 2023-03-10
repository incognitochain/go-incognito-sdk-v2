package rpc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

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

func (server *RPCServer) GetBeaconStaker(beaconHeight uint64, beaconStakerPublicKey string) (interface{}, error) {
	params := make([]interface{}, 0)
	params = append(params, beaconHeight)
	params = append(params, beaconStakerPublicKey)

	return server.SendQuery(getBeaconStaker, params)
}

func (server *RPCServer) GetShardStaker(beaconHeight uint64, shardStakerPublicKey string) (interface{}, error) {
	params := make([]interface{}, 0)
	params = append(params, beaconHeight)
	params = append(params, shardStakerPublicKey)

	return server.SendQuery(getShardStaker, params)
}

func (server *RPCServer) GetBeaconCommitteeState(beaconHeight uint64) (interface{}, error) {
	params := make([]interface{}, 0)
	params = append(params, beaconHeight)

	return server.SendQuery(getBeaconCommitteeState, params)
}

func (server *RPCServer) GetBeaconCandidateUID(beaconStakerPublicKey string) (interface{}, error) {
	params := make([]interface{}, 0)
	params = append(params, beaconStakerPublicKey)

	res, err := server.SendQuery(getBeaconCandidateUID, params)
	if err != nil {
		return nil, err
	}
	responseMap := map[string]interface{}{}
	err = json.Unmarshal(res, &responseMap)
	if err != nil {
		return nil, err
	}
	result, ok := responseMap["Result"]
	if !ok {
		return "", errors.Errorf("Can not get beacon candidate UID, response %+v", responseMap)
	}
	if beaconUID, ok := result.(string); ok {
		return beaconUID, nil
	}
	return "", errors.Errorf("Can not get beacon candidate UID, response %+v", responseMap)
}
