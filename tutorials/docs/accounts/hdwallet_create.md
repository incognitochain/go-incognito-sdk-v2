---
Description: Tutorial on how to create HD wallets.
---
# HD Wallets
A **hierarchical deterministic** (HD) wallet is a cryptocurrency wallet that generates all addresses from a single master seed phrase and automatically creates a hierarchical (tree-like) structure to govern your private keys.

Only with one mnemonic phrase, you'll be able to back up and recover all of your keys. You've probably used this feature if you've used a wallet like Trust Wallet.

## Incognito HD wallets
All Incognito HD wallets use a variant of the standard 12-word master seed key, and each time this seed is extended at the end by a counter value, allowing an infinite number of new addresses to be generated automatically.
These wallets following the [BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) standard with the `coinType = 587` (see this [post](https://github.com/satoshilabs/slips/blob/master/slip-0044.md)). An example of a key derivation path is 
<pre>
m/44'/587'/0'/0/1
</pre>

## Create a new HD wallet
HD wallets are implemented in the `wallet` package. To use HD wallets, call the function [`NewMasterKey`](../../../wallet/hdwallet.go).
```go
// create a master wallet
masterWallet, mnemonic, err := wallet.NewMasterKey()
if err != nil {
	log.Fatal(err)
}
fmt.Printf("mnemonic: `%v`\n", mnemonic)
```
It returns a mnemonic string, and the master wallet for this mnemonic. This walllet will be used to derive all other keys in the chain. 

## Derive child keys
Call the [`DeriveChild`](../../../wallet/hdwallet.go) with the child index to generate a new wallet. For every mnemonic, the first account is at `index = 1`. 
```go
// derive the first wallet from the master (i.e., the `Anon` account on the Incognito wallet)
childIdx := uint32(1)
firstWallet, err := masterWallet.DeriveChild(childIdx)
if err != nil {
	log.Fatal(err)
}
```

## Example
[hdwallet_create.go](../../code/accounts/hdwallet_create/hdwallet_create.go)

```go
package main

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

func main() {
	// create a master wallet
	masterWallet, mnemonic, err := wallet.NewMasterKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("mnemonic: `%v`\n", mnemonic)

	// derive the first wallet from the master (i.e., the `Anon` account on the Incognito wallet)
	childIdx := uint32(1)
	firstWallet, err := masterWallet.DeriveChild(childIdx)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := firstWallet.GetPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("firstPrivateKey: %v\n", privateKey)

}
```
---
Return to [the table of contents](../../../README.md).
