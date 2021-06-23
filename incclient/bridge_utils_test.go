package incclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIncClient_GetEVMTxByHash(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	txHash := "0x392bf44aeac2c395fc4ed7ba425f1fc61b7b62d98a96c2a2d5e22c5ec8cd8f23"

	txDetail, err := ic.GetEVMTxByHash(txHash)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(txDetail, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetEVMBlockByHash(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	blockHash := "0xe886b93d341bb6cf1f4e24a2ffa40c0a6107adb6214e8f7e43fce04d07fc3f1f"

	blockDetail, err := ic.GetEVMBlockByHash(blockHash)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(blockDetail, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetEVMTxReceipt(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	txHash := "0x392bf44aeac2c395fc4ed7ba425f1fc61b7b62d98a96c2a2d5e22c5ec8cd8f23"

	receipt, err := ic.GetEVMTxReceipt(txHash)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(receipt, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsb))
}

func TestIncClient_GetEVMDepositProof(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	txHash := "0x392bf44aeac2c395fc4ed7ba425f1fc61b7b62d98a96c2a2d5e22c5ec8cd8f23"

	depositProof, amount, err := ic.GetEVMDepositProof(txHash)
	if err != nil {
		panic(err)
	}

	fmt.Println(amount)
	fmt.Println(depositProof.BlockNumber(), depositProof.BlockHash().String(), depositProof.TxIdx())
	fmt.Println(depositProof.NodeList())
}

func TestIncClient_GetMostRecentEVMBlockNumber(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	mostRecentBlock, err := ic.GetMostRecentEVMBlockNumber()
	if err != nil {
		panic(err)
	}

	fmt.Printf("mostRecentBlock: %v\n", mostRecentBlock)
}

func TestIncClient_GetEVMTransactionStatus(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	txHash := "0x392bf44aeac2c395fc4ed7ba425f1fc61b7b62d98a96c2a2d5e22c5ec8cd8f23"

	status, err := ic.GetEVMTransactionStatus(txHash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("status: %v\n", status)
}