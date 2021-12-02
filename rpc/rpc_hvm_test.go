package rpc

import (
	"fmt"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/hvm"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const (
	NoArgs = iota
	JavaLangInteger
	Int
	JavaLangShort
	Short
	JavaLangLong
	Long
	JavaLangByte
	Byte
	JavaLangBoolean
	Boolean
	JavaLangCharacter
	Char
	JavaLangFloat
	Float
	JavaLangDouble
	Double
	String
	JavaLangObject
)

func TestEncode(t *testing.T) {
	// read abi
	abiPath := "../hvmtestfile/hvm1.abi"
	abiJson, err := common.ReadFileAsString(abiPath)
	if err != nil {
		t.Error(err)
		return
	}

	abi, err := hvm.GenAbi(abiJson)
	if err != nil {
		t.Error(err)
		return
	}

	easyBean := "cn.test.contract.invoke.EasyInvoke"
	beanAbi, err := abi.GetBeanAbi(easyBean)
	if err != nil {
		t.Error(err)
		return
	}

	// encode
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
	_, sysErr := hvm.GenPayload(beanAbi, "true", "c", "20", "100", "1000", "10000", "1.1", "1.11", "string",
		//`["1f","2f","3f"]`, `[1,2,3]`,
		//`{789:{456:12.2},234:{345:12.2}}`, `[[1,2],[2,4]]`,
		`{"name":"tom","age":21}`, `{"beanName":"bean1","person":{"name":"tom","age":21}}`,
		`["strList1","strList2"]`,
		`[{"name":"tom","age":21},{"name":"jack","age":18}]`,
		`{"person1":{"name":"tom","age":21},"person2":{"name":"jack","age":18}}`,
		`{"bean1":{"beanName":"bean1","person":{"name":"tom","age":21}},"bean2":{"beanName":"bean2","person":{"name":"jack","age":18}}}`)

	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	_, sysErr = hvm.GenPayload(beanAbi, hvm.Convert("true"), hvm.Convert("c"), hvm.Convert("20"), hvm.Convert("20"), hvm.Convert("20"), hvm.Convert("20"), hvm.Convert("20.1"), hvm.Convert("20.2"), hvm.Convert("asasd"),
		//hvm.Convert(fclist), hvm.Convert(fcarray),
		//hvm.Convert(fcmap),
		//hvm.Convert(fclist1),
		hvm.Convert(person3), bean1,
		[]interface{}{"strList1", "strList2"},
		[]interface{}{person1, person2},
		hvm.Convert(fcmaptruct),
		[]interface{}{[]interface{}{"bean1", bean1}, []interface{}{"bean2", bean2}})
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}

}

func TestRPC_Hvm(t *testing.T) {
	t.Skip()
	deployJar, err := DecompressFromJar("../hvmtestfile/fibonacci/fibonacci-1.0-fibonacci.jar")
	if err != nil {
		t.Error(err)
	}

	accountJson, sysErr := account.NewAccountED25519("12345678")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	key, sysErr := account.GenKeyFromAccountJson(accountJson, "12345678")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}

	newAddress := key.(*account.ED25519Key).GetAddress()

	transaction := NewTransaction(newAddress.Hex()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
	transaction.Sign(key)
	receipt, err := rpc.DeployContract(transaction)
	assert.Nil(t, err)
	t.Log("contract address:", receipt.ContractAddress)

	abiPath := "../hvmtestfile/fibonacci/hvm.abi"
	abiJson, rerr := common.ReadFileAsString(abiPath)
	assert.Nil(t, rerr)
	abi, gerr := hvm.GenAbi(abiJson)
	if gerr != nil {
		logger.Error(gerr)
	}

	easyBean := "invoke.InvokeFibonacci"
	beanAbi, err := abi.GetBeanAbi(easyBean)
	if err != nil {
		logger.Error(err)
	}

	payload, err := hvm.GenPayload(beanAbi)
	if err != nil {
		logger.Error(err)
	}

	transaction1 := NewTransaction(newAddress.Hex()).Invoke(receipt.ContractAddress, payload).VMType(HVM)
	transaction1.Sign(key)
	invokeContract, err := rpc.InvokeContract(transaction1)
	if err != nil {
		t.Error(err)
	}
	t.Log(invokeContract)
}

func TestRPC_DirectInvoke(t *testing.T) {
	t.Skip()

	accountJson, sysErr := account.NewAccount("")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	key, sysErr := account.NewAccountFromAccountJSON(accountJson, "")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}

	deployJar, err := DecompressFromJar("../hvmtestfile/directInvoke/directInvokeTest-1.0.jar")
	if err != nil {
		t.Error(err)
	}

	deployTX := NewTransaction(key.GetAddress().Hex()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
	deployTX.Sign(key)
	receipt, err := rpc.DeployContract(deployTX)
	assert.Nil(t, err)
	t.Log("contract address:", receipt.ContractAddress)

	contractAddress := receipt.ContractAddress

	// test method without args
	testInvokeArgs(key, "hello", contractAddress, NoArgs, t)

	// test method with a String type arg
	testInvokeArgs(key, "printString", contractAddress, String, t)

	// test method with an Integer type arg
	testInvokeArgs(key, "printInteger", contractAddress, JavaLangInteger, t)

	// test method with an int type arg
	testInvokeArgs(key, "printInt", contractAddress, Int, t)

	// test method with an Short type arg
	testInvokeArgs(key, "printShort", contractAddress, JavaLangShort, t)

	// test method with an short type arg
	testInvokeArgs(key, "printshort", contractAddress, Short, t)

	// test method with an Long type arg
	testInvokeArgs(key, "printLong", contractAddress, JavaLangLong, t)

	// test method with an long type arg
	testInvokeArgs(key, "printlong", contractAddress, Long, t)

	// test method with an Byte type arg
	testInvokeArgs(key, "printByte", contractAddress, JavaLangByte, t)

	// test method with an byte type arg
	testInvokeArgs(key, "printbyte", contractAddress, Byte, t)

	// test method with an Boolean type arg
	testInvokeArgs(key, "printBoolean", contractAddress, JavaLangBoolean, t)

	// test method with an boolean type arg
	testInvokeArgs(key, "printboolean", contractAddress, Boolean, t)

	// test method with an Character type arg
	testInvokeArgs(key, "printCharacter", contractAddress, JavaLangCharacter, t)

	// test method with an char type arg
	testInvokeArgs(key, "printChar", contractAddress, Char, t)

	// test method with an Float type arg
	testInvokeArgs(key, "printFloat", contractAddress, JavaLangFloat, t)

	// test method with an float type arg
	testInvokeArgs(key, "printfloat", contractAddress, Float, t)

	// test method with an Double type arg
	testInvokeArgs(key, "printDouble", contractAddress, JavaLangDouble, t)

	// test method with an double type arg
	testInvokeArgs(key, "printdouble", contractAddress, Double, t)

	// test method with an User type arg
	testInvokeArgs(key, "printUser", contractAddress, JavaLangObject, t)
}

func testInvokeArgs(key *account.ECDSAKey, methodName, contractAddress string, argType int, t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	paramBuilder := hvm.NewParamBuilder(methodName)
	switch argType {
	case NoArgs:
	case JavaLangInteger:
		paramBuilder.AddInteger(2333)
	case Int:
		paramBuilder.Addint(233)
	case JavaLangShort:
		paramBuilder.AddShort(int16(111))
	case Short:
		paramBuilder.Addshort(int16(11))
	case JavaLangLong:
		paramBuilder.AddLong(int64(99999999))
	case Long:
		paramBuilder.Addlong(int64(9999999))
	case JavaLangByte:
		paramBuilder.AddByte(57)
	case Byte:
		paramBuilder.Addbyte(48)
	case JavaLangBoolean:
		paramBuilder.AddBoolean(true)
	case Boolean:
		paramBuilder.Addbool(false)
	case JavaLangCharacter:
		paramBuilder.AddCharacter('C')
	case Char:
		paramBuilder.Addchar('c')
	case JavaLangFloat:
		paramBuilder.AddFloat(float32(1.0))
	case Float:
		paramBuilder.Addfloat(float32(0.5))
	case JavaLangDouble:
		paramBuilder.AddDouble(1.000000000)
	case Double:
		paramBuilder.Adddouble(0.500000000)
	case JavaLangObject:
		paramBuilder.AddObject("logic.User", User{
			Name: "user1",
			Age:  20,
		})
	case String:
		paramBuilder.AddString("test printString")
	}
	payload := paramBuilder.Build()
	invokeTX := NewTransaction(key.GetAddress().Hex()).Invoke(contractAddress, payload).VMType(HVM)
	invokeTX.Sign(key)
	invokeContract, err := rpc.InvokeContract(invokeTX)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(common.Hex2Bytes(invokeContract.Ret[2:])))
}
