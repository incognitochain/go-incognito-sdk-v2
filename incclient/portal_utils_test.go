package incclient

import (
	"encoding/json"
	"testing"
)

func TestIncClient_GetPortalUnShieldingRequestStatus(t *testing.T) {
	var err error
	ic, err = NewMainNetClientWithCache()
	if err != nil {
		panic(err)
	}

	unShieldID := "decc21f35ed8f9edc5167e1f7b3622e46f95216d0218fe2991d5cf1e4e491511"
	status, err := ic.GetPortalUnShieldingRequestStatus(unShieldID)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(jsb))
}
