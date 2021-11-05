package common

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"math/big"
	"testing"
)

func TestLeftPadStringWithChar(t *testing.T) {
	str := "1111"
	l1 := 5
	l2 := 2
	c := "s"
	assert.Equal(t, LeftPadStringWithChar(str, l1, c), "s1111")
	assert.Equal(t, LeftPadStringWithChar(str, l2, c), "1111")
}

func TestRightPadStringWithChar(t *testing.T) {
	str := "1111"
	l1 := 5
	l2 := 2
	c := "s"
	assert.Equal(t, RightPadStringWithChar(str, l1, c), "1111s")
	assert.Equal(t, RightPadStringWithChar(str, l2, c), "1111")

}

func TestSplitStringByInterval(t *testing.T) {
	str := "1111"
	l2 := 2
	assert.Equal(t, SplitStringByInterval(str, l2), []string{"11", "11"})
}

func TestStrip(t *testing.T) {
	str := "1111"
	l2 := "2"
	//assert.Equal(t, (Strip([]byte(str), l2), []byte{49, 49, 49, 49})
	fmt.Println(Strip([]byte(str), l2))
}

func TestRandomChoice(t *testing.T) {

	fmt.Println(RandomChoice(5 / 6))
}

func TestRandomString(t *testing.T) {

	fmt.Println(RandomString(1))
}

func TestRandomInt64(t *testing.T) {

	fmt.Println(RandomInt64(1, 3))
}

func TestRandomNonce(t *testing.T) {

	fmt.Println(RandomNonce())
}

func TestRandomAddress(t *testing.T) {

	fmt.Println(RandomAddress())
}

func TestStringToHex(t *testing.T) {

	assert.Equal(t, StringToHex("1"), "0x1")
}

func TestParseDataT(t *testing.T) {

	fmt.Println(ParseData([]byte("22")))
	fmt.Println(ParseData("0x123"))

}

func TestIsHexT(t *testing.T) {

	assert.Equal(t, IsHex("1"), false)
	assert.Equal(t, IsHex("0x123"), false)

}

func TestBytesToBig(t *testing.T) {

	assert.Equal(t, BytesToBig([]byte("22")), big.NewInt(12850))
}

func TestString2Big(t *testing.T) {

	assert.Equal(t, String2Big("ss"), big.NewInt(0))
}

func TestU256(t *testing.T) {

	assert.Equal(t, U256(BytesToBig([]byte("11"))), big.NewInt(12593))
}

func TestHexToString(t *testing.T) {

	assert.Equal(t, HexToString("0x11"), "11")
	assert.Equal(t, HexToString("1"), "1")

}

func TestRemoveSubString(t *testing.T) {

	_, ans := RemoveSubString("222", 1, 2)
	assert.Equal(t, ans, "2")
}

func TestConvertStringSliceToIntSlice(t *testing.T) {
	assert.Equal(t, ConvertStringSliceToIntSlice([]string{"2", "2"}), []int{2, 2})
}

func TestConvertStringSliceToBoolSlice(t *testing.T) {
	assert.Equal(t, ConvertStringSliceToBoolSlice([]string{"1", "1"}), []bool{true, true})
}
