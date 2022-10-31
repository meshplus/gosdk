package scale

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/meshplus/gosdk/common"
	"io"
)

type Contract struct {
	Name        string          `json:"name"`
	Constructor ConstructorSpec `json:"constructor"`
}

type ConstructorSpec struct {
	Input []TypeInfo `json:"input"`
}

type TypeInfo struct {
	TypeId uint32 `json:"type_id"`
}

type InvokeBean struct {
	MethodName string    `json:"method_name"`
	Params     []Compact `json:"params"`
}

type Method struct {
	Name   string     `json:"name"`
	Input  []TypeInfo `json:"input"`
	Output []TypeInfo `json:"output"`
}

type Type struct {
	Id        uint32       `json:"id"`
	Type      string       `json:"type"`
	Fields    []TypeInfo   `json:"fields"`
	ArrayLen  int          `json:"array_len"`
	Primitive string       `json:"primitive"`
	Variants  [][]TypeInfo `json:"variants"`
}

type Abi struct {
	Contract    Contract `json:"contract"`
	Methods     []Method `json:"methods"`
	Types       []Type   `json:"types"`
	methodIndex map[string]int
}

func JSON(reader io.Reader) (Abi, error) {
	dec := json.NewDecoder(reader)

	var abi Abi
	if err := dec.Decode(&abi); err != nil {
		return Abi{}, err
	}
	abi.initMethodMap()
	return abi, nil
}

func (a *Abi) initMethodMap() {
	a.methodIndex = make(map[string]int)
	for i, v := range a.Methods {
		a.methodIndex[v.Name] = i
	}
}

func (a *Abi) checkType(tp Type, params interface{}) (bool, error) {
	switch changeStringToType(TypeString(tp.Type)) {
	case Vec:
		if val, ok := params.(*CompactVec); ok {
			nextType := a.Types[tp.Fields[0].TypeId]
			checkType, err := a.checkType(nextType, val.Val[0])
			if err != nil {
				return false, err
			}
			return checkType, nil
		} else {
			return false, errors.New("invalid type")
		}
	}
	return true, nil
}

func (a *Abi) getVec(tp Type, res []TypeString) []TypeString {
	switch changeStringToType(TypeString(tp.Type)) {
	case Vec:
		res = append(res, VecName)
		return a.getVec(a.Types[tp.Fields[0].TypeId], res)
	case Primitive:
		res = append(res, formatTypeString(TypeString(tp.Primitive)))
		return res
	case Array:
		res = append(res, ArrayName)
		return a.getVec(a.Types[tp.Fields[0].TypeId], res)
	case Struct:
		res = append(res, StructName)
	}

	return res
}

func (a *Abi) getArray(tp Type, res []TypeString, len int) (int, []TypeString) {
	switch changeStringToType(TypeString(tp.Type)) {
	case Vec:
		res = append(res, VecName)
		return a.getArray(a.Types[tp.Fields[0].TypeId], res, len)
	case Array:
		res = append(res, ArrayName)
		len = tp.ArrayLen
		return a.getArray(a.Types[tp.Fields[0].TypeId], res, len)
	case Primitive:
		res = append(res, TypeString(tp.Primitive))
		return len, res
	}

	return len, res
}

func (a *Abi) Encode(methodName string, params ...interface{}) ([]byte, error) {
	if methodName == "" {
		return a.encodeConstruct(params...)
	}
	return a.encodeMethod(methodName, params...)
}

func (a *Abi) encodeMethod(methodName string, params ...interface{}) ([]byte, error) {
	if methodId, ok := a.methodIndex[methodName]; ok {
		method := a.Methods[methodId]
		if len(method.Input) != len(params) {
			return nil, errors.New(fmt.Sprintf("need %d params, got %d", len(method.Input), len(params)))
		}
		var buf bytes.Buffer
		// encode method name
		name := &CompactString{Val: methodName}
		nameEncode, err := name.Encode()
		if err != nil {
			return nil, err
		}
		buf.Write(nameEncode)

		// encode input params
		for i := 0; i < len(method.Input); i++ {
			tp := method.Input[i]
			current, err := a.convert(a.Types[tp.TypeId], params[i])
			if err != nil {
				return nil, err
			}
			tmp, err := encode(current)
			if err != nil {
				return nil, err
			}
			buf.Write(tmp)
		}
		return buf.Bytes(), nil
	} else {
		return nil, errors.New("method not exist")
	}
}

func (a *Abi) encodeConstruct(params ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	// section type
	buf.WriteByte(0)

	sectionLen := len(CustomParamsSection) + 1

	var buf2 bytes.Buffer
	for i := 0; i < len(a.Contract.Constructor.Input); i++ {
		tp := a.Contract.Constructor.Input[i]
		current, err := a.convert(a.Types[tp.TypeId], params[i])
		if err != nil {
			return nil, err
		}
		tmp, err := encode(current)
		if err != nil {
			return nil, err
		}
		buf2.Write(tmp)
		sectionLen += len(tmp)
	}

	// section len
	buf.Write(common.EncodeInt32(int32(sectionLen)))

	// section name
	buf.WriteByte(byte(len(CustomParamsSection)))
	buf.Write([]byte(CustomParamsSection))
	// section content
	head := &CompactString{Val: "N"}
	headRaw, err := head.Encode()
	if err != nil {
		return nil, err
	}
	buf.Write(headRaw)
	buf.Write(buf2.Bytes())
	return buf.Bytes(), nil
}

func (a *Abi) encodeConstructCompact(params ...Compact) ([]byte, error) {
	var buf bytes.Buffer
	// section type
	buf.WriteByte(0)

	sectionLen := len(CustomParamsSection) + 1

	var buf2 bytes.Buffer
	for i := 0; i < len(a.Contract.Constructor.Input); i++ {
		tp := a.Contract.Constructor.Input[i]
		_, err := a.checkType(a.Types[tp.TypeId], params[i])
		if err != nil {
			return nil, err
		}

		tmp, err := encode(params[i])
		if err != nil {
			return nil, err
		}
		buf2.Write(tmp)
		sectionLen += len(tmp)
	}

	// section len
	buf.Write(common.EncodeInt32(int32(sectionLen)))

	// section name
	buf.WriteByte(byte(len(CustomParamsSection)))
	buf.Write([]byte(CustomParamsSection))

	// section content
	buf.Write(buf2.Bytes())
	return buf.Bytes(), nil
}

func (a *Abi) encodeMethodCompact(methodName string, params ...Compact) ([]byte, error) {
	if methodId, ok := a.methodIndex[methodName]; ok {
		method := a.Methods[methodId]
		if len(method.Input) != len(params) {
			return nil, errors.New(fmt.Sprintf("need %d params, got %d", len(method.Input), len(params)))
		}
		var buf bytes.Buffer
		// encode method name
		name := &CompactString{Val: methodName}
		nameEncode, err := name.Encode()
		if err != nil {
			return nil, err
		}
		buf.Write(nameEncode)

		// encode input params
		for i := 0; i < len(method.Input); i++ {
			tp := method.Input[i]
			_, err := a.checkType(a.Types[tp.TypeId], params[i])
			if err != nil {
				return nil, err
			}

			tmp, err := encode(params[i])
			if err != nil {
				return nil, err
			}
			buf.Write(tmp)
		}
		return buf.Bytes(), nil
	} else {
		return nil, errors.New("method not exist")
	}
}

func (a *Abi) EncodeCompact(methodName string, params ...Compact) ([]byte, error) {
	if methodName == "" {
		return a.encodeConstructCompact(params...)
	}
	return a.encodeMethodCompact(methodName, params...)
}

func (a *Abi) DecodeInput(methodName string, val []byte) (*InvokeBean, error) {
	if methodId, ok := a.methodIndex[methodName]; ok {
		result := &InvokeBean{}
		method := a.Methods[methodId]
		var params []Compact
		methodCompact := &CompactString{}
		offset, err := methodCompact.Decode(val)
		if err != nil {
			return nil, err
		}
		result.MethodName = methodCompact.Val
		if result.MethodName != methodName {
			return nil, errors.New("error method name")
		}
		val = val[offset:]
		for _, v := range method.Input {
			res, offset, err := a.decodeType(a.Types[v.TypeId], val)
			if err != nil {
				return nil, err
			}
			val = val[offset:]
			params = append(params, res)
		}
		result.Params = params
		return result, nil
	} else {
		return nil, errors.New("method not exist")
	}
}

func (a *Abi) DecodeRet(ret []byte, methodName string) (*InvokeBean, error) {
	if methodId, ok := a.methodIndex[methodName]; ok {
		result := &InvokeBean{}
		method := a.Methods[methodId]
		result.MethodName = methodName
		var params []Compact
		for _, v := range method.Output {
			res, offset, err := a.decodeType(a.Types[v.TypeId], ret)
			if err != nil {
				return nil, err
			}
			ret = ret[offset:]
			params = append(params, res)
		}
		result.Params = params
		return result, nil
	} else {
		return nil, errors.New("method not exist")
	}
}

func (a *Abi) convert(curType Type, param interface{}) (Compact, error) {
	var cur = curType.Type
	switch changeStringToType(TypeString(cur)) {
	case Primitive:
		return a.convertToPrimitive(curType, param)
	case Vec:
		return a.convertToVec(curType, param)
	case Struct:
		return a.convertToStruct(curType, param)
	case Array:
		return a.convertToArray(curType, param)
	case Tuple:
		return a.convertToTuple(curType, param)
	case Enum:
		return a.convertToEnum(curType, param)
	default:
		return nil, errors.New("unsupported type")
	}
}

func (a *Abi) convertToPrimitive(curType Type, param interface{}) (Compact, error) {
	cur := curType.Primitive
	temp, err := convertPrimitive(cur, param)
	if err != nil {
		return nil, err
	}
	return temp, nil
}

func (a *Abi) convertToVec(curType Type, param interface{}) (*CompactVec, error) {
	var depth []TypeString
	depth = a.getVec(a.Types[curType.Fields[0].TypeId], depth)
	cc := &CompactVec{
		Val:      nil,
		NextList: depth,
	}
	if val, ok := param.([]interface{}); ok {
		for _, k := range val {
			tmp, err := a.convert(a.Types[curType.Fields[0].TypeId], k)
			if err != nil {
				return nil, err
			}
			cc.Val = append(cc.Val, tmp)
		}
		return cc, nil
	} else {
		return nil, errors.New("param not vec")
	}
}

func (a *Abi) convertToStruct(curType Type, param interface{}) (*CompactStruct, error) {
	cc := &CompactStruct{}
	if val, ok := param.([]interface{}); ok {
		for i, v := range curType.Fields {
			tmp, err := a.convert(a.Types[v.TypeId], val[i])
			if err != nil {
				return nil, err
			}
			cc.Val = append(cc.Val, tmp)
		}
		return cc, nil
	} else {
		return nil, errors.New("param not struct json string")
	}
}

func (a *Abi) convertToArray(curType Type, param interface{}) (*CompactArray, error) {
	var depth []TypeString
	depth = a.getVec(curType, depth)
	cc := &CompactArray{
		Val:      nil,
		NextList: depth[1:],
		Len:      curType.ArrayLen,
	}
	if val, ok := param.([]interface{}); ok {
		if len(val) != cc.Len {
			return nil, errors.New("invalid array params length")
		}
		for _, k := range val {
			tmp, err := a.convert(a.Types[curType.Fields[0].TypeId], k)
			if err != nil {
				return nil, err
			}
			cc.Val = append(cc.Val, tmp)
		}
		return cc, nil
	} else {
		return nil, errors.New("param not array")
	}
}

func (a *Abi) convertToTuple(curType Type, param interface{}) (*CompactTuple, error) {
	cc := &CompactTuple{}
	if val, ok := param.([]interface{}); ok {
		for i, v := range curType.Fields {
			tmp, err := a.convert(a.Types[v.TypeId], val[i])
			if err != nil {
				return nil, err
			}
			cc.Val = append(cc.Val, tmp)
		}
		return cc, nil
	} else {
		return nil, errors.New("param not tuple")
	}
}

func (a *Abi) convertToEnum(curType Type, param interface{}) (*CompactEnum, error) {
	cc := &CompactEnum{}
	if val, ok := param.([]interface{}); ok {
		index, err := convertToUint8(val[0])
		if err != nil {
			return nil, err
		}
		cc.index = index.Val
		for i, v := range curType.Variants[cc.index] {
			tmp, err := a.convert(a.Types[v.TypeId], val[i+1])
			if err != nil {
				return nil, err
			}
			cc.Val = tmp
		}
		return cc, nil
	} else {
		return nil, errors.New("param not enum")
	}
}

func (a *Abi) GetMethod(name string) (*Method, error) {
	if methodId, ok := a.methodIndex[name]; ok {
		return &a.Methods[methodId], nil
	} else {
		return nil, errors.New("method not exist")
	}
}
