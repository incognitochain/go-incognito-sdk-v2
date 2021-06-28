package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"time"
)

// ConsolidatePRVs consolidates the list of UTXOs of an account, with the given version into a smaller group
// whose size is at most maxUTXOs.
func (client *IncClient) ConsolidatePRVs(privateKey string, version int8, maxUTXOs int) ([]string, error) {
	var txList []string
	if maxUTXOs < 1 {
		return txList, fmt.Errorf("maxUTXOs cannot be less than 1")
	}
	if version > 2 || version < 1 {
		return txList, fmt.Errorf("version %v not supported", version)
	}

	utxoList, idxList, err := client.getUTXOsListByVersion(privateKey, common.PRVIDStr, uint8(version))
	if err != nil {
		return txList, err
	}

	if len(utxoList) <= maxUTXOs {
		log.Printf("already consolidated\n")
		return txList, nil
	}

	timeOut := time.After(30 * time.Minute)
	errCh := make(chan error)
	txDoneCh := make(chan string)
	txList = make([]string, 0)
	for len(utxoList) > maxUTXOs {
		log.Printf("#numUTXOs: %v\n", len(utxoList))
		numWorker := 0
		for current := 0; current < len(utxoList); current += MaxInputSize {
			next := current + MaxInputSize
			if next > len(utxoList) {
				next = len(utxoList)
			}
			if next-current < 2 {
				break
			}

			tmpUTXOList := utxoList[current:next]
			var tmpIdxList []uint64
			if idxList != nil {
				tmpIdxList = idxList[current:next]
			}
			go client.consolidatePRVs(numWorker, privateKey, tmpUTXOList, tmpIdxList, txDoneCh, errCh)

			numWorker++
			if numWorker >= 30 {
				break
			}
		}

		log.Printf("numWorkers: %v\n", numWorker)

		allDone := false
		numErr := 0
		numDone := 0
		for {
			select {
			case txHash := <-txDoneCh:
				numDone++
				txList = append(txList, txHash)
				log.Printf("Finished tx %v, numDone %v, numErr %v\n", txHash, numDone, numErr)
			case err = <-errCh:
				numErr++
				log.Printf("%v\n", err)
			case <-timeOut:
				log.Printf("Timeout!!!!\n")
				return txList, fmt.Errorf("time-out")
			default:
				if numDone+numErr == numWorker {
					log.Printf("All WORKERs FINISHED, numDone %v, numErr %v\n", numDone, numErr)
					allDone = true
					time.Sleep(5 * time.Second)
					break
				}
			}
			if allDone {
				break
			}
		}

		utxoList, idxList, err = client.getUTXOsListByVersion(privateKey, common.PRVIDStr, uint8(version))
		if err != nil {
			return txList, err
		}
	}

	return txList, nil
}

// ConsolidateTokenV1s consolidates the list of UTXOs of an account, with the given version into a smaller group
// whose size is at most maxUTXOs.
func (client *IncClient) ConsolidateTokenV1s(privateKey, tokenIDStr string, maxUTXOs int) ([]string, error) {
	var txList []string
	if maxUTXOs < 1 {
		return txList, fmt.Errorf("maxUTXOs cannot be less than 1")
	}

	utxoList, idxList, err := client.getUTXOsListByVersion(privateKey, tokenIDStr, 1)
	if err != nil {
		return txList, err
	}

	if len(utxoList) <= maxUTXOs {
		log.Printf("already consolidated\n")
		return txList, nil
	}

	timeOut := time.After(30 * time.Minute)
	errCh := make(chan error)
	txDoneCh := make(chan string)
	txList = make([]string, 0)
	for len(utxoList) > maxUTXOs {
		log.Printf("#numUTXOs: %v\n", len(utxoList))
		numWorker := 0
		for current := 0; current < len(utxoList); current += MaxInputSize {
			next := current + MaxInputSize
			if next > len(utxoList) {
				next = len(utxoList)
			}
			if next-current < 2 {
				break
			}

			tmpUTXOList := utxoList[current:next]
			var tmpIdxList []uint64
			if idxList != nil {
				tmpIdxList = idxList[current:next]
			}
			go client.consolidateTokenV1s(numWorker, privateKey, tokenIDStr, tmpUTXOList, tmpIdxList, txDoneCh, errCh)

			numWorker++
			if numWorker >= 30 {
				break
			}
			time.Sleep(3 * time.Second)
		}

		log.Printf("numWorkers: %v\n", numWorker)

		allDone := false
		numErr := 0
		numDone := 0
		for {
			select {
			case txHash := <-txDoneCh:
				numDone++
				txList = append(txList, txHash)
				log.Printf("Finished tx %v, numDone %v, numErr %v\n", txHash, numDone, numErr)
			case err = <-errCh:
				numErr++
				log.Printf("%v\n", err)
			case <-timeOut:
				log.Printf("Timeout!!!!\n")
				return txList, fmt.Errorf("time-out")
			default:
				if numDone+numErr == numWorker {
					log.Printf("All WORKERs FINISHED, numDone %v, numErr %v\n", numDone, numErr)
					allDone = true
					time.Sleep(5 * time.Second)
					break
				}
			}
			if allDone {
				break
			}
		}

		utxoList, idxList, err = client.getUTXOsListByVersion(privateKey, tokenIDStr, 1)
		if err != nil {
			return txList, err
		}
	}

	return txList, nil
}

// consolidatePRVs creates a transaction that consolidates a list of PRV UTXOs into a single UTXO.
func (client *IncClient) consolidatePRVs(id int, privateKey string,
	inputCoins []coin.PlainCoin,
	indices []uint64,
	txDoneCh chan string,
	errCh chan error,
) {
	log.Printf("[ID %v] CONSOLIDATING %v UTXOs, %v INDICES\n", id, len(inputCoins), len(indices))
	totalAmount := uint64(0)
	for _, c := range inputCoins {
		totalAmount += c.GetValue()
	}
	if totalAmount <= DefaultPRVFee {
		errCh <- fmt.Errorf("[ID %v] not enough PRV, got %v, want at least %v", id, totalAmount, DefaultPRVFee+1)
		return
	}

	addr := PrivateKeyToPaymentAddress(privateKey, -1)
	txParam := NewTxParam(privateKey, []string{addr}, []uint64{totalAmount - DefaultPRVFee}, DefaultPRVFee, nil, nil, nil)

	encodedTx, txHash, err := client.CreateRawTransactionWithInputCoins(txParam, inputCoins, indices)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}
	err = client.waitingCheckTxInBlock(txHash)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}

	txDoneCh <- txHash
	log.Printf("[ID %v] FINISHED\n\n", id)
	return
}

// consolidatePRVs creates a transaction that consolidates a list of PRV UTXOs into a single UTXO.
func (client *IncClient) consolidateTokenV1s(id int, privateKey, tokenIDStr string,
	inputCoins []coin.PlainCoin,
	indices []uint64,
	txDoneCh chan string,
	errCh chan error,
) {
	log.Printf("[ID %v] CONSOLIDATING %v TOKEN UTXOs, %v INDICES\n", id, len(inputCoins), len(indices))
	totalAmount := uint64(0)
	for _, c := range inputCoins {
		totalAmount += c.GetValue()
	}

	// estimate token fee
	shardID := GetShardIDFromPrivateKey(privateKey)
	tokenFee, err := client.GetTokenFee(shardID, tokenIDStr)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] cannot estimate token fee: %v", id, err)
	}
	tokenFee = (MaxInputSize * tokenFee) / 10
	log.Printf("[ID %v] tokenFee %v\n", id, tokenFee)
	if totalAmount <= tokenFee {
		errCh <- fmt.Errorf("[ID %v] not enough PRV, got %v, want at least %v", id, totalAmount, tokenFee+1)
		return
	}

	addr := PrivateKeyToPaymentAddress(privateKey, -1)
	txTokenParam := NewTxTokenParam(tokenIDStr, 1, []string{addr}, []uint64{totalAmount - tokenFee}, true, tokenFee, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, txTokenParam, nil, nil)

	encodedTx, txHash, err := client.CreateRawTokenTransactionWithInputCoins(txParam, inputCoins, indices, nil, nil)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}

	log.Printf("[ID %v] TxHash %v\n", id, txHash)
	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}

	err = client.waitingCheckTxInBlock(txHash)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}

	txDoneCh <- txHash
	log.Printf("[ID %v] FINISHED\n\n", id)
	return
}
