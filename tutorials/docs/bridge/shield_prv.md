---
Description: Tutorial on how to shield pegged-PRV into the Incognito network
---

# Shielding pegged-PRV

This is the same as shielding an EVM token except for the
function [`CreateAndSendIssuingPRVPeggingRequestTransaction`](../../../incclient/prv_pegging.go) is used instead of the
function [`CreateAndSendIssuingEVMRequestTransaction`](../../../incclient/bridge.go).

## Example

[shield_prv.go](../../code/bridge/shield_prv/shield_prv.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	ic, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	evmTxHash := "0xb31d963b3f183d60532ca60d534e0113ca56070af795fde450dd456945a7be42"
	isBSC := false

	evmProof, depositAmount, err := ic.GetEVMDepositProof(evmTxHash, isBSC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deposited amount: %v\n", depositAmount)

	txHashStr, err := ic.CreateAndSendIssuingPRVPeggingRequestTransaction(privateKey, *evmProof, isBSC)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TxHash: %v\n", txHashStr)

	time.Sleep(10 * time.Second)

	fmt.Printf("check shielding status\n")
	for {
		status, err := ic.CheckShieldStatus(txHashStr)
		if err != nil {
			log.Fatal(err)
		}
		if status == 1 || status == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Printf("shielding status: %v\n", status)
		break
	}
}
```
---
Return to [the table of contents](../../../README.md).
