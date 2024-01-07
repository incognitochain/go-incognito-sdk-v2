package jsonresult

import (
	"math/big"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
)

type BridgeAggState struct {
	BeaconTimeStamp     int64                                                `json:"BeaconTimeStamp"`
	UnifiedTokenVaults  map[common.Hash]map[common.Hash]*BridgeAggVaultState `json:"UnifiedTokenVaults"`
	WaitingUnshieldReqs map[common.Hash][]*BridgeAggWaitingUnshieldReq       `json:"WaitingUnshieldReqs"`
	Param               *BridgeAggParamState                                 `json:"Param"`
	BaseDecimal         uint8                                                `json:"BaseDecimal"`
	MaxLenOfPath        uint8                                                `json:"MaxLenOfPath"`
}

type BridgeAggVaultState struct {
	Amount                uint64      `json:"Amount"`
	LockedAmount          uint64      `json:"LockedAmount"`
	WaitingUnshieldAmount uint64      `json:"WaitingUnshieldAmount"`
	WaitingUnshieldFee    uint64      `json:"WaitingUnshieldFee"`
	ExtDecimal            uint8       `json:"ExtDecimal"`
	NetworkID             uint8       `json:"NetworkID"`
	IncTokenID            common.Hash `json:"IncTokenID"`
}

type BridgeAggWaitingUnshieldReq struct {
	UnshieldID   common.Hash              `json:"UnshieldID"`
	Data         []WaitingUnshieldReqData `json:"Data"`
	BeaconHeight uint64                   `json:"BeaconHeight"`
}

type WaitingUnshieldReqData struct {
	IncTokenID             common.Hash `json:"IncTokenID"`
	BurningAmount          uint64      `json:"BurningAmount"`
	RemoteAddress          string      `json:"RemoteAddress"`
	Fee                    uint64      `json:"Fee"`
	ExternalTokenID        []byte      `json:"ExternalTokenID"`
	ExternalReceivedAmt    *big.Int    `json:"ExternalReceivedAmt"`
	BurningConfirmMetaType uint        `json:"BurningConfirmMetaType"`
}

type BridgeAggParamState struct {
	PercentFeeWithDec uint64
}

type BridgeAggEstimateFeeByBurntAmount struct {
	BurntAmount    uint64 `json:"BurntAmount"`
	Fee            uint64 `json:"Fee"`
	ReceivedAmount uint64 `json:"ReceivedAmount"`

	MaxFee            uint64 `json:"MaxFee"`
	MinReceivedAmount uint64 `json:"MinReceivedAmount"`
}

type BridgeAggEstimateFeeByReceivedAmount struct {
	ReceivedAmount uint64 `json:"ReceivedAmount"`
	Fee            uint64 `json:"Fee"`
	BurntAmount    uint64 `json:"BurntAmount"`

	MaxFee         uint64 `json:"MaxFee"`
	MaxBurntAmount uint64 `json:"MaxBurntAmount"`
}

type BridgeAggEstimateReward struct {
	ReceivedAmount uint64 `json:"ReceivedAmount"`
	Reward         uint64 `json:"Reward"`
}
