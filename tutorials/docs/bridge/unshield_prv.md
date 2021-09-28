---
Description: Tutorial on how to withdraw PRV to EVM networks
---
# Withdraw PRV to EVM networks
By the end of September 2021, Incognito allows users to withdraw their PRV to the Ethereum/Binance Smart Chain networks. This is enabled by the implementations of two pegged-PRV tokens ([ERC20](https://etherscan.io/address/0xB64fde8f199F073F41c132B9eC7aD5b61De0B1B7#code) / [BEP20](https://bscscan.com/address/0xB64fde8f199F073F41c132B9eC7aD5b61De0B1B7)).
The withdrawing procedure is pretty much the same as that of an EVM token:
* The first step is to burn the PRV inside the Incognito network. This is done using the function [`CreateAndSendBurningPRVPeggingRequestTransaction`](../../../incclient/prv_pegging.go). This step also needs to specify which network will the PRV be withdrawn to (using param `isBSC`, defaults to `false`).
* The second step is to retrieve the burn proof from the beacon chain. This is done via the function [`GetBurnPRVPeggingProof`](../../../incclient/prv_pegging.go).
* Finally, we submit the burn proof to the designated pegged-PRV smart contract.

## Example
[unshield_prv.go](../../code/bridge/unshield_prv/unshield_prv.go)
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
	remoteAddr := "b446151522b8f1c9d27cacedce93f398a016f84337c1b79fc54c8436af5f7900"
	burnedAmount := uint64(50000000)
	isBSC := false

	// burn PRV
	txHash, err := ic.CreateAndSendBurningPRVPeggingRequestTransaction(privateKey, remoteAddr, burnedAmount, isBSC)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TxHash: %v\n", txHash)

	// wait for the above tx to reach the beacon chain
	time.Sleep(50 * time.Second)

	// retrieve the burn proof
	prvBurnProof, err := ic.GetBurnPRVPeggingProof(txHash, isBSC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prvBurnProof)
}
```

After the burning process is successful, the next step is to retrieve the burn proof and submit it to the corresponding smart contract. In this step, an EVM account is required to interact with the smart contract.
This step is out of the scope of this tutorial series. See example on the full flow of un-shielding [here](https://github.com/incognitochain/incognito-cli/blob/development/bridge_evm.go).

---
Return to [the table of contents](../../../README.md).