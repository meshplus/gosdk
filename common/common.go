package common

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
)

// Padding char char at the begin of the string
func LeftPad(str string, l int, char string) string {
	if l < len(str) {
		return str
	}
	zero := ""
	for i := 0; i < l-len(str); i += 1 {
		zero = zero + char
	}
	return zero + str
}

// Padding 0 char at the begin of the string
func LeftPadString(str string, l int) string {
	if l < len(str) {
		return str
	}
	zero := ""
	for i := 0; i < l-len(str); i += 1 {
		zero = zero + "0"
	}
	return zero + str
}

// Padding 0 char at the end of the string
func RightPadString(str string, l int) string {
	if l < len(str) {
		return str
	}
	zero := ""
	for i := 0; i < l-len(str); i += 1 {
		zero = zero + "0"
	}
	return str + zero
}

// Padding 0 char at the begin of the string
func LeftPadStringWithChar(str string, l int, c string) string {
	if l < len(str) {
		return str
	}
	zero := ""
	for i := 0; i < l-len(str); i += 1 {
		zero = zero + c
	}
	return zero + str
}

// Padding 0 char at the end of the string
func RightPadStringWithChar(str string, l int, c string) string {
	if l < len(str) {
		return str
	}
	zero := ""
	for i := 0; i < l-len(str); i += 1 {
		zero = zero + c
	}
	return str + zero
}

func SplitStringByInterval(str string, interval int) []string {
	var ret []string
	idx := 0
	length := len(str)
	for {
		if idx < length {
			if idx+interval <= length {
				ret = append(ret, str[idx:idx+interval])
			} else {
				ret = append(ret, str[idx:length])
			}
			idx += interval
		} else {
			break
		}
	}
	return ret
}

// Strip returns a slice of the byte s with all leading and
// trailing Unicode code points contained in cutset removed.
func Strip(in []byte, cutset string) []byte {
	str := string(in)
	out := strings.Trim(str, cutset)
	return []byte(out)
}

/*
	Random utils
*/
// Generate random string with specific length
func RandomString(length int) string {
	return fastRandomString(uint(length))
}

// Return a random choice with specific ratio, true of false
func RandomChoice(ratio float64) int {
	var tmp = make([]int, 0, 10)
	for i := 0; i < int(ratio*10); i += 1 {
		tmp = append(tmp, 0)
	}
	for i := int(ratio * 10); i < 10; i += 1 {
		tmp = append(tmp, 1)
	}
	return tmp[fastRandomIntn(len(tmp))]
}

// Return a random int in the specific range
func RandomInt64(lowerLimit, upperLimit int) int64 {
	return int64(fastRandomInt(lowerLimit, upperLimit))
}

func RandomNonce() int64 {
	return fastRandomInt63()
}

func RandomAddress() string {
	addr := StringToHex(fastRandomAddr())
	return addr
}

// other
var tt256m1 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

/*
	Common
*/

func StringToHex(s string) string {
	if len(s) >= 2 && s[:2] == "0x" {
		return s
	} else {
		return "0x" + s
	}
}

func ParseData(data ...interface{}) (ret []byte) {
	for _, item := range data {
		switch t := item.(type) {
		case string:
			var str []byte
			if IsHex(t) {
				str = Hex2Bytes(t[2:])
			} else {
				str = []byte(t)
			}

			ret = append(ret, RightPadBytes(str, 32)...)
		case []byte:
			ret = append(ret, LeftPadBytes(t, 32)...)
		}
	}

	return
}

func IsHex(str string) bool {
	l := len(str)
	return l >= 4 && l%2 == 0 && str[0:2] == "0x"
}

func BytesToBig(data []byte) *big.Int {
	n := new(big.Int)
	n.SetBytes(data)

	return n
}

func String2Big(num string) *big.Int {
	n := new(big.Int)
	n.SetString(num, 0)
	return n
}

func U256(x *big.Int) *big.Int {
	//if x.Cmp(Big0) < 0 {
	//		return new(big.Int).Add(tt256, x)
	//	}

	x.And(x, tt256m1)

	return x
}

func Bytes2Big(data []byte) *big.Int { return BytesToBig(data) }

func BigD(data []byte) *big.Int { return BytesToBig(data) }

func HexToString(s string) string {
	if len(s) >= 2 && s[:2] == "0x" {
		return s[2:]
	} else {
		return s
	}
}

func ConvertStringSliceToIntSlice(input []string) []int {
	var ret []int
	for _, elem := range input {
		e, err := strconv.ParseInt(elem, 10, 0)
		if err != nil {
			//logger.Error("invalid config args")
			return nil
		}
		ret = append(ret, int(e))
	}
	return ret
}

func ConvertStringSliceToBoolSlice(input []string) []bool {
	var ret []bool
	for _, elem := range input {
		e, err := strconv.ParseBool(elem)
		if err != nil {
			//logger.Error("invalid config args")
			return nil
		}
		ret = append(ret, e)
	}
	return ret
}

func RemoveSubString(str string, begin, end int) (error, string) {
	length := len(str)
	if length <= end || length <= begin || begin < 0 || begin > end {
		return errors.New("invalid params"), ""
	}
	return nil, str[0:begin] + str[end+1:length]
}

func DelHex(str string) string {
	strCopy := str
	if len(strCopy) > 1 {
		if strCopy[0:2] == "0x" {
			strCopy = strCopy[2:]
		}
	}
	return strCopy
}

func RandInt(num int) int {
	return fastRandomIntn(num)
}
