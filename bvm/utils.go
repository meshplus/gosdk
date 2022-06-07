package bvm

import (
	"encoding/json"
	"strconv"
)

// NsFilterRule is the rule of tx api filter
type NsFilterRule struct {
	// AllowAnyone determines whether the resources can be accessed freely by anyone
	AllowAnyone bool `json:"allow_anyone" mapstructure:"allow_anyone"`

	// AuthorizedRoles determine who can access the resource if the resources can not be accessed freely
	AuthorizedRoles []string `json:"authorized_roles" mapstructure:"authorized_roles"`

	// ForbiddenRoles determine who can not access the resources though he has the authorized roles
	ForbiddenRoles []string `json:"forbidden_roles" mapstructure:"forbidden_roles"`

	// ID is the identity sequence number for priority
	ID int `json:"id" mapstructure:"id"`

	// Name is the identity string for reading
	Name string `json:"name" mapstructure:"name"`

	// To is  the `to` address used to define resources of tx api
	To []string `json:"to" mapstructure:"to"`

	// VM is the `vmType` used to define resources of tx api
	VM []string `json:"vm" mapstructure:"vm"`
}

// GenesisNode define the filed for genesis node.
type GenesisNode struct {
	Hostname    string `json:"genesisNode" mapstructure:"genesisNode"`
	CertContent string `json:"certContent" mapstructure:"certContent"`
}

// GenesisInfo define the filed in genesis info.
type GenesisInfo struct {
	GenesisAccount map[string]string `json:"genesisAccount,omitempty"`
	GenesisNodes   []*GenesisNode    `json:"genesisNodes,omitempty"`
}

// ContractManagerOptions contract options
// vmType must have value
// when call DeployContract, source, bin must have value yet
// when call UpgradeContract, source, bin, addr must have value yet
// when call MaintainContract, addr and opCode must have value yet
type ContractManagerOptions struct {
	VMType     string            `json:"vmType,omitempty"`
	Source     []byte            `json:"source,omitempty"`
	Bin        []byte            `json:"bin,omitempty"`
	Addr       string            `json:"addr,omitempty"`
	Name       string            `json:"name,omitempty"`
	OpCode     int               `json:"opCode,omitempty"`
	CompileOpt map[string]string `json:"compileOpt,omitempty"`
}

type AlgoSet struct {
	HashAlgo    string `json:"hash_algo"`
	EncryptAlgo string `json:"encrypt_algo"`
}

func boolToString(b bool) string {
	return strconv.FormatBool(b)
}

func intToString(i int) string {
	return strconv.Itoa(i)
}

func rulesToString(rs []*NsFilterRule) string {
	bs, _ := json.Marshal(rs)
	return string(bs)
}
