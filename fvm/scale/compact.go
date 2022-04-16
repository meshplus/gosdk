package scale

import (
	"bytes"
	"errors"
	"math"
	"math/big"
)

// Compact base interface for encoding and decoding. We redefine
// the way to encode and decode all the types in abi file.
// For more details see [SCALE Codec](https://docs.substrate.io/v3/advanced/scale-codec/)
type Compact interface {
	Encode() ([]byte, error)
	Decode(val []byte) (int, error)
	GetVal() interface{}
	GetType() TypeString
	CloneNew() Compact
}

// CompactU8 one byte uint in compact mode
type CompactU8 struct {
	Val uint8
}

func (c *CompactU8) Encode() ([]byte, error) {
	if c.Val > 255 {
		return nil, U8OutOfRange
	}
	if c.Val > 63 && int(c.Val) <= MaxDouble {
		temp := uint16(c.Val)
		temp = temp<<MoveBit + DoubleMode
		return []byte{byte(temp), byte(temp >> 8)}, nil
	}
	return []byte{c.Val<<MoveBit + SingleMode}, nil
}

func (c *CompactU8) Decode(val []byte) (int, error) {
	switch val[0] % 4 {
	case SingleMode:
		c.Val = val[0] >> MoveBit
		return 1, nil
	case DoubleMode:
		temp := uint16(val[0]>>MoveBit) | (uint16(val[1])<<8)>>MoveBit
		if temp > 63 && temp <= 255 {
			c.Val = uint8(temp)
		} else {
			return 0, U8OutOfRange
		}
		return 2, nil
	default:
		return 0, errors.New("unexpected prefix decoding Compact<u8>")
	}
}

func (c *CompactU8) GetVal() interface{} {
	return c.Val
}

func (c *CompactU8) GetType() TypeString {
	return CompactUInt8Name
}

func (c *CompactU8) CloneNew() Compact {
	return &CompactU8{Val: c.Val}
}

// CompactU16 two bytes uint in compact mode
type CompactU16 struct {
	Val uint16
}

func (u *CompactU16) GetVal() interface{} {
	return u.Val
}

func (u *CompactU16) Encode() ([]byte, error) {
	if u.Val > math.MaxUint16 {
		return nil, U16OutOfRange
	}
	if u.Val <= MaxSingle {
		s := uint8(u.Val)
		return []byte{s<<MoveBit + SingleMode}, nil
	} else if u.Val <= MaxDouble {
		u.Val = u.Val<<MoveBit + DoubleMode
		return []byte{byte(u.Val), byte(u.Val >> 8)}, nil
	} else {
		s := uint32(u.Val)
		s = s<<MoveBit + FourByteMode
		return []byte{byte(s), byte(s >> 8), byte(s >> 16), byte(s >> 24)}, nil
	}
}

func (u *CompactU16) Decode(val []byte) (int, error) {
	switch val[0] % 4 {
	case SingleMode:
		u.Val = uint16(val[0]) >> MoveBit
		return 1, nil
	case DoubleMode:
		u.Val = uint16(val[0]>>MoveBit) | (uint16(val[1])<<8)>>MoveBit
		return 2, nil
	case FourByteMode:
		temp := uint32(val[0]>>MoveBit) | (uint32(val[1])<<8|uint32(val[2])<<16|(uint32(val[3])<<24))>>MoveBit
		if temp > MaxDouble && temp <= math.MaxUint16 {
			u.Val = uint16(temp)
			return 4, nil
		} else {
			return 0, U16OutOfRange
		}
	default:
		return 0, errors.New("unexpected prefix decoding Compact<u16>")
	}
}

func (u *CompactU16) GetType() TypeString {
	return CompactUint16Name
}

func (u *CompactU16) CloneNew() Compact {
	return &CompactU16{Val: u.Val}
}

// CompactU32 four bytes uint
type CompactU32 struct {
	Val uint32
}

func (u *CompactU32) GetVal() interface{} {
	return u.Val
}

func (u *CompactU32) Encode() ([]byte, error) {
	if u.Val > math.MaxUint32 {
		return nil, U32OutOfRange
	}
	if u.Val <= MaxSingle {
		s := uint8(u.Val)
		return []byte{s<<MoveBit + SingleMode}, nil
	} else if u.Val <= MaxDouble {
		s := uint16(u.Val)
		s = s<<MoveBit + DoubleMode
		return []byte{byte(s), byte(s >> 8)}, nil
	} else if u.Val <= MaxFour {
		u.Val = u.Val<<MoveBit + FourByteMode
		return []byte{byte(u.Val), byte(u.Val >> 8), byte(u.Val >> 16), byte(u.Val >> 24)}, nil
	} else {
		return []byte{BigIntMode, byte(u.Val), byte(u.Val >> 8), byte(u.Val >> 16), byte(u.Val >> 24)}, nil
	}
}

func (u *CompactU32) Decode(val []byte) (int, error) {
	switch val[0] % 4 {
	case SingleMode:
		u.Val = uint32(val[0] >> MoveBit)
		return 1, nil
	case DoubleMode:
		temp := uint16(val[0]>>MoveBit) | (uint16(val[1])<<8)>>MoveBit
		if temp > 63 && temp < MaxDouble {
			u.Val = uint32(temp)
			return 2, nil
		} else {
			return 0, U32OutOfRange
		}
	case FourByteMode:
		u.Val = uint32(val[0]>>MoveBit) | (uint32(val[1])<<8|uint32(val[2])<<16|uint32(val[3])<<24)>>MoveBit
		return 4, nil
	case BigIntMode:
		if val[0]>>MoveBit == 0 {
			temp := uint32(val[1]) | uint32(val[2])<<8 | uint32(val[3])<<16 | uint32(val[4])<<24
			if temp > math.MaxUint32>>2 {
				u.Val = temp
				return 5, nil
			} else {
				return 0, U32OutOfRange
			}
		}
	}
	return 0, nil
}

func (u *CompactU32) GetType() TypeString {
	return CompactUint32Name
}

func (u *CompactU32) CloneNew() Compact {
	return &CompactU32{Val: u.Val}
}

type CompactU64 struct {
	Val uint64
}

func (c *CompactU64) GetVal() interface{} {
	return c.Val
}

func (c *CompactU64) Encode() ([]byte, error) {
	if c.Val > uint64(math.MaxUint64) {
		return nil, U64OutOfRange
	}
	if c.Val <= MaxSingle {
		s := uint8(c.Val)
		return []byte{s<<MoveBit + SingleMode}, nil
	} else if c.Val <= MaxDouble {
		s := uint16(c.Val)
		s = s<<MoveBit + DoubleMode
		return []byte{byte(s), byte(s >> 8)}, nil
	} else if c.Val <= MaxFour {
		s := uint32(c.Val)
		s = s<<MoveBit + FourByteMode
		return []byte{byte(s), byte(s >> 8), byte(s >> 16), byte(s >> 24)}, nil
	} else {
		res := BaseIntEncode.EncodeUInt64(c.Val)
		var buf bytes.Buffer
		end := len(res)
		for i := len(res) - 1; i >= 0; i-- {
			if res[i] != 0 {
				break
			}
			end--
		}
		buf.WriteByte(byte((end-4)<<MoveBit + BigIntMode))
		buf.Write(res[:end])
		return buf.Bytes(), nil
	}
}

func (c *CompactU64) Decode(val []byte) (int, error) {
	switch val[0] % 4 {
	case SingleMode:
		c.Val = uint64(val[0] >> MoveBit)
		return 1, nil
	case DoubleMode:
		temp := uint16(val[0]>>MoveBit) | (uint16(val[1])<<8)>>MoveBit
		if temp > 63 && temp < MaxDouble {
			c.Val = uint64(temp)
			return 2, nil
		} else {
			return 0, U32OutOfRange
		}
	case FourByteMode:
		temp := uint32(val[0]>>MoveBit) | (uint32(val[1])<<8|uint32(val[2])<<16|uint32(val[3])<<24)>>MoveBit
		if temp > MaxDouble && temp < math.MaxUint32>>2 {
			c.Val = uint64(temp)
			return 4, nil
		} else {
			return 0, U64OutOfRange
		}
	case BigIntMode:
		ll := int(val[0]>>MoveBit) + 4
		if ll == 4 {
			temp := BaseIntEncode.UInt32(val[1:])
			if temp > math.MaxUint32>>2 {
				c.Val = uint64(temp)
			} else {
				return 0, U64OutOfRange
			}
		} else if ll == 8 {
			temp := BaseIntEncode.UInt64(val[1:])
			if temp > math.MaxUint64>>2 {
				c.Val = temp
			} else {
				return 0, U64OutOfRange
			}
		} else if ll > 8 {
			return 0, errors.New("unexpected prefix decoding Compact<u64>")
		} else {
			var buf bytes.Buffer
			buf.Write(val[1:])
			if ll < 8 {
				before := make([]byte, 8-ll)
				buf.Write(before)
			}
			c.Val = BaseIntEncode.UInt64(buf.Bytes())
		}
		return ll + 1, nil
	}
	return 0, nil
}

func (c *CompactU64) GetType() TypeString {
	return CompactUint64Name
}

func (c *CompactU64) CloneNew() Compact {
	return &CompactU64{Val: c.Val}
}

type CompactU128 struct {
	Val *big.Int
}

func (c *CompactU128) GetVal() interface{} {
	return c.Val
}

func (c *CompactU128) GetType() TypeString {
	return CompactBigIntName
}

func (c *CompactU128) Encode() ([]byte, error) {
	temp := c.Val.Uint64()
	if temp <= MaxSingle {
		s := uint8(temp)
		return []byte{s<<MoveBit + SingleMode}, nil
	} else if temp <= MaxDouble {
		s := uint16(temp)
		s = s<<MoveBit + DoubleMode
		return []byte{byte(s), byte(s >> 8)}, nil
	} else if temp <= MaxFour {
		s := uint32(temp)
		s = s<<MoveBit + FourByteMode
		return []byte{byte(s), byte(s >> 8), byte(s >> 16), byte(s >> 24)}, nil
	} else {
		res := BaseIntEncode.EncodeU128(c.Val)
		var buf bytes.Buffer
		end := len(res)
		for i := len(res) - 1; i >= 0; i-- {
			if res[i] != 0 {
				break
			}
			end--
		}
		buf.WriteByte(byte((end-4)<<MoveBit + BigIntMode))
		buf.Write(res[:end])
		return buf.Bytes(), nil
	}
}

func (c *CompactU128) Decode(val []byte) (int, error) {
	switch val[0] % 4 {
	case SingleMode:
		c.Val = new(big.Int).SetBytes([]byte{val[0] >> MoveBit})
		return 1, nil
	case DoubleMode:
		temp := uint16(val[0]>>MoveBit) | (uint16(val[1])<<8)>>MoveBit
		if temp > 63 && temp < MaxDouble {
			c.Val = new(big.Int).SetUint64(uint64(temp))
			return 2, nil
		} else {
			return 0, U128OutOfRange
		}
	case FourByteMode:
		temp := uint32(val[0]>>MoveBit) | (uint32(val[1])<<8|uint32(val[2])<<16|uint32(val[3])<<24)>>MoveBit
		if temp > MaxDouble && temp < math.MaxUint32>>2 {
			c.Val = new(big.Int).SetUint64(uint64(temp))
			return 4, nil
		} else {
			return 0, U128OutOfRange
		}
	case BigIntMode:
		ll := int(val[0]>>MoveBit) + 4
		if ll == 4 {
			temp := BaseIntEncode.UInt32(val[1:])
			if temp > math.MaxUint32>>2 {
				c.Val = new(big.Int).SetUint64(uint64(temp))
			} else {
				return 0, U128OutOfRange
			}
		} else if ll == 8 {
			temp := BaseIntEncode.UInt64(val[1:])
			if temp > math.MaxUint64>>2 {
				c.Val = new(big.Int).SetUint64(temp)
			} else {
				return 0, U128OutOfRange
			}
		} else if ll > 16 {
			return 0, errors.New("unexpected prefix decoding Compact<u128>")
		} else {
			var buf bytes.Buffer
			buf.Write(val[1:])
			if ll < 16 {
				before := make([]byte, 16-ll)
				buf.Write(before)
			}
			c.Val = BaseIntEncode.U128(buf.Bytes())

		}
		return ll + 1, nil
	}

	return 0, nil
}

func (c *CompactU128) CloneNew() Compact {
	return &CompactU128{Val: new(big.Int).SetBytes(c.Val.Bytes())}
}
