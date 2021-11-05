package common

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestBigToBytes(t *testing.T) {
	assert.Equal(t, BigToBytes(BytesToBig([]byte("22")), 10), []byte{50, 50})
	assert.Equal(t, BigToBytes(BytesToBig([]byte("2222222")), 10), []byte{50, 50, 50, 50, 50, 50, 50})

}
