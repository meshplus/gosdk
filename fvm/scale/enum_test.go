package scale

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnum(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		a := CompactEnum{Val: &FixU8{Val: 42}, index: 0}
		ans, err := a.Encode()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, []byte{0, 42}, ans)
	})

	t.Run("decode", func(t *testing.T) {
		a := CompactEnum{Val: &FixU8{}, index: 0}
		ans, err := a.Decode([]byte{0, 42})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, 2, ans)
		assert.Equal(t, uint8(42), a.GetVal())
	})

	t.Run("decode2-optionNone", func(t *testing.T) {
		a := CompactEnum{
			index: 0,
		}
		ans, err := a.Decode([]byte{0})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, 1, ans)
		assert.Nil(t, a.Val)
	})
}
