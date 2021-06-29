package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"log"
	"testing"
	"time"
)

func TestIncClient_ConsolidatePRVs(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	masterPrivateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"

	minInitialUTXOs := 100 + common.RandInt()%3900
	numThreads := 20
	for i := 0; i < numTests; i++ {
		version := int8(1 + common.RandInt()%2)
		log.Printf("==================== TEST %v, VERSION %v ====================\n", i, version)
		expectedNumUTXOs := 1 + common.RandInt()%29
		log.Printf("ExpectedNumUTXOs: %v, minInitial: %v\n", expectedNumUTXOs, minInitialUTXOs)

		// preparing UTXOs
		err = prepareUTXOs(masterPrivateKey, privateKey, common.PRVIDStr, minInitialUTXOs, version)
		if err != nil {
			panic(err)
		}

		// checking UTXOs
		utxo, _, err := ic.getUTXOsListByVersion(privateKey, common.PRVIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		if len(utxo) < minInitialUTXOs {
			panic(fmt.Sprintf("require at least %v UTXOs, got %v", minInitialUTXOs, len(utxo)))
		}
		log.Printf("minInitialUTXOs: %v\n", len(utxo))

		// consolidating UTXOs
		_, err = ic.ConsolidatePRVs(privateKey, version, numThreads)
		if err != nil {
			panic(err)
		}

		utxo, _, err = ic.getUTXOsListByVersion(privateKey, common.PRVIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		if len(utxo) > expectedNumUTXOs {
			panic(fmt.Sprintf("require at most %v UTXOs, got %v", expectedNumUTXOs, len(utxo)))
		}
		log.Printf("numUTXOs: %v\n", len(utxo))

		minInitialUTXOs = 100 + common.RandInt()%3900
		log.Printf("==================== FINISHED TEST %v ====================\n\n", i)
	}
}

func TestIncClient_ConsolidateTokenV1s(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	masterPrivateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"
	tokenIDStr := "f3e586e281d275ea2059e35ae434d0431947d2b49466b6d2479808378268f822"
	minInitialUTXOs := 100 + common.RandInt()%100
	numThreads := 20
	for i := 0; i < numTests; i++ {
		version := int8(1)
		log.Printf("==================== TEST %v, VERSION %v ====================\n", i, version)
		expectedNumUTXOs := 1 + common.RandInt()%29
		log.Printf("ExpectedNumUTXOs: %v, minInitial: %v\n", expectedNumUTXOs, minInitialUTXOs)

		// preparing UTXOs
		err = prepareUTXOs(masterPrivateKey, privateKey, tokenIDStr, minInitialUTXOs, version)
		if err != nil {
			panic(err)
		}

		// checking UTXOs
		utxo, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		if len(utxo) < minInitialUTXOs {
			panic(fmt.Sprintf("require at least %v UTXOs, got %v", minInitialUTXOs, len(utxo)))
		}
		log.Printf("minInitialUTXOs: %v\n", len(utxo))

		balance, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		log.Printf("oldBalance: %v\n", balance)

		// consolidating UTXOs
		_, err = ic.ConsolidateTokenV1s(privateKey, tokenIDStr, numThreads)
		if err != nil {
			panic(err)
		}

		utxo, _, err = ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		if len(utxo) > expectedNumUTXOs {
			panic(fmt.Sprintf("require at most %v UTXOs, got %v", expectedNumUTXOs, len(utxo)))
		}
		log.Printf("numUTXOs: %v\n", len(utxo))

		balance, err = getBalanceByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		log.Printf("newBalance: %v\n", balance)

		minInitialUTXOs = 100 + common.RandInt()%3900
		log.Printf("==================== FINISHED TEST %v ====================\n\n", i)
	}
}

func prepareUTXOs(senderPrivateKey, receiverPrivateKey, tokenIDStr string, numUTXOs int, version int8) error {
	utxo, _, err := ic.getUTXOsListByVersion(receiverPrivateKey, tokenIDStr, uint8(version))
	if err != nil {
		return err
	}

	addr := PrivateKeyToPaymentAddress(receiverPrivateKey, -1)
	addrList := make([]string, 0)
	for i := 0; i < MaxOutputSize; i++ {
		addrList = append(addrList, addr)
	}

	shardID := GetShardIDFromPrivateKey(senderPrivateKey)
	tmpPrivateKeys := make([]string, 0)
	tmpAddresses := make([]string, 0)
	for i := 0; i < MaxOutputSize; i++ {
		randWallet, err := wallet.GenRandomWalletForShardID(shardID)
		if err != nil {
			return err
		}
		randPrivateKey := randWallet.Base58CheckSerialize(wallet.PrivateKeyType)
		randAddr := randWallet.Base58CheckSerialize(wallet.PaymentAddressType)
		tmpPrivateKeys = append(tmpPrivateKeys, randPrivateKey)
		tmpAddresses = append(tmpAddresses, randAddr)

		err = ic.SubmitKey(randPrivateKey)
		if err != nil {
			return err
		}
	}

	for len(utxo) < numUTXOs {
		log.Printf("\n\ninitialNumUTXOs: %v\n", len(utxo))
		log.Printf("Sending funds to tmp accounts...\n")
		amountForEach := common.RandUint64() % 10000
		log.Printf("Amount for each: %v\n", amountForEach)
		tmpAmountList := make([]uint64, 0)
		amountList := make([]uint64, 0)
		feeAmountList := make([]uint64, 0)
		for i := 0; i < MaxOutputSize; i++ {
			tmpAmountList = append(tmpAmountList, amountForEach*30)
			amountList = append(amountList, amountForEach)
			feeAmountList = append(feeAmountList, DefaultPRVFee)
		}
		var txHash string
		if tokenIDStr == common.PRVIDStr {
			txHash, err = ic.CreateAndSendRawTransaction(senderPrivateKey, tmpAddresses, tmpAmountList, version, nil)
		} else {
			txHash, err = ic.CreateAndSendRawTokenTransaction(senderPrivateKey, tmpAddresses, tmpAmountList, tokenIDStr, version, nil)
		}
		if err != nil {
			return err
		}
		log.Printf("txHash %v\n", txHash)
		err = waitingCheckTxInBlock(txHash)
		if err != nil {
			return err
		}
		numPassed := 0
		for numPassed < 5 {
			r := common.RandInt() % MaxOutputSize
			utxo, _, err := ic.getUTXOsListByVersion(tmpPrivateKeys[r], tokenIDStr, uint8(version))
			if err != nil {
				return err
			}
			if len(utxo) != 0 {
				numPassed++
				log.Printf("numPasssed %v\n", numPassed)
			} else {
				time.Sleep(5 * time.Second)
			}
		}

		txHash, err = ic.CreateAndSendRawTransaction(senderPrivateKey, tmpAddresses, feeAmountList, version, nil)
		if err != nil {
			return err
		}
		log.Printf("txHash %v\n", txHash)
		err = waitingCheckTxInBlock(txHash)
		if err != nil {
			return err
		}
		numPassed = 0
		for numPassed < 5 {
			r := common.RandInt() % MaxOutputSize
			utxo, _, err := ic.getUTXOsListByVersion(tmpPrivateKeys[r], tokenIDStr, uint8(version))
			if err != nil {
				return err
			}
			if len(utxo) != 0 {
				numPassed++
				log.Printf("numPasssed %v\n", numPassed)
			} else {
				time.Sleep(5 * time.Second)
			}
		}

		errCh := make(chan error)
		doneCh := make(chan string)
		for i := 0; i < MaxOutputSize; i++ {
			go send(i, tmpPrivateKeys[i], tokenIDStr, addrList, amountList, version, doneCh, errCh)
		}

		allDone := false
		numErr := 0
		numDone := 0
		for {
			select {
			case err = <-errCh:
				log.Println(err)
				numErr++
			case txHash = <-doneCh:
				numDone++
				log.Printf("TxHash %v DONE, numDone %v, numErr %v\n", txHash, numDone, numErr)
			default:
				if numErr == MaxOutputSize {
					return fmt.Errorf("ALL FAILED")
				}
				if numDone == MaxOutputSize {
					log.Printf("ALL SUCCESS\n")
					allDone = true
					break
				}
				if numDone+numErr == MaxOutputSize {
					log.Printf("ALL FINISHED!!! numDone %v, numErr %v\n", numDone, numErr)
					allDone = true
					break
				}
				time.Sleep(10 * time.Second)
			}
			if allDone {
				break
			}
		}

		time.Sleep(20 * time.Second)
		utxo, _, err = ic.getUTXOsListByVersion(receiverPrivateKey, tokenIDStr, uint8(version))
		if err != nil {
			return err
		}
	}

	return nil
}

func send(id int, privateKey, tokenIDStr string, addrList []string, amountList []uint64,
	version int8, doneCh chan string, errCh chan error) {
	var txHash string
	var err error

	log.Printf("[ID %v] version %v\n", id, version)

	if tokenIDStr == common.PRVIDStr {
		txHash, err = ic.CreateAndSendRawTransaction(privateKey, addrList, amountList, version, nil)
	} else {
		txHash, err = ic.CreateAndSendRawTokenTransaction(privateKey, addrList, amountList, tokenIDStr, version, nil)
	}
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}
	log.Printf("[ID %v] TxHash %v\n", id, txHash)
	err = waitingCheckTxInBlock(txHash)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}

	doneCh <- txHash
	return
}
