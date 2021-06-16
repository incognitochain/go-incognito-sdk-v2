package metadata

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

type BurningRequest struct {
	BurnerAddress key.PaymentAddress
	BurningAmount uint64 // must be equal to vOut value
	TokenID       common.Hash
	TokenName     string
	RemoteAddress string
	MetadataBase
}

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

func (bReq *BurningRequest) CalculateSize() uint64 {
	return calculateSize(bReq)
}
