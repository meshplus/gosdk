package abi

import (
	"fmt"
	"github.com/meshplus/gosdk/common"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// funcName is the function name of contract
// args is the arguments of function, should be string or []interface{}
func (abi ABI) Encode(funcName string, args ...interface{}) ([]byte, error) {

	var method Method

	if funcName == "" {
		method = abi.Constructor
	} else {
		methodPtr, err := abi.GetMethod(funcName)
		if err != nil {
			return nil, err
		}
		method = *methodPtr
	}

	if len(method.Inputs) > len(args) {
		return nil, fmt.Errorf("the num of inputs is %v, expectd %v", len(method.Inputs), len(args))
	}

	typedArgs := make([]interface{}, len(method.Inputs))

	for idx, input := range method.Inputs {
		typedArgs[idx] = convert(input.Type, args[idx])
	}

	payload, err := abi.Pack(funcName, typedArgs...)
	return payload, err
}

// funcName is the function name of contract
// result is the return value of invocation solidity contract
func (abi ABI) Decode(funcName string, result []byte) (interface{}, error) {
	method, err := abi.GetMethod(funcName)

	if err != nil {
		return nil, err
	}

	outputLen := len(method.Outputs)
	var elem interface{}

	if outputLen == 1 {
		// single return
		elem = reflect.New(method.Outputs[0].Type.Type).Elem().Interface()
		err := abi.Unpack(&elem, funcName, result)
		return elem, err
	} else {
		// tuple return
		ptrs := make([]interface{}, outputLen)
		values := make([]reflect.Value, outputLen)
		elems := make([]interface{}, outputLen)
		for idx, val := range method.Outputs {
			value := reflect.New(val.Type.Type)
			values[idx] = value
			ptrs[idx] = value.Interface()
		}
		err := abi.Unpack(&ptrs, funcName, result)

		for idx, val := range values {
			elems[idx] = val.Elem().Interface()
		}
		return elems, err
	}
}

// convert val into target type through certain method
// support: array / slice / bytesN / basic type
// not support: nested array or slice
func convert(t Type, input interface{}) interface{} {
	// array or slice
	switch t.T {
	case ArrayTy:
		// make sure that the length of input equals to the t.Size
		var (
			fmtVal = make([]interface{}, t.Size)
			idx    int
		)
		reflectInput := reflect.ValueOf(input)
		switch reflectInput.Kind() {
		case reflect.String:
			if t.Size >= 1 {
				fmtVal[idx] = input
				idx++
			}
		case reflect.Slice:
			valLen := reflectInput.Len()
			var formatLen int
			if valLen < t.Size {
				formatLen = valLen
			} else {
				formatLen = t.Size
			}
			for idx = 0; idx < formatLen; idx++ {
				fmtVal[idx] = reflectInput.Index(idx).Interface()
			}
		}

		// complete input with default "" (empty string)
		for i := idx; i < t.Size; i++ {
			fmtVal[idx] = ""
		}
		// build the array (not slice)
		data := reflect.New(t.Type).Elem()
		for idx, val := range fmtVal {
			elem := convert(*t.Elem, val)
			data.Index(idx).Set(reflect.ValueOf(elem))
		}
		return data.Interface()

	case SliceTy:
		// todo: reflect
		var fmtVal []interface{}
		reflectInput := reflect.ValueOf(input)
		switch reflectInput.Kind() {
		case reflect.String:
			fmtVal = []interface{}{input}
		case reflect.Slice:
			inputLen := reflectInput.Len()
			fmtVal = make([]interface{}, inputLen)
			for i := 0; i < inputLen; i++ {
				fmtVal[i] = reflectInput.Index(i).Interface()
			}
		}

		data := reflect.MakeSlice(t.Type, len(fmtVal), len(fmtVal))
		for idx, val := range fmtVal {
			elem := convert(*t.Elem, val)
			data.Index(idx).Set(reflect.ValueOf(elem))
		}
		return data.Interface()

	case FixedBytesTy:
		if str, ok := input.(string); ok {
			return newFixedBytes(t.Size, str)
		}
	default:
		if str, ok := input.(string); ok {
			return newElement(t, str)
		}

	}
	return nil
}

// convert from string to basic type element
func newElement(t Type, val string) interface{} {
	if t.T == SliceTy || t.T == ArrayTy {
		return nil
	}
	var UNIT = 64
	var elem interface{}
	switch t.stringKind {
	case "uint8":
		num, _ := strconv.ParseUint(val, 10, UNIT)
		elem = uint8(num)
	case "uint16":
		num, _ := strconv.ParseUint(val, 10, UNIT)
		elem = uint16(num)
	case "uint32":
		num, _ := strconv.ParseUint(val, 10, UNIT)
		elem = uint32(num)
	case "uint64":
		num, _ := strconv.ParseUint(val, 10, UNIT)
		elem = uint64(num)
	case "int8":
		num, _ := strconv.ParseInt(val, 10, UNIT)
		elem = int8(num)
	case "int16":
		num, _ := strconv.ParseInt(val, 10, UNIT)
		elem = int16(num)
	case "int32":
		num, _ := strconv.ParseInt(val, 10, UNIT)
		elem = int32(num)
	case "int64":
		num, _ := strconv.ParseInt(val, 10, UNIT)
		elem = int64(num)
	case "bool":
		v, _ := strconv.ParseBool(val)
		elem = v
	case "address":
		elem = common.HexToAddress(val)
	case "string":
		elem = val
	case "bytes":
		elem = common.Hex2Bytes(val)
	default:
		if strings.HasPrefix(t.stringKind, "int") || strings.HasPrefix(t.stringKind, "uint") {
			var size int64
			if strings.HasPrefix(t.stringKind, "int") {
				size, _ = strconv.ParseInt(strings.TrimPrefix(t.stringKind, "int"), 10, UNIT)
			} else {
				size, _ = strconv.ParseInt(strings.TrimPrefix(t.stringKind, "uint"), 10, UNIT)
			}
			if size > 0 && size <= 256 && size%8 == 0 {
				var num *big.Int
				if val == "" {
					num = big.NewInt(0)
				} else {
					num, _ = big.NewInt(0).SetString(val, 10)
				}
				elem = num
			} else {
				elem = reflect.New(t.Type).Elem().Interface()
			}
		} else {
			// default use reflect but do not use val
			// because it's impossible to know how to convert from string to target type
			elem = reflect.New(t.Type).Elem().Interface()
		}
	}

	return elem
}

var byteTy = reflect.TypeOf(byte(0))

// the return val is a byte array, not slice
func newFixedBytes(size int, val string) interface{} {
	// pre-define size 1,2,3...32 and 64, other size use reflect
	switch size {
	case 1:
		var data [1]byte
		copy(data[:], []byte(val))
		return data
	case 2:
		var data [2]byte
		copy(data[:], []byte(val))
		return data
	case 3:
		var data [3]byte
		copy(data[:], []byte(val))
		return data
	case 4:
		var data [4]byte
		copy(data[:], []byte(val))
		return data
	case 5:
		var data [5]byte
		copy(data[:], []byte(val))
		return data
	case 6:
		var data [6]byte
		copy(data[:], []byte(val))
		return data
	case 7:
		var data [7]byte
		copy(data[:], []byte(val))
		return data
	case 8:
		var data [8]byte
		copy(data[:], []byte(val))
		return data
	case 9:
		var data [9]byte
		copy(data[:], []byte(val))
		return data
	case 10:
		var data [10]byte
		copy(data[:], []byte(val))
		return data
	case 11:
		var data [11]byte
		copy(data[:], []byte(val))
		return data
	case 12:
		var data [12]byte
		copy(data[:], []byte(val))
		return data
	case 13:
		var data [13]byte
		copy(data[:], []byte(val))
		return data
	case 14:
		var data [14]byte
		copy(data[:], []byte(val))
		return data
	case 15:
		var data [15]byte
		copy(data[:], []byte(val))
		return data
	case 16:
		var data [16]byte
		copy(data[:], []byte(val))
		return data
	case 17:
		var data [17]byte
		copy(data[:], []byte(val))
		return data
	case 18:
		var data [18]byte
		copy(data[:], []byte(val))
		return data
	case 19:
		var data [19]byte
		copy(data[:], []byte(val))
		return data
	case 20:
		var data [20]byte
		copy(data[:], []byte(val))
		return data
	case 21:
		var data [21]byte
		copy(data[:], []byte(val))
		return data
	case 22:
		var data [22]byte
		copy(data[:], []byte(val))
		return data
	case 23:
		var data [23]byte
		copy(data[:], []byte(val))
		return data
	case 24:
		var data [24]byte
		copy(data[:], []byte(val))
		return data
	case 25:
		var data [25]byte
		copy(data[:], []byte(val))
		return data
	case 26:
		var data [26]byte
		copy(data[:], []byte(val))
		return data
	case 27:
		var data [27]byte
		copy(data[:], []byte(val))
		return data
	case 28:
		var data [28]byte
		copy(data[:], []byte(val))
		return data
	case 29:
		var data [29]byte
		copy(data[:], []byte(val))
		return data
	case 30:
		var data [30]byte
		copy(data[:], []byte(val))
		return data
	case 31:
		var data [31]byte
		copy(data[:], []byte(val))
		return data
	case 32:
		var data [32]byte
		copy(data[:], []byte(val))
		return data
	case 64:
		var data [64]byte
		copy(data[:], []byte(val))
		return data
	default:
		return newFixedBytesWithReflect(size, val)
	}
}

//! NOTICE: newFixedBytesWithReflect take more 15 times of time than newFixedBytes
//! So it is just use for those fixed bytes which are not commonly used.
func newFixedBytesWithReflect(size int, val string) interface{} {
	data := reflect.New(reflect.ArrayOf(size, byteTy)).Elem()
	bytes := reflect.ValueOf([]byte(val))
	reflect.Copy(data, bytes)
	return data.Interface()
}
