package incclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIncClient_GetActiveShards(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	shards, err := ic.GetActiveShard()
	if err != nil {
		panic(err)
	}

	fmt.Printf("#shards: %v\n", shards)
}

func TestIncClient_GetBestBlocks(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	bestBlocks, err := ic.GetBestBlock()
	if err != nil {
		panic(err)
	}

	fmt.Printf("bestBlocks: %v\n", bestBlocks)
}

func TestIncClient_GetListToken(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	listTokens, err := ic.GetListToken()
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(listTokens, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetRawMemPool(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	rawMemPool, err := ic.GetRawMemPool()
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(rawMemPool, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetCommitteeStateByShard(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	committeeState, err := ic.GetCommitteeStateByShard(1, "")
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(committeeState, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetShardBestState(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	state, err := ic.GetShardBestState(0)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(state, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetBeaconBestState(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	state, err := ic.GetBeaconBestState(0)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(state, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}
