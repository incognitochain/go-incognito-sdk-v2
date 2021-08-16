package incclient

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"io/ioutil"
	"os"
)

type cachedOutCoins struct {
	Data map[uint64]jsonresult.ICoinInfo `json:"Data"`
}

func NewCachedOutCoins() *cachedOutCoins {
	data := make(map[uint64]jsonresult.ICoinInfo)

	return &cachedOutCoins{Data: data}
}

func (co *cachedOutCoins) UnmarshalJSON(data []byte) error {
	tmpRes := make(map[string]interface{})
	err := json.Unmarshal(data, &tmpRes)
	if err != nil {
		return err
	}

	tmpData, ok := tmpRes["Data"]
	if !ok {
		return fmt.Errorf("`Data` not found")
	}
	jsb, _ := json.Marshal(tmpData)

	cachedData := make(map[uint64]string)
	err = json.Unmarshal(jsb, &cachedData)

	outCoinData := make(map[uint64]jsonresult.ICoinInfo)
	for idx, encodedOutCoin := range cachedData {
		rawOutCoins, _, err := base58.Base58Check{}.Decode(encodedOutCoin)
		if err != nil {
			return err
		}
		tmpCoinInfo, err := coin.NewCoinFromByte(rawOutCoins)
		if err != nil {
			return err
		}
		coinInfo, ok := tmpCoinInfo.(jsonresult.ICoinInfo)
		if !ok {
			return fmt.Errorf("cannot parse coin as a ICoinInfo")
		}

		outCoinData[idx] = coinInfo
	}

	co.Data = outCoinData

	return nil
}

// tokenCache keeps track of the cached coins of a specific tokenID.
type tokenCache struct {
	LatestIndex uint64         `json:"LatestIndex"`
	OutCoins    cachedOutCoins `json:"OutCoins"`
}

func newTokenCache() *tokenCache {
	outCoins := NewCachedOutCoins()

	return &tokenCache{
		LatestIndex: 0,
		OutCoins:    *outCoins,
	}
}

func (tc *tokenCache) UnmarshalJSON(data []byte) error {
	var tmpRes map[string]interface{}
	err := json.Unmarshal(data, &tmpRes)
	if err != nil {
		return err
	}

	var latestIndex float64
	tmpIdx, ok := tmpRes["LatestIndex"]
	if !ok {
		return fmt.Errorf("cannot parse `LatestIndex`")
	}
	latestIndex, ok = tmpIdx.(float64)
	if !ok {
		return fmt.Errorf("expect the `LatestIndex` to be an `uint64`, got %v", tmpIdx)
	}

	OutCoins, ok := tmpRes["OutCoins"]
	if !ok {
		return fmt.Errorf("cannot parse `OutCoins`")
	}
	jsb, err := json.Marshal(OutCoins)
	if err != nil {
		return err
	}

	var tmpOutCoins = NewCachedOutCoins()
	err = json.Unmarshal(jsb, &tmpOutCoins)
	if err != nil {
		return err
	}

	tc.OutCoins = *tmpOutCoins
	tc.LatestIndex = uint64(latestIndex)

	return nil
}

// // tokenCache keeps track of the cached coins for an account.
type accountCache struct {
	// the ota secret key of this account
	OtaKey string `json:"OtaKey"`

	// a mapping from tokenIDs to theirs cached utxo.
	CachedTokens map[string]*tokenCache `json:"CachedTokens"`
}

func newAccountCache(otaKey string) *accountCache {
	cachedTokens := make(map[string]*tokenCache, 0)

	return &accountCache{
		OtaKey:       otaKey,
		CachedTokens: cachedTokens,
	}
}

func (ac accountCache) bytes() ([]byte, error) {
	return json.MarshalIndent(ac, "", "\t")
}

func (ac *accountCache) fromBytes(data []byte) error {
	err := json.Unmarshal(data, &ac)
	if err != nil {
		return err
	}

	return nil
}

// update re-updates the cached Data.
func (ac *accountCache) update(tokenIDStr string, latestIndex uint64, outCoins cachedOutCoins) map[uint64]jsonresult.ICoinInfo {
	updatedRecords := make(map[uint64]jsonresult.ICoinInfo)

	tokenCached := ac.CachedTokens[tokenIDStr]
	if tokenCached == nil {
		tokenCached = &tokenCache{
			LatestIndex: latestIndex,
			OutCoins:    outCoins,
		}
		updatedRecords = tokenCached.OutCoins.Data
	} else if len(outCoins.Data) != 0 {
		Logger.Printf("Adding %v OutCoins to cached %v, LatestIndex %v\n", len(outCoins.Data), tokenIDStr, latestIndex)
		for idx, outCoin := range outCoins.Data {
			if _, ok := tokenCached.OutCoins.Data[idx]; !ok {
				tokenCached.OutCoins.Data[idx] = outCoin
				updatedRecords[idx] = outCoin
			}
		}
	}

	tokenCached.LatestIndex = latestIndex
	ac.CachedTokens[tokenIDStr] = tokenCached

	Logger.Printf("Updated %v records for token %v\n", len(updatedRecords), tokenIDStr)
	Logger.Printf("Current cached size for token %v: %v\n", tokenIDStr, len(tokenCached.OutCoins.Data))
	return updatedRecords
}

// store stores cached output coins to file.
func (ac *accountCache) store(cacheDirectory string) error {
	if cacheDirectory == "" {
		cacheDirectory = defaultCacheDirectory
	}

	filePath := fmt.Sprintf("%v/%v", cacheDirectory, ac.OtaKey)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			Logger.Println(err)
		}
	}()

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	dataToWrite, err := ac.bytes()
	if err != nil {
		return err
	}

	_, err = f.Write(dataToWrite)
	if err != nil {
		return err
	}
	Logger.Printf("Saved utxoCache for OtaKey %v successfully!\n", ac.OtaKey)
	return nil
}

// load re-loads cached data from file.
func (ac *accountCache) load(cacheDirectory string) error {
	if cacheDirectory == "" {
		cacheDirectory = defaultCacheDirectory
	}

	filePath := fmt.Sprintf("%v/%v", cacheDirectory, ac.OtaKey)
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			Logger.Println(err)
		}
	}()

	rawData, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = ac.fromBytes(rawData)
	if err != nil {
		return err
	}

	Logger.Printf("Loaded utxoCache for OtaKey %v successfully!\n", ac.OtaKey)
	return nil
}

func (ac *accountCache) updateAllTokens(latestIndex uint64, data cachedOutCoins, rawAssetTags map[string]*common.Hash) error {
	// update for general confidential assets.
	updateRecords := ac.update(common.ConfidentialAssetID.String(), latestIndex, data)

	// update for each token
	w, err := wallet.Base58CheckDeserialize(ac.OtaKey)
	if err != nil {
		return err
	}
	keySet := w.KeySet
	for idx, coinData := range updateRecords {
		tmpCoinV2, ok := coinData.(*coin.CoinV2)
		if !ok {
			return fmt.Errorf("cannot parse coin %v as a CoinV2", idx)
		}
		tokenId, err := tmpCoinV2.GetTokenId(&keySet, rawAssetTags)
		if err != nil {
			return err
		}

		tokenCached := ac.CachedTokens[tokenId.String()]
		if tokenCached == nil {
			tokenCached = newTokenCache()
			tokenCached.LatestIndex = latestIndex
		}
		if _, ok = tokenCached.OutCoins.Data[idx]; !ok {
			tokenCached.OutCoins.Data[idx] = coinData
			tokenCached.LatestIndex = latestIndex
		}

		ac.CachedTokens[tokenId.String()] = tokenCached
	}

	return nil
}
