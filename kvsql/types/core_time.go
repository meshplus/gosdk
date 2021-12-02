package types

import (
	"errors"
	"fmt"
	gotime "time"
)

// CoreTime is the internal struct type for Time.
type CoreTime uint64

// ZeroCoreTime is the zero value for TimeInternal type.
var ZeroCoreTime = CoreTime(0)

// IsLeapYear returns if it's leap year.
func (t CoreTime) IsLeapYear() bool {
	return isLeapYear(t.getYear())
}

func isLeapYear(year uint16) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// String implements fmt.Stringer.
func (t CoreTime) String() string {
	return fmt.Sprintf("{%d %d %d %d %d %d %d}", t.getYear(), t.getMonth(), t.getDay(), t.getHour(), t.getMinute(), t.getSecond(), t.getMicrosecond())
}

func (t CoreTime) getYear() uint16 {
	return uint16((uint64(t) & yearBitFieldMask) >> yearBitFieldOffset)
}

func (t *CoreTime) setYear(year uint16) {
	*(*uint64)(t) &= ^yearBitFieldMask
	*(*uint64)(t) |= (uint64(year) << yearBitFieldOffset) & yearBitFieldMask
}

// Year returns the year value.
func (t CoreTime) Year() int {
	return int(t.getYear())
}

func (t CoreTime) getMonth() uint8 {
	return uint8((uint64(t) & monthBitFieldMask) >> monthBitFieldOffset)
}

func (t *CoreTime) setMonth(month uint8) {
	*(*uint64)(t) &= ^monthBitFieldMask
	*(*uint64)(t) |= (uint64(month) << monthBitFieldOffset) & monthBitFieldMask
}

// Month returns the month value.
func (t CoreTime) Month() int {
	return int(t.getMonth())
}

func (t CoreTime) getDay() uint8 {
	return uint8((uint64(t) & dayBitFieldMask) >> dayBitFieldOffset)
}

func (t *CoreTime) setDay(day uint8) {
	*(*uint64)(t) &= ^dayBitFieldMask
	*(*uint64)(t) |= (uint64(day) << dayBitFieldOffset) & dayBitFieldMask
}

// Day returns the day value.
func (t CoreTime) Day() int {
	return int(t.getDay())
}

func (t CoreTime) getHour() uint8 {
	return uint8((uint64(t) & hourBitFieldMask) >> hourBitFieldOffset)
}

func (t *CoreTime) setHour(hour uint8) {
	*(*uint64)(t) &= ^hourBitFieldMask
	*(*uint64)(t) |= (uint64(hour) << hourBitFieldOffset) & hourBitFieldMask
}

// Hour returns the hour value.
func (t CoreTime) Hour() int {
	return int(t.getHour())
}

func (t CoreTime) getMinute() uint8 {
	return uint8((uint64(t) & minuteBitFieldMask) >> minuteBitFieldOffset)
}

func (t *CoreTime) setMinute(minute uint8) {
	*(*uint64)(t) &= ^minuteBitFieldMask
	*(*uint64)(t) |= (uint64(minute) << minuteBitFieldOffset) & minuteBitFieldMask
}

// Minute returns the minute value.
func (t CoreTime) Minute() int {
	return int(t.getMinute())
}

func (t CoreTime) getSecond() uint8 {
	return uint8((uint64(t) & secondBitFieldMask) >> secondBitFieldOffset)
}

func (t *CoreTime) setSecond(second uint8) {
	*(*uint64)(t) &= ^secondBitFieldMask
	*(*uint64)(t) |= (uint64(second) << secondBitFieldOffset) & secondBitFieldMask
}

// Second returns the second value.
func (t CoreTime) Second() int {
	return int(t.getSecond())
}

func (t CoreTime) getMicrosecond() uint32 {
	return uint32((uint64(t) & microsecondBitFieldMask) >> microsecondBitFieldOffset)
}

func (t *CoreTime) setMicrosecond(microsecond uint32) {
	*(*uint64)(t) &= ^microsecondBitFieldMask
	*(*uint64)(t) |= (uint64(microsecond) << microsecondBitFieldOffset) & microsecondBitFieldMask
}

// Microsecond returns the microsecond value.
func (t CoreTime) Microsecond() int {
	return int(t.getMicrosecond())
}

// GoTime converts Time to GoTime.
func (t CoreTime) GoTime(loc *gotime.Location) (gotime.Time, error) {
	// gotime.Time can't represent month 0 or day 0, date contains 0 would be converted to a nearest date,
	// For example, 2006-12-00 00:00:00 would become 2015-11-30 23:59:59.
	year, month, day, hour, minute, second, microsecond := t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Microsecond()
	tm := gotime.Date(year, gotime.Month(month), day, hour, minute, second, microsecond*1000, loc)
	year2, month2, day2 := tm.Date()
	hour2, minute2, second2 := tm.Clock()
	microsec2 := tm.Nanosecond() / 1000
	// This function will check the result, and return an error if it's not the same with the origin input .
	if year2 != year || int(month2) != month || day2 != day ||
		hour2 != hour || minute2 != minute || second2 != second ||
		microsec2 != microsecond {
		return tm, errors.New("value error")
	}
	return tm, nil
}

// datetimeToUint64 converts time value to integer in YYYYMMDDHHMMSS format.
func datetimeToUint64(t CoreTime) uint64 {
	return uint64(t.Year())*1e10 +
		uint64(t.Month())*1e8 +
		uint64(t.Day())*1e6 +
		uint64(t.Hour())*1e4 +
		uint64(t.Minute())*1e2 +
		uint64(t.Second())
}

// compareTime compare two Time.
// return:
//  0: if a == b
//  1: if a > b
// -1: if a < b
func compareTime(a, b CoreTime) int {
	ta := datetimeToUint64(a)
	tb := datetimeToUint64(b)

	switch {
	case ta < tb:
		return -1
	case ta > tb:
		return 1
	}

	switch {
	case a.Microsecond() < b.Microsecond():
		return -1
	case a.Microsecond() > b.Microsecond():
		return 1
	}

	return 0
}

// calcDaynr calculates days since 0000-00-00.
func calcDaynr(year, month, day int) int {
	if year == 0 && month == 0 {
		return 0
	}

	delsum := 365*year + 31*(month-1) + day
	if month <= 2 {
		year--
	} else {
		delsum -= (month*4 + 23) / 10
	}
	temp := ((year/100 + 1) * 3) / 4
	return delsum + year/4 - temp
}

// calcDaysInYear calculates days in one year, it works with 0 <= year <= 99.
func calcDaysInYear(year int) int {
	if (year&3) == 0 && (year%100 != 0 || (year%400 == 0 && (year != 0))) {
		return 366
	}
	return 365
}

// calcWeekday calculates weekday from daynr, returns 0 for Monday, 1 for Tuesday ...
// nolint
func calcWeekday(daynr int, sundayFirstDayOfWeek bool) int {
	daynr += 5
	if sundayFirstDayOfWeek {
		daynr++
	}
	return daynr % 7
}
