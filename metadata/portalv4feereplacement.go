package metadata

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type PortalReplacementFeeRequest struct {
	MetadataBaseWithSignature
	TokenID string
	BatchID string
	Fee     uint
}

func NewPortalReplacementFeeRequest(metaType int, tokenID, batchID string, fee uint) (*PortalReplacementFeeRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}

	portalUnshieldReq := &PortalReplacementFeeRequest{
		TokenID: tokenID,
		BatchID: batchID,
		Fee:     fee,
	}

	portalUnshieldReq.MetadataBase = metadataBase

	return portalUnshieldReq, nil
}

func (repl PortalReplacementFeeRequest) Hash() *common.Hash {
	record := repl.MetadataBase.Hash().String()
	record += repl.TokenID
	record += repl.BatchID
	record += strconv.FormatUint(uint64(repl.Fee), 10)

	if repl.Sig != nil && len(repl.Sig) != 0 {
		record += string(repl.Sig)
	}
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (repl PortalReplacementFeeRequest) HashWithoutSig() *common.Hash {
	record := repl.MetadataBaseWithSignature.Hash().String()
	record += repl.TokenID
	record += repl.BatchID
	record += strconv.FormatUint(uint64(repl.Fee), 10)
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (repl *PortalReplacementFeeRequest) CalculateSize() uint64 {
	return calculateSize(repl)
}
