package bridge

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// whoever can send this type of tx
type BurningRequest struct {
	BurnerAddress key.PaymentAddress
	BurningAmount uint64 // must be equal to vout value
	TokenID       common.Hash
	TokenName     string
	RemoteAddress string
	metadataCommon.MetadataBase
}

func NewBurningRequest(
	burnerAddress key.PaymentAddress,
	burningAmount uint64,
	tokenID common.Hash,
	tokenName string,
	remoteAddress string,
	metaType int,
) (*BurningRequest, error) {
	metadataBase := metadataCommon.MetadataBase{
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

func (bReq BurningRequest) ValidateMetadataByItself() bool {
	return bReq.Type == metadataCommon.BurningRequestMeta || bReq.Type == metadataCommon.BurningForDepositToSCRequestMeta || bReq.Type == metadataCommon.BurningRequestMetaV2 ||
		bReq.Type == metadataCommon.BurningForDepositToSCRequestMetaV2 || bReq.Type == metadataCommon.BurningPBSCRequestMeta ||
		bReq.Type == metadataCommon.BurningPRVERC20RequestMeta || bReq.Type == metadataCommon.BurningPRVBEP20RequestMeta ||
		bReq.Type == metadataCommon.BurningPBSCForDepositToSCRequestMeta ||
		bReq.Type == metadataCommon.BurningPLGRequestMeta || bReq.Type == metadataCommon.BurningPLGForDepositToSCRequestMeta ||
		bReq.Type == metadataCommon.BurningFantomRequestMeta || bReq.Type == metadataCommon.BurningFantomForDepositToSCRequestMeta
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
	return metadataCommon.CalculateSize(bReq)
}
