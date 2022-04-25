package scale

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompactArray_Encode(t *testing.T) {
	v := &CompactArray{Val: []Compact{&FixU16{
		Val: uint16(1),
	}, &FixU16{
		Val: uint16(64),
	}}, NextList: []TypeString{Uint16Name}, Len: 2}
	res, err := v.Encode()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []byte{1, 0, 64, 0}, res)
	assert.Equal(t, ArrayName, v.GetType())
}

func TestCompactArray_Decode(t *testing.T) {
	s := &CompactArray{NextList: []TypeString{Uint16Name}, Len: 2}
	_, err := s.Decode([]byte{1, 0, 64, 0})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, uint16(1), s.GetVal().([]Compact)[0].GetVal())
	assert.Equal(t, uint16(64), s.Val[1].GetVal())
}
