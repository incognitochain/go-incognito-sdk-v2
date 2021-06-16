package range_proof

type RangeProof interface {
	Init()
	IsNil() bool
	Bytes() []byte
	SetBytes([]byte) error
	Verify() (bool, error)
}

// type AggregatedRangeProofV1 = aggregatedrange.RangeProof
// type AggregatedRangeProofV2 = bulletproofs.RangeProof
