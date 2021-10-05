package coin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/hybridencryption"
)

// PlainCoinV1 represents a decrypted CoinV1. It is mainly used as the input of a transaction v1 or a conversion transaction.
//
// This struct implements the PlainCoin interface.
type PlainCoinV1 struct {
	publicKey    *crypto.Point
	commitment   *crypto.Point
	snDerivator  *crypto.Scalar
	serialNumber *crypto.Point
	randomness   *crypto.Scalar
	value        uint64
	info         []byte //256 bytes
}

// Init initializes an empty PlainCoinV1 object.
func (pc *PlainCoinV1) Init() *PlainCoinV1 {
	pc.value = 0
	pc.randomness = new(crypto.Scalar)
	pc.publicKey = new(crypto.Point).Identity()
	pc.serialNumber = new(crypto.Point).Identity()
	pc.snDerivator = new(crypto.Scalar).FromUint64(0)
	pc.commitment = nil
	return pc
}

// GetVersion returns the version of a PlainCoinV1.
func (*PlainCoinV1) GetVersion() uint8 { return 1 }

// GetShardID returns the shardID in which a PlainCoinV1 belongs to.
func (pc *PlainCoinV1) GetShardID() (uint8, error) {
	if pc.publicKey == nil {
		return 255, fmt.Errorf("cannot get ShardID because PublicKey of PlainCoinV1 is concealed")
	}
	pubKeyBytes := pc.publicKey.ToBytes()
	lastByte := pubKeyBytes[crypto.Ed25519KeySize-1]
	shardID := common.GetShardIDFromLastByte(lastByte)
	return shardID, nil
}

// GetCommitment returns the commitment of a PlainCoinV1.
func (pc PlainCoinV1) GetCommitment() *crypto.Point { return pc.commitment }

// GetPublicKey returns the public key of a PlainCoinV1. For a PlainCoinV1, its public key is the public key of the owner.
func (pc PlainCoinV1) GetPublicKey() *crypto.Point { return pc.publicKey }

// GetSNDerivator returns the serial number derivator of a PlainCoinV1.
func (pc PlainCoinV1) GetSNDerivator() *crypto.Scalar { return pc.snDerivator }

// GetKeyImage returns the serial number of a PlainCoinV1.
func (pc PlainCoinV1) GetKeyImage() *crypto.Point { return pc.serialNumber }

// GetRandomness returns the randomness of a PlainCoinV1.
func (pc PlainCoinV1) GetRandomness() *crypto.Scalar { return pc.randomness }

// GetValue returns the value of a PlainCoinV1.
func (pc PlainCoinV1) GetValue() uint64 { return pc.value }

// GetInfo returns the info of a PlainCoinV1.
func (pc PlainCoinV1) GetInfo() []byte { return pc.info }

// GetAssetTag returns the asset tag of a PlainCoinV1. For a PlainCoinV1, this value is nil.
func (pc PlainCoinV1) GetAssetTag() *crypto.Point { return nil }

// GetTxRandom returns the transaction randomness of a PlainCoinV1. For a PlainCoinV1, this value is nil.
func (pc PlainCoinV1) GetTxRandom() *TxRandom { return nil }

// GetSharedRandom returns the shared random of a PlainCoinV1. For a PlainCoinV1, this value is nil.
func (pc PlainCoinV1) GetSharedRandom() *crypto.Scalar { return nil }

// GetSharedConcealRandom returns the shared random when concealing of a PlainCoinV1. For a PlainCoinV1, this value is nil.
func (pc PlainCoinV1) GetSharedConcealRandom() *crypto.Scalar { return nil }

// IsEncrypted checks if whether a PlainCoinV1 is encrypted. This value is always false.
func (pc PlainCoinV1) IsEncrypted() bool { return false }

// GetCoinDetailEncrypted returns the encrypted detail of a PlainCoinV1. For a PlainCoinV1, this value is always nil.
func (pc PlainCoinV1) GetCoinDetailEncrypted() []byte {
	return nil
}

// SetPublicKey sets v as the public key of a PlainCoinV1.
func (pc *PlainCoinV1) SetPublicKey(v *crypto.Point) { pc.publicKey = v }

// SetCommitment sets v as the commitment of a PlainCoinV1.
func (pc *PlainCoinV1) SetCommitment(v *crypto.Point) { pc.commitment = v }

// SetSNDerivator sets v as the serial number derivator of a PlainCoinV1.
func (pc *PlainCoinV1) SetSNDerivator(v *crypto.Scalar) { pc.snDerivator = v }

// SetKeyImage sets v as the serial number of a PlainCoinV1.
func (pc *PlainCoinV1) SetKeyImage(v *crypto.Point) { pc.serialNumber = v }

// SetRandomness sets v as the randomness of a PlainCoinV1.
func (pc *PlainCoinV1) SetRandomness(v *crypto.Scalar) { pc.randomness = v }

// SetValue sets v as the value of a PlainCoinV1.
func (pc *PlainCoinV1) SetValue(v uint64) { pc.value = v }

// SetInfo sets v as the info of a PlainCoinV1.
func (pc *PlainCoinV1) SetInfo(v []byte) {
	pc.info = make([]byte, len(v))
	copy(pc.info, v)
}

// ParsePrivateKeyOfCoin sets privateKey as the private key of a PlainCoinV1.
func (pc PlainCoinV1) ParsePrivateKeyOfCoin(privateKey key.PrivateKey) (*crypto.Scalar, error) {
	return new(crypto.Scalar).FromBytesS(privateKey), nil
}

// ParseKeyImageWithPrivateKey derives the key image of a PlainCoinV1 from its private key.
func (pc PlainCoinV1) ParseKeyImageWithPrivateKey(privateKey key.PrivateKey) (*crypto.Point, error) {
	k, err := pc.ParsePrivateKeyOfCoin(privateKey)
	if err != nil {
		return nil, err
	}
	keyImage := new(crypto.Point).Derive(
		crypto.PedCom.G[crypto.PedersenPrivateKeyIndex],
		k,
		pc.GetSNDerivator())
	pc.SetKeyImage(keyImage)

	return pc.GetKeyImage(), nil
}

// ConcealOutputCoin conceals the true value of a PlainCoinV1, leaving only the serial number unchanged.
func (pc *PlainCoinV1) ConcealOutputCoin(_ interface{}) error {
	pc.SetCommitment(nil)
	pc.SetValue(0)
	pc.SetSNDerivator(nil)
	pc.SetPublicKey(nil)
	pc.SetRandomness(nil)
	return nil
}

// MarshalJSON converts coin to a byte-array.
func (pc PlainCoinV1) MarshalJSON() ([]byte, error) {
	data := pc.Bytes()
	temp := base58.Base58Check{}.Encode(data, common.ZeroByte)
	return json.Marshal(temp)
}

// UnmarshalJSON converts a slice of bytes (was Marshalled before) into a PlainCoinV1 objects.
func (pc *PlainCoinV1) UnmarshalJSON(data []byte) error {
	dataStr := ""
	_ = json.Unmarshal(data, &dataStr)
	temp, _, err := base58.Base58Check{}.Decode(dataStr)
	if err != nil {
		return err
	}
	err = pc.SetBytes(temp)
	if err != nil {
		return err
	}
	return nil
}

// HashH returns the SHA3-256 hashing of a PlainCoinV1.
func (pc *PlainCoinV1) HashH() *common.Hash {
	hash := common.HashH(pc.Bytes())
	return &hash
}

// CommitAll calculates the Pedersen commitment of a PlainCoinV1 from its attributes.
//
// The commitment includes 5 attributes with 5 different bases:
//	- The public key
//	- The value
//	- The serial number derivator
//	- The shardID
//	- The randomness.
func (pc *PlainCoinV1) CommitAll() error {
	shardID, err := pc.GetShardID()
	if err != nil {
		return err
	}
	values := []*crypto.Scalar{
		new(crypto.Scalar).FromUint64(0),
		new(crypto.Scalar).FromUint64(pc.value),
		pc.snDerivator,
		new(crypto.Scalar).FromUint64(uint64(shardID)),
		pc.randomness,
	}
	pc.commitment, err = crypto.PedCom.CommitAll(values)
	if err != nil {
		return err
	}
	pc.commitment.Add(pc.commitment, pc.publicKey)

	return nil
}

// Bytes converts a PlainCoinV1 into a slice of bytes.
//
// This conversion is unique for each PlainCoinV1.
func (pc *PlainCoinV1) Bytes() []byte {
	var coinBytes []byte

	if pc.publicKey != nil {
		publicKey := pc.publicKey.ToBytesS()
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, publicKey...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if pc.commitment != nil {
		commitment := pc.commitment.ToBytesS()
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, commitment...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if pc.snDerivator != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, pc.snDerivator.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if pc.serialNumber != nil {
		serialNumber := pc.serialNumber.ToBytesS()
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, serialNumber...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if pc.randomness != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, pc.randomness.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if pc.value > 0 {
		value := new(big.Int).SetUint64(pc.value).Bytes()
		coinBytes = append(coinBytes, byte(len(value)))
		coinBytes = append(coinBytes, value...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if len(pc.info) > 0 {
		byteLengthInfo := byte(getMin(len(pc.info), MaxSizeInfoCoin))
		coinBytes = append(coinBytes, byteLengthInfo)
		infoBytes := pc.info[0:byteLengthInfo]
		coinBytes = append(coinBytes, infoBytes...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	return coinBytes
}

// SetBytes parses a slice of bytes into a PlainCoinV1.
func (pc *PlainCoinV1) SetBytes(coinBytes []byte) error {
	if len(coinBytes) == 0 {
		return fmt.Errorf("coinBytes is empty")
	}
	var err error

	offset := 0
	pc.publicKey, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV1 publicKey error: " + err.Error())
	}
	pc.commitment, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV1 commitment error: " + err.Error())
	}
	pc.snDerivator, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV1 snDerivator error: " + err.Error())
	}
	pc.serialNumber, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV1 serialNumber error: " + err.Error())
	}
	pc.randomness, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV1 serialNumber error: " + err.Error())
	}

	if offset >= len(coinBytes) {
		return fmt.Errorf("SetBytes CoinV1: out of range Parse value")
	}
	lenField := coinBytes[offset]
	offset++
	if lenField != 0 {
		if offset+int(lenField) > len(coinBytes) {
			// out of range
			return fmt.Errorf("out of range Parse PublicKey")
		}
		pc.value = new(big.Int).SetBytes(coinBytes[offset : offset+int(lenField)]).Uint64()
		offset += int(lenField)
	}

	pc.info, err = parseInfoForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV1 info error: " + err.Error())
	}
	return nil
}

// DoesCoinBelongToKeySet checks if a PlainCoinV1 belongs to the given key set.
func (pc *PlainCoinV1) DoesCoinBelongToKeySet(keySet *key.KeySet) (bool, *crypto.Point) {
	return crypto.IsPointEqual(keySet.PaymentAddress.GetPublicSpend(), pc.GetPublicKey()), nil
}

// CoinV1 implements the Coin interface. It is mainly used as an output coin of a transaction v1.
//
// It contains CoinDetails and CoinDetailsEncrypted (encrypted value and randomness).
// CoinDetailsEncrypted is nil when you send tx without privacy.
type CoinV1 struct {
	CoinDetails          *PlainCoinV1
	CoinDetailsEncrypted *hybridencryption.HybridCipherText
}

// Init initializes a new CoinV1.
func (c *CoinV1) Init() *CoinV1 {
	c.CoinDetails = new(PlainCoinV1).Init()
	c.CoinDetailsEncrypted = new(hybridencryption.HybridCipherText)
	return c
}

// GetVersion returns the version of a CoinV1.
func (c CoinV1) GetVersion() uint8 { return 1 }

// GetPublicKey returns the public key of a CoinV1. For a CoinV1, its public key is the public key of the owner.
func (c CoinV1) GetPublicKey() *crypto.Point { return c.CoinDetails.GetPublicKey() }

// GetCommitment returns the commitment of a CoinV1.
func (c CoinV1) GetCommitment() *crypto.Point { return c.CoinDetails.GetCommitment() }

// GetKeyImage returns the serial number of a CoinV1.
func (c CoinV1) GetKeyImage() *crypto.Point { return c.CoinDetails.GetKeyImage() }

// GetRandomness returns the randomness of a CoinV1.
func (c CoinV1) GetRandomness() *crypto.Scalar { return c.CoinDetails.GetRandomness() }

// GetSNDerivator returns the serial number derivator of a CoinV1.
func (c CoinV1) GetSNDerivator() *crypto.Scalar { return c.CoinDetails.GetSNDerivator() }

// GetShardID returns the shardID in which a CoinV1 belongs to.
func (c CoinV1) GetShardID() (uint8, error) { return c.CoinDetails.GetShardID() }

// GetValue returns the value of a CoinV1.
func (c CoinV1) GetValue() uint64 { return c.CoinDetails.GetValue() }

// GetInfo returns the info of a CoinV1.
func (c CoinV1) GetInfo() []byte { return c.CoinDetails.GetInfo() }

// IsEncrypted checks if whether a CoinV1 is encrypted.
func (c CoinV1) IsEncrypted() bool { return c.CoinDetailsEncrypted != nil }

// GetTxRandom returns the transaction randomness of a CoinV1. For a CoinV1, this value is nil.
func (c CoinV1) GetTxRandom() *TxRandom { return nil }

// GetSharedRandom returns the shared random of a CoinV1. For a CoinV1, this value is nil.
func (c CoinV1) GetSharedRandom() *crypto.Scalar { return nil }

// GetSharedConcealRandom returns the shared random when concealing of a CoinV1. For a CoinV1, this value is nil.
func (c CoinV1) GetSharedConcealRandom() *crypto.Scalar { return nil }

// GetAssetTag returns the asset tag of a CoinV1. For a CoinV1, this value is nil.
func (c CoinV1) GetAssetTag() *crypto.Point { return nil }

// GetCoinDetailEncrypted returns the encrypted detail of a CoinV1. For a non-private transaction, this value is always nil.
func (c CoinV1) GetCoinDetailEncrypted() []byte {
	if c.CoinDetailsEncrypted != nil {
		return c.CoinDetailsEncrypted.Bytes()
	}
	return nil
}

// SetValue sets v as the value of a CoinV1.
func (c CoinV1) SetValue(v uint64) { c.CoinDetails.value = v }

// Bytes converts a CoinV1 into a slice of bytes.
//
// This conversion is unique for each CoinV1.
func (c *CoinV1) Bytes() []byte {
	var outCoinBytes []byte

	if c.CoinDetailsEncrypted != nil {
		coinDetailsEncryptedBytes := c.CoinDetailsEncrypted.Bytes()
		outCoinBytes = append(outCoinBytes, byte(len(coinDetailsEncryptedBytes)))
		outCoinBytes = append(outCoinBytes, coinDetailsEncryptedBytes...)
	} else {
		outCoinBytes = append(outCoinBytes, byte(0))
	}

	coinDetailBytes := c.CoinDetails.Bytes()

	lenCoinDetailBytes := make([]byte, 0)
	if len(coinDetailBytes) <= 255 {
		lenCoinDetailBytes = []byte{byte(len(coinDetailBytes))}
	} else {
		lenCoinDetailBytes = common.IntToBytes(len(coinDetailBytes))
	}

	outCoinBytes = append(outCoinBytes, lenCoinDetailBytes...)
	outCoinBytes = append(outCoinBytes, coinDetailBytes...)
	return outCoinBytes
}

// SetBytes parses a slice of bytes into a CoinV1.
func (c *CoinV1) SetBytes(bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("coinBytes is empty")
	}

	offset := 0
	lenCoinDetailEncrypted := int(bytes[0])
	offset += 1

	if lenCoinDetailEncrypted > 0 {
		if offset+lenCoinDetailEncrypted > len(bytes) {
			// out of range
			return fmt.Errorf("out of range Parse CoinDetailsEncrypted")
		}
		c.CoinDetailsEncrypted = new(hybridencryption.HybridCipherText)
		err := c.CoinDetailsEncrypted.SetBytes(bytes[offset : offset+lenCoinDetailEncrypted])
		if err != nil {
			return err
		}
		offset += lenCoinDetailEncrypted
	}

	// try get 1-byte for len
	if offset >= len(bytes) {
		// out of range
		return fmt.Errorf("out of range Parse CoinDetails")
	}
	lenOutputCoin := int(bytes[offset])
	c.CoinDetails = new(PlainCoinV1)
	if lenOutputCoin != 0 {
		offset += 1
		if offset+lenOutputCoin > len(bytes) {
			// out of range
			return fmt.Errorf("out of range Parse output coin details 1")
		}
		err := c.CoinDetails.SetBytes(bytes[offset : offset+lenOutputCoin])
		if err != nil {
			// 1-byte is wrong
			// try get 2-byte for len
			if offset+1 > len(bytes) {
				// out of range
				return fmt.Errorf("out of range Parse output coin details 2 ")
			}
			lenOutputCoin = common.BytesToInt(bytes[offset-1 : offset+1])
			offset += 1
			if offset+lenOutputCoin > len(bytes) {
				// out of range
				return fmt.Errorf("out of range Parse output coin details 3 ")
			}
			err1 := c.CoinDetails.SetBytes(bytes[offset : offset+lenOutputCoin])
			return err1
		}
	} else {
		// 1-byte is wrong
		// try get 2-byte for len
		if offset+2 > len(bytes) {
			// out of range
			return fmt.Errorf("out of range Parse output coin details 4")
		}
		lenOutputCoin = common.BytesToInt(bytes[offset : offset+2])
		offset += 2
		if offset+lenOutputCoin > len(bytes) {
			// out of range
			return fmt.Errorf("out of range Parse output coin details 5")
		}
		err1 := c.CoinDetails.SetBytes(bytes[offset : offset+lenOutputCoin])
		return err1
	}

	return nil
}

// Encrypt returns a ciphertext encrypting for a coin using a hybrid crypto-system,
// in which the  AES encryption scheme is used as a data encapsulation scheme,
// and the ElGamal crypto system is used as a key encapsulation scheme.
func (c *CoinV1) Encrypt(recipientTK key.TransmissionKey) error {
	// 32-byte first: Randomness, the rest of msg is value of coin
	msg := append(c.CoinDetails.randomness.ToBytesS(), new(big.Int).SetUint64(c.CoinDetails.value).Bytes()...)

	pubKeyPoint, err := new(crypto.Point).FromBytesS(recipientTK)
	if err != nil {
		return err
	}

	c.CoinDetailsEncrypted, err = hybridencryption.HybridEncrypt(msg, pubKeyPoint)
	if err != nil {
		return err
	}

	return nil
}

// Decrypt decrypts a CoinV1 into a PlainCoinV1 using the given key set.
func (c CoinV1) Decrypt(keySet *key.KeySet) (PlainCoin, error) {
	if keySet == nil {
		err := fmt.Errorf("cannot decrypt coinv1 with empty key")
		return nil, err
	}

	if len(keySet.ReadonlyKey.Rk) == 0 && len(keySet.PrivateKey) == 0 {
		err := fmt.Errorf("cannot Decrypt CoinV1: Keyset does not contain viewkey or privatekey")
		return nil, err
	}

	if bytes.Equal(c.GetPublicKey().ToBytesS(), keySet.PaymentAddress.Pk[:]) {
		result := &CoinV1{
			CoinDetails:          c.CoinDetails,
			CoinDetailsEncrypted: c.CoinDetailsEncrypted,
		}
		if result.CoinDetailsEncrypted != nil && !result.CoinDetailsEncrypted.IsNil() {
			if len(keySet.ReadonlyKey.Rk) > 0 {
				msg, err := hybridencryption.HybridDecrypt(c.CoinDetailsEncrypted, new(crypto.Scalar).FromBytesS(keySet.ReadonlyKey.Rk))
				if err != nil {
					return nil, err
				}
				// Assign randomness and value to outputCoin details
				result.CoinDetails.randomness = new(crypto.Scalar).FromBytesS(msg[0:crypto.Ed25519KeySize])
				result.CoinDetails.value = new(big.Int).SetBytes(msg[crypto.Ed25519KeySize:]).Uint64()

				//re-calculate commitment
				shardID, err := result.CoinDetails.GetShardID()
				if err != nil {
					return nil, fmt.Errorf("cannot get shardID of coin")
				}
				tmpCmt := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenValueIndex], new(crypto.Scalar).FromUint64(result.CoinDetails.value))
				tmpCmt = tmpCmt.Add(tmpCmt, new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenSndIndex], result.CoinDetails.GetSNDerivator()))
				tmpCmt = tmpCmt.Add(tmpCmt, new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenShardIDIndex], new(crypto.Scalar).FromUint64(uint64(shardID))))
				tmpCmt = tmpCmt.Add(tmpCmt, new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenRandomnessIndex], result.CoinDetails.GetRandomness()))
				tmpCmt = tmpCmt.Add(tmpCmt, result.CoinDetails.GetPublicKey())

				if !crypto.IsPointEqual(tmpCmt, result.CoinDetails.GetCommitment()) {
					tmpCmtStr := base58.Base58Check{}.Encode(tmpCmt.ToBytesS(), common.ZeroByte)
					cmtStr := base58.Base58Check{}.Encode(result.CoinDetails.GetCommitment().ToBytesS(), common.ZeroByte)
					return nil, fmt.Errorf("expected commitment %v, got %v", cmtStr, tmpCmtStr)
				}
			}
		}
		if len(keySet.PrivateKey) > 0 {
			// check spent with private key
			keyImage := new(crypto.Point).Derive(
				crypto.PedCom.G[crypto.PedersenPrivateKeyIndex],
				new(crypto.Scalar).FromBytesS(keySet.PrivateKey),
				result.CoinDetails.GetSNDerivator())
			result.CoinDetails.SetKeyImage(keyImage)
		}
		return result.CoinDetails, nil
	}
	err := fmt.Errorf("coin publicKey does not equal keyset paymentAddress")
	return nil, err
}

// CheckCoinValid checks if a CoinV1 is valid for its amount and payment address.
func (c *CoinV1) CheckCoinValid(paymentAdd key.PaymentAddress, _ []byte, amount uint64) bool {
	return bytes.Equal(c.GetPublicKey().ToBytesS(), paymentAdd.GetPublicSpend().ToBytesS()) && amount == c.GetValue()
}

// DoesCoinBelongToKeySet checks if a CoinV1 belongs to the given key set.
func (c *CoinV1) DoesCoinBelongToKeySet(keySet *key.KeySet) (bool, *crypto.Point) {
	return crypto.IsPointEqual(keySet.PaymentAddress.GetPublicSpend(), c.GetPublicKey()), nil
}

// MarshalJSON converts coin to a byte-array.
func (c CoinV1) MarshalJSON() ([]byte, error) {
	data := c.Bytes()
	temp := base58.Base58Check{}.Encode(data, common.ZeroByte)
	return json.Marshal(temp)
}

// UnmarshalJSON converts a slice of bytes (was Marshalled before) into a CoinV1 objects.
func (c *CoinV1) UnmarshalJSON(data []byte) error {
	dataStr := ""
	_ = json.Unmarshal(data, &dataStr)
	temp, _, err := base58.Base58Check{}.Decode(dataStr)
	if err != nil {
		return err
	}
	err = c.SetBytes(temp)
	if err != nil {
		return err
	}
	return nil
}
