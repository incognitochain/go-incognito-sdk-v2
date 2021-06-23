package key

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// KeySet consists of the following fields
// - PrivateKey: used to spend UTXOs.
// - PaymentAddress: used to receive UTXOs.
// - ReadonlyKey: used to decrypt UTXOs.
// - OTAKey: used to check the owner of UTXOs.
type KeySet struct {
	PrivateKey     PrivateKey     //Master Private key
	PaymentAddress PaymentAddress //Payment address for sending coins
	ReadonlyKey    ViewingKey     //ViewingKey for retrieving the amount of coin (both V1 + V2) as well as the asset tag (V2 ONLY)
	OTAKey         OTAKey         //OTAKey is for recovering one time addresses: ONLY in V2
}

// GenerateKey generates key set from seed in byte array.
func (keySet *KeySet) GenerateKey(seed []byte) *KeySet {
	keySet.PrivateKey = GeneratePrivateKey(seed)
	keySet.PaymentAddress = GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = GenerateViewingKey(keySet.PrivateKey[:])
	keySet.OTAKey = GenerateOTAKey(keySet.PrivateKey[:])
	return keySet
}

// InitFromPrivateKeyByte receives private key in bytes array,
// and re-generates its payment address and other related keys.
// It returns an Error if the private key is invalid.
func (keySet *KeySet) InitFromPrivateKeyByte(privateKey []byte) error {
	if len(privateKey) != common.PrivateKeySize {
		return NewError(InvalidPrivateKeyErr, nil)
	}

	keySet.PrivateKey = privateKey
	keySet.PaymentAddress = GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = GenerateViewingKey(keySet.PrivateKey[:])
	keySet.OTAKey = GenerateOTAKey(keySet.PrivateKey[:])
	return nil
}

// InitFromPrivateKey receives private key in PrivateKey type,
// and re-generates the payment address and other related keys.
// It returns an Error if private key is invalid.
func (keySet *KeySet) InitFromPrivateKey(privateKey *PrivateKey) error {
	if privateKey == nil || len(*privateKey) != common.PrivateKeySize {
		return NewError(InvalidPrivateKeyErr, nil)
	}

	keySet.PrivateKey = *privateKey
	keySet.PaymentAddress = GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = GenerateViewingKey(keySet.PrivateKey[:])
	keySet.OTAKey = GenerateOTAKey(keySet.PrivateKey[:])

	return nil
}
