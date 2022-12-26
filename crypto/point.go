package crypto

import (
	"crypto/subtle"
	"fmt"

	C25519 "github.com/incognitochain/go-incognito-sdk-v2/crypto/curve25519"
)

// Point represents an elliptic curve point. It only needs 32 bytes to represent a point.
type Point struct {
	key C25519.Key
}

// RandomPoint returns a random Point on the elliptic curve.
func RandomPoint() *Point {
	sc := RandomScalar()
	return new(Point).ScalarMultBase(sc)
}

// PointValid checks if a Point is valid.
func (p Point) PointValid() bool {
	var point C25519.ExtendedGroupElement
	return point.FromBytes(&p.key)
}

// GetKey returns the key of a Point.
func (p Point) GetKey() C25519.Key {
	return p.key
}

// SetKey sets v as the key of a Point.
func (p *Point) SetKey(v *C25519.Key) (*Point, error) {
	p.key = *v

	var point C25519.ExtendedGroupElement
	if !point.FromBytes(&p.key) {
		return nil, fmt.Errorf("invalid point value")
	}
	return p, nil
}

// Set sets p = q and returns p.
func (p *Point) Set(q *Point) *Point {
	p.key = q.key
	return p
}

// String returns the hex-encoded string of a Point.
func (p Point) String() string {
	return fmt.Sprintf("%x", p.key[:])
}

// ToBytes returns an 32-byte long array from a Point.
func (p Point) ToBytes() [Ed25519KeySize]byte {
	return p.key.ToBytes()
}

// ToBytesS returns a slice of bytes from a Point.
func (p Point) ToBytesS() []byte {
	slice := p.key.ToBytes()
	return slice[:]
}

// FromBytes sets an array of 32 bytes to a Point.
func (p *Point) FromBytes(b [C25519.KeyLength]byte) (*Point, error) {
	p.key.FromBytes(b)

	var point C25519.ExtendedGroupElement
	if !point.FromBytes(&p.key) {
		return nil, fmt.Errorf("invalid point value")
	}

	return p, nil
}

// FromBytesS sets a slice of bytes to a Point.
func (p *Point) FromBytesS(b []byte) (*Point, error) {
	if len(b) != Ed25519KeySize {
		return nil, fmt.Errorf("invalid Ed25519 Key Size")
	}

	var array [Ed25519KeySize]byte
	copy(array[:], b)
	p.key.FromBytes(array)

	var point C25519.ExtendedGroupElement
	if !point.FromBytes(&p.key) {
		return nil, fmt.Errorf("invalid point value")
	}

	return p, nil
}

// Identity sets p to the identity point and returns p.
func (p *Point) Identity() *Point {
	p.key = C25519.Identity
	return p
}

// IsIdentity checks if p is the identity point.
func (p Point) IsIdentity() bool {
	if p.key == C25519.Identity {
		return true
	}
	return false
}

// ScalarMultBase set p = a * G, where a is a scalar and G is the curve base point and returns p.
func (p *Point) ScalarMultBase(a *Scalar) *Point {
	key := C25519.ScalarmultBase(&a.key)
	p.key = *key
	return p
}

// ScalarMult sets p = a * pa and returns p.
func (p *Point) ScalarMult(pa *Point, a *Scalar) *Point {
	key := C25519.ScalarMultKey(&pa.key, &a.key)
	p.key = *key
	return p
}

// MultiScalarMult sets p = sum(sList[i] * pList[i]) and returns p.
func (p *Point) MultiScalarMult(sList []*Scalar, pList []*Point) *Point {
	nSc := len(sList)
	nPoint := len(pList)

	if nSc != nPoint {
		panic("cannot Multi-ScalarMult with different size inputs")
	}

	scalarKeyLs := make([]*C25519.Key, nSc)
	pointKeyLs := make([]*C25519.Key, nSc)
	for i := 0; i < nSc; i++ {
		scalarKeyLs[i] = &sList[i].key
		pointKeyLs[i] = &pList[i].key
	}
	key := C25519.MultiScalarMultKey(pointKeyLs, scalarKeyLs)

	res, _ := new(Point).SetKey(key)
	return res
}

// InvertScalarMultBase sets p = (1/a) * G and returns p.
func (p *Point) InvertScalarMultBase(a *Scalar) *Point {
	inv := new(Scalar).Invert(a)
	p.ScalarMultBase(inv)
	return p
}

// InvertScalarMult sets p = (1/a) * pa and returns p.
func (p *Point) InvertScalarMult(pa *Point, a *Scalar) *Point {
	inv := new(Scalar).Invert(a)
	p.ScalarMult(pa, inv)
	return p
}

// Derive sets p = 1/(a+b) * pa and returns p.
func (p *Point) Derive(pa *Point, a *Scalar, b *Scalar) *Point {
	c := new(Scalar).Add(a, b)
	return p.InvertScalarMult(pa, c)
}

// Add sets p = pa + pb and returns p.
func (p *Point) Add(pa, pb *Point) *Point {
	res := p.key
	C25519.AddKeys(&res, &pa.key, &pb.key)
	p.key = res
	return p
}

// AddPedersen sets p = aA + bB and returns p.
func (p *Point) AddPedersen(a *Scalar, A *Point, b *Scalar, B *Point) *Point {
	var precomputedA [8]C25519.CachedGroupElement
	Ae := new(C25519.ExtendedGroupElement)
	Ae.FromBytes(&A.key)
	C25519.GePrecompute(&precomputedA, Ae)

	var precomputedB [8]C25519.CachedGroupElement
	Be := new(C25519.ExtendedGroupElement)
	Be.FromBytes(&B.key)
	C25519.GePrecompute(&precomputedB, Be)

	var key C25519.Key
	C25519.AddKeys3_3(&key, &a.key, &precomputedA, &b.key, &precomputedB)
	p.key = key
	return p
}

// Sub sets p = pa - pb and returns p.
func (p *Point) Sub(pa, pb *Point) *Point {
	res := p.key
	C25519.SubKeys(&res, &pa.key, &pb.key)
	p.key = res
	return p
}

// IsPointEqual checks if pa = pb.
func IsPointEqual(pa *Point, pb *Point) bool {
	tmpA := pa.ToBytesS()
	tmpB := pb.ToBytesS()
	return subtle.ConstantTimeCompare(tmpA, tmpB) == 1
}

// HashToPointFromIndex returns the hash of the concatenation of padStr and index.
func HashToPointFromIndex(index int64, padStr string) *Point {
	array := C25519.GBASE.ToBytes()
	msg := array[:]
	msg = append(msg, []byte(padStr)...)
	msg = append(msg, []byte(string(index))...)

	keyHash := C25519.Key(C25519.Keccak256(msg))
	keyPoint := keyHash.HashToPoint()

	p, _ := new(Point).SetKey(keyPoint)
	return p
}

// HashToPoint returns the Point value of the hash of b.
func HashToPoint(b []byte) *Point {
	keyHash := C25519.Key(C25519.Keccak256(b))
	keyPoint := keyHash.HashToPoint()

	p, _ := new(Point).SetKey(keyPoint)
	return p
}
