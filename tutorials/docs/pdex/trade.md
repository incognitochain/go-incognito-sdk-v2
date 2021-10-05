---
Description: Tutorial on how to creat trading transactions in pDEX.
---
# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

# pDEX Trades
The Incognito pDEX borrows heavily from Nick Johnson's [reddit post](https://www.reddit.com/r/ethereum/comments/54l32y/euler_the_simplest_exchange_and_currency/) in 2016, Vitalik Buterin's [reddit post](https://www.reddit.com/r/ethereum/comments/55m04x/lets_run_onchain_decentralized_exchanges_the_way/) in 2016, Hayden Adam's [Uniswap implementation](https://github.com/Uniswap/contracts-vyper/blob/master/contracts/uniswap_exchange.vy) in 2018. 

pDEX does not use an order book.  Instead, it implements a novel Automated Market Making algorithm that provides instant matching, no matter how large the order size is or how tiny the liquidity pool is.

The main idea is to replace the traditional order book with a bonding curve mechanism known as constant product. On a typical exchange such as Coinbase or Binance, market makers supply liquidity at various price points. pDEX takes everyone's bids and asks and pool them into two giant buckets. Market makers no longer specify at which prices they are willing to buy or sell. Instead, pDEX automatically makes markets based on a [Automated Market Making algorithm](https://github.com/runtimeverification/verified-smart-contracts/blob/uniswap/uniswap/x-y-k.pdf). 

In this tutorial, we'll see how to create trading transactions using the `go-sdk`. For more information about how the pDEX works, please see [this post](https://raw.githubusercontent.com/incognitochain/incognito-chain/production/specs/pdex.md).

## Create trading transactions with `go-sdk`
It turns out that creating a trading transaction with the `go-sdk` is very simple. All we need is to call the function [`CreateAndSendPDETradeTransaction`](../../../incclient/pdex.go#L139) with the following input parameters:

* `privateKey`: the private key to sign the transaction.
* `tokenIDToSell`: the tokenID we wish to sell
* `tokenIDToBuy`: the tokenID we wish to buy
* `sellAmount`: the selling amount.
* `expectedAmount`: the expected amount we wish to receive.
* `tradingFee`: the trading fee (paid in PRV). The higher the trading fee, the more likely our transaction will be successful.


## Example
[trade.go](../../code/pdex/trade/trade.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	//Trade PRV to tokens
	tokenToSell := common.PRVIDStr
	tokenToBuy := "0000000000000000000000000000000000000000000000000000000000000100"
	sellAmount := uint64(500000000)
	expectedAmount, err := client.CheckXPrice(tokenToSell, tokenToBuy, sellAmount)
	if err != nil {
		log.Fatal(err)
	}
	tradingFee := uint64(10)

	txHash, err := client.CreateAndSendPDETradeTransaction(privateKey, tokenToSell, tokenToBuy, sellAmount, expectedAmount, tradingFee)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
```

---
Return to [the table of contents](../../../README.md).
