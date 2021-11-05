package rpc

import (
	"github.com/magiconair/properties/assert"
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

	setTxVersion("1.8")
	expect18 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM"
	assert.Equal(t, expect18, needHashString(tax))

	setTxVersion("2.0")
	expect20 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.0"
	assert.Equal(t, expect20, needHashString(tax))

	setTxVersion("2.1")
	expect21 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.1&extraid="
	assert.Equal(t, expect21, needHashString(tax))

	setTxVersion("2.2")
	expect22 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.2&extraid=&cname="
	assert.Equal(t, expect22, needHashString(tax))

	setTxVersion("2.3")
	expect23 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.3&extraid=&cname="
	assert.Equal(t, expect23, needHashString(tax))
}
