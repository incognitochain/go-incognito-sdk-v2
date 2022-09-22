package coin

import (
	"fmt"
	"testing"
	"time"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

var (
	numTests = 100
)

func TestNewCoinFromPaymentInfo(t *testing.T) {
	for i := 0; i < numTests; i++ {
		prefix := fmt.Sprintf("[TEST %v]", i)
		common.MaxShardNumber = 1 + common.RandInt()%7
		senderShard := common.RandInt() % common.MaxShardNumber
		receiverShard := common.RandInt() % common.MaxShardNumber
		coinType := common.RandInt() % 2
		if coinType == PrivacyTypeMint {
			senderShard = receiverShard
		}
		fmt.Printf("%v STARTED\n", prefix)
		fmt.Printf("%v numShards: %v, senderShard: %v, receiverShard: %v, coinType: %v\n", prefix,
			common.MaxShardNumber, senderShard, receiverShard, coinType)

		w, err := wallet.GenRandomWalletForShardID(byte(receiverShard))
		if err != nil {
			panic(fmt.Sprintf("%v %v", prefix, err))
		}
		paymentInfo := &key.PaymentInfo{
			PaymentAddress: w.KeySet.PaymentAddress,
			Amount:         0,
			Message:        []byte{},
		}

		var coinParam *CoinParams
		if coinType == PrivacyTypeTransfer {
			coinParam = NewTransferCoinParams(paymentInfo, byte(senderShard))
		} else {
			coinParam = NewMintCoinParams(paymentInfo)
		}

		start := time.Now()
		c, err := NewCoinFromPaymentInfo(coinParam)
		if err != nil {
			panic(fmt.Sprintf("%v %v", prefix, err))
		}

		tmpSenderShard, tmpReceiverShard, tmpCoinType, err := DeriveShardInfoFromCoin(c.GetPublicKey().ToBytesS())
		if err != nil {
			panic(fmt.Sprintf("%v %v", prefix, err))
		}

		if tmpSenderShard != senderShard {
			panic(fmt.Sprintf("%v expect senderShard to be %v, got %v", prefix, senderShard, tmpSenderShard))
		}

		if tmpReceiverShard != receiverShard {
			panic(fmt.Sprintf("%v expect receiverShard to be %v, got %v", prefix, receiverShard, tmpReceiverShard))
		}

		if tmpCoinType != coinType {
			panic(fmt.Sprintf("%v expect coinType to be %v, got %v", prefix, coinType, tmpCoinType))
		}

		fmt.Printf("%v FINISHED: %v\n\n", prefix, time.Since(start).Seconds())
	}
}
