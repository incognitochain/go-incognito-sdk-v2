package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"testing"
	"time"
)

func TestTxHistoryProcessor_GetTxsIn(t *testing.T) {
	var err error
	ic, err = NewMainNetClientWithCache()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""

	p := NewTxHistoryProcessor(ic, 30)

	txIns, err := p.GetTxsIn(privateKey, tokenIDStr, 2)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v\n", len(txIns))

	for _, txIn := range txIns {
		log.Println(txIn.String())
	}

}

func TestTxHistoryProcessor_GetTxsOut(t *testing.T) {
	var err error
	ic, err = NewMainNetClientWithCache()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""

	p := NewTxHistoryProcessor(ic, 15)

	start := time.Now()
	txOuts, err := p.GetTxsOut(privateKey, tokenIDStr, 2)
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
	ic, err = NewMainNetClientWithCache()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""

	p := NewTxHistoryProcessor(ic, 30)

	h, err := p.GetTokenHistory(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	err = SaveTxHistory(h, "tmp.csv")
	if err != nil {
		panic(err)
	}
}
