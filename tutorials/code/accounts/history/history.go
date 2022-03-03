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
	client, err := incclient.NewMainNetClient()
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

	err = incclient.SaveTxHistory(h, "history.csv")
	if err != nil {
		log.Fatal(err)
	}
}

// GetHistoryFaster helps retrieve the history faster by running parallel workers.
func GetHistoryFaster() {
	// For main-net
	client, err := incclient.NewMainNetClient()
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

	err = incclient.SaveTxHistory(h, "history.csv")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// comment one of these functions.
	//GetHistory()
	GetHistoryFaster()
}
