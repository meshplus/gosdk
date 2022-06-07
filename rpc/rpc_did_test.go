package rpc

import (
	"fmt"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/bvm"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDID_SetChainID(t *testing.T) {
	t.Skip("skip this test")
	rpc := NewRPC()
	// account for test
	var genesisAccountJson = `{"address":"0x000","version":"4.0", "algo":"0x03","publicKey":"","privateKey":""}`
	key, _ := account.GenKeyFromAccountJson(genesisAccountJson, "")
	opt := bvm.NewDIDSetChainIDOperation("chainID_01")
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	_, err := rpc.SignAndInvokeContract(tx, key)
	assert.Nil(t, err)
}

func TestDID_GetChainID(t *testing.T) {
	t.Skip("skip this test")
	rpc := NewRPC()
	res, _ := rpc.GetNodeChainID()
	fmt.Println(res)

}

func registerDIDAccount(rpc *RPC, admins []string) *account.DIDKey {
	rpc.SetLocalChainID()
	randNum := common.RandInt(100)
	var accountJson string
	password := "hyper"
	switch randNum % 3 {
	case 0:
		accountJson, _ = account.NewAccountED25519(password)
	case 1:
		accountJson, _ = account.NewAccountSm2(password)
	case 2:
		accountJson, _ = account.NewAccount(password)
	}
	key, _ := account.GenKeyFromAccountJson(accountJson, password)
	suffix := common.RandomString(10)
	didKey := account.NewDIDAccount(key.(account.Key), rpc.chainID, suffix)

	puKey, _ := GenDIDPublicKeyFromDIDKey(didKey)
	document := NewDIDDocument(didKey.GetAddress(), puKey, admins)
	tx := NewTransaction(didKey.GetAddress()).Register(document)
	_, err := rpc.SendDIDTransaction(tx, didKey)
	if err == nil {
		return didKey
	}
	return nil
}

func TestDID_Account(t *testing.T) {
	t.Skip("skip this test")
	rpc := NewRPC()
	adminAccount := registerDIDAccount(rpc, nil)
	ac := registerDIDAccount(rpc, []string{adminAccount.GetAddress()})

	tx := NewTransaction(adminAccount.GetAddress()).MaintainDID(ac.GetAddress(), DID_FREEZE)
	_, err := rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)
	tx = NewTransaction(adminAccount.GetAddress()).MaintainDID(ac.GetAddress(), DID_UNFREEZE)
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)

	tempAc := registerDIDAccount(rpc, nil)
	puKey, _ := GenDIDPublicKeyFromDIDKey(tempAc)
	tx = NewTransaction(adminAccount.GetAddress()).UpdatePublicKey(ac.GetAddress(), puKey)
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)

	tx = NewTransaction(adminAccount.GetAddress()).UpdateAdmins(ac.GetAddress(), []string{tempAc.GetAddress()})
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)

	adminAccount = tempAc
	tx = NewTransaction(adminAccount.GetAddress()).MaintainDID(ac.GetAddress(), DID_UNFREEZE)
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)
	tx = NewTransaction(adminAccount.GetAddress()).MaintainDID(ac.GetAddress(), DID_ABANDON)
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)
}

func TestDID_Credential(t *testing.T) {
	t.Skip("skip this test")
	rpc := NewRPC()
	holder := registerDIDAccount(rpc, nil)
	issuer := registerDIDAccount(rpc, nil)

	cred := NewDIDCredential("type", issuer.GetAddress(), holder.GetAddress(), "", time.Now().UnixNano(), time.Now().UnixNano()+1e11)
	cred.Sign(issuer)
	tx := NewTransaction(holder.GetAddress()).UploadCredential(cred)
	_, err := rpc.SendDIDTransaction(tx, holder)
	assert.Nil(t, err)

	tx = NewTransaction(issuer.GetAddress()).DownloadCredential(cred.ID)
	_, err = rpc.SendDIDTransaction(tx, issuer)
	assert.Nil(t, err)

}

func TestDID_EXTransaction(t *testing.T) {
	t.Skip("skip this test")
	rpc := NewRPC()
	ac := registerDIDAccount(rpc, nil)
	_, err := rpc.GetDIDDocument(ac.GetAddress())
	assert.Nil(t, err)

	_, err = rpc.GetNodeChainID()
	assert.Nil(t, err)

	holder := registerDIDAccount(rpc, nil)
	issuer := registerDIDAccount(rpc, nil)
	cred := NewDIDCredential("type", issuer.GetAddress(), holder.GetAddress(), "", time.Now().UnixNano(), time.Now().UnixNano()+1e11)
	cred.Sign(issuer)
	tx := NewTransaction(holder.GetAddress()).UploadCredential(cred)
	_, err = rpc.SendDIDTransaction(tx, holder)
	assert.Nil(t, err)

	_, err = rpc.GetCredentialPrimaryMessage(cred.ID)
	assert.Nil(t, err)

	valid, _ := rpc.CheckCredentialValid(cred.ID)
	assert.True(t, valid)

	abandoned, _ := rpc.CheckCredentialAbandoned(cred.ID)
	assert.True(t, !abandoned)
}
