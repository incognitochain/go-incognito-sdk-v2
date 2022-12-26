// Package incclient provides access to almost all functions needed to create transactions, become a node validator,
// retrieve information from full-nodes, shield or un-shield access, etc. It is the main focus of this sdk.
package incclient

import (
	"fmt"
	"strings"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
)

// IncClient defines the environment with which users want to interact.
type IncClient struct {
	// the Incognito-RPC server
	rpcServer *rpc.RPCServer

	// the EVM-RPC servers
	evmServers map[int]*rpc.RPCServer

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
	evmServers := map[int]*rpc.RPCServer{
		rpc.ETHNetworkID:    rpc.NewRPCServer(TestNetETHHost),
		rpc.BSCNetworkID:    rpc.NewRPCServer(TestNetBSCHost),
		rpc.PLGNetworkID:    rpc.NewRPCServer(TestNetPLGHost),
		rpc.FTMNetworkID:    rpc.NewRPCServer(TestNetFTMHost),
		rpc.AURORANetworkID: rpc.NewRPCServer(TestNetAURORAHost),
		rpc.AVAXNetworkID:   rpc.NewRPCServer(TestNetAVAXHost),
	}

	incClient := IncClient{
		rpcServer:       rpcServer,
		evmServers:      evmServers,
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

	return incClient, nil
}

// NewTestNet1Client creates a new IncClient with the test-net 1 environment.
func NewTestNet1Client() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(TestNet1FullNode)
	evmServers := map[int]*rpc.RPCServer{
		rpc.ETHNetworkID:    rpc.NewRPCServer(TestNet1ETHHost),
		rpc.BSCNetworkID:    rpc.NewRPCServer(TestNet1BSCHost),
		rpc.PLGNetworkID:    rpc.NewRPCServer(TestNet1PLGHost),
		rpc.FTMNetworkID:    rpc.NewRPCServer(TestNet1FTMHost),
		rpc.AURORANetworkID: rpc.NewRPCServer(TestNet1AURORAHost),
		rpc.AVAXNetworkID:   rpc.NewRPCServer(TestNet1AVAXHost),
	}

	incClient := IncClient{
		rpcServer:       rpcServer,
		evmServers:      evmServers,
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

	return incClient, nil
}

// NewMainNetClient creates a new IncClient with the main-net environment.
func NewMainNetClient() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(MainNetFullNode)
	evmServers := map[int]*rpc.RPCServer{
		rpc.ETHNetworkID:    rpc.NewRPCServer(MainNetETHHost),
		rpc.BSCNetworkID:    rpc.NewRPCServer(MainNetBSCHost),
		rpc.PLGNetworkID:    rpc.NewRPCServer(MainNetPLGHost),
		rpc.FTMNetworkID:    rpc.NewRPCServer(MainNetFTMHost),
		rpc.AURORANetworkID: rpc.NewRPCServer(MainNetAURORAHost),
		rpc.AVAXNetworkID:   rpc.NewRPCServer(MainNetAVAXHost),
	}

	incClient := IncClient{
		rpcServer:       rpcServer,
		evmServers:      evmServers,
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

	return incClient, nil
}

// NewLocalClient creates a new IncClient with the local environment.
func NewLocalClient(port string) (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(LocalFullNode)
	evmServers := map[int]*rpc.RPCServer{
		rpc.ETHNetworkID: rpc.NewRPCServer(LocalETHHost),
	}

	incClient := IncClient{
		rpcServer:       rpcServer,
		evmServers:      evmServers,
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
	evmServers := map[int]*rpc.RPCServer{
		rpc.ETHNetworkID: rpc.NewRPCServer(ethNode),
		rpc.BSCNetworkID: rpc.NewRPCServer(MainNetBSCHost),
		rpc.PLGNetworkID: rpc.NewRPCServer(MainNetPLGHost),
	}

	incClient := IncClient{
		rpcServer:       rpcServer,
		evmServers:      evmServers,
		btcPortalParams: &mainNetBTCPortalV4Params,
		version:         version,
	}
	if len(networks) > 0 {
		switch strings.ToLower(networks[0]) {
		case "testnet":
			incClient.btcPortalParams = &testNetBTCPortalV4Params
			incClient.evmServers[rpc.BSCNetworkID] = rpc.NewRPCServer(TestNetBSCHost)
			incClient.evmServers[rpc.PLGNetworkID] = rpc.NewRPCServer(TestNetPLGHost)
		case "testnet1":
			incClient.btcPortalParams = &testNet1BTCPortalV4Params
			incClient.evmServers[rpc.BSCNetworkID] = rpc.NewRPCServer(TestNet1BSCHost)
			incClient.evmServers[rpc.PLGNetworkID] = rpc.NewRPCServer(TestNet1PLGHost)
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

	return incClient, nil
}
