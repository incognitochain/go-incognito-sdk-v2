package pdexv3

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type Pdexv3Params struct {
	DefaultFeeRateBPS               uint            // the default value if fee rate is not specific in FeeRateBPS (default 0.3% ~ 30 BPS)
	FeeRateBPS                      map[string]uint // map: pool ID -> fee rate (0.1% ~ 10 BPS)
	PRVDiscountPercent              uint            // percent of fee that will be discounted if using PRV as the trading token fee (default: 25%)
	TradingProtocolFeePercent       uint            // percent of fees that is rewarded for the core team (default: 0%)
	TradingStakingPoolRewardPercent uint            // percent of fees that is distributed for staking pools (PRV, PDEX, ..., default: 10%)
	PDEXRewardPoolPairsShare        map[string]uint // map: pool pair ID -> PDEX reward share weight
	StakingPoolsShare               map[string]uint // map: staking tokenID -> pool staking share weight
	StakingRewardTokens             []common.Hash   // list of staking reward tokens
	MintNftRequireAmount            uint64          // amount prv for depositing to pdex
	MaxOrdersPerNft                 uint            // max orders per nft
	AutoWithdrawOrderLimitAmount    uint            // max orders will be auto withdraw each shard for each blocks
	MinPRVReserveTradingRate        uint64          // min prv reserve for checking price of trading fee paid by PRV
	OrderMiningRewardRatioBPS       map[string]uint // map: pool ID -> ratio of LOP rewards compare with LP rewards (0.1% ~ 10 BPS)
}

type ParamsModifyingRequest struct {
	metadataCommon.MetadataBaseWithSignature
	Pdexv3Params `json:"Pdexv3Params"`
}

type ParamsModifyingContent struct {
	Content  Pdexv3Params `json:"Content"`
	ErrorMsg string       `json:"ErrorMsg"`
	TxReqID  common.Hash  `json:"TxReqID"`
	ShardID  byte         `json:"ShardID"`
}

type ParamsModifyingRequestStatus struct {
	Status       int    `json:"Status"`
	ErrorMsg     string `json:"ErrorMsg"`
	Pdexv3Params `json:"Pdexv3Params"`
}

func (paramsModifying ParamsModifyingRequest) Hash() *common.Hash {
	record := paramsModifying.MetadataBaseWithSignature.Hash().String()
	if paramsModifying.Sig != nil && len(paramsModifying.Sig) != 0 {
		record += string(paramsModifying.Sig)
	}
	contentBytes, _ := json.Marshal(paramsModifying.Pdexv3Params)
	hashParams := common.HashH(contentBytes)
	record += hashParams.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (paramsModifying ParamsModifyingRequest) HashWithoutSig() *common.Hash {
	record := paramsModifying.MetadataBaseWithSignature.Hash().String()
	contentBytes, _ := json.Marshal(paramsModifying.Pdexv3Params)
	hashParams := common.HashH(contentBytes)
	record += hashParams.String()

	// final hash
	hash := common.HashH([]byte(record))
	return &hash
}

func (paramsModifying *ParamsModifyingRequest) CalculateSize() uint64 {
	return metadataCommon.CalculateSize(paramsModifying)
}
