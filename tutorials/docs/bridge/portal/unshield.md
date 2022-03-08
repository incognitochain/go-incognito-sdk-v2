---
Description: Tutorial on how to create a Portal un-shielding transaction
---

# Before Going Further

Please read through the tutorials on [key submission](../../accounts/submit_key.md)
and [UTXO cache](../../accounts/utxo_cache.md) for proper balance and UTXO retrieval. Skip these parts if you're familiar
with these notions.

# Withdrawing an Portal token from the Incognito Network

The very first step in withdrawing a Portal token from the Incognito network to the public networks is to burn
the corresponding pPortal token inside the Incognito network. 
This is done using the function [`CreateAndSendPortalUnShieldTransaction`](../../../../incclient/bridge.go) supplied with a corresponding
external address. See the below example.

## Example

[unshield.go](../../../code/bridge/portal/unshield/unshield.go)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	ic, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "YOUR_PRIVATE_KEY"
	externalAddr := "EXTERNAL_ADDRESS"
	tokenIDStr := "PORTAL_TOKEN"
	unShieldAmount := uint64(1000000)

	txHash, err := ic.CreateAndSendPortalUnShieldTransaction(
		privateKey, tokenIDStr, externalAddr, unShieldAmount, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TxHash: %v\n", txHash)

	time.Sleep(100 * time.Second)
	status, err := ic.GetPortalUnShieldingRequestStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}

	jsb, _ := json.Marshal(status)
	fmt.Println(string(jsb))
}
```

After the burning process is successful, wait for about 30 minutes for the fund to be released. One can check 
the status of the un-shielding request via function `GetPortalUnShieldingRequestStatus`.

---
Return to [the table of contents](../../../../README.md).
