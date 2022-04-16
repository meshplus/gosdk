package scale

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBool(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		a := CompactBool{Val: true}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{1}, res)

		b := CompactBool{
			Val: false,
		}
		res2, err := b.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0}, res2)
		assert.Equal(t, BoolName, a.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		a := CompactBool{}
		num, err := a.Decode([]byte{1})
		assert.Nil(t, err)
		assert.Equal(t, 1, num)
		assert.Equal(t, true, a.GetVal())

		num, err = a.Decode([]byte{0})
		assert.Nil(t, err)
		assert.Equal(t, 1, num)
		assert.Equal(t, false, a.GetVal())
	})
}
