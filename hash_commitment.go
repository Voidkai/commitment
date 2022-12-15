package hash_commitment

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

var b []byte

func commit(x []byte) []byte {
	h := sha256.New()
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
	}
	// The slice should now contain random bytes instead of only zeroes.
	h.Write(append(x, b...))
	c := h.Sum(nil)
	fmt.Printf("%x", h.Sum(nil))
	return c
}

func verify(x []byte, c []byte) bool {

	h := sha256.New()

	h.Write(append(x, b...))
	cc := h.Sum(nil)

	if bytes.Compare(c, cc) == 0 {
		return true
	} else {
		return false
	}

}
