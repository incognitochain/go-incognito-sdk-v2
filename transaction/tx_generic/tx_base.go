package tx_generic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"math"
	"sort"
	"strconv"
	"time"
)

type TxBase struct {
	// Basic data, required
	Version  int8   `json:"Version"`
	Type     string `json:"Type"` // Transaction type
	LockTime int64  `json:"LockTime"`
	Fee      uint64 `json:"Fee"` // Fee applies: always consant
	Info     []byte // 512 bytes
	// Sign and Privacy proof, required
	SigPubKey            []byte `json:"SigPubKey, omitempty"`
	Sig                  []byte `json:"Sig, omitempty"`
	Proof                privacy.Proof
	PubKeyLastByteSender byte
	// Metadata, optional
	Metadata metadata.Metadata
	// private field, not use for json parser, only use as temp variable
	sigPrivKey       []byte       // is ALWAYS private property of struct, if privacy: 64 bytes, and otherwise, 32 bytes
	cachedHash       *common.Hash // cached hash data of tx
	cachedActualSize *uint64      // cached actualsize data for tx
}

type TxPrivacyInitParams struct {
	SenderSK    *key.PrivateKey
	PaymentInfo []*key.PaymentInfo
	InputCoins  []coin.PlainCoin
	Fee         uint64
	HasPrivacy  bool
	TokenID     *common.Hash // default is nil -> use for prv coin
	MetaData    metadata.Metadata
	Info        []byte // 512 bytes
	Kvargs      map[string]interface{}
}

func NewTxPrivacyInitParams(
	senderSK *key.PrivateKey,
	paymentInfo []*key.PaymentInfo,
	inputCoins []coin.PlainCoin,
	fee uint64,
	hasPrivacy bool,
	tokenID *common.Hash, // default is nil -> use for prv coin
	metaData metadata.Metadata,
	info []byte,
	kvargs map[string]interface{}) *TxPrivacyInitParams {
	if info == nil {
		info = []byte{}
	}
	params := &TxPrivacyInitParams{
		TokenID:     tokenID,
		HasPrivacy:  hasPrivacy,
		InputCoins:  inputCoins,
		Fee:         fee,
		MetaData:    metaData,
		PaymentInfo: paymentInfo,
		SenderSK:    senderSK,
		Info:        info,
		Kvargs:      kvargs,
	}
	return params
}

func GetTxInfo(paramInfo []byte) ([]byte, error) {
	if lenTxInfo := len(paramInfo); lenTxInfo > utils.MaxSizeInfo {
		return []byte{}, fmt.Errorf("length info (%v) exceeds max size (%v)", lenTxInfo, utils.MaxSizeInfo)
	}
	return paramInfo, nil
}

func CalculateSentBackInfo(params *TxPrivacyInitParams, senderPaymentAddress key.PaymentAddress) error {
	// Calculate sum of all output coins' value
	sumOutputValue := uint64(0)
	for _, p := range params.PaymentInfo {
		sumOutputValue += p.Amount
	}

	// Calculate sum of all input coins' value
	sumInputValue := uint64(0)
	for _, coin := range params.InputCoins {
		sumInputValue += coin.GetValue()
	}

	overBalance := int64(sumInputValue - sumOutputValue - params.Fee)
	// Check if sum of input coins' value is at least sum of output coins' value and tx fee
	if overBalance < 0 {
		return fmt.Errorf("sum of inputs less than outputs %v: sumInputValue=%d sumOutputValue=%d fee=%d", params.TokenID.String(), sumInputValue, sumOutputValue, params.Fee)
	}
	// Create a new payment to sender's pk where amount is overBalance if > 0
	if overBalance > 0 {
		// Should not check error because have checked before
		changePaymentInfo := new(key.PaymentInfo)
		changePaymentInfo.Amount = uint64(overBalance)
		changePaymentInfo.PaymentAddress = senderPaymentAddress
		params.PaymentInfo = append(params.PaymentInfo, changePaymentInfo)
	}

	return nil
}

func GetTxVersionFromCoins(inputCoins []coin.PlainCoin) (int8, error) {
	// If this is nonPrivacyNonInputCoins (maybe)
	if len(inputCoins) == 0 {
		return utils.CurrentTxVersion, nil
	}
	check := [3]bool{false, false, false}
	for i := 0; i < len(inputCoins); i += 1 {
		check[inputCoins[i].GetVersion()] = true
	}

	// If inputCoins contain 2 versions
	if check[1] && check[2] {
		return 0, errors.New("cannot get tx version because there are 2 versions of input coins")
	}

	// If somehow no version is checked???
	if !check[1] && !check[2] {
		return 0, errors.New("cannot get tx version, something is wrong with coins.version, it should be 1 or 2 only")
	}

	if check[2] {
		return 2, nil
	} else {
		return 1, nil
	}
}

// return bool indicates whether we should continue "Init" function or not
func (tx *TxBase) InitializeTxAndParams(params *TxPrivacyInitParams) error {
	var err error
	// Get Keyset from param
	senderKeySet := key.KeySet{}
	if err := senderKeySet.InitFromPrivateKey(params.SenderSK); err != nil {
		return fmt.Errorf("cannot parse Private Key. Err: %v", err)
	}
	tx.sigPrivKey = *params.SenderSK
	// Tx: initialize some values
	if tx.LockTime == 0 {
		tx.LockTime = time.Now().Unix()
	}
	tx.Fee = params.Fee
	tx.Type = common.TxNormalType
	tx.Metadata = params.MetaData
	tx.PubKeyLastByteSender = common.GetShardIDFromLastByte(senderKeySet.PaymentAddress.Pk[len(senderKeySet.PaymentAddress.Pk)-1])

	if tx.Version, err = GetTxVersionFromCoins(params.InputCoins); err != nil {
		return err
	}
	if tx.Info, err = GetTxInfo(params.Info); err != nil {
		return err
	}

	// Params: update balance if overbalance
	if err = CalculateSentBackInfo(params, senderKeySet.PaymentAddress); err != nil {
		return err
	}
	return nil
}

// =================== PARSING JSON FUNCTIONS ===================

func (tx *TxBase) UnmarshalJSON(data []byte) error {
	// For rolling version
	type Alias TxBase
	temp := &struct {
		Metadata *json.RawMessage
		Proof    *json.RawMessage
		*Alias
	}{
		Alias: (*Alias)(tx),
	}
	err := json.Unmarshal(data, temp)
	if err != nil {
		return err
	}

	if temp.Metadata == nil {
		tx.SetMetadata(nil)
	} else {
		metaInBytes, err := json.Marshal(temp.Metadata)
		if err != nil {
			return err
		}

		meta, parseErr := metadata.ParseMetadata(metaInBytes)
		if parseErr != nil {
			return parseErr
		}
		tx.SetMetadata(meta)
	}

	proofType := tx.Type
	if proofType == common.TxTokenConversionType{
		proofType = common.TxNormalType
	}

	if temp.Proof == nil {
		tx.SetProof(nil)
	} else {
		proof, proofErr := utils.ParseProof(temp.Proof, tx.Version, proofType)
		if proofErr != nil {
			return proofErr
		}
		tx.SetProof(proof)
	}
	return nil
}

// =================== GET/SET FUNCTIONS ===================

func (tx TxBase) GetVersion() int8 { return tx.Version }

func (tx *TxBase) SetVersion(version int8) { tx.Version = version }

func (tx TxBase) GetMetadataType() int {
	if tx.Metadata != nil {
		return tx.Metadata.GetType()
	}
	return metadata.InvalidMeta
}

func (tx TxBase) GetType() string { return tx.Type }

func (tx *TxBase) SetType(t string) { tx.Type = t }

func (tx TxBase) GetLockTime() int64 { return tx.LockTime }

func (tx *TxBase) SetLockTime(locktime int64) { tx.LockTime = locktime }

func (tx TxBase) GetSenderAddrLastByte() byte { return tx.PubKeyLastByteSender }

func (tx *TxBase) SetGetSenderAddrLastByte(b byte) { tx.PubKeyLastByteSender = b }

func (tx TxBase) GetTxFee() uint64 { return tx.Fee }

func (tx *TxBase) SetTxFee(fee uint64) { tx.Fee = fee }

func (tx TxBase) GetTxFeeToken() uint64 { return uint64(0) }

func (tx TxBase) GetInfo() []byte { return tx.Info }

func (tx *TxBase) SetInfo(info []byte) { tx.Info = info }

func (tx TxBase) GetSigPubKey() []byte { return tx.SigPubKey }

func (tx *TxBase) SetSigPubKey(sigPubkey []byte) { tx.SigPubKey = sigPubkey }

func (tx TxBase) GetSig() []byte { return tx.Sig }

func (tx *TxBase) SetSig(sig []byte) { tx.Sig = sig }

func (tx TxBase) GetProof() privacy.Proof { return tx.Proof }

func (tx *TxBase) SetProof(proof privacy.Proof) { tx.Proof = proof }

func (tx TxBase) GetTokenID() *common.Hash { return &common.PRVCoinID }

func (tx TxBase) GetMetadata() metadata.Metadata { return tx.Metadata }

func (tx *TxBase) SetMetadata(meta metadata.Metadata) { tx.Metadata = meta }

func (tx TxBase) GetPrivateKey() []byte{
	return tx.sigPrivKey
}

func (tx *TxBase) SetPrivateKey(sk []byte){
	tx.sigPrivKey = sk
}

func (tx TxBase) GetCachedActualSize() *uint64{
	return tx.cachedActualSize
}

func (tx *TxBase) SetCachedActualSize(sz *uint64){
	tx.cachedActualSize = sz
}

func (tx TxBase) GetCachedHash() *common.Hash{
	return tx.cachedHash
}

func (tx *TxBase) SetCachedHash(h *common.Hash){
	tx.cachedHash = h
}

// =================== FUNCTIONS THAT GET STUFF AND REQUIRE SOME CODING ===================

func (tx TxBase) GetTxActualSize() uint64 {
	//txBytes, _ := json.Marshal(tx)
	//txSizeInByte := len(txBytes)
	//
	//return uint64(math.Ceil(float64(txSizeInByte) / 1024))
	if tx.cachedActualSize != nil {
		return *tx.cachedActualSize
	}
	sizeTx := uint64(1)                // int8
	sizeTx += uint64(len(tx.Type) + 1) // string
	sizeTx += uint64(8)                // int64
	sizeTx += uint64(8)

	sigPubKey := uint64(len(tx.SigPubKey))
	sizeTx += sigPubKey
	sig := uint64(len(tx.Sig))
	sizeTx += sig
	if tx.Proof != nil {
		proof := uint64(len(tx.Proof.Bytes()))
		sizeTx += proof
	}

	sizeTx += uint64(1)
	info := uint64(len(tx.Info))
	sizeTx += info

	meta := tx.Metadata
	if meta != nil {
		metaSize := meta.CalculateSize()
		sizeTx += metaSize
	}
	result := uint64(math.Ceil(float64(sizeTx) / 1024))
	tx.cachedActualSize = &result
	return *tx.cachedActualSize
}

func (tx TxBase) GetReceivers() ([][]byte, []uint64) {
	pubkeys := [][]byte{}
	amounts := []uint64{}
	if tx.Proof != nil && len(tx.Proof.GetOutputCoins()) > 0 {
		for _, coin := range tx.Proof.GetOutputCoins() {
			added := false
			coinPubKey := coin.GetPublicKey().ToBytesS()
			for i, key := range pubkeys {
				if bytes.Equal(coinPubKey, key) {
					added = true
					amounts[i] += coin.GetValue()
					break
				}
			}
			if !added {
				pubkeys = append(pubkeys, coinPubKey)
				amounts = append(amounts, coin.GetValue())
			}
		}
	}
	return pubkeys, amounts
}

func (tx TxBase) GetTransferData() (bool, []byte, uint64, *common.Hash) {
	pubkeys, amounts := tx.GetReceivers()
	if len(pubkeys) == 0 {
		return false, nil, 0, &common.PRVCoinID
	}
	if len(pubkeys) > 1 {
		return false, nil, 0, &common.PRVCoinID
	}
	return true, pubkeys[0], amounts[0], &common.PRVCoinID
}

func (tx TxBase) ListSerialNumbersHashH() []common.Hash {
	result := []common.Hash{}
	if tx.Proof != nil {
		for _, d := range tx.Proof.GetInputCoins() {
			hash := common.HashH(d.GetKeyImage().ToBytesS())
			result = append(result, hash)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})
	return result
}

func (tx TxBase) String() string {
	record := strconv.Itoa(int(tx.Version))
	record += strconv.FormatInt(tx.LockTime, 10)
	record += strconv.FormatUint(tx.Fee, 10)
	if tx.Proof != nil {
		//tmp := base58.Base58Check{}.Encode(tx.Proof.Bytes(), 0x00)
		record += base64.StdEncoding.EncodeToString(tx.Proof.Bytes())
		// fmt.Printf("Proof check base 58: %v\n",tmp)
	}
	if tx.Metadata != nil {
		metadataHash := tx.Metadata.Hash()
		record += metadataHash.String()
		//Logger.log.Debugf("\n\n\n\n test metadata after hashing: %v\n", metadataHash.GetBytes())
	}

	// TODO: To be uncomment
	// record += string(tx.Info)

	return record
}

func (tx *TxBase) Hash() *common.Hash {
	if tx.cachedHash != nil {
		return tx.cachedHash
	}
	inBytes := []byte(tx.String())
	hash := common.HashH(inBytes)
	return &hash
}

func (tx *TxBase) HashWithoutMetadataSig() *common.Hash {
	// hashing to sign metadata is version-specific
	return nil
}

func (tx TxBase) CalculateTxValue() uint64 {
	if tx.Proof == nil {
		return 0
	}

	outputCoins := tx.Proof.GetOutputCoins()
	inputCoins := tx.Proof.GetInputCoins()
	if outputCoins == nil || len(outputCoins) == 0 {
		return 0
	}
	if inputCoins == nil || len(inputCoins) == 0 { // coinbase tx
		txValue := uint64(0)
		for _, outCoin := range outputCoins {
			txValue += outCoin.GetValue()
		}
		return txValue
	}

	senderPKBytes := inputCoins[0].GetPublicKey().ToBytesS()
	txValue := uint64(0)
	for _, outCoin := range outputCoins {
		outPKBytes := outCoin.GetPublicKey().ToBytesS()
		if bytes.Equal(senderPKBytes, outPKBytes) {
			continue
		}
		txValue += outCoin.GetValue()
	}
	return txValue
}

// =================== FUNCTION THAT CHECK STUFFS  ===================

func (tx TxBase) CheckTxVersion(maxTxVersion int8) bool {
	return !(tx.Version > maxTxVersion)
}

func (tx *TxBase) IsNonPrivacyNonInput(params *TxPrivacyInitParams) (bool, error) {
	var err error
	//Logger.Log.Debugf("len(inputCoins), fee, hasPrivacy: %d, %d, %v\n", len(params.InputCoins), params.Fee, params.HasPrivacy)
	if len(params.InputCoins) == 0 && params.Fee == 0 && !params.HasPrivacy {
		//Logger.Log.Debugf("len(inputCoins) == 0 && fee == 0 && !hasPrivacy\n")
		tx.sigPrivKey = *params.SenderSK
		if tx.Sig, tx.SigPubKey, err = SignNoPrivacy(params.SenderSK, tx.Hash()[:]); err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (tx TxBase) IsSalaryTx() bool {
	if tx.GetType() != common.TxRewardType {
		return false
	}
	if len(tx.Proof.GetInputCoins()) > 0 {
		return false
	}
	return true
}

func (tx TxBase) IsPrivacy() bool {
	// In the case of NonPrivacyNonInput, we do not have proof
	if tx.Proof == nil {
		return false
	}
	return tx.Proof.IsPrivacy()
}
// =================== FUNCTIONS THAT VALIDATE STUFFS ===================
func (tx TxBase) ListOTAHashH() []common.Hash {
	return []common.Hash{}
}
