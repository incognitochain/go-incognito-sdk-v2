package jsonresult

// CreateTransactionResult describes an RPC-result of creating PRV transactions.
type CreateTransactionResult struct {
	Base58CheckData string
	TxID            string
	ShardID         byte
}

// CreateTransactionTokenResult describes an RPC-result of creating token transactions.
type CreateTransactionTokenResult struct {
	Base58CheckData string
	ShardID         byte   `json:"ShardID"`
	TxID            string `json:"TxID"`
	TokenID         string `json:"TokenID"`
	TokenName       string `json:"TokenName"`
	TokenAmount     uint64 `json:"TokenAmount"`
}
