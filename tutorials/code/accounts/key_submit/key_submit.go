package main

import (
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

func main() {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}

	privateOTAKey := "14yBChbLDg42noBQHonDR5mj3FMD9CPCCNfPoa68jeE8bE2LsyfCKcNkgupEsm6pW4BZFnDHmay9XjDGE1iTaTEcEpN7UUaPoU344g2"

	// for regular cache
	err = client.SubmitKey(privateOTAKey)
	if err != nil {
		log.Fatal(err)
	}

	// at this point, if you submit again, it will throw an error
	err = client.SubmitKey(privateOTAKey)
	if err != nil {
		log.Println(err) // should throw an error: OTAKey has been submitted and status = 2
	}

	// However, you can override the regular submission by the enhanced mode.
	accessToken := "0c3d46946bbf9339c8213dd7f6c640ed643003bdc056a5b68e7e80f5ef5aa0dd"
	fromHeight := uint64(0)
	isReset := true
	err = client.AuthorizedSubmitKey(privateOTAKey, accessToken, fromHeight, isReset)
	if err != nil {
		log.Fatal(err)
	}
}
