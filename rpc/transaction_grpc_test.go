package rpc

import (
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestTransactionGrpc_SendTransaction(t *testing.T) {
	t.Skip()
	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
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
	transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
	transaction.Sign(guomiKey)
	ans, err := tg.SendTransaction(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestTransactionGrpc_SendTransactionTLS(t *testing.T) {
	t.Skip()
	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
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
	transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
	transaction.Sign(guomiKey)
	ans, err := tg.SendTransaction(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestTransactionGrpc_SendTransaction2(t *testing.T) {
	t.Skip()
	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
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
	transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
	transaction.Sign(guomiKey)
	ans, err := tg.SendTransaction(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestTransactionGrpc_SendTransactionReturnReceipt(t *testing.T) {
	t.Skip()
	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
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
	transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
	transaction.Sign(guomiKey)
	ans, err := tg.SendTransactionReturnReceipt(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

// BenchmarkSendTransaction-8   	    1270	    829267 ns/op
func BenchmarkSendTransaction(b *testing.B) {
	b.Skip()
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}

	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
		StreamNumber: 10,
	})
	assert.Nil(b, err)
	defer tg.Close()
	for i := 0; i < b.N; i++ {
		transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
		transaction.Sign(guomiKey)
		_, err = tg.SendTransaction(transaction)
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkSendTransaction2-8   	    1314	   1011307 ns/op
func BenchmarkSendTransaction2(b *testing.B) {
	b.Skip()
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}

	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(b, err)
	defer tg.Close()
	for i := 0; i < b.N; i++ {
		transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
		transaction.Sign(guomiKey)
		_, err = tg.SendTransaction(transaction)
		if err != nil {
			b.Error(err)
		}

	}
}

// BenchmarkSendTransactionReturnReceipt-8   	       2	 501816208 ns/op
func BenchmarkSendTransactionReturnReceipt(b *testing.B) {
	b.Skip()
	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
		StreamNumber: 10,
	})
	assert.Nil(b, err)
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

	for i := 0; i < b.N; i++ {
		transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
		transaction.Sign(guomiKey)
		_, err = tg.SendTransactionReturnReceipt(transaction)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestSyncGo(t *testing.T) {
	t.Skip()
	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
		StreamNumber: 10,
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

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
			transaction.Sign(guomiKey)
			_, err = tg.SendTransactionReturnReceipt(transaction)
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	t.Log("pass")
}

// BenchmarkSendTransactionReturnReceipt2-8   	       3	 499972806 ns/op
func BenchmarkSendTransactionReturnReceipt2(b *testing.B) {
	b.Skip()
	tg, err := NewGRPC(BindNodes(0)).NewTransactionGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(b, err)
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

	for i := 0; i < b.N; i++ {
		transaction := NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
		transaction.Sign(guomiKey)
		_, err = tg.SendTransactionReturnReceipt(transaction)
		if err != nil {
			b.Error(err)
		}
	}
}
