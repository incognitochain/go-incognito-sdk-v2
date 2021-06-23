package key

import (
	"math/big"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// BLSKeyGen takes an input seed and returns a BLS Key.
func BLSKeyGen(seed []byte) (*big.Int, *bn256.G2) {
	sk := BLSSKGen(seed)
	return sk, BLSPKGen(sk)
}

// BLSSKGen takes a seed and returns a BLS secret key.
func BLSSKGen(seed []byte) *big.Int {
	sk := big.NewInt(0)
	sk.SetBytes(common.HashB(seed))
	for {
		if sk.Cmp(bn256.Order) == -1 {
			break
		}
		sk.SetBytes(common.Hash4Bls(sk.Bytes()))
	}
	return sk
}

// BLSPKGen takes a secret key and returns a BLS public key.
func BLSPKGen(sk *big.Int) *bn256.G2 {
	pk := new(bn256.G2)
	pk = pk.ScalarBaseMult(sk)
	return pk
}

// PKBytes takes as input a public key point and returns the corresponding slice of bytes.
func PKBytes(pk *bn256.G2) PublicKey {
	return pk.Marshal()
}
