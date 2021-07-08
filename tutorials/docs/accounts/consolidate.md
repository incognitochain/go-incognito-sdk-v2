---
Description: Tutorial on how to consolidate UTXOs of an account.
---
Many people have been experiencing issues when sending large transactions. The cause is primarily due to a large number of UTXOs, resulting in the sending transactions exceeding the maximum allowable size. Spending cryptocurrency results in the accumulation of UTXOs in the same way that spending hard cash results in the accumulation of spare change. When it comes to spending money, rummaging through a pocket full of change is inefficient, and excessive UTXOs cause undesirable behaviors such as slow or unsuccessful transactions. We'll show you how to condense your UTXOs and free up more space for transactions in this tutorial, similar to exchanging loose coins for larger bills.

## Consolidate UTXOs
Consolidation is done by calling the function [`Consolidate`](../../../incclient/tx_consolidate.go) in the [`incclient`](../../../incclient) package.
The first step is to prepare the parameters. `Consolidate` uses a number of threads working simultaneously to boost up the consolidating process. For each thread, it combines 30 UTXOs (maximum input size of a transaction) at a time into a single UTXO. They will stop when the number of UTXOs is less than 10.

To consolidate your UTXOs, we first prepare the parameters.
```go
privateKey := "YOUR_PRIVATE_KEY" // input your private key
tokenIDStr := common.PRVIDStr
version := int8(1) // current version of the main-net is 1
numThreads := 20
```
Finally, just call the [`Consolidate`](../../../incclient/tx_consolidate.go) function and wait. It takes time, be patient.
```go
txList, err := client.Consolidate(privateKey, tokenIDStr, version, numThreads)
if err != nil {
    log.Printf("txList: %v\n", txList)
    log.Fatal(err)
}
log.Printf("txList: %v\n", txList)
```

See the full example below.
## Example
[consolidate.go](../../code/accounts/consolidate/consolidate.go)
```go
package main

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewMainNetClient()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "YOUR_PRIVATE_KEY" // input your private key
	tokenIDStr := common.PRVIDStr
	version := int8(1)
	numThreads := 20

	txList, err := client.Consolidate(privateKey, tokenIDStr, version, numThreads)
	if err != nil {
		log.Printf("txList: %v\n", txList)
		log.Fatal(err)
	}
	log.Printf("txList: %v\n", txList)
}
```
---
Return to [the table of contents](../../../README.md).