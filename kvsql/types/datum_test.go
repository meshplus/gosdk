package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollation(t *testing.T) {
	d := Datum{}
	d.SetCollation("utf8")
	assert.Equal(t, "utf8", d.Collation())
}

func TestDatum_Frac(t *testing.T) {
	d := Datum{}
	d.SetFrac(1)
	assert.Equal(t, 1, d.Frac())
}

func TestLength(t *testing.T) {
	d := Datum{}
	d.SetLength(16)
	assert.Equal(t, 16, d.Length())
}

func TestKind(t *testing.T) {
	d := Datum{
		k: KindNull,
	}
	assert.Equal(t, KindNull, d.Kind())
}

func TestDatum_Interface(t *testing.T) {
	d := Datum{}
	d.SetInterface(20)
	assert.Equal(t, 20, d.GetInterface())
}

func TestDatum_SetNull(t *testing.T) {
	d := Datum{k: KindBytes}
	d.SetNull()
	assert.Equal(t, KindNull, d.Kind())
}

func TestDatum_int64(t *testing.T) {
	d := Datum{k: KindBytes}
	d.SetInt64(10)
	assert.Equal(t, int64(10), d.GetInt64())
	assert.Equal(t, KindInt64, d.Kind())
}

func TestDatum_uint64(t *testing.T) {
	d := Datum{k: KindBytes}
	d.SetUint64(10)
	assert.Equal(t, uint64(10), d.GetUint64())
	assert.Equal(t, KindUint64, d.Kind())
}

func TestDatum_float32(t *testing.T) {
	d := Datum{k: KindBytes}
	d.SetFloat32(10)
	assert.Equal(t, float32(10), d.GetFloat32())
	assert.Equal(t, KindFloat32, d.Kind())
}

func TestDatum_float64(t *testing.T) {
	d := Datum{k: KindBytes}
	d.SetFloat64(10)
	assert.Equal(t, float64(10), d.GetFloat64())
	assert.Equal(t, KindFloat64, d.Kind())
}

func TestDatum_String(t *testing.T) {
	d := Datum{k: KindBytes}
	d.SetString("hello", "")
	assert.Equal(t, "hello", d.GetString())
	assert.Equal(t, KindString, d.Kind())
}

func TestDatum_bytes(t *testing.T) {
	d := Datum{k: KindString}
	d.SetBytes([]byte("hello"))
	assert.Equal(t, []byte("hello"), d.GetBytes())
	assert.Equal(t, KindBytes, d.Kind())
}

func TestDatum_MysqlTime(t *testing.T) {
	d := Datum{k: KindString}
	curTime := Time{}
	curTime.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)
	d.SetMysqlTime(curTime)
	assert.Equal(t, "2004-01-01 01:11:01", d.GetMysqlTime().String())
	assert.Equal(t, KindMysqlTime, d.Kind())
}

func TestDatum_MysqlDuration(t *testing.T) {
	d := Datum{k: KindString}
	d.SetMysqlDuration(Duration{
		Duration: time.Hour,
		Fsp:      6,
	})
	assert.Equal(t, Duration{
		Duration: time.Hour,
		Fsp:      6,
	}, d.GetMysqlDuration())
}

func TestDatum_MysqlEnum(t *testing.T) {
	d := Datum{}

	d.SetMysqlEnum(Enum{
		Name: "hello",
	}, "")
	assert.Equal(t, Enum{
		Name: "hello",
	}, d.GetMysqlEnum())

}

func TestDatum_IsNull(t *testing.T) {
	d := Datum{
		k:         KindInt64,
		decimal:   0,
		length:    0,
		i:         64,
		collation: "",
		b:         nil,
		x:         nil,
	}
	assert.False(t, d.IsNull())
}
