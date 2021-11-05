package hvm

import (
	"errors"
	"github.com/meshplus/gosdk/common"
)

func GenPayload(beanAbi *BeanAbi, params ...interface{}) ([]byte, error) {
	switch beanAbi.BeanType {
	case MethodBean:
		return methodBeanPayload(beanAbi, params...)
	case InvokeBean:
		fallthrough
	default:
		return invokeBeanPayload(beanAbi, params...)
	}
}

// | class length(4B) | name length(2B) | class | class name | bin |
func invokeBeanPayload(beanAbi *BeanAbi, params ...interface{}) ([]byte, error) {
	classBytes := beanAbi.classBytes()

	if len(classBytes) > 0xffff {
		return nil, errors.New("the bean class is too large") // 64k
	}

	beanName := []byte(beanAbi.BeanName)
	isJson := true
	for _, str := range params {
		if _, ok := str.(string); !ok {
			isJson = false
			break
		}
	}
	var bin string
	var err error
	if isJson {
		bin, err = beanAbi.encodeJson(params...)
	} else {
		bin, err = beanAbi.encode(params...)
	}

	if err != nil {
		return nil, err
	}
	binBytes := []byte(bin)

	result := make([]byte, 0)
	classLenByte := common.IntToBytes4(len(classBytes))
	nameLenByte := common.IntToBytes2(len(beanName))
	result = append(result, classLenByte[:]...)
	result = append(result, nameLenByte[:]...)
	result = append(result, classBytes...)
	result = append(result, beanName...)
	result = append(result, binBytes...)

	return result, nil
}

func methodBeanPayload(methodAbi *BeanAbi, params ...interface{}) ([]byte, error) {
	if len(params) != len(methodAbi.Inputs) {
		return nil, errors.New("param count is not enough")
	}

	paramBuilder := NewParamBuilder(methodAbi.BeanName)

	for i, input := range methodAbi.Inputs {
		param := Convert(params[i])
		if input.Name == "" {
			return nil, errors.New("input StructName is empty")
		}
		switch input.EntryType {
		case Bool, Byte, Short, Int, Long, Float, Double, Char, String:
			if str, ok := param.(string); ok {
				paramBuilder.appendPayload([]byte(input.StructName), []byte(str))
			}
		case Struct, Array, List, Map:
			paramBuilder.AddObject(input.StructName, param)
		}
	}
	return paramBuilder.Build(), nil
}
