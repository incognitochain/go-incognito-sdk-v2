---
description: Tutorial on key types in Incognito.
---

## Keys

The followings are all types of keys a user might possess. Each key used in this `go-sdk` is encoded with base58. 
* `privateKey`: the master private key which is used to spend UTXOs, sign transactions, and check the spending status of TXOs. An example of a private key is `1111111ChQjKmu2wZUAYoHWyD5ZiXFNKSWV8hyfMZkGqr74Lj3KGjpFgY18Voum9B2ADfXiYJhoaSdBC8D7u6XSPRAxw8sUBxtEjaY1DJdQ`.
* `ReadOnlyKey`: this key is used to decrypt the amount of each TXO belonging to the user. An example of a readonly key is `13hdbKkqXZWYRYk9k4LxzTEgqXznAUvGwp6nAUa2fA1UjPuA6cKL23pinwQFBaPDUZAP5H8CDhUHPd7wgP1bDp43q2RmSswgRyTVAPG`.
* `PrivateOTAKey`: this key is used to check if a TXO belong to the user or not (in Privacy V2 only). An example of a private OTA key is `14yFBV8kHMawSGQdN8XKq9DcmG2yNVuo42JbRTEck78iemqZLrHGt6Sm8idTS96kWCN6njHLxtk9BuhKXfuQqxmHQs6nBbc6LkgDPPC`.
* `PaymentAddress`: this is the address to receive funds of a user. An example of a payment address is `12S21ANvmmS9Qq5g7zAc9mFkuoxaxTvkqbtxuVuSaCtEp1xkma1iSw7bvUzqeZXympSjiqFTyHcPJdZRfbDHqqW8hQAG8oEA7vhTBWP`.
* `MiningKey`: this is the key used to stake nodes. An example of a payment address is `1csptwXqLVZFMb2jPkEwiqQmJUUKkEjeWTm1WAQXiYjnfxG8tH`.
* `MiningPubKey`: this is the public key corresponding to the `MiningKey` above. An example of a payment address is `121VhftSAygpEJZ6i9jGkKtR4pAHfzZ88HCwNKX6SpoNwGZSNRWcU424si7zKNQQbbEHx9T17Vrk321NeMif3XNaDimZhVgwp8mgv1aRbzhqqSVaX2sRtsymeLGqM3bNMcsNeSHWsf5ZQKU5RtvnHSCPwQV5vWkqJjpQq5dTsgEmmzvizXPUNPuysArp2kNTKJmF2nDRsta3kBFfu5YUwWrwyq23x9LnSpmeByLAfJaa4Bu1W47gNzPhQ3A29JPWKb7ikY3BnYt8pZ4Tao8HzyxfzKha6JPTzW3zZRm4ygwyRQghY31cH15o8p4JTGRL8X6uDbJCeMuTz4vtuJud3NtsGxy7yq8aSfPHRzbdt54nv3GpNVS34SBw4AoFMRVTNBW5HqUAUM9K56CKELRMn5YdkPhscA5RbPTAG9xJCAQjs43n`.

To generate an account on Incognito using this SDK, import the `wallet` package.
```go
keyWallet, err := wallet.NewMasterKeyFromSeed(common.RandBytes(32))
if err != nil {
log.Fatal(err)
}
```

To retrieve all information of a private key, call the [`GetAccountInfoFromPrivateKey`](../../../incclient/account_common.go) function.
## Example
[keys.go](../../code/accounts/keys/keys.go)

```go
package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func main() {
	keyWallet, err := wallet.NewMasterKeyFromSeed(common.RandBytes(32))
	if err != nil {
		log.Fatal(err)
	}

	privateKey := keyWallet.Base58CheckSerialize(wallet.PrivateKeyType)
	accInfo, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AccountInfo: %v\n", accInfo.String())
}


```
---
Return to [the table of contents](../../../README.md).
