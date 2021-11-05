package abi

import (
	"bytes"
	"fmt"
	"github.com/meshplus/gosdk/common"
	"reflect"
)

// Unpack log output in v according to the abi specification
// param data means the Data field of TxLog, params topics means the Topics field of TxLog
//
// warn: try to decode a dynamic type indexed event param(like: bytes in solidity) will cause panic
func (abi ABI) UnpackLog(v interface{}, name string, data string, topics []string) (err error) {
	if err = abi.Unpack(v, name, common.Hex2Bytes(data)); err != nil {
		return
	}
	event, err := abi.GetEvent(name)
	if err != nil {
		return err
	}
	if len(event.Inputs.Indexed()) == 0 {
		return nil
	}
	if err = abi.unPackLogTopics(v, name, topics); err != nil {
		return
	}
	return
}

func (abi ABI) unPackLogTopics(v interface{}, name string, topics []string) (err error) {
	// make sure the passed value is arguments pointer
	if reflect.Ptr != reflect.ValueOf(v).Kind() {
		return fmt.Errorf("abi: Unpack(non-pointer %T)", v)
	}

	var data []byte
	if len(topics) > 1 {
		for i := 1; i < len(topics); i++ {
			data = append(data, common.Hex2Bytes(topics[i][2:])...)
		}
	}

	event, err := abi.GetEvent(name)
	if err != nil {
		return err
	}
	args := event.Inputs

	if event, err := abi.GetEvent(name); err == nil {
		indexeds := event.Inputs.Indexed()
		retval := make([]interface{}, 0, len(indexeds))
		virtualArgs := 0
		for index, arg := range indexeds {
			marshalledValue, err := toGoType((index+virtualArgs)*32, arg.Type, data)
			if arg.Type.T == ArrayTy {
				virtualArgs += getArraySize(&arg.Type) - 1
			}
			if err != nil {
				return err
			}
			retval = append(retval, marshalledValue)
		}
		// Tuple
		if indexeds.isTuple() {
			var (
				value = reflect.ValueOf(v).Elem()
				typ   = value.Type()
				kind  = value.Kind()
			)

			if err := requireUnpackKind(value, typ, kind, indexeds); err != nil {
				return err
			}

			// If the interface is a struct, get of abi->struct_field mapping

			var abi2struct map[string]string
			if kind == reflect.Struct {
				var err error
				abi2struct, err = mapAbiToStructFields(args, value)
				if err != nil {
					return err
				}
			}
			for i, arg := range indexeds {

				reflectValue := reflect.ValueOf(retval[i])

				switch kind {
				case reflect.Struct:
					if structField, ok := abi2struct[arg.Name]; ok {
						if err := set(value.FieldByName(structField), reflectValue, arg); err != nil {
							return err
						}
					}
				case reflect.Slice, reflect.Array:
					if value.Len() < i {
						return fmt.Errorf("abi: insufficient number of arguments for unpack, want %d, got %d", len(indexeds), value.Len())
					}
					v := value.Index(i)
					if err := requireAssignable(v, reflectValue); err != nil {
						return err
					}

					if err := set(v.Elem(), reflectValue, arg); err != nil {
						return err
					}
				default:
					return fmt.Errorf("abi:[2] cannot unmarshal tuple in to %v", typ)
				}
			}
			return nil
		} else {
			// Atomic
			if len(retval) != 1 {
				return fmt.Errorf("abi: wrong length, expected single value, got %d", len(retval))
			}

			elem := reflect.ValueOf(v).Elem()
			kind := elem.Kind()
			reflectValue := reflect.ValueOf(retval[0])

			var abi2struct map[string]string
			if kind == reflect.Struct {
				var err error
				if abi2struct, err = mapAbiToStructFields(args, elem); err != nil {
					return err
				}
				arg := indexeds[0]
				if structField, ok := abi2struct[arg.Name]; ok {
					return set(elem.FieldByName(structField), reflectValue, arg)
				}
				return nil
			}

			return set(elem, reflectValue, indexeds[0])
		}
	}
	return fmt.Errorf("abi: could not locate named event")
}

// NonIndexed returns the arguments with indexed arguments filtered out
func (arguments Arguments) Indexed() Arguments {
	var ret []Argument
	for _, arg := range arguments {
		if arg.Indexed {
			ret = append(ret, arg)
		}
	}
	return ret
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

// UnpackResult use abi to decode data from solidity return
func (abi *ABI) UnpackResult(v interface{}, name, data string) (err error) {
	if err := abi.Unpack(v, name, common.FromHex(data)); err != nil {
		return err
	}
	return nil
}

// UnpackInput unpack method input
func (abi *ABI) UnpackInput(v interface{}, methodName string, data []byte) (err error) {
	if len(data) == 0 {
		return fmt.Errorf("abi: unmarshalling empty output")
	}
	if method, err := abi.GetMethod(methodName); err == nil {
		if len(data[4:])%32 != 0 {
			return fmt.Errorf("abi: improperly formatted output")
		}
		return method.Inputs.Unpack(v, data[4:])
	} else if abi.Constructor.Name == methodName {
		return abi.Constructor.Inputs.Unpack(v, data)
	}
	return fmt.Errorf("abi: could not locate named method")
}

// UnpackInputWithoutMethod unpack input
func (abi *ABI) UnpackInputWithoutMethod(v interface{}, data []byte) (err error) {
	if len(data) == 0 {
		return fmt.Errorf("abi: unmarshalling empty output")
	}
	methodId := data[:4]
	method, err := abi.MethodById(methodId)
	// constructor
	if err != nil && method == nil {
		if err = abi.Constructor.Inputs.Unpack(v, data); err != nil {
			return err
		}
	} else { // method
		if err = method.Inputs.Unpack(v, data[4:]); err != nil {
			return err
		}
	}
	return nil
}
