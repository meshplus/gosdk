package rpc

import (
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/bvm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootCA(t *testing.T) {
	t.Skip("skip")
	accountJsons := []string{""}

	t.Run("add_root_ca", func(t *testing.T) {

		rpc := NewRPC()
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		rootCA := `-----BEGIN CERTIFICATE-----
MIICODCCAeSgAwIBAgIBATAKBggqhkjOPQQDAjB0MQkwBwYDVQQIEwAxCTAHBgNV
BAcTADEJMAcGA1UECRMAMQkwBwYDVQQREwAxDjAMBgNVBAoTBWZsYXRvMQkwBwYD
VQQLEwAxDjAMBgNVBAMTBW5vZGU1MQswCQYDVQQGEwJaSDEOMAwGA1UEKhMFZWNl
cnQwIBcNMjAwNTIyMDUyNTQ4WhgPMjEyMDA0MjgwNjI1NDhaMHQxCTAHBgNVBAgT
ADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAxCTAHBgNVBBETADEOMAwGA1UEChMFZmxh
dG8xCTAHBgNVBAsTADEOMAwGA1UEAxMFbm9kZTUxCzAJBgNVBAYTAlpIMQ4wDAYD
VQQqEwVlY2VydDBWMBAGByqGSM49AgEGBSuBBAAKA0IABLs/ih4HBDxs0nHNliXt
sZKAaW2fdgi51H07eMNtHib/8R55GYFinvjmeJa9OvbyYsWAbCTQsBu3LV7aPpOT
0QGjaDBmMA4GA1UdDwEB/wQEAwIChDAmBgNVHSUEHzAdBggrBgEFBQcDAgYIKwYB
BQUHAwEGAioDBgOBCwEwDwYDVR0TAQH/BAUwAwEB/zANBgNVHQ4EBgQEAQIDBDAM
BgMqVgEEBWVjZXJ0MAoGCCqGSM49BAMCA0IA/9MkrtXUMni8THIhym5t2manIffB
xpW51rRuxPjB/0wZvgKPDWuOcq8eqsOOX/d4wFtMwJ4e7MqLZu6RWSc2MQA=
-----END CERTIFICATE-----
`
		opt := bvm.NewRootCAAddOperation(rootCA)
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		res, serr := rpc.SignAndInvokeContract(tx, key)
		assert.Nil(t, serr)
		ret := bvm.Decode(res.Ret)
		assert.True(t, ret.Success)
		t.Log(string(ret.Ret))
	})

	t.Run("get_root_ca", func(t *testing.T) {
		rpc := NewRPC()
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		opt := bvm.NewRootCAGetOperation()
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		res, serr := rpc.SignAndInvokeContract(tx, key)
		assert.Nil(t, serr)
		ret := bvm.Decode(res.Ret)
		t.Log(ret)
	})

	t.Run("delete_node_for_none", func(t *testing.T) {
		//t.Skip("skip this test")
		rpc := NewRPC()
		newAccount, err := account.NewAccount(pwd)
		assert.Nil(t, err)
		key, err := account.NewAccountFromAccountJSON(newAccount, pwd)
		assert.Nil(t, err)
		opt := bvm.NewProposalDirectOperationForNode(bvm.NewNodeRemoveVPOperation("node2", "global"))
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		res, serr := rpc.SignAndInvokeContract(tx, key)
		assert.Nil(t, serr)
		ret := bvm.Decode(res.Ret)
		//assert.True(t, ret.Success)
		t.Log(ret)
	})

	t.Run("delete_node_for_center", func(t *testing.T) {
		//t.Skip("skip this test")
		rpc := NewRPC()
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		opt := bvm.NewProposalDirectOperationForNode(bvm.NewNodeRemoveVPOperation("node3", "global"))
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		res, serr := rpc.SignAndInvokeContract(tx, key)
		assert.Nil(t, serr)
		ret := bvm.Decode(res.Ret)
		//assert.True(t, ret.Success)
		t.Log(ret)
	})

	t.Run("get_ca_mode", func(t *testing.T) {
		//t.Skip("skip this test")
		rpc := NewRPC()
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		opt := bvm.NewProposalDirectOperationForCA(bvm.NewCAGetCAModeOperation())
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		t.Logf(tx.GetPayload())
		res, serr := rpc.SignAndInvokeContract(tx, key)
		assert.Nil(t, serr)
		ret := bvm.Decode(res.Ret)
		//assert.True(t, ret.Success)
		t.Log(ret)
		t.Log(string(ret.Ret))
	})

	t.Run("set_ca_mode", func(t *testing.T) {
		//t.Skip("skip this test")
		rpc := NewRPC()
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		mode, _ := bvm.NewCASetCAModeOperation("center")
		opt := bvm.NewProposalCreateOperationForCA(mode)
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		t.Log(tx.GetPayload())
		res, serr := rpc.SignAndInvokeContract(tx, key)
		assert.Nil(t, serr)
		ret := bvm.Decode(res.Ret)
		//assert.True(t, ret.Success)
		t.Log(ret)
		t.Log(string(ret.Ret))
	})

	t.Run("set_ca_mode_first", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		mode, _ := bvm.NewCASetCAModeOperation("center")
		res := completeProposal(t, creator, 6, bvm.NewProposalCreateOperationForCA(mode))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		t.Log(res[0].Msg)
	})

	t.Run("revoke_cert", func(t *testing.T) {
		//t.Skip("skip this test")
		rpc := NewRPC()
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		cert := []byte(`-----BEGIN CERTIFICATE-----
MIICSTCCAfWgAwIBAgIBATAKBggqhkjOPQQDAjB0MQkwBwYDVQQIEwAxCTAHBgNV
BAcTADEJMAcGA1UECRMAMQkwBwYDVQQREwAxDjAMBgNVBAoTBWZsYXRvMQkwBwYD
VQQLEwAxDjAMBgNVBAMTBW5vZGUyMQswCQYDVQQGEwJaSDEOMAwGA1UEKhMFZWNl
cnQwIBcNMjAwNTIxMDU1ODU2WhgPMjEyMDA0MjcwNjU4NTZaMHQxCTAHBgNVBAgT
ADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAxCTAHBgNVBBETADEOMAwGA1UEChMFZmxh
dG8xCTAHBgNVBAsTADEOMAwGA1UEAxMFbm9kZTQxCzAJBgNVBAYTAlpIMQ4wDAYD
VQQqEwVlY2VydDBWMBAGByqGSM49AgEGBSuBBAAKA0IABBI3ewNK21vHNOPG6U3X
mKJohSNNz72QKDxUpRt0fCJHwaGYfSvY4cnqkbliclfckUTpCkFSRr4cqN6PURCF
zkWjeTB3MA4GA1UdDwEB/wQEAwIChDAmBgNVHSUEHzAdBggrBgEFBQcDAgYIKwYB
BQUHAwEGAioDBgOBCwEwDwYDVR0TAQH/BAUwAwEB/zANBgNVHQ4EBgQEAQIDBDAP
BgNVHSMECDAGgAQBAgMEMAwGAypWAQQFZWNlcnQwCgYIKoZIzj0EAwIDQgDJibFh
a1tZ3VhL3WIs36DqOS22aetvcn2dXHH9Pw5/s2XI70Mr3ow3RKqJmdmi0PsmLr+K
pCFkuMv2bHnkWuiZAQ==
-----END CERTIFICATE-----
`)
		opt, _ := bvm.NewCertRevokeOperation(cert, nil)
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		t.Log(tx.GetPayload())
		res, serr := rpc.SignAndInvokeContract(tx, key)
		assert.Nil(t, serr)
		ret := bvm.Decode(res.Ret)
		//assert.True(t, ret.Success)
		t.Log(ret)
		t.Log(string(ret.Ret))
	})
}
