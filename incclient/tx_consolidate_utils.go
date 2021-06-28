package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"log"
	"strings"
	"time"
)

// getUTXOsListByVersion returns the list of UTXOs and indices of a privateKey with the given version.
func (client *IncClient) getUTXOsListByVersion(privateKey string,
	tokenIDStr string,
	version uint8,
) ([]coin.PlainCoin, []uint64, error) {
	allUTXOList, allIdxList, err := client.GetUnspentOutputCoinsFromCache(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, nil, err
	}

	utxoList := make([]coin.PlainCoin, 0)
	idxList := make([]uint64, 0)
	for i, utxo := range allUTXOList {
		if utxo.GetVersion() == version { // discard UTXOs with zero value or different version
			utxoList = append(utxoList, utxo)
			idxList = append(idxList, allIdxList[i].Uint64())
		}
	}

	return utxoList, idxList, nil
}

// waitingCheckTxInBlock waits and checks until a transaction has been included in a block.
//
// In case the transaction is invalid, it stops.
func (client *IncClient) waitingCheckTxInBlock(txHash string) error {
	timeOut := time.After(5 * time.Minute)
	for {
		isInBlock, err := client.CheckTxInBlock(txHash)
		if err != nil {
			if !strings.Contains(err.Error(), "-m") {
				log.Printf("CheckTxInBlock of %v error: %v\n", txHash, err)
				return err
			}
		}

		if isInBlock {
			log.Printf("Tx %v is in block\n", txHash)
			return nil
		}

		select {
		case <-timeOut:
			return fmt.Errorf("time-out")
		default:
			time.Sleep(10 * time.Second)
			break
		}
	}
}

func estimateNumTxs(initialNumUTXOs, expectedNumUTXos int) int {
	if initialNumUTXOs <= expectedNumUTXos {
		return 0
	}

	if initialNumUTXOs < expectedNumUTXos+MaxOutputSize {
		return 1
	}

	numTxs := 0
	for initialNumUTXOs > expectedNumUTXos {
		numTxs += initialNumUTXOs / MaxOutputSize
		if initialNumUTXOs%MaxOutputSize == 0 {
			initialNumUTXOs = initialNumUTXOs / MaxOutputSize
		} else {
			initialNumUTXOs = initialNumUTXOs/MaxOutputSize + 1
		}
	}

	return numTxs
}
