# Develop Your Own Applications with Go-Incognito SDK

This project serves as a general sdk for anyone wanting to develop Incognito applications using the Go programming language. You will be able to perform general operations such as querying, creating and submitting transactions, etc,.

It also comes with a series of examples that we believe anyone will encounter when they develop an application on the Incognito blockchain. It will help you walk through most things that you should be aware of in order to create an effect Incognito-related application.

## Read the code
The main focus of this project lies in the [**incclient**](../../../incclient) package. It provides access to almost all functions needed to create transactions,
  become a node validator, retrieve information from full-nodes, shield or un-shield access, etc. One can find the instruction of how to use this go-sdk in the [**tutorials**](../../../tutorials) package.
Besides, if you want to understand the code, try other packages:

* [**coin**](../../../coin) implements input and output coins, which are basic elements of a transactions. Incognito is currently supporting two version of coins.
* [**common**](../../../common) handles all common functions and variables used in this project.
* [**crypto**](../../../crypto) implements the `Edwards25519` elliptic curve as well as the [Pedersen commitment scheme](https://en.wikipedia.org/wiki/Commitment_scheme).
* [**eth_bridge**](../../../eth_bridge) consists of smart contracts, mainly responsible for shielding or un-shielding ETH, ERC20.
* [**key**](../../../key) handles all key-related things including private keys, payment addresses, read-only keys, OTA private keys or mining keys.
* [**metadata**](../../../metadata) contains additional information enclosed with a transaction to specify their special functions. They include, but not limited to
  * **pDEX**
  * **Staking**
  * **Bridge**
  * **Portal**
* [**privacy**](../../../privacy) is the most important layer of the Incognito network. It helps protect user anonymity, transaction confidentiality. The `privacy` layer features
  * **Ring Signatures.** For privacy, Incognito implements MLSAG signatures with stealth addresses.
  * **Zero-Knowledge Proofs.** We use zero-knowledge proofs (ZKPs) to hide the amount of a transaction.
  * **Confidential Asset.** ZKPs hide the amount of the transaction, but they don't hide the type of transferred asset. Confidential Asset solves that.
* [**rpchandler**](../../../rpchandler) is responsible for interacting with full-nodes to retrieve information.
* [**transaction**](../../../transaction) helps create various types of transactions ranging from regular transferring transactions, to pDEX trading, or staking transactions.
* [**wallet**](../../../wallet) helps generate new key-sets, restore keys, etc.

## Contribution
Any issues are welcomed to be submitted at [issues](https://github.com/incognitochain/go-incognito-sdk-v2/issues). And feel free to make a [pull request](https://github.com/incognitochain/go-incognito-sdk-v2/pulls) if you observe things that should be improved.
