package metadata

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/pkg/errors"
	"strconv"
)

type WithDrawRewardRequest struct {
	MetadataBaseWithSignature
	PaymentAddress key.PaymentAddress
	TokenID        common.Hash
	Version        int
}

func (withdrawRequestMetadata *WithDrawRewardRequest) UnmarshalJSON(data []byte) error {
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

	withdrawRequestMetadata.MetadataBase = tmp.MetadataBase
	withdrawRequestMetadata.PaymentAddress = tmp.PaymentAddress
	withdrawRequestMetadata.TokenID = tmp.TokenID
	withdrawRequestMetadata.Version = tmp.Version
	return nil
}

func (withDrawRewardRequest WithDrawRewardRequest) Hash() *common.Hash {
	if withDrawRewardRequest.Version == 1 {
		bArr := append(withDrawRewardRequest.PaymentAddress.Bytes(), withDrawRewardRequest.TokenID.GetBytes()...)
		if withDrawRewardRequest.Sig != nil && len(withDrawRewardRequest.Sig) != 0 {
			bArr = append(bArr, withDrawRewardRequest.Sig...)
		}
		txReqHash := common.HashH(bArr)
		return &txReqHash
	} else {
		record := strconv.Itoa(withDrawRewardRequest.Type)
		data := []byte(record)
		hash := common.HashH(data)
		return &hash
	}
}

func (withDrawRewardRequest WithDrawRewardRequest) HashWithoutSig() *common.Hash {
	if withDrawRewardRequest.Version == 1 {
		bArr := append(withDrawRewardRequest.PaymentAddress.Bytes(), withDrawRewardRequest.TokenID.GetBytes()...)
		txReqHash := common.HashH(bArr)
		return &txReqHash
	} else {
		record := strconv.Itoa(withDrawRewardRequest.Type)
		data := []byte(record)
		hash := common.HashH(data)
		return &hash
	}
}

func (*WithDrawRewardRequest) ShouldSignMetaData() bool { return true }

func NewWithDrawRewardRequest(tokenIDStr string, paymentAddStr string, version float64, metaType int) (*WithDrawRewardRequest, error) {
	metadataBase := NewMetadataBaseWithSignature(metaType)
	tokenID, err := common.Hash{}.NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, errors.New("token ID is invalid")
	}
	paymentAddWallet, err := wallet.Base58CheckDeserialize(paymentAddStr)
	if err != nil {
		return nil, errors.New("payment address is invalid")
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

func (withDrawRewardRequest WithDrawRewardRequest) GetType() int {
	return withDrawRewardRequest.Type
}

func (withDrawRewardRequest *WithDrawRewardRequest) CalculateSize() uint64 {
	return calculateSize(withDrawRewardRequest)
}
