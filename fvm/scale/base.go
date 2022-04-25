package scale

import (
	"bytes"
	"math/big"
)

type littleEncode struct {
	Val  interface{}
	Type TypeString
}

var BaseIntEncode littleEncode

func (l *littleEncode) EncodeInt8(val int8) []byte {
	return []byte{byte(val)}
}

func (l *littleEncode) EncodeUInt8(val uint8) []byte {
	return []byte{val}
}

func (l *littleEncode) Int8(val []byte) int8 {
	return int8(val[0])
}

func (l *littleEncode) UInt8(val []byte) uint8 {
	return val[0]
}

func (l *littleEncode) EncodeInt16(val int16) []byte {
	return []byte{byte(val), byte(val >> 8)}
}

func (l *littleEncode) EncodeUInt16(val uint16) []byte {
	return []byte{byte(val), byte(val >> 8)}
}

func (l *littleEncode) Int16(val []byte) int16 {
	return int16(val[0]) | int16(val[1])<<8
}

func (l *littleEncode) UInt16(val []byte) uint16 {
	return uint16(val[0]) | uint16(val[1])<<8
}

func (l *littleEncode) EncodeInt32(val int32) []byte {
	return []byte{byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24)}
}

func (l *littleEncode) EncodeUInt32(val uint32) []byte {
	return []byte{byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24)}
}

func (l *littleEncode) Int32(val []byte) int32 {
	return int32(val[0]) | int32(val[1])<<8 | int32(val[2])<<16 | int32(val[3])<<24
}

func (l *littleEncode) UInt32(val []byte) uint32 {
	return uint32(val[0]) | uint32(val[1])<<8 | uint32(val[2])<<16 | uint32(val[3])<<24
}

func (l *littleEncode) EncodeInt64(val int64) []byte {
	return []byte{byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24), byte(val >> 32), byte(val >> 40), byte(val >> 48), byte(val >> 56)}
}

func (l *littleEncode) EncodeUInt64(val uint64) []byte {
	return []byte{byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24), byte(val >> 32), byte(val >> 40), byte(val >> 48), byte(val >> 56)}
}

func (l *littleEncode) Int64(val []byte) int64 {
	return int64(val[0]) | int64(val[1])<<8 | int64(val[2])<<16 | int64(val[3])<<24 | int64(val[4])<<32 | int64(val[5])<<40 | int64(val[6])<<48 | int64(val[7])<<56
}

func (l *littleEncode) UInt64(val []byte) uint64 {
	return uint64(val[0]) | uint64(val[1])<<8 | uint64(val[2])<<16 | uint64(val[3])<<24 | uint64(val[4])<<32 | uint64(val[5])<<40 | uint64(val[6])<<48 | uint64(val[7])<<56
}

func (l *littleEncode) EncodeU128(val *big.Int) []byte {
	var buf bytes.Buffer
	origin := val.Bytes()
	for i := len(origin) - 1; i >= 0; i-- {
		buf.WriteByte(origin[i])
	}
	if len(origin) < 16 {
		ap := make([]byte, 16-len(origin))
		buf.Write(ap)
	}
	return buf.Bytes()
}

func (l *littleEncode) U128(val []byte) *big.Int {
	var v big.Int
	v.SetBytes(reverse(val))
	return &v
}
