package hvm

import (
	"bytes"
	"encoding/json"
	"github.com/meshplus/gosdk/classfile"
	"github.com/meshplus/gosdk/common"
	"strings"
)

type invokeArg struct {
	ParamName  string      `json:"name,omitempty"`
	ParamType  interface{} `json:"type"`
	ParamValue interface{} `json:"value"`
}

type PayLoad struct {
	InvokeBeanName string   `json:"invokeBeanName"`
	InvokeArgs     string   `json:"invokeArgs"`
	InvokeMethods  []string `json:"invokeMethods"`
}

func DecodePayload(pl string) (*PayLoad, error) {
	if strings.HasPrefix(pl, "0xfefffbce") || strings.HasPrefix(pl, "fefffbce") {
		return decodePayloadInvokeDirectly(pl)
	}
	return decodePayloadInvoke(pl)
}

func decodePayloadInvokeDirectly(pl string) (*PayLoad, error) {
	s := common.Hex2Bytes(pl)
	nameLen := common.BytesToInt32(s[4:6])
	name := s[6 : nameLen+6]
	begin := 6 + nameLen
	var args []*invokeArg
	for begin < len(s) {
		paramTypeLen := common.BytesToInt32(s[begin : begin+2])
		paramLen := common.BytesToInt32(s[begin+2 : begin+6])
		paramType := string(s[begin+6 : begin+6+paramTypeLen])
		param := string(s[begin+6+paramTypeLen : begin+6+paramTypeLen+paramLen])
		begin = begin + 6 + paramLen + paramTypeLen
		args = append(args, &invokeArg{
			ParamType:  paramType,
			ParamValue: param,
		})
	}
	methodName := string(name)
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(args)
	if err != nil {
		return nil, err
	}
	return &PayLoad{
		InvokeBeanName: "",
		InvokeArgs:     bf.String(),
		InvokeMethods:  []string{methodName},
	}, nil
}

func decodePayloadInvoke(pl string) (*PayLoad, error) {
	s := common.Hex2Bytes(pl)
	sLen := len(s)
	classLen := common.BytesToInt32(s[0:4])
	nameLen := common.BytesToInt32(s[4:6])
	classBytes := s[6 : classLen+6]
	cf, err := classfile.Parse(classBytes)
	if err != nil {
		return nil, err
	}
	var invokeMethods []string
	for _, m := range cf.Methods {
		desc := strings.Split(m.Descriptor(), ";)")
		if strings.ToLower(m.Name()) == "invoke" && len(desc) >= 2 && strings.Index(desc[1], "Object") == -1 {
			methodArgument := m.ArgumentTypes()
			invokeMethods = m.GetInvokeMethods(methodArgument)
		}
	}
	name := s[6+classLen : 6+classLen+nameLen]
	bin := s[6+classLen+nameLen : sLen]
	args := string(bin)
	return &PayLoad{
		InvokeBeanName: string(name),
		InvokeArgs:     args,
		InvokeMethods:  invokeMethods,
	}, nil
}
