package key

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

const (
	CBridgeSigSz = 65
)

// BridgeKeyGen generates a pair of ecdsa.PrivateKey and ecdsa.PublicKey from the given seed.
func BridgeKeyGen(seed []byte) (ecdsa.PrivateKey, ecdsa.PublicKey) {
	priKey := new(ecdsa.PrivateKey)
	priKey.Curve = ethCrypto.S256()
	priKey.D = common.B2ImN(seed)
	priKey.PublicKey.X, priKey.PublicKey.Y = priKey.Curve.ScalarBaseMult(priKey.D.Bytes())
	return *priKey, priKey.PublicKey
}

// BridgePKBytes returns the compressed version of a ecdsa.PublicKey.
func BridgePKBytes(pubKey *ecdsa.PublicKey) []byte {
	return ethCrypto.CompressPubkey(pubKey)
}

// DecodeECDSASig decodes an ecdsa signature given its string representation.
func DecodeECDSASig(sigStr string) (v byte, r []byte, s []byte, err error) {
	sig, err := hex.DecodeString(sigStr)
	if (len(sig) != CBridgeSigSz) || (err != nil) {
		err = fmt.Errorf("signature size is invalid")
		return
	}
	v = sig[64] + 27
	r = sig[:32]
	s = sig[32:64]
	return
}
