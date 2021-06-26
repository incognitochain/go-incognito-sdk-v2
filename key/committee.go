package key

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
)

// CommitteePublicKey consists of public keys of a user used in the consensus protocol.
// A CommitteePublicKey has
//	- IncPubKey: the public key of the user.
//	- MiningPubKey: a list of keys used in the consensus protocol.
//		+ BLS: used to sign blocks, create votes inside the Incognito network.
//		+ ECDSA: used to sign blocks for interacting with outside blockchain networks.
type CommitteePublicKey struct {
	IncPubKey    PublicKey
	MiningPubKey map[string][]byte
}

// IsEqualMiningPubKey checks if a CommitteePublicKey is equal to another CommitteePublicKey w.r.t the given consensus name.
func (pubKey *CommitteePublicKey) IsEqualMiningPubKey(consensusName string, k *CommitteePublicKey) bool {
	u, _ := pubKey.GetMiningKey(consensusName)
	b, _ := k.GetMiningKey(consensusName)
	return reflect.DeepEqual(u, b)
}

func NewCommitteePublicKey() *CommitteePublicKey {
	return &CommitteePublicKey{
		IncPubKey:    PublicKey{},
		MiningPubKey: make(map[string][]byte),
	}
}

// CheckSanityData checks sanity of a CommitteePublicKey.
func (pubKey *CommitteePublicKey) CheckSanityData() bool {
	if (len(pubKey.IncPubKey) != common.PublicKeySize) ||
		(len(pubKey.MiningPubKey[common.BlsConsensus]) != common.BLSPublicKeySize) ||
		(len(pubKey.MiningPubKey[common.BridgeConsensus]) != common.BriPublicKeySize) {
		return false
	}
	return true
}

// FromString sets a CommitteePublicKey from a string.
func (pubKey *CommitteePublicKey) FromString(keyString string) error {
	keyBytes, ver, err := base58.Base58Check{}.Decode(keyString)
	if (ver != common.ZeroByte) || (err != nil) {
		return NewError(B58DecodePubKeyErr, errors.New(ErrCodeMessage[B58DecodePubKeyErr].Message))
	}
	err = json.Unmarshal(keyBytes, pubKey)
	if err != nil {
		return NewError(JSONError, errors.New(ErrCodeMessage[JSONError].Message))
	}
	return nil
}

// NewCommitteeKeyFromSeed creates a new NewCommitteeKeyFromSeed given a seed and a public key.
func NewCommitteeKeyFromSeed(seed, incPubKey []byte) (CommitteePublicKey, error) {
	CommitteePublicKey := new(CommitteePublicKey)
	CommitteePublicKey.IncPubKey = incPubKey
	CommitteePublicKey.MiningPubKey = map[string][]byte{}
	_, blsPubKey := BLSKeyGen(seed)
	blsPubKeyBytes := PKBytes(blsPubKey)
	CommitteePublicKey.MiningPubKey[common.BlsConsensus] = blsPubKeyBytes
	_, briPubKey := BridgeKeyGen(seed)
	briPubKeyBytes := BridgePKBytes(&briPubKey)
	CommitteePublicKey.MiningPubKey[common.BridgeConsensus] = briPubKeyBytes
	return *CommitteePublicKey, nil
}

// FromBytes sets raw-data to a CommitteePublicKey.
func (pubKey *CommitteePublicKey) FromBytes(keyBytes []byte) error {
	err := json.Unmarshal(keyBytes, pubKey)
	if err != nil {
		return NewError(JSONError, err)
	}
	return nil
}

// RawBytes returns the raw-byte data of a CommitteePublicKey.
func (pubKey *CommitteePublicKey) RawBytes() ([]byte, error) {
	keys := make([]string, 0, len(pubKey.MiningPubKey))
	for k := range pubKey.MiningPubKey {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	res := pubKey.IncPubKey
	for _, k := range keys {
		res = append(res, pubKey.MiningPubKey[k]...)
	}
	return res, nil
}

// Bytes returns the JSON-marshalled data of a CommitteePublicKey.
func (pubKey *CommitteePublicKey) Bytes() ([]byte, error) {
	res, err := json.Marshal(pubKey)
	if err != nil {
		return []byte{0}, NewError(JSONError, err)
	}
	return res, nil
}

// GetNormalKey returns the public key of a CommitteePublicKey.
func (pubKey *CommitteePublicKey) GetNormalKey() []byte {
	return pubKey.IncPubKey
}

// GetMiningKey returns the mining key of a CommitteePublicKey given the consensus scheme.
func (pubKey *CommitteePublicKey) GetMiningKey(schemeName string) ([]byte, error) {
	allKey := map[string][]byte{}
	var ok bool
	allKey[schemeName], ok = pubKey.MiningPubKey[schemeName]
	if !ok {
		return nil, errors.New("this schemeName doesn't exist")
	}
	allKey[common.BridgeConsensus], ok = pubKey.MiningPubKey[common.BridgeConsensus]
	if !ok {
		return nil, errors.New("this lightweight schemeName doesn't exist")
	}
	result, err := json.Marshal(allKey)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetMiningKeyBase58 returns the base58-encoded mining key of a CommitteePublicKey given the consensus scheme.
func (pubKey *CommitteePublicKey) GetMiningKeyBase58(schemeName string) string {
	b, _ := pubKey.RawBytes()
	key := schemeName + string(b)
	value, exist := getMiningKeyBase58Cache.Get(key)
	if exist {
		return value.(string)
	}
	keyBytes, ok := pubKey.MiningPubKey[schemeName]
	if !ok {
		return ""
	}
	encodeData := base58.Base58Check{}.Encode(keyBytes, common.Base58Version)
	getMiningKeyBase58Cache.Add(key, encodeData)
	return encodeData
}

// GetIncKeyBase58 returns the base58-encoded public key of a CommitteePublicKey.
func (pubKey *CommitteePublicKey) GetIncKeyBase58() string {
	return base58.Base58Check{}.Encode(pubKey.IncPubKey, common.Base58Version)
}

// ToBase58 returns the base58-encoded representation of a CommitteePublicKey
func (pubKey *CommitteePublicKey) ToBase58() (string, error) {
	if pubKey == nil {
		result, err := json.Marshal(pubKey)
		if err != nil {
			return "", err
		}
		return base58.Base58Check{}.Encode(result, common.Base58Version), nil
	}

	b, _ := pubKey.RawBytes()
	key := string(b)
	value, exist := toBase58Cache.Get(key)
	if exist {
		return value.(string), nil
	}
	result, err := json.Marshal(pubKey)
	if err != nil {
		return "", err
	}
	encodeData := base58.Base58Check{}.Encode(result, common.Base58Version)
	toBase58Cache.Add(key, encodeData)
	return encodeData, nil
}

// FromBase58 recovers the CommitteePublicKey from its base58-representation.
func (pubKey *CommitteePublicKey) FromBase58(keyString string) error {
	keyBytes, ver, err := base58.Base58Check{}.Decode(keyString)
	if (ver != common.ZeroByte) || (err != nil) {
		return errors.New("wrong input")
	}
	return json.Unmarshal(keyBytes, pubKey)
}

// CommitteeKeyString is the string alternative to a CommitteePublicKey.
type CommitteeKeyString struct {
	IncPubKey    string
	MiningPubKey map[string]string
}

// IsEqual checks if a CommitteePublicKey is equal to the input CommitteePublicKey.
func (pubKey *CommitteePublicKey) IsEqual(target CommitteePublicKey) bool {
	if bytes.Compare(pubKey.IncPubKey[:], target.IncPubKey[:]) != 0 {
		return false
	}
	if pubKey.MiningPubKey == nil && target.MiningPubKey != nil {
		return false
	}
	for key, value := range pubKey.MiningPubKey {
		if targetValue, ok := target.MiningPubKey[key]; !ok {
			return false
		} else {
			if bytes.Compare(targetValue, value) != 0 {
				return false
			}
		}
	}
	return true
}

var getMiningKeyBase58Cache, _ = lru.New(2000)
var toBase58Cache, _ = lru.New(2000)
