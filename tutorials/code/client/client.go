package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	incClient, err := incclient.NewMainNetClient() // or use incclient.NewMainNetClientWithCache()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Let's get yourself into the Incognito network!")
	_ = incClient
}
