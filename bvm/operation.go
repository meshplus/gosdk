package bvm

import (
	"encoding/hex"
	"encoding/json"
	"github.com/meshplus/gosdk/common/bvmcom"
	"strconv"
)

// ContractMethod bvm contract method name
type ContractMethod string

const (
	// ContractDeployContract
	ContractDeployContract ContractMethod = "DeployContract"
	// ContractUpgradeContract
	ContractUpgradeContract ContractMethod = "UpgradeContract"
	// ContractMaintainContract
	ContractMaintainContract ContractMethod = "MaintainContract"

	// CNSSetCName cns contract method `SetCName`
	CNSSetCName ContractMethod = "SetCName"

	// NodeAddNode node contract method `AddNode`
	NodeAddNode ContractMethod = "AddNode"
	// NodeAddVP node contract method `AddVP`
	NodeAddVP ContractMethod = "AddVP"
	// NodeRemoveVP node contract method `RemoveVP`
	NodeRemoveVP ContractMethod = "RemoveVP"

	// PermissionCreateRole permission contract method `CreateRole`
	PermissionCreateRole ContractMethod = "CreateRole"
	// PermissionDeleteRole permission contract method `DeleteRole`
	PermissionDeleteRole ContractMethod = "DeleteRole"
	// PermissionGrant permission contract method `Grant`
	PermissionGrant ContractMethod = "Grant"
	// PermissionRevoke permission contract method `Revoke`
	PermissionRevoke ContractMethod = "Revoke"

	// CASetCAMode ca contract method `SetCAMode`
	CASetCAMode ContractMethod = "SetCAMode"
	// CAGetCAMode ca contract method `GetCAMode`
	CAGetCAMode ContractMethod = "GetCAMode"

	// ProposalCreate proposal contract method `Create`
	ProposalCreate ContractMethod = "Create"
	// ProposalDirect proposal contract method `Direct`
	ProposalDirect ContractMethod = "Direct"
	// ProposalVote proposal contract method `Vote`
	ProposalVote ContractMethod = "Vote"
	// ProposalGrant proposal contract method `Cancel`
	ProposalCancel ContractMethod = "Cancel"
	// ProposalExecute proposal contract method `Execute`
	ProposalExecute ContractMethod = "Execute"

	// HashSet hash contract method `Set`
	HashSet ContractMethod = "Set"
	// HashGet hash contract method `Get`
	HashGet ContractMethod = "Get"

	// AccountRegister account contract method `Register`
	AccountRegister ContractMethod = "Register"
	// AccountAbandon account contract method `Logout`
	AccountAbandon ContractMethod = "Abandon"

	SetchainID ContractMethod = "SetChainID"

	//CertRevoke cert contract method `CertRevoke`
	CertRevoke ContractMethod = "CertRevoke"
	//CertCheck cert contract method `CheckRevoke`
	CertCheck ContractMethod = "CertCheck"

	//CertFreeze cert contract method `CertFreeze`
	CertFreeze ContractMethod = "CertFreeze"
	//CertUnfreeze cert contract method `CertUnfreeze`
	CertUnfreeze ContractMethod = "CertUnfreeze"

	SRSInfo    ContractMethod = "GetSRSInfo"
	SRSBeacon  ContractMethod = "Beacon"
	SRSHistory ContractMethod = "GetHistory"

	AddRootCA  ContractMethod = "AddRootCA"
	GetRootCAs ContractMethod = "GetRootCAs"

	ChangeHashAlgo     ContractMethod = "ChangeHashAlgo"
	GetHashAlgo        ContractMethod = "GetHashAlgo"
	GetSupportHashAlgo ContractMethod = "GetSupportHashAlgo"

	RegisterAnchor   ContractMethod = "RegisterAnchorNode"
	UnRegisterAnchor ContractMethod = "UnregisterAnchorNode"
	ReplaceAnchor    ContractMethod = "ReplaceAnchorNode"
	ReadCrossChain   ContractMethod = "ReadCrossChainTransaction"
	ReadAnchor       ContractMethod = "ReadAnchorStatus"
	Timeout          ContractMethod = "Timeout"

	accountAddress      = "0x0000000000000000000000000000000000ffff04"
	proposalAddress     = "0x0000000000000000000000000000000000ffff02"
	hashAddress         = "0x0000000000000000000000000000000000ffff01"
	certAddress         = "0x0000000000000000000000000000000000ffff05"
	didAddress          = "0x0000000000000000000000000000000000ffff06"
	NormalAnchorAddress = "0x0000000000000000000000000000000000ffff08"
	mpcAddress          = "0x0000000000000000000000000000000000ffff09"
	SystemAnchorAddress = "0x0000000000000000000000000000000000ffff0a"
	rootCAAddress       = "0x0000000000000000000000000000000000ffff0b"
	hashAlgoAddress     = "0x0000000000000000000000000000000000ffff0d"

	// GenesisInfoKey the key for store genesis info in state db.
	genesisInfoKey = "the_key_for_genesis_info"
)

// ProposalType proposal type
type ProposalType int

const (
	// ProposalTypeConfig proposal of config type
	ProposalTypeConfig ProposalType = iota
	// ProposalTypePermission proposal of permission type
	ProposalTypePermission
	// ProposalTypeNode proposal of node type
	ProposalTypeNode
	// ProposalTypeCNS proposal of cns type
	ProposalTypeCNS
	// ProposalTypeContract proposal of contract type
	ProposalTypeContract
	// ProposalTypeCA proposal of ca type
	ProposalTypeCA
)

func (pt ProposalType) String() string {
	switch pt {
	case ProposalTypeNode:
		return "NODE"
	case ProposalTypeConfig:
		return "CONFIG"
	case ProposalTypePermission:
		return "PERMISSION"
	case ProposalTypeCNS:
		return "CNS"
	case ProposalTypeContract:
		return "CONTRACT"
	case ProposalTypeCA:
		return "CA"
	default:
		return "unknown proposal type"
	}
}

// Operation define the operation for proposal
type Operation interface {
	// Method return method name
	Method() ContractMethod

	// Args return args for call method
	Args() []string

	SetMethod(ContractMethod)

	SetArgs([]string)
}

type BuiltinOperation interface {
	Operation
	Address() string
}

type ProposalOperationImpl struct {
	operationImpl
}

func (po *ProposalOperationImpl) Address() string {
	return proposalAddress
}

type HashOperationImpl struct {
	operationImpl
}

func (po *HashOperationImpl) Address() string {
	return hashAddress
}

type AccountOperationImpl struct {
	operationImpl
}

func (ao *AccountOperationImpl) Address() string {
	return accountAddress
}

type CertOperationImpl struct {
	operationImpl
}

func (po *CertOperationImpl) Address() string {
	return certAddress
}

type HashAlgoOperationImpl struct {
	operationImpl
}

func (mo *HashAlgoOperationImpl) Address() string {
	return hashAlgoAddress
}

type MPCOperationImpl struct {
	operationImpl
}

func (mo *MPCOperationImpl) Address() string {
	return mpcAddress
}

type RootCAOperationImpl struct {
	operationImpl
}

func (ro *RootCAOperationImpl) Address() string {
	return rootCAAddress
}

type ProposalContentOperation interface {
	Operation
	ProposalType()
}

type PermissionOperation interface {
	ProposalContentOperation
	PermissionType()
}

type permissionOperationImpl struct {
	operationImpl
}

func (po *permissionOperationImpl) ProposalType()   {}
func (po *permissionOperationImpl) PermissionType() {}

type CAOperation interface {
	ProposalContentOperation
	CAType()
}

type caOperationImpl struct {
	operationImpl
}

func (co *caOperationImpl) ProposalType() {}

func (co *caOperationImpl) CAType() {}

type NodeOperation interface {
	ProposalContentOperation
	NodeType()
}

type nodeOperationImpl struct {
	operationImpl
}

func (po *nodeOperationImpl) ProposalType() {}
func (po *nodeOperationImpl) NodeType()     {}

// CNSOperation cns operation
type CNSOperation interface {
	ProposalContentOperation
	CNSType()
}

type cnsOperationImpl struct {
	operationImpl
}

func (co *cnsOperationImpl) ProposalType() {}

func (co *cnsOperationImpl) CNSType() {}

type ContractOperation interface {
	ProposalContentOperation
	ContractType()
}

type contractOperationImpl struct {
	operationImpl
}

func (co *contractOperationImpl) ProposalType() {}

func (co *contractOperationImpl) ContractType() {}

//DIDOperationImpl used for set chainID
type DIDOperationImpl struct {
	operationImpl
}

//Address return did contract address
func (did *DIDOperationImpl) Address() string {
	return didAddress
}

type NormalAnchorOperationImpl struct {
	operationImpl
}

func (no *NormalAnchorOperationImpl) Address() string {
	return NormalAnchorAddress
}

type SystemAnchorOperationImpl struct {
	operationImpl
}

func (ao *SystemAnchorOperationImpl) Address() string {
	return SystemAnchorAddress
}

// NewProposalCreateOperationForContract new proposal create operation for contract operation
func NewProposalCreateOperationForContract(ops ...ContractOperation) BuiltinOperation {
	data := encodeProposalContentOperation(convertToProposalContentOperations(ops))
	return newProposalCreateOperation(data, ProposalTypeContract)
}

// NewProposalCreateOperationForCNS new proposal create operation for cns operation
func NewProposalCreateOperationForCNS(ops ...CNSOperation) BuiltinOperation {
	data := encodeProposalContentOperation(convertToProposalContentOperations(ops))
	return newProposalCreateOperation(data, ProposalTypeCNS)
}

// NewProposalCreateOperationForNode new proposal create operation for node operation
func NewProposalCreateOperationForNode(ops ...NodeOperation) BuiltinOperation {
	data := encodeProposalContentOperation(convertToProposalContentOperations(ops))
	return newProposalCreateOperation(data, ProposalTypeNode)
}

// NewProposalCreateOperationForCA new proposal create operation for ca mode operation
func NewProposalCreateOperationForCA(ops ...CAOperation) BuiltinOperation {
	data := encodeProposalContentOperation(convertToProposalContentOperations(ops))
	return newProposalCreateOperation(data, ProposalTypeCA)
}

// NewProposalDirectOperationForNode new proposal direct operation for node operation
func NewProposalDirectOperationForNode(ops ...NodeOperation) BuiltinOperation {
	data := encodeProposalContentOperation(convertToProposalContentOperations(ops))
	return newProposalDirectOperation(data, ProposalTypeNode)
}

// NewProposalDirectOperationForCA new proposal direct operation for ca mode operation
func NewProposalDirectOperationForCA(ops ...CAOperation) BuiltinOperation {
	data := encodeProposalContentOperation(convertToProposalContentOperations(ops))
	return newProposalDirectOperation(data, ProposalTypeCA)
}

// convertToProposalContentOperations convert the slice of ProposalContentOperations's drive interface to ProposalContentOperations
func convertToProposalContentOperations(item interface{}) []ProposalContentOperation {

	var operations []ProposalContentOperation

	switch ops := item.(type) {
	case []PermissionOperation:
		for _, op := range ops {
			operations = append(operations, op)
		}
	case []NodeOperation:
		for _, op := range ops {
			operations = append(operations, op)
		}
	case []ConfigOperation:
		for _, op := range ops {
			operations = append(operations, op)
		}
	case []CNSOperation:
		for _, op := range ops {
			operations = append(operations, op)
		}
	case []ContractOperation:
		for _, op := range ops {
			operations = append(operations, op)
		}
	case []CAOperation:
		for _, op := range ops {
			operations = append(operations, op)
		}
	}
	return operations
}

// NewProposalCreateOperationForPermission new proposal create operation for permission operation
func NewProposalCreateOperationForPermission(ops ...PermissionOperation) BuiltinOperation {
	data := encodeProposalContentOperation(convertToProposalContentOperations(ops))
	return newProposalCreateOperation(data, ProposalTypePermission)
}

// NewProposalCreateOperation create a new ProposalCreate operation and return
func newProposalCreateOperation(data []byte, pType ProposalType) *ProposalOperationImpl {
	return newProposalOperation(ProposalCreate, string(data), pType.String())
}

// NewProposalCreateOperation create a new ProposalCreate operation and return
func newProposalDirectOperation(data []byte, pType ProposalType) *ProposalOperationImpl {
	return newProposalOperation(ProposalDirect, string(data), pType.String())
}

func newProposalOperation(method ContractMethod, args ...string) *ProposalOperationImpl {
	return &ProposalOperationImpl{operationImpl{method: method, args: args}}
}

// NewProposalVoteOperation create a new ProposalVote operation and return
func NewProposalVoteOperation(proposalID int, vote bool) BuiltinOperation {
	return newProposalOperation(ProposalVote, strconv.Itoa(proposalID), strconv.FormatBool(vote))
}

// NewProposalCancelOperation create a new ProposalCancel operation and return
func NewProposalCancelOperation(proposalID int) BuiltinOperation {
	return newProposalOperation(ProposalCancel, strconv.Itoa(proposalID))
}

// NewProposalExecuteOperation create a new ProposalExecute operation and return
func NewProposalExecuteOperation(proposalID int) BuiltinOperation {
	return newProposalOperation(ProposalExecute, strconv.Itoa(proposalID))
}

func newPermissionOperation(method ContractMethod, args ...string) PermissionOperation {
	return &permissionOperationImpl{operationImpl{method: method, args: args}}
}

// NewPermissionCreateRoleOperation create a new PermissionCreateRole operation and return
func NewPermissionCreateRoleOperation(role string) PermissionOperation {
	return newPermissionOperation(PermissionCreateRole, role)
}

// NewPermissionDeleteRoleOperation create a new PermissionDeleteRole operation and return
func NewPermissionDeleteRoleOperation(role string) PermissionOperation {
	return newPermissionOperation(PermissionDeleteRole, role)
}

// NewPermissionDeleteRoleOperation create a new PermissionGrant operation and return
func NewPermissionGrantOperation(role string, address string) PermissionOperation {
	return newPermissionOperation(PermissionGrant, role, address)
}

// NewPermissionRevokeOperation create a new PermissionRevoke operation and return
func NewPermissionRevokeOperation(role string, address string) PermissionOperation {
	return newPermissionOperation(PermissionRevoke, role, address)
}

func newCAOperation(method ContractMethod, args ...string) CAOperation {
	return &caOperationImpl{operationImpl{method: method, args: args}}
}

func NewCASetCAModeOperation(mode string) (CAOperation, error) {
	caMode, err := bvmcom.ConvertCAMode(mode)
	if err != nil {
		return nil, err
	}
	return newCAOperation(CASetCAMode, intToString(int(caMode))), nil
}

func NewCAGetCAModeOperation() CAOperation {
	return newCAOperation(CAGetCAMode)
}

func newNodeOperation(method ContractMethod, args ...string) NodeOperation {
	return &nodeOperationImpl{operationImpl{method: method, args: args}}
}

// NewNodeAddNodeOperation create a new NodeAddNode operation and return
func NewNodeAddNodeOperation(pub []byte, hostname, role, namespace string) NodeOperation {
	return newNodeOperation(NodeAddNode, string(pub), hostname, role, namespace)
}

// NewNodeAddVPOperation create a new NodeAddVP operation and return
func NewNodeAddVPOperation(hostname, namespace string) NodeOperation {
	return newNodeOperation(NodeAddVP, hostname, namespace)
}

// NewNodeRemoveVPOperation create a new NodeRemoveVP operation and return
func NewNodeRemoveVPOperation(hostname, namespace string) NodeOperation {
	return newNodeOperation(NodeRemoveVP, hostname, namespace)
}

func newCNSOperation(methodName ContractMethod, args ...string) CNSOperation {
	return &cnsOperationImpl{operationImpl{
		method: methodName,
		args:   args,
	}}
}

// NewCNSSetCNameOperation create a new CNSSetCName operation and return
func NewCNSSetCNameOperation(address string, cnsName string) CNSOperation {
	return newCNSOperation(CNSSetCName, address, cnsName)
}

func newContractOperation(methodName ContractMethod, opt *ContractManagerOptions) ContractOperation {
	bytes, _ := json.Marshal(opt)
	return &contractOperationImpl{operationImpl{
		method: methodName,
		args:   []string{string(bytes)},
	}}
}

// NewContractDeployContractOperation create a new ContractDeployContract operation and return
func NewContractDeployContractOperation(source, bin []byte, vmType string, compileOpts map[string]string) ContractOperation {
	opt := &ContractManagerOptions{
		Source:     source,
		Bin:        bin,
		VMType:     vmType,
		CompileOpt: compileOpts,
	}
	return newContractOperation(ContractDeployContract, opt)
}

// NewContractUpgradeContractOperation create a new ContractUpgradeContract operation and return
func NewContractUpgradeContractOperation(source, bin []byte, vmType, contractAddress string, compileOpts map[string]string) ContractOperation {
	opt := &ContractManagerOptions{
		Source:     source,
		Bin:        bin,
		VMType:     vmType,
		Addr:       contractAddress,
		CompileOpt: compileOpts,
	}
	return newContractOperation(ContractUpgradeContract, opt)
}

// NewContractUpgradeOperationByName create a new ContractUpgradeContract operation by contract name and return
func NewContractUpgradeOperationByName(source, bin []byte, vmType, contractName string, compileOpts map[string]string) ContractOperation {
	opt := &ContractManagerOptions{
		Source:     source,
		Bin:        bin,
		VMType:     vmType,
		Name:       contractName,
		CompileOpt: compileOpts,
	}
	return newContractOperation(ContractUpgradeContract, opt)
}

// NewContractMaintainContractOperation create a new ContractMaintainContract operation and return
func NewContractMaintainContractOperation(contractAddress, vmType string, opcode int) ContractOperation {
	opt := &ContractManagerOptions{
		Addr:   contractAddress,
		OpCode: opcode,
		VMType: vmType,
	}
	return newContractOperation(ContractMaintainContract, opt)
}

// NewContractMaintainOperationByName create a new ContractMaintainContract operation by contract name and return
func NewContractMaintainOperationByName(contractName, vmType string, opcode int) ContractOperation {
	opt := &ContractManagerOptions{
		Name:   contractName,
		OpCode: opcode,
		VMType: vmType,
	}
	return newContractOperation(ContractMaintainContract, opt)
}

func newHashOperation(method ContractMethod, args ...string) *HashOperationImpl {
	return &HashOperationImpl{operationImpl{method: method, args: args}}
}

// NewSetGenesisInfoForHpcOperation create a new HashGenesisInfo operation for hyperchain and return
func NewSetGenesisInfoForHpcOperation(genesisInfo *GenesisInfo) BuiltinOperation {
	genesisBytes, _ := json.Marshal(genesisInfo)
	return newHashOperation(HashSet, genesisInfoKey, string(genesisBytes))
}

// NewHashSetOperation create a new HashSet operation and return
func NewHashSetOperation(key, value string) BuiltinOperation {
	return newHashOperation(HashSet, key, value)
}

// NewHashGetOperation create a new HashGet operation and return
func NewHashGetOperation(key string) BuiltinOperation {
	return newHashOperation(HashGet, key)
}

func newAccountOperation(method ContractMethod, args ...string) *AccountOperationImpl {
	return &AccountOperationImpl{operationImpl{method: method, args: args}}
}

// NewAccountRegisterOperation create a new AccountRegister operation and return
func NewAccountRegisterOperation(address string, cert []byte) BuiltinOperation {
	return newAccountOperation(AccountRegister, address, string(cert))
}

// NewAccountAbandonOperation create a new AccountAbandon operation and return
func NewAccountAbandonOperation(address string, sdkCert []byte) BuiltinOperation {
	return newAccountOperation(AccountAbandon, address, string(sdkCert))
}

func newCertOperation(method ContractMethod, args ...string) *CertOperationImpl {
	return &CertOperationImpl{operationImpl{method: method, args: args}}
}

func newMPCOperation(method ContractMethod, args ...string) *MPCOperationImpl {
	return &MPCOperationImpl{operationImpl{method: method, args: args}}
}

func newRootCAOperation(method ContractMethod, args ...string) *RootCAOperationImpl {
	return &RootCAOperationImpl{operationImpl{method: method, args: args}}
}

func newHashAlgoOperation(method ContractMethod, args ...string) *HashAlgoOperationImpl {
	return &HashAlgoOperationImpl{operationImpl{method: method, args: args}}
}

func NewMPCInfoOperation(tag, ct string) BuiltinOperation {
	return newMPCOperation(SRSInfo, tag, ct)
}

func NewMPCBeaconOperation(ptau []byte) BuiltinOperation {
	return newMPCOperation(SRSBeacon, hex.EncodeToString(ptau))
}

func NewMPCHistoryOperation(ct string) BuiltinOperation {
	return newMPCOperation(SRSHistory, ct)
}

func NewRootCAAddOperation(rootCA string) BuiltinOperation {
	return newRootCAOperation(AddRootCA, rootCA)
}

func NewRootCAGetOperation() BuiltinOperation {
	return newRootCAOperation(GetRootCAs)
}

// NewCertRevokeOperation create a new CertRevoke operation and return
func NewCertRevokeOperation(cert, priv []byte) (BuiltinOperation, error) {
	if priv == nil {
		return newCertOperation(CertRevoke, string(cert), "", ""), nil
	}
	msg := "revoke"
	sign := []byte(msg)
	if priv != nil {
		key, err := ParsePriv(priv)
		if err != nil {
			return nil, err
		}
		sign, err = key.Sign([]byte(msg))
		if err != nil {
			return nil, err
		}
	}
	return newCertOperation(CertRevoke, string(cert), msg, hex.EncodeToString(sign)), nil
}

// NewCertFreezeOperation create a new CertFreeze operation and return
func NewCertFreezeOperation(cert, priv []byte) (BuiltinOperation, error) {
	msg := "freeze"
	sign := []byte(msg)
	if priv != nil {
		key, err := ParsePriv(priv)
		if err != nil {
			return nil, err
		}
		sign, err = key.Sign([]byte(msg))
		if err != nil {
			return nil, err
		}
	}
	return newCertOperation(CertFreeze, string(cert), msg, hex.EncodeToString(sign)), nil
}

// NewCertUnfreezeOperation create a new CertUnfreeze operation and return
func NewCertUnfreezeOperation(cert, priv []byte) (BuiltinOperation, error) {
	msg := "unfreeze"
	sign := []byte(msg)
	if priv != nil {
		key, err := ParsePriv(priv)
		if err != nil {
			return nil, err
		}
		sign, err = key.Sign([]byte(msg))
		if err != nil {
			return nil, err
		}
	}
	return newCertOperation(CertUnfreeze, string(cert), msg, hex.EncodeToString(sign)), nil
}

// NewCertCheckOperation create a new CertCheck operation and return
func NewCertCheckOperation(cert []byte) BuiltinOperation {
	return newCertOperation(CertCheck, string(cert))
}

//NewDIDSetChainIDOperation create a setChainID operation
func NewDIDSetChainIDOperation(chainID string) BuiltinOperation {
	return &DIDOperationImpl{
		operationImpl{method: SetchainID, args: []string{chainID}},
	}
}

func NewNormalAnchorOperation(method ContractMethod, args []string) BuiltinOperation {
	return &NormalAnchorOperationImpl{
		operationImpl{method: method, args: args},
	}
}

func NewSystemAnchorOperation(method ContractMethod, args ...string) BuiltinOperation {
	return &SystemAnchorOperationImpl{
		operationImpl{method: method, args: args},
	}
}

func NewHashGetSupportOperation() BuiltinOperation {
	return newHashAlgoOperation(GetSupportHashAlgo)
}

func NewHashGetAlgoOperation() BuiltinOperation {
	return newHashAlgoOperation(GetHashAlgo)
}

func NewHashChangeHashAlgo(data []byte) BuiltinOperation {
	return newHashAlgoOperation(ChangeHashAlgo, string(data))
}

// operationImpl the implementation of Operation
type operationImpl struct {
	method ContractMethod
	args   []string
}

func NewOperation() Operation {
	return &operationImpl{
		method: "",
		args:   nil,
	}
}

func (o *operationImpl) Method() ContractMethod {
	return o.method
}

func (o *operationImpl) Args() []string {
	return o.args
}

func (o *operationImpl) SetArgs(params []string) {
	o.args = params
}

func (o *operationImpl) SetMethod(method ContractMethod) {
	o.method = method
}
