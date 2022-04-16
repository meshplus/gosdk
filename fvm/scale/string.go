package scale

import (
	"bytes"
	"errors"
	"math/big"
)

type CompactString struct {
	Val  string
	Type TypeString
}

func (c *CompactString) GetVal() interface{} {
	return c.Val
}

func getStringLenCompactType(length int) Compact {
	mode := getCompactIntEncodeMode(length)
	switch mode {
	case SingleMode:
		return &CompactU8{Val: uint8(length)}
	case DoubleMode:
		return &CompactU16{Val: uint16(length)}
	case FourByteMode:
		return &CompactU32{Val: uint32(length)}
	case BigIntMode:
		if length < 1<<62-1 {
			return &CompactU64{
				Val: uint64(length),
			}
		} else {
			return &CompactU128{Val: new(big.Int).SetUint64(uint64(length))}
		}
	default:
		return nil
	}
}

func (c *CompactString) Encode() ([]byte, error) {
	v := []byte(c.Val)
	l := getStringLenCompactType(len(v))
	if l == nil {
		return nil, errors.New("string's length is too long")
	}
	var buf bytes.Buffer
	lengthEncode, err := l.Encode()
	if err != nil {
		return nil, err
	}
	buf.Write(lengthEncode)
	buf.Write(v)
	return buf.Bytes(), nil
}

func (c *CompactString) Decode(value []byte) (int, error) {
	ss := &CompactU128{}
	offset, err := ss.Decode(value)
	if err != nil {
		return 0, err
	}
	end := offset + int(ss.Val.Uint64())
	c.Val = string(value[offset:end])
	return offset + int(ss.Val.Uint64()), nil
}

func (c *CompactString) GetType() TypeString {
	return StringName
}

func (c *CompactString) CloneNew() Compact {
	return &CompactString{
		Val:  c.Val,
		Type: StringName,
	}
}
