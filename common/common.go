package common

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"reflect"
)

// SliceExists checks if an item exists in a slice.
func SliceExists(slice interface{}, item interface{}) (bool, error) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return false, fmt.Errorf("non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		tmpItem := s.Index(i).Interface()
		if tmpItem == item {
			return true, nil
		}
	}

	return false, nil
}

// GetShardIDFromLastByte returns the shardID from the last byte b of a public key.
// The shardID is calculated by taking the remainder of b % MaxShardNumber.
func GetShardIDFromLastByte(b byte) byte {
	return byte(int(b) % MaxShardNumber)
}

// IntToBytes converts an integer number to 2-byte array in big endian.
func IntToBytes(n int) []byte {
	if n == 0 {
		return []byte{0, 0}
	}

	a := big.NewInt(int64(n))

	if len(a.Bytes()) > 2 {
		return []byte{}
	}

	if len(a.Bytes()) == 1 {
		return []byte{0, a.Bytes()[0]}
	}

	return a.Bytes()
}

// BytesToInt reverts an integer number from 2-byte array.
func BytesToInt(bytesArr []byte) int {
	if len(bytesArr) != 2 {
		return 0
	}

	numInt := new(big.Int).SetBytes(bytesArr)
	return int(numInt.Int64())
}

// BytesToUint32 converts big endian 4-byte array to uint32 number.
func BytesToUint32(b []byte) (uint32, error) {
	if len(b) != Uint32Size {
		return 0, fmt.Errorf("invalid length of input BytesToUint32")
	}
	return binary.BigEndian.Uint32(b), nil
}

// Uint32ToBytes converts uint32 number to big endian 4-byte array.
func Uint32ToBytes(value uint32) []byte {
	b := make([]byte, Uint32Size)
	binary.BigEndian.PutUint32(b, value)
	return b
}

// AddPaddingBigInt adds padding to big int to it is fixed size.
func AddPaddingBigInt(numInt *big.Int, fixedSize int) []byte {
	numBytes := numInt.Bytes()
	lenNumBytes := len(numBytes)
	zeroBytes := make([]byte, fixedSize-lenNumBytes)
	numBytes = append(zeroBytes, numBytes...)
	return numBytes
}

// Has0xPrefix validates str begins with '0x' or '0X'.
func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// B2ImN is Bytes to Int mod N, with N is Secp256k1 curve order
func B2ImN(bytes []byte) *big.Int {
	x := big.NewInt(0)
	x.SetBytes(ethCrypto.Keccak256Hash(bytes).Bytes())
	for x.Cmp(ethCrypto.S256().Params().N) != -1 {
		x.SetBytes(ethCrypto.Keccak256Hash(x.Bytes()).Bytes())
	}
	return x
}

var (
	MaxTxSize = uint64(100) // unit KB = 100KB
)

var (
	PRVCoinID           = Hash{4}
	ConfidentialAssetID = Hash{5}
	MaxShardNumber      = 8 //programmatically config based on networkID
	AddressVersion      = 1
)
