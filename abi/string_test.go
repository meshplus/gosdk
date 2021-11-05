package abi

import (
	"encoding/hex"
	"fmt"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestABI_Encode(t *testing.T) {
	abiJSON := `
[
	{"constant":false,"inputs":[{"name":"p1","type":"int256"},{"name":"p2","type":"int256[]"},{"name":"p3","type":"int256[3]"}],"name":"typeInt256","outputs":[{"name":"","type":"int256"},{"name":"","type":"int256[]"},{"name":"","type":"int256[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
    {"constant":false,"inputs":[{"name":"p1","type":"uint32"},{"name":"p2","type":"uint32[]"},{"name":"p3","type":"uint32[3]"}],"name":"typeUint32","outputs":[{"name":"","type":"uint32"},{"name":"","type":"uint32[]"},{"name":"","type":"uint32[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
    {"constant":false,"inputs":[{"name":"p1","type":"int64"},{"name":"p2","type":"int64[]"},{"name":"p3","type":"int64[3]"}],"name":"typeInt64","outputs":[{"name":"","type":"int64"},{"name":"","type":"int64[]"},{"name":"","type":"int64[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
    {"constant":false,"inputs":[{"name":"p1","type":"bytes1"},{"name":"p2","type":"bytes1[]"},{"name":"p3","type":"bytes1[3]"}],"name":"typeBytes1","outputs":[{"name":"","type":"bytes1"},{"name":"","type":"bytes1[]"},{"name":"","type":"bytes1[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bytes"}],"name":"typeBytes","outputs":[{"name":"","type":"bytes"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"address"},{"name":"p2","type":"address[]"},{"name":"p3","type":"address[3]"}],"name":"typeAddress","outputs":[{"name":"","type":"address"},{"name":"","type":"address[]"},{"name":"","type":"address[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
 	{"constant":false,"inputs":[{"name":"p1","type":"int8"},{"name":"p2","type":"int8[]"},{"name":"p3","type":"int8[3]"}],"name":"typeInt8","outputs":[{"name":"","type":"int8"},{"name":"","type":"int8[]"},{"name":"","type":"int8[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"uint16"},{"name":"p2","type":"uint16[]"},{"name":"p3","type":"uint16[3]"}],"name":"typeUint16","outputs":[{"name":"","type":"uint16"},{"name":"","type":"uint16[]"},{"name":"","type":"uint16[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bytes24"},{"name":"p2","type":"bytes24[]"},{"name":"p3","type":"bytes24[3]"}],"name":"typeBytes24","outputs":[{"name":"","type":"bytes24"},{"name":"","type":"bytes24[]"},{"name":"","type":"bytes24[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"uint64"},{"name":"p2","type":"uint64[]"},{"name":"p3","type":"uint64[3]"}],"name":"typeUint64","outputs":[{"name":"","type":"uint64"},{"name":"","type":"uint64[]"},{"name":"","type":"uint64[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"int16"},{"name":"p2","type":"int16[]"},{"name":"p3","type":"int16[3]"}],"name":"typeInt16","outputs":[{"name":"","type":"int16"},{"name":"","type":"int16[]"},{"name":"","type":"int16[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"string"}],"name":"typeString","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bytes2"},{"name":"p2","type":"bytes2[]"},{"name":"p3","type":"bytes2[3]"}],"name":"typeBytes2","outputs":[{"name":"","type":"bytes2"},{"name":"","type":"bytes2[]"},{"name":"","type":"bytes2[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"int32"},{"name":"p2","type":"int32[]"},{"name":"p3","type":"int32[3]"}],"name":"typeInt32","outputs":[{"name":"","type":"int32"},{"name":"","type":"int32[]"},{"name":"","type":"int32[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"uint8"},{"name":"p2","type":"uint8[]"},{"name":"p3","type":"uint8[3]"}],"name":"typeUint8","outputs":[{"name":"","type":"uint8"},{"name":"","type":"uint8[]"},{"name":"","type":"uint8[2]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"uint256"},{"name":"p2","type":"uint256[]"},{"name":"p3","type":"uint256[3]"}],"name":"typeUint256","outputs":[{"name":"","type":"uint256"},{"name":"","type":"uint256[]"},{"name":"","type":"uint256[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bytes7"},{"name":"p2","type":"bytes7[]"},{"name":"p3","type":"bytes7[3]"}],"name":"typeBytes7","outputs":[{"name":"","type":"bytes7"},{"name":"","type":"bytes7[]"},{"name":"","type":"bytes7[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bytes32"},{"name":"p2","type":"bytes32[]"},{"name":"p3","type":"bytes32[3]"}],"name":"typeBytes32","outputs":[{"name":"","type":"bytes32"},{"name":"","type":"bytes32[]"},{"name":"","type":"bytes32[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bytes127"},{"name":"p2","type":"bytes127[]"},{"name":"p3","type":"bytes127[3]"}],"name":"typeBytes127","outputs":[{"name":"","type":"bytes127"},{"name":"","type":"bytes127[]"},{"name":"","type":"bytes127[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
    {"constant":false,"inputs":[{"name":"p1","type":"bool"},{"name":"p2","type":"bool[]"},{"name":"p3","type":"bool[3]"}],"name":"typeBool","outputs":[{"name":"","type":"bool"},{"name":"","type":"bool[]"},{"name":"","type":"bool[3]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"key","type":"bytes"},{"name":"value","type":"bytes"}],"name":"setHash","outputs":[{"name":"","type":"bytes"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"key","type":"string"},{"name":"value","type":"string"}],"name":"setHash","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"int256[2][3]"},{"name":"p2","type":"int256[2][]"}],"name":"nestedInt256","outputs":[{"name":"","type":"int256[2][3]"},{"name":"","type":"int256[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"uint256[2][3]"},{"name":"p2","type":"uint256[2][]"}],"name":"nestedUint256","outputs":[{"name":"","type":"uint256[2][3]"},{"name":"","type":"uint256[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"uint8[2][3]"},{"name":"p2","type":"uint8[2][]"}],"name":"nestedUint8","outputs":[{"name":"","type":"uint8[2][3]"},{"name":"","type":"uint8[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"int8[2][3]"},{"name":"p2","type":"int8[2][]"}],"name":"nestedInt8","outputs":[{"name":"","type":"int8[2][3]"},{"name":"","type":"int8[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bytes32[2][3]"},{"name":"p2","type":"bytes32[2][]"}],"name":"nestedBytes32","outputs":[{"name":"","type":"bytes32[2][3]"},{"name":"","type":"bytes32[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bool[2][3]"},{"name":"p2","type":"bool[2][]"}],"name":"nestedBool","outputs":[{"name":"","type":"bool[2][3]"},{"name":"","type":"bool[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"uint64[2][3]"},{"name":"p2","type":"uint64[2][]"}],"name":"nestedUint64","outputs":[{"name":"","type":"uint64[2][3]"},{"name":"","type":"uint64[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"int64[2][3]"},{"name":"p2","type":"int64[2][]"}],"name":"nestedInt64","outputs":[{"name":"","type":"int64[2][3]"},{"name":"","type":"int64[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"address[2][3]"},{"name":"p2","type":"address[2][]"}],"name":"nestedAddress","outputs":[{"name":"","type":"address[2][3]"},{"name":"","type":"address[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"name":"p1","type":"bytes1[2][3]"},{"name":"p2","type":"bytes1[2][]"}],"name":"nestedBytes1","outputs":[{"name":"","type":"bytes1[2][3]"},{"name":"","type":"bytes1[2][]"}],"payable":false,"stateMutability":"nonpayable","type":"function"}
]`
	abi, err := JSON(strings.NewReader(abiJSON))
	if err != nil {
		t.Fatal(err)
	}

	var encodeTests = []struct {
		funcName string
		equal    []interface{}
		input    []interface{}
	}{
		{
			funcName: "typeUint8",
			equal:    []interface{}{uint8(1), []uint8{1, 1}, [3]uint8{1, 1, 1}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeUint16",
			equal:    []interface{}{uint16(1), []uint16{1, 1}, [3]uint16{1, 1, 1}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeUint32",
			equal:    []interface{}{uint32(1), []uint32{1, 1}, [3]uint32{1, 1, 1}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeUint64",
			equal:    []interface{}{uint64(1), []uint64{1, 1}, [3]uint64{1, 1, 1}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeUint256",
			equal:    []interface{}{big.NewInt(1), []*big.Int{big.NewInt(1), big.NewInt(1)}, [3]*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeInt8",
			equal:    []interface{}{int8(1), []int8{1, 1}, [3]int8{1, 1, 1}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeInt16",
			equal:    []interface{}{int16(1), []int16{1, 1}, [3]int16{1, 1, 1}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeInt32",
			equal:    []interface{}{int32(1), []int32{1, 1}, [3]int32{1, 1, 1}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeInt64",
			equal:    []interface{}{int64(1), []int64{1, 1}, [3]int64{1, 1, 1}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeInt256",
			equal:    []interface{}{big.NewInt(1), []*big.Int{big.NewInt(1), big.NewInt(1)}, [3]*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeAddress",
			equal:    []interface{}{common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, []common.Address{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}, [3]common.Address{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeBytes1",
			equal:    []interface{}{[1]byte{'1'}, [][1]byte{{'1'}, {'1'}}, [3][1]byte{{'1'}, {'1'}, {'1'}}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeBytes2",
			equal:    []interface{}{[2]byte{'1'}, [][2]byte{{'1'}, {'1'}}, [3][2]byte{{'1'}, {'1'}, {'1'}}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeBytes7",
			equal:    []interface{}{[7]byte{'1'}, [][7]byte{{'1'}, {'1'}}, [3][7]byte{{'1'}, {'1'}, {'1'}}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeBytes24",
			equal:    []interface{}{[24]byte{'1'}, [][24]byte{{'1'}, {'1'}}, [3][24]byte{{'1'}, {'1'}, {'1'}}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeBytes32",
			equal:    []interface{}{[32]byte{'1'}, [][32]byte{{'1'}, {'1'}}, [3][32]byte{{'1'}, {'1'}, {'1'}}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeBytes127",
			equal:    []interface{}{[127]byte{'1'}, [][127]byte{{'1'}, {'1'}}, [3][127]byte{{'1'}, {'1'}, {'1'}}},
			input:    []interface{}{"1", []interface{}{"1", "1"}, []interface{}{"1", "1", "1"}},
		}, {
			funcName: "typeBool",
			equal:    []interface{}{true, []bool{true, false}, [3]bool{false, true, false}},
			input:    []interface{}{"true", []interface{}{"true", "false"}, []interface{}{"false", "true", "false"}},
		}, {
			funcName: "typeString",
			equal:    []interface{}{"test"},
			input:    []interface{}{"test"},
		}, {
			funcName: "typeBytes",
			equal:    []interface{}{[]byte("test")},
			input:    []interface{}{hex.EncodeToString([]byte("test"))},
		}, {
			funcName: "nestedUint8",
			equal:    []interface{}{[3][2]uint8{{127, 0}, {255, 0}, {127, 255}}, [][2]uint8{{127, 0}, {255, 0}, {127, 255}}},
			input:    []interface{}{[][]interface{}{{"127", "0"}, {"255", "0"}, {"127", "255"}}, [][]interface{}{{"127", "0"}, {"255", "0"}, {"127", "255"}}},
		}, {
			funcName: "nestedUint64",
			equal:    []interface{}{[3][2]uint64{{127, 0}, {255, 0}, {127, 255}}, [][2]uint64{{127, 0}, {255, 0}, {127, 255}}},
			input:    []interface{}{[][]interface{}{{"127", "0"}, {"255", "0"}, {"127", "255"}}, [][]interface{}{{"127", "0"}, {"255", "0"}, {"127", "255"}}},
		}, {
			funcName: "nestedUint256",
			equal:    []interface{}{[3][2]*big.Int{{big.NewInt(127), big.NewInt(0)}, {big.NewInt(255), big.NewInt(0)}, {big.NewInt(127), big.NewInt(255)}}, [][2]*big.Int{{big.NewInt(127), big.NewInt(0)}, {big.NewInt(255), big.NewInt(0)}, {big.NewInt(127), big.NewInt(255)}}},
			input:    []interface{}{[][]interface{}{{"127", "0"}, {"255", "0"}, {"127", "255"}}, [][]interface{}{{"127", "0"}, {"255", "0"}, {"127", "255"}}},
		}, {
			funcName: "nestedInt8",
			equal:    []interface{}{[3][2]int8{{127, 0}, {-127, 0}, {127, -127}}, [][2]int8{{127, 0}, {-127, 0}, {127, -127}}},
			input:    []interface{}{[][]interface{}{{"127", "0"}, {"-127", "0"}, {"127", "-127"}}, [][]interface{}{{"127", "0"}, {"-127", "0"}, {"127", "-127"}}},
		}, {
			funcName: "nestedInt64",
			equal:    []interface{}{[3][2]int64{{127, 0}, {-127, 0}, {127, -127}}, [][2]int64{{127, 0}, {-127, 0}, {127, -127}}},
			input:    []interface{}{[][]interface{}{{"127", "0"}, {"-127", "0"}, {"127", "-127"}}, [][]interface{}{{"127", "0"}, {"-127", "0"}, {"127", "-127"}}},
		}, {
			funcName: "nestedInt256",
			equal:    []interface{}{[3][2]*big.Int{{big.NewInt(127), big.NewInt(0)}, {big.NewInt(-127), big.NewInt(0)}, {big.NewInt(127), big.NewInt(-127)}}, [][2]*big.Int{{big.NewInt(127), big.NewInt(0)}, {big.NewInt(-127), big.NewInt(0)}, {big.NewInt(127), big.NewInt(-127)}}},
			input:    []interface{}{[][]interface{}{{"127", "0"}, {"-127", "0"}, {"127", "-127"}}, [][]interface{}{{"127", "0"}, {"-127", "0"}, {"127", "-127"}}},
		}, {
			funcName: "nestedBool",
			equal:    []interface{}{[3][2]bool{{true, false}, {false, true}, {true, false}}, [][2]bool{{true, false}, {false, true}, {true, false}}},
			input:    []interface{}{[][]interface{}{{"true", "false"}, {"false", "true"}, {"true", "false"}}, [][]interface{}{{"true", "false"}, {"false", "true"}, {"true", "false"}}},
		}, {
			funcName: "nestedAddress",
			equal: []interface{}{
				[3][2]common.Address{
					{
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					},
					{
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					},
					{
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					},
				},
				[][2]common.Address{
					{
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					},
					{
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					},
					{
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
						common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					},
				}},
			input: []interface{}{
				[][]interface{}{{"1", "1"}, {"1", "1"}, {"1", "1"}},
				[][]interface{}{{"1", "1"}, {"1", "1"}, {"1", "1"}},
			},
		}, {
			funcName: "nestedBytes1",
			equal: []interface{}{
				[3][2][1]byte{{{'1'}, {'1'}}, {{'1'}, {'1'}}, {{'1'}, {'1'}}},
				[][2][1]byte{{{'1'}, {'1'}}, {{'1'}, {'1'}}, {{'1'}, {'1'}}},
			},
			input: []interface{}{
				[][]interface{}{{"1", "1"}, {"1", "1"}, {"1", "1"}},
				[][]interface{}{{"1", "1"}, {"1", "1"}, {"1", "1"}},
			},
		}, {
			funcName: "nestedBytes32",
			equal: []interface{}{
				[3][2][32]byte{{{'1'}, {'1'}}, {{'1'}, {'1'}}, {{'1'}, {'1'}}},
				[][2][32]byte{{{'1'}, {'1'}}, {{'1'}, {'1'}}, {{'1'}, {'1'}}},
			},
			input: []interface{}{
				[][]interface{}{{"1", "1"}, {"1", "1"}, {"1", "1"}},
				[][]interface{}{{"1", "1"}, {"1", "1"}, {"1", "1"}},
			},
		},
	}

	for _, test := range encodeTests {
		t.Run(test.funcName, func(t *testing.T) {
			actual, err := abi.Encode(test.funcName, test.input...)
			assert.NoError(t, err)

			expect, err := abi.Pack(test.funcName, test.equal...)
			assert.NoError(t, err)

			assert.Equal(t, expect, actual)
		})
	}

}

var decodeTests = []unpackTest{
	{
		def:  `[{ "type": "bool" }]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000001",
		want: true,
	},
	{
		def:  `[{ "type": "bool" }]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000000",
		want: false,
	},
	{
		def:  `[{ "type": "bool" }]`,
		enc:  "0000000000000000000000000000000000000000000000000001000000000001",
		want: false,
		err:  "abi: improperly encoded boolean value",
	},
	{
		def:  `[{ "type": "bool" }]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000003",
		want: false,
		err:  "abi: improperly encoded boolean value",
	},
	{
		def:  `[{"type": "uint32"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000001",
		want: uint32(1),
	},
	{
		def:  `[{"type": "uint17"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000001",
		want: big.NewInt(1),
	},
	{
		def:  `[{"type": "int32"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000001",
		want: int32(1),
	},
	{
		def:  `[{"type": "int17"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000001",
		want: big.NewInt(1),
	},
	{
		def:  `[{"type": "int256"}]`,
		enc:  "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		want: big.NewInt(-1),
	},
	{
		def:  `[{"type": "address"}]`,
		enc:  "0000000000000000000000000100000000000000000000000000000000000000",
		want: common.Address{1},
	},
	{
		def:  `[{"type": "bytes32"}]`,
		enc:  "0100000000000000000000000000000000000000000000000000000000000000",
		want: [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		def:  `[{"type": "bytes"}]`,
		enc:  "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200100000000000000000000000000000000000000000000000000000000000000",
		want: common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000"),
	},
	{
		def:  `[{"type": "bytes32"}]`,
		enc:  "0100000000000000000000000000000000000000000000000000000000000000",
		want: [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		def:  `[{"type": "function"}]`,
		enc:  "0100000000000000000000000000000000000000000000000000000000000000",
		want: [24]byte{1},
	},
	// slices
	{
		def:  `[{"type": "uint8[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []uint8{1, 2},
	},
	{
		def:  `[{"type": "uint8[2]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2]uint8{1, 2},
	},
	// multi dimensional, if these pass, all types that don't require length prefix should pass
	{
		def:  `[{"type": "uint8[][]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [][]uint8{{1, 2}, {1, 2}},
	},
	{
		def:  `[{"type": "uint8[2][2]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2][2]uint8{{1, 2}, {1, 2}},
	},
	{
		def:  `[{"type": "uint8[][2]"}]`,
		enc:  "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001",
		want: [2][]uint8{{1}, {1}},
	},
	{
		def:  `[{"type": "uint8[2][]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [][2]uint8{{1, 2}},
	},
	{
		def:  `[{"type": "uint8[2][]"}]`,
		enc:  "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [][2]uint8{{1, 2}, {1, 2}},
	},
	{
		def:  `[{"type": "uint16[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []uint16{1, 2},
	},
	{
		def:  `[{"type": "uint16[2]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2]uint16{1, 2},
	},
	{
		def:  `[{"type": "uint32[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []uint32{1, 2},
	},
	{
		def:  `[{"type": "uint32[2]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2]uint32{1, 2},
	},
	{
		def:  `[{"type": "uint32[2][3][4]"}]`,
		enc:  "000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000700000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000b000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000d000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000f000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001300000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000015000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000170000000000000000000000000000000000000000000000000000000000000018",
		want: [4][3][2]uint32{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}, {{13, 14}, {15, 16}, {17, 18}}, {{19, 20}, {21, 22}, {23, 24}}},
	},
	{
		def:  `[{"type": "uint64[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []uint64{1, 2},
	},
	{
		def:  `[{"type": "uint64[2]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2]uint64{1, 2},
	},
	{
		def:  `[{"type": "uint256[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []*big.Int{big.NewInt(1), big.NewInt(2)},
	},
	{
		def:  `[{"type": "uint256[3]"}]`,
		enc:  "000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003",
		want: [3]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
	},
	{
		def:  `[{"type": "int8[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []int8{1, 2},
	},
	{
		def:  `[{"type": "int8[2]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2]int8{1, 2},
	},
	{
		def:  `[{"type": "int16[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []int16{1, 2},
	},
	{
		def:  `[{"type": "int16[2]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2]int16{1, 2},
	},
	{
		def:  `[{"type": "int32[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []int32{1, 2},
	},
	{
		def:  `[{"type": "int32[2]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2]int32{1, 2},
	},
	{
		def:  `[{"type": "int64[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []int64{1, 2},
	},
	{
		def:  `[{"type": "int64[2]"}]`,
		enc:  "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: [2]int64{1, 2},
	},
	{
		def:  `[{"type": "int256[]"}]`,
		enc:  "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		want: []*big.Int{big.NewInt(1), big.NewInt(2)},
	},
	{
		def:  `[{"type": "int256[3]"}]`,
		enc:  "000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003",
		want: [3]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
	},
	// multiple return value
	{
		def:  `[{"type": "uint64"}, {"type": "uint64[]"}, {"type": "uint64[3]"}]`,
		enc:  "000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001",
		want: []interface{}{uint64(1), []uint64{1}, [3]uint64{1, 1, 1}},
	},
	{
		def:  `[{"type": "int256"}, {"type": "int256[]"}, {"type": "int256[3]"}]`,
		enc:  "000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001",
		want: []interface{}{big.NewInt(1), []*big.Int{big.NewInt(1)}, [3]*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)}},
	},
}

func TestABI_Decode(t *testing.T) {
	for i, test := range decodeTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			def := fmt.Sprintf(`[{ "name" : "method", "outputs": %s}]`, test.def)
			abi, err := JSON(strings.NewReader(def))
			if err != nil {
				t.Fatalf("invalid ABI definition %s: %v", def, err)
			}
			encb, err := hex.DecodeString(test.enc)
			if err != nil {
				t.Fatalf("invalid hex: %s" + test.enc)
			}
			var output interface{}
			method, err := abi.GetMethod("method")
			if err != nil {
				t.Error(err)
			}
			switch len(method.Outputs) {
			case 0:
				// do nothing
			case 1:
				// single return
				output, err = abi.Decode("method", encb)
			default:
				// tuple return
				output, err = abi.Decode("method", encb)
			}

			if err := test.checkError(err); err != nil {
				t.Errorf("test %d (%v) failed: %v", i, test.def, err)
				return
			}
			if !reflect.DeepEqual(test.want, output) {
				t.Errorf("test %d (%v) failed: expected %v, got %v", i, test.def, test.want, output)
			}
		})
	}
}

func TestNewFixedBytes(t *testing.T) {
	for n := 0; n < 65; n++ {
		newFixedBytes(n, "123")
	}
}

func BenchmarkString_newFixedBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		newFixedBytes(32, "123")
	}
}

func BenchmarkString_newFixedBytesWithReflect(b *testing.B) {
	for n := 0; n < b.N; n++ {
		newFixedBytesWithReflect(32, "123")
	}
}
