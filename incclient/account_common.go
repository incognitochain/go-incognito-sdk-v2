package incclient

import (
	"encoding/json"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// KeyInfo contains all key-related information of an account.
type KeyInfo struct {
	PrivateKey         string
	PublicKey          string
	PaymentAddressV1   string
	PaymentAddress     string
	ReadOnlyKey        string
	OTAPrivateKey      string
	MiningKey          string
	MiningPublicKey    string
	ValidatorPublicKey string
	ShardID            byte
}

func (k KeyInfo) String() string {
	jsb, err := json.MarshalIndent(k, "", "\t")
	if err != nil {
		return ""
	}

	return string(jsb)
}

// GetAccountInfoFromPrivateKey returns all fields related to a private key.
func GetAccountInfoFromPrivateKey(privateKey string) (*KeyInfo, error) {
	w, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, err
	}

	if len(w.KeySet.PrivateKey) != 32 {
		return nil, fmt.Errorf("privateKey is invalid")
	}

	pubKey := PrivateKeyToPublicKey(privateKey)
	addr := PrivateKeyToPaymentAddress(privateKey, -1)
	addrV1 := PrivateKeyToPaymentAddress(privateKey, 0)
	readonlyKey := PrivateKeyToReadonlyKey(privateKey)
	otaKey := PrivateKeyToPrivateOTAKey(privateKey)
	miningKey := PrivateKeyToMiningKey(privateKey)
	shardID := GetShardIDFromPrivateKey(privateKey)

	miningKeyBytes, _, err := base58.Base58Check{}.Decode(miningKey)
	if err != nil {
		return nil, err
	}

	committeeKey, err := key.NewCommitteeKeyFromSeed(miningKeyBytes, pubKey)
	if err != nil {
		return nil, err
	}
	miningPubKeyStr, err := committeeKey.ToBase58()
	if err != nil {
		return nil, err
	}
	validatorPublicKey := committeeKey.GetMiningKeyBase58(common.BlsConsensus)

	return &KeyInfo{
		PrivateKey:         privateKey,
		PublicKey:          base58.Base58Check{}.Encode(pubKey, common.ZeroByte),
		PaymentAddressV1:   addrV1,
		PaymentAddress:     addr,
		ReadOnlyKey:        readonlyKey,
		OTAPrivateKey:      otaKey,
		MiningKey:          miningKey,
		MiningPublicKey:    miningPubKeyStr,
		ValidatorPublicKey: validatorPublicKey,
		ShardID:            shardID,
	}, nil
}
