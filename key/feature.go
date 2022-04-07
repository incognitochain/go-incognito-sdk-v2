package key

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"math/big"
)

// OTDepositKey represents a pair of one-time depositing key for shielding.
type OTDepositKey struct {
	// PrivateKey is used to sign shielding requests when this OTDepositKey is employed.
	PrivateKey []byte

	// PublicKey serves as a chain-code to generate a unique multi-sig address (depositAddr) for shielding request. It is
	// derived from the PrivateKey, and is used to verify signatures signed by the PrivateKey to authorize shielding requests.
	// It is used to replace the Incognito address for better privacy. Different PublicKey results in a different depositAddr.
	// Note that one can re-use the OTPubKey in many shielding requests. However, this is NOT RECOMMENDED because it
	// will lower the privacy level and allow an observer to link his shields.
	PublicKey []byte

	// Index is the index of the current OTDepositKey. Since most of the time, an OTDepositKey is generated from a master key,
	// this Index serves as a tool for the ease of key management. More detail about the Index can be found here: https://we.incognito.org/t/work-in-progress-one-time-shielding-addresses/15677
	Index uint64
}

func (k OTDepositKey) MarshalJSON() ([]byte, error) {
	type holder struct {
		PrivateKey string
		PublicKey  string
		Index      uint64
	}

	privateKeyStr := base58.Base58Check{}.NewEncode(k.PrivateKey, 0)
	pubKeyStr := base58.Base58Check{}.NewEncode(k.PublicKey, 0)
	h := holder{
		PrivateKey: privateKeyStr,
		PublicKey:  pubKeyStr,
		Index:      k.Index,
	}

	return json.Marshal(h)
}

func (k *OTDepositKey) UnmarshalJSON(data []byte) error {
	type holder struct {
		PrivateKey string
		PublicKey  string
		Index      uint64
	}
	var tmpH holder
	err := json.Unmarshal(data, &tmpH)
	if err != nil {
		return err
	}

	privateKey, _, err := base58.Base58Check{}.Decode(tmpH.PrivateKey)
	if err != nil {
		return err
	}
	pubKey, _, err := base58.Base58Check{}.Decode(tmpH.PublicKey)
	if err != nil {
		return nil
	}

	k.Index = tmpH.Index
	k.PublicKey = pubKey
	k.PrivateKey = privateKey

	return nil
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
	pubKey := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], privateKey)

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
	pubKey := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], privateKey)

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
	pubKey := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenPrivateKeyIndex], privateKey)

	return &OTDepositKey{
		PrivateKey: privateKey.ToBytesS(),
		PublicKey:  pubKey.ToBytesS(),
		Index:      index,
	}, nil
}
