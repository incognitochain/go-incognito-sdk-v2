package pdexv3

import (
	"encoding/json"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type WithdrawalProtocolFeeResponse struct {
	metadataCommon.MetadataBase
	Address      string      `json:"Address"`
	TokenID      common.Hash `json:"TokenID"`
	Amount       uint64      `json:"Amount"`
	ReqTxID      common.Hash `json:"ReqTxID"`
	SharedRandom []byte      `json:"SharedRandom"`
}

func (withdrawalResponse WithdrawalProtocolFeeResponse) Hash() *common.Hash {
	rawBytes, _ := json.Marshal(struct {
		Type    int         `json:"Type"`
		Address string      `json:"Address"`
		TokenID common.Hash `json:"TokenID"`
		Amount  uint64      `json:"Amount"`
		ReqTxID common.Hash `json:"ReqTxID"`
	}{
		Type:    metadataCommon.Pdexv3WithdrawProtocolFeeResponseMeta,
		Address: withdrawalResponse.Address,
		TokenID: withdrawalResponse.TokenID,
		Amount:  withdrawalResponse.Amount,
		ReqTxID: withdrawalResponse.ReqTxID,
	})

	hash := common.HashH([]byte(rawBytes))
	return &hash
}

func (withdrawalResponse *WithdrawalProtocolFeeResponse) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(withdrawalResponse)
}

func (withdrawalResponse *WithdrawalProtocolFeeResponse) SetSharedRandom(r []byte) {
	withdrawalResponse.SharedRandom = r
}
