//Package incclient provides a client for the Incognito RPC API.
package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// TxParam describes the parameters needed to create a transaction in general.
//
// For creating a token transaction, txTokenParam must not be nil. Otherwise, it should be nil.
type TxParam struct {
	senderPrivateKey string
	receiverList     []string
	amountList       []uint64
	fee              uint64
	txTokenParam     *TxTokenParam
	md               metadata.Metadata

	// additional parameters for special functions
	//	- "PRVInputCoins": a coinParams consisting of PRV input coins and indices used to create a transaction with given
	//input coins.
	//	- "TokenInputCoins": a coinParams consisting of token input coins and indices used to create a transaction with given
	//input coins..
	kArgs map[string]interface{}
}

// TxTokenParam describes the parameters needed for creating a token transaction.
type TxTokenParam struct {
	tokenID      string
	tokenType    int
	receiverList []string
	amountList   []uint64
	hasTokenFee  bool
	tokenFee     uint64
	kArgs        map[string]interface{}
}

// CustomToken represents information of a token.
type CustomToken struct {
	tokenID   string
	tokenName string
	amount    uint64
}

// ToString returns the string-representation of a CustomToken.
func (ct CustomToken) ToString() string {
	return fmt.Sprintf("tokenID: %v, tokenName: %v, amount: %v", ct.tokenID, ct.tokenName, ct.tokenID)
}

// NewTxParam creates a new TxParam.
func NewTxParam(privateKey string, receiverList []string, amountList []uint64, prvFee uint64,
	tokenParam *TxTokenParam, md metadata.Metadata, kArgs map[string]interface{}) *TxParam {
	return &TxParam{
		senderPrivateKey: privateKey,
		receiverList:     receiverList,
		amountList:       amountList,
		fee:              prvFee,
		txTokenParam:     tokenParam,
		md:               md,
		kArgs:            kArgs,
	}
}

// NewTxTokenParam creates a new TxTokenParam.
func NewTxTokenParam(tokenID string, tokenType int, receiverList []string, amountList []uint64, hasTokenFee bool, tokenFee uint64,
	kArgs map[string]interface{}) *TxTokenParam {
	return &TxTokenParam{
		tokenID:      tokenID,
		tokenType:    tokenType,
		receiverList: receiverList,
		amountList:   amountList,
		hasTokenFee:  hasTokenFee,
		tokenFee:     tokenFee,
		kArgs:        kArgs,
	}
}

// PrivateKeyToPaymentAddress returns the payment address for its private key corresponding to the key type.
// KeyType should be -1, 0, 1 where
//	- -1: payment address of version 2
//	- 0: payment address of version 1 with old encoding
//	- 1: payment address of version 1 with new encoding
func PrivateKeyToPaymentAddress(privateKey string, keyType int) string {
	keyWallet, _ := wallet.Base58CheckDeserialize(privateKey)
	err := keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	if err != nil {
		return ""
	}
	paymentAddStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	switch keyType {
	case 0: //Old address, old encoding
		addr, _ := wallet.GetPaymentAddressV1(paymentAddStr, false)
		return addr
	case 1:
		addr, _ := wallet.GetPaymentAddressV1(paymentAddStr, true)
		return addr
	default:
		return paymentAddStr
	}
}

// PrivateKeyToPublicKey returns the public key of a private key.
//
// If the private key is invalid, it returns nil.
func PrivateKeyToPublicKey(privateKey string) []byte {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil
	}

	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	if err != nil {
		return nil
	}
	return keyWallet.KeySet.PaymentAddress.Pk
}

// PrivateKeyToPrivateOTAKey returns the private OTA key of a private key.
//
// If the private key is invalid, it returns an empty string.
func PrivateKeyToPrivateOTAKey(privateKey string) string {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		Logger.Println(err)
		return ""
	}

	if len(keyWallet.KeySet.PrivateKey) == 0 {
		Logger.Println("no private key found")
		return ""
	}

	return keyWallet.Base58CheckSerialize(wallet.OTAKeyType)
}

// PrivateKeyToReadonlyKey returns the readonly key of a private key.
//
// If the private key is invalid, it returns an empty string.
func PrivateKeyToReadonlyKey(privateKey string) string {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		Logger.Println(err)
		return ""
	}

	if len(keyWallet.KeySet.PrivateKey) == 0 {
		Logger.Println("no private key found")
		return ""
	}

	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	return keyWallet.Base58CheckSerialize(wallet.ReadonlyKeyType)
}

// PrivateKeyToMiningKey returns the mining key of a private key.
func PrivateKeyToMiningKey(privateKey string) string {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		Logger.Println(err)
		return ""
	}

	if len(keyWallet.KeySet.PrivateKey) == 0 {
		return ""
	}

	miningKey := base58.Base58Check{}.Encode(common.HashB(common.HashB(keyWallet.KeySet.PrivateKey)), common.ZeroByte)

	return miningKey
}

// GetShardIDFromPrivateKey returns the shardID where the private key resides in.
//
// If the private key is invalid, it returns 255.
func GetShardIDFromPrivateKey(privateKey string) byte {
	pubKey := PrivateKeyToPublicKey(privateKey)
	if pubKey == nil {
		return 0
	}
	return common.GetShardIDFromLastByte(pubKey[len(pubKey)-1])
}

// GetShardIDFromPaymentAddress returns the shardID where the payment address resides in.
//
// If the private key is invalid, it returns 255.
func GetShardIDFromPaymentAddress(addrStr string) (byte, error) {
	keyWallet, err := wallet.Base58CheckDeserialize(addrStr)
	if err != nil {
		return 255, err
	}

	pubKey := keyWallet.KeySet.PaymentAddress.Pk
	if pubKey == nil {
		return 255, fmt.Errorf("publicKey is nil")
	}
	return common.GetShardIDFromLastByte(pubKey[len(pubKey)-1]), nil
}

// AssertPaymentAddressAndTxVersion checks if a string payment address is supported by the underlying transaction.
func AssertPaymentAddressAndTxVersion(paymentAddress interface{}, version int8) (key.PaymentAddress, error) {
	var addr key.PaymentAddress
	var ok bool
	//try to parse the payment address
	if addr, ok = paymentAddress.(key.PaymentAddress); !ok {
		//try the pointer
		if tmpAddr, ok := paymentAddress.(*key.PaymentAddress); !ok {
			//try the string one
			addrStr, ok := paymentAddress.(string)
			if !ok {
				return key.PaymentAddress{}, fmt.Errorf("cannot parse payment address - %v: Not a payment address or string address (txversion %v)", paymentAddress, version)
			}
			keyWallet, err := wallet.Base58CheckDeserialize(addrStr)
			if err != nil {
				return key.PaymentAddress{}, err
			}
			addr = keyWallet.KeySet.PaymentAddress
		} else {
			addr = *tmpAddr
		}
	}

	//Always check public spend and public view keys
	if addr.GetPublicSpend() == nil || addr.GetPublicView() == nil {
		return key.PaymentAddress{}, fmt.Errorf("PublicSpend or PublicView not found")
	}

	//If tx is in version 1, PublicOTAKey must be nil
	if version == 1 {
		if addr.GetOTAPublicKey() != nil {
			return key.PaymentAddress{}, fmt.Errorf("PublicOTAKey must be nil")
		}
	}

	//If tx is in version 2, PublicOTAKey must not be nil
	if version == 2 {
		if addr.GetOTAPublicKey() == nil {
			return key.PaymentAddress{}, fmt.Errorf("PublicOTAKey not found")
		}
	}

	return addr, nil
}
