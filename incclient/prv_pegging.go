package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"strings"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
)

// CreateIssuingPRVPeggingRequestTransaction creates a shielding trading transaction. By EVM, it means either ETH or BSC.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateIssuingPRVPeggingRequestTransaction(
	privateKey string, proof EVMDepositProof, evmNetworkID ...int,
) ([]byte, string, error) {
	tokenIDStr := common.PRVIDStr
	tokenID, err := new(common.Hash).NewHashFromStr(tokenIDStr)
	if err != nil {
		return nil, "", err
	}

	mdType := metadata.IssuingPRVERC20RequestMeta
	if len(evmNetworkID) > 0 {
		switch evmNetworkID[0] {
		case rpc.BSCNetworkID:
			mdType = metadata.IssuingPRVBEP20RequestMeta
		case rpc.PLGNetworkID:
			return nil, "", rpc.EVMNetworkNotFoundError(evmNetworkID[0])
		}
	}

	var issuingPRVPeggingRequestMeta *metadata.IssuingEVMRequest
	issuingPRVPeggingRequestMeta, err = metadata.NewIssuingEVMRequest(proof.blockHash, proof.txIdx, proof.nodeList, *tokenID, "", nil, mdType)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init issue eth request for %v, tokenID %v: %v", proof, tokenIDStr, err)
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, DefaultPRVFee, nil, issuingPRVPeggingRequestMeta, nil)
	return client.CreateRawTransaction(txParam, -1)
}

// CreateAndSendIssuingPRVPeggingRequestTransaction creates a PRV pegging shielding transaction,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendIssuingPRVPeggingRequestTransaction(
	privateKey string, proof EVMDepositProof, evmNetworkIDs ...int) (string, error) {
	encodedTx, txHash, err := client.CreateIssuingPRVPeggingRequestTransaction(privateKey, proof, evmNetworkIDs...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// CreateBurningPRVPeggingRequestTransaction creates a PRV pegging burning transaction for exiting the Incognito network.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateBurningPRVPeggingRequestTransaction(
	privateKey, remoteAddress string, burnedAmount uint64, evmNetworkIDs ...int,
) ([]byte, string, error) {
	tokenIDStr := common.PRVIDStr
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

	mdType := metadata.BurningPRVERC20RequestMeta
	if len(evmNetworkIDs) > 0 {
		switch evmNetworkIDs[0] {
		case rpc.BSCNetworkID:
			mdType = metadata.BurningPRVBEP20RequestMeta
		case rpc.PLGNetworkID:
			return nil, "", rpc.EVMNetworkNotFoundError(evmNetworkIDs[0])
		}
	}

	var md *metadata.BurningRequest
	md, err = metadata.NewBurningRequest(burnerAddress, burnedAmount, *tokenID, tokenIDStr, remoteAddress, mdType)
	if err != nil {
		return nil, "", fmt.Errorf("cannot init burning request with tokenID %v, burnedAmount %v, remoteAddress %v: %v",
			tokenIDStr, burnedAmount, remoteAddress, err)
	}

	txParam := NewTxParam(privateKey, []string{common.BurningAddress2}, []uint64{burnedAmount}, DefaultPRVFee, nil, md, nil)

	return client.CreateRawTransaction(txParam, -1)
}

// CreateAndSendBurningPRVPeggingRequestTransaction creates a PRV pegging burning transaction for exiting the Incognito network,
// and submits it to the network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendBurningPRVPeggingRequestTransaction(
	privateKey, remoteAddress string, burnedAmount uint64, evmNetworkIDs ...int,
) (string, error) {
	encodedTx, txHash, err := client.CreateBurningPRVPeggingRequestTransaction(privateKey, remoteAddress, burnedAmount, evmNetworkIDs...)
	if err != nil {
		return "", err
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// GetBurnPRVPeggingProof retrieves the burning proof for the Incognito network for submitting to the smart contract later.
func (client *IncClient) GetBurnPRVPeggingProof(txHash string, evmNetworkIDs ...int) (*jsonresult.InstructionProof, error) {
	responseInBytes, err := client.rpcServer.GetBurnPRVPeggingProof(txHash, evmNetworkIDs...)
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
