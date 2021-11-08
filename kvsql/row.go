package kvsql

import (
	"github.com/meshplus/gosdk/kvsql/types"
)

// Row represents a row of data, can be used to access values.
type Row struct {
	c   *Chunk
	idx int
}

// Chunk returns the Chunk which the row belongs to.
func (r Row) Chunk() *Chunk {
	return r.c
}

// IsEmpty returns true if the Row is empty.
func (r Row) IsEmpty() bool {
	return r == Row{}
}

// Idx returns the row index of Chunk.
func (r Row) Idx() int {
	return r.idx
}

// Len returns the number of values in the row.
func (r Row) Len() int {
	return r.c.NumCols()
}

/**************** these should not be used in normal, supported for decode recode set ****************/

// GetInt8 returns the int8 value with the colIdx.
func (r Row) GetInt8(colIdx int) int8 {
	return r.c.columns[colIdx].getInt8(r.idx)
}

// GetUint8 returns the uint8 value with the colIdx.
func (r Row) GetUint8(colIdx int) uint8 {
	return r.c.columns[colIdx].getUint8(r.idx)
}

// GetInt16 returns the int16 value with the colIdx.
func (r Row) GetInt16(colIdx int) int16 {
	return r.c.columns[colIdx].getInt16(r.idx)
}

// GetUint16 returns the uint16 value with the colIdx.
func (r Row) GetUint16(colIdx int) uint16 {
	return r.c.columns[colIdx].getUint16(r.idx)
}

// GetInt32 returns the int32 value with the colIdx.
func (r Row) GetInt32(colIdx int) int32 {
	return r.c.columns[colIdx].getInt32(r.idx)
}

// GetUint32 returns the uint8 value with the colIdx.
func (r Row) GetUint32(colIdx int) uint32 {
	return r.c.columns[colIdx].getUint32(r.idx)
}

// GetInt64 returns the int64 value with the colIdx.
func (r Row) GetInt64(colIdx int) int64 {
	return r.c.columns[colIdx].GetInt64(r.idx)
}

// GetUint64 returns the uint64 value with the colIdx.
func (r Row) GetUint64(colIdx int) uint64 {
	return r.c.columns[colIdx].GetUint64(r.idx)
}

// GetFloat32 returns the float32 value with the colIdx.
func (r Row) GetFloat32(colIdx int) float32 {
	return r.c.columns[colIdx].GetFloat32(r.idx)
}

// GetFloat64 returns the float64 value with the colIdx.
func (r Row) GetFloat64(colIdx int) float64 {
	return r.c.columns[colIdx].GetFloat64(r.idx)
}

// GetString returns the string value with the colIdx.
func (r Row) GetString(colIdx int) string {
	return r.c.columns[colIdx].GetString(r.idx)
}

// GetBytes returns the bytes value with the colIdx.
func (r Row) GetBytes(colIdx int) []byte {
	return r.c.columns[colIdx].GetBytes(r.idx)
}

// GetTime returns the Time value with the colIdx.
func (r Row) GetTime(colIdx int) types.Time {
	return r.c.columns[colIdx].GetTime(r.idx)
}

// GetDuration returns the Duration value with the colIdx.
func (r Row) GetDuration(colIdx int, fillFsp int) types.Duration {
	return r.c.columns[colIdx].GetDuration(r.idx, fillFsp)
}

func (r Row) getNameValue(colIdx int) (string, uint64) {
	return r.c.columns[colIdx].getNameValue(r.idx)
}

// GetEnum returns the Enum value with the colIdx.
func (r Row) GetEnum(colIdx int) types.Enum {
	return r.c.columns[colIdx].GetEnum(r.idx)
}

// IsNull returns if the datum in the chunk.Row is null.
func (r Row) IsNull(colIdx int) bool {
	return r.c.columns[colIdx].IsNull(r.idx)
}

// GetDecimal returns the MyDecimal value with the colIdx.
func (r Row) GetDecimal(colIdx int) *types.MyDecimal {
	return r.c.columns[colIdx].GetDecimal(r.idx)
}

// GetSet returns the Set value with the colIdx.
func (r Row) GetSet(colIdx int) types.Set {
	return r.c.columns[colIdx].GetSet(r.idx)
}
