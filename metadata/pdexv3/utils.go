package pdexv3

import (
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
)

type ReceiverInfo struct {
	Address coin.OTAReceiver `json:"Address"`
	Amount  uint64           `json:"Amount"`
}
