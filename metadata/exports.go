package metadata

import (
	// "github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

// export interfaces
type Metadata = metadataCommon.Metadata
type MetadataBase = metadataCommon.MetadataBase
type MetadataBaseWithSignature = metadataCommon.MetadataBaseWithSignature
type Transaction = metadataCommon.Transaction
// type ChainRetriever = metadataCommon.ChainRetriever
// type ShardViewRetriever = metadataCommon.ShardViewRetriever
// type BeaconViewRetriever = metadataCommon.BeaconViewRetriever
// type MempoolRetriever = metadataCommon.MempoolRetriever
// type ValidationEnviroment = metadataCommon.ValidationEnviroment
// type TxDesc = metadataCommon.TxDesc

// export structs
type OTADeclaration = metadataCommon.OTADeclaration
// type MintData = metadataCommon.MintData
// type AccumulatedValues = metadataCommon.AccumulatedValues

var AcceptedWithdrawRewardRequestVersion = metadataCommon.AcceptedWithdrawRewardRequestVersion

// export functions
var AssertPaymentAddressAndTxVersion = metadataCommon.AssertPaymentAddressAndTxVersion
var GenTokenIDFromRequest = metadataCommon.GenTokenIDFromRequest
var NewMetadataBase = metadataCommon.NewMetadataBase
var NewMetadataBaseWithSignature = metadataCommon.NewMetadataBaseWithSignature
// var ValidatePortalExternalAddress = metadataCommon.ValidatePortalExternalAddress
// var NewMetadataTxError = metadataCommon.NewMetadataTxError
var IsAvailableMetaInTxType = metadataCommon.IsAvailableMetaInTxType
var NoInputNoOutput = metadataCommon.NoInputNoOutput
var NoInputHasOutput = metadataCommon.NoInputHasOutput
var IsPortalRelayingMetaType = metadataCommon.IsPortalRelayingMetaType
var IsPortalMetaTypeV3 = metadataCommon.IsPortalMetaTypeV3
var GetMetaAction = metadataCommon.GetMetaAction
// var IsPDEType = metadataCommon.IsPDEType
var GetLimitOfMeta = metadataCommon.GetLimitOfMeta
// var IsPDETx = metadataCommon.IsPDETx
// var IsPdexv3Tx = metadataCommon.IsPdexv3Tx
// var ConvertPrivacyTokenToNativeToken = metadataCommon.ConvertPrivacyTokenToNativeToken
// var ConvertNativeTokenToPrivacyToken = metadataCommon.ConvertNativeTokenToPrivacyToken
var HasBridgeInstructions = metadataCommon.HasBridgeInstructions
var HasPortalInstructions = metadataCommon.HasPortalInstructions

var calculateSize = metadataCommon.CalculateSize

// export package constants
const (
	InvalidMeta                  = metadataCommon.InvalidMeta
	IssuingRequestMeta           = metadataCommon.IssuingRequestMeta
	IssuingResponseMeta          = metadataCommon.IssuingResponseMeta
	ContractingRequestMeta       = metadataCommon.ContractingRequestMeta
	IssuingETHRequestMeta        = metadataCommon.IssuingETHRequestMeta
	IssuingETHResponseMeta       = metadataCommon.IssuingETHResponseMeta
	ShardBlockReward             = metadataCommon.ShardBlockReward
	ShardBlockSalaryResponseMeta = metadataCommon.ShardBlockSalaryResponseMeta
	BeaconRewardRequestMeta      = metadataCommon.BeaconRewardRequestMeta
	BeaconSalaryResponseMeta     = metadataCommon.BeaconSalaryResponseMeta
	ReturnStakingMeta            = metadataCommon.ReturnStakingMeta
	IncDAORewardRequestMeta      = metadataCommon.IncDAORewardRequestMeta
	WithDrawRewardRequestMeta    = metadataCommon.WithDrawRewardRequestMeta
	WithDrawRewardResponseMeta   = metadataCommon.WithDrawRewardResponseMeta
	//staking
	ShardStakingMeta    = metadataCommon.ShardStakingMeta
	StopAutoStakingMeta = metadataCommon.StopAutoStakingMeta
	BeaconStakingMeta   = metadataCommon.BeaconStakingMeta
	UnStakingMeta       = metadataCommon.UnStakingMeta
	// Incognito -> Ethereum bridge
	BeaconSwapConfirmMeta = metadataCommon.BeaconSwapConfirmMeta
	BridgeSwapConfirmMeta = metadataCommon.BridgeSwapConfirmMeta
	BurningRequestMeta    = metadataCommon.BurningRequestMeta
	BurningRequestMetaV2  = metadataCommon.BurningRequestMetaV2
	BurningConfirmMeta    = metadataCommon.BurningConfirmMeta
	BurningConfirmMetaV2  = metadataCommon.BurningConfirmMetaV2
	// pde
	PDEContributionMeta                   = metadataCommon.PDEContributionMeta
	PDETradeRequestMeta                   = metadataCommon.PDETradeRequestMeta
	PDETradeResponseMeta                  = metadataCommon.PDETradeResponseMeta
	PDEWithdrawalRequestMeta              = metadataCommon.PDEWithdrawalRequestMeta
	PDEWithdrawalResponseMeta             = metadataCommon.PDEWithdrawalResponseMeta
	PDEContributionResponseMeta           = metadataCommon.PDEContributionResponseMeta
	PDEPRVRequiredContributionRequestMeta = metadataCommon.PDEPRVRequiredContributionRequestMeta
	PDECrossPoolTradeRequestMeta          = metadataCommon.PDECrossPoolTradeRequestMeta
	PDECrossPoolTradeResponseMeta         = metadataCommon.PDECrossPoolTradeResponseMeta
	PDEFeeWithdrawalRequestMeta           = metadataCommon.PDEFeeWithdrawalRequestMeta
	PDEFeeWithdrawalResponseMeta          = metadataCommon.PDEFeeWithdrawalResponseMeta
	PDETradingFeesDistributionMeta        = metadataCommon.PDETradingFeesDistributionMeta
	// pDEX v3
	Pdexv3TradeRequestMeta          = metadataCommon.Pdexv3TradeRequestMeta
	Pdexv3TradeResponseMeta         = metadataCommon.Pdexv3TradeResponseMeta
	Pdexv3AddOrderRequestMeta       = metadataCommon.Pdexv3AddOrderRequestMeta
	Pdexv3AddOrderResponseMeta      = metadataCommon.Pdexv3AddOrderResponseMeta
	Pdexv3WithdrawOrderRequestMeta  = metadataCommon.Pdexv3WithdrawOrderRequestMeta
	Pdexv3WithdrawOrderResponseMeta = metadataCommon.Pdexv3WithdrawOrderResponseMeta
	// portal
	PortalCustodianDepositMeta                  = metadataCommon.PortalCustodianDepositMeta
	PortalRequestPortingMeta                    = metadataCommon.PortalRequestPortingMeta
	PortalUserRequestPTokenMeta                 = metadataCommon.PortalUserRequestPTokenMeta
	PortalCustodianDepositResponseMeta          = metadataCommon.PortalCustodianDepositResponseMeta
	PortalUserRequestPTokenResponseMeta         = metadataCommon.PortalUserRequestPTokenResponseMeta
	PortalExchangeRatesMeta                     = metadataCommon.PortalExchangeRatesMeta
	PortalRedeemRequestMeta                     = metadataCommon.PortalRedeemRequestMeta
	PortalRedeemRequestResponseMeta             = metadataCommon.PortalRedeemRequestResponseMeta
	PortalRequestUnlockCollateralMeta           = metadataCommon.PortalRequestUnlockCollateralMeta
	PortalCustodianWithdrawRequestMeta          = metadataCommon.PortalCustodianWithdrawRequestMeta
	PortalCustodianWithdrawResponseMeta         = metadataCommon.PortalCustodianWithdrawResponseMeta
	PortalLiquidateCustodianMeta                = metadataCommon.PortalLiquidateCustodianMeta
	PortalLiquidateCustodianResponseMeta        = metadataCommon.PortalLiquidateCustodianResponseMeta
	PortalLiquidateTPExchangeRatesMeta          = metadataCommon.PortalLiquidateTPExchangeRatesMeta
	PortalExpiredWaitingPortingReqMeta          = metadataCommon.PortalExpiredWaitingPortingReqMeta
	PortalRewardMeta                            = metadataCommon.PortalRewardMeta
	PortalRequestWithdrawRewardMeta             = metadataCommon.PortalRequestWithdrawRewardMeta
	PortalRequestWithdrawRewardResponseMeta     = metadataCommon.PortalRequestWithdrawRewardResponseMeta
	PortalRedeemFromLiquidationPoolMeta         = metadataCommon.PortalRedeemFromLiquidationPoolMeta
	PortalRedeemFromLiquidationPoolResponseMeta = metadataCommon.PortalRedeemFromLiquidationPoolResponseMeta
	PortalCustodianTopupMeta                    = metadataCommon.PortalCustodianTopupMeta
	PortalCustodianTopupResponseMeta            = metadataCommon.PortalCustodianTopupResponseMeta
	PortalTotalRewardCustodianMeta              = metadataCommon.PortalTotalRewardCustodianMeta
	PortalPortingResponseMeta                   = metadataCommon.PortalPortingResponseMeta
	PortalReqMatchingRedeemMeta                 = metadataCommon.PortalReqMatchingRedeemMeta
	PortalPickMoreCustodianForRedeemMeta        = metadataCommon.PortalPickMoreCustodianForRedeemMeta
	PortalCustodianTopupMetaV2                  = metadataCommon.PortalCustodianTopupMetaV2
	PortalCustodianTopupResponseMetaV2          = metadataCommon.PortalCustodianTopupResponseMetaV2
	// Portal v3
	PortalCustodianDepositMetaV3                  = metadataCommon.PortalCustodianDepositMetaV3
	PortalCustodianWithdrawRequestMetaV3          = metadataCommon.PortalCustodianWithdrawRequestMetaV3
	PortalRewardMetaV3                            = metadataCommon.PortalRewardMetaV3
	PortalRequestUnlockCollateralMetaV3           = metadataCommon.PortalRequestUnlockCollateralMetaV3
	PortalLiquidateCustodianMetaV3                = metadataCommon.PortalLiquidateCustodianMetaV3
	PortalLiquidateByRatesMetaV3                  = metadataCommon.PortalLiquidateByRatesMetaV3
	PortalRedeemFromLiquidationPoolMetaV3         = metadataCommon.PortalRedeemFromLiquidationPoolMetaV3
	PortalRedeemFromLiquidationPoolResponseMetaV3 = metadataCommon.PortalRedeemFromLiquidationPoolResponseMetaV3
	PortalCustodianTopupMetaV3                    = metadataCommon.PortalCustodianTopupMetaV3
	PortalTopUpWaitingPortingRequestMetaV3        = metadataCommon.PortalTopUpWaitingPortingRequestMetaV3
	PortalRequestPortingMetaV3                    = metadataCommon.PortalRequestPortingMetaV3
	PortalRedeemRequestMetaV3                     = metadataCommon.PortalRedeemRequestMetaV3
	PortalUnlockOverRateCollateralsMeta           = metadataCommon.PortalUnlockOverRateCollateralsMeta
	// Incognito => Ethereum's SC for portal
	PortalCustodianWithdrawConfirmMetaV3         = metadataCommon.PortalCustodianWithdrawConfirmMetaV3
	PortalRedeemFromLiquidationPoolConfirmMetaV3 = metadataCommon.PortalRedeemFromLiquidationPoolConfirmMetaV3
	PortalLiquidateRunAwayCustodianConfirmMetaV3 = metadataCommon.PortalLiquidateRunAwayCustodianConfirmMetaV3
	//Note: don't use this metadata type for others
	PortalResetPortalDBMeta = metadataCommon.PortalResetPortalDBMeta
	// relaying
	RelayingBNBHeaderMeta                 = metadataCommon.RelayingBNBHeaderMeta
	RelayingBTCHeaderMeta                 = metadataCommon.RelayingBTCHeaderMeta
	PortalTopUpWaitingPortingRequestMeta  = metadataCommon.PortalTopUpWaitingPortingRequestMeta
	PortalTopUpWaitingPortingResponseMeta = metadataCommon.PortalTopUpWaitingPortingResponseMeta
	// incognito mode for smart contract
	BurningForDepositToSCRequestMeta   = metadataCommon.BurningForDepositToSCRequestMeta
	BurningForDepositToSCRequestMetaV2 = metadataCommon.BurningForDepositToSCRequestMetaV2
	BurningConfirmForDepositToSCMeta   = metadataCommon.BurningConfirmForDepositToSCMeta
	BurningConfirmForDepositToSCMetaV2 = metadataCommon.BurningConfirmForDepositToSCMetaV2
	InitTokenRequestMeta               = metadataCommon.InitTokenRequestMeta
	InitTokenResponseMeta              = metadataCommon.InitTokenResponseMeta
	// incognito mode for bsc
	IssuingBSCRequestMeta    = metadataCommon.IssuingBSCRequestMeta
	IssuingBSCResponseMeta   = metadataCommon.IssuingBSCResponseMeta
	BurningPBSCRequestMeta   = metadataCommon.BurningPBSCRequestMeta
	BurningBSCConfirmMeta    = metadataCommon.BurningBSCConfirmMeta
	AllShards                = metadataCommon.AllShards
	BeaconOnly               = metadataCommon.BeaconOnly
	StopAutoStakingAmount    = metadataCommon.StopAutoStakingAmount
	EVMConfirmationBlocks    = metadataCommon.EVMConfirmationBlocks
	NoAction                 = metadataCommon.NoAction
	MetaRequestBeaconMintTxs = metadataCommon.MetaRequestBeaconMintTxs
	MetaRequestShardMintTxs  = metadataCommon.MetaRequestShardMintTxs
)
