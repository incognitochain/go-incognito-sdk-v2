package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"testing"
)

func TestIncClient_GetOutputCoins(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "" // input the private key
	tokenID := common.PRVIDStr

	outCoinKey, err := NewOutCoinKeyFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	outCoinList, idxList, err := ic.GetOutputCoins(outCoinKey, tokenID, 0)
	if err != nil {
		panic(err)
	}

	for i, outCoin := range outCoinList {
		fmt.Printf("ver: %v, idx: %v, pubKey: %v, cmt: %v, value %v\n", outCoin.GetVersion(), idxList[i],
			base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), common.ZeroByte),
			base58.Base58Check{}.Encode(outCoin.GetCommitment().ToBytesS(), common.ZeroByte),
			outCoin.GetValue())
	}
}

func TestIncClient_GetListDecryptedOutCoin(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "" // input the private key
	tokenID := common.PRVIDStr

	utxoList, err := ic.GetListDecryptedOutCoin(privateKey, tokenID, 0)
	if err != nil {
		panic(err)
	}

	for serialNumber, utxo := range utxoList {
		fmt.Printf("ver: %v, sn: %v, pubKey: %v, cmt: %v, value: %v\n", utxo.GetVersion(), serialNumber,
			base58.Base58Check{}.Encode(utxo.GetPublicKey().ToBytesS(), common.ZeroByte),
			base58.Base58Check{}.Encode(utxo.GetCommitment().ToBytesS(), common.ZeroByte),
			utxo.GetValue())
	}
}

func TestIncClient_GetUnspentOutputCoins(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "" // input the private key
	tokenID := common.PRVIDStr

	utxoList, idxList, err := ic.GetUnspentOutputCoins(privateKey, tokenID, 0)
	if err != nil {
		panic(err)
	}

	for i, utxo := range utxoList {
		fmt.Printf("ver: %v, idx: %v, pubKey: %v, cmt: %v, value: %v\n", utxo.GetVersion(), idxList[i],
			base58.Base58Check{}.Encode(utxo.GetPublicKey().ToBytesS(), common.ZeroByte),
			base58.Base58Check{}.Encode(utxo.GetCommitment().ToBytesS(), common.ZeroByte),
			utxo.GetValue())
	}
}

func TestIncClient_GetSpentOutputCoins(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "" // input the private key
	tokenID := common.PRVIDStr

	spentCoins, idxList, err := ic.GetSpentOutputCoins(privateKey, tokenID, 0)
	if err != nil {
		panic(err)
	}

	for i, spentCoin := range spentCoins {
		fmt.Printf("ver: %v, idx: %v, pubKey: %v, cmt: %v, value: %v\n", spentCoin.GetVersion(), idxList[i],
			base58.Base58Check{}.Encode(spentCoin.GetPublicKey().ToBytesS(), common.ZeroByte),
			base58.Base58Check{}.Encode(spentCoin.GetCommitment().ToBytesS(), common.ZeroByte),
			spentCoin.GetValue())
	}
}

func TestIncClient_BuildAssetTags(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	assetTags, err := ic.GetAllAssetTags()
	if err != nil {
		panic(err)
	}

	Logger.Println(len(assetTags), assetTags)
}
