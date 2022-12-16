package primitives

import (
	"bytes"
	"fmt"
	"github.com/drand/kyber/group/mod"
	"math/big"
)

type Polynomial struct {
	Coefficient []*mod.Int
	Degree      int
}

func (p *Polynomial) Init(coeffs []*mod.Int) *Polynomial {
	return &Polynomial{coeffs, len(coeffs)}
}

func (p *Polynomial) InitFromZerosArray(n int) *Polynomial {
	r := make([]*mod.Int, n)
	for i := 0; i < n; i++ {
		r[i] = new(mod.Int).Init(big.NewInt(int64(0)), Q)
	}
	return &Polynomial{r[:], n}
}

// newPolZeroAt generates a new polynomial that has value zero at the given value
func (p *Polynomial) InitPolZeroAt(pointPos, totalPoints int, height *mod.Int) *Polynomial {
	fac := 1
	for i := 1; i < totalPoints+1; i++ {
		if i != pointPos {
			fac = fac * (pointPos - i)
		}
	}
	facBig := new(mod.Int).Init(big.NewInt(int64(fac)), Q)
	hf := new(mod.Int).Div(height, facBig)
	r := new(Polynomial).Init([]*mod.Int{hf.(*mod.Int)})
	for i := 1; i < totalPoints+1; i++ {
		if i != pointPos {
			ineg := new(mod.Int).Init(big.NewInt(int64(-i)), Q)
			b1 := new(mod.Int).Init(big.NewInt(int64(1)), Q)
			r = p.Mul(r, &Polynomial{[]*mod.Int{ineg, b1}, 2})
		}
	}
	return r
}

// polynomial operation.

func (p *Polynomial) Add(a, b *Polynomial) *Polynomial {
	p = p.InitFromZerosArray(max(a.Degree, b.Degree))
	for i := 0; i < a.Degree; i++ {
		p.Coefficient[i] = new(mod.Int).Add(p.Coefficient[i], a.Coefficient[i]).(*mod.Int)
	}
	for i := 0; i < b.Degree; i++ {
		p.Coefficient[i] = new(mod.Int).Add(p.Coefficient[i], b.Coefficient[i]).(*mod.Int)
	}
	return p
}

func (p *Polynomial) Sub(a, b *Polynomial) *Polynomial {
	p = p.InitFromZerosArray(max(a.Degree, b.Degree))
	for i := 0; i < a.Degree; i++ {
		p.Coefficient[i] = new(mod.Int).Add(p.Coefficient[i], a.Coefficient[i]).(*mod.Int)
	}
	for i := 0; i < b.Degree; i++ {
		p.Coefficient[i] = new(mod.Int).Sub(p.Coefficient[i], a.Coefficient[i]).(*mod.Int)
	}
	return p
}

func (p *Polynomial) Mul(a, b *Polynomial) *Polynomial {
	p = new(Polynomial).InitFromZerosArray(a.Degree + b.Degree - 1)
	for i := 0; i < a.Degree; i++ {
		for j := 0; j < b.Degree; j++ {
			p.Coefficient[i+j] = new(mod.Int).Add(p.Coefficient[i+j], new(mod.Int).Mul(a.Coefficient[i], b.Coefficient[j])).(*mod.Int)
		}
	}
	return p
}

func (p *Polynomial) Div(a, b *Polynomial) (*Polynomial, *Polynomial) {
	// https://en.wikipedia.org/wiki/Division_algorithm
	p = new(Polynomial).InitFromZerosArray(a.Degree - b.Degree + 1)
	rem := a
	for rem.Degree >= b.Degree {
		l := new(mod.Int).Div(rem.Coefficient[rem.Degree-1], b.Coefficient[b.Degree-1]).(*mod.Int)
		pos := rem.Degree - b.Degree
		p.Coefficient[pos] = l
		aux := new(Polynomial).InitFromZerosArray(pos)
		aux1 := append(aux.Coefficient, l)
		tempoly := new(Polynomial).Init(aux1)
		aux2 := new(Polynomial).Sub(rem, new(Polynomial).Mul(b, tempoly))
		rem.Coefficient = aux2.Coefficient[:aux2.Degree-1]
	}
	return p, rem
}

func (p *Polynomial) MulByConstant(a *Polynomial, c *mod.Int) *Polynomial {
	for i := 0; i < a.Degree; i++ {
		a.Coefficient[i] = new(mod.Int).Mul(a.Coefficient[i], c).(*mod.Int)
	}
	return a
}
func (p *Polynomial) DivByConstant(c *mod.Int) *Polynomial {
	for i := 0; i < p.Degree; i++ {
		p.Coefficient[i] = new(mod.Int).Div(p.Coefficient[i], c).(*mod.Int)
	}
	return p
}

// polynomialEval evaluates the polinomial over the Finite Field at the given value x
func (p *Polynomial) Eval(x *mod.Int) *mod.Int {
	r := new(mod.Int).Init(big.NewInt(int64(0)), Q)
	for i := 0; i < p.Degree; i++ {
		xi := new(mod.Int).Exp(x, big.NewInt(int64(i)))
		elem := new(mod.Int).Mul(p.Coefficient[i], xi)
		r = new(mod.Int).Add(r, elem).(*mod.Int)
	}
	return r
}

func (p *Polynomial) ToString() string {
	s := ""
	for i := p.Degree - 1; i >= 1; i-- {
		if bytes.Equal(p.Coefficient[i].V.Bytes(), big.NewInt(1).Bytes()) {
			s += fmt.Sprintf("x%s + ", intToSNum(i))
		} else if !bytes.Equal(p.Coefficient[i].V.Bytes(), big.NewInt(0).Bytes()) {
			s += fmt.Sprintf("%sx%s + ", p.Coefficient[i], intToSNum(i))
		}
	}
	s += p.Coefficient[0].String()
	return s
}

// zeroPolynomial returns the zero polynomial:
// z(x) = (x - z_0) (x - z_1) ... (x - z_{k-1})
func (p *Polynomial) ZeroPolynomial(zs []*big.Int) Poly {
	z := Poly{[]*big.Int{FieldNeg(zs[0]), big.NewInt(1)}, 2} // (x - z0)
	for i := 1; i < len(zs); i++ {
		z = PolynomialMul(z, Poly{[]*big.Int{FieldNeg(zs[i]), big.NewInt(1)}, 2}) // (x - zi)
	}
	return z
}
