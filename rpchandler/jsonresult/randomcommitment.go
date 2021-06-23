package jsonresult

// RandomCommitmentResult represents the response to a request retrieving random decoys for transactions of version 1.
type RandomCommitmentResult struct {
	CommitmentIndices   []uint64 `json:"CommitmentIndices"`
	MyCommitmentIndices []uint64 `json:"MyCommitmentIndexs"`
	Commitments         []string `json:"Commitments"`
}
