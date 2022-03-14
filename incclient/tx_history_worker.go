package incclient

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// TxHistoryProcessor implements a processor for retrieving transaction history in a parallel manner.
//
// Each TxHistoryProcessor consists of several TxHistoryWorker's that help retrieve transaction history faster.
type TxHistoryProcessor struct {
	client    *IncClient
	mtx       *sync.RWMutex
	history   *TxHistory
	errChan   chan error
	txChan    chan TxHistory
	workers   []*TxHistoryWorker
	cachedTxs map[string]metadata.Transaction
}

// NewTxHistoryProcessor creates a TxHistoryProcess with a number of TxHistoryWorker's.
func NewTxHistoryProcessor(client *IncClient, numWorkers int) *TxHistoryProcessor {
	mtx := new(sync.RWMutex)
	errChan := make(chan error, numWorkers)
	txChan := make(chan TxHistory, numWorkers)
	workers := make([]*TxHistoryWorker, 0)
	history := new(TxHistory)
	history.TxInList = make([]TxIn, 0)
	history.TxOutList = make([]TxOut, 0)

	for i := 0; i < numWorkers; i++ {
		worker := NewTxHistoryWorker(i, client)
		workers = append(workers, worker)
	}

	return &TxHistoryProcessor{
		client:  client,
		mtx:     mtx,
		history: history,
		errChan: errChan,
		txChan:  txChan,
		workers: workers,
	}
}

func (p *TxHistoryProcessor) addHistory(history TxHistory) {
	p.mtx.Lock()

	mapTxIns := make(map[string]bool)
	for _, txIn := range p.history.TxInList {
		mapTxIns[txIn.TxHash] = true
	}

	mapTxOuts := make(map[string]bool)
	for _, txOut := range p.history.TxOutList {
		mapTxOuts[txOut.TxHash] = true
	}

	for _, txIn := range history.TxInList {
		if mapTxIns[txIn.TxHash] {
			continue
		}
		p.history.TxInList = append(p.history.TxInList, txIn)
	}

	for _, txOut := range history.TxOutList {
		if mapTxOuts[txOut.TxHash] {
			continue
		}
		p.history.TxOutList = append(p.history.TxOutList, txOut)
	}
	Logger.Printf("Added %v TxsIn, %v TxsOut\n", len(history.TxInList), len(history.TxOutList))
	Logger.Printf("Current history: #TxsIn = %v, #TxsOut = %v\n", len(p.history.TxInList), len(p.history.TxOutList))

	p.mtx.Unlock()
}

// GetTxsIn returns the list of in-coming transactions in a parallel manner.
func (p *TxHistoryProcessor) GetTxsIn(privateKey string, tokenIDStr string, version int8) ([]TxIn, error) {
	kWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot deserialize private key %v: %v", privateKey, err)
	}
	addrStr := PrivateKeyToPaymentAddress(privateKey, -1)

	listDecryptedCoins, err := p.client.GetListDecryptedOutCoin(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, err
	}

	numWorkers := 0
	if version == 1 {
		txList, err := p.client.GetTransactionHashesByReceiver(addrStr)
		if err != nil {
			return nil, err
		}

		if len(txList) < len(p.workers) {
			numWorkers = 1
			go p.workers[0].getTxsInV1(&kWallet.KeySet, listDecryptedCoins,
				txList, tokenIDStr, p.txChan, p.errChan)
		} else {
			numWorkers = len(p.workers)
			//calculate the number of txs each worker has to handle
			numForEach := len(txList) / len(p.workers)

			//call each worker to retrieve history
			for i := 0; i < len(p.workers); i++ {
				start := i * numForEach
				end := (i + 1) * numForEach
				if i == len(p.workers)-1 {
					end = len(txList)
				}
				go p.workers[i].getTxsInV1(&kWallet.KeySet, listDecryptedCoins,
					txList[start:end], tokenIDStr, p.txChan, p.errChan)
			}
		}
	} else if version == 2 {
		pubKeys := make([]string, 0)
		for _, outCoin := range listDecryptedCoins {
			pubKeys = append(pubKeys, base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), 0))
		}

		if len(pubKeys) < len(p.workers) {
			numWorkers = 1
			go p.workers[0].getTxsInV2(&kWallet.KeySet, listDecryptedCoins,
				pubKeys, tokenIDStr, p.txChan, p.errChan)
		} else {
			numWorkers = len(p.workers)
			//calculate the number of txs each worker has to handle
			numForEach := len(pubKeys) / len(p.workers)

			//call each worker to retrieve history
			for i := 0; i < len(p.workers); i++ {
				start := i * numForEach
				end := (i + 1) * numForEach
				if i == len(p.workers)-1 {
					end = len(pubKeys)
				}
				go p.workers[i].getTxsInV2(&kWallet.KeySet, listDecryptedCoins,
					pubKeys[start:end], tokenIDStr, p.txChan, p.errChan)
			}
		}
	}

	numSuccess := 0
	for {
		select {
		case err := <-p.errChan:
			return p.history.TxInList, err
		case txHistory := <-p.txChan:
			numSuccess++
			p.addHistory(txHistory)
			if numSuccess == numWorkers {
				sort.Slice(p.history.TxInList, func(i, j int) bool {
					return p.history.TxInList[i].LockTime > p.history.TxInList[j].LockTime
				})
				return p.history.TxInList, nil
			}
			Logger.Printf("Receive new data, numSuccess = %v\n", numSuccess)
		}
	}

}

// GetTxsOut returns the list of out-going transactions in a parallel manner.
func (p *TxHistoryProcessor) GetTxsOut(privateKey string, tokenIDStr string, version int8) ([]TxOut, error) {
	kWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot deserialize private key %v: %v", privateKey, err)
	}

	listSpentCoins, _, err := p.client.GetSpentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return nil, err
	}

	// Create a map from serial numbers to coins
	mapSpentCoins := make(map[string]coin.PlainCoin)
	snList := make([]string, 0)
	for _, spentCoin := range listSpentCoins {
		if spentCoin.GetVersion() == uint8(version) {
			snStr := base58.Base58Check{}.Encode(spentCoin.GetKeyImage().ToBytesS(), common.ZeroByte)
			mapSpentCoins[snStr] = spentCoin
			snList = append(snList, snStr)
		}
	}

	if len(snList) == 0 {
		return nil, nil
	}

	Logger.Printf("len(snList) = %v\n", len(snList))

	numWorkers := 0
	if len(snList) < len(p.workers) {
		numWorkers = 1
		go p.workers[0].getTxsOut(&kWallet.KeySet, mapSpentCoins,
			snList, tokenIDStr, p.txChan, p.errChan)
	} else {
		numWorkers = len(p.workers)
		//calculate the number of txs each worker has to handle
		numForEach := len(snList) / len(p.workers)

		//call each worker to retrieve history
		for i := 0; i < len(p.workers); i++ {
			start := i * numForEach
			end := (i + 1) * numForEach
			if i == len(p.workers)-1 {
				end = len(snList)
			}
			go p.workers[i].getTxsOut(&kWallet.KeySet, mapSpentCoins,
				snList[start:end], tokenIDStr, p.txChan, p.errChan)
		}
	}

	numSuccess := 0
	for {
		select {
		case err := <-p.errChan:
			return p.history.TxOutList, err
		case txHistory := <-p.txChan:
			numSuccess++
			p.addHistory(txHistory)
			if numSuccess == numWorkers {
				sort.Slice(p.history.TxOutList, func(i, j int) bool {
					return p.history.TxOutList[i].LockTime > p.history.TxOutList[j].LockTime
				})
				return p.history.TxOutList, nil
			}
			Logger.Printf("Receive new data, numSuccess = %v\n", numSuccess)
		}
	}

}

// GetTokenHistory returns the history of a private key w.r.t a tokenID in a parallel manner.
func (p *TxHistoryProcessor) GetTokenHistory(privateKey string, tokenIDStr string) (*TxHistory, error) {
	Logger.Printf("GETTING in-coming v1 txs for token %v\n", tokenIDStr)
	txsIn, err := p.GetTxsIn(privateKey, tokenIDStr, 1)
	if err != nil {
		return nil, err
	}
	Logger.Printf("FINISHED in-coming v1 txs for token %v\n\n", tokenIDStr)

	Logger.Printf("GETTING out-going v1 txs for token %v\n", tokenIDStr)
	txsOut, err := p.GetTxsOut(privateKey, tokenIDStr, 1)
	if err != nil {
		return nil, err
	}
	Logger.Printf("FINISHED out-going v1 txs for token %v\n\n", tokenIDStr)

	Logger.Printf("GETTING in-coming v2 txs for token %v\n", tokenIDStr)
	txsInV2, err := p.GetTxsIn(privateKey, tokenIDStr, 2)
	if err != nil {
		return nil, err
	}
	Logger.Printf("FINISHED in-coming v2 txs for token %v\n\n", tokenIDStr)
	for _, txInV2 := range txsInV2 {
		txsIn = append(txsIn, txInV2)
	}
	sort.Slice(txsIn, func(i, j int) bool {
		return txsIn[i].LockTime > txsIn[j].LockTime
	})

	Logger.Printf("GETTING out-going v2 txs for token %v\n", tokenIDStr)
	txsOutV2, err := p.GetTxsOut(privateKey, tokenIDStr, 2)
	if err != nil {
		return nil, err
	}
	Logger.Printf("FINISHED out-going v2 txs for token %v\n\n", tokenIDStr)
	for _, txOutV2 := range txsOutV2 {
		txsOut = append(txsOut, txOutV2)
	}
	sort.Slice(txsOut, func(i, j int) bool {
		return txsOut[i].LockTime > txsOut[j].LockTime
	})

	return &TxHistory{
		TxInList:  txsIn,
		TxOutList: txsOut,
	}, nil

}

// GetAllHistory returns all the history of an account in a parallel manner.
func (p *TxHistoryProcessor) GetAllHistory(privateKeyStr string) (map[string]*TxHistory, error) {
	prefix := "[GetAllHistory]"
	Logger.Printf("%v STARTING...\n", prefix)
	res := make(map[string]*TxHistory)

	tokenIDs, err := p.client.getAllTokens(privateKeyStr)
	if err != nil {
		return nil, err
	}

	Logger.Printf("%v #TokenIDs: %v\n", prefix, len(tokenIDs))
	finishedCount := 0
	for _, tokenID := range tokenIDs {
		res[tokenID], err = p.GetTokenHistory(privateKeyStr, tokenID)
		if err != nil {
			return nil, err
		}
		finishedCount++
		Logger.Printf("%v Finished token %v, count: %v/%v\n", prefix, tokenID, finishedCount, len(tokenIDs))
	}

	Logger.Printf("%v FINISHED ALL\n\n", prefix)

	return res, nil
}

// TxHistoryWorker implements a worker for retrieving transaction history.
type TxHistoryWorker struct {
	id        int
	client    *IncClient
	cachedTxs map[string]metadata.Transaction
}

// NewTxHistoryWorker creates a new TxHistoryWorker.
func NewTxHistoryWorker(id int, client *IncClient) *TxHistoryWorker {
	return &TxHistoryWorker{
		id:     id,
		client: client,
	}
}

// getListTxs returns a list of transactions (in object) on input a list of transaction hashes.
func (worker TxHistoryWorker) getListTxs(txList []string) (map[string]metadata.Transaction, error) {
	count := 0
	start := time.Now()

	res := make(map[string]metadata.Transaction)
	for current := 0; current < len(txList); current += pageSize {
		next := current + pageSize
		if next > len(txList) {
			next = len(txList)
		}

		txMap, err := worker.client.GetTxs(txList[current:next])
		if err != nil {
			return nil, err
		}

		for txHash, tx := range txMap {
			res[txHash] = tx
		}
		count += len(txMap)
		Logger.Printf("[WORKER %v], count %v, timeElapsed %v\n", worker.id, count, time.Since(start).Seconds())
	}

	return res, nil
}

// getTxsInV1 returns the list of in-coming transactions of version 1.
//
// It only returns the list of transactions whose value is greater than 0.
func (worker TxHistoryWorker) getTxsInV1(keySet *key.KeySet, listDecryptedCoins map[string]coin.PlainCoin, txList []string, tokenIDStr string, txChan chan TxHistory, errChan chan error) {
	Logger.Printf("[WORKER %v] getTxsInV1, #TXS: %v\n", worker.id, len(txList))
	mapCmt := makeMapCMToPlainCoin(listDecryptedCoins)

	//retrieve transactions in object
	txMap, err := worker.getListTxs(txList)
	if err != nil {
		errChan <- err
		return
	}

	res := make([]TxIn, 0)
	for txHash, tx := range txMap {
		if isOut, err := isTxOut(tx, tokenIDStr, listDecryptedCoins); err != nil {
			errChan <- err
			return
		} else if isOut {
			continue
		}

		outCoins, err := getTxOutputCoinsByKeySet(tx, tokenIDStr, keySet)
		if err != nil {
			errChan <- err
			return
		}

		amount := uint64(0)
		for cmtStr := range outCoins {
			if outCoin, ok := mapCmt[cmtStr]; ok {
				amount += outCoin.GetValue()
				continue
			}
		}
		if amount > 0 {
			newTxIn := TxIn{
				Version:  tx.GetVersion(),
				LockTime: tx.GetLockTime(),
				TxHash:   txHash,
				TokenID:  tx.GetTokenID().String(),
				Metadata: tx.GetMetadata(),
				Note:     txMetadataNote[tx.GetMetadataType()],
			}
			newTxIn.Amount = amount
			res = append(res, newTxIn)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})

	txChan <- TxHistory{
		TxInList:  res,
		TxOutList: nil,
	}
	Logger.Printf("[WORKER %v] FINISHED getTxsInV1, #TXS: %v!!\n\n", worker.id, len(res))
}

// getTxsInV2 returns the list of in-coming transactions of version 2.
//
// It only returns the list of transactions whose value is greater than 0.
func (worker TxHistoryWorker) getTxsInV2(keySet *key.KeySet, listDecryptedCoins map[string]coin.PlainCoin, publicKeys []string, tokenIDStr string, txChan chan TxHistory, errChan chan error) {
	res := make([]TxIn, 0)
	if worker.client.version != 2 {
		txChan <- TxHistory{
			TxInList:  res,
			TxOutList: nil,
		}
	}

	Logger.Printf("[WORKER %v] getTxsInV2, #No: %v\n", worker.id, len(publicKeys))
	mapCmt := makeMapCMToPlainCoin(listDecryptedCoins)

	//retrieve transactions in object
	txMap, err := worker.client.GetTransactionsByPublicKeys(publicKeys)
	if err != nil {
		errChan <- err
		return
	}

	mapRes := make(map[string]TxIn)
	for _, tmpTxMap := range txMap {
		for txHash, tx := range tmpTxMap {
			if _, ok := mapRes[txHash]; ok {
				continue
			}
			if isOut, err := isTxOut(tx, tokenIDStr, listDecryptedCoins); err != nil {
				errChan <- err
				return
			} else if isOut {
				continue
			}

			outCoins, err := getTxOutputCoinsByKeySet(tx, tokenIDStr, keySet)
			if err != nil {
				errChan <- err
				return
			}

			pubKeys := make(map[string]uint64)
			amount := uint64(0)
			for cmtStr := range outCoins {
				if outCoin, ok := mapCmt[cmtStr]; ok {
					amount += outCoin.GetValue()
					pubKeys[base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), 0)] = outCoin.GetValue()
					continue
				}
			}
			if amount > 0 {
				newTxIn := TxIn{
					Version:  tx.GetVersion(),
					LockTime: tx.GetLockTime(),
					OutCoins: pubKeys,
					TxHash:   txHash,
					TokenID:  tx.GetTokenID().String(),
					Amount:   amount,
					Metadata: tx.GetMetadata(),
					Note:     txMetadataNote[tx.GetMetadataType()],
				}
				mapRes[txHash] = newTxIn
			}
		}
	}

	for _, txIn := range mapRes {
		res = append(res, txIn)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})

	txChan <- TxHistory{
		TxInList:  res,
		TxOutList: nil,
	}
	Logger.Printf("[WORKER %v] FINISHED getTxsInV2, #TXS: %v!!\n\n", worker.id, len(res))
}

// getTxsOut returns the list of out-going transactions of version 1.
//
// It only returns the list of transactions whose value is greater than 0.
func (worker TxHistoryWorker) getTxsOut(keySet *key.KeySet, mapSpentCoins map[string]coin.PlainCoin, snList []string, tokenIDStr string, txChan chan TxHistory, errChan chan error) {
	Logger.Printf("[WORKER %v] getTxsOut, #No: %v\n", worker.id, len(snList))

	shardID := common.GetShardIDFromLastByte(keySet.PaymentAddress.Pk[len(keySet.PaymentAddress.Pk)-1])

	// Retrieve the list of transactions which spent these coins
	mapSpentTxs, err := worker.client.GetTxHashBySerialNumbers(snList, tokenIDStr, shardID)
	if err != nil {
		if strings.Contains(err.Error(), "Method not found") {
			errChan <- fmt.Errorf("method not supported by the remote node configurations")
		}
		errChan <- err
		return
	}

	// Create a list of txs
	txHashList := make([]string, 0)
	for _, txHash := range mapSpentTxs {
		txHashList = append(txHashList, txHash)
	}
	// Get txs from hashes
	txs, err := worker.getListTxs(txHashList)
	if err != nil {
		errChan <- err
		return
	}

	mapRes := make(map[string]TxOut)
	res := make([]TxOut, 0)

	var ok bool
	for _, txHash := range mapSpentTxs {
		// check if the txHash has been processed
		if _, ok = mapRes[txHash]; ok {
			continue
		}

		var tx metadata.Transaction
		tx, ok = txs[txHash]
		if !ok {
			errChan <- fmt.Errorf("tx %v not found", txHash)
			return
		}

		//get transaction fee
		fee, isPRVFee := getTxFeeBy(tx)

		//calculate transaction's amount
		inputAmount, spentCoins, err := getTxInputAmount(tx, tokenIDStr, mapSpentCoins)
		if err != nil {
			errChan <- err
			return
		}
		outputAmount, err := getTxOutputAmountByKeySet(tx, tokenIDStr, keySet)
		if err != nil {
			errChan <- err
			return
		}
		amount := inputAmount - outputAmount
		if isPRVFee && tokenIDStr == common.PRVIDStr {
			amount -= fee
		}
		if !isPRVFee && tokenIDStr != common.PRVIDStr {
			amount -= fee
		}

		//get list of receivers' public keys
		receivers, err := getTxReceivers(tx, tokenIDStr)
		if err != nil {
			errChan <- err
			return
		}

		if amount > 0 || tokenIDStr == common.PRVIDStr {
			note := txMetadataNote[tx.GetMetadataType()]
			if tokenIDStr == common.PRVIDStr && amount == 0 {
				note += " (Tx Fee)"
			}
			note = strings.TrimSpace(note)

			newTxOut := TxOut{
				Version:    tx.GetVersion(),
				LockTime:   tx.GetLockTime(),
				TxHash:     txHash,
				TokenID:    tx.GetTokenID().String(),
				SpentCoins: spentCoins,
				Receivers:  receivers,
				Amount:     amount,
				Metadata:   tx.GetMetadata(),
				PRVFee:     fee,
				Note:       note,
			}

			if !isPRVFee {
				newTxOut.PRVFee = 0
				newTxOut.TokenFee = fee
			}

			mapRes[txHash] = newTxOut
			res = append(res, newTxOut)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].LockTime > res[j].LockTime
	})
	txChan <- TxHistory{
		TxInList:  nil,
		TxOutList: res,
	}
	Logger.Printf("[WORKER %v] FINISHED getTxsOut, #TXS: %v!!\n\n", worker.id, len(res))
}
