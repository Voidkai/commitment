package Polynomial_commitment

import (
	"commitment/primitives"
	"crypto/rand"
	"github.com/drand/kyber/group/mod"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrustedSetup(t *testing.T) {
	ts, _ := NewTrustedSetup(10)
	t.Log(len(ts.Tau2))
}

func BenchmarkNewTrustedSetup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTrustedSetup(100000)
	}
}

func TestBn256PairOperation(t *testing.T) {
	a, _ := rand.Int(rand.Reader, primitives.Q)
	b, _ := rand.Int(rand.Reader, primitives.Q)

	pa := new(bn256.G1).ScalarBaseMult(new(big.Int).Neg(a))
	qb := new(bn256.G2).ScalarBaseMult(b)
	tab := bn256.Pair(pa, qb)

	pb := new(bn256.G1).ScalarBaseMult(new(big.Int).Neg(b))
	qa := new(bn256.G2).ScalarBaseMult(a)
	tba := bn256.Pair(pb, qa)

	assert.True(t, tab.String() == tba.String())
}

func TestSimpleFlow(t *testing.T) {
	// p(x) = x^3 + x + 5
	p := new(primitives.Polynomial).Init([]*mod.Int{
		mod.NewInt64(5, primitives.Q),
		mod.NewInt64(1, primitives.Q), // x^1
		mod.NewInt64(0, primitives.Q), // x^2
		mod.NewInt64(1, primitives.Q), // x^3
	})
	assert.Equal(t, "x³ + x¹ + 5", p.ToString())

	// TrustedSetup
	ts, err := NewTrustedSetup(p.Degree)
	assert.Nil(t, err)

	// Commit
	c := Commit(ts, p)

	// p(z)=y --> p(3)=35
	z := mod.NewInt64(3, primitives.Q)
	y := mod.NewInt64(35, primitives.Q)

	// z & y: to prove an evaluation p(z)=y
	proof, err := EvaluationProof(ts, p, z, y)
	assert.Nil(t, err)

	v := Verify(ts, c, proof, z, y)
	assert.True(t, v)

	v = Verify(ts, c, proof, new(mod.Int).Init(big.NewInt(4), primitives.Q), y)
	assert.False(t, v)
}

func TestBatchProof(t *testing.T) {
	// p(x) = 10x^4+x^3 + x + 5
	p := new(primitives.Polynomial).Init([]*mod.Int{
		mod.NewInt64(5, primitives.Q),
		mod.NewInt64(1, primitives.Q),  // x^1
		mod.NewInt64(0, primitives.Q),  // x^2
		mod.NewInt64(1, primitives.Q),  // x^3
		mod.NewInt64(10, primitives.Q), // x^4
	})
	assert.Equal(t, "10x⁴ + x³ + x¹ + 5", p.ToString())

	// TrustedSetup
	ts, err := NewTrustedSetup(p.Degree)
	assert.Nil(t, err)

	// Commit
	c := Commit(ts, p)

	// 1st point: p(z)=y --> p(3)=35
	z0 := mod.NewInt64(3, primitives.Q)
	y0 := p.Eval(z0)

	// 2nd point: p(10)=1015
	z1 := mod.NewInt64(10, primitives.Q)
	y1 := p.Eval(z1)

	// 3nd point: p(256)=16777477
	z2 := mod.NewInt64(256, primitives.Q)
	y2 := p.Eval(z2)

	zs := []*mod.Int{z0, z1, z2}
	ys := []*mod.Int{y0, y1, y2}

	// prove an evaluation of the multiple z_i & y_i
	proof, err := EvaluationBatchProof(ts, p, zs, ys)
	assert.Nil(t, err)

	// batch proof verification
	v := VerifyBatchProof(ts, c, proof, zs, ys)
	assert.True(t, v)

	// changing order of the points to be verified
	zs[0], zs[1], zs[2] = zs[1], zs[2], zs[0]
	ys[0], ys[1], ys[2] = ys[1], ys[2], ys[0]
	v = VerifyBatchProof(ts, c, proof, zs, ys)
	assert.True(t, v)

	// change a value of zs and check that verification fails
	zs[0] = mod.NewInt64(2, primitives.Q)
	v = VerifyBatchProof(ts, c, proof, zs, ys)
	assert.False(t, v)

	// using a value that is not in the evaluation proof should generate a
	// proof that will not correctly be verified
	zs = []*mod.Int{z0, z1, z2}
	ys = []*mod.Int{y0, y1, y2}
	proof, err = EvaluationBatchProof(ts, p, zs, ys)
	assert.Nil(t, err)
	zs[2] = mod.NewInt64(2500, primitives.Q)
	ys[2] = p.Eval(zs[2])
	v = VerifyBatchProof(ts, c, proof, zs, ys)
	assert.False(t, v)
}
