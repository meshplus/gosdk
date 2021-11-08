package kvsql

import (
	"github.com/meshplus/gosdk/kvsql/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumnNew(t *testing.T) {
	col := NewColumnWithData([]byte("hello"), []byte{1, 2}, []int64{1})
	assert.Equal(t, []byte("hello"), col.data)
}

func TestAppendMultiSameNullBitmap(t *testing.T) {
	col := NewColumnWithData([]byte("hello"), []byte{1, 2}, []int64{1})
	col.appendMultiSameNullBitmap(false, 5)
	assert.Equal(t, []byte{1, 2}, col.nullBitmap)

	col.appendMultiSameNullBitmap(true, 5)
	assert.Equal(t, []byte{255, 2}, col.nullBitmap)

	col.appendMultiSameNullBitmap(true, 25)
	assert.Equal(t, []byte{255, 2, 255, 1}, col.nullBitmap)
}

func TestAppendNullBitmap(t *testing.T) {
	col := NewColumnWithData([]byte("hello"), []byte{1, 2}, []int64{1})
	col.length = 16
	col.appendNullBitmap(true)
	assert.Equal(t, []byte{1, 2, 1}, col.nullBitmap)
}

func TestColumn_GetInt8(t *testing.T) {
	//[104 101 108 108 111]
	col := NewColumnWithData([]byte("hello"), []byte{1, 2}, []int64{1})
	assert.Equal(t, int8(101), col.getInt8(1))
}

func TestColumn_GetUint8(t *testing.T) {
	col := NewColumnWithData([]byte("hello"), []byte{1, 2}, []int64{1})
	assert.Equal(t, uint8(101), col.getUint8(1))
}

func TestColumn_GetInt16(t *testing.T) {
	col := NewColumnWithData([]byte{0, 224, 0, 1}, []byte{1}, []int64{1})
	assert.Equal(t, int16(256), col.getInt16(1))
}

func TestColumn_GetUint16(t *testing.T) {
	col := NewColumnWithData([]byte{0, 224, 0, 1}, []byte{1, 2}, []int64{1})
	assert.Equal(t, uint16(256), col.getUint16(1))
}

func TestColumn_GetInt32(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 0, 0, 1, 0, 0}, []byte{1, 2}, []int64{1})
	assert.Equal(t, int32(256), col.getInt32(1))
}

func TestColumn_GetUint32(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 0, 0, 1, 0, 0}, []byte{1, 2}, []int64{1})
	assert.Equal(t, uint32(256), col.getUint32(1))
}

func TestColumn_GetInt64(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, []byte{1, 2}, []int64{1})
	assert.Equal(t, int64(1), col.GetInt64(1))
}

func TestColumn_GetUint64(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, []byte{1, 2}, []int64{1})
	assert.Equal(t, uint64(1), col.GetUint64(1))
}

func TestColumn_GetFloat32(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 0, 0, 0, 0, 64}, []byte{1, 2}, []int64{1})
	assert.Equal(t, float32(2), col.GetFloat32(1))
}

func TestColumn_GetFloat64(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64}, []byte{1, 2}, []int64{1})
	assert.Equal(t, float64(2), col.GetFloat64(1))
}

func TestColumn_GetString(t *testing.T) {
	col := NewColumnWithData([]byte("hello"), []byte{1}, []int64{0, 5})
	assert.Equal(t, "hello", col.GetString(0))
}

func TestColumn_GetBytes(t *testing.T) {
	col := NewColumnWithData([]byte("hello"), []byte{1}, []int64{0, 5})
	assert.Equal(t, []byte("hello"), col.GetBytes(0))
}

func TestColumn_GetTime(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, []byte{1}, []int64{0, 5})
	s := types.FromDate(0, 0, 0, 0, 0, 1, 0)
	curT := types.Time{}
	curT.SetCoreTime(s)
	assert.Equal(t, curT, col.GetTime(0))
}

func TestColumn_GetDuration(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 0, 0, 0, 0, 0, 60, 0, 0, 0, 0, 0, 0, 0}, []byte{1}, []int64{0, 5})
	assert.Equal(t, types.Duration{
		Duration: 60,
		Fsp:      0,
	}, col.GetDuration(1, 0))
}

func TestColumn_GetEnum(t *testing.T) {
	col := NewColumnWithData([]byte{104, 101, 108, 108, 111}, []byte{1}, []int64{0, 5})
	assert.Equal(t, "hello", col.GetEnum(0).String())
}

func TestColumn_reset(t *testing.T) {
	col := NewColumnWithData([]byte{0, 0, 0, 0, 0, 0, 0, 0, 104, 101, 108, 108, 111}, []byte{1}, []int64{0, 13})
	col.reset()
	assert.Equal(t, []byte{}, col.nullBitmap)
}
