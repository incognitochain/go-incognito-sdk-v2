package rpc

// OutCoinKey is used to retrieve output coins via RPC.
//
// The payment addresses is required in all cases. For retrieving output coins V2, the ota key is required.
// Readonly keys are optional.
type OutCoinKey struct {
	paymentAddress string
	otaKey         string
	readonlyKey    string
}

// PaymentAddress returns the payment address of an OutCoinKey.
func (outCoinKey OutCoinKey) PaymentAddress() string {
	return outCoinKey.paymentAddress
}

// OtaKey returns the ota key of an OutCoinKey.
func (outCoinKey OutCoinKey) OtaKey() string {
	return outCoinKey.otaKey
}

// ReadonlyKey returns the read-only of an OutCoinKey.
func (outCoinKey OutCoinKey) ReadonlyKey() string {
	return outCoinKey.readonlyKey
}

// SetOTAKey sets v as the ota key of an OutCoinKey.
func (outCoinKey *OutCoinKey) SetOTAKey(v string) {
	outCoinKey.otaKey = v
}

// SetPaymentAddress sets v as the payment address of an OutCoinKey.
func (outCoinKey *OutCoinKey) SetPaymentAddress(v string) {
	outCoinKey.paymentAddress = v
}

// SetReadonlyKey sets v as the read-only key of an OutCoinKey.
func (outCoinKey *OutCoinKey) SetReadonlyKey(v string) {
	outCoinKey.readonlyKey = v
}

// NewOutCoinKey create a new OutCoinKey with the given parameters.
func NewOutCoinKey(paymentAddress, otaKey, readonlyKey string) *OutCoinKey {
	return &OutCoinKey{paymentAddress: paymentAddress, otaKey: otaKey, readonlyKey: readonlyKey}
}
