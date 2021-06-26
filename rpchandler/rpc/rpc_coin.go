package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
)

//===================== OUTPUT COINS RPC =====================//
//These RPCs return raw JSON bytes.

// GetListOutputCoinsByRPCV1 retrieves list of output coins of an OutCoinKey and returns the result in raw json bytes.
func (server *RPCServer) GetListOutputCoinsByRPCV1(outCoinKey *OutCoinKey, tokenID string, h uint64) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	query := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"method": "listoutputcoins",
		"params": [
			0,
			999999,
			[
				{
			  "PaymentAddress": "%s",
			  "OTASecretKey": "%s",
			  "ReadonlyKey" : "%s",
			  "StartHeight": %d
				}
			],
		  "%s"
		  ],
		"id": 1
	}`, outCoinKey.paymentAddress, outCoinKey.otaKey, outCoinKey.readonlyKey, h, tokenID)

	return server.SendPostRequestWithQuery(query)
}

// GetListOutputCoinsByRPCV2 retrieves list of output coins of an OutCoinKey and returns the result in raw json bytes.
func (server *RPCServer) GetListOutputCoinsByRPCV2(outCoinKey *OutCoinKey, tokenID string, toHeight uint64) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	query := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"method": "listoutputcoinsfromcache",
		"params": [
			0,
			%v,
			[
				{
			  "PaymentAddress": "%s",
			  "OTASecretKey": "%s",
			  "ReadonlyKey" : "%s",
			  "StartHeight": 0
				}
			],
		  "%s"
		  ],
		"id": 1
	}`, toHeight, outCoinKey.paymentAddress, outCoinKey.otaKey, outCoinKey.readonlyKey, tokenID)

	return server.SendPostRequestWithQuery(query)
}

// ListUnspentOutputCoinsByRPC retrieves list of output coins of an OutCoinKey and returns the result in raw json bytes.
//
// NOTE: PrivateKey must be supplied.
func (server *RPCServer) ListUnspentOutputCoinsByRPC(privKeyStr string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	query := fmt.Sprintf(`{
	   "jsonrpc":"1.0",
	   "method":"listunspentoutputcoins",
	   "params":[
		  0,
		  999999,
		  [
			 {
				"PrivateKey":"%s",
				"StartHeight": 0
			 }

		  ]
	   ],
	   "id":1
	}`, privKeyStr)

	return server.SendPostRequestWithQuery(query)
}

// ListPrivacyCustomTokenByRPC lists all tokens currently present on the blockchain.
func (server *RPCServer) ListPrivacyCustomTokenByRPC() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	query := `{
		"id": 1,
		"jsonrpc": "1.0",
		"method": "listprivacycustomtoken",
		"params": []
	}`
	return server.SendPostRequestWithQuery(query)
}

// ListBridgeTokenByRPC lists all bridge-tokens currently present on the blockchain.
func (server *RPCServer) ListBridgeTokenByRPC() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getAllBridgeTokens

	params := make([]interface{}, 0)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

// GetTokenByRPC retrieves all the token's information on the blockchain.
func (server *RPCServer) GetTokenByRPC(tokenID string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getPrivacyCustomToken

	params := make([]interface{}, 0)
	params = append(params, tokenID)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

// HasSerialNumberByRPC checks if the provided serial numbers have been spent or not.
//
// Returned result in raw json bytes.
func (server *RPCServer) HasSerialNumberByRPC(shardID byte, tokenID string, snList []string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	if len(snList) == 0 {
		return nil, errors.New("no serial number provided to be checked")
	}
	snQueryList := make([]string, 0)
	for _, sn := range snList {
		snQueryList = append(snQueryList, fmt.Sprintf(`"%s"`, sn))
	}

	addr := rpchandler.CreatePaymentAddress(shardID) // use a random payment address for anonymity

	method := hasSerialNumbers

	params := make([]interface{}, 0)
	params = append(params, addr)
	params = append(params, snList)
	params = append(params, tokenID)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))

}

func (server *RPCServer) HasSerialNumberInMemPool(snList []string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	if len(snList) == 0 {
		return nil, errors.New("no serial number provided to be checked")
	}
	snQueryList := make([]string, 0)
	for _, sn := range snList {
		snQueryList = append(snQueryList, fmt.Sprintf(`"%s"`, sn))
	}

	method := hasSerialNumbersInMempool

	params := make([]interface{}, 0)
	params = append(params, snList)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetBalanceByPrivatekey(privKeyStr string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	query := fmt.Sprintf(`{
	   "jsonrpc":"1.0",
	   "method":"getbalancebyprivatekey",
	   "params":["%s"],
	   "id":1
	}`, privKeyStr)

	return server.SendPostRequestWithQuery(query)
}

func (server *RPCServer) SubmitKey(otaStr string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := submitKey
	params := make([]interface{}, 0)
	params = append(params, otaStr)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

// AuthorizedSubmitKey submits an OTA Key in an authorized manner for more privileges.
func (server *RPCServer) AuthorizedSubmitKey(otaStr string, accessToken string, fromHeight uint64, isReset bool) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := authorizedSubmitKey
	params := make([]interface{}, 0)
	params = append(params, otaStr)
	params = append(params, accessToken)
	params = append(params, fromHeight)
	params = append(params, isReset)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) RandomCommitments(shardID byte, inputCoins []jsonresult.OutCoin, tokenID string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	addr := rpchandler.CreatePaymentAddress(shardID)

	method := randomCommitments

	params := make([]interface{}, 0)
	params = append(params, addr)
	params = append(params, inputCoins)
	params = append(params, tokenID)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) RandomCommitmentsAndPublicKeys(shardID byte, tokenID string, lenDecoy int) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := randomCommitmentsAndPublicKeys

	params := make([]interface{}, 0)
	params = append(params, shardID)
	params = append(params, lenDecoy)
	params = append(params, tokenID)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

//===================== END OF OUTPUT COINS RPC =====================//
