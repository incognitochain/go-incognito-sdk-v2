package metadata

const (
	InvalidMeta = 1

	IssuingRequestMeta     = 24
	IssuingResponseMeta    = 25
	ContractingRequestMeta = 26
	IssuingETHRequestMeta  = 80
	IssuingETHResponseMeta = 81

	ShardBlockReward             = 36
	AcceptedBlockRewardInfoMeta  = 37
	ShardBlockSalaryResponseMeta = 38
	BeaconRewardRequestMeta      = 39
	BeaconSalaryResponseMeta     = 40
	ReturnStakingMeta            = 41
	IncDAORewardRequestMeta      = 42
	ShardBlockRewardRequestMeta  = 43
	WithDrawRewardRequestMeta    = 44
	WithDrawRewardResponseMeta   = 45

	ShardStakingMeta    = 63
	StopAutoStakingMeta = 127
	BeaconStakingMeta   = 64
	UnStakingMeta       = 210

	BeaconSwapConfirmMeta = 70
	BridgeSwapConfirmMeta = 71
	BurningRequestMeta    = 27
	BurningRequestMetaV2  = 240
	BurningConfirmMeta    = 72
	BurningConfirmMetaV2  = 241

	PDEContributionMeta                   = 90
	PDETradeRequestMeta                   = 91
	PDETradeResponseMeta                  = 92
	PDEWithdrawalRequestMeta              = 93
	PDEWithdrawalResponseMeta             = 94
	PDEContributionResponseMeta           = 95
	PDEPRVRequiredContributionRequestMeta = 204
	PDECrossPoolTradeRequestMeta          = 205
	PDECrossPoolTradeResponseMeta         = 206
	PDEFeeWithdrawalRequestMeta           = 207
	PDEFeeWithdrawalResponseMeta          = 208
	PDETradingFeesDistributionMeta        = 209

	PortalCustodianDepositMeta                  = 100
	PortalRequestPortingMeta                    = 101
	PortalUserRequestPTokenMeta                 = 102
	PortalCustodianDepositResponseMeta          = 103
	PortalUserRequestPTokenResponseMeta         = 104
	PortalExchangeRatesMeta                     = 105
	PortalRedeemRequestMeta                     = 106
	PortalRedeemRequestResponseMeta             = 107
	PortalRequestUnlockCollateralMeta           = 108
	PortalCustodianWithdrawRequestMeta          = 110
	PortalCustodianWithdrawResponseMeta         = 111
	PortalLiquidateCustodianMeta                = 112
	PortalLiquidateCustodianResponseMeta        = 113
	PortalLiquidateTPExchangeRatesMeta          = 114
	PortalExpiredWaitingPortingReqMeta          = 116
	PortalRewardMeta                            = 117
	PortalRequestWithdrawRewardMeta             = 118
	PortalRequestWithdrawRewardResponseMeta     = 119
	PortalRedeemFromLiquidationPoolMeta         = 120
	PortalRedeemFromLiquidationPoolResponseMeta = 121
	PortalCustodianTopUpMeta                    = 122
	PortalCustodianTopUpResponseMeta            = 123
	PortalTotalRewardCustodianMeta              = 124
	PortalPortingResponseMeta                   = 125
	PortalReqMatchingRedeemMeta                 = 126
	PortalPickMoreCustodianForRedeemMeta        = 128
	PortalCustodianTopUpMetaV2                  = 129
	PortalCustodianTopUpResponseMetaV2          = 130

	PortalCustodianDepositMetaV3                  = 131
	PortalCustodianWithdrawRequestMetaV3          = 132
	PortalRewardMetaV3                            = 133
	PortalRequestUnlockCollateralMetaV3           = 134
	PortalLiquidateCustodianMetaV3                = 135
	PortalLiquidateByRatesMetaV3                  = 136
	PortalRedeemFromLiquidationPoolMetaV3         = 137
	PortalRedeemFromLiquidationPoolResponseMetaV3 = 138
	PortalCustodianTopUpMetaV3                    = 139
	PortalTopUpWaitingPortingRequestMetaV3        = 140
	PortalRequestPortingMetaV3                    = 141
	PortalRedeemRequestMetaV3                     = 142
	PortalUnlockOverRateCollateralsMeta           = 143

	PortalV4ShieldingRequestMeta      = 1250
	PortalV4UnShieldingRequestMeta    = 1252
	PortalV4FeeReplacementRequestMeta = 1255
	PortalV4SubmitConfirmedTxMeta     = 1256
	PortalV4ConvertVaultRequestMeta   = 1257

	PortalCustodianWithdrawConfirmMetaV3         = 170
	PortalRedeemFromLiquidationPoolConfirmMetaV3 = 171
	PortalLiquidateRunAwayCustodianConfirmMetaV3 = 172

	PortalResetPortalDBMeta = 199

	RelayingBNBHeaderMeta = 200
	RelayingBTCHeaderMeta = 201

	PortalTopUpWaitingPortingRequestMeta  = 202
	PortalTopUpWaitingPortingResponseMeta = 203

	BurningForDepositToSCRequestMeta   = 96
	BurningForDepositToSCRequestMetaV2 = 242
	BurningConfirmForDepositToSCMeta   = 97
	BurningConfirmForDepositToSCMetaV2 = 243

	InitTokenRequestMeta  = 244
	InitTokenResponseMeta = 245

	IssuingBSCRequestMeta  = 250
	IssuingBSCResponseMeta = 251
	BurningPBSCRequestMeta = 252
	BurningBSCConfirmMeta  = 253
)

var minerCreatedMetaTypes = []int{
	ShardBlockReward,
	BeaconSalaryResponseMeta,
	IssuingResponseMeta,
	IssuingETHResponseMeta,
	ReturnStakingMeta,
	WithDrawRewardResponseMeta,
	PDETradeResponseMeta,
	PDECrossPoolTradeResponseMeta,
	PDEWithdrawalResponseMeta,
	PDEFeeWithdrawalResponseMeta,
	PDEContributionResponseMeta,
	PortalUserRequestPTokenResponseMeta,
	PortalCustodianDepositResponseMeta,
	PortalRedeemRequestResponseMeta,
	PortalCustodianWithdrawResponseMeta,
	PortalLiquidateCustodianResponseMeta,
	PortalRequestWithdrawRewardResponseMeta,
	PortalRedeemFromLiquidationPoolResponseMeta,
	PortalCustodianTopUpResponseMeta,
	PortalCustodianTopUpResponseMetaV2,
	PortalPortingResponseMeta,
	PortalTopUpWaitingPortingResponseMeta,
	PortalRedeemFromLiquidationPoolResponseMetaV3,
	InitTokenResponseMeta,
}

var AcceptedWithdrawRewardRequestVersion = []int{0, 1}
