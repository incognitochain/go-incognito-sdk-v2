package incclient

import (
	"bytes"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	batchSize = 5000

	// MaxGetCoinThreads is the maximum number of threads running simultaneously to retrieve output coins in the cache layer.
	// By default, it is set to the maximum of 4 and the number of CPUs of the running machine.
	MaxGetCoinThreads = int(math.Max(float64(runtime.NumCPU()), 4))
)

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

// getCoinStatus manages the status of a thread for getting output coins by indices.
type getCoinStatus struct {
	data      map[uint64]jsonresult.ICoinInfo
	err       error
	fromIndex uint64
	toIndex   uint64
}

// newUTXOCache creates a new utxoCache instance.
func newUTXOCache(cacheDirectory string) (*utxoCache, error) {
	cachedData := make(map[string]*accountCache)
	mtx := new(sync.Mutex)

	// if the cache directory does not exist, create one.
	if _, err := os.Stat(cacheDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(cacheDirectory, os.ModePerm)
		if err != nil {
			Logger.Printf("make directory %v error: %v\n", cacheDirectory, err)
			return nil, err
		}
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Printf("cacheDirectory: %v/%v\n", currentDir, cacheDirectory)

	return &utxoCache{
		cacheDirectory: cacheDirectory,
		cachedData:     cachedData,
		mtx:            mtx,
	}, nil
}

func (uc *utxoCache) start() {
	uc.isRunning = true
}

// saveAndStop saves the current cache and stops.
func (uc *utxoCache) saveAndStop() error {
	uc.isRunning = false
	err := uc.save()
	if err != nil {
		return err
	}
	return nil
}

// save either backs up the whole cache or a specific otaKey.
// Only the first value of `otaKeys` is processed.
func (uc *utxoCache) save(otaKeys ...string) error {
	Logger.Println("Storing cached data...")

	var otaKeyStr string
	if len(otaKeys) > 0 {
		otaKeyStr = otaKeys[0]
		if otaKeyStr == "" {
			return fmt.Errorf("invalid OTAKey")
		}
	}

	var err error
	uc.mtx.Lock()
	defer func() {
		uc.mtx.Unlock()
	}()
	for otaKey, cachedData := range uc.cachedData {
		if otaKeyStr != "" && otaKey != otaKeyStr {
			continue
		}
		err = cachedData.store(uc.cacheDirectory)
		if err != nil {
			return err
		}
	}

	return nil
}

// load either loads the whole cache or a specific otaKey.
// Only the first value of `otaKeys` is processed.
func (uc *utxoCache) load(otaKeys ...string) error {
	var otaKeyStr string
	if len(otaKeys) > 0 {
		otaKeyStr = otaKeys[0]
		if otaKeyStr == "" {
			return fmt.Errorf("invalid OTAKey")
		}
	}

	if uc.cachedData != nil && otaKeyStr != "" {
		if _, ok := uc.cachedData[otaKeyStr]; ok {
			Logger.Printf("otaKey %v has already been loaded\n", otaKeyStr)
			return nil
		}
	}

	files, err := ioutil.ReadDir(uc.cacheDirectory)
	if err != nil {
		return err
	}

	cachedData := make(map[string]*accountCache)

	uc.mtx.Lock()
	defer func() {
		uc.mtx.Unlock()
	}()
	for _, f := range files {
		fileNameSplit := strings.Split(f.Name(), "/")
		otaKey := fileNameSplit[len(fileNameSplit)-1]
		if otaKeyStr != "" && otaKey != otaKeyStr {
			continue
		}
		ac := newAccountCache(otaKey)
		err = ac.load(uc.cacheDirectory)
		if err != nil {
			Logger.Printf("loadCacheUTXO fail for ota %v: %v", fileNameSplit, err)
			return err
		}
		cachedData[otaKey] = ac
	}
	if otaKeyStr != "" {
		if _, ok := cachedData[otaKeyStr]; !ok {
			Logger.Printf("otaKey %v not found in cache\n", otaKeyStr)
		}
	}
	uc.cachedData = cachedData

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

// addAccount adds an account to the cache, and saves it into a temp file if needed.
func (uc *utxoCache) addAccount(otaKey string, cachedAccount *accountCache, save bool) {
	uc.mtx.Lock()
	defer func() {
		uc.mtx.Unlock()
	}()
	uc.cachedData[otaKey] = cachedAccount
	if save {
		err := cachedAccount.store(uc.cacheDirectory)
		if err != nil {
			Logger.Printf("save file %v failed: %v\n", otaKey, err)
			delete(uc.cachedData, otaKey)
		}
	}
}

// syncOutCoinV2 syncs v2 output coins of an account w.r.t the given tokenIDStr.
func (client *IncClient) syncOutCoinV2(outCoinKey *rpc.OutCoinKey, tokenIDStr string) error {
	if tokenIDStr != common.PRVIDStr {
		tokenIDStr = common.ConfidentialAssetID.String()
	}

	// load the cache for the otaKey (if possible)
	err := client.cache.load(outCoinKey.OtaKey())
	if err != nil {
		return err
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
	start := time.Now()

	currentIndex := cachedToken.LatestIndex + 1
	if currentIndex == 1 {
		currentIndex = 0
	}
	Logger.Printf("Current LatestIndex for token %v: %v\n", tokenIDStr, cachedToken.LatestIndex)
	var rawAssetTags map[string]*common.Hash
	if currentIndex < coinLength {
		Logger.Printf("MaxGetCoinThreads: %v\n", MaxGetCoinThreads)
		statusChan := make(chan getCoinStatus, MaxGetCoinThreads)
		doneCount := 0
		mtx := new(sync.Mutex)
		numWorking := 0
		numThreads := math.Ceil(float64(coinLength-currentIndex) / float64(batchSize))
		for {
			select {
			case status := <-statusChan:
				if status.err != nil {
					err = fmt.Errorf("getCoinsByIndices FAILED at indices [%v-%v]: %v", status.fromIndex, status.toIndex, status.err)
					Logger.Println(err)
					return err
				} else {
					mtx.Lock()
					for idx, tmpCoin := range status.data {
						res.Data[idx] = tmpCoin
					}
					doneCount++
					numWorking--
					Logger.Printf("syncOutCoinV2 doneCount: %v/%v\n", doneCount, numThreads)
					mtx.Unlock()
				}
			default:
				if doneCount == int(numThreads) {
					break
				}
				if numWorking < MaxGetCoinThreads && currentIndex < coinLength {
					nextIndex := currentIndex + uint64(batchSize)
					if nextIndex > coinLength {
						nextIndex = coinLength
					}
					go client.getCoinsByIndices(keySet, shardID, tokenIDStr, currentIndex, nextIndex-1, statusChan)
					currentIndex = nextIndex
					mtx.Lock()
					numWorking++
					mtx.Unlock()
					Logger.Printf("syncOutCoinV2 timeElapsed: %v\n", time.Since(start).Seconds())
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}
			if doneCount == int(numThreads) {
				break
			}
		}

		if tokenIDStr != common.PRVIDStr && rawAssetTags == nil && len(res.Data) > 0 {
			// update cached data for each token
			rawAssetTags, err = client.GetAllAssetTags()
			if err != nil {
				return err
			}
		}

		Logger.Printf("newOutCoins: %v\n", len(res.Data))

		if tokenIDStr == common.PRVIDStr {
			cachedAccount.update(common.PRVIDStr, coinLength-1, *res)
		} else {
			err = cachedAccount.updateAllTokens(coinLength-1, *res, rawAssetTags)
			if err != nil {
				return err
			}
		}

		// add account to cache and save to file.
		client.cache.addAccount(outCoinKey.OtaKey(), cachedAccount, true)
	}

	Logger.Printf("FINISHED SYNCING OUTPUT COINS OF TOKEN %v AFTER %v SECOND\n", tokenIDStr, time.Since(start).Seconds())

	return nil
}

func (client *IncClient) getCoinsByIndices(
	keySet key.KeySet,
	shardID byte,
	tokenIDStr string,
	fromIndex, toIndex uint64,
	statusChan chan getCoinStatus,
) {
	res := make(map[uint64]jsonresult.ICoinInfo)
	start := time.Now()
	Logger.Printf("Get output coins of indices from %v to %v\n", fromIndex, toIndex)

	status := getCoinStatus{fromIndex: fromIndex, toIndex: toIndex}
	idxList := make([]uint64, 0)
	for i := fromIndex; i <= toIndex; i++ {
		idxList = append(idxList, i)
	}

	tmpOutCoins, err := client.GetOTACoinsByIndices(shardID, tokenIDStr, idxList)
	if err != nil {
		status.err = err
		statusChan <- status
		return
	}

	found := 0
	burningPubKey := wallet.GetBurningPublicKey()
	for idx, outCoin := range tmpOutCoins {
		if bytes.Equal(outCoin.Bytes(), burningPubKey) {
			continue
		}
		belongs, _ := outCoin.DoesCoinBelongToKeySet(&keySet)
		if belongs {
			res[idx] = outCoin
			found += 1
		}
	}

	Logger.Printf("Found %v output coins (%v) for heights from %v to %v with time %v\n", found, tokenIDStr, fromIndex, toIndex, time.Since(start).Seconds())
	status.data = res
	statusChan <- status
}

// GetAndCacheOutCoins retrieves the list of output coins and caches them for faster retrieval later.
// This function should only be called after the cache is initialized.
func (client *IncClient) GetAndCacheOutCoins(outCoinKey *rpc.OutCoinKey, tokenID string) ([]jsonresult.ICoinInfo, []*big.Int, error) {
	if client.cache == nil || !client.cache.isRunning {
		return nil, nil, fmt.Errorf("utxoCache is not running")
	}

	// sync v2 output coins from the remote node
	err := client.syncOutCoinV2(outCoinKey, tokenID)
	if err != nil {
		return nil, nil, err
	}

	outCoins := make([]jsonresult.ICoinInfo, 0)
	indices := make([]*big.Int, 0)

	// query v2 output coins
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

	// query v1 output coins
	otaKey := outCoinKey.OtaKey()
	outCoinKey.SetOTAKey("") // set this to empty so that the full-node only query v1 output coins.
	v1OutCoins, _, err := client.GetOutputCoinsV1(outCoinKey, tokenID, 0)
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
	outCoinKey.SetOTAKey(otaKey)
	Logger.Printf("Found %v v1 output coins\n", v1Count)

	return outCoins, indices, nil
}
