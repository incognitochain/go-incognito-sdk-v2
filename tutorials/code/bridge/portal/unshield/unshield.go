package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"time"
)

func main() {
	ic, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := "YOUR_PRIVATE_KEY"
	externalAddr := "EXTERNAL_ADDRESS"
	tokenIDStr := "PORTAL_TOKEN"
	unShieldAmount := uint64(1000000)

	txHash, err := ic.CreateAndSendPortalUnShieldTransaction(
		privateKey, tokenIDStr, externalAddr, unShieldAmount, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TxHash: %v\n", txHash)

	time.Sleep(100 * time.Second)
	status, err := ic.GetPortalUnShieldingRequestStatus(txHash)
	if err != nil {
		log.Fatal(err)
	}

	jsb, _ := json.Marshal(status)
	fmt.Println(string(jsb))
}
