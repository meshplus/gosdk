package abi

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
)

func TestTypeConversion(t *testing.T) {
	var gvalue string
	var gtype string

	gvalue = "11"
	gtype = "uint8"
	uint8Bin := TypeConversion(gvalue, gtype)
	if uint8Bin == common.LeftPadString(strconv.FormatUint(uint64(11), 16), UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "1"
	gtype = "uint32"
	uint32Bin := TypeConversion(gvalue, gtype)
	if uint32Bin == common.LeftPadString("1", UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "991"
	gtype = "uint64"
	uint64Bin := TypeConversion(gvalue, gtype)
	if uint64Bin == common.LeftPadString(strconv.FormatUint(991, 16), UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "12345"
	gtype = "int32"
	int32Bin := TypeConversion(gvalue, gtype)
	if int32Bin == common.LeftPadString(strconv.FormatInt(12345, 16), UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "91191"
	gtype = "int64"
	int64Bin := TypeConversion(gvalue, gtype)
	if int64Bin == common.LeftPadString(strconv.FormatInt(91191, 16), UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "1"
	gtype = "int128"
	int128Bin := TypeConversion(gvalue, gtype)
	if int128Bin == common.LeftPadString("1", UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "9999999999999999999999"
	gtype = "int256"
	bigIntBin := TypeConversion(gvalue, gtype)
	if bigIntBin == common.LeftPadString("21e19e0c9bab23fffff", UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "-11111119999999999999999999999"
	gtype = "int256"
	bigInt2Bin := TypeConversion(gvalue, gtype)
	if bigInt2Bin == common.LeftPad("dc1918d1dcd62d98de000001", UNIT, "f") {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "-1"
	gtype = "int256"
	intNegBin := TypeConversion(gvalue, gtype)
	if intNegBin == common.LeftPad("f", UNIT, "f") {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "1"
	gtype = "int256"
	intBin := TypeConversion(gvalue, gtype)
	if intBin == common.LeftPadString("1", UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "true"
	gtype = "bool"
	boolTBin := TypeConversion(gvalue, gtype)
	if boolTBin == common.LeftPadString("1", UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "false"
	gtype = "bool"
	boolFBin := TypeConversion(gvalue, gtype)
	if boolFBin == common.LeftPadString("0", UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "0x24602722816b6cad0e143ce9fabf31f6026ec622"
	gtype = "address"
	addrBin := TypeConversion(gvalue, gtype)
	if addrBin == common.LeftPadString(gvalue[2:], UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "Hello world !"
	gtype = "string"
	strBin := TypeConversion(gvalue, gtype)
	l := len(gvalue)
	lhex := strconv.FormatInt(int64(l), 16)
	payload := common.LeftPadString("20", UNIT)
	payload += common.LeftPadString(lhex, UNIT)
	payload += common.RightPadString(common.Bytes2Hex([]byte(gvalue)), UNIT)
	if strBin == payload {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "abc123"
	gtype = "bytes"
	l = len(gvalue)
	lhex = strconv.FormatInt(int64(l), 16)
	byteBin := TypeConversion(gvalue, gtype)
	payload = common.LeftPadString("20", UNIT)
	payload += common.LeftPadString(lhex, UNIT)
	payload += common.RightPadString(common.Bytes2Hex([]byte("abc123")), UNIT)
	if byteBin == payload {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "abc"
	gtype = "bytes4"
	byte4Bin := TypeConversion(gvalue, gtype)
	if byte4Bin == common.RightPadString(common.Bytes2Hex([]byte(gvalue)), UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "abcabcabcabcabcabcabcabc"
	gtype = "bytes32"
	byte32Bin := TypeConversion(gvalue, gtype)
	if byte32Bin == common.RightPadString(common.Bytes2Hex([]byte(gvalue)), UNIT) {
		t.Log("ok")
	} else {
		t.Error("error: invalid binary ")
	}

	gvalue = "1, 11"
	gtype = "uint32[]"
	arrayBin := TypeConversion(gvalue, gtype)
	encodedToUNIT32 := `00000000000000000000000000000000
                        00000000000000000000000000000020
                        00000000000000000000000000000000
                        00000000000000000000000000000002
                        00000000000000000000000000000000
                        00000000000000000000000000000001
                        00000000000000000000000000000000
                        0000000000000000000000000000000b`

	encodedToUNIT64 := `00000000000000000000000000000000
                        00000000000000000000000000000000
                        00000000000000000000000000000000
                        00000000000000000000000000000020
                        00000000000000000000000000000000
                        00000000000000000000000000000000
                        00000000000000000000000000000000
                        00000000000000000000000000000002
                        00000000000000000000000000000000
                        00000000000000000000000000000000
                        00000000000000000000000000000000
                        00000000000000000000000000000001
                        00000000000000000000000000000000
                        00000000000000000000000000000000
                        00000000000000000000000000000000
                        0000000000000000000000000000000b`

	if UNIT == 32 {
		assert.Equal(t, arrayBin, strings.Replace(strings.Replace(encodedToUNIT32, " ", "", -1), "\n", "", -1))
	}

	if UNIT == 64 {
		assert.Equal(t, arrayBin, strings.Replace(strings.Replace(encodedToUNIT64, " ", "", -1), "\n", "", -1))
	}
}

func TestFuncSelector(t *testing.T) {
	typeList := []string{"uint256"}
	if FuncSelector("set", typeList) == "60fe47b1" {
		t.Log("ok")
	} else {
		t.Error("error: func selector error")
	}
}

func TestNeg2Hex(t *testing.T) {
	fmt.Print(neg2Hex("111"))
}

func TestInt2Inverted(t *testing.T) {
	fmt.Print(int2Inverted(1))
}
