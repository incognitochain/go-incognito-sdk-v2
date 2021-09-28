package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UTILS
var PRVPeggingIncClient *IncClient
var tokenId = common.PRVIDStr

func initPRVPeggingIncClient() error {
	var err error
	PRVPeggingIncClient, err = NewTestNetClient()
	if err != nil {
		return fmt.Errorf("cannot init new incognito client")
	}

	return nil
}

// END UTILS

// TEST FUNCTIONS

type TestCaseShieldPRVPegging struct {
	externalTxID string
	isBSC        bool
	shieldAmt    uint64
}

func TestIncClient_ShieldPRVPegging(t *testing.T) {
	// init testcases
	// INPUT YOUR TESTCASE
	tcs := []TestCaseShieldPRVPegging{
		{
			externalTxID: "",
			isBSC:        false,
			shieldAmt:    uint64(0),
		},
		{
			externalTxID: "",
			isBSC:        true,
			shieldAmt:    uint64(0),
		},
	}

	// Incognito keys
	privateKey := ""

	err := initPRVPeggingIncClient()
	if err != nil {
		panic(err)
	}

	for _, tc := range tcs {
		oldBalance, err := PRVPeggingIncClient.GetBalance(privateKey, tokenId)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Current balance of token %v: %v\n", tokenId, oldBalance)

		fmt.Printf("Start shielding token...\n")

		externalTxHash := tc.externalTxID // an ETH transaction
		ethProof, _, err := PRVPeggingIncClient.GetEVMDepositProof(externalTxHash)
		if err != nil {
			panic(err)
		}

		txHashStr, err := PRVPeggingIncClient.CreateAndSendIssuingPRVPeggingRequestTransaction(privateKey, *ethProof, tc.isBSC)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Finish shielding: %v\n", txHashStr)
		time.Sleep(10 * time.Second)

		fmt.Printf("Check shielding status\n")
		for {
			status, err := PRVPeggingIncClient.CheckShieldStatus(txHashStr)
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
		newBalance, err := PRVPeggingIncClient.GetBalance(privateKey, tokenId)
		if err != nil {
			panic(err)
		}
		fmt.Printf("New balance of token %v: %v\n", tokenId, newBalance)

		assert.Equal(t, oldBalance+tc.shieldAmt, newBalance)
	}
}

type TestCaseUnShieldPRVPegging struct {
	isBSC           bool
	unshieldAmt     uint64
	externalAddress string
}

func TestIncClient_UnShieldPRVPegging(t *testing.T) {
	// init testcases
	// INPUT YOUR TESTCASE
	tcs := []TestCaseUnShieldPRVPegging{
		{
			isBSC:           false,
			unshieldAmt:     2 * 1e9,
			externalAddress: "0xF91cEe2DE943733e338891Ef602c962eF4D7Eb81",
		},
		{
			isBSC:           true,
			unshieldAmt:     1 * 1e9,
			externalAddress: "0xF91cEe2DE943733e338891Ef602c962eF4D7Eb81",
		},
	}
	privateKey := ""

	err := initPRVPeggingIncClient()
	if err != nil {
		panic(err)
	}

	for _, tc := range tcs {
		burnedTxHash, err := PRVPeggingIncClient.CreateAndSendBurningPRVPeggingRequestTransaction(
			privateKey, tc.externalAddress, tc.unshieldAmt)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Wait for tx %v to be confirmed\n", burnedTxHash)
		for {
			isInBlock, err := PRVPeggingIncClient.CheckTxInBlock(burnedTxHash)
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
		burningProofResult, err := PRVPeggingIncClient.GetBurnPRVPeggingProof(burnedTxHash)
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
}
