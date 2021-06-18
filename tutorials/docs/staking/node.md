---
description: Tutorial on how to monitor a validator node.
---
# Node Monitoring
Beside the official [tool](https://monitor.incognito.org/node-monitor), node owners can use this go-sdk to monitor their nodes, collect statistics
about them and more. The go-sdk supports the following functions for node owners.

Name | Description | Host
-------------|-------------|-------------
ListReward| List all current  staking rewards on the blockchain | Full-node
GetReward | Get the rewards of a payment address | Full-node
GetMiningInfo | Get the staking information of a node | Validator node
GetSyncStats | Get the statistics of data-synchronizing status of a node | Validator node

All of these functions are implemented in this [file](../../../incclient/staking_utils.go). Next, we give an example of each function.

## List rewards
It's easy to list all the current staking rewards by calling the function `ListReward`.
```go
client, err := incclient.NewTestNet1Client()
if err != nil {
    log.Fatal(err)
}

listRewards, err := client.ListReward()
if err != nil {
    log.Fatal(err)
}

jsb, err := json.MarshalIndent(listRewards, "", "\t")
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(jsb))
```

Here is an example of the expected result. The result consists of a mapping from a public key to another mapping that maps a tokenID to the corresponding reward amount.
```json
{
	"1j6w4uScn5HSDcnHxeXd2NR9XppXDmyfRu9XsjcYNKrEmzYb29": {
		"0000000000000000000000000000000000000000000000000000000000000004": 2021674877
	},
	"1j9pF6dvBxZSCw3fv1U1zWTbVw2Ynt9hViCx6unS11y3sCzAEJ": {
		"0000000000000000000000000000000000000000000000000000000000000004": 4358274501652,
		"01f7587311227070c61cc736b1746534689dc81045cd1020ecdb0fb6d5ee55ec": 31500,
		"4946b16a08a9d4afbdf416edf52ef15073db0fc4a63e78eb9de80f94f6c0852a": 3,
		"86c45a9fdddc5546e3b4f09dba211b836aefc5d08ed22e7d33cff7f9b8b39c10": 1,
		"880ea0787f6c1555e59e3958a595086b7802fc7a38276bcd80d4525606557fbc": 2,
		"9fca0a0947f4393994145ef50eecd2da2aa15da2483b310c2c0650301c59b17d": 0,
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854": 2
	},
	"1jCr2G2MGvnX4CH4wnYt4Jr2EDhUiESpM6QpzbfUQmsdJg2qjv": {
		"0000000000000000000000000000000000000000000000000000000000000004": 9820367405601,
		"42f4bee6e1c14f94697fb35b0b0bd7e08da1b3ab8a0311563a6793175e31e93b": 70714,
		"4946b16a08a9d4afbdf416edf52ef15073db0fc4a63e78eb9de80f94f6c0852a": 1,
		"5c562893dc38c3c2899143ec32cf67051912fc5b6cf8a8c8c7f8d7397fa64418": 109090,
		"880ea0787f6c1555e59e3958a595086b7802fc7a38276bcd80d4525606557fbc": 0,
		"93573e0e59f687e5c95fb787aa85d94144e688ba2bc6c76fc1fc658b99eea99f": 66666,
		"961179e5a1c6b354e3544cb7e3c74d1cd1625e59d1138fdafca7b0f9c0c9eaad": 24827,
		"9fca0a0947f4393994145ef50eecd2da2aa15da2483b310c2c0650301c59b17d": 0,
		"f6c3b18679aff8d307b08d4724697bb8dca123a536b863cbe55dc59c110f5c10": 81817,
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854": 1
	}
}
```
## Get reward amounts
This time, we use the function [`GetReward`](../../../incclient/staking_utils.go)
```go
client, err := incclient.NewTestNet1Client()
if err != nil {
    panic(err)
}

addr := "12S6wUBvQ2wbRjj3VQWYKdiJLznts4SnMpm42XgzBRqNesUW1PMq8hhFJQRQi889u6pk9XGG6SfaBmSU6TGGHDcsS8w52iXCqPp4eLT"
listRewards, err := client.GetRewardAmount(addr)
if err != nil {
    panic(err)
}

jsb, err := json.MarshalIndent(listRewards, "", "\t")
if err != nil {
    panic(err)
}

fmt.Println(string(jsb))
```
and the result looks like
```json
{
	"0000000000000000000000000000000000000000000000000000000000000004": 9925598355659,
	"002ffd86f6b6d0342ebb641e7b89748ba44075db1765173b7d4e77289fbf28fd": 48000,
	"0fa3e49c7d01a3df067c55293705844ae7d41befd3dfc2f231ab763e9c7daa04": 5,
	"42f4bee6e1c14f94697fb35b0b0bd7e08da1b3ab8a0311563a6793175e31e93b": 70714,
	"4946b16a08a9d4afbdf416edf52ef15073db0fc4a63e78eb9de80f94f6c0852a": 2,
	"5c562893dc38c3c2899143ec32cf67051912fc5b6cf8a8c8c7f8d7397fa64418": 29032,
	"880ea0787f6c1555e59e3958a595086b7802fc7a38276bcd80d4525606557fbc": 4,
	"8ba3466c61cbcdd895be8ccbdcc74e7f56a764d6bf390a9abdc8bfe1322e67d6": 4,
	"961179e5a1c6b354e3544cb7e3c74d1cd1625e59d1138fdafca7b0f9c0c9eaad": 54545,
	"96d4ee94024abb55c0f000978f73dee078682f94b93a9fa67afdcdc11b79e4ef": 189000,
	"a37469618aa6e768e6d511db6414fcfe8668b914651976b9509a01ce9e855f58": 94500,
	"f6c3b18679aff8d307b08d4724697bb8dca123a536b863cbe55dc59c110f5c10": 27272,
	"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854": 9
}
```

## Get mining info
The function `GetMiningInfo` is dedicated to node owners. Unlike the two previous functions, it points to the node address instead of a full-node address.
This function returns the basic information of a staked node, including the ShardID, the least shard height, beacon height, the current role, etc. Let's see an example.
```go
// create a new IncClient instance pointing to a validator node
client, err := incclient.NewIncClient("http://139.162.55.124:10335", "", 1)
if err != nil {
    panic(err)
}

miningInfo, err := client.GetMiningInfo()
if err != nil {
    panic(err)
}

jsb, err := json.MarshalIndent(miningInfo, "", "\t")
if err != nil {
    panic(err)
}

fmt.Println(string(jsb))
```
Notice that this time, we initiate a new IncClient instance that points to a staked node, instead of a full-node. An expected result should contain the following information
```json
{
	"ShardHeight": 279926,
	"BeaconHeight": 289040,
	"CurrentShardBlockTx": 0,
	"PoolSize": 0,
	"Chain": "testnet-1",
	"ShardID": 1,
	"Layer": "shard",
	"Role": "pending",
	"MiningPublickey": "1TgJD9emZHRxga1AMFjvt3dck6g8Rtc3623RhkAgshpfTgK3YbaL6cK8J82qMcUDogEw1ATLADMvnTjDtZTMy9BPLMV5kvNbeNkNtFAPp8nzDGsGH6uGYsUEPntGTW4qnBchQgP3Cb9wra7JZexb9AQWJEwKNgjpN6dypnRmhMJqRGiRHzcts",
	"IsEnableMining": true
}
```
## Get data-synchronizing statistics
Like `GetMiningInfo`, `GetSyncStats` is also dedicated to node owners. An example is
```go
// create a new IncClient instance pointing to a validator node
client, err = incclient.NewIncClient("http://139.162.55.124:10335", "", 1)
if err != nil {
    panic(err)
}

stats, err := client.GetSyncStats()
if err != nil {
    panic(err)
}

jsb, err := json.MarshalIndent(stats, "", "\t")
if err != nil {
    panic(err)
}

fmt.Println(string(jsb))
```

and the expected result:
```json
{
	"Beacon": {
		"IsSync": true,
		"LastInsert": "2021-06-17T10:00:22+0000",
		"BlockHeight": 289086,
		"BlockTime": "2021-06-17T10:00:21+0000",
		"BlockHash": "afd3b639fecc6a3f699603302a6f1713e2abd5686d60b858f8a1a05c3db4278d"
	},
	"Shard": {
		"0": {
			"IsSync": false,
			"LastInsert": "",
			"BlockHeight": 19621,
			"BlockTime": "2021-05-20T15:07:10+0000",
			"BlockHash": "0eacb42d43bcab68ab79dbed0b5d26969dea31b7161b5bbb355069882d1eeaac"
		},
		"1": {
			"IsSync": true,
			"LastInsert": "2021-06-17T10:00:25+0000",
			"BlockHeight": 279972,
			"BlockTime": "2021-06-17T10:00:18+0000",
			"BlockHash": "1be48d78f2a2b9e6f11ccb4208f368f97f1bd0225c23eb0a8ecacfd92bfe24b7"
		}
	}
}
```

## Example
[node_monitor.go](../../code/staking/node/node_monitor.go)
```go
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
```

---
Return to [the table of contents](../../../README.md).