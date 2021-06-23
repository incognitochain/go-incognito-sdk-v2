package utils

import (
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
)

// GenerateChallenge returns the hash of n points in G appended with the input values
//
// blake_2b(G[0]||G[1]||...||G[CM_CAPACITY-1]||<values>);
// G[i] is list of all generator point of Curve
func GenerateChallenge(values [][]byte) *crypto.Scalar {
	bytes := make([]byte, 0)
	for i := 0; i < len(crypto.PedCom.G); i++ {
		bytes = append(bytes, crypto.PedCom.G[i].ToBytesS()...)
	}

	for i := 0; i < len(values); i++ {
		bytes = append(bytes, values[i]...)
	}

	hash := crypto.HashToScalar(bytes)
	//res := new(big.Int).SetBytes(hash)
	//res.Mod(res, crypto.Curve.Params().N)
	return hash
}
