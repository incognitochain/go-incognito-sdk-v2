package common

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"math/big"
	"reflect"
	"sort"
)

// SliceExists receives a slice and a item in interface type
// checks whether the slice contain the item or not
func SliceExists(slice interface{}, item interface{}) (bool, error) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return false, errors.New("SliceExists() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		interfacea := s.Index(i).Interface()
		if interfacea == item {
			return true, nil
		}
	}

	return false, nil
}

// GetShardIDFromLastByte receives a last byte of public key and
// returns a corresponding shardID
func GetShardIDFromLastByte(b byte) byte {
	return byte(int(b) % MaxShardNumber)
}

// RandBigIntMaxRange generates a big int with maximum value
func RandBigIntMaxRange(max *big.Int) (*big.Int, error) {
	return rand.Int(rand.Reader, max)
}

// RandBytes generates random bytes with length
func RandBytes(length int) []byte {
	rbytes := make([]byte, length)
	rand.Read(rbytes)
	return rbytes
}

// CompareStringArray receives 2 arrays of string
// and check whether 2 arrays is the same or not
func CompareStringArray(src []string, dst []string) bool {
	if len(src) != len(dst) {
		return false
	}
	for idx, val := range src {
		if dst[idx] != val {
			return false
		}
	}
	return true
}

// BytesToInt32 converts little endian 4-byte array to int32 number
func BytesToInt32(b []byte) (int32, error) {
	if len(b) != Int32Size {
		return 0, errors.New("invalid length of input BytesToInt32")
	}

	return int32(binary.LittleEndian.Uint32(b)), nil
}

// Int32ToBytes converts int32 number to little endian 4-byte array
func Int32ToBytes(value int32) []byte {
	b := make([]byte, Int32Size)
	binary.LittleEndian.PutUint32(b, uint32(value))
	return b
}

// IntToBytes converts an integer number to 2-byte array in big endian
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

// BytesToInt reverts an integer number from 2-byte array
func BytesToInt(bytesArr []byte) int {
	if len(bytesArr) != 2 {
		return 0
	}

	numInt := new(big.Int).SetBytes(bytesArr)
	return int(numInt.Int64())
}

// BytesToUint32 converts big endian 4-byte array to uint32 number
func BytesToUint32(b []byte) (uint32, error) {
	if len(b) != Uint32Size {
		return 0, errors.New("invalid length of input BytesToUint32")
	}
	return binary.BigEndian.Uint32(b), nil
}

// Uint32ToBytes converts uint32 number to big endian 4-byte array
func Uint32ToBytes(value uint32) []byte {
	b := make([]byte, Uint32Size)
	binary.BigEndian.PutUint32(b, value)
	return b
}

// BytesToUint64 converts little endian 8-byte array to uint64 number
func BytesToUint64(b []byte) (uint64, error) {
	if len(b) != Uint64Size {
		return 0, errors.New("invalid length of input BytesToUint64")
	}
	return binary.LittleEndian.Uint64(b), nil
}

// Uint64ToBytes converts uint64 number to little endian 8-byte array
func Uint64ToBytes(value uint64) []byte {
	b := make([]byte, Uint64Size)
	binary.LittleEndian.PutUint64(b, value)
	return b
}

// Int64ToBytes converts int64 number to little endian 8-byte array
func Int64ToBytes(value int64) []byte {
	return Uint64ToBytes(uint64(value))
}

// BoolToByte receives a value in bool
// and returns a value in byte
func BoolToByte(value bool) byte {
	var bitSetVar byte
	if value {
		bitSetVar = 1
	}
	return bitSetVar
}

// AddPaddingBigInt adds padding to big int to it is fixed size
// and returns bytes array
func AddPaddingBigInt(numInt *big.Int, fixedSize int) []byte {
	numBytes := numInt.Bytes()
	lenNumBytes := len(numBytes)
	zeroBytes := make([]byte, fixedSize-lenNumBytes)
	numBytes = append(zeroBytes, numBytes...)
	return numBytes
}

// AppendSliceString is a variadic function,
// receives some lists of array of strings
// and appends them to one list of array of strings
func AppendSliceString(arrayStrings ...[][]string) [][]string {
	res := [][]string{}
	for _, arrayString := range arrayStrings {
		res = append(res, arrayString...)
	}
	return res
}

func Uint16ToBytes(v uint16) [2]byte {
	var res [2]byte
	res[0] = uint8(v >> 8)
	res[1] = uint8(v & 0xff)
	return res
}

func BytesToUint16(b [2]byte) uint16 {
	return uint16(b[0])<<8 + uint16(b[1])
}

func BytesSToUint16(b []byte) (uint16, error) {
	if len(b) != 2 {
		return 0, errors.New("Cannot convert BytesSToUint16: length of byte is not 2")
	}
	var bytes [2]byte
	copy(bytes[:], b[:2])
	return BytesToUint16(bytes), nil
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

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// B2I convert byte array to big int which belong to Fp
func B2I(bytes []byte) *big.Int {
	res := big.NewInt(0)
	res.SetBytes(bytes)
	for res.Cmp(bn256.Order) != -1 {
		bytes = Hash4Bls(bytes)
		res.SetBytes(bytes)
	}
	return res
}

// B2ImN is Bytes to Int mod N, with N is secp256k1 curve order
func B2ImN(bytes []byte) *big.Int {
	x := big.NewInt(0)
	x.SetBytes(ethcrypto.Keccak256Hash(bytes).Bytes())
	for x.Cmp(ethcrypto.S256().Params().N) != -1 {
		x.SetBytes(ethcrypto.Keccak256Hash(x.Bytes()).Bytes())
	}
	return x
}

var PortalSupportedIncTokenIDs = []string{
	PortalBTCIDStr, // pBTC
	PortalBNBIDStr, // pBNB
}

// set MinAmountPortalPToken to avoid attacking with amount is less than smallest unit of cryptocurrency
// such as satoshi in BTC
var MinAmountPortalPToken = map[string]uint64{
	PortalBTCIDStr: 10,
	PortalBNBIDStr: 10,
}

var (
	EmptyRoot = HexToHash(HexEmptyRoot)
)

var EthContractAddressStr = MainETHContractAddressStr

var TIMESLOT = uint64(0) //need to be set when init chain

var (
	MaxTxSize    = uint64(100)  // unit KB = 100KB
	MaxBlockSize = uint64(2000) //unit kilobytes = 2 Megabyte
)

// special token ids (aka. PropertyID in custom token)
var (
	PRVCoinID             = Hash{4} // To send PRV in custom token
	PRVCoinName           = "PRV"   // To send PRV in custom token
	ConfidentialAssetID   = Hash{5}
	ConfidentialAssetName = "CA"
	MaxShardNumber        = 8 //programmatically config based on networkID
	AddressVersion        = 0
)

type PoolInfo struct {
	Token1IDStr     string
	Token1PoolValue uint64
	Token2IDStr     string
	Token2PoolValue uint64
}

type Share struct {
	TokenID1Str string
	TokenID2Str string
	ShareAmount uint64
}

var SupportedTokenID map[string]string
var ReversedSupportedTokenID map[string]string

func MakeReversedTokenID() {
	tmp := make(map[string]string)
	for name, id := range SupportedTokenID {
		tmp[id] = name
	}
	ReversedSupportedTokenID = tmp
}

var PDEPoolPairs map[string]*PoolInfo
var CurrentBeaconHeight uint64

func BuildPDEPoolForPairKey(
	token1IDStr string,
	token2IDStr string,
) []byte {
	beaconHeightBytes := []byte(fmt.Sprintf("%d-", CurrentBeaconHeight))
	pdePoolForPairByBCHeightPrefix := append([]byte("pdepool-"), beaconHeightBytes...)
	tokenIDStrs := []string{token1IDStr, token2IDStr}
	sort.Strings(tokenIDStrs)
	return append(pdePoolForPairByBCHeightPrefix, []byte(tokenIDStrs[0]+"-"+tokenIDStrs[1])...)
}

func GetPDEPoolPair(token1IDStr string, token2IDStr string, poolPairs map[string]*PoolInfo) (*PoolInfo, error) {
	if token1IDStr == token2IDStr {
		return nil, fmt.Errorf("tokenIDs are the same: %v", token1IDStr)
	}
	tmpPoolPairs := PDEPoolPairs
	if poolPairs != nil {
		tmpPoolPairs = poolPairs
	}
	keyPool := BuildPDEPoolForPairKey(token1IDStr, token2IDStr)

	if poolPair, ok := tmpPoolPairs[string(keyPool)]; ok {
		return poolPair, nil
	}

	return nil, fmt.Errorf("cannot found pool pair for tokenID %v and %v", token1IDStr, token2IDStr)
}

var ALPHABET = "abcdefghijklmnopqrstvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandChars(length int) string {
	res := ""
	for i := 0; i < length; i++ {
		r := RandInt() % len(ALPHABET)
		res += string(ALPHABET[r])
	}

	return res
}
