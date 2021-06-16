package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"math/big"
	"strings"
	"time"
)

// Utils
var ALPHABET = "abcdefghijklmnopqrstvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	MaxAttempts = 30
)

// randChars returns a string consisting of n alphabet characters.
func randChars(n int) string {
	res := ""
	for i := 0; i < n; i++ {
		r := common.RandInt() % len(ALPHABET)
		res += string(ALPHABET[r])
	}

	return res
}

func calculatePoolAmount(pool *common.PoolInfo, totalShare uint64, shareAmount uint64) (uint64, uint64) {
	shareBig := new(big.Int).SetUint64(shareAmount)
	totalShareBig := new(big.Int).SetUint64(totalShare)

	value1 := new(big.Int).SetUint64(pool.Token1PoolValue)
	value1 = value1.Mul(value1, shareBig)
	value1 = value1.Div(value1, totalShareBig)

	value2 := new(big.Int).SetUint64(pool.Token2PoolValue)
	value2 = value2.Mul(value2, shareBig)
	value2 = value2.Div(value2, totalShareBig)

	return value1.Uint64(), value2.Uint64()
}

// waitingCheckTxInBlock waits and checks until a transaction has been included in a block.
//
// In case the transaction is invalid, it stops.
func waitingCheckTxInBlock(txHash string) error {
	for {
		isInBlock, err := ic.CheckTxInBlock(txHash)
		if err != nil {
			if !strings.Contains(err.Error(), "-m") {
				fmt.Printf("CheckTxInBlock of %v error: %v\n", txHash, err)
				return err
			} else {
				time.Sleep(10 * time.Second)
				continue
			}
		}
		if isInBlock {
			fmt.Printf("Tx %v is in block\n", txHash)
			return nil
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

// END Utils
