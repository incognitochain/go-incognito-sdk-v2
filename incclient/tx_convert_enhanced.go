package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"time"
)

// ConvertAllUTXOs converts all UTXOs of an account for the given tokenIDStr.
// It uses a number of threads working simultaneously to boost up the speed. All of UTXOs v1 will be converted into UTXOs
// v2, therefore, this process is time-consuming and should be run with care. In case you only want to run a single conversion
// transaction at a time, consider using CreateAndSendRawConversionTransaction for better performance.
//
// Parameters:
//		- privateKey: your private key.
//		- tokenIDStr: the id of the asset being converted.
//		- numThreads: the number of workers working simultaneously to convert UTXOs.
func (client *IncClient) ConvertAllUTXOs(privateKey, tokenIDStr string, numThreads int) ([]string, error) {
	if tokenIDStr == common.PRVIDStr {
		return client.convertAllPRVs(privateKey, numThreads)
	} else {
		return client.convertAllTokens(privateKey, tokenIDStr, numThreads)
	}
}

// convertAllPRVs converts all PRV UTXOs v1 of an account into UTXOs v2 in a parallelized manner.
func (client *IncClient) convertAllPRVs(privateKey string, numThreads int) ([]string, error) {
	var txList []string

	utxoV1List, _, err := client.getUTXOsListByVersion(privateKey, common.PRVIDStr, uint8(1))
	if err != nil {
		return txList, err
	}
	if len(utxoV1List) == 0 {
		return nil, fmt.Errorf("no UTXOs to convert")
	} else if len(utxoV1List) <= MaxInputSize {
		txHash, err := client.CreateAndSendRawConversionTransaction(privateKey, common.PRVIDStr)
		if err != nil {
			return nil, err
		}
		return []string{txHash}, nil
	}

	timeOut := time.After(30 * time.Minute)
	errCh := make(chan error)
	txDoneCh := make(chan string)
	txList = make([]string, 0)
	for len(utxoV1List) > 0 {
		Logger.Printf("#numUTXOs: %v\n", len(utxoV1List))
		numWORKERS := 0
		for current := 0; current < len(utxoV1List); current += MaxInputSize {
			next := current + MaxInputSize
			if next > len(utxoV1List) {
				next = len(utxoV1List)
			}
			if next-current < 2 {
				break
			}

			tmpUTXOList := utxoV1List[current:next]
			go client.convertPRVs(numWORKERS, privateKey, tmpUTXOList, txDoneCh, errCh)

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

		utxoV1List, _, err = client.getUTXOsListByVersion(privateKey, common.PRVIDStr, uint8(1))
		if err != nil {
			return txList, err
		}
	}

	return txList, nil
}

// convertAllTokens converts all token UTXOs V2 of an account in a parallelized manner.
func (client *IncClient) convertAllTokens(privateKey, tokenIDStr string, numThreads int) ([]string, error) {
	var txList []string

	utxoV1List, _, err := client.getUTXOsListByVersion(privateKey, tokenIDStr, 1)
	if err != nil {
		return txList, err
	}
	if len(utxoV1List) == 0 {
		return nil, fmt.Errorf("no UTXOs to convert")
	} else if len(utxoV1List) <= MaxInputSize {
		txHash, err := client.CreateAndSendRawConversionTransaction(privateKey, tokenIDStr)
		if err != nil {
			return nil, err
		}
		return []string{txHash}, nil
	}

	timeOut := time.After(30 * time.Minute)
	errCh := make(chan error)
	txDoneCh := make(chan string)
	txList = make([]string, 0)
	for len(utxoV1List) > 0 {
		Logger.Printf("#numUTXOs: %v\n", len(utxoV1List))
		Logger.Printf("Splitting PRVs for paying transaction fees...\n")
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
		Logger.Printf("FINISHED splitting PRVs\n")

		currentPRVIdx := 0
		numWORKERS := 0
		for current := 0; current < len(utxoV1List); current += MaxInputSize {
			next := current + MaxInputSize
			if next > len(utxoV1List) {
				next = len(utxoV1List)
			}
			if next-current < 2 {
				break
			}

			tmpUTXOList := utxoV1List[current:next]
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

			go client.convertTokens(numWORKERS, privateKey, tokenIDStr,
				tmpUTXOList, tmpPRV, tmpPRVIdx, txDoneCh, errCh)

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

		utxoV1List, _, err = client.getUTXOsListByVersion(privateKey, tokenIDStr, 1)
		if err != nil {
			return txList, err
		}
	}

	return txList, nil
}

// convertPRVs creates a transaction that converts a batch of at most MaxInputSize PRV UTXOs v1 into a single UTXO v2.
func (client *IncClient) convertPRVs(id int, privateKey string,
	inputCoins []coin.PlainCoin,
	txDoneCh chan string,
	errCh chan error,
) {
	Logger.Printf("[ID %v] CONVERTING %v UTXOs\n", id, len(inputCoins))
	totalAmount := uint64(0)
	for _, c := range inputCoins {
		totalAmount += c.GetValue()
	}
	if totalAmount <= DefaultPRVFee {
		errCh <- fmt.Errorf("[ID %v] not enough PRV, got %v, want at least %v", id, totalAmount, DefaultPRVFee+1)
		return
	}

	encodedTx, txHash, err := client.CreateConversionTransactionWithInputCoins(privateKey, inputCoins)
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

// convertTokens creates a transaction that converts a batch of at most MaxInputSize token UTXOs V2 into a single UTXO v2.
func (client *IncClient) convertTokens(id int, privateKey, tokenIDStr string,
	inputCoins []coin.PlainCoin,
	prvInputCoins []coin.PlainCoin,
	prvIndices []uint64,
	txDoneCh chan string,
	errCh chan error,
) {
	Logger.Printf("[ID %v] CONVERTING %v TOKEN UTXOs, USING %v PRV UTXOs\n", id, len(inputCoins), len(prvInputCoins))
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

	encodedTx, txHash, err := client.CreateTokenConversionTransactionWithInputCoins(privateKey, tokenIDStr,
		inputCoins, prvInputCoins, prvIndices)
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
