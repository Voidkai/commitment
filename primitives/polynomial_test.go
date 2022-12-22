package primitives

import (
	"github.com/drand/kyber/group/mod"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"testing"
)

func negf(a *mod.Int) *mod.Int {
	return new(mod.Int).Neg(a).(*mod.Int)
}

func TestPolynomial_ToString(t *testing.T) {
	// polynomialEval
	// p(x) = x^3 + x + 5
	p := new(Polynomial).Init([]*mod.Int{
		mod.NewInt64(5, Q),
		mod.NewInt64(1, Q), // x^1
		mod.NewInt64(0, Q), // x^2
		mod.NewInt64(1, Q), // x^3
	})
	assert.Equal(t, "x³ + x¹ + 5", p.ToString())
}

func TestPolynomial_Add(t *testing.T) {
	b0 := new(mod.Int).Init64(int64(0), Q)
	b1 := new(mod.Int).Init64(int64(1), Q)
	b3 := new(mod.Int).Init64(int64(3), Q)
	b4 := new(mod.Int).Init64(int64(4), Q)
	b5 := new(mod.Int).Init64(int64(5), Q)
	b6 := new(mod.Int).Init64(int64(6), Q)

	a := &Polynomial{[]*mod.Int{b1, b0, b5}, 3}
	b := &Polynomial{[]*mod.Int{b3, b0, b1}, 3}

	// polynomial addition
	o := new(Polynomial).Add(a, b)
	assert.Equal(t, o.Coefficient, []*mod.Int{b4, b0, b6})

}

func TestPolynomial_Sub(t *testing.T) {
	b0 := new(mod.Int).Init64(int64(0), Q)
	b1 := new(mod.Int).Init64(int64(1), Q)
	//b2 := new(mod.Int).Init64(int64(2), Q)
	b3 := new(mod.Int).Init64(int64(3), Q)
	b4 := new(mod.Int).Init64(int64(4), Q)
	b5 := new(mod.Int).Init64(int64(5), Q)
	b6 := new(mod.Int).Init64(int64(6), Q)
	//b16 := new(mod.Int).Init64(int64(16), Q)
	bn2 := new(mod.Int).Init64(int64(-2), Q)

	a := &Polynomial{[]*mod.Int{b1, b0, b5}, 3} // 1 + 5x^2
	b := &Polynomial{[]*mod.Int{b3, b0, b1}, 3} // 3 + x^2

	// polynomial subtraction
	o := new(Polynomial).Sub(a, b)
	println(o.ToString())
	assert.Equal(t, o.Coefficient, []*mod.Int{bn2, b0, b4})

	c := new(Polynomial).Init([]*mod.Int{b5, b6, b1})
	d := new(Polynomial).Init([]*mod.Int{b6, b3})
	o = new(Polynomial).Sub(c, d)
	assert.Equal(t, o.Coefficient, []*mod.Int{new(mod.Int).Neg(b1).(*mod.Int), b3, b1})
}

func TestPolynomial_Mul(t *testing.T) {
	b0 := new(mod.Int).Init64(int64(0), Q)
	b1 := new(mod.Int).Init64(int64(1), Q)
	//b2 := new(mod.Int).Init64(int64(2), Q)
	b3 := new(mod.Int).Init64(int64(3), Q)
	//b4 := new(mod.Int).Init64(int64(4), Q)
	b5 := new(mod.Int).Init64(int64(5), Q)
	//b6 := new(mod.Int).Init64(int64(6), Q)
	b16 := new(mod.Int).Init64(int64(16), Q)

	a := &Polynomial{[]*mod.Int{b1, b0, b5}, 3} // 1 + 5x^2
	b := &Polynomial{[]*mod.Int{b3, b0, b1}, 3} // 3 + x^2

	// polynomial multiplication
	o := new(Polynomial).Mul(a, b)
	assert.Equal(t, o.Coefficient, []*mod.Int{b3, b0, b16, b0, b5})
}

func TestPolynomial_Div(t *testing.T) {
	b0 := new(mod.Int).Init64(int64(0), Q)
	b1 := new(mod.Int).Init64(int64(1), Q)
	b2 := new(mod.Int).Init64(int64(2), Q)
	b3 := new(mod.Int).Init64(int64(3), Q)
	b4 := new(mod.Int).Init64(int64(4), Q)
	b5 := new(mod.Int).Init64(int64(5), Q)
	//b6 := new(mod.Int).Init64(int64(6), Q)
	b14 := new(mod.Int).Init64(int64(14), Q)
	//b16 := new(mod.Int).Init64(int64(16), Q)

	a := &Polynomial{[]*mod.Int{b1, b0, b5}, 3} // 1 + 5x^2
	b := &Polynomial{[]*mod.Int{b3, b0, b1}, 3} // 3 + x^2

	// polynomial division
	quo, rem := new(Polynomial).Div(a, b)
	assert.Equal(t, quo.Coefficient[0], b5)
	// check the rem result without modulo
	assert.Equal(t, rem.Coefficient[0], negf(b14))
	c := new(Polynomial).Init([]*mod.Int{negf(b4), b0, negf(b2), b1})
	d := new(Polynomial).Init([]*mod.Int{negf(b3), b1})
	quo2, rem2 := new(Polynomial).Div(c, d)
	assert.Equal(t, quo2, new(Polynomial).Init([]*mod.Int{b3, b1, b1}))
	assert.Equal(t, rem2.Coefficient[0].Int64(), int64(5))
}

func TestPolynomial_InitPolZeroAt(t *testing.T) {
	b2 := new(mod.Int).Init64(int64(2), Q)
	b4 := new(mod.Int).Init64(int64(4), Q)
	b3 := new(mod.Int).Init64(int64(3), Q)
	// NewPolZeroAt
	o := new(Polynomial).InitPolZeroAt(3, 4, b4)
	assert.Equal(t, o.Eval(b3), b4)
	o = new(Polynomial).InitPolZeroAt(2, 4, b3)
	assert.Equal(t, o.Eval(b2), b3)
}

func TestPolynomial_Eval(t *testing.T) {
	// polynomialEval
	// p(x) = x^3 + x + 5
	p := new(Polynomial).Init([]*mod.Int{
		new(mod.Int).Init(big.NewInt(5), Q),
		new(mod.Int).Init(big.NewInt(1), Q), // x^1
		new(mod.Int).Init(big.NewInt(0), Q), // x^2
		new(mod.Int).Init(big.NewInt(1), Q), // x^3
	})
	assert.Equal(t, "x³ + x¹ + 5", p.ToString())
	o, _ := strconv.ParseInt(p.Eval(new(mod.Int).Init(big.NewInt(3), Q)).String(), 16, 64)
	assert.Equal(t, "35", strconv.Itoa(int(o)))
	o, _ = strconv.ParseInt(p.Eval(new(mod.Int).Init(big.NewInt(10), Q)).String(), 16, 64)
	assert.Equal(t, "1015", strconv.Itoa(int(o)))
	o, _ = strconv.ParseInt(p.Eval(new(mod.Int).Init(big.NewInt(256), Q)).String(), 16, 64)
	assert.Equal(t, "16777477", strconv.Itoa(int(o)))
	o, _ = strconv.ParseInt(p.Eval(new(mod.Int).Init(big.NewInt(50), Q)).String(), 16, 64)
	assert.Equal(t, "125055", strconv.Itoa(int(o)))
	o, _ = strconv.ParseInt(p.Eval(new(mod.Int).Init(big.NewInt(1), Q)).String(), 16, 64)
	assert.Equal(t, "7", strconv.Itoa(int(o)))

}

func TestPolynomial_Eval2(t *testing.T) {
	p := new(Polynomial).Init([]*mod.Int{
		new(mod.Int).Init(big.NewInt(5), Q),
		new(mod.Int).Init(big.NewInt(1), Q), // x^1
		new(mod.Int).Init(big.NewInt(0), Q), // x^2
		new(mod.Int).Init(big.NewInt(1), Q), // x^3
	})
	n := new(Polynomial).Sub(p, new(Polynomial).Init([]*mod.Int{mod.NewInt64(35, Q)})) // p-y
	println(p.ToString())
	println("n", n.ToString())
	// n := p // we can omit y (p(z))
	d := new(Polynomial).Init([]*mod.Int{mod.NewInt64(-3, Q), mod.NewInt64(1, Q)}) // x-z
	q, rem := new(Polynomial).Div(n, d)
	println("mult of q x d : ", new(Polynomial).Mul(q, d).ToString())
	println(q.ToString())
	println(rem.ToString())
}
