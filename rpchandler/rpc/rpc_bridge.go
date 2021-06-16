package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

func (server *RPCServer) GetBurnProof(txHash string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getBurnProof

	params := make([]interface{}, 0)
	params = append(params, txHash)
	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetBurnProofForSC(txHash string) ([]byte, error) {
	method := getBurnProofForDepositToSC

	params := make([]interface{}, 0)
	params = append(params, txHash)
	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) CheckShieldStatus(txHash string) ([]byte, error){
	method := getBridgeReqWithStatus

	tmpParams := make(map[string]interface{})
	tmpParams["TxReqID"] = txHash

	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}
