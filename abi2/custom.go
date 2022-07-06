package abi2

import (
	"bytes"
	"fmt"
	"github.com/meshplus/gosdk/common"
	"reflect"
)

// UnpackLog log output in v according to the abi specification
// param data means the Data field of TxLog, params topics means the Topics field of TxLog
//
// warn: try to decode a dynamic type indexed event param(like: bytes in solidity) will cause panic
func (abi *ABI) UnpackLog(v interface{}, name string, data string, topics []string) (err error) {
	if len(data) > 0 {
		if err := abi.UnpackIntoInterface(v, name, common.Hex2Bytes(data)); err != nil {
			return err
		}
	}
	var indexed Arguments
	for _, arg := range abi.Events[name].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	var hashTopic []common.Hash
	for _, arg := range topics {
		hashTopic = append(hashTopic, common.HexToHash(arg))
	}
	return ParseTopics(v, indexed, hashTopic[1:])
}

// UnpackResult use abi to decode data from solidity return
func (abi *ABI) UnpackResult(v interface{}, name, data string) (err error) {
	return abi.UnpackIntoInterface(v, name, common.FromHex(data))
}

// ByteArrayToString transfer a byte array to string
func ByteArrayToString(v interface{}) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Array {
		return "", fmt.Errorf("only support byte array")
	}
	if val.Len() == 0 {
		return "", nil
	}
	if val.Index(0).Kind() != reflect.Uint8 {
		return "", fmt.Errorf("only support byte array")
	}

	var (
		buffer   bytes.Buffer
		endIndex int
	)

	for i := val.Len() - 1; i >= 0; i-- {
		if b, _ := val.Index(i).Interface().(byte); b != 0 {
			endIndex = i
			break
		}
	}

	for i := 0; i <= endIndex; i++ {
		b, _ := val.Index(i).Interface().(byte)
		err := buffer.WriteByte(b)
		if err != nil {
			return "", nil
		}
	}

	return buffer.String(), nil
}
