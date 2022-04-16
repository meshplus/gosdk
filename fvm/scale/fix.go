package scale

import (
	"bytes"
	"errors"
	"math"
	"math/big"
)

type FixU8 struct {
	Val uint8
}

func (f *FixU8) Encode() ([]byte, error) {
	return []byte{f.Val}, nil
}

func (f *FixU8) Decode(val []byte) (int, error) {
	f.Val = val[0]
	return 1, nil
}

func (f *FixU8) GetVal() interface{} {
	return f.Val
}

func (f *FixU8) GetType() TypeString {
	return Uint8Name
}

func (f *FixU8) CloneNew() Compact {
	return &FixU8{Val: f.Val}
}

type FixI8 struct {
	Val int8
}

func (f *FixI8) Encode() ([]byte, error) {
	return []byte{byte(f.Val)}, nil
}

func (f *FixI8) Decode(val []byte) (int, error) {
	f.Val = int8(val[0])
	return 1, nil
}

func (f *FixI8) GetVal() interface{} {
	return f.Val
}

func (f *FixI8) GetType() TypeString {
	return Int8Name
}

func (f *FixI8) CloneNew() Compact {
	return &FixI8{Val: f.Val}
}

type FixU16 struct {
	Val uint16
}

func (f *FixU16) Encode() ([]byte, error) {
	return []byte{byte(f.Val), byte(f.Val >> 8)}, nil
}

func (f *FixU16) Decode(val []byte) (int, error) {
	f.Val = uint16(val[0]) | uint16(val[1])<<8
	return 2, nil
}

func (f *FixU16) GetVal() interface{} {
	return f.Val
}

func (f *FixU16) GetType() TypeString {
	return Uint16Name
}

func (f *FixU16) CloneNew() Compact {
	return &FixU16{Val: f.Val}
}

type FixI16 struct {
	Val int16
}

func (f *FixI16) Encode() ([]byte, error) {
	return []byte{byte(f.Val), byte(f.Val >> 8)}, nil
}

func (f *FixI16) Decode(val []byte) (int, error) {
	f.Val = int16(val[0]) | int16(val[1])<<8
	return 2, nil
}

func (f *FixI16) GetVal() interface{} {
	return f.Val
}

func (f *FixI16) GetType() TypeString {
	return Int16Name
}

func (f *FixI16) CloneNew() Compact {
	return &FixI16{Val: f.Val}
}

type FixU32 struct {
	Val uint32
}

func (f *FixU32) Encode() ([]byte, error) {
	return []byte{byte(f.Val), byte(f.Val >> 8), byte(f.Val >> 16), byte(f.Val >> 24)}, nil
}

func (f *FixU32) Decode(val []byte) (int, error) {
	f.Val = uint32(val[0]) | uint32(val[1])<<8 | uint32(val[2])<<16 | uint32(val[3])<<24
	return 4, nil
}

func (f *FixU32) GetVal() interface{} {
	return f.Val
}

func (f *FixU32) GetType() TypeString {
	return Uint32Name
}

func (f *FixU32) CloneNew() Compact {
	return &FixU32{Val: f.Val}
}

type FixI32 struct {
	Val int32
}

func (f *FixI32) Encode() ([]byte, error) {
	return []byte{byte(f.Val), byte(f.Val >> 8), byte(f.Val >> 16), byte(f.Val >> 24)}, nil
}

func (f *FixI32) Decode(val []byte) (int, error) {
	f.Val = int32(val[0]) | int32(val[1])<<8 | int32(val[2])<<16 | int32(val[3])<<24
	return 4, nil
}

func (f *FixI32) GetVal() interface{} {
	return f.Val
}

func (f *FixI32) GetType() TypeString {
	return Int32Name
}

func (f *FixI32) CloneNew() Compact {
	return &FixI32{Val: f.Val}
}

type FixU64 struct {
	Val uint64
}

func (f *FixU64) Encode() ([]byte, error) {
	return []byte{
		byte(f.Val),
		byte(f.Val >> 8),
		byte(f.Val >> 16),
		byte(f.Val >> 24),
		byte(f.Val >> 32),
		byte(f.Val >> 40),
		byte(f.Val >> 48),
		byte(f.Val >> 56),
	}, nil
}

func (f *FixU64) Decode(val []byte) (int, error) {
	f.Val = uint64(val[0]) | uint64(val[1])<<8 | uint64(val[2])<<16 | uint64(val[3])<<24 | uint64(val[4])<<32 | uint64(val[5])<<40 | uint64(val[6])<<48 | uint64(val[7])<<56
	return 8, nil
}

func (f *FixU64) GetVal() interface{} {
	return f.Val
}

func (f *FixU64) GetType() TypeString {
	return Uint64Name
}

func (f *FixU64) CloneNew() Compact {
	return &FixU64{Val: f.Val}
}

type FixI64 struct {
	Val int64
}

func (f *FixI64) Encode() ([]byte, error) {
	return []byte{
		byte(f.Val),
		byte(f.Val >> 8),
		byte(f.Val >> 16),
		byte(f.Val >> 24),
		byte(f.Val >> 32),
		byte(f.Val >> 40),
		byte(f.Val >> 48),
		byte(f.Val >> 56),
	}, nil
}

func (f *FixI64) Decode(val []byte) (int, error) {
	f.Val = int64(val[0]) | int64(val[1])<<8 | int64(val[2])<<16 | int64(val[3])<<24 | int64(val[4])<<32 | int64(val[5])<<40 | int64(val[6])<<48 | int64(val[7])<<56
	return 8, nil
}

func (f *FixI64) GetVal() interface{} {
	return f.Val
}

func (f *FixI64) GetType() TypeString {
	return Int64Name
}

func (f *FixI64) CloneNew() Compact {
	return &FixI64{Val: f.Val}
}

type FixU128 struct {
	Val *big.Int
}

func (f *FixU128) Encode() ([]byte, error) {
	var buf bytes.Buffer
	origin := f.Val.Bytes()
	if len(origin) > 16 {
		return nil, errors.New("out of u128 range")
	}
	for i := len(origin) - 1; i >= 0; i-- {
		buf.WriteByte(origin[i])
	}
	if len(origin) < 16 {
		ap := make([]byte, 16-len(origin))
		buf.Write(ap)
	}
	return buf.Bytes(), nil
}

func (f *FixU128) Decode(val []byte) (int, error) {
	var v big.Int
	if len(val) != 16 {
		return 0, errors.New("out of u128 range")
	}
	v.SetBytes(reverse(val))
	f.Val = &v
	return 16, nil
}

func (f *FixU128) GetVal() interface{} {
	return f.Val
}

func (f *FixU128) GetType() TypeString {
	return BigUIntName
}

func (f *FixU128) CloneNew() Compact {
	return &FixU128{Val: new(big.Int).SetBytes(f.Val.Bytes())}
}

type FixI128 struct {
	Val *big.Int
}

func (f *FixI128) Encode() ([]byte, error) {
	var buf bytes.Buffer
	val := f.Val.Bytes()
	if len(val) > 16 {
		return nil, errors.New("out of u128 range")
	}
	sign := f.Val.Sign()
	if sign < 0 {
		negchange(val)
	}
	val = reverse(val)
	buf.Write(val)
	if len(val) < 16 {
		for i := 0; i < 16-len(val); i++ {
			if sign < 0 {
				buf.WriteByte(255)
			} else {
				buf.WriteByte(0)
			}
		}
	}
	return buf.Bytes(), nil
}

func (f *FixI128) Decode(val []byte) (int, error) {
	v := new(big.Int)
	if len(val) != 16 {
		return 0, errors.New("out of range i128")
	}
	res := reverse(val)
	if res[0] > math.MaxInt8 {
		negchange(res)
		v.SetBytes(res)
		v.Neg(v)
	} else {
		v.SetBytes(res)
	}

	f.Val = v

	return 16, nil
}

func (f *FixI128) GetVal() interface{} {
	return f.Val
}

func (f *FixI128) GetType() TypeString {
	return BigIntName
}

func (f *FixI128) CloneNew() Compact {
	return &FixI128{Val: new(big.Int).SetBytes(f.Val.Bytes())}
}
