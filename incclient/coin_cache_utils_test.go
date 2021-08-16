package incclient

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"testing"
)

func TestCachedOutCoin_Unmarshal(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	shardID := byte(0)
	tokenID := common.PRVIDStr

	length, err := ic.GetOTACoinLengthByShard(shardID, tokenID)
	if err != nil {
		panic(err)
	}

	r := 1 + common.RandInt() % 100
	idxList := make([]uint64, 0)
	for len(idxList) < r {
		idxList = append(idxList, common.RandUint64() % length)
	}

	res, err := ic.GetOTACoinsByIndices(shardID, tokenID, idxList)
	if err != nil {
		panic(err)
	}

	initCachedOutCoins := NewCachedOutCoins()
	initCachedOutCoins.Data = res

	jsb, err := json.Marshal(initCachedOutCoins)
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))

	var tmpRes map[string]interface{}
	err = json.Unmarshal(jsb, &tmpRes)
	if err != nil {
		panic(err)
	}
	Logger.Println(tmpRes)

	tmpCachedOutCoin := NewCachedOutCoins()
	err = json.Unmarshal(jsb, tmpCachedOutCoin)
	if err != nil {
		panic(err)
	}

	jsb, err = json.Marshal(tmpCachedOutCoin)
	if err != nil {
		panic(err)
	}

	Logger.Println(string(jsb))
}