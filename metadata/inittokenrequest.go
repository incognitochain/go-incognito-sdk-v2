package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

type InitTokenRequest struct {
	OTAStr      string
	TxRandomStr string
	Amount      uint64
	TokenName   string
	TokenSymbol string
	MetadataBase
}

type InitTokenReqAction struct {
	Meta    InitTokenRequest `json:"meta"`
	TxReqID common.Hash      `json:"txReqID"`
	TokenID common.Hash      `json:"tokenID"`
}

type InitTokenAcceptedInst struct {
	OTAStr        string      `json:"OTAStr"`
	TxRandomStr   string      `json:"TxRandomStr"`
	Amount        uint64      `json:"Amount"`
	TokenID       common.Hash `json:"TokenID"`
	TokenName     string      `json:"TokenName"`
	TokenSymbol   string      `json:"TokenSymbol"`
	TokenType     int         `json:"TokenType"`
	ShardID       byte        `json:"ShardID"`
	RequestedTxID common.Hash `json:"txReqID"`
}

func NewInitTokenRequest(otaStr string, txRandomStr string, amount uint64, tokenName, tokenSymbol string, metaType int) (*InitTokenRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	initPTokenMeta := &InitTokenRequest{
		OTAStr:      otaStr,
		TxRandomStr: txRandomStr,
		TokenName:   tokenName,
		TokenSymbol: tokenSymbol,
		Amount:      amount,
	}
	initPTokenMeta.MetadataBase = metadataBase
	return initPTokenMeta, nil
}

//Hash returns the hash of all components in the request.
func (iReq InitTokenRequest) Hash() *common.Hash {
	record := iReq.MetadataBase.Hash().String()
	record += iReq.OTAStr
	record += iReq.TxRandomStr
	record += iReq.TokenName
	record += iReq.TokenSymbol
	record += strconv.FormatUint(iReq.Amount, 10)

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

//genTokenID generates a (deterministically) random tokenID for the request transaction.
//From now on, users cannot generate their own tokenID.
//The generated tokenID is calculated as the hash of the following components:
//	- The Tx hash
//	- The shardID at which the request is sent
func (iReq *InitTokenRequest) genTokenID(tx Transaction, shardID byte) *common.Hash {
	record := tx.Hash().String()
	record += strconv.FormatUint(uint64(shardID), 10)

	tokenID := common.HashH([]byte(record))
	return &tokenID
}

//CalculateSize returns the size of the request.
func (iReq *InitTokenRequest) CalculateSize() uint64 {
	return calculateSize(iReq)
}
