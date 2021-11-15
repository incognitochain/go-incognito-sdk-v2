---
Description: Tutorial on how to create a staking transaction.
---
# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

# Node Staking
A validator is someone who is responsible for verifying transactions on the Incognito network. Once transactions are verified, they are added to the database of each validator. The present of validators helps secure the network by participating in the consensus decisions. Incognito proposes and implements a variant of [pBFT](http://pmg.csail.mit.edu/papers/osdi99.pdf) at the consensus layer. We further improve its efficiency by employing the BLS-based aggregate multi-signature scheme [AMSP](https://eprint.iacr.org/2018/483.pdf). 

Anyone and everyone can become a validator for the network’s consensus, as long as they have 1,750 PRV. To create a staking transaction, we use the function [`CreateAndSendShardStakingTransaction`](../../../incclient/staking.go) with the following input parameters:
* `privateKey`: the private key to sign the transaction.
* `miningKey`: the mining private key of the validator as described [here](../accounts/keys.md).
* `candidateAddress`: the payment address of the candidate.
* `rewardAddress`: the payment address of the reward receivers. The `rewardAddress` is not necessarily the same as the `candidateAddress` (e.g., funded staking).
* `autoReStake`: the indicator whether the node will re-stake or un-stake after finish its role as a validator.

## Example
[stake.go](../../code/staking/stake/stake.go)

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
	rewardAddress := candidateAddress //NOTE: the reward receiver can either be the same as the candidate address or be different
	autoReStake := true

	txHash, err := client.CreateAndSendShardStakingTransaction(privateKey, privateSeed, candidateAddress, rewardAddress, autoReStake)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}
```
---
Return to [the table of contents](../../../README.md).
