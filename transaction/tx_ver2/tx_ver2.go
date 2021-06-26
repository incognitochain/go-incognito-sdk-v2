package tx_ver2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v2/mlsag"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"math"
	"math/big"
	"sort"
	"strconv"
	"time"
)

// TxSigPubKey of ver2 is array of Indexes in database
type SigPubKey struct {
	Indexes [][]*big.Int
}

type Tx struct {
	tx_generic.TxBase
}

func (sigPub SigPubKey) Bytes() ([]byte, error) {
	n := len(sigPub.Indexes)
	if n == 0 {
		return nil, errors.New("TxSigPublicKeyVer2.ToBytes: Indexes is empty")
	}
	if n > utils.MaxSizeByte {
		return nil, errors.New("TxSigPublicKeyVer2.ToBytes: Indexes is too large, too many rows")
	}
	m := len(sigPub.Indexes[0])
	if m > utils.MaxSizeByte {
		return nil, errors.New("TxSigPublicKeyVer2.ToBytes: Indexes is too large, too many columns")
	}
	for i := 1; i < n; i += 1 {
		if len(sigPub.Indexes[i]) != m {
			return nil, errors.New("TxSigPublicKeyVer2.ToBytes: Indexes is not a rectangle array")
		}
	}

	b := make([]byte, 0)
	b = append(b, byte(n))
	b = append(b, byte(m))
	for i := 0; i < n; i += 1 {
		for j := 0; j < m; j += 1 {
			currentByte := sigPub.Indexes[i][j].Bytes()
			lengthByte := len(currentByte)
			if lengthByte > utils.MaxSizeByte {
				return nil, errors.New("TxSigPublicKeyVer2.ToBytes: IndexesByte is too large")
			}
			b = append(b, byte(lengthByte))
			b = append(b, currentByte...)
		}
	}
	return b, nil
}

func (sigPub *SigPubKey) SetBytes(b []byte) error {
	if len(b) < 2 {
		return errors.New("txSigPubKeyFromBytes: cannot parse length of Indexes, length of input byte is too small")
	}
	n := int(b[0])
	m := int(b[1])
	offset := 2
	indexes := make([][]*big.Int, n)
	for i := 0; i < n; i += 1 {
		row := make([]*big.Int, m)
		for j := 0; j < m; j += 1 {
			if offset >= len(b) {
				return errors.New("txSigPubKeyFromBytes: cannot parse byte length of index[i][j], length of input byte is too small")
			}
			byteLength := int(b[offset])
			offset += 1
			if offset+byteLength > len(b) {
				return errors.New("txSigPubKeyFromBytes: cannot parse big int index[i][j], length of input byte is too small")
			}
			currentByte := b[offset : offset+byteLength]
			offset += byteLength
			row[j] = new(big.Int).SetBytes(currentByte)
		}
		indexes[i] = row
	}
	if sigPub == nil {
		sigPub = new(SigPubKey)
	}
	sigPub.Indexes = indexes
	return nil
}

// ========== GET FUNCTION ===========

func (tx *Tx) GetReceiverData() ([]coin.Coin, error) {
	if tx.Proof != nil && len(tx.Proof.GetOutputCoins()) > 0 {
		return tx.Proof.GetOutputCoins(), nil
	}
	return nil, nil
}

// ========== NORMAL INIT FUNCTIONS ==========

func createPrivKeyMlsag(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, senderSK *key.PrivateKey, commitmentToZero *crypto.Point) ([]*crypto.Scalar, error) {
	sumRand := new(crypto.Scalar).FromUint64(0)
	for _, in := range inputCoins {
		sumRand.Add(sumRand, in.GetRandomness())
	}
	for _, out := range outputCoins {
		sumRand.Sub(sumRand, out.GetRandomness())
	}

	privKeyMlsag := make([]*crypto.Scalar, len(inputCoins)+1)
	for i := 0; i < len(inputCoins); i += 1 {
		var err error
		privKeyMlsag[i], err = inputCoins[i].ParsePrivateKeyOfCoin(*senderSK)
		if err != nil {
			return nil, err
		}
	}
	commitmentToZeroRecomputed := new(crypto.Point).ScalarMult(privacy.PedCom.G[privacy.PedersenRandomnessIndex], sumRand)
	match := crypto.IsPointEqual(commitmentToZeroRecomputed, commitmentToZero)
	if !match {
		return nil, utils.NewTransactionErr(utils.SignTxError, errors.New("Error : asset tag sum or commitment sum mismatch"))
	}
	privKeyMlsag[len(inputCoins)] = sumRand
	return privKeyMlsag, nil
}

func (tx *Tx) Init(paramsInterface interface{}) error {
	params, ok := paramsInterface.(*tx_generic.TxPrivacyInitParams)
	if !ok {
		return errors.New("params of tx Init is not TxPrivacyInitParam")
	}

	jsb, _ := json.Marshal(params)
	if err := tx_generic.ValidateTxParams(params); err != nil {
		return err
	}

	// Init tx and params (tx and params will be changed)
	if err := tx.InitializeTxAndParams(params); err != nil {
		return err
	}

	// Check if this tx is nonPrivacyNonInput
	// Case 1: tx ptoken transfer with ptoken fee
	// Case 2: tx Reward
	// If it is non privacy non input then return
	if check, err := tx.IsNonPrivacyNonInput(params); check {
		return err
	}
	if err := tx.prove(params); err != nil {
		return err
	}
	jsb, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("marshal tx error: %v", err)
	}
	tx1 := new(Tx)
	err = json.Unmarshal(jsb, &tx1)
	if err != nil {
		return err
	}
	if tx1.Hash().String() != tx.Hash().String() {
		jsb, err := json.Marshal(tx1)
		if err != nil {
			return fmt.Errorf("marshal tx error: %v", err)
		}
		fmt.Println(string(jsb))
		return fmt.Errorf("txHash changes after unmarshalling: %v, %v", tx.Hash().String(), tx1.Hash().String())
	}

	txSize := tx.GetTxActualSize()
	if txSize > common.MaxTxSize {
		return utils.NewTransactionErr(utils.ExceedSizeTx, nil, strconv.Itoa(int(txSize)))
	}

	return nil
}

func (tx *Tx) signOnMessage(inp []coin.PlainCoin, out []*coin.CoinV2, params *tx_generic.TxPrivacyInitParams, hashedMessage []byte) error {
	if tx.Sig != nil {
		return utils.NewTransactionErr(utils.UnexpectedError, errors.New("input transaction must be an unsigned one"))
	}
	ringSize := privacy.RingSize

	// Generate Ring
	piBig, piErr := common.RandBigIntMaxRange(big.NewInt(int64(ringSize)))
	if piErr != nil {
		return piErr
	}
	var pi = int(piBig.Int64())
	ring, indexes, commitmentToZero, err := generateMLSAGRingWithIndexes(inp, out, params, pi, ringSize)
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
	privKeysMlsag, err := createPrivKeyMlsag(inp, out, params.SenderSK, commitmentToZero)
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
	mlsagSignature, err := sag.Sign(hashedMessage)
	if err != nil {
		fmt.Printf("Cannot signOnMessage mlsagSignature, error %v ", err)
		return err
	}
	// inputCoins already hold keyImage so set to nil to reduce size
	mlsagSignature.SetKeyImages(nil)
	tx.Sig, err = mlsagSignature.ToBytes()

	return err
}

func (tx *Tx) prove(params *tx_generic.TxPrivacyInitParams) error {
	var err error
	outputCoins := make([]*coin.CoinV2, 0)
	for _, paymentInfo := range params.PaymentInfo {
		outputCoin, err := coin.NewCoinFromPaymentInfo(paymentInfo) //We do not mind duplicated OTAs, server will handle them.
		if err != nil {
			return err
		}

		outputCoins = append(outputCoins, outputCoin)
	}

	// inputCoins is plainCoin because it may have coinV1 with coinV2
	inputCoins := params.InputCoins

	tx.Proof, err = privacy.ProveV2(inputCoins, outputCoins, nil, false, params.PaymentInfo)
	if err != nil {
		return err
	}

	if tx.GetMetadata() != nil {
		if err := tx.GetMetadata().Sign(params.SenderSK, tx); err != nil {
			return err
		}
	}

	err = tx.signOnMessage(inputCoins, outputCoins, params, tx.Hash()[:])
	return err
}

// ========== NORMAL VERIFY FUNCTIONS ==========

//Parse params and check their validity for generating a MLSAG ring.
func ParseParamsForRing(kvArgs map[string]interface{}, lenInput, ringSize int) (cmtIndices []uint64, myIndices []uint64, commitments []*crypto.Point, publicKeys []*crypto.Point, assetTags []*crypto.Point, err error) {
	if kvArgs == nil {
		fmt.Println("kvargs is nil: need more params to proceed")
		return nil, nil, nil, nil, nil, errors.New("kvargs is nil: need more params to proceed")
	}

	//Get list of decoy indices.
	tmp, ok := kvArgs[utils.CommitmentIndices]
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("decoy commitment indices not found: %v", kvArgs)
	}

	cmtIndices, ok = tmp.([]uint64)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot parse commitment indices: %v", tmp)
	}
	if len(cmtIndices) < lenInput*(ringSize-1) {
		return nil, nil, nil, nil, nil, fmt.Errorf("not enough decoy commitment indices: have %v, need at least %v (%v input coins)", len(cmtIndices), lenInput*(ringSize-1), lenInput)
	}

	//Get list of decoy commitments.
	tmp, ok = kvArgs[utils.Commitments]
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("decoy commitment list not found: %v", kvArgs)
	}

	commitments, ok = tmp.([]*crypto.Point)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot parse decoy commitment indices: %v", tmp)
	}
	if len(commitments) < lenInput*(ringSize-1) {
		return nil, nil, nil, nil, nil, fmt.Errorf("not enough decoy commitments: have %v, need at least %v (%v input coins)", len(commitments), lenInput*(ringSize-1), lenInput)
	}

	//Get list of decoy public keys
	tmp, ok = kvArgs[utils.PublicKeys]
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("decoy public keys not found: %v", kvArgs)
	}

	publicKeys, ok = tmp.([]*crypto.Point)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot parse decoy public keys: %v", tmp)
	}
	if len(publicKeys) < lenInput*(ringSize-1) {
		return nil, nil, nil, nil, nil, fmt.Errorf("not enough decoy public keys: have %v, need at least %v (%v input coins)", len(publicKeys), lenInput*(ringSize-1), lenInput)
	}

	//Get list of decoy asset tags
	tmp, ok = kvArgs[utils.AssetTags]
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("decoy asset tags not found: %v", kvArgs)
	}

	assetTags, ok = tmp.([]*crypto.Point)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot parse decoy asset tags: %v", tmp)
	}
	//if len(assetTags) < lenInput*(ringSize-1) {
	//	return nil, nil, nil, nil, nil, fmt.Errorf("not enough decoy asset tags: have %v, need at least %v (%v input coins).", len(assetTags), lenInput*(ringSize-1), lenInput))
	//}

	//Get list of inputcoin indices
	tmp, ok = kvArgs[utils.MyIndices]
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("inputCoin commitment indices not found: %v", kvArgs)
	}

	myIndices, ok = tmp.([]uint64)
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot parse inputCoin commitment indices: %v", tmp)
	}
	if len(myIndices) != lenInput {
		return nil, nil, nil, nil, nil, fmt.Errorf("not enough indices for input coins: have %v, want %v.", len(myIndices), lenInput)
	}

	return
}

//Generate an MLSAG ring with input decoy commitments, public keys, and indices (params.Kvargs).
func generateMLSAGRingWithIndexes(inputCoins []coin.PlainCoin, outputCoins []*coin.CoinV2, params *tx_generic.TxPrivacyInitParams, pi int, ringSize int) (*mlsag.Ring, [][]*big.Int, *crypto.Point, error) {
	lenInput := len(inputCoins)
	kvArgs := params.Kvargs

	//Retrieve decoys' info from kvArgs
	cmtIndices, myIndices, commitments, publicKeys, _, err := ParseParamsForRing(kvArgs, lenInput, ringSize)
	if err != nil {
		return nil, nil, nil, err
	}

	outputCoinsAsGeneric := make([]coin.Coin, len(outputCoins))
	for i := 0; i < len(outputCoins); i++ {
		outputCoinsAsGeneric[i] = outputCoins[i]
	}
	sumOutputsWithFee := tx_generic.CalculateSumOutputsWithFee(outputCoinsAsGeneric, params.Fee)
	indices := make([][]*big.Int, ringSize)
	ring := make([][]*crypto.Point, ringSize)
	var commitmentToZero *crypto.Point

	currentIndex := 0
	for i := 0; i < ringSize; i += 1 {
		sumInputs := new(crypto.Point).Identity()
		sumInputs.Sub(sumInputs, sumOutputsWithFee)

		row := make([]*crypto.Point, len(inputCoins))
		rowIndexes := make([]*big.Int, len(inputCoins))
		if i == pi {
			for j := 0; j < len(inputCoins); j += 1 {
				row[j] = inputCoins[j].GetPublicKey()
				rowIndexes[j] = new(big.Int).SetUint64(myIndices[j])
				sumInputs.Add(sumInputs, inputCoins[j].GetCommitment())
			}
		} else {
			for j := 0; j < len(inputCoins); j += 1 {
				rowIndexes[j] = new(big.Int).SetUint64(cmtIndices[currentIndex])
				row[j] = publicKeys[currentIndex]
				sumInputs.Add(sumInputs, commitments[currentIndex])

				currentIndex += 1
			}
		}
		row = append(row, sumInputs)
		if i == pi {
			commitmentToZero = sumInputs
		}
		ring[i] = row
		indices[i] = rowIndexes
	}
	return mlsag.NewRing(ring), indices, commitmentToZero, nil
}

// ========== SALARY FUNCTIONS: INIT AND VALIDATE  ==========

func (tx *Tx) InitTxSalary(otaCoin *coin.CoinV2, privateKey *key.PrivateKey, metaData metadata.Metadata) error {
	tokenID := &common.Hash{}
	if err := tokenID.SetBytes(common.PRVCoinID[:]); err != nil {
		return utils.NewTransactionErr(utils.TokenIDInvalidError, err, tokenID.String())
	}

	tx.Version = utils.TxVersion2Number
	tx.Type = common.TxRewardType
	if tx.LockTime == 0 {
		tx.LockTime = time.Now().Unix()
	}

	tempOutputCoin := []coin.Coin{otaCoin}
	proof := new(privacy.ProofV2)
	proof.Init()
	proof.SetOutputCoins(tempOutputCoin)
	tx.Proof = proof

	publicKeyBytes := otaCoin.GetPublicKey().ToBytesS()
	tx.PubKeyLastByteSender = common.GetShardIDFromLastByte(publicKeyBytes[len(publicKeyBytes)-1])

	// signOnMessage Tx using ver1 schnorr
	tx.SetPrivateKey(*privateKey)
	tx.SetMetadata(metaData)

	var err error
	if tx.Sig, tx.SigPubKey, err = tx_generic.SignNoPrivacy(privateKey, tx.Hash()[:]); err != nil {
		return utils.NewTransactionErr(utils.SignTxError, err)
	}
	return nil
}

func (tx Tx) Hash() *common.Hash {
	// leave out signature & its public key when hashing tx
	tx.Sig = []byte{}
	tx.SigPubKey = []byte{}
	inBytes, err := json.Marshal(tx)
	if err != nil {
		return nil
	}
	hash := common.HashH(inBytes)
	// after this returns, tx is restored since the receiver is not a pointer
	return &hash
}

// HashWithoutMetadataSig returns the hash of this transaction, but it leaves out the metadata's own signature field. It is used to verify that metadata signature.
func (tx Tx) HashWithoutMetadataSig() *common.Hash {
	md := tx.GetMetadata()
	mdHash := md.HashWithoutSig()
	tx.SetMetadata(nil)
	txHash := tx.Hash()
	if mdHash == nil || txHash == nil {
		return nil
	}
	// tx.SetMetadata(md)
	inBytes := append(mdHash[:], txHash[:]...)
	hash := common.HashH(inBytes)
	return &hash
}

// ========== SHARED FUNCTIONS ============

func (tx Tx) GetTxMintData() (bool, coin.Coin, *common.Hash, error) {
	return tx_generic.GetTxMintData(&tx, &common.PRVCoinID)
}

func (tx Tx) GetTxBurnData() (bool, coin.Coin, *common.Hash, error) {
	return tx_generic.GetTxBurnData(&tx)
}

func (tx Tx) GetTxFullBurnData() (bool, coin.Coin, coin.Coin, *common.Hash, error) {
	isBurn, burnedCoin, burnedToken, err := tx.GetTxBurnData()
	return isBurn, burnedCoin, nil, burnedToken, err
}

func (tx Tx) GetTxActualSize() uint64 {
	jsb, err := json.Marshal(tx)
	if err != nil {
		return 0
	}
	return uint64(math.Ceil(float64(len(jsb)) / 1024))
}

func (tx Tx) ListOTAHashH() []common.Hash {
	result := make([]common.Hash, 0)
	if tx.Proof != nil {
		for _, outputCoin := range tx.Proof.GetOutputCoins() {
			//Discard coins sent to the burning address
			if wallet.IsPublicKeyBurningAddress(outputCoin.GetPublicKey().ToBytesS()) {
				continue
			}
			hash := common.HashH(outputCoin.GetPublicKey().ToBytesS())
			result = append(result, hash)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})
	return result
}
