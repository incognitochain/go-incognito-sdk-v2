package metadata

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"strconv"
)

// PDECrossPoolTradeRequest - privacy dex cross pool trade
type PDECrossPoolTradeRequest struct {
	TokenIDToBuyStr     string
	TokenIDToSellStr    string
	SellAmount          uint64 // must be equal to vout value
	MinAcceptableAmount uint64
	TradingFee          uint64
	TraderAddressStr    string
	TxRandomStr         string `json:"TxRandomStr,omitempty"`
	SubTraderAddressStr string `json:"SubTraderAddressStr,omitempty"`
	SubTxRandomStr      string `json:"SubTxRandomStr,omitempty"`
	MetadataBase
}

type PDECrossPoolTradeRequestAction struct {
	Meta    PDECrossPoolTradeRequest
	TxReqID common.Hash
	ShardID byte
}

type PDECrossPoolTradeAcceptedContent struct {
	TraderAddressStr         string
	TxRandomStr              string `json:"TxRandomStr,omitempty"`
	TokenIDToBuyStr          string
	ReceiveAmount            uint64
	Token1IDStr              string
	Token2IDStr              string
	Token1PoolValueOperation TokenPoolValueOperation
	Token2PoolValueOperation TokenPoolValueOperation
	ShardID                  byte
	RequestedTxID            common.Hash
	AddingFee                uint64
}

type PDERefundCrossPoolTrade struct {
	TraderAddressStr string
	TxRandomStr      string `json:"TxRandomStr,omitempty"`
	TokenIDStr       string
	Amount           uint64
	ShardID          byte
	TxReqID          common.Hash
}

func NewPDECrossPoolTradeRequest(
	tokenIDToBuyStr string,
	tokenIDToSellStr string,
	sellAmount uint64,
	minAcceptableAmount uint64,
	tradingFee uint64,
	traderAddressStr string,
	txRandomStr string,
	subTraderAddressStr string,
	subTxRandomStr string,
	metaType int,
) (*PDECrossPoolTradeRequest, error) {
	metadataBase := MetadataBase{
		Type: metaType,
	}
	pdeCrossPoolTradeRequest := &PDECrossPoolTradeRequest{
		TokenIDToBuyStr:     tokenIDToBuyStr,
		TokenIDToSellStr:    tokenIDToSellStr,
		SellAmount:          sellAmount,
		MinAcceptableAmount: minAcceptableAmount,
		TradingFee:          tradingFee,
		TraderAddressStr:    traderAddressStr,
		TxRandomStr:         txRandomStr,
		SubTraderAddressStr: subTraderAddressStr,
		SubTxRandomStr:      subTxRandomStr,
	}
	pdeCrossPoolTradeRequest.MetadataBase = metadataBase
	return pdeCrossPoolTradeRequest, nil
}

func (pc PDECrossPoolTradeRequest) Hash() *common.Hash {
	record := pc.MetadataBase.Hash().String()
	record += pc.TokenIDToBuyStr
	record += pc.TokenIDToSellStr
	record += pc.TraderAddressStr
	if len(pc.TxRandomStr) > 0 {
		record += pc.TxRandomStr
	}
	if len(pc.SubTraderAddressStr) > 0 {
		record += pc.SubTraderAddressStr
	}
	if len(pc.SubTxRandomStr) > 0 {
		record += pc.SubTxRandomStr
	}
	record += strconv.FormatUint(pc.SellAmount, 10)
	record += strconv.FormatUint(pc.MinAcceptableAmount, 10)
	record += strconv.FormatUint(pc.TradingFee, 10)
	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (pc *PDECrossPoolTradeRequest) CalculateSize() uint64 {
	return calculateSize(pc)
}
