package incclient

import (
	"encoding/json"
	"fmt"
	"strings"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

type ETHDepositProof struct {
	blockNumber uint
	blockHash   ethCommon.Hash
	txIdx       uint
	nodeList    []string
}

func (E ETHDepositProof) TxIdx() uint {
	return E.txIdx
}

func (E ETHDepositProof) BlockNumber() uint {
	return E.blockNumber
}

func (E ETHDepositProof) BlockHash() ethCommon.Hash {
	return E.blockHash
}

func (E ETHDepositProof) NodeList() []string {
	return E.nodeList
}

func NewETHDepositProof(blockNumber uint, blockHash ethCommon.Hash, txIdx uint, nodeList []string) *ETHDepositProof {
	proof := ETHDepositProof{
		blockNumber: blockNumber,
		blockHash:   blockHash,
		txIdx:       txIdx,
		nodeList:    nodeList,
	}

	return &proof
}

func (client *IncClient) CreateIssuingETHRequestTransaction(privateKey, tokenIDStr string, proof ETHDepositProof) ([]byte, string, error) {
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	var issuingETHRequestMeta *metadata.IssuingETHRequest
	issuingETHRequestMeta, err = metadata.NewIssuingETHRequest(proof.blockHash, proof.txIdx, proof.nodeList, *tokenID, metadata.IssuingETHRequestMeta)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init issue eth request for %v, tokenID %v: %v", proof, tokenIDStr, err)
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, issuingETHRequestMeta, nil)

	return client.CreateRawTransaction(txParam, -1)
}
func (client *IncClient) CreateAndSendIssuingETHRequestTransaction(privateKey, tokenIDStr string, proof ETHDepositProof) (string, error) {
	encodedTx, txHash, err := client.CreateIssuingETHRequestTransaction(privateKey, tokenIDStr, proof)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

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
func (client *IncClient) GetBurnProof(txHash string) (*jsonresult.GetInstructionProof, error) {
	responseInBytes, err := client.rpcServer.GetBurnProof(txHash)
	if err != nil {
		return nil, err
	}

	response, err := rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return nil, err
	}

	var tmp *jsonresult.GetInstructionProof
	err = json.Unmarshal(response.Result, &tmp)
	if err != nil {
		return nil, err
	}

	return tmp, nil
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

	response, err := rpchandler.ParseResponse(responseInBytes)
	if err != nil {
		return -1, err
	}

	var status int
	err = json.Unmarshal(response.Result, &status)

	return status, err
}
