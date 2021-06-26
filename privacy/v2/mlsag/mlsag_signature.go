package mlsag

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
)

type MlsagSig struct {
	c         *crypto.Scalar     // 32 bytes
	keyImages []*crypto.Point    // 32 * size bytes
	r         [][]*crypto.Scalar // 32 * size_1 * size_2 bytes
}

func NewMlsagSig(c *crypto.Scalar, keyImages []*crypto.Point, r [][]*crypto.Scalar) (*MlsagSig, error) {
	if len(r) == 0 {
		return nil, fmt.Errorf("cannot create new mlsag signature, length of r is not correct")
	}
	if len(keyImages) != len(r[0]) {
		return nil, fmt.Errorf("cannot create new mlsag signature, length of keyImages is not correct")
	}
	res := new(MlsagSig)
	res.SetC(c)
	res.SetR(r)
	res.SetKeyImages(keyImages)
	return res, nil
}

func (this MlsagSig) GetC() *crypto.Scalar          { return this.c }
func (this MlsagSig) GetKeyImages() []*crypto.Point { return this.keyImages }
func (this MlsagSig) GetR() [][]*crypto.Scalar      { return this.r }

func (this *MlsagSig) SetC(c *crypto.Scalar)                  { this.c = c }
func (this *MlsagSig) SetKeyImages(keyImages []*crypto.Point) { this.keyImages = keyImages }
func (this *MlsagSig) SetR(r [][]*crypto.Scalar)              { this.r = r }

func (this *MlsagSig) ToBytes() ([]byte, error) {
	b := []byte{MlsagPrefix}

	if this.c != nil {
		b = append(b, crypto.Ed25519KeySize)
		b = append(b, this.c.ToBytesS()...)
	} else {
		b = append(b, 0)
	}

	if this.keyImages != nil {
		if len(this.keyImages) > MaxSizeByte {
			return nil, fmt.Errorf("length of key image is too large > 255")
		}
		lenKeyImage := byte(len(this.keyImages) & 0xFF)
		b = append(b, lenKeyImage)
		for i := 0; i < int(lenKeyImage); i += 1 {
			b = append(b, this.keyImages[i].ToBytesS()...)
		}
	} else {
		b = append(b, 0)
	}

	if this.r != nil {
		n := len(this.r)
		if n == 0 {
			b = append(b, 0)
			b = append(b, 0)
			return b, nil
		}
		m := len(this.r[0])
		if n > MaxSizeByte || m > MaxSizeByte {
			return nil, fmt.Errorf("length of R of mlsagSig is too large > 255")
		}
		b = append(b, byte(n&0xFF))
		b = append(b, byte(m&0xFF))
		for i := 0; i < n; i += 1 {
			if m != len(this.r[i]) {
				return []byte{}, fmt.Errorf("error in MLSAG MlsagSig ToBytes: the signature is broken (size of keyImages and r differ)")
			}
			for j := 0; j < m; j += 1 {
				b = append(b, this.r[i][j].ToBytesS()...)
			}
		}
	} else {
		b = append(b, 0)
		b = append(b, 0)
	}

	return b, nil
}

// Get from byte and store to signature
func (this *MlsagSig) FromBytes(b []byte) (*MlsagSig, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("length of byte is empty, cannot setbyte mlsagSig")
	}
	if b[0] != MlsagPrefix {
		return nil, fmt.Errorf("The signature byte is broken (first byte is not mlsag)")
	}

	offset := 1
	if b[offset] != crypto.Ed25519KeySize {
		return nil, fmt.Errorf("cannot parse value C, byte length of C is wrong")
	}
	offset += 1
	if offset+crypto.Ed25519KeySize > len(b) {
		return nil, fmt.Errorf("cannot parse value C, byte is too small")
	}
	C := new(crypto.Scalar).FromBytesS(b[offset : offset+crypto.Ed25519KeySize])
	if !C.ScalarValid() {
		return nil, fmt.Errorf("cannot parse value C, invalid scalar")
	}
	offset += crypto.Ed25519KeySize

	if offset >= len(b) {
		return nil, fmt.Errorf("cannot parse length of keyimage, byte is too small")
	}
	lenKeyImages := int(b[offset])
	offset += 1
	keyImages := make([]*crypto.Point, lenKeyImages)
	for i := 0; i < lenKeyImages; i += 1 {
		if offset+crypto.Ed25519KeySize > len(b) {
			return nil, fmt.Errorf("cannot parse keyimage of mlsagSig, byte is too small")
		}
		var err error
		keyImages[i], err = new(crypto.Point).FromBytesS(b[offset : offset+crypto.Ed25519KeySize])
		if err != nil {
			return nil, fmt.Errorf("cannot convert byte to crypto point keyimage")
		}
		offset += crypto.Ed25519KeySize
	}

	if offset+2 > len(b) {
		return nil, fmt.Errorf("cannot parse length of R, byte is too small")
	}
	n := int(b[offset])
	m := int(b[offset+1])
	offset += 2

	R := make([][]*crypto.Scalar, n)
	for i := 0; i < n; i += 1 {
		R[i] = make([]*crypto.Scalar, m)
		for j := 0; j < m; j += 1 {
			if offset+crypto.Ed25519KeySize > len(b) {
				return nil, fmt.Errorf("cannot parse R of mlsagSig, byte is too small")
			}
			sc := new(crypto.Scalar).FromBytesS(b[offset : offset+crypto.Ed25519KeySize])
			if !sc.ScalarValid() {
				return nil, fmt.Errorf("cannot parse R of mlsagSig, invalid scalar")
			}
			R[i][j] = sc
			offset += crypto.Ed25519KeySize
		}
	}

	if this == nil {
		this = new(MlsagSig)
	}
	this.SetC(C)
	this.SetKeyImages(keyImages)
	this.SetR(R)
	return this, nil
}
