package rpc

import (
	"fmt"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/bvm"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/hvm"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRPC_InspectorRole(t *testing.T) {
	t.Skip()
	//newAccount, err := account.NewAccount("")
	//assert.Nil(t, err)
	//t.Log(newAccount)
	defaultPwd := ""
	newAccount := `{"address":"0x430d637bb98033e1a2406ab8c24949e142fa4e75","algo":"0x03","version":"4.0","publicKey":"0x0454bbfe47a39685ac16d8263b0f7e3e0406996dc218543399411dfff12e304768945d0afe07c6b932fb203d1fb670fc46839adce54efed8e153fd21779c29d02b","privateKey":"e549225831723fdf9a7d4530b17b746e2c9da37fc1a982f4927f94d2b2972ac2"}`
	ac, err := account.NewAccountFromAccountJSON(newAccount, defaultPwd)
	assert.Nil(t, err)

	role := "user"
	role1 := "account"
	t.Run("add_role_success", func(t *testing.T) {
		stdError := rpc.AddRoleForNode(ac.GetAddress().Hex(), role, role1)
		assert.Nil(t, stdError)
	})

	t.Run("add_role_already_exist", func(t *testing.T) {
		stdError := rpc.AddRoleForNode(ac.GetAddress().Hex(), role)
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address:%s already has roles:[%s]", ac.GetAddress().Hex(), role)))
	})

	t.Run("delete_role_success", func(t *testing.T) {
		stdError := rpc.DeleteRoleFromNode(ac.GetAddress().Hex(), role1)
		assert.Nil(t, stdError)
	})

	t.Run("delete_role_success", func(t *testing.T) {
		stdError := rpc.DeleteRoleFromNode(ac.GetAddress().Hex(), role1)
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address:%s has not roles:[%s]", ac.GetAddress().Hex(), role1)))
	})

	t.Run("get_role", func(t *testing.T) {
		roles, err := rpc.GetRoleFromNode(ac.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Equal(t, roles, []string{role})
	})

	t.Run("get_all_role", func(t *testing.T) {
		roles, err := rpc.GetAllRolesFromNode()
		assert.Nil(t, err)
		assert.NotNil(t, roles)
	})

	// inspector enable
	t.Run("auth_forbidden_and_authorized_and_methods", func(t *testing.T) {
		// set rules
		user := "user"
		admin := "admin"
		rule := &InspectorRule{
			AllowAnyone:     false,
			AuthorizedRoles: []string{admin},
			ForbiddenRoles:  []string{user},
			ID:              0,
			Method:          []string{"node_*"},
			Name:            "node",
		}
		stdError := rpc.SetRulesInNode([]*InspectorRule{rule})
		assert.Nil(t, stdError)

		// add user to account A
		accountA, _ := account.NewAccount(defaultPwd)
		acA, _ := account.NewAccountFromAccountJSON(accountA, defaultPwd)
		stdError = rpc.AddRoleForNode(acA.GetAddress().Hex(), user)
		assert.Nil(t, stdError)
		// use account A call node_getNodeStates has not permission
		rpc.SetAccount(acA)
		_, stdError = rpc.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", acA.GetAddress().Hex())))
		// use account A call node_getNodes has not permission
		rpc.SetAccount(acA)
		_, stdError = rpc.GetNodes()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodes", acA.GetAddress().Hex())))

		// add admin to account B
		accountB, _ := account.NewAccount(defaultPwd)
		acB, _ := account.NewAccountFromAccountJSON(accountB, defaultPwd)
		stdError = rpc.AddRoleForNode(acB.GetAddress().Hex(), admin)
		assert.Nil(t, stdError)
		// use account B call node_getNodeStates has permission
		rpc.SetAccount(acB)
		nodeState, stdError := rpc.GetNodeStates()
		assert.Nil(t, stdError)
		assert.NotNil(t, nodeState)
		t.Log(nodeState)
	})

	t.Run("auth_allow_anyone_false", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = [""]
		//id = 0
		//methods = ["node_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rule := &InspectorRule{
			AllowAnyone:     false,
			AuthorizedRoles: []string{admin},
			ForbiddenRoles:  []string{},
			ID:              0,
			Method:          []string{"node_*"},
			Name:            "node",
		}
		stdError := rpc.SetRulesInNode([]*InspectorRule{rule})
		assert.Nil(t, stdError)

		// add user to account B
		accountB, _ := account.NewAccount(defaultPwd)
		acB, _ := account.NewAccountFromAccountJSON(accountB, defaultPwd)
		stdError = rpc.AddRoleForNode(acB.GetAddress().Hex(), user)
		assert.Nil(t, stdError)
		// use account B call node_getNodeStates has not permission
		rpc.SetAccount(acB)
		_, stdError = rpc.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", acB.GetAddress().Hex())))
	})

	t.Run("auth_rules_id", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = [""]
		//id = 0
		//methods = ["node_*"]
		//name = "node"
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = ["user"]
		//id = 1
		//methods = ["node_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ID:              0,
				Method:          []string{"node_*"},
				Name:            "node",
			},
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"node_*"},
				Name:            "node",
			},
		}
		stdError := rpc.SetRulesInNode(rules)
		assert.Nil(t, stdError)
		inspectorRules, stdError := rpc.GetRulesFromNode()
		assert.Nil(t, stdError)
		assert.Len(t, inspectorRules, len(rules))

		// add user to account C
		accountC, _ := account.NewAccount(defaultPwd)
		acC, _ := account.NewAccountFromAccountJSON(accountC, defaultPwd)
		stdError = rpc.AddRoleForNode(acC.GetAddress().Hex(), user)
		assert.Nil(t, stdError)
		// use account C call node_getNodeStates has not permission
		rpc.SetAccount(acC)
		_, stdError = rpc.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", acC.GetAddress().Hex())))

	})

	t.Run("auth_priority", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = true
		//authorized_roles = ["admin"]
		//forbidden_roles = ["user"]
		//id = 1
		//methods = ["node_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     true,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"node_*"},
				Name:            "node",
			},
		}
		stdError := rpc.SetRulesInNode(rules)
		assert.Nil(t, stdError)

		// add user to account D
		accountD, _ := account.NewAccount(defaultPwd)
		acD, _ := account.NewAccountFromAccountJSON(accountD, defaultPwd)
		stdError = rpc.AddRoleForNode(acD.GetAddress().Hex(), user, admin)
		assert.Nil(t, stdError)
		// use account D call node_getNodeStates has not permission
		rpc.SetAccount(acD)
		_, stdError = rpc.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", acD.GetAddress().Hex())))
	})

	t.Run("auth_transaction", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = ["","user"]
		//id = 1
		//methods = ["tx_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"contract_*"},
				Name:            "node",
			},
		}
		stdError := rpc.SetRulesInNode(rules)
		assert.Nil(t, stdError)

		// add user,admin to account H
		accountH, _ := account.NewAccount(defaultPwd)
		acH, _ := account.NewAccountFromAccountJSON(accountH, defaultPwd)
		stdError = rpc.AddRoleForNode(acH.GetAddress().Hex(), user, admin)
		assert.Nil(t, stdError)

		// send tx
		transaction := NewTransaction(acH.GetAddress().Hex()).Transfer(address, int64(0))
		transaction.Sign(acH)
		receipt, stdError := rpc.SendTx(transaction)
		assert.Nil(t, stdError)
		assert.NotNil(t, receipt)

		// deploy contract
		deployJar, err := DecompressFromJar("../hvmtestfile/fibonacci/fibonacci-1.0-fibonacci.jar")
		assert.Nil(t, err)
		transaction = NewTransaction(acH.GetAddress().Hex()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
		transaction.Sign(acH)
		receipt, err = rpc.DeployContract(transaction)
		assert.Nil(t, err)
		assert.NotNil(t, receipt)

		// invoke contract
		abiPath := "../hvmtestfile/fibonacci/hvm.abi"
		abiJson, rerr := common.ReadFileAsString(abiPath)
		assert.Nil(t, rerr)
		abi, gerr := hvm.GenAbi(abiJson)
		assert.Nil(t, gerr)
		easyBean := "invoke.InvokeFibonacci"
		beanAbi, err := abi.GetBeanAbi(easyBean)
		assert.Nil(t, err)
		payload, err := hvm.GenPayload(beanAbi)
		assert.Nil(t, err)
		transaction1 := NewTransaction(acH.GetAddress().Hex()).Invoke(receipt.ContractAddress, payload).VMType(HVM)
		transaction1.Sign(acH)
		_, err = rpc.InvokeContract(transaction1)
		assert.Nil(t, err)

		// manager contract by vote
		ope := bvm.NewContractDeployContractOperation([]byte("source"), deployJar, "hvm", nil)
		contractOpt := bvm.NewProposalCreateOperationForContract(ope)
		payload = bvm.EncodeOperation(contractOpt)
		tx := NewTransaction(acH.GetAddress().Hex()).Invoke(contractOpt.Address(), payload).VMType(BVM)
		tx.Sign(acH)
		_, err = rpc.ManageContractByVote(tx)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "ManagerContractByVote] is not enable"))
	})

	t.Run("auth_more_methods_and_revoke", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = ["","user"]
		//id = 1
		//methods = ["contract_*", "block_*","archive_*", "sub_*", "node_*", "cert_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"contract_*", "block_*", "archive_*", "sub_*", "node_*", "cert_*"},
				Name:            "node",
			},
		}
		stdError := rpc.SetRulesInNode(rules)
		assert.Nil(t, stdError)

		// add user,admin to account D
		accountD, _ := account.NewAccount(defaultPwd)
		acD, _ := account.NewAccountFromAccountJSON(accountD, defaultPwd)
		stdError = rpc.AddRoleForNode(acD.GetAddress().Hex(), user, admin)
		assert.Nil(t, stdError)

		rpc.SetAccount(acD)
		// get contract_*
		_, err := rpc.GetContractCountByAddr(acD.GetAddress().Hex())
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", acD.GetAddress().Hex())))

		// get block_*
		_, err = rpc.GetChainHeight()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", acD.GetAddress().Hex())))

		// get archive_*
		_, err = rpc.ListSnapshot()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", acD.GetAddress().Hex())))

		// get node_*
		_, err = rpc.GetNodes()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", acD.GetAddress().Hex())))

		// sub_*, cert_* has not get api

		// delete user from account
		// add user,admin to account D
		stdError = rpc.DeleteRoleFromNode(acD.GetAddress().Hex(), user)
		assert.Nil(t, stdError)

		rpc.SetAccount(acD)
		// get contract_*
		_, err = rpc.GetContractCountByAddr(acD.GetAddress().Hex())
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "Account dose not exist or account balance is 0"))

		// get block_*
		_, err = rpc.GetChainHeight()
		assert.Nil(t, err)

		// get archive_*
		_, err = rpc.ListSnapshot()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "The process of snapshot or archive happened error"))

		// get node_*
		_, err = rpc.GetNodes()
		assert.Nil(t, err)
	})

}

func TestRPC_InspectorNotEnable(t *testing.T) {
	t.Skip()
	// set rules
	//[inspector]
	//enable = true
	//[[inspector.rules]]
	//allow_anyone = false
	//authorized_roles = ["admin"]
	//forbidden_roles = ["user"]
	//id = 0
	//methods = ["node_*"]
	//name = "node"
	user := "user"
	admin := "admin"
	rules := []*InspectorRule{
		{
			AllowAnyone:     false,
			AuthorizedRoles: []string{admin},
			ForbiddenRoles:  []string{user},
			ID:              0,
			Method:          []string{"node_*"},
			Name:            "node",
		},
	}
	stdError := rpc.SetRulesInNode(rules)
	assert.Nil(t, stdError)

	// add user, admin to account
	accountD, _ := account.NewAccount("")
	acD, _ := account.NewAccountFromAccountJSON(accountD, "")
	stdError = rpc.AddRoleForNode(acD.GetAddress().Hex(), user, admin)
	assert.Nil(t, stdError)

	rpc.SetAccount(acD)
	states, err := rpc.GetNodeStates()
	assert.Nil(t, err)
	assert.NotNil(t, states)
}
