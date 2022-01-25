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

// IsValid() checks the validity of this OTAReceiver (all referenced Points must be valid).
// Note that some sanity checks are already done when unmarshalling
func (recv OTAReceiver) IsValid() bool {
	_, err := recv.TxRandom.GetTxConcealRandomPoint()
	if err != nil {
		return false
	}
	_, err = recv.TxRandom.GetTxOTARandomPoint()
	if err != nil {
		return false
	}
	return recv.PublicKey.PointValid()
}

func (recv *OTAReceiver) FromAddress(addr key.PaymentAddress) error {
	if recv == nil {
		return errors.New("OTAReceiver not initialized")
	}

	targetShardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])
	otaRand := crypto.RandomScalar()
	concealRand := crypto.RandomScalar()

	// Increase index until have the right shardID
	index := uint32(0)
	publicOTA := addr.GetOTAPublicKey()
	if publicOTA == nil {
		return errors.New("Missing public OTA in payment address")
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
			recv.PublicKey = *publicKey
			recv.TxRandom = *NewTxRandom()
			recv.TxRandom.SetTxOTARandomPoint(otaRandomPoint)
			recv.TxRandom.SetTxConcealRandomPoint(concealRandomPoint)
			recv.TxRandom.SetIndex(index)
			return nil
		}
	}
	return fmt.Errorf("Cannot generate OTAReceiver after %d attempts", MaxTriesOTA)
}

// FromString() returns a new OTAReceiver parsed from the input string,
// or error on failure
func (recv *OTAReceiver) FromString(data string) error {
	raw, _, err := base58.Base58Check{}.Decode(data)
	if err != nil {
		return err
	}
	err = recv.SetBytes(raw)
	if err != nil {
		return err
	}
	return nil
}

// String() marshals the OTAReceiver, then encodes it with base58
func (recv OTAReceiver) String() string {
	rawBytes := recv.Bytes()
	return base58.Base58Check{}.NewEncode(rawBytes, common.ZeroByte)
}

func (recv OTAReceiver) Bytes() []byte {
	rawBytes := []byte{byte(wallet.PrivateReceivingAddressType)}
	rawBytes = append(rawBytes, recv.PublicKey.ToBytesS()...)
	rawBytes = append(rawBytes, recv.TxRandom.Bytes()...)
	return rawBytes
}

func (recv *OTAReceiver) SetBytes(b []byte) error {
	if len(b) == 0 {
		return errors.New("Not enough bytes to parse ReceivingAddress")
	}
	if recv == nil {
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
		recv.PublicKey = *pk
		txr := NewTxRandom()
		// SetBytes() will perform length check
		err = txr.SetBytes(b[33:])
		if err != nil {
			return err
		}
		recv.TxRandom = *txr
		return nil
	default:
		return errors.New("Unrecognized prefix for ReceivingAddress")
	}
}

func (recv OTAReceiver) MarshalJSON() ([]byte, error) {
	s := recv.String()
	return json.Marshal(s)
}

func (recv *OTAReceiver) UnmarshalJSON(raw []byte) error {
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
	*recv = temp
	return nil
}

func (recv OTAReceiver) GetShardID() byte {
	pkb := recv.PublicKey.ToBytes()
	lastByte := pkb[crypto.Ed25519KeySize-1]
	shardID := common.GetShardIDFromLastByte(lastByte)
	return shardID
}
