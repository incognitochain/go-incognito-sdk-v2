package metadata

import (
	rCommon "github.com/ethereum/go-ethereum/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// IssuingEVMRequest represents an EVM shielding request. Users create transactions with this metadata after
// sending public tokens to the corresponding smart contract. There are two ways to use this metadata,
// depending on which data has been enclosed with the depositing transaction:
// 	- payment address: Receiver and Signature must be empty;
//	- using one-time depositing public key: Receiver must be an OTAReceiver, a signature is required.
type IssuingEVMRequest struct {
	// BlockHash is the hash of the block where the public depositing transaction resides in.
	BlockHash rCommon.Hash

	// TxIndex is the index of the public transaction in the BlockHash.
	TxIndex uint

	// ProofStrs is the generated proof for this shielding request.
	ProofStrs []string

	// IncTokenID is the Incognito tokenID of the shielding token.
	IncTokenID common.Hash

	// Signature is the signature for validating the authenticity of the request. This signature is different from a
	// MetadataBaseWithSignature type since it is signed with the tx privateKey.
	Signature []byte `json:"Signature,omitempty"`

	// Receiver is the recipient of this shielding request. It is an OTAReceiver if OTDepositPubKey is not empty.
	Receiver string `json:"Receiver,omitempty"`

	MetadataBase
}

func NewIssuingEVMRequest(
	blockHash rCommon.Hash,
	txIndex uint,
	proofStrs []string,
	incTokenID common.Hash,
	receiver string,
	signature []byte,
	metaType int,
) (*IssuingEVMRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	issuingEVMReq := &IssuingEVMRequest{
		BlockHash:  blockHash,
		TxIndex:    txIndex,
		ProofStrs:  proofStrs,
		IncTokenID: incTokenID,
		Receiver:   receiver,
		Signature:  signature,
	}
	issuingEVMReq.MetadataBase = metadataBase
	return issuingEVMReq, nil
}

// Hash overrides MetadataBase.Hash().
func (iReq IssuingEVMRequest) Hash() *common.Hash {
	record := iReq.BlockHash.String()
	record += string(iReq.TxIndex)
	proofs := iReq.ProofStrs
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
