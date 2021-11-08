package mysql

const (
	TypeTiny      byte = 1
	TypeShort     byte = 2
	TypeLong      byte = 3
	TypeFloat     byte = 4
	TypeDouble    byte = 5
	TypeTimestamp byte = 7
	TypeLonglong  byte = 8
	TypeInt24     byte = 9
	TypeDate      byte = 10
	/* TypeDuration original name was TypeTime, renamed to TypeDuration to resolve the conflict with Go type Time.*/
	TypeDuration byte = 11
	TypeDatetime byte = 12
	TypeYear     byte = 13
	TypeVarchar  byte = 15
	TypeBit      byte = 16

	TypeJSON       byte = 0xf5
	TypeNewDecimal byte = 0xf6
	TypeEnum       byte = 0xf7
	TypeSet        byte = 0xf8
	TypeTinyBlob   byte = 0xf9
	TypeMediumBlob byte = 0xfa
	TypeLongBlob   byte = 0xfb
	TypeBlob       byte = 0xfc
	TypeVarString  byte = 0xfd
	TypeString     byte = 0xfe
)

// nolint
// MySQL type maximum length.
const (
	// For arguments that have no fixed number of decimals, the decimals value is set to 31,
	// which is 1 more than the maximum number of decimals permitted for the DECIMAL, FLOAT, and DOUBLE data types.
	NotFixedDec = 31

	MaxIntWidth              = 20
	MaxRealWidth             = 23
	MaxFloatingTypeScale     = 30
	MaxFloatingTypeWidth     = 255
	MaxDecimalScale          = 30
	MaxDecimalWidth          = 65
	MaxDateWidth             = 10 // YYYY-MM-DD.
	MaxDatetimeWidthNoFsp    = 19 // YYYY-MM-DD HH:MM:SS
	MaxDatetimeWidthWithFsp  = 26 // YYYY-MM-DD HH:MM:SS[.fraction]
	MaxDatetimeFullWidth     = 29 // YYYY-MM-DD HH:MM:SS.###### AM
	MaxDurationWidthNoFsp    = 10 // HH:MM:SS
	MaxDurationWidthWithFsp  = 15 // HH:MM:SS[.fraction]
	MaxBlobWidth             = 16777216
	MaxBitDisplayWidth       = 64
	MaxFloatPrecisionLength  = 24
	MaxDoublePrecisionLength = 53
)
