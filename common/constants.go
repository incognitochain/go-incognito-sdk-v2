package common

// for common
const (
	EmptyString       = ""
	ZeroByte          = byte(0x00)
	DateOutputFormat  = "2006-01-02T15:04:05.999999"
	BigIntSize        = 32 // bytes
	CheckSumLen       = 4  // bytes
	AESKeySize        = 32 // bytes
	Int32Size         = 4  // bytes
	Uint32Size        = 4  // bytes
	Uint64Size        = 8  // bytes
	HashSize          = 32 // bytes
	MaxHashStringSize = HashSize * 2
	Base58Version     = 0
)

const (
	PrivateKeySize   = 32  // bytes
	PublicKeySize    = 32  // bytes
	BLSPublicKeySize = 128 // bytes
	BriPublicKeySize = 33  // bytes

	// for signature size
	// it is used for both privacy and no privacy
	SigPubKeySize    = 32
	SigNoPrivacySize = 64
	SigPrivacySize   = 96
)

// For all Transaction information
const (
	TxNormalType             = "n"   // normal tx
	TxRewardType             = "s"   // reward tx
	TxReturnStakingType      = "rs"  // return-staking tx
	TxConversionType         = "cv"  // Convert 1 - 2 normal tx
	TxTokenConversionType    = "tcv" // Convert 1 - 2 token tx
	TxCustomTokenPrivacyType = "tp"  // token  tx with supporting privacy
)

const (
	BlsConsensus    = "bls"
	BridgeConsensus = "dsa"
	IncKeyType      = "inc"
)

const PRVIDStr = "0000000000000000000000000000000000000000000000000000000000000004"

const (
	BurningAddress  = "15pABFiJVeh9D5uiQEhQX4SVibGGbdAVipQxBdxkmDqAJaoG1EdFKHBrNfs"
	BurningAddress2 = "12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA"
)

const (
	SalaryVerFixHash = 1
)
