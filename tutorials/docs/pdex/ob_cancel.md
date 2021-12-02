---
Description: Tutorial on how to withdraw an order from the pDEX.
---

# Before Going Further

Please read through the tutorials on [key submission](../accounts/submit_key.md)
and [UTXO cache](../accounts/utxo_cache.md) for proper balance and UTXO retrieval. Skip these parts if you're familiar
with these notions.


# Withdraw an Order


## Place an Order Book
To with an order book, a user must create a transaction with the following metadata:
```go
type WithdrawOrderRequest struct {
// PoolPairID is the ID of the target pool from which the user wants to withdraw his order.
PoolPairID string                           `json:"PoolPairID"`

// OrderID is the ID of the added order.
OrderID    string                           `json:"OrderID"`

// Amount is the amount in which we want to withdraw (0 for all).
Amount     uint64                           `json:"Amount"`

// Receiver is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
Receiver   map[common.Hash]coin.OTAReceiver `json:"Receiver"`

// NftID is the ID of the NFT associated with the order.
NftID      common.Hash                      `json:"NftID"`

metadataCommon.MetadataBase
}
```
It is required that the `OrderID` must be specified. This value can be retrieved by getting the status of the adding request (see [the previous tutorial](./ob_add.md)). Other parameters are the same as other pDEX actions. We use the method `CreateAndSendPdexv3WithdrawOrderTransaction` to perform an order withdrawal.

```go
privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
// set withdrawn amount to 0 to withdraw all remaining balance
withdrawAmount := uint64(0)
poolPairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
orderID := "4d033bad4ae9ef2104feda1712e2b7b7ef215b25a4e58103e6f5a29bb63fd387"
// specify which token(s) in this pool to withdraw, leave it empty if withdrawing all tokens.
withdrawTokenIDs := make([]string, 0)

txHash, err := client.CreateAndSendPdexv3WithdrawOrderTransaction(privateKey, poolPairID, orderID, nftIDStr, withdrawAmount, withdrawTokenIDs...)
if err != nil {
    log.Fatal(err)
}
```
Note that we are withdrawing all tokens in the order. If you only want to withdraw a single token in the order, specify it in the `withdrawTokenIDs` params. For example, if you only want to withdraw PRV, set
```go
withdrawTokenIDs = []string{common.PRVStr}
```

## Retrieve the Status
After submit a withdrawal request, we can use the method `CheckOrderWithdrawalStatus` to retrieve its status. An example result will look like the following,
```go
{
    "Status": 1,
    "TokenID": "0000000000000000000000000000000000000000000000000000000000000004",
    "Amount": 100000
}
```
where `Status = 1` indicates the order has been successfully added, etc.

## Example
[cancel_order.go](../../code/pdex/ob_withdraw/withdraw.go)

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	// set withdrawn amount to 0 to withdraw all remaining balance
	withdrawAmount := uint64(0)
	poolPairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	orderID := "4d033bad4ae9ef2104feda1712e2b7b7ef215b25a4e58103e6f5a29bb63fd387"
	// specify which token(s) in this pool to withdraw, leave it empty if withdrawing all tokens.
	withdrawTokenIDs := make([]string, 0)

	txHash, err := client.CreateAndSendPdexv3WithdrawOrderTransaction(privateKey, poolPairID, orderID, nftIDStr, withdrawAmount, withdrawTokenIDs...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdraw Order %s...\nSubmitted in TX %v\n", orderID, txHash)

	time.Sleep(100 * time.Second)
	status, err := client.CheckOrderWithdrawalStatus(txHash)
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

Another interesting feature of the new pDEX is `Staking` where users can provide liquidity only with a single token. However, only a limited
number of tokens are available with this feature. Let's see [how it works](./stake.md).

---
Return to [the table of contents](../../../README.md).
