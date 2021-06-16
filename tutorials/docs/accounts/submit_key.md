---
description: Tutorial on how to have a full-node cache our output coins.
---
# Full-nodes' Cache
The benefits of increased privacy are not without costs. Retrieving output coins is one of the most prominent. In Privacy V1, a user's output coins all have the same public key, which is the user's public key. In this approach, a full node's database can effortlessly aggregate these output coins. When a user requests all of his or her output coins, the full-node only needs to look for those that have the same public key.
Because we used one-time addresses to enhance receiver anonymity in Privacy V2, this no longer holds true. Different public keys will now be assigned to each of a user's output coins. As a result, the full-node is unable to determine if two output coins belong to the same or separate users. As a consequence, the full-node will be unable to respond to a user's request for his/her output coins.

One approach is for full-nodes to cache a user's output coins once the user submits his or her privateOTA key, allowing for quicker retrieval. The full-node may then scan each output coin and identify if it belongs to this user or not. However, determining whether a coin belongs to a user is costly. As a result, the full-node is designed with two modes of operation.
* **Default mode**: after the user's keys are submitted, the full-node only caches the user's output coins. As a result, any output coins V2 received prior to the submission of the privateOTA key will not be stored in the full-node's cache. If a user wishes to use this mode, he/she must submit his/her key before the first output coin arrives.
* **Enhanced mode**: the full-node will re-scan the database and add all output coins of a user into its cache.

Next, we describe in detail each of these modes of operation and give a brief comparison between the two.

## Default mode
In this mode of operation, the full-node only starts caching a user’s output coin after the user submits the privateOTA key.
Consider the following example, suppose that Bob sends an UTXO (say UTXO A) to Alice at the block height 10, and at this time, Alice has not submitted her privateOTA key to the full node. As a result, the full node has no idea that these coins belong to Alice. 10 blocks later, i.e, block 20, Alice tells the full node to cache her coins. And right after that, David sends another UTXO (say B) to Alice. Now, if Alice queries her output coins from the full-node, the returned result only acknowledges the UTXO B (since the UTXO A has been “lost” in the view of the full-node).

Here is the summary of this mode.
* The full node only does its work after Alice submits her key.
* Every time a new block arrives, the full node will try to check if any output coins (of this block) belong to Alice. If yes, it caches these coins for Alice. This caching process is called **“passive caching”** and it does not take much effort of the full node.
* All output coins arriving before the key is submitted will be “lost”.
* To avoid losing UTXOs, users are **RECOMMENDED** to submit their key before any UTXO arrives.

### Submit keys via RPC
```json
{
        "Jsonrpc": "1.0",
        "Method": "submitkey",
        "Params": [
                "14yBChbLDg42noBQHonDR5mj3FMD9CPCCNfPoa68jeE8bE2LsyfCKcNkgupEsm6pW4BZFnDHmay9XjDGE1iTaTEcEpN7UUaPoU344g2"
        ],
        "Id": 1
}
```
* **Method**: `submitkey`.
* **Params**: the only parameter is a base58-encoded privateOTA key.
* **Errors**: the following are some of the error messages that might be returned by this query.

**Error Message** | Description
-----------|-----------
OTAKey has been submitted and status = 1(2,3) | The privateOTA key has been submitted before and has status. If status = 1, the indexing process is in progress, if status = 2, the regular indexing process has finished, if status = 3, the enhanced indexing process has finished.
OTA key submission not supported by this node configuration | The current node does not have the cache layer.

### Submit keys via go-sdk
```go
// for regular cache
err = client.SubmitKey(privateOTAKey)
if err != nil {
log.Fatal(err)
}
```

## Enhanced mode
As we can see, the previous mode does not provide much flexibility and creates a poor user experience. That is, if a user fails to submit the key, he/she will be unable to retrieve the balance through the full-node. As a result, we introduce the enhanced mode to assist the full-node owner or anyone else authorized in retrieving their total balance.
Consider the previous Alice example. In this case, she will be able to retrieve both UTXOs A and B if she utilizes the enhanced mode.

Here is the summary of this mode.
* After a key is submitted, the full-node will try to check and cache all of the output coins of this key from the beginning.
* Because the full-node has to re-scan from the beginning, this mode is very expensive and takes quite a long time. During the time the full-node is re-indexing output coins, any operations to check balance, retrieve output coins will be stalled. The longer the blockchain, the more expensive this process is. That’s why it is called **“active caching”**. Therefore, we recommend you run this mode for the full-node **ONLY**.
* Authorization is required since the fullnode is easily DDoS'ed. As a result, only a limited number of users should be allowed to use this mode.
* We estimate that with more than 300 requests at the same time, the full-node will be out of order. Therefore, **DO NOT** share the access token if it is not necessary. Additionally, connections with **HTTPS** are recommended to make sure the access token is not stolen.
* Also, if the authorization fails, the basic mode will be employed.

### Submit authorized keys via RPC
```json
{
        "Jsonrpc": "1.0",
        "Method": "authorizedsubmitkey",
        "Params": [
                "14y8spKEPrqLndpwjrQsfdX4y8VWrSwAhLPmKF2GpLocEh3pvuaDoug5T7gEgifV8amh9RBs1MKa4fSvXLwL4iAovHQPLwbGcEjJ3A2",
                "0c3d46946bbf9339c8213dd7f6c640ed643003bdc056a5b68e7e80f5ef5aa0dd",
                0,
                false
        ],
        "Id": 1
}
```
* **Method**: `authorizedsubmitkey`.
* **Params**: There are 4 parameters needed
    * The first parameter is  a base58-encoded privateOTA key.
    * The second is an access token. This access token is generated by the full-node’s owner, and is compulsory in this RPC.
    * The third one is the block height at which the full-node will re-scan from. If this parameter is set to 0, the full-node will rescan from the beginning.
    * The final parameter is a boolean indicating the flag reset, and it is optional. In case the privateOTA key has been indexed before and this flag is set to true, the full-node will re-index all output coins for this key.
* **Errors**: the following are some error messages that might be returned by this query.

**Error Message** | Description
-----------|-----------
OTAKey has been submitted and status = 1(2,3) | The privateOTA key has been submitted before and has status. If status = 1, the indexing process is in progress, if status = 2, the regular indexing process has finished, if status = 3, the enhanced indexing process has finished.
OTA key submission not supported by this node configuration | The current node does not have the cache layer.
enhanced caching not supported by this node configuration | The current node only operates with the basic mode.
the current authorized queue is full, please check back later | The cache layer only supports a limited number of users at the same time to reduce the risk of being DDoS’ed.
fromHeight is larger than the current shard height  | The third parameter is larger than the current block height on the blockchain.
  
### Submit authorized keys via go-sdk
```go
// for enhanced cache
accessToken := "0c3d46946bbf9339c8213dd7f6c640ed643003bdc056a5b68e7e80f5ef5aa0dd"
fromHeight := uint64(0)
isReset := true
err = client.AuthorizedSubmitKey(privateOTAKey, accessToken, fromHeight, isReset)
if err != nil {
log.Fatal(err)
}
```

## Comparison
Here is a comparison of the two modes.

**Property**  | **Passive Caching** | **Active Caching**
--------|------------|-----------
When to cache coins | After the key is submitted | From the beginning
UX | Bad | Better
Full-node's load | Low | Very high
#Users | Unlimited | Limited, only a few
Access token | Not required | Required
Node-friendly | Validators/Shard/Full-nodes | Full-nodes

## Retrieve output coins from the full-node's cache
* Via RPC: `listoutputcoinsfromcache`.
* Via go-sdk: `GetOutputCoinsV2()` or `GetListOutputCoinsByRPCV2()`.

## Example
[key_submit.go](../../code/accounts/key_submit/key_submit.go)
```go
package main

import (
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateOTAKey := "14yBChbLDg42noBQHonDR5mj3FMD9CPCCNfPoa68jeE8bE2LsyfCKcNkgupEsm6pW4BZFnDHmay9XjDGE1iTaTEcEpN7UUaPoU344g2"

	// for regular cache
	err = client.SubmitKey(privateOTAKey)
	if err != nil {
		log.Fatal(err)
	}

	// for enhanced cache
	accessToken := "0c3d46946bbf9339c8213dd7f6c640ed643003bdc056a5b68e7e80f5ef5aa0dd"
	fromHeight := uint64(0)
	isReset := true
	err = client.AuthorizedSubmitKey(privateOTAKey, accessToken, fromHeight, isReset)
	if err != nil {
		log.Fatal(err)
	}
}
```

Next, let's see how we can get our [balances](../accounts/balances.md).

--- 
Return to [the table of contents](../../../README.md).