package config

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"io/ioutil"
	"math"
	"os"
	"runtime"
)

// ClientConfig consists of all necessary configurations for a client.
type ClientConfig struct {
	// Network indicates the type of the Incognito network (mainnet, testnet, testnet1, local).
	Network string `json:"Network"`

	// RPCHost is the address of the remote Incognito full-node.
	RPCHost string `json:"RPCHost"`

	// Version is the version of the client.
	Version int `json:"Version"`

	// DefaultPRVFee is the default value for the Incognito transaction fee.
	DefaultPRVFee uint64 `json:"DefaultPRVFee"`

	// EVMNetworks consists of configurations for EVM networks.
	EVMNetworks map[string]EVMNetworkConfig `json:"EVMNetworks"`

	// UTXOCache consists of caching configurations.
	UTXOCache UTXOCacheConfig `json:"UTXOCache"`

	// PortalParams consists of the configurations of Portal V4.
	PortalParams *BTCPortalV4Params `json:"PortalParams,omitempty"`

	// LogConfig consist of the configurations of the logger.
	LogConfig *LoggerConfig `json:"LogConfig"`
}

// EVMNetworkConfig consists of all EVM configurations.
type EVMNetworkConfig struct {
	ContractAddress    string `json:"ContractAddress"`
	PRVContractAddress string `json:"PRVContractAddress,omitempty"`
	FullNodeHost       string `json:"FullNodeHost"`
}

// UTXOCacheConfig specifies configurations for the UTXO cache system.
type UTXOCacheConfig struct {
	Enable            bool   `json:"Enable"`
	MaxGetCoinThreads int    `json:"MaxGetCoinThreads"`
	CacheLocation     string `json:"CacheLocation"`
}

// LoggerConfig specifies configurations for the logger.
type LoggerConfig struct {
	Enable bool `json:"Enable"`
}

// MainNetConfig is the default configuration for interacting with the Incognito main-net.
var MainNetConfig = &ClientConfig{
	RPCHost:       "https://beta-fullnode.incognito.org/fullnode",
	Version:       2,
	DefaultPRVFee: 100,
	EVMNetworks: map[string]EVMNetworkConfig{
		"ETH": {
			ContractAddress:    "0x43D037A562099A4C2c95b1E2120cc43054450629",
			PRVContractAddress: "0xB64fde8f199F073F41c132B9eC7aD5b61De0B1B7",
			FullNodeHost:       "https://mainnet.infura.io/v3/34918000975d4374a056ed78fe21c517",
		},
		"BSC": {
			ContractAddress:    "0x43D037A562099A4C2c95b1E2120cc43054450629",
			PRVContractAddress: "0xB64fde8f199F073F41c132B9eC7aD5b61De0B1B7",
			FullNodeHost:       "https://bsc-dataseed.binance.org",
		},
		"PLG": {
			ContractAddress: "0x43D037A562099A4C2c95b1E2120cc43054450629",
			FullNodeHost:    "https://matic-mainnet.chainstacklabs.com",
		},
		"FTM": {
			ContractAddress: "0x43D037A562099A4C2c95b1E2120cc43054450629",
			FullNodeHost:    "https://rpc.ftm.tools/",
		},
	},
	UTXOCache: UTXOCacheConfig{
		Enable:            false,
		MaxGetCoinThreads: int(math.Max(float64(runtime.NumCPU()), 4)),
		CacheLocation:     ".cache",
	},
	PortalParams: &mainNetBTCPortalV4Params,
	LogConfig:    &LoggerConfig{Enable: false},
}

// TestNetConfig is the default configuration for interacting with the Incognito test-net.
var TestNetConfig = &ClientConfig{
	RPCHost:       "https://testnet.incognito.org/fullnode",
	Version:       2,
	DefaultPRVFee: 100,
	EVMNetworks: map[string]EVMNetworkConfig{
		"ETH": {
			ContractAddress:    "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7",
			PRVContractAddress: "0xaE61fEFD69BacF3951F1C86c9A4D3F006810Ac21",
			FullNodeHost:       "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef",
		},
		"BSC": {
			ContractAddress:    "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7",
			PRVContractAddress: "0xB49E8844a72CF1ce885aEf13F82BeeAEEFc01527",
			FullNodeHost:       "https://data-seed-prebsc-2-s2.binance.org:8545",
		},
		"PLG": {
			ContractAddress: "0x4fF5c88cD1FD773436C2aBcFE175fe4ba6a2eB68",
			FullNodeHost:    "https://matic-mumbai.chainstacklabs.com",
		},
		"FTM": {
			ContractAddress: "0x9cb4baf1b60DaBB6B22BcFf07cc0e10395423aed",
			FullNodeHost:    "https://rpc.testnet.fantom.network/",
		},
	},
	UTXOCache: UTXOCacheConfig{
		Enable:            false,
		MaxGetCoinThreads: int(math.Max(float64(runtime.NumCPU()), 4)),
		CacheLocation:     ".cache",
	},
	PortalParams: &testNetBTCPortalV4Params,
	LogConfig:    &LoggerConfig{Enable: false},
}

// TestNet1Config is the default configuration for interacting with the Incognito test-net1.
var TestNet1Config = &ClientConfig{
	RPCHost:       "https://testnet1.incognito.org/fullnode",
	Version:       2,
	DefaultPRVFee: 100,
	EVMNetworks: map[string]EVMNetworkConfig{
		"ETH": {
			ContractAddress:    "0xE0D5e7217c6C4bc475404b26d763fAD3F14D2b86",
			PRVContractAddress: "0x917637E3E1ee531231747690189e22C5FA38D88C",
			FullNodeHost:       "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef",
		},
		"BSC": {
			ContractAddress:    "0x1ce57B254DC2DBB41e1aeA296Dc7dBD6fb549241",
			PRVContractAddress: "0x3d2E0c1b0b2d81D1A76544E6cC08670af9b86531",
			FullNodeHost:       "https://data-seed-prebsc-2-s1.binance.org:8545",
		},
		"PLG": {
			ContractAddress: "0xaCc76d988a1Ad9322069d9999D04b49b85E77a99",
			FullNodeHost:    "https://polygon-mumbai.g.alchemy.com/v2/CBQ1SQRLf3eQGbTXd_aA3LU7hvmwR33K",
		},
		"FTM": {
			ContractAddress: "0x9cb4baf1b60DaBB6B22BcFf07cc0e10395423aed",
			FullNodeHost:    "https://rpc.testnet.fantom.network/",
		},
	},
	UTXOCache: UTXOCacheConfig{
		Enable:            false,
		MaxGetCoinThreads: int(math.Max(float64(runtime.NumCPU()), 4)),
		CacheLocation:     ".cache",
	},
	PortalParams: &testNet1BTCPortalV4Params,
	LogConfig:    &LoggerConfig{Enable: false},
}

// LocalConfig is the default configuration for interacting with the Incognito local-net.
var LocalConfig = &ClientConfig{
	RPCHost:       "http://127.0.0.1:8334",
	Version:       2,
	DefaultPRVFee: 100,
	EVMNetworks: map[string]EVMNetworkConfig{
		"ETH": {
			ContractAddress:    "0xE0D5e7217c6C4bc475404b26d763fAD3F14D2b86",
			PRVContractAddress: "0x917637E3E1ee531231747690189e22C5FA38D88C",
			FullNodeHost:       "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef",
		},
	},
	UTXOCache: UTXOCacheConfig{
		Enable:            false,
		MaxGetCoinThreads: int(math.Max(float64(runtime.NumCPU()), 4)),
		CacheLocation:     ".cache",
	},
	PortalParams: &localBTCPortalV4Params,
	LogConfig:    &LoggerConfig{Enable: false},
}

// LoadConfig returns a ClientConfig from the given json config file.
func LoadConfig(configFile string) (*ClientConfig, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(configFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			fmt.Println("Cannot load config:", err)
		}
	}()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result ClientConfig
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, err
	}

	switch result.Network {
	case "mainnet":
		if result.EVMNetworks == nil {
			result.EVMNetworks = MainNetConfig.EVMNetworks
		}
		if result.PortalParams.ChainParams == nil {
			result.PortalParams = &BTCPortalV4Params{
				NetworkName:       result.PortalParams.NetworkName,
				MasterPubKeys:     mainNetBTCPortalV4Params.MasterPubKeys,
				NumRequiredSigs:   mainNetBTCPortalV4Params.NumRequiredSigs,
				MinUnShieldAmount: mainNetBTCPortalV4Params.MinUnShieldAmount,
				ChainParams:       mainNetBTCPortalV4Params.ChainParams,
				TokenID:           result.PortalParams.TokenID,
			}
		}
	case "testnet":
		if result.EVMNetworks == nil {
			result.EVMNetworks = TestNetConfig.EVMNetworks
		}
		if result.PortalParams.ChainParams == nil {
			result.PortalParams = &BTCPortalV4Params{
				NetworkName:       result.PortalParams.NetworkName,
				MasterPubKeys:     testNetBTCPortalV4Params.MasterPubKeys,
				NumRequiredSigs:   testNetBTCPortalV4Params.NumRequiredSigs,
				MinUnShieldAmount: testNetBTCPortalV4Params.MinUnShieldAmount,
				ChainParams:       testNetBTCPortalV4Params.ChainParams,
				TokenID:           result.PortalParams.TokenID,
			}
		}
	case "testnet1":
		if result.EVMNetworks == nil {
			result.EVMNetworks = TestNet1Config.EVMNetworks
		}
		if result.PortalParams.ChainParams == nil {
			result.PortalParams = &BTCPortalV4Params{
				NetworkName:       result.PortalParams.NetworkName,
				MasterPubKeys:     testNet1BTCPortalV4Params.MasterPubKeys,
				NumRequiredSigs:   testNet1BTCPortalV4Params.NumRequiredSigs,
				MinUnShieldAmount: testNet1BTCPortalV4Params.MinUnShieldAmount,
				ChainParams:       testNet1BTCPortalV4Params.ChainParams,
				TokenID:           result.PortalParams.TokenID,
			}
		}
	case "local":
		if result.EVMNetworks == nil {
			result.EVMNetworks = LocalConfig.EVMNetworks
		}
		if result.PortalParams.ChainParams == nil {
			result.PortalParams = &BTCPortalV4Params{
				NetworkName:       result.PortalParams.NetworkName,
				MasterPubKeys:     localBTCPortalV4Params.MasterPubKeys,
				NumRequiredSigs:   localBTCPortalV4Params.NumRequiredSigs,
				MinUnShieldAmount: localBTCPortalV4Params.MinUnShieldAmount,
				ChainParams:       localBTCPortalV4Params.ChainParams,
				TokenID:           result.PortalParams.TokenID,
			}
		}
	case "custom":
		if result.PortalParams.ChainParams == nil {
			if result.PortalParams.MasterPubKeys == nil {
				return nil, fmt.Errorf("portal MasterPubKeys not found from the config file")
			}
			result.PortalParams.ChainParams = &chaincfg.MainNetParams
			switch result.PortalParams.NetworkName {
			case "testnet3":
				result.PortalParams.ChainParams = &chaincfg.TestNet3Params
			case "regtest":
				result.PortalParams.ChainParams = &chaincfg.RegressionNetParams
			case "simnet":
				result.PortalParams.ChainParams = &chaincfg.SimNetParams
			case "mainnet":
			default:
				return nil, fmt.Errorf("BTC network `%v` not found", result.PortalParams.NetworkName)
			}
		}
	default:
		return nil, fmt.Errorf("network %v not supported", result.Network)
	}

	return &result, nil
}
