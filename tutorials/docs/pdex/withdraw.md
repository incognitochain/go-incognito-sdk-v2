---
Description: Tutorial on how to withdraw liquidity from the pDEX.
---

# Before Going Further

Please read through the tutorials on [key submission](../accounts/submit_key.md)
and [UTXO cache](../accounts/utxo_cache.md) for proper balance and UTXO retrieval. Skip these parts if you're familiar
with these notions.

# pDEX Withdrawal
In the previous tutorial, we have learned how to create a new pool in the pDEX and contribute liquidity to it. In this tutorial,
we'll walk through the withdrawal process. Liquidity providers can withdraw their contributions at any time they want by submitting a transaction
which consists of the following metadata:
```go
type WithdrawLiquidityRequest struct {
    metadataCommon.MetadataBase
    
    // poolPairID is the ID of the target pool in which the user wants to withdraw his contribution from.
    poolPairID   string
    
    // nftID is the ID of the NFT which he used to contribute with.
    nftID        string
    
    // otaReceivers is a mapping from a tokenID to the corresponding one-time address for receiving back the funds.
    otaReceivers map[string]string
    
    // shareAmount is the amount of share he wants to withdraw from the target pool.
    shareAmount  uint64
}
```

An LP can create this type of transactions using the method `CreateAndSendPdexv3WithdrawLiquidityTransaction`. See the full example below.

## Example
[withdraw.go](../../code/pdex/pdex_withdraw/withdraw.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	firstToken := common.PRVIDStr
	secondToken := "00000000000000000000000000000000000000000000000000000000000115d7"
	pairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	sharedAmount := uint64(5000)

	txHash, err := client.CreateAndSendPdexv3WithdrawLiquidityTransaction(privateKey, pairID, firstToken, secondToken, nftIDStr, sharedAmount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdrawal transaction %v\n", txHash)
}
```

Next, we will see how to [withdraw LP fees](./lp_fee_withdraw.md) from the pDEX.

---
Return to [the table of contents](../../../README.md).
