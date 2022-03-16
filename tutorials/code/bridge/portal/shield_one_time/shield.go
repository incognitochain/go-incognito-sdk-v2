package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
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
	depositKeyIndex := uint64(0)

	// Generate depositKey from the privateKey and a depositKeyIndex
	depositKey, err := ic.GenerateDepositKeyFromPrivateKey(privateKey, tokenIDStr, depositKeyIndex)
	if err != nil {
		log.Fatal(err)
	}
	jsb, _ := json.Marshal(depositKey)
	fmt.Printf("depositKey: %v\n", string(jsb))

	// Generate an OTAReceiver and sign it
	tokenID, _ := common.Hash{}.NewHashFromStr(tokenIDStr)
	w, err := wallet.Base58CheckDeserialize(addr)
	if err != nil {
		log.Fatal(err)
	}
	otaReceivers, err := incclient.GenerateOTAReceivers([]common.Hash{*tokenID}, w.KeySet.PaymentAddress)
	if err != nil {
		log.Fatal(err)
	}
	otaReceiver := otaReceivers[*tokenID]
	signature, err := incclient.SignDepositData(depositKey, otaReceiver.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// Generate depositAddr from the PublicKey of the depositKey
	depositPubKeyStr := base58.Base58Check{}.NewEncode(depositKey.PublicKey, 0)
	depositAddr, err := ic.GeneratePortalShieldingAddress(depositPubKeyStr, tokenIDStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("depositAddr: %v\n", depositAddr)

	// SEND SOME PUBLIC TOKENS TO depositAddr, AND THEN RETRIEVE THE SHIELDING PROOF.
	// SEE HOW TO GET THE SHIELD PROOF: https://github.com/incognitochain/incognito-cli/blob/main/portal.go#L77
	depositProof := "DEPOSIT_PROOF"

	depositParam := incclient.PortalDepositParams{
		TokenID:       tokenIDStr,
		ShieldProof:   depositProof,
		DepositPubKey: depositPubKeyStr,
		Receiver:      otaReceiver.String(),
		Signature:     base58.Base58Check{}.Encode(signature, 0),
	}

	// Create the shielding transaction
	txHashStr, err := ic.CreateAndSendPortalShieldTransactionWithDepositKey(
		privateKey,
		depositParam,
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
