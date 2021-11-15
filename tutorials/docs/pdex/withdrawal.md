---
Description: Tutorial on how to withdraw pairs from pDEX.
---

# Before Going Further

Please read through the tutorials on [key submission](../accounts/submit_key.md)
and [UTXO cache](../accounts/utxo_cache.md) for proper balance and UTXO retrieval. Skip these parts if you're familiar
with these notions.

# pDEX Withdrawal
Liquidity providers can withdraw their contributions at any time they want. The Client implements a transaction [`CreateAndSendPDEWithdrawalTransaction`](../../../incclient/pdex.go) to facilitate this operation.

## Prepare our inputs
```go
privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

firstToken := common.PRVIDStr
secondToken := "0000000000000000000000000000000000000000000000000000000000000100"
addr := incclient.PrivateKeyToPaymentAddress(privateKey, -1)
sharedAmount, err := client.GetShareAmount(0, firstToken, secondToken, addr) // get our current shared amount
```
We need to specify the two tokenIDs we wish to withdraw, our payment address, and the withdrawal shared amount. In this case, we are withdrawing everything we have in the pDEX.

## Create and send the withdrawal transaction
```go
txHash, err := client.CreateAndSendPDEWithdrawalTransaction(privateKey, firstToken, secondToken, sharedAmount, 2)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Withdrawal transaction %v\n", txHash)
```

## Example
[withdraw.go](../../code/pdex/withdrawal/withdraw.go)

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

	firstToken := common.PRVIDStr
	secondToken := "0000000000000000000000000000000000000000000000000000000000000100"
	addr := incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	sharedAmount, err := client.GetShareAmount(0, firstToken, secondToken, addr) // get our current shared amount

	txHash, err := client.CreateAndSendPDEWithdrawalTransaction(privateKey, firstToken, secondToken, sharedAmount, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Withdrawal transaction %v\n", txHash)
}
```

---
Return to [the table of contents](../../../README.md).
