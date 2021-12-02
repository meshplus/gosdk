package hvm

import (
	"fmt"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/rpc"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRPC_Encoder(t *testing.T) {

	logger := common.GetLogger("main")
	accountJson, sysErr := account.NewAccount("12345678")
	if sysErr != nil {
		logger.Error(sysErr)
		return

	}

	logger.Debugf(accountJson)
	key, sysErr := account.GenKeyFromAccountJson(accountJson, "12345678")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	hrpc := rpc.NewRPCWithPath("../conf")
	jarPath := "../hvmtestfile/hvm-first-1.0.jar"
	//jarPath = "../hvmtestfile/hvm-first.jar"
	payload, sysErr := rpc.DecompressFromJar(jarPath)
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	ecKey := key.(*account.ECDSAKey)

	transaction := rpc.NewTransaction(ecKey.GetAddress().Hex()).Deploy(common.Bytes2Hex(payload)).VMType(rpc.HVM)
	transaction.Sign(key)
	receipt, sysErr := hrpc.DeployContract(transaction)
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	// contractAddr := receipt.ContractAddress
	//read abi
	abiPath := "../hvmtestfile/hvm1.abi"
	abiJson, sysErr := common.ReadFileAsString(abiPath)
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	abi, sysErr := GenAbi(abiJson)
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	// fmt.Println(abi)
	easyBean := "cn.test.contract.invoke.EasyInvoke"
	beanAbi, sysErr := abi.GetBeanAbi(easyBean)
	// fmt.Println(beanAbi)
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}

	person1 := []interface{}{"tom", "21"}
	person2 := []interface{}{"jack", "18"}
	bean1 := []interface{}{"bean1", person1}
	bean2 := []interface{}{"bean2", person2}
	fmt.Println(reflect.TypeOf(person1))
	fmt.Println(reflect.TypeOf(bean1))
	fmap := make(map[int]string)
	fmt.Println(reflect.TypeOf(fmap))
	fmt.Println([]interface{}{[]interface{}{"789", []interface{}{[]interface{}{"456", "12.2"}}}, []interface{}{"234", []interface{}{[]interface{}{"345", "12.2"}}}})
	fmt.Println(len(beanAbi.Inputs))
	//fcarray := [3]string{"1", "2", "3"}
	//fclist := []string{"1f", "2f", "3f"}
	fcmap := make(map[string]map[string]string)
	person3 := []string{"tom", "123"}
	fcmap["789"] = make(map[string]string)
	fcmap["234"] = make(map[string]string)
	//father := make(map[string][]string[])
	fcmap["789"]["456"] = "12.2"
	fcmap["234"]["345"] = "12.2"
	fcmaptruct := make(map[string][]string)
	fcmaptruct["person1"] = []string{"tom", "123"}
	fcmaptruct["person2"] = []string{"fuc", "234"}

	//fclist1 := [][]string{{"1", "2"}, {"2", "4"}}
	//nolint
	invokePayload, gerr := GenPayload(beanAbi, "true", "c", "20", "100", "1000", "10000", "1.1", "1.11", "string",
		//`["1f","2f","3f"]`, `[1,2,3]`,
		//`{789:{456:12.2},234:{345:12.2}}`, `[[1,2],[2,4]]`,
		`{"name":"tom","age":21}`, `{"beanName":"bean1","person":{"name":"tom","age":21}}`,
		`["strList1","strList2"]`,
		`[{"name":"tom","age":21},{"name":"jack","age":18}]`,
		`{"person1":{"name":"tom","age":21},"person2":{"name":"jack","age":18}}`,
		`{"bean1":{"beanName":"bean1","person":{"name":"tom","age":21}},"bean2":{"beanName":"bean2","person":{"name":"jack","age":18}}}`)

	if gerr != nil {
		logger.Error(sysErr)
		return
	}
	invokePayload, sysErr = GenPayload(beanAbi, Convert("true"), Convert("c"), Convert("20"), Convert("20"), Convert("20"), Convert("20"), Convert("20.1"), Convert("20.2"), Convert("asasd"),
		//Convert(fclist), Convert(fcarray),
		//Convert(fcmap),
		//Convert(fclist1),
		Convert(person3), bean1,
		[]interface{}{"strList1", "strList2"},
		[]interface{}{person1, person2},
		Convert(fcmaptruct),
		[]interface{}{[]interface{}{"bean1", bean1}, []interface{}{"bean2", bean2}})
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	invokeTx := rpc.NewTransaction(ecKey.GetAddress().Hex()).Invoke(receipt.ContractAddress, invokePayload).VMType(rpc.HVM)
	invokeTx.Sign(key)
	_, sysErr = hrpc.InvokeContract(invokeTx)
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
}

func TestRPC_MethodAbi(t *testing.T) {
	t.Skip()
	logger := common.GetLogger("main")
	methodInvokeJar := "../hvmtestfile/methodInvoke/share-2.0.jar"
	methodInvokeAbi := "../hvmtestfile/methodInvoke/hvm1.abi"
	hrpc := rpc.NewRPCWithPath("../conf")
	accountJson, err := account.NewAccount("12345678")
	assert.Nil(t, err)

	logger.Debugf(accountJson)
	key, err := account.GenKeyFromAccountJson(accountJson, "12345678")
	assert.Nil(t, err)
	ecKey := key.(*account.ECDSAKey)

	payload, err := rpc.DecompressFromJar(methodInvokeJar)
	assert.Nil(t, err)

	transaction := rpc.NewTransaction(ecKey.GetAddress().Hex()).Deploy(common.Bytes2Hex(payload)).VMType(rpc.HVM)
	transaction.Sign(key)
	receipt, err := hrpc.DeployContract(transaction)
	assert.Nil(t, err)

	contractAddress := receipt.ContractAddress

	abiJson, err := common.ReadFileAsString(methodInvokeAbi)
	assert.Nil(t, err)

	abi, err := GenAbi(abiJson)
	assert.Nil(t, err)

	// invokeBean
	beanAbi, err := abi.GetBeanAbi("invoke.ShareInvoke")
	assert.Nil(t, err)
	invokePayload, err := GenPayload(beanAbi, "friend1", "100", "[\"friend2\",\"friend3\"]")
	assert.Nil(t, err)
	invokeTx := rpc.NewTransaction(ecKey.GetAddress().Hex()).Invoke(contractAddress, invokePayload).VMType(rpc.HVM)
	invokeTx.Sign(key)
	receipt, err = hrpc.InvokeContract(invokeTx)
	assert.Nil(t, err)

	// method displayFriends
	//beanAbi1, err := abi.GetMethodAbi("displayMan")
	//assert.Nil(t, err)
	//invokePayload1, err := GenPayload(beanAbi1, `{"m":{1:"111",2:"222"},"name":"Jack","number":111}`)
	//assert.Nil(t, err)
	//invokeTx1 := rpc.NewTransaction(ecKey.GetAddress().Hex()).Invoke(contractAddress, invokePayload1).VMType(rpc.HVM)
	//invokeTx1.Sign(key)
	//receipt, err = hrpc.InvokeContract(invokeTx1)
	//assert.Nil(t, err)
	//
	//logger.Debug(receipt)
	//
	//// method shareMoney, not support func which params is interface
	//beanAbi2, err := abi.GetMethodAbi("shareMoney")
	//assert.Nil(t, err)
	//
	//invokePayload2, err := GenPayload(beanAbi2, "friend1", "100", "[\"friend2\",\"friend3\"]")
	//assert.Nil(t, err)
	//invokeTx2 := rpc.NewTransaction(ecKey.GetAddress().Hex()).Invoke(contractAddress, invokePayload2).VMType(rpc.HVM)
	//invokeTx2.Sign(key)
	//receipt, err = hrpc.InvokeContract(invokeTx2)
	//assert.Nil(t, err)
	//
	//logger.Debug(receipt)
	//
	//// method printInt
	//beanAbi3, err := abi.GetMethodAbi("printInt")
	//assert.Nil(t, err)
	//invokePayload3, err := GenPayload(beanAbi3, "1")
	//assert.Nil(t, err)
	//invokeTx3 := rpc.NewTransaction(ecKey.GetAddress().Hex()).Invoke(contractAddress, invokePayload3).VMType(rpc.HVM)
	//invokeTx3.Sign(key)
	//receipt, err = hrpc.InvokeContract(invokeTx3)
	//assert.Nil(t, err)
	//
	//logger.Debug(receipt)
}
