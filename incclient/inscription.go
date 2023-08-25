package incclient

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
	metadataInsc "github.com/incognitochain/go-incognito-sdk-v2/metadata/inscription"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// CreateInscribeRequestTx creates a inscribing transaction.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateInscribeRequestTx(
	privateKey, data string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) ([]byte, string, error) {
	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	address, _ := senderWallet.GetPaymentAddress()
	fmt.Println("address: ", address)

	otaReceiver := coin.OTAReceiver{}
	paymentInfo := &key.PaymentInfo{PaymentAddress: senderWallet.KeySet.PaymentAddress, Message: []byte{}}
	err = otaReceiver.FromCoinParams(coin.NewMintCoinParams(paymentInfo))
	if err != nil {
		return nil, "", err
	}

	inscribeMd, err := metadataInsc.NewInscribeRequest(
		data,
		otaReceiver,
		metadataCommon.InscribeRequestMeta,
	)
	if err != nil {
		return nil, "", err
	}

	// const txSize = await estimateTxSize(numInputs, prvPayments.length, burnReqMetadata, null);
	//     let estFee = txSize * INSC_MIN_FEE_PER_KB;
	//     if (estFee < INSC_MIN_FEE_PER_TX) {
	//         estFee = INSC_MIN_FEE_PER_TX;
	//     }

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{100}, 0, nil, inscribeMd, nil)
	if len(inputCoins) > 0 {
		return client.CreateRawTransactionWithInputCoins(txParam, inputCoins, coinIndices)
	}
	return client.CreateRawTransaction(txParam, 2)
}

// CreateAndSendInscribeRequestTx creates a inscribing transaction,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendInscribeRequestTx(
	privateKey, data string, inputCoins []coin.PlainCoin, coinIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreateInscribeRequestTx(privateKey, data, inputCoins, coinIndices)
	if err != nil {
		return "", err
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}
