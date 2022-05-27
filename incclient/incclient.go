// Package incclient provides access to almost all functions needed to create transactions, become a node validator,
// retrieve information from full-nodes, shield or un-shield access, etc. It is the main focus of this sdk.
package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient/config"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"strings"
)

// IncClient defines the environment with which users want to interact.
type IncClient struct {
	// the Incognito-RPC server
	rpcServer *rpc.RPCServer

	// the EVM-RPC servers
	evmServers map[int]*rpc.RPCServer

	// the version of the client
	version int

	// cfg consists of all configurations of the IncClient.
	cfg *config.ClientConfig

	// the utxoCache of the client
	cache *utxoCache
}

// GetConfig returns the current configuration of the client.
func (client IncClient) GetConfig() config.ClientConfig {
	return *client.cfg
}

// NewTestNetClient creates a new IncClient with the test-net environment.
func NewTestNetClient() (*IncClient, error) {
	return NewIncClientWithConfig(config.TestNetConfig)
}

// NewTestNetClientWithCache creates a new IncClient with the test-net environment.
// It also creates a cache instance for locally saving UTXOs.
func NewTestNetClientWithCache() (*IncClient, error) {
	tmpConfig := *config.TestNetConfig
	tmpConfig.UTXOCache.Enable = true

	return NewIncClientWithConfig(&tmpConfig)
}

// NewTestNet1Client creates a new IncClient with the test-net 1 environment.
func NewTestNet1Client() (*IncClient, error) {
	return NewIncClientWithConfig(config.TestNet1Config)
}

// NewTestNet1ClientWithCache creates a new IncClient with the test-net-1 environment.
// It also creates a cache instance for locally saving UTXOs.
func NewTestNet1ClientWithCache() (*IncClient, error) {
	tmpConfig := *config.TestNet1Config
	tmpConfig.UTXOCache.Enable = true

	return NewIncClientWithConfig(&tmpConfig)
}

// NewMainNetClient creates a new IncClient with the main-net environment.
func NewMainNetClient() (*IncClient, error) {
	return NewIncClientWithConfig(config.MainNetConfig)
}

// NewMainNetClientWithCache creates a new IncClient with the main-net environment.
// It also creates a cache instance for locally saving UTXOs.
func NewMainNetClientWithCache() (*IncClient, error) {
	tmpConfig := *config.MainNetConfig
	tmpConfig.UTXOCache.Enable = true

	return NewIncClientWithConfig(&tmpConfig)
}

// NewLocalClient creates a new IncClient with the local environment.
func NewLocalClient(_ ...string) (*IncClient, error) {
	return NewIncClientWithConfig(config.LocalConfig)
}

// NewLocalClientWithCache creates a new IncClient with the local environment.
// It also creates a cache instance for locally saving UTXOs.
func NewLocalClientWithCache() (*IncClient, error) {
	tmpConfig := *config.LocalConfig
	tmpConfig.UTXOCache.Enable = true

	return NewIncClientWithConfig(&tmpConfig)
}

// NewIncClient creates a new IncClient from given parameters.
//
// Specify which network the client is interacting with by the parameter `networks`.
// A valid network is one of the following: mainnet, testnet, testnet1, local. By default, this function will initialize
// a main-net client if no value is assigned to `networks`.
// Note that only the first value passed to `networks` is processed.
//
// Deprecated: use NewIncClientWithConfig instead.
func NewIncClient(fullNode, ethNode string, version int, networks ...string) (*IncClient, error) {
	tmpConfig := *config.MainNetConfig
	var network = "mainnet"
	if len(networks) > 0 {
		network = strings.ToLower(networks[0])
	}
	switch network {
	case "testnet":
		tmpConfig = *config.TestNetConfig
	case "testnet1":
		tmpConfig = *config.TestNet1Config
	case "local":
		tmpConfig = *config.LocalConfig
	case "mainnet":
	default:
		return nil, fmt.Errorf("network %v not supported", network)
	}

	tmpConfig.RPCHost = fullNode
	tmpConfig.Version = version

	return NewIncClientWithConfig(&tmpConfig)
}

// NewIncClientWithCache creates a new IncClient from given parameters.
// It also creates a cache instance for locally saving UTXOs.
//
// Specify which network the client is interacting with by the parameter `networks`.
// A valid network is one of the following: mainnet, testnet, testnet1, local. By default, this function will initialize
// a main-net client if no value is assigned to `networks`.
// Note that only the first value passed to `networks` is processed.
func NewIncClientWithCache(fullNode, _ string, version int, networks ...string) (*IncClient, error) {
	tmpConfig := *config.MainNetConfig
	var network = "mainnet"
	if len(networks) > 0 {
		network = strings.ToLower(networks[0])
	}
	switch network {
	case "testnet":
		tmpConfig = *config.TestNetConfig
	case "testnet1":
		tmpConfig = *config.TestNet1Config
	case "local":
		tmpConfig = *config.LocalConfig
	case "mainnet":
	default:
		return nil, fmt.Errorf("network %v not supported", network)
	}

	tmpConfig.RPCHost = fullNode
	tmpConfig.Version = version
	tmpConfig.UTXOCache.Enable = true

	return NewIncClientWithConfig(&tmpConfig)
}

// NewIncClientWithConfig creates a new IncClient from the given ClientConfig.
func NewIncClientWithConfig(cfg *config.ClientConfig) (*IncClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("empty config")
	}

	// Init the main RPC server
	rpcServer := rpc.NewRPCServer(cfg.RPCHost)

	// Init the EVM hosts
	evmServers := make(map[int]*rpc.RPCServer)
	for evmNetwork, evmConfig := range cfg.EVMNetworks {
		var evmNetworkID int
		switch evmNetwork {
		case "ETH":
			evmNetworkID = rpc.ETHNetworkID
		case "BSC":
			evmNetworkID = rpc.BSCNetworkID
		case "PLG":
			evmNetworkID = rpc.PLGNetworkID
		case "FTM":
			evmNetworkID = rpc.FTMNetworkID
		default:
			return nil, fmt.Errorf("EVMNetwork %v not supported", evmNetwork)
		}

		evmServers[evmNetworkID] = rpc.NewRPCServer(evmConfig.FullNodeHost)
	}

	// Init the client
	incClient := &IncClient{
		rpcServer:  rpcServer,
		evmServers: evmServers,
		version:    TestNet1PrivacyVersion,
		cfg:        cfg,
	}

	// Init the logger
	if cfg.LogConfig.Enable {
		Logger.IsEnable = true
	}

	// Get the current number of shards
	activeShards, err := incClient.GetActiveShard()
	if err != nil {
		return nil, err
	}
	common.MaxShardNumber = activeShards
	if incClient.version == 1 {
		common.AddressVersion = 0
	} else if incClient.version == 2 {
		common.AddressVersion = 1
	}
	Logger.Printf("Init to %v, activeShards: %v\n", cfg.RPCHost, activeShards)

	// Init the cache (if needed)
	if cfg.UTXOCache.Enable {
		incClient.cache, err = newUTXOCache(fmt.Sprintf("%v/%v", cfg.UTXOCache.CacheLocation, cfg.Network))
		if err != nil {
			return nil, err
		}
		incClient.cache.start()
	}

	return incClient, nil
}
