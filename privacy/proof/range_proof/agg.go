package range_proof

// RangeProof represents a range proof, which is used to prove a number lies with-in a specific
// interval without revealing the number.
type RangeProof interface {
	Init()
	IsNil() bool
	Bytes() []byte
	SetBytes([]byte) error
}
