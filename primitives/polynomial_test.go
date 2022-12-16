package primitives

import (
	"bytes"
	"github.com/drand/kyber/group/mod"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func negf(a *mod.Int) *mod.Int {
	return new(mod.Int).Neg(a).(*mod.Int)
}

func TestPolynomial(t *testing.T) {
	b0 := new(mod.Int).Init64(int64(0), Q)
	b1 := new(mod.Int).Init64(int64(1), Q)
	b2 := new(mod.Int).Init64(int64(2), Q)
	b3 := new(mod.Int).Init64(int64(3), Q)
	b4 := new(mod.Int).Init64(int64(4), Q)
	b5 := new(mod.Int).Init64(int64(5), Q)
	b6 := new(mod.Int).Init64(int64(6), Q)
	b16 := new(mod.Int).Init64(int64(16), Q)

	a := &Polynomial{[]*mod.Int{b1, b0, b5}, 3}
	b := &Polynomial{[]*mod.Int{b3, b0, b1}, 3}

	// new Finite Field
	r, ok := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10) //nolint:lll
	fr := new(mod.Int).Init(r, Q)
	assert.True(nil, ok)

	// polynomial multiplication
	o := new(Polynomial).Mul(a, b)

	assert.Equal(t, o, &Polynomial{[]*mod.Int{b3, b0, b16, b0, b5}, 5})

	// polynomial division
	quo, rem := new(Polynomial).Div(a, b)
	assert.Equal(t, quo.Coefficient[0].Int64(), int64(5))
	// check the rem result without modulo
	assert.Equal(t, new(mod.Int).Sub(rem.Coefficient[0], fr).(*mod.Int).V.Int64(), int64(-14))

	c := new(Polynomial).Init([]*mod.Int{negf(b4), b0, negf(b2), b1})
	d := new(Polynomial).Init([]*mod.Int{negf(b3), b1})
	quo2, rem2 := new(Polynomial).Div(c, d)
	assert.Equal(t, quo2, new(Polynomial).Init([]*mod.Int{b3, b1, b1}))
	assert.Equal(t, rem2.Coefficient[0].Int64(), int64(5))

	// polynomial addition
	o = new(Polynomial).Add(a, b)
	assert.Equal(t, o, []*mod.Int{b4, b0, b6})

	// polynomial subtraction
	o1 := new(Polynomial).Sub(a, b)
	o2 := new(Polynomial).Sub(b, a)
	o = new(Polynomial).Add(o1, o2)
	assert.True(t, bytes.Equal(b0.V.Bytes(), o.Coefficient[0].V.Bytes()))
	assert.True(t, bytes.Equal(b0.V.Bytes(), o.Coefficient[1].V.Bytes()))
	assert.True(t, bytes.Equal(b0.V.Bytes(), o.Coefficient[2].V.Bytes()))

	//c = new(Polynomial).Init([]*big.Int{b5, b6, b1})
	//d = new(Polynomial).Init([]*big.Int{b1, b3})
	//o = PolynomialSub(c, d)
	//assert.Equal(t, o, []*big.Int{b4, b3, b1})
	//
	//// NewPolZeroAt
	//o = NewPolZeroAt(3, 4, b4)
	//assert.Equal(t, PolynomialEval(o, big.NewInt(3)), b4)
	//o = NewPolZeroAt(2, 4, b3)
	//assert.Equal(t, PolynomialEval(o, big.NewInt(2)), b3)

	// polynomialEval
	// p(x) = x^3 + x + 5
	p := new(Polynomial).Init([]*mod.Int{
		new(mod.Int).Init(big.NewInt(5), Q),
		new(mod.Int).Init(big.NewInt(1), Q), // x^1
		new(mod.Int).Init(big.NewInt(0), Q), // x^2
		new(mod.Int).Init(big.NewInt(1), Q), // x^3
	})
	assert.Equal(t, "x³ + x¹ + 5", p.ToString())
	assert.Equal(t, "35", p.Eval(new(mod.Int).Init(big.NewInt(3), Q)).String())
	assert.Equal(t, "1015", p.Eval(new(mod.Int).Init(big.NewInt(10), Q)).String())
	assert.Equal(t, "16777477", p.Eval(new(mod.Int).Init(big.NewInt(256), Q)).String())
	assert.Equal(t, "125055", p.Eval(new(mod.Int).Init(big.NewInt(50), Q)).String())
	assert.Equal(t, "7", p.Eval(new(mod.Int).Init(big.NewInt(1), Q)).String())
}
