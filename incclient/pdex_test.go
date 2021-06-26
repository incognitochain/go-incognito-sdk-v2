package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"math"
	"testing"
	"time"
)

func TestIncClient_CreateAndSendPDETradeTransaction(t *testing.T) {
	ic, err := NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := ""

	//Trade PRV to tokens
	tokenToSell := common.PRVIDStr
	tokenToBuy := "02a41194b536aa20960fd62bd4937a895fcc7c7d84a83bf212a349df2b6ea1f2"
	sellAmount := uint64(5000000)
	expectedAmount, err := ic.CheckXPrice(tokenToSell, tokenToBuy, sellAmount)
	if err != nil {
		panic(err)
	}
	tradingFee := uint64(100)

	txHash, err := ic.CreateAndSendPDETradeTransaction(privateKey, tokenToSell, tokenToBuy, sellAmount, expectedAmount, tradingFee)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}

func TestIncClient_CreateAndSendCrossPDETradeTransaction(t *testing.T) {
	ic, err := NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	privateKey := ""

	//Trade token to token
	tokenToBuy := "0795495cb9eb84ae7bd8c8494420663b9a1642c7bbc99e57b04d536db9001d0e"
	tokenToSell := "02a41194b536aa20960fd62bd4937a895fcc7c7d84a83bf212a349df2b6ea1f2"
	sellAmount := uint64(50000)
	expectedAmount, err := ic.CheckXPrice(tokenToSell, tokenToBuy, sellAmount)
	if err != nil {
		panic(err)
	}
	tradingFee := uint64(1000)

	txHash, err := ic.CreateAndSendPDETradeTransaction(privateKey, tokenToSell, tokenToBuy, sellAmount, expectedAmount, tradingFee)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %v\n", txHash)
}

func TestIncClient_CreateAndSendPDEContributeTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	// init params
	privateKey := ""
	addr := PrivateKeyToPaymentAddress(privateKey, -1)
	tokenID1 := common.PRVIDStr
	tokenID2 := "00000000000000000000000000000000000000000000000000000000000000ff"
	pairID := "INC" + randChars(10)

	oldBalance1, err := ic.GetBalance(privateKey, tokenID1)
	if err != nil {
		panic(err)
	}

	oldBalance2, err := ic.GetBalance(privateKey, tokenID2)
	if err != nil {
		panic(err)
	}

	// if balances are insufficient, stop
	if oldBalance1 == 0 || oldBalance2 == 0 {
		panic(fmt.Errorf("balances insuffient: %v, %v\n", oldBalance1, oldBalance2))
	}

	fmt.Printf("oldBalance1: %v, oldBalance2: %v\n", oldBalance1, oldBalance2)

	// retrieve the total shared amount for the pool
	oldTotalShares, err := ic.GetTotalSharesAmount(0, tokenID1, tokenID2)
	if err != nil {
		panic(err)
	}

	// get the current shared amount of the user
	oldShare, err := ic.GetShareAmount(0, tokenID1, tokenID2, addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current totalShare: %v, myShare: %v\n", oldTotalShares, oldShare)

	// calculate the amount for each
	var token1Amount, token2Amount uint64
	pool, err := ic.GetPDEPoolPair(0, tokenID1, tokenID2)
	contributedShare := uint64(0)
	attempt := 0
	for attempt < MaxAttempts {
		minAmount := uint64(math.Min(float64(oldBalance1), float64(oldBalance2)))
		minAmount = uint64(math.Min(float64(minAmount), float64(oldTotalShares)))
		contributedShare = 1 + common.RandUint64()%(minAmount/10)
		if pool == nil {
			token1Amount = contributedShare
			token2Amount = common.RandUint64() % oldBalance2
		} else {
			token1Amount, token2Amount = calculatePoolAmount(pool, oldTotalShares, contributedShare)
			if pool.Token2IDStr == tokenID1 {
				token1Amount, token2Amount = token2Amount, token1Amount
			}
		}

		if token1Amount < oldBalance1 && token2Amount < oldBalance2 {
			fmt.Printf("Contributed share: %v\n", contributedShare)
			break
		}
		attempt += 1
	}
	if attempt > MaxAttempts {
		panic("cannot calculate contributed amounts")
	}

	// Contributed Tx 1
	txHash, err := ic.CreateAndSendPDEContributeTransaction(privateKey, pairID, tokenID1, token1Amount, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contributed TxHash1: %v, token: %v, pairID: %v, amount: %v\n", txHash, tokenID1, pairID, token1Amount)
	err = waitingCheckTxInBlock(txHash)
	if err != nil {
		panic(err)
	}

	// Contributed Tx 2
	txHash, err = ic.CreateAndSendPDEContributeTransaction(privateKey, pairID, tokenID2, token2Amount, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contributed TxHash2: %v, token: %v, pairID: %v, amount: %v\n", txHash, tokenID2, pairID, token2Amount)
	err = waitingCheckTxInBlock(txHash)
	if err != nil {
		panic(err)
	}

	attempt = 0
	for attempt < MaxAttempts {
		newTotalShares, err := ic.GetTotalSharesAmount(0, tokenID1, tokenID2)
		if err != nil {
			panic(err)
		}
		fmt.Printf("newTotalShares: %v\n", newTotalShares)

		if newTotalShares != oldTotalShares {
			diff := float64(newTotalShares) - float64(oldTotalShares)
			if math.Abs(diff-float64(contributedShare)) < 10 {
				newShare, err := ic.GetShareAmount(0, tokenID1, tokenID2, addr)
				if err != nil {
					panic(err)
				}

				if newShare-oldShare == uint64(diff) {
					fmt.Printf("Contributed share success!\n\n")
					break
				} else {
					panic(fmt.Errorf("expected newShare - oldShare = %v, got %v", contributedShare, newShare-oldShare))
				}
			} else {
				panic(fmt.Errorf("expect newTotalShares - oldTotalShares = %v, got %v", contributedShare, diff))
			}
		}
		attempt += 1
		time.Sleep(10 * time.Second)
	}
}

func TestIncClient_CreateAndSendPDEWithdrawalTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	// init params
	privateKey := ""
	addr := PrivateKeyToPaymentAddress(privateKey, -1)
	tokenID1 := common.PRVIDStr
	tokenID2 := "00000000000000000000000000000000000000000000000000000000000000ff"

	oldBalance1, err := ic.GetBalance(privateKey, tokenID1)
	if err != nil {
		panic(err)
	}

	oldBalance2, err := ic.GetBalance(privateKey, tokenID2)
	if err != nil {
		panic(err)
	}

	// if balances are insufficient, stop
	if oldBalance1 == 0 || oldBalance2 == 0 {
		panic(fmt.Errorf("balances insuffient: %v, %v\n", oldBalance1, oldBalance2))
	}

	fmt.Printf("oldBalance1: %v, oldBalance2: %v\n", oldBalance1, oldBalance2)

	// retrieve the total shared amount for the pool
	oldTotalShares, err := ic.GetTotalSharesAmount(0, tokenID1, tokenID2)
	if err != nil {
		panic(err)
	}

	// get the current shared amount of the user
	oldShare, err := ic.GetShareAmount(0, tokenID1, tokenID2, addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current totalShare: %v, myShare: %v\n", oldTotalShares, oldShare)

	if oldShare <= 1 {
		panic("not enough share")
	}

	shareAmount := 1 + common.RandUint64()%(oldShare-1)
	fmt.Printf("WithdrawShare: %v\n", shareAmount)

	// get current pool information
	pool, err := ic.GetPDEPoolPair(0, tokenID1, tokenID2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current pool: %v\n", *pool)

	// calculate the expected amounts
	var expectedValue1, expectedValue2 uint64
	expectedValue1, expectedValue2 = calculatePoolAmount(pool, oldTotalShares, shareAmount)
	if pool.Token1IDStr == tokenID2 {
		expectedValue1, expectedValue2 = expectedValue2, expectedValue1
	}

	fmt.Printf("expectedValue1: %v, expectedValue2: %v\n", expectedValue1, expectedValue2)

	txHash, err := ic.CreateAndSendPDEWithdrawalTransaction(privateKey, tokenID1, tokenID2, shareAmount, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TxHash: %v\n", txHash)

	err = waitingCheckTxInBlock(txHash)
	if err != nil {
		panic(err)
	}

	attempt := 0
	for attempt < MaxAttempts {
		newBalance1, err := ic.GetBalance(privateKey, tokenID1)
		if err != nil {
			panic(err)
		}

		newBalance2, err := ic.GetBalance(privateKey, tokenID2)
		if err != nil {
			panic(err)
		}

		fmt.Printf("newBalances: %v, %v\n", newBalance1, newBalance2)

		diff1 := float64(newBalance1) - float64(oldBalance1)
		diff2 := float64(newBalance2) - float64(oldBalance2)

		if tokenID1 == common.PRVIDStr {
			diff1 += float64(DefaultPRVFee)
		} else if tokenID2 == common.PRVIDStr {
			diff2 += float64(DefaultPRVFee)
		}

		if diff1 == 0 || diff2 == 0 {
			attempt += 1
			time.Sleep(10 * time.Second)
			continue
		} else {
			fmt.Printf("diffs: %v, %v\n", diff1, diff2)
			if math.Abs(diff1-float64(expectedValue1)) > 10 {
				panic(fmt.Errorf("token %v: expected received %v, got %v", tokenID1, expectedValue1, diff1))
			}

			if math.Abs(diff2-float64(expectedValue2)) > 10 {
				panic(fmt.Errorf("token %v: expected received %v, got %v", tokenID2, expectedValue2, diff2))
			}

			break
		}
	}
}
