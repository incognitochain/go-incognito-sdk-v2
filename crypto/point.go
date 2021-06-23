package crypto

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	C25519 "github.com/incognitochain/go-incognito-sdk-v2/crypto/curve25519"
)

type Point struct {
	key C25519.Key
}

func RandomPoint() *Point {
	sc := RandomScalar()
	return new(Point).ScalarMultBase(sc)
}

func (p Point) PointValid() bool {
	var point C25519.ExtendedGroupElement
	return point.FromBytes(&p.key)
}

func (p Point) GetKey() C25519.Key {
	return p.key
}

func (p *Point) SetKey(a *C25519.Key) (*Point, error) {
	p.key = *a

	var point C25519.ExtendedGroupElement
	if !point.FromBytes(&p.key) {
		return nil, fmt.Errorf("invalid point value")
	}
	return p, nil
}

func (p *Point) Set(q *Point) *Point {
	p.key = q.key
	return p
}

func (p Point) String() string {
	return fmt.Sprintf("%x", p.key[:])
}

func (p Point) MarshalText() []byte {
	return []byte(fmt.Sprintf("%x", p.key[:]))
}

func (p *Point) UnmarshalText(data []byte) (*Point, error) {
	byteSlice, _ := hex.DecodeString(string(data))
	if len(byteSlice) != Ed25519KeySize {
		return nil, fmt.Errorf("incorrect key size")
	}
	copy(p.key[:], byteSlice)
	return p, nil
}

func (p Point) ToBytes() [Ed25519KeySize]byte {
	return p.key.ToBytes()
}

func (p Point) ToBytesS() []byte {
	slice := p.key.ToBytes()
	return slice[:]
}

func (p *Point) FromBytes(b [Ed25519KeySize]byte) (*Point, error) {
	p.key.FromBytes(b)

	var point C25519.ExtendedGroupElement
	if !point.FromBytes(&p.key) {
		return nil, fmt.Errorf("invalid point value")
	}

	return p, nil
}

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

func (p *Point) Identity() *Point {
	p.key = C25519.Identity
	return p
}

func (p Point) IsIdentity() bool {
	if p.key == C25519.Identity {
		return true
	}
	return false
}

// ScalarMultBase does a * G where a is a scalar and G is the curve base point.
func (p *Point) ScalarMultBase(a *Scalar) *Point {
	key := C25519.ScalarmultBase(&a.key)
	p.key = *key
	return p
}

func (p *Point) ScalarMult(pa *Point, a *Scalar) *Point {
	key := C25519.ScalarMultKey(&pa.key, &a.key)
	p.key = *key
	return p
}

func (p *Point) MultiScalarMultCached(scalarLs []*Scalar, pointPreComputedLs [][8]C25519.CachedGroupElement) *Point {
	nSc := len(scalarLs)

	if nSc != len(pointPreComputedLs) {
		panic("cannot Multi-ScalarMult with different size inputs")
	}

	scalarKeyLs := make([]*C25519.Key, nSc)
	for i := 0; i < nSc; i++ {
		scalarKeyLs[i] = &scalarLs[i].key
	}
	key := C25519.MultiScalarMultKeyCached(pointPreComputedLs, scalarKeyLs)
	res, _ := new(Point).SetKey(key)
	return res
}

func (p *Point) MultiScalarMult(scalarLs []*Scalar, pointLs []*Point) *Point {
	nSc := len(scalarLs)
	nPoint := len(pointLs)

	if nSc != nPoint {
		panic("cannot Multi-ScalarMult with different size inputs")
	}

	scalarKeyLs := make([]*C25519.Key, nSc)
	pointKeyLs := make([]*C25519.Key, nSc)
	for i := 0; i < nSc; i++ {
		scalarKeyLs[i] = &scalarLs[i].key
		pointKeyLs[i] = &pointLs[i].key
	}
	key := C25519.MultiScalarMultKey(pointKeyLs, scalarKeyLs)

	res, _ := new(Point).SetKey(key)
	return res
}

func (p *Point) InvertScalarMultBase(a *Scalar) *Point {
	inv := new(Scalar).Invert(a)
	p.ScalarMultBase(inv)
	return p
}

func (p *Point) InvertScalarMult(pa *Point, a *Scalar) *Point {
	inv := new(Scalar).Invert(a)
	p.ScalarMult(pa, inv)
	return p
}

func (p *Point) Derive(pa *Point, a *Scalar, b *Scalar) *Point {
	c := new(Scalar).Add(a, b)
	return p.InvertScalarMult(pa, c)
}

func (p *Point) Add(pa, pb *Point) *Point {
	res := p.key
	C25519.AddKeys(&res, &pa.key, &pb.key)
	p.key = res
	return p
}

// AddPedersen returns aA + bB.
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

func (p *Point) AddPedersenCached(a *Scalar, APreCompute [8]C25519.CachedGroupElement, b *Scalar, BPreCompute [8]C25519.CachedGroupElement) *Point {
	var key C25519.Key
	C25519.AddKeys3_3(&key, &a.key, &APreCompute, &b.key, &BPreCompute)
	p.key = key
	return p
}

func (p *Point) Sub(pa, pb *Point) *Point {
	res := p.key
	C25519.SubKeys(&res, &pa.key, &pb.key)
	p.key = res
	return p
}

func IsPointEqual(pa *Point, pb *Point) bool {
	tmpA := pa.ToBytesS()
	tmpB := pb.ToBytesS()
	return subtle.ConstantTimeCompare(tmpA, tmpB) == 1
}

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

func HashToPoint(b []byte) *Point {
	keyHash := C25519.Key(C25519.Keccak256(b))
	keyPoint := keyHash.HashToPoint()

	p, _ := new(Point).SetKey(keyPoint)
	return p
}
