package bridge

import (
	rCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type IssuingEVMRequest struct {
	BlockHash  rCommon.Hash
	TxIndex    uint
	ProofStrs  []string
	IncTokenID common.Hash
	NetworkID  uint `json:"NetworkID,omitempty"`
	metadataCommon.MetadataBase
}

type IssuingEVMReqAction struct {
	Meta       IssuingEVMRequest `json:"meta"`
	TxReqID    common.Hash       `json:"txReqId"`
	EVMReceipt *types.Receipt    `json:"ethReceipt"` // don't update the jsontag to make it compatible with the old shielding eth tx
}

type IssuingEVMAcceptedInst struct {
	ShardID         byte        `json:"shardId"`
	IssuingAmount   uint64      `json:"issuingAmount"`
	ReceiverAddrStr string      `json:"receiverAddrStr"`
	IncTokenID      common.Hash `json:"incTokenId"`
	TxReqID         common.Hash `json:"txReqId"`
	UniqTx          []byte      `json:"uniqETHTx"` // don't update the jsontag to make it compatible with the old shielding eth tx
	ExternalTokenID []byte      `json:"externalTokenId"`
}

const (
	LegacyTxType = iota
	AccessListTxType
	DynamicFeeTxType
)

func NewIssuingEVMRequest(
	blockHash rCommon.Hash,
	txIndex uint,
	proofStrs []string,
	incTokenID common.Hash,
	networkID uint,
	metaType int,
) (*IssuingEVMRequest, error) {
	metadataBase := metadataCommon.MetadataBase{
		Type: metaType,
	}
	issuingEVMReq := &IssuingEVMRequest{
		BlockHash:  blockHash,
		TxIndex:    txIndex,
		ProofStrs:  proofStrs,
		IncTokenID: incTokenID,
		NetworkID:  networkID,
	}
	issuingEVMReq.MetadataBase = metadataBase
	return issuingEVMReq, nil
}

func NewIssuingEVMRequestFromMap(
	data map[string]interface{},
	networkID uint,
	metatype int,
) (*IssuingEVMRequest, error) {
	blockHash := rCommon.HexToHash(data["BlockHash"].(string))
	txIdx := uint(data["TxIndex"].(float64))
	proofsRaw := data["ProofStrs"].([]interface{})
	proofStrs := []string{}
	for _, item := range proofsRaw {
		proofStrs = append(proofStrs, item.(string))
	}

	incTokenID, err := common.Hash{}.NewHashFromStr(data["IncTokenID"].(string))
	if err != nil {
		return nil, errors.Errorf("TokenID incorrect")
	}

	req, _ := NewIssuingEVMRequest(
		blockHash,
		txIdx,
		proofStrs,
		*incTokenID,
		networkID,
		metatype,
	)
	return req, nil
}

type EVMProof struct {
	BlockHash rCommon.Hash `json:"BlockHash"`
	TxIndex   uint         `json:"TxIndex"`
	Proof     []string     `json:"Proof"`
}

func NewIssuingEVMRequestFromProofData(proofData EVMProof, networkID uint, incTokenID common.Hash) (*IssuingEVMRequest, error) {
	evmShieldRequest, _ := NewIssuingEVMRequest(
		proofData.BlockHash, proofData.TxIndex, proofData.Proof, incTokenID, networkID,
		metadataCommon.IssuingUnifiedTokenRequestMeta,
	) // error always null
	return evmShieldRequest, nil
}

func (iReq IssuingEVMRequest) Hash() *common.Hash {
	record := iReq.BlockHash.String()
	record += string(iReq.TxIndex)
	proofStrs := iReq.ProofStrs
	for _, proofStr := range proofStrs {
		record += proofStr
	}
	record += iReq.MetadataBase.Hash().String()
	record += iReq.IncTokenID.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iReq *IssuingEVMRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(iReq)
}
