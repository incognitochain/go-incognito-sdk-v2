package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

var invalidMaxHashSizeErr = fmt.Errorf("invalid max hash size")
var invalidHashSizeErr = fmt.Errorf("invalid hash size")
var nilHashErr = fmt.Errorf("input hash is nil")

type Hash [HashSize]byte

// MarshalText converts hashObj string to bytes array.
func (hashObj Hash) MarshalText() ([]byte, error) {
	return []byte(hashObj.String()), nil
}

// UnmarshalText reverts bytes array to hashObj.
func (hashObj Hash) UnmarshalText(text []byte) error {
	copy(hashObj[:], text)
	return nil
}

// UnmarshalJSON unmarshal json data to hashObj.
func (hashObj *Hash) UnmarshalJSON(data []byte) error {
	hashString := ""
	err := json.Unmarshal(data, &hashString)
	if err != nil {
		return err
	}
	return hashObj.Decode(hashObj, hashString)
}

// SetBytes sets the bytes array which represent the hash.
func (hashObj *Hash) SetBytes(newHash []byte) error {
	l := len(newHash)
	if l != HashSize {
		return invalidHashSizeErr
	}
	copy(hashObj[:], newHash)

	return nil
}

// GetBytes returns bytes array of hashObj
func (hashObj *Hash) GetBytes() []byte {
	newBytes := make([]byte, 0)
	newBytes = make([]byte, len(hashObj))
	copy(newBytes, hashObj[:])
	return newBytes
}

// NewHash receives a bytes array and returns a corresponding object Hash
func (hashObj Hash) NewHash(newHash []byte) (*Hash, error) {
	err := hashObj.SetBytes(newHash)
	if err != nil {
		return nil, err
	}
	return &hashObj, err
}

// String returns the Hash as the hexadecimal string of the byte-reversed hash.
func (hashObj Hash) String() string {
	for i := 0; i < HashSize/2; i++ {
		hashObj[i], hashObj[HashSize-1-i] = hashObj[HashSize-1-i], hashObj[i]
	}
	return hex.EncodeToString(hashObj[:])
}

// IsEqual returns true if target is the same as hashObj.
func (hashObj *Hash) IsEqual(target *Hash) bool {
	if hashObj == nil && target == nil {
		return true
	}
	if hashObj == nil || target == nil {
		return false
	}
	return *hashObj == *target
}

// NewHashFromStr creates a Hash from a hash string.  The string should be
// the hexadecimal string of a byte-reversed hash, but any missing characters
// result in zero padding at the end of the Hash.
func (hashObj Hash) NewHashFromStr(hash string) (*Hash, error) {
	err := hashObj.Decode(&hashObj, hash)
	if err != nil {
		return nil, err
	}
	return &hashObj, nil
}

// Decode decodes the byte-reversed hexadecimal string encoding of a Hash to a
// destination.
func (hashObj *Hash) Decode(dst *Hash, src string) error {
	// Return error if hash string is too long.
	if len(src) > MaxHashStringSize {
		return invalidMaxHashSizeErr
	}

	// Hex decoder expects the hash to be a multiple of two.  When not, pad
	// with a leading zero.
	var srcBytes []byte
	if len(src)%2 == 0 {
		srcBytes = []byte(src)
	} else {
		srcBytes = make([]byte, 1+len(src))
		srcBytes[0] = '0'
		copy(srcBytes[1:], src)
	}

	// Hex decode the source bytes to a temporary destination.
	var reversedHash Hash
	_, err := hex.Decode(reversedHash[HashSize-hex.DecodedLen(len(srcBytes)):], srcBytes)
	if err != nil {
		return err
	}

	// Reverse copy from the temporary hash to destination.  Because the
	// temporary was zeroed, the written result will be correctly padded.
	for i, b := range reversedHash[:HashSize/2] {
		dst[i], dst[HashSize-1-i] = reversedHash[HashSize-1-i], b
	}

	return nil
}

// Cmp compare two hashes, if
// - hash = target: return 0.
// - hash > target: return 1.
// - hash < target: return -1.
func (hashObj *Hash) Cmp(target *Hash) (int, error) {
	if hashObj == nil || target == nil {
		return 0, nilHashErr
	}
	for i := 0; i < HashSize; i++ {
		if hashObj[i] > target[i] {
			return 1, nil
		}
		if hashObj[i] < target[i] {
			return -1, nil
		}
	}
	return 0, nil
}

// Keccak256 returns Keccak256 hash as a Hash object for storing and comparing.
func Keccak256(data ...[]byte) Hash {
	h := crypto.Keccak256(data...)
	r := Hash{}
	copy(r[:], h)
	return r
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h Hash) {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}
