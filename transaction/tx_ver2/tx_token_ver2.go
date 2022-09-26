package tx_ver2

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// TxTokenDataVersion2 represents all data of a token transaction v2.
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

// Hash calculates the hash of a TxTokenDataVersion2.
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

// ToCompatTokenData creates a TxTokenData from a TxTokenDataVersion2 and the given transaction.
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

// TxToken represents a token transaction of version 2. A token transaction v2 consists of 2 sub-transactions
//   - TxBase: PRV sub-transaction for paying the transaction fee. All transactions v2 pay fees in PRV.
//   - TxNormal: the token sub-transaction to transfer token.
type TxToken struct {
	Tx             Tx                  `json:"Tx"`
	TokenData      TxTokenDataVersion2 `json:"TxTokenPrivacyData"`
	cachedTxNormal *Tx
}

// GetTxBase returns the PRV sub-transaction of a TxToken.
func (txToken *TxToken) GetTxBase() metadata.Transaction {
	return &txToken.Tx
}

// GetTxNormal returns the token sub-transaction of a TxToken.
func (txToken *TxToken) GetTxNormal() metadata.Transaction {
	if txToken.cachedTxNormal != nil {
		return txToken.cachedTxNormal
	}
	result := makeTxToken(&txToken.Tx, txToken.TokenData.SigPubKey, txToken.TokenData.Sig, txToken.TokenData.Proof)
	// tx.cachedTxNormal = result
	return result
}

// GetVersion returns the version of a TxToken.
func (txToken TxToken) GetVersion() int8 { return txToken.Tx.Version }

// GetTokenID returns the tokenID of a TxToken.
func (txToken TxToken) GetTokenID() *common.Hash { return &txToken.TokenData.PropertyID }

// GetMetadata returns the metadata of a TxToken.
func (txToken TxToken) GetMetadata() metadata.Metadata { return txToken.Tx.Metadata }

// GetMetadataType returns the metadata type of a TxToken.
func (txToken TxToken) GetMetadataType() int {
	if txToken.Tx.Metadata != nil {
		return txToken.Tx.Metadata.GetType()
	}
	return metadata.InvalidMeta
}

// GetType returns the transaction type of a TxToken.
func (txToken TxToken) GetType() string { return txToken.Tx.Type }

// GetLockTime returns the lock-time of a TxToken.
func (txToken TxToken) GetLockTime() int64 { return txToken.Tx.LockTime }

// GetTxActualSize returns the size of a TxToken in kb.
func (txToken TxToken) GetTxActualSize() uint64 {
	jsb, err := json.Marshal(txToken)
	if err != nil {
		return 0
	}
	return uint64(math.Ceil(float64(len(jsb)) / 1024))
}

// GetSenderAddrLastByte returns the last byte of the sender of a TxToken.
func (txToken TxToken) GetSenderAddrLastByte() byte { return txToken.Tx.PubKeyLastByteSender }

// GetTxFee returns the PRV fee of a TxToken.
func (txToken TxToken) GetTxFee() uint64 { return txToken.Tx.Fee }

// GetTxFeeToken returns the token fee of a TxToken.
// All transactions v2 pay fees in PRV, so it returns 0.
func (txToken TxToken) GetTxFeeToken() uint64 { return uint64(0) }

// GetInfo returns the info of a TxToken.
func (txToken TxToken) GetInfo() []byte { return txToken.Tx.Info }

// GetSigPubKey not supported.
func (txToken TxToken) GetSigPubKey() []byte { return []byte{} }

// GetSig not supported.
func (txToken TxToken) GetSig() []byte { return []byte{} }

// GetProof not supported.
func (txToken TxToken) GetProof() privacy.Proof { return nil }

// GetReceivers not supported.
func (txToken TxToken) GetReceivers() ([][]byte, []uint64) {
	return nil, nil
}

// GetReceiverData returns the output coins of a TxToken.
func (txToken *TxToken) GetReceiverData() ([]coin.Coin, error) {
	if txToken.Tx.Proof != nil && len(txToken.Tx.Proof.GetOutputCoins()) > 0 {
		return txToken.Tx.Proof.GetOutputCoins(), nil
	}
	return nil, nil
}

// GetTransferData returns the transferred data of a TxToken.
func (txToken *TxToken) GetTransferData() (bool, []byte, uint64, *common.Hash) {
	pubKeys, amounts := txToken.GetTxNormal().GetReceivers()
	if len(pubKeys) == 0 {
		log.Printf("GetTransferData receive 0 output, it should has exactly 1 output")
		return false, nil, 0, &txToken.TokenData.PropertyID
	}
	if len(pubKeys) > 1 {
		log.Printf("GetTransferData receiver: More than 1 receiver")
		return false, nil, 0, &txToken.TokenData.PropertyID
	}
	return true, pubKeys[0], amounts[0], &txToken.TokenData.PropertyID
}

// GetTxTokenData returns the token data of a TxToken.
func (txToken *TxToken) GetTxTokenData() tx_generic.TxTokenData {
	return txToken.TokenData.ToCompatTokenData(txToken.GetTxNormal())
}

// GetTxMintData returns the minting data of a TxToken.
func (txToken *TxToken) GetTxMintData() (bool, coin.Coin, *common.Hash, error) {
	tokenID := txToken.TokenData.PropertyID
	return tx_generic.GetTxMintData(txToken.GetTxNormal(), &tokenID)
}

// GetTxBurnData returns the burning (token only) of a TxToken.
func (txToken *TxToken) GetTxBurnData() (bool, coin.Coin, *common.Hash, error) {
	tokenID := txToken.TokenData.PropertyID
	isBurn, burnCoin, _, err := txToken.GetTxNormal().GetTxBurnData()
	return isBurn, burnCoin, &tokenID, err
}

// GetTxFullBurnData returns the full burning data (both PRV and token) of a TxToken.
func (txToken *TxToken) GetTxFullBurnData() (bool, coin.Coin, coin.Coin, *common.Hash, error) {
	isBurnToken, burnToken, burnedTokenID, errToken := txToken.GetTxBurnData()
	isBurnPrv, burnPrv, _, errPrv := txToken.GetTxBase().GetTxBurnData()

	if errToken != nil && errPrv != nil {
		return false, nil, nil, nil, fmt.Errorf("%v and %v", errPrv, errToken)
	}

	return isBurnPrv || isBurnToken, burnPrv, burnToken, burnedTokenID, nil
}

// CheckTxVersion checks if the version of a TxToken is valid.
func (txToken TxToken) CheckTxVersion(maxTxVersion int8) bool {
	return !(txToken.Tx.Version > maxTxVersion)
}

// IsSalaryTx checks if a TxToken is a salary transaction.
func (txToken TxToken) IsSalaryTx() bool {
	if txToken.Tx.GetType() != common.TxRewardType {
		return false
	}
	if len(txToken.TokenData.Proof.GetInputCoins()) > 0 {
		return false
	}
	return true
}

// IsPrivacy checks if a TxToken is a private transaction.
func (txToken TxToken) IsPrivacy() bool {
	// In the case of NonPrivacyNonInput, we do not have proof
	if txToken.Tx.Proof == nil {
		return false
	}
	return txToken.Tx.Proof.IsPrivacy()
}

// ListSerialNumbersHashH returns the hash list of all serial numbers in a TxToken.
func (txToken TxToken) ListSerialNumbersHashH() []common.Hash {
	result := make([]common.Hash, 0)
	if txToken.Tx.GetProof() != nil {
		for _, d := range txToken.Tx.GetProof().GetInputCoins() {
			hash := common.HashH(d.GetKeyImage().ToBytesS())
			result = append(result, hash)
		}
	}
	if txToken.GetTxNormal().GetProof() != nil {
		for _, d := range txToken.GetTxNormal().GetProof().GetInputCoins() {
			hash := common.HashH(d.GetKeyImage().ToBytesS())
			result = append(result, hash)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})
	return result
}

// ListOTAHashH returns the hash list of all OTA keys in a TxToken.
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

// GetPrivateKey returns the sender's private key.
func (txToken TxToken) GetPrivateKey() []byte {
	return txToken.Tx.GetPrivateKey()
}

// Hash calculates the hash of a TxToken.
func (txToken *TxToken) Hash() *common.Hash {
	firstHash := txToken.Tx.Hash()
	secondHash, err := txToken.TokenData.Hash()
	if err != nil {
		return nil
	}
	result := common.HashH(append(firstHash[:], secondHash[:]...))
	return &result
}

// HashWithoutMetadataSig calculates the hash of a TxToken with out adding the signature of its metadata.
func (txToken TxToken) HashWithoutMetadataSig() *common.Hash {
	return txToken.Tx.HashWithoutMetadataSig()
}

// String returns the string-representation of a TxToken.
func (txToken TxToken) String() string {
	jsb, err := json.Marshal(txToken)
	if err != nil {
		return ""
	}
	return string(jsb)
}

// SetTxBase sets v as the TxBase of a TxToken.
func (txToken *TxToken) SetTxBase(v metadata.Transaction) error {
	temp, ok := v.(*Tx)
	if !ok {
		return fmt.Errorf("cannot set TxBase: wrong type")
	}
	txToken.Tx = *temp
	return nil
}

// SetTxNormal sets v as the TxNormal of a TxToken.
func (txToken *TxToken) SetTxNormal(v metadata.Transaction) error {
	temp, ok := v.(*Tx)
	if !ok {
		return utils.NewTransactionErr(utils.UnexpectedError, fmt.Errorf("cannot set TxNormal: wrong type"))
	}
	txToken.TokenData.SigPubKey = temp.SigPubKey
	txToken.TokenData.Sig = temp.Sig
	txToken.TokenData.Proof = temp.Proof
	txToken.cachedTxNormal = temp
	return nil
}

// SetTxTokenData sets v as the token data of a TxToken.
func (txToken *TxToken) SetTxTokenData(v tx_generic.TxTokenData) error {
	td, txN, err := decomposeTokenData(v)
	if err == nil {
		txToken.TokenData = *td
		return txToken.SetTxNormal(txN)
	}
	return err
}

// SetVersion sets v as the version of a TxToken.
func (txToken *TxToken) SetVersion(v int8) { txToken.Tx.Version = v }

// SetType sets v as the transaction type of a TxToken.
func (txToken *TxToken) SetType(v string) { txToken.Tx.Type = v }

// SetLockTime sets v as the lock-time of a TxToken.
func (txToken *TxToken) SetLockTime(v int64) { txToken.Tx.LockTime = v }

// SetGetSenderAddrLastByte sets v as the sender's last byte of a TxToken.
func (txToken *TxToken) SetGetSenderAddrLastByte(v byte) { txToken.Tx.PubKeyLastByteSender = v }

// SetTxFee sets v as the PRV fee of a TxToken.
func (txToken *TxToken) SetTxFee(v uint64) { txToken.Tx.Fee = v }

// SetInfo sets v as the info of a TxToken.
func (txToken *TxToken) SetInfo(v []byte) { txToken.Tx.Info = v }

// SetSigPubKey not supported.
func (txToken *TxToken) SetSigPubKey([]byte) {}

// SetSig not supported.
func (txToken *TxToken) SetSig([]byte) {}

// SetProof not supported.
func (txToken *TxToken) SetProof(privacy.Proof) {}

// SetMetadata sets v the metadata of a TxToken.
func (txToken *TxToken) SetMetadata(v metadata.Metadata) { txToken.Tx.Metadata = v }

// SetPrivateKey sets v as the private key of a TxToken.
func (txToken *TxToken) SetPrivateKey(v []byte) {
	txToken.Tx.SetPrivateKey(v)
}

// Init creates a token transaction version 2 from the given parameter.
// The input parameter should be a *tx_generic.TxTokenParams.
func (txToken *TxToken) Init(txTokenParams interface{}) error {
	params, ok := txTokenParams.(*tx_generic.TxTokenParams)
	if !ok {
		return fmt.Errorf("cannot parse the input as a TxTokenParams")
	}

	if params.TokenParams.Fee > 0 || params.FeeNativeCoin == 0 {
		log.Printf("only accept tx fee in PRV")
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

	if check, err := tx.IsNonPrivacyNonInput(txPrivacyParams); check {
		return err
	}

	// Init PRV Fee
	ins, outs, err := txToken.initPRV(tx, txPrivacyParams)
	if err != nil {
		log.Printf("Cannot init PRV fee for tokenver2: err %v", err)
		return err
	}

	txn := makeTxToken(tx, nil, nil, nil)
	// Init, prove and sign(CA) Token
	if err := txToken.initToken(txn, params); err != nil {
		log.Printf("Cannot init token ver2: err %v", err)
		return err
	}
	tdh, err := txToken.TokenData.Hash()
	if err != nil {
		return err
	}
	message := common.HashH(append(tx.Hash()[:], tdh[:]...))
	err = tx.signOnMessage(ins, outs, txPrivacyParams, message[:])
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

// CalculateTxValue calculates total output values.
func (txToken *TxToken) CalculateTxValue() uint64 {
	proof := txToken.GetTxNormal().GetProof()
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

	if txToken.GetTxNormal().IsPrivacy() {
		return 0
	}

	txValue := uint64(0)
	for _, outCoin := range proof.GetOutputCoins() {
		txValue += outCoin.GetValue()
	}
	return txValue
}

// UnmarshalJSON does the JSON-unmarshalling operation for a TxToken.
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
			return utils.NewTransactionErr(utils.PrivacyTokenJsonError, fmt.Errorf("error while unmarshalling TX token v2: wrong proof version"))
		}
		txToken.TokenData.Proof = &privacy.ProofForConversion{}
		txToken.TokenData.Proof.Init()
	case common.TxCustomTokenPrivacyType:
		if txToken.Tx.Version != utils.TxVersion2Number {
			return utils.NewTransactionErr(utils.PrivacyTokenJsonError, fmt.Errorf("error while unmarshalling TX token v2: wrong proof version"))
		}
		txToken.TokenData.Proof = &privacy.ProofV2{}
		txToken.TokenData.Proof.Init()
	default:
		return utils.NewTransactionErr(utils.PrivacyTokenJsonError, fmt.Errorf("error while unmarshalling TX token v2: wrong proof type"))
	}

	err = json.Unmarshal(holder.TxTokenPrivacyData, &txToken.TokenData)
	if err != nil {
		fmt.Println(err)
		return utils.NewTransactionErr(utils.PrivacyTokenJsonError, err)
	}
	// proof := txToken.TokenData.Proof.(*privacy.ProofV2).GetAggregatedRangeProof().(*privacy.AggregatedRangeProofV2)
	// log.Printf("Unmarshalled proof into token data: %v\n", agg)
	txToken.cachedTxNormal = makeTxToken(&txToken.Tx, txToken.TokenData.SigPubKey, txToken.TokenData.Sig, txToken.TokenData.Proof)
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
		return utils.NewTransactionErr(utils.PrivacyTokenTxTypeNotHandleError, fmt.Errorf("can't handle this TokenTxType"))
	}
}

// this signs on the hash of both sub TXs
func (tx *Tx) provePRV(params *tx_generic.TxPrivacyInitParams) ([]coin.PlainCoin, []*coin.CoinV2, error) {
	var err error
	outputCoins := make([]*coin.CoinV2, 0)
	for _, paymentInfo := range params.PaymentInfo {
		// We do not mind duplicated OTAs, server will handle them.
		outputCoin, seal, err := coin.NewCoinFromPaymentInfo(coin.NewTransferCoinParams(paymentInfo, params.GetSenderShard()))
		if err != nil {
			log.Printf("Cannot parse outputCoinV2 to outputCoins, error %v\n", err)
			return nil, nil, err
		}
		_ = seal //TODO: export
		outputCoins = append(outputCoins, outputCoin)
	}

	// inputCoins is plainCoin because it may have coinV1 with coinV2
	inputCoins := params.InputCoins

	tx.Proof, err = privacy.ProveV2(inputCoins, outputCoins, nil, false, params.PaymentInfo)
	if err != nil {
		log.Printf("Error in privacy_v2.Prove, error %v ", err)
		return nil, nil, err
	}

	if tx.GetMetadata() != nil {
		if err := tx.GetMetadata().Sign(params.SenderSK, tx); err != nil {
			log.Printf("Cannot signOnMessage txMetadata in shouldSignMetadata")
			return nil, nil, err
		}
	}

	// Get Hash of the whole txToken then sign on it
	// message := common.HashH(append(tx.Hash()[:], hashedTokenMessage...))

	return inputCoins, outputCoins, nil
}

func (txToken *TxToken) initPRV(feeTx *Tx, params *tx_generic.TxPrivacyInitParams) ([]coin.PlainCoin, []*coin.CoinV2, error) {
	feeTx.SetType(common.TxCustomTokenPrivacyType)
	ins, outs, err := feeTx.provePRV(params)
	if err != nil {
		return nil, nil, utils.NewTransactionErr(utils.PrivacyTokenInitPRVError, err)
	}
	// override TxCustomTokenPrivacyType type
	// txToken.SetTxBase(feeTx)

	return ins, outs, nil
}

// makeTxToken creates the token sub-transaction given its proof, sig, and the PRV sub-transaction.
func makeTxToken(txPRV *Tx, pubKey, sig []byte, proof privacy.Proof) *Tx {
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
			log.Printf("Cannot parse proof %s using version %v - type %v", string(jsb), txPRV.Version, txPRV.Type)
			return nil
		}
	}
	var clonedSig []byte = nil
	if sig != nil {
		clonedSig = make([]byte, len(sig))
		copy(clonedSig, sig)
	}
	var clonedPk []byte = nil
	if pubKey != nil {
		clonedPk = make([]byte, len(pubKey))
		copy(clonedPk, pubKey)
	}
	result.Info = clonedInfo
	result.Proof = clonedProof
	result.Sig = clonedSig
	result.SigPubKey = clonedPk
	result.Info = clonedInfo

	return result
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
		return nil, nil, fmt.Errorf("error while casting a transaction to v2")
	}
	return &result, tx, nil
}
