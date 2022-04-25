package scale

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestLittleEncode_EncodeUInt8(t *testing.T) {
	assert.Equal(t, []byte{3}, BaseIntEncode.EncodeUInt8(3))
}

func TestLittleEncode_UInt8(t *testing.T) {
	assert.Equal(t, uint8(3), BaseIntEncode.UInt8([]byte{3}))
}

func TestLittleEncode_EncodeInt8(t *testing.T) {
	assert.Equal(t, []byte{3}, BaseIntEncode.EncodeInt8(3))
}

func TestLittleEncode_Int8(t *testing.T) {
	assert.Equal(t, int8(-1), BaseIntEncode.Int8([]byte{255}))
}

func TestLittleEncode_EncodeUInt16(t *testing.T) {
	assert.Equal(t, []byte{0x18, 0x0}, BaseIntEncode.EncodeUInt16(24))
}

func TestLittleEncode_UInt16(t *testing.T) {
	assert.Equal(t, uint16(24), BaseIntEncode.UInt16([]byte{0x18, 0x0}))
}

func TestLittleEncode_EncodeInt16(t *testing.T) {
	assert.Equal(t, []byte{0xff, 0xff}, BaseIntEncode.EncodeInt16(-1))
}

func TestLittleEncode_Int16(t *testing.T) {
	assert.Equal(t, int16(-1), BaseIntEncode.Int16([]byte{0xff, 0xff}))
}

func TestLittleEncode_EncodeUInt32(t *testing.T) {
	assert.Equal(t, []byte{0xc8, 0x1, 0x0, 0x0}, BaseIntEncode.EncodeUInt32(456))
}

func TestLittleEncode_UInt32(t *testing.T) {
	assert.Equal(t, uint32(456), BaseIntEncode.UInt32([]byte{0xc8, 0x1, 0x0, 0x0}))
}

func TestLittleEncode_EncodeInt32(t *testing.T) {
	assert.Equal(t, []byte{0xc8, 0x1, 0x0, 0x0}, BaseIntEncode.EncodeInt32(456))
}

func TestLittleEncode_Int32(t *testing.T) {
	assert.Equal(t, int32(456), BaseIntEncode.Int32([]byte{0xc8, 0x1, 0x0, 0x0}))
}

func TestLittleEncode_EncodeUInt64(t *testing.T) {
	assert.Equal(t, []byte{0xc8, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, BaseIntEncode.EncodeUInt64(456))
}

func TestLittleEncode_UInt64(t *testing.T) {
	assert.Equal(t, uint64(456), BaseIntEncode.UInt64([]byte{0xc8, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}))
}

func TestLittleEncode_EncodeInt64(t *testing.T) {
	assert.Equal(t, []byte{0xc8, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, BaseIntEncode.EncodeInt64(456))
}

func TestLittleEncode_Int64(t *testing.T) {
	assert.Equal(t, int64(456), BaseIntEncode.Int64([]byte{0xc8, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}))
}

func TestLittleEncode_EncodeU128(t *testing.T) {
	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, BaseIntEncode.EncodeU128(new(big.Int).SetUint64(1<<63)))
	assert.Equal(t, uint64(1<<63), BaseIntEncode.U128([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).Uint64())
}
