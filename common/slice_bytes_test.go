package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestSliceBytes_readByte(t *testing.T) {
	bt := []byte{1, 2, 3, 4, 5, 6}
	sbs := NewSliceBytes(bt)
	b, err := sbs.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(1), b)
	b2, err := sbs.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(2), b2)
	b3, err := sbs.ReadByteN(3)
	assert.Nil(t, err)
	b3[0] = 10
	assert.Equal(t, []byte{10, 4, 5}, b3)
}

func TestSliceBytes_readByte2(t *testing.T) {
	bt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	sbs := NewSliceBytes(bt)
	assert.Equal(t, 8, sbs.Len())
	b, err := sbs.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(1), b)
	sbs.MarkTeeRead()
	b2, err := sbs.ReadByte()
	assert.Nil(t, err)
	assert.Equal(t, byte(2), b2)
	b3, err := sbs.ReadByteN(3)
	assert.Nil(t, err)
	teeRead := sbs.TeeRead()
	b3[0] = 10
	assert.Equal(t, []byte{10, 4, 5}, b3)
	assert.Equal(t, []byte{1, 2, 10, 4, 5, 6, 7, 8}, bt)
	assert.Equal(t, []byte{2, 3, 4, 5}, teeRead)
}

func TestSliceBytes_readByteN(t *testing.T) {
	err := fmt.Errorf("read section id: %w", io.EOF)
	assert.ErrorIs(t, err, io.EOF)
}
