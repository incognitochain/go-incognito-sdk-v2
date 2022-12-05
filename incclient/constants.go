package incclient

// MainNet config
const (
	MainNetETHContractAddressStr      = "0x43D037A562099A4C2c95b1E2120cc43054450629"
	MainNetBSCContractAddressStr      = "0x43D037A562099A4C2c95b1E2120cc43054450629"
	MainNetPLGContractAddressStr      = "0x43D037A562099A4C2c95b1E2120cc43054450629"
	MainNetFTMContractAddressStr      = "0x43D037A562099A4C2c95b1E2120cc43054450629"
	MainNetPRVERC20ContractAddressStr = "0xB64fde8f199F073F41c132B9eC7aD5b61De0B1B7"
	MainNetPRVBEP20ContractAddressStr = "0xB64fde8f199F073F41c132B9eC7aD5b61De0B1B7"
	MainNetFullNode                   = "https://beta-fullnode.incognito.org/fullnode"
	MainNetETHHost                    = "https://mainnet.infura.io/v3/34918000975d4374a056ed78fe21c517"
	MainNetBSCHost                    = "https://bsc-dataseed.binance.org"
	MainNetPLGHost                    = "https://matic-mainnet.chainstacklabs.com"
	MainNetFTMHost                    = "https://rpc.ftm.tools/"
	MainNetPrivacyVersion             = 2
)

// TestNet config
const (
	TestNetETHContractAddressStr      = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	TestNetBSCContractAddressStr      = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	TestNetPLGContractAddressStr      = "0x4fF5c88cD1FD773436C2aBcFE175fe4ba6a2eB68"
	TestNetFTMContractAddressStr      = "0x9cb4baf1b60DaBB6B22BcFf07cc0e10395423aed"
	TestNetPRVERC20ContractAddressStr = "0xaE61fEFD69BacF3951F1C86c9A4D3F006810Ac21"
	TestNetPRVBEP20ContractAddressStr = "0xB49E8844a72CF1ce885aEf13F82BeeAEEFc01527"
	TestNetFullNode                   = "https://testnet.incognito.org/fullnode"
	TestNetETHHost                    = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	TestNetBSCHost                    = "https://data-seed-prebsc-2-s2.binance.org:8545"
	TestNetPLGHost                    = "https://matic-mumbai.chainstacklabs.com"
	TestNetFTMHost                    = "https://rpc.testnet.fantom.network/"
	TestNetPrivacyVersion             = 2
)

// TestNet1 config
const (
	TestNet1ETHContractAddressStr      = "0xE0D5e7217c6C4bc475404b26d763fAD3F14D2b86"
	TestNet1BSCContractAddressStr      = "0x1ce57B254DC2DBB41e1aeA296Dc7dBD6fb549241"
	TestNet1PLGContractAddressStr      = "0xaCc76d988a1Ad9322069d9999D04b49b85E77a99"
	TestNet1FTMContractAddressStr      = "0x9cb4baf1b60DaBB6B22BcFf07cc0e10395423aed"
	TestNet1PRVERC20ContractAddressStr = "0x917637E3E1ee531231747690189e22C5FA38D88C"
	TestNet1PRVBEP20ContractAddressStr = "0x3d2E0c1b0b2d81D1A76544E6cC08670af9b86531"
	TestNet1FullNode                   = "https://testnet1.incognito.org/fullnode"
	TestNet1ETHHost                    = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	TestNet1BSCHost                    = "https://data-seed-prebsc-2-s1.binance.org:8545"
	TestNet1PLGHost                    = "https://polygon-mumbai.g.alchemy.com/v2/CBQ1SQRLf3eQGbTXd_aA3LU7hvmwR33K"
	TestNet1FTMHost                    = "https://rpc.testnet.fantom.network/"
	TestNet1PrivacyVersion             = 2
)

// Local config
const (
	LocalETHContractAddressStr = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	LocalFullNode              = "http://127.0.0.1:8334"
	LocalETHHost               = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	LocalPrivacyVersion        = 2
)

const (
	DefaultPRVFee            = uint64(100000000) // 0.1 PRV
	defaultNftRequiredAmount = 100
	MaxInputSize             = 30
	MaxOutputSize            = 30
	prvInCoinKey             = "PRVInputCoins"
	tokenInCoinKey           = "TokenInputCoins"
	defaultCacheDirectory    = ".cache"
)
