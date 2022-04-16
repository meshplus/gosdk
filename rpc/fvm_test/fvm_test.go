package fvm_test

import (
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/fvm"
	"github.com/meshplus/gosdk/rpc"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestDeploy(t *testing.T) {
	t.Skip()
	rp := rpc.NewRPCWithPath("../../conf")
	wasmPath1 := "./wasm/SetHash-gc.wasm"
	buf, err := ioutil.ReadFile(wasmPath1)
	if err != nil {
		t.Fatal(err)
	}
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}

	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(common.Bytes2Hex(buf)).VMType(rpc.FVM)
	transaction.Sign(guomiKey)
	rep, err := rp.SignAndDeployContract(transaction, guomiKey)
	if err != nil {
		t.Error(err)
	}

	invokeInput := []byte{32, 115, 101, 116, 95, 104, 97, 115, 104, 24, 107, 101, 121, 48, 48, 49, 100, 116, 104, 105, 115, 32, 105, 115, 32, 116, 104, 101, 32, 118, 97, 108, 117, 101, 32, 111, 102, 32, 48, 48, 48, 49}
	invokeTrans := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput).VMType(rpc.FVM)
	invokeTrans.Sign(guomiKey)
	_, err = rp.SignAndInvokeContract(invokeTrans, guomiKey)
	if err != nil {
		t.Error(err)
	}

	invokeInput2 := []byte{32, 103, 101, 116, 95, 104, 97, 115, 104, 24, 107, 101, 121, 48, 48, 49}
	invokeTrans2 := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput2).VMType(rpc.FVM)
	invokeTrans2.Sign(guomiKey)
	recipt2, err := rp.SignAndInvokeContract(invokeTrans2, guomiKey)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "dthis is the value of 0001", string(common.Hex2Bytes(recipt2.Ret)))
}

func TestDemo(t *testing.T) {
	t.Skip()
	const currentABI = `{
  "contract": {
    "name": "SetHash",
    "constructor": {
      "input": []
    }
  },
  "methods": [
    {
      "name": "set_hash",
      "input": [
        {
          "type_id": 0
        },
        {
          "type_id": 0
        }
      ],
      "output": []
    },
    {
      "name": "get_hash",
      "input": [
        {
          "type_id": 0
        }
      ],
      "output": [
        {
          "type_id": 1
        }
      ]
    }
  ],
  "types": [
    {
      "id": 0,
      "type": "primitive",
      "primitive": "str"
    },
    {
      "id": 1,
      "type": "primitive",
      "primitive": "str"
    }
  ]
}`
	rp := rpc.NewRPCWithPath("../../conf")
	wasmPath1 := "./wasm/SetHash-gc.wasm"
	buf, err := ioutil.ReadFile(wasmPath1)
	if err != nil {
		t.Fatal(err)
	}
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}

	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(common.Bytes2Hex(buf)).VMType(rpc.FVM)
	transaction.Sign(guomiKey)
	rep, err := rp.SignAndDeployContract(transaction, guomiKey)
	if err != nil {
		t.Error(err)
	}
	a, err := fvm.GenAbi(strings.NewReader(currentABI))
	if err != nil {
		t.Error(err)
	}
	invokeInput, err := fvm.Encode(a, "set_hash", "key", "value")
	if err != nil {
		t.Error(err)
	}
	invokeTrans := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput).VMType(rpc.FVM)
	invokeTrans.Sign(guomiKey)
	_, err = rp.SignAndInvokeContract(invokeTrans, guomiKey)
	if err != nil {
		t.Error(err)
	}

	invokeInput2, err := fvm.Encode(a, "get_hash", "key")
	if err != nil {
		t.Error(err)
	}
	invokeTrans2 := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput2).VMType(rpc.FVM)
	invokeTrans2.Sign(guomiKey)
	recipt2, err := rp.SignAndInvokeContract(invokeTrans2, guomiKey)
	if err != nil {
		t.Error(err)
	}

	res, err := fvm.DecodeRet(a, "get_hash", common.Hex2Bytes(recipt2.Ret))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "value", res.Params[0].GetVal())
}
