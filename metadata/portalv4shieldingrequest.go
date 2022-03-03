package metadata

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// PortalV4ShieldingRequest represents a shielding request of Portal V4. Users create transactions with this metadata after
// sending public tokens to multi-sig wallets. There are two ways to use this metadata, depending on how the corresponding
// multi-sig wallet (a.k.a. depositing address) is generated:
// 	- using payment address: Receiver must be a payment address, OTDepositPubKey, Signature must be empty and the corresponding
//	deposit address must be built with Receiver as the chain-code;
//	- using one-time depositing public key: Receiver must be an OTAReceiver, OTDepositPubKey must not be empty,
// 	a signature is required and the corresponding deposit address must be built with OTDepositPubKey as the chain-code.
type PortalV4ShieldingRequest struct {
	MetadataBase

	// TokenID is the Incognito tokenID of the shielding token.
	TokenID string // pTokenID in incognito chain

	// OTDepositPubKey is the base58-encoded public key for this shielding request, used to validate the authenticity of the request.
	// This field is only used with one-time depositing addresses.
	// If set to empty, Receiver must be a payment address. Otherwise, Receiver must be an OTAReceiver.
	OTDepositPubKey string `json:"OTDepositPubKey,omitempty"`

	// Signature is the signature for validating the authenticity of the request. This signature is different from a
	// MetadataBaseWithSignature type since it is signed with the tx privateKey.
	Signature []byte `json:"Signature,omitempty"`

	// Receiver is the recipient of this shielding request.
	// Receiver is
	//	- an Incognito payment address if OTDepositPubKey is empty;
	//	- an OTAReceiver if OTDepositPubKey is not empty.
	Receiver string `json:"IncogAddressStr"` // the json-tag is required for backward-compatibility.

	// ShieldingProof is the generated proof for this shielding request.
	ShieldingProof string
}

// PortalShieldingRequestStatus is used for beacon to track the status of a shielding request.
type PortalShieldingRequestStatus struct {
	Status          byte
	Error           string
	TokenID         string
	OTDepositPubKey string `json:"OTDepositPubKey,omitempty"`
	Receiver        string `json:"IncogAddressStr"` // the json-tag is required for backward-compatibility.
	ProofHash       string
	MintingAmount   uint64
	TxReqID         common.Hash
	ExternalTxID    string
}

// NewPortalShieldingRequest creates a new PortalV4ShieldingRequest based on given data.
// If depositPubKey is not nil or empty, it will create a request with a signature.
func NewPortalShieldingRequest(
	metaType int,
	tokenID string,
	receiver string,
	shieldingProof string,
	depositPubKey string,
	signature []byte) (*PortalV4ShieldingRequest, error) {
	shieldingRequestMeta := &PortalV4ShieldingRequest{
		TokenID:        tokenID,
		Receiver:       receiver,
		ShieldingProof: shieldingProof,
		MetadataBase:   MetadataBase{Type: metaType},
	}

	if len(depositPubKey) != 0 {
		shieldingRequestMeta.Signature = signature
		shieldingRequestMeta.OTDepositPubKey = depositPubKey
	}

	return shieldingRequestMeta, nil
}

func (req PortalV4ShieldingRequest) Hash() *common.Hash {
	var record string
	if req.OTDepositPubKey != "" {
		jsb, _ := json.Marshal(req)
		hash := common.HashH(jsb)
		return &hash
	}

	// old shielding request
	record = req.MetadataBase.Hash().String()
	record += req.TokenID
	record += req.Receiver
	record += req.ShieldingProof
	hash := common.HashH([]byte(record))

	return &hash
}

func (req *PortalV4ShieldingRequest) CalculateSize() uint64 {
	return calculateSize(req)
}
