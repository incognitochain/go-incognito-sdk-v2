---
Description: Tutorial on how to shield pegged-PRV into the Incognito network
---

# Before Going Further

Please read through the tutorials on [key submission](../../accounts/submit_key.md)
and [UTXO cache](../../accounts/utxo_cache.md) for proper balance and UTXO retrieval. Skip these parts if you're familiar
with these notions.

# Shielding pegged-PRV

This is the same as shielding an EVM token except for the
function [`CreateAndSendIssuingPRVPeggingRequestTransaction`](../../../../incclient/prv_pegging.go) is used instead of the
function [`CreateAndSendIssuingEVMRequestTransaction`](../../../../incclient/bridge.go).

## Example

[shield_prv.go](../../../code/bridge/evm/shield_prv/shield_prv.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"log"
	"time"
)

func main() {
	ic, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "YOUR_PRIVATE_KEY_HERE"
	evmTxHash := "" //the PRV deposit transaction hash on the EVM network
	// specify which EVM network we are interacting with. evmNetworkID could be one of the following:
	// 	- rpc.ETHNetworkID
	//	- rpc.BSCNetworkID
	evmNetworkID := rpc.ETHNetworkID

	evmProof, depositAmount, err := ic.GetEVMDepositProof(evmTxHash, evmNetworkID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deposited amount: %v\n", depositAmount)

	txHashStr, err := ic.CreateAndSendIssuingPRVPeggingRequestTransaction(privateKey, *evmProof, evmNetworkID)
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
Return to [the table of contents](../../../../README.md).
