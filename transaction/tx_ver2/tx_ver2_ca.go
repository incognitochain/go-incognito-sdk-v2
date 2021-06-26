package tx_ver2

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v2/mlsag"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"log"
	"math/big"
)

// Create unique OTA coin without the help of the db
func createUniqueOTACoinCA(paymentInfo *key.PaymentInfo, tokenID *common.Hash) (*coin.CoinV2, *crypto.Point, error) {
	if tokenID == nil {
		tokenID = &common.PRVCoinID
	}
	c, sharedSecret, err := coin.NewCoinCA(paymentInfo, tokenID)
	if err != nil {
		log.Printf("Cannot parse coin based on payment info err: %v", err)
		return nil, nil, err
	}
	// If previously created coin is burning address
	if sharedSecret == nil {
		// assetTag := privacy.HashToPoint(tokenID[:])
		// c.SetAssetTag(assetTag)
		return c, nil, nil // No need to check db
	}
	return c, sharedSecret, nil
}

func createPrivateKeyMlsagCA(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, outputSharedSecrets []*crypto.Point, params *tx_generic.TxPrivacyInitParams, _ byte, commitmentsToZero []*crypto.Point) ([]*crypto.Scalar, error) {
	senderSK := params.SenderSK
	// db := params.StateDB
	tokenID := params.TokenID
	if tokenID == nil {
		tokenID = &common.PRVCoinID
	}
	rehashed := crypto.HashToPoint(tokenID[:])
	sumRand := new(crypto.Scalar).FromUint64(0)

	privateKeyMlsag := make([]*crypto.Scalar, len(inputCoins)+2)
	sumInputAssetTagBlinders := new(crypto.Scalar).FromUint64(0)
	numOfInputs := new(crypto.Scalar).FromUint64(uint64(len(inputCoins)))
	numOfOutputs := new(crypto.Scalar).FromUint64(uint64(len(outputCoins)))
	mySkBytes := (*senderSK)[:]
	for i := 0; i < len(inputCoins); i += 1 {
		var err error
		privateKeyMlsag[i], err = inputCoins[i].ParsePrivateKeyOfCoin(*senderSK)
		if err != nil {
			log.Printf("Cannot parse private key of coin %v\n", err)
			return nil, err
		}

		tmpInCoin, ok := inputCoins[i].(*coin.CoinV2)
		if !ok || tmpInCoin.GetAssetTag() == nil {
			return nil, fmt.Errorf("cannot cast a coin as v2-CA")
		}

		isUnBlinded := crypto.IsPointEqual(rehashed, tmpInCoin.GetAssetTag())

		sharedSecret := new(crypto.Point).Identity()
		bl := new(crypto.Scalar).FromUint64(0)
		if !isUnBlinded {
			sharedSecret, err = tmpInCoin.RecomputeSharedSecret(mySkBytes)
			if err != nil {
				log.Printf("cannot recompute shared secret: %v\n", err)
				return nil, err
			}

			bl, err = coin.ComputeAssetTagBlinder(sharedSecret)
			if err != nil {
				return nil, err
			}
		}

		v := tmpInCoin.GetAmount()
		effectiveRCom := new(crypto.Scalar).Mul(bl, v)
		effectiveRCom.Add(effectiveRCom, tmpInCoin.GetRandomness())

		sumInputAssetTagBlinders.Add(sumInputAssetTagBlinders, bl)
		sumRand.Add(sumRand, effectiveRCom)
	}
	sumInputAssetTagBlinders.Mul(sumInputAssetTagBlinders, numOfOutputs)

	sumOutputAssetTagBlinders := new(crypto.Scalar).FromUint64(0)

	var err error
	for i, oc := range outputCoins {
		if oc.GetAssetTag() == nil {
			return nil, fmt.Errorf("cannot cast a coin as v2-CA")
		}
		// lengths between 0 and len(outputCoins) were rejected before
		bl := new(crypto.Scalar).FromUint64(0)
		isUnBlinded := crypto.IsPointEqual(rehashed, oc.GetAssetTag())
		if !isUnBlinded {
			bl, err = coin.ComputeAssetTagBlinder(outputSharedSecrets[i])
			if err != nil {
				return nil, err
			}
		}

		v := oc.GetAmount()
		effectiveRCom := new(crypto.Scalar).Mul(bl, v)
		effectiveRCom.Add(effectiveRCom, oc.GetRandomness())
		sumOutputAssetTagBlinders.Add(sumOutputAssetTagBlinders, bl)
		sumRand.Sub(sumRand, effectiveRCom)
	}
	sumOutputAssetTagBlinders.Mul(sumOutputAssetTagBlinders, numOfInputs)

	// 2 final elements in `private keys` for MLSAG
	assetSum := new(crypto.Scalar).Sub(sumInputAssetTagBlinders, sumOutputAssetTagBlinders)
	firstCommitmentToZeroRecomputed := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenRandomnessIndex], assetSum)
	secondCommitmentToZeroRecomputed := new(crypto.Point).ScalarMult(crypto.PedCom.G[crypto.PedersenRandomnessIndex], sumRand)
	if len(commitmentsToZero) != 2 {
		log.Printf("Received %d points to check when signing MLSAG\n", len(commitmentsToZero))
		return nil, utils.NewTransactionErr(utils.UnexpectedError, fmt.Errorf("need exactly 2 points for MLSAG double-checking"))
	}
	match1 := crypto.IsPointEqual(firstCommitmentToZeroRecomputed, commitmentsToZero[0])
	match2 := crypto.IsPointEqual(secondCommitmentToZeroRecomputed, commitmentsToZero[1])
	if !match1 || !match2 {
		return nil, utils.NewTransactionErr(utils.UnexpectedError, fmt.Errorf("asset tag sum or commitment sum mismatch"))
	}

	privateKeyMlsag[len(inputCoins)] = assetSum
	privateKeyMlsag[len(inputCoins)+1] = sumRand
	return privateKeyMlsag, nil
}

func generateMlsagRingWithIndexesCA(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, params *tx_generic.TxPrivacyInitParams, pi int, _ byte, ringSize int) (*mlsag.Ring, [][]*big.Int, []*crypto.Point, error) {
	cmtIndices, myIndices, commitments, publicKeys, assetTags, err := parseParamsForRing(params.KvArgs, len(inputCoins), ringSize)
	if err != nil {
		return nil, nil, nil, utils.NewTransactionErr(utils.UnexpectedError, fmt.Errorf("parseParamsForRing error: %v", err))
	}
	if len(assetTags) < len(inputCoins)*(ringSize-1) {
		return nil, nil, nil, fmt.Errorf("not enough decoy asset tags: have %v, need at least %v (%v input coins)", len(assetTags), len(inputCoins)*(ringSize-1), len(inputCoins))
	}

	outputCoinsAsGeneric := make([]coin.Coin, len(outputCoins))
	for i := 0; i < len(outputCoins); i++ {
		outputCoinsAsGeneric[i] = outputCoins[i]
	}
	sumOutputsWithFee := tx_generic.CalculateSumOutputsWithFee(outputCoinsAsGeneric, params.Fee)
	inCount := new(crypto.Scalar).FromUint64(uint64(len(inputCoins)))
	outCount := new(crypto.Scalar).FromUint64(uint64(len(outputCoins)))

	sumOutputAssetTags := new(crypto.Point).Identity()
	for _, oc := range outputCoins {
		if oc.GetAssetTag() == nil {
			log.Printf("CA error: missing asset tag for signing in output coin - %v\n", oc.Bytes())
			err := utils.NewTransactionErr(utils.SignTxError, fmt.Errorf("cannot sign CA token : an output coin does not have asset tag"))
			return nil, nil, nil, err
		}
		sumOutputAssetTags.Add(sumOutputAssetTags, oc.GetAssetTag())
	}
	sumOutputAssetTags.ScalarMult(sumOutputAssetTags, inCount)

	indexes := make([][]*big.Int, ringSize)
	ring := make([][]*crypto.Point, ringSize)
	var lastTwoColumnsCommitmentToZero []*crypto.Point
	currentIndex := 0
	for i := 0; i < ringSize; i += 1 {
		sumInputs := new(crypto.Point).Identity()
		sumInputs.Sub(sumInputs, sumOutputsWithFee)
		sumInputAssetTags := new(crypto.Point).Identity()

		row := make([]*crypto.Point, len(inputCoins))
		rowIndexes := make([]*big.Int, len(inputCoins))
		if i == pi {
			for j := 0; j < len(inputCoins); j += 1 {
				row[j] = inputCoins[j].GetPublicKey()
				rowIndexes[j] = new(big.Int).SetUint64(myIndices[j])
				sumInputs.Add(sumInputs, inputCoins[j].GetCommitment())
				tmpInCoin, ok := inputCoins[j].(*coin.CoinV2)
				if !ok {
					return nil, nil, nil, fmt.Errorf("cannot cast a coin as v2")
				}
				if tmpInCoin.GetAssetTag() == nil {
					log.Printf("CA error: missing asset tag for signing in input coin - %v\n", tmpInCoin.Bytes())
					err := utils.NewTransactionErr(utils.SignTxError, fmt.Errorf("cannot sign CA token : an input coin does not have asset tag"))
					return nil, nil, nil, err
				}
				sumInputAssetTags.Add(sumInputAssetTags, tmpInCoin.GetAssetTag())
			}
		} else {
			for j := 0; j < len(inputCoins); j += 1 {
				rowIndexes[j] = new(big.Int).SetUint64(cmtIndices[currentIndex])
				row[j] = publicKeys[currentIndex]
				sumInputs.Add(sumInputs, commitments[currentIndex])
				if assetTags[currentIndex] == nil {
					log.Printf("CA error: missing asset tag for signing in DB coin - %v\n", currentIndex)
					err := utils.NewTransactionErr(utils.SignTxError, fmt.Errorf("cannot sign CA token : a CA coin in DB does not have asset tag"))
					return nil, nil, nil, err
				}
				sumInputAssetTags.Add(sumInputAssetTags, assetTags[currentIndex])
				currentIndex += 1
			}
		}
		sumInputAssetTags.ScalarMult(sumInputAssetTags, outCount)

		assetSum := new(crypto.Point).Sub(sumInputAssetTags, sumOutputAssetTags)
		row = append(row, assetSum)
		row = append(row, sumInputs)
		if i == pi {
			lastTwoColumnsCommitmentToZero = []*crypto.Point{assetSum, sumInputs}
		}

		ring[i] = row
		indexes[i] = rowIndexes
	}
	return mlsag.NewRing(ring), indexes, lastTwoColumnsCommitmentToZero, nil
}

func (tx *Tx) proveCA(params *tx_generic.TxPrivacyInitParams) (bool, error) {
	var err error
	var outputCoins []*coin.CoinV2
	var sharedSecrets []*crypto.Point
	// log.Printf("tokenID is %v\n",params.TokenID)
	var numOfCoinsBurned uint = 0
	var isBurning = false
	for _, inf := range params.PaymentInfo {
		c, ss, err := createUniqueOTACoinCA(inf, params.TokenID)
		if err != nil {
			log.Printf("Cannot parse outputCoinV2 to outputCoins, error %v\n", err)
			return false, err
		}
		// the only way err!=nil but ss==nil is a coin meant for burning address
		if ss == nil {
			isBurning = true
			numOfCoinsBurned += 1
		}
		sharedSecrets = append(sharedSecrets, ss)
		outputCoins = append(outputCoins, c)
	}
	// first, reject the invalid case. After this, isBurning will correctly determine if TX is burning
	if numOfCoinsBurned > 1 {
		log.Printf("Cannot burn multiple coins")
		return false, utils.NewTransactionErr(utils.UnexpectedError, fmt.Errorf("output must not have more than 1 burned coin"))
	}
	// outputCoins, err := newCoinV2ArrayFromPaymentInfoArray(params.PaymentInfo, params.TokenID, params.StateDB)

	// inputCoins is plainCoin because it may have coinV1 with coinV2
	inputCoins := params.InputCoins
	tx.Proof, err = privacy.ProveV2(inputCoins, outputCoins, sharedSecrets, true, params.PaymentInfo)
	if err != nil {
		log.Printf("Error in privacy_v2.Prove, error %v ", err)
		return false, err
	}

	err = tx.signCA(inputCoins, outputCoins, sharedSecrets, params, tx.Hash()[:])
	return isBurning, err
}

func (tx *Tx) signCA(inp []coin.PlainCoin, out []*coin.CoinV2, outputSharedSecrets []*crypto.Point, params *tx_generic.TxPrivacyInitParams, hashedMessage []byte) error {
	if tx.Sig != nil {
		return utils.NewTransactionErr(utils.UnexpectedError, fmt.Errorf("input transaction must be an unsigned one"))
	}
	ringSize := privacy.RingSize

	// Generate Ring
	piBig, piErr := common.RandBigIntMaxRange(big.NewInt(int64(ringSize)))
	if piErr != nil {
		return piErr
	}
	var pi = int(piBig.Int64())
	shardID := common.GetShardIDFromLastByte(tx.PubKeyLastByteSender)
	ring, indexes, commitmentsToZero, err := generateMlsagRingWithIndexesCA(inp, out, params, pi, shardID, ringSize)
	if err != nil {
		log.Printf("generateMLSAGRingWithIndexes got error %v ", err)
		return err
	}

	// Set SigPubKey
	txSigPubKey := new(SigPubKey)
	txSigPubKey.Indexes = indexes
	tx.SigPubKey, err = txSigPubKey.Bytes()
	if err != nil {
		log.Printf("tx.SigPubKey cannot parse from Bytes, error %v ", err)
		return err
	}

	privateKeysMlsag, err := createPrivateKeyMlsagCA(inp, out, outputSharedSecrets, params, shardID, commitmentsToZero)
	if err != nil {
		log.Printf("Cannot create private key of mlsag: %v", err)
		return err
	}
	sag := mlsag.NewMlsag(privateKeysMlsag, ring, pi)
	sk, err := privacy.ArrayScalarToBytes(&privateKeysMlsag)
	if err != nil {
		log.Printf("tx.SigPrivKey cannot parse arrayScalar to Bytes, error %v ", err)
		return err
	}
	tx.SetPrivateKey(sk)

	// Set Signature
	mlsagSignature, err := sag.SignConfidentialAsset(hashedMessage)
	if err != nil {
		log.Printf("Cannot signOnMessage mlsagSignature, error %v ", err)
		return err
	}
	// inputCoins already hold keyImage so set to nil to reduce size
	mlsagSignature.SetKeyImages(nil)
	tx.Sig, err = mlsagSignature.ToBytes()

	return err
}
