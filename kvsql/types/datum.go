package types

import (
	"github.com/meshplus/gosdk/kvsql/util/hack"
	"math"
	"time"
)

// Kind constants.
const (
	KindNull          byte = 0
	KindInt64         byte = 1
	KindUint64        byte = 2
	KindFloat32       byte = 3
	KindFloat64       byte = 4
	KindString        byte = 5
	KindBytes         byte = 6
	KindMysqlDuration byte = 9
	KindMysqlEnum     byte = 10
	KindMysqlTime     byte = 13
	KindInterface     byte = 14
)

// Datum is a data box holds different kind of data.
// It has better performance and is easier to use than `interface{}`.
// nolint
type Datum struct {
	k       byte   // datum kind.
	decimal uint16 // decimal can hold uint16 values.
	length  uint32 // length can hold uint32 values.
	i       int64  // i can hold int64 uint64 float64 values.
	// 需要规定好支持的字符集
	collation string      // collation hold the collation information for string value.
	b         []byte      // b can hold string or []byte values.
	x         interface{} // x hold all other types.
}

// Collation gets the collation of the datum.
func (d *Datum) Collation() string {
	return d.collation
}

// SetCollation sets the collation of the datum.
func (d *Datum) SetCollation(collation string) {
	d.collation = collation
}

// Frac gets the frac of the datum.
func (d *Datum) Frac() int {
	return int(d.decimal)
}

// SetFrac sets the frac of the datum.
func (d *Datum) SetFrac(frac int) {
	d.decimal = uint16(frac)
}

// Length gets the length of the datum.
func (d *Datum) Length() int {
	return int(d.length)
}

// SetLength sets the length of the datum.
func (d *Datum) SetLength(l int) {
	d.length = uint32(l)
}

// Kind gets the kind of the datum.
func (d *Datum) Kind() byte {
	return d.k
}

// GetInterface gets interface value.
func (d *Datum) GetInterface() interface{} {
	return d.x
}

// SetInterface sets interface to datum.
func (d *Datum) SetInterface(x interface{}) {
	d.k = KindInterface
	d.x = x
}

// SetNull sets datum to nil.
func (d *Datum) SetNull() {
	d.k = KindNull
	d.x = nil
}

// GetInt64 gets int64 value.
func (d *Datum) GetInt64() int64 {
	return d.i
}

// SetInt64 sets int64 value.
func (d *Datum) SetInt64(i int64) {
	d.k = KindInt64
	d.i = i
}

// GetUint64 gets uint64 value.
func (d *Datum) GetUint64() uint64 {
	return uint64(d.i)
}

// SetUint64 sets uint64 value.
func (d *Datum) SetUint64(i uint64) {
	d.k = KindUint64
	d.i = int64(i)
}

// GetFloat32 gets float32 value.
func (d *Datum) GetFloat32() float32 {
	return float32(math.Float64frombits(uint64(d.i)))
}

// SetFloat32 sets float32 value.
func (d *Datum) SetFloat32(f float32) {
	d.k = KindFloat32
	d.i = int64(math.Float64bits(float64(f)))
}

// GetFloat64 gets float64 value.
func (d *Datum) GetFloat64() float64 {
	return math.Float64frombits(uint64(d.i))
}

// SetFloat64 sets float64 value.
func (d *Datum) SetFloat64(f float64) {
	d.k = KindFloat64
	d.i = int64(math.Float64bits(f))
}

// GetString gets string value.
func (d *Datum) GetString() string {
	return string(hack.String(d.b))
}

// SetString sets string value.
func (d *Datum) SetString(s string, collation string) {
	d.k = KindString
	sink(s)
	d.b = hack.Slice(s)
	d.collation = collation
}

// sink prevents s from being allocated on the stack.
var sink = func(s string) {
}

// GetBytes gets bytes value.
func (d *Datum) GetBytes() []byte {
	return d.b
}

// SetBytes sets bytes value to datum.
func (d *Datum) SetBytes(b []byte) {
	d.k = KindBytes
	d.b = b
}

// GetMysqlTime gets types.Time value
func (d *Datum) GetMysqlTime() Time {
	return d.x.(Time)
}

// SetMysqlTime sets types.Time value
func (d *Datum) SetMysqlTime(b Time) {
	d.k = KindMysqlTime
	d.x = b
}

// GetMysqlDuration gets Duration value
func (d *Datum) GetMysqlDuration() Duration {
	return Duration{Duration: time.Duration(d.i), Fsp: int8(d.decimal)}
}

// SetMysqlDuration sets Duration value
func (d *Datum) SetMysqlDuration(b Duration) {
	d.k = KindMysqlDuration
	d.i = int64(b.Duration)
	d.decimal = uint16(b.Fsp)
}

// GetMysqlEnum gets Enum value
func (d *Datum) GetMysqlEnum() Enum {
	str := string(hack.String(d.b))
	return Enum{Name: str}
}

// SetMysqlEnum sets Enum value
func (d *Datum) SetMysqlEnum(b Enum, collation string) {
	d.k = KindMysqlEnum
	sink(b.Name)
	d.collation = collation
	d.b = hack.Slice(b.Name)
}

// IsNull checks if datum is null.
func (d *Datum) IsNull() bool {
	return d.k == KindNull
}
