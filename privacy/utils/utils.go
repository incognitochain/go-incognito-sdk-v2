package utils

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto/curve25519"
	"math/big"
)

func ScalarToBigInt(sc *crypto.Scalar) *big.Int {
	keyR := crypto.Reverse(sc.GetKey())
	keyRByte := keyR.ToBytes()
	bi := new(big.Int).SetBytes(keyRByte[:])
	return bi
}

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

// ConvertIntToBinary represents a integer number in binary array with little endian with size n
func ConvertIntToBinary(inum int, n int) []byte {
	binary := make([]byte, n)

	for i := 0; i < n; i++ {
		binary[i] = byte(inum % 2)
		inum = inum / 2
	}

	return binary
}

// ConvertIntToBinary represents a integer number in binary
func ConvertUint64ToBinary(number uint64, n int) []*crypto.Scalar {
	if number == 0 {
		res := make([]*crypto.Scalar, n)
		for i := 0; i < n; i++ {
			res[i] = new(crypto.Scalar).FromUint64(0)
		}
		return res
	}

	binary := make([]*crypto.Scalar, n)

	for i := 0; i < n; i++ {
		binary[i] = new(crypto.Scalar).FromUint64(number % 2)
		number = number / 2
	}
	return binary
}


func SliceToArray(slice []byte) [crypto.Ed25519KeySize]byte {
	var array [crypto.Ed25519KeySize]byte
	copy(array[:], slice)
	return array
}

func ArrayToSlice(array [crypto.Ed25519KeySize]byte) []byte {
	var slice []byte
	slice = array[:]
	return slice
}