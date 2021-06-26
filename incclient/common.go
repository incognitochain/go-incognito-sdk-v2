//Package incclient provides a client for the Incognito RPC API.
package incclient

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

const DefaultPRVFee = uint64(100)

type TxParam struct {
	senderPrivateKey string
	receiverList     []string
	amountList       []uint64
	fee              uint64
	txTokenParam     *TxTokenParam
	md               metadata.Metadata
	kArgs            map[string]interface{}
}

type TxTokenParam struct {
	tokenID      string
	tokenType    int
	receiverList []string
	amountList   []uint64
	hasTokenFee  bool
	tokenFee     uint64
	kArgs        map[string]interface{}
}

type CustomToken struct {
	tokenID   string
	tokenName string
	amount    uint64
}

func (ct CustomToken) ToString() string {
	return fmt.Sprintf("tokenID: %v, tokenName: %v, amount: %v", ct.tokenID, ct.tokenName, ct.tokenID)
}

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

//PrivateKeyToPaymentAddress returns the payment address for its private key corresponding to the key type.
//KeyType should be -1, 0, 1 where
//	* -1: payment address of version 2
//	* 0: payment address of version 1 with old encoding
//	* 1: payment address of version 1 with new encoding
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

//PrivateKeyToPublicKey returns the public key of a private key.
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

//PrivateKeyToPrivateOTAKey returns the private OTA key of a private key.
func PrivateKeyToPrivateOTAKey(privateKey string) string {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		panic(err)
	}

	if len(keyWallet.KeySet.PrivateKey) == 0 {
		log.Println("no private key found")
		return ""
	}

	return keyWallet.Base58CheckSerialize(wallet.OTAKeyType)
}

//PrivateKeyToReadonlyKey returns the readonly key of a private key.
func PrivateKeyToReadonlyKey(privateKey string) string {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		panic(err)
	}
	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	return keyWallet.Base58CheckSerialize(wallet.ReadonlyKeyType)
}

//PrivateKeyToMiningKey returns the mining key of a private key.
func PrivateKeyToMiningKey(privateKey string) string {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		panic(err)
	}

	if len(keyWallet.KeySet.PrivateKey) == 0 {
		return ""
	}

	miningKey := base58.Base58Check{}.Encode(common.HashB(common.HashB(keyWallet.KeySet.PrivateKey)), common.ZeroByte)

	return miningKey
}

//GetShardIDFromPrivateKey returns the shardID where the private key resides in.
func GetShardIDFromPrivateKey(privateKey string) byte {
	pubKey := PrivateKeyToPublicKey(privateKey)
	return common.GetShardIDFromLastByte(pubKey[len(pubKey)-1])
}
