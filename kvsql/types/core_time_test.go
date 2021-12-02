package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCoreTime_String(t *testing.T) {
	assert.Equal(t, "{0 0 0 0 0 0 0}", ZeroCoreTime.String())
}

func TestIsLeapYear(t *testing.T) {
	assert.True(t, isLeapYear(2004))
	assert.False(t, isLeapYear(2005))
	y := FromDate(2004, 1, 1, 0, 0, 0, 0)
	assert.True(t, y.IsLeapYear())

	y = FromDate(2005, 1, 1, 0, 0, 0, 0)
	assert.False(t, y.IsLeapYear())
}

func TestGetYear(t *testing.T) {
	y := FromDate(2004, 1, 1, 0, 0, 0, 0)
	assert.Equal(t, 2004, y.Year())
	assert.Equal(t, 1, y.Month())
	assert.Equal(t, 1, y.Day())
	assert.Equal(t, 0, y.Hour())
	assert.Equal(t, 0, y.Minute())
	assert.Equal(t, 0, y.Second())
	assert.Equal(t, 0, y.Microsecond())
	s := datetimeToUint64(y)
	assert.Equal(t, uint64(20040101000000), s)

	y.setHour(1)
	assert.Equal(t, 1, y.Hour())
	y.setDay(2)
	assert.Equal(t, 2, y.Day())

	y.setYear(2005)
	assert.Equal(t, 2005, y.Year())

	y.setMonth(3)
	assert.Equal(t, 3, y.Month())

	y.setMinute(20)
	assert.Equal(t, 20, y.Minute())

	y.setSecond(30)
	assert.Equal(t, 30, y.Second())

	y.setMicrosecond(90)
	assert.Equal(t, 90, y.Microsecond())
}

func TestCompare(t *testing.T) {
	a := FromDate(2004, 1, 1, 0, 0, 0, 0)
	b := FromDate(2004, 1, 1, 0, 0, 0, 1)

	assert.Equal(t, -1, compareTime(a, b))
	assert.Equal(t, 1, compareTime(b, a))

	a = FromDate(2004, 1, 1, 0, 0, 0, 0)
	b = FromDate(2004, 1, 1, 0, 0, 0, 0)
	assert.Equal(t, 0, compareTime(a, b))

	a = FromDate(2004, 2, 1, 0, 0, 0, 0)
	b = FromDate(2004, 1, 1, 0, 0, 0, 0)
	assert.Equal(t, 1, compareTime(a, b))
	assert.Equal(t, -1, compareTime(b, a))

}

func TestCalcDaynr(t *testing.T) {
	assert.Equal(t, 731946, calcDaynr(2004, 1, 1))
	assert.Equal(t, 0, calcDaynr(0, 0, 1))
	assert.Equal(t, 732037, calcDaynr(2004, 4, 1))
}

func TestCalcDaysInYear(t *testing.T) {
	assert.Equal(t, 366, calcDaysInYear(2004))
	assert.Equal(t, 365, calcDaysInYear(2005))
}

func TestCalcWeekday(t *testing.T) {
	assert.Equal(t, 4, calcWeekday(731946, true))
	assert.Equal(t, 3, calcWeekday(731946, false))
}

func TestCoreTime_GoTime(t *testing.T) {
	a := FromDate(2004, 1, 1, 0, 0, 0, 0)
	res, err := a.GoTime(time.Local)
	assert.Nil(t, err)
	assert.Equal(t, 2004, res.Year())

	a = FromDate(2004, 12, 0, 0, 0, 0, 0)
	res, err = a.GoTime(time.Local)
	assert.NotNil(t, err)
	assert.Equal(t, 2004, res.Year())
	assert.Equal(t, "November", res.Month().String())
}
