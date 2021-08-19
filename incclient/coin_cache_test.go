package incclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"math/big"
	"testing"
	"time"
)

func TestIncClient_GetAndCacheOutCoins(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	masterPrivateKey := ""

	privateKey := ""
	tokenIDStr := "75b4045a68b30ab04eb7077a5a972b6ec92fdf24ec3993d685b0c4657dfce948"
	address := PrivateKeyToPaymentAddress(privateKey, -1)

	outCoinKey, err := NewOutCoinKeyFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	outCoinKey.SetReadonlyKey("")

	testTokenIDStr := tokenIDStr
	for i := 0; i < 100; i++ {
		Logger.Printf("TEST %v\n", i)

		isPRV := (common.RandInt() % 2) == 1
		Logger.Printf("isPRV %v\n", isPRV)
		if isPRV {
			testTokenIDStr = common.PRVIDStr
		} else {
			testTokenIDStr = tokenIDStr
		}

		// send some token to the designated address
		addrList := make([]string, 0)
		amtList := make([]uint64, 0)
		numOutCoins := 1 + common.RandInt()%9
		Logger.Printf("#numOutCoins: %v\n", numOutCoins)
		for i := 0; i < numOutCoins; i++ {
			amount := 1 + common.RandUint64()%50
			addrList = append(addrList, address)
			amtList = append(amtList, amount)
		}

		var txHash string
		if isPRV {
			txHash, err = ic.CreateAndSendRawTransaction(
				masterPrivateKey,
				addrList,
				amtList,
				2, nil)
		} else {
			txHash, err = ic.CreateAndSendRawTokenTransaction(
				masterPrivateKey,
				addrList,
				amtList,
				testTokenIDStr, 2, nil)
		}

		if err != nil {
			panic(err)
		}
		Logger.Printf("TxHash %v\n", txHash)

		err = waitingCheckTxInBlock(txHash)
		if err != nil {
			panic(err)
		}
		time.Sleep(40 * time.Second)

		start := time.Now()
		secondOutCoins, secondIndices, err := ic.GetAndCacheOutCoins(outCoinKey, testTokenIDStr)
		if err != nil {
			panic(err)
		}
		Logger.Printf("GetOutputCoinsFromLocalCache time %v\n", time.Since(start).Seconds())

		start = time.Now()
		firstOutCoins, firstIndices, err := ic.GetOutputCoins(outCoinKey, testTokenIDStr, 0)
		if err != nil {
			panic(err)
		}
		Logger.Printf("GetOutputCoins time %v\n", time.Since(start).Seconds())

		isEqual, err := compareUTXOs(firstOutCoins, secondOutCoins, firstIndices, secondIndices)
		if err != nil {
			panic(err)
		}

		Logger.Println("isEqual", isEqual)
		Logger.Printf("FINISHED TEST %v\n\n", i)
		time.Sleep(10 * time.Second)
	}
}

func compareUTXOs(firstOutCoins, secondOutCoins []jsonresult.ICoinInfo, firstIndices, secondIndices []*big.Int) (bool, error) {
	if len(firstIndices) != len(secondIndices) {
		return false, fmt.Errorf("idx lengths mismatch: %v != %v", len(firstIndices), len(secondIndices))
	}

	if len(firstOutCoins) != len(secondOutCoins) {
		return false, fmt.Errorf("outCoin lengths mismatch: %v != %v", len(firstOutCoins), len(secondOutCoins))
	}

	if len(firstOutCoins) != len(firstIndices) {
		return false, fmt.Errorf("outCoin and Idx lengths mismatch: %v != %v", len(firstOutCoins), len(firstIndices))
	}

	type outCoinIdx struct {
		outCoin jsonresult.ICoinInfo
		idx     *big.Int
	}
	mapOutCoins := make(map[string]outCoinIdx)
	for i, outCoin := range firstOutCoins {
		cmtStr := base58.Base58Check{}.Encode(outCoin.GetCommitment().ToBytesS(), 0x00)
		mapOutCoins[cmtStr] = outCoinIdx{
			outCoin: outCoin,
			idx:     firstIndices[i],
		}
	}

	for i, outCoin := range secondOutCoins {
		cmtStr := base58.Base58Check{}.Encode(outCoin.GetCommitment().ToBytesS(), 0x00)
		if res, ok := mapOutCoins[cmtStr]; !ok {
			return false, fmt.Errorf("expect an output coin with commitment %v but get none", cmtStr)
		} else {
			if secondIndices[i].Uint64() != res.idx.Uint64() {
				return false, fmt.Errorf("[%v] expect Idx %v, got %v", cmtStr, secondIndices[i].Uint64(), res.idx.Uint64())
			}

			firstJsb, _ := json.Marshal(res.outCoin)
			secondJsb, _ := json.Marshal(outCoin)
			if !bytes.Equal(firstJsb, secondJsb) {
				Logger.Printf("firstJsb: %v\n", string(firstJsb))
				Logger.Printf("secondJsb: %v\n", string(secondJsb))
				return false, fmt.Errorf("content mismatch")
			}

		}
	}

	return true, nil
}
