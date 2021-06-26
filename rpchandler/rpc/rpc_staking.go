package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

func (server *RPCServer) GetListRewardAmount() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := listRewardAmount
	params := make([]interface{}, 0)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetRewardAmount(paymentAddress string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getRewardAmount
	params := make([]interface{}, 0)
	params = append(params, paymentAddress)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetMiningInfo() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getMiningInfo
	params := make([]interface{}, 0)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetSyncStats() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getSyncStats
	params := make([]interface{}, 0)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}
