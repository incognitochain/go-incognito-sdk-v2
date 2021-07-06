package incclient

import (
	"fmt"
	"testing"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func TestIncClient_AuthorizedSubmitKey(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKeys := make([]string, 0)

	for i := 0; i < 4; i++ {
		randomWallet, err := wallet.GenRandomWalletForShardID(byte(common.RandInt() % common.MaxShardNumber))
		if err != nil {
			panic(err)
		}
		tmpPrivateKey := randomWallet.Base58CheckSerialize(wallet.PrivateKeyType)

		privateKeys = append(privateKeys, tmpPrivateKey)
	}

	accessToken := ""
	count := 0
	for i, privateKey := range privateKeys {
		otaKey := PrivateKeyToPrivateOTAKey(privateKey)
		err = ic.AuthorizedSubmitKey(otaKey, accessToken, 0, true)
		if err != nil {
			panic(fmt.Errorf("failed at index %v: %v", i, err))
		}

		fmt.Printf("%v: Submited OTAKey %v success!\n", i, otaKey)
		count += 1
	}
}

func TestIncClient_NewRPCCall(t *testing.T) {
	var err error
	ic, err = NewIncClient("http://51.79.76.38:8334", "", 2)
	if err != nil {
		panic(err)
	}

	method := "listunspentoutputcoinsfromcache"
	params := make([]interface{}, 0)
	params = append(params, 0)
	params = append(params, 99999999)
	keyParams := make([]interface{}, 0)
	keyParams = append(keyParams, map[string]interface{}{"PrivateKey": "112t8rnX6USJnBzswUeuuanesuEEUGsxE8Pj3kkxkqvGRedUUPyocmtsqETX2WMBSvfBCwwsmMpxonhfQm2N5wy3SrNk11eYx6pMsmsic4Vz"})
	params = append(params, keyParams)

	resp, err := ic.NewRPCCall("1.0", method, params, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp))
}

func TestIncClient_NewRPCCall2(t *testing.T) {
	var err error
	ic, err = NewIncClient("http://51.79.76.38:8334", "", 2)
	if err != nil {
		panic(err)
	}

	method := "getbalancebyprivatekey"
	params := make([]interface{}, 0)
	params = append(params, "112t8rnX6USJnBzswUeuuanesuEEUGsxE8Pj3kkxkqvGRedUUPyocmtsqETX2WMBSvfBCwwsmMpxonhfQm2N5wy3SrNk11eYx6pMsmsic4Vz")

	resp, err := ic.NewRPCCall("1.0", method, params, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp))
}
