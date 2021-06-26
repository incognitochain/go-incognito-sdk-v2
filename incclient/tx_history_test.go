package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"testing"
)

func TestIncClient_GetListTxsInV1(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""
	txList, err := ic.GetListTxsInV1(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v\n", len(txList))

	for _, txIn := range txList {
		log.Printf("%v\n", txIn.String())
	}
}

func TestIncClient_GetListTxsInV2(t *testing.T) {
	var err error
	ic, err = NewLocalClient("")
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""
	txList, err := ic.GetListTxsInV2(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v\n", len(txList))

	for _, txIn := range txList {
		log.Printf("%v\n", txIn.String())
	}
}

func TestIncClient_GetListTxsIn(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""
	txList, err := ic.GetListTxsIn(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v\n", len(txList))

	for _, txIn := range txList {
		log.Printf("%v\n", txIn.String())
	}
}

func TestIncClient_GetListTxsOutV1(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""
	txList, err := ic.GetListTxsOutV1(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxOuts: %v\n", len(txList))

	for _, txOut := range txList {
		log.Printf("%v\n", txOut.String())
	}
}

func TestIncClient_GetListTxsOutV2(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""
	txList, err := ic.GetListTxsOutV2(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxOuts: %v\n", len(txList))

	for _, txOut := range txList {
		log.Printf("%v\n", txOut.String())
	}
}

func TestIncClient_GetTxHistoryV1(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""
	h, err := ic.GetTxHistoryV1(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v\n", len(h.TxInList))
	for _, txIn := range h.TxInList {
		log.Printf("%v\n", txIn.String())
	}
	log.Printf("\n#TxOuts: %v\n", len(h.TxOutList))
	for _, txOut := range h.TxOutList {
		log.Printf("%v\n", txOut.String())
	}
}

func TestIncClient_GetTxHistoryV2(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	tokenIDStr := common.PRVIDStr
	privateKey := ""
	h, err := ic.GetTxHistoryV2(privateKey, tokenIDStr)
	if err != nil {
		panic(err)
	}

	log.Printf("#TxIns: %v\n", len(h.TxInList))
	for _, txIn := range h.TxInList {
		log.Printf("%v\n", txIn.String())
	}
	log.Printf("\n#TxOuts: %v\n", len(h.TxOutList))
	for _, txOut := range h.TxOutList {
		log.Printf("%v\n", txOut.String())
	}
}
