package metadata

import "github.com/incognitochain/go-incognito-sdk-v2/common"

// PortalConvertVaultRequest is a request to convert a centralized vault into a PortalV4 vault.
// This Metadata should ONLY be enclosed with a normal (PRV) transaction.
//
// @@NOTE: This tx is created only one time when migrating centralized bridge to portal v4.
type PortalConvertVaultRequest struct {
	MetadataBaseWithSignature
	TokenID          string // pTokenID in incognito chain
	ConvertProof     string
	IncognitoAddress string `json:"IncognitoAddress,omitempty"`
}

// NewPortalConvertVaultRequest creates a new PortalConvertVaultRequest.
func NewPortalConvertVaultRequest(
	metaType int,
	tokenID string,
	convertingProof string,
	incognitoAddress string) (*PortalConvertVaultRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	convertRequestMeta := &PortalConvertVaultRequest{
		TokenID:          tokenID,
		ConvertProof:     convertingProof,
		IncognitoAddress: incognitoAddress,
	}
	convertRequestMeta.MetadataBase = metadataBase
	return convertRequestMeta, nil
}

// Hash overrides MetadataBase.Hash().
func (convertVaultReq PortalConvertVaultRequest) Hash() *common.Hash {
	record := convertVaultReq.MetadataBase.Hash().String()
	record += convertVaultReq.TokenID
	record += convertVaultReq.ConvertProof
	if convertVaultReq.Sig != nil && len(convertVaultReq.Sig) != 0 {
		record += string(convertVaultReq.Sig)
	}

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// HashWithoutSig overrides MetadataBase.HashWithoutSig().
func (convertVaultReq PortalConvertVaultRequest) HashWithoutSig() *common.Hash {
	record := convertVaultReq.MetadataBaseWithSignature.Hash().String()
	record += convertVaultReq.TokenID
	record += convertVaultReq.ConvertProof
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (convertVaultReq *PortalConvertVaultRequest) CalculateSize() uint64 {
	return calculateSize(convertVaultReq)
}
