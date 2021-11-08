package errortype

import (
	"errors"
)

var (
	ErrOverflow  = errors.New("value is out of range")
	ErrTruncated = errors.New("data truncated")
	ErrBadNumber = errors.New("bad number")
	ErrDivByZero = errors.New("div by zero")
)
