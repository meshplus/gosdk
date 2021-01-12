package hvm

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

const (
	ObjINT    = "java.lang.Integer"
	INT       = "int"
	ObjSHORT  = "java.lang.Short"
	SHORT     = "short"
	ObjLONG   = "java.lang.Long"
	LONG      = "long"
	ObjBYTE   = "java.lang.Byte"
	BYTE      = "byte"
	ObjFLOAT  = "java.lang.Float"
	FLOAT     = "float"
	ObjDOUBLE = "java.lang.Double"
	DOUBLE    = "double"
	ObjCHAR   = "java.lang.Character"
	CHAR      = "char"
	ObjBOOL   = "java.lang.Boolean"
	BOOL      = "boolean"
	STRING    = "java.lang.String"
	OBJECT    = "java.lang.Object"
)

type ParamBuilder struct {
	buf bytes.Buffer
}

func NewParamBuilder(s string) *ParamBuilder {
	var p = new(ParamBuilder)
	methodName := []byte(s)
	p.buf.Write([]byte("fefffbce"))
	p.buf.Write(get2Length(len(methodName)))
	p.buf.Write(methodName)
	return p

}

func (p *ParamBuilder) CreateMethod(s string) *ParamBuilder {

	methodName := []byte(s)
	p.buf.Write(get2Length(len(methodName)))
	p.buf.Write(methodName)
	return p

}

func (p *ParamBuilder) AddInteger(s int) *ParamBuilder {
	clazzName := []byte(ObjINT)
	param := []byte(strconv.Itoa(s))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) Addint(s int) *ParamBuilder {
	clazzName := []byte(INT)
	param := []byte(strconv.Itoa(s))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) AddShort(s int16) *ParamBuilder {
	clazzName := []byte(ObjSHORT)
	param := []byte(strconv.Itoa(int(s)))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) Addshort(s int16) *ParamBuilder {
	clazzName := []byte(SHORT)
	param := []byte(strconv.Itoa(int(s)))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) AddLong(s int64) *ParamBuilder {
	clazzName := []byte(ObjLONG)
	param := []byte(strconv.FormatInt(s, 10))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) Addlong(s int64) *ParamBuilder {
	clazzName := []byte(LONG)
	param := []byte(strconv.FormatInt(s, 10))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) AddByte(s byte) *ParamBuilder {
	clazzName := []byte(ObjBYTE)
	param := []byte{s}
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) Addbyte(s byte) *ParamBuilder {
	clazzName := []byte(BYTE)
	param := []byte{s}
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) AddFloat(s float32) *ParamBuilder {
	clazzName := []byte(ObjFLOAT)
	param := []byte(fmt.Sprintf("%f", s))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) Addfloat(s float32) *ParamBuilder {
	clazzName := []byte(FLOAT)
	param := []byte(fmt.Sprintf("%f", s))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) AddDouble(s float64) *ParamBuilder {
	clazzName := []byte(ObjDOUBLE)
	param := []byte(fmt.Sprintf("%f", s))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) Adddouble(s float64) *ParamBuilder {
	clazzName := []byte(DOUBLE)
	param := []byte(fmt.Sprintf("%f", s))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) AddCharacter(s uint16) *ParamBuilder {
	clazzName := []byte(ObjCHAR)
	param := []byte(string(rune(s)))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) Addchar(s uint16) *ParamBuilder {
	clazzName := []byte(CHAR)
	param := []byte(string(rune(s)))
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) AddBoolean(s bool) *ParamBuilder {
	clazzName := []byte(ObjBOOL)
	param := []byte(strconv.FormatBool(s))
	p.appendPayload(clazzName, param)
	return p
}

func (p *ParamBuilder) Addbool(s bool) *ParamBuilder {
	clazzName := []byte(BOOL)
	param := []byte(strconv.FormatBool(s))
	p.appendPayload(clazzName, param)
	return p
}

func (p *ParamBuilder) AddString(s string) *ParamBuilder {
	clazzName := []byte(STRING)
	param := []byte(s)
	p.appendPayload(clazzName, param)
	return p

}

func (p *ParamBuilder) AddObject(clazz string, s interface{}) *ParamBuilder {
	clazzName := []byte(clazz)
	var param []byte
	if reflect.TypeOf(s).Kind() == reflect.String {
		param = []byte(s.(string))
	} else {
		param, _ = json.Marshal(s)
	}
	p.appendPayload(clazzName, param)
	return p
}

func (p *ParamBuilder) Build() []byte {

	return p.buf.Bytes()
}

func (p *ParamBuilder) appendPayload(clazzName []byte, param []byte) {
	p.buf.Write(get2Length(len(clazzName)))
	p.buf.Write(get4Length(len(param)))
	p.buf.Write(clazzName)
	p.buf.Write(param)

}

func get2Length(length int) []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(length))
	return bs

}

func get4Length(length int) []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(length))
	return bs
}
