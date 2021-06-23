package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/pkg/errors"
	"strconv"
)

// WithDrawRewardRequest is a request to withdraw staking rewards of a user.
// The user needs to sign this request to make sure he/she is authorized to withdraw the rewards.
type WithDrawRewardRequest struct {
	MetadataBaseWithSignature
	PaymentAddress key.PaymentAddress
	TokenID        common.Hash
	Version        int
}

// NewWithDrawRewardRequest creates a new WithDrawRewardRequest.
func NewWithDrawRewardRequest(tokenIDStr string, paymentAddStr string, version float64, metaType int) (*WithDrawRewardRequest, error) {
	metadataBase := NewMetadataBaseWithSignature(metaType)
	tokenID, err := common.Hash{}.NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, fmt.Errorf("token ID is invalid")
	}
	paymentAddWallet, err := wallet.Base58CheckDeserialize(paymentAddStr)
	if err != nil {
		return nil, fmt.Errorf("payment address is invalid")
	}
	ok, err := common.SliceExists(AcceptedWithdrawRewardRequestVersion, int(version))
	if !ok || err != nil {
		return nil, errors.Errorf("Invalid version %v", version)
	}

	withdrawRewardRequest := &WithDrawRewardRequest{
		MetadataBaseWithSignature: *metadataBase,
		TokenID:                   *tokenID,
		PaymentAddress:            paymentAddWallet.KeySet.PaymentAddress,
		Version:                   int(version),
	}
	return withdrawRewardRequest, nil
}

// UnmarshalJSON does the JSON-unmarshalling operation for a WithDrawRewardRequest.
func (req *WithDrawRewardRequest) UnmarshalJSON(data []byte) error {
	tmp := &struct {
		MetadataBase
		PaymentAddress key.PaymentAddress
		TokenID        common.Hash
		Version        int
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if tmp.PaymentAddress.Pk == nil && tmp.PaymentAddress.Tk == nil {
		tmpOld := &struct {
			MetadataBase
			key.PaymentAddress
			TokenID common.Hash
			Version int
		}{}
		if err := json.Unmarshal(data, &tmpOld); err != nil {
			return err
		}

		tmp.PaymentAddress.Tk = tmpOld.Tk
		tmp.PaymentAddress.Pk = tmpOld.Pk
	}

	req.MetadataBase = tmp.MetadataBase
	req.PaymentAddress = tmp.PaymentAddress
	req.TokenID = tmp.TokenID
	req.Version = tmp.Version
	return nil
}

// Hash overrides MetadataBase.Hash().
func (req WithDrawRewardRequest) Hash() *common.Hash {
	if req.Version == 1 {
		bArr := append(req.PaymentAddress.Bytes(), req.TokenID.GetBytes()...)
		if req.Sig != nil && len(req.Sig) != 0 {
			bArr = append(bArr, req.Sig...)
		}
		txReqHash := common.HashH(bArr)
		return &txReqHash
	} else {
		record := strconv.Itoa(req.Type)
		data := []byte(record)
		hash := common.HashH(data)
		return &hash
	}
}

// HashWithoutSig overrides MetadataBase.HashWithoutSig().
func (req WithDrawRewardRequest) HashWithoutSig() *common.Hash {
	if req.Version == 1 {
		bArr := append(req.PaymentAddress.Bytes(), req.TokenID.GetBytes()...)
		txReqHash := common.HashH(bArr)
		return &txReqHash
	} else {
		record := strconv.Itoa(req.Type)
		data := []byte(record)
		hash := common.HashH(data)
		return &hash
	}
}

// ShouldSignMetaData returns true.
func (*WithDrawRewardRequest) ShouldSignMetaData() bool { return true }

// GetType overrides MetadataBase.GetType().
func (req WithDrawRewardRequest) GetType() int {
	return req.Type
}

// CalculateSize overrides MetadataBase.CalculateSize().
func (req *WithDrawRewardRequest) CalculateSize() uint64 {
	return calculateSize(req)
}
