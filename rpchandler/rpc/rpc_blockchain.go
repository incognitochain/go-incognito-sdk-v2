package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

func (server *RPCServer) GetBlockchainInfo() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	query := `{
		"jsonrpc":"1.0",
		"method":"getblockchaininfo",
		"params": "",
		"id":1
	}`
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) GetBestBlock() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	if len(server.GetURL()) == 0 {
		return []byte{}, errors.New("Server has not set mainnet or testnet")
	}
	query := `{
		"jsonrpc":"1.0",
		"method":"getbestblock",
		"params": "",
		"id":1
	}`
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) GetBestBlockHash() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	query := `{
		"jsonrpc":"1.0",
		"method":"getbestblockhash",
		"params": "",
		"id":1
	}`
	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) RetrieveBlock(blockHash string, verbosity string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := retrieveBlock

	params := make([]interface{}, 0)
	params = append(params, blockHash)
	params = append(params, verbosity)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetShardBestState(shardID byte) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getShardBestState
	params := make([]interface{}, 0)
	params = append(params, shardID)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetBeaconBestState() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getBeaconBestState
	params := make([]interface{}, 0)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetRawMempool() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}
	query := `{
		"jsonrpc": "1.0",
		"method": "getmempoolinfo",
		"params": "",
		"id": 1
	}`
	return server.SendPostRequestWithQuery(query)
}
