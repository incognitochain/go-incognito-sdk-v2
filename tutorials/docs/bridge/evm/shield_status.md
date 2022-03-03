---
description: Tutorial on how to check the status of a shield request
---
# Checking Shielding Requests
This sdk allows you to check the status of a (BSC/ETH/ERC20) shielding transaction via the function [`CheckShieldStatus`](../../../../incclient/bridge.go). All you need to do is to supply it with an Incognito transaction hash. Possile returned values are as follows.
* -1: an error has occurred, see the the error for more detail.
* 0: the transaction is not found or it is not a shielding transaction.
* 1: the transaction is pending.
* 2: the shielding transaction is accepted.
* 3: the shielding transaction is rejected.

See the following example.

## Example
[shield_status](../../../code/bridge/evm/status/status.go)

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
	
	txHash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	status, err := client.CheckShieldStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("status: %v\n", status)
}
```
---
Return to [the table of contents](../../../../README.md).