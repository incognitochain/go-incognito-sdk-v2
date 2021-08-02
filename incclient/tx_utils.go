package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"math/big"
	"sort"
	"time"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

const pageSize = 100

type coinParams struct {
	coinList []coin.PlainCoin
	idxList  []uint64
}

func (cp coinParams) Bytes() []byte {
	if cp.coinList == nil || len(cp.coinList) == 0 {
		return nil
	}
	resBytes := make([]byte, 0)

	// first byte is the number of coins
	resBytes = append(resBytes, byte(len(cp.coinList)))
	for _, c := range cp.coinList {
		cBytes := c.Bytes()
		resBytes = append(resBytes, byte(len(cBytes)))
		resBytes = append(resBytes, cBytes...)
	}

	// next byte is the length of indices
	resBytes = append(resBytes, byte(len(cp.idxList)))
	for _, idx := range cp.idxList {
		idxBytes := common.IntToBytes(int(idx))
		resBytes = append(resBytes, byte(len(idxBytes)))
		resBytes = append(resBytes, idxBytes...)
	}

	return resBytes
}

func (cp *coinParams) SetBytes(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("length data is zero")
	}

	var err error
	offSet := 0

	// get num input coins
	if offSet >= len(data) {
		return fmt.Errorf("out of range  numInCoins")
	}
	numInCoins := int(data[offSet])
	offSet++
	cp.coinList = make([]coin.PlainCoin, numInCoins)
	for i := 0; i < numInCoins; i++ {
		if offSet >= len(data) {
			return fmt.Errorf("out of range lenCoin")
		}
		lenCoin := int(data[offSet])
		offSet++

		if offSet+lenCoin > len(data) {
			return fmt.Errorf("out of range input coins")
		}
		coinBytes := data[offSet : offSet+lenCoin]
		cp.coinList[i], err = coin.NewPlainCoinFromByte(coinBytes)
		if err != nil {
			return fmt.Errorf("set byte to inputCoin got error")
		}
		offSet += lenCoin
	}

	if offSet >= len(data) {
		return fmt.Errorf("out of range numIndices")
	}
	numIndices := int(data[offSet])
	offSet++
	if numIndices == 0 {
		return nil
	}

	cp.idxList = make([]uint64, numIndices)
	for i := 0; i < numIndices; i++ {
		if offSet >= len(data) {
			return fmt.Errorf("out of range lenIdx")
		}
		lenIdx := int(data[offSet])
		offSet++

		if offSet+lenIdx > len(data) {
			return fmt.Errorf("out of range index")
		}
		idxBytes := data[offSet : offSet+lenIdx]
		cp.idxList[i] = uint64(common.BytesToInt(idxBytes))
		offSet += lenIdx
	}

	if len(cp.idxList) > 0 && len(cp.idxList) != len(cp.coinList) {
		return fmt.Errorf("lengths of idxList and coinList mismatch: %v != %v", len(cp.idxList), len(cp.coinList))
	}

	return nil
}

// createPaymentInfos creates a list of key.PaymentInfo based on the provided address list and corresponding amount list.
func createPaymentInfos(addrList []string, amountList []uint64) ([]*key.PaymentInfo, error) {
	if len(addrList) != len(amountList) {
		return nil, fmt.Errorf("length of payment address (%v) and length amount (%v) mismatch", len(addrList), len(amountList))
	}

	paymentInfos := make([]*key.PaymentInfo, 0)
	for i, addr := range addrList {
		receiverWallet, err := wallet.Base58CheckDeserialize(addr)
		if err != nil {
			return nil, fmt.Errorf("cannot deserialize key %v: %v", addr, err)
		}
		paymentInfo := key.PaymentInfo{PaymentAddress: receiverWallet.KeySet.PaymentAddress, Amount: amountList[i], Message: []byte{}}
		paymentInfos = append(paymentInfos, &paymentInfo)
	}

	return paymentInfos, nil
}

// chooseBestCoinsByAmount chooses best UTXOs to spend depending on the provided amount.
//
// Assume that the input coins have be sorted in the descending order.
func chooseBestCoinsByAmount(coinList []coin.PlainCoin, requiredAmount uint64) ([]coin.PlainCoin, []uint64, error) {
	totalInputAmount := uint64(0)
	for _, inputCoin := range coinList {
		totalInputAmount += inputCoin.GetValue()
	}

	if totalInputAmount < requiredAmount {
		return nil, nil, fmt.Errorf("total unspent amount (%v) is less than the required amount (%v)", totalInputAmount, requiredAmount)
	}

	if totalInputAmount == requiredAmount {
		chosenIndexList := make([]uint64, 0)
		for i := 0; i < len(coinList); i++ {
			chosenIndexList = append(chosenIndexList, uint64(i))
		}
		return coinList, chosenIndexList, nil
	}

	coinsToSpend := make([]coin.PlainCoin, 0)
	chosenIndexList := make([]uint64, 0)
	remainAmount := requiredAmount
	totalChosenAmount := uint64(0)
	//TODO: find a better solution for this.
	for i := 0; i < len(coinList)-1; i++ {
		if coinList[i].GetValue() > remainAmount {
			if coinList[i+1].GetValue() >= remainAmount {
				continue
			} else {
				coinsToSpend = append(coinsToSpend, coinList[i])
				chosenIndexList = append(chosenIndexList, uint64(i))
				totalChosenAmount += coinList[i].GetValue()
				break
			}
		} else {
			coinsToSpend = append(coinsToSpend, coinList[i])
			chosenIndexList = append(chosenIndexList, uint64(i))
			remainAmount -= coinList[i].GetValue()
			totalChosenAmount += coinList[i].GetValue()
		}
	}

	if totalChosenAmount < requiredAmount {
		totalChosenAmount += coinList[len(coinList)-1].GetValue()
		coinsToSpend = append(coinsToSpend, coinList[len(coinList)-1])
		chosenIndexList = append(chosenIndexList, uint64(len(coinList)-1))
		if totalChosenAmount < requiredAmount {
			return nil, nil, fmt.Errorf("not enough coin to spend")
		}
	}

	return coinsToSpend, chosenIndexList, nil
}

// divideCoins divides the list of coins w.r.t their version and sort them by values if needed.
func divideCoins(coinList []coin.PlainCoin, idxList []*big.Int, needSorted bool) ([]coin.PlainCoin, []coin.PlainCoin, []uint64, error) {
	if idxList != nil {
		if len(coinList) != len(idxList) {
			return nil, nil, nil, fmt.Errorf("cannot divide coins: length of coin (%v) != length of index (%v)", len(coinList), len(idxList))
		}
	}

	coinV1List := make([]coin.PlainCoin, 0)
	coinV2List := make([]coin.PlainCoin, 0)
	idxV2List := make([]uint64, 0)
	for i, inputCoin := range coinList {
		if inputCoin.GetVersion() == 2 {
			tmpCoin, ok := inputCoin.(*coin.CoinV2)
			if !ok {
				return nil, nil, nil, fmt.Errorf("cannot parse coinV2")
			}

			coinV2List = append(coinV2List, tmpCoin)
			if idxList != nil {
				if idxList[i] == nil {
					return nil, nil, nil, fmt.Errorf("idx of coinV2 %v is nil: (idxList: %v)", i, idxList)
				}
				idxV2List = append(idxV2List, idxList[i].Uint64())
			}
		} else {
			tmpCoin, ok := inputCoin.(*coin.PlainCoinV1)
			if !ok {
				return nil, nil, nil, fmt.Errorf("cannot parse coinV2")
			}

			coinV1List = append(coinV1List, tmpCoin)
		}
	}

	if needSorted {
		sort.Slice(coinV1List, func(i, j int) bool {
			return coinV1List[i].GetValue() > coinV1List[j].GetValue()
		})

		sort.Slice(coinV2List, func(i, j int) bool {
			return coinV2List[i].GetValue() > coinV2List[j].GetValue()
		})

		var err error
		idxV2List, err = getListIdx(coinV2List, coinList, idxList)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return coinV1List, coinV2List, idxV2List, nil
}

func getListIdx(inCoins []coin.PlainCoin, allCoins []coin.PlainCoin, allIdx []*big.Int) ([]uint64, error) {
	if len(allIdx) == 0 {
		return []uint64{}, nil
	}
	res := make([]uint64, 0)
	for _, inCoin := range inCoins {
		for i, c := range allCoins {
			if c.GetVersion() != 2 {
				continue
			}
			if c.GetPublicKey().String() == inCoin.GetPublicKey().String() {
				res = append(res, allIdx[i].Uint64())
				break
			}
		}
	}

	if len(res) != len(inCoins) {
		return nil, fmt.Errorf("some coin cannot be retrieved")
	}

	return res, nil
}

// getVersionFromInputCoins checks if all of the given input coins have the same version, and return the version.
func getVersionFromInputCoins(inputCoins []coin.PlainCoin) (uint8, error) {
	if len(inputCoins) == 0 {
		return 0, fmt.Errorf("no coin to check")
	}

	version := inputCoins[0].GetVersion()
	for i := 1; i < len(inputCoins); i++ {
		if inputCoins[i].GetVersion() != version {
			return 0, fmt.Errorf("expect input coin %v to have version %v, but got %v", i, version, inputCoins[i].GetVersion())
		}
	}

	return version, nil
}

func (client *IncClient) getRandomCommitmentV1(inputCoins []coin.PlainCoin, tokenID string) (map[string]interface{}, error) {
	if len(inputCoins) == 0 {
		return nil, fmt.Errorf("no input coin to retrieve random commitments, tokenID: %v", tokenID)
	}
	outCoinList := make([]jsonresult.OutCoin, 0)
	for _, inputCoin := range inputCoins {
		outCoin := jsonresult.NewOutCoin(inputCoin)
		outCoin.Conceal()
		outCoinList = append(outCoinList, outCoin)
	}

	lastByte := inputCoins[0].GetPublicKey().ToBytesS()[len(inputCoins[0].GetPublicKey().ToBytesS())-1]
	shardID := common.GetShardIDFromLastByte(lastByte)

	responseInBytes, err := client.rpcServer.RandomCommitments(shardID, outCoinList, tokenID)
	if err != nil {
		return nil, err
	}

	var randomCommitment jsonresult.RandomCommitmentResult
	err = rpchandler.ParseResponse(responseInBytes, &randomCommitment)
	if err != nil {
		return nil, err
	}

	commitmentList := make([]*crypto.Point, 0)
	for _, commitmentStr := range randomCommitment.Commitments {
		cmtBytes, _, err := base58.Base58Check{}.Decode(commitmentStr)
		if err != nil {
			return nil, fmt.Errorf("cannot decode commitment %v: %v", commitmentStr, err)
		}

		commitment, err := new(crypto.Point).FromBytesS(cmtBytes)
		if err != nil {
			return nil, fmt.Errorf("cannot parse commitment %v: %v", cmtBytes, err)
		}

		commitmentList = append(commitmentList, commitment)
	}

	result := make(map[string]interface{})
	result[utils.CommitmentIndices] = randomCommitment.CommitmentIndices
	result[utils.MyIndices] = randomCommitment.MyCommitmentIndices
	result[utils.Commitments] = commitmentList

	return result, nil
}

func (client *IncClient) getRandomCommitmentV2(shardID byte, tokenID string, lenDecoy int) (map[string]interface{}, error) {
	if lenDecoy == 0 {
		return nil, fmt.Errorf("no input coin to retrieve random commitments")
	}

	responseInBytes, err := client.rpcServer.RandomCommitmentsAndPublicKeys(shardID, tokenID, lenDecoy)
	if err != nil {
		return nil, err
	}

	var randomCmtAndPk jsonresult.RandomCommitmentAndPublicKeyResult
	err = rpchandler.ParseResponse(responseInBytes, &randomCmtAndPk)
	if err != nil {
		return nil, err
	}

	commitmentList := make([]*crypto.Point, 0)
	for _, commitmentStr := range randomCmtAndPk.Commitments {
		cmtBytes, _, err := base58.Base58Check{}.Decode(commitmentStr)
		if err != nil {
			return nil, fmt.Errorf("cannot decode commitment %v: %v", commitmentStr, err)
		}

		commitment, err := new(crypto.Point).FromBytesS(cmtBytes)
		if err != nil {
			return nil, fmt.Errorf("cannot parse commitment %v: %v", cmtBytes, err)
		}

		commitmentList = append(commitmentList, commitment)
	}

	pkList := make([]*crypto.Point, 0)
	for _, pubKeyStr := range randomCmtAndPk.PublicKeys {
		pkBytes, _, err := base58.Base58Check{}.Decode(pubKeyStr)
		if err != nil {
			return nil, fmt.Errorf("cannot decode public key %v: %v", pubKeyStr, err)
		}

		pk, err := new(crypto.Point).FromBytesS(pkBytes)
		if err != nil {
			return nil, fmt.Errorf("cannot parse public key %v: %v", pkBytes, err)
		}

		pkList = append(pkList, pk)
	}

	assetTagList := make([]*crypto.Point, 0)
	for _, assetStr := range randomCmtAndPk.AssetTags {
		assetBytes, _, err := base58.Base58Check{}.Decode(assetStr)
		if err != nil {
			return nil, fmt.Errorf("cannot decode assetTag %v: %v", assetStr, err)
		}

		assetTag, err := new(crypto.Point).FromBytesS(assetBytes)
		if err != nil {
			return nil, fmt.Errorf("cannot parse assetTag %v: %v", assetBytes, err)
		}

		assetTagList = append(assetTagList, assetTag)
	}

	result := make(map[string]interface{})
	result[utils.CommitmentIndices] = randomCmtAndPk.CommitmentIndices
	result[utils.Commitments] = commitmentList
	result[utils.PublicKeys] = pkList
	result[utils.AssetTags] = assetTagList

	return result, nil
}

// initParams queries and chooses coins to spend + init random params.
func (client *IncClient) initParams(privateKey string, tokenIDStr string, totalAmount uint64, hasPrivacy bool, version int) ([]coin.PlainCoin, map[string]interface{}, error) {
	_, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, nil, err
	}
	//Create sender private key from string
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	lastByteSender := senderWallet.KeySet.PaymentAddress.Pk[len(senderWallet.KeySet.PaymentAddress.Pk)-1]
	shardID := common.GetShardIDFromLastByte(lastByteSender)

	//fmt.Printf("Getting UTXOs for tokenID %v...\n", tokenIDStr)
	//Get list of UTXOs
	utxoList, idxList, err := client.GetUnspentOutputCoinsFromCache(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, nil, err
	}

	//fmt.Printf("Finish getting UTXOs for %v of %v. Length of UTXOs: %v\n", totalAmount, tokenIDStr, len(utxoList))
	coinV1List, coinV2List, idxV2List, err := divideCoins(utxoList, idxList, true)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot divide coin: %v", err)
	}

	var coinsToSpend []coin.PlainCoin
	var kvArgs = make(map[string]interface{})
	if version == 1 {
		//Choose best coins for creating transactions
		coinsToSpend, _, err = chooseBestCoinsByAmount(coinV1List, totalAmount)
		if err != nil {
			return nil, nil, err
		}

		if hasPrivacy {
			//fmt.Printf("Getting random commitments for %v.\n", tokenIDStr)
			//Retrieve commitments and indices
			kvArgs, err = client.getRandomCommitmentV1(coinsToSpend, tokenIDStr)
			if err != nil {
				return nil, nil, err
			}
			//fmt.Printf("Finish getting random commitments.\n")
		}

		return coinsToSpend, kvArgs, nil
	} else {
		var chosenIdxList []uint64
		coinsToSpend, chosenIdxList, err = chooseBestCoinsByAmount(coinV2List, totalAmount)
		if err != nil {
			return nil, nil, err
		}

		//fmt.Printf("Getting random commitments for %v.\n", tokenIDStr)
		//Retrieve commitments and indices
		kvArgs, err = client.getRandomCommitmentV2(shardID, tokenIDStr, len(coinsToSpend)*(privacy.RingSize-1))
		if err != nil {
			return nil, nil, err
		}
		//fmt.Printf("Finish getting random commitments.\n")
		idxToSpendPRV := make([]uint64, 0)
		for _, idx := range chosenIdxList {
			idxToSpendPRV = append(idxToSpendPRV, idxV2List[idx])
		}
		kvArgs[utils.MyIndices] = idxToSpendPRV

		return coinsToSpend, kvArgs, nil
	}
}

// initParamsV1 queries and chooses coins to spend + init random params v1.
func (client *IncClient) initParamsV1(txParam *TxParam, tokenIDStr string, totalAmount uint64, hasPrivacy bool) ([]coin.PlainCoin, map[string]interface{}, error) {
	_, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, nil, err
	}

	//Create sender private key from string
	privateKey := txParam.senderPrivateKey

	var coinsToSpend []coin.PlainCoin
	if txParam.kArgs != nil { // in case we use provided input coins to init the transaction.
		var ok bool
		var cpInterface interface{}
		if tokenIDStr == common.PRVIDStr {
			cpInterface, ok = txParam.kArgs[prvInCoinKey]
		} else {
			cpInterface, ok = txParam.kArgs[tokenInCoinKey]
		}

		if ok {
			cp, ok := cpInterface.(coinParams)
			if ok {
				v, _ := getVersionFromInputCoins(cp.coinList)
				if v == 1 {
					coinsToSpend = cp.coinList
				}
			}
		}
	}

	if coinsToSpend == nil {
		//Get list of UTXOs
		utxoList, idxList, err := client.GetUnspentOutputCoinsFromCache(privateKey, tokenIDStr, 0)
		if err != nil {
			return nil, nil, err
		}

		//fmt.Printf("Finish getting UTXOs for %v of %v. Length of UTXOs: %v\n", totalAmount, tokenIDStr, len(utxoList))
		coinV1List, _, _, err := divideCoins(utxoList, idxList, true)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot divide coin: %v", err)
		}

		//Choose best coins for creating transactions
		coinsToSpend, _, err = chooseBestCoinsByAmount(coinV1List, totalAmount)
		if err != nil {
			return nil, nil, err
		}
	}

	var kvArgs = make(map[string]interface{})
	if hasPrivacy {
		//Retrieve commitments and indices
		kvArgs, err = client.getRandomCommitmentV1(coinsToSpend, tokenIDStr)
		if err != nil {
			return nil, nil, err
		}
	}

	return coinsToSpend, kvArgs, nil

}

// initParamsV2 queries and chooses coins to spend + init random params v2.
func (client *IncClient) initParamsV2(txParam *TxParam, tokenIDStr string, totalAmount uint64) ([]coin.PlainCoin, map[string]interface{}, error) {
	_, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, nil, err
	}

	//Create sender private key from string
	privateKey := txParam.senderPrivateKey
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot init private key %v: %v", privateKey, err)
	}

	lastByteSender := senderWallet.KeySet.PaymentAddress.Pk[len(senderWallet.KeySet.PaymentAddress.Pk)-1]
	shardID := common.GetShardIDFromLastByte(lastByteSender)

	var coinsToSpend []coin.PlainCoin
	var myIndices []uint64
	if txParam.kArgs != nil { // in case we use provided input coins to init the transaction.
		var ok bool
		var cpInterface interface{}
		if tokenIDStr == common.PRVIDStr {
			cpInterface, ok = txParam.kArgs[prvInCoinKey]
		} else {
			cpInterface, ok = txParam.kArgs[tokenInCoinKey]
		}

		if ok {
			cp, ok := cpInterface.(coinParams)
			if ok {
				v, _ := getVersionFromInputCoins(cp.coinList)
				if v == 2 {
					coinsToSpend = cp.coinList
					myIndices = cp.idxList
				}
			}
		}
	} // in case we use provided input coins to init the transaction.

	if coinsToSpend == nil {
		//Get list of UTXOs
		utxoList, idxList, err := client.GetUnspentOutputCoinsFromCache(privateKey, tokenIDStr, 0)
		if err != nil {
			return nil, nil, err
		}

		_, coinV2List, idxV2List, err := divideCoins(utxoList, idxList, true)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot divide coin: %v", err)
		}

		var chosenIdxList []uint64
		coinsToSpend, chosenIdxList, err = chooseBestCoinsByAmount(coinV2List, totalAmount)
		if err != nil {
			return nil, nil, err
		}

		myIndices = make([]uint64, 0)
		for _, idx := range chosenIdxList {
			myIndices = append(myIndices, idxV2List[idx])
		}
	}

	//Retrieve commitments and indices
	var kvArgs = make(map[string]interface{})
	kvArgs, err = client.getRandomCommitmentV2(shardID, tokenIDStr, len(coinsToSpend)*(privacy.RingSize-1))
	if err != nil {
		return nil, nil, err
	}
	kvArgs[utils.MyIndices] = myIndices

	return coinsToSpend, kvArgs, nil

}

// GetTokenFee returns the token fee per kb.
func (client *IncClient) GetTokenFee(shardID byte, tokenIDStr string) (uint64, error) {
	if tokenIDStr == common.PRVIDStr {
		return DefaultPRVFee, nil
	}
	responseInBytes, err := client.rpcServer.EstimateFeeWithEstimator(-1, shardID, 10, tokenIDStr)
	if err != nil {
		return 0, err
	}

	var feeEstimateResult rpc.EstimateFeeResult
	err = rpchandler.ParseResponse(responseInBytes, &feeEstimateResult)
	if err != nil {
		return 0, err
	}

	return feeEstimateResult.EstimateFeeCoinPerKb, nil

}

// GetTxDetail retrieves the transaction detail from its hash.
func (client *IncClient) GetTxDetail(txHash string) (*jsonresult.TransactionDetail, error) {
	responseInBytes, err := client.rpcServer.GetTransactionByHash(txHash)
	if err != nil {
		return nil, err
	}

	var txDetail jsonresult.TransactionDetail
	err = rpchandler.ParseResponse(responseInBytes, &txDetail)
	if err != nil {
		return nil, err
	}

	return &txDetail, err
}

// GetTx retrieves the transaction detail and parses it to a transaction object.
func (client *IncClient) GetTx(txHash string) (metadata.Transaction, error) {
	txDetail, err := client.GetTxDetail(txHash)
	if err != nil {
		return nil, err
	}

	return jsonresult.ParseTxDetail(*txDetail)
}

// GetTxs retrieves transactions and parses them to transaction objects given their hashes.
func (client *IncClient) GetTxs(txHashList []string) (map[string]metadata.Transaction, error) {
	responseInBytes, err := client.rpcServer.GetEncodedTransactionsByHashes(txHashList)
	if err != nil {
		return nil, err
	}

	mapRes := make(map[string]string)
	err = rpchandler.ParseResponse(responseInBytes, &mapRes)
	if err != nil {
		panic(err)
	}

	res := make(map[string]metadata.Transaction)
	for txHash, encodedTx := range mapRes {
		txBytes, _, err := base58.Base58Check{}.Decode(encodedTx)
		if err != nil {
			Logger.Printf("base58-decode failed: %v\n", string(txBytes))
			return nil, err
		}

		txChoice, err := transaction.DeserializeTransactionJSON(txBytes)
		if err != nil {
			Logger.Printf("unMarshal failed: %v\n", string(txBytes))
			return nil, err
		}
		tx := txChoice.ToTx()

		if tx.Hash().String() != txHash {
			Logger.Printf("txParseFail: %v\n", string(txBytes))
			return nil, fmt.Errorf("txHash changes after unmarshalling, expect %v, got %v", txHash, tx.Hash().String())
		}
		res[txHash] = tx
	}

	return res, nil
}

// GetTransactionHashesByReceiver retrieves the list of all transactions received by a payment address.
func (client *IncClient) GetTransactionHashesByReceiver(paymentAddress string) ([]string, error) {
	responseInBytes, err := client.rpcServer.GetTxHashByReceiver(paymentAddress)
	if err != nil {
		return nil, err
	}

	var tmpRes map[string][]string
	err = rpchandler.ParseResponse(responseInBytes, &tmpRes)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, txList := range tmpRes {
		res = append(res, txList...)
	}

	return res, nil
}

// GetTransactionsByReceiver retrieves the list of all transactions (in object) received by a payment address.
//
// Notice that this function is time-consuming since it has to parse every single transaction into an object.
func (client *IncClient) GetTransactionsByReceiver(paymentAddress string) (map[string]metadata.Transaction, error) {
	txList, err := client.GetTransactionHashesByReceiver(paymentAddress)
	if err != nil {
		return nil, err
	}

	fmt.Printf("#Txs: %v\n", len(txList))

	count := 0
	start := time.Now()
	res := make(map[string]metadata.Transaction)
	for _, txHash := range txList {
		tx, err := client.GetTx(txHash)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve tx %v: %v", txHash, err)
		}
		res[txHash] = tx
		count += 1
		if count%5 == 0 {
			Logger.Printf("count %v, timeElapsed: %v\n", count, time.Since(start).Seconds())
		}
	}

	return res, nil
}

// GetTxHashByPublicKeys retrieves the list of all transactions' hash sent to a list of public keys.
func (client *IncClient) GetTxHashByPublicKeys(publicKeys []string) (map[string][]string, error) {
	responseInBytes, err := client.rpcServer.GetTxHashByPublicKey(publicKeys)
	if err != nil {
		return nil, err
	}

	tmpRes := make(map[string]map[byte][]string)
	err = rpchandler.ParseResponse(responseInBytes, &tmpRes)
	if err != nil {
		return nil, err
	}

	res := make(map[string][]string)
	for publicKeyStr, txMap := range tmpRes {
		txList := make([]string, 0)
		for _, tmpTxList := range txMap {
			txList = append(txList, tmpTxList...)
		}
		res[publicKeyStr] = txList
	}

	return res, nil
}

// GetTransactionsByPublicKeys retrieves the list of all transactions (in object) sent to a list of base58-encoded public keys.
//
// Notice that this function is time-consuming since it has to parse every single transaction into an object.
func (client *IncClient) GetTransactionsByPublicKeys(publicKeys []string) (map[string]map[string]metadata.Transaction, error) {
	txMap, err := client.GetTxHashByPublicKeys(publicKeys)
	if err != nil {
		return nil, err
	}

	res := make(map[string]map[string]metadata.Transaction)
	for publicKeyStr, txList := range txMap {
		tmpRes := make(map[string]metadata.Transaction)
		for current := 0; current < len(txList); current += pageSize {
			next := current + pageSize
			if next > len(txList) {
				next = len(txList)
			}

			mapRes, err := client.GetTxs(txList[current:next])
			if err != nil {
				return nil, err
			}

			for txHash, tx := range mapRes {
				tmpRes[txHash] = tx
			}
		}

		res[publicKeyStr] = tmpRes
	}

	return res, nil
}

// GetTxHashBySerialNumbers retrieves the list of tokenIDStr transactions in which serial numbers have been spent.
//
// Set shardID = 255 to retrieve in all shards.
func (client *IncClient) GetTxHashBySerialNumbers(snList []string, tokenIDStr string, shardID byte) (map[string]string, error) {
	responseInBytes, err := client.rpcServer.GetTxHashBySerialNumber(snList, tokenIDStr, shardID)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CheckTxInBlock checks if a transaction has been included in a block or not.
func (client *IncClient) CheckTxInBlock(txHash string) (bool, error) {
	txDetail, err := client.GetTxDetail(txHash)
	if err != nil {
		return false, err
	}

	if txDetail.IsInMempool {
		return false, nil
	}

	return txDetail.IsInBlock, nil
}

// GetReceivingInfo verifies if a transaction is sent to the given `otaKey` and the transacted tokenIds.
// Furthermore, in case a read-only key is given, it will
// decrypt the received output coins and return the total amounts. If there are multiple read-only keys,
// only the first one is used.
func (client *IncClient) GetReceivingInfo(
	txHash string,
	otaKey string,
	readonlyKey ...string,
) (received bool, mapResult map[string]uint64, err error) {
	mapResult = make(map[string]uint64)

	// deserialize the ota key
	w, err := wallet.Base58CheckDeserialize(otaKey)
	if err != nil || w.KeySet.OTAKey.GetOTASecretKey() == nil || w.KeySet.OTAKey.GetPublicSpend() == nil {
		err = fmt.Errorf("otaKey is invalid: %v", err)
		return
	}
	keySet := w.KeySet
	keySet.PaymentAddress = key.PaymentAddress{Pk: keySet.OTAKey.GetPublicSpend().ToBytesS()}

	// deserialize the ota key (if have)
	if len(readonlyKey) > 0 {
		tmpWallet, tmpErr := wallet.Base58CheckDeserialize(readonlyKey[0])
		if tmpErr != nil ||
			tmpWallet.KeySet.ReadonlyKey.GetPublicSpend() == nil ||
			tmpWallet.KeySet.ReadonlyKey.GetPrivateView() == nil {
			err = fmt.Errorf("readonlyKey is invalid: %v", tmpErr)
			return
		}
		keySet.ReadonlyKey = tmpWallet.KeySet.ReadonlyKey
	}

	// get the transaction detail
	tmpTxs, err := client.GetTxs([]string{txHash})
	if err != nil {
		return
	}
	tx := tmpTxs[txHash]
	tokenIdStr := tx.GetTokenID().String()

	// get the output coins
	outCoins := make([]coin.Coin, 0)
	switch tx.GetType() {
	case common.TxCustomTokenPrivacyType, common.TxTokenConversionType:
		txToken, ok := tx.(tx_generic.TransactionToken)
		if !ok {
			err = fmt.Errorf("cannot parse tx as a token transaction")
			return
		}
		// get the PRV amount (if have)
		if txToken.GetTxBase() != nil {
			prvAmount, err := getTxOutputAmountByKeySet(txToken, common.PRVIDStr, &keySet)
			if err != nil {
				Logger.Printf("get PRV amount error: %v\n", err)
			}
			if prvAmount > 0 {
				received = true
			}
			mapResult[common.PRVIDStr] = prvAmount
		}

		txNormal := txToken.GetTxNormal()
		if txNormal.GetProof() != nil && txNormal.GetProof().GetOutputCoins() != nil {
			outCoins = append(outCoins, txNormal.GetProof().GetOutputCoins()...)
		}
	case common.TxNormalType, common.TxRewardType, common.TxReturnStakingType, common.TxConversionType:
		prvAmount, err := getTxOutputAmountByKeySet(tx, common.PRVIDStr, &keySet)
		if err != nil {
			Logger.Printf("get PRV amount error: %v\n", err)
		}
		if prvAmount > 0 {
			received = true
		}
		mapResult[common.PRVIDStr] = prvAmount
	default:
		err = fmt.Errorf("transaction type `%v` is invalid", tx.GetType())
	}

	if len(outCoins) == 0 {
		err = fmt.Errorf("transaction does not have output coins")
	}

	// getAssetTags
	assetTags, err := client.GetAllAssetTags()
	if err != nil {
		return
	}

	// check if there is any output coins belong to the `keySet`, and decrypt it if there is a read-only key.
	var plainCoin coin.PlainCoin
	var tmpTokenId *common.Hash
	for _, outCoin := range outCoins {
		belong, _ := outCoin.DoesCoinBelongToKeySet(&keySet)
		if belong {
			received = true

			// try to decrypt first
			amount := uint64(0)
			plainCoin, _ = outCoin.Decrypt(&keySet)
			if plainCoin != nil {
				amount = plainCoin.GetValue()
			}

			switch tokenIdStr {
			case common.ConfidentialAssetID.String():
				if tmpTokenId == nil {
					tmpTokenId, err = outCoin.(*coin.CoinV2).GetTokenId(&keySet, assetTags)
					if err != nil {
						return
					}
				}
				mapResult[tmpTokenId.String()] += amount
			default:
				mapResult[tokenIdStr] += amount
			}
		}
	}

	return
}