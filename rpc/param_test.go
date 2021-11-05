package rpc

import (
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerlizerToString(t *testing.T) {
	qtr := &QueryTxRange{
		From: 1,
		To:   2,
	}
	param := qtr.SerializeToString()
	assert.Equal(t, param, "")
}

func TestNnewMapParam(t *testing.T) {
	mapParam := newMapParam("1", 2)
	ans := mapParam.SerializeToString()
	assert.Equal(t, ans, "{\"1\":2}")

}

func TestNewLogsFilter(t *testing.T) {
	lf := NewLogsFilter()
	NewBlockEventFilter().SetBlockInfo(true)
	lf.SetFromBlock(1)
	lf.SetToBlock(2)
	lf.AddAddress("efefw")
	lf.SetTopic(1, *new(common.Hash))
	lf.GetEventType()
	lf.Serialize()
}

func TestAarrayToString(t *testing.T) {
	arrayToString([]string{"1", "1"})
	arrayToString([]int{1, 1})
	arrayToString([]common.Address{common.BytesToAddress([]byte{1}), common.BytesToAddress([]byte{1})})
	arrayToString([]float32{1.1})
	arrayToString([]common.Hash{common.BytesToHash([]byte{1}), common.BytesToHash([]byte{1})})
}

func TestNewSystemStatusFilter(t *testing.T) {
	ssf := NewSystemStatusFilter()
	ssf.AddModules("module1")
	ssf.AddModulesExclude("me1")
	ssf.AddSubtypes("subtype1")
	ssf.AddSubtypesExclude("AddSubtypesExclude1")
	ssf.AddErrorCode("AddErrorCode1")
	ssf.AddErrorCodeExclude("AddErrorCodeExclude")
	ssf.GetEventType()
	ssf.Serialize()

}
