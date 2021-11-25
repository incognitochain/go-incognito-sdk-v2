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
	ic, err := NewTestNetClient()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key
	myNFTs, err := ic.GetAllNFTs(privateKey)
	if err != nil {
		panic(err)
	}
	jsb, err := json.MarshalIndent(myNFTs, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))
}

func TestGetAccountInfoFromPrivateKey(t *testing.T) {
	privateKey := "" // input the private key

	keyInfo, err := GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", keyInfo.String())
}

func TestIncClient_GetAllBalances(t *testing.T) {
	ic, err := NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j" // input the private key
	start := time.Now()
	allBalances, err := ic.GetAllBalances(privateKey, false)
	if err != nil {
		panic(err)
	}
	jsb, _ := json.MarshalIndent(allBalances, "", "\t")
	Logger.Log.Printf("AllBalances: %v\n", string(jsb))
	Logger.Log.Printf("TimeElapsed without v1: %v\n", time.Since(start).Seconds())

	start = time.Now()
	allBalances, err = ic.GetAllBalances(privateKey, true)
	if err != nil {
		panic(err)
	}
	jsb, _ = json.MarshalIndent(allBalances, "", "\t")
	Logger.Log.Printf("AllBalances: %v\n", string(jsb))
	Logger.Log.Printf("TimeElapsed with v1: %v\n", time.Since(start).Seconds())
}
