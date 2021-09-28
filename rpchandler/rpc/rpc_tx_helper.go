package rpc

// ListCustomToken represents the custom-token listing result returned by the remote server.
type ListCustomToken struct {
	ID     int `json:"Id"`
	Result struct {
		ListCustomToken []struct {
			ID                 string        `json:"ID"`
			Name               string        `json:"Name"`
			Symbol             string        `json:"Symbol"`
			Image              string        `json:"Image"`
			Amount             float64       `json:"Amount"`
			IsPrivacy          bool          `json:"IsPrivacy"`
			IsBridgeToken      bool          `json:"IsBridgeToken"`
			ListTxs            []interface{} `json:"ListTxs"`
			CountTxs           int           `json:"CountTxs"`
			InitiatorPublicKey string        `json:"InitiatorPublicKey"`
			TxInfo             string        `json:"TxInfo"`
		} `json:"ListCustomToken"`
	} `json:"Result"`
	Error   interface{}   `json:"Error"`
	Params  []interface{} `json:"Params"`
	Method  string        `json:"Method"`
	Jsonrpc string        `json:"JsonRPC"`
}

// TokenInitParam represents the parameters needed for the RPC createAndSendTokenInitTransaction.
type TokenInitParam struct {
	PrivateKey  string `json:"PrivateKey"`
	TokenName   string `json:"TokenName"`
	TokenSymbol string `json:"TokenSymbol"`
	Amount      uint64 `json:"Amount"`
}
