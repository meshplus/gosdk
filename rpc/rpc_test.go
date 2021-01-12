package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/bvm"
	"io/ioutil"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/meshplus/gosdk/abi"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
)

const (
	//contractAddress = "0x421a1fb06bd9c9fae9b8cdaf8a662cf3c41ffa10"
	abiContract = `[{"constant":false,"inputs":[{"name":"num1","type":"uint32"},{"name":"num2","type":"uint32"}],"name":"add","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"archiveSum","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"getSum","outputs":[{"name":"","type":"uint32"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"increment","outputs":[],"payable":false,"type":"function"}]`
	binContract = "0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029"
)

var (
	//rpc           = NewRPCWithPath("../conf")
	rpc           = NewRPC()
	address       = "bfa5bd992e3eb123c8b86ebe892099d4e9efb783"
	privateKey, _ = account.NewAccountFromPriv("a1fd6ed6225e76aac3884b5420c8cdbb4fde1db01e9ef773415b8f2b5a9b77d4")

	//guomiPub      = "02739518af5e065b22dabb35ea5369a4c64d4865565874a006399bbb0e62e18004"
	guomiPri = "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri      = new(gm.SM2PrivateKey).FromBytes(common.FromHex(guomiPri))
	guomiKey = &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
)

type TestAsyncHandler struct {
	t        *testing.T
	IsCalled bool
}

func (tah *TestAsyncHandler) OnSuccess(receipt *TxReceipt) {
	tah.IsCalled = true
	fmt.Println(receipt.Ret)
}

func (tah *TestAsyncHandler) OnFailure(err StdError) {
	tah.t.Error(err.String())
}

func TestRPC_New(t *testing.T) {
	rpc := NewRPC()
	logger.Info(rpc)
}

func TestRPC_NewWithPath(t *testing.T) {
	rpc := NewRPCWithPath("../conf")
	logger.Info(rpc)
}

func TestRPC_DefaultRPC(t *testing.T) {
	rpc := DefaultRPC(NewNode("localhost", "8081", "11001")).Https("../conf/certs/tls/tlsca.ca", "../conf/certs/tls/tls_peer.cert", "../conf/certs/tls/tls_peer.priv").Tcert(true, "../conf/certs/sdkcert.cert", "../conf/certs/sdkcert.priv", "../conf/certs/unique.pub", "../conf/certs/unique.priv")
	logger.Info(rpc)
}

func TestRPC_GetNodes(t *testing.T) {
	nodeInfo, err := rpc.GetNodes()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(nodeInfo))
	logger.Info(nodeInfo)
}

/*---------------------------------- contract ----------------------------------*/

func TestRPC_CompileContract(t *testing.T) {
	//nolint
	compileContract("../conf/contract/Accumulator.sol")
}

func TestRPC_ManageContractByVote(t *testing.T) {
	t.Skip()
	source, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")
	bin := `6060604052341561000f57600080fd5b5b6104c78061001f6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680635b6beeb914610049578063e15fe02314610120575b600080fd5b341561005457600080fd5b6100a4600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061023a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e55780820151818401525b6020810190506100c9565b50505050905090810190601f1680156101125780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561012b57600080fd5b6101be600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061034f565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101ff5780820151818401525b6020810190506101e3565b50505050905090810190601f16801561022c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102426103e2565b6000826040518082805190602001908083835b60208310151561027b57805182525b602082019150602081019050602083039250610255565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103425780601f1061031757610100808354040283529160200191610342565b820191906000526020600020905b81548152906001019060200180831161032557829003601f168201915b505050505090505b919050565b6103576103e2565b816000846040518082805190602001908083835b60208310151561039157805182525b60208201915060208101905060208303925061036b565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906103d79291906103f6565b508190505b92915050565b602060405190810160405280600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061043757805160ff1916838001178555610465565b82800160010185558215610465579182015b82811115610464578251825591602001919060010190610449565b5b5090506104729190610476565b5090565b61049891905b8082111561049457600081600090555060010161047c565b5090565b905600a165627a7a723058208ac1d22e128cf8381d7ac66b4c438a6a906ccf5ee583c3a9e46d4cdf7b3f94580029`
	//cr, _ := rpc.CompileContract(string(source))
	ope := bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(bin), "evm", nil)
	contractOpt := bvm.NewProposalCreateOperationForContract(ope)
	payload := bvm.EncodeOperation(contractOpt)
	tx := NewTransaction(privateKey.GetAddress().Hex()).Invoke(contractOpt.Address(), payload).VMType(BVM)
	tx.Sign(privateKey)
	re, err := rpc.ManageContractByVote(tx)
	assert.NoError(t, err)
	result := bvm.Decode(re.Ret)
	assert.True(t, result.Success)
	var proposal bvm.ProposalData
	_ = proto.Unmarshal([]byte(result.Ret), &proposal)
	t.Log(proposal.String())
}

func TestRPC_DeployContract(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0])
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)
}

func TestRPC_DeployContractAsync(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0])
	transaction.Sign(guomiKey)
	asyncHandler := TestAsyncHandler{t: t}
	rpc.DeployContractAsync(transaction, &asyncHandler)
	time.Sleep(3 * time.Second)
	assert.EqualValues(t, true, asyncHandler.IsCalled, "回调未被执行")
}

func TestRPC_InvokeContract(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	contractAddress, err := deployContract(cr.Bin[0], cr.Abi[0])
	ABI, serr := abi.JSON(strings.NewReader(cr.Abi[0]))
	if err != nil {
		t.Error(serr)
		return
	}
	packed, serr := ABI.Pack("add", uint32(1), uint32(2))
	if serr != nil {
		t.Error(serr)
		return
	}
	transaction := NewTransaction(address).Invoke(contractAddress, packed)
	transaction.Sign(privateKey)
	receipt, _ := rpc.InvokeContract(transaction)
	fmt.Println("ret:", receipt.Ret)
}

func TestRPC_InvokeContractAsync(t *testing.T) {
	contractAddress, _ := deployContract(binContract, abiContract)
	ABI, err := abi.JSON(strings.NewReader(abiContract))
	if err != nil {
		t.Error(err)
		return
	}
	packed, err := ABI.Pack("add", uint32(1), uint32(2))
	if err != nil {
		t.Error(err)
		return
	}
	transaction := NewTransaction(address).Invoke(contractAddress, packed)
	transaction.Sign(privateKey)
	asyncHandler := TestAsyncHandler{t: t}
	rpc.InvokeContractAsync(transaction, &asyncHandler)
	time.Sleep(3 * time.Second)
	assert.EqualValues(t, true, asyncHandler.IsCalled, "回调未被执行")
}

func TestRPC_GetCode(t *testing.T) {
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)

	code, err := rpc.GetCode(receipt.ContractAddress)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(code)
}

func TestRPC_GetContractCountByAddr(t *testing.T) {
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	count, err := rpc.GetContractCountByAddr(common.BytesToAddress(newAddress).Hex())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}

func TestRPC_EncryptoMessage(t *testing.T) {
	count, err := rpc.EncryptoMessage(100, 10, "123456")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}

func TestRPC_CheckHmValue(t *testing.T) {
	count, err := rpc.CheckHmValue([]uint64{1, 2}, []string{"123", "456"}, "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}

func TestRPC_DeployContractWithArgs(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator2.sol")
	var arg [32]byte
	copy(arg[:], "test")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).DeployArgs(cr.Abi[0], uint32(10), arg)
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)

	fmt.Println("-----------------------------------")

	ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
	packed, _ := ABI.Pack("getMul")
	transaction1 := NewTransaction(address).Invoke(receipt.ContractAddress, packed)
	transaction1.Sign(privateKey)
	receipt1, err := rpc.InvokeContract(transaction1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("ret:", receipt1.Ret)

	var p0 []byte
	var p1 int64
	var p2 common.Address
	testV := []interface{}{&p0, &p1, &p2}
	fmt.Println(reflect.TypeOf(testV))
	decode(ABI, &testV, "getMul", receipt1.Ret)
	fmt.Println(string(p0), p1, p2.Hex())
}

func TestRPC_DeployContractWithStringArgs(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator2.sol")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).DeployStringArgs(cr.Abi[0], "10", "test")
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)
	fmt.Println("-----------------------------------")
	ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
	packed, _ := ABI.Encode("getMul")
	invokeTx := NewTransaction(address).Invoke(receipt.ContractAddress, packed)
	invokeTx.Sign(privateKey)
	invokeReceipt, err := rpc.InvokeContract(invokeTx)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("ret:", invokeReceipt.Ret)
	ret, e := ABI.Decode("getMul", common.FromHex(invokeReceipt.Ret))
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("%v\n", ret)
}

func TestRPC_UnpackLog(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator2.sol")
	var arg [32]byte
	copy(arg[:], "test")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).DeployArgs(cr.Abi[0], uint32(10), arg)
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)

	fmt.Println("-----------------------------------")

	ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
	packed, _ := ABI.Pack("getHello")
	transaction1 := NewTransaction(address).Invoke(receipt.ContractAddress, packed)
	transaction1.Sign(privateKey)
	receipt1, err := rpc.InvokeContract(transaction1)
	if err != nil {
		t.Error(err)
		return
	}
	test := struct {
		Addr int64   `abi:"addr1"`
		Msg1 [8]byte `abi:"msg"`
	}{}

	// testLog
	sysErr := ABI.UnpackLog(&test, "sayHello", receipt1.Log[0].Data, receipt1.Log[0].Topics)
	if sysErr != nil {
		t.Error(sysErr)
		return
	}
	msg, sysErr := abi.ByteArrayToString(test.Msg1)
	if sysErr != nil {
		t.Error(sysErr)
		return
	}
	assert.Equal(t, int64(1), test.Addr, "解码失败")
	assert.Equal(t, "test", msg, "解码失败")
}

func TestRPC_SendTx(t *testing.T) {
	guomiKey, _ := gm.GenerateSM2Key()
	pubKey := &account.SM2Key{SM2PrivateKey: guomiKey}
	newAddress := pubKey.GetAddress()

	transaction := NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
	transaction.Sign(pubKey)
	receipt, err := rpc.SendTx(transaction)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(receipt.Ret)
}

// {"method":"tx_sendTransaction","jsonrpc":"2.0","id":1,"namespace":"global","params":[{"from":"0xadadb15c50f457413f63a153df733e39645915ce","nonce":8573472654983015045,"signature":"0x01040739b9f7e43ce8d0956fa3cb8e96ac532052c7bb8940e069427480f2b3a02ec8cf2d3848bb7a19daa2aae621b94612056fdd02f0e299e0d21353ff100ebc4a9430440220734f5dc59df7958ceaec01b60142698ea94f09c10fa622a069a7643332c03b970220046fdf52229511b8053a4c339ef74fa742aa7da71de797af64483938bda20d98",  "simulate":false,"timestamp":1605171211407332000,"to":"0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783","type":"EVM","value":0}]}
// {"method":"tx_sendTransaction","jsonrpc":"2.0","id":1,"namespace":"global","params":[{"from":"0x1200fa001e97ab8ea420a629260aac2b0b5c887a","nonce":3587644803921885232,"signature":"0x01048a3aee19db96a3d276f8081bd1b5a4a60d02d09becbcdf873dff8f032aa5333fb3852da16df4067a90712041ca77ea660c715f6a9fc54de40b48237d367f1ed53045022100a142884bada29a675bfd8114a664698d91fb5635f40898d55b9d2795430cffff02206e4a2402095f2a348293d47db59ed6b70ac1b412d0d979ff42dc16b66e5fe90a","simulate":false,"timestamp":1605171324319517000,"to":"0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783","type":"EVM","value":0}]}

// {"params":[{"from":"0x4fd0258745a5a840af1ea8909548a2476893c9bd","nonce":6726120488343232460,"signature":"0x010429ea047e2d2a86bd5fe8c3da9c191ad29b69267edf98ecb91807bc6199cd7f8927212420a9fe73d6f8fc913551335adb10ee2086e302035fd52ca5ef94ce1cf23045022100c8d4cd6525df46dece359d20c2497197a8d9491649466f8f52c96aec2464486c022039bd6370cbd607e34254e622eaea45109426ea8a679f344ec60e6aa778b529b0","simulate":false,"timestamp":1605171388312113000,"to":"0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783","type":"EVM","value":0}]}
// {           "from":"0x4fd0258745a5a840af1ea8909548a2476893c9bd","to":"0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783","value":"0x0","payload":"","signature":"0x010429ea047e2d2a86bd5fe8c3da9c191ad29b69267edf98ecb91807bc6199cd7f8927212420a9fe73d6f8fc913551335adb10ee2086e302035fd52ca5ef94ce1cf23045022100c8d4cd6525df46dece359d20c2497197a8d9491649466f8f52c96aec2464486c022039bd6370cbd607e34254e622eaea45109426ea8a679f344ec60e6aa778b529b0","timestamp":1605171388312113000,"simulate":false,"nonce":6726120488343232460,"extra":"","type":"EVM","opcode":0,"snapshotId":"","extraIdInt64":null,"extraIdString":null,"cName":""}

// {"params":[{"from":"0x7573dd00ccfa47258fc0bf92361c459cadc3e7a4","nonce":4316686915858542454,"signature":"0x01043cb005963cfc56fc2a538b6d4d4883e1ecb31220db12d4273dbeb90e0a4855d519bf33978173391ee916d96bced484d61e4f342d1cd6e0e2d12798d198ab6a60304502203ab9bc9d165a057220897ad2156a1a40309b267cf05ad3670cd95e138750677a022100d67266f80cc92b8128addce6a64b0c11ccf658218dd801e675b02bf8e58136ff","simulate":false,"timestamp":1605171476123870000,"to":"0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783","type":"EVM","value":0}]}
// {           "from":"0x7573dd00ccfa47258fc0bf92361c459cadc3e7a4","to":"0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783","value":"0x0","payload":"","signature":"0x01043cb005963cfc56fc2a538b6d4d4883e1ecb31220db12d4273dbeb90e0a4855d519bf33978173391ee916d96bced484d61e4f342d1cd6e0e2d12798d198ab6a60304502203ab9bc9d165a057220897ad2156a1a40309b267cf05ad3670cd95e138750677a022100d67266f80cc92b8128addce6a64b0c11ccf658218dd801e675b02bf8e58136ff","timestamp":1605171476123870000,"simulate":false,"nonce":4316686915858542454,"extra":"","type":"EVM","opcode":0,"snapshotId":"","extraIdInt64":null,"extraIdString":null,"cName":""}
func TestRPC_SendED25519Tx(t *testing.T) {
	accountJSON, _ := account.NewAccountED25519("12345678")
	t.Log("account", accountJSON)
	ekey, err := account.GenKeyFromAccountJson(accountJSON, "12345678")
	assert.Nil(t, err)
	newAddress := ekey.(*account.ED25519Key).GetAddress()

	transaction := NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
	transaction.Sign(ekey)
	receipt, err := rpc.SendTx(transaction)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, err)
	fmt.Println(receipt.Ret)
}

func TestRPC_SendTx_With_ExtraID(t *testing.T) {
	addr := guomiKey.GetAddress()
	transaction := NewTransaction(hex.EncodeToString(addr[:])).Transfer(address, int64(0))
	transaction.SetExtraIDInt64(123, 456)
	transaction.SetExtraIDString("abc")
	transaction.Sign(guomiKey)
	receipt, err := rpc.SendTx(transaction)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(receipt.Ret)
}

func TestRPC_SM2Account(t *testing.T) {
	//accountJson, _ := account.NewAccountSm2("")
	strAcc := `{"address":"0x8485147cbf02dec93ee84f81824a3b60e355f5cd","publicKey":"04a1b4c82a2a13e15a11e3ee9316504de0c3b54d46f5c189ae42603c9cd07a50fdca2ac35d0ceef4a8466ccb182f52403d9a58b573e1bf6fd4f52c31493bf7241b","privateKey":"f67136bf3caa4197a1cfaf38a5392ff94dae91bda700f8898b11cf49891a47bb","privateKeyEncrypted":false}`
	key, _ := account.NewAccountSm2FromAccountJSON(strAcc, "")
	//fmt.Println(accountJson)
	//key, _ := account.NewAccountSm2FromAccountJSON(accountJson, "")
	pubKey, _ := key.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Transfer(address, int64(0))
	transaction.Sign(key)
	receipt, err := rpc.SendTx(transaction)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 66, len(receipt.TxHash))

	accountJSON, _ := account.NewAccountSm2("12345678")
	aKey, syserr := account.NewAccountSm2FromAccountJSON(accountJSON, "12345678")
	if syserr != nil {
		t.Error(syserr)
	}
	newAddress2 := aKey.GetAddress()

	transaction1 := NewTransaction(newAddress2.Hex()).Transfer(address, int64(0))
	transaction1.Sign(aKey)
	receipt1, err := rpc.SendTx(transaction1)
	if err != nil {
		t.Error(err)
		return
	}
	//assert.EqualValues(t, transaction1.GetTransactionHash(DefaultTxGasLimit), receipt1.TxHash)
	assert.EqualValues(t, 66, len(receipt1.TxHash))
}

func TestRPC_SendTxAsync(t *testing.T) {
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Transfer(address, int64(0))
	transaction.Sign(guomiKey)
	asyncHandler := TestAsyncHandler{t: t}
	rpc.SendTxAsync(transaction, &asyncHandler)
	time.Sleep(3 * time.Second)
	assert.EqualValues(t, true, asyncHandler.IsCalled, "回调未被执行")
}

// maintain contract by opcode 1
func TestRPC_MaintainContract(t *testing.T) {
	t.Skip("solc")
	contractOriginFile := "../conf/contract/Accumulator.sol"
	contractUpdateFile := "../conf/contract/AccumulatorUpdate.sol"
	compileOrigin, _ := compileContract(contractOriginFile)
	compileUpdate, _ := compileContract(contractUpdateFile)
	contractAddress, err := deployContract(compileOrigin.Bin[0], compileOrigin.Abi[0])
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("contractAddress:", contractAddress)

	// test invoke before update
	ABIBefore, serr := abi.JSON(strings.NewReader(compileOrigin.Abi[0]))
	assert.Nil(t, serr)
	packed, serr := ABIBefore.Pack("add", uint32(11), uint32(1))
	if serr != nil {
		t.Error(serr)
		return
	}
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transactionInvokeBe := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
	transactionInvokeBe.Sign(guomiKey)
	receiptBe, err := rpc.InvokeContract(transactionInvokeBe)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(receiptBe.Ret)

	var result1 uint32
	decode(ABIBefore, &result1, "add", receiptBe.Ret)
	fmt.Println(result1)

	fmt.Println("-----------------------------")

	transactionUpdate := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(1, contractAddress, compileUpdate.Bin[0])
	//transactionUpdate, err := NewMaintainTransaction(guomiKey.GetAddress(), contractAddress, compileUpdate.Bin[0], 1, EVM)
	transactionUpdate.Sign(guomiKey)
	receiptUpdate, err := rpc.MaintainContract(transactionUpdate)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(receiptUpdate.ContractAddress)

	// test invoke after update
	ABI, serr := abi.JSON(strings.NewReader(compileUpdate.Abi[0]))
	if serr != nil {
		t.Error(err)
		return
	}
	packed2, serr := ABI.Pack("addUpdate", uint32(1), uint32(2))
	if serr != nil {
		t.Error(serr)
		return
	}
	transactionInvoke := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed2)
	//transactionInvoke, err := NewInvokeTransaction(guomiKey.GetAddress(), contractAddress, common.ToHex(packed2), false, EVM)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	transactionInvoke.Sign(guomiKey)
	receiptInvoke, err := rpc.InvokeContract(transactionInvoke)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(receiptInvoke.Ret)
	var result2 uint32
	decode(ABI, &result2, "addUpdate", receiptInvoke.Ret)
	fmt.Println(result2)
}

// maintain contract by opcode 2 and 3
func TestRPC_MaintainContract2(t *testing.T) {
	contractAddress, _ := deployContract(binContract, abiContract)
	ABI, _ := abi.JSON(strings.NewReader(abiContract))
	// invoke first
	packed, _ := ABI.Pack("getSum")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction1 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
	//transaction1, _ := NewInvokeTransaction(guomiKey.GetAddress(), contractAddress, common.ToHex(packed), false, EVM)
	transaction1.Sign(guomiKey)
	receipt1, err := rpc.InvokeContract(transaction1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("invoke first:", receipt1.Ret)

	// freeze contract
	transactionFreeze := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(2, contractAddress, "")
	//transactionFreeze, _ := NewMaintainTransaction(guomiKey.GetAddress(), contractAddress, "", 2, EVM)
	transactionFreeze.Sign(guomiKey)
	receiptFreeze, err := rpc.MaintainContract(transactionFreeze)
	assert.Nil(t, err)
	fmt.Println(receiptFreeze.TxHash)
	status, err := rpc.GetContractStatus(contractAddress)
	assert.Nil(t, err)
	fmt.Println("contract status >>", status)

	// invoke after freeze
	transaction2 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
	//transaction2, _ := NewInvokeTransaction(guomiKey.GetAddress(), contractAddress, common.ToHex(packed), false, EVM)
	transaction2.Sign(guomiKey)
	receipt2, err := rpc.InvokeContract(transaction2)
	if err != nil {
		fmt.Println("invoke second receipt2 is null ", receipt2 == nil)
		fmt.Println(err)
	}

	// unfreeze contract
	transactionUnfreeze := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(3, contractAddress, "")
	//transactionUnfreeze, _ := NewMaintainTransaction(guomiKey.GetAddress(), contractAddress, "", 3, EVM)
	transactionUnfreeze.Sign(guomiKey)
	receiptUnFreeze, err := rpc.MaintainContract(transactionUnfreeze)
	assert.Nil(t, err)
	fmt.Println(receiptUnFreeze.TxHash)
	status, _ = rpc.GetContractStatus(contractAddress)
	fmt.Println("contract status >>", status)

	// invoke after unfreeze
	transaction3 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
	//transaction3, _ := NewInvokeTransaction(guomiKey.GetAddress(), contractAddress, common.ToHex(packed), false, EVM)
	transaction3.Sign(guomiKey)
	receipt3, err := rpc.InvokeContract(transaction3)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("invoke third:", receipt3.Ret)
}

func TestRPC_GetContractStatus(t *testing.T) {
	//t.Skip("the node can get the account")
	contractAddress, _ := deployContract(binContract, abiContract)
	statu, err := rpc.GetContractStatus(contractAddress)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(statu)
}

func TestRPC_GetContractStatusByName(t *testing.T) {
	t.Skip("set contract name `HashContract`")
	dateTime, err := rpc.GetContractStatusByName("HashContract")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dateTime)
}

func TestRPC_GetCreator(t *testing.T) {
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)

	accountAddress, err := rpc.GetCreator(receipt.ContractAddress)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(accountAddress)
}

func TestRPC_GetCreatorByName(t *testing.T) {
	t.Skip("set contract name `HashContract`")
	dateTime, err := rpc.GetCreatorByName("HashContract")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dateTime)
}

func TestRPC_GetCreateTime(t *testing.T) {
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)

	dateTime, err := rpc.GetCreateTime(receipt.ContractAddress)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dateTime)
}

func TestRPC_GetCreateTimeByName(t *testing.T) {
	t.Skip("set contract name `HashContract`")
	dateTime, err := rpc.GetCreateTimeByName("HashContract")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dateTime)
}

func TestRPC_GetDeployedList(t *testing.T) {
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	list, err := rpc.GetDeployedList(common.BytesToAddress(newAddress).Hex())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(list))
}

func TestRPC_InvokeContractReturnHash(t *testing.T) {
	t.Skip("pressure test, do not put this test in CI")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	contractAddress, err := deployContract(cr.Bin[0], cr.Abi[0])
	ABI, serr := abi.JSON(strings.NewReader(cr.Abi[0]))
	if err != nil {
		t.Error(serr)
		return
	}
	packed, serr := ABI.Pack("add", uint32(1), uint32(2))
	if serr != nil {
		t.Error(serr)
		return
	}
	transaction := NewTransaction(address).Invoke(contractAddress, packed)
	transaction.Sign(privateKey)
	tt := time.After(1 * time.Minute)
	counter := 0
	for {
		_, err = rpc.InvokeContractReturnHash(transaction)
		if err != nil {
			t.Error(err)
			return
		}
		select {
		case <-tt:
			fmt.Println(counter)
			return
		default:
			counter++
		}
	}
}

func TestRPC_InvokeContract2(t *testing.T) {
	to := `0x12345678901234567890123456789012345678901234567890`
	rawABI := `[{"constant":false,"inputs":[{"name":"num1","type":"uint32"},{"name":"num2","type":"uint32"}],"name":"add","outputs":[{"name":"","type":"uint32"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"getSum","outputs":[{"name":"","type":"uint32"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getHello","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"increment","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"addr1","type":"address"},{"indexed":false,"name":"msg","type":"bytes32"}],"name":"sayHello","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"msg","type":"bytes32"},{"indexed":false,"name":"sum","type":"uint32"}],"name":"saySum","type":"event"}]`

	t.Run("normal input", func(t *testing.T) {
		tx1 := NewTransaction(address).InvokeContract(to, rawABI, "add", "111", "111")

		ABI, err := abi.JSON(strings.NewReader(rawABI))
		assert.NoError(t, err)
		payload, err := ABI.Pack("add", uint32(111), uint32(111))
		tx2 := NewTransaction(address).Invoke(to, payload)
		assert.NoError(t, err)
		assert.Equal(t, tx2.payload, tx1.payload)
	})

	t.Run("error input", func(t *testing.T) {
		errTx := NewTransaction(address).InvokeContract(to, strings.Replace(rawABI, "[", "{", -1), "add", "111", "111")
		assert.Nil(t, errTx)

		errTx = NewTransaction(address).InvokeContract(to, rawABI, "123", "111", "111")
		assert.Nil(t, errTx)
	})
}

/*---------------------------------- archive ----------------------------------*/

func TestRPC_Snapshot(t *testing.T) {
	t.Skip()
	res, err := rpc.Snapshot(1)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_QuerySnapshotExist(t *testing.T) {
	t.Skip()
	res, err := rpc.QuerySnapshotExist("0x5d86cce7e537cd0e0346468889801196")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_CheckSnapshot(t *testing.T) {
	t.Skip()
	res, err := rpc.CheckSnapshot("0x5d86cce7e537cd0e0346468889801196")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_Archive(t *testing.T) {
	t.Skip()
	res, err := rpc.Archive("0x5d86cce7e537cd0e0346468889801196", false)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_Restore(t *testing.T) {
	t.Skip()
	res, err := rpc.Restore("0x5d86cce7e537cd0e0346468889801196", false)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}
func TestRPC_RestoreAll(t *testing.T) {
	t.Skip()
	res, err := rpc.RestoreAll(false)
	if err != nil {
		t.Error(err)

	}

	fmt.Println(res)
}

func TestRPC_QueryArchiveExist(t *testing.T) {
	t.Skip()
	res, err := rpc.QueryArchiveExist("0x5d86cce7e537cd0e0346468889801196")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_Pending(t *testing.T) {
	t.Skip()
	res, err := rpc.QueryArchiveExist("0x5d86cce7e537cd0e0346468889801196")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

/*---------------------------------- node ----------------------------------*/

func TestRPC_GetNodeHash(t *testing.T) {
	hash, err := rpc.GetNodeHash()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(hash)
}

func TestRPC_GetNodeHashById(t *testing.T) {
	id := 1
	hash, err := rpc.GetNodeHashByID(id)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(hash)
}

func TestRPC_DeleteNodeVP(t *testing.T) {
	t.Skip("do not delete VP in CI")
	hash1, _ := rpc.GetNodeHashByID(1)
	success, _ := rpc.DeleteNodeVP(hash1)
	assert.Equal(t, true, success)

	hash11, _ := rpc.GetNodeHashByID(1)
	fmt.Println(hash11)
}

func TestRPC_DeleteNodeNVP(t *testing.T) {
	t.Skip("do not delete NVP in CI")
	hash1, _ := rpc.GetNodeHashByID(1)
	success, _ := rpc.DeleteNodeNVP(hash1)
	assert.Equal(t, true, success)
}

func TestRPC_GetNodeStates(t *testing.T) {
	states, err := rpc.GetNodeStates()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 4, len(states))
}

/*---------------------------------- block ----------------------------------*/

func TestRPC_GetLatestBlock(t *testing.T) {
	block, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(block)
}

func TestRPC_GetBlocks(t *testing.T) {
	latestBlock, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	blocks, err := rpc.GetBlocks(latestBlock.Number-1, latestBlock.Number, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(blocks)
}

func TestRPC_GetBlocksWithLimit(t *testing.T) {
	t.Skip()
	latestBlock, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	metadata := &Metadata{
		PageSize: 5,
	}

	pageResult, err := rpc.GetBlocksWithLimit(latestBlock.Number-1, latestBlock.Number, true, metadata)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(pageResult)
}

func TestRPC_GetBlockByHash(t *testing.T) {
	latestBlock, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	block, err := rpc.GetBlockByHash(latestBlock.Hash, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(block)
}

func TestRPC_GetBatchBlocksByHash(t *testing.T) {
	latestBlock, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	blocks, err := rpc.GetBatchBlocksByHash([]string{latestBlock.Hash}, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(blocks)
}

func TestRPC_GetBlockByNumber(t *testing.T) {
	latestBlock, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	//nolint
	rpc.GetBlockByNumber("latest", false)
	block, err := rpc.GetBlockByNumber(latestBlock.Number, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(block)
}

func TestRPC_GetBatchBlocksByNumber(t *testing.T) {
	latestBlock, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	blocks, err := rpc.GetBatchBlocksByNumber([]uint64{latestBlock.Number}, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(blocks)
}

func TestRPC_GetAvgGenTimeByBlkNum(t *testing.T) {
	block, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	avgTime, err := rpc.GetAvgGenTimeByBlockNum(block.Number-2, block.Number)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(avgTime)
}

func TestRPC_GetBlockByTime(t *testing.T) {
	start := time.Now().Unix() - 1
	end := time.Now().Unix()
	blockInterval, err := rpc.GetBlocksByTime(uint64(start), uint64(end))
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(blockInterval)
}

func TestRPC_QueryTPS(t *testing.T) {
	tpsInfo, err := rpc.QueryTPS(1, 1778959217012956575)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(tpsInfo)
}

func TestRPC_GetGenesisBlock(t *testing.T) {
	blkNum, err := rpc.GetGenesisBlock()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, true, strings.HasPrefix(blkNum, "0x"))
}

func TestRPC_GetChainHeight(t *testing.T) {
	blkNum, err := rpc.GetChainHeight()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, true, strings.HasPrefix(blkNum, "0x"))
}

/*---------------------------------- transaction ----------------------------------*/

func TestRPC_GetTransactions(t *testing.T) {
	block, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	txs, err := rpc.GetTransactionsByBlkNum(block.Number-1, block.Number)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txs[0].Invalid)
}

func TestRPC_GetTransactionsWithLimit(t *testing.T) {
	t.Skip()
	block, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	metadata := &Metadata{
		PageSize: 1,
		Bookmark: &Bookmark{
			BlockNumber: 1,
			TxIndex:     0,
		},
		Backward: false,
	}

	pageResult, err := rpc.GetTransactionsByBlkNumWithLimit(block.Number-1, block.Number, metadata)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pageResult)
}

func TestRPC_GetDiscardTx(t *testing.T) {
	txs, err := rpc.GetDiscardTx()
	if err != nil {
		//t.Error(err)
		return
	}
	fmt.Println(len(txs))
	if len(txs) > 0 {
		fmt.Println(txs[len(txs)-1].Hash)
	}
}

func TestRPC_GetTransactionByHash(t *testing.T) {
	//t.Skip()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)
	fmt.Println("txhash:", receipt.TxHash)

	hash := receipt.TxHash
	tx, err := rpc.GetTransactionByHash(hash)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tx.Hash)
	assert.Equal(t, receipt.TxHash, tx.Hash)
}

func TestRPC_GetBatchTxByHash(t *testing.T) {
	//t.Skip()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction1 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	transaction1.Sign(guomiKey)
	receipt1, _ := rpc.DeployContract(transaction1)
	fmt.Println("txhash1:", receipt1.TxHash)

	// 模拟一个无法查询到交易的hash
	txHash2 := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	txhashes := make([]string, 0)
	txhashes = append(txhashes, receipt1.TxHash, txHash2)

	txs, err := rpc.GetBatchTxByHash(txhashes)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(txs))
	fmt.Println(txs[0].Hash, txs[1].Hash)
	assert.Equal(t, receipt1.TxHash, txs[0].Hash)
	assert.Equal(t, txHash2, txs[1].Hash)
}

func TestRPC_GetTxByBlkHashAndIdx(t *testing.T) {
	block, _ := rpc.GetLatestBlock()
	info, err := rpc.GetTxByBlkHashAndIdx(block.Hash, 0)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(info)
	assert.EqualValues(t, 66, len(info.Hash))
}

func TestRPC_GetTxByBlkNumAndIdx(t *testing.T) {
	block, _ := rpc.GetLatestBlock()
	info, err := rpc.GetTxByBlkNumAndIdx(block.Number, 0)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(info)
	assert.EqualValues(t, 66, len(info.Hash))
}

func TestRPC_GetTxAvgTimeByBlockNumber(t *testing.T) {
	block, _ := rpc.GetLatestBlock()
	time, err := rpc.GetTxAvgTimeByBlockNumber(block.Number-2, block.Number)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(time)
}

func TestRPC_GetBatchReceipt(t *testing.T) {
	//t.Skip()
	block, _ := rpc.GetLatestBlock()
	trans, _ := rpc.GetTransactionsByBlkNum(block.Number-2, block.Number)
	// 模拟一个无法查询到回执的hash
	txHash2 := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	hashes := []string{trans[0].Hash, txHash2}
	txs, err := rpc.GetBatchReceipt(hashes)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, len(txs))
}

func TestRPC_GetTransactionsCountByTime(t *testing.T) {
	count, err := rpc.GetTransactionsCountByTime(1, uint64(time.Now().UnixNano()))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}

func TestRPC_GetTxCountByContractAddr(t *testing.T) {
	cAddress, _ := deployContract(binContract, abiContract)
	ABI, _ := abi.JSON(strings.NewReader(abiContract))
	packed, _ := ABI.Pack("getSum")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction.Sign(guomiKey)
	//nolint
	rpc.InvokeContract(transaction)
	transaction2 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction2.Sign(guomiKey)
	//nolint
	rpc.InvokeContract(transaction2)

	block, _ := rpc.GetLatestBlock()
	count, err := rpc.GetTxCountByContractAddr(block.Number-1, block.Number, cAddress, false)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 2, count.Count)
}

func TestRPC_GetTransactionsCountByMethodID(t *testing.T) {
	cAddress, _ := deployContract(binContract, abiContract)
	abi, _ := abi.JSON(strings.NewReader(abiContract))
	methodID := string(abi.Constructor.Id())

	block, _ := rpc.GetLatestBlock()
	count, err := rpc.GetTransactionsCountByMethodID(block.Number-1, block.Number, cAddress, methodID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
}

func TestRPC_GetTxByTime(t *testing.T) {
	//t.Skip("the length of result is too long")
	infos, err := rpc.GetTxByTime(1, uint64(time.Now().UnixNano()))
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(infos)
	assert.EqualValues(t, true, len(infos) > 0)
}

func TestRPC_GetTxByTimeWithLimit(t *testing.T) {
	t.Skip()
	metadata := &Metadata{
		PageSize: 1,
		Bookmark: &Bookmark{
			BlockNumber: 1,
			TxIndex:     0,
		},
		Backward: false,
	}

	pageResult, err := rpc.GetTxByTimeWithLimit(1, uint64(time.Now().UnixNano()), metadata)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pageResult)
}

func TestRPC_GetDiscardTransactionsByTime(t *testing.T) {
	//t.Skip()
	infos, err := rpc.GetDiscardTransactionsByTime(1, uint64(time.Now().UnixNano()))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(infos)
}

func TestRPC_GetNextPageTxs(t *testing.T) {
	//t.Skip("hyperchain snapshot will case error")
	cAddress, _ := deployContract(binContract, abiContract)
	ABI, _ := abi.JSON(strings.NewReader(abiContract))
	packed, _ := ABI.Pack("getSum")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction.Sign(guomiKey)
	//nolint
	rpc.InvokeContract(transaction)
	transaction2 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction2.Sign(guomiKey)
	//nolint
	rpc.InvokeContract(transaction2)

	block, _ := rpc.GetLatestBlock()

	infos, err := rpc.GetNextPageTxs(block.Number-10, 0, 1, block.Number, 0, 10, false, cAddress)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 3, len(infos))

	t.Skip()
	txs, err := rpc.GetNextPageTxs(block.Number-10, 0, 1, block.Number, 0, 10, false, "")
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 10, len(txs))
}

func TestRPC_GetPrevPageTxs(t *testing.T) {
	//t.Skip("hyperchain snapshot will case error")
	cAddress, _ := deployContract(binContract, abiContract)
	ABI, _ := abi.JSON(strings.NewReader(abiContract))
	packed, _ := ABI.Pack("getSum")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction.Sign(guomiKey)
	//nolint
	rpc.InvokeContract(transaction)
	transaction2 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction2.Sign(guomiKey)
	//nolint
	rpc.InvokeContract(transaction2)

	block, _ := rpc.GetLatestBlock()

	infos, err := rpc.GetPrevPageTxs(block.Number, 0, 1, block.Number, 0, 10, false, cAddress)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 2, len(infos))

	t.Skip()
	txs, err := rpc.GetPrevPageTxs(block.Number-10, 0, 1, block.Number, 0, 10, false, "")
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 10, len(txs))
}

func TestRPC_GetBlkTxCountByHash(t *testing.T) {
	block, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	count, err := rpc.GetBlkTxCountByHash(block.Hash)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
}

func TestRPC_GetBlkTxCountByNumber(t *testing.T) {
	block, err := rpc.GetLatestBlock()
	hex := "0x" + strconv.FormatUint(block.Number, 16)
	fmt.Println("=====", block, hex)
	if err != nil {
		t.Error(err)
		return
	}
	count, err := rpc.GetBlkTxCountByNumber(hex)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
}

//func TestRPC_GetSignHash(t *testing.T) {
//	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
//	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
//	newAddress := h[12:]
//	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex())
//	transaction.to = common.BytesToAddress(newAddress).Hex()
//	transaction.Sign(guomiKey)
//	count, err := rpc.GetSignHash(transaction)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	fmt.Println(count)
//}

func TestRPC_GetTxCount(t *testing.T) {
	txCount, err := rpc.GetTxCount()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txCount.Count)
}

/*---------------------------------- cert ----------------------------------*/

//func TestRPC_GetTCert(t *testing.T) {
//	tCert, err := rpc.GetTCert(1)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	fmt.Println(tCert)
//}

/*---------------------------------- account ----------------------------------*/

func TestRPC_GetBalance(t *testing.T) {
	account := "0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd"
	balance, err := rpc.GetBalance(account)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(balance)
}

func TestRPC_GetRoles(t *testing.T) {
	t.Skip()
	account := "0x2a307e1e5b53863242a465bf99ca6e94947da898"
	roles, err := rpc.GetRoles(account)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(roles)
}

func TestRPC_GetAccountsByRole(t *testing.T) {
	t.Skip()
	role := "admin"
	roles, err := rpc.GetAccountsByRole(role)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(roles)
}

/*---------------------------------- config ----------------------------------*/

func TestRPC_GetProposal(t *testing.T) {
	t.Skip()
	proposal, err := rpc.GetProposal()
	if err != nil {
		t.Error(err)
		return
	}
	bytes, _ := json.Marshal(proposal)
	fmt.Println(string(bytes))
}

func TestRPC_GetHosts(t *testing.T) {
	t.Skip()
	data, err := rpc.GetHosts("vp")
	if err != nil {
		t.Error(err)
		return
	}
	bytes, _ := json.Marshal(data)
	fmt.Println(string(bytes))
}

func TestRPC_GetVSet(t *testing.T) {
	t.Skip()
	data, err := rpc.GetVSet()
	if err != nil {
		t.Error(err)
		return
	}
	bytes, _ := json.Marshal(data)
	fmt.Println(string(bytes))
}

func TestRPC_GetConfig(t *testing.T) {
	t.Skip()
	config, err := rpc.GetConfig()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(config)
}

var accountJsons = []string{
	//ecdsa
	`{"address":"0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd","version":"4.0", "algo":"0x03","publicKey":"0400ddbadb932a0d276e257c6df50599a425804a3743f40942d031f806bf14ab0c57aed6977b1ad14646672f9b9ce385f2c98c4581267b611f48f4b7937de386ac","privateKey":"16acbf6b4f09a476a35ebd4c01e337238b5dceceb6ff55ff0c4bd83c4f91e11b"}`,
	`{"address":"0x6201cb0448964ac597faf6fdf1f472edf2a22b89","version":"4.0", "algo":"0x03","publicKey":"04e482f140d70a1b8ec8185cc699db5b391ea5a7b8e93e274b9f706be9efdaec69542eb32a61421ba6219230b9cf87bf849fa01c1d10a8d298cbe3dcfa5954134c","privateKey":"21ff03a654c939f0c9b83e969aaa9050484aa4108028094ee2e927ba7e7d1bbb"}`,
	`{"address":"0xb18c8575e3284e79b92100025a31378feb8100d6","version":"4.0", "algo":"0x03","publicKey":"042169a7260acaff308228579aab2a2c6b3a790922c6a6b58b218cdd7ce0b1db0fbfa6f68737a452010b9d138187b8321288cae98f07fc758bb67bb818292cab9b","privateKey":"aa9c83316f68c17bcc21cf20a4733ae2b2bf76ad1c745f634c0ebf7d5094500e"}`,
	`{"address":"0xe93b92f1da08f925bdee44e91e7768380ae83307","version":"4.0","algo":"0x03","publicKey":"047196daf5d4d1fe339da58e2fe0543bbfec9a464b76546f180facdcc56315b8eeeca50474100f15fb17606695ce24a1f8e5a990600c1c4ea9787ba4dd65c8ce3e","privateKey":"8cdfbe86deb690e331453a84a98c956f0422dd1e783c3a02aed9180b1f4516a9"}`,
	//sm2
	`{"address":"0xfbca6a7e9e29728773b270d3f00153c75d04e1ad","version":"4.0","algo":"0x13","publicKey":"049c330d0aea3d9c73063db339b4a1a84d1c3197980d1fb9585347ceeb40a5d262166ee1e1cb0c29fd9b2ef0e4f7a7dfb1be6c5e759bf411c520a616863ee046a4","privateKey":"5f0a3ea6c1d3eb7733c3170f2271c10c1206bc49b6b2c7e550c9947cb8f098e3"}`,
	`{"address":"0x856e2b9a5fa82fd1b031d1ff6863864dbac7995d","publicKey":"047ea464762c333762d3be8a04536b22955d97231062442f81a3cff46cb009bbdbb0f30e61ade5705254d4e4e0c0745fb3ba69006d4b377f82ecec05ed094dbe87","privateKey":"71b9acc4ee2b32b3d2c79b5abe9e118e5f73765aee5e7755d6aa31f12945036d"}`,
}

func TestRPC_GetAllRoles(t *testing.T) {
	t.Skip()
	roles, err := rpc.GetAllRoles()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(roles)
}

func TestRPC_IsRoleExist(t *testing.T) {
	t.Skip()
	exist, err := rpc.IsRoleExist("admin")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(exist)
}

func TestRPC_GetAddressByName(t *testing.T) {
	t.Skip()
	addr, err := rpc.GetAddressByName("HashContract")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(addr)
}

func TestRPC_GetNameByAddress(t *testing.T) {
	t.Skip()
	name, err := rpc.GetNameByAddress("0x0000000000000000000000000000000000ffff01")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(name)
}

func TestRPC_GetAllCNS(t *testing.T) {
	t.Skip()
	all, err := rpc.GetAllCNS()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(all)
}

func TestRPC_AddRoleForNode(t *testing.T) {
	t.Skip()
	//err := rpc.AddRoleForNode("0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd", "accountManager")
	//err := rpc.AddRoleForNode("0x6201cb0448964ac597faf6fdf1f472edf2a22b89", "accountManager")
	rpc.SetAccount(privateKey)
	err := rpc.AddRoleForNode(privateKey.GetAddress().Hex(), "accountManager")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRPC_GetRoleFromNode(t *testing.T) {
	t.Skip()
	//roles, err := rpc.GetRoleFromNode("0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd")
	roles, err := rpc.GetRoleFromNode("0x6201cb0448964ac597faf6fdf1f472edf2a22b89")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(roles)
}

func TestRPC_GetAddressFromNode(t *testing.T) {
	t.Skip()
	roles, err := rpc.GetAddressFromNode("accountManager")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(roles)
}

func TestRPC_GetAllRolesFromNode(t *testing.T) {
	t.Skip()
	roles, err := rpc.GetAllRolesFromNode()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(roles)
}

func TestRPC_DeleteRoles(t *testing.T) {
	t.Skip()
	err := rpc.DeleteRoleFromNode("0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd", "accountManager")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRPC_SetRulesInNode(t *testing.T) {
	t.Skip()
	rule := &InspectorRule{
		Name:            "rule",
		ID:              1,
		AllowAnyone:     false,
		AuthorizedRoles: []string{"accountManager"},
		Method:          []string{"account_*"},
	}
	err := rpc.SetRulesInNode([]*InspectorRule{rule})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRPC_GetRulesFromNode(t *testing.T) {
	t.Skip()
	rpc.SetAccount(privateKey)
	rules, err := rpc.GetRulesFromNode()
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(rules)
	t.Log(string(marshal))
}

func TestRPC_SetAccount(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	rpc.SetAccount(key)
	balance, stdError := rpc.GetBalance(privateKey.GetAddress().Hex())
	t.Log(stdError)
	t.Log(balance)
}

/**************************** self function ******************************/

func compileContract(path string) (*CompileResult, error) {
	contract, _ := common.ReadFileAsString(path)
	cr, err := rpc.CompileContract(contract)
	if err != nil {
		logger.Error("can not get compile return, ", err.String())
		return nil, err
	}
	fmt.Println("abi:", cr.Abi[0])
	fmt.Println("bin:", cr.Bin[0])
	fmt.Println("type:", cr.Types[0])

	return cr, err
}

func decode(contractAbi abi.ABI, v interface{}, method string, ret string) (result interface{}) {
	if err := contractAbi.UnpackResult(v, method, ret); err != nil {
		logger.Error(NewSystemError(err).String())
	}
	result = v
	return result
}

func deployContract(bin, abi string, params ...interface{}) (string, StdError) {
	var transaction *Transaction
	var err StdError
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	if len(params) == 0 {

		transaction = NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(bin)
		//transaction, err = NewDeployTransaction(guomiKey.GetAddress(), bin, false, EVM)
	} else {
		transaction = NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(bin).DeployArgs(abi, params)
		//transaction, err = NewDeployTransactionWithArgs(guomiKey.GetAddress(), bin, false, EVM, abi, params)
	}
	transaction.Sign(guomiKey)
	txReceipt, err := rpc.DeployContract(transaction)
	if err != nil {
		logger.Error(err)
	}
	return txReceipt.ContractAddress, nil
}

//func Test_Bass(t *testing.T) {
//	contract, _ := common.ReadFileAsString("../conf/contract/Digitalpoint.sol")
//	cr, _ := rpc.CompileContract(contract)
//	abis := cr.Abi[0]
//	bin := cr.Bin[0]
//
//	deployTx := NewTransaction(guomiKey.GetAddress()).Deploy(bin)
//	deployTx.Sign(guomiKey)
//	deployRe, _ := rpc.DeployContract(deployTx)
//	contractAddress := deployRe.ContractAddress
//
//	ABI, _ := abi.JSON(strings.NewReader(abis))
//	packed, err := ABI.Pack("newMarket", "a", "aa", "1")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	invokeTx := NewTransaction(guomiKey.GetAddress()).Invoke(contractAddress, packed)
//	invokeTx.Sign(guomiKey)
//	invokeRe, _ := rpc.InvokeContract(invokeTx)
//	fmt.Println(invokeRe.Ret)
//}

func Test_TypeCheck(t *testing.T) {
	t.Skip("solc")
	contract, _ := common.ReadFileAsString("../conf/contract/TypeCheck.sol")
	cr, _ := rpc.CompileContract(contract)
	abiStr := cr.Abi[0]
	bin := cr.Bin[0]
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	deployTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(bin)
	deployTx.Sign(guomiKey)
	deployRe, _ := rpc.DeployContract(deployTx)
	contractAddress := deployRe.ContractAddress

	ABI, _ := abi.JSON(strings.NewReader(abiStr))

	// invoke fun1
	{
		var data32 [32]byte
		copy(data32[:], "data32")
		var data8 [8]byte
		copy(data8[:], "byte8")
		packed1, _ := ABI.Pack("fun1", []byte("data1"), data32, data8)
		invokeTx1 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed1)
		invokeTx1.Sign(guomiKey)
		invokeRe1, _ := rpc.InvokeContract(invokeTx1)

		var p0 []byte
		var p1 [32]byte
		var p2 [8]byte
		testV := []interface{}{&p0, &p1, &p2}
		if err := ABI.UnpackResult(&testV, "fun1", invokeRe1.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(string(p0), string(p1[:]), string(p2[:]))
	}

	// invoke fun2
	{
		bigInt1 := big.NewInt(-100001)
		bigInt2 := big.NewInt(-1000001)
		bigInt3 := big.NewInt(10000001)
		int1 := int64(-10001)
		int2 := int8(101)
		packed, _ := ABI.Pack("fun2", bigInt1, bigInt2, bigInt3, int1, int2)
		invokeTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
		invokeTx.Sign(guomiKey)
		invokeRe, _ := rpc.InvokeContract(invokeTx)

		var p0 interface{}
		var p1 *big.Int
		var p2 *big.Int
		var p3 interface{}
		var p4 int8
		testV := []interface{}{&p0, &p1, &p2, &p3, &p4}
		if err := ABI.UnpackResult(&testV, "fun2", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0, p1.Int64(), p2, p3, p4)
	}

	// invoke fun3
	{
		bigInt1 := big.NewInt(100001)
		bigInt2 := big.NewInt(1000001)
		bigInt3 := big.NewInt(10000001)
		int1 := uint64(10001)
		int2 := uint8(101)
		packed, _ := ABI.Pack("fun3", bigInt1, bigInt2, bigInt3, int1, int2)
		invokeTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
		invokeTx.Sign(guomiKey)
		invokeRe, _ := rpc.InvokeContract(invokeTx)

		var p0 interface{}
		var p1 *big.Int
		var p2 *big.Int
		var p3 interface{}
		var p4 uint8
		testV := []interface{}{&p0, &p1, &p2, &p3, &p4}
		if err := ABI.UnpackResult(&testV, "fun3", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0, p1, p2, p3, p4)
	}

	// invoke fun4
	{
		bigInt1 := big.NewInt(-100001)
		a16int := int16(-10001)
		bigInt3 := big.NewInt(10001)
		bigInt4 := big.NewInt(1111111)
		a16uint := uint16(10001)
		bigInt6 := big.NewInt(111111)
		packed, _ := ABI.Pack("fun4", bigInt1, a16int, bigInt3, bigInt4, a16uint, bigInt6)
		invokeTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
		invokeTx.Sign(guomiKey)
		invokeRe, _ := rpc.InvokeContract(invokeTx)

		var p0 interface{}
		var p1 int16
		var p2 *big.Int
		var p3 interface{}
		var p4 uint16
		var p5 *big.Int
		testV := []interface{}{&p0, &p1, &p2, &p3, &p4, &p5}
		if err := ABI.UnpackResult(&testV, "fun4", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0, p1, p2, p3, p4, p5)
	}

	// invoke fun5
	{
		address := common.Address{}
		address.SetString("2312321312")
		packed, _ := ABI.Pack("fun5", "data1", address)
		invokeTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
		invokeTx.Sign(guomiKey)
		invokeRe, _ := rpc.InvokeContract(invokeTx)

		var p0 string
		var p1 common.Address
		testV := []interface{}{&p0, &p1}
		if err := ABI.UnpackResult(&testV, "fun5", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0, p1)
	}
}

func TestRPC_ListenContract(t *testing.T) {
	t.Skip("hyperchain needs to start the radar service first")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0])
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)
	srcCode, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")
	result, err := rpc.ListenContract(string(srcCode), receipt.ContractAddress)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
}

func TestRPC_Setter(t *testing.T) {
	rpc := NewRPC()
	rpc.Namespace("global").ResendTimes(int64(1)).FirstPollTime(int64(1))
	rpc.FirstPollInterval(int64(1)).SecondPollInterval(int64(1)).SecondPollTime(int64(1))
	rpc.ReConnTime(int64(1))
	rpc.AddNode("127.0.0.1", "8081", "10001")
	//nolint
	rpc.BindNodes(1)
	fmt.Println(rpc)
}

func TestRPC_Simulate(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).Simulate(true)
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)
}
