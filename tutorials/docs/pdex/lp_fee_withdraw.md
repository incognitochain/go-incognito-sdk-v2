---
Description: Tutorial on how to withdraw LP fees from the pDEX.
---

# Before Going Further

Please read through the tutorials on [key submission](../accounts/submit_key.md)
and [UTXO cache](../accounts/utxo_cache.md) for proper balance and UTXO retrieval. Skip these parts if you're familiar
with these notions.

## Withdraw LP Fees
One downside of the previous pDEX is that trading fees were not enforced. Therefore, a liquidity provider has to suffer from very high impermanent losses and basically has no interest in providing liquidity. With the new design, trading fees now
become considerable, and thus there must be a method for an LP to withdraw them.

For this, the LP can create a transaction with the following metadata:
```go
type WithdrawalLPFeeRequest struct {
    metadataCommon.MetadataBase
    
    // PoolPairID is the ID of the target pool pair.
    PoolPairID string                           `json:"PoolPairID"`
    
    // NftID is the ID of the NFT which he used to make contribution.
    NftID      common.Hash                      `json:"NftID"`
    
    // is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
    Receivers  map[common.Hash]coin.OTAReceiver `json:"Receivers"`
}
```
This is made possible via the method `CreateAndSendPdexv3WithdrawLPFeeTransaction`. See the example below.

### Example
[lp_fee_withdraw.go](../../code/pdex/lp_fee_withdraw/lp_fee.go)

```go
package main

import (
   "fmt"
   "log"

   "github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
   client, err := incclient.NewTestNetClient()
   if err != nil {
      log.Fatal(err)
   }

   // replace with your network's data
   privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
   withdrawTokenIDs := make([]string, 0) // leave it empty if you want to withdraw all fees in the pool. Otherwise, specify which token.
   pairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
   nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"

   txHash, err := client.CreateAndSendPdexv3WithdrawLPFeeTransaction(privateKey, pairID, nftIDStr, withdrawTokenIDs...)
   if err != nil {
      log.Fatal(err)
   }

   fmt.Printf("Withdraw Liquidity-Provider Fee submitted in TX %v\n", txHash)
}
```
In this example, we are withdrawing all the fees related to the target pool. If you only want to withdraw the fees related to only one token in the pool, specify it in the `withdrawTokenIDs` params. For example, if you only want to withdraw the PRV fee, set
```go
withdrawTokenIDs = []string{common.PRVStr}
```

So far, we have seen how to manipulate liquidity in the new pDEX. The subsequent tutorials will explain order books in detail, beginning with how to [add an order book](./ob_add.md).

---
Return to [the table of contents](../../../README.md).
