// Package incclient provides access to almost all functions needed to create transactions, become a node validator,
// retrieve information from full-nodes, shield or un-shield access, etc. It is the main focus of this sdk.
package incclient

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
)

// IncClient defines the environment with which users want to interact.
type IncClient struct {
	rpcServer *rpc.RPCServer
	ethServer *rpc.RPCServer
	version   int
}

// NewTestNetClient creates a new IncClient with the test-net environment.
func NewTestNetClient() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(TestNetFullNode)
	ethServer := rpc.NewRPCServer(TestNetETHHost)

	incClient := IncClient{rpcServer: rpcServer, ethServer: ethServer, version: TestNetPrivacyVersion}

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

// NewTestNet1Client creates a new IncClient with the test-net 1 environment.
func NewTestNet1Client() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(TestNet1FullNode)
	ethServer := rpc.NewRPCServer(TestNet1ETHHost)

	incClient := IncClient{rpcServer: rpcServer, ethServer: ethServer, version: TestNet1PrivacyVersion}

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

// NewMainNetClient creates a new IncClient with the main-net environment.
func NewMainNetClient() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(MainNetFullNode)
	ethServer := rpc.NewRPCServer(MainNetETHHost)

	incClient := IncClient{rpcServer: rpcServer, ethServer: ethServer, version: MainNetPrivacyVersion}

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

// NewLocalClient creates a new IncClient with the local environment.
func NewLocalClient(port string) (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(LocalFullNode)
	ethServer := rpc.NewRPCServer(LocalETHHost)

	incClient := IncClient{rpcServer: rpcServer, ethServer: ethServer, version: LocalPrivacyVersion}
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

// NewDevNetClient creates a new IncClient with the dev-net environment.
func NewDevNetClient() (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(DevNetFullNode)
	ethServer := rpc.NewRPCServer(DevNetETHHost)

	incClient := IncClient{rpcServer: rpcServer, ethServer: ethServer, version: DevNetPrivacyVersion}

	activeShards, err := incClient.GetActiveShard()
	if err != nil {
		return nil, err
	}

	Logger.Printf("Init to %v, activeShards: %v\n", DevNetFullNode, activeShards)

	common.MaxShardNumber = activeShards
	if incClient.version == 1 {
		common.AddressVersion = 0
	} else if incClient.version == 2 {
		common.AddressVersion = 1
	}

	return &incClient, nil
}

// NewIncClient creates a new IncClient from given parameters.
func NewIncClient(fullNode, ethNode string, version int) (*IncClient, error) {
	rpcServer := rpc.NewRPCServer(fullNode)
	ethServer := rpc.NewRPCServer(ethNode)

	incClient := IncClient{rpcServer: rpcServer, ethServer: ethServer, version: version}

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
