---
Description: Tutorial on how to create an un-staking transaction.
---
# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

# Node UnStaking
To un-stake a node, call the [`CreateAndSendUnStakingTransaction`](../../../incclient/staking.go) function with our private key, mining private key, and the candidate address.

## Example
[unstake.go](../../code/staking/unstake/unstake.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	privateSeed :=  incclient.PrivateKeyToMiningKey(privateKey) //NOTE: the private seed (a.k.a the mining key) can be randomly generated and not be dependent on the private key
	candidateAddress := incclient.PrivateKeyToPaymentAddress(privateKey, -1)

	txHash, err := client.CreateAndSendUnStakingTransaction(privateKey, privateSeed, candidateAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
```
---
Return to [the table of contents](../../../README.md).
