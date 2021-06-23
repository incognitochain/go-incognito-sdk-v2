package metadata

import (
	rCommon "github.com/ethereum/go-ethereum/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// IssuingEVMRequest is a request to mint an amount of an ETH/BSC token in the Incognito network
// after the same amount has been locked in the smart contracts.
type IssuingEVMRequest struct {
	BlockHash  rCommon.Hash
	TxIndex    uint
	Proofs     []string
	IncTokenID common.Hash
	MetadataBase
}

// NewIssuingEVMRequest creates a new IssuingEVMRequest.
func NewIssuingEVMRequest(
	blockHash rCommon.Hash,
	txIndex uint,
	proofs []string,
	incTokenID common.Hash,
	metaType int,
) (*IssuingEVMRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	issuingETHReq := &IssuingEVMRequest{
		BlockHash:  blockHash,
		TxIndex:    txIndex,
		Proofs:     proofs,
		IncTokenID: incTokenID,
	}
	issuingETHReq.MetadataBase = metadataBase
	return issuingETHReq, nil
}

// Hash overrides MetadataBase.Hash().
func (iReq IssuingEVMRequest) Hash() *common.Hash {
	record := iReq.BlockHash.String()
	record += string(iReq.TxIndex)
	proofs := iReq.Proofs
	for _, proofStr := range proofs {
		record += proofStr
	}
	record += iReq.MetadataBase.Hash().String()
	record += iReq.IncTokenID.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (iReq *IssuingEVMRequest) CalculateSize() uint64 {
	return calculateSize(iReq)
}
