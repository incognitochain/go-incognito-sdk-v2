[![Go Report Card](https://goreportcard.com/badge/github.com/incognitochain/go-incognito-sdk-v2)](https://goreportcard.com/report/github.com/incognitochain/go-incognito-sdk-v2) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/LICENSE)

# go-incognito

## Introduction

A Golang SDK for interacting with the Incognito network.

## Installation

```go
go get github.com/incognitochain/go -incognito-sdk-v2
```

## Dependencies

See [go.mod](./go.mod)

## Import

```go
import (
github.com/incognitochain/go -incognito-sdk-v2
)
```

## Tutorials

Following is a series of examples to help you get familiar with the Incognito network. The list does not cover all the
capabilities of the SDK, we will try to cover them as much as possible.

* [Introduction](tutorials/docs/intro/intro.md)
* [Client](tutorials/docs/client)
    * [Setting up the Client](tutorials/docs/client/client.md)
* [Accounts](tutorials/docs/accounts)
    * [Keys](tutorials/docs/accounts/keys.md)
    * [Creating Accounts with HD Wallets](tutorials/docs/accounts/hdwallet_create.md)
    * [Importing Accounts with Mnemonic Strings](tutorials/docs/accounts/hdwallet_import.md)
    * [UTXOs](tutorials/docs/accounts/utxo.md)
        * [Retrieving Output Coins V1](tutorials/docs/accounts/utxo_retrieve.md)
        * [Key Submission](tutorials/docs/accounts/submit_key.md)
        * [UTXO Cache](tutorials/docs/accounts/utxo_cache.md)
        * [Consolidating](tutorials/docs/accounts/consolidate.md)
    * [Account Balances](tutorials/docs/accounts/balances.md)
    * [Account History](tutorials/docs/accounts/tx_history.md)
* [Transactions](tutorials/docs/transactions)
    * [Transaction Parameters](tutorials/docs/transactions/params.md)
    * [Transferring PRV](tutorials/docs/transactions/raw_tx.md)
    * [Transferring Token](tutorials/docs/transactions/raw_tx_token.md)
    * [Initializing Custom Tokens](tutorials/docs/transactions/init_token.md)
    * [Converting UTXOs](tutorials/docs/transactions/convert.md)
* [pDEX](tutorials/docs/pdex/intro.md)
    * [Querying the pDEX](tutorials/docs/pdex/query.md)
    * [NFT Minting](tutorials/docs/pdex/nft.md)
    * AMM Liquidity
        * [Contribution](tutorials/docs/pdex/contribute.md)
        * [Withdrawal](tutorials/docs/pdex/withdraw.md)
        * [LP Fee Withdrawal](tutorials/docs/pdex/lp_fee_withdraw.md)
    * Order Books
        * [Adding an Order Book](tutorials/docs/pdex/ob_add.md)
        * [Canceling an Order Book](tutorials/docs/pdex/ob_cancel.md)
    * AMM One-sided Liquidity (Staking)
        * [Staking](tutorials/docs/pdex/stake.md)
        * [UnStaking](tutorials/docs/pdex/unstake.md)
        * [Withdrawing the Staking Reward](tutorials/docs/pdex/staking_reward_withdraw.md)
    * [Trading](tutorials/docs/pdex/trade.md)
* [Staking](tutorials/docs/staking)
    * [Creating Staking Transactions](tutorials/docs/staking/stake.md)
    * [Creating UnStaking Transactions](tutorials/docs/staking/unstake.md)
    * [Creating Reward Withdrawal Transactions](tutorials/docs/staking/withdraw_reward.md)
    * [Node Monitoring](tutorials/docs/staking/node.md)
* Bridges
    * [EVM/PRV Shielding/UnShielding](tutorials/docs/bridge/evm/bridge.md)
        * [Shielding Transactions](tutorials/docs/bridge/evm/shield.md)
        * [Checking Shielding Status](tutorials/docs/bridge/evm/shield_status.md)
        * [Un-Shielding Transactions](tutorials/docs/bridge/evm/unshield.md)
        * [Un-Shielding PRV](tutorials/docs/bridge/evm/unshield_prv.md)
        * [Shielding PRV](tutorials/docs/bridge/evm/shield_prv.md)
    * [Portal Shielding/UnShielding](tutorials/docs/bridge/portal/portal.md)
        * [Terms and Functions](tutorials/docs/bridge/portal/terms_functions.md)
        * [Shielding](tutorials/docs/bridge/portal/shield.md)
        * [UnShielding](tutorials/docs/bridge/portal/unshield.md)
* [Calling RPCs](tutorials/docs/rpc/rpc.md)
