---
Description: Tutorial on how to retrieve balances of an Incognito account.
---

## Balances

To retrieve the balance of either PRV or any token on the Incognito network, we first need to get ourselves connected to the network
using the `incclient` package as described in [Client](../client/client.md).

```go
client, err := incclient.NewTestNet1Client()
```

Now, reading the balance of a token is pretty simple, just simply call the `GetBalance` function of the client with the inputs consisting of your private key and the tokenID.

```go
balance, err := client.GetBalance(privateKey, common.PRVStr)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("balance: %v\n", balance)
```

Note that the private key is required to check if a TXO has been spent or not. The private key will never leave your local machine.

## Example
[balances.go](../../code/accounts/balances/balances.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	incClient, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	balancePRV, err := incClient.GetBalance(privateKey, common.PRVIDStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("balancePRV: %v\n", balancePRV)

	tokenID := "0000000000000000000000000000000000000000000000000000000000000100"
	balanceToken, err := incClient.GetBalance(privateKey, tokenID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("balanceToken: %v\n", balanceToken)
}
```
---
Return to [the table of contents](../../../README.md).
