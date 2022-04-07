package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"testing"
	"time"
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
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key
	tokenID := common.ConfidentialAssetID.String()

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

func TestIncClient_GetOTACoinsByIndices(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	shardID := byte(0)
	tokenID := common.PRVIDStr

	for i := 0; i < 1; i++ {
		lengths, err := ic.GetOTACoinLength()
		if err != nil {
			panic(err)
		}

		length := lengths[tokenID][shardID]
		r := 1 + common.RandInt()%100
		idxList := make([]uint64, 0)
		for len(idxList) < r {
			idxList = append(idxList, common.RandUint64()%length)
		}

		res, err := ic.GetOTACoinsByIndices(shardID, tokenID, idxList)
		if err != nil {
			panic(err)
		}

		Logger.Println(res)
	}
}

func TestIncClient_GetOTACoinLength(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		lengths, err := ic.GetOTACoinLength()
		if err != nil {
			panic(err)
		}

		Logger.Println(lengths)
	}
}

func TestIncClient_GetOTACoinLengthByShard(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		shardID := byte(common.RandInt() % common.MaxShardNumber)
		isPRV := (common.RandInt() % 2) == 1
		tokenID := common.PRVIDStr
		if !isPRV {
			tokenID = common.ConfidentialAssetID.String()
		}

		length, err := ic.GetOTACoinLengthByShard(shardID, tokenID)
		if err != nil {
			panic(err)
		}

		Logger.Println(shardID, tokenID, length)
	}
}

func TestIncClient_GetAllAssetTags(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	start := time.Now()
	_, err = ic.GetAllAssetTags()
	if err != nil {
		panic(err)
	}

	//Logger.Println(len(assetTags), assetTags)
	Logger.Printf("timeElapsed: %v\n", time.Since(start).Seconds())

	start = time.Now()
	_, err = ic.GetAllAssetTags()
	if err != nil {
		panic(err)
	}

	//Logger.Println(len(assetTags), assetTags)
	Logger.Printf("timeElapsed (second call): %v\n", time.Since(start).Seconds())
}

func TestIncClient_GetAllUTXOsV2(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key

	allUTXOs, allIndices, err := ic.GetAllUTXOsV2(privateKey)
	if err != nil {
		panic(err)
	}

	for tokenID, utxoList := range allUTXOs {
		fmt.Printf("\ntokenID %v: %v\n", tokenID, len(utxoList))
		for i, utxo := range utxoList {
			fmt.Printf("ver: %v, idx: %v, pubKey: %v, cmt: %v, value: %v\n", utxo.GetVersion(), allIndices[tokenID][i],
				base58.Base58Check{}.Encode(utxo.GetPublicKey().ToBytesS(), common.ZeroByte),
				base58.Base58Check{}.Encode(utxo.GetCommitment().ToBytesS(), common.ZeroByte),
				utxo.GetValue())
		}

	}
}

func TestIncClient_HasOTAPubKey(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	exists := ic.HasOTAPubKey("12hi3bHaJhda63RS1yfA2H3eXazfDaY1UBGs1wzQzhqVSbq2UY1")
	fmt.Println(exists)
}
