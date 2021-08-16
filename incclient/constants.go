package incclient

// MainNet config
const (
	MainNetETHContractAddressStr = "0x43D037A562099A4C2c95b1E2120cc43054450629"
	MainNetFullNode              = "https://beta-fullnode.incognito.org/fullnode"
	MainNetETHHost               = "//https://mainnet.infura.io/v3/34918000975d4374a056ed78fe21c517"
	MainNetPrivacyVersion        = 2
)

// TestNet config
const (
	TestNetETHContractAddressStr = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	TestNetFullNode              = "https://testnet.incognito.org/fullnode"
	TestNetETHHost               = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	TestNetPrivacyVersion        = 2
)

// TestNet1 config
const (
	TestNet1ETHContractAddressStr = "0xa63705AA35Ca2F3273d44B275252332750c6B8B4"
	TestNet1FullNode              = "https://testnet1.incognito.org/fullnode"
	TestNet1ETHHost               = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	TestNet1PrivacyVersion        = 2
)

// Local config
const (
	LocalETHContractAddressStr = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	LocalFullNode              = "http://127.0.0.1:8334"
	LocalETHHost               = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	LocalPrivacyVersion        = 2
)

// DevNet config
const (
	DevNetETHContractAddressStr = "0x2f6F03F1b43Eab22f7952bd617A24AB46E970dF7"
	DevNetFullNode              = "http://139.162.55.124:8334"
	DevNetETHHost               = "https://kovan.infura.io/v3/93fe721349134964aa71071a713c5cef"
	DevNetPrivacyVersion        = 2
)

const (
	DefaultPRVFee         = uint64(100)
	MaxInputSize          = 30
	MaxOutputSize         = 30
	prvInCoinKey          = "PRVInputCoins"
	tokenInCoinKey        = "TokenInputCoins"
	defaultCacheDirectory = ".cache"
)
