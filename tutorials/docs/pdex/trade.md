**--- Description: Tutorial on how to create trading transactions in pDEX.
---

# Before Going Further

Please read through the tutorials on [key submission](../accounts/submit_key.md)
and [UTXO cache](../accounts/utxo_cache.md) for proper balance and UTXO retrieval. Skip these parts if you're familiar
with these notions.

# pDEX Trades

Incognito [has recently introduced a new version of its pDEX](https://we.incognito.org/t/introducing-the-new-pdex-pdex-v3/13026)
with the promise of addressing the obstacles of the old (and somewhat inefficient)
instance, and the centralized Provide. Unlike the previous version of the pDEX, in this version, Incognito uses a hybrid
architecture allowing both AMMs and Order Books to empower the best of both worlds. More detail about the design can be
found in this post. In this tutorial, we quickly go through how a trade works in this new design.

## How a Trade is Processed

Suppose we have a pool of ETH/USDT with an AMM pool size of 50 ETH + 100,000 USDT (amplifier = 2) so the amplified pool
size is 100 ETH + 200,000 and the rate is 2000 USDT/ETH.

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
* Then match 0.4 ETH for 799.2 USDT at the rate of 1998 with the 1st order in the Buy orders list (the 1st order is
  filled 100%)
* Then swap 1.06849 ETH for 2,112 USDT with AMM pool, the rate changes from 1998 to 1956.
* Then match 0.98151 ETH left for 1,919.83 USDT at the rate of 1956 with the 2nd order in the Buy orders list (the 2nd
  order is partially filled)

After completing the swaps above, the virtual AMM pool will have 101.11849 ETH + 197,788.05 USDT with a rate of 1956 and
the order book should look like this:

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

So the user will get 4,930.98 USDT from the swap of 2.5 ETH at the rate of 1972.39. If solely swapped with the AMM pool,
the user will only get 4,878 USDT. The hybrid approach of Amplified AMM and Order-Books will help significantly reduce
the slippage.

## Create Trading Transactions

It turns out that creating a trading transaction with the `go-sdk` is very simple. All we need is to call the
function [`CreateAndSendPdexv3TradeTransaction`](../../../incclient/pdex.go) with the following input parameters:

* `privateKey`: the private key to sign the transaction.
* `tradingPath`: a list of pool pairs for this trade.
* `tokenIDToSell`: the tokenID we wish to sell
* `tokenIDToBuy`: the tokenID we wish to buy
* `sellAmount`: the selling amount.
* `expectedAmount`: the expected amount we wish to receive.
* `tradingFee`: the trading fee (paid in PRV). The higher the trading fee, the more likely our transaction will be
  successful.
* `feeInPRV`: whether the trading fee is calculated in PRV.

Here, `tradingPath` is a fairly new term that has just been introduced in pDEX v3. In pDEX v2, for a pair of tokens,
there was at most one pool for it. A trade would never consume more than 2 pools, and thus the pDEX would be able know
which pools to calculate the trade information. For example, if there was a trade from USDT to ETH, `PRV-USDT`
and `PRV-ETH`
pools would be consumed. This will change in pDEX v3. Because there isn't any constraint on the number of pools for each
token pair (e.g, a pair USDT-ETH will have pools USDT-ETH-A, USDT-ETH-B, USDT-ETH-C, etc.), and it doesn't require a
pool to have PRV, the pDEX will not be able to know which pool a trade is targeting. Therefore, the `tradingPath`
parameter is required. This parameter is a list of pools, that a trade consumes. Assuming that we have the following
pools:

* USDT-ETH
    * USDT-ETH-A
    * USDT-ETH-B
    * USDT-ETH-C
* ETH-PRV
    * ETH-PRV-A
    * ETH-PRV-B
* PRV-BTC
    * PRV-BTC-A
* USDT-BTC
    * USDT-BTC-A

To trade from USDT to BTC, a trading path could simply be `[USDT-BTC-A]`, or it could
be `[USDT-ETH-A, ETH-PRV-B, PRV-BTC-A]` depending on which part gives a better receiving amount. Note that ordering of
pools matters because it requires the next pool must contain the buying token of the previous pool. For
example, `[USDT-ETH-A, PRV-BTC-A, ETH-PRV-B]` is not a valid trading path. Furthermore, the more pools in a trading path,
the more trading fee you have to pay. For the above example, the trading path with the only pool `[USDT-BTC-A]` will 
pay less trading fee than the one consisting of `[USDT-ETH-A, ETH-PRV-B, PRV-BTC-A]`.

To check status of a trade, we use the function `CheckTradeStatus`. See more detail in the following example.

## Example

[trade.go](../../code/pdex/trade/trade.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
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
	// Trade between some tokens
	tokenToSell := "00000000000000000000000000000000000000000000000000000000000115d7"
	tokenToBuy := common.PRVIDStr
	sellAmount := uint64(10000)
	expectedAmount := uint64(7000000)
	tradePath := []string{"00000000000000000000000000000000000000000000000000000000000115d7-00000000000000000000000000000000000000000000000000000000000115dc-aeb37b2be73b62b6b5b95086e47687767950e66772e14db6daeef01e40344dd5", "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115dc-03696365b2ff79bb9ef35bf43a74e655ffadae0fa139b8016148d7a036716c5c"}
	tradingFee := uint64(50)
	feeInPRV := false

	txHash, err := client.CreateAndSendPdexv3TradeTransaction(privateKey, tradePath, tokenToSell, tokenToBuy, sellAmount, expectedAmount, tradingFee, feeInPRV)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("txHash: %v\n", txHash)

	time.Sleep(100 * time.Second)
	status, err := client.CheckTradeStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	common.PrintJson(status, "TradeStatus")
}
```

---
Return to [the table of contents](../../../README.md).
