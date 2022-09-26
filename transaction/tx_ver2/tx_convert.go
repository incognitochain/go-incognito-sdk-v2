package tx_ver2

import (
	"fmt"
	"log"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/conversion"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy/v1/zkp/serialnumbernoprivacy"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_generic"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/utils"

	"strconv"
	"time"
)

// TxConvertVer1ToVer2InitParams consists of parameters used to create a new PRV conversion transaction.
type TxConvertVer1ToVer2InitParams struct {
	senderSK    *key.PrivateKey
	paymentInfo []*coin.PaymentInfo
	inputCoins  []coin.PlainCoin
	fee         uint64
	tokenID     *common.Hash // default is nil -> use for prv coin
	metaData    metadata.Metadata
	info        []byte // 512 bytes
	kvArgs      map[string]interface{}
}

// NewTxConvertVer1ToVer2InitParams creates a new TxConvertVer1ToVer2InitParams from the given parameters.
func NewTxConvertVer1ToVer2InitParams(senderSK *key.PrivateKey,
	paymentInfo []*coin.PaymentInfo,
	inputCoins []coin.PlainCoin,
	fee uint64,
	tokenID *common.Hash, // default is nil -> use for prv coin
	metaData metadata.Metadata,
	info []byte,
	kvArgs map[string]interface{}) *TxConvertVer1ToVer2InitParams {
	// make sure info is not nil ; zero value for it is []byte{}

	if info == nil {
		info = []byte{}
	}

	return &TxConvertVer1ToVer2InitParams{
		tokenID:     tokenID,
		inputCoins:  inputCoins,
		fee:         fee,
		metaData:    metaData,
		paymentInfo: paymentInfo,
		senderSK:    senderSK,
		info:        info,
		kvArgs:      kvArgs,
	}
}

// InitConversion creates a conversion transaction that converts PRV UTXOs v1 to v2. A conversion transaction is
// a special PRV transaction of version 2. It is non-private, meaning that all details of the transaction are publicly visible.
//   - InputCoins: PlainCoin V1
//   - OutputCoins: CoinV2
//   - Signature: Schnorr signature with no privacy.
func InitConversion(tx *Tx, params *TxConvertVer1ToVer2InitParams) error {
	// validate again
	if err := validateTxConvertVer1ToVer2Params(params); err != nil {
		return err
	}
	if err := initializeTxConversion(tx, params); err != nil {
		return err
	}
	if err := proveConversion(tx, params); err != nil {
		return err
	}
	txSize := tx.GetTxActualSize()
	if txSize > common.MaxTxSize {
		return utils.NewTransactionErr(utils.ExceedSizeTx, nil, strconv.Itoa(int(txSize)))
	}
	return nil
}

func validateTxConvertVer1ToVer2Params(params *TxConvertVer1ToVer2InitParams) error {
	if len(params.inputCoins) > 255 {
		return utils.NewTransactionErr(utils.InputCoinIsVeryLargeError, nil, strconv.Itoa(len(params.inputCoins)))
	}
	if len(params.paymentInfo) > 254 {
		return utils.NewTransactionErr(utils.PaymentInfoIsVeryLargeError, nil, strconv.Itoa(len(params.paymentInfo)))
	}

	sumInput, sumOutput := uint64(0), uint64(0)
	for _, c := range params.inputCoins {
		if c.GetVersion() != 1 {
			err := fmt.Errorf("TxConversion should only have inputCoins version 1")
			return utils.NewTransactionErr(utils.InvalidInputCoinVersionErr, err)
		}

		//Verify if input coins have been concealed
		if c.GetRandomness() == nil || c.GetSNDerivator() == nil || c.GetPublicKey() == nil || c.GetCommitment() == nil {
			err := fmt.Errorf("input coins should not be concealed")
			return utils.NewTransactionErr(utils.InvalidInputCoinVersionErr, err)
		}
		sumInput += c.GetValue()
	}
	for _, c := range params.paymentInfo {
		sumOutput += c.Amount
	}
	if sumInput != sumOutput+params.fee {
		err := fmt.Errorf("TxConversion's sum input coin and output coin (with fee) is not the same")
		return utils.NewTransactionErr(utils.SumInputCoinsAndOutputCoinsError, err)
	}

	if params.tokenID == nil {
		// using default PRV
		params.tokenID = &common.Hash{}
		if err := params.tokenID.SetBytes(common.PRVCoinID[:]); err != nil {
			return utils.NewTransactionErr(utils.TokenIDInvalidError, err, params.tokenID.String())
		}
	}
	return nil
}

func initializeTxConversion(tx *Tx, params *TxConvertVer1ToVer2InitParams) error {
	var err error
	senderKeySet := key.KeySet{}
	if err := senderKeySet.InitFromPrivateKey(params.senderSK); err != nil {
		return utils.NewTransactionErr(utils.PrivateKeySenderInvalidError, err)
	}

	// Tx: initialize some values
	tx.Fee = params.fee
	tx.Version = utils.TxConversionVersion12Number
	tx.Type = common.TxConversionType
	tx.Metadata = params.metaData
	tx.PubKeyLastByteSender = common.GetShardIDFromLastByte(senderKeySet.PaymentAddress.Pk[len(senderKeySet.PaymentAddress.Pk)-1])

	if tx.LockTime == 0 {
		tx.LockTime = time.Now().Unix()
	}
	if tx.Info, err = tx_generic.GetTxInfo(params.info); err != nil {
		return err
	}
	return nil
}

func createOutputCoins(paymentInfos []*coin.PaymentInfo, tokenID *common.Hash) ([]*coin.CoinV2, error) {
	var err error
	isPRV := (tokenID == nil) || (*tokenID == common.PRVCoinID)
	c := make([]*coin.CoinV2, len(paymentInfos))

	for i := 0; i < len(paymentInfos); i += 1 {
		if isPRV {
			c[i], _, err = coin.NewCoinFromPaymentInfo(coin.NewTransferCoinParams(paymentInfos[i]))
			if err != nil {
				log.Printf("TxConversion cannot create new coin unique OTA, got error %v\n", err)
				return nil, err
			}
		} else {
			createdCACoin, _, err := createUniqueOTACoinCA(coin.NewTransferCoinParams(paymentInfos[i]), tokenID)
			if err != nil {
				log.Printf("TxConversion cannot create new CA coin - %v\n", err)
				return nil, err
			}
			err = createdCACoin.SetPlainTokenID(tokenID)
			if err != nil {
				return nil, err
			}
			c[i] = createdCACoin
		}
	}
	return c, nil
}

func proveConversion(tx *Tx, params *TxConvertVer1ToVer2InitParams) error {
	inputCoins := params.inputCoins
	outputCoins, err := createOutputCoins(params.paymentInfo, params.tokenID)
	if err != nil {
		log.Printf("TxConversion cannot get output coins from payment info got error %v\n", err)
		return err
	}
	lenInputs := len(inputCoins)
	snWitness := make([]*serialnumbernoprivacy.SNNoPrivacyWitness, lenInputs)
	for i := 0; i < len(inputCoins); i++ {
		/***** Build witness for proving that serial number is derived from the committed derivator *****/
		snWitness[i] = new(serialnumbernoprivacy.SNNoPrivacyWitness)
		snWitness[i].Set(inputCoins[i].GetKeyImage(), inputCoins[i].GetPublicKey(),
			inputCoins[i].GetSNDerivator(), new(crypto.Scalar).FromBytesS(*params.senderSK))
	}
	tx.Proof, err = conversion.ProveConversion(inputCoins, outputCoins, snWitness)
	if err != nil {
		log.Printf("error in privacy_v2.Prove, error %v\n", err)
		return err
	}

	// sign tx
	if tx.Sig, tx.SigPubKey, err = tx_generic.SignNoPrivacy(params.senderSK, tx.Hash()[:]); err != nil {
		fmt.Println("error signNoPrivacy", err)
		return utils.NewTransactionErr(utils.SignTxError, err)
	}
	return nil
}

// CustomTokenConversionParams describes the parameters needed to create a token conversion transaction.
type CustomTokenConversionParams struct {
	tokenID       *common.Hash
	tokenInputs   []coin.PlainCoin
	tokenPayments []*coin.PaymentInfo
}

// TxTokenConvertVer1ToVer2InitParams consists of parameters used to create a new token conversion transaction.
type TxTokenConvertVer1ToVer2InitParams struct {
	senderSK    *key.PrivateKey
	feeInputs   []coin.PlainCoin
	feePayments []*coin.PaymentInfo
	fee         uint64
	tokenParams *CustomTokenConversionParams
	metaData    metadata.Metadata
	info        []byte // 512 bytes
	kvArgs      map[string]interface{}
}

// NewTxTokenConvertVer1ToVer2InitParams creates a new TxTokenConvertVer1ToVer2InitParams from the given parameters.
func NewTxTokenConvertVer1ToVer2InitParams(senderSK *key.PrivateKey,
	feeInputs []coin.PlainCoin,
	feePayments []*coin.PaymentInfo,
	tokenInputs []coin.PlainCoin,
	tokenPayments []*coin.PaymentInfo,
	fee uint64,
	tokenID *common.Hash, // tokenID of the conversion coin
	metaData metadata.Metadata,
	info []byte,
	kvArgs map[string]interface{}) *TxTokenConvertVer1ToVer2InitParams {

	if info == nil {
		info = []byte{}
	}

	tokenParams := &CustomTokenConversionParams{
		tokenID:       tokenID,
		tokenPayments: tokenPayments,
		tokenInputs:   tokenInputs,
	}
	return &TxTokenConvertVer1ToVer2InitParams{
		feeInputs:   feeInputs,
		fee:         fee,
		tokenParams: tokenParams,
		metaData:    metaData,
		feePayments: feePayments,
		senderSK:    senderSK,
		info:        info,
		kvArgs:      kvArgs,
	}
}

// InitTokenConversion creates a token conversion transaction that converts token UTXOs v1 to v2. A token conversion
// transaction is a special token transaction of version 2. It pays the transaction fee in PRV and it is required that
// the account has enough PRV v2 to pay the fee. This transaction is non-private, meaning that all details of the
// transaction are publicly visible.
//   - TxBase: A PRV transaction V2
//   - TxNormal:
//   - InputCoins: PlainCoin V1
//   - OutputCoins: CoinV2
//   - Signature: Schnorr signature with no privacy.
func InitTokenConversion(txToken *TxToken, params *TxTokenConvertVer1ToVer2InitParams) error {
	if err := validateTxTokenConvertVer1ToVer2Params(params); err != nil {
		return err
	}

	txPrivacyParams := tx_generic.NewTxPrivacyInitParams(
		params.senderSK, params.feePayments, params.feeInputs, params.fee,
		false, nil, params.metaData, params.info, params.kvArgs)
	if err := tx_generic.ValidateTxParams(txPrivacyParams); err != nil {
		return err
	}
	// Init tx and params (tx and params will be changed)
	tx := new(Tx)
	if err := tx.InitializeTxAndParams(txPrivacyParams); err != nil {
		return err
	}

	// Init PRV Fee
	ins, outs, err := txToken.initPRVFeeConversion(tx, txPrivacyParams)
	if err != nil {
		log.Printf("Cannot init token ver2: err %v", err)
		return err
	}
	txn := makeTxToken(tx, nil, nil, nil)

	// Init Token
	if err := txToken.initTokenConversion(txn, params); err != nil {
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
	txSize := txToken.GetTxActualSize()
	if txSize > common.MaxTxSize {
		return utils.NewTransactionErr(utils.ExceedSizeTx, nil, strconv.Itoa(int(txSize)))
	}
	return nil
}

func (txToken *TxToken) initTokenConversion(txNormal *Tx, params *TxTokenConvertVer1ToVer2InitParams) error {
	txToken.TokenData.Type = utils.CustomTokenTransfer
	txToken.TokenData.PropertyName = ""
	txToken.TokenData.PropertySymbol = ""
	txToken.TokenData.Mintable = false
	txToken.TokenData.PropertyID = *params.tokenParams.tokenID

	txConvertParams := NewTxConvertVer1ToVer2InitParams(
		params.senderSK,
		params.tokenParams.tokenPayments,
		params.tokenParams.tokenInputs,
		0,
		params.tokenParams.tokenID,
		nil,
		params.info,
		params.kvArgs)

	if err := validateTxConvertVer1ToVer2Params(txConvertParams); err != nil {
		return utils.NewTransactionErr(utils.PrivacyTokenInitTokenDataError, err)
	}
	if err := initializeTxConversion(txNormal, txConvertParams); err != nil {
		return utils.NewTransactionErr(utils.PrivacyTokenInitTokenDataError, err)
	}
	txNormal.SetType(common.TxTokenConversionType)
	if err := proveConversion(txNormal, txConvertParams); err != nil {
		return utils.NewTransactionErr(utils.PrivacyTokenInitTokenDataError, err)
	}

	err := txToken.SetTxNormal(txNormal)
	return err
}

func (txToken *TxToken) initPRVFeeConversion(feeTx *Tx, params *tx_generic.TxPrivacyInitParams) ([]coin.PlainCoin, []*coin.CoinV2, error) {
	feeTx.SetVersion(utils.TxConversionVersion12Number)
	feeTx.SetType(common.TxTokenConversionType)
	ins, outs, err := feeTx.provePRV(params)
	if err != nil {
		return nil, nil, utils.NewTransactionErr(utils.PrivacyTokenInitPRVError, err)
	}

	return ins, outs, nil
}

func validateTxTokenConvertVer1ToVer2Params(params *TxTokenConvertVer1ToVer2InitParams) error {
	if len(params.feeInputs) > 255 {
		return fmt.Errorf("FeeInput is too large, feeInputs length = " + strconv.Itoa(len(params.feeInputs)))
	}
	if len(params.feePayments) > 255 {
		return fmt.Errorf("FeePayment is too large, feePayments length = " + strconv.Itoa(len(params.feePayments)))
	}
	if len(params.tokenParams.tokenPayments) > 255 {
		return fmt.Errorf("tokenPayments is too large, tokenPayments length = " + strconv.Itoa(len(params.tokenParams.tokenPayments)))
	}
	if len(params.tokenParams.tokenInputs) > 255 {
		return fmt.Errorf("tokenInputs length = " + strconv.Itoa(len(params.tokenParams.tokenInputs)))
	}

	for _, c := range params.feeInputs {
		if c.GetVersion() != utils.TxVersion2Number {
			return fmt.Errorf("TxConversion should only have fee input coins version 2")
		}
	}
	tokenParams := params.tokenParams
	if tokenParams.tokenID == nil {
		return utils.NewTransactionErr(utils.TokenIDInvalidError, fmt.Errorf("TxTokenConversion should have its tokenID not null"))
	}
	sumInput := uint64(0)
	for _, c := range tokenParams.tokenInputs {
		sumInput += c.GetValue()
	}
	if sumInput != tokenParams.tokenPayments[0].Amount {
		return utils.NewTransactionErr(utils.SumInputCoinsAndOutputCoinsError, fmt.Errorf("sumInput and sum TokenPayment amount is not equal"))
	}
	return nil
}
