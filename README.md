# go-incognito
## Introduction
A Golang SDK for interacting with the Incognito network.

[![Go Report Card](https://goreportcard.com/badge/github.com/incognitochain/go-incognito-sdk-v2)](https://goreportcard.com/report/github.com/incognitochain/go-incognito-sdk-v2)
## Installation
```go
go get github.com/incognitochain/go-incognito-sdk-v2
```

## Dependencies
* `go-ethereum`: version 1.9 

## Import
```go
import (
    github.com/incognitochain/go-incognito-sdk-v2
)
```

## Tutorials
Following is a series of examples to help you get familiar with the Incognito network. The list does not cover all the capabilities of the SDK, we will try to cover them as much as possible. 

* [Introduction](tutorials/docs/intro/intro.md)
* [Client](tutorials/docs/client)
    * [Setting up the Client](tutorials/docs/client/client.md)
* [Accounts](tutorials/docs/accounts)
    * [Keys](tutorials/docs/accounts/keys.md)
    * [Creating Accounts with HD Wallets](tutorials/docs/accounts/hdwallet_create.md)
    * [Importing Accounts with Mnemonic Strings](tutorials/docs/accounts/hdwallet_import.md)
    * [UTXOs](tutorials/docs/accounts/utxo.md)
      * [Retrieving output coins V1](tutorials/docs/accounts/utxo_retrieve.md)
      * [Key Submitting](tutorials/docs/accounts/submit_key.md)
    * [Account Balances](tutorials/docs/accounts/balances.md)
    * [Account History](tutorials/docs/accounts/tx_history.md)
* [Transactions](tutorials/docs/transactions)
    * [Transaction Parameters](tutorials/docs/transactions/params.md)
    * [Transferring PRV](tutorials/docs/transactions/raw_tx.md)
    * [Transferring Token](tutorials/docs/transactions/raw_tx_token.md)
    * [Initializing Custom Tokens](tutorials/docs/transactions/init_token.md)
    * [Converting UTXOs](tutorials/docs/transactions/convert.md)
* [pDEX](tutorials/docs/pdex)
    * [Querying pDEX](tutorials/docs/pdex/query.md)
    * [Creating pDEX Contribution Transactions](tutorials/docs/pdex/contribute.md)
    * [Creating pDEX Withdrawal Transactions](tutorials/docs/pdex/withdrawal.md)
    * [Creating pDEX Trading Transactions](tutorials/docs/pdex/trade.md)
* [Staking](tutorials/docs/staking)
    * [Creating Staking Transactions](tutorials/docs/staking/stake.md)
    * [Creating UnStaking Transactions](tutorials/docs/staking/unstake.md)
    * [Creating Reward Withdrawal Transactions](tutorials/docs/staking/withdraw_reward.md)
    * [Node Monitoring](tutorials/docs/staking/node.md)
* [Shielding/UnShielding](tutorials/docs/bridge/bridge.md)
    * [Creating Shielding Transactions](tutorials/docs/bridge/shield.md)
    * [Creating Un-Shielding Transactions](tutorials/docs/bridge/unshield.md)
* [Calling RPCs](tutorials/docs/rpc/rpc.md)
  
## TODOs

- [X] UTXOs
- [X] Balance
- [X] PRV + Token transactions
- [X] pDEX
- [X] Stake
- [X] Shield ETH, ERC20
- [X] UnShield ETH, ERC20
- [X] HD Wallet
- [ ] History
- [ ] ...
