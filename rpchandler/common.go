package rpchandler

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// RPCError represents an error that is used as a part of a JSON-RPC JsonResponse
// object.
type RPCError struct {
	Code       int    `json:"Code,omitempty"`
	Message    string `json:"Message,omitempty"`
	StackTrace string `json:"StackTrace"`
	Err        error  `json:"Err"`
}

// JsonRequest represents a JSON-RPC request.
type JsonRequest struct {
	JsonRPC string      `json:"Jsonrpc"`
	Method  string      `json:"Method"`
	Params  interface{} `json:"Params"`
	Id      interface{} `json:"Id"`
}

// JsonResponse represents a JSON-RPC response.
type JsonResponse struct {
	Id      *interface{}    `json:"Id"`
	Result  json.RawMessage `json:"Result"`
	Error   *RPCError       `json:"Error"`
	Params  interface{}     `json:"Params"`
	Method  string          `json:"Method"`
	JsonRPC string          `json:"Jsonrpc"`
}

// OldParseResponse parses a raw JSON-RPC response into a JsonResponse.
//
// Deprecated: use ParseResponse instead.
func OldParseResponse(respondInBytes []byte) (*JsonResponse, error) {
	var respond JsonResponse
	err := json.Unmarshal(respondInBytes, &respond)
	if err != nil {
		return nil, err
	}

	if respond.Error != nil {
		return nil, fmt.Errorf("RPC returns an error: %v", respond.Error)
	}

	return &respond, nil
}

// ParseResponse parses a JSON-RPC response to val.
func ParseResponse(respondInBytes []byte, val interface{}) error {
	var respond JsonResponse
	err := json.Unmarshal(respondInBytes, &respond)
	if err != nil {
		return err
	}

	if respond.Error != nil {
		return fmt.Errorf("RPC returns an error: %v", respond.Error)
	}

	if val == nil {
		return nil
	}

	err = json.Unmarshal(respond.Result, val)
	if err != nil {
		return err
	}

	return nil
}

// CreateJsonRequest creates a new JsonRequest given the method and parameters.
func CreateJsonRequest(jsonRPC, method string, params []interface{}, id interface{}) *JsonRequest {
	request := new(JsonRequest)
	request.JsonRPC = jsonRPC
	request.Method = method
	request.Id = id
	request.Params = params

	return request
}

// CreatePaymentAddress is a temp function that creates a payment address of a specific shard.
func CreatePaymentAddress(shardID byte) string {
	pk := common.RandBytes(31)
	tk := common.RandBytes(32)

	//Set last byte of pk to be the shardID
	pk = append(pk, shardID)

	addr := key.PaymentAddress{Pk: pk, Tk: tk, OTAPublic: nil}

	keyWallet := new(wallet.KeyWallet)
	keyWallet.KeySet = key.KeySet{PaymentAddress: addr}

	return keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)
}
