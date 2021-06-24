package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"testing"
	"time"
)

func TestTxHistoryProcessor_GetTxsIn(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""

	p := NewTxHistoryProcessor(ic, 15)

	start := time.Now()
	txIns, err := p.GetTxsIn(privateKey, tokenIDStr, 1)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v\n", len(txIns))

	totalIn := uint64(0)
	for _, txIn := range txIns {
		totalIn += txIn.Amount
		log.Printf("%v\n", txIn.String())
	}
	log.Printf("TotalIn: %v\n", totalIn)

	log.Printf("\nTime elapsed: %v\n", time.Since(start).Seconds())

}

func TestTxHistoryProcessor_GetTxsOut(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""

	p := NewTxHistoryProcessor(ic, 15)

	start := time.Now()
	txOuts, err := p.GetTxsOut(privateKey, tokenIDStr, 1)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v\n", len(txOuts))

	totalOut := uint64(0)
	for _, txOut := range txOuts {
		totalOut += txOut.Amount
		log.Printf("%v\n", txOut.String())
	}
	log.Printf("TotalOut: %v\n", totalOut)

	log.Printf("\nTime elapsed: %v\n", time.Since(start).Seconds())

}

func TestTxHistoryProcessor_GetTokenHistory(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""

	p := NewTxHistoryProcessor(ic, 15)

	start := time.Now()
	h, err := p.GetTokenHistory(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v, #TxsOut: %v\n", len(h.TxInList), len(h.TxOutList))

	totalIn := uint64(0)
	log.Printf("TxsIn\n")
	for _, txIn := range h.TxInList {
		totalIn += txIn.Amount
		log.Printf("%v\n", txIn.String())
	}
	log.Printf("Finished TxsIn\n\n")

	totalOut := uint64(0)
	log.Printf("TxsOut\n")
	for _, txOut := range h.TxOutList {
		totalOut += txOut.Amount
		log.Printf("%v\n", txOut.String())
	}
	log.Printf("Finished TxsOut\n\n")

	balance, err := ic.GetBalance(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}
	log.Printf("currentBalance: %v, totalIn: %v, totalOut: %v\n", balance, totalIn, totalOut)

	log.Printf("\nTime elapsed: %v\n", time.Since(start).Seconds())

}
