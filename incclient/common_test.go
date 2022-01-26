package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"log"
	"strings"
	"time"
)

const (
	maxAttempts = 30
	numTests    = 10
)

func init() {
	Logger.IsEnable = true
	Logger.Println("This runs before test!!")
}

// waitingCheckTxInBlock waits and checks until a transaction has been included in a block.
//
// In case the transaction is invalid, it stops.
func waitingCheckTxInBlock(txHash string) error {
	for {
		isInBlock, err := ic.CheckTxInBlock(txHash)
		if err != nil {
			if !strings.Contains(err.Error(), "-m") {
				log.Printf("CheckTxInBlock of %v error: %v\n", txHash, err)
				return err
			} else {
				time.Sleep(10 * time.Second)
				continue
			}
		}
		if isInBlock {
			log.Printf("Tx %v is in block\n", txHash)
			return nil
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

// waitingCheckTxInBlock waits and checks until a transaction has been included in a block.
//
// In case the transaction is invalid, it stops.
func waitingCheckBalanceUpdated(privateKey, tokenID string, oldAmount, expectedNewAmount uint64, version uint8) error {
	for {
		balance, err := getBalanceByVersion(privateKey, tokenID, version)
		if err != nil {
			return err
		}

		if balance == oldAmount {
			log.Printf("balance not updated\n")
			time.Sleep(10 * time.Second)
			continue
		}

		if balance != expectedNewAmount {
			return fmt.Errorf("expect balance to be %v, got %v", expectedNewAmount, balance)
		} else {
			log.Printf("balance updated correctly\n")
			return nil
		}
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
		Logger.Printf("\n\ninitialNumUTXOs: %v\n", len(utxo))
		Logger.Printf("Sending funds to tmp accounts...\n")
		amountForEach := common.RandUint64() % 10000
		Logger.Printf("Amount for each: %v\n", amountForEach)
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
		Logger.Printf("txHash %v\n", txHash)
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
				Logger.Printf("numPassed %v\n", numPassed)
			} else {
				time.Sleep(5 * time.Second)
			}
		}

		txHash, err = ic.CreateAndSendRawTransaction(senderPrivateKey, tmpAddresses, feeAmountList, version, nil)
		if err != nil {
			return err
		}
		Logger.Printf("txHash for sending PRV fees %v\n", txHash)
		err = waitingCheckTxInBlock(txHash)
		if err != nil {
			return err
		}
		numPassed = 0
		for numPassed < 5 {
			r := common.RandInt() % MaxOutputSize
			utxo, _, err := ic.getUTXOsListByVersion(tmpPrivateKeys[r], common.PRVIDStr, uint8(version))
			if err != nil {
				return err
			}
			if len(utxo) != 0 {
				numPassed++
				Logger.Printf("numPassed %v\n", numPassed)
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
				Logger.Printf("TxHash %v DONE, numDone %v, numErr %v\n", txHash, numDone, numErr)
			default:
				if numErr == MaxOutputSize {
					return fmt.Errorf("ALL FAILED")
				}
				if numDone == MaxOutputSize {
					Logger.Printf("ALL SUCCESS\n")
					allDone = true
					break
				}
				if numDone+numErr == MaxOutputSize {
					Logger.Printf("ALL FINISHED!!! numDone %v, numErr %v\n", numDone, numErr)
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

	Logger.Printf("[ID %v] sending token %v with version %v\n", id, tokenIDStr, version)

	if tokenIDStr == common.PRVIDStr {
		txHash, err = ic.CreateAndSendRawTransaction(privateKey, addrList, amountList, version, nil)
	} else {
		txHash, err = ic.CreateAndSendRawTokenTransaction(privateKey, addrList, amountList, tokenIDStr, version, nil)
	}
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}
	Logger.Printf("[ID %v] TxHash %v\n", id, txHash)
	err = waitingCheckTxInBlock(txHash)
	if err != nil {
		errCh <- fmt.Errorf("[ID %v] %v", id, err)
		return
	}

	doneCh <- txHash
	return
}

func jsonPrint(val interface{}) error {
	jsb, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		return err
	}

	Logger.Println(string(jsb))
	return nil
}
