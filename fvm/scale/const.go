package scale

const (
	// MoveBit for the highest six bit is has special means, in order to get origin num, must move right 2bit
	// also encode a number must move left 2bit
	MoveBit = 2

	// SingleMode the higher six bit is LE for value, which is 0 to 63
	SingleMode = 0b00
	MaxSingle  = 63
	// DoubleMode the higher six bit and next byte is LE for value, which is 64 to 1<<14 - 1
	DoubleMode = 0b01
	MaxDouble  = 1<<14 - 1
	// FourByteMode the higher six bit and next three bytes is LE for value, which is 1<<14 -1 to 1<<30-1
	FourByteMode = 0b10
	MaxFour      = 1<<30 - 1
	// BigIntMode the highest six bit represent the bytes num in the next, which less than 4.
	// then, the next bytes contain the LE for value. the highest one can not be 0, which is
	// 1<<30-1 to 1<<536-1
	BigIntMode = 0b11
)

type PrimitiveType int

const (
	None PrimitiveType = iota
	String
	Struct
	Uint8
	Int8
	CompactUInt8
	Uint16
	Int16
	CompactUint16
	Uint32
	Int32
	CompactUint32
	Uint64
	Int64
	CompactUint64
	BigInt
	BigUint
	CompactBigInt
	Bool
	Vec
	Primitive
	Array
	Tuple
	Enum
)

type TypeString string

const (
	NoneName          TypeString = TypeString("None")
	StringName        TypeString = "String"
	StructName        TypeString = "struct"
	Uint8Name         TypeString = "u8"
	Int8Name          TypeString = "i8"
	CompactUInt8Name  TypeString = "Compact < u8 >"
	Uint16Name        TypeString = "u16"
	Int16Name         TypeString = "i16"
	CompactUint16Name TypeString = "Compact < u16 >"
	Uint32Name        TypeString = "u32"
	Int32Name         TypeString = "i32"
	CompactUint32Name TypeString = "Compact < u32 >"
	Uint64Name        TypeString = "u64"
	Int64Name         TypeString = "i64"
	CompactUint64Name TypeString = "Compact < u64 >"
	BigIntName        TypeString = "i128"
	BigUIntName       TypeString = "u128"
	CompactBigIntName TypeString = "Compact < u128 >"
	BoolName          TypeString = "bool"
	VecName           TypeString = "Vec"
	PrimitiveName     TypeString = "primitive"
	ArrayName         TypeString = "Array"
	TupleName         TypeString = "tuple"
	EnumName          TypeString = "enum"
)

func (t TypeString) String() string {
	return string(t)
}

// CustomParamsSection custom section named `params` in wasm defined by us.
const CustomParamsSection = "params"
