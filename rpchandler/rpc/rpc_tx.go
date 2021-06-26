package rpc

// GetTransactionByHash retrieves the transaction detail given its hash.
func (server *RPCServer) GetTransactionByHash(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)

	return server.SendQuery(getTransactionByHash, params)
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
	params := make([]interface{}, 0)
	params = append(params, encodedTx)

	return server.SendQuery(sendRawTransaction, params)
}

// SendRawTokenTx broadcasts a base58-encoded token transaction to the network.
func (server *RPCServer) SendRawTokenTx(encodedTx string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, encodedTx)

	return server.SendQuery(sendRawPrivacyCustomTokenTransaction, params)
}

// GetTxHashBySerialNumber returns the list of transactions which have spent the given serial numbers.
func (server *RPCServer) GetTxHashBySerialNumber(snList []string, tokenID string, shardID byte) ([]byte, error) {
	paramList := make(map[string]interface{})
	paramList["SerialNumbers"] = snList
	paramList["TokenID"] = tokenID
	paramList["ShardID"] = shardID

	params := make([]interface{}, 0)
	params = append(params, paramList)

	return server.SendQuery(gettransactionbyserialnumber, params)
}

// GetTxHashByReceiver returns the list of transactions V1 sent to a payment address.
func (server *RPCServer) GetTxHashByReceiver(paymentAddress string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, paymentAddress)

	return server.SendQuery(gettransactionhashbyreceiver, params)
}

// GetTxHashByPublicKey returns the list of transactions sent to a public key.
func (server *RPCServer) GetTxHashByPublicKey(publicKeys []string) ([]byte, error) {
	paramList := make(map[string]interface{})
	paramList["PublicKeys"] = publicKeys

	params := make([]interface{}, 0)
	params = append(params, paramList)

	return server.SendQuery(gettransactionbypublickey, params)
}

// CreateAndSendTokenInitTransaction has the server create and broadcast a transaction that initializes a new token on
// the network.
//
// NOTE: PrivateKey must be supplied and sent to the server.
func (server *RPCServer) CreateAndSendTokenInitTransaction(privateKey string,
	tokenName,
	tokenSymbol string,
	initAmount uint64) ([]byte, error) {
	initParam := TokenInitParam{
		PrivateKey:  privateKey,
		TokenName:   tokenName,
		TokenSymbol: tokenSymbol,
		Amount:      initAmount,
	}

	params := make([]interface{}, 0)
	params = append(params, initParam)

	return server.SendQuery(createAndSendTokenInitTransaction, params)
}
