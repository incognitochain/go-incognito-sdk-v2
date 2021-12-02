---
Description: Tutorial on how to query the pDEX.
---

To query information from the pDEX, the SDK supports the following functions:

Name | Description | Note
-------------|-------------|-------------
GetPdexState| Get the state of pDEX | `beaconHeight = 0` will return the latest information
GetAllPdexPoolPairs | Get all pairs in pDEX | `beaconHeight = 0` will return the latest information
GetPdexPoolPair | Get all pools for a pair of tokens | `beaconHeight = 0` will return the latest information
GetPoolPairStateByID | Get the information of a pool by its ID | `beaconHeight = 0` will return the latest information
GetListNftIDs | Get all nftIDs on the pDEX | `beaconHeight = 0` will return the latest information
GetMyNFTs | Get all nftIDs owned by a private key |
GetOrderByID | Get the information of an order book given its ID |
GetPoolShareAmount | Get the share amount in a pool of an nftID |
CheckPrice | Calculate the receiving amount in a pool |
GetListStakingPoolShares | Get the list of tokens allowed to stake to the pDEX |
GetListStakingRewardTokens | Get the list of all available staking reward tokens |
GetEstimatedDEXStakingReward | Get the estimated pDEX staking rewards for an nftID with the given staking pool | `beaconHeight = 0` will return the latest information
CheckNFTMintingStatus | Get the status of an nft-minting transaction |
CheckTradeStatus | Get the status of a trade |
CheckDEXLiquidityContributionStatus | Get the status of a liquidity-contributing transaction |
CheckDEXLiquidityWithdrawalStatus | Get the status of a liquidity-withdrawal transaction |
CheckOrderAddingStatus | Get the status of an order-book adding transaction |
CheckOrderWithdrawalStatus | Get the status of an order-book withdrawal transaction |
CheckDEXStakingStatus | Get the status of a pDEX staking transaction |
CheckDEXUnStakingStatus | Get the status of a pDEX un-staking transaction |
CheckDEXStakingRewardWithdrawalStatus | Get the status of a pDEX staking-reward withdrawal transaction |
CheckDEXLPFeeWithdrawalStatus | Get the status of a pDEX LP fee withdrawal transaction |
CheckDEXProtocolFeeWithdrawalStatus | Get the status of a pDEX protocol fee withdrawal transaction |

In the next tutorial, we'll how to [mint a pDEX NFT](./nft.md).

---
Return to [the table of contents](../../../README.md).