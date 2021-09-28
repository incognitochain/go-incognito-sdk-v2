package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"testing"
)

func TestIncClient_ConsolidatePRVs(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	masterPrivateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"
	MaxUTXO := 200
	minInitialUTXOs := 100 + common.RandInt()%(MaxUTXO-100)
	numThreads := 20
	for i := 0; i < numTests; i++ {
		//version := int8(1 + common.RandInt()%2)
		version := int8(2)
		Logger.Printf("==================== TEST %v, VERSION %v ====================\n", i, version)
		expectedNumUTXOs := 1 + common.RandInt()%29
		Logger.Printf("ExpectedNumUTXOs: %v, minInitial: %v\n", expectedNumUTXOs, minInitialUTXOs)

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
		Logger.Printf("minInitialUTXOs: %v\n", len(utxo))

		balance, err := getBalanceByVersion(privateKey, common.PRVIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		Logger.Printf("oldBalance: %v\n", balance)

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
		Logger.Printf("numUTXOs: %v\n", len(utxo))
		balance, err = getBalanceByVersion(privateKey, common.PRVIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		Logger.Printf("newBalance: %v\n", balance)

		minInitialUTXOs = 100 + common.RandInt()%(MaxUTXO-100)
		Logger.Printf("==================== FINISHED TEST %v ====================\n\n", i)
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
	MaxUTXO := 200
	minInitialUTXOs := 100 + common.RandInt()%(MaxUTXO-100)
	numThreads := 20
	for i := 0; i < numTests; i++ {
		version := int8(1)
		Logger.Printf("==================== TEST %v, VERSION %v ====================\n", i, version)
		expectedNumUTXOs := 1 + common.RandInt()%29
		Logger.Printf("ExpectedNumUTXOs: %v, minInitial: %v\n", expectedNumUTXOs, minInitialUTXOs)

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
		Logger.Printf("minInitialUTXOs: %v\n", len(utxo))

		balance, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		Logger.Printf("oldBalance: %v\n", balance)

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
		Logger.Printf("numUTXOs: %v\n", len(utxo))

		balance, err = getBalanceByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		Logger.Printf("newBalance: %v\n", balance)

		minInitialUTXOs = 100 + common.RandInt()%(MaxUTXO-100)
		Logger.Printf("==================== FINISHED TEST %v ====================\n\n", i)
	}
}

func TestIncClient_ConsolidateTokenV2s(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	//masterPrivateKey := "112t8roafGgHL1rhAP9632Yef3sx5k8xgp8cwK4MCJsCL1UWcxXvpzg97N4dwvcD735iKf31Q2ZgrAvKfVjeSUEvnzKJyyJD3GqqSZdxN4or"
	privateKey := "11111117yu4WAe9fiqmRR4GTxocW6VUKD4dB58wHFjbcQXeDSWQMNyND6Ms3x136EfGcfL7rk3L83BZBzUJLSczmmNi1ngra1WW5Wsjsu5P"
	tokenIDStr := "f3e586e281d275ea2059e35ae434d0431947d2b49466b6d2479808378268f822"
	MaxUTXO := 200
	minInitialUTXOs := 100 + common.RandInt()%(MaxUTXO-100)
	numThreads := 20
	for i := 0; i < numTests; i++ {
		version := int8(2)
		Logger.Printf("==================== TEST %v, VERSION %v ====================\n", i, version)
		Logger.Printf("ExpectedNumUTXOs: %v, minInitial: %v\n", maxUTXOsAfterConsolidated, minInitialUTXOs)

		//// preparing UTXOs
		//err = prepareUTXOs(masterPrivateKey, privateKey, tokenIDStr, minInitialUTXOs, version)
		//if err != nil {
		//	panic(err)
		//}

		// checking UTXOs
		utxo, _, err := ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		if len(utxo) < minInitialUTXOs {
			panic(fmt.Sprintf("require at least %v UTXOs, got %v", minInitialUTXOs, len(utxo)))
		}
		Logger.Printf("minInitialUTXOs: %v\n", len(utxo))

		balance, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		Logger.Printf("oldBalance: %v\n", balance)

		// consolidating UTXOs
		_, err = ic.ConsolidateTokenV2s(privateKey, tokenIDStr, numThreads)
		if err != nil {
			panic(err)
		}

		utxo, _, err = ic.getUTXOsListByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		if len(utxo) > maxUTXOsAfterConsolidated {
			panic(fmt.Sprintf("require at most %v UTXOs, got %v", maxUTXOsAfterConsolidated, len(utxo)))
		}
		Logger.Printf("numUTXOs: %v\n", len(utxo))

		newBalance, err := getBalanceByVersion(privateKey, tokenIDStr, uint8(version))
		if err != nil {
			panic(err)
		}
		Logger.Printf("newBalance: %v\n", newBalance)
		if newBalance != balance {
			panic(fmt.Errorf("expect newBalance to be %v, got %v", balance, newBalance))
		}

		minInitialUTXOs = 100 + common.RandInt()%(MaxUTXO-100)
		Logger.Printf("==================== FINISHED TEST %v ====================\n\n", i)
	}
}
