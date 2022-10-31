package scale

import (
	"errors"
)

func (a *Abi) decodeType(tp Type, val []byte) (Compact, int, error) {
	switch changeStringToType(TypeString(tp.Type)) {
	case Primitive:
		return a.decodePrimitive(tp, val)
	case Vec:
		return a.decodeVec(tp, val)
	case Struct:
		return a.decodeStruct(tp, val)
	case Array:
		return a.decodeArray(tp, val)
	case Tuple:
		return a.decodeTuple(tp, val)
	case Enum:
		return a.decodeEnum(tp, val)
	default:
		return nil, 0, errors.New("not supported")
	}
}

func (a *Abi) decodePrimitive(tp Type, val []byte) (Compact, int, error) {
	var result Compact
	curType := tp.Primitive
	switch changeStringToType(TypeString(curType)) {
	case String:
		result = &CompactString{Type: StringName}
	case Uint8:
		result = &FixU8{}
	case Int8:
		result = &FixI8{}
	case CompactUInt8:
		result = &CompactU8{}
	case Uint16:
		result = &FixU16{}
	case Int16:
		result = &FixI16{}
	case CompactUint16:
		result = &CompactU16{}
	case Uint32:
		result = &FixU32{}
	case Int32:
		result = &FixI32{}
	case CompactUint32:
		result = &CompactU32{}
	case Uint64:
		result = &FixU64{}
	case Int64:
		result = &FixI64{}
	case CompactUint64:
		result = &CompactU64{}
	case BigInt:
		result = &FixI128{}
	case BigUint:
		result = &FixU128{}
	case CompactBigInt:
		result = &CompactU128{}
	case Bool:
		result = &CompactBool{}
	default:
		return nil, 0, errors.New("not supported type")
	}
	offset, err := result.Decode(val)
	if err != nil {
		return nil, 0, err
	}
	val = val[offset:]
	return result, offset, nil
}

func (a *Abi) decodeVec(tp Type, val []byte) (Compact, int, error) {
	var offset = 0
	cc := &CompactVec{}
	ss := &CompactU128{}
	tempOffset, err := ss.Decode(val)
	if err != nil {
		return nil, 0, err
	}
	val = val[tempOffset:]
	offset += tempOffset
	for i := 0; i < int(ss.Val.Uint64()); i++ {
		res, tempOffset, err := a.decodeType(a.Types[tp.Fields[0].TypeId], val)
		if err != nil {
			return nil, 0, err
		}
		val = val[tempOffset:]
		offset += tempOffset
		cc.Val = append(cc.Val, res)
	}

	return cc, offset, nil
}

func (a *Abi) decodeStruct(tp Type, val []byte) (Compact, int, error) {
	cc := &CompactStruct{}
	var offset int = 0
	for _, k := range tp.Fields {
		res, tempOffset, err := a.decodeType(a.Types[k.TypeId], val)
		if err != nil {
			return nil, 0, err
		}
		val = val[tempOffset:]
		offset += tempOffset
		cc.Val = append(cc.Val, res)
	}

	return cc, offset, nil
}

func (a *Abi) decodeArray(tp Type, val []byte) (Compact, int, error) {
	cc := &CompactArray{
		Len: tp.ArrayLen,
	}
	var offset = 0
	for i := 0; i < tp.ArrayLen; i++ {
		res, tempOffset, err := a.decodeType(a.Types[tp.Fields[0].TypeId], val)
		if err != nil {
			return nil, 0, err
		}
		val = val[tempOffset:]
		offset += tempOffset
		cc.Val = append(cc.Val, res)
	}

	return cc, offset, nil
}

func (a *Abi) decodeTuple(tp Type, val []byte) (Compact, int, error) {
	cc := &CompactTuple{}
	var offset int = 0
	for _, k := range tp.Fields {
		res, tempOffset, err := a.decodeType(a.Types[k.TypeId], val)
		if err != nil {
			return nil, 0, err
		}
		val = val[tempOffset:]
		offset += tempOffset
		cc.Val = append(cc.Val, res)
	}

	return cc, offset, nil
}

func (a *Abi) decodeEnum(tp Type, val []byte) (Compact, int, error) {
	cc := &CompactEnum{}
	cc.index = val[0]
	val = val[1:]
	var offset int = 1
	for _, k := range tp.Variants[cc.index] {
		res, tempOffset, err := a.decodeType(a.Types[k.TypeId], val)
		if err != nil {
			return nil, 0, err
		}
		val = val[tempOffset:]
		offset += tempOffset
		cc.Val = res
	}

	return cc, offset, nil
}
