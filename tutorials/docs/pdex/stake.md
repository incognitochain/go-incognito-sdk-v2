---
Description: Tutorial on how to stake a single token.
---
# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

## pDEX Staking Pools
Apart from contributing to AMM pools, there is another way to earn trading fees from the pDEX, which is staking pools. In this way, a user only needs to contribute a single token (instead of two as in the case of AMM), and earns trading fees as rewards.
The more a user stakes, the more portion of trading fees he will earn. We could consider a staking pool as Provide + Decentralization, except that staked amounts will not be used to provide liquidity.

However, the number of tokens allowed for staking, as well as the number of reward tokens are limited. In other words, a user is only allowed
to stake a token in the staking token list, and when a trade pays the trading fee in a token in the list of reward tokens, the user receives a portion of the trading fee.

### Available Staking Tokens
We can see the list of available staking tokens by calling the function `GetListStakingPoolShares` on input the beacon height.
If the beacon height is set to 0, it will retrieve the latest information.
```go
res, err := ic.GetListStakingPoolShares(0)
if err != nil {
    panic(err)
}
```
And the result will look like this:
```json
{
    "0000000000000000000000000000000000000000000000000000000000000004": 200,
    "0000000000000000000000000000000000000000000000000000000000000006": 100
}
```
The numbers indicate how much the trading fees are distributed to the pools.

### Available Reward Tokens
To list all the reward tokens, we use the method `GetListStakingRewardTokens` supplied with a beacon height.
If the beacon height is set to 0, it will return the latest information. For example,
```go
res, err := ic.GetListStakingRewardTokens(0)
if err != nil {
    panic(err)
}
```
and the result will look like this:
```json
[
    "0000000000000000000000000000000000000000000000000000000000000004",
    "a7e1e12fab9fdee4d96ee5c930f75c608ef3e96cd7c0468f2033533b5cb12a8f",
    "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
]
```

What we can understand from these examples is that a user is only allowed to stake tokens with ID `0000000000000000000000000000000000000000000000000000000000000004` or `0000000000000000000000000000000000000000000000000000000000000006`;
and whenever there is a trade that pays the trading fee in one of the token `(0000000000000000000000000000000000000000000000000000000000000004, a7e1e12fab9fdee4d96ee5c930f75c608ef3e96cd7c0468f2033533b5cb12a8f, ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854)`, he will receive
an amount of the trading fee. This amount is calculated based on his staking amount, and the percentage of the pool. For example, if a user stakes PRV and his staking accounts for 20% of the staking pool, and the trading fee is 100 PRV then:
- 2 PRV (2%) is distributed to the PRV staking pool; and
- 0.4 PRV (20%) is distributed to the user.

## Stake
Just like other types of transactions, in this case, we need to create a transaction with the following metadata:
```go
type StakingRequest struct {
    metadataCommon.MetadataBase
   
   // tokenID is the token we wish to stake. This token must be in the list of allowed staking tokens.
    tokenID     string
   
   // otaReceiver is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
    otaReceiver string
   
   // nftID is the ID of the NFT associated with the staking request.
    nftID       string
   
   // tokenAmount is the staking amount.
    tokenAmount uint64
}
```

Creating a staking transaction is done via the method `CreateAndSendPdexv3StakingTransaction` while checking the result can be done via the method `CheckDEXStakingStatus`.
The status consists of the following information:
```go
// DEXStakeStatus represents the status of a pDEX staking transaction.
type DEXStakeStatus struct {
    // Status represents the status of the transaction, and should be understood as follows:
    //	- 0: the request is rejected;
    //	- 1: the request is accepted.
    Status int `json:"Status"`
    
    // NftID is the ID of the NFT associated with the action.
    NftID string `json:"NftID"`
    
    // StakingPoolID is the ID of the pool.
    StakingPoolID string `json:"StakingPoolID"`
    
    // Liquidity is the staked amount.
    Liquidity uint64 `json:"Liquidity"`
}
```

## Example
[stake.go](../../code/pdex/stake/stake.go)

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	tokenIDStr := common.PRVIDStr
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	amount := uint64(4300000)

	txHash, err := client.CreateAndSendPdexv3StakingTransaction(privateKey, tokenIDStr, nftIDStr, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Staking TX submitted %v\n", txHash)

	time.Sleep(100 * time.Second)
	status, err := client.CheckDEXStakingStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status: %v\n", string(jsb))
}
```

Next, let's see how we can [un-stake](./unstake.md).

---
Return to [the table of contents](../../../README.md).
