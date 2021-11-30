package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type WithdrawalLPFeeResponse struct {
	metadataCommon.MetadataBase
	ReqTxID common.Hash `json:"ReqTxID"`
}

func (withdrawalResponse WithdrawalLPFeeResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(withdrawalResponse)
	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (withdrawalResponse *WithdrawalLPFeeResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(withdrawalResponse)
}
