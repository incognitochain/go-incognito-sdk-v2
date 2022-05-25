package tx_generic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"math"
	"sort"
	"strconv"
	"time"
)

// TxBase represents a PRV transaction field. It is used in both TxVer1 and TxVer2.
type TxBase struct {
	Version  int8   `json:"Version"`
	Type     string `json:"Type"`
	LockTime int64  `json:"LockTime"`
	Fee      uint64 `json:"Fee"`
	Info     []byte // 512 bytes
	// Sign and Privacy proof, required
	SigPubKey            []byte `json:"SigPubKey, omitempty"`
	Sig                  []byte `json:"Sig, omitempty"`
	Proof                privacy.Proof
	PubKeyLastByteSender byte
	Metadata             metadata.Metadata

	sigPrivateKey    []byte
	cachedHash       *common.Hash
	cachedActualSize *uint64
}

// TxPrivacyInitParams consists of parameters used to create a new PRV transaction.
type TxPrivacyInitParams struct {
	SenderSK    *key.PrivateKey
	PaymentInfo []*key.PaymentInfo
	InputCoins  []coin.PlainCoin
	Fee         uint64
	HasPrivacy  bool
	TokenID     *common.Hash // default is nil -> use for prv coin
	MetaData    metadata.Metadata
	Info        []byte // 512 bytes
	KvArgs      map[string]interface{}
}

// NewTxPrivacyInitParams creates a new TxPrivacyInitParams based on the given inputs.
func NewTxPrivacyInitParams(
	senderSK *key.PrivateKey,
	paymentInfo []*key.PaymentInfo,
	inputCoins []coin.PlainCoin,
	fee uint64,
	hasPrivacy bool,
	tokenID *common.Hash, // default is nil -> use for prv coin
	metaData metadata.Metadata,
	info []byte,
	kvaArgs map[string]interface{}) *TxPrivacyInitParams {
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
		KvArgs:      kvaArgs,
	}
	return params
}

// GetSenderShard returns the shardID of the sender of a TxPrivacyInitParams.
func (param *TxPrivacyInitParams) GetSenderShard() byte {
	pubKey := new(crypto.Point).ScalarMultBase(new(crypto.Scalar).FromBytesS(*param.SenderSK))
	pubKeyBytes := pubKey.ToBytesS()

	return common.GetShardIDFromLastByte(pubKeyBytes[len(pubKeyBytes)-1])
}

// GetTxInfo checks and returns valid info.
func GetTxInfo(paramInfo []byte) ([]byte, error) {
	if lenTxInfo := len(paramInfo); lenTxInfo > utils.MaxSizeInfo {
		return []byte{}, fmt.Errorf("length info (%v) exceeds max size (%v)", lenTxInfo, utils.MaxSizeInfo)
	}
	return paramInfo, nil
}

// CalculateSentBackInfo calculates the remaining amount to send back to the sender.
func CalculateSentBackInfo(params *TxPrivacyInitParams, senderPaymentAddress key.PaymentAddress) error {
	// Calculate sum of all output coins' value
	sumOutputValue := uint64(0)
	for _, p := range params.PaymentInfo {
		sumOutputValue += p.Amount
	}

	// Calculate sum of all input coins' value
	sumInputValue := uint64(0)
	for _, c := range params.InputCoins {
		sumInputValue += c.GetValue()
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

// GetTxVersionFromCoins returns the version of a list of input coins.
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

// InitializeTxAndParams initializes a new TxBase with values, prepared for the next steps.
func (tx *TxBase) InitializeTxAndParams(params *TxPrivacyInitParams) error {
	var err error
	senderKeySet := key.KeySet{}
	if err := senderKeySet.InitFromPrivateKey(params.SenderSK); err != nil {
		return fmt.Errorf("cannot parse Private Key. Err: %v", err)
	}

	tx.sigPrivateKey = *params.SenderSK
	// Tx: initialize some values
	if tx.LockTime == 0 {
		tx.LockTime = time.Now().Unix() - (1 + common.RandInt64()%100)
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

// UnmarshalJSON does the JSON-unmarshalling operation for a TxBase.
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
	if proofType == common.TxTokenConversionType {
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

// GetVersion returns the version of a TxBase.
func (tx TxBase) GetVersion() int8 { return tx.Version }

// GetMetadataType returns the metadata type of a TxBase.
func (tx TxBase) GetMetadataType() int {
	if tx.Metadata != nil {
		return tx.Metadata.GetType()
	}
	return metadata.InvalidMeta
}

// GetType returns the transaction type of a TxBase.
func (tx TxBase) GetType() string { return tx.Type }

// GetLockTime returns the lock-time of a TxBase.
func (tx TxBase) GetLockTime() int64 { return tx.LockTime }

// GetSenderAddrLastByte returns the last byte of the sender of a transaction.
func (tx TxBase) GetSenderAddrLastByte() byte { return tx.PubKeyLastByteSender }

// GetTxFee returns the PRV fee of a TxBase.
func (tx TxBase) GetTxFee() uint64 { return tx.Fee }

// GetTxFeeToken returns the token fee of a TxBase. For a TxBase, it returns 0.
func (tx TxBase) GetTxFeeToken() uint64 { return uint64(0) }

// GetInfo returns the info of a TxBase.
func (tx TxBase) GetInfo() []byte { return tx.Info }

// GetSigPubKey returns the sigPubKey of a TxBase.
func (tx TxBase) GetSigPubKey() []byte { return tx.SigPubKey }

// GetSig returns the signature of a TxBase.
func (tx TxBase) GetSig() []byte { return tx.Sig }

// GetProof returns the payment proof of a TxBase.
func (tx TxBase) GetProof() privacy.Proof { return tx.Proof }

// GetTokenID returns the tokenID of a TxBase. For a TxBase, it returns the tokenID of PRV.
func (tx TxBase) GetTokenID() *common.Hash { return &common.PRVCoinID }

// GetMetadata returns the metadata of a TxBase.
func (tx TxBase) GetMetadata() metadata.Metadata { return tx.Metadata }

// GetPrivateKey returns the sigPrivateKey of a TxBase.
func (tx TxBase) GetPrivateKey() []byte {
	return tx.sigPrivateKey
}

// GetTxActualSize returns the size of a TxBase in kb.
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

// GetReceivers returns a list of receivers (public keys) and a list of corresponding amounts of a TxBase.
func (tx TxBase) GetReceivers() ([][]byte, []uint64) {
	pubKeys := make([][]byte, 0)
	amounts := make([]uint64, 0)
	if tx.Proof != nil && len(tx.Proof.GetOutputCoins()) > 0 {
		for _, c := range tx.Proof.GetOutputCoins() {
			added := false
			coinPubKey := c.GetPublicKey().ToBytesS()
			for i, k := range pubKeys {
				if bytes.Equal(coinPubKey, k) {
					added = true
					amounts[i] += c.GetValue()
					break
				}
			}
			if !added {
				pubKeys = append(pubKeys, coinPubKey)
				amounts = append(amounts, c.GetValue())
			}
		}
	}
	return pubKeys, amounts
}

// GetTransferData returns the transferred data (receivers and amounts) of a TxBase.
func (tx TxBase) GetTransferData() (bool, []byte, uint64, *common.Hash) {
	pubKeys, amounts := tx.GetReceivers()
	if len(pubKeys) == 0 {
		return false, nil, 0, &common.PRVCoinID
	}
	if len(pubKeys) > 1 {
		return false, nil, 0, &common.PRVCoinID
	}
	return true, pubKeys[0], amounts[0], &common.PRVCoinID
}

// SetVersion sets v as the version of a TxBase.
func (tx *TxBase) SetVersion(v int8) { tx.Version = v }

// SetType sets v as the type of a TxBase.
func (tx *TxBase) SetType(v string) { tx.Type = v }

// SetLockTime sets v as the lock-time of a TxBase.
func (tx *TxBase) SetLockTime(v int64) { tx.LockTime = v }

// SetGetSenderAddrLastByte sets v as the last byte of the sender of a transaction.
func (tx *TxBase) SetGetSenderAddrLastByte(v byte) { tx.PubKeyLastByteSender = v }

// SetTxFee sets v as the PRV fee of a TxBase.
func (tx *TxBase) SetTxFee(v uint64) { tx.Fee = v }

// SetInfo sets v as the info of a TxBase.
func (tx *TxBase) SetInfo(v []byte) { tx.Info = v }

// SetSigPubKey sets v as the sigPubKey of a TxBase.
func (tx *TxBase) SetSigPubKey(v []byte) { tx.SigPubKey = v }

// SetSig sets v as the signature of a TxBase.
func (tx *TxBase) SetSig(v []byte) { tx.Sig = v }

// SetProof sets v as the payment proof of a TxBase.
func (tx *TxBase) SetProof(v privacy.Proof) { tx.Proof = v }

// SetMetadata sets v as the metadata of a TxBase.
func (tx *TxBase) SetMetadata(v metadata.Metadata) { tx.Metadata = v }

// SetPrivateKey sets v as the sigPrivateKey of a TxBase.
func (tx *TxBase) SetPrivateKey(v []byte) {
	tx.sigPrivateKey = v
}

// ListSerialNumbersHashH returns the hash list of all serial numbers in a TxBase.
func (tx TxBase) ListSerialNumbersHashH() []common.Hash {
	result := make([]common.Hash, 0)
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

// String returns the string-representation of a TxBase.
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
	}

	// TODO: To be uncomment
	// record += string(tx.Info)

	return record
}

// Hash calculates the hash of a TxBase.
func (tx *TxBase) Hash() *common.Hash {
	if tx.cachedHash != nil {
		return tx.cachedHash
	}
	inBytes := []byte(tx.String())
	hash := common.HashH(inBytes)
	return &hash
}

// HashWithoutMetadataSig calculates the hash of a TxBase with out adding the signature of its metadata.
func (tx *TxBase) HashWithoutMetadataSig() *common.Hash {
	// hashing to sign metadata is version-specific
	return nil
}

// CalculateTxValue calculates total output values (not including the coins which are sent back to the sender).
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

// CheckTxVersion checks if the version of a TxBase is valid.
func (tx TxBase) CheckTxVersion(maxTxVersion int8) bool {
	return !(tx.Version > maxTxVersion)
}

// IsNonPrivacyNonInput checks if a TxBase is a non-private and non-input transaction.
func (tx *TxBase) IsNonPrivacyNonInput(params *TxPrivacyInitParams) (bool, error) {
	var err error
	if len(params.InputCoins) == 0 && params.Fee == 0 && !params.HasPrivacy {
		tx.sigPrivateKey = *params.SenderSK
		if tx.Sig, tx.SigPubKey, err = SignNoPrivacy(params.SenderSK, tx.Hash()[:]); err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

// IsSalaryTx checks if a TxBase is a salary transaction.
// A salary transaction is a transaction with 0 input and at least 1 output.
func (tx TxBase) IsSalaryTx() bool {
	if tx.GetType() != common.TxRewardType {
		return false
	}
	if len(tx.Proof.GetInputCoins()) > 0 {
		return false
	}
	return true
}

// IsPrivacy checks if a TxBase is a private transaction.
func (tx TxBase) IsPrivacy() bool {
	// In the case of NonPrivacyNonInput, we do not have proof
	if tx.Proof == nil {
		return false
	}
	return tx.Proof.IsPrivacy()
}

// ListOTAHashH returns the hash list of all OTA keys in a TxBase.
func (tx TxBase) ListOTAHashH() []common.Hash {
	return []common.Hash{}
}
