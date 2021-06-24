package main

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

// GetHistory retrieves the history of an account in a normal way.
// It is suitable for accounts that have a few transactions.
func GetHistory() {
	// For main-net
	client, err := incclient.NewIncClient("https://beta-fullnode.incognito.org/fullnode", incclient.MainNetETHHost, 1)
	if err != nil {
		log.Fatal(err)
	}

	tokenIDStr := common.PRVIDStr    // input the tokenID in which you want to retrieve the history of.
	privateKey := "YOUR_PRIVATE_KEY" // input your private key here

	// get the history in a normal way.
	h, err := client.GetTxHistoryV1(privateKey, tokenIDStr)
	if err != nil {
		log.Fatal(err)
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

// GetHistoryFaster helps retrieve the history faster by running parallel workers.
func GetHistoryFaster() {
	// For main-net
	client, err := incclient.NewIncClient("https://beta-fullnode.incognito.org/fullnode", incclient.MainNetETHHost, 1)
	if err != nil {
		log.Fatal(err)
	}

	tokenIDStr := common.PRVIDStr    // input the tokenID in which you want to retrieve the history of.
	privateKey := "YOUR_PRIVATE_KEY" // input your private key here

	numWorkers := 15
	p := incclient.NewTxHistoryProcessor(client, numWorkers)

	h, err := p.GetTokenHistory(privateKey, tokenIDStr)
	if err != nil {
		log.Fatal(err)
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
}

func main() {
	// comment one of these functions.
	GetHistory()
	GetHistoryFaster()
}
