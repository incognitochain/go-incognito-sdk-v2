---
Description: Tutorial on how to retrieve output coins V1.
---

# Retrieve Output Coins
In this tutorial, we only examine the way to get our output coins V1. The next [post](../accounts/submit_key.md) will be dedicated for output coins V2.
As we have seen in the previous [tutorial](../accounts/utxo.md), all output coins V1 of a user (whether spent or unspent) have the same public key, which is the public key of the user. Therefore, it is easy to store and retrieve these coins with just the user's public key.
We now present two ways a user can retrieve all of his output coins V1. 

* Via JSON-RPC
```json
{
	"jsonrpc": "1.0",
	"method": "listoutputcoins",
	"params": [
			0,
			9999999,
			[
				{
				"PaymentAddress": "12sm5BLDevJsUbevkJd7eaU9zAiuNixEeUsNnYddnsqEYTrcFMfE7aSS2J6mK3GeHbdT7LMm4VcRETaJCRzzU8xKKa1Tn2t9XcGiqWSDpG7jewQkDeRDY3czMHVEgwWGfUWMvkd2pWr1QpMw1i4s",
				"ReadonlyKey" : "13hZcYDRTsydn5WvfjbraPno7XK1wBPg6ATaYbRYeh6tfr5wgLhma1545K8TPDCLrS4G9GF4AGRzwP7sd4vPvv3XP2WRvAt8Y5YUJcD",
				}
			],
		"0000000000000000000000000000000000000000000000000000000000000004"
		],
	"id": 1
}
```
The first two parameters are deprecated, but they are still needed for old full-node. So you better do not change them. `PaymentAddress` is required to get all TXOs of the user while the `ReadonlyKey` is optional, which is used to decrypt these TXOs. For better privacy, we recommend you leave it empty.
* Via the go-sdk

First, create a new client instance.
```go
client, err := incclient.NewTestNet1Client()
if err != nil {
	log.Fatal(err)
}
```
Then, create a new [`OutCoinKey`](../../../rpchandler/rpc/rpc_coin.go) instance from your private key. The OutCoinKey consists of a payment address, a readonly key, and a privateOTA key.
```go
// create a new OutCoinKey
outCoinKey, err := incclient.NewOutCoinKeyFromPrivateKey(privateKey)
if err != nil {
	log.Fatal(err)
}
outCoinKey.SetReadonlyKey("") // call this if you do not want the remote full-node to decrypt your coin
```

## Example
[utxo.go](../../code/accounts/utxo/utxo.go)
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
	client, err := incclient.NewTestNet1Client()
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

Pretty simple, right? Now, let's [move on](../accounts/submit_key.md) to output coins V2.

---
Return to [the table of contents](../../../README.md).