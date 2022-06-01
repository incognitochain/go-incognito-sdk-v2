package coin

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

const (
	MaxSizeInfoCoin   = 255
	JsonMarshalFlag   = 34
	CoinVersion1      = 1
	CoinVersion2      = 2
	TxRandomGroupSize = 68
)

const (
	PedersenPrivateKeyIndex = crypto.PedersenPrivateKeyIndex
	PedersenValueIndex      = crypto.PedersenValueIndex
	PedersenSndIndex        = crypto.PedersenSndIndex
	PedersenShardIDIndex    = crypto.PedersenShardIDIndex
	PedersenRandomnessIndex = crypto.PedersenRandomnessIndex
)

// getMin returns the min of a and b.
func getMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// parseScalarForSetBytes parses a slice of bytes beginning with the given offset into a crypto.Scalar object.
func parseScalarForSetBytes(coinBytes *[]byte, offset *int) (*crypto.Scalar, error) {
	b := *coinBytes
	if *offset >= len(b) {
		return nil, fmt.Errorf("offset is larger than len(bytes), cannot parse scalar")
	}
	var sc *crypto.Scalar = nil
	lenField := b[*offset]
	*offset += 1
	if lenField != 0 {
		if *offset+int(lenField) > len(b) {
			return nil, fmt.Errorf("Offset+curLen is larger than len(bytes), cannot parse scalar for set bytes")
		}
		data := b[*offset : *offset+int(lenField)]
		sc = new(crypto.Scalar).FromBytesS(data)
		*offset += int(lenField)
	}
	return sc, nil
}

// parsePointForSetBytes parses a slice of bytes beginning with the given offset into a crypto.Point object.
func parsePointForSetBytes(coinBytes *[]byte, offset *int) (*crypto.Point, error) {
	b := *coinBytes
	if *offset >= len(b) {
		return nil, fmt.Errorf("offset is larger than len(bytes), cannot parse point")
	}
	var point *crypto.Point = nil
	var err error
	lenField := b[*offset]
	*offset += 1
	if lenField != 0 {
		if *offset+int(lenField) > len(b) {
			return nil, fmt.Errorf("offset+curLen is larger than len(bytes), cannot parse point for set bytes")
		}
		data := b[*offset : *offset+int(lenField)]
		point, err = new(crypto.Point).FromBytesS(data)
		if err != nil {
			return nil, err
		}
		*offset += int(lenField)
	}
	return point, nil
}

// parseInfoForSetBytes parses a slice of bytes at the given offset into the coin info.
func parseInfoForSetBytes(coinBytes *[]byte, offset *int) ([]byte, error) {
	b := *coinBytes
	if *offset >= len(b) {
		return []byte{}, fmt.Errorf("offset is larger than len(bytes), cannot parse info")
	}
	var info []byte
	lenField := b[*offset]
	*offset += 1
	if lenField != 0 {
		if *offset+int(lenField) > len(b) {
			return []byte{}, fmt.Errorf("offset+curLen is larger than len(bytes), cannot parse info for set bytes")
		}
		info = make([]byte, lenField)
		copy(info, b[*offset:*offset+int(lenField)])
		*offset += int(lenField)
	}
	return info, nil
}

// NewCoinFromOTAReceiver creates a new CoinV2 from a given OTA receiver, amount and info.
func NewCoinFromOTAReceiver(otaReceiver OTAReceiver, amount uint64, info []byte) *CoinV2 {
	c := new(CoinV2).Init()
	c.SetPublicKey(&otaReceiver.PublicKey)
	c.SetAmount(new(crypto.Scalar).FromUint64(amount))
	c.SetRandomness(crypto.RandomScalar())
	c.SetTxRandom(&otaReceiver.TxRandom)
	c.SetCommitment(crypto.PedCom.CommitAtIndex(c.GetAmount(), c.GetRandomness(), crypto.PedersenValueIndex))
	c.SetSharedRandom(nil)
	c.SetSharedConcealRandom(nil)
	c.SetInfo(info)
	return c
}

// NewCoinFromPaymentInfo creates a new CoinV2 from the given payment info.
func NewCoinFromPaymentInfo(info *key.PaymentInfo) (*CoinV2, error) {
	if info.OTAReceiver != "" {
		otaReceiver := new(OTAReceiver)
		err := otaReceiver.FromString(info.OTAReceiver)
		if err != nil {
			return nil, fmt.Errorf("invalid otaReceiver %v: %v", info.OTAReceiver, err)
		}

		return NewCoinFromOTAReceiver(*otaReceiver, info.Amount, info.Message), nil
	}

	receiverPublicKey, err := new(crypto.Point).FromBytesS(info.PaymentAddress.Pk)
	if err != nil {
		errStr := fmt.Sprintf("Cannot parse outputCoinV2 from PaymentInfo when parseByte PublicKey, error %v ", err)
		return nil, fmt.Errorf(errStr)
	}
	receiverPublicKeyBytes := receiverPublicKey.ToBytesS()
	targetShardID := common.GetShardIDFromLastByte(receiverPublicKeyBytes[len(receiverPublicKeyBytes)-1])

	c := new(CoinV2).Init()
	// Amount, Randomness, SharedRandom are transparency until we call concealData
	c.SetAmount(new(crypto.Scalar).FromUint64(info.Amount))
	c.SetRandomness(crypto.RandomScalar())
	c.SetSharedRandom(crypto.RandomScalar())        // shared randomness for creating one-time-address
	c.SetSharedConcealRandom(crypto.RandomScalar()) //shared randomness for concealing amount and blinding asset tag
	c.SetInfo(info.Message)
	c.SetCommitment(crypto.PedCom.CommitAtIndex(c.GetAmount(), c.GetRandomness(), crypto.PedersenValueIndex))

	// If this is going to burning address then dont need to create ota
	if wallet.IsPublicKeyBurningAddress(info.PaymentAddress.Pk) {
		publicKey, err := new(crypto.Point).FromBytesS(info.PaymentAddress.Pk)
		if err != nil {
			return nil, fmt.Errorf("something is wrong with info.paymentAddress.pk, burning address should be a valid point")
		}
		c.SetPublicKey(publicKey)
		return c, nil
	}

	// Increase index until have the right shardID
	index := uint32(0)
	publicOTA := info.PaymentAddress.GetOTAPublicKey()
	if publicOTA == nil {
		return nil, fmt.Errorf("public OTA from payment address is nil")
	}
	publicSpend := info.PaymentAddress.GetPublicSpend()
	rK := new(crypto.Point).ScalarMult(publicOTA, c.GetSharedRandom())
	for {
		index += 1

		// Get publicKey
		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := new(crypto.Point).ScalarMultBase(hash)
		publicKey := new(crypto.Point).Add(HrKG, publicSpend)
		c.SetPublicKey(publicKey)

		currentShardID, err := c.GetShardID()
		if err != nil {
			return nil, err
		}
		if currentShardID == targetShardID {
			otaRandomPoint := new(crypto.Point).ScalarMultBase(c.GetSharedRandom())
			concealRandomPoint := new(crypto.Point).ScalarMultBase(c.GetSharedConcealRandom())
			c.SetTxRandomDetail(concealRandomPoint, otaRandomPoint, index)
			break
		}
	}
	return c, nil
}
