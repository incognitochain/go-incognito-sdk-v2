---
Description: Tutorial on how to create an EVM un-shielding transaction
---
# Withdrawing an EVM token from the Incognito Network
The very first step in withdrawing an EVM token from the Incognito network to the main Ethereum/BSC networks is to burn the corresponding pEVM token inside the Incognito network and have the beacon generate a burn proof for us.
This is done using the function [`CreateAndSendBurningRequestTransaction`](../../../incclient/bridge.go) supplied with an Ethereum address.

## Example
[unshield.go](../../code/bridge/unshield/unshield.go)
```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	ic, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	tokenIDStr := "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
	remoteAddr := "b446151522b8f1c9d27cacedce93f398a016f84337c1b79fc54c8436af5f7900"
	burnedAmount :=  uint64(50000000)

	txHash, err := ic.CreateAndSendBurningRequestTransaction(privateKey, remoteAddr, tokenIDStr, burnedAmount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TxHash: %v\n", txHash)
}
```

After the burning process is successful, the next step is to retrieve the burn proof and submit it to the corresponding smart contract. In this step, an EVM account is required to interact with the smart contract. 
This step is out of the scope of this tutorial series. See example on the full flow of un-shielding [here](https://github.com/incognitochain/incognito-cli/blob/development/bridge_evm.go).

---
Return to [the table of contents](../../../README.md).
