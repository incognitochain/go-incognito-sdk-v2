package metadata

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func calculateSize(meta Metadata) uint64 {
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return 0
	}
	return uint64(len(metaBytes))
}

func ParseMetadata(metaInBytes []byte) (Metadata, error) {
	if len(metaInBytes) == 0 {
		return nil, nil
	}

	mtTemp := map[string]interface{}{}
	err := json.Unmarshal(metaInBytes, &mtTemp)
	if err != nil {
		return nil, err
	}

	var md Metadata
	typeFloat, ok := mtTemp["Type"].(float64)
	if !ok {
		return nil, errors.Errorf("Could not parse metadata with type: %v", mtTemp["Type"])
	}
	theType := int(typeFloat)
	switch theType {
	case InitTokenRequestMeta:
		md = &InitTokenRequest{}
	case InitTokenResponseMeta:
		md = &InitTokenResponse{}
	case IssuingRequestMeta:
		md = &IssuingRequest{}
	case IssuingResponseMeta:
		md = &IssuingResponse{}
	case ContractingRequestMeta:
		md = &ContractingRequest{}
	case IssuingETHRequestMeta:
		md = &IssuingEVMRequest{}
	case IssuingETHResponseMeta:
		md = &IssuingEVMResponse{}
	case IssuingBSCRequestMeta:
		md = &IssuingEVMRequest{}
	case IssuingBSCResponseMeta:
		md = &IssuingEVMResponse{}
	case BurningRequestMeta:
		md = &BurningRequest{}
	case BurningRequestMetaV2:
		md = &BurningRequest{}
	case BurningPBSCRequestMeta:
		md = &BurningRequest{}
	case ShardStakingMeta:
		md = &StakingMetadata{}
	case BeaconStakingMeta:
		md = &StakingMetadata{}
	case ReturnStakingMeta:
		md = &ReturnStakingMetadata{}
	case WithDrawRewardRequestMeta:
		md = &WithDrawRewardRequest{}
	case WithDrawRewardResponseMeta:
		md = &WithDrawRewardResponse{}
	case UnStakingMeta:
		md = &UnStakingMetadata{}
	case StopAutoStakingMeta:
		md = &StopAutoStakingMetadata{}
	case PDEContributionMeta:
		md = &PDEContribution{}
	case PDEPRVRequiredContributionRequestMeta:
		md = &PDEContribution{}
	case PDETradeRequestMeta:
		md = &PDETradeRequest{}
	case PDETradeResponseMeta:
		md = &PDETradeResponse{}
	case PDECrossPoolTradeRequestMeta:
		md = &PDECrossPoolTradeRequest{}
	case PDECrossPoolTradeResponseMeta:
		md = &PDECrossPoolTradeResponse{}
	case PDEWithdrawalRequestMeta:
		md = &PDEWithdrawalRequest{}
	case PDEWithdrawalResponseMeta:
		md = &PDEWithdrawalResponse{}
	case PDEFeeWithdrawalRequestMeta:
		md = &PDEFeeWithdrawalRequest{}
	case PDEFeeWithdrawalResponseMeta:
		md = &PDEFeeWithdrawalResponse{}
	case PDEContributionResponseMeta:
		md = &PDEContributionResponse{}
	case BurningForDepositToSCRequestMeta:
		md = &BurningRequest{}
	case BurningForDepositToSCRequestMetaV2:
		md = &BurningRequest{}
	case RelayingBTCHeaderMeta:
		md = &RelayingHeader{}
	case PortalV4ShieldingRequestMeta:
		md = &PortalShieldingRequest{}
	case PortalV4UnShieldingRequestMeta:
		md = &PortalUnshieldRequest{}
	case PortalV4FeeReplacementRequestMeta:
		md = &PortalReplacementFeeRequest{}
	case PortalV4SubmitConfirmedTxMeta:
		md = &PortalSubmitConfirmedTxRequest{}
	case PortalV4ConvertVaultRequestMeta:
		md = &PortalConvertVaultRequest{}
	default:
		return nil, errors.Errorf("Could not parse metadata with type: %d", theType)
	}

	err = json.Unmarshal(metaInBytes, &md)
	if err != nil {
		return nil, err
	}

	return md, nil
}
