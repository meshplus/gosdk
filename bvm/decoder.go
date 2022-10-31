package bvm

import (
	"encoding/json"
	"fmt"
	"github.com/meshplus/gosdk/common"
)

// Result represent the execute result of BVM
type Result struct {
	Success bool
	Ret     []byte
	Err     string
}

func (r *Result) String() string {
	return fmt.Sprintf(`{"Success":%v, "Ret":%s, "Err":%s}`, r.Success, string(r.Ret), r.Err)
}

// Decode decode ret to result
func Decode(ret string) *Result {
	if len(ret) < 2 {
		return &Result{}
	}
	var result Result
	_ = json.Unmarshal(common.FromHex(ret), &result)
	return &result
}

type ProposalCode struct {
	MethodName string
	Params     []string
}

// DecodeProposalCode decode proposal.code
// Example:
// 	a := common.Hex2Bytes(proposalHexString)
//	var pro ProposalData
//	proto.Unmarshal(a, &pro)
//	code, err := DecodeProposalCode(pro.Code)
func DecodeProposalCode(code []byte) (string, error) {
	const defaultLen = 4
	operationLen := common.BytesToInt32(code[0:defaultLen])
	index := defaultLen
	var result []ProposalCode
	for i := 0; i < operationLen; i++ {
		nameLen := common.BytesToInt32(code[index : index+defaultLen])
		index += defaultLen
		methodName := string(code[index : index+nameLen])
		index += nameLen
		paramCount := common.BytesToInt32(code[index : index+defaultLen])
		index += defaultLen
		var params []string
		for j := 0; j < paramCount; j++ {
			paramLen := common.BytesToInt32(code[index : index+defaultLen])
			index += defaultLen
			param := string(code[index : index+paramLen])
			index += paramLen
			params = append(params, param)
		}

		result = append(result, ProposalCode{
			MethodName: methodName,
			Params:     params,
		})
	}

	res, err := json.Marshal(result)
	return string(res), err
}

// DecodePayload
// payload struct as below
// methodName length(4 bytes) | methodName | params count(4 bytes) | (paramX length(4 bytes) | paramX)...
func DecodePayload(payloadBytes []byte) (Operation, error) {
	const defaultLen = 4
	nameLen := common.BytesToInt32(payloadBytes[0:defaultLen])
	index := defaultLen
	methodName := string(payloadBytes[index : index+nameLen])
	index += nameLen
	paramsCount := common.BytesToInt32(payloadBytes[index : index+defaultLen])
	index += defaultLen
	var params []string
	for j := 0; j < paramsCount; j++ {
		paramLen := common.BytesToInt32(payloadBytes[index : index+defaultLen])
		index += defaultLen
		param := string(payloadBytes[index : index+paramLen])
		index += paramLen
		params = append(params, param)
	}
	s := NewOperation()
	s.SetMethod(ContractMethod(methodName))
	s.SetArgs(params)
	return s, nil
}
