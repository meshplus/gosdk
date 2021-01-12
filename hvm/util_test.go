package hvm

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestBytesToInt32(t *testing.T) {
	bytes := make([]byte, 2)
	bytes[0] = byte(1)
	bytes[1] = byte(2)

	assert.Equal(t, int32(258), BytesToInt32(bytes))
}
