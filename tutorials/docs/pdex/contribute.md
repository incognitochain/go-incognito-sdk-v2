---
Description: Tutorial on how to add pairs to the pDEX.
---
# Before Going Further
Please read through the tutorials on [key submission](../accounts/submit_key.md) and [UTXO cache](../accounts/utxo_cache.md) for proper
balance and UTXO retrieval. Skip these parts if you're familiar with these notions.

# pDEX Contribution
Liquidity providers play an essential role in pDEX. They provide liquidity to various pools on pDEX and earn trading fees. The current pDEX consists of several pairs of tokens that help accelerate trading activities. The more liquidity in the pDEX, the better experience the trading process gets. For a pair with high liquidity, the slippage rate will be small. On the other hand, trading with low-liquidity pair will result in a high slippage rate.

In this tutorial, we will see how we can provide liquidity for a pair in the pDEX. Please see this [post](https://github.com/incognitochain/incognito-chain/blob/production/specs/pdex.md) to understand how the pDEX works. There are 3 ways a user can provide liquidity for the pDEX:
* 2-sided liquidity adding;
* 1-sided liquidity contribution; and
* order-book placing.

This tutorial only focuses on the first method. A two-sided contribution contributes both tokens to the liquidity pool. A liquidity provider must create two separate metadata transactions (one for each token) to burn the corresponding amount of the token.
This process is exactly like that in the previous pDEX version, except that the metadata is now slightly different. Instead of using the payment address as the identification, the LP now uses [a so-called NFT ID](./nft.md) for the contribution. In this way, his contributions will no longer be linkable if they are created with different nftIDs.

Here is the description of the metadata.
```go
type AddLiquidityRequest struct {
    // poolPairID is the ID of the target pool in which the LP wants to add liquidity to. 
    // If this is the first contribution (i.e, pool-initialization), the poolPairID must be left empty.
    poolPairID  string
    
    // pairHash is a string for matching the two contributing transactions. It can be anything as long as it is the same in
    // both contributing transaction.
    pairHash    string
    
    // otaReceiver is a one-time address for receiving back the token in case of being refunded.
    otaReceiver string
    
    // tokenID is the ID of the contributing token.
    tokenID     string
    
    // nftID is the ID of the NFT associated with this contribution. This value must be the same in both contributing transactions.
    nftID       string
    
    // tokenAmount is the contributing amount of this token.
    tokenAmount uint64
    
    // amplifier is the amplifier of the pool. In the case of contributing to an existing pool, this value must match that of the existing pool. 
    // The detail of this param can be found in Uniswap's White-paper (https://uniswap.org/whitepaper-v3.pdf).
    amplifier   uint
    
    metadataCommon.MetadataBase
}
```

Now, let's try to create a new pDEX pool using the SDK.
## Create a new pDEX Pool
### Prepare our inputs
As usual, we need to specify our private key.
```go
privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
```
Then, we must specify the contributing information.
```go
poolPairID := ""       // for pool-initializing, leave it empty. Otherwise, input the poolPairID of the existing pool
pairHash := "JUSTARANDOMSTRING" // a string to match the two transactions of the contribution
firstToken := common.PRVIDStr
secondToken := "00000000000000000000000000000000000000000000000000000000000115d7"
firstAmount := uint64(10000)
secondAmount := uint64(10000)
nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
amplifier := uint64(15000)
```
So we are contributing a pair of tokens (PRV and `00000000000000000000000000000000000000000000000000000000000115d7`) with the amount of `10000` for each, using the NFT ID `54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22`.
Since we are creating a new pool, the `poolPairID` is left empty.

### Create contributing transactions
To create the contributing transactions, we use the method `CreateAndSendPdexv3ContributeTransaction` provided by the SDK supplied with the above-specified paramters.
```go
firstTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, poolPairID, pairHash, firstToken, nftIDStr, firstAmount, amplifier)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("firstTx: %v\n", firstTx)

// wait for the first transaction to be confirmed, so the nftID has been re-minted to proceed.
time.Sleep(60 * time.Second)
secondTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, poolPairID, pairHash, secondToken, nftIDStr, secondAmount, amplifier)
if err != nil {
    log.Fatal(err)
}
```

Note that both of the transactions use the same nftID. Therefore, we must wait for the previous transaction to finish before making the second one.

### Check contribution status
Checking the contributing status is done via the function `CheckDEXLiquidityContributionStatus`, supplied with the previously-created transaction hash.
An example of the status is as follows:
```go
{
    "Status": 2,
    "Token0ID": "",
    "Token0ContributedAmount": 0,
    "Token0ReturnedAmount": 0,
    "Token1ID": "",
    "Token1ContributedAmount": 0,
    "Token1ReturnedAmount": 0,
    "PoolPairID": ""
}
```
Here, there's not much we can get from this except for `Status = 2` meaning that the contributing request has been fully accepted and this contribution
is the first contribution of the pool (see more about this [here](./query.md)).

## Add Liquidity to an Existing Pool
This case is somewhat similar to the case of initializing a pool, except:
* Now we need to specify which pool we want to contribute liquidity to.
* We must specify the amplifier exactly as that of the target pool.

Suppose that the previous contribution created a new pool with the `poolPairID` `0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2`.
We can add liquidity to this pool as follows.
### Prepare our inputs
```go
poolPairID := "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"       // for pool-initializing, leave it empty. Otherwise, input the poolPairID of the existing pool
pairHash := "JUSTANOTHERARANDOMSTRING" // a string to match the two transactions of the contribution
firstToken := common.PRVIDStr
secondToken := "00000000000000000000000000000000000000000000000000000000000115d7"
firstAmount := uint64(10000)
secondAmount := uint64(10000)
nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
amplifier := uint64(15000)
```
### Create contributing transactions
This is done exactly like before. And we can see the status looks like the following:
```go
{
    "Status": 4,
    "Token0ID": "0000000000000000000000000000000000000000000000000000000000000004",
    "Token0ContributedAmount": 10000,
    "Token0ReturnedAmount": 0,
    "Token1ID": "00000000000000000000000000000000000000000000000000000000000115d7",
    "Token1ContributedAmount": 10000,
    "Token1ReturnedAmount": 0,
    "PoolPairID": "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115d7-0868e6a074566d77c2ebdce49949352efbe69b0eda7da839bfc8985e7ed300f2"
}
```
`Status = 4` indicates that the contribution request has been accepted with associated information. In many cases, one of the returned amounts will be non-zero.
This is because the beacon will calculate the contributing amounts based on the current pool rate, and therefore there might be some leftover for one of the tokens.


## NOTE
In this tutorial, we use the same nftID for both of the contributions. Although this is allowed by the protocol,
it is not RECOMMENDED because it results in the possibility of linking these two contributions together. In practice, the LP
might want to generate a bunch of NFTs for his account, and for each contribution, he might use a different NFT. However, please take minting fees into account since there is an amount of PRV required for minting a new NFT.

## Example
[contribute](../../code/pdex/pdex_contribute/contribute.go)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	// replace with your network's data
	privateKey := "112t8rneWAhErTC8YUFTnfcKHvB1x6uAVdehy1S8GP2psgqDxK3RHouUcd69fz88oAL9XuMyQ8mBY5FmmGJdcyrpwXjWBXRpoWwgJXjsxi4j"
	poolPairID := "" // for pool-initializing, leave it empty. Otherwise, input the poolPairID of the existing pool
	pairHash := "JUSTARANDOMSTRING" // a string to match the two transactions of the contribution
	firstToken := common.PRVIDStr
	secondToken := "00000000000000000000000000000000000000000000000000000000000115d7"
	firstAmount := uint64(3000)
	secondAmount := uint64(3000)
	nftIDStr := "54d488dae373d2dc4c7df4d653037c8d80087800cade4e961efb857c68b91a22"
	amplifier := uint64(15000)

	firstTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, poolPairID, pairHash, firstToken, nftIDStr, firstAmount, amplifier)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("firstTx: %v\n", firstTx)

	// wait for the first transaction to be confirmed, so the nftID has been re-minted to proceed.
	//time.Sleep(60 * time.Second)
	secondTx, err := client.CreateAndSendPdexv3ContributeTransaction(privateKey, poolPairID, pairHash, secondToken, nftIDStr, secondAmount, amplifier)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("secondTx: %v\n", secondTx)

	// check the minting status
	time.Sleep(100 * time.Second)
	status, err := client.CheckDEXLiquidityContributionStatus(firstTx)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("firstTxstatus: %v\n", string(jsb))

	status, err = client.CheckDEXLiquidityContributionStatus(secondTx)
	if err != nil {
		log.Fatal(err)
	}
	jsb, err = json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("secondTxstatus: %v\n", string(jsb))
}

```

We have seen how to add a pair to the pDEX. Let's see how we can [withdraw liquidity](./withdrawal.md) from the pDEX.

---
Return to [the table of contents](../../../README.md).
