package rpc

import (
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/bvm"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/hvm"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestRPC_BVMCrossChainAnchorTx1(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON("", "")
	assert.Nil(t, err)

	operation1 := bvm.NewSystemAnchorOperation(bvm.RegisterAnchor, "node1", "ns2")
	payload1 := bvm.EncodeOperation(operation1)
	tx1 := NewTransaction(key.GetAddress().Hex()).Invoke(operation1.Address(), payload1).VMType(BVM)
	//tx1.Sign(key)
	rpc.namespace = "global"
	re1, err := rpc.SignAndInvokeCrossChainContract(tx1, "invokeAnchorContract", key)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re1.Ret))
}

func TestRPC_BVMCrossChainAnchorTx2(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON("", "")
	assert.Nil(t, err)

	operation2 := bvm.NewNormalAnchorOperation(bvm.RegisterAnchor, []string{"node1"})
	payload2 := bvm.EncodeOperation(operation2)
	tx2 := NewTransaction(key.GetAddress().Hex()).Invoke(operation2.Address(), payload2).VMType(BVM)
	rpc.namespace = "ns2"
	re2, err := rpc.SignAndInvokeCrossChainContract(tx2, "invokeAnchorContract", key)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re2.Ret))
}

func TestRPC_BVMCrossChainAnchorUnregister(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON("", "")
	assert.Nil(t, err)

	operation2 := bvm.NewSystemAnchorOperation(bvm.UnRegisterAnchor, "node1", "ns2")
	payload2 := bvm.EncodeOperation(operation2)
	tx2 := NewTransaction(key.GetAddress().Hex()).Invoke(operation2.Address(), payload2).VMType(BVM)
	rpc.namespace = "global"
	re2, err := rpc.SignAndInvokeCrossChainContract(tx2, "invokeAnchorContract", key)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re2.Ret))
}

func TestRPC_BVMCrossChainAnchorTxReplace(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON("", "")
	assert.Nil(t, err)

	operation := bvm.NewSystemAnchorOperation(bvm.ReplaceAnchor, "node2", "ns1", "node1", "ns1")
	payload := bvm.EncodeOperation(operation)
	tx2 := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	rpc.namespace = "global"
	re2, err := rpc.SignAndInvokeCrossChainContract(tx2, "invokeAnchorContract", key)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re2.Ret))
}

func TestRPC_BVMGetAnchorStatus(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON("", "")
	assert.Nil(t, err)

	operation1 := bvm.NewSystemAnchorOperation(bvm.ReadAnchor, "ns3")
	payload1 := bvm.EncodeOperation(operation1)
	tx1 := NewTransaction(key.GetAddress().Hex()).Invoke(operation1.Address(), payload1).VMType(BVM)
	//tx1.Sign(key)
	rpc.namespace = "global"
	re1, err := rpc.SignAndInvokeContract(tx1, key)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re1.Ret))
}

func TestRPC_BVMGetCrossChainTx(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON("", "")
	assert.Nil(t, err)

	operation1 := bvm.NewSystemAnchorOperation(bvm.ReadCrossChain, "ns1__@__0x16b22229e93fdad80cef0920cc26d2fb6acb32945f8f5fe2248f0adcc3aa2c29__@__ns2")
	payload1 := bvm.EncodeOperation(operation1)
	tx1 := NewTransaction(key.GetAddress().Hex()).Invoke(operation1.Address(), payload1).VMType(BVM)
	//tx1.Sign(key)
	rpc.namespace = "global"
	re1, err := rpc.SignAndInvokeContract(tx1, key)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re1.Ret))
}

func TestRPC_Hvm_Cross_Chain_Dep(t *testing.T) {
	t.Skip()
	deployJar, err := ioutil.ReadFile("../hvmtestfile/crosschain/cctest-1.0-SNAPSHOT-cross_chain.jar")
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
	rpc.namespace = "ns2"
	receipt, err := rpc.SignAndDeployContract(transaction, key)
	assert.Nil(t, err)
	t.Log("contract address:", receipt.ContractAddress)
}

func TestRPC_Hvm_Cross_Chain_SetTargetAddress(t *testing.T) {
	// ns1: 0xd8b9d000225f5f636859e0bb5599ea833f349326
	// ns2: 0x1d1f50c19b7d9d2b0a8dfd5ea89d96aff2d37f82
	t.Skip()

	a, sysErr := account.NewAccountED25519("12345678")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	key, sysErr := account.GenKeyFromAccountJson(a, "12345678")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}

	newAddress := key.(*account.ED25519Key).GetAddress()

	abiPath := "../hvmtestfile/crosschain/hvm.abi"
	abiJson, rerr := common.ReadFileAsString(abiPath)
	assert.Nil(t, rerr)
	abi, gerr := hvm.GenAbi(abiJson)
	if gerr != nil {
		logger.Error(gerr)
	}
	addAbi, err := abi.GetMethodAbi("SetTargetAddress")
	if err != nil {
		logger.Error(err)
	}

	payload, err := hvm.GenPayload(addAbi, "0x1d1f50c19b7d9d2b0a8dfd5ea89d96aff2d37f82")
	if err != nil {
		logger.Error(err)
	}

	transaction1 := NewTransaction(newAddress.Hex()).Invoke("0xd8b9d000225f5f636859e0bb5599ea833f349326", payload).VMType(HVM)
	rpc.namespace = "ns1"
	invokeContract, err := rpc.SignAndInvokeContract(transaction1, key)
	if err != nil {
		t.Error(err)
	}
	t.Log(invokeContract)
	t.Log(invokeContract.Ret)
}

func TestRPC_Hvm_Cross_Chain_Invoke(t *testing.T) {
	// ns1: 0xd8b9d000225f5f636859e0bb5599ea833f349326
	// ns2: 0x1d1f50c19b7d9d2b0a8dfd5ea89d96aff2d37f82
	t.Skip()

	a, sysErr := account.NewAccountED25519("12345678")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	key, sysErr := account.GenKeyFromAccountJson(a, "12345678")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}

	newAddress := key.(*account.ED25519Key).GetAddress()

	abiPath := "../hvmtestfile/crosschain/hvm.abi"
	abiJson, rerr := common.ReadFileAsString(abiPath)
	assert.Nil(t, rerr)
	abi, gerr := hvm.GenAbi(abiJson)
	if gerr != nil {
		logger.Error(gerr)
	}
	addAbi, err := abi.GetMethodAbi("Add")
	//addAbi, err := abi.GetMethodAbi("SetTargetAddress")
	if err != nil {
		logger.Error(err)
	}

	payload, err := hvm.GenPayload(addAbi, "11")
	if err != nil {
		logger.Error(err)
	}

	transaction1 := NewTransaction(newAddress.Hex()).Invoke("0xd8b9d000225f5f636859e0bb5599ea833f349326", payload).VMType(HVM)
	rpc.namespace = "ns1"
	invokeContract, err := rpc.SignAndInvokeCrossChainContract(transaction1, "invokeContract", key)
	if err != nil {
		t.Error(err)
	}
	t.Log(invokeContract)
	t.Log(invokeContract.Ret)
}

func TestRPC_Cross_Chain_Timeout(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON("", "")
	assert.Nil(t, err)
	operation := bvm.NewSystemAnchorOperation(bvm.Timeout,
		"ns1__@__0x16b22229e93fdad80cef0920cc26d2fb6acb32945f8f5fe2248f0adcc3aa2c29__@__ns2")
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	rpc.namespace = "global"
	invoke, err := rpc.SignAndInvokeCrossChainContract(tx, "invokeTimeoutContract", key)
	if err != nil {
		t.Error(err)
	}
	t.Log(invoke)
}
