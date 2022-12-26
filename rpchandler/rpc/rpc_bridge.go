package rpc

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
)

const (
	ETHNetworkID = iota
	BSCNetworkID
	PLGNetworkID
	FTMNetworkID
	AURORANetworkID
	AVAXNetworkID
)

// EVMIssuingMetadata keeps track of EVM issuing metadata types based on the EVM networkIDs.
var EVMIssuingMetadata = map[int]int{
	ETHNetworkID:    metadata.IssuingETHRequestMeta,
	BSCNetworkID:    metadata.IssuingBSCRequestMeta,
	PLGNetworkID:    metadata.IssuingPLGRequestMeta,
	FTMNetworkID:    metadata.IssuingFantomRequestMeta,
	AURORANetworkID: metadata.IssuingAuroraRequestMeta,
	AVAXNetworkID:   metadata.IssuingAvaxRequestMeta,
}

// EVMBurningMetadata keeps track of EVM burning metadata types based on the EVM networkIDs.
var EVMBurningMetadata = map[int]int{
	ETHNetworkID:    metadata.BurningRequestMetaV2,
	BSCNetworkID:    metadata.BurningPBSCRequestMeta,
	PLGNetworkID:    metadata.BurningPLGRequestMeta,
	FTMNetworkID:    metadata.BurningFantomRequestMeta,
	AURORANetworkID: metadata.BurningAuroraRequestMeta,
	AVAXNetworkID:   metadata.BurningAvaxRequestMeta,
}

var burnProofRPCMethod = map[int]string{
	ETHNetworkID:    getBurnProof,
	BSCNetworkID:    getBSCBurnProof,
	PLGNetworkID:    getPLGBurnProof,
	FTMNetworkID:    getFTMBurnProof,
	AURORANetworkID: getAURORABurnProof,
	AVAXNetworkID:   getAVAXBurnProof,
}

// EVMNetworkNotFoundError returns an error indicating that the given EVM networkID is not supported.
func EVMNetworkNotFoundError(evmNetworkID int) error {
	return fmt.Errorf("EVMNetworkID %v not supported", evmNetworkID)
}

// GetBurnProof retrieves the burning proof of a transaction with the given target evmNetworkID.
// evmNetworkID can be one of the following:
//   - ETHNetworkID: the Ethereum network
//   - BSCNetworkID: the Binance Smart Chain network
//   - PLGNetworkID: the Polygon network
//   - FTMNetworkID: the Fantom network
//
// If set empty, evmNetworkID defaults to ETHNetworkID. NOTE that only the first value of evmNetworkID is used.
func (server *RPCServer) GetBurnProof(txHash string, evmNetworkID ...int) ([]byte, error) {
	networkID := ETHNetworkID
	if len(evmNetworkID) > 0 {
		networkID = evmNetworkID[0]
	}

	if _, ok := burnProofRPCMethod[networkID]; !ok {
		return nil, EVMNetworkNotFoundError(networkID)
	}
	method := burnProofRPCMethod[networkID]
	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(method, params)
}

// GetBurnProofForSC retrieves the burning proof of a transaction for depositing to smart contracts.
func (server *RPCServer) GetBurnProofForSC(txHash string) ([]byte, error) {
	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(getBurnProofForDepositToSC, params)
}

// GetBurnPRVPeggingProof retrieves the burning prv pegging proof of a transaction.
func (server *RPCServer) GetBurnPRVPeggingProof(txHash string, evmNetworkIDs ...int) ([]byte, error) {
	method := getPRVERC20BurnProof
	if len(evmNetworkIDs) > 0 {
		switch evmNetworkIDs[0] {
		case BSCNetworkID:
			method = getPRVBEP20BurnProof
		case PLGNetworkID, FTMNetworkID, AURORANetworkID, AVAXNetworkID:
			return nil, EVMNetworkNotFoundError(evmNetworkIDs[0])
		}
	}
	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(method, params)
}

// CheckShieldStatus checks the status of a decentralized shielding transaction.
func (server *RPCServer) CheckShieldStatus(txHash string) ([]byte, error) {
	tmpParams := make(map[string]interface{})
	tmpParams["TxReqID"] = txHash

	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	return server.SendQuery(getBridgeReqWithStatus, params)
}

// GetAllBridgeTokens retrieves the list of bridge tokens in the network.
func (server *RPCServer) GetAllBridgeTokens() ([]byte, error) {
	return server.SendQuery(getAllBridgeTokens, nil)
}

// GetBridgeAggState get bridge aggregator state
func (server *RPCServer) GetBridgeAggState(beaconHeight uint64) ([]byte, error) {
	tmpParams := make(map[string]interface{})
	tmpParams["BeaconHeight"] = beaconHeight

	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	return server.SendQuery(bridgeaggState, params)
}

// CheckUnshieldUnifiedStatus checks the status of a decentralized unshielding transaction.
func (server *RPCServer) CheckUnshieldUnifiedStatus(txHash string) ([]byte, error) {

	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(bridgeaggStatusUnshield, params)
}

// CheckShieldUnifiedStatus checks the status of a decentralized shielding transaction.
func (server *RPCServer) CheckShieldUnifiedStatus(txHash string) ([]byte, error) {

	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(bridgeaggStatusShield, params)
}

// CheckConvertStatuspUnifiedStatus checks the status of a decentralized convert transaction.
func (server *RPCServer) CheckConvertStatuspUnifiedStatus(txHash string) ([]byte, error) {

	params := make([]interface{}, 0)
	params = append(params, txHash)
	return server.SendQuery(bridgeaggStatusConvert, params)
}

func (server *RPCServer) GetBridgeAggEstimateFeeByExpectedAmount(pUnifiedTokenID, tokenID string, expectedAmount uint64) ([]byte, error) {
	tmpParams := make(map[string]interface{})
	tmpParams["UnifiedTokenID"] = pUnifiedTokenID
	tmpParams["TokenID"] = tokenID
	tmpParams["ExpectedAmount"] = expectedAmount

	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	return server.SendQuery(bridgeaggEstimateFeeByExpectedAmount, params)
}

func (server *RPCServer) GetBridgeAggEstimateFeeByBurntAmount(pUnifiedTokenID, tokenID string, burnAmount uint64) ([]byte, error) {
	tmpParams := make(map[string]interface{})
	tmpParams["UnifiedTokenID"] = pUnifiedTokenID
	tmpParams["TokenID"] = tokenID
	tmpParams["BurntAmount"] = burnAmount

	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	return server.SendQuery(bridgeaggEstimateFeeByBurntAmount, params)
}

func (server *RPCServer) GetBridgeAggEstimateReward(pUnifiedTokenID, tokenID string, amount uint64) ([]byte, error) {
	tmpParams := make(map[string]interface{})
	tmpParams["UnifiedTokenID"] = pUnifiedTokenID
	tmpParams["TokenID"] = tokenID
	tmpParams["Amount"] = amount

	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	return server.SendQuery(bridgeaggEstimateReward, params)
}

func (server *RPCServer) GetBridgeAggGetBurnProof(txReqID string, dataIndex ...int) ([]byte, error) {
	tmpParams := make(map[string]interface{})
	tmpParams["TxReqID"] = txReqID
	if len(dataIndex) > 0 {
		tmpParams["DataIndex"] = dataIndex[0]
	}
	params := make([]interface{}, 0)
	params = append(params, tmpParams)
	return server.SendQuery(bridgeaggGetBurnProof, params)
}
