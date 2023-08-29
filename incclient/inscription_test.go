package incclient

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

var ICInscription *IncClient

func initICInscription() error {
	var err error
	ICInscription, err = NewTestNetClientWithCache()
	if err != nil {
		return fmt.Errorf("cannot init new incognito client")
	}

	return nil
}

func TestIncClient_CreateAndSendInscribeRequestTransaction(t *testing.T) {
	err := initICInscription()
	if err != nil {
		panic(err)
	}

	// Input your private key - to pay transaction fee
	privateKey := ""

	path := "./images/a.svg"
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	data := base64.StdEncoding.EncodeToString(fileData)

	txID, err := ICInscription.CreateAndSendInscribeRequestTx(privateKey, data, nil, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Inscribe request sent and waiting for complete\n")
	fmt.Println("TxID: ", txID)
	time.Sleep(50 * time.Second)

	// 	status, err := ICPortal.GetPortalShieldingRequestStatus(shieldID)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if status.Status == 1 {
	// 		fmt.Printf("Shield completed\n")
	// 	} else {
	// 		fmt.Printf("Shield reject with error %v\n", status.Error)
	// 	}
	// }
}
