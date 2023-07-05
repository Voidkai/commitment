package commitment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {

	hc := new(hash_commiter)
	r := hc.Setup()
	value := []byte("")
	c := hc.Commit(value, r)
	assert.Equal(t, true, hc.Verify(value, r, c))

}
