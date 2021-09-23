package incclient

import (
	"fmt"
	"sort"
	// "strings"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// Share represents a pDEX contribution share.
type Share struct {
	TokenID1Str string
	TokenID2Str string
	ShareAmount uint64
}

// GetPDEState retrieves the state of pDEX at the provided beacon height.
// If the beacon height is set to 0, it returns the latest pDEX state.
func (client *IncClient) GetPDEState(beaconHeight uint64) (*jsonresult.CurrentPdexState, error) {
	if beaconHeight == 0 {
		bestBlocks, err := client.GetBestBlock()
		if err != nil {
			return nil, fmt.Errorf("cannot get best blocks: %v", err)
		}
		beaconHeight = bestBlocks[-1]
	}

	responseInBytes, err := client.rpcServer.GetPDEState(beaconHeight)
	if err != nil {
		return nil, err
	}

	var pdeState jsonresult.CurrentPdexState
	err = rpchandler.ParseResponse(responseInBytes, &pdeState)
	if err != nil {
		return nil, err
	}

	return &pdeState, nil
}

// GetAllPDEPoolPairs retrieves all pools in pDEX at the provided beacon height.
// If the beacon height is set to 0, it returns the latest pDEX pool pairs.
func (client *IncClient) GetAllPDEPoolPairs(beaconHeight uint64) (map[string]*jsonresult.Pdexv3PoolPairState, error) {
	pdeState, err := client.GetPDEState(beaconHeight)
	if err != nil {
		return nil, err
	}

	return pdeState.PoolPairs, nil
}

// GetPDEPoolPair retrieves the pDEX pool information for pair tokenID1-tokenID2 at the provided beacon height.
// If the beacon height is set to 0, it returns the latest information.
func (client *IncClient) GetPDEPoolPair(beaconHeight uint64, tokenID1, tokenID2 string) (*jsonresult.Pdexv3PoolPairState, error) {
	if beaconHeight == 0 {
		bestBlocks, err := client.GetBestBlock()
		if err != nil {
			return nil, fmt.Errorf("cannot get best blocks: %v", err)
		}
		beaconHeight = bestBlocks[-1]
	}

	allPoolPairs, err := client.GetAllPDEPoolPairs(beaconHeight)
	if err != nil {
		return nil, err
	}

	keyPool := jsonresult.BuildPDEPoolForPairKey(beaconHeight, tokenID1, tokenID2)
	if poolPair, ok := allPoolPairs[string(keyPool)]; ok {
		return poolPair, nil
	}

	return nil, fmt.Errorf("cannot found pool pair for tokenID %v and %v", tokenID1, tokenID2)
}

// CheckPrice gets the remote server to check price for trading things.
func (client *IncClient) CheckPrice(tokenToSell, TokenToBuy string, sellAmount uint64) (uint64, error) {
	responseInBytes, err := client.rpcServer.ConvertPDEPrice(tokenToSell, TokenToBuy, sellAmount)
	if err != nil {
		return 0, err
	}

	var convertedPrice []*rpc.ConvertedPrice
	err = rpchandler.ParseResponse(responseInBytes, &convertedPrice)
	if err != nil {
		return 0, err
	}

	if len(convertedPrice) == 0 {
		return 0, fmt.Errorf("no convertedPrice found for %v", tokenToSell)
	}

	return convertedPrice[0].Price, nil
}

// CheckXPrice gets the remote server to check cross price for trading things (for cross-pool tokens).
func (client *IncClient) CheckXPrice(tokenToSell, TokenToBuy string, sellAmount uint64) (uint64, error) {
	if tokenToSell == common.PRVIDStr || TokenToBuy == common.PRVIDStr {
		return client.CheckPrice(tokenToSell, TokenToBuy, sellAmount)
	}

	expectedPRV, err := client.CheckPrice(tokenToSell, common.PRVIDStr, sellAmount)
	if err != nil {
		return 0, err
	}

	return client.CheckPrice(common.PRVIDStr, TokenToBuy, expectedPRV)
}

// GetShareAmount retrieves the share amount of a payment address in pDEX pool of tokenID1 and tokenID2.
// func (client *IncClient) GetShareAmount(beaconHeight uint64, tokenID1, tokenID2, paymentAddress string) (uint64, error) {
// 	if beaconHeight == 0 {
// 		bestBlocks, err := client.GetBestBlock()
// 		if err != nil {
// 			return 0, fmt.Errorf("cannot get best blocks: %v", err)
// 		}
// 		beaconHeight = bestBlocks[-1]
// 	}

// 	pdeState, err := client.GetPDEState(beaconHeight)
// 	if err != nil {
// 		return 0, err
// 	}

// 	allShares := pdeState.Shares
// 	shareKey, err := BuildPDEShareKey(beaconHeight, tokenID1, tokenID2, paymentAddress)
// 	if err != nil {
// 		return 0, fmt.Errorf("cannot build the pDEX share key")
// 	}

// 	if amount, ok := allShares[string(shareKey)]; ok {
// 		return amount, nil
// 	} else {
// 		return 0, nil
// 	}

// }

// GetAllShares retrieves all shares in pDEX a user has contributed.
// func (client *IncClient) GetAllShares(beaconHeight uint64, paymentAddress string) ([]*Share, error) {
// 	if beaconHeight == 0 {
// 		bestBlocks, err := client.GetBestBlock()
// 		if err != nil {
// 			return nil, fmt.Errorf("cannot get best blocks: %v", err)
// 		}
// 		beaconHeight = bestBlocks[-1]
// 	}

// 	pdeState, err := client.GetPDEState(beaconHeight)
// 	if err != nil {
// 		return nil, err
// 	}

// 	allShares := pdeState.PDEShares
// 	keyAddr, err := wallet.GetPaymentAddressV1(paymentAddress, false)
// 	if err != nil {
// 		return nil, err
// 	}

// 	res := make([]*Share, 0)
// 	for key, value := range allShares {
// 		if strings.Contains(key, keyAddr) {
// 			sliceStrings := strings.Split(key, "-")
// 			res = append(res, &Share{
// 				TokenID1Str: sliceStrings[2],
// 				TokenID2Str: sliceStrings[3],
// 				ShareAmount: value,
// 			})
// 		}
// 	}

// 	return res, nil
// }

// GetTotalSharesAmount retrieves the total shares' amount of a pDEX pool.
// func (client *IncClient) GetTotalSharesAmount(beaconHeight uint64, tokenID1, tokenID2 string) (uint64, error) {
// 	pdeState, err := client.GetPDEState(beaconHeight)
// 	if err != nil {
// 		return 0, err
// 	}

// 	totalSharesAmount := uint64(0)

// 	allShares := pdeState.PDEShares
// 	poolKey := BuildPDEPoolKey(tokenID1, tokenID2)
// 	for shareKey, amount := range allShares {
// 		if strings.Contains(shareKey, poolKey) {
// 			totalSharesAmount += amount
// 		}
// 	}

// 	return totalSharesAmount, nil
// }

// CheckTradeStatus checks the status of a trading transaction.
// It returns
//	- -1: if an error occurred;
//	- 1: if the trade is accepted;
//	- 2: if the trade is not accepted.
func (client *IncClient) CheckTradeStatus(txHash string) (int, error) {
	responseInBytes, err := client.rpcServer.CheckTradeStatus(txHash)
	if err != nil {
		return -1, err
	}

	var tradeStatus int
	err = rpchandler.ParseResponse(responseInBytes, &tradeStatus)
	if err != nil {
		return -1, err
	}

	return tradeStatus, err
}

// BuildPDEShareKey constructs a key for retrieving contributed shares in pDEX.
func BuildPDEShareKey(beaconHeight uint64, token1ID string, token2ID string, contributorAddress string) ([]byte, error) {
	pdeSharePrefix := []byte("pdeshare-")
	prefix := append(pdeSharePrefix, []byte(fmt.Sprintf("%d-", beaconHeight))...)
	tokenIDs := []string{token1ID, token2ID}
	sort.Strings(tokenIDs)

	var keyAddr string
	var err error
	if len(contributorAddress) == 0 {
		keyAddr = contributorAddress
	} else {
		//Always parse the contributor address into the oldest version for compatibility
		keyAddr, err = wallet.GetPaymentAddressV1(contributorAddress, false)
		if err != nil {
			return nil, err
		}
	}
	return append(prefix, []byte(tokenIDs[0]+"-"+tokenIDs[1]+"-"+keyAddr)...), nil
}

// BuildPDEPoolKey constructs a key for a pool in pDEX.
func BuildPDEPoolKey(token1ID string, token2ID string) string {
	tokenIDs := []string{token1ID, token2ID}
	sort.Strings(tokenIDs)

	return fmt.Sprintf("%v-%v", tokenIDs[0], tokenIDs[1])
}
