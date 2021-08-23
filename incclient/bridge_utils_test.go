package incclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"log"
	"testing"
)

func TestIncClient_GetEVMTxByHash(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	txHash := "0x392bf44aeac2c395fc4ed7ba425f1fc61b7b62d98a96c2a2d5e22c5ec8cd8f23"

	txDetail, err := ic.GetEVMTxByHash(txHash)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(txDetail, "", "\t")
	if err != nil {
		panic(err)
	}

	Logger.Println(string(jsb))
}

func TestIncClient_GetEVMBlockByHash(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	blockHash := "0xe886b93d341bb6cf1f4e24a2ffa40c0a6107adb6214e8f7e43fce04d07fc3f1f"

	blockDetail, err := ic.GetEVMBlockByHash(blockHash)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(blockDetail, "", "\t")
	if err != nil {
		panic(err)
	}

	Logger.Println(string(jsb))
}

func TestIncClient_GetEVMTxReceipt(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	txHash := "0x392bf44aeac2c395fc4ed7ba425f1fc61b7b62d98a96c2a2d5e22c5ec8cd8f23"

	receipt, err := ic.GetEVMTxReceipt(txHash)
	if err != nil {
		panic(err)
	}

	jsb, err := json.MarshalIndent(receipt, "", "\t")
	if err != nil {
		panic(err)
	}

	Logger.Println(string(jsb))
}

func TestIncClient_GetEVMDepositProof(t *testing.T) {
	var err error
	ic, err = NewTestNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "0x44f01c88fe21ed42408b70312a3899893497fdb6d215f87c9f038adae978a484"

	depositProof, amount, err := ic.GetEVMDepositProof(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Println(amount)
	Logger.Println(depositProof.BlockNumber(), depositProof.BlockHash().String(), depositProof.TxIdx())
	Logger.Println(depositProof.NodeList())
}

func TestIncClient_GetMostRecentEVMBlockNumber(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	mostRecentBlock, err := ic.GetMostRecentEVMBlockNumber()
	if err != nil {
		panic(err)
	}

	Logger.Printf("mostRecentBlock: %v\n", mostRecentBlock)
}

func TestIncClient_GetEVMTransactionStatus(t *testing.T) {
	var err error
	ic, err = NewTestNet1Client()
	if err != nil {
		panic(err)
	}

	txHash := "0x392bf44aeac2c395fc4ed7ba425f1fc61b7b62d98a96c2a2d5e22c5ec8cd8f23"

	status, err := ic.GetEVMTransactionStatus(txHash)
	if err != nil {
		panic(err)
	}

	Logger.Printf("status: %v\n", status)
}

func TestIncClient_GetBurnProof(t *testing.T) {
	var err error
	ic, err = NewDevNetClient()
	if err != nil {
		panic(err)
	}

	txHash := "c87985f9b09012dc182dffffc8630d7396aa34d0c541c265e5c9d777755e0754"
	burnProof, err := ic.GetBurnProof(txHash)
	if err != nil {
		panic(err)
	}

	log.Println(burnProof)
}

func TestIncClient_CreateAndSendIssuingEVMRequestTransaction(t *testing.T) {
	var err error
	ic, err = NewTestNetClientWithCache()
	if err != nil {
		panic(err)
	}

	privateKey := ""
	txHash := "0x44f01c88fe21ed42408b70312a3899893497fdb6d215f87c9f038adae978a484"

	depositProof, _, err := ic.GetEVMDepositProof(txHash)
	if err != nil {
		panic(err)
	}

	encodedTx, incTxHash, err := ic.CreateIssuingEVMRequestTransaction(privateKey, pEthID, *depositProof)
	if err != nil {
		panic(err)
	}
	Logger.Printf("incTxHash: %v\n", incTxHash)

	tx := new(tx_ver2.Tx)
	rawTxData, _, err := base58.Base58Check{}.Decode(string(encodedTx))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(rawTxData, &tx)
	if err != nil {
		panic(err)
	}
	md := tx.GetMetadata().(*metadata.IssuingEVMRequest)
	jsb, _ := json.Marshal(md)
	Logger.Println(string(jsb))
	_, err = verifyProofAndParseReceipt(md)
	if err != nil {
		panic(err)
	}
	Logger.Println("Verify proof SUCCEEDED!!!")

	err = ic.SendRawTx(encodedTx)
	if err != nil {
		panic(err)
	}
	Logger.Println("SendRawTX SUCCEEDED!!")
}

func verifyProofAndParseReceipt(iReq *metadata.IssuingEVMRequest) (*types.Receipt, error) {
	ethClient, err := ethclient.Dial(TestNetETHHost)
	evmHeader, err := ethClient.HeaderByHash(context.Background(), iReq.BlockHash)
	if err != nil {
		return nil, err
	}
	if evmHeader == nil {
		return nil, fmt.Errorf("WARNING: Could not find out the EVM block header with the hash: %s", iReq.BlockHash.String())
	}

	log.Println("evmHeader:", evmHeader.ReceiptHash.String())

	keyBuf := new(bytes.Buffer)
	keyBuf.Reset()
	err = rlp.Encode(keyBuf, iReq.TxIndex)
	if err != nil {
		return nil, err
	}

	log.Println("keyBuf:", keyBuf.Bytes())

	nodeList := new(light.NodeList)
	for i, proofStr := range iReq.ProofStrs {
		proofBytes, err := base64.StdEncoding.DecodeString(proofStr)
		if err != nil {
			return nil, err
		}
		log.Println(i, proofBytes)
		err = nodeList.Put([]byte{}, proofBytes)
		if err != nil {
			return nil, err
		}
	}
	proof := nodeList.NodeSet()
	val, err := trie.VerifyProof(evmHeader.ReceiptHash, keyBuf.Bytes(), proof)
	if err != nil {
		return nil, err
	}

	// Decode value from VerifyProof into Receipt
	constructedReceipt := new(types.Receipt)
	err = rlp.DecodeBytes(val, constructedReceipt)
	if err != nil {
		return nil, err
	}

	if constructedReceipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("the constructedReceipt's status is not success")
	}

	return constructedReceipt, nil
}
