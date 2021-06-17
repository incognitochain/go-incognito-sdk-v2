package incclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	rCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

// GetETHTxByHash retrieves an Ethereum transaction from its hash.
func (client *IncClient) GetETHTxByHash(tx string) (map[string]interface{}, error) {
	method := "eth_getTransactionByHash"
	params := []interface{}{tx}

	request := rpchandler.CreateJsonRequest("2.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	responseInBytes, err := client.ethServer.SendPostRequestWithQuery(string(query))

	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetETHBlockByHash retrieves an Ethereum block from its hash.
func (client *IncClient) GetETHBlockByHash(blockHash string) (map[string]interface{}, error) {
	method := "eth_getBlockByHash"
	params := []interface{}{blockHash, false}

	request := rpchandler.CreateJsonRequest("2.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	responseInBytes, err := client.ethServer.SendPostRequestWithQuery(string(query))
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetETHTxReceipt retrieves a Ethereum transaction receipt from its hash.
func (client *IncClient) GetETHTxReceipt(txHash string) (*types.Receipt, error) {
	method := "eth_getTransactionReceipt"
	params := []interface{}{txHash}

	request := rpchandler.CreateJsonRequest("2.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	responseInBytes, err := client.ethServer.SendPostRequestWithQuery(string(query))
	if err != nil {
		return nil, err
	}

	var res types.Receipt
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetETHDepositProof retrieves an Ethereum-depositing proof of a transaction hash.
func (client *IncClient) GetETHDepositProof(txHash string) (*ETHDepositProof, uint64, error) {
	// Get tx content
	txContent, err := client.GetETHTxByHash(txHash)
	if err != nil {
		fmt.Println("cannot get eth by hash", err)
		return nil, 0, err
	}

	_, ok := txContent["value"]
	if !ok {
		return nil, 0, fmt.Errorf("cannot find value in %v", txContent)
	}
	valueStr, ok := txContent["value"].(string)
	if !ok {
		return nil, 0, fmt.Errorf("cannot parse value in %v", txContent)
	}
	amtBigInt, ok := new(big.Int).SetString(valueStr[2:], 16)
	if !ok {
		return nil, 0, fmt.Errorf("cannot set bigInt value in %v", txContent)
	}
	var amount uint64
	//If ETH, divide ETH to 10^9
	amount = big.NewInt(0).Div(amtBigInt, big.NewInt(1000000000)).Uint64()

	_, ok = txContent["blockHash"]
	if !ok {
		return nil, 0, fmt.Errorf("cannot find blockHash in %v", txContent)
	}
	blockHashStr, ok := txContent["blockHash"].(string)
	if !ok {
		return nil, 0, fmt.Errorf("cannot parse blockHash in %v", txContent)
	}
	blockHash := rCommon.HexToHash(blockHashStr)

	_, ok = txContent["transactionIndex"]
	if !ok {
		return nil, 0, fmt.Errorf("cannot find transactionIndex in %v", txContent)
	}
	txIndexStr, ok := txContent["transactionIndex"].(string)
	if !ok {
		return nil, 0, fmt.Errorf("cannot parse transactionIndex in %v", txContent)
	}

	txIndex, err := strconv.ParseUint(txIndexStr[2:], 16, 64)
	if err != nil {
		return nil, 0, err
	}

	// Get txs block for constructing receipt trie
	_, ok = txContent["blockNumber"]
	if !ok {
		return nil, 0, fmt.Errorf("cannot find blockNumber in %v", txContent)
	}
	blockNumString, ok := txContent["blockNumber"].(string)
	if !ok {
		return nil, 0, fmt.Errorf("cannot parse blockNumber in %v", txContent)
	}
	blockNumber, err := strconv.ParseInt(blockNumString[2:], 16, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot convert blockNumber into integer")
	}

	blockHeader, err := client.GetETHBlockByHash(blockHashStr)
	if err != nil {
		return nil, 0, err
	}

	// Get all sibling Txs
	_, ok = blockHeader["transactions"]
	if !ok {
		return nil, 0, fmt.Errorf("cannot find transactions in %v", txContent)
	}
	siblingTxs, ok := blockHeader["transactions"].([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("cannot parse transactions in %v", txContent)
	}

	fmt.Println("length of transactions in block", len(siblingTxs))

	// Constructing the receipt trie (source: go-ethereum/core/types/derive_sha.go)
	keyBuf := new(bytes.Buffer)
	receiptTrie := new(trie.Trie)
	fmt.Println("Start creating receipt trie...")
	for i, tx := range siblingTxs {
		txStr, ok := tx.(string)
		if !ok {
			return nil, 0, fmt.Errorf("cannot parse sibling tx: %v", tx)
		}
		siblingReceipt, err := client.GetETHTxReceipt(txStr)
		if err != nil {
			return nil, 0, err
		}
		keyBuf.Reset()
		err = rlp.Encode(keyBuf, uint(i))
		if err != nil {
			return nil, 0, fmt.Errorf("rlp encode returns an error: %v", err)
		}
		encodedReceipt, err := rlp.EncodeToBytes(siblingReceipt)
		if err != nil {
			return nil, 0, err
		}
		receiptTrie.Update(keyBuf.Bytes(), encodedReceipt)
	}

	fmt.Println("Finish creating receipt trie.")

	// Constructing the proof for the current receipt (source: go-ethereum/trie/proof.go)
	proof := light.NewNodeSet()
	keyBuf.Reset()
	err = rlp.Encode(keyBuf, uint(txIndex))
	if err != nil {
		return nil, 0, fmt.Errorf("rlp encode returns an error: %v", err)
	}
	fmt.Println("Start proving receipt trie...")
	err = receiptTrie.Prove(keyBuf.Bytes(), 0, proof)
	if err != nil {
		return nil, 0, err
	}
	fmt.Println("Finish proving receipt trie.")

	nodeList := proof.NodeList()
	encNodeList := make([]string, 0)
	for _, node := range nodeList {
		str := base64.StdEncoding.EncodeToString(node)
		encNodeList = append(encNodeList, str)
	}

	return NewETHDepositProof(uint(blockNumber), blockHash, uint(txIndex), encNodeList), amount, nil
}

// GetMostRecentETHBlockNumber retrieves the most recent Ethereum block number.
func (client *IncClient) GetMostRecentETHBlockNumber() (uint64, error) {
	method := "eth_blockNumber"
	params := make([]interface{}, 0)

	request := rpchandler.CreateJsonRequest("2.0", method, params, 1)
	query, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	responseInBytes, err := client.ethServer.SendPostRequestWithQuery(string(query))

	if err != nil {
		return 0, err
	}

	var hexResult string
	err = rpchandler.ParseResponse(responseInBytes, &hexResult)
	if err != nil {
		return 0, err
	}

	res, ok := new(big.Int).SetString(hexResult[2:], 16)
	if !ok {
		return 0, fmt.Errorf("cannot set hex to big: %v", hexResult)
	}

	return res.Uint64(), nil
}

// GetETHTransactionStatus returns the status of an Ethereum transaction.
func (client *IncClient) GetETHTransactionStatus(txHash string) (int, error) {
	receipt, err := client.GetETHTxReceipt(txHash)
	if err != nil {
		return -1, err
	}

	return int(receipt.Status), nil
}
