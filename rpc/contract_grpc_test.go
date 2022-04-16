package rpc

import (
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/gosdk/abi"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
)

func TestContractGrpc_DeployContract(t *testing.T) {
	t.Skip()
	g := NewGRPC()
	tg, err := g.NewContractGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	transaction := NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
	transaction.Sign(guomiKey)
	ans, err := tg.DeployContract(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestContractGrpc_DeployContractReturnReceipt(t *testing.T) {
	t.Skip()
	g := NewGRPC()
	tg, err := g.NewContractGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	transaction := NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
	transaction.Sign(guomiKey)
	ans, err := tg.DeployContractReturnReceipt(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestContractGrpc_InvokeContract(t *testing.T) {
	t.Skip()
	g := NewGRPC()
	tg, err := g.NewContractGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	addrTransaction := NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
	addrTransaction.Sign(guomiKey)
	addr, err := tg.DeployContractReturnReceipt(addrTransaction)
	if err != nil {
		t.Error(err)
	}

	ABI, err := abi.JSON(strings.NewReader(abiContract))
	if err != nil {
		t.Error(err)
	}
	packed, err := ABI.Pack("getSum")
	if err != nil {
		t.Error(err)
	}
	transaction := NewTransaction(guomiKey.GetAddress().Hex()).Invoke(addr.ContractAddress, packed)
	transaction.Sign(guomiKey)
	ans, err := tg.InvokeContract(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestContractGrpc_InvokeContractReturnReceipt(t *testing.T) {
	t.Skip()
	g := NewGRPC()
	tg, err := g.NewContractGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	addrTransaction := NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
	addrTransaction.Sign(guomiKey)
	addr, err := tg.DeployContractReturnReceipt(addrTransaction)
	if err != nil {
		t.Error(err)
	}
	ABI, err := abi.JSON(strings.NewReader(abiContract))
	if err != nil {
		t.Error(err)
	}
	packed, err := ABI.Pack("getSum")
	if err != nil {
		t.Error(err)
	}
	transaction := NewTransaction(guomiKey.GetAddress().Hex()).Invoke(addr.ContractAddress, packed)
	transaction.Sign(guomiKey)
	ans, err := tg.InvokeContractReturnReceipt(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestContractGrpc_MaintainContract(t *testing.T) {
	t.Skip()
	g := NewGRPC()
	tg, err := g.NewContractGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	addrTransaction := NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
	addrTransaction.Sign(guomiKey)
	addr, err := tg.DeployContractReturnReceipt(addrTransaction)
	if err != nil {
		t.Error(err)
	}
	// freeze contract
	transactionFreeze := NewTransaction(guomiKey.GetAddress().Hex()).Maintain(2, addr.ContractAddress, "")
	transactionFreeze.Sign(guomiKey)

	ans, err := tg.MaintainContract(transactionFreeze)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)

	// unfreeze
	transactionUnFreeze := NewTransaction(guomiKey.GetAddress().Hex()).Maintain(3, addr.ContractAddress, "")
	transactionUnFreeze.Sign(guomiKey)

	ans2, err := tg.MaintainContractReturnReceipt(transactionUnFreeze)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans2)
}

func TestContractGrpc_DeployContractMulti(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(1))
	tg, err := g.NewTransactionGrpc(ClientOption{
		StreamNumber: 100,
	})
	assert.Nil(t, err)

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	var sg sync.WaitGroup
	for i := 1; i <= 30000; i++ {
		go func(sg *sync.WaitGroup, i int) {
			sg.Add(1)
			defer sg.Done()
			//transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", 0)
			transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", 0)
			transaction.Sign(guomiKey)
			ans, err := tg.SendTransaction(transaction)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(ans)
			t.Log(ans)
			t.Log(i)
		}(&sg, i)
	}
	sg.Wait()
	tg.Close()
}
