package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
)

func main() {
	client, err := incclient.NewTestNet1Client()
	if err != nil {
		log.Fatal(err)
	}

	method := "getshardbeststate"
	params := make([]interface{}, 0)
	params = append(params, 1)

	resp, err := client.NewRPCCall("1.0", method, params, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp))
}
