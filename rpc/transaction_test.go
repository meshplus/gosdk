package rpc

import (
	"github.com/magiconair/properties/assert"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common/hexutil"
	"github.com/meshplus/gosdk/kvsql"
	"testing"
)

func TestTransaction(t *testing.T) {
	tax := NewTransaction("0x0000000000000000")
	tax.setTxVersion("1.0")
	tax.SetNonce(int64(1))
	tax.SetExtra("extra")
	tax.SetFrom("0x0000000000000000")
	tax.SetHasExtra(true)
	tax.SetIsDeploy(false)
	tax.SetIsInvoke(false)
	tax.SetIsMaintain(true)
	tax.SetIsPrivateTxm(false)
	tax.Simulate(false)
	tax.SetIsValue(false)
	tax.SetOpcode(1)
	tax.SetParticipants([]string{})
	tax.SetTo("0x0000000000000000")
	tax.SetVmType("HVM")
	tax.Nonce(int64(1))
	tax.Timestamp(int64(1))
	tax.Value(int64(1))
	tax.SetPayload("nothing")
	tax.SetValue(int64(123))
	tax.SetTimestamp(int64(1))
	tax.SetSignature("signature")
	tax.SerializeToString()
}

func TestNeedHashString(t *testing.T) {
	tax := NewTransaction("0x0000000000000000")
	tax.SetNonce(int64(1))
	tax.SetExtra("extra")
	tax.SetFrom("0x0000000000000000")
	tax.SetHasExtra(true)
	tax.SetIsDeploy(false)
	tax.SetIsInvoke(false)
	tax.SetIsMaintain(true)
	tax.SetIsPrivateTxm(false)
	tax.Simulate(false)
	tax.SetIsValue(false)
	tax.SetOpcode(1)
	tax.SetParticipants([]string{})
	tax.SetTo("0x0000000000000000")
	tax.SetVmType("HVM")
	tax.Nonce(int64(1))
	tax.Timestamp(int64(1))
	tax.Value(int64(1))
	tax.SetPayload("nothing")
	tax.SetValue(int64(123))
	tax.SetTimestamp(int64(1))

	tax.setTxVersion("1.8")
	expect18 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM"
	assert.Equal(t, expect18, needHashString(tax))

	tax.setTxVersion("2.0")
	expect20 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.0"
	assert.Equal(t, expect20, needHashString(tax))

	tax.setTxVersion("2.1")
	expect21 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.1&extraid="
	assert.Equal(t, expect21, needHashString(tax))

	tax.setTxVersion("2.2")
	expect22 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.2&extraid=&cname="
	assert.Equal(t, expect22, needHashString(tax))

	tax.setTxVersion("2.3")
	expect23 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.3&extraid=&cname="
	assert.Equal(t, expect23, needHashString(tax))
}

func TestAddKVSQLType(t *testing.T) {
	// todo open after kvsql
	t.Skip("latest flato not support kvsql")
	hrpc := NewRPCWithPath("../conf")

	js, err := account.NewAccountSm2("12345678")
	assert.Equal(t, nil, err)

	gmAcc, err := account.GenKeyFromAccountJson(js, "12345678")

	assert.Equal(t, nil, err)
	newAddress := gmAcc.(*account.SM2Key).GetAddress()
	transaction := NewTransaction(newAddress.Hex())
	transaction.VMType(KVSQL)
	transaction.Deploy(hexutil.Encode([]byte("KVSQL")))

	transaction.Sign(gmAcc)

	// 建库
	txReceipt, err := hrpc.SignAndDeployContract(transaction, gmAcc)
	assert.Equal(t, nil, err)

	// 建表
	str := "CREATE TABLE IF NOT EXISTS testTable (id bigint(20) NOT NULL, name varchar(32) NOT NULL, exp bigint(20), money double(16,2) NOT NULL DEFAULT '99', primary key (id), unique key name (name));"
	tranInvoke := NewTransaction(newAddress.Hex()).InvokeSql(txReceipt.ContractAddress, []byte(str))
	tranInvoke.VMType(KVSQL)
	tranInvoke.Sign(gmAcc)
	res, err := hrpc.SignAndInvokeContract(tranInvoke, gmAcc)
	assert.Equal(t, nil, err)
	t.Log(res.Ret)

	// 插入
	in := "insert into testTable (id, name, exp, money) values (1, 'test', 1, 1.1);"
	tranInvoke2 := NewTransaction(newAddress.Hex()).InvokeSql(txReceipt.ContractAddress, []byte(in))
	tranInvoke2.VMType(KVSQL)
	tranInvoke2.Sign(gmAcc)
	_, err = hrpc.SignAndInvokeContract(tranInvoke2, gmAcc)
	assert.Equal(t, nil, err)

	// 获取
	sl := "select * from testTable where id = 1;"
	tranInvoke3 := NewTransaction(newAddress.Hex()).InvokeSql(txReceipt.ContractAddress, []byte(sl))
	tranInvoke3.VMType(KVSQL)
	tranInvoke3.Sign(gmAcc)
	res, err = hrpc.SignAndInvokeContract(tranInvoke3, gmAcc)
	assert.Equal(t, nil, err)
	t.Log(res.Ret)
	b, _ := hexutil.Decode(res.Ret)
	rs := kvsql.DecodeRecordSet(b)
	t.Log(rs)
}
