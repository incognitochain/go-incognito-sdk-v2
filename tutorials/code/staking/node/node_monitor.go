package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	// list all rewards
	listRewards, err := client.ListReward()
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(listRewards, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsb))

	// get rewards of a user
	addr := "12S6wUBvQ2wbRjj3VQWYKdiJLznts4SnMpm42XgzBRqNesUW1PMq8hhFJQRQi889u6pk9XGG6SfaBmSU6TGGHDcsS8w52iXCqPp4eLT"
	userRewards, err := client.GetRewardAmount(addr)
	if err != nil {
		panic(err)
	}
	jsb, err = json.MarshalIndent(userRewards, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsb))

	// create a new IncClient instance pointing to a validator node
	client, err = incclient.NewIncClient("http://139.162.55.124:10335", "", 1)
	if err != nil {
		panic(err)
	}

	// retrieve the mining info of a node
	miningInfo, err := client.GetMiningInfo()
	if err != nil {
		panic(err)
	}
	jsb, err = json.MarshalIndent(miningInfo, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsb))

	// let's see the sync progress
	stats, err := client.GetSyncStats()
	if err != nil {
		panic(err)
	}
	jsb, err = json.MarshalIndent(stats, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsb))
}
