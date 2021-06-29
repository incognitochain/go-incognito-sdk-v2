---
Description: Tutorial on how to convert UTXOs from V1 to V2 in Incognito.
---

# UTXOs Conversion
In this tutorial, we'll learn how to convert UTXOs from version 1 to version 2. This process is mandatory
for those who still have UTXOs version 1 and wish to use them in version 2 to enhance privacy. Besides, the Incognito
network will soon not accept transactions of version 1. Therefore, UTXOs conversion is a must.

Please notice that conversion transactions are (by default) non-private. All information (sender, amount, tokenID) will be publicly visible.

Please make changes to the following example to convert your UTXOs.

## Example
[convert.go](../../code/transactions/convert/convert.go)

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
	tokenID := "00000000000000000000000000000000000000000000000000000000000000ff"
	txHash, err := client.CreateAndSendRawConversionTransaction(privateKey, tokenID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Initialize a token successfully, txHash: %v\n", txHash)
}

```
---
Return to [the table of contents](../../../README.md).
