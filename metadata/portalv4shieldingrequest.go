package metadata

import "github.com/incognitochain/go-incognito-sdk-v2/common"

// PortalShieldingRequest is a Metadata that a portal user requests to mint pToken (after sending public tokens to a multi-sig wallet)
// This Metadata should ONLY be enclosed with a normal (PRV) transaction.
type PortalShieldingRequest struct {
	MetadataBase
	TokenID        string // pTokenID in incognito chain
	IncAddressStr  string `json:"IncogAddressStr"`
	ShieldingProof string
}

// PortalShieldingRequestStatus represents the status of an un-shield request on the Portal.
type PortalShieldingRequestStatus struct {
	Status        byte
	Error         string
	TokenID       string
	IncAddressStr string `json:"IncogAddressStr"`
	ProofHash     string
	MintingAmount uint64
	TxReqID       common.Hash
	ExternalTxID  string
}

// NewPortalShieldingRequest creates a new PortalShieldingRequest.
func NewPortalShieldingRequest(
	metaType int,
	tokenID string,
	incAddressStr string,
	shieldingProof string) (*PortalShieldingRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	shieldingRequestMeta := &PortalShieldingRequest{
		TokenID:        tokenID,
		IncAddressStr:  incAddressStr,
		ShieldingProof: shieldingProof,
	}
	shieldingRequestMeta.MetadataBase = metadataBase
	return shieldingRequestMeta, nil
}

// Hash overrides MetadataBase.Hash().
func (shieldingReq PortalShieldingRequest) Hash() *common.Hash {
	record := shieldingReq.MetadataBase.Hash().String()
	record += shieldingReq.TokenID
	record += shieldingReq.IncAddressStr
	record += shieldingReq.ShieldingProof
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (shieldingReq *PortalShieldingRequest) CalculateSize() uint64 {
	return calculateSize(shieldingReq)
}
