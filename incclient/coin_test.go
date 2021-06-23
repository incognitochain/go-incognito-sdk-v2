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

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key
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

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key
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

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key
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
