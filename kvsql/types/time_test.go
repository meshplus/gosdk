package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromDate(t *testing.T) {
	s := FromDate(2004, 1, 1, 1, 11, 1, 0)
	assert.Equal(t, 2004, s.Year())
	assert.Equal(t, 1, s.YearDay())
	assert.Equal(t, 0, s.Week(0))
	assert.Equal(t, 1, s.Week(1))
	assert.Equal(t, "Thursday", s.Weekday().String())
	a, b := s.YearWeek(0)
	assert.Equal(t, 2003, a)
	assert.Equal(t, 52, b)
	s = FromDate(2004, 0, 1, 1, 11, 1, 0)
	assert.Equal(t, 0, s.Week(0))
	assert.Equal(t, 0, s.YearDay())
	assert.Equal(t, "Monday", s.Weekday().String())
}

func TestTime(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)
	assert.Equal(t, false, a.IsZero())

	assert.Equal(t, "2004-01-01 01:11:01", a.String())

	y, m, s := a.Clock()
	assert.Equal(t, 1, y)
	assert.Equal(t, 11, m)
	assert.Equal(t, 1, s)

}

func TestDateFormat(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)
	dateStr, err := a.DateFormat("%H:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %H:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%b-%d %H:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-Jan-01 01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%M-%d %H:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-January-01 01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%c-%d %H:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-1-01 01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%c-%D %H:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-1-1st 01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%c-%e %H:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-1-1 01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%c-%j %H:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-1-001 01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %k:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 1:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %h:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 01:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %l:%i:%s %p")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 1:11:01 AM", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %r")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 01:11:01 AM", dateStr)

	a.coreTime = FromDate(2004, 1, 1, 13, 11, 1, 0)
	dateStr, err = a.DateFormat("%Y-%m-%d %h:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 01:11:01", dateStr)
	dateStr, err = a.DateFormat("%Y-%m-%d %l:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 1:11:01", dateStr)
	dateStr, err = a.DateFormat("%Y-%m-%d %r")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 01:11:01 PM", dateStr)

	a.coreTime = FromDate(2004, 1, 1, 0, 11, 1, 0)
	dateStr, err = a.DateFormat("%Y-%m-%d %r")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01 AM", dateStr)

	a.coreTime = FromDate(2004, 1, 1, 12, 11, 1, 0)
	dateStr, err = a.DateFormat("%Y-%m-%d %h:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01", dateStr)
	dateStr, err = a.DateFormat("%Y-%m-%d %l:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %l:%i:%s %p")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01 PM", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %r")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01 PM", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %T")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %T.%f")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01.000000", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %T.%g")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01.g", dateStr)

	dateStr, err = a.DateFormat("%y-%m-%d %T")
	assert.Nil(t, err)
	assert.Equal(t, "04-01-01 12:11:01", dateStr)

	dateStr, err = a.DateFormat("%x-%m-%d %l:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01", dateStr)

	dateStr, err = a.DateFormat("%X-%m-%d %l:%i:%s")
	assert.Nil(t, err)
	assert.Equal(t, "2003-01-01 12:11:01", dateStr)

	dateStr, err = a.DateFormat("%Y-%m-%d %T %w")
	assert.Nil(t, err)
	assert.Equal(t, "2004-01-01 12:11:01 4", dateStr)

	dateStr, err = a.DateFormat("%W")
	assert.Nil(t, err)
	assert.Equal(t, "Thursday", dateStr)

	dateStr, err = a.DateFormat("%a")
	assert.Nil(t, err)
	assert.Equal(t, "Thu", dateStr)

	dateStr, err = a.DateFormat("%v")
	assert.Nil(t, err)
	assert.Equal(t, "01", dateStr)

	dateStr, err = a.DateFormat("%V")
	assert.Nil(t, err)
	assert.Equal(t, "52", dateStr)

	dateStr, err = a.DateFormat("%u")
	assert.Nil(t, err)
	assert.Equal(t, "01", dateStr)

	dateStr, err = a.DateFormat("%U")
	assert.Nil(t, err)
	assert.Equal(t, "00", dateStr)

	a.coreTime = FromDate(2004, 14, 1, 1, 11, 1, 0)
	dateStr, err = a.DateFormat("%Y-%b-%d %H:%i:%s")
	assert.NotNil(t, err)
	assert.Equal(t, "", dateStr)
	dateStr, err = a.DateFormat("%Y-%M-%d %H:%i:%s")
	assert.NotNil(t, err)
	assert.Equal(t, "", dateStr)

}

func TestInvalidZero(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 0, 1, 1, 11, 1, 0)
	assert.Equal(t, true, a.InvalidZero())
}

func TestAbbrDayOfMonth(t *testing.T) {
	assert.Equal(t, "st", abbrDayOfMonth(1))
	assert.Equal(t, "nd", abbrDayOfMonth(2))
	assert.Equal(t, "rd", abbrDayOfMonth(3))
	assert.Equal(t, "th", abbrDayOfMonth(4))
}

func TestSetFspTt(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)
	assert.Equal(t, uint8(0), a.getFspTt())
	a.setFspTt(20)
	assert.Equal(t, uint8(4), a.getFspTt())
}

func TestType(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)
	assert.Equal(t, uint8(12), a.Type())

	a.setFspTt(14)
	assert.Equal(t, int8(0), a.Fsp())
	assert.Equal(t, uint8(10), a.Type())
}

func TestSetType(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)

	a.SetType(7)
	assert.Equal(t, uint8(7), a.Type())

	a.SetType(12)
	assert.Equal(t, uint8(12), a.Type())

	a.SetType(10)
	assert.Equal(t, uint8(10), a.Type())

	a.setFspTt(14)
	a.SetType(12)
	assert.Equal(t, uint8(0), a.getFspTt())
}

func TestSetFsp(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)

	a.SetFsp(1)
	assert.Equal(t, int8(1), a.Fsp())

	a.SetFsp(-1)
	assert.Equal(t, int8(0), a.Fsp())

	a.setFspTt(14)
	a.SetFsp(20)
	assert.Equal(t, int8(0), a.Fsp())
}

func TestString(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)
	assert.Equal(t, "2004-01-01 01:11:01", a.String())

	a.SetFsp(10)
	assert.Equal(t, "2004-01-01 01:11:01.00", a.String())

	a.SetType(10)
	assert.Equal(t, "2004-01-01", a.String())
}

func TestCoreTime(t *testing.T) {
	a := Time{}
	a.SetCoreTime(FromDate(2004, 1, 1, 1, 11, 1, 0))
	assert.Equal(t, FromDate(2004, 1, 1, 1, 11, 1, 0), a.CoreTime())
}

func TestToPackedUint(t *testing.T) {
	a := Time{}
	a.coreTime = FromDate(2004, 1, 1, 1, 11, 1, 0)
	res, err := a.ToPackedUint()
	assert.Nil(t, err)
	assert.Equal(t, uint64(1833319171631349760), res)

	a.SetCoreTime(FromDate(0, 0, 0, 0, 0, 0, 0))
	res, err = a.ToPackedUint()
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), res)
}

func TestFromPackedUint(t *testing.T) {
	a := Time{}
	err := a.FromPackedUint(1833319171631349760)
	assert.Nil(t, err)
	assert.Equal(t, FromDate(2004, 1, 1, 1, 11, 1, 0), a.CoreTime())

	err = a.FromPackedUint(0)
	assert.Nil(t, err)
	assert.Equal(t, FromDate(0, 0, 0, 0, 0, 0, 0), a.CoreTime())
}

func TestDuration(t *testing.T) {
	d := Duration{
		Duration: time.Hour,
		Fsp:      6,
	}
	assert.Equal(t, "01:00:00.000000", d.String())
	assert.Equal(t, 1, d.Hour())

	d = Duration{
		Duration: -1,
		Fsp:      6,
	}
	assert.Equal(t, "-00:00:00.000000", d.String())
}
