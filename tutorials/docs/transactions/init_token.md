---
description: Tutorial on how to create a custom token in Incognito.
---

# Create Custom Tokens
Assume that we have already connected to the Incognito network. 
We first describe the information of the token being created such as `tokenName`, `tokenSymbol` and the total market cap.
```go
privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
tokenName := "INC"
tokenSymbol := "INC"
tokenCap := uint64(1000000000000)
tokenVersion := 1
```
The Client has 3 functions that support us to create our new token:
* [`CreateTokenInitTransactionV1`](../../../incclient/txtoken.go): this function creates a token version 1; it receives as input a private key, the token name, the token symbol and the total market cap. It returns a base58-encoded token transaction.
* [`CreateTokenInitTransactionV2`](../../../incclient/txtoken.go): its parameters are the same as version 1, however, the returned transaction is a PRV transaction. In addition, it requires us to have PRV of version 2 to create such a transaction (merely to pay the transaction fee).
* [`CreateTokenInitTransaction`](../../../incclient/txtoken.go): this is wrapper function that will automatically call the appropriate function based on the provided additional parameter `tokenVersion`. Set `tokenVersion = -1` if we are not sure about which version we want.

Notice that we don't need to create a [TxParam](../../../incclient/common.go#L14) in this case because of the help of these functions.

Here, we are creating a token of version `1`, with name `INC`, symbol `INC` and the total cap is `1000000000000`.
```go
encodedTx, txHash, err := client.CreateTokenInitTransaction(privateKey, tokenName, tokenSymbol, tokenCap, tokenVersion)
if err != nil {
	log.Fatal(err)
}
```

After we have the encoded transaction, just simply call the right broadcasting function. If the token version is `1`, call `SendRawTokenTx`; otherwise, call `SendRawTx`.
```go
if tokenVersion == 1 {
		err = client.SendRawTokenTx(encodedTx)
	} else {
		err = client.SendRawTx(encodedTx)
	}
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example
[init_token.go](../../code/transactions/init_token/init_token.go)

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
	tokenName := "INC"
	tokenSymbol := "INC"
	tokenCap := uint64(1000000000000)
	tokenVersion := 1
	
	encodedTx, txHash, err := client.CreateTokenInitTransaction(privateKey, tokenName, tokenSymbol, tokenCap, tokenVersion)
	if err != nil {
		log.Fatal(err)
	}
	
	if tokenVersion == 1 {
		err = client.SendRawTokenTx(encodedTx)
	} else {
		err = client.SendRawTx(encodedTx)
	}
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Initialize a token successfully, txHash: %v\n", txHash)
}
```

Enough with basic transactions, let's move on to a more interesting part, [pDEX](../pdex/query.md).

---
Return to [the table of contents](../../../README.md).
