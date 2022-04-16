package scale

import (
	"bytes"
	"errors"
)

type CompactVec struct {
	Val      []Compact
	NextList []TypeString
}

func (c *CompactVec) Encode() ([]byte, error) {
	if len(c.NextList) == 0 {
		return nil, errors.New("must point next type list for vec")
	}
	l := getStringLenCompactType(len(c.Val))
	if l == nil {
		return nil, errors.New("vector length is too long")
	}
	var buf bytes.Buffer
	lengthEncode, err := l.Encode()
	if err != nil {
		return nil, err
	}
	buf.Write(lengthEncode)
	for _, v := range c.Val {
		if v.GetType() != c.NextList[0] {
			return nil, errors.New("invalid vec")
		}
		res, err := v.Encode()
		if err != nil {
			return nil, err
		}
		buf.Write(res)
	}
	return buf.Bytes(), nil
}

func (c *CompactVec) Decode(value []byte) (int, error) {
	if len(c.NextList) == 0 {
		return 0, errors.New("must point next type list for vec")
	}
	ss := &CompactU128{}
	offset, err := ss.Decode(value)
	if err != nil {
		return 0, err
	}
	if c.Val == nil {
		c.Val = make([]Compact, ss.Val.Uint64())
	}
	value = value[offset:]
	if ss.Val.Uint64() == 0 {
		c.Clear()
	}
	for i := 0; i < int(ss.Val.Uint64()); i++ {
		temp, err := c.getNextCompact()
		if err != nil {
			return 0, err
		}
		tempOffset, err := temp.Decode(value)
		if err != nil {
			return 0, err
		}
		if len(c.Val)-1 < i {
			c.Val = append(c.Val, temp)
		} else {
			c.Val[i] = temp
		}
		offset += tempOffset
		value = value[tempOffset:]
	}
	return offset, nil
}

func (c *CompactVec) GetVal() interface{} {
	return c.Val
}

func (c *CompactVec) GetType() TypeString {
	return VecName
}

func (c *CompactVec) getNextCompact() (Compact, error) {
	switch changeStringToType(c.NextList[0]) {
	case String:
		return &CompactString{Type: StringName}, nil
	case Uint8:
		return &FixU8{}, nil
	case Int8:
		return &FixI8{}, nil
	case CompactUInt8:
		return &CompactU8{}, nil
	case Uint16:
		return &FixU16{}, nil
	case Int16:
		return &FixI16{}, nil
	case CompactUint16:
		return &CompactU16{}, nil
	case Uint32:
		return &FixU32{}, nil
	case Int32:
		return &FixI32{}, nil
	case CompactUint32:
		return &CompactU32{}, nil
	case Uint64:
		return &FixU64{}, nil
	case Int64:
		return &FixI64{}, nil
	case CompactUint64:
		return &CompactU64{}, nil
	case BigInt:
		return &FixI128{}, nil
	case BigUint:
		return &FixU128{}, nil
	case CompactBigInt:
		return &CompactU128{}, nil
	case Bool:
		return &CompactBool{}, nil
	case Vec:
		return c.Val[0].(*CompactVec).CloneNew(), nil
	case Struct:
		return c.Val[0].(*CompactStruct).CloneNew(), nil
	case Array:
		return c.Val[0].(*CompactArray).CloneNew(), nil
	default:
		return nil, errors.New("not supported type")
	}
}

func (c *CompactVec) CloneNew() Compact {
	temp := &CompactVec{
		Val:      nil,
		NextList: nil,
	}
	for _, v := range c.NextList {
		temp.NextList = append(temp.NextList, v)
	}
	for _, v := range c.Val {
		temp.Val = append(temp.Val, v.CloneNew())
	}
	return temp
}

func (c *CompactVec) Clear() {
	c.Val = nil
}
