package jsonresult

import "github.com/incognitochain/go-incognito-sdk-v2/common/base58"

// EncodeBase58Check returns the base58-encoded string for a given byte slice.
func EncodeBase58Check(b []byte) string {
	if b == nil || len(b) == 0 {
		return ""
	}
	return base58.Base58Check{}.Encode(b, 0x0)
}
