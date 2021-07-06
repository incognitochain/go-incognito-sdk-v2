package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"time"
)

const (
	maxUTXOsAfterConsolidated = 10
)

// Consolidate consolidates the list of UTXOs of an account for the given tokenIDStr.
// It uses a number of threads working simultaneously to boost up the consolidating speed.
func (client *IncClient) Consolidate(privateKey, tokenIDStr string, version int8, numThreads int) ([]string, error) {
	if tokenIDStr == common.PRVIDStr {
		return client.ConsolidatePRVs(privateKey, version, numThreads)
	}
	if version == 1 {
		return client.ConsolidateTokenV1s(privateKey, tokenIDStr, numThreads)
	} else {
		return client.ConsolidateTokenV2s(privateKey, tokenIDStr, numThreads)
	}
}

// ConsolidatePRVs consolidates the list of UTXOs of an account, with the given version into a smaller group
// whose size is at most 10.
func (client *IncClient) ConsolidatePRVs(privateKey string, version int8, numThreads int) ([]string, error) {
	var txList []string
	if version > 2 || version < 1 {
		return txList, fmt.Errorf("version %v not supported", version)
	}

	utxoList, idxList, err := client.getUTXOsListByVersion(privateKey, common.PRVIDStr, uint8(version))
	if err != nil {
		return txList, err
	}

	if len(utxoList) <= maxUTXOsAfterConsolidated {
		Logger.Printf("already consolidated\n")
		return txList, nil
	}

	timeOut := time.After(30 * time.Minute)
	errCh := make(chan error)
	txDoneCh := make(chan string)
	txList = make([]string, 0)
	for len(utxoList) > maxUTXOsAfterConsolidated {
		Logger.Printf("#numUTXOs: %v\n", len(utxoList))
		numWORKERS := 0
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
			go client.consolidatePRVs(numWORKERS, privateKey, tmpUTXOList, tmpIdxList, txDoneCh, errCh)

			numWORKERS++
			if numWORKERS >= numThreads {
				break
			}
		}

		Logger.Printf("numWORKERS: %v\n", numWORKERS)

		allDone := false
		numErr := 0
		numDone := 0
		for {
			select {
			case txHash := <-txDoneCh:
				numDone++
				txList = append(txList, txHash)
				Logger.Printf("Finished tx %v, numDone %v, numErr %v\n", txHash, numDone, numErr)
			case err = <-errCh:
				numErr++
				Logger.Printf("%v\n", err)
			case <-timeOut:
				Logger.Printf("Timeout!!!!\n")
				return txList, fmt.Errorf("time-out")
			default:
				if numDone == numWORKERS {
					Logger.Printf("ALL SUCCEEDED\n")
					allDone = true
					break
				}
				if numErr == numWORKERS {
					Logger.Printf("ALL FAILED\n")
					return txList, fmt.Errorf("all thread fails, please try again later")
				}
				if numDone+numErr == numWORKERS {
					Logger.Printf("All WORKERS FINISHED, numDone %v, numErr %v\n", numDone, numErr)
					allDone = true
					break
				}
				time.Sleep(5 * time.Second)
			}
			if allDone {
				time.Sleep(5 * time.Second)
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

// ConsolidateTokenV1s consolidates the list of token UTXOs V1 of an account, with the given version into a smaller group
// whose size is at most 10.
func (client *IncClient) ConsolidateTokenV1s(privateKey, tokenIDStr string, numThreads int) ([]string, error) {
	var txList []string

	utxoList, idxList, err := client.getUTXOsListByVersion(privateKey, tokenIDStr, 1)
	if err != nil {
		return txList, err
	}

	if len(utxoList) <= maxUTXOsAfterConsolidated {
		Logger.Printf("token %v v1 already consolidated, numUTXOs: %v\n", tokenIDStr, len(utxoList))
		return txList, nil
	}

	timeOut := time.After(30 * time.Minute)
	errCh := make(chan error)
	txDoneCh := make(chan string)
	txList = make([]string, 0)
	for len(utxoList) > maxUTXOsAfterConsolidated {
		Logger.Printf("#numUTXOs: %v\n", len(utxoList))
		numWORKERS := 0
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
			go client.consolidateTokenV1s(numWORKERS, privateKey, tokenIDStr, tmpUTXOList, tmpIdxList, txDoneCh, errCh)

			numWORKERS++
			if numWORKERS >= numThreads {
				break
			}
			time.Sleep(3 * time.Second)
		}

		Logger.Printf("numWORKERS: %v\n", numWORKERS)

		allDone := false
		numErr := 0
		numDone := 0
		for {
			select {
			case txHash := <-txDoneCh:
				numDone++
				txList = append(txList, txHash)
				Logger.Printf("Finished tx %v, numDone %v, numErr %v\n", txHash, numDone, numErr)
			case err = <-errCh:
				numErr++
				Logger.Printf("%v\n", err)
			case <-timeOut:
				Logger.Printf("Timeout!!!!\n")
				return txList, fmt.Errorf("time-out")
			default:
				if numDone == numWORKERS {
					Logger.Printf("ALL SUCCEEDED\n")
					allDone = true
					break
				}
				if numErr == numWORKERS {
					Logger.Printf("ALL FAILED\n")
					return txList, fmt.Errorf("all thread fails, please try again later")
				}
				if numDone+numErr == numWORKERS {
					Logger.Printf("All WORKERS FINISHED, numDone %v, numErr %v\n", numDone, numErr)
					allDone = true
					break
				}
				time.Sleep(5 * time.Second)
			}
			if allDone {
				time.Sleep(5 * time.Second)
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

// ConsolidateTokenV2s consolidates the list of token UTXOs V2 of an account, with the given version into a smaller group
// whose size is at most 10.
func (client *IncClient) ConsolidateTokenV2s(privateKey, tokenIDStr string, numThreads int) ([]string, error) {
	var txList []string

	utxoList, idxList, err := client.getUTXOsListByVersion(privateKey, tokenIDStr, 2)
	if err != nil {
		return txList, err
	}

	if len(utxoList) <= maxUTXOsAfterConsolidated {
		Logger.Printf("token %v v2 already consolidated, numUTXOs: %v\n", tokenIDStr, len(utxoList))
		return txList, nil
	}

	timeOut := time.After(30 * time.Minute)
	errCh := make(chan error)
	txDoneCh := make(chan string)
	txList = make([]string, 0)
	for len(utxoList) > maxUTXOsAfterConsolidated {
		Logger.Printf("#numUTXOs: %v\n", len(utxoList))
		txHash, err := client.splitPRVForFees(privateKey, 2, numThreads)
		if err != nil {
			return txList, err
		}
		if txHash != "" {
			txList = append(txList, txHash)
		}
		prvUTXOList, prvIndices, err := client.getUTXOsListByVersion(privateKey, common.PRVIDStr, 2)
		if err != nil {
			return txList, err
		}

		currentPRVIdx := 0
		numWORKERS := 0
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
			var tmpPRV []coin.PlainCoin
			var tmpPRVIdx []uint64
			for i := currentPRVIdx; i < len(prvUTXOList); i++ {
				if prvUTXOList[i].GetValue() >= DefaultPRVFee {
					tmpPRV = []coin.PlainCoin{prvUTXOList[i]}
					tmpPRVIdx = []uint64{prvIndices[i]}
					currentPRVIdx = i + 1
					break
				}
			}

			if tmpPRV == nil || tmpPRVIdx == nil {
				return txList, fmt.Errorf("cannot get PRV UTXO to payfee for index %v", current)
			}

			Logger.Printf("[ID %v] PRV idx: %v, pubKey: %v\n", numWORKERS, tmpPRVIdx[0],
				base58.Base58Check{}.Encode(tmpPRV[0].GetPublicKey().ToBytesS(), 0x00))

			go client.consolidateTokenV2s(numWORKERS, privateKey, tokenIDStr,
				tmpUTXOList, tmpIdxList, tmpPRV, tmpPRVIdx, txDoneCh, errCh)

			numWORKERS++
			if numWORKERS >= numThreads {
				break
			}
			time.Sleep(3 * time.Second)
		}

		Logger.Printf("numWORKERS: %v\n", numWORKERS)

		allDone := false
		numErr := 0
		numDone := 0
		for {
			select {
			case txHash := <-txDoneCh:
				numDone++
				txList = append(txList, txHash)
				Logger.Printf("Finished tx %v, numDone %v, numErr %v\n", txHash, numDone, numErr)
			case err = <-errCh:
				numErr++
				Logger.Printf("%v\n", err)
			case <-timeOut:
				Logger.Printf("Timeout!!!!\n")
				return txList, fmt.Errorf("time-out")
			default:
				if numDone == numWORKERS {
					Logger.Printf("ALL SUCCEEDED\n")
					allDone = true
					break
				}
				if numErr == numWORKERS {
					Logger.Printf("ALL FAILED\n")
					return txList, fmt.Errorf("all thread fails, please try again later")
				}
				if numDone+numErr == numWORKERS {
					Logger.Printf("All WORKERS FINISHED, numDone %v, numErr %v\n", numDone, numErr)
					allDone = true
					break
				}
				time.Sleep(5 * time.Second)
			}
			if allDone {
				time.Sleep(5 * time.Second)
				break
			}
		}

		utxoList, idxList, err = client.getUTXOsListByVersion(privateKey, tokenIDStr, 2)
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
	Logger.Printf("[ID %v] CONSOLIDATING %v UTXOs, %v INDICES\n", id, len(inputCoins), len(indices))
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
	Logger.Printf("[ID %v] TxHash %v\n", id, txHash)
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
	Logger.Printf("[ID %v] FINISHED\n\n", id)
	return
}

// consolidateTokenV1s creates a transaction that consolidates a list of token UTXOs v1 into a single UTXO.
func (client *IncClient) consolidateTokenV1s(id int, privateKey, tokenIDStr string,
	inputCoins []coin.PlainCoin,
	indices []uint64,
	txDoneCh chan string,
	errCh chan error,
) {
	Logger.Printf("[ID %v] CONSOLIDATING %v TOKEN UTXOs, %v INDICES\n", id, len(inputCoins), len(indices))
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
	Logger.Printf("[ID %v] tokenFee %v\n", id, tokenFee)
	if totalAmount <= tokenFee {
		errCh <- fmt.Errorf("[ID %v] not enough token, got %v, want at least %v", id, totalAmount, tokenFee+1)
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
	Logger.Printf("[ID %v] TxHash %v\n", id, txHash)
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
	Logger.Printf("[ID %v] FINISHED\n\n", id)
	return
}

// consolidateTokenV2s creates a transaction that consolidates a list of token UTXOs V2 into a single UTXO.
func (client *IncClient) consolidateTokenV2s(id int, privateKey, tokenIDStr string,
	inputCoins []coin.PlainCoin,
	indices []uint64,
	prvInputCoins []coin.PlainCoin,
	prvIndices []uint64,
	txDoneCh chan string,
	errCh chan error,
) {
	Logger.Printf("[ID %v] CONSOLIDATING %v TOKEN UTXOs, %v INDICES\n", id, len(inputCoins), len(indices))
	totalAmount := uint64(0)
	for _, c := range inputCoins {
		totalAmount += c.GetValue()
	}

	totalPRVAmount := uint64(0)
	for _, c := range prvInputCoins {
		totalPRVAmount += c.GetValue()
	}
	if totalPRVAmount < DefaultPRVFee {
		errCh <- fmt.Errorf("[ID %v] not enough PRV, got %v, want at least %v", id, totalAmount, DefaultPRVFee)
		return
	}

	addr := PrivateKeyToPaymentAddress(privateKey, -1)
	txTokenParam := NewTxTokenParam(tokenIDStr, 1, []string{addr}, []uint64{totalAmount}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, txTokenParam, nil, nil)

	encodedTx, txHash, err := client.CreateRawTokenTransactionWithInputCoins(txParam, inputCoins, indices, prvInputCoins, prvIndices)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}
	Logger.Printf("[ID %v] TxHash %v\n", id, txHash)
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
	Logger.Printf("[ID %v] FINISHED\n\n", id)
	return
}
