package main

import (
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
	addr := "PAYMENT_ADDRESS"
	tokenIDStr := "PORTAL_TOKEN"
	shieldProof := "SHIELD_PROOF"

	txHashStr, err := ic.CreateAndSendPortalShieldTransaction(
		privateKey,
		tokenIDStr,
		addr,
		shieldProof,
		nil, nil,
	)
	fmt.Printf("TxHash: %v\n", txHashStr)

	time.Sleep(10 * time.Second)

	fmt.Printf("check shielding status\n")
	for {
		status, err := ic.GetPortalShieldingRequestStatus(txHashStr)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Printf("shielding status: %v\n", status)
		break
	}
}
