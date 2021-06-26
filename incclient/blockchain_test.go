package incclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

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
