package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

type PortalShieldingResponse struct {
	MetadataBase

	// RequestStatus is the status of the shielding request.
	RequestStatus string

	// ReqTxID is the hash of the shielding request transaction.
	ReqTxID common.Hash

	// Receiver is the same as in the request.
	// If Receiver is an Incognito payment address, SharedRandom must not be empty.
	// If Receiver is an OTAReceiver, SharedRandom is not required.
	Receiver string `json:"RequesterAddrStr"` // the json-tag is required for backward-compatibility.

	// MintingAmount is the shielding amount.
	MintingAmount uint64

	// IncTokenID is the Incognito ID of the shielding token.
	IncTokenID string

	// SharedRandom is combined with Receiver to make sure the minting amount is for the eligible party.
	SharedRandom []byte `json:"SharedRandom,omitempty"`
}

func NewPortalShieldingResponse(
	depositStatus string,
	reqTxID common.Hash,
	receiver string,
	amount uint64,
	tokenID string,
	metaType int,
) *PortalShieldingResponse {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	return &PortalShieldingResponse{
		RequestStatus: depositStatus,
		ReqTxID:       reqTxID,
		MetadataBase:  metadataBase,
		Receiver:      receiver,
		MintingAmount: amount,
		IncTokenID:    tokenID,
	}
}

func (iRes PortalShieldingResponse) Hash() *common.Hash {
	record := iRes.MetadataBase.Hash().String()
	record += iRes.RequestStatus
	record += iRes.ReqTxID.String()
	record += iRes.Receiver
	record += strconv.FormatUint(iRes.MintingAmount, 10)
	record += iRes.IncTokenID
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iRes *PortalShieldingResponse) CalculateSize() uint64 {
	return calculateSize(iRes)
}

func (iRes *PortalShieldingResponse) SetSharedRandom(r []byte) {
	iRes.SharedRandom = r
}
