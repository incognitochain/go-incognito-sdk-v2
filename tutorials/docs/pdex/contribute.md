---
Description: Tutorial on how to add pairs to the pDEX.
---
# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

# pDEX Contribution
Liquidity providers play an essential role in pDEX. They provide liquidity to various pools on pDEX and earn trading fees. The current pDEX consists of several pairs of tokens that help accelerate trading activities. The more liquidity in the pDEX, the better experience the trading process gets. For a pair with high liquidity, the slippage rate will be small. On the other hand, trading with low-liquidity pair will result in a high slippage rate.

In this tutorial, we will see how we can provide liquidity for a pair in the pDEX. Please see this [post](https://github.com/incognitochain/incognito-chain/blob/production/specs/pdex.md) to understand how the pDEX works.

## Prepare our inputs
As usual, we need to specify the private key.
```go
privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
```
Next, specify the two contributed tokenIDs, the amount for each token.
```go
pairID := "newPairID"
firstToken := common.PRVIDStr
secondToken := "0000000000000000000000000000000000000000000000000000000000000100"
firstAmount := uint64(1000000000)
secondAmount := uint64(1000000000)

expectedSecondAmount, err := client.CheckPrice(firstToken, secondToken, firstAmount)
if err == nil {
    secondAmount = expectedSecondAmount
} else {
    log.Println("pool has not been initialized")
}
```

Direct contribution to a pair in the pDEX requires two transaction for the corresponding two tokens. For the pDEX to match these transactions, a `pairID` is also needed. This `pairID` can be anything as long as the two transactions have the same `pairID`. Two transactions with the same `pairID` will be grouped and added to the pDEX.

If the pDEX does have a pool for the two tokenIDs, a new pool pair will be created, and the rate is calculated based on the provided amount (in this case is 1:1). For an existing pair, the pDEX will calculate the contributing amounts based on the current rate, and the remaining amount will be returned to the liquidity provider.

## Create contributing transactions
Now, we have two create two separate transactions to *add* the above tokens to the pDEX. This time, we use the [`CreateAndSendPDEContributeTransaction`](../../../incclient/pdex.go) function.
```go
firstTx, err := client.CreateAndSendPDEContributeTransaction(privateKey, pairID, firstToken, firstAmount, 2)
if err != nil {
	log.Fatal(err)
}

secondTx, err := client.CreateAndSendPDEContributeTransaction(privateKey, pairID, secondToken, secondAmount, 2)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("%v - %v\n", firstTx, secondTx)
```

It is required that these two transactions must all be successful. Any failed transaction will result in the other fund being locked in the pDEX.

## Example
[contribute.go](../../code/pdex/contribution/contribute.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	addr := incclient.PrivateKeyToPaymentAddress(privateKey, -1)

	pairID := "newPairID"
	firstToken := common.PRVIDStr
	secondToken := "0000000000000000000000000000000000000000000000000000000000000100"
	firstAmount := uint64(1000000000)
	secondAmount := uint64(1000000000)

	expectedSecondAmount, err := client.CheckPrice(firstToken, secondToken, firstAmount)
	if err == nil {
		secondAmount = expectedSecondAmount
	} else {
		log.Println("pool has not been initialized")
	}

	firstTx, err := client.CreateAndSendPDEContributeTransaction(privateKey, pairID, firstToken, firstAmount, 2)
	if err != nil {
		log.Fatal(err)
	}

	secondTx, err := client.CreateAndSendPDEContributeTransaction(privateKey, pairID, secondToken, secondAmount, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v - %v\n", firstTx, secondTx)

	fmt.Println("checking share updated...")
	time.Sleep(60 * time.Second)

	// retrieve contributed share
	myShare, err := client.GetShareAmount(0, firstToken, secondToken, addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("myShare: %v\n", myShare)
}
```

We have seen how to add a pair to the pDEX. The next section is instructions on [withdrawing our pair](withdrawal.md).

---
Return to [the table of contents](../../../README.md).
