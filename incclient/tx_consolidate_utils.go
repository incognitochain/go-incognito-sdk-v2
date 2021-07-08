package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
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

func (client *IncClient) splitPRVForFees(privateKey string, version uint8, numThreads int) (string, error) {
	Logger.Printf("Splitting PRV for numThreads %v\n", numThreads)
	addr := PrivateKeyToPaymentAddress(privateKey, -1)
	if len(addr) == 0 {
		return "", fmt.Errorf("private key is invalid")
	}

	utxoList, _, err := client.getUTXOsListByVersion(privateKey, common.PRVIDStr, version)
	if err != nil {
		return "", err
	}
	totalAmount := uint64(0)
	numRequiredUTXOs := 0 // a required UTXO is an UTXO whose value is greater than the DefaultPRVFee.
	for _, c := range utxoList {
		if c.GetValue() >= DefaultPRVFee {
			numRequiredUTXOs++
		}
		totalAmount += c.GetValue()
	}
	if totalAmount < uint64(numThreads+1)*DefaultPRVFee {
		return "", fmt.Errorf("require at least %v nano PRV of version %v, got %v", uint64(numThreads+1)*DefaultPRVFee, version, totalAmount)
	}
	if numRequiredUTXOs >= numThreads {
		Logger.Log.Printf("Already have enough UTXOs\n")
		return "", nil
	}

	// create a sample of addresses and amounts to split PRVs.
	addrList := make([]string, 0)
	amountList := make([]uint64, 0)
	for i := 0; i < numThreads; i++ {
		addrList = append(addrList, addr)
		amountList = append(amountList, DefaultPRVFee)
	}

	txHash, err := client.CreateAndSendRawTransaction(privateKey, addrList, amountList, int8(version), nil)
	if err != nil {
		return "", err
	}
	Logger.Printf("TxHash for splitting PRV fees %v\n", txHash)
	err = client.waitingCheckTxInBlock(txHash)
	if err != nil {
		return txHash, err
	}

	// check if we have enough PRV UTXOs
	Logger.Printf("Checking UTXOs updated...\n")
	for {
		numRequiredUTXOs = 0
		utxoList, _, err := client.getUTXOsListByVersion(privateKey, common.PRVIDStr, version)
		if err != nil {
			return "", err
		}

		for _, c := range utxoList {
			if c.GetValue() >= DefaultPRVFee {
				numRequiredUTXOs++
			}
		}
		if numRequiredUTXOs >= numThreads {
			break
		}

		time.Sleep(10 * time.Second)
	}
	Logger.Printf("UTXOs updated\n\n")

	return txHash, nil
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
				Logger.Printf("CheckTxInBlock of %v error: %v\n", txHash, err)
				return err
			}
		}

		if isInBlock {
			Logger.Printf("Tx %v is in block\n", txHash)
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
