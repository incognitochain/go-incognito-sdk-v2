package metadata

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// BurningRequest is a request to burn UTXOs (mostly to un-shield).
// Anyone can send this request.
type BurningRequest struct {
	BurnerAddress key.PaymentAddress
	BurningAmount uint64 // must be equal to vOut value
	TokenID       common.Hash
	TokenName     string
	RemoteAddress string
	MetadataBase
}

// NewBurningRequest creates a new BurningRequest.
func NewBurningRequest(
	burnerAddress key.PaymentAddress,
	burningAmount uint64,
	tokenID common.Hash,
	tokenName string,
	remoteAddress string,
	metaType int,
) (*BurningRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	burningReq := &BurningRequest{
		BurnerAddress: burnerAddress,
		BurningAmount: burningAmount,
		TokenID:       tokenID,
		TokenName:     tokenName,
		RemoteAddress: remoteAddress,
	}
	burningReq.MetadataBase = metadataBase
	return burningReq, nil
}

// Hash overrides MetadataBase.Hash().
func (bReq BurningRequest) Hash() *common.Hash {
	record := bReq.MetadataBase.Hash().String()
	record += bReq.BurnerAddress.String()
	record += bReq.TokenID.String()
	record += strconv.FormatUint(bReq.BurningAmount, 10)
	record += bReq.TokenName
	record += bReq.RemoteAddress
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// HashWithoutSig overrides MetadataBase.HashWithoutSig().
func (bReq BurningRequest) HashWithoutSig() *common.Hash {
	record := bReq.MetadataBase.Hash().String()
	record += bReq.BurnerAddress.String()
	record += bReq.TokenID.String()
	record += strconv.FormatUint(bReq.BurningAmount, 10)
	record += bReq.TokenName
	record += bReq.RemoteAddress

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (bReq *BurningRequest) CalculateSize() uint64 {
	return calculateSize(bReq)
}
