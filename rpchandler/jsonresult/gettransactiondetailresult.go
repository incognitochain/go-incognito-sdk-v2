package jsonresult

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver1"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"time"
)

type TransactionDetail struct {
	BlockHash   string `json:"BlockHash"`
	BlockHeight uint64 `json:"BlockHeight"`
	TxSize      uint64 `json:"TxSize"`
	Index       uint64 `json:"Index"`
	ShardID     byte   `json:"ShardID"`
	Hash        string `json:"Hash"`
	Version     int8   `json:"Version"`
	Type        string `json:"Type"` // Transaction type
	LockTime    string `json:"LockTime"`
	RawLockTime int64  `json:"RawLockTime,omitempty"`
	Fee         uint64 `json:"Fee"` // Fee applies: always consant
	Image       string `json:"Image"`

	IsPrivacy bool   `json:"IsPrivacy"`
	Proof     string `json:"Proof"`
	//ProofDetail     ProofDetail   `json:"ProofDetail"`
	InputCoinPubKey string `json:"InputCoinPubKey"`
	SigPubKey       string `json:"SigPubKey,omitempty"`    // 64 bytes
	RawSigPubKey    []byte `json:"RawSigPubKey,omitempty"` // 64 bytes
	Sig             string `json:"Sig,omitempty"`          // 64 bytes

	Metadata                 string `json:"Metadata"`
	CustomTokenData          string `json:"CustomTokenData"`
	PrivacyCustomTokenID     string `json:"PrivacyCustomTokenID"`
	PrivacyCustomTokenName   string `json:"PrivacyCustomTokenName"`
	PrivacyCustomTokenSymbol string `json:"PrivacyCustomTokenSymbol"`
	PrivacyCustomTokenData   string `json:"PrivacyCustomTokenData"`
	//PrivacyCustomTokenProofDetail ProofDetail `json:"PrivacyCustomTokenProofDetail"`
	PrivacyCustomTokenIsPrivacy bool   `json:"PrivacyCustomTokenIsPrivacy"`
	PrivacyCustomTokenFee       uint64 `json:"PrivacyCustomTokenFee"`

	IsInMempool bool `json:"IsInMempool"`
	IsInBlock   bool `json:"IsInBlock"`

	Info string `json:"Info"`
}

type TxNormalRPC struct {
	Version              int8
	Type                 string
	LockTime             int64
	Fee                  uint64
	Info                 []byte
	SigPubKey            string
	Sig                  string
	Proof                string
	PubKeyLastByteSender byte
	Metadata             string
}

// ParseTxDetail parses a transaction detail into a transaction object.
func ParseTxDetail(txDetail TransactionDetail) (metadata.Transaction, error) {
	version := txDetail.Version
	txType := txDetail.Type
	txFee := txDetail.Fee
	var info []byte
	if txDetail.Info == "null" {
		info = nil
	} else {
		info = []byte(txDetail.Info)
	}

	//Parse lock time
	lockTime := txDetail.RawLockTime
	if lockTime == 0 {
		tmpLockTime, err := time.Parse(common.DateOutputFormat, txDetail.LockTime)
		if err != nil {
			return nil, fmt.Errorf("decode locktime error: %v", err)
		}
		lockTime = tmpLockTime.Unix()
	}

	//Parse sig
	sig, _, err := base58.Base58Check{}.Decode(txDetail.Sig)
	if err != nil {
		return nil, fmt.Errorf("decode sig error: %v", err)
	}

	//Parse sig pubkey
	sigPubKey := txDetail.RawSigPubKey
	if len(sigPubKey) == 0 {
		if txDetail.Version != 2 {
			sigPubKey, _, err = base58.Base58Check{}.Decode(txDetail.SigPubKey)
			if err != nil {
				return nil, fmt.Errorf("decode sig pubkey %v error: %v", txDetail.SigPubKey, err)
			}
		} else {
			return nil, fmt.Errorf("cannot decode sig pubkey for version 2 without rawSigPubKey")
		}
	}

	//Parse metadata
	var meta metadata.Metadata
	if len(txDetail.Metadata) != 0 {
		meta, err = metadata.ParseMetadata([]byte(txDetail.Metadata))
		if err != nil {
			return nil, fmt.Errorf("parse metadata error: %v", err)
		}
	}

	//Parse proof
	var proof privacy.Proof
	if len(txDetail.Proof) > 0 {
		proof = MakeProof(txType, version, false)
		var proofBytes []byte
		proofBytes, err = base64.StdEncoding.DecodeString(txDetail.Proof)
		if err != nil {
			return nil, fmt.Errorf("decode proof error: %v", err)
		}

		err = proof.SetBytes(proofBytes)
		if err != nil {
			return nil, fmt.Errorf("parse proof error: %v", err)
		}
	}

	var tx metadata.Transaction
	switch txType {
	case common.TxNormalType, common.TxRewardType, common.TxReturnStakingType, common.TxConversionType:
		if version == 1 {
			tx = new(tx_ver1.Tx)
		} else {
			tx = new(tx_ver2.Tx)
		}

		tx.SetType(txType)
		tx.SetTxFee(txFee)
		tx.SetInfo(info)
		tx.SetLockTime(lockTime)
		tx.SetMetadata(meta)
		tx.SetProof(proof)
		tx.SetSig(sig)
		tx.SetSigPubKey(sigPubKey)
		tx.SetVersion(version)
		tx.SetGetSenderAddrLastByte(txDetail.ShardID)

	case common.TxCustomTokenPrivacyType, common.TxTokenConversionType:
		//Parse txTokenData
		var txTokenData tx_generic.TxTokenData
		if txDetail.PrivacyCustomTokenData != "" {
			txTokenData, err = ParseTxTokenData([]byte(txDetail.PrivacyCustomTokenData))
			if err != nil {
				return nil, fmt.Errorf("parse txTokenData error: %v", err)
			}
		}

		var tmpTxToken tx_generic.TransactionToken
		var txBase metadata.Transaction
		switch version {
		case 1:
			tmpTxToken = new(tx_ver1.TxToken)
			txBase = new(tx_ver1.Tx)
		case 2:
			tmpTxToken = new(tx_ver2.TxToken)
			txBase = new(tx_ver2.Tx)
		default:
			tmpTxToken = new(tx_ver1.TxToken)
			txBase = new(tx_ver1.Tx)
		}

		txBase.SetType(txType)
		txBase.SetTxFee(txFee)
		txBase.SetInfo(info)
		txBase.SetLockTime(lockTime)
		txBase.SetMetadata(meta)
		txBase.SetProof(proof)
		txBase.SetSig(sig)
		txBase.SetSigPubKey(sigPubKey)
		txBase.SetVersion(version)
		txBase.SetGetSenderAddrLastByte(txDetail.ShardID)

		err = tmpTxToken.SetTxBase(txBase)
		if err != nil {
			return nil, fmt.Errorf("set txBase error: %v", err)
		}

		err = tmpTxToken.SetTxTokenData(txTokenData)
		if err != nil {
			return nil, fmt.Errorf("set txTokenData error: %v", err)
		}

		tx = tmpTxToken
	default:
		return nil, fmt.Errorf("transaction type %v not found", txType)
	}

	//txHash := tx.Hash().String()
	//if txHash != txDetail.Hash && txDetail.Version == 2 {
	//	return nil, fmt.Errorf("expect txHash to be %v, got %v", txDetail.Hash, txHash)
	//}

	return tx, nil
}

// ParseTxTokenData parses a RPC token detail into a TxTokenData object.
func ParseTxTokenData(txTokenData []byte) (tx_generic.TxTokenData, error) {
	var res tx_generic.TxTokenData
	var tmpTxTokenData struct {
		TxNormal       TxNormalRPC `json:"TxNormal"`
		PropertyID     string      `json:"PropertyID"`
		PropertyName   string      `json:"PropertyName"`
		PropertySymbol string      `json:"PropertySymbol"`
		Type           int         `json:"Type"`
		Mintable       bool        `json:"Mintable"`
		Amount         uint64      `json:"Amount"`
	}

	err := json.Unmarshal(txTokenData, &tmpTxTokenData)
	if err != nil {
		return res, err
	}

	res.Type = tmpTxTokenData.Type
	res.PropertyName = tmpTxTokenData.PropertyName
	res.PropertySymbol = tmpTxTokenData.PropertySymbol
	res.Amount = tmpTxTokenData.Amount
	res.Mintable = tmpTxTokenData.Mintable

	tmpTokenID, err := new(common.Hash).NewHashFromStr(tmpTxTokenData.PropertyID)
	if err != nil {
		return res, fmt.Errorf("cannot decode tokenID %v: %v", tmpTxTokenData.PropertyID, err)
	}
	res.PropertyID = *tmpTokenID

	res.TxNormal, err = NewTxNormalFromTxNormalRPC(tmpTxTokenData.TxNormal)
	if err != nil {
		return res, err
	}

	return res, nil
}

// NewTxNormalFromTxNormalRPC creates a txNormal object from RPC-resulting data.
func NewTxNormalFromTxNormalRPC(data TxNormalRPC) (metadata.Transaction, error) {
	var tx metadata.Transaction
	if data.Version == 2 {
		tx = new(tx_ver2.Tx)
	} else {
		tx = new(tx_ver1.Tx)
	}

	tx.SetVersion(data.Version)
	tx.SetType(data.Type)
	tx.SetLockTime(data.LockTime)
	tx.SetTxFee(data.Fee)
	tx.SetInfo([]byte(data.Info))

	sigBytes, err := base64.StdEncoding.DecodeString(data.Sig)
	if err != nil {
		return nil, fmt.Errorf("decode sig error: %v", err)
	}

	sigPubKeyBytes, err := base64.StdEncoding.DecodeString(data.SigPubKey)
	if err != nil {
		return nil, fmt.Errorf("decode sigPubKey error: %v", err)
	}

	proof := MakeProof(data.Type, data.Version, true)
	if len(data.Proof) > 0 {
		var proofBytes []byte
		proofBytes, err = base64.StdEncoding.DecodeString(data.Proof)
		if err != nil {
			return nil, fmt.Errorf("decode proof error: %v", err)
		}

		err = proof.SetBytes(proofBytes)
		if err != nil {
			return nil, fmt.Errorf("parse proof error: %v", err)
		}
	}

	tx.SetSig(sigBytes)
	tx.SetSigPubKey(sigPubKeyBytes)
	tx.SetProof(proof)

	return tx, nil
}

// MakeProof returns an empty proof associated with the input version and tx type.
func MakeProof(txType string, version int8, isTxNormal bool) privacy.Proof {
	switch txType {
	case common.TxConversionType:
		return new(privacy.ProofForConversion)
	case common.TxTokenConversionType:
		if isTxNormal {
			return new(privacy.ProofForConversion)
		}
		return new(privacy.ProofV2)
	default:
		if version == 2 {
			return new(privacy.ProofV2)
		} else {
			return new(privacy.ProofV1)
		}
	}
}
