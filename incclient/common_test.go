package incclient

import (
	"fmt"
	"strings"
	"time"
)

const (
	MaxAttempts = 30
)

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
