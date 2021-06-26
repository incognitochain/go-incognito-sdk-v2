package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// only centralized website can send this type of tx
type IssuingRequest struct {
	ReceiverAddress key.PaymentAddress
	DepositedAmount uint64
	TokenID         common.Hash
	TokenName       string
	MetadataBaseWithSignature
}

type IssuingReqAction struct {
	Meta    IssuingRequest `json:"meta"`
	TxReqID common.Hash    `json:"txReqId"`
}

type IssuingAcceptedInst struct {
	ShardID         byte               `json:"shardId"`
	DepositedAmount uint64             `json:"issuingAmount"`
	ReceiverAddr    key.PaymentAddress `json:"receiverAddrStr"`
	IncTokenID      common.Hash        `json:"incTokenId"`
	IncTokenName    string             `json:"incTokenName"`
	TxReqID         common.Hash        `json:"txReqId"`
}

func (iReq IssuingRequest) Hash() *common.Hash {
	record := iReq.ReceiverAddress.String()
	record += iReq.TokenID.String()
	// TODO: @hung change to record += fmt.Sprint(iReq.DepositedAmount)
	record += string(iReq.DepositedAmount)
	record += iReq.TokenName
	record += iReq.MetadataBaseWithSignature.Hash().String()
	if iReq.Sig != nil && len(iReq.Sig) != 0 {
		record += string(iReq.Sig)
	}
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iReq IssuingRequest) HashWithoutSig() *common.Hash {
	record := iReq.ReceiverAddress.String()
	record += iReq.TokenID.String()
	record += string(iReq.DepositedAmount)
	record += iReq.TokenName
	record += iReq.MetadataBaseWithSignature.Hash().String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iReq *IssuingRequest) CalculateSize() uint64 {
	return calculateSize(iReq)
}
