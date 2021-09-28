// Package incclient provides access to almost all functions needed to create transactions, become a node validator,
// retrieve information from full-nodes, shield or un-shield access, etc. It is the main focus of this sdk.
package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"strings"
)

// IncClient defines the environment with which users want to interact.
type IncClient struct {
	// the Incognito-RPC server
	rpcServer *rpc.RPCServer

	// the Ethereum-RPC server
	ethServer *rpc.RPCServer

	// the BSC-RPC server
	bscServer *rpc.RPCServer

	// the parameters used in the v4 portal for BTC
	btcPortalParams *BTCPortalV4Params

	// the version of the client
	version int

	// the utxoCache of the client
	cache *utxoCache
}

// NewTestNetClient creates a new IncClient with the test-net environment.
func NewTestNetClient() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(TestNetFullNode)
	ethServer := rpc.NewRPCServer(TestNetETHHost)
	bscServer := rpc.NewRPCServer(TestNetBSCHost)

	incClient := IncClient{
		rpcServer:       rpcServer,
		ethServer:       ethServer,
		bscServer:       bscServer,
		btcPortalParams: &testNetBTCPortalV4Params,
		version:         TestNetPrivacyVersion,
	}

	activeShards, err := incClient.GetActiveShard()
	if err != nil {
		return nil, err
	}

	Logger.Printf("Init to %v, activeShards: %v\n", TestNetFullNode, activeShards)

	common.MaxShardNumber = activeShards
	if incClient.version == 1 {
		common.AddressVersion = 0
	} else if incClient.version == 2 {
		common.AddressVersion = 1
	}

	return &incClient, nil
}

// NewTestNetClientWithCache creates a new IncClient with the test-net environment.
// It also creates a cache instance for locally saving UTXOs.
func NewTestNetClientWithCache() (*IncClient, error) {
	incClient, err := NewTestNetClient()
	if err != nil {
		return nil, err
	}

	incClient.cache, err = newUTXOCache(fmt.Sprintf("%v/%v", defaultCacheDirectory, "testnet"))
	if err != nil {
		return nil, err
	}
	incClient.cache.start()
	rawAssetTags, err = incClient.GetAllAssetTags()
	if err != nil {
		Logger.Printf("Cannot get raw asset tags: %v\n", err)
	}

	return incClient, nil
}

// NewTestNet1Client creates a new IncClient with the test-net 1 environment.
func NewTestNet1Client() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(TestNet1FullNode)
	ethServer := rpc.NewRPCServer(TestNet1ETHHost)
	bscServer := rpc.NewRPCServer(TestNet1BSCHost)

	incClient := IncClient{
		rpcServer:       rpcServer,
		ethServer:       ethServer,
		bscServer:       bscServer,
		btcPortalParams: &testNet1BTCPortalV4Params,
		version:         TestNet1PrivacyVersion}

	activeShards, err := incClient.GetActiveShard()
	if err != nil {
		return nil, err
	}

	Logger.Printf("Init to %v, activeShards: %v\n", TestNet1FullNode, activeShards)

	common.MaxShardNumber = activeShards
	if incClient.version == 1 {
		common.AddressVersion = 0
	} else if incClient.version == 2 {
		common.AddressVersion = 1
	}

	return &incClient, nil
}

// NewTestNet1ClientWithCache creates a new IncClient with the test-net-1 environment.
// It also creates a cache instance for locally saving UTXOs.
func NewTestNet1ClientWithCache() (*IncClient, error) {
	incClient, err := NewTestNetClient()
	if err != nil {
		return nil, err
	}

	incClient.cache, err = newUTXOCache(fmt.Sprintf("%v/%v", defaultCacheDirectory, "testnet1"))
	if err != nil {
		return nil, err
	}
	incClient.cache.start()
	rawAssetTags, err = incClient.GetAllAssetTags()
	if err != nil {
		Logger.Printf("Cannot get raw asset tags: %v\n", err)
	}

	return incClient, nil
}

// NewMainNetClient creates a new IncClient with the main-net environment.
func NewMainNetClient() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(MainNetFullNode)
	ethServer := rpc.NewRPCServer(MainNetETHHost)
	bscServer := rpc.NewRPCServer(MainNetBSCHost)

	incClient := IncClient{
		rpcServer:       rpcServer,
		ethServer:       ethServer,
		bscServer:       bscServer,
		btcPortalParams: &mainNetBTCPortalV4Params,
		version:         MainNetPrivacyVersion}

	activeShards, err := incClient.GetActiveShard()
	if err != nil {
		return nil, err
	}

	Logger.Printf("Init to %v, activeShards: %v\n", MainNetFullNode, activeShards)

	common.MaxShardNumber = activeShards
	if incClient.version == 1 {
		common.AddressVersion = 0
	} else if incClient.version == 2 {
		common.AddressVersion = 1
	}

	return &incClient, nil
}

// NewMainNetClientWithCache creates a new IncClient with the main-net environment.
// It also creates a cache instance for locally saving UTXOs.
func NewMainNetClientWithCache() (*IncClient, error) {
	incClient, err := NewMainNetClient()
	if err != nil {
		return nil, err
	}

	incClient.cache, err = newUTXOCache(fmt.Sprintf("%v/%v", defaultCacheDirectory, "mainnet"))
	if err != nil {
		return nil, err
	}
	incClient.cache.start()
	rawAssetTags, err = incClient.GetAllAssetTags()
	if err != nil {
		Logger.Printf("Cannot get raw asset tags: %v\n", err)
	}

	return incClient, nil
}

// NewLocalClient creates a new IncClient with the local environment.
func NewLocalClient(port string) (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(LocalFullNode)
	ethServer := rpc.NewRPCServer(LocalETHHost)

	incClient := IncClient{
		rpcServer:       rpcServer,
		ethServer:       ethServer,
		btcPortalParams: &localBTCPortalV4Params,
		version:         LocalPrivacyVersion}
	if port != "" {
		incClient.rpcServer = rpc.NewRPCServer(fmt.Sprintf("http://127.0.0.1:%v", port))
	}

	activeShards, err := incClient.GetActiveShard()
	if err != nil {
		return nil, err
	}

	Logger.Printf("Init to %v, activeShards: %v\n", LocalFullNode, activeShards)

	common.MaxShardNumber = activeShards
	if incClient.version == 1 {
		common.AddressVersion = 0
	} else if incClient.version == 2 {
		common.AddressVersion = 1
	}

	return &incClient, nil
}

// NewLocalClientWithCache creates a new IncClient with the local environment.
// It also creates a cache instance for locally saving UTXOs.
func NewLocalClientWithCache() (*IncClient, error) {
	incClient, err := NewLocalClient("")
	if err != nil {
		return nil, err
	}

	incClient.cache, err = newUTXOCache(fmt.Sprintf("%v/%v", defaultCacheDirectory, "local"))
	if err != nil {
		return nil, err
	}
	incClient.cache.start()
	rawAssetTags, err = incClient.GetAllAssetTags()
	if err != nil {
		Logger.Printf("Cannot get raw asset tags: %v\n", err)
	}

	return incClient, nil
}

// NewIncClient creates a new IncClient from given parameters.
//
// Specify which network the client is interacting with by the parameter `networks`.
// A valid network is one of the following: mainnet, testnet, testnet1, local. By default, this function will initialize
// a main-net client if no value is assigned to `networks`.
// Note that only the first value passed to `networks` is processed.
func NewIncClient(fullNode, ethNode string, version int, networks ...string) (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(fullNode)
	ethServer := rpc.NewRPCServer(ethNode)

	incClient := IncClient{
		rpcServer:       rpcServer,
		ethServer:       ethServer,
		bscServer:       rpc.NewRPCServer(MainNetBSCHost),
		btcPortalParams: &mainNetBTCPortalV4Params,
		version:         version,
	}
	if len(networks) > 0 {
		switch strings.ToLower(networks[0]) {
		case "testnet":
			incClient.btcPortalParams = &testNetBTCPortalV4Params
			incClient.bscServer = rpc.NewRPCServer(TestNetBSCHost)
		case "testnet1":
			incClient.btcPortalParams = &testNet1BTCPortalV4Params
			incClient.bscServer = rpc.NewRPCServer(TestNet1BSCHost)
		case "local":
			incClient.btcPortalParams = &localBTCPortalV4Params
		case "mainnet":
		default:
			return nil, fmt.Errorf("network %v not valid", networks[0])
		}
	}

	activeShards, err := incClient.GetActiveShard()
	if err != nil {
		return nil, err
	}

	Logger.Printf("Init to %v, activeShards: %v\n", fullNode, activeShards)

	common.MaxShardNumber = activeShards
	if incClient.version == 1 {
		common.AddressVersion = 0
	} else if incClient.version == 2 {
		common.AddressVersion = 1
	} else {
		return nil, fmt.Errorf("version %v not supported", version)
	}

	return &incClient, nil
}

// NewIncClientWithCache creates a new IncClient from given parameters.
// It also creates a cache instance for locally saving UTXOs.
//
// Specify which network the client is interacting with by the parameter `networks`.
// A valid network is one of the following: mainnet, testnet, testnet1, local. By default, this function will initialize
// a main-net client if no value is assigned to `networks`.
// Note that only the first value passed to `networks` is processed.
func NewIncClientWithCache(fullNode, ethNode string, version int, networks ...string) (*IncClient, error) {
	incClient, err := NewIncClient(fullNode, ethNode, version, networks...)
	if err != nil {
		return nil, err
	}

	incClient.cache, err = newUTXOCache(fmt.Sprintf("%v/%v", defaultCacheDirectory, "custom"))
	if err != nil {
		return nil, err
	}
	incClient.cache.start()
	rawAssetTags, err = incClient.GetAllAssetTags()
	if err != nil {
		Logger.Printf("Cannot get raw asset tags: %v\n", err)
	}

	return incClient, nil
}
