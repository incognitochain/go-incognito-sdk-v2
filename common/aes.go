package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

var plainTextIsEmptyErr = fmt.Errorf("plaintext is empty")
var cipherTextIsEmptyErr = fmt.Errorf("ciphertext is empty")
var invalidAESKeyErr = fmt.Errorf("aes key is invalid")

// AES consists of the symmetric key used in the aes encryption scheme.
type AES struct {
	Key []byte
}

// Encrypt encrypts a message. The encryption operation uses the CRT mode.
func (aesObj *AES) Encrypt(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return []byte{}, plainTextIsEmptyErr
	}

	block, err := aes.NewCipher(aesObj.Key)
	if err != nil {
		return nil, invalidAESKeyErr
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

// Decrypt decrypts a cipher text.
func (aesObj *AES) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return []byte{}, cipherTextIsEmptyErr
	}

	plaintext := make([]byte, len(ciphertext[aes.BlockSize:]))

	block, err := aes.NewCipher(aesObj.Key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	return plaintext, nil
}
