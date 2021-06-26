package polynomial

import (
	"fmt"
	"math/big"
)

// Poly is a data structure representing a polynomial. A Poly is just an array in the reverse order.
//
// For example, f(x) = 3x^3 + 2x + 1 => [1 2 0 3].
type Poly []*big.Int

// NewPoly returns the polynomial with given integers.
func NewPoly(coEffs ...int) (p Poly) {
	p = make([]*big.Int, len(coEffs))
	for i := 0; i < len(coEffs); i++ {
		p[i] = big.NewInt(int64(coEffs[i]))
	}
	p.Trim()
	return
}

// Trim makes sure that the highest coefficient never has zero value
// When we Add or subtract two polynomials, sometimes, the highest coefficient is zero,
// if we don't remove the highest and zero coefficient, GetDegree() will return the wrong result.
func (p *Poly) Trim() {
	var last = 0
	for i := p.GetDegree(); i > 0; i-- { // why i > 0, not i >=0? do not remove the constant
		if (*p)[i].Sign() != 0 {
			last = i
			break
		}
	}
	*p = (*p)[:(last + 1)]
}

// IsZero checks if p = 0.
func (p *Poly) IsZero() bool {
	if p.GetDegree() == 0 && (*p)[0].Cmp(big.NewInt(0)) == 0 {
		return true
	}
	return false
}

// GetDegree returns the degree of a Poly.
//
// For example, if p = x^3 + 2x^2 + 5, GetDegree() returns 3.
func (p Poly) GetDegree() int {
	return len(p) - 1
}

// String returns a beautified string-representation of a Poly.
func (p Poly) String() (s string) {
	s = "["
	for i := len(p) - 1; i >= 0; i-- {
		switch p[i].Sign() {
		case -1:
			if i == len(p)-1 {
				s += "-"
			} else {
				s += " - "
			}
			if i == 0 || p[i].Int64() != -1 {
				s += p[i].String()[1:]
			}
		case 0:
			continue
		case 1:
			if i < len(p)-1 {
				s += " + "
			}
			if i == 0 || p[i].Int64() != 1 {
				s += p[i].String()
			}
		}
		if i > 0 {
			s += "x"
			if i > 1 {
				s += "^" + fmt.Sprintf("%d", i)
			}
		}
	}
	if s == "[" {
		s += "0"
	}
	s += "]"
	return
}

// Compare returns -1 if p < 1; 0 if p == q; and 1 if p > q.
func (p *Poly) Compare(q *Poly) int {
	switch {
	case p.GetDegree() > q.GetDegree():
		return 1
	case p.GetDegree() < q.GetDegree():
		return -1
	}
	for i := 0; i <= p.GetDegree(); i++ {
		switch (*p)[i].Cmp((*q)[i]) {
		case 1:
			return 1
		case -1:
			return -1
		}
	}
	return 0
}

// Add returns (p + q) % m.
func (p Poly) Add(q Poly, m *big.Int) Poly {
	if p.Compare(&q) < 0 {
		return q.Add(p, m)
	}
	var r Poly = make([]*big.Int, len(p))
	for i := 0; i < len(q); i++ {
		a := new(big.Int)
		a.Add(p[i], q[i])
		r[i] = a
	}
	for i := len(q); i < len(p); i++ {
		a := new(big.Int)
		a.Set(p[i])
		r[i] = a
	}
	if m != nil {
		for i := 0; i < len(p); i++ {
			r[i].Mod(r[i], m)
		}
	}
	r.Trim()
	return r
}

// Neg returns -p.
func (p *Poly) Neg() Poly {
	var q Poly = make([]*big.Int, len(*p))
	for i := 0; i < len(*p); i++ {
		b := new(big.Int)
		b.Neg((*p)[i])
		q[i] = b
	}
	return q
}

// Sub returns (p - q) % m.
func (p Poly) Sub(q Poly, m *big.Int) Poly {
	r := q.Neg()
	return p.Add(r, m)
}

// Mul returns (p * q) % m.
func (p Poly) Mul(q Poly, m *big.Int) Poly {
	if m != nil {
		p.sanitize(m)
		q.sanitize(m)
	}
	var r Poly = make([]*big.Int, p.GetDegree()+q.GetDegree()+1)
	for i := 0; i < len(r); i++ {
		r[i] = big.NewInt(0)
	}
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(q); j++ {
			a := new(big.Int)
			a.Mul(p[i], q[j])
			a.Add(a, r[i+j])
			if m != nil {
				a = new(big.Int).Mod(a, m)
			}
			r[i+j] = a
		}
	}
	r.Trim()
	return r
}

// Div returns (p / q, p % q).
func (p Poly) Div(q Poly, m *big.Int) (quo, rem Poly) {
	if m != nil {
		p.sanitize(m)
		q.sanitize(m)
	}
	if p.GetDegree() < q.GetDegree() || q.IsZero() {
		quo = NewPoly(0)
		rem = p.clone(0)
		return
	}
	quo = make([]*big.Int, p.GetDegree()-q.GetDegree()+1)
	rem = p.clone(0)
	for i := 0; i < len(quo); i++ {
		quo[i] = big.NewInt(0)
	}
	t := p.clone(0)
	qd := q.GetDegree()
	for {
		td := t.GetDegree()
		rd := td - qd
		if rd < 0 || t.IsZero() {
			rem = t
			break
		}
		r := new(big.Int)
		if m != nil {
			r.ModInverse(q[qd], m)
			r.Mul(r, t[td])
			r.Mod(r, m)
		} else {
			r.Div(t[td], q[qd])
		}
		// if r == 0, it means that the highest coefficient of the result is not an integer
		// this polynomial library handles integer coefficients
		if r.Cmp(big.NewInt(0)) == 0 {
			quo = NewPoly(0)
			rem = p.clone(0)
			return
		}
		u := q.clone(rd)
		for i := rd; i < len(u); i++ {
			u[i].Mul(u[i], r)
			if m != nil {
				u[i].Mod(u[i], m)
			}
		}
		t = t.Sub(u, m)
		t.Trim()
		quo[rd] = r
	}
	quo.Trim()
	rem.Trim()
	return
}

// GCD returns the greatest common divisor(GCD) of p and q (Euclidean algorithm).
func (p Poly) GCD(q Poly, m *big.Int) Poly {
	if p.Compare(&q) < 0 {
		return q.GCD(p, m)
	}
	if q.IsZero() {
		return p
	} else {
		_, rem := p.Div(q, m)
		return q.GCD(rem, m)
	}
}

// Eval returns p(x) % m.
func (p Poly) Eval(x *big.Int, m *big.Int) (y *big.Int) {
	y = big.NewInt(0)
	accumulatedX := big.NewInt(1)
	xd := new(big.Int)
	for i := 0; i <= p.GetDegree(); i++ {
		xd.Mul(accumulatedX, p[i])
		y.Add(y, xd)
		accumulatedX.Mul(accumulatedX, x)
		if m != nil {
			y.Mod(y, m)
			accumulatedX.Mod(accumulatedX, m)
		}
	}
	return y
}

// clone returns a copy of a Poly whose degree is adjusted by the given value. The adjusted value cannot be negative.
//
// For example, if p = x + 1 and adjust = 2, Clone() will return x^3 + x^2.
func (p Poly) clone(adjust int) Poly {
	var q Poly = make([]*big.Int, len(p)+adjust)
	if adjust < 0 {
		return NewPoly(0)
	}
	for i := 0; i < adjust; i++ {
		q[i] = big.NewInt(0)
	}
	for i := adjust; i < len(p)+adjust; i++ {
		b := new(big.Int)
		b.Set(p[i-adjust])
		q[i] = b
	}
	return q
}

// sanitize does modular arithmetic with m, i.e. p % m.
func (p *Poly) sanitize(m *big.Int) {
	if m == nil {
		return
	}
	for i := 0; i <= (*p).GetDegree(); i++ {
		(*p)[i].Mod((*p)[i], m)
	}
	p.Trim()
}
