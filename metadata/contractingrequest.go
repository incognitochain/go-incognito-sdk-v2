package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// ContractingRequest is a request to burn centralized UTXOs (mostly to un-shield).
// Anyone can send this request.
type ContractingRequest struct {
	BurnerAddress key.PaymentAddress
	BurnedAmount  uint64
	TokenID       common.Hash
	MetadataBase
}

// NewContractingRequest creates a new ContractingRequest.
func NewContractingRequest(
	burnerAddress key.PaymentAddress,
	burnedAmount uint64,
	tokenID common.Hash,
	metaType int,
) (*ContractingRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	contractingReq := &ContractingRequest{
		TokenID:       tokenID,
		BurnedAmount:  burnedAmount,
		BurnerAddress: burnerAddress,
	}
	contractingReq.MetadataBase = metadataBase
	return contractingReq, nil
}

// Hash overrides MetadataBase.Hash().
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

// CalculateSize overrides MetadataBase.CalculateSize().
func (cReq *ContractingRequest) CalculateSize() uint64 {
	return calculateSize(cReq)
}
