---
Description: Tutorial on transaction parameters when using the `incclient`.
---

## Transaction Parameters
When creating a transaction, we need to know the receivers, amount for each, the transaction fee, etc,.
The transaction must then be signed with the sender's private key. 

The following [TxParam](../../../incclient/common.go) object specifies these pieces of information.

```go
type TxParam struct {
	senderPrivateKey string
	receiverList     []string
	amountList       []uint64
	fee              uint64
	txTokenParam     *TxTokenParam
	md               metadata.Metadata
	kArgs            map[string]interface{}
}
```
* `senderPrivateKey`: the private key of the sender.
* `receiverList`: the list of payment addresses of corresponding PRV receivers.
* `amountList`: the list of PRV amounts for each receiver.
* `fee`: the transaction fee, paid in PRV; if `fee` is set to `0`, the default fee will be used.
* `txTokenParam`: the information describing token transactions (see below);
    * if this is a PRV transaction, `txTokenParam` is set to `nil`;
    * otherwise, it must not be nil.
* `md`: a metadata for indicating special transactions (trading, staking, etc,.).
* `kArgs`: a redundant parameter which is used for the sake of later extension; right now, it is usually set to `nil`

The [TxTokenParam](../../../incclient/common.go) struct is the description of token information.
```go
type TxTokenParam struct {
	tokenID      string
	tokenType    int
	receiverList []string
	amountList   []uint64
	hasTokenFee  bool
	tokenFee     uint64
	kArgs        map[string]interface{}
}
```
* `tokenID`: the ID of the token.
* `tokenType`: type of the token transaction, `0` is token initialization, `1` is token transferring.
* `receiverList`: the list of payment addresses of corresponding token receivers.
* `amountList`: the list of token amounts for each token receiver.
* `hasTokenFee`: the indicator of whether the transaction should pay fee in PRV or in the token. Notice that paying the transaction fee in token is only supported in transactions of version `1`.
* `tokenFee`: the transaction fee, paid in token; in case `tokenFee` is set to `0`, 
    * if `hasTokenFee = false`, the transaction pays fee in PRV
    * if `hasTokenFee = true`, the `incclient` will automatically calculate the token fee based on the current price. 
* `kArgs`: a redundant parameter which is used for the sake of later extension; right now, it is usually set to `nil`
---
Return to [the table of contents](../../../README.md).