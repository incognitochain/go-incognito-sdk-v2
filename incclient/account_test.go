package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"testing"
	"time"
)

func TestIncClient_GetBalance(t *testing.T) {
	ic, err := NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := "" // input the private key
	tokenID := common.PRVIDStr

	balance, err := ic.GetBalance(privateKey, tokenID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Balance: %v\n", balance)
}

func TestIncClient_GetAllNFTs(t *testing.T) {
	ic, err := NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key
	start := time.Now()
	myNFTs, err := ic.GetMyNFTs(privateKey)
	if err != nil {
		panic(err)
	}
	jsb, err := json.MarshalIndent(myNFTs, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Log.Println(string(jsb))
	Logger.Printf("timeElapsed: %v\n", time.Since(start).Seconds())

	start = time.Now()
	myNFTs, err = ic.GetMyNFTs(privateKey)
	if err != nil {
		panic(err)
	}
	jsb, err = json.MarshalIndent(myNFTs, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Log.Println(string(jsb))
	Logger.Printf("timeElapsed (second call): %v\n", time.Since(start).Seconds())
}

func TestGetAccountInfoFromPrivateKey(t *testing.T) {
	privateKey := "" // input the private key

	keyInfo, err := GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", keyInfo.String())
}

func TestIncClient_GetAllBalancesV2(t *testing.T) {
	ic, err := NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key
	start := time.Now()
	allBalances, err := ic.GetAllBalancesV2(privateKey)
	if err != nil {
		panic(err)
	}
	jsb, _ := json.MarshalIndent(allBalances, "", "\t")
	Logger.Log.Printf("AllBalances: %v\n", string(jsb))
	Logger.Log.Printf("TimeElapsed with cache: %v\n", time.Since(start).Seconds())

	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	start = time.Now()
	allBalances, err = ic.GetAllBalancesV2(privateKey)
	if err != nil {
		panic(err)
	}
	jsb, _ = json.MarshalIndent(allBalances, "", "\t")
	Logger.Log.Printf("AllBalances: %v\n", string(jsb))
	Logger.Log.Printf("TimeElapsed without cache: %v\n", time.Since(start).Seconds())
}
