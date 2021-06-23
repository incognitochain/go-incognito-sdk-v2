package tx_ver2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"math"
	"sort"
	"strconv"
)

type TxTokenDataVersion2 struct {
	PropertyID     common.Hash
	PropertyName   string
	PropertySymbol string
	SigPubKey      []byte `json:"SigPubKey,omitempty"` // 33 bytes
	Sig            []byte `json:"Sig,omitempty"`       //
	Proof          privacy.Proof
	Type           int
	Mintable       bool
}

func (td TxTokenDataVersion2) Hash() (*common.Hash, error) {
	// leave out signature & its public key when hashing tx
	td.Sig = []byte{}
	td.SigPubKey = []byte{}
	inBytes, err := json.Marshal(td)

	if err != nil {
		return nil, err
	}
	// after this returns, tx is restored since the receiver is not a pointer
	hash := common.HashH(inBytes)
	return &hash, nil
}

func makeTxToken(txPRV *Tx, pubkey, sig []byte, proof privacy.Proof) *Tx {
	result := &Tx{
		TxBase: tx_generic.TxBase{
			Version:              txPRV.Version,
			Type:                 txPRV.Type,
			LockTime:             txPRV.LockTime,
			Fee:                  0,
			PubKeyLastByteSender: common.GetShardIDFromLastByte(txPRV.PubKeyLastByteSender),
			Metadata:             nil,
		},
	}
	var clonedInfo []byte = nil
	var err error
	if txPRV.Info != nil {
		clonedInfo = make([]byte, len(txPRV.Info))
		copy(clonedInfo, txPRV.Info)
	}
	var clonedProof privacy.Proof = nil
	// feed the type to parse proof 
	proofType := txPRV.Type
	if proofType == common.TxTokenConversionType {
		proofType = common.TxConversionType
	}
	if proof != nil {
		clonedProof, err = utils.ParseProof(proof, txPRV.Version, proofType)
		if err != nil {
			jsb, _ := json.Marshal(proof)
			fmt.Printf("Cannot parse proof %s using version %v - type %v", string(jsb), txPRV.Version, txPRV.Type)
			return nil
		}
	}
	var clonedSig []byte = nil
	if sig != nil {
		clonedSig = make([]byte, len(sig))
		copy(clonedSig, sig)
	}
	var clonedPk []byte = nil
	if pubkey != nil {
		clonedPk = make([]byte, len(pubkey))
		copy(clonedPk, pubkey)
	}
	result.Info = clonedInfo
	result.Proof = clonedProof
	result.Sig = clonedSig
	result.SigPubKey = clonedPk
	result.Info = clonedInfo

	return result
}

type TxToken struct {
	Tx             Tx                  `json:"Tx"`
	TokenData      TxTokenDataVersion2 `json:"TxTokenPrivacyData"`
	cachedTxNormal *Tx
}

func (tx *TxToken) Hash() *common.Hash {
	firstHash := tx.Tx.Hash()
	secondHash, err := tx.TokenData.Hash()
	if err != nil {
		return nil
	}
	result := common.HashH(append(firstHash[:], secondHash[:]...))
	return &result
}
func (txToken TxToken) HashWithoutMetadataSig() *common.Hash { return txToken.Tx.HashWithoutMetadataSig() }
func (td TxTokenDataVersion2) ToCompatTokenData(ttx metadata.Transaction) tx_generic.TxTokenData {
	return tx_generic.TxTokenData{
		TxNormal:       ttx,
		PropertyID:     td.PropertyID,
		PropertyName:   td.PropertyName,
		PropertySymbol: td.PropertySymbol,
		Type:           td.Type,
		Mintable:       td.Mintable,
		Amount:         0,
	}
}
func decomposeTokenData(td tx_generic.TxTokenData) (*TxTokenDataVersion2, *Tx, error) {
	result := TxTokenDataVersion2{
		PropertyID:     td.PropertyID,
		PropertyName:   td.PropertyName,
		PropertySymbol: td.PropertySymbol,
		Type:           td.Type,
		Mintable:       td.Mintable,
	}
	tx, ok := td.TxNormal.(*Tx)
	if !ok {
		return nil, nil, errors.New("Error while casting a transaction to v2")
	}
	return &result, tx, nil
}
func (tx *TxToken) GetTxBase() metadata.Transaction {
	return &tx.Tx
}
func (tx *TxToken) SetTxBase(inTx metadata.Transaction) error {
	temp, ok := inTx.(*Tx)
	if !ok {
		return errors.New("Cannot set TxBase : wrong type")
	}
	tx.Tx = *temp
	return nil
}
func (tx *TxToken) GetTxNormal() metadata.Transaction {
	if tx.cachedTxNormal != nil {
		return tx.cachedTxNormal
	}
	result := makeTxToken(&tx.Tx, tx.TokenData.SigPubKey, tx.TokenData.Sig, tx.TokenData.Proof)
	// tx.cachedTxNormal = result
	return result
}
func (tx *TxToken) SetTxNormal(inTx metadata.Transaction) error {
	temp, ok := inTx.(*Tx)
	if !ok {
		return utils.NewTransactionErr(utils.UnexpectedError, errors.New("Cannot set TxNormal : wrong type"))
	}
	tx.TokenData.SigPubKey = temp.SigPubKey
	tx.TokenData.Sig = temp.Sig
	tx.TokenData.Proof = temp.Proof
	tx.cachedTxNormal = temp
	return nil
}

// this signs only on the hash of the data in it
func (tx *Tx) proveToken(params *tx_generic.TxPrivacyInitParams) (bool, error) {
	if err := tx_generic.ValidateTxParams(params); err != nil {
		return false, err
	}

	// Init tx and params (tx and params will be changed)
	if err := tx.InitializeTxAndParams(params); err != nil {
		return false, err
	}
	tx.SetType(common.TxCustomTokenPrivacyType)
	isBurning, err := tx.proveCA(params)
	if err != nil {
		return false, err
	}
	return isBurning, nil
}

func (txToken *TxToken) initToken(txNormal *Tx, params *tx_generic.TxTokenParams) error {
	txToken.TokenData.Type = params.TokenParams.TokenTxType
	txToken.TokenData.PropertyName = params.TokenParams.PropertyName
	txToken.TokenData.PropertySymbol = params.TokenParams.PropertySymbol
	txToken.TokenData.Mintable = params.TokenParams.Mintable

	switch params.TokenParams.TokenTxType {
	case utils.CustomTokenInit:
		return fmt.Errorf("wrong method for initializing a token, use metadata instead")
	case utils.CustomTokenTransfer:
		{
			propertyID, err := new(common.Hash).NewHashFromStr(params.TokenParams.PropertyID)
			if err != nil {
				return utils.NewTransactionErr(utils.TokenIDInvalidError, err)
			}
			dbFacingTokenID := common.ConfidentialAssetID

			// fee in pToken is not supported
			feeToken := uint64(0)
			txParams := tx_generic.NewTxPrivacyInitParams(
				params.SenderKey,
				params.TokenParams.Receiver,
				params.TokenParams.TokenInput,
				feeToken,
				params.HasPrivacyToken,
				propertyID,
				nil,
				nil,
				params.TokenParams.KvArgs)
			isBurning, err := txNormal.proveToken(txParams)
			if err != nil {
				return utils.NewTransactionErr(utils.PrivacyTokenInitTokenDataError, err)
			}

			if isBurning {
				// show plain tokenID if this is a burning TX
				txToken.TokenData.PropertyID = *propertyID
			} else {
				// tokenID is already hidden in asset tags in coin, here we use the umbrella ID
				txToken.TokenData.PropertyID = dbFacingTokenID
			}

			err = txToken.SetTxNormal(txNormal)
			if err != nil {
				return utils.NewTransactionErr(utils.UnexpectedError, err)
			}
			return nil
		}
	default:
		return utils.NewTransactionErr(utils.PrivacyTokenTxTypeNotHandleError, errors.New("can't handle this TokenTxType"))
	}
}

// this signs on the hash of both sub TXs
func (tx *Tx) provePRV(params *tx_generic.TxPrivacyInitParams) ([]coin.PlainCoin, []*coin.CoinV2, error) {
	var err error
	outputCoins := make([]*coin.CoinV2, 0)
	for _, paymentInfo := range params.PaymentInfo {
		outputCoin, err := coin.NewCoinFromPaymentInfo(paymentInfo) //We do not mind duplicated OTAs, server will handle them.
		if err != nil {
			fmt.Printf("Cannot parse outputCoinV2 to outputCoins, error %v\n", err)
			return nil, nil, err
		}

		outputCoins = append(outputCoins, outputCoin)
	}

	// inputCoins is plainCoin because it may have coinV1 with coinV2
	inputCoins := params.InputCoins

	tx.Proof, err = privacy.ProveV2(inputCoins, outputCoins, nil, false, params.PaymentInfo)
	if err != nil {
		fmt.Printf("Error in privacy_v2.Prove, error %v ", err)
		return nil, nil, err
	}

	if tx.GetMetadata() != nil {
		if err := tx.GetMetadata().Sign(params.SenderSK, tx); err != nil {
			fmt.Printf("Cannot signOnMessage txMetadata in shouldSignMetadata")
			return nil, nil, err
		}
	}

	// Get Hash of the whole txToken then sign on it
	// message := common.HashH(append(tx.Hash()[:], hashedTokenMessage...))

	return inputCoins, outputCoins, nil
}

func (txToken *TxToken) initPRV(feeTx *Tx, params *tx_generic.TxPrivacyInitParams) ([]coin.PlainCoin, []*coin.CoinV2, error) {
	// txTokenDataHash, err := txToken.TokenData.Hash()
	// if err != nil {
	// 	fmt.Printf("Cannot calculate txPrivacyTokenData Hash, err %v", err)
	// 	return nil, nil, err
	// }
	feeTx.SetType(common.TxCustomTokenPrivacyType)
	inps, outs, err := feeTx.provePRV(params)
	if err != nil {
		return nil, nil, utils.NewTransactionErr(utils.PrivacyTokenInitPRVError, err)
	}
	// override TxCustomTokenPrivacyType type
	// txToken.SetTxBase(feeTx)

	return inps, outs, nil
}

func (txToken *TxToken) Init(paramsInterface interface{}) error {
	params, ok := paramsInterface.(*tx_generic.TxTokenParams)
	if !ok {
		return errors.New("Cannot init TxCustomTokenPrivacy because params is not correct")
	}

	if params.TokenParams.Fee > 0 || params.FeeNativeCoin == 0 {
		fmt.Printf("only accept tx fee in PRV")
		return utils.NewTransactionErr(utils.PrivacyTokenInitFeeParamsError, nil, strconv.Itoa(int(params.TokenParams.Fee)))
	}

	txPrivacyParams := tx_generic.NewTxPrivacyInitParams(
		params.SenderKey,
		params.PaymentInfo,
		params.InputCoin,
		params.FeeNativeCoin,
		params.HasPrivacyCoin,
		nil,
		params.MetaData,
		params.Info,
		params.KvArgs)
	if err := tx_generic.ValidateTxParams(txPrivacyParams); err != nil {
		return err
	}
	// Init tx and params (tx and params will be changed)
	tx := new(Tx)
	if err := tx.InitializeTxAndParams(txPrivacyParams); err != nil {
		return err
	}

	// Check if this tx is nonPrivacyNonInput
	// Case 1: tx ptoken transfer with ptoken fee
	// Case 2: tx Reward
	// If it is non privacy non input then return
	if check, err := tx.IsNonPrivacyNonInput(txPrivacyParams); check {
		return err
	}

	// Init PRV Fee
	inps, outs, err := txToken.initPRV(tx, txPrivacyParams)
	if err != nil {
		fmt.Printf("Cannot init PRV fee for tokenver2: err %v", err)
		return err
	}

	txn := makeTxToken(tx, nil, nil, nil)
	// Init, prove and sign(CA) Token
	if err := txToken.initToken(txn, params); err != nil {
		fmt.Printf("Cannot init token ver2: err %v", err)
		return err
	}
	tdh, err := txToken.TokenData.Hash()
	if err != nil {
		return err
	}
	message := common.HashH(append(tx.Hash()[:], tdh[:]...))
	err = tx.signOnMessage(inps, outs, txPrivacyParams, message[:])
	if err != nil {
		return err
	}

	err = txToken.SetTxBase(tx)
	if err != nil {
		return err
	}
	// check tx size
	txSize := txToken.GetTxActualSize()
	if txSize > common.MaxTxSize {
		return utils.NewTransactionErr(utils.ExceedSizeTx, nil, strconv.Itoa(int(txSize)))
	}
	return nil
}

func (txToken *TxToken) InitTxTokenSalary(otaCoin *coin.CoinV2, privKey *key.PrivateKey, metaData metadata.Metadata, coinID *common.Hash, coinName string) error {
	var err error
	// Check validate params
	txPrivacyParams := tx_generic.NewTxPrivacyInitParams(
		privKey, []*key.PaymentInfo{}, nil, 0, false, nil, metaData, nil, nil)
	if err := tx_generic.ValidateTxParams(txPrivacyParams); err != nil {
		return err
	}

	// Create TxToken
	var propertyID [common.HashSize]byte
	copy(propertyID[:], coinID[:])
	txToken.TokenData.PropertyID = propertyID
	txToken.TokenData.Type = utils.CustomTokenInit
	txToken.TokenData.PropertyName = coinName
	txToken.TokenData.PropertySymbol = coinName
	txToken.TokenData.Mintable = true

	tempOutputCoin := []coin.Coin{otaCoin}
	proof := new(privacy.ProofV2)
	proof.Init()
	if err = proof.SetOutputCoins(tempOutputCoin); err != nil {
		fmt.Printf("Init customPrivacyToken cannot set outputCoins")
		return err
	}

	// Init tx fee params
	tx := new(Tx)
	if err := tx.InitializeTxAndParams(txPrivacyParams); err != nil {
		return err
	}
	tx.SetType(common.TxCustomTokenPrivacyType)
	tx.SetPrivateKey(*txPrivacyParams.SenderSK)
	temp := makeTxToken(tx, []byte{}, []byte{}, proof)
	err = txToken.SetTxNormal(temp)
	if err != nil {
		return utils.NewTransactionErr(utils.UnexpectedError, err)
	}

	hashedTokenMessage, err := txToken.TokenData.Hash()
	if err != nil {
		return utils.NewTransactionErr(utils.SignTxError, err)
	}

	message := common.HashH(append(tx.Hash()[:], hashedTokenMessage[:]...))
	if tx.Sig, tx.SigPubKey, err = tx_generic.SignNoPrivacy(privKey, message[:]); err != nil {
		fmt.Println(fmt.Errorf("Cannot signOnMessage tx %v\n", err))
		return utils.NewTransactionErr(utils.SignTxError, err)
	}

	err = txToken.SetTxBase(tx)
	if err != nil {
		return utils.NewTransactionErr(utils.UnexpectedError, err)
	}
	return nil
}

func (txToken TxToken) GetTxActualSize() uint64 {
	jsb, err := json.Marshal(txToken)
	if err != nil {
		return 0
	}
	return uint64(math.Ceil(float64(len(jsb)) / 1024))
}

func (tx TxToken) GetVersion() int8 { return tx.Tx.Version }

func (tx *TxToken) SetVersion(version int8) { tx.Tx.Version = version }

func (tx TxToken) GetMetadataType() int {
	if tx.Tx.Metadata != nil {
		return tx.Tx.Metadata.GetType()
	}
	return metadata.InvalidMeta
}

func (tx TxToken) GetType() string { return tx.Tx.Type }

func (tx *TxToken) SetType(t string) { tx.Tx.Type = t }

func (tx TxToken) GetLockTime() int64 { return tx.Tx.LockTime }

func (tx *TxToken) SetLockTime(locktime int64) { tx.Tx.LockTime = locktime }

func (tx TxToken) GetSenderAddrLastByte() byte { return tx.Tx.PubKeyLastByteSender }

func (tx *TxToken) SetGetSenderAddrLastByte(b byte) { tx.Tx.PubKeyLastByteSender = b }

func (tx TxToken) GetTxFee() uint64 { return tx.Tx.Fee }

func (tx *TxToken) SetTxFee(fee uint64) { tx.Tx.Fee = fee }

func (tx TxToken) GetTxFeeToken() uint64 { return uint64(0) }

func (tx TxToken) GetInfo() []byte { return tx.Tx.Info }

func (tx *TxToken) SetInfo(info []byte) { tx.Tx.Info = info }

// not supported
func (tx TxToken) GetSigPubKey() []byte           { return []byte{} }
func (tx *TxToken) SetSigPubKey(sigPubkey []byte) {}
func (tx TxToken) GetSig() []byte                 { return []byte{} }
func (tx *TxToken) SetSig(sig []byte)             {}
func (tx TxToken) GetProof() privacy.Proof        { return nil }
func (tx *TxToken) SetProof(proof privacy.Proof)  {}
func (tx TxToken) GetCachedActualSize() *uint64 {
	return nil
}
func (tx *TxToken) SetCachedActualSize(sz *uint64) {}

func (tx TxToken) GetCachedHash() *common.Hash {
	return nil
}
func (tx *TxToken) SetCachedHash(h *common.Hash) {}

func (tx TxToken) GetTokenID() *common.Hash { return &tx.TokenData.PropertyID }

func (tx TxToken) GetMetadata() metadata.Metadata { return tx.Tx.Metadata }

func (tx *TxToken) SetMetadata(meta metadata.Metadata) { tx.Tx.Metadata = meta }
func (tx TxToken) GetPrivateKey() []byte {
	return tx.Tx.GetPrivateKey()
}
func (tx *TxToken) SetPrivateKey(sk []byte) {
	tx.Tx.SetPrivateKey(sk)
}

func (tx TxToken) GetReceivers() ([][]byte, []uint64) {
	return nil, nil
}

func (tx TxToken) ListSerialNumbersHashH() []common.Hash {
	result := []common.Hash{}
	if tx.Tx.GetProof() != nil {
		for _, d := range tx.Tx.GetProof().GetInputCoins() {
			hash := common.HashH(d.GetKeyImage().ToBytesS())
			result = append(result, hash)
		}
	}
	if tx.GetTxNormal().GetProof() != nil {
		for _, d := range tx.GetTxNormal().GetProof().GetInputCoins() {
			hash := common.HashH(d.GetKeyImage().ToBytesS())
			result = append(result, hash)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})
	return result
}

func (tx TxToken) String() string {
	jsb, err := json.Marshal(tx)
	if err != nil {
		return ""
	}
	return string(jsb)
	// record := strconv.Itoa(int(tx.Tx.Version))
	// record += strconv.FormatInt(tx.Tx.LockTime, 10)
	// record += strconv.FormatUint(tx.Tx.Fee, 10)
	// if tx.Proof != nil {
	// 	record += base64.StdEncoding.EncodeToString(tx.Tx.Proof.Bytes())
	// }
	// if tx.Metadata != nil {
	// 	metadataHash := tx.Metadata.Hash()
	// 	record += metadataHash.String()
	// }
	// return record
}
func (tx *TxToken) CalculateTxValue() uint64 {
	proof := tx.GetTxNormal().GetProof()
	if proof == nil {
		return 0
	}
	if proof.GetOutputCoins() == nil || len(proof.GetOutputCoins()) == 0 {
		return 0
	}
	if proof.GetInputCoins() == nil || len(proof.GetInputCoins()) == 0 { // coinbase tx
		txValue := uint64(0)
		for _, outCoin := range proof.GetOutputCoins() {
			txValue += outCoin.GetValue()
		}
		return txValue
	}

	if tx.GetTxNormal().IsPrivacy() {
		return 0
	}

	senderPKBytes := proof.GetInputCoins()[0].GetPublicKey().ToBytesS()
	txValue := uint64(0)
	for _, outCoin := range proof.GetOutputCoins() {
		outPKBytes := outCoin.GetPublicKey().ToBytesS()
		if bytes.Equal(senderPKBytes, outPKBytes) {
			continue
		}
		txValue += outCoin.GetValue()
	}
	return txValue
}

func (tx TxToken) CheckTxVersion(maxTxVersion int8) bool {
	return !(tx.Tx.Version > maxTxVersion)
}

func (tx TxToken) IsSalaryTx() bool {
	if tx.Tx.GetType() != common.TxRewardType {
		return false
	}
	if len(tx.TokenData.Proof.GetInputCoins()) > 0 {
		return false
	}
	return true
}

func (tx TxToken) IsPrivacy() bool {
	// In the case of NonPrivacyNonInput, we do not have proof
	if tx.Tx.Proof == nil {
		return false
	}
	return tx.Tx.Proof.IsPrivacy()
}


func (tx *TxToken) GetReceiverData() ([]coin.Coin, error) {
	if tx.Tx.Proof != nil && len(tx.Tx.Proof.GetOutputCoins()) > 0 {
		return tx.Tx.Proof.GetOutputCoins(), nil
	}
	return nil, nil
}

func (txToken *TxToken) GetTransferData() (bool, []byte, uint64, *common.Hash) {
	pubkeys, amounts := txToken.GetTxNormal().GetReceivers()
	if len(pubkeys) == 0 {
		fmt.Printf("GetTransferData receive 0 output, it should has exactly 1 output")
		return false, nil, 0, &txToken.TokenData.PropertyID
	}
	if len(pubkeys) > 1 {
		fmt.Printf("GetTransferData receiver: More than 1 receiver")
		return false, nil, 0, &txToken.TokenData.PropertyID
	}
	return true, pubkeys[0], amounts[0], &txToken.TokenData.PropertyID
}

func (txToken *TxToken) GetTxTokenData() tx_generic.TxTokenData {
	return txToken.TokenData.ToCompatTokenData(txToken.GetTxNormal())
}

func (txToken *TxToken) SetTxTokenData(data tx_generic.TxTokenData) error {
	td, txN, err := decomposeTokenData(data)
	if err == nil {
		txToken.TokenData = *td
		return txToken.SetTxNormal(txN)
	}
	return err
}

func (txToken *TxToken) GetTxMintData() (bool, coin.Coin, *common.Hash, error) {
	tokenID := txToken.TokenData.PropertyID
	return tx_generic.GetTxMintData(txToken.GetTxNormal(), &tokenID)
}

func (txToken *TxToken) GetTxBurnData() (bool, coin.Coin, *common.Hash, error) {
	tokenID := txToken.TokenData.PropertyID
	isBurn, burnCoin, _, err := txToken.GetTxNormal().GetTxBurnData()
	return isBurn, burnCoin, &tokenID, err
}

func (txToken *TxToken) GetTxFullBurnData() (bool, coin.Coin, coin.Coin, *common.Hash, error) {
	isBurnToken, burnToken, burnedTokenID, errToken := txToken.GetTxBurnData()
	isBurnPrv, burnPrv, _, errPrv := txToken.GetTxBase().GetTxBurnData()

	if errToken != nil && errPrv != nil {
		return false, nil, nil, nil, fmt.Errorf("%v and %v", errPrv, errToken)
	}

	return isBurnPrv || isBurnToken, burnPrv, burnToken, burnedTokenID, nil
}

func (txToken *TxToken) UnmarshalJSON(data []byte) error {
	var err error
	type TxTokenHolder struct {
		Tx                 json.RawMessage
		TxTokenPrivacyData json.RawMessage
	}
	var holder TxTokenHolder
	if err = json.Unmarshal(data, &holder); err != nil {
		return err
	}

	if err = json.Unmarshal(holder.Tx, &txToken.Tx); err != nil {
		return err
	}

	switch txToken.Tx.Type {
	case common.TxTokenConversionType:
		if txToken.Tx.Version != utils.TxConversionVersion12Number {
			return utils.NewTransactionErr(utils.PrivacyTokenJsonError, errors.New("Error while unmarshalling TX token v2 : wrong proof version"))
		}
		txToken.TokenData.Proof = &privacy.ProofForConversion{}
		txToken.TokenData.Proof.Init()
	case common.TxCustomTokenPrivacyType:
		if txToken.Tx.Version != utils.TxVersion2Number {
			return utils.NewTransactionErr(utils.PrivacyTokenJsonError, errors.New("Error while unmarshalling TX token v2 : wrong proof version"))
		}
		txToken.TokenData.Proof = &privacy.ProofV2{}
		txToken.TokenData.Proof.Init()
	default:
		return utils.NewTransactionErr(utils.PrivacyTokenJsonError, errors.New("Error while unmarshalling TX token v2 : wrong proof type"))
	}

	err = json.Unmarshal(holder.TxTokenPrivacyData, &txToken.TokenData)
	if err != nil {
		fmt.Println(err)
		return utils.NewTransactionErr(utils.PrivacyTokenJsonError, err)
	}
	// proof := txToken.TokenData.Proof.(*privacy.ProofV2).GetAggregatedRangeProof().(*privacy.AggregatedRangeProofV2)
	// fmt.Printf("Unmarshalled proof into token data: %v\n", agg)
	txToken.cachedTxNormal = makeTxToken(&txToken.Tx, txToken.TokenData.SigPubKey, txToken.TokenData.Sig, txToken.TokenData.Proof)
	return nil
}

func (txToken TxToken) ListOTAHashH() []common.Hash {
	result := make([]common.Hash, 0)

	//Retrieve PRV output coins
	if txToken.GetTxBase().GetProof() != nil {
		for _, outputCoin := range txToken.GetTxBase().GetProof().GetOutputCoins() {
			//Discard coins sent to the burning address
			if wallet.IsPublicKeyBurningAddress(outputCoin.GetPublicKey().ToBytesS()) {
				continue
			}
			hash := common.HashH(outputCoin.GetPublicKey().ToBytesS())
			result = append(result, hash)
		}
	}

	//Retrieve token output coins
	if txToken.GetTxNormal().GetProof() != nil {
		for _, outputCoin := range txToken.GetTxNormal().GetProof().GetOutputCoins() {
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
