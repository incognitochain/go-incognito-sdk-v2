---
Description: Tutorial on how to add order books with the new pDEX
---

# Before Going Further

Please read through the tutorials on [key submission](../accounts/submit_key.md)
and [UTXO cache](../accounts/utxo_cache.md) for proper balance and UTXO retrieval. Skip these parts if you're familiar
with these notions.

# Order Books
An AMM model, on the one hand, offers a good user experience as traders will always get a price without getting too much into the why's and how's.
However, on the other hand, the main disadvantage of an AMM model is high slippage. Just like in the previous pDEX,
users suffered a bad rate for a big trade (say trade token A for token B) when trading in a small pool since such a trade's result
depends solely on liquidity added to the pool. Similarly, another big trade (say trade token B for token A) had this issue as well.
Although the new AMM model allows [an amplifier for each pool](./contribute.md) to amplify the pool liquidity, the slippage is still considerable.

This is why Incognito has introduced a hybrid approach combining both the Multiplied (or Amplified) AMM and Order-Books to utilize
the best of both worlds. Moreover, from the user experience perspective, users do definitely have a demand of placing limit orders on desired prices.

## How the Hybrid Approach Works
So it seems like a hybrid approach of order-book and AMM will be a rescue for both issues above because:
- slippage reduction by aggregating both the investing capital (i.e, the amplified liquidity) of LPs and the trading capital of traders, this also helps reduce the dependency of AMM liquidity where trades will still be possible even in the case lacking liquidity provided by LPs;
- a familiar trading experience (e.g. placing limit orders) just like what a traditional exchange did.

Looking into the following example to see how a trade works in the hybrid approach:

Suppose we have a pool of ETH/USDT with an AMM pool size of 50 ETH + 100,000 USDT (amplifier = 2) so the amplified pool size is 100 ETH + 200,000 and the rate is 2000 USDT/ETH.

An order-book with existing limit orders placed by users as follows:

**Sell orders**
- rate (quantity in ETH)
- 2011 (2)
- 2019 (1.5)
- 2025 (4.2)
- 2034 (3.7)

**Buy orders**
- rate (quantity in ETH)
- 1998 (0.4)
- 1956 (3.8)
- 1912 (2.2)

When a user makes a swap of 2.5 ETH for USDT, it will be executed by:

* Swap 0.05 ETH for 99.95 USDT with AMM pool, the rate changes from 2000 to 1998.
* Then match 0.4 ETH for 799.2 USDT at the rate of 1998 with the 1st order in the Buy orders list (the 1st order is filled 100%)
* Then swap 1.06849 ETH for 2,112 USDT with AMM pool, the rate changes from 1998 to 1956.
* Then match 0.98151 ETH left for 1,919.83 USDT at the rate of 1956 with the 2nd order in the Buy orders list (the 2nd order is partially filled)

After completing the swaps above, the virtual AMM pool will have 101.11849 ETH + 197,788.05 USDT with a rate of 1956 and the order book should look like this:

**Sell orders**
- rate (quantity in ETH)
- 2011 (2)
- 2019 (1.5)
- 2025 (4.2)
- 2034 (3.7)

**Buy orders**
- rate (quantity in ETH)
- 1956 (2.81849)
- 1912 (2.2)

So the user will get 4,930.98 USDT from the swap of 2.5 ETH at the rate of 1972.39. If solely swapped with the AMM pool, the user will only get 4,878 USDT. The hybrid approach of Amplified AMM and Order-Books will help significantly reduce the slippage.

## Place an Order Book
In this section, we'll learn how to add new order to an existing pool. To add an order book, a user must create a transaction with the following metadata:
```go
type AddOrderRequest struct {
    // TokenToSell is the ID of the selling token.
    TokenToSell         common.Hash                      `json:"TokenToSell"`
    
    // PoolPairID is the ID of the pool pair where the order belongs to. In Incognito, an order book is subject to a specific pool.
    PoolPairID          string                           `json:"PoolPairID"`
    
    // SellAmount is the amount of the `TokenToSell` the user wished to sell.
    SellAmount          uint64                           `json:"SellAmount"`
    
    // MinAcceptableAmount is the minimum amount of the buying token the user wished to receive.
    MinAcceptableAmount uint64                           `json:"MinAcceptableAmount"`
    
    // Receiver is a mapping from a tokenID to the corresponding one-time address for receiving back the funds (different OTAs for different tokens).
    Receiver            map[common.Hash]coin.OTAReceiver `json:"Receiver"`
    
    // is the ID of the NFT associated with the order.
    NftID               common.Hash                      `json:"NftID"`
    
    metadataCommon.MetadataBase
}
```

It's fortunate that users do not need to create this metadata themselves. Instead, they just need to use the provided method
`CreateAndSendPdexv3AddOrderTransaction` with the following (sample) parameters:
```go
privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
poolPairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
tokenToSell := common.PRVIDStr
tokenToBuy := "00000000000000000000000000000000000000000000000000000000000115d7"
nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
sellAmount := uint64(100000)
minAcceptableAmount := uint64(100000)

txHash, err := client.CreateAndSendPdexv3AddOrderTransaction(privateKey,
    poolPairID, tokenToSell, tokenToBuy, nftIDStr,
    sellAmount, minAcceptableAmount,
)
```
Please note that the `poolPairID` must consist of both the `tokenToSell` and `tokenToBuy`.

## Retrieve the Status
After placing an order, we wish to know its status to see if it's accepted. For this, we use the method `CheckOrderAddingStatus`. An example result will look like the following,
```go
{
    "Status": 1,
    "OrderID": "4d033bad4ae9ef2104feda1712e2b7b7ef215b25a4e58103e6f5a29bb63fd387"
}
```
where `Status = 1` indicates the order has been successfully added, and `OrderID` is the unique ID of the order.

## Example
[add_order](../../code/pdex/ob_add/add_order.go)

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
   poolPairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
   tokenToSell := common.PRVIDStr
   tokenToBuy := "00000000000000000000000000000000000000000000000000000000000115d7"
   nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
   sellAmount := uint64(100000)
   minAcceptableAmount := uint64(100000)

   txHash, err := client.CreateAndSendPdexv3AddOrderTransaction(privateKey,
      poolPairID, tokenToSell, tokenToBuy, nftIDStr,
      sellAmount, minAcceptableAmount,
   )
   if err != nil {
      log.Fatal(err)
   }
   fmt.Printf("txHash: %v\n", txHash)

   time.Sleep(100 * time.Second)
   status, err := client.CheckOrderAddingStatus(txHash)
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

Adding an order is done. What if we want to cancel it?. Let's move on to the next [tutorial](./ob_cancel.md).

---
Return to [the table of contents](../../../README.md).
