package rpc

import "encoding/json"

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

type AutoTxByHash struct {
	ID     int `json:"Id"`
	Result struct {
		BlockHash   string `json:"BlockHash"`
		BlockHeight int    `json:"BlockHeight"`
		TxSize      int    `json:"TxSize"`
		Index       int    `json:"Index"`
		ShardID     int    `json:"ShardID"`
		Hash        string `json:"Hash"`
		Version     int    `json:"Version"`
		Type        string `json:"Type"`
		LockTime    string `json:"LockTime"`
		Fee         int    `json:"Fee"`
		Image       string `json:"Image"`
		IsPrivacy   bool   `json:"IsPrivacy"`
		Proof       string `json:"Proof"`
		ProofDetail struct {
			InputCoins []struct {
				CoinDetails struct {
					PublicKey      string `json:"PublicKey"`
					CoinCommitment string `json:"CoinCommitment"`
					SNDerivator    struct {
					} `json:"SNDerivator"`
					SerialNumber string `json:"SerialNumber"`
					Randomness   struct {
					} `json:"Randomness"`
					Value int    `json:"Value"`
					Info  string `json:"Info"`
				} `json:"CoinDetails"`
				CoinDetailsEncrypted string `json:"CoinDetailsEncrypted"`
			} `json:"InputCoins"`
			OutputCoins []struct {
				CoinDetails struct {
					PublicKey      string `json:"PublicKey"`
					CoinCommitment string `json:"CoinCommitment"`
					SNDerivator    struct {
					} `json:"SNDerivator"`
					SerialNumber string `json:"SerialNumber"`
					Randomness   struct {
					} `json:"Randomness"`
					Value int    `json:"Value"`
					Info  string `json:"Info"`
				} `json:"CoinDetails"`
				CoinDetailsEncrypted string `json:"CoinDetailsEncrypted"`
			} `json:"OutputCoins"`
		} `json:"ProofDetail"`
		InputCoinPubKey               string `json:"InputCoinPubKey"`
		SigPubKey                     string `json:"SigPubKey"`
		Sig                           string `json:"Sig"`
		Metadata                      string `json:"Metadata"`
		CustomTokenData               string `json:"CustomTokenData"`
		PrivacyCustomTokenID          string `json:"PrivacyCustomTokenID"`
		PrivacyCustomTokenName        string `json:"PrivacyCustomTokenName"`
		PrivacyCustomTokenSymbol      string `json:"PrivacyCustomTokenSymbol"`
		PrivacyCustomTokenData        string `json:"PrivacyCustomTokenData"`
		PrivacyCustomTokenProofDetail struct {
			InputCoins  interface{} `json:"InputCoins"`
			OutputCoins interface{} `json:"OutputCoins"`
		} `json:"PrivacyCustomTokenProofDetail"`
		PrivacyCustomTokenIsPrivacy bool   `json:"PrivacyCustomTokenIsPrivacy"`
		PrivacyCustomTokenFee       int    `json:"PrivacyCustomTokenFee"`
		IsInMempool                 bool   `json:"IsInMempool"`
		IsInBlock                   bool   `json:"IsInBlock"`
		Info                        string `json:"Info"`
	} `json:"Result"`
	Error   interface{} `json:"Error"`
	Params  []string    `json:"Params"`
	Method  string      `json:"Method"`
	Jsonrpc string      `json:"JsonRPC"`
}

type TokenInitParam struct {
	PrivateKey  string `json:"PrivateKey"`
	TokenName   string `json:"TokenName"`
	TokenSymbol string `json:"TokenSymbol"`
	Amount      uint64 `json:"Amount"`
}

var privIndicator string = "1"

// Parse from byte to AutoTxByHash
func ParseAutoTxHashFromBytes(b []byte) (*AutoTxByHash, error) {
	data := new(AutoTxByHash)
	err := json.Unmarshal(b, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
