package jsonresult

// RandomCommitmentAndPublicKeyResult represents the response to a request retrieving random decoys for transactions of version 2.
type RandomCommitmentAndPublicKeyResult struct {
	CommitmentIndices []uint64 `json:"CommitmentIndices"`
	PublicKeys        []string `json:"PublicKeys"`
	Commitments       []string `json:"Commitments"`
	AssetTags         []string `json:"AssetTags"`
}
