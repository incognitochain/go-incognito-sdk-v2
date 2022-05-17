package incclient

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

//UTILS
const (
	pTokenID = "c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0"
	pEthID   = "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
)

var ic *IncClient

func initClients() error {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		return fmt.Errorf("cannot init new incognito client")
	}

	return nil
}

//END UTILS

//TEST FUNCTIONS
func TestIncClient_ShieldETH(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	//Incognito keys
	privateKey := ""

	oldBalance, err := ic.GetBalance(privateKey, pEthID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current balance of token %v: %v\n", pEthID, oldBalance)

	ethTxHash := "" // input the ETH transaction hash here.
	fmt.Printf("Start shielding eth...")

	ethProof, pETHAmount, err := ic.GetEVMDepositProof(ethTxHash)
	if err != nil {
		panic(err)
	}

	txHashStr, err := ic.CreateAndSendIssuingEVMRequestTransaction(privateKey, pEthID, *ethProof)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finish shielding: %v\n", txHashStr)
	time.Sleep(10 * time.Second)

	fmt.Printf("Check shielding status\n")
	for {
		status, err := ic.CheckShieldStatus(txHashStr)
		if err != nil {
			panic(err)
		}
		if status == 1 || status == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		if status == 2 {
			fmt.Printf("Shielding accepted, start checking balance\n")
			break
		} else {
			panic("Shield rejected!")
		}
	}

	for {
		newBalance, err := ic.GetBalance(privateKey, pEthID)
		if err != nil {
			panic(err)
		}
		updatedAmount := newBalance - oldBalance
		if updatedAmount != 0 {
			if updatedAmount != pETHAmount {
				panic(fmt.Sprintf("expected %v, got %v\n", pETHAmount, updatedAmount))
			}
			fmt.Printf("Balance updated!\nnewBalance %v, increasedAmount %v, ethAmount %v\n", newBalance, updatedAmount, pETHAmount)
			break
		}
		fmt.Printf("Balance not updated, sleeping for more...\n")
		time.Sleep(5 * time.Second)
	}
}

func TestIncClient_ShieldERC20(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	//Incognito keys
	privateKey := ""

	oldBalance, err := ic.GetBalance(privateKey, pTokenID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current balance of token %v: %v\n", pTokenID, oldBalance)

	fmt.Printf("Start shielding token...\n")

	ethTxHash := "" // an ETH transaction
	ethProof, _, err := ic.GetEVMDepositProof(ethTxHash)
	if err != nil {
		panic(err)
	}

	txHashStr, err := ic.CreateAndSendIssuingEVMRequestTransaction(privateKey, pTokenID, *ethProof)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finish shielding: %v\n", txHashStr)
	time.Sleep(10 * time.Second)

	fmt.Printf("Check shielding status\n")
	for {
		status, err := ic.CheckShieldStatus(txHashStr)
		if err != nil {
			panic(err)
		}
		if status == 1 || status == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		if status == 2 {
			fmt.Printf("Shielding accepted, start checking balance\n")
			break
		} else {
			panic(fmt.Sprintf("Shield rejected, status: %v\n", status))
		}
	}
}

func TestIncClient_UnShieldETH(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	remoteAddr := "" // an ETH address
	burnedAmount := uint64(50000000)
	burnedTxHash, err := ic.CreateAndSendBurningRequestTransaction(privateKey, remoteAddr, pEthID, burnedAmount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wait for tx %v to be confirmed\n", burnedTxHash)
	for {
		isInBlock, err := ic.CheckTxInBlock(burnedTxHash)
		if err != nil {
			panic(err)
		}

		if !isInBlock {
			fmt.Printf("Tx %v is currently in mempool\n", burnedTxHash)
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("Tx %v is in block\n", burnedTxHash)
			fmt.Printf("Sleep 40 seconds for getting burning proof\n")
			time.Sleep(40 * time.Second)
			break
		}
	}

	fmt.Printf("Start to retrieve the burning proof\n")
	burningProofResult, err := ic.GetBurnProof(burnedTxHash)
	if err != nil {
		panic(err)
	}

	burnProof, err := DecodeBurnProof(burningProofResult)
	if err != nil {
		panic(err)
	}

	jsb, _ := json.Marshal(burnProof)

	fmt.Printf("Burn proof from Incog: %v\n", string(jsb))

	fmt.Printf("Finish getting the burning proof\n")
}

func TestIncClient_UnShieldERC20(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	tokenIDStr := "c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0" //incognito tokenID for pDAI
	remoteAddr := ""                                                                 // an ETH address
	burnedAmount := uint64(10000000)
	burnedTxHash, err := ic.CreateAndSendBurningRequestTransaction(privateKey, remoteAddr, tokenIDStr, burnedAmount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wait for tx %v to be confirmed\n", burnedTxHash)
	for {
		isInBlock, err := ic.CheckTxInBlock(burnedTxHash)
		if err != nil {
			panic(err)
		}

		if !isInBlock {
			fmt.Printf("Tx %v is currently in mempool\n", burnedTxHash)
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("Tx %v is in block\n", burnedTxHash)
			fmt.Printf("Sleep 40 seconds for getting burning proof\n")
			time.Sleep(40 * time.Second)
			break
		}
	}

	fmt.Printf("Start to retrieve the burning proof\n")
	burningProofResult, err := ic.GetBurnProof(burnedTxHash)
	if err != nil {
		panic(err)
	}

	burnProof, err := DecodeBurnProof(burningProofResult)
	if err != nil {
		panic(err)
	}

	jsb, _ := json.Marshal(burnProof)

	fmt.Printf("Burn proof from Incog: %v\n", string(jsb))
	fmt.Printf("Finish getting the burning proof\n")
}

func TestIncClient_GetETHTxReceipt(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	txHash := "0xc400656111f353ef021f3f65711461679e4e1227071411c2789cac762e8948bb"

	receipt, err := ic.GetEVMTxReceipt(txHash)
	if err != nil {
		panic(err)
	}

	jsb, err := json.Marshal(receipt)
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(jsb))
}

func TestIncClient_GetMostRecentETHBlockNumber(t *testing.T) {
	err := initClients()
	if err != nil {
		panic(err)
	}

	blockNum, err := ic.GetMostRecentEVMBlockNumber()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current blockNum: %v\n", blockNum)
}

func TestIncClient_GetBridgeTokens(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	allBridgeTokens, err := ic.GetBridgeTokens()
	if err != nil {
		panic(err)
	}

	for _, token := range allBridgeTokens {
		Logger.Printf("%v, %v, %v, %x\n", token.IsCentralized, token.Network, token.TokenID.String(), token.ExternalTokenID)
	}
}

func TestIncClient_CheckShieldStatus(t *testing.T) {
	var err error
	ic, err = NewMainNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "5b3eb00a96aafe4b477b91d0e78051b58c44385fae07f2a862f3f68195ad7db3"
	status, err := ic.CheckShieldStatus(txHash)
	if err != nil {
		panic(err)
	}
	Logger.Println(status)
}

func TestGenerateTokenID(t *testing.T) {
	tokenID, err := GenerateTokenID("ETH", "USDT")
	if err != nil {
		panic(err)
	}

	Logger.Printf("tokenID: %v\n", tokenID.String())
}

//END TEST FUNCTIONS
