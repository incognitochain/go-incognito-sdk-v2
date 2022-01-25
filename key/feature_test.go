package key

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"testing"
)

func TestOTDepositKey_UnmarshalJSON(t *testing.T) {
	privateKey := common.RandBytes(32)
	tokenID := common.HashH(common.RandBytes(32)).String()

	for index := uint64(0); index < 100000; index++ {
		depositKey, err := GenerateOTDepositKeyFromPrivateKey(privateKey, tokenID, index)
		if err != nil {
			panic(err)
		}
		jsb, err := json.Marshal(depositKey)
		if err != nil {
			panic(err)
		}

		recoveredKey := new(OTDepositKey)
		err = json.Unmarshal(jsb, recoveredKey)
		if err != nil {
			panic(err)
		}

		if recoveredKey.Index != depositKey.Index ||
			!bytes.Equal(recoveredKey.PublicKey, depositKey.PublicKey) ||
			!bytes.Equal(recoveredKey.PrivateKey, depositKey.PrivateKey) {
			panic(fmt.Sprintf("recovered: %v, actual %v\n", recoveredKey, depositKey))
		}
	}
}
