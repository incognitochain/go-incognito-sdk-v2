package wallet

import (
	"bytes"
	"fmt"

	"github.com/wemeetagain/go-hdwallet"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// KeyWallet defines the components of a hierarchical deterministic wallet.
type KeyWallet struct {
	Depth       byte               // 1 bytes
	ChildNumber []byte             // 4 bytes
	ChainCode   []byte             // 32 bytes
	HDKey       *hdwallet.HDWallet // Actual HD wallet key
	KeySet      key.KeySet         // Incognito key-set
}

// NewMasterKey returns a new KeyWallet, and a mnemonic that generates the KeyWallet.
func NewMasterKey() (*KeyWallet, string, error) {
	mnemonic, err := NewMnemonic(128)
	if err != nil {
		return nil, mnemonic, err
	}

	seed, err := NewSeedFromMnemonic(mnemonic)
	if err != nil {
		return nil, mnemonic, err
	}

	k := hdwallet.MasterKey(seed)

	keySet := (&key.KeySet{}).GenerateKey(k.Key)

	w := &KeyWallet{
		HDKey:       k,
		ChainCode:   k.Chaincode,
		KeySet:      *keySet,
		Depth:       0x00,
		ChildNumber: []byte{0x00, 0x00, 0x00, 0x00},
	}

	return w, mnemonic, nil
}

// NewMasterKeyFromSeed returns a new KeyWallet given a random seed.
func NewMasterKeyFromSeed(seed []byte) (*KeyWallet, error) {
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, fmt.Errorf("invalid seed length")
	}

	k := hdwallet.MasterKey(seed)

	keySet := (&key.KeySet{}).GenerateKey(k.Key)

	w := &KeyWallet{
		HDKey:       k,
		ChainCode:   k.Chaincode,
		KeySet:      *keySet,
		Depth:       0x00,
		ChildNumber: []byte{0x00, 0x00, 0x00, 0x00},
	}

	return w, nil
}

// NewMasterKeyFromMnemonic returns a new KeyWallet given a BIP39 mnemonic string.
func NewMasterKeyFromMnemonic(mnemonic string) (*KeyWallet, error) {
	seed, err := NewSeedFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return NewMasterKeyFromSeed(seed)
}

// DeriveChild returns the i-th child of wallet w following the BIP-44 standard.
// Call this function with increasing i to create as many wallets as you want.
func (w *KeyWallet) DeriveChild(i uint32) (*KeyWallet, error) {
	if w.HDKey == nil {
		return nil, fmt.Errorf("cannot dereive child key: HDKey not found")
	}

	purposeKey, err := w.HDKey.Child(HardenedKeyZeroIndex + BIP44Purpose)
	if err != nil {
		return nil, err
	}

	coinTypeKey, err := purposeKey.Child(HardenedKeyZeroIndex + Bip44CoinType)
	if err != nil {
		return nil, err
	}

	accountKey, err := coinTypeKey.Child(HardenedKeyZeroIndex)
	if err != nil {
		return nil, err
	}

	changeKey, err := accountKey.Child(0)
	if err != nil {
		return nil, err
	}

	childKey, err := changeKey.Child(i)
	if err != nil {
		return nil, err
	}

	childKeySet := new(key.KeySet)
	childKeySet.GenerateKey(childKey.Key[1:]) // the first byte is 0x00, so we skip

	// Create Child KeySet with data common to all both scenarios
	childWallet := &KeyWallet{
		HDKey:       childKey,
		ChildNumber: common.Uint32ToBytes(i),
		ChainCode:   childKey.Chaincode,
		Depth:       w.Depth + 1,
		KeySet:      *childKeySet,
	}

	return childWallet, nil
}

// Serialize receives keyType and serializes key which has keyType to bytes array
// and append 4-byte checksum into bytes array
func (w *KeyWallet) Serialize(keyType byte, isNewCheckSum bool) ([]byte, error) {
	// Write fields to buffer in order
	buffer := new(bytes.Buffer)
	buffer.WriteByte(keyType)
	if keyType == PrivateKeyType {
		buffer.WriteByte(w.Depth)
		buffer.Write(w.ChildNumber)
		buffer.Write(w.ChainCode)

		// Private keys should be prepended with a single null byte
		keyBytes := make([]byte, 0)
		keyBytes = append(keyBytes, byte(len(w.KeySet.PrivateKey))) // set length
		keyBytes = append(keyBytes, w.KeySet.PrivateKey[:]...)      // set pri-w
		buffer.Write(keyBytes)
	} else if keyType == PaymentAddressType {
		keyBytes := make([]byte, 0)
		keyBytes = append(keyBytes, byte(len(w.KeySet.PaymentAddress.Pk)))
		keyBytes = append(keyBytes, w.KeySet.PaymentAddress.Pk[:]...)

		keyBytes = append(keyBytes, byte(len(w.KeySet.PaymentAddress.Tk)))
		keyBytes = append(keyBytes, w.KeySet.PaymentAddress.Tk[:]...)

		if isNewCheckSum && len(w.KeySet.PaymentAddress.OTAPublic) > 0 {
			keyBytes = append(keyBytes, byte(len(w.KeySet.PaymentAddress.OTAPublic)))
			keyBytes = append(keyBytes, w.KeySet.PaymentAddress.OTAPublic[:]...)
		}

		buffer.Write(keyBytes)
	} else if keyType == ReadonlyKeyType {
		keyBytes := make([]byte, 0)
		keyBytes = append(keyBytes, byte(len(w.KeySet.ReadonlyKey.Pk)))
		keyBytes = append(keyBytes, w.KeySet.ReadonlyKey.Pk[:]...)

		keyBytes = append(keyBytes, byte(len(w.KeySet.ReadonlyKey.Rk)))
		keyBytes = append(keyBytes, w.KeySet.ReadonlyKey.Rk[:]...)
		buffer.Write(keyBytes)
	} else if keyType == OTAKeyType {
		keyBytes := make([]byte, 0)
		keyBytes = append(keyBytes, byte(len(w.KeySet.OTAKey.GetPublicSpend().ToBytesS())))
		keyBytes = append(keyBytes, w.KeySet.OTAKey.GetPublicSpend().ToBytesS()[:]...)

		keyBytes = append(keyBytes, byte(len(w.KeySet.OTAKey.GetOTASecretKey().ToBytesS())))
		keyBytes = append(keyBytes, w.KeySet.OTAKey.GetOTASecretKey().ToBytesS()[:]...)
		buffer.Write(keyBytes)
	} else {
		return []byte{}, NewWalletError(InvalidKeyTypeErr, nil)
	}

	checkSum := base58.ChecksumFirst4Bytes(buffer.Bytes(), isNewCheckSum)

	serializedKey := append(buffer.Bytes(), checkSum...)
	return serializedKey, nil
}

// Base58CheckSerialize encodes the key corresponding to keyType in KeySet
// in the standard Incognito base58 encoding
// It returns the encoding string of the key
func (w *KeyWallet) Base58CheckSerialize(keyType byte) string {
	isNewEncoding := common.AddressVersion == 1
	serializedKey, err := w.Serialize(keyType, isNewEncoding) //Must use the new checksum from now on
	if err != nil {
		return ""
	}

	if isNewEncoding {
		return base58.Base58Check{}.NewEncode(serializedKey, 0)
	}
	return base58.Base58Check{}.Encode(serializedKey, 0) //Must use the new encoding algorithm from now on
}

// Deserialize receives a byte array and deserializes into KeySet
// because data contains keyType and serialized data of corresponding key
// it returns KeySet just contain corresponding key
func deserialize(data []byte) (*KeyWallet, error) {
	var k = &KeyWallet{}
	if len(data) < 2 {
		return nil, NewWalletError(InvalidKeyTypeErr, nil)
	}
	keyType := data[0]
	if keyType == PrivateKeyType {
		if len(data) != privateKeySerializedBytesLen {
			return nil, NewWalletError(InvalidSeserializedKey, nil)
		}

		k.Depth = data[1]
		k.ChildNumber = data[2:6]
		k.ChainCode = data[6:38]
		keyLength := int(data[38])
		k.KeySet.PrivateKey = make([]byte, keyLength)
		copy(k.KeySet.PrivateKey[:], data[39:39+keyLength])
		err := k.KeySet.InitFromPrivateKey(&k.KeySet.PrivateKey)
		if err != nil {
			return nil, err
		}
	} else if keyType == PaymentAddressType {
		if !bytes.Equal(burnAddress1BytesDecode, data) {
			if len(data) != paymentAddrSerializedBytesLen && len(data) != paymentAddrSerializedBytesLen+1+crypto.Ed25519KeySize {
				return nil, NewWalletError(InvalidSeserializedKey, fmt.Errorf("length ota public k not valid: %v", len(data)))
			}
		}
		apkKeyLength := int(data[1])
		publicViewKeyLength := int(data[apkKeyLength+2])
		k.KeySet.PaymentAddress.Pk = make([]byte, apkKeyLength)
		k.KeySet.PaymentAddress.Tk = make([]byte, publicViewKeyLength)
		copy(k.KeySet.PaymentAddress.Pk[:], data[2:2+apkKeyLength])
		copy(k.KeySet.PaymentAddress.Tk[:], data[3+apkKeyLength:3+apkKeyLength+publicViewKeyLength])

		//Deserialize OTAPublic Key
		if len(data) > paymentAddrSerializedBytesLen {
			publicOTALength := int(data[apkKeyLength+publicViewKeyLength+3])
			if publicOTALength != crypto.Ed25519KeySize {
				return nil, NewWalletError(InvalidSeserializedKey, fmt.Errorf("length ota public k not valid: %v", publicOTALength))
			}
			k.KeySet.PaymentAddress.OTAPublic = append([]byte{}, data[apkKeyLength+publicViewKeyLength+4:apkKeyLength+publicViewKeyLength+publicOTALength+4]...)
		}

	} else if keyType == ReadonlyKeyType {
		if len(data) != readOnlyKeySerializedBytesLen {
			return nil, NewWalletError(InvalidSeserializedKey, nil)
		}

		apkKeyLength := int(data[1])
		if len(data) < apkKeyLength+3 {
			return nil, NewWalletError(InvalidKeyTypeErr, nil)
		}

		publicViewKeyLength := int(data[apkKeyLength+2])
		k.KeySet.ReadonlyKey.Pk = make([]byte, apkKeyLength)
		k.KeySet.ReadonlyKey.Rk = make([]byte, publicViewKeyLength)
		copy(k.KeySet.ReadonlyKey.Pk[:], data[2:2+apkKeyLength])
		copy(k.KeySet.ReadonlyKey.Rk[:], data[3+apkKeyLength:3+apkKeyLength+publicViewKeyLength])
	} else if keyType == OTAKeyType {
		if len(data) != otaKeySerializedBytesLen {
			return nil, NewWalletError(InvalidSeserializedKey, nil)
		}

		pkKeyLength := int(data[1])
		if len(data) < pkKeyLength+3 {
			return nil, NewWalletError(InvalidKeyTypeErr, nil)
		}
		skKeyLength := int(data[pkKeyLength+2])

		k.KeySet.OTAKey.SetPublicSpend(data[2 : 2+pkKeyLength])
		k.KeySet.OTAKey.SetOTASecretKey(data[3+pkKeyLength : 3+pkKeyLength+skKeyLength])
	} else {
		return nil, NewWalletError(InvalidKeyTypeErr, fmt.Errorf("cannot detect k type"))
	}

	// validate checksum: allowing both new- and old-encoded strings
	// try to verify in the new way first
	cs1 := base58.ChecksumFirst4Bytes(data[0:len(data)-4], true)
	cs2 := data[len(data)-4:]
	if !bytes.Equal(cs1, cs2) { // try to compare old checksum
		oldCS1 := base58.ChecksumFirst4Bytes(data[0:len(data)-4], false)
		if !bytes.Equal(oldCS1, cs2) {
			return nil, NewWalletError(InvalidChecksumErr, nil)
		}
	}

	return k, nil
}

// Base58CheckDeserialize deserializes the keySet of a KeyWallet encoded in base58
// because data contains keyType and serialized data of corresponding key
// it returns KeySet just contain corresponding key
func Base58CheckDeserialize(data string) (*KeyWallet, error) {
	b, _, err := base58.Base58Check{}.Decode(data)
	if err != nil {
		return nil, err
	}
	return deserialize(b)
}

// GetPrivateKey returns the base58-encoded private key of a KeyWallet.
func (w *KeyWallet) GetPrivateKey() (string, error) {
	var privateKey string
	if len(w.KeySet.PrivateKey) == 0 {
		return privateKey, fmt.Errorf("private key not found")
	}

	privateKey = w.Base58CheckSerialize(PrivateKeyType)
	return privateKey, nil
}

// GetPublicKey returns the base58-encoded public key of a KeyWallet.
func (w *KeyWallet) GetPublicKey() (string, error) {
	var pubKey string
	if len(w.KeySet.PaymentAddress.Pk) == 0 {
		return pubKey, fmt.Errorf("public key not found")
	}

	pubKey = base58.Base58Check{}.NewEncode(w.KeySet.PaymentAddress.Pk, common.ZeroByte)
	return pubKey, nil
}

// GetPaymentAddress returns the base58-encoded payment address of a KeyWallet.
func (w *KeyWallet) GetPaymentAddress() (string, error) {
	var addr string
	if len(w.KeySet.PaymentAddress.Bytes()) == 0 {
		return addr, fmt.Errorf("payment address not found")
	}

	addr = w.Base58CheckSerialize(PaymentAddressType)
	return addr, nil
}

// GetReadonlyKey returns the base58-encoded readonly key of a KeyWallet.
func (w *KeyWallet) GetReadonlyKey() (string, error) {
	var readonlyKey string
	if w.KeySet.ReadonlyKey.GetPublicSpend() == nil || w.KeySet.ReadonlyKey.GetPrivateView() == nil {
		return readonlyKey, fmt.Errorf("read-only key not found")
	}

	readonlyKey = w.Base58CheckSerialize(ReadonlyKeyType)
	return readonlyKey, nil
}

// GetOTAPrivateKey returns the base58-encoded privateOTA key of a KeyWallet.
func (w *KeyWallet) GetOTAPrivateKey() (string, error) {
	var readonlyKey string
	if w.KeySet.OTAKey.GetOTASecretKey() == nil || w.KeySet.OTAKey.GetPublicSpend() == nil {
		return readonlyKey, fmt.Errorf("privateOTA key not found")
	}

	readonlyKey = w.Base58CheckSerialize(OTAKeyType)
	return readonlyKey, nil
}
