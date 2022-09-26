package incclient

import (
	"encoding/json"
	"fmt"
	"strings"

	rCommon "github.com/ethereum/go-ethereum/common"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	metadataBridge "github.com/incognitochain/go-incognito-sdk-v2/metadata/bridge"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// EVMDepositProof represents a proof for depositing tokens to the smart contracts.
type EVMDepositProof struct {
	blockNumber uint
	blockHash   ethCommon.Hash
	txIdx       uint
	nodeList    []string
}

// TxIdx returns the transaction index of an EVMDepositProof.
func (E EVMDepositProof) TxIdx() uint {
	return E.txIdx
}

// BlockNumber returns the block number of an EVMDepositProof.
func (E EVMDepositProof) BlockNumber() uint {
	return E.blockNumber
}

// BlockHash returns the block hash of an EVMDepositProof.
func (E EVMDepositProof) BlockHash() ethCommon.Hash {
	return E.blockHash
}

// NodeList returns the node list of an EVMDepositProof.
func (E EVMDepositProof) NodeList() []string {
	return E.nodeList
}

// NewETHDepositProof creates a new EVMDepositProof with the given parameters.
func NewETHDepositProof(blockNumber uint, blockHash ethCommon.Hash, txIdx uint, nodeList []string) *EVMDepositProof {
	proof := EVMDepositProof{
		blockNumber: blockNumber,
		blockHash:   blockHash,
		txIdx:       txIdx,
		nodeList:    nodeList,
	}

	return &proof
}

// CreateIssuingRequestTransaction creates a centralized shielding transaction.
// This function should only be called along with the privateKey of the centralized account.
func (client *IncClient) CreateIssuingRequestTransaction(privateKey, receiver, tokenIDStr, tokenName string, depositAmount uint64) ([]byte, string, error) {
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	w, err := wallet.Base58CheckDeserialize(receiver)
	if err != nil {
		return nil, "", err
	}
	addr := w.KeySet.PaymentAddress
	if _, err := AssertPaymentAddressAndTxVersion(addr, 2); err != nil {
		return nil, "", fmt.Errorf("invalid receiver address")
	}

	var issuingRequestMeta *metadata.IssuingRequest
	issuingRequestMeta, err = metadata.NewIssuingRequest(addr, depositAmount, *tokenID, tokenName, metadata.IssuingRequestMeta)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init issue request for %v, tokenID %v: %v", receiver, tokenIDStr, err)
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, issuingRequestMeta, nil)
	return client.CreateRawTransaction(txParam, 2)
}

// CreateAndSendIssuingRequestTransaction creates a centralized shielding transaction, and submits it to the Incognito network.
func (client *IncClient) CreateAndSendIssuingRequestTransaction(privateKey,
	receiver,
	tokenIDStr,
	tokenName string,
	depositAmount uint64,
) (string, error) {
	encodedTx, txHash, err := client.CreateIssuingRequestTransaction(privateKey, receiver, tokenIDStr, tokenName, depositAmount)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateIssuingEVMRequestTransaction creates an EVM shielding trading transaction. By EVM, it means either ETH or BSC.
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
func (client *IncClient) CreateIssuingEVMRequestTransaction(privateKey, tokenIDStr string, proof EVMDepositProof, evmNetworkID ...int) ([]byte, string, error) {
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	networkID := rpc.ETHNetworkID
	if len(evmNetworkID) > 0 {
		networkID = evmNetworkID[0]
	}
	if _, ok := rpc.EVMIssuingMetadata[networkID]; !ok {
		return nil, "", fmt.Errorf("networkID %v not found", networkID)
	}
	mdType := rpc.EVMIssuingMetadata[networkID]

	var issuingETHRequestMeta *metadata.IssuingEVMRequest
	issuingETHRequestMeta, err = metadata.NewIssuingEVMRequest(proof.blockHash, proof.txIdx, proof.nodeList, *tokenID, mdType)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init issue eth request for %v, tokenID %v: %v", proof, tokenIDStr, err)
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, issuingETHRequestMeta, nil)
	return client.CreateRawTransaction(txParam, -1)
}

// CreateAndSendIssuingEVMRequestTransaction creates an EVM shielding transaction, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
//
// An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
//   - rpc.ETHNetworkID: the Ethereum network
//   - rpc.BSCNetworkID: the Binance Smart Chain network
//   - rpc.PLGNetworkID: the Polygon network
//   - rpc.FTMNetworkID: the Fantom network
//
// If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (client *IncClient) CreateAndSendIssuingEVMRequestTransaction(privateKey, tokenIDStr string, proof EVMDepositProof, evmNetworkID ...int) (string, error) {
	encodedTx, txHash, err := client.CreateIssuingEVMRequestTransaction(privateKey, tokenIDStr, proof, evmNetworkID...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateBurningRequestTransaction creates an EVM burning transaction for exiting the Incognito network.
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
func (client *IncClient) CreateBurningRequestTransaction(privateKey, remoteAddress, tokenIDStr string, burnedAmount uint64, evmNetworkID ...int) ([]byte, string, error) {
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

	if strings.Contains(remoteAddress, "0x") {
		remoteAddress = remoteAddress[2:]
	}

	networkID := rpc.ETHNetworkID
	if len(evmNetworkID) > 0 {
		networkID = evmNetworkID[0]
	}
	if _, ok := rpc.EVMBurningMetadata[networkID]; !ok {
		return nil, "", fmt.Errorf("networkID %v not found", networkID)
	}
	mdType := rpc.EVMBurningMetadata[networkID]

	var md *metadata.BurningRequest
	md, err = metadata.NewBurningRequest(burnerAddress, burnedAmount, *tokenID, tokenIDStr, remoteAddress, mdType)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init burning request with tokenID %v, burnedAmount %v, remoteAddress %v: %v", tokenIDStr, burnedAmount, remoteAddress, err)
	}

	tokenParam := NewTxTokenParam(tokenIDStr, 1, []string{common.BurningAddress2}, []uint64{burnedAmount}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, tokenParam, md, nil)

	return client.CreateRawTokenTransaction(txParam, -1)
}

// CreateAndSendBurningRequestTransaction creates an EVM burning transaction for exiting the Incognito network, and submits it to the network.
//
// It returns the transaction's hash, and an error (if any).
//
// An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
//   - rpc.ETHNetworkID: the Ethereum network
//   - rpc.BSCNetworkID: the Binance Smart Chain network
//   - rpc.PLGNetworkID: the Polygon network
//   - rpc.FTMNetworkID: the Fantom network
//
// If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (client *IncClient) CreateAndSendBurningRequestTransaction(privateKey, remoteAddress, tokenIDStr string, burnedAmount uint64, evmNetworkID ...int) (string, error) {
	encodedTx, txHash, err := client.CreateBurningRequestTransaction(privateKey, remoteAddress, tokenIDStr, burnedAmount, evmNetworkID...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateAndSendBurningpUnifiedRequestTransaction creates an EVM pUnified burning transaction for exiting the Incognito network, and submits it to the network.
//
// It returns the transaction's hash, and an error (if any).
//
// An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
//   - rpc.ETHNetworkID: the Ethereum network
//   - rpc.BSCNetworkID: the Binance Smart Chain network
//   - rpc.PLGNetworkID: the Polygon network
//   - rpc.FTMNetworkID: the Fantom network
//
// If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (client *IncClient) CreateAndSendBurningpUnifiedRequestTransaction(privateKey, remoteAddress, tokenIDStr string, burnedAmount uint64, evmNetworkID ...int) (string, error) {
	encodedTx, txHash, err := client.CreateBurningpUnifiedRequestTransaction(privateKey, remoteAddress, tokenIDStr, burnedAmount, evmNetworkID...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

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
func (client *IncClient) CreateBurningpUnifiedRequestTransaction(privateKey, remoteAddress, tokenIDStr string, burnedAmount uint64, evmNetworkID ...int) ([]byte, string, error) {
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

	if strings.Contains(remoteAddress, "0x") {
		remoteAddress = remoteAddress[2:]
	}

	// networkID := rpc.ETHNetworkID
	// if len(evmNetworkID) > 0 {
	// 	networkID = evmNetworkID[0]
	// }
	// if _, ok := rpc.EVMBurningMetadata[networkID]; !ok {
	// 	return nil, "", fmt.Errorf("networkID %v not found", networkID)
	// }
	// mdType := rpc.EVMBurningMetadata[networkID]
	mdType := metadata.BurningUnifiedTokenRequestMeta
	var md *metadataBridge.BurningRequest
	md, err = metadataBridge.NewBurningRequest(burnerAddress, burnedAmount, *tokenID, tokenIDStr, remoteAddress, mdType)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init burning request with tokenID %v, burnedAmount %v, remoteAddress %v: %v", tokenIDStr, burnedAmount, remoteAddress, err)
	}

	tokenParam := NewTxTokenParam(tokenIDStr, 1, []string{common.BurningAddress2}, []uint64{burnedAmount}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, tokenParam, md, nil)

	return client.CreateRawTokenTransaction(txParam, -1)
}

// CreateIssuingpUnifiedRequestTransaction creates an EVM pUnified shielding trading transaction. By EVM, it means either ETH or BSC.
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
func (client *IncClient) CreateIssuingpUnifiedRequestTransaction(privateKey, tokenIDStr string, pUnifiedTokenIDStr string, proof EVMDepositProof, evmNetworkID ...int) ([]byte, string, error) {
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	pUnifiedTokenID, err := new(common.Hash).NewHashFromStr(pUnifiedTokenIDStr)
	if err != nil {
		return nil, "", err
	}

	networkID := rpc.ETHNetworkID
	if len(evmNetworkID) > 0 {
		networkID = evmNetworkID[0]
	}
	// if _, ok := rpc.EVMIssuingMetadata[networkID]; !ok {
	// 	return nil, "", fmt.Errorf("networkID %v not found", networkID)
	// }

	type EVMProof struct {
		BlockHash rCommon.Hash `json:"BlockHash"`
		TxIndex   uint         `json:"TxIndex"`
		Proof     []string     `json:"Proof"`
	}

	proofData := EVMProof{
		BlockHash: proof.blockHash,
		TxIndex:   proof.txIdx,
		Proof:     proof.nodeList,
	}
	proofBytes, err := json.Marshal(proofData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal proof")
	}

	var issuingETHRequestMeta *metadataBridge.ShieldRequest
	shieldRequestData := metadataBridge.ShieldRequestData{
		IncTokenID: *tokenID,
		NetworkID:  uint8(networkID),
		Proof:      proofBytes,
	}
	issuingETHRequestMeta = metadataBridge.NewShieldRequestWithValue([]metadataBridge.ShieldRequestData{shieldRequestData}, *pUnifiedTokenID)

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, issuingETHRequestMeta, nil)
	return client.CreateRawTransaction(txParam, -1)
}

// CreateAndSendIssuingpUnifiedRequestTransaction creates an EVM pUnified shielding transaction, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
//
// An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
//   - rpc.ETHNetworkID: the Ethereum network
//   - rpc.BSCNetworkID: the Binance Smart Chain network
//   - rpc.PLGNetworkID: the Polygon network
//   - rpc.FTMNetworkID: the Fantom network
//
// If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (client *IncClient) CreateAndSendIssuingpUnifiedRequestTransaction(privateKey, tokenIDStr string, pUnifiedTokenIDStr string, proof EVMDepositProof, evmNetworkID ...int) (string, error) {
	encodedTx, txHash, err := client.CreateIssuingpUnifiedRequestTransaction(privateKey, tokenIDStr, pUnifiedTokenIDStr, proof, evmNetworkID...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

func (client *IncClient) CreateBridgeAggConvertTokenToUnifiedTokenRequestTransaction(privateKey, tokenIDStr, pUnifiedTokenIDStr string, burnedAmount uint64) ([]byte, string, error) {
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	pUnifiedTokenID, err := new(common.Hash).NewHashFromStr(pUnifiedTokenIDStr)
	if err != nil {
		return nil, "", err
	}

	// mdType := metadata.BridgeAggConvertTokenToUnifiedTokenRequestMeta

	var md *metadataBridge.ConvertTokenToUnifiedTokenRequest

	senderWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("cannot deserialize the sender private key")
	}
	burnerAddress := senderWallet.KeySet.PaymentAddress
	if common.AddressVersion == 0 {
		burnerAddress.OTAPublic = nil
	}

	otaReceiver := coin.OTAReceiver{}
	paymentInfo := &coin.PaymentInfo{PaymentAddress: &senderWallet.KeySet.PaymentAddress, Message: []byte{}}
	err = otaReceiver.FromCoinParams(coin.NewMintCoinParams(paymentInfo))
	if err != nil {
		return nil, "", err
	}
	md = metadataBridge.NewConvertTokenToUnifiedTokenRequestWithValue(*tokenID, *pUnifiedTokenID, burnedAmount, otaReceiver)

	tokenParam := NewTxTokenParam(tokenIDStr, 1, []string{common.BurningAddress2}, []uint64{burnedAmount}, false, 0, nil)
	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, tokenParam, md, nil)

	return client.CreateRawTokenTransaction(txParam, -1)
}

func (client *IncClient) CreateAndSendBridgeAggConvertTokenToUnifiedTokenRequestTransaction(privateKey, tokenIDStr, pUnifiedTokenIDStr string, burnedAmount uint64, evmNetworkID ...int) (string, error) {
	encodedTx, txHash, err := client.CreateBridgeAggConvertTokenToUnifiedTokenRequestTransaction(privateKey, tokenIDStr, pUnifiedTokenIDStr, burnedAmount)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

func (client *IncClient) CreateBridgeAggModifyParamTransaction(privateKey string, percentFeeWithDec uint64) ([]byte, string, error) {

	md := metadataBridge.NewModifyBridgeAggParamReqWithValue(percentFeeWithDec)

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, md, nil)

	return client.CreateRawTokenTransaction(txParam, -1)
}

func (client *IncClient) CreateAndSendBridgeAggModifyParamTransaction(privateKey string, percentFeeWithDec uint64) (string, error) {
	encodedTx, txHash, err := client.CreateBridgeAggModifyParamTransaction(privateKey, percentFeeWithDec)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// GetBurnProof retrieves the burning proof for the Incognito network for submitting to the smart contract later.
//
// An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
//   - rpc.ETHNetworkID: the Ethereum network
//   - rpc.BSCNetworkID: the Binance Smart Chain network
//   - rpc.PLGNetworkID: the Polygon network
//   - rpc.FTMNetworkID: the Fantom network
//
// If set empty, evmNetworkID defaults to rpc.ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (client *IncClient) GetBurnProof(txHash string, evmNetworkID ...int) (*jsonresult.InstructionProof, error) {
	responseInBytes, err := client.rpcServer.GetBurnProof(txHash, evmNetworkID...)
	if err != nil {
		return nil, err
	}

	var tmp jsonresult.InstructionProof
	err = rpchandler.ParseResponse(responseInBytes, &tmp)
	if err != nil {
		return nil, err
	}

	return &tmp, nil
}

// GetBridgeTokens returns all bridge tokens in the network.
func (client *IncClient) GetBridgeTokens() ([]*BridgeTokenInfo, error) {
	responseInBytes, err := client.rpcServer.GetAllBridgeTokens()
	if err != nil {
		return nil, err
	}

	res := make([]*BridgeTokenInfo, 0)
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CheckShieldStatus returns the status of an eth-shielding request.
//   - -1: error
//   - 0: tx not found
//   - 1: tx is pending
//   - 2: tx is accepted
//   - 3: tx is rejected
func (client *IncClient) CheckShieldStatus(txHash string) (int, error) {
	responseInBytes, err := client.rpcServer.CheckShieldStatus(txHash)
	if err != nil {
		return -1, err
	}

	var status int
	err = rpchandler.ParseResponse(responseInBytes, &status)
	if err != nil {
		return -1, err
	}

	return status, err
}

// GenerateTokenID generates an Incognito tokenID for a bridge token.
func GenerateTokenID(network, tokenName string) (common.Hash, error) {
	point := crypto.HashToPoint([]byte(network + "-" + tokenName))
	hash := new(common.Hash)
	err := hash.SetBytes(point.ToBytesS())
	if err != nil {
		return common.Hash{}, err
	}
	return *hash, nil
}

func (client *IncClient) CheckUnifiedShieldStatus(txHash string) (*ShieldStatus, error) {
	responseInBytes, err := client.rpcServer.CheckShieldUnifiedStatus(txHash)
	if err != nil {
		return nil, err
	}

	var status ShieldStatus
	err = rpchandler.ParseResponse(responseInBytes, &status)
	if err != nil {
		return nil, err
	}

	return &status, err
}

func (client *IncClient) CheckUnifiedUnshieldStatus(txHash string) (*UnshieldStatus, error) {
	responseInBytes, err := client.rpcServer.CheckUnshieldUnifiedStatus(txHash)
	if err != nil {
		return nil, err
	}

	var status UnshieldStatus
	err = rpchandler.ParseResponse(responseInBytes, &status)
	if err != nil {
		return nil, err
	}

	return &status, err
}

type ShieldStatusData struct {
	Amount uint64 `json:"Amount"`
	Reward uint64 `json:"Reward"`
}

type ShieldStatus struct {
	Status    byte               `json:"Status"`
	Data      []ShieldStatusData `json:"Data,omitempty"`
	ErrorCode int                `json:"ErrorCode,omitempty"`
}

type UnshieldStatusData struct {
	ReceivedAmount uint64 `json:"ReceivedAmount"`
	Fee            uint64 `json:"Fee"`
}

type UnshieldStatus struct {
	Status    byte                 `json:"Status"`
	Data      []UnshieldStatusData `json:"Data,omitempty"`
	ErrorCode int                  `json:"ErrorCode,omitempty"`
}
