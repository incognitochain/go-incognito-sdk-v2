package key

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"math/big"
)

// OTDepositKey represents a pair of one-time depositing key for shielding.
type OTDepositKey struct {
	PrivateKey []byte
	PublicKey  []byte
	Index      uint64
}

// GenerateOTDepositKey generates a new OTDepositKey from the keySet with the given tokenID and index.
func (keySet *KeySet) GenerateOTDepositKey(tokenIDStr string, index uint64) (*OTDepositKey, error) {
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, err
	}

	tmp := append([]byte(common.PortalV4DepositKeyGenSeed), tokenID[:]...)
	masterDepositSeed := common.SHA256(append(keySet.PrivateKey[:], tmp...))
	indexBig := new(big.Int).SetUint64(index)

	privateKey := crypto.HashToScalar(append(masterDepositSeed, indexBig.Bytes()...))
	pubKey := new(crypto.Point).ScalarMultBase(privateKey)

	return &OTDepositKey{
		PrivateKey: privateKey.ToBytesS(),
		PublicKey:  pubKey.ToBytesS(),
		Index:      index,
	}, nil
}

// GenerateOTDepositKeyFromPrivateKey generates a new OTDepositKey from the given privateKey, tokenID and index.
func GenerateOTDepositKeyFromPrivateKey(incPrivateKey []byte, tokenIDStr string, index uint64) (*OTDepositKey, error) {
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, err
	}

	tmp := append([]byte(common.PortalV4DepositKeyGenSeed), tokenID[:]...)
	masterDepositSeed := common.SHA256(append(incPrivateKey[:], tmp...))
	indexBig := new(big.Int).SetUint64(index)

	privateKey := crypto.HashToScalar(append(masterDepositSeed, indexBig.Bytes()...))
	pubKey := new(crypto.Point).ScalarMultBase(privateKey)

	return &OTDepositKey{
		PrivateKey: privateKey.ToBytesS(),
		PublicKey:  pubKey.ToBytesS(),
		Index:      index,
	}, nil
}

// GenerateOTDepositKeyFromMasterDepositSeed generates a new OTDepositKey from the given masterDepositSeed, tokenID and index.
func GenerateOTDepositKeyFromMasterDepositSeed(masterDepositSeed []byte, index uint64) (*OTDepositKey, error) {
	indexBig := new(big.Int).SetUint64(index)

	privateKey := crypto.HashToScalar(append(masterDepositSeed, indexBig.Bytes()...))
	pubKey := new(crypto.Point).ScalarMultBase(privateKey)

	return &OTDepositKey{
		PrivateKey: privateKey.ToBytesS(),
		PublicKey:  pubKey.ToBytesS(),
		Index:      index,
	}, nil
}
