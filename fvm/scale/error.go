package scale

import "errors"

var (
	U8OutOfRange   = errors.New("out of range Compact<u8>")
	U16OutOfRange  = errors.New("out of range Compact<u16>")
	U32OutOfRange  = errors.New("out of range Compact<u32>")
	U64OutOfRange  = errors.New("out of range Compact<u64>")
	U128OutOfRange = errors.New("out of range Compact<u128>")
)
