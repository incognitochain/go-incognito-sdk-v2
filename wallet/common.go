package wallet

const (
	privateKeySerializedBytesLen  = 75 // bytes
	paymentAddrSerializedBytesLen = 71 // bytes
	readOnlyKeySerializedBytesLen = 71 // bytes
	otaKeySerializedBytesLen      = 71 // bytes
)

const (
	PrivateKeyType     = byte(0x0) // Serialize wallet account key into string with only PRIVATE KEY of account key set
	PaymentAddressType = byte(0x1) // Serialize wallet account key into string with only PAYMENT ADDRESS of account key set
	ReadonlyKeyType    = byte(0x2) // Serialize wallet account key into string with only READONLY KEY of account key set
	OTAKeyType         = byte(0x3) // Serialize wallet account key into string with only OTA KEY of account key set
)

const (
	DefaultPassword      = "12345678"
	HardenedKeyZeroIndex = 0x80000000
	BIP44Purpose         = 44
	Bip44CoinType        = 587
	MinSeedBytes         = 16
	MaxSeedBytes         = 64
)

var (
	burnAddress1BytesDecode = []byte{1, 32, 99, 183, 246, 161, 68, 172, 228, 222, 153, 9, 172, 39, 208, 245, 167, 79, 11, 2, 114, 65, 241, 69, 85, 40, 193, 104, 199, 79, 70, 4, 53, 0, 0, 163, 228, 236, 208}
)
