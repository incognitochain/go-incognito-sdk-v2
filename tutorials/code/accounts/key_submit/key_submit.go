package main

import (
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	privateOTAKey := "14yBChbLDg42noBQHonDR5mj3FMD9CPCCNfPoa68jeE8bE2LsyfCKcNkgupEsm6pW4BZFnDHmay9XjDGE1iTaTEcEpN7UUaPoU344g2"

	// for regular cache
	err = client.SubmitKey(privateOTAKey)
	if err != nil {
		log.Fatal(err)
	}

	// for enhanced cache
	accessToken := "0c3d46946bbf9339c8213dd7f6c640ed643003bdc056a5b68e7e80f5ef5aa0dd"
	fromHeight := uint64(0)
	isReset := true
	err = client.AuthorizedSubmitKey(privateOTAKey, accessToken, fromHeight, isReset)
	if err != nil {
		log.Fatal(err)
	}
}
