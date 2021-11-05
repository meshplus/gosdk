package abi

import (
	"fmt"
	"github.com/meshplus/crypto-standard/hash/sha3"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/meshplus/gosdk/common"
)

const UNIT = 64

var logger = common.GetLogger("abi")

// TypeConversion conversion go type to binary
func TypeConversion(gvalue string, gtype string) string {

	var payload string

	// Bytes<m>
	// 0 < m <= 32
	if strings.HasPrefix(gtype, "bytes") && gtype != "bytes" {
		// parse m
		m, err := strconv.ParseInt(gtype[len("bytes"):], 10, UNIT)
		if err != nil || len(gvalue) > int(m) {
			logger.Error("parse bytes error")
			return ""
		}
		payload += common.RightPadString(common.Bytes2Hex([]byte(gvalue)), UNIT)
		return payload
	}

	if isArray(gtype) {
		// value -> 1,2,3,4,5,6
		eleType := fetchArrayType(gtype)
		vList := strings.Split(strings.Replace(gvalue, " ", "", -1), ",")
		payload += common.LeftPadString("0", UNIT)
		payload += common.LeftPadString("20", UNIT)
		payload += common.LeftPadString("0", UNIT)
		payload += common.LeftPadString(string(strconv.FormatUint(uint64(len(vList)), 16)), UNIT)
		for _, v := range vList {
			payload += common.LeftPadString("0", UNIT)
			payload += TypeConversion(v, eleType)
		}
		return payload
	}

	switch gtype {
	case "uint8", "uint16", "uint32", "uint64":
		num, _ := strconv.ParseUint(gvalue, 10, UNIT)
		hex := strconv.FormatUint(num, 16)
		payload = common.LeftPadString(hex, UNIT)
	case "uint128", "uint256":
		num, _ := big.NewInt(0).SetString(gvalue, 10)
		hex := num.Text(16)
		if len(hex) >= UNIT {
			logger.Error("uint256 overflow")
			return ""
		}
		payload = common.LeftPadString(hex, UNIT)
	case "int8", "int16", "int32", "int64":
		num, _ := strconv.ParseInt(gvalue, 10, UNIT)
		if num < 0 {
			v := int(math.Abs(float64(num)))
			strInverted := int2Inverted(v)
			hex := neg2Hex(strInverted)
			payload = common.LeftPad(hex, UNIT, "f")
		} else {
			hex := strconv.FormatUint(uint64(num), 16)
			payload = common.LeftPadString(hex, UNIT)
		}
	case "int128", "int256":
		num, _ := big.NewInt(0).SetString(gvalue, 10)
		if num.Cmp(big.NewInt(0)) < 0 {
			v := num.Abs(num)
			bstr := []byte(v.Text(2))
			bint := new(big.Int)
			bc := binaryComplement(string(bstr))
			binaryStr, _ := bint.SetString(common.LeftPad(bc, UNIT*4, "1"), 2)
			hex := fmt.Sprintf("%x", binaryStr)
			if len(hex) > UNIT {
				logger.Error("int256 overflow")
				return ""
			}
			payload = common.LeftPad(hex, UNIT, "f")
		} else {
			payload = common.LeftPadString(num.Text(16), UNIT)
		}
	case "bool":
		if v, ok := strconv.ParseBool(gvalue); v && ok == nil {
			payload = common.LeftPadString("1", UNIT)
		} else {
			payload = common.LeftPadString("0", UNIT)
		}
	case "address":
		if common.IsHexAddress(gvalue) {
			gvalue = strings.TrimPrefix(gvalue, "0x")
			payload = common.LeftPadString(gvalue, UNIT)
		} else {
			logger.Error("invalid address")
			return ""
		}
	case "string":
		l := len(gvalue)
		lhex := strconv.FormatInt(int64(l), 16)
		payload += common.LeftPadString("20", UNIT)
		payload += common.LeftPadString(lhex, UNIT)
		payload += common.RightPadString(common.Bytes2Hex([]byte(gvalue)), UNIT)
	case "bytes":
		var bitSize int
		l := len([]byte(gvalue))
		lhex := strconv.FormatInt(int64(l), 16)
		payload += common.LeftPadString("20", UNIT)
		payload += common.LeftPadString(lhex, UNIT)
		if l%UNIT == 0 {
			bitSize = l / UNIT
		} else {
			bitSize = l/UNIT + 1
		}
		payload += common.RightPadString(common.Bytes2Hex([]byte(gvalue)), UNIT*bitSize)
	default:
		logger.Error("invalid type")
		return ""
	}
	return payload
}

// FuncSelector convert func and args to ABI
// set(uint256) -> 60fe47b1
func FuncSelector(funcName string, argsType []string) string {
	funcString := funcName + "(" + strings.Join(argsType, ",") + ")"
	hw := sha3.NewKeccak256()
	_, _ = hw.Write([]byte(funcString))
	return common.Bytes2Hex(hw.Sum([]byte{})[0:4])
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func binaryComplement(a string) string {
	bstr := []byte(a)
	for i := 0; i < len(bstr); i++ {
		if bstr[i] == '0' {
			bstr[i] = '1'
		} else {
			bstr[i] = '0'
		}
	}
	last := len(bstr) - 1
	if bstr[last] == '1' {
		bstr[last] = '0'
		for i := last; i > 0; i-- {
			if bstr[i-1] == '1' {
				bstr[i-1] = '0'
			} else {
				bstr[i-1] = '1'
				break
			}
		}
	} else {
		bstr[last] = '1'
	}
	return string(bstr)
}

func neg2Hex(binaryStr string) string {
	// 求补码
	bc := binaryComplement(binaryStr)
	bytes := common.LeftPad(bc, 32, "1")
	vvv, err := strconv.ParseInt(bytes, 2, 64)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return strconv.FormatInt(int64(vvv), 16)
}

func int2Inverted(v int) string {
	rbinaryStr := ""
	// 转二进制（需要反转才是正确的二进制）
	for q := v; q > 0; q = q / 2 {
		rbinaryStr += strconv.Itoa(q % 2)
	}
	return reverseString(rbinaryStr)
}

func isArray(gtype string) bool {
	// gtype like this []int32
	return strings.HasSuffix(gtype, "[]")
}

func fetchArrayType(gtype string) string {
	if strings.HasSuffix(gtype, "[]") {
		return gtype[:len(gtype)-2]
	}
	return ""
}
