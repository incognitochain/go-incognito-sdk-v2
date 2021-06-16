package rpchandler

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// RPCError represents an error that is used as a part of a JSON-RPC JsonResponse
// object.
type RPCError struct {
	Code       int    `json:"Code,omitempty"`
	Message    string `json:"Message,omitempty"`
	StackTrace string `json:"StackTrace"`

	err 	   error  `json:"Err"`
}

type JsonRequest struct {
	Jsonrpc string      `json:"Jsonrpc"`
	Method  string      `json:"Method"`
	Params  interface{} `json:"Params"`
	Id      interface{} `json:"Id"`
}

type JsonResponse struct {
	Id      *interface{}         `json:"Id"`
	Result  json.RawMessage      `json:"Result"`
	Error   *RPCError 			 `json:"Error"`
	Params  interface{}          `json:"Params"`
	Method  string               `json:"Method"`
	Jsonrpc string               `json:"Jsonrpc"`
}

var Server = new(RPCServer)
var EthServer = new(RPCServer)

func EncodeBase58Check(data []byte) string {
	b := base58.Base58Check{}.Encode(data, 0)
	return b
}

func DecodeBase58Check(s string) ([]byte, error) {
	b, _, err := base58.Base58Check{}.Decode(s)
	return b, err
}

/*Common functions*/
// RandIntInterval returns a random int in range [L; R]
func RandIntInterval(L, R int) int {
	length := R - L + 1
	r := common.RandInt() % length
	return L + r
}

func ParseResponse(respondInBytes []byte) (*JsonResponse, error) {
	var respond JsonResponse
	err := json.Unmarshal(respondInBytes, &respond)
	if err != nil {
		return nil, err
	}

	if respond.Error != nil{
		return nil, fmt.Errorf("RPC returns an error: %v", respond.Error)
	}

	return &respond, nil
}

func NewParseResponse(respondInBytes []byte, val interface{}) error {
	var respond JsonResponse
	err := json.Unmarshal(respondInBytes, &respond)
	if err != nil {
		return err
	}

	if respond.Error != nil{
		return fmt.Errorf("RPC returns an error: %v", respond.Error)
	}

	err = json.Unmarshal(respond.Result, &val)
	if err != nil {
		return err
	}

	return nil
}

func CreateJsonRequest(jsonRPC, method string, params []interface{}, id interface{}) *JsonRequest{
	request := new(JsonRequest)
	request.Jsonrpc = jsonRPC
	request.Method = method
	request.Id = id
	request.Params = params

	return request
}

//Temp function that creates a payment address of a specific shard.
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