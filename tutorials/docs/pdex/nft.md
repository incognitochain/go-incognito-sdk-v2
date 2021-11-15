---
Description: Introduction to NFTs in pDEX v3
---

## What is an NFT?
An NFT is a pDEX Access Token, introduced in the new pDEX v3 to enhance privacy when interacting with the pDEX.
Do you wonder if this NFT is the same as the Non-Fungible Token in the wild? The answer is no. In the context of pDEX v3, an NFT
mainly represents the ownership of a person on the pDEX (e.g, liquidity, order books, etc.). It is required that the user must supply
an NFT when performing one of these actions.

## Why NFTs?
To see what an NFT is used for, consider the following example. In the previous version of pDEX, a liquidity contribution looks
like the following,
```json
{
    "ContributorAddressStr": "12RpxrQRf2RrEJ4GzfPLfL9EgmtGUAbHLScz2fUGn8PCp7BLQMGWBPbenMd7U4wdrPV9SjJYueNNQ8iZuExb3j3Fhb389kE5kmozRfU",
    "TokenIDStr": "0000000000000000000000000000000000000000000000000000000000000004",
    "Amount": 1234000000000,
    "TxReqID": "ff8d575ae9fa60d745fd1f7f27680927480b230234804cb5d96111bc5b166b29"
}
```
where `ContributorAddressStr` is the payment address of the contributor. This address is required for two main purposes:
* to prove the ownership of the contribution;
* to receive back the token when withdrawing.

And if this person makes another contribution, for example,
```json
{
    "ContributorAddressStr": "12RpxrQRf2RrEJ4GzfPLfL9EgmtGUAbHLScz2fUGn8PCp7BLQMGWBPbenMd7U4wdrPV9SjJYueNNQ8iZuExb3j3Fhb389kE5kmozRfU",
    "TokenIDStr": "0000000000000000000000000000000000000000000000000000000000000004",
    "Amount": 2134000000000,
    "TxReqID": "dff913ce31224ad682b83d9804c037e4b2330b98b2c7f936d510a1734370a7ee"
}
```
the `ContributorAddressStr` stays the same. In this way, an observer can easily tell that these contributions are of the same person. As a consequence,
he can easily link any contributions (if they are related).

Let's see how NFTs can help in this scenario. In the new pDEX, whenever a user wants to perform a pDEX action (except for trading), he must send along an NFT, and a one-time address (OTA). If he wants to
make another contribution, he just needs to create another NFT and another OTA for it.

Now, let's see how contributions look like in the new pDEX with the following two examples.
```json
{
    "PoolPairID": "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115dc-03696365b2ff79bb9ef35bf43a74e655ffadae0fa139b8016148d7a036716c5c",
    "OtaReceiver": "15qTAohch6DZpk5Nzu7WBF17PLefNKngUmfx8HKkHKqnLJJsTPPopgBhCEQDPk82caHaPEb8WfvdP33WCwK3uG5WWHi1FykHf7VVK8LqpdJYy4wrTMz1o378c16fx5TVYyA2TX4LFJUGjCYa",
    "TokenID": "00000000000000000000000000000000000000000000000000000000000115dc",
    "Amount": 20000000000,
    "Amplifier": 20000,
    "TxReqID": "09bad3a7010bc1715ae5b6978fee0bd51f9b34d536cac36ec28ef0440c4ef59c",
    "NftID": "43b3d15a02d90e820f93d3de7fa682213836c46bd3d92f979848c872d5449bc9",
    "ShardID": 0
}
```

```json
{
    "PoolPairID": "0000000000000000000000000000000000000000000000000000000000000004-00000000000000000000000000000000000000000000000000000000000115dc-03696365b2ff79bb9ef35bf43a74e655ffadae0fa139b8016148d7a036716c5c",
    "OtaReceiver": "16CUGzz3KxMvHGEse34Jbj4tpAKGbjcac1XA62fTGac4sS4inqPZts14xUUMyeX4tcgiC5oUgq3oe3n9qMKzDdpXGk9AnUDd5bHbxEsJq7xapCNKhz8D6Qtc8RyEfmxtM3oPa4oUKWnx4dGE",
    "TokenID": "0000000000000000000000000000000000000000000000000000000000000004",
    "Amount": 2000000000,
    "Amplifier": 20000,
    "TxReqID": "c7a9d0617bfa835ee50c58a4556a781cfa59a7eef1f18e38afb206df8ab2c658",
    "NftID": "301c7ac6f9e6de3a67013a4fdb4f1fadf493f5571549e352fffcb5621881ca25",
    "ShardID": 0
}
```
The `ContributorAddressStr` in the previous example is now replaced with an `OtaReceiver` and an`NftID`. Here, the `NftID` represents the ownership while the `OtaReceiver` is used to receive back the contributed token when withdrawing.
As we can see, although the two contributions are of the same person, there's no link between them.

## How to Mint an NFT?
To mint a new NFT token, we need to create a transaction enclosed with the following metadata:
```go
type UserMintNftRequest struct {
   metadataCommon.MetadataBase
   otaReceiver string
   amount      uint64
}
```
in which
* `otaReceiver`: the OTA address for receiving the NFT.
* `amount`: the amount of burned PRV to mint the NFT. To prevent one from creating an infinite number of NFTs, an amount of PRV must be burned when minting a new NFT. This value is currently set to 1 PRV.

Luckily, you don't have to create the metadata yourself, the SDK provides a function called `CreateAndSendPdexv3UserMintNFTransaction`. With this function, all you need to do is to supply is your private key. And to check the status of an NFT-minting transaction, we use the function `CheckNFTMintingStatus` supplied with the created hash.
The status consists of the following information.
```go
// MintNFTStatus represents the status of a pDEX nft minting transaction.
type MintNFTStatus struct {
    // Status represents the status of the transaction, and should be understood as follows:
    // - 1: the request is accepted;
    // - 2: the request is rejected.
    Status int `json:"Status"`
    
    // BurntAmount is the amount of PRV that was burned to mint this NFT.
    BurntAmount uint64 `json:"BurntAmount"`
    
    // NftID is the ID of the minted NFT.
    NftID string `json:"NftID"`
}
```
where
See the following example ([mint.go](../../code/pdex/nft/mint.go)).

```go
package main

import (
   "encoding/json"
   "fmt"
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
   // burn some PRV to get your NFTID to use in pdex operations
   privateKey := ""

   txHash, err := client.CreateAndSendPdexv3UserMintNFTransaction(privateKey)
   if err != nil {
      log.Fatal(err)
   }

   fmt.Printf("Mint-NFT submitted in TX %v\n", txHash)

   // check the minting status
   time.Sleep(100 * time.Second)
   status, err := client.CheckNFTMintingStatus(txHash)
   if err != nil {
      log.Fatal(err)
   }

   jsb, err := json.MarshalIndent(status, "", "\t")
   if err != nil {
      log.Fatal(err)
   }
   fmt.Printf("status: %v\n", string(jsb))
}

```

We have seen how to mint an NFT and check its status on the new pDEX. Next, we see how to use this NFT to perform a new pDEX action, i.e, [pDEX contribution](./contribute.md).

---
Return to [the table of contents](../../../README.md).
