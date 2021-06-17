---
description: Tutorial on how to call any RPC requests with the go-sdk
---
To make an RPC request to a full-node, use the function [`NewRPCCall`](../../../incclient/general.go) in the [`incclient`](../../../incclient) package. It returns the raw reponse in bytes. See the following example.

## Example
```go
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

```
