package bridge

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	metadataCommon "github.com/incognitochain/go-incognito-sdk-v2/metadata/common"
)

type Vault struct {
	ExternalDecimal uint8  `mapstructure:"external_decimal"`
	ExternalTokenID string `mapstructure:"external_token_id"`
	NetworkID       uint8  `mapstructure:"network_id"`
}

type AddToken struct {
	NewListTokens map[common.Hash]map[common.Hash]Vault `json:"NewListTokens"`
}

func (a *AddToken) StringSlice() ([]string, error) {
	contentBytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return []string{
		strconv.Itoa(metadataCommon.BridgeAggAddTokenMeta),
		base64.StdEncoding.EncodeToString(contentBytes),
	}, nil
}

func (a *AddToken) FromStringSlice(source []string) error {
	if len(source) != 2 {
		return fmt.Errorf("len of instruction need to be 2 but get %d", len(source))
	}
	if strconv.Itoa(metadataCommon.BridgeAggAddTokenMeta) != source[0] {
		return fmt.Errorf("metaType need to be %d but get %s", metadataCommon.BridgeAggAddTokenMeta, source[0])
	}
	contentBytes, err := base64.StdEncoding.DecodeString(source[1])
	if err != nil {
		return err
	}
	err = json.Unmarshal(contentBytes, &a)
	if err != nil {
		return err
	}
	return nil
}
