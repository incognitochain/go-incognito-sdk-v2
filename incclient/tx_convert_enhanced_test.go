package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"testing"
)

func TestIncClient_ConvertAllPRVs(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}
	//Enable the logger
	Logger.IsEnable = true

	masterPrivateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	privateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	MaxUTXO := 2000
	initialUTXOs := 100 + common.RandInt()%(MaxUTXO-100)
	numThreads := 20
	tokenIDStr := common.PRVIDStr
	for i := 0; i < numTests; i++ {
		Logger.Printf("==================== TEST %v ====================\n", i)
		Logger.Printf("RequiredUTXOs: %v\n", initialUTXOs)

		// preparing UTXOs
		err = prepareUTXOs(masterPrivateKey, privateKey, tokenIDStr, initialUTXOs, 1)
		if err != nil {
			panic(err)
		}

		Logger.Printf("\nBEGIN CONVERTING UTXOs...\n")

		// checking UTXOs
		oldUTXOsV1, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(1))
		if err != nil {
			panic(err)
		}
		if len(oldUTXOsV1) < initialUTXOs {
			panic(fmt.Sprintf("require at least %v UTXOs, got %v", initialUTXOs, len(oldUTXOsV1)))
		}
		utxoV2, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(2))
		if err != nil {
			panic(err)
		}
		Logger.Printf("#oldUTXOsV1: %v, #oldUTXOsV2: %v\n", len(oldUTXOsV1), len(utxoV2))

		oldBalanceV1, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(1))
		if err != nil {
			panic(err)
		}
		oldBalanceV2, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(2))
		if err != nil {
			panic(err)
		}
		Logger.Printf("oldBalanceV1: %v, oldBalanceV2: %v\n", oldBalanceV1, oldBalanceV2)

		// consolidating UTXOs
		_, err = ic.ConvertAllUTXOs(privateKey, tokenIDStr, numThreads)
		if err != nil {
			panic(err)
		}

		newUTXOsV1, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(1))
		if err != nil {
			panic(err)
		}
		if len(newUTXOsV1) > 0 {
			panic(fmt.Sprintf("require no UTXOs v1, got %v", len(newUTXOsV1)))
		}
		newUTXOsV2, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(2))
		if err != nil {
			panic(err)
		}
		Logger.Printf("#newUTXOsV1: %v, #newUTXOsV2: %v\n", len(newUTXOsV1), len(newUTXOsV2))
		newBalanceV1, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(1))
		if err != nil {
			panic(err)
		}
		newBalanceV2, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(2))
		if err != nil {
			panic(err)
		}
		Logger.Printf("newBalanceV1: %v, newBalanceV2: %v\n", newBalanceV1, newBalanceV2)

		initialUTXOs = 100 + common.RandInt()%(MaxUTXO-100)
		Logger.Printf("FINISHED CONVERTING UTXOs\n")
		Logger.Printf("==================== FINISHED TEST %v ====================\n\n", i)
	}
}

func TestIncClient_ConvertAllTokens(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}
	//Enable the logger
	Logger.IsEnable = true

	masterPrivateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	privateKey := "11111113iP7vLqNpK2RPPmwkQgaXf4c6dzto5RfyNYTsk8L1hNLajtcPRMihKpD9Tg8N8UkGrGso3iAUHaDbDDT2rrf7QXwAGADHkuV5A1U"
	MaxUTXO := 2000
	initialUTXOs := 100 + common.RandInt()%(MaxUTXO-100)
	numThreads := 20
	tokenIDStr := "f3e586e281d275ea2059e35ae434d0431947d2b49466b6d2479808378268f822"
	for i := 0; i < numTests; i++ {
		Logger.Printf("==================== TEST %v ====================\n", i)
		Logger.Printf("RequiredUTXOs: %v\n", initialUTXOs)

		// preparing UTXOs
		err = prepareUTXOs(masterPrivateKey, privateKey, tokenIDStr, initialUTXOs, 1)
		if err != nil {
			panic(err)
		}

		Logger.Printf("\nBEGIN CONVERTING UTXOs...\n")

		// checking UTXOs
		oldUTXOsV1, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(1))
		if err != nil {
			panic(err)
		}
		if len(oldUTXOsV1) < initialUTXOs {
			panic(fmt.Sprintf("require at least %v UTXOs, got %v", initialUTXOs, len(oldUTXOsV1)))
		}
		utxoV2, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(2))
		if err != nil {
			panic(err)
		}
		Logger.Printf("#oldUTXOsV1: %v, #oldUTXOsV2: %v\n", len(oldUTXOsV1), len(utxoV2))

		oldBalanceV1, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(1))
		if err != nil {
			panic(err)
		}
		oldBalanceV2, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(2))
		if err != nil {
			panic(err)
		}
		Logger.Printf("oldBalanceV1: %v, oldBalanceV2: %v\n", oldBalanceV1, oldBalanceV2)

		// consolidating UTXOs
		_, err = ic.ConvertAllUTXOs(privateKey, tokenIDStr, numThreads)
		if err != nil {
			panic(err)
		}

		newUTXOsV1, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(1))
		if err != nil {
			panic(err)
		}
		if len(newUTXOsV1) > 0 {
			panic(fmt.Sprintf("require no UTXOs v1, got %v", len(newUTXOsV1)))
		}
		newUTXOsV2, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(2))
		if err != nil {
			panic(err)
		}
		Logger.Printf("#newUTXOsV1: %v, #newUTXOsV2: %v\n", len(newUTXOsV1), len(newUTXOsV2))
		newBalanceV1, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(1))
		if err != nil {
			panic(err)
		}
		newBalanceV2, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(2))
		if err != nil {
			panic(err)
		}
		Logger.Printf("newBalanceV1: %v, newBalanceV2: %v\n", newBalanceV1, newBalanceV2)

		initialUTXOs = 100 + common.RandInt()%(MaxUTXO-100)
		Logger.Printf("FINISHED CONVERTING UTXOs\n")
		Logger.Printf("==================== FINISHED TEST %v ====================\n\n", i)
	}
}
