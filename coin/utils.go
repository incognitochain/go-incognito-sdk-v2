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

func NewCoinFromPaymentInfo(p *CoinParams) (*CoinV2, error) {
	receiverPublicKey, err := new(crypto.Point).FromBytesS(p.PaymentAddress.Pk)
	if err != nil {
		errStr := fmt.Sprintf("Cannot parse outputCoinV2 from PaymentInfo when parseByte PublicKey, error %v ", err)
		return nil, fmt.Errorf(errStr)
	}
	receiverPublicKeyBytes := receiverPublicKey.ToBytesS()
	targetShardID := common.GetShardIDFromLastByte(receiverPublicKeyBytes[len(receiverPublicKeyBytes)-1])

	c := new(CoinV2).Init()
	// Amount, Randomness, SharedRandom are transparency until we call concealData
	c.SetAmount(new(crypto.Scalar).FromUint64(p.Amount))
	c.SetRandomness(crypto.RandomScalar())
	c.SetSharedRandom(crypto.RandomScalar())        // shared randomness for creating one-time-address
	c.SetSharedConcealRandom(crypto.RandomScalar()) // shared randomness for concealing amount and blinding asset tag
	c.SetInfo(p.Message)
	c.SetCommitment(crypto.PedCom.CommitAtIndex(c.GetAmount(), c.GetRandomness(), crypto.PedersenValueIndex))

	// If this is going to burning address then don't need to create ota
	if wallet.IsPublicKeyBurningAddress(p.PaymentAddress.Pk) {
		publicKey, err := new(crypto.Point).FromBytesS(p.PaymentAddress.Pk)
		if err != nil {
			panic("Something is wrong with info.paymentAddress.pk, burning address should be a valid point")
		}
		c.SetPublicKey(publicKey)
		return c, nil
	}

	// Increase index until have the right shardID
	index := uint32(0)
	publicOTA := p.PaymentAddress.GetOTAPublicKey()
	if publicOTA == nil {
		return nil, fmt.Errorf("public OTA from payment address is nil")
	}
	publicSpend := p.PaymentAddress.GetPublicSpend()
	rK := new(crypto.Point).ScalarMult(publicOTA, c.GetSharedRandom())
	for {
		index++

		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := new(crypto.Point).ScalarMultBase(hash)
		publicKey := new(crypto.Point).Add(HrKG, publicSpend)
		c.SetPublicKey(publicKey)

		senderShardID, receivingShardID, coinPrivacyType, _ := DeriveShardInfoFromCoin(publicKey.ToBytesS())
		if receivingShardID == int(targetShardID) && senderShardID == p.SenderShardID && coinPrivacyType == p.CoinPrivacyType {
			otaRandomPoint := new(crypto.Point).ScalarMultBase(c.GetSharedRandom())
			concealRandomPoint := new(crypto.Point).ScalarMultBase(c.GetSharedConcealRandom())
			c.SetTxRandomDetail(concealRandomPoint, otaRandomPoint, index)
			break
		}
	}
	return c, nil
}

const (
	PrivacyTypeTransfer = iota
	PrivacyTypeMint
)

// DeriveShardInfoFromCoin returns the sender origin & receiver shard of a coin based on the
// PublicKey on that coin (encoded inside its last byte).
// Does not support MaxShardNumber > 8.
func DeriveShardInfoFromCoin(coinPubKey []byte) (int, int, int, error) {
	if common.MaxShardNumber > 8 {
		return -1, -1, -1, fmt.Errorf("cannot derive shardID with MaxShardNumber = %v", common.MaxShardNumber)
	}
	numShards := common.MaxShardNumber
	n := int(coinPubKey[len(coinPubKey)-1]) % 128 // use 7 bits
	receiverShardID := n % numShards
	n /= numShards
	senderShardID := n % numShards
	coinPrivacyType := n / numShards

	if coinPrivacyType > PrivacyTypeMint {
		return -1, -1, -1, fmt.Errorf("coin %x has unsupported PrivacyType %d", coinPubKey, coinPrivacyType)
	}
	return senderShardID, receiverShardID, coinPrivacyType, nil
}

// CoinParams contains the necessary data to create a new coin.
type CoinParams struct {
	key.PaymentInfo
	SenderShardID   int
	CoinPrivacyType int
}

// NewCoinParams returns an empty CoinParams.
func NewCoinParams() *CoinParams { return &CoinParams{} }

// NewTransferCoinParams returns a new CoinParams for the transferring purpose.
// If `senderShardParams` is not given, `senderShard` will default to the shardID of the given payment info.
// Otherwise, the first value of `senderShardParams` will be set as the `senderShard`.
func NewTransferCoinParams(paymentInfo *key.PaymentInfo, senderShardParams ...byte) *CoinParams {
	var senderShard byte
	if len(senderShardParams) > 0 {
		senderShard = senderShardParams[0]
	} else {
		receiverPublicKeyBytes := paymentInfo.PaymentAddress.Pk
		senderShard = common.GetShardIDFromLastByte(receiverPublicKeyBytes[len(receiverPublicKeyBytes)-1])
	}

	return &CoinParams{
		PaymentInfo:     *paymentInfo,
		SenderShardID:   int(senderShard),
		CoinPrivacyType: PrivacyTypeTransfer,
	}
}

// NewMintCoinParams returns a new CoinParams for the minting purpose.
func NewMintCoinParams(paymentInfo *key.PaymentInfo) *CoinParams {
	var senderShard byte
	receiverPublicKeyBytes := paymentInfo.PaymentAddress.Pk
	senderShard = common.GetShardIDFromLastByte(receiverPublicKeyBytes[len(receiverPublicKeyBytes)-1])

	return &CoinParams{
		PaymentInfo:     *paymentInfo,
		SenderShardID:   int(senderShard),
		CoinPrivacyType: PrivacyTypeMint,
	}
}
