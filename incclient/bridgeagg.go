package incclient

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	rCommon "github.com/ethereum/go-ethereum/common"
// 	"github.com/incognitochain/go-incognito-sdk-v2/coin"
// 	"github.com/incognitochain/go-incognito-sdk-v2/key"
// 	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
// 	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"

// 	"github.com/incognitochain/go-incognito-sdk-v2/common"
// 	metadataBridge "github.com/incognitochain/go-incognito-sdk-v2/metadata/bridge"
// 	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
// 	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
// )

// // GetBridgeAggState returns bridge agg vaults state in the network.
// // An additional parameter `beaconHeightParam`:
// // If set empty, beaconHeightParam defaults to 0 (the latest state)
// func (client *IncClient) GetBridgeAggState(beaconHeightParam ...uint64) (*jsonresult.BridgeAggState, error) {
// 	beaconHeight := uint64(0)
// 	if len(beaconHeightParam) > 0 {
// 		beaconHeight = beaconHeightParam[0]
// 	}
// 	responseInBytes, err := client.rpcServer.GetBridgeAggState(beaconHeight)
// 	if err != nil {
// 		return nil, err
// 	}

// 	res := jsonresult.BridgeAggState{}
// 	err = rpchandler.ParseResponse(responseInBytes, &res)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &res, nil
// }

// // CreateAndSendBurningPUnifiedRequestTransaction creates an EVM pUnified burning transaction for exiting the Incognito network, and submits it to the network.
// //
// // It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
// //
// // prvTransferReceivers, prvTransferAmts contain PRV payment infos (dont include burn info)
// // tokenTransferReceivers, tokenTransferAmts contain token payment infos (dont include burn info)
// //
// // An additional parameter `isDepositToSCParam`:
// //	- false: exit from Incognito chain
// //	- true: interact with dApps and re-shield to Incognito
// // If set empty, isDepositToSCParam defaults to FALSE.
// // NOTE that only the first value of evmNetworkID is used.
// func (client *IncClient) CreateAndSendBurningPUnifiedRequestTransaction(
// 	privateKey, unifiedTokenIDStr string,
// 	unshieldDatas []metadataBridge.UnshieldRequestData,
// 	prvTransferReceivers []string, prvTransferAmts []uint64,
// 	tokenTransferReceivers []string, tokenTransferAmts []uint64,
// 	isDepositToSCParam ...bool,
// ) (string, error) {
// 	encodedTx, txHash, err := client.CreateBurningPUnifiedRequestTransaction(
// 		privateKey,
// 		unifiedTokenIDStr,
// 		unshieldDatas,
// 		prvTransferReceivers, prvTransferAmts,
// 		tokenTransferReceivers, tokenTransferAmts,
// 		isDepositToSCParam...)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = client.SendRawTokenTx(encodedTx)
// 	if err != nil {
// 		return "", err
// 	}

// 	return txHash, nil
// }

// // CreateBurningPUnifiedRequestTransaction creates an EVM pUnified burning transaction for exiting the Incognito network.
// //
// // It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
// //
// // prvTransferReceivers, prvTransferAmts contain PRV payment infos (dont include burn info)
// // tokenTransferReceivers, tokenTransferAmts contain token payment infos (dont include burn info)
// //
// // An additional parameter `isDepositToSCParam`:
// //	- false: exit from Incognito chain
// //	- true: interact with dApps and re-shield to Incognito
// // If set empty, isDepositToSCParam defaults to FALSE.
// // NOTE that only the first value of evmNetworkID is used.
// func (client *IncClient) CreateBurningPUnifiedRequestTransaction(
// 	privateKey, unifiedTokenIDStr string,
// 	unshieldDatas []metadataBridge.UnshieldRequestData,
// 	prvTransferReceivers []string, prvTransferAmts []uint64,
// 	tokenTransferReceivers []string, tokenTransferAmts []uint64,
// 	isDepositToSCParam ...bool,
// ) ([]byte, string, error) {
// 	if unifiedTokenIDStr == common.PRVIDStr {
// 		return nil, "", fmt.Errorf("cannot burn PRV in a burning request transaction")
// 	}

// 	unifiedTokenID, err := new(common.Hash).NewHashFromStr(unifiedTokenIDStr)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("cannot deserialize the sender private key")
// 	}

// 	isDepositToSC := false
// 	if len(isDepositToSCParam) > 0 {
// 		isDepositToSC = isDepositToSCParam[0]
// 	}

// 	totalBurnAmt := uint64(0)
// 	for i, data := range unshieldDatas {
// 		totalBurnAmt += data.BurningAmount
// 		if totalBurnAmt < data.BurningAmount {
// 			return nil, "", fmt.Errorf("Out of range total burning amount")
// 		}
// 		if strings.Contains(data.RemoteAddress, "0x") {
// 			unshieldDatas[i].RemoteAddress = data.RemoteAddress[2:]
// 		}
// 	}

// 	otaReceivers, err := GenerateOTAReceivers([]common.Hash{*unifiedTokenID}, senderWallet.KeySet.PaymentAddress)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	md := metadataBridge.NewUnshieldRequestWithValue(*unifiedTokenID, unshieldDatas, otaReceivers[*unifiedTokenID], isDepositToSC)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("cannot init unshield unified request with unifiedTokenID %v: %v", unifiedTokenIDStr, err)
// 	}

// 	tokenParam := NewTxTokenParam(unifiedTokenIDStr, 1, append(tokenTransferReceivers, common.BurningAddress2), append(tokenTransferAmts, totalBurnAmt), false, 0, nil)
// 	txParam := NewTxParam(privateKey, prvTransferReceivers, prvTransferAmts, DefaultPRVFee, tokenParam, md, nil)

// 	return client.CreateRawTokenTransaction(txParam, -1)
// }

// // CreateIssuingpUnifiedRequestTransaction creates an EVM pUnified shielding trading transaction. By EVM, it means either ETH or BSC.
// //
// // It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
// //
// // An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
// //	- rpc.ETHNetworkID: the Ethereum network
// //	- rpc.BSCNetworkID: the Binance Smart Chain network
// //	- rpc.PLGNetworkID: the Polygon network
// //	- rpc.FTMNetworkID: the Fantom network
// // If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
// func (client *IncClient) CreateIssuingpUnifiedRequestTransaction(privateKey, tokenIDStr string, pUnifiedTokenIDStr string, proof EVMDepositProof, evmNetworkID ...int) ([]byte, string, error) {
// 	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	pUnifiedTokenID, err := new(common.Hash).NewHashFromStr(pUnifiedTokenIDStr)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	networkID := rpc.ETHNetworkID
// 	if len(evmNetworkID) > 0 {
// 		networkID = evmNetworkID[0]
// 	}
// 	// if _, ok := rpc.EVMIssuingMetadata[networkID]; !ok {
// 	// 	return nil, "", fmt.Errorf("networkID %v not found", networkID)
// 	// }

// 	type EVMProof struct {
// 		BlockHash rCommon.Hash `json:"BlockHash"`
// 		TxIndex   uint         `json:"TxIndex"`
// 		Proof     []string     `json:"Proof"`
// 	}

// 	proofData := EVMProof{
// 		BlockHash: proof.blockHash,
// 		TxIndex:   proof.txIdx,
// 		Proof:     proof.nodeList,
// 	}
// 	proofBytes, err := json.Marshal(proofData)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to marshal proof")
// 	}

// 	var issuingETHRequestMeta *metadataBridge.ShieldRequest
// 	shieldRequestData := metadataBridge.ShieldRequestData{
// 		IncTokenID: *tokenID,
// 		NetworkID:  uint8(networkID),
// 		Proof:      proofBytes,
// 	}
// 	issuingETHRequestMeta = metadataBridge.NewShieldRequestWithValue([]metadataBridge.ShieldRequestData{shieldRequestData}, *pUnifiedTokenID)

// 	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, issuingETHRequestMeta, nil)
// 	return client.CreateRawTransaction(txParam, -1)
// }

// // CreateAndSendIssuingpUnifiedRequestTransaction creates an EVM pUnified shielding transaction, and submits it to the Incognito network.
// //
// // It returns the transaction's hash, and an error (if any).
// //
// // An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
// //	- rpc.ETHNetworkID: the Ethereum network
// //	- rpc.BSCNetworkID: the Binance Smart Chain network
// //	- rpc.PLGNetworkID: the Polygon network
// //	- rpc.FTMNetworkID: the Fantom network
// // If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
// func (client *IncClient) CreateAndSendIssuingpUnifiedRequestTransaction(privateKey, tokenIDStr string, pUnifiedTokenIDStr string, proof EVMDepositProof, evmNetworkID ...int) (string, error) {
// 	encodedTx, txHash, err := client.CreateIssuingpUnifiedRequestTransaction(privateKey, tokenIDStr, pUnifiedTokenIDStr, proof, evmNetworkID...)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = client.SendRawTx(encodedTx)
// 	if err != nil {
// 		return "", err
// 	}

// 	return txHash, nil
// }

// func (client *IncClient) CreateBridgeAggConvertTokenToUnifiedTokenRequestTransaction(privateKey, tokenIDStr, pUnifiedTokenIDStr string, burnedAmount uint64) ([]byte, string, error) {
// 	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	pUnifiedTokenID, err := new(common.Hash).NewHashFromStr(pUnifiedTokenIDStr)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	// mdType := metadata.BridgeAggConvertTokenToUnifiedTokenRequestMeta

// 	var md *metadataBridge.ConvertTokenToUnifiedTokenRequest

// 	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("cannot deserialize the sender private key")
// 	}
// 	burnerAddress := senderWallet.KeySet.PaymentAddress
// 	if common.AddressVersion == 0 {
// 		burnerAddress.OTAPublic = nil
// 	}

// 	otaReceiver := coin.OTAReceiver{}
// 	paymentInfo := &key.PaymentInfo{PaymentAddress: senderWallet.KeySet.PaymentAddress, Message: []byte{}}
// 	err = otaReceiver.FromCoinParams(coin.NewMintCoinParams(paymentInfo))
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	md = metadataBridge.NewConvertTokenToUnifiedTokenRequestWithValue(*tokenID, *pUnifiedTokenID, burnedAmount, otaReceiver)

// 	tokenParam := NewTxTokenParam(tokenIDStr, 1, []string{common.BurningAddress2}, []uint64{burnedAmount}, false, 0, nil)
// 	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, tokenParam, md, nil)

// 	return client.CreateRawTokenTransaction(txParam, -1)
// }

// func (client *IncClient) CreateAndSendBridgeAggConvertTokenToUnifiedTokenRequestTransaction(privateKey, tokenIDStr, pUnifiedTokenIDStr string, burnedAmount uint64, evmNetworkID ...int) (string, error) {
// 	encodedTx, txHash, err := client.CreateBridgeAggConvertTokenToUnifiedTokenRequestTransaction(privateKey, tokenIDStr, pUnifiedTokenIDStr, burnedAmount)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = client.SendRawTx(encodedTx)
// 	if err != nil {
// 		return "", err
// 	}

// 	return txHash, nil
// }

// func (client *IncClient) CreateBridgeAggModifyParamTransaction(privateKey string, percentFeeWithDec uint64) ([]byte, string, error) {

// 	md := metadataBridge.NewModifyBridgeAggParamReqWithValue(percentFeeWithDec)

// 	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, md, nil)

// 	return client.CreateRawTokenTransaction(txParam, -1)
// }

// func (client *IncClient) CreateAndSendBridgeAggModifyParamTransaction(privateKey string, percentFeeWithDec uint64) (string, error) {
// 	encodedTx, txHash, err := client.CreateBridgeAggModifyParamTransaction(privateKey, percentFeeWithDec)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = client.SendRawTx(encodedTx)
// 	if err != nil {
// 		return "", err
// 	}

// 	return txHash, nil
// }

// func (client *IncClient) CheckUnifiedShieldStatus(txHash string) (*ShieldStatus, error) {
// 	responseInBytes, err := client.rpcServer.CheckShieldUnifiedStatus(txHash)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var status ShieldStatus
// 	err = rpchandler.ParseResponse(responseInBytes, &status)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &status, err
// }

// func (client *IncClient) CheckUnifiedUnshieldStatus(txHash string) (*UnshieldStatus, error) {
// 	responseInBytes, err := client.rpcServer.CheckUnshieldUnifiedStatus(txHash)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var status UnshieldStatus
// 	err = rpchandler.ParseResponse(responseInBytes, &status)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &status, err
// }

// type ShieldStatusData struct {
// 	Amount uint64 `json:"Amount"`
// 	Reward uint64 `json:"Reward"`
// }

// type ShieldStatus struct {
// 	Status    byte               `json:"Status"`
// 	Data      []ShieldStatusData `json:"Data,omitempty"`
// 	ErrorCode int                `json:"ErrorCode,omitempty"`
// }

// type UnshieldStatusData struct {
// 	ReceivedAmount uint64 `json:"ReceivedAmount"`
// 	Fee            uint64 `json:"Fee"`
// }

// type UnshieldStatus struct {
// 	Status    byte                 `json:"Status"`
// 	Data      []UnshieldStatusData `json:"Data,omitempty"`
// 	ErrorCode int                  `json:"ErrorCode,omitempty"`
// }
