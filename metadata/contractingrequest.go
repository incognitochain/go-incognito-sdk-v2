package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// whoever can send this type of tx
type ContractingRequest struct {
	BurnerAddress key.PaymentAddress
	BurnedAmount  uint64 // must be equal to vout value
	TokenID       common.Hash
	MetadataBase
}

func (cReq ContractingRequest) Hash() *common.Hash {
	record := cReq.MetadataBase.Hash().String()
	record += cReq.BurnerAddress.String()
	record += cReq.TokenID.String()
	// TODO: @hung change to record += fmt.Sprint(cReq.BurnedAmount)
	record += string(cReq.BurnedAmount)

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (cReq *ContractingRequest) CalculateSize() uint64 {
	return calculateSize(cReq)
}
