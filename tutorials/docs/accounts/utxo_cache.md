---
Description: Tutorial on how to retrieve output coins without the need of key submission and cache them.
---

One painful problem with [key submission](./submit_key.md) is that some amount of privacy is sacrificed for faster output coins retrieval. Furthermore,
if a user forgets to submit his key before any v2 output coins arrives, retrieving these output coins becomes much more difficult. These are why
a [new cache layer](../../../incclient/coin_cache.go) is embedded within an [IncClient](../../../incclient/incclient.go). With this cache, the user does not need to submit his [OTAKey](./keys.md) to a full-node to get his
output coins. Instead, he will pull all output coins from the node and check if any output coin belongs to him locally.

To initialize an IncClient with a cache, simply use `NewIncClientWithCache` or `NewMainNetClientWithCache`, etc. See the following example for more detail.

## Example
[utxo_cache.go](../../code/accounts/utxo_cache/utxo_cache.go)
```go
package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	// create a new client
	client, err := incclient.NewTestNet1ClientWithCache()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	tokenID := "0000000000000000000000000000000000000000000000000000000000000004"

	// create a new OutCoinKey
	outCoinKey, err := incclient.NewOutCoinKeyFromPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	outCoinKey.SetReadonlyKey("") // call this if you do not want the remote full-node to decrypt your coin

	outCoinsV1, idxList, err := client.GetOutputCoins(outCoinKey, tokenID, 0)
	for i, outCoin := range outCoinsV1 {
		fmt.Printf("idx: %v, version: %v, isEncrypted: %v, publicKey: %v\n", idxList[i].Uint64(),
			outCoin.GetVersion(), outCoin.IsEncrypted(), base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), 0x00))
	}
}
```
Please note that it will be slow for the first time you retrieve output coins. And the cached output coins will be stored at the cached directory (something like `.cache/mainnet`).

---
Return to [the table of contents](../../../README.md).