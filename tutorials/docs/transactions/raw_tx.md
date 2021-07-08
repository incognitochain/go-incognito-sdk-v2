---
Description: Tutorial on how to create a PRV transaction in Incognito.
---

# Transferring PRV
In this tutorial, we'll learn how to transfer PRV from one account to another account. 
A PRV transaction consists of the receiver addresses, amount for each receiver, and optional data. 
The transaction must then be signed with the private key of the sender before it's broadcast to the network. Subsequently, we'll walk through each
step to create a PRV transaction. 

## Get ourselves connected to the network
```go
client, err := incclient.NewTestNet1Client()
if err != nil {
    log.Fatal(err)
}
```

## Prepare our inputs
In this step, we must prepare a [TxParam](../transactions/params.md) for creating the transaction.
```go
privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
receiverList := []string{"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H"}
amountList := []uint64{10000000}
txVersion := int8(1) //txVersion should be -1, 1, or 2
txFee := uint64(10)

txParam := incclient.NewTxParam(privateKey, receiverList, amountList, txFee, nil, nil, nil)

```
It is required that the number of elements of the `receiverList` must be equal to the number of elements of the `amountList`. Notice that in this case, the `txTokenParam` is set to `nil` because we are transferring PRV. For transaction fee, if it is set to `0`, the default fee will be used.


## Call the creating function
The next step is to call the creating function with prepared inputs.
```go
encodedTx, txHash, err := client.CreateRawTransaction(txParam, txVersion)
if err != nil {
    log.Fatal(err)
}
```
If you are not sure which version of the transaction you want, just simply set it to `-1`.
The function [`CreateRawTransaction`](../../../incclient/tx.go) returns a *base58-encoded* transaction, along with its hash. 

## Submit the transaction
Now we are ready to submit the transaction to the Incognito network. This time, we just need to call the
[`SendRawTx`](../../../incclient/tx.go) function using our client on input the encoded transaction created previously.
```go
err = client.SendRawTx(encodedTx)
if err != nil {
    log.Fatal(err)
}
```

## Example
[raw_tx.go](../../code/transactions/raw_tx/raw_tx.go)

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
	receiverList := []string{"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H"}
	amountList := []uint64{10000000}
	txVersion := int8(1) //txVersion should be -1, 1, or 2
	txFee := uint64(10)

	txParam := incclient.NewTxParam(privateKey, receiverList, amountList, txFee, nil, nil, nil)

	encodedTx, txHash, err := client.CreateRawTransaction(txParam, txVersion)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create and send tx successfully, txHash: %v\n", txHash)
}
```
---
Return to [the table of contents](../../../README.md).
