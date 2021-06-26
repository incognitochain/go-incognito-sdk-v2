package rpc

//OutCoinKey is used to retrieve output coins via RPC.
//
//Payment address must always be present. Other fields are optional
type OutCoinKey struct {
	paymentAddress string
	otaKey         string
	readonlyKey    string
}

func (outCoinKey OutCoinKey) PaymentAddress() string {
	return outCoinKey.paymentAddress
}

func (outCoinKey OutCoinKey) OtaKey() string {
	return outCoinKey.otaKey
}

func (outCoinKey OutCoinKey) ReadonlyKey() string {
	return outCoinKey.readonlyKey
}

func (outCoinKey *OutCoinKey) SetOTAKey(otaKey string) {
	outCoinKey.otaKey = otaKey
}

func (outCoinKey *OutCoinKey) SetPaymentAddress(paymentAddress string) {
	outCoinKey.paymentAddress = paymentAddress
}

func (outCoinKey *OutCoinKey) SetReadonlyKey(readonlyKey string) {
	outCoinKey.readonlyKey = readonlyKey
}

func NewOutCoinKey(paymentAddress, otaKey, readonlyKey string) *OutCoinKey {
	return &OutCoinKey{paymentAddress: paymentAddress, otaKey: otaKey, readonlyKey: readonlyKey}
}