package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

//========== GET RPCs ==========

// Query the RPC server then return the AutoTxByHash
func (server *RPCServer) getAutoTxByHash(txHash string) (*AutoTxByHash, error) {
	if len(server.GetURL()) == 0 {
		return nil, errors.New("Server has not set mainnet or testnet")
	}
	query := fmt.Sprintf(`{
		"jsonrpc":"1.0",
		"method":"gettransactionbyhash",
		"params":["%s"],
		"id":1
	}`, txHash)
	b, err := server.SendPostRequestWithQuery(query)
	if err != nil {
		return nil, err
	}
	autoTx, txError := ParseAutoTxHashFromBytes(b)
	if txError != nil {
		return nil, err
	}
	return autoTx, nil
}

// Get only the proof of transaction requiring the txHash
func (server *RPCServer) GetProofTransactionByHash(txHash string) (string, error) {
	tx, err := server.getAutoTxByHash(txHash)
	if err != nil {
		return "", err
	}
	return tx.Result.Proof, nil
}

// Get only the Sig of transaction requiring the txHash
func (server *RPCServer) GetSigTransactionByHash(txHash string) (string, error) {
	tx, err := server.getAutoTxByHash(txHash)
	if err != nil {
		return "", err
	}
	return tx.Result.Sig, nil
}

// Get only the BlockHash of transaction requiring the txHash
func (server *RPCServer) GetBlockHashTransactionByHash(txHash string) (string, error) {
	tx, err := server.getAutoTxByHash(txHash)
	if err != nil {
		return "", err
	}
	return tx.Result.BlockHash, nil
}

// Get only the BlockHeight of transaction requiring the txHash
func (server *RPCServer) GetBlockHeightTransactionByHash(txHash string) (int, error) {
	tx, err := server.getAutoTxByHash(txHash)
	if err != nil {
		return -1, err
	}
	return tx.Result.BlockHeight, nil
}

// Get the whole result of rpc call 'gettransactionbyhash'
func (server *RPCServer) GetTransactionByHash(txHash string) ([]byte, error) {
	if len(server.GetURL()) == 0 {
		return []byte{}, errors.New("Server has not set mainnet or testnet")
	}
	query := fmt.Sprintf(`{
		"jsonrpc":"1.0",
		"method":"gettransactionbyhash",
		"params":["%s"],
		"id":1
	}`, txHash)
	return server.SendPostRequestWithQuery(query)
}

//========== END GET RPCs ==========

//========== CREATE TX RPCs ==========

func (server *RPCServer) CreateAndSendTransaction() ([]byte, error) {
	if len(server.GetURL()) == 0 {
		return []byte{}, errors.New("Server has not set mainnet or testnet")
	}
	query := `{
		"jsonrpc": "1.0",
		"method": "createandsendtransaction",
		"params": [
			"112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or",
			{
				"12RuhVZQtGgYmCVzVi49zFZD7gR8SQx8Uuz8oHh6eSZ8PwB2MwaNE6Kkhd6GoykfkRnHNSHz1o2CzMiQBCyFPikHmjvvrZkLERuhcVE":200000000000000,
				"12RxDSnQVjPojzf7uju6dcgC2zkKkg85muvQh347S76wKSSsKPAqXkvfpSeJzyEH3PREHZZ6SKsXLkDZbs3BSqwEdxqprqih4VzANK9":200000000000000,
				"12S6m2LpzN17jorYnLb2ApNKaV2EVeZtd6unvrPT1GH8yHGCyjYzKbywweQDZ7aAkhD31gutYAgfQizb2JhJTgBb3AJ8aB4hyppm2ax":200000000000000,
				"12S42y9fq8xWXx1YpZ6KVDLGx6tLjeFWqbSBo6zGxwgVnPe1rMGxtLs87PyziCzYPEiRGdmwU1ewWFXwjLwog3X71K87ApNUrd3LQB3":200000000000000,
				"12S3yvTvWUJfubx3whjYLv23NtaNSwQMGWWScSaAkf3uQg8xdZjPFD4fG8vGvXjpRgrRioS5zuyzZbkac44rjBfs7mEdgoL4pwKu87u":200000000000000,
				"12S6mGbnS3Df5bGBaUfBTh56NRax4PvFPDhUnxvP9D6cZVjnTx9T4FsVdFT44pFE8KXTGYaHSAmb2MkpnUJzkrAe49EPHkBULM8N2ZJ":200000000000000,
				"12Rs5tQTYkWGzEdPNo2GRA1tjZ5aDCTYUyzXf6SJFq89QnY3US3ZzYSjWHVmmLUa6h8bdHHUuVYoR3iCVRoYDCNn1AfP6pxTz5YL8Aj":200000000000000,
				"12S33dTF3aVsuSxY7iniK3UULUYyLMZumExKm6DPfsqnNepGjgDZqkQCDp1Z7Te9dFKQp7G2WeeYqCr5vcDCfrA3id4x5UvL4yyLrrT":200000000000000
			},
			1,
			1
		],
		"id": 1
	}`
	return server.SendPostRequestWithQuery(query)
}

// GetEncodedTransactionsByHashes retrieves base58-encoded transactions given their hashes.
func (server *RPCServer) GetEncodedTransactionsByHashes(txHashList []string) ([]byte, error) {
	mapParams := make(map[string][]string)
	mapParams["TxHashList"] = txHashList

	params := make([]interface{}, 0)
	params = append(params, mapParams)

	return server.SendQuery(getEncodedTransactionsByHashes, params)
}

// SendRawTx broadcasts a base58-encoded PRV transaction to the network.
func (server *RPCServer) SendRawTx(encodedTx string) ([]byte, error) {
	method := sendRawTransaction
	params := make([]interface{}, 0)
	params = append(params, encodedTx)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) SendRawTokenTx(encodedTx string) ([]byte, error) {
	method := sendRawPrivacyCustomTokenTransaction
	params := make([]interface{}, 0)
	params = append(params, encodedTx)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) GetTxHashBySerialNumber(snList []string, tokenID string, shardID byte) ([]byte, error) {
	if len(server.GetURL()) == 0 {
		return []byte{}, errors.New("server has not set mainnet or testnet")
	}
	method := gettransactionbyserialnumber
	params := make([]interface{}, 0)

	paramList := make(map[string]interface{})
	paramList["SerialNumbers"] = snList
	paramList["TokenID"] = tokenID
	paramList["ShardID"] = shardID

	params = append(params, paramList)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

// GetTxHashByReceiver returns the list of transaction V1 sent to a payment address.
func (server *RPCServer) GetTxHashByReceiver(paymentAddress string) ([]byte, error) {
	if len(server.GetURL()) == 0 {
		return []byte{}, errors.New("server has not set mainnet or testnet")
	}
	method := gettransactionhashbyreceiver
	params := make([]interface{}, 0)

	params = append(params, paymentAddress)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

// GetTxHashByPublicKey returns the list of transactions sent to a public key.
func (server *RPCServer) GetTxHashByPublicKey(publicKeys []string) ([]byte, error) {
	if len(server.GetURL()) == 0 {
		return []byte{}, errors.New("server has not set mainnet or testnet")
	}
	method := gettransactionbypublickey
	params := make([]interface{}, 0)

	paramList := make(map[string]interface{})
	paramList["PublicKeys"] = publicKeys

	params = append(params, paramList)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func (server *RPCServer) CreateAndSendTokenInitTransaction(privateKey string, tokenName, tokenSymbol string, initAmount uint64) ([]byte, error) {
	method := createAndSendTokenInitTransaction
	params := make([]interface{}, 0)

	initParam := TokenInitParam{
		PrivateKey:  privateKey,
		TokenName:   tokenName,
		TokenSymbol: tokenSymbol,
		Amount:      initAmount,
	}

	params = append(params, initParam)

	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)
	query, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}
