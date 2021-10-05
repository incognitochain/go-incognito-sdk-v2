---
Description: Tutorial on how to create a token transaction in Incognito.
---
# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

# Transferring Token
The processing of creating a token transaction is quite similar to that of creating a PRV transaction.

## Get ourselves connected to the network
```go
client, err := incclient.NewTestNet1Client() // use `NewTestNet1ClientWithCache()` if you prefer the local UTX cache
if err != nil {
    log.Fatal(err)
}
```

## Prepare our inputs
First, we need to create a [TxParam](../transactions/params.md). 

### Load the private key
As usual, we need a private key.

```go
privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
```

### Describe token information
Unlike the previous [tutorial](), this time, we aslo need a [TxTokenParam]() to describe the token information. 
```go
tokenID := "0000000000000000000000000000000000000000000000000000000000000100"
tokenReceivers := []string{"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H"}
tokenAmounts := []uint64{10000000}
hasTokenFee := false
tokenFee := uint64(0)
txVersion := int8(1)

txTokenParam := incclient.NewTxTokenParam(tokenID, 1, tokenReceivers, tokenAmounts, hasTokenFee, tokenFee, nil)
```

### Create transaction params
The next step is to create a `TxParam` with the token information above.
```go
txParam := incclient.NewTxParam(privateKey, nil, nil, 0, txTokenParam, nil, nil)
```
Notice that in this case, the PRV receivers and PRV amounts are set to `nil` because we are transferring token. 
We set the transaction fee to `0` to indicate that the default fee is used.

## Call the creating function
The next step is to call the creating function with prepared inputs.
```go
encodedTx, txHash, err := client.CreateRawTokenTransaction(txParam, txVersion)
if err != nil {
    log.Fatal(err)
}
```
If you are not sure which version of the transaction you want, just simply set it to `-1`.
The function [`CreateRawTokenTransaction`](../../../incclient/txtoken.go) returns a base58-encoded transaction, along with its hash. 

## Submit the transaction
Now we are ready to submit the transaction to the Incognito network. This time, we call the
[`SendRawTokenTx`](../../../incclient/txtoken.go) function instead of [`SendRawTx`](../../../incclient/tx.go).
```go
err = client.SendRawTokenTx(encodedTx)
if err != nil {
    log.Fatal(err)
}
```

## Example
[raw_tx_token.go](../../code/transactions/raw_tx_token/raw_tx_token.go)

```go
package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"

	tokenID := "0000000000000000000000000000000000000000000000000000000000000100"
	tokenReceivers := []string{"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H"}
	tokenAmounts := []uint64{10000000}
	hasTokenFee := false
	tokenFee := uint64(0)
	txVersion := int8(1)

	txTokenParam := incclient.NewTxTokenParam(tokenID, 1, tokenReceivers, tokenAmounts, hasTokenFee, tokenFee, nil)
	txParam := incclient.NewTxParam(privateKey, nil, nil, 0, txTokenParam, nil, nil)

	encodedTx, txHash, err := client.CreateRawTokenTransaction(txParam, txVersion)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create and send tx token successfull, txhash: %v\n", txHash)
}
```

Next, we'll see how we can [create our own tokens](../transactions/init_token.md).

---
Return to [the table of contents](../../../README.md).
