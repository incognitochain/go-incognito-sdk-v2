package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"testing"
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

func TestIncClient_GetBalanceAll(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := "" // input the private key

	Logger.IsEnable = true
	balances, err := ic.GetBalanceAll(privateKey)
	if err != nil {
		panic(err)
	}

	jsb, _ := json.MarshalIndent(balances, "", "\t")
	fmt.Printf(string(jsb))
}

func TestGetAccountInfoFromPrivateKey(t *testing.T) {
	privateKey := "" // input the private key

	keyInfo, err := GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", keyInfo.String())
}
