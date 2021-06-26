package crypto

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"

	C25519 "github.com/incognitochain/incognito-chain/privacy/curve25519"
)

type Scalar struct {
	key C25519.Key
}

func (sc Scalar) GetKey() C25519.Key {
	return sc.key
}

func (sc Scalar) String() string {
	return fmt.Sprintf("%x", sc.key[:])
}

func (sc Scalar) MarshalText() []byte {
	return []byte(fmt.Sprintf("%x", sc.key[:]))
}

func (sc *Scalar) UnmarshalText(data []byte) (*Scalar, error) {
	byteSlice, _ := hex.DecodeString(string(data))
	if len(byteSlice) != Ed25519KeySize {
		return nil, fmt.Errorf("incorrect key size")
	}
	copy(sc.key[:], byteSlice)
	return sc, nil
}

func (sc Scalar) ToBytes() [Ed25519KeySize]byte {
	return sc.key.ToBytes()
}

func (sc Scalar) ToBytesS() []byte {
	slice := sc.key.ToBytes()
	return slice[:]
}

func (sc *Scalar) FromBytes(b [Ed25519KeySize]byte) *Scalar {
	sc.key.FromBytes(b)
	return sc
}

func (sc *Scalar) FromBytesS(b []byte) *Scalar {
	var array [Ed25519KeySize]byte
	copy(array[:], b)
	sc.key.FromBytes(array)

	return sc
}

func (sc *Scalar) SetKeyUnsafe(a *C25519.Key) *Scalar {
	sc.key = *a
	return sc
}

func (sc *Scalar) SetKey(a *C25519.Key) (*Scalar, error) {
	sc.key = *a
	if sc.ScalarValid() == false {
		return nil, fmt.Errorf("invalid key value")
	}
	return sc, nil
}

// Set sets v to a Scalar.
func (sc *Scalar) Set(v *Scalar) *Scalar {
	sc.key = v.key
	return sc
}

// RandomScalar returns a random Scalar.
func RandomScalar() *Scalar {
	sc := new(Scalar)
	key := C25519.RandomScalar()
	sc.key = *key
	return sc
}

// HashToScalar returns the hash of msg in the form of a scalar.
func HashToScalar(msg []byte) *Scalar {
	key := C25519.HashToScalar(msg)
	sc, err := new(Scalar).SetKey(key)
	if err != nil {
		return nil
	}
	return sc
}

// FromUint64 sets the value of a Scalar to v.
func (sc *Scalar) FromUint64(v uint64) *Scalar {
	sc, err := sc.SetKey(d2h(v))
	if err != nil {
		return nil
	}
	return sc
}

// ToUint64Little returns the uint64 value of a Scalar.
func (sc *Scalar) ToUint64Little() uint64 {
	reversedKey := Reverse(sc.key)
	keyBN := new(big.Int).SetBytes(reversedKey[:])
	return keyBN.Uint64()
}

// Add returns (a + b) % CurveOrder.
func (sc *Scalar) Add(a, b *Scalar) *Scalar {
	var res C25519.Key
	C25519.ScAdd(&res, &a.key, &b.key)
	sc.key = res
	return sc
}

// Sub returns (a - b) % CurveOrder.
func (sc *Scalar) Sub(a, b *Scalar) *Scalar {
	var res C25519.Key
	C25519.ScSub(&res, &a.key, &b.key)
	sc.key = res
	return sc
}

// Mul returns a * b % CurveOrder.
func (sc *Scalar) Mul(a, b *Scalar) *Scalar {
	var res C25519.Key
	C25519.ScMul(&res, &a.key, &b.key)
	sc.key = res
	return sc
}

// MulAdd return (a * b + c) % CurveOrder.
func (sc *Scalar) MulAdd(a, b, c *Scalar) *Scalar {
	var res C25519.Key
	C25519.ScMulAdd(&res, &a.key, &b.key, &c.key)
	sc.key = res
	return sc
}

// Exp returns a^v % CurveOrder.
func (sc *Scalar) Exp(a *Scalar, v uint64) *Scalar {
	var res C25519.Key
	C25519.ScMul(&res, &a.key, &a.key)
	for i := 0; i < int(v)-2; i++ {
		C25519.ScMul(&res, &res, &a.key)
	}

	sc.key = res
	return sc
}

// ScalarValid checks if a Scalar is valid.
func (sc *Scalar) ScalarValid() bool {
	return C25519.ScValid(&sc.key)
}

// IsOne checks if a Scalar equals to 1.
func (sc *Scalar) IsOne() bool {
	s := sc.key
	return ((int(s[0]|s[1]|s[2]|s[3]|s[4]|s[5]|s[6]|s[7]|s[8]|
		s[9]|s[10]|s[11]|s[12]|s[13]|s[14]|s[15]|s[16]|s[17]|
		s[18]|s[19]|s[20]|s[21]|s[22]|s[23]|s[24]|s[25]|s[26]|
		s[27]|s[28]|s[29]|s[30]|s[31])-1)>>8)+1 == 1
}

// IsScalarEqual checks if two scalars are the same.
func IsScalarEqual(sc1, sc2 *Scalar) bool {
	tmpA := sc1.ToBytesS()
	tmpB := sc2.ToBytesS()

	return subtle.ConstantTimeCompare(tmpA, tmpB) == 1
}

// Compare returns -1 if a < b, 0 if a = b, and 1 if a > b.
func Compare(a, b *Scalar) int {
	tmpA := a.ToBytesS()
	tmpB := b.ToBytesS()

	for i := Ed25519KeySize - 1; i >= 0; i-- {
		if uint64(tmpA[i]) > uint64(tmpB[i]) {
			return 1
		}

		if uint64(tmpA[i]) < uint64(tmpB[i]) {
			return -1
		}
	}
	return 0
}

// IsZero checks if a scalar equals to 0.
func (sc *Scalar) IsZero() bool {
	if sc == nil {
		return false
	}
	return C25519.ScIsZero(&sc.key)
}

// Invert returns the a^-1.
func (sc *Scalar) Invert(a *Scalar) *Scalar {
	var inverseResult C25519.Key
	x := a.key

	reversedX := Reverse(x)
	bigX := new(big.Int).SetBytes(reversedX[:])

	reverseL := Reverse(C25519.CurveOrder()) // as speed improvements it can be made constant
	bigL := new(big.Int).SetBytes(reverseL[:])

	var inverse big.Int
	inverse.ModInverse(bigX, bigL)

	inverseBytes := inverse.Bytes()

	if len(inverseBytes) > Ed25519KeySize {
		panic("Inverse cannot be more than Ed25519KeySize bytes in this domain")
	}

	for i, j := 0, len(inverseBytes)-1; i < j; i, j = i+1, j-1 {
		inverseBytes[i], inverseBytes[j] = inverseBytes[j], inverseBytes[i]
	}
	copy(inverseResult[:], inverseBytes[:]) // copy the bytes  as they should be

	sc.key = inverseResult
	return sc
}

// Reverse returns the reverse byte-array of a C25519.Key x.
func Reverse(x C25519.Key) (result C25519.Key) {
	result = x

	// A key is in little-endian, but the big package wants the bytes in
	// big-endian, so Reverse them.
	lenB := len(x) // its hardcoded 32 bytes, so why do len but lets do it
	for i := 0; i < lenB/2; i++ {
		result[i], result[lenB-1-i] = result[lenB-1-i], result[i]
	}
	return
}

func d2h(val uint64) *C25519.Key {
	key := new(C25519.Key)
	for i := 0; val > 0; i++ {
		key[i] = byte(val & 0xFF)
		val /= 256
	}
	return key
}
