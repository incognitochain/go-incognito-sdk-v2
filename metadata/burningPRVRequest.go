package metadata

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type BurningPRVRequest struct {
	BurnerAddress     key.PaymentAddress // unused
	BurningAmount     uint64             // must be equal to vout value
	TokenID           common.Hash
	TokenName         string // unused
	RemoteAddress     string
	RedepositReceiver *coin.OTAReceiver `json:"RedepositReceiver,omitempty"`
	metadataCommon.MetadataBase
}

func NewBurningPRVRequest(
	burnerAddress key.PaymentAddress,
	burningAmount uint64,
	tokenID common.Hash,
	tokenName string,
	remoteAddress string,
	redepositReceiver coin.OTAReceiver,
	metaType int,
) (*BurningPRVRequest, error) {
	metadataBase := metadataCommon.MetadataBase{
		Type: metaType,
	}
	burningReq := &BurningPRVRequest{
		BurnerAddress:     burnerAddress,
		BurningAmount:     burningAmount,
		TokenID:           tokenID,
		TokenName:         tokenName,
		RemoteAddress:     remoteAddress,
		RedepositReceiver: &redepositReceiver,
	}
	burningReq.MetadataBase = metadataBase
	return burningReq, nil
}

func (bReq BurningPRVRequest) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(bReq)
	hash := common.HashH(rawBytes)
	return &hash
}

func (bReq *BurningPRVRequest) GetOTADeclarations() []metadataCommon.OTADeclaration {
	var result []metadataCommon.OTADeclaration
	result = append(result, metadataCommon.OTADeclaration{
		PublicKey: bReq.RedepositReceiver.PublicKey.ToBytes(), TokenID: common.PRVCoinID,
	})
	return result
}
