package rpc

import (
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/bvm"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func iniDID(t *testing.T) {
	accountJson := `{"address":"0xfbca6a7e9e29728773b270d3f00153c75d04e1ad","version":"4.0","algo":"0x13","publicKey":"049c330d0aea3d9c73063db339b4a1a84d1c3197980d1fb9585347ceeb40a5d262166ee1e1cb0c29fd9b2ef0e4f7a7dfb1be6c5e759bf411c520a616863ee046a4","privateKey":"5f0a3ea6c1d3eb7733c3170f2271c10c1206bc49b6b2c7e550c9947cb8f098e3"}`
	key, _ := account.GenKeyFromAccountJson(accountJson, "")
	opt := bvm.NewDIDSetChainIDOperation("chainID_01")
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	_, err := rpc.SignAndInvokeContract(tx, key)
	assert.Nil(t, err)
}

func TestDidGrpc_SendDIDTransaction(t *testing.T) {
	t.Skip()
	iniDID(t)
	g := NewGRPC()
	tg, err := g.NewDidGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()
	var accountJson string
	password := "hyper"
	accountJson, _ = account.NewAccountSm2(password)
	key, _ := account.GenDIDKeyFromAccountJson(accountJson, password)
	suffix := common.RandomString(10)
	didKey := account.NewDIDAccount(key.(account.Key), "chainID_01", suffix)
	puKey, _ := GenDIDPublicKeyFromDIDKey(didKey)
	document := NewDIDDocument(didKey.GetAddress(), puKey, nil)
	transaction := NewTransaction(didKey.GetAddress()).Register(document)
	transaction.Sign(didKey)
	ans, err := tg.SendDIDTransaction(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestDidGrpc_SendDIDTransactionReturnReceipt(t *testing.T) {
	t.Skip()
	iniDID(t)
	g := NewGRPC()
	tg, err := g.NewDidGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()
	var accountJson string
	password := "hyper"
	accountJson, _ = account.NewAccountSm2(password)
	key, _ := account.GenDIDKeyFromAccountJson(accountJson, password)
	suffix := common.RandomString(10)
	didKey := account.NewDIDAccount(key.(account.Key), "chainID_01", suffix)
	puKey, _ := GenDIDPublicKeyFromDIDKey(didKey)
	document := NewDIDDocument(didKey.GetAddress(), puKey, nil)
	transaction := NewTransaction(didKey.GetAddress()).Register(document)
	transaction.Sign(didKey)
	ans, err := tg.SendDIDTransactionReturnReceipt(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}
