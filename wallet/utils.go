package wallet

import (
	"bytes"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"

	"github.com/tyler-smith/go-bip39"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

// NewMnemonic generates a mnemonic string given the entropy bitSize.
func NewMnemonic(bitSize int) (string, error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}

// NewMnemonicFromEntropy generates a mnemonic string given the entropy
func NewMnemonicFromEntropy(entropy []byte) (string, error) {
	return bip39.NewMnemonic(entropy)
}

// NewSeedFromMnemonic creates a hashed seed output given a provided mnemonic string.
// The mnemonic is validated against the BIP39 standard before the seed is generated.
func NewSeedFromMnemonic(mnemonic string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("mnemonic is invalid")
	}

	return bip39.NewSeed(mnemonic, ""), nil
}

// GetBurningPublicKey returns the public key of the burning address.
func GetBurningPublicKey() []byte {
	// get burning address
	w, err := Base58CheckDeserialize(common.BurningAddress)
	if err != nil {
		return nil
	}

	return w.KeySet.PaymentAddress.Pk
}

// IsPublicKeyBurningAddress checks if a public key is a burning address in the Incognito network.
func IsPublicKeyBurningAddress(publicKey []byte) bool {
	// get burning address
	keyWalletBurningAdd1, err := Base58CheckDeserialize(common.BurningAddress)
	if err != nil {
		return false
	}
	if bytes.Equal(publicKey, keyWalletBurningAdd1.KeySet.PaymentAddress.Pk) {
		return true
	}

	keyWalletBurningAdd2, err := Base58CheckDeserialize(common.BurningAddress2)
	if err != nil {
		return false
	}
	if bytes.Equal(publicKey, keyWalletBurningAdd2.KeySet.PaymentAddress.Pk) {
		return true
	}

	return false
}

// GetPaymentAddressV1 retrieves the payment address ver 1 from the payment address ver 2.
//	- Payment Address V1 consists of: PK + TK
//	- Payment Address V2 consists of: PK + TK + PublicOTA
//
// If the input is a payment address ver 2, try to retrieve the corresponding payment address ver 1.
// Otherwise, return the input.
func GetPaymentAddressV1(addr string, isNewEncoding bool) (string, error) {
	newWallet, err := Base58CheckDeserialize(addr)
	if err != nil {
		return "", err
	}

	if len(newWallet.KeySet.PaymentAddress.Pk) == 0 || len(newWallet.KeySet.PaymentAddress.Pk) == 0 {
		return "", fmt.Errorf("something must be wrong with the provided payment address: %v", addr)
	}

	//Remove the publicOTA key and try to deserialize
	newWallet.KeySet.PaymentAddress.OTAPublic = nil

	if isNewEncoding {
		addrV1 := newWallet.Base58CheckSerialize(PaymentAddressType)
		if len(addrV1) == 0 {
			return "", fmt.Errorf("cannot decode new payment address: %v", addr)
		}

		return addrV1, nil
	} else {
		addr1InBytes, err := newWallet.Serialize(PaymentAddressType, false)
		if err != nil {
			return "", fmt.Errorf("cannot decode new payment address: %v", addr)
		}

		addrV1 := base58.Base58Check{}.Encode(addr1InBytes, common.ZeroByte)
		if len(addrV1) == 0 {
			return "", fmt.Errorf("cannot decode new payment address: %v", addr)
		}

		return addrV1, nil
	}
}

// ComparePaymentAddresses checks if two payment addresses are generated from the same private key.
//
// Just need to compare PKs and TKs.
func ComparePaymentAddresses(addr1, addr2 string) (bool, error) {
	//If these address strings are the same, just try to deserialize one of them
	if addr1 == addr2 {
		_, err := Base58CheckDeserialize(addr1)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	//If their lengths are the same, just compare the inputs
	keyWallet1, err := Base58CheckDeserialize(addr1)
	if err != nil {
		return false, err
	}

	keyWallet2, err := Base58CheckDeserialize(addr2)
	if err != nil {
		return false, err
	}

	pk1 := keyWallet1.KeySet.PaymentAddress.Pk
	tk1 := keyWallet1.KeySet.PaymentAddress.Tk

	pk2 := keyWallet2.KeySet.PaymentAddress.Pk
	tk2 := keyWallet2.KeySet.PaymentAddress.Tk

	if !bytes.Equal(pk1, pk2) {
		return false, fmt.Errorf("public keys mismatch: %v, %v", pk1, pk2)
	}

	if !bytes.Equal(tk1, tk2) {
		return false, fmt.Errorf("transmission keys mismatch: %v, %v", tk1, tk2)
	}

	return true, nil
}

// GenRandomWalletForShardID generates a random wallet for a specific shardID.
func GenRandomWalletForShardID(shardID byte) (*KeyWallet, error) {
	numTries := 100000
	for numTries > 0 {
		tmpWallet, err := NewMasterKeyFromSeed(common.RandBytes(32))
		if err != nil {
			return nil, err
		}

		pk := tmpWallet.KeySet.PaymentAddress.Pk

		lastByte := pk[len(pk)-1]
		if lastByte == shardID {
			return tmpWallet, nil
		}

		numTries--
	}

	return nil, fmt.Errorf("failed after %v tries", numTries)
}
