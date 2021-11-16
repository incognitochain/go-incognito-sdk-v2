---
Description: Tutorial on how to withdraw pDEX staking rewards.
---

# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

## Check pDEX Staking Rewards
To see how much reward we are having, we use the method `GetEstimatedDEXStakingReward`:
```go
res, err := client.GetEstimatedDEXStakingReward(0, tokenIDStr, nftIDStr)
if err != nil {
    log.Fatal(err)
}
```
Here is an example result:
```json
{
  "0000000000000000000000000000000000000000000000000000000000000004": 1199307924
}
```

## Withdraw pDEX Staking Rewards
Similar to other pDEX-related transaction types, withdrawing the pDEX staking rewards requires the following metadata to be enclosed:
```go
type WithdrawalStakingRewardRequest struct {
    metadataCommon.MetadataBase
    
    // StakingPoolID 
    StakingPoolID string                           `json:"StakingPoolID"`
    
    // NftID is theID of the NFT associated with the staking request.
    NftID         common.Hash                      `json:"NftID"`
    
    // Receivers is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
    Receivers     map[common.Hash]coin.OTAReceiver `json:"Receivers"`
}
```

To create a transaction with this metadata, we employ the function `CreateAndSendPdexv3WithdrawStakeRewardTransaction`, and for checking the status,
we use the function `CheckDEXStakingRewardWithdrawalStatus`.

The status has the following form:
```go
// DEXWithdrawStakingRewardStatus represents the status of a pDEX staking reward withdrawal transaction.
type DEXWithdrawStakingRewardStatus struct {
    // Status represents the status of the transaction, and should be understood as follows:
    // - 0: the request is rejected;
    // - 1: the request is accepted.
    Status int `json:"Status"`
    
    // Receivers is the receiving information.
    Receivers map[string]struct {
    Address string `json:"Address"`
    Amount  uint64 `json:"Amount"`
    } `json:"Receivers"`
}
```
See the following example.

## Example
[staking_reward_withdraw.go](../../code/pdex/staking_reward_withdraw/staking_reward_withdraw.go)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
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
	// specify which token(s) in this pool to withdraw, leave it empty if withdrawing all tokens.
	withdrawTokenIDs := make([]string, 0)

	// check the current rewards
	res, err := client.GetEstimatedDEXStakingReward(0, tokenIDStr, nftIDStr)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	// withdraw the rewards
	txHash, err := client.CreateAndSendPdexv3WithdrawStakeRewardTransaction(privateKey, tokenIDStr, nftIDStr, withdrawTokenIDs...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Staking TX submitted %v\n", txHash)

	// check the withdrawing status
	time.Sleep(100 * time.Second)
	status, err := client.CheckDEXStakingRewardWithdrawalStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err = json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status: %v\n", string(jsb))
}
```
In this example, we are withdrawing everything, and hence `withdrawTokenIDs` is left empty. In case you want to withdraw a specific token, say PRV, change
`withdrawTokenIDs` to the following:
```go
    withdrawTokenIDs := []string{common.PRVStr}
```

---
Return to [the table of contents](../../../README.md).
