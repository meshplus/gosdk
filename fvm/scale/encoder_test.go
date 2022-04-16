package scale

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeVec(t *testing.T) {
	res, err := encode(&CompactVec{Val: []Compact{
		&FixU16{
			Val: uint16(1),
		}, &FixU16{
			Val: uint16(2),
		}, &FixU16{
			Val: uint16(3),
		}}, NextList: []TypeString{Uint16Name}})
	assert.Nil(t, err)
	assert.Equal(t, []byte{12, 1, 0, 2, 0, 3, 0}, res)
}

func TestEncodeVecErr(t *testing.T) {
	_, err := encode("string")
	assert.NotNil(t, err)
}
