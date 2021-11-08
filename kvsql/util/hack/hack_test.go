package hack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHack_String(t *testing.T) {
	assert.Equal(t, MutableString(""), String([]byte{}))
	assert.Equal(t, MutableString("hello"), String([]byte("hello")))
}

func TestHack_Slice(t *testing.T) {
	assert.Equal(t, []byte("hello"), Slice("hello"))
}
