package bridge

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
	"github.com/pkg/errors"
)

type IssuingEVMAuroraRequest struct {
	TxHash     common.Hash
	IncTokenID common.Hash
	NetworkID  uint `json:"NetworkID,omitempty"`
	metadataCommon.MetadataBase
}

func NewIssuingEVMAuroraRequest(
	txHash common.Hash,
	incTokenId common.Hash,
	networkID uint,
	metaType int,
) (*IssuingEVMAuroraRequest, error) {
	metadataBase := metadataCommon.MetadataBase{
		Type: metaType,
	}
	issuingEVMReq := &IssuingEVMAuroraRequest{
		TxHash:     txHash,
		IncTokenID: incTokenId,
		NetworkID:  networkID,
	}
	issuingEVMReq.MetadataBase = metadataBase
	return issuingEVMReq, nil
}

func NewIssuingEVMAuroraRequestFromMap(
	data map[string]interface{},
	networkID uint,
	metatype int,
) (*IssuingEVMAuroraRequest, error) {
	txHash, err := common.Hash{}.NewHashFromStr(data["TxHash"].(string))
	if err != nil {
		return nil, errors.Errorf("TxHash incorrect")
	}

	incTokenID, err := common.Hash{}.NewHashFromStr(data["IncTokenID"].(string))
	if err != nil {
		return nil, errors.Errorf("TokenID incorrect")
	}

	req, _ := NewIssuingEVMAuroraRequest(
		*txHash,
		*incTokenID,
		networkID,
		metatype,
	)
	return req, nil
}

func (iReq IssuingEVMAuroraRequest) Hash() *common.Hash {
	record := iReq.MetadataBase.Hash().String()
	record += iReq.TxHash.String()
	record += iReq.IncTokenID.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}
