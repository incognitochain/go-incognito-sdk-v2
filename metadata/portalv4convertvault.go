package metadata

import "github.com/incognitochain/go-incognito-sdk-v2/common"

// @@NOTE: This tx is created only one time when migration centralized bridge to portal v4
// PortalConvertVaultRequest
// metadata - portal centralized incognito address convert vault request - create normal tx with this metadata
type PortalConvertVaultRequest struct {
	MetadataBaseWithSignature
	TokenID          string // pTokenID in incognito chain
	ConvertProof     string
	IncognitoAddress string
}

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

func (convertVaultReq PortalConvertVaultRequest) HashWithoutSig() *common.Hash {
	record := convertVaultReq.MetadataBaseWithSignature.Hash().String()
	record += convertVaultReq.TokenID
	record += convertVaultReq.ConvertProof
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (convertVaultReq *PortalConvertVaultRequest) CalculateSize() uint64 {
	return calculateSize(convertVaultReq)
}
