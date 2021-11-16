---
Description: Tutorial on how to query the pDEX.
---

The Incognito client supports the following functions:

Name | Description | Status
-------------|-------------|-------------
GetPDEState| Get the state of pDEX | Finished
GetAllPDEPoolPairs | Get all pairs in pDEX | Finished
GetPDEPoolPair | Get the information of a pair in pDEX | Finished
CheckXPrice | Get the price of a trade based in the current state of pDEX| Finished
CheckTradeStatus | Check the status of a trading request | Finished
GetShareAmount | Get the share amount of a user for a pool on the pDEX | Finished
CreatePDEContributeTransaction | Create transactions contributing tokens to pDEX | Finished
CreatePDEWithdrawalTransaction | Create transactions withdrawing pDEX contribution | Finished
CreatePDETradeTransaction | Create trading transactions | Finished

# Querying pDEX

In this tutorial, we will learn how to perform basic pDEX querying operations with the Incognito SDK.

## Get the pDEX State

On input the beacon height, the function `GetPDEState` returns the state of pDEX at this beacon height. If the beacon
height is set to `0`, it returns the latest state.

```go
pdeState, err := client.GetPDEState(0)
if err != nil {
log.Fatal(err)
}

fmt.Printf("pdeState: \n%v\n", pdeState)
```

A state of the pDEX consists of

* all pairs in the pDEX;
* all shares in the pDEX;
* total trading fees at the beacon height; and
* the list of pending contributions.

## Get all pDEX Pairs

To retrieve all pairs in the pDEX, we call the function `GetAllPDEPoolPairs` of the client. Similar to above, it receives as input a
beacon height.

```go
allPairs, err := client.GetAllPDEPoolPairs(0)
if err != nil {
log.Fatal(err)
}

fmt.Printf("pdeState: \n%v\n", allPairs)
```

For each pair, it returns

* IDs of the two tokens;
* amount of each token.

## Get Pair Information

In case we want to retrieve the information of a specific pair, use `GetPDEPoolPair`.

```go
tokenID1 := common.PRVIDStr
tokenID2 := "0000000000000000000000000000000000000000000000000000000000000100"
pair, err := client.GetPDEPoolPair(0, tokenID1, tokenID2)
if err != nil {
log.Fatal(err)
}
fmt.Printf("pair: %v\n", pair)
```

The first parameter is the beacon height, just like previous; followed by the 2 tokenIDs (any order is fine).

## Check pDEX Price

Checking prices is one of the most important functions of pDEX. What we need is to specify which token we want to sell,
which token we want to buy, and the selling amount. Notice that the selling amount is required to calculate the exact
amount since pDEX used the AMM algorithm.

```go
tokenToSell := common.PRVIDStr
tokenToBuy := "0000000000000000000000000000000000000000000000000000000000000100"
sellAmount := uint64(1000000000)
expectedAmount, err := client.CheckXPrice(tokenToSell, tokenToBuy, sellAmount)
if err != nil {
log.Fatal(err)
}
fmt.Printf("Expected amount: %v\n", expectedAmount)
```

Here, we are selling 1 PRV to buy `0000000000000000000000000000000000000000000000000000000000000100`.

## Query the Status

### Mint NFTs
```go
status, err := client.CheckNFTMintingStatus(txHash string)
if err != nil {
    log.Fatal(err)
}
```
The status will consist of the following information.
```go
// MintNFTStatus represents the status of a pDEX nft minting transaction.
type MintNFTStatus struct {
    // Status represents the status of the transaction, and should be understood as follows:
    //	- 1: the request is accepted;
    //	- 2: the request is rejected.
    Status int `json:"Status"`
    
    // BurntAmount is the amount of PRV that was burned to mint this NFT.
    BurntAmount uint64 `json:"BurntAmount"`
    
    // NftID is the ID of the minted NFT.
    NftID string `json:"NftID"`
}
```

### Trade
```go
status, err := client.CheckTradeStatus(txHash string)
if err != nil {
    log.Fatal(err)
}
```
The status will consist of the following information.
```go
// DEXTradeStatus represents the status of a pDEX v3 trade.
type DEXTradeStatus struct {
    // Status represents the status of the trade, and should be understood as follows:
    // 	- 0: the trade request is refunded;
    //	- 1: the trade request is accepted.
    Status int `json:"Status"`
    
    // BuyAmount is the receiving amount of the trade (in case of failure, it equals to 0).
    BuyAmount uint64 `json:"BuyAmount"`
    
    // TokenToBuy is the buying tokenId.
    TokenToBuy string `json:"TokenToBuy"`
}
```