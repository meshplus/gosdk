package rpc

import (
	"encoding/json"
	"github.com/meshplus/gosdk/abi"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/bvm"
	"github.com/meshplus/gosdk/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestRPC_Contract_Smock(t *testing.T) {
	t.Skip()
	adminCount := 6
	source, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")
	//bin := `6060604052341561000f57600080fd5b5b6104c78061001f6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680635b6beeb914610049578063e15fe02314610120575b600080fd5b341561005457600080fd5b6100a4600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061023a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e55780820151818401525b6020810190506100c9565b50505050905090810190601f1680156101125780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561012b57600080fd5b6101be600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061034f565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101ff5780820151818401525b6020810190506101e3565b50505050905090810190601f16801561022c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102426103e2565b6000826040518082805190602001908083835b60208310151561027b57805182525b602082019150602081019050602083039250610255565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103425780601f1061031757610100808354040283529160200191610342565b820191906000526020600020905b81548152906001019060200180831161032557829003601f168201915b505050505090505b919050565b6103576103e2565b816000846040518082805190602001908083835b60208310151561039157805182525b60208201915060208101905060208303925061036b565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906103d79291906103f6565b508190505b92915050565b602060405190810160405280600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061043757805160ff1916838001178555610465565b82800160010185558215610465579182015b82811115610464578251825591602001919060010190610449565b5b5090506104729190610476565b5090565b61049891905b8082111561049457600081600090555060010161047c565b5090565b905600a165627a7a723058208ac1d22e128cf8381d7ac66b4c438a6a906ccf5ee583c3a9e46d4cdf7b3f94580029`
	//abiRaw := `[{"constant":false,"inputs":[{"name":"key","type":"string"}],"name":"getHash","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"key","type":"string"},{"name":"value","type":"string"}],"name":"setHash","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

	t.Run("set_proposal.contract.vote.enable_false", func(t *testing.T) {
		setProposalContractVoteEnable(t, adminCount, false)
	})

	// proposal.contract.vote.enable = false
	t.Run("deploy_success", func(t *testing.T) {
		// use deployContract deploy contract success
		contractAddr := deploySuccess(t)

		// invoke contract
		invokeContractSuccess(t, contractAddr)

		// get proposal.contract.vote.enable value
		assertContractVote(t, false)
	})

	t.Run("deploy_by_vote_fail", func(t *testing.T) {
		ope := bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil)
		createProposalForContractFail(t, ope)
	})

	t.Run("upgrade_success", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false)
		// deploy
		contractAddr := deploySuccess(t)
		// upgrade
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(1, contractAddr, binContract).VMType(EVM)
		tx.Sign(privateKey)
		re, err := rpc.MaintainContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re)

		// contract status is normal
		assertContractStatusSuccess(t, contractAddr, `"normal"`)
	})

	t.Run("upgrade_by_vote_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false)
		// deploy
		contractAddr := deploySuccess(t)
		// upgrade fail
		ope := bvm.NewContractUpgradeContractOperation(source, common.Hex2Bytes(binContract), "evm", contractAddr, nil)
		createProposalForContractFail(t, ope)
	})

	t.Run("freeze_success", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false)
		// deploy
		contractAddr := deploySuccess(t)
		// freeze success
		maintainSuccess(t, contractAddr, 2)

		// contract status is frozen
		assertContractStatusSuccess(t, contractAddr, `"frozen"`)
	})

	t.Run("freeze_by_vote_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false)
		// deploy
		contractAddr := deploySuccess(t)
		// freeze fail
		ope := bvm.NewContractMaintainContractOperation(contractAddr, "evm", 2)
		createProposalForContractFail(t, ope)
	})

	t.Run("unfreeze_success", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false)
		// deploy
		contractAddr := deploySuccess(t)
		// freeze
		maintainSuccess(t, contractAddr, 2)
		// unfreeze
		maintainSuccess(t, contractAddr, 3)

		// contract status is normal
		assertContractStatusSuccess(t, contractAddr, `"normal"`)
	})

	t.Run("unfreeze_by_vote_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false)
		// deploy
		contractAddr := deploySuccess(t)
		// freeze
		maintainSuccess(t, contractAddr, 2)
		// unfreeze fail
		ope := bvm.NewContractMaintainContractOperation(contractAddr, "evm", 3)
		createProposalForContractFail(t, ope)
	})

	t.Run("destroy_success", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false)
		// deploy
		contractAddr := deploySuccess(t)
		// destroy_success
		maintainSuccess(t, contractAddr, 5)
		// contract status is 'destroy'
		assertContractStatusSuccess(t, contractAddr, `"destroy"`)
	})

	t.Run("destroy_by_vote_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false)
		// deploy
		contractAddr := deploySuccess(t)
		// destroy fail
		ope := bvm.NewContractMaintainContractOperation(contractAddr, "evm", 5)
		createProposalForContractFail(t, ope)
	})

	t.Run("set_proposal.contract.vote.enable_true", func(t *testing.T) {
		setProposalContractVoteEnable(t, adminCount, true)
	})

	// proposal for contract
	t.Run("create_by_contractManager", func(t *testing.T) {
		// grant contractManager to a new account
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, adminCount, bvm.NewPermissionGrantOperation("contractManager", newKey.GetAddress().Hex()))

		// use new account create proposal for contract
		createProposalForContractSuccess(t, newKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		// cancel proposal
		proposal, _ := rpc.GetProposal()
		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
	})

	t.Run("create_by_admin", func(t *testing.T) {
		// grant admin to a new account
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, adminCount, bvm.NewPermissionGrantOperation("admin", newKey.GetAddress().Hex()))

		// use new account create proposal for contract
		createProposalForContractSuccess(t, newKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		// cancel proposal
		proposal, _ := rpc.GetProposal()
		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
	})

	t.Run("create_by_normal", func(t *testing.T) {
		newKey := genNewAccountKey(t)

		// use new account create proposal for contract
		createProposalForContractSuccess(t, newKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		// cancel proposal
		proposal, _ := rpc.GetProposal()
		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
	})

	t.Run("vote_by_contractManager", func(t *testing.T) {
		createProposalForContractSuccess(t, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		proposal, _ := rpc.GetProposal()
		key, _ := account.NewAccountFromAccountJSON("", password)
		invokeProposalContractSuccess(bvm.NewProposalVoteOperation(int(proposal.ID), true), key, t)

		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	t.Run("vote_by_normal", func(t *testing.T) {
		createProposalForContractSuccess(t, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		newKey := genNewAccountKey(t)
		t.Log(newKey.GetAddress().Hex())
		proposal, _ := rpc.GetProposal()
		invokeProposalContractFail(bvm.NewProposalVoteOperation(int(proposal.ID), true), newKey, t)

		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	t.Run("execute_by_creator", func(t *testing.T) {
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		assertContractStatusSuccess(t, res[0].Msg, `"normal"`)
	})

	t.Run("execute_by_other", func(t *testing.T) {
		createProposalForContractSuccess(t, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rpc.GetProposal()
		voteProposalByAdminCount(int(proposal.ID), adminCount, t, 0)
		newKey := genNewAccountKey(t)
		invokeProposalContractFail(bvm.NewProposalExecuteOperation(int(proposal.ID)), newKey, t)

		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	t.Run("cancel_when_voting_by_creator", func(t *testing.T) {
		createProposalForContractSuccess(t, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rpc.GetProposal()
		assert.Equal(t, "VOTING", proposal.Status)
		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
		proposal, _ = rpc.GetProposal()
		assert.Equal(t, "CANCEL", proposal.Status)
	})

	t.Run("cancel_when_waitingExe_by_creator", func(t *testing.T) {
		createProposalForContractSuccess(t, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rpc.GetProposal()
		voteProposalByAdminCount(int(proposal.ID), adminCount, t, 0)
		proposal, _ = rpc.GetProposal()
		assert.Equal(t, "WAITING_EXE", proposal.Status)
		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
		proposal, _ = rpc.GetProposal()
		assert.Equal(t, "CANCEL", proposal.Status)
	})

	t.Run("cancel_when_voting_by_other", func(t *testing.T) {
		createProposalForContractSuccess(t, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rpc.GetProposal()
		assert.Equal(t, "VOTING", proposal.Status)
		newKey := genNewAccountKey(t)
		invokeProposalContractFail(bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
		proposal, _ = rpc.GetProposal()
		assert.Equal(t, "VOTING", proposal.Status)

		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	t.Run("cancel_when_waitingExe_by_other", func(t *testing.T) {
		createProposalForContractSuccess(t, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rpc.GetProposal()
		voteProposalByAdminCount(int(proposal.ID), adminCount, t, 0)
		proposal, _ = rpc.GetProposal()
		assert.Equal(t, "WAITING_EXE", proposal.Status)
		newKey := genNewAccountKey(t)
		invokeProposalContractFail(bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
		proposal, _ = rpc.GetProposal()
		assert.Equal(t, "WAITING_EXE", proposal.Status)

		invokeProposalContractSuccess(bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	// proposal for permission
	t.Run("create_contractManager_role", func(t *testing.T) {
		key, _ := account.NewAccountFromAccountJSON("", password)
		invokeProposalContractFail(bvm.NewProposalCreateOperationForPermission(bvm.NewPermissionCreateRoleOperation("contractManager")), key, t)
	})

	t.Run("delete_admin_role", func(t *testing.T) {
		key, _ := account.NewAccountFromAccountJSON("", password)
		invokeProposalContractFail(bvm.NewProposalCreateOperationForPermission(bvm.NewPermissionDeleteRoleOperation("admin")), key, t)
	})

	t.Run("delete_contractManager_role", func(t *testing.T) {
		key, _ := account.NewAccountFromAccountJSON("", password)
		invokeProposalContractFail(bvm.NewProposalCreateOperationForPermission(bvm.NewPermissionDeleteRoleOperation("contractManager")), key, t)
	})

	t.Run("grant_admin_to_normal", func(t *testing.T) {
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, adminCount, bvm.NewPermissionGrantOperation("admin", newKey.GetAddress().Hex()))
		roles, err := rpc.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, "admin", roles[0])
	})

	t.Run("grant_contractManager_to_normal", func(t *testing.T) {
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, adminCount, bvm.NewPermissionGrantOperation("contractManager", newKey.GetAddress().Hex()))
		roles, err := rpc.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, "contractManager", roles[0])
	})

	t.Run("revoke_contractManager", func(t *testing.T) {
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, adminCount, bvm.NewPermissionGrantOperation("contractManager", newKey.GetAddress().Hex()))
		roles, err := rpc.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, "contractManager", roles[0])

		completePermissionProposal(t, adminCount, bvm.NewPermissionRevokeOperation("contractManager", newKey.GetAddress().Hex()))
		roles, err = rpc.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 0)
	})

	t.Run("revoke_admin", func(t *testing.T) {
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, adminCount, bvm.NewPermissionGrantOperation("admin", newKey.GetAddress().Hex()))
		roles, err := rpc.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, "admin", roles[0])

		completePermissionProposal(t, adminCount, bvm.NewPermissionRevokeOperation("admin", newKey.GetAddress().Hex()))
		roles, err = rpc.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 0)
	})

	// proposal.contract.vote.enable = true
	t.Run("deploy_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy by vote success
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg
		// invoke
		invokeContractSuccess(t, contractAddr)
	})

	t.Run("deploy_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Deploy(binContract).VMType(EVM)
		tx.Sign(privateKey)
		_, err := rpc.DeployContract(tx)
		assert.Error(t, err)
	})

	t.Run("upgrade_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// update by vote success
		res = manageContractByVote(t, privateKey, adminCount, bvm.NewContractUpgradeContractOperation(source, common.Hex2Bytes(binContract), "evm", contractAddr, nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
	})

	t.Run("upgrade_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// update fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(1, contractAddr, binContract).VMType(EVM)
		tx.Sign(privateKey)
		_, err := rpc.MaintainContract(tx)
		assert.Error(t, err)
	})

	t.Run("freeze_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// freeze by vote success
		res = manageContractByVote(t, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 2))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		assertContractStatusSuccess(t, contractAddr, `"frozen"`)
	})

	t.Run("freeze_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// freeze fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(2, contractAddr, "").VMType(EVM)
		tx.Sign(privateKey)
		_, err := rpc.MaintainContract(tx)
		assert.Error(t, err)

		assertContractStatusSuccess(t, contractAddr, `"normal"`)

	})

	t.Run("unfreeze_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// freeze
		res = manageContractByVote(t, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 2))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		assertContractStatusSuccess(t, contractAddr, `"frozen"`)
		// unfreeze by vote success
		res = manageContractByVote(t, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 3))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		assertContractStatusSuccess(t, contractAddr, `"normal"`)

	})

	t.Run("unfreeze_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// freeze
		res = manageContractByVote(t, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 2))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		assertContractStatusSuccess(t, contractAddr, `"frozen"`)
		// unfreeze fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(3, contractAddr, "").VMType(EVM)
		tx.Sign(privateKey)
		_, err := rpc.MaintainContract(tx)
		assert.Error(t, err)

		assertContractStatusSuccess(t, contractAddr, `"frozen"`)

	})

	t.Run("destroy_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		assertContractStatusSuccess(t, contractAddr, `"normal"`)
		// destroy by vote success
		res = manageContractByVote(t, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 5))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		assertContractStatusSuccess(t, contractAddr, `"destroy"`)
	})

	t.Run("destroy_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true)

		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		assertContractStatusSuccess(t, contractAddr, `"normal"`)
		// unfreeze fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(5, contractAddr, "").VMType(EVM)
		tx.Sign(privateKey)
		_, err := rpc.MaintainContract(tx)
		assert.Error(t, err)
		assertContractStatusSuccess(t, contractAddr, `"normal"`)

	})

}

func TestRPC_ThresholdAndVP(t *testing.T) {
	t.Skip()
	t.Run("init_node", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON("", password)
		nodes := []string{"node1", "node2", "node3", "node4"}
		ns := "global"
		role := "vp"
		pub := []byte("pub1")
		var ops []bvm.NodeOperation
		for _, n := range nodes {
			ops = append(ops, bvm.NewNodeAddNodeOperation(pub, n, role, ns))
			ops = append(ops, bvm.NewNodeAddVPOperation(n, ns))
		}
		res := completeProposal(t, creator, 6, bvm.NewProposalCreateOperationForNode(ops...))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
	})

	t.Run("remove_vp", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON("", password)
		res := completeProposal(t, creator, 6, bvm.NewProposalCreateOperationForNode(bvm.NewNodeRemoveVPOperation("node1", "global")))
		assert.NotEqual(t, bvm.SuccessCode, res[0].Code)
		t.Log(res[0].Msg)
	})

	t.Run("SetProposalThreshold", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON("", password)
		invokeProposalContractFail(bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetProposalThreshold(7)), creator, t)
	})

	t.Run("revoke", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON("", password)
		ops := []bvm.PermissionOperation{
			bvm.NewPermissionRevokeOperation("admin", creator.GetAddress().Hex()),
			bvm.NewPermissionRevokeOperation("contractManager", creator.GetAddress().Hex()),
		}
		res := completeProposal(t, creator, 6, bvm.NewProposalCreateOperationForPermission(ops...))
		assert.Len(t, res, 2)
		//assert.NotEqual(t, bvm.SuccessCode,res[0].Code)
		//assert.NotEqual(t, bvm.SuccessCode,res[1].Code)
	})

}

func TestRPC_ContractByName(t *testing.T) {
	t.Skip()
	adminCount := 6
	source, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")

	t.Run("set_proposal.contract.vote.enable_false", func(t *testing.T) {
		setProposalContractVoteEnable(t, adminCount, false)
	})

	t.Run("maintain_by_name_success", func(t *testing.T) {
		contractAddr := deploySuccess(t)
		ri := rand.Int()
		contractName := "name" + strconv.Itoa(ri)
		setCName(contractAddr, contractName, t)

		// freeze success
		tx := NewTransaction(privateKey.GetAddress().Hex()).MaintainByName(2, contractName, "").VMType(EVM)
		tx.Sign(privateKey)
		re, err := rpc.MaintainContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re)

		status, stdError := rpc.GetContractStatus(contractAddr)
		assert.Nil(t, stdError)
		statu, stdError := rpc.GetContractStatusByName(contractName)
		assert.Nil(t, stdError)
		assert.Equal(t, status, statu)
		assert.Equal(t, `"frozen"`, status)
	})

	t.Run("set_proposal.contract.vote.enable_true", func(t *testing.T) {
		setProposalContractVoteEnable(t, adminCount, true)
	})

	t.Run("maintain_by_name_by_vote_success", func(t *testing.T) {
		// deploy
		res := manageContractByVote(t, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg
		ri := rand.Int()
		contractName := "name" + strconv.Itoa(ri)
		setCName(contractAddr, contractName, t)

		// update by vote success
		res = manageContractByVote(t, privateKey, adminCount, bvm.NewContractUpgradeOperationByName(source, common.Hex2Bytes(binContract), "evm", contractName, nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		// freeze by vote success
		res = manageContractByVote(t, privateKey, adminCount, bvm.NewContractMaintainOperationByName(contractName, "evm", 2))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		status, stdError := rpc.GetContractStatus(contractAddr)
		assert.Nil(t, stdError)
		statu, stdError := rpc.GetContractStatusByName(contractName)
		assert.Nil(t, stdError)
		assert.Equal(t, status, statu)
		assert.Equal(t, `"frozen"`, status)
	})
}

func completeProposal(t *testing.T, creatorKey account.Key, adminCount int, opt bvm.BuiltinOperation) []*bvm.OpResult {
	invokeProposalContractSuccess(opt, creatorKey, t)
	proposal, _ := rpc.GetProposal()
	voteProposalByAdminCount(int(proposal.ID), adminCount, t, 1)
	ret := invokeProposalContractSuccess(bvm.NewProposalExecuteOperation(int(proposal.ID)), creatorKey, t)
	var res []*bvm.OpResult
	_ = json.Unmarshal(ret, &res)
	return res
}

func manageContractByVote(t *testing.T, creatorKey account.Key, adminCount int, ops ...bvm.ContractOperation) []*bvm.OpResult {
	createProposalForContractSuccess(t, creatorKey, ops...)
	proposal, _ := rpc.GetProposal()
	voteProposalByAdminCount(int(proposal.ID), adminCount, t, 0)
	ret := invokeProposalContractSuccess(bvm.NewProposalExecuteOperation(int(proposal.ID)), creatorKey, t)
	var res []*bvm.OpResult
	_ = json.Unmarshal(ret, &res)
	assert.Len(t, res, len(ops))
	return res
}

func invokeContractSuccess(t *testing.T, contractAddr string) {
	ABI, er := abi.JSON(strings.NewReader(abiContract))
	assert.Nil(t, er)
	packed, er := ABI.Pack("add", uint32(1), uint32(2))
	assert.Nil(t, er)
	tx := NewTransaction(privateKey.GetAddress().Hex()).Invoke(contractAddr, packed)
	tx.Sign(privateKey)
	_, err := rpc.InvokeContract(tx)
	assert.Nil(t, err)
}

func completePermissionProposal(t *testing.T, adminCount int, ops ...bvm.PermissionOperation) {
	key, _ := account.NewAccountFromAccountJSON("", password)
	invokeProposalContractSuccess(bvm.NewProposalCreateOperationForPermission(ops...), key, t)
	proposal, _ := rpc.GetProposal()
	voteProposalByAdminCount(int(proposal.ID), adminCount, t, 1)
	invokeProposalContractSuccess(bvm.NewProposalExecuteOperation(int(proposal.ID)), key, t)
}

func genNewAccountKey(t *testing.T) account.Key {
	pwd := "12347890"
	newAccount, err := account.NewAccount(pwd)
	assert.Nil(t, err)
	newKey, err := account.NewAccountFromAccountJSON(newAccount, pwd)
	assert.Nil(t, err)
	return newKey
}

func setProposalContractVoteEnable(t *testing.T, adminCount int, enable bool) {
	// create proposal
	key, _ := account.NewAccountFromAccountJSON("", password)
	invokeProposalContractSuccess(bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetContactVoteEnable(enable)), key, t)

	proposal, _ := rpc.GetProposal()
	// vote
	voteProposalByAdminCount(int(proposal.ID), adminCount, t, 1)
	// execute
	invokeProposalContractSuccess(bvm.NewProposalExecuteOperation(int(proposal.ID)), key, t)
}

func voteProposalByAdminCount(pid, adminCount int, t *testing.T, startCount int) {
	// vote
	for i := startCount; i < adminCount; i++ {
		var (
			key account.Key
		)
		if i < 4 {
			key, _ = account.NewAccountFromAccountJSON("", pwd)
		} else {
			key, _ = account.NewAccountSm2FromAccountJSON("", pwd)
		}
		operation := bvm.NewProposalVoteOperation(pid, true)
		invokeProposalContractSuccess(operation, key, t)
	}
}

func invokeBVMSuccess(operation bvm.BuiltinOperation, key account.Key, t *testing.T) *bvm.Result {
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rpc.InvokeContract(tx)
	assert.Nil(t, err)
	assert.NotNil(t, re)
	result := bvm.Decode(re.Ret)
	return result
}

func invokeProposalContractSuccess(operation bvm.BuiltinOperation, key account.Key, t *testing.T) []byte {
	result := invokeBVMSuccess(operation, key, t)
	assert.True(t, result.Success)
	t.Log(result)
	return []byte(result.Ret)
}

func invokeProposalContractFail(operation bvm.BuiltinOperation, key account.Key, t *testing.T) {
	result := invokeBVMSuccess(operation, key, t)
	assert.False(t, result.Success)
	t.Log(result)
}

func assertContractStatusSuccess(t *testing.T, contractAddr string, expect string) {
	status, err := rpc.GetContractStatus(contractAddr)
	assert.Nil(t, err)
	assert.Equal(t, expect, status)
}

func createProposalForContractFail(t *testing.T, ops ...bvm.ContractOperation) {
	_, err := createProposalForManagerContract(t, privateKey, ops...)
	assert.Error(t, err)
}

func createProposalForContractSuccess(t *testing.T, key account.Key, ops ...bvm.ContractOperation) {
	re, err := createProposalForManagerContract(t, key, ops...)
	assert.Nil(t, err)
	assert.NotNil(t, re)
	result := bvm.Decode(re.Ret)
	assert.True(t, result.Success)
	t.Log(result)
}

func createProposalForManagerContract(t *testing.T, key account.Key, ops ...bvm.ContractOperation) (*TxReceipt, StdError) {
	contractOpt := bvm.NewProposalCreateOperationForContract(ops...)
	payload := bvm.EncodeOperation(contractOpt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(contractOpt.Address(), payload).VMType(BVM)
	tx.Sign(key)
	return rpc.ManageContractByVote(tx)
}

func maintainSuccess(t *testing.T, contractAddr string, op int64) {
	tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(op, contractAddr, "").VMType(EVM)
	tx.Sign(privateKey)
	re, err := rpc.MaintainContract(tx)
	assert.Nil(t, err)
	assert.NotNil(t, re)
}

func deploySuccess(t *testing.T) string {
	tx := NewTransaction(privateKey.GetAddress().Hex()).Deploy(binContract).VMType(EVM)
	tx.Sign(privateKey)
	re, err := rpc.DeployContract(tx)
	assert.Nil(t, err)
	assert.NotNil(t, re)
	t.Log(re.ContractAddress)
	return re.ContractAddress
}

func assertContractVote(t *testing.T, expect bool) {
	config, err := rpc.GetConfig()
	assert.Nil(t, err)
	v := viper.New()
	v.SetConfigType("toml")
	er := v.ReadConfig(strings.NewReader(config))
	assert.Nil(t, er)
	assert.Equal(t, expect, v.GetBool("proposal.contract.vote.enable"))
}

func TestInt(t *testing.T) {
	bigThreshold := new(big.Int)
	_ = json.Unmarshal(nil, bigThreshold)
	t.Log(bigThreshold.Int64())
}
