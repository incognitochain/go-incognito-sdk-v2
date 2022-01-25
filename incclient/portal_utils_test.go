package incclient

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"testing"
)

func TestIncClient_GetPortalShieldingRequestStatus(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	shieldID := "0104babba81c7be00e8628ba5e0f72f7ebb2d0e15244dabd009316b5e6952319"
	status, err := ic.GetPortalShieldingRequestStatus(shieldID)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))
}

func TestIncClient_GeneratePortalShieldingAddress(t *testing.T) {
	var err error
	ic, err = NewMainNetClientWithCache()
	if err != nil {
		panic(err)
	}

	paymentAddress := "12sdVuLAbKAetr7zaS4nQKHrZ3wxqqSFiyiXDnar4gMj552wNbXVZFTXAQuQ9wUyZuMV6ZZuWwGnKM43162ctwqe3U4rmjxmk4Ng8nFVeGH2e5TjVMACvjvWsrVd2wgmvwYtUgrMvp9eMwU2rJJn"
	tokenIDStr := "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696"
	status, err := ic.GeneratePortalShieldingAddress(paymentAddress, tokenIDStr)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))
}

func TestIncClient_GetPortalUnShieldingRequestStatus(t *testing.T) {
	var err error
	ic, err = NewMainNetClientWithCache()
	if err != nil {
		panic(err)
	}

	unShieldID := "decc21f35ed8f9edc5167e1f7b3622e46f95216d0218fe2991d5cf1e4e491511"
	status, err := ic.GetPortalUnShieldingRequestStatus(unShieldID)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))
}

func TestIncClient_GenerateDepositPubKeyFromPrivateKey(t *testing.T) {
	var err error
	ic, err = NewIncClient("http://51.222.43.133:9334", "", 2, "mainnet")
	if err != nil {
		panic(err)
	}

	privateKeyStr := "11111113mea9j9z4QogdaVFQ2VXGQNK2Y6hLHFZGD42kJ1J8FSvQLogdCHuhQbvxLpGWtcwiJLHQHm4yqSetTnUBWG8wusWHAqnTJGpHdJD"
	tokenIdStr := "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696"
	for index := uint64(0); index < 100; index++ {
		depositKey, err := ic.GenerateDepositPubKeyFromPrivateKey(privateKeyStr, tokenIdStr, index)
		if err != nil {
			panic(err)
		}
		jsb, _ := json.Marshal(depositKey)
		Logger.Printf("Index: %v, DepositKey: %v\n\n", index, string(jsb))
	}
}

func TestIncClient_GetNextOTDepositKey(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	tokenIdStr := "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696"
	for attempt := uint64(0); attempt < 100; attempt++ {
		w, err := wallet.GenRandomWalletForShardID(byte(common.RandInt() % 8))
		if err != nil {
			panic(err)
		}
		privateKeyStr := w.Base58CheckSerialize(0)
		depositKey, depositAddr, err := ic.GetNextOTDepositKey(privateKeyStr, tokenIdStr)
		if err != nil {
			panic(err)
		}
		jsb, _ := json.Marshal(depositKey)
		Logger.Printf("Attempt: %v, DepositAddr: %v, DepositKey: %v\n\n", attempt, depositAddr, string(jsb))
	}
}
