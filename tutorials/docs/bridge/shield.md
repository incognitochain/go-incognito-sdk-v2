---
Description: Tutorial on how to create an ETH-shielding transaction
---
# Depositing ETH/ERC20 to Incognito
Suppose that we already have a transaction that deposited some ETH/ERC20 to the smart contract. To mint the same amount of pETH/pERC20 inside the Incognito network, we use the function [`CreateAndSendIssuingETHRequestTransaction`](../../../incclient/bridge.go) with the following inputs:
* `privateKey`: our private key to sign the transaction.
* `tokenID`: the pETH/pERC20 tokenID.
* `depositProof`: the Ethereum receipt for the depositing transaction.

## Example
[shield.go](../../code/bridge/shield/shield.go)
(see more at [bridge_test.go](../../../incclient/bridge_test.go))

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
	tokenIDStr := "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
	ethTxHash := "0xb31d963b3f183d60532ca60d534e0113ca56070af795fde450dd456945a7be42"

	ethProof, depositAmount, err := ic.GetETHDepositProof(ethTxHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deposited amount: %v\n", depositAmount)

	txHashStr, err := ic.CreateAndSendIssuingETHRequestTransaction(privateKey, tokenIDStr, *ethProof)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TxHash: %v\n", txHashStr)

	time.Sleep(10 * time.Second)

	fmt.Printf("check shielding status\n")
	for {
		status, err := ic.CheckShieldStatus(txHashStr)
		if err != nil{
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
