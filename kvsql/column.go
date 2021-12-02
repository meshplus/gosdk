package kvsql

import (
	"github.com/meshplus/gosdk/kvsql/types"
	"github.com/meshplus/gosdk/kvsql/util/hack"
	"time"
	"unsafe"
)

// nolint
// will use later
const (
	SizeInt64      = int(unsafe.Sizeof(int64(0)))
	SizeUint64     = int(unsafe.Sizeof(uint64(0)))
	SizeFloat32    = int(unsafe.Sizeof(float32(0)))
	SizeFloat64    = int(unsafe.Sizeof(float64(0)))
	SizeGoDuration = int(unsafe.Sizeof(time.Duration(0)))
	SizeTime       = int(unsafe.Sizeof(types.ZeroTime))
)

// Column stores one column of data in Apache Arrow format.
// See https://arrow.apache.org/docs/memory_layout.html
type Column struct {
	length     int
	nullBitmap []byte // bit 0 is null, 1 is not null
	offsets    []int64
	data       []byte
	elemBuf    []byte
}

// NewColumnWithData create a new column with data、nullBitmap and offsets
func NewColumnWithData(data []byte, nullBitmap []byte, offsets []int64) *Column {
	return &Column{
		data:       data,
		nullBitmap: nullBitmap,
		offsets:    offsets,
	}
}

// appendMultiSameNullBitmap appends multiple same bit value to `nullBitMap`.
// notNull means not null.
// num means the number of bits that should be appended.
func (c *Column) appendMultiSameNullBitmap(notNull bool, num int) {
	numNewBytes := ((c.length + num + 7) >> 3) - len(c.nullBitmap)
	b := byte(0)
	if notNull {
		b = 0xff
	}
	for i := 0; i < numNewBytes; i++ {
		c.nullBitmap = append(c.nullBitmap, b)
	}
	if !notNull {
		return
	}
	// 1. Set all the remaining bits in the last slot of old c.numBitMap to 1.
	numRemainingBits := uint(c.length % 8)
	bitMask := byte(^((1 << numRemainingBits) - 1))
	c.nullBitmap[c.length/8] |= bitMask
	// 2. Set all the redundant bits in the last slot of new c.numBitMap to 0.
	numRedundantBits := uint(len(c.nullBitmap)*8 - c.length - num)
	bitMask = byte(1<<(8-numRedundantBits)) - 1
	c.nullBitmap[len(c.nullBitmap)-1] &= bitMask
}

func (c *Column) appendNullBitmap(notNull bool) {
	idx := c.length >> 3
	if idx >= len(c.nullBitmap) {
		c.nullBitmap = append(c.nullBitmap, 0)
	}
	if notNull {
		pos := uint(c.length) & 7
		c.nullBitmap[idx] |= byte(1 << pos)
	}
}

/**************** these should not be used in normal, supported for decode recode set ****************/

// getInt8 returns the int8 in the specific row.
func (c *Column) getInt8(rowID int) int8 {
	return *(*int8)(unsafe.Pointer(&c.data[rowID]))
}

// getUint8 returns the uint8 in the specific row.
func (c *Column) getUint8(rowID int) uint8 {
	return *(*uint8)(unsafe.Pointer(&c.data[rowID]))
}

// getInt16 returns the int16 in the specific row.
func (c *Column) getInt16(rowID int) int16 {
	return *(*int16)(unsafe.Pointer(&c.data[rowID*2]))
}

// getUint16 returns the uint16 in the specific row.
func (c *Column) getUint16(rowID int) uint16 {
	return *(*uint16)(unsafe.Pointer(&c.data[rowID*2]))
}

// getInt32 returns the int32 in the specific row.
func (c *Column) getInt32(rowID int) int32 {
	return *(*int32)(unsafe.Pointer(&c.data[rowID*4]))
}

// getUint32 returns the uint32 in the specific row.
func (c *Column) getUint32(rowID int) uint32 {
	return *(*uint32)(unsafe.Pointer(&c.data[rowID*4]))
}

/*************************************************************************************************/

// GetInt64 returns the int64 in the specific row.
func (c *Column) GetInt64(rowID int) int64 {
	return *(*int64)(unsafe.Pointer(&c.data[rowID*8]))
}

// GetUint64 returns the uint64 in the specific row.
func (c *Column) GetUint64(rowID int) uint64 {
	return *(*uint64)(unsafe.Pointer(&c.data[rowID*8]))
}

// GetFloat32 returns the float32 in the specific row.
func (c *Column) GetFloat32(rowID int) float32 {
	return *(*float32)(unsafe.Pointer(&c.data[rowID*4]))
}

// GetFloat64 returns the float64 in the specific row.
func (c *Column) GetFloat64(rowID int) float64 {
	return *(*float64)(unsafe.Pointer(&c.data[rowID*8]))
}

// GetString returns the string in the specific row.
func (c *Column) GetString(rowID int) string {
	return string(hack.String(c.data[c.offsets[rowID]:c.offsets[rowID+1]]))
}

// GetBytes returns the byte slice in the specific row.
func (c *Column) GetBytes(rowID int) []byte {
	return c.data[c.offsets[rowID]:c.offsets[rowID+1]]
}

// GetTime returns the Time in the specific row.
func (c *Column) GetTime(rowID int) types.Time {
	return *(*types.Time)(unsafe.Pointer(&c.data[rowID*SizeTime]))
}

// GetDuration returns the Duration in the specific row.
func (c *Column) GetDuration(rowID int, fillFsp int) types.Duration {
	dur := *(*int64)(unsafe.Pointer(&c.data[rowID*8]))
	return types.Duration{Duration: time.Duration(dur), Fsp: int8(fillFsp)}
}

// GetDecimal returns the decimal in the specific row.
func (c *Column) GetDecimal(rowID int) *types.MyDecimal {
	d := new(types.MyDecimal)
	d.FromString(c.data[c.offsets[rowID]:c.offsets[rowID+1]])
	return d
}

// GetEnum returns the Enum in the specific row.
func (c *Column) GetEnum(rowID int) types.Enum {
	name, _ := c.getNameValue(rowID)
	return types.Enum{Name: name}
}

// GetSet returns the Set in the specific row.
func (c *Column) GetSet(rowID int) types.Set {
	name, _ := c.getNameValue(rowID)
	return types.Set{Name: name}
}

func (c *Column) getNameValue(rowID int) (string, uint64) {
	start, end := c.offsets[rowID], c.offsets[rowID+1]
	if start == end {
		return "", 0
	}
	return string(hack.String(c.data[start:end])), 0
}

// reset resets the underlying data of this Column but doesn't modify its data type.
func (c *Column) reset() {
	c.length = 0
	c.nullBitmap = c.nullBitmap[:0]
	if len(c.offsets) > 0 {
		// The first offset is always 0, it makes slicing the data easier, we need to keep it.
		c.offsets = c.offsets[:1]
	}
	c.data = c.data[:0]
}

// IsNull returns if this row is null.
func (c *Column) IsNull(rowIdx int) bool {
	nullByte := c.nullBitmap[rowIdx/8]
	return nullByte&(1<<(uint(rowIdx)&7)) == 0
}

func (c *Column) isFixed() bool {
	return c.elemBuf != nil
}

func (c *Column) appendToData(b []byte) {
	c.data = append(c.data, b...)
}
