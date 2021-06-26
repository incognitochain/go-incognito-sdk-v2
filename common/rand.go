package common

import (
	cRand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

var alphabet = "abcdefghijklmnopqrstvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandInt returns a random int number using math/rand
func RandInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}

func RandUint64() uint64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint64()
}

// RandIntInterval returns a random int in range [L; R]
func RandIntInterval(L, R int) int {
	length := R - L + 1
	r := RandInt() % length
	return L + r
}

// RandInt64 returns a random int64 number using math/rand
func RandInt64() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63()
}

// RandBigIntMaxRange generates a random big.Int whose value is less than a given max value.
func RandBigIntMaxRange(max *big.Int) (*big.Int, error) {
	return cRand.Int(cRand.Reader, max)
}

// RandBytes generates a random l-byte long slice.
func RandBytes(l int) []byte {
	randBytes := make([]byte, l)
	_, err := rand.Read(randBytes)
	if err != nil {
		return randBytes
	}

	return randBytes
}

// RandChars returns a random l-character long string.
func RandChars(l int) string {
	res := ""
	for i := 0; i < l; i++ {
		r := RandInt() % len(alphabet)
		res += string(alphabet[r])
	}

	return res
}
