package hybridencryption

import (
	"encoding/json"
	"errors"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/utils"
)

// HybridCipherText represents a cipher text for the hybrid encryption scheme.
// The hybrid encryption schemes uses the AES scheme to encrypt a message with arbitrary size,
// and uses Elgamal encryption to encrypt the AES key.
type HybridCipherText struct {
	msgEncrypted    []byte
	symKeyEncrypted []byte
}

// GetMsgEncrypted returns the encrypted message of a HybridCipherText.
func (c HybridCipherText) GetMsgEncrypted() []byte {
	return c.msgEncrypted
}

// GetSymKeyEncrypted returns the encrypted key of a HybridCipherText.
func (c HybridCipherText) GetSymKeyEncrypted() []byte {
	return c.symKeyEncrypted
}

// IsNil checks if a HybridCipherText is empty.
func (c HybridCipherText) IsNil() bool {
	if len(c.msgEncrypted) == 0 {
		return true
	}

	return len(c.symKeyEncrypted) == 0
}

// MarshalJSON returns the JSON-marshalled data of a HybridCipherText.
func (c HybridCipherText) MarshalJSON() ([]byte, error) {
	data := c.Bytes()
	temp := base58.Base58Check{}.Encode(data, common.ZeroByte)
	return json.Marshal(temp)
}

// UnmarshalJSON un-marshals raw-data into a HybridCipherText.
func (c *HybridCipherText) UnmarshalJSON(data []byte) error {
	dataStr := ""
	_ = json.Unmarshal(data, &dataStr)
	temp, _, err := base58.Base58Check{}.Decode(dataStr)
	if err != nil {
		return err
	}
	c.SetBytes(temp)
	return nil
}

// Bytes converts ciphertext to bytes array.
// If ciphertext is nil, return empty byte array.
func (c HybridCipherText) Bytes() []byte {
	if c.IsNil() {
		return []byte{}
	}

	res := make([]byte, 0)
	res = append(res, c.symKeyEncrypted...)
	res = append(res, c.msgEncrypted...)

	return res
}

// SetBytes sets byte-representation data into a HybridCipherText.
func (c *HybridCipherText) SetBytes(bytes []byte) error {
	if len(bytes) == 0 {
		return utils.NewPrivacyErr(utils.InvalidInputToSetBytesErr, nil)
	}

	if len(bytes) < elGamalCiphertextSize {
		// out of range
		return errors.New("out of range Parse c")
	}
	c.symKeyEncrypted = bytes[0:elGamalCiphertextSize]
	c.msgEncrypted = bytes[elGamalCiphertextSize:]
	return nil
}

// HybridEncrypt encrypts a message with arbitrary size, and returns a HybridCipherText.
// HybridEncrypt generates an AES key by getting the X-coordinate of a randomized elliptic point.
// It then uses this AES key to encrypt the message. The key is then encrypted using the ElGamal encryption scheme
// with the given publicKey.
func HybridEncrypt(msg []byte, publicKey *crypto.Point) (ciphertext *HybridCipherText, err error) {
	ciphertext = new(HybridCipherText)

	// Generate a AES key bytes
	sKeyPoint := crypto.RandomPoint()
	sKeyByte := sKeyPoint.ToBytes()
	// Encrypt msg using aesKeyByte

	aesKey := sKeyByte[:]
	aesScheme := &common.AES{
		Key: aesKey,
	}
	ciphertext.msgEncrypted, err = aesScheme.Encrypt(msg)
	if err != nil {
		return nil, err
	}

	// Using ElGamal cryptosystem for encrypting AES sym key
	pubKey := new(elGamalPublicKey)
	pubKey.h = publicKey
	ciphertext.symKeyEncrypted = pubKey.encrypt(sKeyPoint).Bytes()

	return ciphertext, nil
}

// HybridDecrypt returns the message by decrypting the given HybridCipherText using the given privateKey.
func HybridDecrypt(ciphertext *HybridCipherText, privateKey *crypto.Scalar) (msg []byte, err error) {
	// Validate ciphertext
	if ciphertext.IsNil() {
		return []byte{}, errors.New("ciphertext must not be nil")
	}

	// Get receiving key, which is a private key of ElGamal cryptosystem
	privKey := new(elGamalPrivateKey)
	privKey.set(privateKey)

	// Parse encrypted AES key encoded as an elliptic point from EncryptedSymKey
	encryptedAESKey := new(elGamalCipherText)
	err = encryptedAESKey.SetBytes(ciphertext.symKeyEncrypted)
	if err != nil {
		return []byte{}, err
	}

	// Decrypt encryptedAESKey using recipient's receiving key
	aesKeyPoint, err := privKey.decrypt(encryptedAESKey)
	if err != nil {
		return []byte{}, err
	}

	// Get AES key
	aesKeyByte := aesKeyPoint.ToBytes()
	aesKey := aesKeyByte[:]
	aesScheme := &common.AES{
		Key: aesKey,
	}

	// Decrypt encrypted coin randomness using AES keysatt
	msg, err = aesScheme.Decrypt(ciphertext.msgEncrypted)
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}
