package incclient

// MainNet config
const (
	MainNetETHContractAddressStr = "0x43D037A562099A4C2c95b1E2120cc43054450629"
	MainNetBSCContractAddressStr = "0x43D037A562099A4C2c95b1E2120cc43054450629"
	MainNetFullNode              = "https://beta-fullnode.incognito.org/fullnode"
	MainNetETHHost               = "https://mainnet.infura.io/v3/34918000975d4374a056ed78fe21c517"
	MainNetBSCHost               = "https://bsc-dataseed.binance.org"
	MainNetPrivacyVersion        = 2
)

// TestNet config
const (
	TestNetETHContractAddressStr = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	TestNetBSCContractAddressStr = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	TestNetFullNode              = "https://testnet.incognito.org/fullnode"
	TestNetETHHost               = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	TestNetBSCHost               = "https://data-seed-prebsc-2-s1.binance.org:8545"
	TestNetPrivacyVersion        = 2
)

// TestNet1 config
const (
	TestNet1ETHContractAddressStr = "0xE0D5e7217c6C4bc475404b26d763fAD3F14D2b86"
	TestNet1BSCContractAddressStr = "0x1ce57B254DC2DBB41e1aeA296Dc7dBD6fb549241"
	TestNet1FullNode              = "https://testnet1.incognito.org/fullnode"
	TestNet1ETHHost               = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	TestNet1BSCHost               = "https://data-seed-prebsc-2-s1.binance.org:8545"
	TestNet1PrivacyVersion        = 2
)

// Local config
const (
	LocalETHContractAddressStr = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	LocalFullNode              = "http://127.0.0.1:8334"
	LocalETHHost               = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	LocalPrivacyVersion        = 2
)

const (
	DefaultPRVFee         = uint64(100)
	MaxInputSize          = 30
	MaxOutputSize         = 30
	prvInCoinKey          = "PRVInputCoins"
	tokenInCoinKey        = "TokenInputCoins"
	defaultCacheDirectory = ".cache"
)
