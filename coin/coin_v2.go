package coin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
)

// TxRandom is a struct implementing for transactions of version 2.
//
// A TxRandom consists of 3 elements, represented an array of TxRandomGroupSize bytes:
//	- An OTA random point
//	- A concealing random point
//	- An index
type TxRandom [TxRandomGroupSize]byte

// NewTxRandom initializes a new TxRandom.
func NewTxRandom() *TxRandom {
	txRandom := new(crypto.Point).Identity()
	index := uint32(0)

	res := new(TxRandom)
	res.SetTxConcealRandomPoint(txRandom)
	res.SetIndex(index)
	return res
}

// GetTxConcealRandomPoint returns the conceal random point of a TxRandom.
func (t TxRandom) GetTxConcealRandomPoint() (*crypto.Point, error) {
	return new(crypto.Point).FromBytesS(t[crypto.Ed25519KeySize+4:])
}

// GetTxOTARandomPoint returns the OTA random point of a TxRandom.
func (t TxRandom) GetTxOTARandomPoint() (*crypto.Point, error) {
	return new(crypto.Point).FromBytesS(t[:crypto.Ed25519KeySize])
}

// GetIndex returns the index of a TxRandom.
func (t TxRandom) GetIndex() (uint32, error) {
	return common.BytesToUint32(t[crypto.Ed25519KeySize : crypto.Ed25519KeySize+4])
}

// SetTxConcealRandomPoint sets v as the conceal random point of a TxRandom.
func (t *TxRandom) SetTxConcealRandomPoint(v *crypto.Point) {
	txRandomBytes := v.ToBytesS()
	copy(t[crypto.Ed25519KeySize+4:], txRandomBytes)
}

// SetTxOTARandomPoint sets v as the OTA random point of a TxRandom.
func (t *TxRandom) SetTxOTARandomPoint(v *crypto.Point) {
	txRandomBytes := v.ToBytesS()
	copy(t[:crypto.Ed25519KeySize], txRandomBytes)
}

// SetIndex sets v as the index of a TxRandom.
func (t *TxRandom) SetIndex(v uint32) {
	indexBytes := common.Uint32ToBytes(v)
	copy(t[crypto.Ed25519KeySize:], indexBytes)
}

// Bytes converts a TxRandom into a slice of bytes.
func (t TxRandom) Bytes() []byte {
	return t[:]
}

// SetBytes sets the content of a slice of bytes b to a TxRandom.
func (t *TxRandom) SetBytes(v []byte) error {
	if v == nil || len(v) != TxRandomGroupSize {
		return fmt.Errorf("cannot SetByte to TxRandom. Input is invalid")
	}
	_, err := new(crypto.Point).FromBytesS(v[:crypto.Ed25519KeySize])
	if err != nil {
		errStr := fmt.Sprintf("cannot TxRandomGroupSize.SetBytes: bytes is not crypto.Point err: %v", err)
		return fmt.Errorf(errStr)
	}
	_, err = new(crypto.Point).FromBytesS(v[crypto.Ed25519KeySize+4:])
	if err != nil {
		errStr := fmt.Sprintf("cannot TxRandomGroupSize.SetBytes: bytes is not crypto.Point err: %v", err)
		return fmt.Errorf(errStr)
	}
	copy(t[:], v)
	return nil
}

// CoinV2 implements both the PlainCoin and Coin interfaces. It is mainly used as inputs and outputs of a transaction v2.
//
// If not privacy, mask and amount will be the original randomness and value
// If has privacy, mask and amount will be as described in the Monero paper.
type CoinV2 struct {
	// Public
	version    uint8
	info       []byte
	publicKey  *crypto.Point
	commitment *crypto.Point
	keyImage   *crypto.Point

	// sharedRandom and txRandom is shared secret between receiver and giver
	// sharedRandom is only visible when creating coins, when it is broadcast to network, it will be set to null
	sharedConcealRandom *crypto.Scalar //rConceal: shared random when concealing output coin and blinding assetTag
	sharedRandom        *crypto.Scalar // rOTA: shared random when creating one-time-address
	txRandom            *TxRandom      // rConceal*G + rOTA*G + index

	// mask = randomness
	// amount = value
	mask   *crypto.Scalar
	amount *crypto.Scalar
	// tag is nil unless confidential asset
	assetTag *crypto.Point
	// the hash of the tokenID
	rawAssetTag *crypto.Point
}

// ParsePrivateKeyOfCoin sets privateKey as the private key of a CoinV2.
func (c CoinV2) ParsePrivateKeyOfCoin(privateKey key.PrivateKey) (*crypto.Scalar, error) {
	keySet := new(key.KeySet)
	if err := keySet.InitFromPrivateKey(&privateKey); err != nil {
		err := fmt.Errorf("cannot init keyset from privateKey")
		return nil, err
	}
	_, txRandomOTAPoint, index, err := c.GetTxRandomDetail()
	if err != nil {
		return nil, err
	}
	rK := new(crypto.Point).ScalarMult(txRandomOTAPoint, keySet.OTAKey.GetOTASecretKey()) //(r_ota*G) * k = r_ota * K
	H := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))       // Hash(r_ota*K, index)

	k := new(crypto.Scalar).FromBytesS(privateKey)
	return new(crypto.Scalar).Add(H, k), nil // Hash(rK, index) + privateSpend
}

// ParseKeyImageWithPrivateKey derives the key image of a CoinV2 from its private key.
func (c CoinV2) ParseKeyImageWithPrivateKey(privateKey key.PrivateKey) (*crypto.Point, error) {
	k, err := c.ParsePrivateKeyOfCoin(privateKey)
	if err != nil {
		return nil, err
	}
	Hp := crypto.HashToPoint(c.GetPublicKey().ToBytesS())
	return new(crypto.Point).ScalarMult(Hp, k), nil
}

// ConcealOutputCoin conceals the amount of coin using the publicView of the receiver
//
//	- AdditionalData: must be the publicView of the receiver.
func (c *CoinV2) ConcealOutputCoin(additionalData interface{}) error {
	// already encrypted
	if c.IsEncrypted() {
		return nil
	}

	paymentInfo, ok := additionalData.(*key.PaymentInfo)
	if !ok {
		return fmt.Errorf("expect additionalData to be a PaymentInfo")
	}

	var rK *crypto.Point
	if paymentInfo.OTAReceiver != "" {
		otaReceiver := new(OTAReceiver)
		_ = otaReceiver.FromString(paymentInfo.OTAReceiver) // error has been handled by callers

		rK = &otaReceiver.SharedSecrets[1]
	} else {
		// created by other person
		if c.GetSharedConcealRandom() == nil {
			return nil
		}

		// re-calculate the sharedConcealSecret
		rK = new(crypto.Point).ScalarMult(paymentInfo.PaymentAddress.GetPublicView(), c.GetSharedConcealRandom())
	}

	hash := crypto.HashToScalar(rK.ToBytesS()) //hash(rK)
	hash = crypto.HashToScalar(hash.ToBytesS())
	mask := new(crypto.Scalar).Add(c.GetRandomness(), hash) //mask = c.mask + hash

	hash = crypto.HashToScalar(hash.ToBytesS())
	amount := new(crypto.Scalar).Add(c.GetAmount(), hash) //amount = c.amount + hash
	c.SetRandomness(mask)
	c.SetAmount(amount)
	c.SetSharedConcealRandom(nil)
	c.SetSharedRandom(nil)
	return nil
}

// ConcealInputCoin conceals the true value of a CoinV2, leaving only the serial number unchanged (mainly when the coin is used as an input of a transaction).
func (c *CoinV2) ConcealInputCoin() {
	c.SetValue(0)
	c.SetRandomness(nil)
	c.SetPublicKey(nil)
	c.SetCommitment(nil)
	c.SetTxRandomDetail(new(crypto.Point).Identity(), new(crypto.Point).Identity(), 0)
}

// Decrypt decrypts a CoinV2 into a PlainCoin using the given key set.
func (c *CoinV2) Decrypt(keySet *key.KeySet) (PlainCoin, error) {
	if keySet == nil {
		err := fmt.Errorf("cannot Decrypt CoinV2: Keyset is empty")
		return nil, err
	}

	// Must parse keyImage first in any situation
	if len(keySet.PrivateKey) > 0 {
		keyImage, err := c.ParseKeyImageWithPrivateKey(keySet.PrivateKey)
		if err != nil {
			errReturn := fmt.Errorf("cannot parse key image with privateKey CoinV2" + err.Error())
			return nil, errReturn
		}
		c.SetKeyImage(keyImage)
	}

	if !c.IsEncrypted() {
		return c, nil
	}

	viewKey := keySet.ReadonlyKey
	if len(viewKey.Rk) == 0 && len(keySet.PrivateKey) == 0 {
		err := fmt.Errorf("cannot Decrypt CoinV2: Keyset does not contain viewkey or privatekey")
		return nil, err
	}

	if viewKey.GetPrivateView() != nil {
		txConcealRandomPoint, err := c.GetTxRandom().GetTxConcealRandomPoint()
		if err != nil {
			return nil, err
		}
		rK := new(crypto.Point).ScalarMult(txConcealRandomPoint, viewKey.GetPrivateView())

		// Hash multiple times
		hash := crypto.HashToScalar(rK.ToBytesS())
		hash = crypto.HashToScalar(hash.ToBytesS())
		randomness := c.GetRandomness().Sub(c.GetRandomness(), hash)

		// Hash 1 more time to get value
		hash = crypto.HashToScalar(hash.ToBytesS())
		value := c.GetAmount().Sub(c.GetAmount(), hash)

		commitment := crypto.PedCom.CommitAtIndex(value, randomness, crypto.PedersenValueIndex)
		// for `confidential asset` coin, we commit differently
		if c.GetAssetTag() != nil {
			com, err := ComputeCommitmentCA(c.GetAssetTag(), randomness, value)
			if err != nil {
				err := fmt.Errorf("cannot recompute commitment when decrypting")
				return nil, err
			}
			commitment = com
		}
		if !crypto.IsPointEqual(commitment, c.GetCommitment()) {
			err := fmt.Errorf("cannot Decrypt CoinV2: Commitment is not the same after decrypt")
			return nil, err
		}
		c.SetRandomness(randomness)
		c.SetAmount(value)
	}
	return c, nil
}

// Init (OutputCoin) initializes a output coin
func (c *CoinV2) Init() *CoinV2 {
	c.version = 2
	c.info = []byte{}
	c.publicKey = new(crypto.Point).Identity()
	c.commitment = new(crypto.Point).Identity()
	c.keyImage = new(crypto.Point).Identity()
	c.sharedRandom = new(crypto.Scalar)
	c.txRandom = NewTxRandom()
	c.mask = new(crypto.Scalar)
	c.amount = new(crypto.Scalar)
	return c
}

// GetSNDerivator returns the serial number derivator of a CoinV2. For a CoinV2, this value is always nil.
func (c CoinV2) GetSNDerivator() *crypto.Scalar { return nil }

// IsEncrypted checks if whether a CoinV2 is encrypted.
func (c CoinV2) IsEncrypted() bool {
	if c.mask == nil || c.amount == nil {
		return true
	}
	tempCommitment := crypto.PedCom.CommitAtIndex(c.amount, c.mask, crypto.PedersenValueIndex)
	if c.GetAssetTag() != nil {
		// err is only for nil parameters, which we already checked, so here it is ignored
		com, _ := c.ComputeCommitmentCA()
		tempCommitment = com
	}
	return !crypto.IsPointEqual(tempCommitment, c.commitment)
}

// GetVersion returns the version of a CoinV2.
func (c CoinV2) GetVersion() uint8 { return 2 }

// GetRandomness returns the randomness of a CoinV2.
func (c CoinV2) GetRandomness() *crypto.Scalar { return c.mask }

func (c CoinV2) GetAmount() *crypto.Scalar { return c.amount }

// GetSharedRandom returns the shared random of a CoinV2.
func (c CoinV2) GetSharedRandom() *crypto.Scalar { return c.sharedRandom }

// GetSharedConcealRandom returns the shared random when concealing of a CoinV2.
func (c CoinV2) GetSharedConcealRandom() *crypto.Scalar { return c.sharedConcealRandom }

// GetPublicKey returns the public key of a CoinV2.
func (c CoinV2) GetPublicKey() *crypto.Point { return c.publicKey }

// GetCommitment returns the commitment of a CoinV2.
func (c CoinV2) GetCommitment() *crypto.Point { return c.commitment }

// GetKeyImage returns the key image of a CoinV2.
func (c CoinV2) GetKeyImage() *crypto.Point { return c.keyImage }

// GetInfo returns the info of a CoinV2.
func (c CoinV2) GetInfo() []byte { return c.info }

// GetAssetTag returns the asset tag of a CoinV2. For a PRV CoinV2, this value is nil.
func (c CoinV2) GetAssetTag() *crypto.Point { return c.assetTag }

// GetValue returns the value of a CoinV2.
func (c CoinV2) GetValue() uint64 {
	if c.IsEncrypted() {
		return 0
	}
	return c.amount.ToUint64Little()
}

// GetTxRandom returns the transaction randomness of a CoinV2.
func (c CoinV2) GetTxRandom() *TxRandom { return c.txRandom }

// GetTxRandomDetail returns the transaction randomness detail of a CoinV2.
func (c CoinV2) GetTxRandomDetail() (*crypto.Point, *crypto.Point, uint32, error) {
	txRandomOTAPoint, err1 := c.txRandom.GetTxOTARandomPoint()
	txRandomConcealPoint, err2 := c.txRandom.GetTxConcealRandomPoint()
	index, err3 := c.txRandom.GetIndex()
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, nil, 0, fmt.Errorf("cannot Get TxRandom: point or index is wrong")
	}
	return txRandomConcealPoint, txRandomOTAPoint, index, nil
}

// GetShardID returns the shardID in which a CoinV2 belongs to.
func (c CoinV2) GetShardID() (uint8, error) {
	if c.publicKey == nil {
		return 255, fmt.Errorf("cannot get GetShardID because GetPublicKey of PlainCoin is concealed")
	}
	pubKeyBytes := c.publicKey.ToBytes()
	lastByte := pubKeyBytes[crypto.Ed25519KeySize-1]
	shardID := common.GetShardIDFromLastByte(lastByte)
	return shardID, nil
}

// GetCoinDetailEncrypted returns the encrypted detail of a CoinV2.
func (c CoinV2) GetCoinDetailEncrypted() []byte {
	return c.GetAmount().ToBytesS()
}

// SetVersion sets the version of a CoinV2 to 2.
func (c *CoinV2) SetVersion(uint8) { c.version = 2 }

// SetRandomness sets v as the randomness of a CoinV2.
func (c *CoinV2) SetRandomness(v *crypto.Scalar) { c.mask = v }

// SetAmount sets v as the amount of a CoinV2.
func (c *CoinV2) SetAmount(v *crypto.Scalar) { c.amount = v }

// SetSharedRandom sets v as the OTA shared random of a CoinV2.
func (c *CoinV2) SetSharedRandom(v *crypto.Scalar) { c.sharedRandom = v }

// SetSharedConcealRandom sets v as the shared conceal random of a CoinV2.
func (c *CoinV2) SetSharedConcealRandom(v *crypto.Scalar) {
	c.sharedConcealRandom = v
}

// SetTxRandom sets v as the TxRandom of a CoinV2.
func (c *CoinV2) SetTxRandom(v *TxRandom) {
	if v == nil {
		c.txRandom = nil
	} else {
		c.txRandom = NewTxRandom()
		err := c.txRandom.SetBytes(v.Bytes())
		if err != nil {
			return
		}
	}
}

// SetTxRandomDetail creates the TxRandom of a CoinV2 based on the three provided parameters.
func (c *CoinV2) SetTxRandomDetail(txConcealRandomPoint, txRandomPoint *crypto.Point, index uint32) {
	res := new(TxRandom)
	res.SetTxConcealRandomPoint(txConcealRandomPoint)
	res.SetTxOTARandomPoint(txRandomPoint)
	res.SetIndex(index)
	c.txRandom = res
}

// SetPublicKey sets v as the public key of a CoinV2. Each CoinV2 of a user will now have different public keys.
func (c *CoinV2) SetPublicKey(v *crypto.Point) { c.publicKey = v }

// SetCommitment sets v as the randomness of a CoinV2.
func (c *CoinV2) SetCommitment(v *crypto.Point) { c.commitment = v }

// SetKeyImage sets v as the key-image (a.k.a serial number) of a CoinV2.
func (c *CoinV2) SetKeyImage(keyImage *crypto.Point) { c.keyImage = keyImage }

// SetValue sets v as the value of a CoinV2.
func (c *CoinV2) SetValue(value uint64) { c.amount = new(crypto.Scalar).FromUint64(value) }

// SetInfo sets v as the info of a CoinV2.
func (c *CoinV2) SetInfo(b []byte) {
	c.info = make([]byte, len(b))
	copy(c.info, b)
}

// SetAssetTag sets v as the asset tag of a CoinV2.
func (c *CoinV2) SetAssetTag(v *crypto.Point) { c.assetTag = v }

// Bytes converts a CoinV2 into a slice of bytes.
//
// This conversion is unique for each CoinV2.
func (c CoinV2) Bytes() []byte {
	coinBytes := []byte{c.GetVersion()}
	info := c.GetInfo()
	byteLengthInfo := byte(getMin(len(info), MaxSizeInfoCoin))
	coinBytes = append(coinBytes, byteLengthInfo)
	coinBytes = append(coinBytes, info[:byteLengthInfo]...)

	if c.publicKey != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, c.publicKey.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.commitment != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, c.commitment.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.keyImage != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, c.keyImage.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.sharedRandom != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, c.sharedRandom.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.sharedConcealRandom != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, c.sharedConcealRandom.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.txRandom != nil {
		coinBytes = append(coinBytes, TxRandomGroupSize)
		coinBytes = append(coinBytes, c.txRandom.Bytes()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.mask != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, c.mask.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.amount != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, c.amount.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.assetTag != nil {
		coinBytes = append(coinBytes, byte(crypto.Ed25519KeySize))
		coinBytes = append(coinBytes, c.assetTag.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	return coinBytes
}

// SetBytes parses a slice of bytes into a CoinV2.
func (c *CoinV2) SetBytes(coinBytes []byte) error {
	var err error
	if c == nil {
		return fmt.Errorf("cannot set bytes for unallocated CoinV2")
	}
	if len(coinBytes) == 0 {
		return fmt.Errorf("coinBytes is empty")
	}
	if coinBytes[0] != 2 {
		return fmt.Errorf("coinBytes version is not 2")
	}
	c.SetVersion(coinBytes[0])

	offset := 1
	c.info, err = parseInfoForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 info error: " + err.Error())
	}

	c.publicKey, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 publicKey error: " + err.Error())
	}
	c.commitment, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 commitment error: " + err.Error())
	}
	c.keyImage, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 keyImage error: " + err.Error())
	}
	c.sharedRandom, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 mask error: " + err.Error())
	}

	c.sharedConcealRandom, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 mask error: " + err.Error())
	}

	if offset >= len(coinBytes) {
		return fmt.Errorf("offset is larger than len(bytes), cannot parse txRandom")
	}
	if coinBytes[offset] != TxRandomGroupSize {
		return fmt.Errorf("SetBytes CoinV2 TxRandomGroup error: length of TxRandomGroup is not correct")
	}
	offset += 1
	if offset+TxRandomGroupSize > len(coinBytes) {
		return fmt.Errorf("SetBytes CoinV2 TxRandomGroup error: length of coinBytes is too small")
	}
	c.txRandom = NewTxRandom()
	err = c.txRandom.SetBytes(coinBytes[offset : offset+TxRandomGroupSize])
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 TxRandomGroup error: " + err.Error())
	}
	offset += TxRandomGroupSize

	c.mask, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 mask error: " + err.Error())
	}
	c.amount, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("SetBytes CoinV2 amount error: " + err.Error())
	}

	if offset >= len(coinBytes) {
		// for parsing old serialization, which does not have assetTag field
		c.assetTag = nil
	} else {
		c.assetTag, err = parsePointForSetBytes(&coinBytes, &offset)
		if err != nil {
			return fmt.Errorf("SetBytes CoinV2 assetTag error: " + err.Error())
		}
	}
	return nil
}

// HashH returns the SHA3-256 hashing of coin bytes array
func (c *CoinV2) HashH() *common.Hash {
	hash := common.HashH(c.Bytes())
	return &hash
}

// MarshalJSON converts coin to a byte-array.
func (c CoinV2) MarshalJSON() ([]byte, error) {
	data := c.Bytes()
	temp := base58.Base58Check{}.Encode(data, common.ZeroByte)
	return json.Marshal(temp)
}

// UnmarshalJSON converts a slice of bytes (was Marshalled before) into a CoinV2 objects.
func (c *CoinV2) UnmarshalJSON(data []byte) error {
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

// CheckCoinValid checks if a CoinV2 is valid for its amount and payment address.
func (c *CoinV2) CheckCoinValid(paymentAdd key.PaymentAddress, sharedRandom []byte, amount uint64) bool {
	if c.GetValue() != amount {
		return false
	}
	// check one-time address is corresponding to the payment address
	r := new(crypto.Scalar).FromBytesS(sharedRandom)
	if !r.ScalarValid() {
		return false
	}

	publicOTA := paymentAdd.GetOTAPublicKey()
	if publicOTA == nil {
		return false
	}
	rK := new(crypto.Point).ScalarMult(publicOTA, r)
	_, txOTARandomPoint, index, err := c.GetTxRandomDetail()
	if err != nil {
		return false
	}
	if !crypto.IsPointEqual(new(crypto.Point).ScalarMultBase(r), txOTARandomPoint) {
		return false
	}

	hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
	HrKG := new(crypto.Point).ScalarMultBase(hash)
	tmpPubKey := new(crypto.Point).Add(HrKG, paymentAdd.GetPublicSpend())
	return bytes.Equal(tmpPubKey.ToBytesS(), c.publicKey.ToBytesS())
}

// DoesCoinBelongToKeySet checks if a CoinV2 belongs to the given key set. If so, it will also try to calculate the raw
// asset tag associated with the coin.
func (c *CoinV2) DoesCoinBelongToKeySet(keySet *key.KeySet) (bool, *crypto.Point) {
	_, txOTARandomPoint, index, err1 := c.GetTxRandomDetail()
	if err1 != nil {
		return false, nil
	}

	//Check if the utxo belong to this one-time-address
	rK := new(crypto.Point).ScalarMult(txOTARandomPoint, keySet.OTAKey.GetOTASecretKey())

	hashed := crypto.HashToScalar(
		append(rK.ToBytesS(), common.Uint32ToBytes(index)...),
	)

	HnG := new(crypto.Point).ScalarMultBase(hashed)
	KCheck := new(crypto.Point).Sub(c.GetPublicKey(), HnG)

	belongs := crypto.IsPointEqual(KCheck, keySet.OTAKey.GetPublicSpend())
	if !belongs {
		return false, nil
	}

	// try to calculate the raw asset tag (i.e, the hash of the real tokenID)
	if c.GetAssetTag() != nil && c.rawAssetTag == nil {
		blinder := crypto.HashToScalar(append(rK.ToBytesS(), []byte("assettag")...))
		rawAssetTag := new(crypto.Point).Sub(
			c.GetAssetTag(),
			new(crypto.Point).ScalarMult(crypto.PedCom.G[PedersenRandomnessIndex], blinder),
		)
		c.rawAssetTag = rawAssetTag
	}

	return true, rK
}

// GetTokenId attempts to retrieve the asset a CoinV2.
// Parameters:
// 	- keySet: the key set of the user, must contain an OTAKey
//	- rawAssetTags: a pre-computed mapping from a raw assetTag to the tokenId (e.g, HashToPoint(PRV) => PRV).
func (c *CoinV2) GetTokenId(keySet *key.KeySet, rawAssetTags map[string]*common.Hash) (*common.Hash, error) {
	if c.rawAssetTag != nil {
		if asset, ok := rawAssetTags[c.rawAssetTag.String()]; ok {
			return asset, nil
		}
	}

	if c.GetAssetTag() == nil {
		return &common.PRVCoinID, nil
	}

	if asset, ok := rawAssetTags[c.GetAssetTag().String()]; ok {
		return asset, nil
	}

	belong, _ := c.DoesCoinBelongToKeySet(keySet)
	if !belong {
		return nil, fmt.Errorf("coin does not belong to the keyset")
	}

	rawAssetTag := c.rawAssetTag
	if rawAssetTag == nil {
		return nil, fmt.Errorf("cannot calculate the raw asset tag")
	}
	if asset, ok := rawAssetTags[rawAssetTag.String()]; ok {
		return asset, nil
	}

	return nil, fmt.Errorf("cannot find the tokenId")
}
