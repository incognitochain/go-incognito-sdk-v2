package coin

import (
	"encoding/json"
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// OTAReceiver holds the data necessary to receive a coin with privacy.
// It is somewhat equivalent in usage with PaymentAddress.
type OTAReceiver struct {
	// PublicKey is the one-time public key of the receiving coin.
	PublicKey crypto.Point

	// TxRandom is for the receiver to recover the receiving information.
	TxRandom TxRandom

	// SharedSecrets are for the sender to mask the amount as well as the asset type of the sending coins.
	// SharedSecrets = []crypto.Point{sharedOTAPoint, sharedConcealPoint}:
	//	- sharedOTAPoint: used for generating the one-time address and concealing the assetID.
	//	- sharedConcealPoint: used for concealing the amount.
	// For non-privacy transactions, this field can be omitted.
	SharedSecrets []crypto.Point `json:"SharedSecrets,omitempty"`
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
	if len(receiver.SharedSecrets) > 0 {
		if len(receiver.SharedSecrets) != 2 {
			return false
		}
		if !receiver.SharedSecrets[0].PointValid() || !receiver.SharedSecrets[1].PointValid() {
			return false
		}
	}
	return receiver.PublicKey.PointValid()
}

// IsConcealable checks if the OTAReceiver supports full privacy.
func (receiver OTAReceiver) IsConcealable() bool {
	return len(receiver.SharedSecrets) == 2
}

// GetPublicKey returns the base58-encoded PublicKey of an OTAReceiver.
func (receiver OTAReceiver) GetPublicKey() string {
	return base58.Base58Check{}.Encode(receiver.PublicKey.ToBytesS(), 0)
}

// GetTxRandom returns the base58-encoded GetTxRandom of an OTAReceiver.
func (receiver OTAReceiver) GetTxRandom() string {
	return base58.Base58Check{}.Encode(receiver.TxRandom.Bytes(), 0)
}

func (receiver *OTAReceiver) FromAddress(addr key.PaymentAddress, sendingShard ...byte) error {
	if receiver == nil {
		return fmt.Errorf("OTAReceiver not initialized")
	}

	targetShardID := common.GetShardIDFromLastByte(addr.Pk[len(addr.Pk)-1])
	fromShard := targetShardID
	if len(sendingShard) > 0 {
		fromShard = sendingShard[0] % byte(common.MaxShardNumber)
	}

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
		tmpSendingShard, tmpReceivingShard := common.GetShardIDsFromPublicKey(pkb)
		if tmpReceivingShard == targetShardID && tmpSendingShard == fromShard {
			otaRandomPoint := (&crypto.Point{}).ScalarMultBase(otaRand)
			concealRandomPoint := (&crypto.Point{}).ScalarMultBase(concealRand)
			sharedOTAPoint := (&crypto.Point{}).ScalarMult(addr.GetOTAPublicKey(), otaRand)
			sharedConcealPoint := (&crypto.Point{}).ScalarMult(addr.GetPublicView(), concealRand)
			receiver.SharedSecrets = []crypto.Point{*sharedOTAPoint, *sharedConcealPoint}

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

// String marshals the OTAReceiver, then encodes it with base58.
// By default, an OTAReceiver will only support receiving assets in a non-private transaction. Set `isConcealable = true`
// to enable receiving assets in a private transaction.
func (receiver OTAReceiver) String(isConcealable ...bool) string {
	return base58.Base58Check{}.NewEncode(receiver.Bytes(isConcealable...), common.ZeroByte)
}

// Bytes returns a byte-encoded form of an OTAReceiver.
// By default, an OTAReceiver will only support receiving assets in a non-private transaction. Set `isConcealable = true`
// to enable receiving assets in a private transaction.
func (receiver OTAReceiver) Bytes(isConcealable ...bool) []byte {
	concealable := len(isConcealable) > 0 && isConcealable[0]
	rawBytes := []byte{wallet.PrivateReceivingAddressType}
	rawBytes = append(rawBytes, receiver.PublicKey.ToBytesS()...)
	rawBytes = append(rawBytes, receiver.TxRandom.Bytes()...)
	if concealable && len(receiver.SharedSecrets) > 0 {
		for _, s := range receiver.SharedSecrets {
			rawBytes = append(rawBytes, s.ToBytesS()...)
		}
	}
	return rawBytes
}

func (receiver *OTAReceiver) SetBytes(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("not enough bytes to parse ReceivingAddress")
	}
	if receiver == nil {
		return fmt.Errorf("OTAReceiver not initialized")
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
		err = txr.SetBytes(b[33:101])
		if err != nil {
			return err
		}
		receiver.TxRandom = *txr

		if len(b) == 165 {
			buf = make([]byte, 32)
			copy(buf, b[101:133])
			s1, err := (&crypto.Point{}).FromBytesS(buf)
			if err != nil {
				return err
			}

			buf = make([]byte, 32)
			copy(buf, b[133:165])
			s2, err := (&crypto.Point{}).FromBytesS(buf)
			if err != nil {
				return err
			}

			receiver.SharedSecrets = []crypto.Point{*s1, *s2}
		}
		return nil
	default:
		return fmt.Errorf("unrecognized prefix for ReceivingAddress")
	}
}

// MarshalJSON returns a non-private byte-sequence representation of an OTAReceiver.
func (receiver OTAReceiver) MarshalJSON() ([]byte, error) {
	return json.Marshal(receiver.String())
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

// GetShardIDs returns a pair of (sendingShard, receivingShard) of an OTAReceiver.
func (receiver OTAReceiver) GetShardIDs() (byte, byte) {
	pkb := receiver.PublicKey.ToBytesS()
	return common.GetShardIDsFromPublicKey(pkb)
}
