package utils

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/incognito-chain/privacy/curve25519"
	"math/big"
)

// ScalarToBigInt converts a scalar into a big.Int.
func ScalarToBigInt(sc *crypto.Scalar) *big.Int {
	keyR := crypto.Reverse(sc.GetKey())
	keyRByte := keyR.ToBytes()
	bi := new(big.Int).SetBytes(keyRByte[:])
	return bi
}

// BigIntToScalar converts a big.Int number into a crypto.Scalar.
func BigIntToScalar(bi *big.Int) *crypto.Scalar {
	biByte := common.AddPaddingBigInt(bi, crypto.Ed25519KeySize)
	var key curve25519.Key
	key.FromBytes(SliceToArray(biByte))
	keyR := crypto.Reverse(key)
	sc, err := new(crypto.Scalar).SetKey(&keyR)
	if err != nil {
		return nil
	}
	return sc
}

// ConvertIntToBinary represents a integer number in binary array with little endian with size n.
func ConvertIntToBinary(iNum int, n int) []byte {
	binary := make([]byte, n)

	for i := 0; i < n; i++ {
		binary[i] = byte(iNum % 2)
		iNum = iNum / 2
	}

	return binary
}

// SliceToArray copies a slice of bytes into an array of 32 bytes.
func SliceToArray(slice []byte) [crypto.Ed25519KeySize]byte {
	var array [crypto.Ed25519KeySize]byte
	copy(array[:], slice)
	return array
}
