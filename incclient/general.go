package incclient

import (
	"encoding/json"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
)

// GetActiveShard returns the number of active shards on the Incognito network.
func (client *IncClient) GetActiveShard() (int, error) {
	responseInBytes, err := client.rpcServer.GetActiveShards()
	if err != nil {
		return 0, err
	}

	response, err := rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return 0, err
	}

	var activeShards int
	err = json.Unmarshal(response.Result, &activeShards)

	return activeShards, err
}

// GetBestBlock returns the best blocks of the beacon chain and each shard.
func (client *IncClient) GetBestBlock() (map[int]uint64, error) {
	responseInBytes, err := client.rpcServer.GetBestBlock()
	if err != nil {
		return nil, err
	}

	response, err := rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return nil, err
	}

	var bestBlocksResult jsonresult.GetBestBlockResult
	err = json.Unmarshal(response.Result, &bestBlocksResult)
	if err != nil {
		return nil, err
	}

	res := make(map[int]uint64)

	for key, value := range bestBlocksResult.BestBlocks {
		res[key] = value.Height
	}

	return res, nil
}

// GetListToken returns all tokens currently on the Incognito network.
func (client *IncClient) GetListToken() (map[string]CustomToken, error) {
	responseInBytes, err := client.rpcServer.ListPrivacyCustomTokenByRPC()
	if err != nil {
		return nil, err
	}
	var res rpc.ListCustomToken
	err = json.Unmarshal(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	tokenCount := 0
	listTokens := make(map[string]CustomToken)
	for _, token := range res.Result.ListCustomToken {
		tmp := CustomToken{
			tokenID:   token.ID,
			tokenName: token.Name,
			amount:    uint64(token.Amount),
		}
		if len(tmp.tokenName) == 0 {
			tmp.tokenName = fmt.Sprintf("%d", tokenCount)
		}

		listTokens[token.ID] = tmp
		tokenCount++
	}

	return listTokens, nil
}

// GetRawMemPool returns a list of transaction hashes currently in the pool.
func (client *IncClient) GetRawMemPool() ([]string, error) {
	responseInBytes, err := client.rpcServer.GetRawMempool()
	if err != nil {
		return nil, err
	}

	response, err := rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return nil, err
	}

	var txHashes map[string][]string
	err = json.Unmarshal(response.Result, &txHashes)
	if err != nil {
		return nil, err
	}

	txList, ok := txHashes["TxHashes"]
	if !ok {
		return nil, fmt.Errorf("TxHashes not found in %v", txHashes)
	}

	return txList, nil

}

// SubmitKey submits an OTAKey to the full node.
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

	_, err = rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return err
	}
	return nil
}
