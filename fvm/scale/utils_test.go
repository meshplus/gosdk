package scale

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/big"
	"testing"
)

func Test_convertToString(t *testing.T) {
	a, err := convertToString("123")
	assert.Nil(t, err)
	assert.Equal(t, "123", a.GetVal())

	_, err = convertToString(1)
	assert.NotNil(t, err)
}

func Test_convertToUint8(t *testing.T) {
	a, err := convertToUint8("123")
	assert.Nil(t, err)
	assert.Equal(t, uint8(123), a.GetVal())

	_, err = convertToUint8("test")
	assert.NotNil(t, err)

	a, err = convertToUint8(uint8(128))
	assert.Nil(t, err)
	assert.Equal(t, uint8(128), a.GetVal())

	a, err = convertToUint8(111)
	assert.Nil(t, err)
	assert.Equal(t, uint8(111), a.GetVal())

	_, err = convertToUint8(256)
	assert.NotNil(t, err)

	_, err = convertToUint8(uint16(23))
	assert.NotNil(t, err)
}

func Test_convertToInt8(t *testing.T) {
	a, err := convertToInt8("123")
	assert.Nil(t, err)
	assert.Equal(t, int8(123), a.GetVal())

	_, err = convertToInt8("test")
	assert.NotNil(t, err)

	a, err = convertToInt8(int8(123))
	assert.Nil(t, err)
	assert.Equal(t, int8(123), a.GetVal())

	_, err = convertToInt8(128)
	assert.NotNil(t, err)

	_, err = convertToInt8(256)
	assert.NotNil(t, err)

	_, err = convertToInt8(uint16(23))
	assert.NotNil(t, err)
}

func Test_convertToCompactU8(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToCompactU8("123")
		assert.Nil(t, err)
		assert.Equal(t, uint8(123), a.GetVal())
	})
	t.Run("string+error", func(t *testing.T) {
		_, err := convertToCompactU8("test")
		assert.NotNil(t, err)
	})
	t.Run("uint8", func(t *testing.T) {
		a, err := convertToCompactU8(uint8(23))
		assert.Nil(t, err)
		assert.Equal(t, uint8(23), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToCompactU8(23)
		assert.Nil(t, err)
		assert.Equal(t, uint8(23), a.GetVal())
	})
	t.Run("int+error", func(t *testing.T) {
		_, err := convertToCompactU8(-1)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToCompactU8(uint16(1))
		assert.NotNil(t, err)
	})
}

func Test_convertToUint16(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToUint16("23456")
		assert.Nil(t, err)
		assert.Equal(t, uint16(23456), a.GetVal())
	})
	t.Run("string+error", func(t *testing.T) {
		_, err := convertToUint16("test")
		assert.NotNil(t, err)
	})
	t.Run("uint16", func(t *testing.T) {
		a, err := convertToUint16(uint16(25))
		assert.Nil(t, err)
		assert.Equal(t, uint16(25), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToUint16(24)
		assert.Nil(t, err)
		assert.Equal(t, uint16(24), a.GetVal())
	})
	t.Run("int+err", func(t *testing.T) {
		_, err := convertToUint16(-1)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToUint16(uint8(1))
		assert.NotNil(t, err)
	})
}

func Test_convertToInt16(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToInt16("123")
		assert.Nil(t, err)
		assert.Equal(t, int16(123), a.GetVal())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToInt16("test")
		assert.NotNil(t, err)
	})
	t.Run("int16", func(t *testing.T) {
		a, err := convertToInt16(int16(30))
		assert.Nil(t, err)
		assert.Equal(t, int16(30), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToInt16(30)
		assert.Nil(t, err)
		assert.Equal(t, int16(30), a.GetVal())
	})
	t.Run("int+err", func(t *testing.T) {
		_, err := convertToInt16(math.MaxInt16 + 20)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToInt16(uint8(20))
		assert.NotNil(t, err)
	})
}

func Test_convertToCompactU16(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToCompactU16("123")
		assert.Nil(t, err)
		assert.Equal(t, uint16(123), a.GetVal())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToCompactU16("test")
		assert.NotNil(t, err)
	})
	t.Run("uint16", func(t *testing.T) {
		a, err := convertToCompactU16(uint16(30))
		assert.Nil(t, err)
		assert.Equal(t, uint16(30), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToCompactU16(30)
		assert.Nil(t, err)
		assert.Equal(t, uint16(30), a.GetVal())
	})
	t.Run("int+err", func(t *testing.T) {
		_, err := convertToCompactU16(math.MaxUint16 + 2)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToCompactU16(uint32(1))
		assert.NotNil(t, err)
	})
}

func Test_convertToUint32(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToUint32("12345")
		assert.Nil(t, err)
		assert.Equal(t, uint32(12345), a.GetVal())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToUint32("test")
		assert.NotNil(t, err)
	})
	t.Run("uint32", func(t *testing.T) {
		a, err := convertToUint32(uint32(45678))
		assert.Nil(t, err)
		assert.Equal(t, uint32(45678), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToUint32(1234)
		assert.Nil(t, err)
		assert.Equal(t, uint32(1234), a.GetVal())
	})
	t.Run("int+err", func(t *testing.T) {
		_, err := convertToUint32(math.MaxInt64)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToUint32(uint8(1))
		assert.NotNil(t, err)
	})
}

func Test_convertToInt32(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToInt32("123")
		assert.Nil(t, err)
		assert.Equal(t, int32(123), a.GetVal())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToInt32("test")
		assert.NotNil(t, err)
	})
	t.Run("int32", func(t *testing.T) {
		a, err := convertToInt32(int32(123))
		assert.Nil(t, err)
		assert.Equal(t, int32(123), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToInt32(123)
		assert.Nil(t, err)
		assert.Equal(t, int32(123), a.GetVal())
	})
	t.Run("int+err", func(t *testing.T) {
		_, err := convertToInt32(math.MaxInt32 + 3)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToInt32(uint16(1))
		assert.NotNil(t, err)
	})
}

func Test_convertToCompactU32(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToCompactU32("123")
		assert.Nil(t, err)
		assert.Equal(t, uint32(123), a.GetVal())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToCompactU32("test")
		assert.NotNil(t, err)
	})
	t.Run("uint32", func(t *testing.T) {
		a, err := convertToCompactU32(uint32(123))
		assert.Nil(t, err)
		assert.Equal(t, uint32(123), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToCompactU32(123)
		assert.Nil(t, err)
		assert.Equal(t, uint32(123), a.GetVal())
	})
	t.Run("int+err", func(t *testing.T) {
		_, err := convertToCompactU32(math.MaxUint32 + 2)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToCompactU32(int16(32))
		assert.NotNil(t, err)
	})
}

func Test_convertToUint64(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToUint64("1234")
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToUint64("test")
		assert.NotNil(t, err)
	})
	t.Run("uint64", func(t *testing.T) {
		a, err := convertToUint64(uint64(123))
		assert.Nil(t, err)
		assert.Equal(t, uint64(123), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToUint64(123)
		assert.Nil(t, err)
		assert.Equal(t, uint64(123), a.GetVal())
	})
	t.Run("int+err", func(t *testing.T) {
		_, err := convertToUint64(-1)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToUint64(int16(-1))
		assert.NotNil(t, err)
	})
}

func Test_convertToInt64(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToInt64("1234")
		assert.Nil(t, err)
		assert.Equal(t, int64(1234), a.GetVal())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToInt64("test")
		assert.NotNil(t, err)
	})
	t.Run("int64", func(t *testing.T) {
		a, err := convertToInt64(int64(-1))
		assert.Nil(t, err)
		assert.Equal(t, int64(-1), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToInt64(123)
		assert.Nil(t, err)
		assert.Equal(t, int64(123), a.GetVal())
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToInt64(int16(-1))
		assert.NotNil(t, err)
	})
}

func Test_convertToCompactU64(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToCompactU64("1234")
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToCompactU64("test")
		assert.NotNil(t, err)
	})
	t.Run("uint64", func(t *testing.T) {
		a, err := convertToCompactU64(uint64(1))
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), a.GetVal())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToCompactU64(123)
		assert.Nil(t, err)
		assert.Equal(t, uint64(123), a.GetVal())
	})
	t.Run("int+err", func(t *testing.T) {
		_, err := convertToCompactU64(-1)
		assert.NotNil(t, err)
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToCompactU64(uint16(1))
		assert.NotNil(t, err)
	})
}

func Test_convertToInt128(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToInt128("1234")
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToInt128("test")
		assert.NotNil(t, err)
	})
	t.Run("int128", func(t *testing.T) {
		a, err := convertToInt128(new(big.Int).SetInt64(1234))
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToInt128(1234)
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToInt128(int16(1234))
		assert.NotNil(t, err)
	})
}

func Test_convertToUint128(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToUint128("1234")
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToUint128("test")
		assert.NotNil(t, err)
	})
	t.Run("int128", func(t *testing.T) {
		a, err := convertToUint128(new(big.Int).SetInt64(1234))
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToUint128(1234)
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToUint128(int16(1234))
		assert.NotNil(t, err)
	})
}

func Test_convertToCompactU128(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertToCompactU128("1234")
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("string+err", func(t *testing.T) {
		_, err := convertToCompactU128("test")
		assert.NotNil(t, err)
	})
	t.Run("int128", func(t *testing.T) {
		a, err := convertToCompactU128(new(big.Int).SetInt64(1234))
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("int", func(t *testing.T) {
		a, err := convertToCompactU128(1234)
		assert.Nil(t, err)
		assert.Equal(t, uint64(1234), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToCompactU128(int16(1234))
		assert.NotNil(t, err)
	})
}

func Test_convertToBool(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		a, err := convertToBool(true)
		assert.Nil(t, err)
		assert.Equal(t, true, a.GetVal())
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertToBool("test")
		assert.NotNil(t, err)
	})
}

func Test_convertPrimitive(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		a, err := convertPrimitive("string", "123")
		assert.Nil(t, err)
		assert.Equal(t, "123", a.GetVal())
	})
	t.Run("u8", func(t *testing.T) {
		a, err := convertPrimitive("u8", "123")
		assert.Nil(t, err)
		assert.Equal(t, uint8(123), a.GetVal())
	})
	t.Run("i8", func(t *testing.T) {
		a, err := convertPrimitive("i8", "123")
		assert.Nil(t, err)
		assert.Equal(t, int8(123), a.GetVal())
	})
	t.Run("compact+u8", func(t *testing.T) {
		a, err := convertPrimitive("compact < u8 >", 123)
		assert.Nil(t, err)
		assert.Equal(t, uint8(123), a.GetVal())
	})
	t.Run("uint16", func(t *testing.T) {
		a, err := convertPrimitive("u16", 23)
		assert.Nil(t, err)
		assert.Equal(t, uint16(23), a.GetVal())
	})
	t.Run("int16", func(t *testing.T) {
		a, err := convertPrimitive("i16", 23)
		assert.Nil(t, err)
		assert.Equal(t, int16(23), a.GetVal())
	})
	t.Run("compact+u16", func(t *testing.T) {
		a, err := convertPrimitive("compact < u16 >", 123)
		assert.Nil(t, err)
		assert.Equal(t, uint16(123), a.GetVal())
	})
	t.Run("uint32", func(t *testing.T) {
		a, err := convertPrimitive("u32", 123)
		assert.Nil(t, err)
		assert.Equal(t, uint32(123), a.GetVal())
	})
	t.Run("int32", func(t *testing.T) {
		a, err := convertPrimitive("i32", 123)
		assert.Nil(t, err)
		assert.Equal(t, int32(123), a.GetVal())
	})
	t.Run("compact+u32", func(t *testing.T) {
		a, err := convertPrimitive("compact < u32 >", 123)
		assert.Nil(t, err)
		assert.Equal(t, uint32(123), a.GetVal())
	})
	t.Run("uint64", func(t *testing.T) {
		a, err := convertPrimitive("u64", 123)
		assert.Nil(t, err)
		assert.Equal(t, uint64(123), a.GetVal())
	})
	t.Run("int64", func(t *testing.T) {
		a, err := convertPrimitive("i64", 123)
		assert.Nil(t, err)
		assert.Equal(t, int64(123), a.GetVal())
	})
	t.Run("compact+u64", func(t *testing.T) {
		a, err := convertPrimitive("compact < u64 >", 123)
		assert.Nil(t, err)
		assert.Equal(t, uint64(123), a.GetVal())
	})
	t.Run("uint128", func(t *testing.T) {
		a, err := convertPrimitive("u128", 123)
		assert.Nil(t, err)
		assert.Equal(t, uint64(123), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("int128", func(t *testing.T) {
		a, err := convertPrimitive("i128", 123)
		assert.Nil(t, err)
		assert.Equal(t, int64(123), a.GetVal().(*big.Int).Int64())
	})
	t.Run("compact+u128", func(t *testing.T) {
		a, err := convertPrimitive("compact < u128 >", 123)
		assert.Nil(t, err)
		assert.Equal(t, uint64(123), a.GetVal().(*big.Int).Uint64())
	})
	t.Run("bool", func(t *testing.T) {
		a, err := convertPrimitive("bool", true)
		assert.Nil(t, err)
		assert.Equal(t, true, a.GetVal())
	})
	t.Run("err", func(t *testing.T) {
		_, err := convertPrimitive("test", true)
		assert.NotNil(t, err)
	})
}

func Test_getCompactIntEncodeMode(t *testing.T) {
	t.Run("0b00", func(t *testing.T) {
		assert.Equal(t, 0, getCompactIntEncodeMode(42))
	})
	t.Run("0b01", func(t *testing.T) {
		assert.Equal(t, 1, getCompactIntEncodeMode(69))
	})
	t.Run("0b10", func(t *testing.T) {
		assert.Equal(t, 2, getCompactIntEncodeMode(1<<14))
	})
	t.Run("0b11", func(t *testing.T) {
		assert.Equal(t, 3, getCompactIntEncodeMode(1<<30))
	})
}

func Test_formatTypeString(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, StringName, formatTypeString("str"))
	})
	t.Run("struct", func(t *testing.T) {
		assert.Equal(t, StructName, formatTypeString("struct"))
	})
	t.Run("u8", func(t *testing.T) {
		assert.Equal(t, Uint8Name, formatTypeString("u8"))
	})
	t.Run("i8", func(t *testing.T) {
		assert.Equal(t, Int8Name, formatTypeString("i8"))
	})
	t.Run("compactu8", func(t *testing.T) {
		assert.Equal(t, CompactUInt8Name, formatTypeString("compact < u8 >"))
	})
	t.Run("i16", func(t *testing.T) {
		assert.Equal(t, Int16Name, formatTypeString("i16"))
	})
	t.Run("u16", func(t *testing.T) {
		assert.Equal(t, Uint16Name, formatTypeString("u16"))
	})
	t.Run("compactu16", func(t *testing.T) {
		assert.Equal(t, CompactUint16Name, formatTypeString("compact < u16 >"))
	})
	t.Run("u32", func(t *testing.T) {
		assert.Equal(t, Uint32Name, formatTypeString("u32"))
	})
	t.Run("i32", func(t *testing.T) {
		assert.Equal(t, Int32Name, formatTypeString("i32"))
	})
	t.Run("compactU32", func(t *testing.T) {
		assert.Equal(t, CompactUint32Name, formatTypeString("compact < u32 >"))
	})
	t.Run("u64", func(t *testing.T) {
		assert.Equal(t, Uint64Name, formatTypeString("u64"))
	})
	t.Run("i64", func(t *testing.T) {
		assert.Equal(t, Int64Name, formatTypeString("i64"))
	})
	t.Run("compactU64", func(t *testing.T) {
		assert.Equal(t, CompactUint64Name, formatTypeString("compact < u64 >"))
	})
	t.Run("u128", func(t *testing.T) {
		assert.Equal(t, BigUIntName, formatTypeString("u128"))
	})
	t.Run("i128", func(t *testing.T) {
		assert.Equal(t, BigIntName, formatTypeString("i128"))
	})
	t.Run("compactU128", func(t *testing.T) {
		assert.Equal(t, CompactBigIntName, formatTypeString("compact < u128 >"))
	})
	t.Run("vec", func(t *testing.T) {
		assert.Equal(t, VecName, formatTypeString("vec"))
	})
	t.Run("array", func(t *testing.T) {
		assert.Equal(t, ArrayName, formatTypeString("array"))
	})
	t.Run("bool", func(t *testing.T) {
		assert.Equal(t, BoolName, formatTypeString("bool"))
	})
	t.Run("primi", func(t *testing.T) {
		assert.Equal(t, PrimitiveName, formatTypeString("primitive"))
	})
	t.Run("none", func(t *testing.T) {
		assert.Equal(t, NoneName, formatTypeString("test"))
	})
}
