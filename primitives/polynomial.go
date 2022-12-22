package primitives

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/drand/kyber/group/mod"
	"math/big"
	"strconv"
)

// Polynomial is the polynomial representation
type Polynomial struct {
	Coefficient []*mod.Int
	Degree      int // Degree is the number of coefficients which is the degree + 1
}

// RandModInt returns a random number between 0 and Q-1
func RandModInt() (*mod.Int, error) {
	maxbits := R.BitLen()
	b := make([]byte, (maxbits/8)-1)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	r := new(big.Int).SetBytes(b)
	rq := new(big.Int).Mod(r, Q)
	return new(mod.Int).Init(rq, Q), nil

}

// Init initializes a polynomial with the given coefficients
func (p *Polynomial) Init(coeffs []*mod.Int) *Polynomial {
	return &Polynomial{coeffs, len(coeffs)}
}

// InitFromCopy creates a new polynomial from a copy of the given one
func (p *Polynomial) InitFromCopy() *Polynomial {
	r := make([]*mod.Int, p.Degree)
	for i := 0; i < p.Degree; i++ {
		r[i] = new(mod.Int).Set(p.Coefficient[i]).(*mod.Int)
	}
	return &Polynomial{r[:], p.Degree}
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

func (p *Polynomial) SetCoefficient(coef []*mod.Int) {
	p.Coefficient = coef
	p.Degree = len(coef)
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
		p.Coefficient[i] = new(mod.Int).Sub(p.Coefficient[i], b.Coefficient[i]).(*mod.Int)
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
	rem := a.InitFromCopy()
	for rem.Degree >= b.Degree {
		l := new(mod.Int).Div(rem.Coefficient[rem.Degree-1], b.Coefficient[b.Degree-1]).(*mod.Int)
		pos := rem.Degree - b.Degree
		p.Coefficient[pos] = l
		aux := new(Polynomial).InitFromZerosArray(pos)
		aux1 := append(aux.Coefficient, l)
		tempoly := new(Polynomial).Init(aux1)
		aux2 := new(Polynomial).Sub(rem, new(Polynomial).Mul(b, tempoly))
		rem.SetCoefficient(aux2.Coefficient[:aux2.Degree-1])
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
		r.Add(r, elem)
	}
	return r
}

func (p *Polynomial) Cmp(a, b *Polynomial) bool {
	if a.Degree != b.Degree {
		return false
	}
	for i := 0; i < a.Degree; i++ {
		if a.Coefficient[i] != b.Coefficient[i] {
			return false
		}
	}
	return true
}

func (p *Polynomial) ToString() string {
	s := ""
	for i := p.Degree - 1; i >= 1; i-- {
		if bytes.Equal(p.Coefficient[i].V.Bytes(), big.NewInt(1).Bytes()) {
			s += fmt.Sprintf("x%s + ", intToSNum(i))
		} else if !bytes.Equal(p.Coefficient[i].V.Bytes(), big.NewInt(0).Bytes()) {
			t, _ := strconv.ParseInt(p.Coefficient[i].String(), 16, 64)
			s += fmt.Sprintf("%sx%s + ", strconv.Itoa(int(t)), intToSNum(i))
		}
	}
	t, _ := strconv.ParseInt(p.Coefficient[0].String(), 16, 64)
	s += strconv.Itoa(int(t))
	return s
}

// zeroPolynomial returns the zero polynomial:
// z(x) = (x - z_0) (x - z_1) ... (x - z_{k-1})
func (p *Polynomial) Zero(zs []*mod.Int) *Polynomial {
	p = &Polynomial{[]*mod.Int{new(mod.Int).Neg(zs[0]).(*mod.Int), new(mod.Int).Init(big.NewInt(1), Q)}, 2} // (x - z0)
	for i := 1; i < len(zs); i++ {
		p = new(Polynomial).Mul(p, &Polynomial{[]*mod.Int{new(mod.Int).Neg(zs[i]).(*mod.Int), new(mod.Int).Init(big.NewInt(1), Q)}, 2}) // (x - zi)
	}
	return p
}

// LagrangeInterpolation implements the Lagrange interpolation:
// https://en.wikipedia.org/wiki/Lagrange_polynomial
func (p *Polynomial) LagrangeInterpolation(x, y []*mod.Int) (*Polynomial, error) {
	// p(x) will be the interpoled polynomial
	// var p []*big.Int
	if len(x) != len(y) {
		return &Polynomial{nil, 0}, fmt.Errorf("len(x)!=len(y): %d, %d", len(x), len(y))
	}
	p = new(Polynomial).InitFromZerosArray(len(x))
	k := len(x)

	for j := 0; j < k; j++ {
		// jPol is the Lagrange basis polynomial for each point
		var jPol = new(Polynomial).InitFromZerosArray(0)
		for m := 0; m < k; m++ {
			// if x[m] == x[j] {
			if m == j {
				continue
			}
			// numerator & denominator of the current iteration
			num := &Polynomial{[]*mod.Int{new(mod.Int).Neg(x[m]).(*mod.Int), new(mod.Int).Init(big.NewInt(1), Q)}, 2} // (x^1 - x_m)
			den := new(mod.Int).Sub(x[j], x[m])                                                                       // x_j-x_m
			mPol := num.DivByConstant(den.(*mod.Int))
			if jPol.Degree == 0 {
				// first j iteration
				jPol = mPol
				continue
			}
			jPol = new(Polynomial).Mul(jPol, mPol)
		}
		p = new(Polynomial).Add(p, new(Polynomial).MulByConstant(jPol, y[j]))
	}

	return p, nil
}
