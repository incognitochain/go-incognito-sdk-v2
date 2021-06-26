package rpc

import (
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

// EstimateFeeResult represents an estimated fee result returned by the remote server.
type EstimateFeeResult struct {
	EstimateFeeCoinPerKb uint64
	EstimateTxSizeInKb   uint64
}

// EstimateFeeWithEstimator retrieves an estimate fee for a tokenID.
func (server *RPCServer) EstimateFeeWithEstimator(defaultFee int, shardID byte, numBlock int, tokenID string) ([]byte, error) {
	addr := rpchandler.CreatePaymentAddress(shardID) // use a random payment address for anonymity

	params := make([]interface{}, 0)
	params = append(params, defaultFee)
	params = append(params, addr)
	params = append(params, numBlock)
	params = append(params, tokenID)

	return server.SendQuery(estimateFeeWithEstimator, params)
}
