package abi

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var abi = `[{"constant":false,"inputs":[{"name":"num1","type":"uint32"},{"name":"num2","type":"uint32"}],"name":"add","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"archiveSum","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"getSum","outputs":[{"name":"","type":"uint32"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"increment","outputs":[],"payable":false,"type":"function"}]`

func TestABI_UnpackInputWithoutMethod(t *testing.T) {
	ABI, _ := JSON(strings.NewReader(abi))
	data, _ := ABI.Pack("add", uint32(1), uint32(2))
	var r1 uint32
	var r2 uint32
	testV := []interface{}{&r1, &r2}
	err := ABI.UnpackInputWithoutMethod(&testV, data)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, uint32(1), r1)
	assert.Equal(t, uint32(2), r2)
}

func TestABI_UnpackInput(t *testing.T) {
	ABI, _ := JSON(strings.NewReader(abi))
	data, _ := ABI.Pack("add", uint32(1), uint32(2))
	var r1 uint32
	var r2 uint32
	testV := []interface{}{&r1, &r2}
	err := ABI.UnpackInput(&testV, "add", data)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, uint32(1), r1)
	assert.Equal(t, uint32(2), r2)
}
