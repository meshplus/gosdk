package rpc

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/bvm"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 新增合约地址和合约命名的映射关系
func TestRPC_CNSSmoke_Mapping(t *testing.T) {
	t.Skip()
	addr1 := "0x0000000000000000000000000000000000ffff01"
	addr2 := "0x0000000000000000000000000000000000ffff02"
	name1 := "HashContract"
	name2 := "HashContract1"
	t.Run("init", func(t *testing.T) {
		initHashContract(t)
	})
	// 校验cns_contract维护合约地址和合约地址别名的映射关系
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5ebe0b3f785e150012748a72/testcase/5ebe0e93785e15001274ab0a
	t.Run("SetCName_normal", func(t *testing.T) {
		// 合约地址addr1无别名
		_, err := rpc.GetNameByAddress(addr1)
		assert.NotNil(t, err)
		// 新增映射关系成功
		result := setCName(addr1, name1, t)
		assert.True(t, result.Success)
		na, err := rpc.GetNameByAddress(addr1)
		assert.Nil(t, err)
		assert.Equal(t, name1, na)
	})

	// 校验cns_contract维护合约地址和合约地址别名的映射关系
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5ebe0b3f785e150012748a72/testcase/5ebe0e93785e15001274ab0c
	t.Run("SetCName_one address_two_name", func(t *testing.T) {
		// 发起新增合约地址addr1和合约地址别名name2的映射关系的提案
		result := setCName(addr1, name2, t)
		assert.True(t, result.Success)
		var res []*bvm.OpResult
		_ = json.Unmarshal([]byte(result.Ret), &res)
		// 新增失败，交易回执会返回error
		assert.Equal(t, bvm.CallErrorCode, res[0].Code)
	})

	// 映射关系的异常场景-重复新增同一映射关系
	t.Run("repeat_SetCName", func(t *testing.T) {
		// 发起新增合约地址地址addr1 -> name1,name1->addr1，整个提案投票通过，执行的时候预期的执行结果
		result := setCName(addr1, name1, t)
		assert.True(t, result.Success)
		var res []*bvm.OpResult
		_ = json.Unmarshal([]byte(result.Ret), &res)
		// 新增失败，交易回执会返回error
		assert.Equal(t, bvm.CallErrorCode, res[0].Code)
	})

	// 映射关系的异常场景-不同合约地址映射同一合约命名
	t.Run("SetCName_used_cName", func(t *testing.T) {
		// 发起新增合约地址：addr2->name1,name1->addr2，整个提案投票通过，执行的时候预期的执行结果
		result := setCName(addr2, name1, t)
		assert.True(t, result.Success)
		var res []*bvm.OpResult
		_ = json.Unmarshal([]byte(result.Ret), &res)
		// 新增失败，交易回执会返回error
		assert.Equal(t, bvm.CallErrorCode, res[0].Code)
	})

	// 校验cns_contract维护合约地址和合约地址别名的映射关系
	t.Run("SetCName_Check_Address_length", func(t *testing.T) {
		re, err := createProposalForSetCName("0x0000fff1", "Check_Address_length", t)
		assert.Nil(t, err)
		result := bvm.Decode(re.Ret)
		t.Log(result)
		assert.False(t, result.Success)
	})

	// 根据cns执行交易
	t.Run("invoke_name_len_big_than_40", func(t *testing.T) {
		// 传入name参数，参数大于40个字符
		name := "HashContract_HashContract_HashContract_HashContract"
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(name, payload).VMType(BVM)
		tx.Sign(key)
		_, err := rpc.InvokeContract(tx)
		assert.NotNil(t, err)
		t.Log(err)
	})

	// 根据cns执行交易
	t.Run("invoke_name_and_to_is_nil", func(t *testing.T) {
		// api层接口，例如tx_sendTransaction的rpc方法，新增合约命名请求参数校验：cnsName和to合约地址参数都为“”或者null（可选参数）
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName("", payload).VMType(BVM)
		tx.Sign(key)
		_, err := rpc.InvokeContract(tx)
		assert.NotNil(t, err)
		t.Log(err)

		tx = NewTransaction(key.GetAddress().Hex()).Invoke(common.Address{}.Hex(), payload).VMType(BVM)
		tx.Sign(key)
		_, err = rpc.InvokeContract(tx)
		assert.NotNil(t, err)
		t.Log(err)
	})

	// 根据cns执行交易
	t.Run("invoke_name_lenth_less_than_2", func(t *testing.T) {
		// 传入name参数，参数少于2个字符
		name := "a"
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(name, payload).VMType(BVM)
		tx.Sign(key)
		_, err := rpc.InvokeContract(tx)
		assert.NotNil(t, err)
		t.Log(err)
	})

	// 根据cns执行交易
	t.Run("invoke_success", func(t *testing.T) {
		// 传入name参数，参数为正确的合约命名
		name := name1
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(name, payload).VMType(BVM)
		tx.Sign(key)
		re, err := rpc.InvokeContract(tx)
		assert.Nil(t, err)
		t.Log(re.TxHash)
		result := bvm.Decode(re.Ret)
		assert.True(t, result.Success)
		assert.Equal(t, "0x456", string(result.Ret))
	})

	// 根据cns执行交易
	t.Run("invoke_name_contain_invalid_char", func(t *testing.T) {
		// 传入name参数，数据格式为“？# abc”或者中文字符
		name := "? # abc"
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(name, payload).VMType(BVM)
		tx.Sign(key)
		_, err := rpc.InvokeContract(tx)
		assert.NotNil(t, err)
		t.Log(err)
	})
}

// 根据合约命名cnsName查询与之对应的合约地址address
func TestRPC_CNSSmoke_QueryTx(t *testing.T) {
	t.Skip()
	addr1 := "0x0000000000000000000000000000000000ffff01"
	name1 := "HashContract"
	nameNotExist := "contract"

	// cnsmanager的查询转换-tx_getTransactionReceipt查询指定交易的回执
	t.Run("get_tx_receipt", func(t *testing.T) {
		// 当用户发交易的时候通过To调用合约，查询回执
		// invoke by to
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(addr1, payload).VMType(BVM)
		tx.Sign(key)
		re, err := rpc.InvokeContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re.TxHash)
		// get receipt
		t.Log(re.TxHash)
		receipt, err := rpc.GetTxReceipt(re.TxHash, false)
		assert.Nil(t, err)
		// contractAddress：addr1, contractName:name1
		assert.Equal(t, addr1, receipt.ContractAddress)
		assert.Equal(t, name1, receipt.ContractName)

		// 当用户发交易的时候通过cName调用合约，查询回执
		// invoke by name
		tx = NewTransaction(key.GetAddress().Hex()).InvokeByName(name1, payload).VMType(BVM)
		tx.Sign(key)
		re, err = rpc.InvokeContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re.TxHash)
		t.Log(re.TxHash)
		// get receipt
		receipt, err = rpc.GetTxReceipt(re.TxHash, false)
		assert.Nil(t, err)
		// contractAddress：addr1, contractName:name1
		assert.Equal(t, addr1, receipt.ContractAddress)
		assert.Equal(t, name1, receipt.ContractName)
	})

	// cnsmanager的查询转换-当用户发交易的时候通过cName调用合约，通过rpc方法：tx_getTransactionsByTime查询指定时间区间内的交易
	t.Run("tx_getTxsByTime", func(t *testing.T) {
		// 当用户发交易的时候通过cName调用合约，通过rpc方法：tx_getTransactionsByTime查询指定时间区间内的交易
		// invoke by name
		start := time.Now().UnixNano()
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashSetOperation("0x111", "0x222")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(name1, payload).VMType(BVM)
		tx.Sign(key)
		re, err := rpc.InvokeContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re.TxHash)
		t.Log(re.TxHash)
		end := time.Now().UnixNano()
		// get txs by time
		pageResult, err := rpc.GetTxByTimeWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		})
		assert.Nil(t, err)
		assert.False(t, pageResult.HasMore)
		assert.Len(t, pageResult.Data, 1)
		assert.Equal(t, re.TxHash, pageResult.Data[0].Hash)
		assert.Equal(t, name1, pageResult.Data[0].CName)
		assert.Equal(t, common.Address{}.Hex(), pageResult.Data[0].To)
		// get txs by time and contract address
		pageResult, err = rpc.GetTxByTimeAndContractAddrWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		}, addr1)
		assert.Nil(t, err)
		assert.False(t, pageResult.HasMore)
		assert.Len(t, pageResult.Data, 1)
		assert.Equal(t, re.TxHash, pageResult.Data[0].Hash)
		assert.Equal(t, name1, pageResult.Data[0].CName)
		assert.Equal(t, common.Address{}.Hex(), pageResult.Data[0].To)
		// get txs by time and contract name
		pageResult, err = rpc.GetTxByTimeAndContractAddrWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		}, addr1)
		assert.Nil(t, err)
		assert.False(t, pageResult.HasMore)
		assert.Len(t, pageResult.Data, 1)
		assert.Equal(t, re.TxHash, pageResult.Data[0].Hash)
		assert.Equal(t, name1, pageResult.Data[0].CName)
		assert.Equal(t, common.Address{}.Hex(), pageResult.Data[0].To)
	})

	// cnsmanager的查询转换-tx_getTransactionReceipt查询指定交易的回执
	t.Run("get_tx_receipt_not_exist", func(t *testing.T) {
		// invoke by name
		start := time.Now().UnixNano()
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(nameNotExist, []byte("transaction payload")).VMType(EVM)
		tx.Sign(key)
		// 执行时有查回执的操作，非法交易，回执返回nil，err为错误信息
		re, err := rpc.InvokeContract(tx)
		assert.Error(t, err)
		assert.Nil(t, re)
		t.Log(err)
		end := time.Now().UnixNano()

		// get txs by time
		pageResult, err := rpc.GetTxByTimeWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		})
		assert.Nil(t, err)
		assert.False(t, pageResult.HasMore)
		assert.Len(t, pageResult.Data, 0)
	})

	// cnsmanager的查询转换-当用户发交易的时候通过to调用合约，通过rpc方法：tx_getTransactionsByTime查询指定时间区间内的交易
	t.Run("tx_getTxsByTime_invokeByTo", func(t *testing.T) {
		// 当用户发交易的时候通过to调用合约，通过rpc方法：tx_getTransactionsByTime查询指定时间区间内的交易
		// invoke by to
		start := time.Now().UnixNano()
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashSetOperation("0x111", "0x222")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(addr1, payload).VMType(BVM)
		tx.Sign(key)
		re, err := rpc.InvokeContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re.TxHash)
		t.Log(re.TxHash)
		end := time.Now().UnixNano()
		// get txs by time
		pageResult, err := rpc.GetTxByTimeWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		})
		assert.Nil(t, err)
		assert.False(t, pageResult.HasMore)
		assert.Len(t, pageResult.Data, 1)
		assert.Equal(t, re.TxHash, pageResult.Data[0].Hash)
		assert.Equal(t, "", pageResult.Data[0].CName)
		assert.Equal(t, addr1, pageResult.Data[0].To)

		// get txs by time and contract address
		pageResult, err = rpc.GetTxByTimeAndContractAddrWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		}, addr1)
		assert.Nil(t, err)
		assert.False(t, pageResult.HasMore)
		assert.Equal(t, 1, len(pageResult.Data))
		assert.Equal(t, re.TxHash, pageResult.Data[0].Hash)
		assert.Equal(t, "", pageResult.Data[0].CName)
		assert.Equal(t, addr1, pageResult.Data[0].To)
		// get txs by time and contract name
		pageResult, err = rpc.GetTxByTimeAndContractNameWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		}, name1)
		assert.Nil(t, err)
		assert.False(t, pageResult.HasMore)
		assert.Equal(t, 1, len(pageResult.Data))
		assert.Equal(t, re.TxHash, pageResult.Data[0].Hash)
		assert.Equal(t, "", pageResult.Data[0].CName)
		assert.Equal(t, addr1, pageResult.Data[0].To)
	})

	// cnsmanager的查询转换-tx_getTransactionsCountByContractAddr查询区块区间交易数量
	t.Run("tx_getTxsCount", func(t *testing.T) {
		block, _ := rpc.GetLatestBlock()
		// 查询类接口使用合约地址address，通过rpc方法tx_getTransactionsCountByContractAddr查询区块区间交易数量
		countByAddr, err := rpc.GetTxCountByContractAddr(1, block.Number, addr1, false)
		assert.Nil(t, err)

		// 查询类接口使用合约地址命名cName，通过rpc方法tx_getTransactionsCountByContractName查询区块区间交易数量
		countByName, err := rpc.GetTxCountByContractName(1, block.Number, name1, false)
		assert.Nil(t, err)

		assert.Equal(t, countByAddr.Count, countByName.Count)
		t.Log(countByName.Count)
	})

}

// 查询合约地址与合约命名的映射关系列表
func TestRPC_CNSSmoke_QueryCNS(t *testing.T) {
	t.Skip()
	addr1 := "0x0000000000000000000000000000000000ffff01"
	name1 := "HashContract"

	// 根据cnsManage做查询转换-查询由合约命名发起的交易
	t.Run("get_tx_by_name", func(t *testing.T) {
		// GetContractAddressByName接口传入参数校验，传入正确的name参数
		addr, err := rpc.GetAddressByName(name1)
		// 查询成功
		assert.Nil(t, err)
		assert.Equal(t, addr1, addr)

		block, _ := rpc.GetLatestBlock()
		min := block.Number
		// invoke by name
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(name1, payload).VMType(BVM)
		tx.Sign(key)
		re, err := rpc.InvokeContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re.TxHash)
		block, _ = rpc.GetLatestBlock()
		max := block.Number
		// 通过rpc方法查询下一页接口tx_getNextPageTransactions使用合约地址addr下发查询
		txs, err := rpc.GetNextPageTxs(min, 0, min, max, 0, 10, false, addr1)
		assert.Nil(t, err)
		marshal, _ := json.Marshal(txs)

		// 通过rpc方法查询下一页接口tx_getNextPageTransactions使用合约命名name下发查询
		txs1, err := rpc.GetNextPageTxsByName(min, 0, min, max, 0, 10, false, name1)
		assert.Nil(t, err)
		marshal1, _ := json.Marshal(txs1)

		assert.Equal(t, 1, len(txs))
		assert.Equal(t, txs, txs1)
		assert.Equal(t, string(marshal), string(marshal1))
	})

	// 根据cnsManage做查询转换-查询由合约地址发起的交易
	t.Run("get_tx_by_to", func(t *testing.T) {
		block, _ := rpc.GetLatestBlock()
		min := block.Number
		// invoke by to
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(addr1, payload).VMType(BVM)
		tx.Sign(key)
		re, err := rpc.InvokeContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re.TxHash)
		block, _ = rpc.GetLatestBlock()
		max := block.Number
		// 通过rpc方法查询下一页接口tx_getNextPageTransactions使用合约地址addr下发查询
		txs, err := rpc.GetNextPageTxs(min, 0, min, max, 0, 10, false, addr1)
		assert.Nil(t, err)
		marshal, _ := json.Marshal(txs)

		// 通过rpc方法查询下一页接口tx_getNextPageTransactions使用合约命名name下发查询
		txs1, err := rpc.GetNextPageTxsByName(min, 0, min, max, 0, 10, false, name1)
		assert.Nil(t, err)
		marshal1, _ := json.Marshal(txs1)

		assert.Equal(t, 1, len(txs))
		assert.Equal(t, txs, txs1)
		assert.Equal(t, string(marshal), string(marshal1))
	})

	// 根据cnsManage做查询转换-查询所有的合约地址address与合约命名cnsName的映射关系
	t.Run("GetAll", func(t *testing.T) {
		allCNS, err := rpc.GetAllCNS()
		assert.Nil(t, err)
		assert.Equal(t, 1, len(allCNS))
		name, ok := allCNS[addr1]
		assert.True(t, ok)
		assert.Equal(t, name1, name)
		t.Log(allCNS)
	})

}

// 遍历查询交易的rpc接口
func TestRPC_CNSSmoke_RPC(t *testing.T) {
	t.Skip()
	addr1 := "0x0000000000000000000000000000000000ffff01"
	name1 := "HashContract"

	// invoke by to
	key, _ := account.NewAccountFromAccountJSON("", pwd)
	opt := bvm.NewHashGetOperation("0x123")
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(addr1, payload).VMType(BVM)
	tx.Sign(key)
	re, err := rpc.InvokeContract(tx)
	assert.Nil(t, err)
	assert.NotNil(t, re.TxHash)
	toHash := re.TxHash

	// invoke by name
	tx = NewTransaction(key.GetAddress().Hex()).InvokeByName(name1, payload).VMType(BVM)
	tx.Sign(key)
	re, err = rpc.InvokeContract(tx)
	assert.Nil(t, err)
	assert.NotNil(t, re.TxHash)
	nameHash := re.TxHash

	// 根据区块号查询交易接口
	t.Run("get_tx_by_blk_number", func(t *testing.T) {
		// get block number
		tx, err := rpc.GetTransactionByHash(toHash)
		assert.Nil(t, err)

		// get tx by block number
		txInfo, err := rpc.GetTxByBlkNumAndIdx(tx.BlockNumber, tx.TxIndex)
		assert.Nil(t, err)
		assert.Equal(t, addr1, txInfo.To)
		assert.Equal(t, "", txInfo.CName)

		// get block number
		tx, err = rpc.GetTransactionByHash(nameHash)
		assert.Nil(t, err)

		// get tx by block number
		txInfo, err = rpc.GetTxByBlkNumAndIdx(tx.BlockNumber, tx.TxIndex)
		assert.Nil(t, err)
		assert.Equal(t, common.Address{}.Hex(), txInfo.To)
		assert.Equal(t, name1, txInfo.CName)

	})

	// 根据区块哈希查询交易
	t.Run("get_tx_by_blk_hash", func(t *testing.T) {
		// get block number
		tx, err := rpc.GetTransactionByHash(toHash)
		assert.Nil(t, err)

		// get tx by block hash
		txInfo, err := rpc.GetTxByBlkHashAndIdx(tx.BlockHash, tx.TxIndex)
		assert.Nil(t, err)
		assert.Equal(t, addr1, txInfo.To)
		assert.Equal(t, "", txInfo.CName)

		// get block number
		tx, err = rpc.GetTransactionByHash(nameHash)
		assert.Nil(t, err)

		// get tx by block hash
		txInfo, err = rpc.GetTxByBlkHashAndIdx(tx.BlockHash, tx.TxIndex)
		assert.Nil(t, err)
		assert.Equal(t, common.Address{}.Hex(), txInfo.To)
		assert.Equal(t, name1, txInfo.CName)
	})

	// 根据交易Hash查询交易
	t.Run("GetTxByHash", func(t *testing.T) {
		// get tx by to hash
		txInfo, stdErr := rpc.GetTransactionByHash(toHash)
		assert.Nil(t, stdErr)
		assert.Equal(t, addr1, txInfo.To)
		assert.Equal(t, "", txInfo.CName)

		// get tx by name hash
		txInfo, stdErr = rpc.GetTransactionByHash(nameHash)
		assert.Nil(t, stdErr)
		assert.Equal(t, common.Address{}.Hex(), txInfo.To)
		assert.Equal(t, name1, txInfo.CName)
	})

	// 根据交易哈希批量查询交易
	t.Run("tx_getBatchTransactions", func(t *testing.T) {
		infos, err := rpc.GetBatchTxByHash([]string{toHash})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(infos))
		assert.Equal(t, addr1, infos[0].To)
		assert.Equal(t, "", infos[0].CName)

		infos, err = rpc.GetBatchTxByHash([]string{nameHash})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(infos))
		assert.Equal(t, common.Address{}.Hex(), infos[0].To)
		assert.Equal(t, name1, infos[0].CName)
	})

	// 查询指定时间区间内的所有交易接口
	//t.Run("")
}

func TestRPC_CNSSmoke_Sync(t *testing.T) {
	t.Skip()
	addr2 := "0x0000000000000000000000000000000000ffff02"
	name2 := "ProposalContract"
	t.Run("setCName", func(t *testing.T) {
		result := setCName(addr2, name2, t)
		assert.True(t, result.Success)
	})

	t.Run("checkSync", func(t *testing.T) {
		// check mapping
		name, err := rpc.GetNameByAddress(addr2)
		assert.Nil(t, err)
		assert.Equal(t, name, name2)
		// invoke by name
		key, _ := account.NewAccountFromAccountJSON("", pwd)
		opt := bvm.NewPermissionCreateRoleOperation("manager")
		permissionOpt := bvm.NewProposalCreateOperationForPermission(opt)
		payload := bvm.EncodeOperation(permissionOpt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(name2, payload).VMType(BVM)
		tx.Sign(key)
		receipt, err := rpc.InvokeContract(tx)
		assert.Nil(t, err)
		result := bvm.Decode(receipt.Ret)
		assert.True(t, result.Success)
	})

	t.Run("getProposal", func(t *testing.T) {
		proposal, _ := rpc.GetProposal()
		marshal, _ := json.Marshal(proposal)
		t.Log(string(marshal))
	})
}

func setCName(address, cName string, t *testing.T) *bvm.Result {
	re, err := createProposalForSetCName(address, cName, t)
	var proposal bvm.ProposalData
	assert.NoError(t, err)
	result := bvm.Decode(re.Ret)
	assert.True(t, result.Success)
	_ = proto.Unmarshal([]byte(result.Ret), &proposal)

	voteProposal(int(proposal.Id), t)
	return executeProposal(int(proposal.Id), t)
}

func createProposalForSetCName(address, cName string, t *testing.T) (*TxReceipt, StdError) {
	key, _ := account.NewAccountFromAccountJSON("", pwd)
	setCNameOpt := bvm.NewCNSSetCNameOperation(address, cName)
	cnsOpt := bvm.NewProposalCreateOperationForCNS(setCNameOpt)
	payload := bvm.EncodeOperation(cnsOpt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(cnsOpt.Address(), payload).VMType(BVM)
	tx.Sign(key)
	return rpc.InvokeContract(tx)
}

func voteProposal(id int, t *testing.T) {
	for i := 1; i < 6; i++ {
		var key account.Key
		if i < 4 {
			key, _ = account.NewAccountFromAccountJSON("", pwd)
		} else {
			key, _ = account.NewAccountSm2FromAccountJSON("", pwd)
		}
		opt := bvm.NewProposalVoteOperation(id, true)
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rpc.InvokeContract(tx)
		assert.NoError(t, err)
		result := bvm.Decode(re.Ret)
		assert.True(t, result.Success)
	}
}

func executeProposal(id int, t *testing.T) *bvm.Result {
	key, _ := account.NewAccountFromAccountJSON("", pwd)
	opt := bvm.NewProposalExecuteOperation(id)
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rpc.InvokeContract(tx)
	assert.NoError(t, err)
	result := bvm.Decode(re.Ret)
	t.Log(result)
	return result
}

func initHashContract(t *testing.T) {
	key, _ := account.NewAccountFromAccountJSON("", pwd)
	opt := bvm.NewHashSetOperation("0x123", "0x456")
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rpc.InvokeContract(tx)
	assert.NoError(t, err)
	result := bvm.Decode(re.Ret)
	assert.True(t, result.Success)
}

func TestRPC_GetNextPageTxs2(t *testing.T) {
	t.Skip()
	txs, stdError := rpc.GetNextPageTxs(24, 0, 24, 26, 0, 10, true, "0x0000000000000000000000000000000000ffff01")
	assert.Nil(t, stdError)
	marshal, _ := json.Marshal(txs)
	t.Log(string(marshal))
	txs, stdError = rpc.GetNextPageTxsByName(24, 0, 24, 26, 0, 10, true, "HashContract")
	assert.Nil(t, stdError)
	marshal1, _ := json.Marshal(txs)
	t.Log(string(marshal1))
	assert.Equal(t, string(marshal), string(marshal1))

	txs, stdError = rpc.GetPrevPageTxs(26, 0, 24, 26, 0, 10, true, "0x0000000000000000000000000000000000ffff01")
	assert.Nil(t, stdError)
	marshal, _ = json.Marshal(txs)
	t.Log(string(marshal))
	txs, stdError = rpc.GetPrevPageTxsByName(26, 0, 24, 26, 0, 10, true, "HashContract")
	assert.Nil(t, stdError)
	marshal1, _ = json.Marshal(txs)
	t.Log(string(marshal1))
	assert.Equal(t, string(marshal), string(marshal1))

}
