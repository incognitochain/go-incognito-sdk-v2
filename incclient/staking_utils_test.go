package incclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIncClient_GetRewardAmount(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	addr := "12S6wUBvQ2wbRjj3VQWYKdiJLznts4SnMpm42XgzBRqNesUW1PMq8hhFJQRQi889u6pk9XGG6SfaBmSU6TGGHDcsS8w52iXCqPp4eLT"
	listRewards, err := ic.GetRewardAmount(addr)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(listRewards, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_ListReward(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	listRewards, err := ic.ListReward()
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(listRewards, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetMiningInfo(t *testing.T) {
	var err error

	// create a new IncClient instance pointing to a validator node
	ic, err = NewIncClient("http://139.162.55.124:10335", "", 1)
	if err != nil {
		panic(err)
	}

	miningInfo, err := ic.GetMiningInfo()
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(miningInfo, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetSyncStats(t *testing.T) {
	var err error

	// create a new IncClient instance pointing to a validator node
	ic, err = NewIncClient("http://139.162.55.124:10335", "", 1)
	if err != nil {
		panic(err)
	}

	stats, err := ic.GetSyncStats()
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(stats, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}