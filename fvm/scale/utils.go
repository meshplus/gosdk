package scale

import (
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func reverse(input []byte) (output []byte) {
	for i := len(input) - 1; i >= 0; i-- {
		output = append(output, input[i])
	}
	return
}

func negchange(input []byte) []byte {
	for i, t := range input {
		input[i] = math.MaxUint8 - t
	}
	incr := 1
	for i := len(input) - 1; i >= 0; i-- {
		temp := int(input[i]) + incr
		if temp > math.MaxUint8 {
			input[i] = 0
		} else {
			input[i] = byte(temp)
			break
		}
	}
	return input
}

func convertToString(param interface{}) (*CompactString, error) {
	if val, ok := param.(string); ok {
		return &CompactString{Val: val}, nil
	} else {
		return nil, errors.New("param not string")
	}
}

func convertToUint8(param interface{}) (*FixU8, error) {
	if val, ok := param.(uint8); ok {
		return &FixU8{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseUint(val2, 10, 8)
		if err != nil {
			return nil, err
		}
		return &FixU8{Val: uint8(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= 0 && val3 <= math.MaxUint8 {
			return &FixU8{Val: uint8(val3)}, nil
		} else {
			return nil, errors.New("out of uint8 range")
		}
	} else {
		return nil, errors.New("param not uint8")
	}
}

func convertToInt8(param interface{}) (*FixI8, error) {
	if val, ok := param.(int8); ok {
		return &FixI8{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseInt(val2, 10, 8)
		if err != nil {
			return nil, err
		}
		return &FixI8{Val: int8(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= math.MinInt8 && val3 <= math.MaxInt8 {
			return &FixI8{Val: int8(val3)}, nil
		} else {
			return nil, errors.New("out of int8 range")
		}
	} else {
		return nil, errors.New("param not int8")
	}
}

func convertToCompactU8(param interface{}) (*CompactU8, error) {
	if val, ok := param.(uint8); ok {
		return &CompactU8{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseUint(val2, 10, 8)
		if err != nil {
			return nil, err
		}
		return &CompactU8{Val: uint8(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= 0 && val3 <= math.MaxUint8 {
			return &CompactU8{Val: uint8(val3)}, nil
		} else {
			return nil, errors.New("out of uint8 range")
		}
	} else {
		return nil, errors.New("param not compact uint8")
	}
}

func convertToUint16(param interface{}) (*FixU16, error) {
	if val, ok := param.(uint16); ok {
		return &FixU16{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseUint(val2, 10, 16)
		if err != nil {
			return nil, err
		}
		return &FixU16{Val: uint16(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= 0 && val3 <= math.MaxUint16 {
			return &FixU16{Val: uint16(val3)}, nil
		} else {
			return nil, errors.New("out of uint16 range")
		}
	} else {
		return nil, errors.New("param not uint16")
	}
}

func convertToInt16(param interface{}) (*FixI16, error) {
	if val, ok := param.(int16); ok {
		return &FixI16{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseInt(val2, 10, 16)
		if err != nil {
			return nil, err
		}
		return &FixI16{Val: int16(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= math.MinInt16 && val3 <= math.MaxInt16 {
			return &FixI16{Val: int16(val3)}, nil
		} else {
			return nil, errors.New("out of int16 range")
		}
	} else {
		return nil, errors.New("param not int16")
	}
}

func convertToCompactU16(param interface{}) (*CompactU16, error) {
	if val, ok := param.(uint16); ok {
		return &CompactU16{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseUint(val2, 10, 16)
		if err != nil {
			return nil, err
		}
		return &CompactU16{Val: uint16(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= 0 && val3 <= math.MaxUint16 {
			return &CompactU16{Val: uint16(val3)}, nil
		} else {
			return nil, errors.New("out of uint16 range")
		}
	} else {
		return nil, errors.New("param not compact uint16")
	}
}

func convertToUint32(param interface{}) (*FixU32, error) {
	if val, ok := param.(uint32); ok {
		return &FixU32{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseUint(val2, 10, 32)
		if err != nil {
			return nil, err
		}
		return &FixU32{Val: uint32(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= 0 && val3 <= math.MaxUint32 {
			return &FixU32{Val: uint32(val3)}, nil
		} else {
			return nil, errors.New("out of uint32 range")
		}
	} else {
		return nil, errors.New("param not uint32")
	}
}

func convertToInt32(param interface{}) (*FixI32, error) {
	if val, ok := param.(int32); ok {
		return &FixI32{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseInt(val2, 10, 32)
		if err != nil {
			return nil, err
		}
		return &FixI32{Val: int32(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= math.MinInt32 && val3 <= math.MaxInt32 {
			return &FixI32{Val: int32(val3)}, nil
		} else {
			return nil, errors.New("out of int32 range")
		}
	} else {
		return nil, errors.New("param not int32")
	}
}

func convertToCompactU32(param interface{}) (*CompactU32, error) {
	if val, ok := param.(uint32); ok {
		return &CompactU32{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseUint(val2, 10, 32)
		if err != nil {
			return nil, err
		}
		return &CompactU32{Val: uint32(s)}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= 0 && val3 <= math.MaxUint32 {
			return &CompactU32{Val: uint32(val3)}, nil
		} else {
			return nil, errors.New("out of uint32 range")
		}
	} else {
		return nil, errors.New("param not compact uint32")
	}
}

func convertToUint64(param interface{}) (*FixU64, error) {
	if val, ok := param.(uint64); ok {
		return &FixU64{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseUint(val2, 10, 64)
		if err != nil {
			return nil, err
		}
		return &FixU64{Val: s}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= 0 && uint64(val3) <= math.MaxUint64 {
			return &FixU64{Val: uint64(val3)}, nil
		} else {
			return nil, errors.New("out of uint64 range")
		}
	} else {
		return nil, errors.New("param not uint64")
	}
}

func convertToInt64(param interface{}) (*FixI64, error) {
	if val, ok := param.(int64); ok {
		return &FixI64{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseInt(val2, 10, 64)
		if err != nil {
			return nil, err
		}
		return &FixI64{Val: s}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= math.MinInt64 && val3 <= math.MaxInt64 {
			return &FixI64{Val: int64(val3)}, nil
		} else {
			return nil, errors.New("out of int64 range")
		}
	} else {
		return nil, errors.New("param not int64")
	}
}

func convertToCompactU64(param interface{}) (*CompactU64, error) {
	if val, ok := param.(uint64); ok {
		return &CompactU64{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		s, err := strconv.ParseUint(val2, 10, 64)
		if err != nil {
			return nil, err
		}
		return &CompactU64{Val: s}, nil
	} else if val3, ok3 := param.(int); ok3 {
		if val3 >= 0 && uint64(val3) <= math.MaxUint64 {
			return &CompactU64{Val: uint64(val3)}, nil
		} else {
			return nil, errors.New("out of uint64 range")
		}
	} else {
		return nil, errors.New("param not compact uint64")
	}
}

func convertToInt128(param interface{}) (*FixI128, error) {
	if val, ok := param.(*big.Int); ok {
		return &FixI128{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		t := new(big.Int)
		_, success := t.SetString(val2, 10)
		if !success {
			return nil, errors.New("invalid param")
		}
		return &FixI128{Val: t}, nil
	} else if val3, ok3 := param.(int); ok3 {
		t := new(big.Int)
		t.SetInt64(int64(val3))
		return &FixI128{Val: t}, nil
	} else {
		return nil, errors.New("param not i128")
	}
}

func convertToUint128(param interface{}) (*FixU128, error) {
	if val, ok := param.(*big.Int); ok {
		return &FixU128{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		t := new(big.Int)
		_, success := t.SetString(val2, 10)
		if !success {
			return nil, errors.New("invalid param")
		}
		return &FixU128{Val: t}, nil
	} else if val3, ok3 := param.(int); ok3 {
		t := new(big.Int)
		t.SetUint64(uint64(val3))
		return &FixU128{Val: t}, nil
	} else {
		return nil, errors.New("param not u128")
	}
}

func convertToCompactU128(param interface{}) (*CompactU128, error) {
	if val, ok := param.(*big.Int); ok {
		return &CompactU128{Val: val}, nil
	} else if val2, ok2 := param.(string); ok2 {
		t := new(big.Int)
		_, success := t.SetString(val2, 10)
		if !success {
			return nil, errors.New("invalid param")
		}
		return &CompactU128{Val: t}, nil
	} else if val3, ok3 := param.(int); ok3 {
		t := new(big.Int)
		t.SetUint64(uint64(val3))
		return &CompactU128{Val: t}, nil
	} else {
		return nil, errors.New("param not compact u128")
	}
}

func convertToBool(param interface{}) (*CompactBool, error) {
	if val, ok := param.(bool); ok {
		return &CompactBool{Val: val}, nil
	} else {
		return nil, errors.New("param not bool")
	}
}

func convertPrimitive(cur string, param interface{}) (Compact, error) {
	switch changeStringToType(TypeString(cur)) {
	case String:
		return convertToString(param)
	case Uint8:
		return convertToUint8(param)
	case Int8:
		return convertToInt8(param)
	case CompactUInt8:
		return convertToCompactU8(param)
	case Uint16:
		return convertToUint16(param)
	case Int16:
		return convertToInt16(param)
	case CompactUint16:
		return convertToCompactU16(param)
	case Uint32:
		return convertToUint32(param)
	case Int32:
		return convertToInt32(param)
	case CompactUint32:
		return convertToCompactU32(param)
	case Uint64:
		return convertToUint64(param)
	case Int64:
		return convertToInt64(param)
	case CompactUint64:
		return convertToCompactU64(param)
	case BigInt:
		return convertToInt128(param)
	case BigUint:
		return convertToUint128(param)
	case CompactBigInt:
		return convertToCompactU128(param)
	case Bool:
		return convertToBool(param)
	default:
		return nil, errors.New("unsupported type")
	}
}

func getCompactIntEncodeMode(a int) int {
	if a >= 0 && a <= MaxSingle {
		return SingleMode
	} else if a <= MaxDouble {
		return DoubleMode
	} else if a <= MaxFour {
		return FourByteMode
	} else {
		return BigIntMode
	}
}

func formatTypeString(tyName TypeString) TypeString {
	switch strings.ToLower(tyName.String()) {
	case "string", "str":
		return StringName
	case "struct":
		return StructName
	case "u8":
		return Uint8Name
	case "i8":
		return Int8Name
	case "compact < u8 >":
		return CompactUInt8Name
	case "u16":
		return Uint16Name
	case "i16":
		return Int16Name
	case "compact < u16 >":
		return CompactUint16Name
	case "u32":
		return Uint32Name
	case "i32":
		return Int32Name
	case "compact < u32 >":
		return CompactUint32Name
	case "u64":
		return Uint64Name
	case "i64":
		return Int64Name
	case "compact < u64 >":
		return CompactUint64Name
	case "u128":
		return BigUIntName
	case "i128":
		return BigIntName
	case "compact < u128 >":
		return CompactBigIntName
	case "bool":
		return BoolName
	case "vec":
		return VecName
	case "primitive":
		return PrimitiveName
	case "array":
		return ArrayName
	case "enum":
		return EnumName
	case "tuple":
		return TupleName
	default:
		return NoneName
	}
}

func changeStringToType(tyName TypeString) PrimitiveType {
	switch strings.ToLower(tyName.String()) {
	case "string", "str":
		return String
	case "struct":
		return Struct
	case "u8":
		return Uint8
	case "i8":
		return Int8
	case "compact < u8 >":
		return CompactUInt8
	case "u16":
		return Uint16
	case "i16":
		return Int16
	case "compact < u16 >":
		return CompactUint16
	case "u32":
		return Uint32
	case "i32":
		return Int32
	case "compact < u32 >":
		return CompactUint32
	case "u64":
		return Uint64
	case "i64":
		return Int64
	case "compact < u64 >":
		return CompactUint64
	case "u128":
		return BigUint
	case "i128":
		return BigInt
	case "compact < u128 >":
		return CompactBigInt
	case "bool":
		return Bool
	case "vec":
		return Vec
	case "primitive":
		return Primitive
	case "array":
		return Array
	case "enum":
		return Enum
	case "tuple":
		return Tuple
	default:
		return None
	}
}

func GetCompactValue(val Compact) interface{} {
	switch val.GetType() {
	case ArrayName, VecName, StructName, TupleName:
		var values []interface{}
		for _, v1 := range val.GetVal().([]Compact) {
			values = append(values, GetCompactValue(v1))
		}
		return values
	default:
		return val.GetVal()
	}
}
