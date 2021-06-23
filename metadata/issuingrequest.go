package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// IssuingRequest is a request to shield a centralized token into the Incognito network.
// Only centralized website can send this metadata, and authorization is required.
type IssuingRequest struct {
	ReceiverAddress key.PaymentAddress
	DepositedAmount uint64
	TokenID         common.Hash
	TokenName       string
	MetadataBaseWithSignature
}

// NewIssuingRequest creates a new IssuingRequest.
func NewIssuingRequest(
	receiverAddress key.PaymentAddress,
	depositedAmount uint64,
	tokenID common.Hash,
	tokenName string,
	metaType int,
) (*IssuingRequest, error) {
	metadataBase := NewMetadataBaseWithSignature(metaType)
	issuingReq := &IssuingRequest{
		ReceiverAddress: receiverAddress,
		DepositedAmount: depositedAmount,
		TokenID:         tokenID,
		TokenName:       tokenName,
	}
	issuingReq.MetadataBaseWithSignature = *metadataBase
	return issuingReq, nil
}

// Hash overrides MetadataBase.Hash().
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

// HashWithoutSig overrides MetadataBase.HashWithoutSig().
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

// CalculateSize overrides MetadataBase.CalculateSize().
func (iReq *IssuingRequest) CalculateSize() uint64 {
	return calculateSize(iReq)
}
