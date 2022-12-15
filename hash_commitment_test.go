package hash_commitment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {

	r := setup()
	value := []byte("")
	c := commit(value, r)
	assert.Equal(t, true, verify(value, r, c))

}
