package base58

import (
	"fmt"
)

// Encode encodes the passed bytes into a base58 encoded string.
func Encode(bin []byte) string {
	return FastBase58Encoding(bin)
}

// FastBase58Encoding encodes the passed bytes into a base58 encoded string.
func FastBase58Encoding(bin []byte) string {
	return FastBase58EncodingAlphabet(bin, BTCAlphabet)
}

// FastBase58EncodingAlphabet encodes the passed bytes into a base58 encoded
// string with the passed alphabet.
func FastBase58EncodingAlphabet(bin []byte, alphabet *Alphabet) string {
	zero := alphabet.encode[0]

	binSize := len(bin)
	var i, j, zCount, high int
	var carry uint32

	for zCount < binSize && bin[zCount] == 0 {
		zCount++
	}

	size := (binSize-zCount)*138/100 + 1
	var buf = make([]byte, size)

	high = size - 1
	for i = zCount; i < binSize; i++ {
		j = size - 1
		for carry = uint32(bin[i]); j > high || carry != 0; j-- {
			carry = carry + 256*uint32(buf[j])
			buf[j] = byte(carry % 58)
			carry /= 58
		}
		high = j
	}

	for j = 0; j < size && buf[j] == 0; j++ {
	}

	var b58 = make([]byte, size-j+zCount)

	if zCount != 0 {
		for i = 0; i < zCount; i++ {
			b58[i] = zero
		}
	}

	for i = zCount; j < size; i++ {
		b58[i] = alphabet.encode[buf[j]]
		j++
	}

	return string(b58)
}

// Decode decodes the base58 encoded bytes.
func Decode(str string) ([]byte, error) {
	return FastBase58Decoding(str)
}

// FastBase58Decoding decodes the base58 encoded bytes.
func FastBase58Decoding(str string) ([]byte, error) {
	return FastBase58DecodingAlphabet(str, BTCAlphabet)
}

// FastBase58DecodingAlphabet decodes the base58 encoded bytes using the given
// b58 alphabet.
func FastBase58DecodingAlphabet(str string, alphabet *Alphabet) ([]byte, error) {
	if len(str) == 0 {
		return nil, fmt.Errorf("zero length string")
	}

	var (
		t        uint64
		zMask, c uint32
		zCount   int

		b58u  = []rune(str)
		b58sz = len(b58u)

		outISize  = (b58sz + 3) / 4 // check to see if we need to change this buffer size to optimize
		binU      = make([]byte, (b58sz+3)*3)
		bytesLeft = b58sz % 4

		zero = rune(alphabet.encode[0])
	)

	if bytesLeft > 0 {
		zMask = 0xffffffff << uint32(bytesLeft*8)
	} else {
		bytesLeft = 4
	}

	var outi = make([]uint32, outISize)

	for i := 0; i < b58sz && b58u[i] == zero; i++ {
		zCount++
	}

	for _, r := range b58u {
		if r > 127 {
			return nil, fmt.Errorf("high-bit set on invalid digit")
		}
		if alphabet.decode[r] == -1 {
			return nil, fmt.Errorf("invalid base58 digit (%q)", r)
		}

		c = uint32(alphabet.decode[r])

		for j := outISize - 1; j >= 0; j-- {
			t = uint64(outi[j])*58 + uint64(c)
			c = uint32(t>>32) & 0x3f
			outi[j] = uint32(t & 0xffffffff)
		}

		if c > 0 {
			return nil, fmt.Errorf("output number too big (carry to the next int32)")
		}

		if outi[0]&zMask != 0 {
			return nil, fmt.Errorf("output number too big (last int32 filled too far)")
		}
	}

	var j, cnt int
	for j, cnt = 0, 0; j < outISize; j++ {
		for mask := byte(bytesLeft-1) * 8; mask <= 0x18; mask, cnt = mask-8, cnt+1 {
			binU[cnt] = byte(outi[j] >> mask)
		}
		if j == 0 {
			bytesLeft = 4 // because it could be less than 4 the first time through
		}
	}

	for n, v := range binU {
		if v > 0 {
			start := n - zCount
			if start < 0 {
				start = 0
			}
			return binU[start:cnt], nil
		}
	}
	return binU[:cnt], nil
}
