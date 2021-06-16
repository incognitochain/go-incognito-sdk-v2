package key

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
)

// KeySet is real raw data of wallet account, which user can use to
// - spend and check double spend coin with private key
// - receive coin with payment address
// - read tx data with readonly key
type KeySet struct {
	PrivateKey     PrivateKey     //Master Private key
	PaymentAddress PaymentAddress //Payment address for sending coins
	ReadonlyKey    ViewingKey     //ViewingKey for retrieving the amount of coin (both V1 + V2) as well as the asset tag (V2 ONLY)
	OTAKey         OTAKey         //OTAKey is for recovering one time addresses: ONLY in V2
}

// GenerateKey generates key set from seed in byte array
func (keySet *KeySet) GenerateKey(seed []byte) *KeySet {
	keySet.PrivateKey = GeneratePrivateKey(seed)
	keySet.PaymentAddress = GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = GenerateViewingKey(keySet.PrivateKey[:])
	keySet.OTAKey = GenerateOTAKey(keySet.PrivateKey[:])
	return keySet
}

// InitFromPrivateKeyByte receives private key in bytes array,
// and regenerates payment address and readonly key
// returns error if private key is invalid
func (keySet *KeySet) InitFromPrivateKeyByte(privateKey []byte) error {
	if len(privateKey) != common.PrivateKeySize {
		return NewCacheError(InvalidPrivateKeyErr, nil)
	}

	keySet.PrivateKey = privateKey
	keySet.PaymentAddress = GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = GenerateViewingKey(keySet.PrivateKey[:])
	keySet.OTAKey = GenerateOTAKey(keySet.PrivateKey[:])
	return nil
}

// InitFromPrivateKey receives private key in PrivateKey type,
// and regenerates payment address and readonly key
// returns error if private key is invalid
func (keySet *KeySet) InitFromPrivateKey(privateKey *PrivateKey) error {
	if privateKey == nil || len(*privateKey) != common.PrivateKeySize {
		return NewCacheError(InvalidPrivateKeyErr, nil)
	}

	keySet.PrivateKey = *privateKey
	keySet.PaymentAddress = GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = GenerateViewingKey(keySet.PrivateKey[:])
	keySet.OTAKey = GenerateOTAKey(keySet.PrivateKey[:])

	return nil
}

// GetPublicKeyInBase58CheckEncode returns the public key which is base58 check encoded
func (keySet KeySet) GetPublicKeyInBase58CheckEncode() string {
	return base58.Base58Check{}.Encode(keySet.PaymentAddress.Pk, common.ZeroByte)
}

func (keySet KeySet) GetReadOnlyKeyInBase58CheckEncode() string {
	return base58.Base58Check{}.Encode(keySet.ReadonlyKey.Rk, common.ZeroByte)
}

func (keySet KeySet) GetOTASecretKeyInBase58CheckEncode() string {
	return base58.Base58Check{}.Encode(keySet.OTAKey.GetOTASecretKey().ToBytesS(), common.ZeroByte)
}
