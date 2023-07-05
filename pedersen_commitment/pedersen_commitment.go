package pedersen_commitment

import (
	"bytes"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"math/big"
)

type pedersen_commiter struct {
	G, H *bn256.G1
}

func (pc pedersen_commiter) Setup() {
	var g, h *bn256.G1
	r := rand.Reader

	_, g, _ = bn256.RandomG1(r)
	_, h, _ = bn256.RandomG1(r)
	pc = pedersen_commiter{H: h, G: g}
}

func (pc pedersen_commiter) Commit(message, blingFactor []byte) []byte {
	m := new(big.Int).SetBytes(message)
	r := new(big.Int).SetBytes(blingFactor)
	commitment := new(bn256.G1).Add(pc.G.ScalarMult(pc.G, m), pc.H.ScalarMult(pc.H, r))

	return commitment.Marshal()
}

func (pc pedersen_commiter) Verify(commitment, message, blingFactor []byte) bool {
	var res bool
	if bytes.Compare(pc.Commit(message, blingFactor), commitment) == 1 {
		res = true
	} else {
		res = false
	}

	return res
}
