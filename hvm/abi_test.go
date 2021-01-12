package hvm

import (
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbi_GetMethodAbi(t *testing.T) {
	methodInvokeAbi := "../hvmtestfile/methodInvoke/hvm.abi"
	abiJson, err := common.ReadFileAsString(methodInvokeAbi)
	assert.Nil(t, err)

	abi, err := GenAbi(abiJson)
	assert.Nil(t, err)

	_, err = abi.GetMethodAbi("Hello")
	assert.Nil(t, err)

	methodAbi2, err := abi.GetMethodAbi("Hello()")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(methodAbi2.Inputs))

	methodAbi3, err := abi.GetMethodAbi("Hello(java.lang.String)")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(methodAbi3.Inputs))

	methodAbi4, err := abi.GetMethodAbi("Hello(int,java.lang.String)")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(methodAbi4.Inputs))
}
