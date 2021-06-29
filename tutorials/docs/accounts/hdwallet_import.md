---
Description: Tutorial on how to restore accounts from a mnemonic string.
---
# Restore Accounts from a Mnemonic String
Suppose that you have created your accounts from somewhere else (e.g., Incognito Wallet) and you want to restore all of them using this `go-sdk`. This time, we use the function [`ImportAccount`](../../../incclient/account.go) in the [`incclient`](../../../incclient) package.


```go
wallets, err := client.ImportAccount(mnemonic)
if err != nil {
	log.Fatal(err)
}
```
This is a special function which returns a list of accounts the user has. In this list, the first account is the master account, which can be used to generate the rest child accounts. The rest are accounts that have user's activities on the blockchain (i.e, user has **AT LEAST** 1 in-coming or 1 out-going transaction). If no accounts has been acknowledged with transactions, **ONLY** the master account is returned.

## Example
[hdwallet_import.go](../../code/accounts/hdwallet_import/hdwallet_import.go)

```go
package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}
	mnemonic := "search trophy awake proud sponsor toe lumber toilet sugar smoke soup joke"

	wallets, err := client.ImportAccount(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	for i, w := range wallets {
		privateKey, err := w.GetPrivateKey()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("idx: %v, privateKey: %v\n", i, privateKey)
	}
}

```
---
Return to [the table of contents](../../../README.md).
