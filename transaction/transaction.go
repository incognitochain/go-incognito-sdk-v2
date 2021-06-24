package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver1"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
)

// TxChoice is a helper struct for parsing transactions of all types from JSON.
// After parsing succeeds, one of its fields will have the TX object; others will be nil.
// This can be used to assert the transaction type.
type TxChoice struct {
	Version1      *tx_ver1.Tx      `json:"TxVersion1,omitempty"`
	TokenVersion1 *tx_ver1.TxToken `json:"TxTokenVersion1,omitempty"`
	Version2      *tx_ver2.Tx      `json:"TxVersion2,omitempty"`
	TokenVersion2 *tx_ver2.TxToken `json:"TxTokenVersion2,omitempty"`
}

// txJsonDataVersion is used to parse json.
type txJsonDataVersion struct {
	Version int8 `json:"Version"`
	Type    string
}

// ToTx returns a generic transaction from a TxChoice object.
// Use this when the underlying TX type is irrelevant.
func (ch *TxChoice) ToTx() metadata.Transaction {
	if ch == nil {
		return nil
	}
	// `choice` struct only ever contains 1 non-nil field
	if ch.Version1 != nil {
		return ch.Version1
	}
	if ch.Version2 != nil {
		return ch.Version2
	}
	if ch.TokenVersion1 != nil {
		return ch.TokenVersion1
	}
	if ch.TokenVersion2 != nil {
		return ch.TokenVersion2
	}
	return nil
}

// DeserializeTransactionJSON parses a transaction from raw JSON into a TxChoice object.
// It covers all transaction types.
func DeserializeTransactionJSON(data []byte) (*TxChoice, error) {
	result := &TxChoice{}
	holder := make(map[string]interface{})
	err := json.Unmarshal(data, &holder)
	if err != nil {
		return nil, err
	}
	_, isTokenTx := holder["TxTokenPrivacyData"]
	_, hasVersionOutside := holder["Version"]
	var verHolder txJsonDataVersion
	err = json.Unmarshal(data, &verHolder)
	if err != nil {
		return nil, err
	}

	if hasVersionOutside {
		switch verHolder.Version {
		case utils.TxVersion1Number:
			if isTokenTx {
				// token ver 1
				result.TokenVersion1 = &tx_ver1.TxToken{}
				err := json.Unmarshal(data, result.TokenVersion1)
				return result, err
			} else {
				// tx ver 1
				result.Version1 = &tx_ver1.Tx{}
				err := json.Unmarshal(data, result.Version1)
				return result, err
			}
		case utils.TxVersion2Number: // the same as utils.TxConversionVersion12Number
			if isTokenTx {
				// rejected
				return nil, fmt.Errorf("error unmarshalling TX from JSON : misplaced version")
			} else {
				// tx ver 2
				result.Version2 = &tx_ver2.Tx{}
				err := json.Unmarshal(data, result.Version2)
				return result, err
			}
		default:
			return nil, fmt.Errorf("error unmarshalling TX from JSON : wrong version of %d", verHolder.Version)
		}
	} else {
		if isTokenTx {
			// token ver 2
			result.TokenVersion2 = &tx_ver2.TxToken{}
			err := json.Unmarshal(data, result.TokenVersion2)
			return result, err
		} else {
			return nil, fmt.Errorf("error unmarshalling TX from JSON")
		}
	}

}
