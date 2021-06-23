package utils

import (
	"encoding/json"
	"errors"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
)

func ParseProof(p interface{}, ver int8, txType string) (privacy.Proof, error) {
	// If transaction is nonPrivacyNonInput then we do not have proof, so parse it as nil
	if p == nil {
		return nil, nil
	}

	proofInBytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	if string(proofInBytes) == "null" {
		return nil, nil
	}

	var res privacy.Proof
	switch txType {
	case common.TxConversionType:
		if ver == TxConversionVersion12Number {
			res = new(privacy.ProofForConversion)
			res.Init()
		} else {
			return nil, errors.New("ParseProof: TxConversion version is incorrect")
		}
	default:
		switch ver {
		case TxVersion1Number, TxVersion0Number:
			res = new(privacy.ProofV1)
		case TxVersion2Number:
			res = new(privacy.ProofV2)
			res.Init()
		default:
			return nil, errors.New("ParseProof: Tx.Version is incorrect")
		}
	}

	err = json.Unmarshal(proofInBytes, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
