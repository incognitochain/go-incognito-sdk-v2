---
Description: Tutorial on how to un-stake a token from the pDEX.
---
# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

## pDEX UnStake
An un-staking transaction must consist of the following metadata:
```go
type UnstakingRequest struct {
    metadataCommon.MetadataBase
    
    // stakingPoolID is the ID of the target staking pool (or the tokenID) wished to un-stake from.
    stakingPoolID   string
    
    // otaReceivers is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
    otaReceivers    map[string]string
    
    // nftID is theID of the NFT associated with the staking request.
    nftID           string
    
    // unstakingAmount is the amount wished to un-stake.
    unstakingAmount uint64
}
```

Creating an un-staking transaction is done via the method `CreateAndSendPdexv3UnstakingTransaction` while checking the result can be done via the method `CheckDEXUnStakingStatus`.
The status consists of the following information:
```go
// DEXUnStakeStatus represents the status of a pDEX un-staking transaction.
type DEXUnStakeStatus struct {
    // Status represents the status of the transaction, and should be understood as follows:
    // - 0: the request is rejected;
    // - 1: the request is accepted.
    Status int `json:"Status"`
    
    // NftID is the ID of the NFT associated with the action.
    NftID string `json:"NftID"`
    
    // StakingPoolID is the ID of the pool.
    StakingPoolID string `json:"StakingPoolID"`
    
    // Liquidity is the un-staked amount.
    Liquidity uint64 `json:"Liquidity"`
}
```

## Example
[unstake.go](../../code/pdex/unstake/unstake.go)

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
	amount := uint64(23000)

	txHash, err := client.CreateAndSendPdexv3UnstakingTransaction(privateKey, tokenIDStr, nftIDStr, amount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Unstaking TX for pool %s submitted %v\n", tokenIDStr, txHash)
	time.Sleep(100 * time.Second)
	status, err := client.CheckDEXUnStakingStatus(txHash)
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

In the next tutorial, we'll learn how to [withdraw the staking reward](./staking_reward_withdraw.md) from the pDEX.

---
Return to [the table of contents](../../../README.md).
