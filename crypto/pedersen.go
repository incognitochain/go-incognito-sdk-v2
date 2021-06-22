package crypto

import (
	"github.com/pkg/errors"
)

const (
	PedersenPrivateKeyIndex = byte(0x00)
	PedersenValueIndex      = byte(0x01)
	PedersenSndIndex        = byte(0x02)
	PedersenShardIDIndex    = byte(0x03)
	PedersenRandomnessIndex = byte(0x04)
)

// PedCom represents the parameters used in the Pedersen commitment scheme.
var PedCom = NewPedersenParams()

// PedersenCommitment represents the base points for the Pedersen commitment scheme.
//	- G[0]: the base point for key-related commitments.
// 	- G[1]: the base point for value commitments.
// 	- G[2]: the base point for SNDerivator commitments.
// 	- G[3]: the base point for shardID commitments.
// 	- G[4]: the base point for randomness factor.
type PedersenCommitment struct {
	G []*Point // generators
}

// GBase is the base point for committing UTXOs' value.
var GBase *Point

// HBase is the base point for randomness factor.
var HBase *Point

// NewPedersenParams creates new PedersenCommitment parameters. It also sets GBase as G[1] and HBase as G[4].
func NewPedersenParams() PedersenCommitment {
	var pcm PedersenCommitment
	const capacity = 5 // fixed value = 5
	pcm.G = make([]*Point, capacity)
	pcm.G[0] = new(Point).ScalarMultBase(new(Scalar).FromUint64(1))

	for i := 1; i < len(pcm.G); i++ {
		pcm.G[i] = HashToPointFromIndex(int64(i), CStringBulletProof)
	}
	GBase = new(Point).Set(pcm.G[1])
	HBase = new(Point).Set(pcm.G[4])
	return pcm
}

// CommitAll commits a list of values using the corresponding base points.
func (com PedersenCommitment) CommitAll(openings []*Scalar) (*Point, error) {
	if len(openings) != len(com.G) {
		return nil, errors.New("invalid length of openings to commit")
	}

	commitment := new(Point).ScalarMult(com.G[0], openings[0])

	for i := 1; i < len(com.G); i++ {
		commitment.Add(commitment, new(Point).ScalarMult(com.G[i], openings[i]))
	}
	return commitment, nil
}

// CommitAtIndex commits specific value using the base point at the given index.
func (com PedersenCommitment) CommitAtIndex(value, rand *Scalar, index byte) *Point {
	return new(Point).AddPedersen(value, com.G[index], rand, com.G[PedersenRandomnessIndex])
}
