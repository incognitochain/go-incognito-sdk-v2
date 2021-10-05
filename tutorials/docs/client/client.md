---
Description: Tutorial on how to set up a client to connect to Incognito with Go.
---

## Setting up the Client

The Incognito network allows anyone to connect by an RPC client. The main-net end-point is available at `https://mainnet.incognito.org/fullnode`, 
while the test-net end-point is located at `https://testnet.incognito.org/fullnode`. 

To interact with the Incognito network, first import the `incclient` package and initialize an Incognito client by calling `NewMainNetClient` which by default connects to the mainnet end-point above. If you wish to connect to the testnet, try `NewTestNetClient`.

To use the local [UTXO cache layer](../accounts/utxo_cache.md), try initializing a client with a post-fix `WithCache` (e.g, `NewMainNetClientWithCache`).

## Examples
[client.go](../../code/client/client.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	incClient, err := incclient.NewMainNetClient() // or use incclient.NewMainNetClientWithCache()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Let's get yourself into the Incognito network!")
	_ = incClient
}
```
---
Return to [the table of contents](../../../README.md).
