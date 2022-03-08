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

	depositAddr, err := ic.GeneratePortalShieldingAddress(addr, tokenIDStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("depositAddr: %v\n", depositAddr)

	// SEND SOME PUBLIC TOKENS TO depositAddr, AND THEN RETRIEVE THE SHIELDING PROOF.
	// SEE HOW TO GET THE SHIELD PROOF: https://github.com/incognitochain/incognito-cli/blob/main/portal.go#L77
	depositProof := "DEPOSIT_PROOF"

	txHashStr, err := ic.CreateAndSendPortalShieldTransaction(
		privateKey,
		tokenIDStr,
		addr,
		depositProof,
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
