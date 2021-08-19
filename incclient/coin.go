package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"math/big"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// GetOutputCoins calls the remote server to get all the output tokens for an output coin key.
// `isFromCache` indicates whether the client should retrieve output tokens from the local cache.
// In case this value is not set, the client uses the regular `GetOutputCoins` method.
// If multiple values are passed to `isFromCache`, only the first one is used.
// For better user experience, if the cache is not running and isFromCache holds true, the client still automatically
// switches to the non-cache method.
//
// The returned result consists of
//	- A list of output coins
//	- A list of corresponding indices. For an output coin v1, its index is -1.
func (client *IncClient) GetOutputCoins(outCoinKey *rpc.OutCoinKey, tokenID string, height uint64, isFromCache ...bool) ([]jsonresult.ICoinInfo, []*big.Int, error) {
	fromCache := false
	if len(isFromCache) != 0 {
		fromCache = isFromCache[0]
	}

	if fromCache && client.cache != nil && client.cache.isRunning {
		return client.GetAndCacheOutCoins(outCoinKey, tokenID)
	}

	if client.version == 1 {
		return client.GetOutputCoinsV1(outCoinKey, tokenID, height)
	} else {
		return client.GetOutputCoinsV2(outCoinKey, tokenID, height)
	}
}

// GetOutputCoinsV1 calls the remote server to get all the output tokens for an output coin key using the old RPC.
func (client *IncClient) GetOutputCoinsV1(outCoinKey *rpc.OutCoinKey, tokenID string, height uint64) ([]jsonresult.ICoinInfo, []*big.Int, error) {
	b, err := client.rpcServer.GetListOutputCoinsByRPCV1(outCoinKey, tokenID, height)
	if err != nil {
		return nil, nil, err
	}

	return ParseCoinFromJsonResponse(b)
}

// GetOutputCoinsV2 calls the remote server to get all the output tokens for an output coin key using the new RPC.
//
// For this function, it is required that the caller has submitted the OTA key to the remote full-node.
func (client *IncClient) GetOutputCoinsV2(outCoinKey *rpc.OutCoinKey, tokenID string, upToHeight uint64) ([]jsonresult.ICoinInfo, []*big.Int, error) {
	b, err := client.rpcServer.GetListOutputCoinsByRPCV2(outCoinKey, tokenID, upToHeight)
	if err != nil {
		return nil, nil, err
	}

	return ParseCoinFromJsonResponse(b)
}

// GetListDecryptedOutCoin retrieves and decrypts all the output tokens for a private key.
// It returns
//	- a map from the serial number to the output coin;
//	- error (if any).
func (client *IncClient) GetListDecryptedOutCoin(privateKey string, tokenID string, height uint64) (map[string]coin.PlainCoin, error) {
	outCoinKey, err := NewOutCoinKeyFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	outCoinKey.SetReadonlyKey("") // call this if you do not want the remote full-node to decrypt your coin

	listOutputCoins, _, err := client.GetOutputCoins(outCoinKey, tokenID, height)
	if err != nil {
		return nil, err
	}

	if len(listOutputCoins) == 0 {
		return nil, nil
	}

	listDecryptedOutCoins, listKeyImages, err := GetListDecryptedCoins(privateKey, listOutputCoins)
	if err != nil {
		return nil, err
	}

	mapOutCoin := make(map[string]coin.PlainCoin)
	if len(listDecryptedOutCoins) != len(listKeyImages) {
		return nil, fmt.Errorf("have %v output coins but %v serial numbers", len(listDecryptedOutCoins), len(listKeyImages))
	}

	for i, outCoin := range listDecryptedOutCoins {
		mapOutCoin[listKeyImages[i]] = outCoin
	}

	return mapOutCoin, nil
}

// CheckCoinsSpent checks if the provided serial numbers have been spent or not.
//
// Returned result in boolean list.
func (client *IncClient) CheckCoinsSpent(shardID byte, tokenID string, snList []string) ([]bool, error) {
	b, err := client.rpcServer.HasSerialNumberByRPC(shardID, tokenID, snList)
	if err != nil {
		return []bool{}, err
	}

	var tmp []bool
	err = rpchandler.ParseResponse(b, &tmp)
	if err != nil {
		return []bool{}, err
	}

	if len(tmp) != len(snList) {
		return []bool{}, fmt.Errorf("length of result and length of snList mismathc: len(Result) = %v, len(snList) = %v; perhaps the shardID was wrong", len(tmp), len(snList))
	}

	return tmp, nil
}

// GetUnspentOutputCoins retrieves all unspent coins of a private key, without sending the private key to the remote full-node.
func (client *IncClient) GetUnspentOutputCoins(privateKey, tokenID string, height uint64) ([]coin.PlainCoin, []*big.Int, error) {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, nil, err
	}

	outCoinKey, err := NewOutCoinKeyFromPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	outCoinKey.SetReadonlyKey("") // call this if you do not want the remote full-node to decrypt your coin

	listOutputCoins, listIndices, err := client.GetOutputCoins(outCoinKey, tokenID, height)
	if err != nil {
		return nil, nil, err
	}

	if len(listOutputCoins) == 0 {
		return nil, nil, nil
	}

	listDecryptedOutCoins, listKeyImages, err := GetListDecryptedCoins(privateKey, listOutputCoins)
	if err != nil {
		return nil, nil, err
	}

	shardID := common.GetShardIDFromLastByte(keyWallet.KeySet.PaymentAddress.Pk[len(keyWallet.KeySet.PaymentAddress.Pk)-1])
	checkSpentList, err := client.CheckCoinsSpent(shardID, tokenID, listKeyImages)
	if err != nil {
		return nil, nil, err
	}

	listUnspentOutputCoins := make([]coin.PlainCoin, 0)
	listUnspentIndices := make([]*big.Int, 0)
	for i, decryptedCoin := range listDecryptedOutCoins {
		if !checkSpentList[i] && decryptedCoin.GetValue() != 0 {
			listUnspentOutputCoins = append(listUnspentOutputCoins, decryptedCoin)
			listUnspentIndices = append(listUnspentIndices, listIndices[i])
		}
	}

	return listUnspentOutputCoins, listUnspentIndices, nil
}

// GetUnspentOutputCoinsFromCache retrieves all unspent coins of a private key, without sending the private key to the remote full-node.
func (client *IncClient) GetUnspentOutputCoinsFromCache(privateKey, tokenID string, height uint64) ([]coin.PlainCoin, []*big.Int, error) {
	return client.GetUnspentOutputCoins(privateKey, tokenID, height)
}

// GetSpentOutputCoins retrieves all spent coins of a private key, without sending the private key to the remote full node.
func (client *IncClient) GetSpentOutputCoins(privateKey, tokenID string, height uint64) ([]coin.PlainCoin, []*big.Int, error) {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, nil, err
	}

	outCoinKey, err := NewOutCoinKeyFromPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	outCoinKey.SetReadonlyKey("") // call this if you do not want the remote full node to decrypt your coin

	listOutputCoins, listIndices, err := client.GetOutputCoins(outCoinKey, tokenID, height)
	if err != nil {
		return nil, nil, err
	}

	Logger.Printf("Len(OutputCoins) = %v\n", len(listOutputCoins))

	if len(listOutputCoins) == 0 {
		return nil, nil, nil
	}

	listDecryptedOutCoins, listKeyImages, err := GetListDecryptedCoins(privateKey, listOutputCoins)
	if err != nil {
		return nil, nil, err
	}

	shardID := common.GetShardIDFromLastByte(keyWallet.KeySet.PaymentAddress.Pk[len(keyWallet.KeySet.PaymentAddress.Pk)-1])
	checkSpentList, err := client.CheckCoinsSpent(shardID, tokenID, listKeyImages)
	if err != nil {
		return nil, nil, err
	}

	listSpentOutputCoins := make([]coin.PlainCoin, 0)
	listSpentIndices := make([]*big.Int, 0)
	for i, decryptedCoin := range listDecryptedOutCoins {
		if checkSpentList[i] && decryptedCoin.GetValue() != 0 {
			listSpentOutputCoins = append(listSpentOutputCoins, decryptedCoin)
			listSpentIndices = append(listSpentIndices, listIndices[i])
		}
	}

	Logger.Printf("Len(spentCoins) = %v\n", len(listSpentOutputCoins))

	return listSpentOutputCoins, listSpentIndices, nil
}

// GetOTACoinsByIndices retrieves output coins given the shardID, tokenID, and a list of indices.
// Callers must make sure the input indices are within the range of output coins ([0 : OTALength - 1]) of the given shard.
// Otherwise, it will throw an error.
func (client *IncClient) GetOTACoinsByIndices(shardID byte, tokenID string, idxList []uint64) (map[uint64]jsonresult.ICoinInfo, error) {
	responseInBytes, err := client.rpcServer.GetOTACoinsByIndices(shardID, tokenID, idxList)
	if err != nil {
		return nil, err
	}

	var tmpRes map[uint64]jsonresult.OutCoin
	err = rpchandler.ParseResponse(responseInBytes, &tmpRes)
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]jsonresult.ICoinInfo)
	for idx, outCoin := range tmpRes {
		coinInfo, _, err := jsonresult.NewCoinFromJsonOutCoin(outCoin)
		if err != nil {
			return nil, err
		}
		res[idx] = coinInfo
	}

	return res, nil
}

// GetOTACoinLength returns the current sizes (number of output coins) of PRV and tokens for each shard.
//
// Sample output:
// map[
//	0000000000000000000000000000000000000000000000000000000000000004:map[0:43635 1:14775 2:11407 3:20724 4:12613 5:10165 6:18216 7:15629]
//	0000000000000000000000000000000000000000000000000000000000000005:map[0:37831 1:7379 2:2349 3:11310 4:2218 5:2247 6:8416 7:8601]
// ]
func (client *IncClient) GetOTACoinLength() (map[string]map[byte]uint64, error) {
	responseInBytes, err := client.rpcServer.GetOTACoinLength()
	if err != nil {
		return nil, err
	}

	var res map[string]map[byte]uint64
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetOTACoinLengthByShard returns the number of output coins of a shard w.r.t to the given tokenID.
//
// Note that all tokens' output coins are grouped together (except PRV), there for it will return the same result for
// two different tokenIDs (jus use the common.ConfidentialAssetID for token retrieval).
func (client *IncClient) GetOTACoinLengthByShard(shardID byte, tokenID string) (uint64, error) {
	if int(shardID) >= common.MaxShardNumber {
		return 0, fmt.Errorf("invalid shardID %v", shardID)
	}

	responseInBytes, err := client.rpcServer.GetOTACoinLength()
	if err != nil {
		return 0, err
	}

	var res map[string]map[byte]uint64
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return 0, err
	}

	switch tokenID {
	case common.PRVIDStr:
		if res[tokenID] == nil {
			return 0, fmt.Errorf("RPC result does not include PRV")
		}
		return res[tokenID][shardID], nil
	default:
		if res[common.ConfidentialAssetID.String()] == nil {
			return 0, fmt.Errorf("RPC result does not include tokens")
		}
		return res[common.ConfidentialAssetID.String()][shardID], nil
	}
}

// GetAllAssetTags computes a mapping from raw assetTags to tokenIds (e.g, HashToPoint(PRV) => PRV).
func (client *IncClient) GetAllAssetTags() (map[string]*common.Hash, error) {
	assetTags := make(map[string]*common.Hash)
	assetTags[crypto.HashToPoint(common.PRVCoinID[:]).String()] = &common.PRVCoinID
	listTokens, err := client.GetListToken()
	if err != nil {
		return nil, err
	}
	for tokenIdStr := range listTokens {
		tokenHash, err := new(common.Hash).NewHashFromStr(tokenIdStr)
		if err != nil {
			return nil, err
		}
		assetTags[crypto.HashToPoint(tokenHash[:]).String()] = tokenHash
	}

	return assetTags, nil
}

// NewOutCoinKeyFromPrivateKey creates a new rpc.OutCoinKey given the private key.
func NewOutCoinKeyFromPrivateKey(privateKey string) (*rpc.OutCoinKey, error) {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, err
	}

	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	if err != nil {
		return nil, err
	}
	paymentAddStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	otaSecretKey := keyWallet.Base58CheckSerialize(wallet.OTAKeyType)
	viewingKeyStr := keyWallet.Base58CheckSerialize(wallet.ReadonlyKeyType)

	return rpc.NewOutCoinKey(paymentAddStr, otaSecretKey, viewingKeyStr), err
}

// ParseCoinFromJsonResponse parses raw coin data returned from an RPC request into a list of ICoinInfo.
func ParseCoinFromJsonResponse(b []byte) ([]jsonresult.ICoinInfo, []*big.Int, error) {
	var tmp jsonresult.ListOutputCoins
	err := rpchandler.ParseResponse(b, &tmp)
	if err != nil {
		return nil, nil, err
	}

	resultOutCoins := make([]jsonresult.ICoinInfo, 0)
	listOutputCoins := tmp.Outputs
	listIndices := make([]*big.Int, 0)
	for _, value := range listOutputCoins {
		for _, outCoin := range value {
			out, idx, err := jsonresult.NewCoinFromJsonOutCoin(outCoin)
			if err != nil {
				return nil, nil, err
			}

			resultOutCoins = append(resultOutCoins, out)
			listIndices = append(listIndices, idx)
		}
	}

	return resultOutCoins, listIndices, nil
}

// GetListDecryptedCoins decrypts a list of ICoinInfo's using the given private key.
func GetListDecryptedCoins(privateKey string, listOutputCoins []jsonresult.ICoinInfo) ([]coin.PlainCoin, []string, error) {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, nil, err
	}
	keySet := keyWallet.KeySet

	listDecryptedOutCoins := make([]coin.PlainCoin, 0)
	listKeyImages := make([]string, 0)
	for _, outCoin := range listOutputCoins {
		if outCoin.GetVersion() == 1 {
			if outCoin.IsEncrypted() {
				tmpCoin, ok := outCoin.(*coin.CoinV1)
				if !ok {
					return nil, nil, fmt.Errorf("invalid CoinV1")
				}

				decryptedCoin, err := tmpCoin.Decrypt(&keySet)
				if err != nil {
					return nil, nil, err
				}
				keyImage, err := decryptedCoin.ParseKeyImageWithPrivateKey(keyWallet.KeySet.PrivateKey)
				if err != nil {
					return nil, nil, err
				}
				decryptedCoin.SetKeyImage(keyImage)

				keyImageString := base58.Base58Check{}.Encode(keyImage.ToBytesS(), common.ZeroByte)

				listKeyImages = append(listKeyImages, keyImageString)
				listDecryptedOutCoins = append(listDecryptedOutCoins, decryptedCoin)
			} else {
				tmpPlainCoinV1, ok := outCoin.(*coin.PlainCoinV1)
				if !ok {
					return nil, nil, fmt.Errorf("invalid PlaincoinV1")
				}

				keyImage, err := tmpPlainCoinV1.ParseKeyImageWithPrivateKey(keyWallet.KeySet.PrivateKey)
				if err != nil {
					return nil, nil, err
				}
				tmpPlainCoinV1.SetKeyImage(keyImage)

				keyImageString := base58.Base58Check{}.Encode(keyImage.ToBytesS(), common.ZeroByte)

				listKeyImages = append(listKeyImages, keyImageString)
				listDecryptedOutCoins = append(listDecryptedOutCoins, tmpPlainCoinV1)
			}
		} else if outCoin.GetVersion() == 2 {
			tmpCoinV2, ok := outCoin.(*coin.CoinV2)
			if !ok {
				return nil, nil, fmt.Errorf("invalid CoinV2")
			}
			decryptedCoin, err := tmpCoinV2.Decrypt(&keyWallet.KeySet)
			if err != nil {
				return nil, nil, err
			}
			keyImage := decryptedCoin.GetKeyImage()
			keyImageString := base58.Base58Check{}.Encode(keyImage.ToBytesS(), common.ZeroByte)

			listKeyImages = append(listKeyImages, keyImageString)
			listDecryptedOutCoins = append(listDecryptedOutCoins, decryptedCoin)
		}
	}

	return listDecryptedOutCoins, listKeyImages, nil
}

// GenerateOTAFromPaymentAddress generates a random one-time address, and TxRandom from a payment address.
func GenerateOTAFromPaymentAddress(paymentAddressStr string) (string, string, error) {
	keyWallet, err := wallet.Base58CheckDeserialize(paymentAddressStr)
	if err != nil {
		return "", "", err
	}

	paymentInfo := key.InitPaymentInfo(keyWallet.KeySet.PaymentAddress, 0, []byte{})
	otaCoin, err := coin.NewCoinFromPaymentInfo(paymentInfo)
	if err != nil {
		return "", "", err
	}

	pubKey := otaCoin.GetPublicKey()
	txRandom := otaCoin.GetTxRandom()

	return base58.Base58Check{}.Encode(pubKey.ToBytesS(), common.ZeroByte), base58.Base58Check{}.Encode(txRandom.Bytes(), common.ZeroByte), nil
}
