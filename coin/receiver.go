package coin

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// OTAReceiver holds the data necessary to send a coin to your receiver with privacy.
// It is somewhat equivalent in usage with PaymentAddress
type OTAReceiver struct {
	PublicKey crypto.Point
	TxRandom  TxRandom
}

// IsValid checks the validity of this OTAReceiver (all referenced Points must be valid).
// Note that some sanity checks are already done when unmarshalling
func (receiver OTAReceiver) IsValid() bool {
	_, err := receiver.TxRandom.GetTxConcealRandomPoint()
	if err != nil {
		return false
	}
	_, err = receiver.TxRandom.GetTxOTARandomPoint()
	if err != nil {
		return false
	}
	return receiver.PublicKey.PointValid()
}

// FromAddress generates an OTAReceiver from the given payment address.
// Note: it does not generate an OTAReceiver matching the description of the new coin-grouping scheme.
// Deprecated: FromCoinParams instead.
func (receiver *OTAReceiver) FromAddress(addr key.PaymentAddress) error {
	if receiver == nil {
		return fmt.Errorf("OTAReceiver not initialized")
	}

	targetShardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])
	otaRand := crypto.RandomScalar()
	concealRand := crypto.RandomScalar()

	// Increase index until have the right shardID
	index := uint32(0)
	publicOTA := addr.GetOTAPublicKey()
	if publicOTA == nil {
		return fmt.Errorf("missing public OTA in payment address")
	}
	publicSpend := addr.GetPublicSpend()
	rK := (&crypto.Point{}).ScalarMult(publicOTA, otaRand)
	for i := MaxTriesOTA; i > 0; i-- {
		index++
		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := (&crypto.Point{}).ScalarMultBase(hash)
		publicKey := (&crypto.Point{}).Add(HrKG, publicSpend)

		pkb := publicKey.ToBytesS()
		currentShardID := common.GetShardIDFromLastByte(pkb[len(pkb)-1])
		if currentShardID == targetShardID {
			otaRandomPoint := (&crypto.Point{}).ScalarMultBase(otaRand)
			concealRandomPoint := (&crypto.Point{}).ScalarMultBase(concealRand)
			receiver.PublicKey = *publicKey
			receiver.TxRandom = *NewTxRandom()
			receiver.TxRandom.SetTxOTARandomPoint(otaRandomPoint)
			receiver.TxRandom.SetTxConcealRandomPoint(concealRandomPoint)
			receiver.TxRandom.SetIndex(index)
			return nil
		}
	}
	return fmt.Errorf("cannot generate OTAReceiver after %d attempts", MaxTriesOTA)
}

// FromCoinParams generates an OTAReceiver from the given CoinParams.
func (receiver *OTAReceiver) FromCoinParams(p *CoinParams) error {
	if receiver == nil {
		return fmt.Errorf("OTAReceiver not initialized")
	}

	addr := p.PaymentInfo.PaymentAddress

	receiverShardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])
	otaRand := crypto.RandomScalar()
	concealRand := crypto.RandomScalar()

	// Increase index until have the right shardID
	index := uint32(0)
	publicOTA := addr.GetOTAPublicKey()
	if publicOTA == nil {
		return fmt.Errorf("missing public OTA in payment address")
	}
	publicSpend := addr.GetPublicSpend()
	rK := (&crypto.Point{}).ScalarMult(publicOTA, otaRand)
	for i := MaxTriesOTA; i > 0; i-- {
		index++
		hash := crypto.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
		HrKG := (&crypto.Point{}).ScalarMultBase(hash)
		publicKey := (&crypto.Point{}).Add(HrKG, publicSpend)

		tmpSenderShardID, tmpReceiverShardID, tmpCoinType, _ := DeriveShardInfoFromCoin(publicKey.ToBytesS())
		if tmpReceiverShardID == int(receiverShardID) && tmpSenderShardID == p.SenderShardID && tmpCoinType == p.CoinPrivacyType {
			otaRandomPoint := (&crypto.Point{}).ScalarMultBase(otaRand)
			concealRandomPoint := (&crypto.Point{}).ScalarMultBase(concealRand)
			receiver.PublicKey = *publicKey
			receiver.TxRandom = *NewTxRandom()
			receiver.TxRandom.SetTxOTARandomPoint(otaRandomPoint)
			receiver.TxRandom.SetTxConcealRandomPoint(concealRandomPoint)
			receiver.TxRandom.SetIndex(index)
			return nil
		}
	}
	return fmt.Errorf("cannot generate OTAReceiver after %d attempts", MaxTriesOTA)
}

// FromString returns a new OTAReceiver parsed from the input string,
// or error on failure
func (receiver *OTAReceiver) FromString(data string) error {
	raw, _, err := base58.Base58Check{}.Decode(data)
	if err != nil {
		return err
	}
	err = receiver.SetBytes(raw)
	if err != nil {
		return err
	}
	return nil
}

// String marshals the OTAReceiver, then encodes it with base58
func (receiver OTAReceiver) String() (string, error) {
	rawBytes, err := receiver.Bytes()
	if err != nil {
		return "", err
	}
	return base58.Base58Check{}.NewEncode(rawBytes, common.ZeroByte), nil
}

func (receiver OTAReceiver) Bytes() ([]byte, error) {
	rawBytes := []byte{byte(wallet.PrivateReceivingAddressType)}
	rawBytes = append(rawBytes, receiver.PublicKey.ToBytesS()...)
	rawBytes = append(rawBytes, receiver.TxRandom.Bytes()...)
	return rawBytes, nil
}

func (receiver *OTAReceiver) SetBytes(b []byte) error {
	if len(b) == 0 {
		return errors.New("Not enough bytes to parse ReceivingAddress")
	}
	if receiver == nil {
		return errors.New("OTAReceiver not initialized")
	}
	keyType := b[0]
	switch keyType {
	case wallet.PrivateReceivingAddressType:
		buf := make([]byte, 32)
		copy(buf, b[1:33])
		pk, err := (&crypto.Point{}).FromBytesS(buf)
		if err != nil {
			return err
		}
		receiver.PublicKey = *pk
		txr := NewTxRandom()
		// SetBytes() will perform length check
		err = txr.SetBytes(b[33:])
		if err != nil {
			return err
		}
		receiver.TxRandom = *txr
		return nil
	default:
		return errors.New("unrecognized prefix for ReceivingAddress")
	}
}

func (receiver OTAReceiver) MarshalJSON() ([]byte, error) {
	s, err := receiver.String()
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

func (receiver *OTAReceiver) UnmarshalJSON(raw []byte) error {
	var encodedString string
	err := json.Unmarshal(raw, &encodedString)
	if err != nil {
		return err
	}
	var temp OTAReceiver
	err = temp.FromString(encodedString)
	if err != nil {
		return err
	}
	*receiver = temp
	return nil
}

func (receiver OTAReceiver) DeriveShardID() (int, int, int, error) {
	return DeriveShardInfoFromCoin(receiver.PublicKey.ToBytesS())
}
