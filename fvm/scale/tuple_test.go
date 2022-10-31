package scale

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTuple(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		a := CompactTuple{Val: []Compact{
			&CompactU32{Val: 3},
			&CompactBool{Val: false},
		}}
		ans, err := a.Encode()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, []byte{12, 0}, ans)
	})

	t.Run("decode", func(t *testing.T) {
		a := CompactTuple{Val: []Compact{
			&CompactU32{},
			&CompactBool{},
		}}
		_, err := a.Decode([]byte{12, 0})
		assert.Nil(t, err)
		assert.Equal(t, uint32(3), a.Val[0].GetVal())
	})
}
