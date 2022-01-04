package rpc

import (
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
)

// GetListOutputCoinsByRPCV1 retrieves list of output coins of an OutCoinKey and returns the result in raw json bytes.
func (server *RPCServer) GetListOutputCoinsByRPCV1(outCoinKey *OutCoinKey, tokenID string, h uint64) ([]byte, error) {
	keyParams := make(map[string]interface{})
	keyParams["PaymentAddress"] = outCoinKey.paymentAddress
	keyParams["OTASecretKey"] = outCoinKey.otaKey
	keyParams["ReadonlyKey"] = outCoinKey.readonlyKey
	keyParams["StartHeight"] = h

	params := make([]interface{}, 0)
	params = append(params, 0)
	params = append(params, h)
	params = append(params, []interface{}{keyParams})
	params = append(params, tokenID)

	return server.SendQuery(listOutputCoins, params)
}

// GetListOutputCoinsByRPCV2 retrieves list of output coins of an OutCoinKey and returns the result in raw json bytes.
func (server *RPCServer) GetListOutputCoinsByRPCV2(outCoinKey *OutCoinKey, tokenID string, _ uint64) ([]byte, error) {
	keyParams := make(map[string]interface{})
	keyParams["PaymentAddress"] = outCoinKey.paymentAddress
	keyParams["OTASecretKey"] = outCoinKey.otaKey
	keyParams["ReadonlyKey"] = outCoinKey.readonlyKey
	keyParams["StartHeight"] = 0

	params := make([]interface{}, 0)
	params = append(params, 0)
	params = append(params, 999999)
	params = append(params, []interface{}{keyParams})
	params = append(params, tokenID)

	return server.SendQuery(listOutputCoinsFromCache, params)
}

// GetOTACoinsByIndices returns the list of output coins given the indices.
func (server *RPCServer) GetOTACoinsByIndices(shardID byte, tokenID string, idxList []uint64) ([]byte, error) {
	mapParams := make(map[string]interface{})
	mapParams["ShardID"] = shardID
	mapParams["TokenID"] = tokenID
	mapParams["Indices"] = idxList

	return server.SendQuery(getOTACoinsByIndices, []interface{}{mapParams})
}

// GetOTACoinLength returns the number of OTA coins for each shard.
func (server *RPCServer) GetOTACoinLength() ([]byte, error) {
	return server.SendQuery(getOTACoinLength, []interface{}{})
}

// ListUnspentOutputCoinsByRPC retrieves list of output coins of an OutCoinKey and returns the result in raw json bytes.
//
// NOTE: PrivateKey must be supplied and sent to the server.
func (server *RPCServer) ListUnspentOutputCoinsByRPC(privateKey string) ([]byte, error) {
	keyParams := make(map[string]interface{})
	keyParams["PrivateKey"] = privateKey
	keyParams["StartHeight"] = 0

	params := make([]interface{}, 0)
	params = append(params, 0)
	params = append(params, 999999)
	params = append(params, []interface{}{keyParams})

	return server.SendQuery(listOutputCoinsFromCache, params)
}

// ListPrivacyCustomTokenByRPC lists all tokens currently present on the blockchain.
func (server *RPCServer) ListPrivacyCustomTokenByRPC() ([]byte, error) {
	return server.SendQuery(listPrivacyCustomToken, nil)
}

// ListPrivacyCustomTokenIDsByRPC lists all token IDs currently present on the blockchain.
func (server *RPCServer) ListPrivacyCustomTokenIDsByRPC() ([]byte, error) {
	return server.SendQuery(listPrivacyCustomTokenIDs, nil)
}

// ListBridgeTokenByRPC lists all bridge-tokens currently present on the blockchain.
func (server *RPCServer) ListBridgeTokenByRPC() ([]byte, error) {
	return server.SendQuery(getAllBridgeTokens, nil)
}

// GetTokenByRPC retrieves all the token's information on the blockchain.
func (server *RPCServer) GetTokenByRPC(tokenID string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, tokenID)

	return server.SendQuery(getPrivacyCustomToken, params)
}

// HasSerialNumberByRPC checks if the provided serial numbers have been spent or not.
//
// Returned result in raw json bytes.
func (server *RPCServer) HasSerialNumberByRPC(shardID byte, tokenID string, snList []string) ([]byte, error) {
	if len(snList) == 0 {
		return nil, errors.New("no serial number provided to be checked")
	}

	snQueryList := make([]string, 0)
	for _, sn := range snList {
		snQueryList = append(snQueryList, fmt.Sprintf(`"%s"`, sn))
	}

	addr := rpchandler.CreatePaymentAddress(shardID) // use a random payment address for anonymity

	params := make([]interface{}, 0)
	params = append(params, addr)
	params = append(params, snList)
	params = append(params, tokenID)

	return server.SendQuery(hasSerialNumbers, params)
}

// HasSerialNumberInMemPool checks if the provided serial numbers are currently in the pool or not.
//
// Returned result in raw json bytes.
func (server *RPCServer) HasSerialNumberInMemPool(snList []string) ([]byte, error) {
	if len(snList) == 0 {
		return nil, errors.New("no serial number provided to be checked")
	}

	snQueryList := make([]string, 0)
	for _, sn := range snList {
		snQueryList = append(snQueryList, fmt.Sprintf(`"%s"`, sn))
	}

	params := make([]interface{}, 0)
	params = append(params, snList)

	return server.SendQuery(hasSerialNumbersInMempool, params)
}

// GetBalanceByPrivateKey retrieves the PRV balance of a private key.
//
// NOTE: PrivateKey must be supplied and sent to the server.
func (server *RPCServer) GetBalanceByPrivateKey(privateKey string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, privateKey)

	return server.SendQuery(getBalanceByPrivatekey, params)
}

// SubmitKey submits an OTA key to use the full-node's cache.
func (server *RPCServer) SubmitKey(otaStr string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, otaStr)

	return server.SendQuery(submitKey, params)
}

// AuthorizedSubmitKey submits an OTA Key in an authorized manner for more privileges.
func (server *RPCServer) AuthorizedSubmitKey(otaStr string, accessToken string, fromHeight uint64, isReset bool) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, otaStr)
	params = append(params, accessToken)
	params = append(params, fromHeight)
	params = append(params, isReset)

	return server.SendQuery(authorizedSubmitKey, params)
}

// GetKeySubmissionInfo returns the information of an OTAKey if it has been submitted.
func (server *RPCServer) GetKeySubmissionInfo(otaStr string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, otaStr)

	return server.SendQuery(getKeySubmissionInfo, params)
}

// RandomCommitments gets a list of random commitments to create transactions of version 1.
func (server *RPCServer) RandomCommitments(shardID byte, inputCoins []jsonresult.OutCoin, tokenID string) ([]byte, error) {
	addr := rpchandler.CreatePaymentAddress(shardID) // use a random payment address for anonymity

	params := make([]interface{}, 0)
	params = append(params, addr)
	params = append(params, inputCoins)
	params = append(params, tokenID)

	return server.SendQuery(randomCommitments, params)
}

// RandomCommitmentsAndPublicKeys gets a list of random commitments to create transactions of version 2.
func (server *RPCServer) RandomCommitmentsAndPublicKeys(shardID byte, tokenID string, lenDecoy int) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, shardID)
	params = append(params, lenDecoy)
	params = append(params, tokenID)

	return server.SendQuery(randomCommitmentsAndPublicKeys, params)
}
