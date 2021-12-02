package hvm

import (
	"errors"
	"github.com/meshplus/gosdk/common"
	"strings"
)

// abi version
type Version string

const (
	Version1 Version = "v1"
)

// abi type
type Type string

// bean type
type BeanType string

const (
	Void   Type = "Void"
	Bool   Type = "Bool"
	Char   Type = "Char"
	Byte   Type = "Byte"
	Short  Type = "Short"
	Int    Type = "Int"
	Long   Type = "Long"
	Float  Type = "Float"
	Double Type = "Double"
	String Type = "String"
	Array  Type = "Array"
	List   Type = "List"
	Map    Type = "Map"
	Struct Type = "Struct"
)

const (
	InvokeBean BeanType = "InvokeBean"
	MethodBean BeanType = "MethodBean"
)

type Entry struct {
	Name       string  `json:"name"`
	EntryType  Type    `json:"type"`
	Properties []Entry `json:"properties,omitempty"`
	StructName string  `json:"structName,omitempty"`
}

type BeanAbi struct {
	BeanVersion Version  `json:"version"`
	BeanName    string   `json:"beanName"`
	Inputs      []Entry  `json:"inputs"`
	Output      Entry    `json:"output"`
	ClassBytes  string   `json:"classBytes"`
	Structs     []Entry  `json:"structs"`
	BeanType    BeanType `json:"beanType"`
}

type Abi []BeanAbi

func (abi Abi) GetBeanAbi(beanName string) (*BeanAbi, error) {
	for _, beanAbi := range abi {
		if beanAbi.BeanName == beanName {
			return &beanAbi, nil
		}
	}
	return nil, errors.New("can not find invoke bean " + beanName)
}

func (abi Abi) GetMethodAbi(methodName string) (*BeanAbi, error) {
	index := strings.IndexByte(methodName, '(')
	if index == -1 {
		for _, beanAbi := range abi {
			if beanAbi.BeanName == methodName && beanAbi.BeanType == MethodBean {
				return &beanAbi, nil
			}
		}
	} else {
		if !strings.HasSuffix(methodName, ")") {
			return nil, errors.New("methodName is not legal")
		}
		name := methodName[0:index]
		paramStr := methodName[index+1 : len(methodName)-1]
		params := make([]string, 0)
		if len(paramStr) != 0 {
			params = strings.Split(paramStr, ",")
		}
		for _, beanAbi := range abi {
			if beanAbi.BeanName == name && beanAbi.BeanType == MethodBean && beanAbi.checkMethodInputs(params) {
				return &beanAbi, nil
			}
		}
	}
	return nil, errors.New("can not find method bean " + methodName)
}

func (beanAbi *BeanAbi) resolveStruct(input Entry, param interface{}) string {
	result := "{"
	s, err := beanAbi.getStruct(input.StructName)
	if err != nil {
		return ""
	}
	properties := s.Properties
	realParams := param.([]interface{})
	for i, prop := range properties {
		p := realParams[i]
		result += beanAbi.resolveEntry(prop, p)

		if i < len(properties)-1 {
			result += ","
		}
	}

	result += "}"
	return result
}

func (beanAbi *BeanAbi) resolveList(input Entry, param interface{}) string {
	result := "["
	entry := input.Properties[0]
	//entryType := entry.EntryType
	realParams := param.([]interface{})
	for i, p := range realParams {
		result += beanAbi.resolveNestListOrMap(entry, p)

		if i < len(realParams)-1 {
			result += ","
		}
	}

	result += "]"
	return result
}

func (beanAbi *BeanAbi) resolveMap(input Entry, param interface{}) string {
	result := "{"
	keyEntry := input.Properties[0]
	valEntry := input.Properties[1]

	realParams := param.([]interface{})
	for i, p := range realParams {

		kv := p.([]interface{})
		result += beanAbi.resolveEntrySingle(keyEntry, kv[0])
		result += ":"
		result += beanAbi.resolveNestListOrMap(valEntry, kv[1])

		if i < len(realParams)-1 {
			result += ","
		}
	}

	result += "}"
	return result
}

func (beanAbi *BeanAbi) resolveNestListOrMap(valEntry Entry, param interface{}) (result string) {
	switch valEntry.EntryType {
	case Bool, Byte, Short, Int, Long, Float, Double:
		result += beanAbi.resolveEntrySingle(valEntry, param)
	case Char, String:
		result += beanAbi.resolveEntrySingle(valEntry, param)
	case Struct:
		result += beanAbi.resolveStruct(valEntry, param)
	case List:
		result += (beanAbi.resolveList(valEntry.Properties[0], param))
	case Map:
		result += (beanAbi.resolveMap(valEntry.Properties[0], param))
	case Array:
		result += (beanAbi.resolveArray(valEntry.Properties[0], param))

	}
	return
}

func (beanAbi *BeanAbi) resolveArray(input Entry, param interface{}) string {
	result := "["
	entry := input.Properties[0]
	//entryType := entry.EntryType
	realParams := param.([]interface{})
	for i, p := range realParams {
		result += beanAbi.resolveEntrySingle(entry, p)

		if i < len(realParams)-1 {
			result += ","
		}
	}

	result += "]"
	return result
}

func (beanAbi *BeanAbi) resolveEntrySingle(entry Entry, param interface{}) (result string) {
	switch entry.EntryType {
	case Bool, Byte, Short, Int, Long, Float, Double:
		if str, ok := param.(string); ok {
			result = str
		}
	case Char, String:
		if str, ok := param.(string); ok {
			result = "\"" + str + "\""
		}
		// TODO list<map<k, v>>, map<k, list<T>>
	case Struct:
		result = beanAbi.resolveStruct(entry, param)
	}
	return
}

func (beanAbi *BeanAbi) resolveEntry(input Entry, param interface{}) (entryResult string) {
	switch input.EntryType {
	case Bool, Byte, Short, Int, Long, Float, Double:
		if str, ok := param.(string); ok {
			entryResult = "\"" + input.Name + "\":" + str
		}
	case Char, String:
		if str, ok := param.(string); ok {
			entryResult = "\"" + input.Name + "\":\"" + str + "\""
		}
	case List:
		entryResult = "\"" + input.Name + "\":" + beanAbi.resolveList(input, param)
	case Map:
		entryResult = "\"" + input.Name + "\":" + beanAbi.resolveMap(input, param)
	case Array:
		entryResult = "\"" + input.Name + "\":" + beanAbi.resolveArray(input, param)
	case Struct:
		entryResult = "\"" + input.Name + "\":" + beanAbi.resolveStruct(input, param)
	}
	return
}
func (beanAbi *BeanAbi) resolveJsonEntry(input Entry, param interface{}) (entryResult string) {
	switch input.EntryType {

	case Char, String:
		if str, ok := param.(string); ok {
			entryResult = "\"" + input.Name + "\":\"" + strings.ReplaceAll(str, `"`, `\"`) + "\""
		}

	default:
		if str, ok := param.(string); ok {
			entryResult = "\"" + input.Name + "\":" + str
		}
	}
	return
}

func (beanAbi BeanAbi) encode(params ...interface{}) (string, error) {
	if len(beanAbi.Inputs) != len(params) {
		return "", errors.New("param count is not enough")
	}
	result := "{"
	for i, input := range beanAbi.Inputs {
		param := Convert(params[i])
		entryResult := beanAbi.resolveEntry(input, param)
		result += entryResult
		if i < len(params)-1 {
			result += ","
		}
	}
	result += "}"

	return result, nil
}

func (beanAbi BeanAbi) encodeJson(params ...interface{}) (string, error) {
	if len(beanAbi.Inputs) != len(params) {
		return "", errors.New("param count is not enough")
	}
	result := "{"
	for i, input := range beanAbi.Inputs {
		param := params[i]

		entryResult := beanAbi.resolveJsonEntry(input, param)
		result += entryResult
		if i < len(params)-1 {
			result += ","
		}
	}
	result += "}"
	return result, nil

}

func (beanAbi BeanAbi) classBytes() []byte {
	return common.Hex2Bytes(beanAbi.ClassBytes)
}

func (beanAbi BeanAbi) getStruct(name string) (*Entry, error) {
	for _, s := range beanAbi.Structs {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, errors.New("can not find struct " + name)
}

func (beanAbi BeanAbi) checkMethodInputs(params []string) bool {
	if len(beanAbi.Inputs) != len(params) || beanAbi.BeanType != MethodBean {
		return false
	}
	for i, input := range beanAbi.Inputs {
		if input.StructName != params[i] {
			return false
		}
	}
	return true
}
