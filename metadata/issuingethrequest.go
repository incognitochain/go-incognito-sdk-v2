package metadata

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	rCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type IssuingETHRequest struct {
	BlockHash  rCommon.Hash
	TxIndex    uint
	ProofStrs  []string
	IncTokenID common.Hash
	MetadataBase
}

type IssuingETHReqAction struct {
	Meta       IssuingETHRequest `json:"meta"`
	TxReqID    common.Hash       `json:"txReqId"`
	ETHReceipt *types.Receipt    `json:"ethReceipt"`
}

type IssuingETHAcceptedInst struct {
	ShardID         byte        `json:"shardId"`
	IssuingAmount   uint64      `json:"issuingAmount"`
	ReceiverAddrStr string      `json:"receiverAddrStr"`
	IncTokenID      common.Hash `json:"incTokenId"`
	TxReqID         common.Hash `json:"txReqId"`
	UniqETHTx       []byte      `json:"uniqETHTx"`
	ExternalTokenID []byte      `json:"externalTokenId"`
}

func ParseETHIssuingInstContent(instContentStr string) (*IssuingETHReqAction, error) {
	contentBytes, err := base64.StdEncoding.DecodeString(instContentStr)
	if err != nil {
		return nil, err
	}
	var issuingETHReqAction IssuingETHReqAction
	err = json.Unmarshal(contentBytes, &issuingETHReqAction)
	if err != nil {
		return nil, err
	}
	return &issuingETHReqAction, nil
}

func NewIssuingETHRequest(
	blockHash rCommon.Hash,
	txIndex uint,
	proofStrs []string,
	incTokenID common.Hash,
	metaType int,
) (*IssuingETHRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	issuingETHReq := &IssuingETHRequest{
		BlockHash:  blockHash,
		TxIndex:    txIndex,
		ProofStrs:  proofStrs,
		IncTokenID: incTokenID,
	}
	issuingETHReq.MetadataBase = metadataBase
	return issuingETHReq, nil
}

func NewIssuingETHRequestFromMap(
	data map[string]interface{},
) (*IssuingETHRequest, error) {
	blockHash := rCommon.HexToHash(data["BlockHash"].(string))
	txIdx := uint(data["TxIndex"].(float64))
	proofsRaw := data["ProofStrs"].([]interface{})
	proofStrs := []string{}
	for _, item := range proofsRaw {
		proofStrs = append(proofStrs, item.(string))
	}

	incTokenID, err := common.Hash{}.NewHashFromStr(data["IncTokenID"].(string))
	if err != nil {
		return nil, errors.Errorf("TokenID incorrect")
	}

	req, _ := NewIssuingETHRequest(
		blockHash,
		txIdx,
		proofStrs,
		*incTokenID,
		IssuingETHRequestMeta,
	)
	return req, nil
}

func (iReq IssuingETHRequest) Hash() *common.Hash {
	record := iReq.BlockHash.String()
	// TODO: @hung change to record += fmt.Sprint(iReq.TxIndex)
	record += string(iReq.TxIndex)
	proofStrs := iReq.ProofStrs
	for _, proofStr := range proofStrs {
		record += proofStr
	}
	record += iReq.MetadataBase.Hash().String()
	record += iReq.IncTokenID.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (iReq *IssuingETHRequest) CalculateSize() uint64 {
	return calculateSize(iReq)
}

func ParseETHLogData(data []byte) (map[string]interface{}, error) {
	abiIns, err := abi.JSON(strings.NewReader(common.AbiJson))
	if err != nil {
		return nil, err
	}
	dataMap := map[string]interface{}{}
	if err = abiIns.UnpackIntoMap(dataMap, "Deposit", data); err != nil {
		return nil, err
	}
	return dataMap, nil
}

func IsETHTxHashUsedInBlock(uniqETHTx []byte, uniqETHTxsUsed [][]byte) bool {
	for _, item := range uniqETHTxsUsed {
		if bytes.Equal(uniqETHTx, item) {
			return true
		}
	}
	return false
}

func PickAndParseLogMapFromReceipt(constructedReceipt *types.Receipt, ethContractAddressStr string) (map[string]interface{}, error) {
	logData := []byte{}
	logLen := len(constructedReceipt.Logs)
	if logLen == 0 {
		fmt.Println("WARNING: LOG data is invalid.")
		return nil, nil
	}
	for _, log := range constructedReceipt.Logs {
		if bytes.Equal(rCommon.HexToAddress(ethContractAddressStr).Bytes(), log.Address.Bytes()) {
			logData = log.Data
			break
		}
	}
	if len(logData) == 0 {
		return nil, nil
	}
	return ParseETHLogData(logData)
}
