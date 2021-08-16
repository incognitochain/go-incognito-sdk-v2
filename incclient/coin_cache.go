package incclient

import (
	"bytes"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"
)

var batchSize = 5000

// utxoCache implements a simple UTXO cache for the incclient.
type utxoCache struct {
	// indicator of whether the cache is running
	isRunning bool

	// the directory where the cached is store.
	cacheDirectory string

	// the mapping from otaKeys to their cached UTXOs.
	cachedData map[string]*accountCache

	// a simple mutex
	mtx *sync.Mutex
}

// newUTXOCache creates a new utxoCache instance.
func newUTXOCache(cacheDirectory string) *utxoCache {
	cachedData := make(map[string]*accountCache)
	mtx := new(sync.Mutex)

	return &utxoCache{
		cacheDirectory: cacheDirectory,
		cachedData:     cachedData,
		mtx:            mtx,
	}
}

func (uc *utxoCache) start() {
	err := uc.load()
	if err != nil {
		log.Fatal(err)
	}

	uc.isRunning = true
}

// saveAndStop saves the current cache and stop.
func (uc *utxoCache) saveAndStop() error {
	uc.isRunning = false
	err := uc.save()
	if err != nil {
		return err
	}
	return nil
}

func (uc *utxoCache) save() error {
	Logger.Println("Storing cached data...")
	var err error
	uc.mtx.Lock()
	for _, cachedData := range uc.cachedData {
		err = cachedData.store(uc.cacheDirectory)
		if err != nil {
			return err
		}
	}
	uc.mtx.Unlock()

	return nil
}

func (uc *utxoCache) load() error {
	files, err := ioutil.ReadDir(uc.cacheDirectory)
	if err != nil {
		return err
	}

	cachedData := make(map[string]*accountCache)

	uc.mtx.Lock()
	for _, f := range files {
		fileNameSplit := strings.Split(f.Name(), "/")
		otaKey := fileNameSplit[len(fileNameSplit)-1]
		ac := newAccountCache(otaKey)
		err = ac.load(uc.cacheDirectory)
		if err != nil {
			return err
		}
		cachedData[otaKey] = ac
	}
	uc.cachedData = cachedData
	uc.mtx.Unlock()

	Logger.Printf("Loading cache successfully!\n")
	Logger.Printf("Current cache size: %v\n", len(uc.cachedData))

	return nil
}

func (uc *utxoCache) getCachedAccount(otaKey string) *accountCache {
	uc.mtx.Lock()
	ac := uc.cachedData[otaKey]
	uc.mtx.Unlock()
	return ac
}

func (uc *utxoCache) addAccount(otaKey string, cachedAccount *accountCache) {
	uc.mtx.Lock()
	uc.cachedData[otaKey] = cachedAccount
	uc.mtx.Unlock()
}

func (client *IncClient) getOutCoinV1(outCoinKey *rpc.OutCoinKey, tokenIDStr string) ([]jsonresult.ICoinInfo, error) {
	otaKey := outCoinKey.OtaKey()
	outCoinKey.SetOTAKey("") // set this to empty so that the full-node only query v1 output coins.

	var tmp jsonresult.ListOutputCoins
	responseInBytes, err := client.rpcServer.GetListOutputCoinsByRPCV1(outCoinKey, tokenIDStr, 0)
	if err != nil {
		return nil, err
	}
	err = rpchandler.ParseResponse(responseInBytes, &tmp)
	if err != nil {
		return nil, err
	}

	outCoins, _, err := ParseCoinFromJsonResponse(responseInBytes)
	if err != nil {
		return nil, fmt.Errorf("get v1 output coins error: %v", err)
	}

	outCoinKey.SetOTAKey(otaKey)
	return outCoins, nil
}

func (client *IncClient) syncOutCoinV2(outCoinKey *rpc.OutCoinKey, tokenIDStr string) error {
	if tokenIDStr != common.PRVIDStr {
		tokenIDStr = common.ConfidentialAssetID.String()
	}

	shardID, err := GetShardIDFromPaymentAddress(outCoinKey.PaymentAddress())
	if err != nil || shardID == 255 {
		return fmt.Errorf("GetShardIDPaymentAddressKey failed: %v", err)
	}

	w, err := wallet.Base58CheckDeserialize(outCoinKey.OtaKey())
	if err != nil {
		return err
	}
	keySet := w.KeySet
	if keySet.OTAKey.GetOTASecretKey() == nil || keySet.OTAKey.GetPublicSpend() == nil {
		return fmt.Errorf("invalid OTAKey")
	}

	coinLength, err := client.GetOTACoinLengthByShard(shardID, tokenIDStr)
	if err != nil {
		return err
	}
	Logger.Printf("Current OTALength for token %v, shard %v: %v\n", tokenIDStr, shardID, coinLength)

	var cachedAccount *accountCache
	var ok bool
	var cachedToken *tokenCache
	if cachedAccount = client.cache.getCachedAccount(outCoinKey.OtaKey()); cachedAccount == nil {
		Logger.Printf("No cache found, creating a new one...\n")
		cachedAccount = newAccountCache(outCoinKey.OtaKey())
		cachedAccount.CachedTokens = make(map[string]*tokenCache)
		cachedToken = newTokenCache()
	} else if cachedToken, ok = cachedAccount.CachedTokens[tokenIDStr]; !ok {
		cachedToken = newTokenCache()
	}

	res := NewCachedOutCoins()
	burningPubKey := wallet.GetBurningPublicKey()

	start := time.Now()
	currentIndex := cachedToken.LatestIndex
	Logger.Printf("Current LatestIndex for token %v: %v\n", tokenIDStr, currentIndex)
	for currentIndex < coinLength {
		idxList := make([]uint64, 0)

		nextIndex := currentIndex + uint64(batchSize)
		if nextIndex > coinLength - 1 {
			nextIndex = coinLength - 1
		}
		for i := currentIndex; i < nextIndex; i++ {
			idxList = append(idxList, i)
		}
		if len(idxList) == 0 {
			break
		}

		Logger.Printf("Get output coins of indices from %v to %v\n", currentIndex, nextIndex - 1)

		tmpOutCoins, err := client.GetOTACoinsByIndices(shardID, tokenIDStr, idxList)
		if err != nil {
			return err
		}
		found := 0
		for idx, outCoin := range tmpOutCoins {
			if bytes.Equal(outCoin.Bytes(), burningPubKey) {
				continue
			}
			belongs, _ := outCoin.DoesCoinBelongToKeySet(&keySet)
			if belongs {
				res.Data[idx] = outCoin
				found += 1
			}
		}
		Logger.Printf("Found %v output coins (%v) for heights from %v to %v with time %v\n", found, tokenIDStr, currentIndex, nextIndex-1, time.Since(start).Seconds())
		currentIndex = nextIndex
	}

	Logger.Printf("newOutCoins: %v\n", len(res.Data))


	if tokenIDStr == common.PRVIDStr {
		cachedAccount.update(common.PRVIDStr, coinLength, *res)
	} else {
		// update cached data for each token
		if rawAssetTags == nil {
			Logger.Printf("dasdasdasdasda\n")
			rawAssetTags, err = client.GetAllAssetTags()
			if err != nil {
				return err
			}
		}

		err = cachedAccount.updateAllTokens(coinLength, *res, rawAssetTags)
		if err != nil {
			return err
		}
	}
	client.cache.addAccount(outCoinKey.OtaKey(), cachedAccount)
	Logger.Printf("FINISHED SYNCING OUTPUT COINS OF TOKEN %v AFTER %v SECOND\n", tokenIDStr, time.Since(start).Seconds())

	return nil
}

// GetOutputCoinsFromLocalCache returns the list of UTXOs from the client's cache.
// This function should only be called after the cache is initialized.
func (client *IncClient) GetOutputCoinsFromLocalCache(outCoinKey *rpc.OutCoinKey, tokenID string, save ...bool) ([]jsonresult.ICoinInfo, []*big.Int, error) {
	if client.cache == nil || !client.cache.isRunning {
		return nil, nil, fmt.Errorf("utxoCache is not running")
	}

	err := client.syncOutCoinV2(outCoinKey, tokenID)
	if err != nil {
		return nil, nil, err
	}

	outCoins := make([]jsonresult.ICoinInfo, 0)
	indices := make([]*big.Int, 0)

	cachedAccount := client.cache.getCachedAccount(outCoinKey.OtaKey())
	if cachedAccount == nil {
		return nil, nil, fmt.Errorf("otaKey %v has not been cached", outCoinKey.OtaKey())
	}
	cached := cachedAccount.CachedTokens[tokenID]
	if cached != nil {
		for idx, outCoin := range cached.OutCoins.Data {
			outCoins = append(outCoins, outCoin)
			idxBig := new(big.Int).SetUint64(idx)
			indices = append(indices, idxBig)
		}
	} else {
		Logger.Printf("No cached found for tokenID %v\n", tokenID)
	}

	v1OutCoins, err := client.getOutCoinV1(outCoinKey, tokenID)
	if err != nil {
		return nil, nil, err
	}
	v1Count := 0
	for _, v1OutCoin := range v1OutCoins {
		if v1OutCoin.GetVersion() != 1 {
			continue
		}
		outCoins = append(outCoins, v1OutCoin)
		idxBig := new(big.Int).SetInt64(-1)
		indices = append(indices, idxBig)
		v1Count++
	}
	Logger.Printf("Found %v v1 output coins\n", v1Count)

	if save != nil {
		isSaved := save[0]
		if isSaved {
			err = client.cache.save()
			if err != nil {
				panic(err)
			}
		}
	}

	return outCoins, indices, nil
}
