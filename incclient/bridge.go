package incclient

import (
	"fmt"
	"strings"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
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

// CreateIssuingEVMRequestTransaction creates an EVM shielding trading transaction. By EVM, it means either ETH or BSC.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateIssuingEVMRequestTransaction(privateKey, tokenIDStr string, proof EVMDepositProof) ([]byte, string, error) {
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	var issuingETHRequestMeta *metadata.IssuingEVMRequest
	issuingETHRequestMeta, err = metadata.NewIssuingEVMRequest(proof.blockHash, proof.txIdx, proof.nodeList, *tokenID, metadata.IssuingETHRequestMeta)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init issue eth request for %v, tokenID %v: %v", proof, tokenIDStr, err)
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, issuingETHRequestMeta, nil)

	return client.CreateRawTransaction(txParam, -1)
}

// CreateAndSendIssuingEVMRequestTransaction creates an EVM shielding transaction, and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendIssuingEVMRequestTransaction(privateKey, tokenIDStr string, proof EVMDepositProof) (string, error) {
	encodedTx, txHash, err := client.CreateIssuingEVMRequestTransaction(privateKey, tokenIDStr, proof)
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
func (client *IncClient) CreateBurningRequestTransaction(privateKey, remoteAddress, tokenIDStr string, burnedAmount uint64) ([]byte, string, error) {
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

	var md *metadata.BurningRequest
	md, err = metadata.NewBurningRequest(burnerAddress, burnedAmount, *tokenID, tokenIDStr, remoteAddress, metadata.BurningRequestMetaV2)
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
func (client *IncClient) CreateAndSendBurningRequestTransaction(privateKey, remoteAddress, tokenIDStr string, burnedAmount uint64) (string, error) {
	encodedTx, txHash, err := client.CreateBurningRequestTransaction(privateKey, remoteAddress, tokenIDStr, burnedAmount)
	if err != nil {
		return "", err
	}

	err = client.SendRawTokenTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// GetBurnProof retrieves the burning proof for the Incognito network for submitting to the smart contract later.
func (client *IncClient) GetBurnProof(txHash string) (*jsonresult.InstructionProof, error) {
	responseInBytes, err := client.rpcServer.GetBurnProof(txHash)
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

// CheckShieldStatus returns the status of an eth-shielding request.
//	* -1: error
//	* 0: tx not found
//	* 1: tx is pending
//	* 2: tx is accepted
//	* 3: tx is rejected
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
