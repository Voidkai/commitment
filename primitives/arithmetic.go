package primitives

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type Poly struct {
	Coefficient []*big.Int
	Degree      int
}

func NewPoly(coeffs []*big.Int) Poly {
	return Poly{coeffs, len(coeffs)}
}

// Q is the order of the integer field (Zq) that fits inside the snark
var Q, _ = new(big.Int).SetString(
	"21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)

// R is the mod of the finite field
var R, _ = new(big.Int).SetString(
	"21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)

func RandBigInt() (*big.Int, error) {
	maxbits := R.BitLen()
	b := make([]byte, (maxbits/8)-1)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	r := new(big.Int).SetBytes(b)
	rq := new(big.Int).Mod(r, R)

	return rq, nil
}

func ArrayOfZeroes(n int) Poly {
	r := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		r[i] = new(big.Int).SetInt64(0)
	}
	return Poly{r[:], n}
}

//nolint:deadcode,unused
func arrayOfZeroesG1(n int) []*bn256.G1 {
	r := make([]*bn256.G1, n)
	for i := 0; i < n; i++ {
		r[i] = new(bn256.G1).ScalarBaseMult(big.NewInt(0))
	}
	return r[:]
}

//nolint:deadcode,unused
func arrayOfZeroesG2(n int) []*bn256.G2 {
	r := make([]*bn256.G2, n)
	for i := 0; i < n; i++ {
		r[i] = new(bn256.G2).ScalarBaseMult(big.NewInt(0))
	}
	return r[:]
}

func ComparePoly(a, b Poly) bool {
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

//nolint:deadcode,unused
func checkArrayOfZeroes(a Poly) bool {
	z := ArrayOfZeroes(a.Degree)
	return ComparePoly(a, z)
}

func FieldAdd(a, b *big.Int) *big.Int {
	ab := new(big.Int).Add(a, b)
	return ab.Mod(ab, R)
}

func FieldSub(a, b *big.Int) *big.Int {
	ab := new(big.Int).Sub(a, b)
	return new(big.Int).Mod(ab, R)
}

func FieldMul(a, b *big.Int) *big.Int {
	ab := new(big.Int).Mul(a, b)
	return ab.Mod(ab, R)
}

func FieldDiv(a, b *big.Int) *big.Int {
	ab := new(big.Int).Mul(a, new(big.Int).ModInverse(b, R))
	return new(big.Int).Mod(ab, R)
}

func FieldNeg(a *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Neg(a), R)
}

//nolint:deadcode,unused
func FieldInv(a *big.Int) *big.Int {
	return new(big.Int).ModInverse(a, R)
}

func FieldExp(base *big.Int, e *big.Int) *big.Int {
	res := big.NewInt(1)
	rem := new(big.Int).Set(e)
	exp := base

	for !bytes.Equal(rem.Bytes(), big.NewInt(int64(0)).Bytes()) {
		// if BigIsOdd(rem) {
		if rem.Bit(0) == 1 { // .Bit(0) returns 1 when is odd
			res = FieldMul(res, exp)
		}
		exp = FieldMul(exp, exp)
		rem.Rsh(rem, 1)
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// polynomial operation.

func PolynomialAdd(a, b Poly) Poly {
	r := ArrayOfZeroes(max(a.Degree, b.Degree))
	for i := 0; i < a.Degree; i++ {
		r.Coefficient[i] = FieldAdd(r.Coefficient[i], a.Coefficient[i])
	}
	for i := 0; i < b.Degree; i++ {
		r.Coefficient[i] = FieldAdd(r.Coefficient[i], b.Coefficient[i])
	}
	return r
}

func PolynomialSub(a, b Poly) Poly {
	r := ArrayOfZeroes(max(a.Degree, b.Degree))
	for i := 0; i < a.Degree; i++ {
		r.Coefficient[i] = FieldAdd(r.Coefficient[i], a.Coefficient[i])
	}
	for i := 0; i < b.Degree; i++ {
		r.Coefficient[i] = FieldSub(r.Coefficient[i], b.Coefficient[i])
	}
	return r
}

func PolynomialMul(a, b Poly) Poly {
	r := ArrayOfZeroes(a.Degree + b.Degree - 1)
	for i := 0; i < a.Degree; i++ {
		for j := 0; j < b.Degree; j++ {
			r.Coefficient[i+j] = FieldAdd(r.Coefficient[i+j], FieldMul(a.Coefficient[i], b.Coefficient[j]))
		}
	}
	return r
}

func PolynomialDiv(a, b Poly) (Poly, Poly) {
	// https://en.wikipedia.org/wiki/Division_algorithm
	r := ArrayOfZeroes(a.Degree - b.Degree + 1)
	rem := a
	for rem.Degree >= b.Degree {
		l := FieldDiv(rem.Coefficient[rem.Degree-1], b.Coefficient[b.Degree-1])
		pos := rem.Degree - b.Degree
		r.Coefficient[pos] = l
		aux := ArrayOfZeroes(pos)
		aux1 := append(aux.Coefficient, l)
		tempoly := Poly{aux1, len(aux1)}
		aux2 := PolynomialSub(rem, PolynomialMul(b, tempoly))
		rem.Coefficient = aux2.Coefficient[:aux2.Degree-1]
	}
	return r, rem
}

func PolynomialMulByConstant(a Poly, c *big.Int) Poly {
	for i := 0; i < a.Degree; i++ {
		a.Coefficient[i] = FieldMul(a.Coefficient[i], c)
	}
	return a
}
func PolynomialDivByConstant(a Poly, c *big.Int) Poly {
	for i := 0; i < a.Degree; i++ {
		a.Coefficient[i] = FieldDiv(a.Coefficient[i], c)
	}
	return a
}

// polynomialEval evaluates the polinomial over the Finite Field at the given value x
func PolynomialEval(p Poly, x *big.Int) *big.Int {
	r := big.NewInt(int64(0))
	for i := 0; i < p.Degree; i++ {
		xi := FieldExp(x, big.NewInt(int64(i)))
		elem := FieldMul(p.Coefficient[i], xi)
		r = FieldAdd(r, elem)
	}
	return r
}

// newPolZeroAt generates a new polynomial that has value zero at the given value
func NewPolZeroAt(pointPos, totalPoints int, height *big.Int) Poly {
	fac := 1
	for i := 1; i < totalPoints+1; i++ {
		if i != pointPos {
			fac = fac * (pointPos - i)
		}
	}
	facBig := big.NewInt(int64(fac))
	hf := FieldDiv(height, facBig)
	r := Poly{[]*big.Int{hf}, 1}
	for i := 1; i < totalPoints+1; i++ {
		if i != pointPos {
			ineg := big.NewInt(int64(-i))
			b1 := big.NewInt(int64(1))
			r = PolynomialMul(r, Poly{[]*big.Int{ineg, b1}, 2})
		}
	}
	return r
}

// zeroPolynomial returns the zero polynomial:
// z(x) = (x - z_0) (x - z_1) ... (x - z_{k-1})
func ZeroPolynomial(zs []*big.Int) Poly {
	z := Poly{[]*big.Int{FieldNeg(zs[0]), big.NewInt(1)}, 2} // (x - z0)
	for i := 1; i < len(zs); i++ {
		z = PolynomialMul(z, Poly{[]*big.Int{FieldNeg(zs[i]), big.NewInt(1)}, 2}) // (x - zi)
	}
	return z
}

var sNums = map[string]string{
	"0": "⁰",
	"1": "¹",
	"2": "²",
	"3": "³",
	"4": "⁴",
	"5": "⁵",
	"6": "⁶",
	"7": "⁷",
	"8": "⁸",
	"9": "⁹",
}

func intToSNum(n int) string {
	s := strconv.Itoa(n)
	sN := ""
	for i := 0; i < len(s); i++ {
		sN += sNums[string(s[i])]
	}
	return sN
}

// PolynomialToString converts a polynomial represented by a *big.Int array,
// into its string human readable representation
func PolynomialToString(p Poly) string {
	s := ""
	for i := p.Degree - 1; i >= 1; i-- {
		if bytes.Equal(p.Coefficient[i].Bytes(), big.NewInt(1).Bytes()) {
			s += fmt.Sprintf("x%s + ", intToSNum(i))
		} else if !bytes.Equal(p.Coefficient[i].Bytes(), big.NewInt(0).Bytes()) {
			s += fmt.Sprintf("%sx%s + ", p.Coefficient[i], intToSNum(i))
		}
	}
	s += p.Coefficient[0].String()
	return s
}

// LagrangeInterpolation implements the Lagrange interpolation:
// https://en.wikipedia.org/wiki/Lagrange_polynomial
func LagrangeInterpolation(x, y []*big.Int) (Poly, error) {
	// p(x) will be the interpoled polynomial
	// var p []*big.Int
	if len(x) != len(y) {
		return Poly{nil, 0}, fmt.Errorf("len(x)!=len(y): %d, %d", len(x), len(y))
	}
	p := ArrayOfZeroes(len(x))
	k := len(x)

	for j := 0; j < k; j++ {
		// jPol is the Lagrange basis polynomial for each point
		var jPol Poly
		for m := 0; m < k; m++ {
			// if x[m] == x[j] {
			if m == j {
				continue
			}
			// numerator & denominator of the current iteration
			num := Poly{[]*big.Int{FieldNeg(x[m]), big.NewInt(1)}, 2} // (x^1 - x_m)
			den := FieldSub(x[j], x[m])                               // x_j-x_m
			mPol := PolynomialDivByConstant(num, den)
			if jPol.Degree == 0 {
				// first j iteration
				jPol = mPol
				continue
			}
			jPol = PolynomialMul(jPol, mPol)
		}
		p = PolynomialAdd(p, PolynomialMulByConstant(jPol, y[j]))
	}

	return p, nil
}

// TODO add method to 'clean' the polynomial, to remove right-zeroes
