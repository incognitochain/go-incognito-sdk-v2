package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

func (server *RPCServer) GetCommitteeState(beaconHeight uint64, beaconRootHash string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getCommitteeState
	params := make([]interface{}, 0)
	params = append(params, beaconHeight)
	params = append(params, beaconRootHash)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetCommitteeStateByShardID(shardID int, shardRootHash string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getCommitteeStateByShard
	params := make([]interface{}, 0)
	params = append(params, shardID)
	params = append(params, shardRootHash)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}
