package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/hash"
	"io/ioutil"
	"math/big"
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
	pri      = new(gm.SM2PrivateKey)
	_        = pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey = &account.SM2Key{
		SM2PrivateKey: &gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
)

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

func TestRPC_GetCode(t *testing.T) {
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	receipt, err := rpc.SignAndDeployContract(transaction, guomiKey)
	assert.Nil(t, err)
	_, err = rpc.GetCode(receipt.ContractAddress)
	assert.Nil(t, err)
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
	t.Skip("the node can get the account")
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

func TestRPC_QueryLatestArchive(t *testing.T) {
	t.Skip()
	res, err := rpc.QueryLatestArchive()
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
	hash1, _ := rpc.GetNodeHashByID(6)
	rpc, _ = rpc.BindNodes(1)
	success, _ := rpc.DeleteNodeNVP(hash1)
	assert.Equal(t, true, success)
}

func TestRPC_DisconnectNodeVP(t *testing.T) {
	t.Skip("do not delete NVP in CI")
	hash1, _ := rpc.GetNodeHashByID(1)
	rpc, _ = rpc.BindNodes(6)
	success, _ := rpc.DisconnectNodeVP(hash1)
	assert.Equal(t, true, success)
}

func TestRPC_ReplaceNodeCerts(t *testing.T) {
	t.Skip("do not delete NVP in CI")
	rpc, _ = rpc.BindNodes(1)
	hash1, err := rpc.ReplaceNodeCerts("node1")
	assert.Nil(t, err)
	fmt.Println(hash1)
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

func TestRPC_GetInvalidTransactionsWithLimit(t *testing.T) {
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

	pageResult, err := rpc.GetInvalidTransactionsByBlkNumWithLimit(block.Number-1, block.Number, metadata)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pageResult)
}

func TestRPC_GetInvalidTxByBlockNumber(t *testing.T) {
	t.Skip()
	block, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	txInfos, err := rpc.GetInvalidTransactionsByBlkNum(block.Number)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txInfos)
}

func TestRPC_GetInvalidTxByBlockHash(t *testing.T) {
	t.Skip()
	block, err := rpc.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	txInfos, err := rpc.GetInvalidTransactionsByBlkHash(block.Hash)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txInfos)
}

func TestRPC_GetInvalidTxsCount(t *testing.T) {
	t.Skip()
	count, err := rpc.GetInvalidTxCount()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
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
	t.Skip()
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

/*---------------------------------- proof ----------------------------------*/

func TestRPC_GetAccountProof(t *testing.T) {
	t.Skip()
	act := "e93b92f1da08f925bdee44e91e7768380ae83307"
	res, err := rpc.GetAccountProof(act)
	if err != nil {
		t.Error(err)
		return
	}

	assert.True(t, ValidateAccountProof(act, res))
}

func TestRPC_GetTxProof(t *testing.T) {
	t.Skip()
	block, _ := rpc.GetLatestBlock()
	info, err := rpc.GetTxByBlkHashAndIdx(block.Hash, 0)
	if err != nil {
		t.Error(err)
		return
	}
	res, err2 := rpc.GetTxProof(info.Hash)
	if err2 != nil {
		t.Error(err2)
		return
	}
	ast := assert.New(t)
	ast.True(ValidateTxProof(info.Hash, block.TxRoot, res))
}

func TestRPC_ArchiveSnapshot(t *testing.T) {
	t.Skip("need flato 1.5.0")
	deployJar, err := DecompressFromJar("../hvmtestfile/fibonacci/fibonacci-1.0-fibonacci.jar")
	if err != nil {
		t.Error(err)
	}
	accountJson, sysErr := account.NewAccountJson(account.SMRAW, "")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	key, sysErr := account.GenKeyFromAccountJson(accountJson, "")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}

	newAddress := key.(*account.SM2Key).GetAddress()
	transaction := NewTransaction(newAddress.Hex()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
	transaction.Sign(key)
	receipt, err := rpc.DeployContract(transaction)
	assert.Nil(t, err)
	t.Log("contract address:", receipt.ContractAddress)

	b, err := rpc.GetLatestBlock()
	assert.Nil(t, err)

	time.Sleep(20 * time.Second)
	id, err := rpc.Snapshot(b.Number)
	assert.Nil(t, err)
	t.Log(b.Number)
	t.Log(b.MerkleRoot)
	t.Log(id)
}

func TestRPC_GetStateProof(t *testing.T) {
	// should send to archiveReader!
	t.Skip("need archiveReader")
	id := "0x5b1a5bb7b10d15bc9d47701eed9c9349"
	seq := 2
	contractAddr := "0x6de31be7a30204189d70bd202340c6d9b395523e"
	merkleRoot := "0xaa2fd673656f4bada6ff6d8588498239eeb3202214a24005d6cf0138a9f30a79"
	proofParam := &ProofParam{
		Meta: &LedgerMetaParam{
			SnapshotID: id,
			SeqNo:      uint64(seq),
		},
		Key: &KeyParam{
			Address:   common.HexToAddress(contractAddr),
			FieldName: "hyperMap1",
			Params:    []string{"key1"},
			VMType:    "HVM",
		},
	}
	proof, err := rpc.GetStateProof(proofParam)
	assert.Nil(t, err)
	t.Log(proof)
	b, _ := json.Marshal(proof)
	t.Log(string(b))

	ok, err := rpc.ValidateStateProof(proofParam, proof, merkleRoot)
	assert.Nil(t, err)
	assert.True(t, ok)
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
	rpc.SetAccount(privateKey)
	err := rpc.AddRoleForNode(privateKey.GetAddress().Hex(), "accountManager")
	if err != nil {
		t.Error(err)
		return
	}
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
	txReceipt, err := rpc.SignAndDeployContract(transaction, guomiKey)
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

func TestRPC_SignAndInvokeContractCombineReturns(t *testing.T) {
	accountJSON, _ := account.NewAccountED25519("12345678")
	ekey, err := account.GenKeyFromAccountJson(accountJSON, "12345678")
	assert.Nil(t, err)
	newAddress := ekey.(*account.ED25519Key).GetAddress()

	transaction := NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
	transaction.Sign(ekey)
	txreceipt, info, err := rpc.SignAndInvokeContractCombineReturns(transaction, ekey)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, info)
	assert.NotNil(t, txreceipt)
	assert.Equal(t, txreceipt.TxHash, info.Hash)
	assert.Nil(t, err)
}
