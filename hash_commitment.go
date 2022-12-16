package polynomial_commitment_kzg

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func setup() []byte {
	r := make([]byte, 32)
	_, err := rand.Read(r)
	if err != nil {
		fmt.Println("error:", err)
	}
	return r
}

func commit(x []byte, r []byte) []byte {

	// The slice should now contain random bytes instead of only zeroes.
	vR := append(x, r...)
	c := sha256.Sum256(vR)
	return c[:]
}

func verify(x []byte, r []byte, c []byte) bool {

	vR := append(x, r...)
	cc := sha256.Sum256(vR)
	if bytes.Compare(c, cc[:]) == 0 {
		return true
	} else {
		return false
	}

}
