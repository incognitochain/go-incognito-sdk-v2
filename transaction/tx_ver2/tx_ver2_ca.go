package tx_ver2

import (
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v2/mlsag"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"math/big"

	// "github.com/incognitochain/incognito-chain/wallet"
)

//Create unique OTA coin without the help of the db
func createUniqueOTACoinCA(paymentInfo *key.PaymentInfo, tokenID *common.Hash) (*coin.CoinV2, *crypto.Point, error) {
	if tokenID == nil {
		tokenID = &common.PRVCoinID
	}
	c, sharedSecret, err := coin.NewCoinCA(paymentInfo, tokenID)
	if err != nil {
		fmt.Printf("Cannot parse coin based on payment info err: %v", err)
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

func createPrivKeyMlsagCA(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, outputSharedSecrets []*crypto.Point, params *tx_generic.TxPrivacyInitParams, shardID byte, commitmentsToZero []*crypto.Point) ([]*crypto.Scalar, error) {
	senderSK := params.SenderSK
	// db := params.StateDB
	tokenID := params.TokenID
	if tokenID == nil {
		tokenID = &common.PRVCoinID
	}
	rehashed := crypto.HashToPoint(tokenID[:])
	sumRand := new(crypto.Scalar).FromUint64(0)

	privKeyMlsag := make([]*crypto.Scalar, len(inputCoins)+2)
	sumInputAssetTagBlinders := new(crypto.Scalar).FromUint64(0)
	numOfInputs := new(crypto.Scalar).FromUint64(uint64(len(inputCoins)))
	numOfOutputs := new(crypto.Scalar).FromUint64(uint64(len(outputCoins)))
	mySkBytes := (*senderSK)[:]
	for i := 0; i < len(inputCoins); i += 1 {
		var err error
		privKeyMlsag[i], err = inputCoins[i].ParsePrivateKeyOfCoin(*senderSK)
		if err != nil {
			fmt.Printf("Cannot parse private key of coin %v", err)
			return nil, err
		}

		inputCoin_specific, ok := inputCoins[i].(*coin.CoinV2)
		if !ok || inputCoin_specific.GetAssetTag() == nil {
			return nil, errors.New("Cannot cast a coin as v2-CA")
		}

		isUnblinded := crypto.IsPointEqual(rehashed, inputCoin_specific.GetAssetTag())

		sharedSecret := new(crypto.Point).Identity()
		bl := new(crypto.Scalar).FromUint64(0)
		if !isUnblinded {
			sharedSecret, err = inputCoin_specific.RecomputeSharedSecret(mySkBytes)
			if err != nil {
				fmt.Printf("Cannot recompute shared secret : %v", err)
				return nil, err
			}

			bl, err = coin.ComputeAssetTagBlinder(sharedSecret)
			if err != nil {
				return nil, err
			}
		}

		v := inputCoin_specific.GetAmount()
		effectiveRCom := new(crypto.Scalar).Mul(bl, v)
		effectiveRCom.Add(effectiveRCom, inputCoin_specific.GetRandomness())

		sumInputAssetTagBlinders.Add(sumInputAssetTagBlinders, bl)
		sumRand.Add(sumRand, effectiveRCom)
	}
	sumInputAssetTagBlinders.Mul(sumInputAssetTagBlinders, numOfOutputs)

	sumOutputAssetTagBlinders := new(crypto.Scalar).FromUint64(0)

	var err error
	for i, oc := range outputCoins {
		if oc.GetAssetTag() == nil {
			return nil, errors.New("Cannot cast a coin as v2-CA")
		}
		// lengths between 0 and len(outputCoins) were rejected before
		bl := new(crypto.Scalar).FromUint64(0)
		isUnblinded := crypto.IsPointEqual(rehashed, oc.GetAssetTag())
		if !isUnblinded {
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
		fmt.Printf("Received %d points to check when signing MLSAG", len(commitmentsToZero))
		return nil, utils.NewTransactionErr(utils.UnexpectedError, errors.New("Error : need exactly 2 points for MLSAG double-checking"))
	}
	match1 := crypto.IsPointEqual(firstCommitmentToZeroRecomputed, commitmentsToZero[0])
	match2 := crypto.IsPointEqual(secondCommitmentToZeroRecomputed, commitmentsToZero[1])
	if !match1 || !match2 {
		return nil, utils.NewTransactionErr(utils.UnexpectedError, errors.New("Error : asset tag sum or commitment sum mismatch"))
	}

	privKeyMlsag[len(inputCoins)] = assetSum
	privKeyMlsag[len(inputCoins)+1] = sumRand
	return privKeyMlsag, nil
}

func generateMlsagRingWithIndexesCA(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, params *tx_generic.TxPrivacyInitParams, pi int, shardID byte, ringSize int) (*mlsag.Ring, [][]*big.Int, []*crypto.Point, error) {
	cmtIndices, myIndices, commitments, publicKeys, assetTags, err := ParseParamsForRing(params.KvArgs, len(inputCoins), ringSize)
	if err != nil {
		return nil, nil, nil, utils.NewTransactionErr(utils.UnexpectedError, fmt.Errorf("ParseParamsForRing error: %v", err))
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
			fmt.Printf("CA error: missing asset tag for signing in output coin - %v", oc.Bytes())
			err := utils.NewTransactionErr(utils.SignTxError, errors.New("Cannot sign CA token : an output coin does not have asset tag"))
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
				inputCoin_specific, ok := inputCoins[j].(*coin.CoinV2)
				if !ok {
					return nil, nil, nil, errors.New("Cannot cast a coin as v2")
				}
				if inputCoin_specific.GetAssetTag() == nil {
					fmt.Printf("CA error: missing asset tag for signing in input coin - %v", inputCoin_specific.Bytes())
					err := utils.NewTransactionErr(utils.SignTxError, errors.New("Cannot sign CA token : an input coin does not have asset tag"))
					return nil, nil, nil, err
				}
				sumInputAssetTags.Add(sumInputAssetTags, inputCoin_specific.GetAssetTag())
			}
		} else {
			for j := 0; j < len(inputCoins); j += 1 {
				rowIndexes[j] = new(big.Int).SetUint64(cmtIndices[currentIndex])
				row[j] = publicKeys[currentIndex]
				sumInputs.Add(sumInputs, commitments[currentIndex])
				if assetTags[currentIndex] == nil {
					fmt.Printf("CA error: missing asset tag for signing in DB coin - %v", currentIndex)
					err := utils.NewTransactionErr(utils.SignTxError, errors.New("Cannot sign CA token : a CA coin in DB does not have asset tag"))
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
	// fmt.Printf("tokenID is %v\n",params.TokenID)
	var numOfCoinsBurned uint = 0
	var isBurning bool = false
	for _, inf := range params.PaymentInfo {
		c, ss, err := createUniqueOTACoinCA(inf, params.TokenID)
		if err != nil {
			fmt.Printf("Cannot parse outputCoinV2 to outputCoins, error %v ", err)
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
		fmt.Printf("Cannot burn multiple coins")
		return false, utils.NewTransactionErr(utils.UnexpectedError, errors.New("output must not have more than 1 burned coin"))
	}
	// outputCoins, err := newCoinV2ArrayFromPaymentInfoArray(params.PaymentInfo, params.TokenID, params.StateDB)

	// inputCoins is plainCoin because it may have coinV1 with coinV2
	inputCoins := params.InputCoins
	tx.Proof, err = privacy.ProveV2(inputCoins, outputCoins, sharedSecrets, true, params.PaymentInfo)
	if err != nil {
		fmt.Printf("Error in privacy_v2.Prove, error %v ", err)
		return false, err
	}

	err = tx.signCA(inputCoins, outputCoins, sharedSecrets, params, tx.Hash()[:])
	return isBurning, err
}

func (tx *Tx) signCA(inp []coin.PlainCoin, out []*coin.CoinV2, outputSharedSecrets []*crypto.Point, params *tx_generic.TxPrivacyInitParams, hashedMessage []byte) error {
	if tx.Sig != nil {
		return utils.NewTransactionErr(utils.UnexpectedError, errors.New("input transaction must be an unsigned one"))
	}
	ringSize := privacy.RingSize

	// Generate Ring
	piBig, piErr := common.RandBigIntMaxRange(big.NewInt(int64(ringSize)))
	if piErr != nil {
		return piErr
	}
	var pi int = int(piBig.Int64())
	shardID := common.GetShardIDFromLastByte(tx.PubKeyLastByteSender)
	ring, indexes, commitmentsToZero, err := generateMlsagRingWithIndexesCA(inp, out, params, pi, shardID, ringSize)
	if err != nil {
		fmt.Printf("generateMLSAGRingWithIndexes got error %v ", err)
		return err
	}

	// Set SigPubKey
	txSigPubKey := new(SigPubKey)
	txSigPubKey.Indexes = indexes
	tx.SigPubKey, err = txSigPubKey.Bytes()
	if err != nil {
		fmt.Printf("tx.SigPubKey cannot parse from Bytes, error %v ", err)
		return err
	}

	// Set sigPrivKey
	privKeysMlsag, err := createPrivKeyMlsagCA(inp, out, outputSharedSecrets, params, shardID, commitmentsToZero)
	if err != nil {
		fmt.Printf("Cannot create private key of mlsag: %v", err)
		return err
	}
	sag := mlsag.NewMlsag(privKeysMlsag, ring, pi)
	sk, err := privacy.ArrayScalarToBytes(&privKeysMlsag)
	if err != nil {
		fmt.Printf("tx.SigPrivKey cannot parse arrayScalar to Bytes, error %v ", err)
		return err
	}
	tx.SetPrivateKey(sk)

	// Set Signature
	mlsagSignature, err := sag.SignConfidentialAsset(hashedMessage)
	if err != nil {
		fmt.Printf("Cannot signOnMessage mlsagSignature, error %v ", err)
		return err
	}
	// inputCoins already hold keyImage so set to nil to reduce size
	mlsagSignature.SetKeyImages(nil)
	tx.Sig, err = mlsagSignature.ToBytes()

	return err
}
