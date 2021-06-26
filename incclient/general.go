package incclient

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

// SubmitKey submits an OTAKey to the full-node.
func (client *IncClient) SubmitKey(otaKey string) error {
	_, err := client.rpcServer.SubmitKey(otaKey)
	return err
}

// AuthorizedSubmitKey handles submitting OTA keys in an authorized manner.
func (client *IncClient) AuthorizedSubmitKey(otaKey string, accessToken string, fromHeight uint64, isReset bool) error {
	responseInBytes, err := client.rpcServer.AuthorizedSubmitKey(otaKey, accessToken, fromHeight, isReset)
	if err != nil {
		return err
	}

	err = rpchandler.ParseResponse(responseInBytes, nil)
	if err != nil {
		return err
	}
	return nil
}

// NewRPCCall creates and sends a new RPC request based on the given method and parameters to the RPC server.
//
// Example call: NewRPCCall("1.0", "getbeaconbeststate", nil, 1)
func (client *IncClient) NewRPCCall(jsonRPC, method string, params []interface{}, id interface{}) ([]byte, error) {
	if jsonRPC == "" {
		jsonRPC = "1.0"
	}

	request := rpchandler.CreateJsonRequest(jsonRPC, method, params, id)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return client.rpcServer.SendPostRequestWithQuery(string(query))
}
