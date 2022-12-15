package hash_commitment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {
	c := commit([]byte("Hello world"))
	assert.Equal(t, true, verify([]byte("hello world"), c))
}
