---
description: Tutorial on how to retrieve the history of an account.
---
Currently, thousands of transactions take place on the Incognito chain every day. A user might have done a few to thousands of transactions. Therefore, the need to retrieve the history of these transaction is inevitable. However, there is not a known method to extract an account's history for users who do not use the Incognito wallet. In this tutorial, we will show you how to perform this task with the go-sdk.

The sdk provides two methods for user to get their account's history.
* Sequential: this method retrieves the history sequentially. It best suits accounts with a few transactions. This method is a part of an [`IncClient`](../../../incclient/tx_history.go) object.
* Parallel: this method creates several workers working simultaneously to boost up the speed. It best suits accounts with many transactions (says 100+). Also, this mode requires more CPU than the previous one, and its performance depends on the device configuration. This method is a part of a [`TxHistoryProcessor`](../../../incclient/tx_history_worker.go) object.

Here are some functions that are used to retrieve the history of an account. Note that all of these functions are for a specific token only. Right now, users have to provide the tokenID for each history request.
## Functions
Name | Description | Object | Status
-------------|-------------|-------------|-----------
GetListTxsInV1| Get all in-coming transactions v1 of a token | IncClient | Finished
GetListTxsInV2 | Get all in-coming transactions v2 of a token | IncClient | In Progress
GetListTxsOutV1 | Get all out-going transactions v1 of a token | IncClient | Finished
GetListTxsOutV2 | Get all out-going transactions v1 of a token | IncClient | In Progress
GetTxHistoryV1 | Get the history V1 of a token | IncClient | Finished
GetTxHistoryV2 | Get the history V2 of a token | IncClient | In Progress
GetTxsIn | Get in-coming transactions of a token (with the given version) | TxHistoryProcessor | Finished (v1) & In Progress (v2)
GetTxsOut | Get out-going transactions of a token (with the given version) | TxHistoryProcessor | Finished (v1) & In Progress (v2)
GetTokenHistory | Get the history of a token (with the given version) | TxHistoryProcessor | Finished (v1) & In Progress (v2)


The followings are the descriptions of an in-coming transaction as well as an out-going transaction, respectively.

```go
// TxIn is an in-coming transaction.
// A transaction is considered to be a TxIn if
//	- it receives at least 1 output coin; and
//	- the input coins are does not belong to the receivers.
// In case a user A sends some coins to a user B and A receives a sent-back output coin, this is not considered to
// be a TxIn.
type TxIn struct {
	Version  int8
	LockTime int64
	TxHash   string
	Amount   uint64
	TokenID  string
	Metadata metadata.Metadata
}
```

```go
// TxOut is an out-going transaction.
// A transaction is considered to be a TxOut if it spends input coins.
type TxOut struct {
	Version    int8
	LockTime   int64
	TxHash     string
	Amount     uint64
	TokenID    string
	SpentCoins map[string]uint64 // map from the coin's serialNumber to its amount
	Receivers  []string
	PRVFee     uint64
	TokenFee   uint64
	Metadata   metadata.Metadata
}
```
And the history looks like
```go
// TxHistory consists of a list of TxIn's and a list of TxOut's.
type TxHistory struct {
	TxInList  []TxIn
	TxOutList []TxOut
}
```

Now, let's get to the detail of how we can use these functions.

## Connect to the network
Because most of currently running full-nodes don't have some required RPCs ad data, the created client in this case must point to the designated address, which is `SOME_URL_HERE`.
```go
// For main-net
client, err := incclient.NewIncClient("https://beta-fullnode.incognito.org/fullnode", incclient.MainNetETHHost, 1)
if err != nil {
	log.Fatal(err)
}
```

## Specify the private key and the tokenID
```go
tokenIDStr := common.PRVIDStr // input the tokenID in which you want to retrieve the history of.
privateKey := "YOUR_PRIVATE_KEY" // input your private key here
```

## Retrive the history
### Sequential method
```go
// get the history in a normal way.
h, err := client.GetTxHistoryV1(privateKey, tokenIDStr)
if err != nil {
	log.Fatal(err)
}
```
If you want to get only the in-coming (out-going) history, consider using `GetListTxsInV1` (`GetListTxsOutV1`).	

### Parallel method
To run in this mode, we have to create a `TxHistoryProcessor`.
```go
numWorkers := 15
p := incclient.NewTxHistoryProcessor(client, numWorkers)
```
The number of workers depends on the local device performance. Keeping it at 5-15 is an advisible choice. If this number is too high, the remote full-node may drop the connection to the device due to limited number of requests at a time. Finally, call the required functions.
```go
h, err := p.GetTokenHistory(privateKey, tokenIDStr)
if err != nil {
	log.Fatal(err)
}
```
If you want to get only the in-coming (out-going) history, consider using `GetTxsIn` (`GetTxsOut`).	

### Print the result
```go
totalIn := uint64(0)
	log.Printf("TxsIn\n")
	for _, txIn := range h.TxInList {
		totalIn += txIn.Amount
		log.Printf("%v\n", txIn.String())
	}
	log.Printf("Finished TxsIn\n\n")

	totalOut := uint64(0)
	log.Printf("TxsOut\n")
	for _, txOut := range h.TxOutList {
		totalOut += txOut.Amount
		log.Printf("%v\n", txOut.String())
	}
	log.Printf("Finished TxsOut\n\n")
```
An example result looks like this.
```
2021/06/24 10:47:31 TxsIn
2021/06/24 10:47:31 Timestamp: 2021-06-06T05:04:02, Detail: {"Version":1,"LockTime":1622952242,"TxHash":"bad93fd599a5cf4807f01e392a2652e71d767ba4173a6e0ed93304f03c9d7040","Amount":2894758,"TokenID":"0000000000000000000000000000000000000000000000000000000000000004","Metadata":{"Type":92,"TradeStatus":"accepted","RequestedTxID":"178565e75ac7788ed14b26796add341a1db3b6ff277423325bb49b9d31513802"}}
2021/06/24 10:47:31 Timestamp: 2021-06-01T16:29:44, Detail: {"Version":1,"LockTime":1622561384,"TxHash":"4ffa23b4e31fd45a09015008e09f42f63205db01b0c271e2181cec1dba96a617","Amount":1648827596,"TokenID":"0000000000000000000000000000000000000000000000000000000000000004","Metadata":null}
...
2021/06/24 10:47:31 Timestamp: 2021-05-18T11:24:24, Detail: {"Version":1,"LockTime":1621333464,"TxHash":"21fb57720e1a795fabcff86ba86b2c5701c39793e781a1a560564bbb1d9ed552","Amount":4139168721,"TokenID":"0000000000000000000000000000000000000000000000000000000000000004","Metadata":null}
2021/06/24 10:47:31 Timestamp: 2021-05-18T11:23:06, Detail: {"Version":1,"LockTime":1621333386,"TxHash":"c314b9c09b94cb4122d3264d7d38a56aacd77b9591f091c94194ef2f2d01836c","Amount":3485021695,"TokenID":"0000000000000000000000000000000000000000000000000000000000000004","Metadata":null}
2021/06/24 10:47:31 Finished TxsIn

2021/06/24 10:47:31 TxsOut
2021/06/24 10:47:31 Timestamp: 2021-06-06T05:01:42, Detail: {"Version":1,"LockTime":1622952102,"TxHash":"257e00e9a6f072f053c1ebf356ec10c1900040269ee9f9b99c156a4525045d75","Amount":500000100,"TokenID":"0000000000000000000000000000000000000000000000000000000000000004","SpentCoins":{"1AEF6JymSaBkxRHzHTGTDNf5QQJuiNzyWPVeqS9eJWpRTH27Ga":648827376},"Receivers":["1y4gnYS1Ns2K7BjQTjgfZ5nTR8JZMkMJ3CTGMj2Pk7CQkSTFgA","1uJC6JzhcdERhthksMM6zaiTT9aXTowRkE2odayUBJspHjncyW"],"PRVFee":60,"TokenFee":0,"Metadata":{"TokenIDToBuyStr":"02a41194b536aa20960fd62bd4937a895fcc7c7d84a83bf212a349df2b6ea1f2","TokenIDToSellStr":"0000000000000000000000000000000000000000000000000000000000000004","SellAmount":500000000,"MinAcceptableAmount":863801896,"TradingFee":100,"TraderAddressStr":"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H","Type":91}}
2021/06/24 10:47:31 Timestamp: 2021-06-06T04:52:14, Detail: {"Version":1,"LockTime":1622951534,"TxHash":"721a854982db05c79b6afe8437d614bb3f133daef2b5ea698f08121146232858","Amount":100,"TokenID":"02a41194b536aa20960fd62bd4937a895fcc7c7d84a83bf212a349df2b6ea1f2","SpentCoins":{"1sA8Phz78Jfr6y3vM2EGoWpzWL7GjmQWiU9a6qpeymfjgNkthW":262392374},"Receivers":["1y4gnYS1Ns2K7BjQTjgfZ5nTR8JZMkMJ3CTGMj2Pk7CQkSTFgA","1uJC6JzhcdERhthksMM6zaiTT9aXTowRkE2odayUBJspHjncyW"],"PRVFee":60,"TokenFee":0,"Metadata":{"TokenIDToBuyStr":"0795495cb9eb84ae7bd8c8494420663b9a1642c7bbc99e57b04d536db9001d0e","TokenIDToSellStr":"02a41194b536aa20960fd62bd4937a895fcc7c7d84a83bf212a349df2b6ea1f2","SellAmount":50000,"MinAcceptableAmount":28935,"TradingFee":100,"TraderAddressStr":"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H","Type":205}}
...
2021/06/24 10:47:31 Timestamp: 2021-05-19T04:26:10, Detail: {"Version":1,"LockTime":1621394770,"TxHash":"9f819ae232006d276d14f14115f524433e5cff2bfd7232521b75f880c9a75ffa","Amount":2428625,"TokenID":"0000000000000000000000000000000000000000000000000000000000000004","SpentCoins":{"12ZieWaLeERdUx2pLjnQd9p1XJSLa3pBaXuk8hML6CGu5GPoJQN":9338862140},"Receivers":["1y4gnYS1Ns2K7BjQTjgfZ5nTR8JZMkMJ3CTGMj2Pk7CQkSTFgA","1uJC6JzhcdERhthksMM6zaiTT9aXTowRkE2odayUBJspHjncyW"],"PRVFee":2,"TokenFee":0,"Metadata":{"TokenIDToBuyStr":"ef80ac984c6367c9c45f8e3b89011d00e76a6f17bd782e939f649fcf95a05b74","TokenIDToSellStr":"0000000000000000000000000000000000000000000000000000000000000004","SellAmount":2234900,"MinAcceptableAmount":3865261,"TradingFee":193725,"TraderAddressStr":"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H","SubTraderAddressStr":"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H","Type":205}}
2021/06/24 10:47:31 Timestamp: 2021-05-19T04:22:07, Detail: {"Version":1,"LockTime":1621394527,"TxHash":"48c24524239c17486b249203dc7bb60d3aaf3f7e4eb8f1ebe70d90e26aa1896b","Amount":196674,"TokenID":"ef80ac984c6367c9c45f8e3b89011d00e76a6f17bd782e939f649fcf95a05b74","SpentCoins":{"1pLBCMNYtxRSZt9GCvuTNan6YvzTTr1eVq93vnauRhsKSxtW1w":9339058818},"Receivers":["1y4gnYS1Ns2K7BjQTjgfZ5nTR8JZMkMJ3CTGMj2Pk7CQkSTFgA","1uJC6JzhcdERhthksMM6zaiTT9aXTowRkE2odayUBJspHjncyW"],"PRVFee":4,"TokenFee":0,"Metadata":{"TokenIDToBuyStr":"0000000000000000000000000000000000000000000000000000000000000004","TokenIDToSellStr":"ef80ac984c6367c9c45f8e3b89011d00e76a6f17bd782e939f649fcf95a05b74","SellAmount":2234900,"MinAcceptableAmount":1292217,"TradingFee":196674,"TraderAddressStr":"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H","SubTraderAddressStr":"12Rx2NqWi5uEmMrT3fRVjhosBoGpjAQ9yxFmHckxZjyekU9YPdN622iVrwL3NwERvepotM6TDxPUo2SV4iDpW3NUukxeNCwJb2QTN9H","Type":205}}
2021/06/24 10:47:31 Finished TxsOut
```
## Example
[history.go](../../code/accounts/history/history.go)

```go
package main

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

// GetHistory retrieves the history of an account in a normal way.
// It is suitable for accounts that have a few transactions.
func GetHistory() {
	// For main-net
	client, err := incclient.NewIncClient("https://beta-fullnode.incognito.org/fullnode", incclient.MainNetETHHost, 1)
	if err != nil {
		log.Fatal(err)
	}

	tokenIDStr := common.PRVIDStr // input the tokenID in which you want to retrieve the history of.
	privateKey := "YOUR_PRIVATE_KEY" // input your private key here

	// get the history in a normal way.
	h, err := client.GetTxHistoryV1(privateKey, tokenIDStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("#TxIns: %v\n", len(h.TxInList))
	for _, txIn := range h.TxInList {
		log.Printf("%v\n", txIn.String())
	}
	log.Printf("\n#TxOuts: %v\n", len(h.TxOutList))
	for _, txOut := range h.TxOutList {
		log.Printf("%v\n", txOut.String())
	}
}

// GetHistoryFaster helps retrieve the history faster by running parallel workers.
func GetHistoryFaster() {
	// For main-net
	client, err := incclient.NewIncClient("https://beta-fullnode.incognito.org/fullnode", incclient.MainNetETHHost, 1)
	if err != nil {
		log.Fatal(err)
	}

	tokenIDStr := common.PRVIDStr // input the tokenID in which you want to retrieve the history of.
	privateKey := "YOUR_PRIVATE_KEY" // input your private key here

	numWorkers := 15
	p := incclient.NewTxHistoryProcessor(client, numWorkers)

	h, err := p.GetTokenHistory(privateKey, tokenIDStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("#TxIns: %v, #TxsOut: %v\n", len(h.TxInList), len(h.TxOutList))

	totalIn := uint64(0)
	log.Printf("TxsIn\n")
	for _, txIn := range h.TxInList {
		totalIn += txIn.Amount
		log.Printf("%v\n", txIn.String())
	}
	log.Printf("Finished TxsIn\n\n")

	totalOut := uint64(0)
	log.Printf("TxsOut\n")
	for _, txOut := range h.TxOutList {
		totalOut += txOut.Amount
		log.Printf("%v\n", txOut.String())
	}
	log.Printf("Finished TxsOut\n\n")
}

func main() {
	// comment one of these functions.
	GetHistory()
	GetHistoryFaster()
}
```
## Disclaimer
* This functionality is ONLY in its BETA phase, there might be some problems with its performance. If you encounter any problem, or find any issue, please report to the team.
* The history-retrieval process is time-consuming. Please take your time.
* Old full-nodes don't have some required RPCs, therefore, you MUST use this the designated full-node in this case.
* Your private key NEVER leave your local device. 

## TODOs

- [X] Get history of each tokenIDs
- [X] Get history of transactions V1
- [ ] Get history of transactions V2 (testing)
- [ ] Get the full history
- [ ] Get history for each feature
- [ ] Export to CSV files
---
Return to [the table of contents](../../../README.md).
