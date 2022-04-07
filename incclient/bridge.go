package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/privacy"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
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

// EVMDepositParams consists of parameters for creating an EVM shielding transaction.
// An EVMDepositParams is valid if at least one of the following conditions hold:
//	- Signature is not empty
//		- Receiver must not be empty
//	- Signature is empty
//		- Either DepositPrivateKey or DepositKeyIndex must not be empty; if DepositPrivateKey is empty, it will be
// 	derived from the DepositKeyIndex
//		- Receiver will be generated from the sender's private key
//		- Signature will be signed using the DepositPrivateKey
type EVMDepositParams struct {
	// MetadataType is the type of the shielding request.
	MetadataType int

	// DepositProof is the proof for the shielding request.
	DepositProof *EVMDepositProof

	// TokenID is the shielding asset ID.
	TokenID string

	// DepositPrivateKey is a base58-encoded deposit privateKey used to sign the request.
	// If set empty, it will be derived from the DepositKeyIndex.
	DepositPrivateKey string

	// DepositKeyIndex is the index of the OTDepositKey. It is used to generate DepositPrivateKey and DepositPubKey when the DepositPrivateKey is not supply.
	DepositKeyIndex uint64

	// Receiver is a base58-encoded OTAReceiver. If set empty, it will be generated from the sender's privateKey.
	Receiver string

	// Signature is a valid signature signed by the owner of the shielding asset.
	// If Signature is not empty, Receiver must not be empty.
	Signature string
}

// IsValid checks if a EVMDepositParams is valid.
func (dp EVMDepositParams) IsValid() (bool, error) {
	var err error
	_, err = common.Hash{}.NewHashFromStr(dp.TokenID)
	if err != nil || dp.TokenID == "" {
		return false, fmt.Errorf("invalid tokenID %v", dp.TokenID)
	}

	if len(dp.DepositProof.NodeList()) == 0 {
		return false, fmt.Errorf("invalid proofs")
	}

	if dp.Signature != "" {
		_, _, err = base58.Base58Check{}.Decode(dp.Signature)
		if err != nil {
			return false, fmt.Errorf("invalid signature")
		}
		if dp.Receiver == "" {
			return false, fmt.Errorf("must have `Receiver`")
		}
	} else {
		if dp.DepositPrivateKey != "" {
			_, _, err = base58.Base58Check{}.Decode(dp.DepositPrivateKey)
			if err != nil {
				return false, fmt.Errorf("invalid DepositPrivateKey")
			}
		}
	}

	if dp.Receiver != "" {
		otaReceiver := new(coin.OTAReceiver)
		err = otaReceiver.FromString(dp.Receiver)
		if err != nil {
			return false, fmt.Errorf("invalid receiver: %v", err)
		}
	}

	return true, nil
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
//
// An additional parameter `evmNetworkID` is introduced to specify the target EVM network. evmNetworkID can be one of the following:
//	- rpc.ETHNetworkID: the Ethereum network
//	- rpc.BSCNetworkID: the Binance Smart Chain network
//	- rpc.PLGNetworkID: the Polygon network
//	- rpc.FTMNetworkID: the Fantom network
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
	issuingETHRequestMeta, err = metadata.NewIssuingEVMRequest(proof.blockHash, proof.txIdx, proof.nodeList, *tokenID, "", nil, mdType)
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
//	- rpc.ETHNetworkID: the Ethereum network
//	- rpc.BSCNetworkID: the Binance Smart Chain network
//	- rpc.PLGNetworkID: the Polygon network
//	- rpc.FTMNetworkID: the Fantom network
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
//	- rpc.ETHNetworkID: the Ethereum network
//	- rpc.BSCNetworkID: the Binance Smart Chain network
//	- rpc.PLGNetworkID: the Polygon network
//	- rpc.FTMNetworkID: the Fantom network
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
//	- rpc.ETHNetworkID: the Ethereum network
//	- rpc.BSCNetworkID: the Binance Smart Chain network
//	- rpc.PLGNetworkID: the Polygon network
//	- rpc.FTMNetworkID: the Fantom network
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

// CreateEVMDepositTxWithDepositKey creates an EVM depositing transaction using one-time deposit key.
// It assumes the corresponding public transaction has been accepted.
//
// It returns the base58-encoded transaction, the transaction's hash, and an error (if any).
func (client *IncClient) CreateEVMDepositTxWithDepositKey(privateKey string, dp EVMDepositParams) ([]byte, string, error) {
	w, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, "", err
	}

	if _, err := dp.IsValid(); err != nil {
		return nil, "", err
	}

	receiver := dp.Receiver
	var sig []byte
	if dp.Signature != "" {
		sig, _, _ = base58.Base58Check{}.Decode(dp.Signature)
	} else {
		if receiver == "" {
			otaReceiver := new(coin.OTAReceiver)
			err = otaReceiver.FromAddress(w.KeySet.PaymentAddress)
			if err != nil {
				return nil, "", err
			}
			receiver = otaReceiver.String()
		}
		otaReceiver := new(coin.OTAReceiver)
		_ = otaReceiver.FromString(receiver)

		var depositPrivateKey *crypto.Scalar
		if dp.DepositPrivateKey != "" {
			tmp, _, _ := base58.Base58Check{}.Decode(dp.DepositPrivateKey)
			depositPrivateKey = new(crypto.Scalar).FromBytesS(tmp)
		} else {
			depositKey, err := client.GenerateDepositKeyFromPrivateKey(privateKey, dp.TokenID, dp.DepositKeyIndex)
			if err != nil {
				return nil, "", err
			}
			depositPrivateKey = new(crypto.Scalar).FromBytesS(depositKey.PrivateKey)
		}

		schnorrPrivateKey := new(privacy.SchnorrPrivateKey)
		r := new(crypto.Scalar).FromUint64(0) // must use r = 0
		schnorrPrivateKey.Set(depositPrivateKey, r)
		metaDataBytes := otaReceiver.Bytes()
		tmpSig, err := schnorrPrivateKey.Sign(common.HashB(metaDataBytes))
		if err != nil {
			return nil, "", err
		}
		sig = tmpSig.Bytes()
	}

	tokenID, _ := new(common.Hash).NewHashFromStr(dp.TokenID)

	md, err := metadata.NewIssuingEVMRequest(
		dp.DepositProof.BlockHash(),
		dp.DepositProof.TxIdx(),
		dp.DepositProof.NodeList(),
		*tokenID,
		receiver,
		sig,
		dp.MetadataType,
	)
	if err != nil {
		return nil, "", err
	}

	txParam := NewTxParam(privateKey, []string{}, []uint64{}, 0, nil, md, nil)
	return client.CreateRawTransaction(txParam, 2)

}

// CreateAndSendEVMDepositTxWithDepositKey creates an EVM shielding transaction using deposit keys,
// and submits it to the Incognito network.
//
// It returns the transaction's hash, and an error (if any).
func (client *IncClient) CreateAndSendEVMDepositTxWithDepositKey(
	privateKey string, dp EVMDepositParams) (string, error) {
	encodedTx, txHash, err := client.CreateEVMDepositTxWithDepositKey(privateKey,
		dp)
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
//	- rpc.ETHNetworkID: the Ethereum network
//	- rpc.BSCNetworkID: the Binance Smart Chain network
//	- rpc.PLGNetworkID: the Polygon network
//	- rpc.FTMNetworkID: the Fantom network
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
