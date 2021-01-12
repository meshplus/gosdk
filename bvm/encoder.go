package bvm

import (
	"bytes"
	"github.com/meshplus/gosdk/common"
)

// EncodeOperation encode operation to payload for bvm
func EncodeOperation(ope Operation) []byte {
	buffer := bytes.NewBuffer([]byte{})
	methodNameLenBytes := common.IntToBytes4(len(ope.Method()))
	paramsLenBytes := common.IntToBytes4(len(ope.Args()))
	buffer.Write(methodNameLenBytes[:])
	buffer.WriteString(string(ope.Method()))
	buffer.Write(paramsLenBytes[:])
	for _, param := range ope.Args() {
		paramLenBytes := common.IntToBytes4(len(param))
		buffer.Write(paramLenBytes[:])
		buffer.WriteString(param)
	}
	return buffer.Bytes()
}

// encodeOperations encode the slice of Operation for create proposal
func encodeProposalContentOperation(ops []ProposalContentOperation) []byte {
	buffer := bytes.NewBuffer([]byte{})
	operationLenBytes := common.IntToBytes4(len(ops))
	buffer.Write(operationLenBytes[:])
	for _, pa := range ops {
		result := EncodeOperation(pa)
		buffer.Write(result)
	}
	return buffer.Bytes()
}
