package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

type EstimateFeeResult struct {
	EstimateFeeCoinPerKb uint64
	EstimateTxSizeInKb   uint64
}

func (server *RPCServer) EstimateFeeWithEstimator(defaultFee int, shardID byte, numBlock int, tokenID string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := estimateFeeWithEstimator
	params := make([]interface{}, 0)
	params = append(params, defaultFee)

	//Generate fake a payment address for a specific shardID
	fakeWallet, err := wallet.GenRandomWalletForShardID(shardID)
	if err != nil {
		return nil, err
	}
	fakeAddress := fakeWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	params = append(params, fakeAddress)

	params = append(params, numBlock)
	params = append(params, tokenID)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}
