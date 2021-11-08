package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/meshplus/gosdk/kvsql/mysql"
	"math"
	"strconv"
	gotime "time"
)

// Time format without fractional seconds precision.
const (
	TimeFormat = "2006-01-02 15:04:05"
	// TimeFSPFormat is time format with fractional seconds precision.
	TimeFSPFormat = "2006-01-02 15:04:05.000000"
)
const (
	// MinYear is the minimum for mysql year type.
	MinYear int16 = 1901
	// MaxYear is the maximum for mysql year type.
	MaxYear int16 = 2155
	// MaxDuration is the maximum for duration.
	MaxDuration int64 = 838*10000 + 59*100 + 59
	// MinTime is the minimum for mysql time type.
	MinTime = -gotime.Duration(838*3600+59*60+59) * gotime.Second
	// MaxTime is the maximum for mysql time type.
	MaxTime = gotime.Duration(838*3600+59*60+59) * gotime.Second
	// ZeroDatetimeStr is the string representation of a zero datetime.
	ZeroDatetimeStr = "0000-00-00 00:00:00"
	// ZeroDateStr is the string representation of a zero date.
	ZeroDateStr = "0000-00-00"

	// TimeMaxHour is the max hour for mysql time type.
	TimeMaxHour = 838
	// TimeMaxMinute is the max minute for mysql time type.
	TimeMaxMinute = 59
	// TimeMaxSecond is the max second for mysql time type.
	TimeMaxSecond = 59
	// TimeMaxValue is the maximum value for mysql time type.
	TimeMaxValue = TimeMaxHour*10000 + TimeMaxMinute*100 + TimeMaxSecond
	// TimeMaxValueSeconds is the maximum second value for mysql time type.
	TimeMaxValueSeconds = TimeMaxHour*3600 + TimeMaxMinute*60 + TimeMaxSecond
)

const (
	// Core time bit fields.
	yearBitFieldOffset, yearBitFieldWidth               uint64 = 50, 14
	monthBitFieldOffset, monthBitFieldWidth             uint64 = 46, 4
	dayBitFieldOffset, dayBitFieldWidth                 uint64 = 41, 5
	hourBitFieldOffset, hourBitFieldWidth               uint64 = 36, 5
	minuteBitFieldOffset, minuteBitFieldWidth           uint64 = 30, 6
	secondBitFieldOffset, secondBitFieldWidth           uint64 = 24, 6
	microsecondBitFieldOffset, microsecondBitFieldWidth uint64 = 4, 20

	// fspTt bit field.
	// `fspTt` format:
	// | fsp: 3 bits | type: 1 bit |
	// When `fsp` is valid (in range [0, 6]):
	// 1. `type` bit 0 represent `DateTime`
	// 2. `type` bit 1 represent `Timestamp`
	//
	// Since s`Date` does not require `fsp`, we could use `fspTt` == 0b1110 to represent it.
	fspTtBitFieldOffset, fspTtBitFieldWidth uint64 = 0, 4

	yearBitFieldMask        uint64 = ((1 << yearBitFieldWidth) - 1) << yearBitFieldOffset
	monthBitFieldMask       uint64 = ((1 << monthBitFieldWidth) - 1) << monthBitFieldOffset
	dayBitFieldMask         uint64 = ((1 << dayBitFieldWidth) - 1) << dayBitFieldOffset
	hourBitFieldMask        uint64 = ((1 << hourBitFieldWidth) - 1) << hourBitFieldOffset
	minuteBitFieldMask      uint64 = ((1 << minuteBitFieldWidth) - 1) << minuteBitFieldOffset
	secondBitFieldMask      uint64 = ((1 << secondBitFieldWidth) - 1) << secondBitFieldOffset
	microsecondBitFieldMask uint64 = ((1 << microsecondBitFieldWidth) - 1) << microsecondBitFieldOffset
	fspTtBitFieldMask       uint64 = ((1 << fspTtBitFieldWidth) - 1) << fspTtBitFieldOffset

	fspTtForDate         uint8  = 0b1110
	fspBitFieldMask      uint64 = 0b1110
	coreTimeBitFieldMask        = ^fspTtBitFieldMask
)

type weekBehaviour uint

const (
	// weekBehaviourMondayFirst set Monday as first day of week; otherwise Sunday is first day of week
	weekBehaviourMondayFirst weekBehaviour = 1 << iota
	// If set, Week is in range 1-53, otherwise Week is in range 0-53.
	// Note that this flag is only relevant if WEEK_JANUARY is not set.
	weekBehaviourYear
	// If not set, Weeks are numbered according to ISO 8601:1988.
	// If set, the week that contains the first 'first-day-of-week' is week 1.
	weekBehaviourFirstWeekday
)

var (

	// ZeroCoreTime is the zero value for Time type.
	ZeroTime = Time{}
)

var (
	// MonthNames lists names of months, which are used in builtin time function `monthname`.
	MonthNames = []string{
		"January", "February",
		"March", "April",
		"May", "June",
		"July", "August",
		"September", "October",
		"November", "December",
	}
)

var abbrevWeekdayName = []string{
	"Sun", "Mon", "Tue",
	"Wed", "Thu", "Fri", "Sat",
}

// coreTime is an alias to CoreTime, embedd in Time.
type coreTime = CoreTime

// Time is the struct for handling datetime, timestamp and date.
type Time struct {
	coreTime
}

// Clock returns the hour, minute, and second within the day specified by t.
func (t Time) Clock() (hour int, minute int, second int) {
	return t.Hour(), t.Minute(), t.Second()
}

func (t Time) getFspTt() uint8 {
	return uint8(uint64(t.coreTime) & fspTtBitFieldMask)
}

func (t *Time) setFspTt(fspTt uint8) {
	*(*uint64)(&t.coreTime) &= ^(fspTtBitFieldMask)
	*(*uint64)(&t.coreTime) |= uint64(fspTt)
}

// Type returns type value.
func (t Time) Type() uint8 {
	if t.getFspTt() == fspTtForDate {
		return mysql.TypeDate
	}
	if uint64(t.coreTime)&1 == 1 {
		return mysql.TypeTimestamp
	}
	return mysql.TypeDatetime
}

// SetType updates the type in Time.
// Only DateTime/Date/Timestamp is valid.
func (t *Time) SetType(tp uint8) {
	fspTt := t.getFspTt()
	if fspTt == fspTtForDate && tp != mysql.TypeDate {
		fspTt = 0
	}
	switch tp {
	case mysql.TypeDate:
		fspTt = fspTtForDate
	case mysql.TypeTimestamp:
		fspTt |= 1
	case mysql.TypeDatetime:
		fspTt &= ^(uint8(1))
	}
	t.setFspTt(fspTt)
}

// Fsp returns fsp value.
func (t Time) Fsp() int8 {
	fspTt := t.getFspTt()
	if fspTt == fspTtForDate {
		return 0
	}
	return int8(fspTt >> 1)
}

// SetFsp updates the fsp in Time.
func (t *Time) SetFsp(fsp int8) {
	if t.getFspTt() == fspTtForDate {
		return
	}
	if fsp == UnspecifiedFsp {
		fsp = DefaultFsp
	}
	*(*uint64)(&t.coreTime) &= ^(fspBitFieldMask)
	*(*uint64)(&t.coreTime) |= (uint64(fsp) << 1)
}

// CoreTime returns core time.
func (t Time) CoreTime() CoreTime {
	return CoreTime(uint64(t.coreTime) & coreTimeBitFieldMask)
}

// SetCoreTime updates core time.
func (t *Time) SetCoreTime(ct CoreTime) {
	*(*uint64)(&t.coreTime) &= ^coreTimeBitFieldMask
	*(*uint64)(&t.coreTime) |= (uint64(ct) & coreTimeBitFieldMask)
}

func (t Time) String() string {
	if t.Type() == mysql.TypeDate {
		// We control the format, so no error would occur.
		str, _ := t.DateFormat("%Y-%m-%d")
		return str
	}

	str, _ := t.DateFormat("%Y-%m-%d %H:%i:%s")
	fsp := t.Fsp()
	if fsp > 0 {
		tmp := fmt.Sprintf(".%06d", t.Microsecond())
		str = str + tmp[:1+fsp]
	}

	return str
}

// IsZero returns a boolean indicating whether the time is equal to ZeroCoreTime.
func (t Time) IsZero() bool {
	return compareTime(t.coreTime, ZeroCoreTime) == 0
}

// InvalidZero returns a boolean indicating whether the month or day is zero.
func (t Time) InvalidZero() bool {
	return t.Month() == 0 || t.Day() == 0
}

// ToPackedUint encodes Time to a packed uint64 value.
//
//    1 bit  0
//   17 bits year*13+month   (year 0-9999, month 0-12)
//    5 bits day             (0-31)
//    5 bits hour            (0-23)
//    6 bits minute          (0-59)
//    6 bits second          (0-59)
//   24 bits microseconds    (0-999999)
//
//   Total: 64 bits = 8 bytes
//
//   0YYYYYYY.YYYYYYYY.YYdddddh.hhhhmmmm.mmssssss.ffffffff.ffffffff.ffffffff
//
func (t Time) ToPackedUint() (uint64, error) {
	tm := t
	if t.IsZero() {
		return 0, nil
	}
	year, month, day := tm.Year(), tm.Month(), tm.Day()
	hour, minute, sec := tm.Hour(), tm.Minute(), tm.Second()
	ymd := uint64(((year*13 + month) << 5) | day)
	hms := uint64(hour<<12 | minute<<6 | sec)
	micro := uint64(tm.Microsecond())
	return ((ymd<<17 | hms) << 24) | micro, nil
}

// FromPackedUint decodes Time from a packed uint64 value.
func (t *Time) FromPackedUint(packed uint64) error {
	if packed == 0 {
		t.SetCoreTime(ZeroCoreTime)
		return nil
	}
	ymdhms := packed >> 24
	ymd := ymdhms >> 17
	day := int(ymd & (1<<5 - 1))
	ym := ymd >> 5
	month := int(ym % 13)
	year := int(ym / 13)

	hms := ymdhms & (1<<17 - 1)
	second := int(hms & (1<<6 - 1))
	minute := int((hms >> 6) & (1<<6 - 1))
	hour := int(hms >> 12)
	microsec := int(packed % (1 << 24))

	t.SetCoreTime(FromDate(year, month, day, hour, minute, second, microsec))

	return nil
}

func (v weekBehaviour) test(flag weekBehaviour) bool {
	return (v & flag) != 0
}

func weekMode(mode int) weekBehaviour {
	weekFormat := weekBehaviour(mode & 7)
	if (weekFormat & weekBehaviourMondayFirst) == 0 {
		weekFormat ^= weekBehaviourFirstWeekday
	}
	return weekFormat
}

// calcWeek calculates week and year for the time.
func calcWeek(t CoreTime, wb weekBehaviour) (year int, week int) {
	var days int
	ty, tm, td := int(t.getYear()), int(t.getMonth()), int(t.getDay())
	daynr := calcDaynr(ty, tm, td)
	firstDaynr := calcDaynr(ty, 1, 1)
	mondayFirst := wb.test(weekBehaviourMondayFirst)
	weekYear := wb.test(weekBehaviourYear)
	firstWeekday := wb.test(weekBehaviourFirstWeekday)

	weekday := calcWeekday(firstDaynr, !mondayFirst)

	year = ty

	if tm == 1 && td <= 7-weekday {
		if !weekYear &&
			((firstWeekday && weekday != 0) || (!firstWeekday && weekday >= 4)) {
			week = 0
			return
		}
		weekYear = true
		year--
		days = calcDaysInYear(year)
		firstDaynr -= days
		weekday = (weekday + 53*7 - days) % 7
	}

	if (firstWeekday && weekday != 0) ||
		(!firstWeekday && weekday >= 4) {
		days = daynr - (firstDaynr + 7 - weekday)
	} else {
		days = daynr - (firstDaynr - weekday)
	}

	if weekYear && days >= 52*7 {
		weekday = (weekday + calcDaysInYear(year)) % 7
		if (!firstWeekday && weekday < 4) ||
			(firstWeekday && weekday == 0) {
			year++
			week = 1
			return
		}
	}
	week = days/7 + 1
	return
}

// Weekday returns weekday value.
func (t CoreTime) Weekday() gotime.Weekday {
	// TODO: Consider time_zone variable.
	t1, err := t.GoTime(gotime.Local)
	// allow invalid dates
	if err != nil {
		return t1.Weekday()
	}
	return t1.Weekday()
}

// YearWeek returns year and week.
func (t CoreTime) YearWeek(mode int) (int, int) {
	behavior := weekMode(mode) | weekBehaviourYear
	return calcWeek(t, behavior)
}

// Week returns week value.
func (t CoreTime) Week(mode int) int {
	if t.getMonth() == 0 || t.getDay() == 0 {
		return 0
	}
	_, week := calcWeek(t, weekMode(mode))
	return week
}

// YearDay returns year and day.
func (t CoreTime) YearDay() int {
	if t.getMonth() == 0 || t.getDay() == 0 {
		return 0
	}
	year, month, day := t.Year(), t.Month(), t.Day()
	return calcDaynr(year, month, day) -
		calcDaynr(year, 1, 1) + 1
}

// FromDate makes a internal time representation from the given date.
func FromDate(year int, month int, day int, hour int, minute int, second int, microsecond int) CoreTime {
	v := uint64(ZeroCoreTime)
	v |= (uint64(microsecond) << microsecondBitFieldOffset) & microsecondBitFieldMask
	v |= (uint64(second) << secondBitFieldOffset) & secondBitFieldMask
	v |= (uint64(minute) << minuteBitFieldOffset) & minuteBitFieldMask
	v |= (uint64(hour) << hourBitFieldOffset) & hourBitFieldMask
	v |= (uint64(day) << dayBitFieldOffset) & dayBitFieldMask
	v |= (uint64(month) << monthBitFieldOffset) & monthBitFieldMask
	v |= (uint64(year) << yearBitFieldOffset) & yearBitFieldMask
	return CoreTime(v)
}

// DateFormat returns a textual representation of the time value formatted
// according to layout.
// See http://dev.mysql.com/doc/refman/5.7/en/date-and-time-functions.html#function_date-format
func (t Time) DateFormat(layout string) (string, error) {
	var buf bytes.Buffer
	inPatternMatch := false
	for _, b := range layout {
		if inPatternMatch {
			if err := t.convertDateFormat(b, &buf); err != nil {
				return "", err
			}
			inPatternMatch = false
			continue
		}

		// It's not in pattern match now.
		if b == '%' {
			inPatternMatch = true
		} else {
			buf.WriteRune(b)
		}
	}
	return buf.String(), nil
}

// FormatIntWidthN uses to format int with width. Insufficient digits are filled by 0.
func FormatIntWidthN(num, n int) string {
	numString := strconv.FormatInt(int64(num), 10)
	if len(numString) >= n {
		return numString
	}
	padBytes := make([]byte, n-len(numString))
	for i := range padBytes {
		padBytes[i] = '0'
	}
	return string(padBytes) + numString
}

func abbrDayOfMonth(day int) string {
	var str string
	switch day {
	case 1, 21, 31:
		str = "st"
	case 2, 22:
		str = "nd"
	case 3, 23:
		str = "rd"
	default:
		str = "th"
	}
	return str
}

func (t Time) convertDateFormat(b rune, buf *bytes.Buffer) error {
	switch b {
	case 'b':
		m := t.Month()
		if m == 0 || m > 12 {
			return errors.New("month error")
		}
		buf.WriteString(MonthNames[m-1][:3])
	case 'M':
		m := t.Month()
		if m == 0 || m > 12 {
			return errors.New("month error")
		}
		buf.WriteString(MonthNames[m-1])
	case 'm':
		buf.WriteString(FormatIntWidthN(t.Month(), 2))
	case 'c':
		buf.WriteString(strconv.FormatInt(int64(t.Month()), 10))
	case 'D':
		buf.WriteString(strconv.FormatInt(int64(t.Day()), 10))
		buf.WriteString(abbrDayOfMonth(t.Day()))
	case 'd':
		buf.WriteString(FormatIntWidthN(t.Day(), 2))
	case 'e':
		buf.WriteString(strconv.FormatInt(int64(t.Day()), 10))
	case 'j':
		fmt.Fprintf(buf, "%03d", t.YearDay())
	case 'H':
		buf.WriteString(FormatIntWidthN(t.Hour(), 2))
	case 'k':
		buf.WriteString(strconv.FormatInt(int64(t.Hour()), 10))
	case 'h', 'I':
		t := t.Hour()
		if t%12 == 0 {
			buf.WriteString("12")
		} else {
			buf.WriteString(FormatIntWidthN(t%12, 2))
		}
	case 'l':
		t := t.Hour()
		if t%12 == 0 {
			buf.WriteString("12")
		} else {
			buf.WriteString(strconv.FormatInt(int64(t%12), 10))
		}
	case 'i':
		buf.WriteString(FormatIntWidthN(t.Minute(), 2))
	case 'p':
		hour := t.Hour()
		if hour/12%2 == 0 {
			buf.WriteString("AM")
		} else {
			buf.WriteString("PM")
		}
	case 'r':
		h := t.Hour()
		h %= 24
		switch {
		case h == 0:
			fmt.Fprintf(buf, "%02d:%02d:%02d AM", 12, t.Minute(), t.Second())
		case h == 12:
			fmt.Fprintf(buf, "%02d:%02d:%02d PM", 12, t.Minute(), t.Second())
		case h < 12:
			fmt.Fprintf(buf, "%02d:%02d:%02d AM", h, t.Minute(), t.Second())
		default:
			fmt.Fprintf(buf, "%02d:%02d:%02d PM", h-12, t.Minute(), t.Second())
		}
	case 'T':
		fmt.Fprintf(buf, "%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
	case 'S', 's':
		buf.WriteString(FormatIntWidthN(t.Second(), 2))
	case 'f':
		fmt.Fprintf(buf, "%06d", t.Microsecond())
	case 'U':
		w := t.Week(0)
		buf.WriteString(FormatIntWidthN(w, 2))
	case 'u':
		w := t.Week(1)
		buf.WriteString(FormatIntWidthN(w, 2))
	case 'V':
		w := t.Week(2)
		buf.WriteString(FormatIntWidthN(w, 2))
	case 'v':
		_, w := t.YearWeek(3)
		buf.WriteString(FormatIntWidthN(w, 2))
	case 'a':
		weekday := t.Weekday()
		buf.WriteString(abbrevWeekdayName[weekday])
	case 'W':
		buf.WriteString(t.Weekday().String())
	case 'w':
		buf.WriteString(strconv.FormatInt(int64(t.Weekday()), 10))
	case 'X':
		year, _ := t.YearWeek(2)
		if year < 0 {
			buf.WriteString(strconv.FormatUint(uint64(math.MaxUint32), 10))
		} else {
			buf.WriteString(FormatIntWidthN(year, 4))
		}
	case 'x':
		year, _ := t.YearWeek(3)
		if year < 0 {
			buf.WriteString(strconv.FormatUint(uint64(math.MaxUint32), 10))
		} else {
			buf.WriteString(FormatIntWidthN(year, 4))
		}
	case 'Y':
		buf.WriteString(FormatIntWidthN(t.Year(), 4))
	case 'y':
		str := FormatIntWidthN(t.Year(), 4)
		buf.WriteString(str[2:])
	default:
		buf.WriteRune(b)
	}

	return nil
}

// Duration is the type for MySQL TIME type.
type Duration struct {
	gotime.Duration
	// Fsp is short for Fractional Seconds Precision.
	// See http://dev.mysql.com/doc/refman/5.7/en/fractional-seconds.html
	Fsp int8
}

func (d Duration) formatFrac(frac int) string {
	s := fmt.Sprintf("%06d", frac)
	return s[0:d.Fsp]
}

// String returns the time formatted using default TimeFormat and fsp.
func (d Duration) String() string {
	var buf bytes.Buffer

	sign, hours, minutes, seconds, fraction := splitDuration(d.Duration)
	if sign < 0 {
		buf.WriteByte('-')
	}

	fmt.Fprintf(&buf, "%02d:%02d:%02d", hours, minutes, seconds)
	if d.Fsp > 0 {
		buf.WriteString(".")
		buf.WriteString(d.formatFrac(fraction))
	}

	p := buf.String()

	return p
}

func splitDuration(t gotime.Duration) (int, int, int, int, int) {
	sign := 1
	if t < 0 {
		t = -t
		sign = -1
	}

	hours := t / gotime.Hour
	t -= hours * gotime.Hour
	minutes := t / gotime.Minute
	t -= minutes * gotime.Minute
	seconds := t / gotime.Second
	t -= seconds * gotime.Second
	fraction := t / gotime.Microsecond

	return sign, int(hours), int(minutes), int(seconds), int(fraction)
}

// Hour returns current hour.
// e.g, hour("11:11:11") -> 11
func (d Duration) Hour() int {
	_, hour, _, _, _ := splitDuration(d.Duration)
	return hour
}
