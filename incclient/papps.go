package incclient

import (
	"fmt"
	"strings"

	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataBridge "github.com/incognitochain/go-incognito-sdk-v2/metadata/bridge"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// CreateBurningpUnifiedRequestTransaction creates an EVM pUnified burning transaction for exiting the Incognito network.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
//
// An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
//   - rpc.ETHNetworkID: the Ethereum network
//   - rpc.BSCNetworkID: the Binance Smart Chain network
//   - rpc.PLGNetworkID: the Polygon network
//   - rpc.FTMNetworkID: the Fantom network
//
// If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (client *IncClient) CreateBurnForCallRequestTransaction(
	privateKey, tokenIDStr string,
	data metadataBridge.BurnForCallRequestData,
	transferTokenReceivers []string, transferTokenAmounts []uint64,
	// it's used to create tx with specific input coins
	tokenInCoins []coin.PlainCoin,
	tokenIndices []uint64,
	prvInCoins []coin.PlainCoin,
	prvIndices []uint64,
) ([]byte, string, error) {
	if tokenIDStr == common.PRVIDStr {
		return nil, "", fmt.Errorf("cannot burn PRV in a burning request transaction")
	}

	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot deserialize the sender private key")
	}
	burnerAddress := senderWallet.KeySet.PaymentAddress
	if common.AddressVersion == 0 {
		burnerAddress.OTAPublic = nil
	}

	if strings.Contains(data.WithdrawAddress, "0x") {
		data.WithdrawAddress = data.WithdrawAddress[2:]
	}

	if len(transferTokenAmounts) != len(transferTokenReceivers) {
		return nil, "", fmt.Errorf("Invalid params transferTokenReceivers and transferTokenAmounts")
	}

	mdType := metadataCommon.BurnForCallRequestMeta
	md := &metadataBridge.BurnForCallRequest{
		BurnTokenID: *tokenID,
		Data:        []metadataBridge.BurnForCallRequestData{data},
		MetadataBase: metadataCommon.MetadataBase{
			Type: mdType,
		},
	}
	burnedAmount, err := md.TotalBurningAmount()
	if err != nil {
		return nil, "", fmt.Errorf("cannot get total burning amount: %v", err)
	}

	transferTokenReceivers = append([]string{common.BurningAddress2}, transferTokenReceivers...)
	transferTokenAmounts = append([]uint64{burnedAmount}, transferTokenAmounts...)

	tokenParam := NewTxTokenParam(tokenIDStr, 1, transferTokenReceivers, transferTokenAmounts, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, tokenParam, md, nil)

	if len(tokenInCoins) > 0 && len(prvInCoins) > 0 {
		return client.CreateRawTokenTransactionWithInputCoins(txParam, tokenInCoins, tokenIndices, prvInCoins, prvIndices)
	} else {
		return client.CreateRawTokenTransaction(txParam, -1)
	}
}

// CreateAndSendBurnForCallRequestTransaction creates an EVM burning transaction for exiting the Incognito network, and submits it to the network.
//
// It returns the transaction's hash, and an error (if any).
//
// If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (client *IncClient) CreateAndSendBurnForCallRequestTransaction(
	privateKey, tokenIDStr string,
	data metadataBridge.BurnForCallRequestData,
	transferTokenReceivers []string, transferTokenAmounts []uint64,
	// it's used to create tx with specific input coins
	tokenInCoins []coin.PlainCoin,
	tokenIndices []uint64,
	prvInCoins []coin.PlainCoin,
	prvIndices []uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreateBurnForCallRequestTransaction(privateKey, tokenIDStr,
		data,
		transferTokenReceivers, transferTokenAmounts,
		tokenInCoins,
		tokenIndices,
		prvInCoins,
		prvIndices)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
