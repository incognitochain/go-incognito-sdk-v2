package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

func (server *RPCServer) GetActiveShards() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getActiveShards

	params := make([]interface{}, 0)
	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) ConvertPaymentAddress(addr string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := "convertpaymentaddress"
	params := make([]interface{}, 0)
	params = append(params, addr)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}
