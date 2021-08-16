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

func TestIncClient_GetOutputCoinsFromLocalCache(t *testing.T) {
	var err error
	ic, err = NewMainNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	tokenIDStr := common.PRVIDStr

	outCoinKey, err := NewOutCoinKeyFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	outCoinKey.SetReadonlyKey("")

	for i := 0; i < 2; i++ {
		Logger.Printf("TEST %v\n", i)

		start := time.Now()
		secondOutCoins, secondIndices, err := ic.GetAndCacheOutCoins(outCoinKey, tokenIDStr, true)
		if err != nil {
			panic(err)
		}
		Logger.Printf("GetOutputCoinsFromLocalCache time %v\n", time.Since(start).Seconds())

		start = time.Now()
		firstOutCoins, firstIndices, err := ic.GetOutputCoins(outCoinKey, tokenIDStr, 0)
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
		idx *big.Int
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



