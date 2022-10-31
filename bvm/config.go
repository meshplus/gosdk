package bvm

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"time"
)

type ConfigOperation interface {
	ProposalContentOperation
	ConfigType()
}

type configOperationImpl struct {
	operationImpl
}

type configPath string

const (
	setFilterEnable          configPath = "SetFilterEnable"
	setFilterRules           configPath = "SetFilterRules"
	setConsensusAlgo         configPath = "SetConsensusAlgo"
	setConsensusSetSize      configPath = "SetConsensusSetSize"
	setProposalTimeout       configPath = "SetProposalTimeout"
	setProposalThreshold     configPath = "SetProposalThreshold"
	setContractVoteThreshold configPath = "SetContractVoteThreshold"
	setContractVoteEnable    configPath = "SetContractVoteEnable"
	setConsensusBatchSize    configPath = "SetConsensusBatchSize"
	setConsensusPoolSize     configPath = "SetConsensusPoolSize"
)

const (
	consensusAlgoPath         = "consensus.algo"
	consensusSetSizePath      = "consensus.set.set_size"
	consensusBatchSizePath    = "consensus.pool.batch_size"
	consensusPoolSizePath     = "consensus.pool.pool_size"
	filterEnablePath          = "filter.enable"
	filterRulesPath           = "filter.rules"
	proposalThresholdPath     = "proposal.threshold"
	proposalTimeoutPath       = "proposal.timeout"
	contractVoteThresholdPath = "proposal.contract.vote.threshold"
	contractVoteEnablePath    = "proposal.contract.vote.enable"
)

func (co *configOperationImpl) ProposalType() {}
func (co *configOperationImpl) ConfigType()   {}

func NewSetFilterEnable(b bool) ConfigOperation {
	return newConfigOperation(setFilterEnable, boolToString(b))
}

func NewSetFilterRules(rules []*NsFilterRule) ConfigOperation {
	return newConfigOperation(setFilterRules, rulesToString(rules))
}

func NewSetConsensusAlgo(algo string) ConfigOperation {
	return newConfigOperation(setConsensusAlgo, algo)
}

func NewSetConsensusSetSize(i int) ConfigOperation {
	return newConfigOperation(setConsensusSetSize, intToString(i))
}

func NewSetConsensusBatchSize(i int) ConfigOperation {
	return newConfigOperation(setConsensusBatchSize, intToString(i))
}

func NewSetConsensusPoolSize(i int) ConfigOperation {
	return newConfigOperation(setConsensusPoolSize, intToString(i))
}

func NewSetProposalTimeout(d time.Duration) ConfigOperation {
	return newConfigOperation(setProposalTimeout, intToString(int(d)))
}

func NewSetProposalThreshold(i int) ConfigOperation {
	return newConfigOperation(setProposalThreshold, intToString(i))
}

func NewSetContactVoteThreshold(i int) ConfigOperation {
	return newConfigOperation(setContractVoteThreshold, intToString(i))
}

func NewSetContactVoteEnable(b bool) ConfigOperation {
	return newConfigOperation(setContractVoteEnable, boolToString(b))
}

func newConfigOperation(method configPath, arg ...string) *configOperationImpl {
	return &configOperationImpl{operationImpl{method: ContractMethod(method), args: arg}}
}

func NewProposalCreateOperationByConfigOps(ops ...ConfigOperation) BuiltinOperation {
	data := encodeProposalContentOperation(convertToProposalContentOperations(ops))
	return newProposalCreateOperation(data, ProposalTypeConfig)
}

// NewProposalCreateOperationForConfig new proposal create operation for config operation
func NewProposalCreateOperationForConfig(config []byte) (BuiltinOperation, error) {
	// split the config file into operations
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		return nil, err
	}
	pathMap := make(map[string]bool)
	for _, k := range v.AllKeys() {
		pathMap[k] = true
	}
	var ops []ConfigOperation
	if pathMap[consensusAlgoPath] {
		ops = append(ops, NewSetConsensusAlgo(v.GetString(consensusAlgoPath)))
	}
	if pathMap[consensusPoolSizePath] {
		ops = append(ops, NewSetConsensusPoolSize(v.GetInt(consensusPoolSizePath)))
	}
	if pathMap[consensusBatchSizePath] {
		ops = append(ops, NewSetConsensusBatchSize(v.GetInt(consensusBatchSizePath)))
	}
	if pathMap[consensusSetSizePath] {
		ops = append(ops, NewSetConsensusSetSize(v.GetInt(consensusSetSizePath)))
	}
	if pathMap[filterEnablePath] {
		ops = append(ops, NewSetFilterEnable(v.GetBool(filterEnablePath)))
	}
	if pathMap[filterRulesPath] {
		m := v.Get(filterRulesPath)
		var rs []*NsFilterRule
		if mapstructure.Decode(m, &rs) == nil {
			ops = append(ops, NewSetFilterRules(rs))
		}
	}
	if pathMap[proposalThresholdPath] {
		ops = append(ops, NewSetProposalThreshold(v.GetInt(proposalThresholdPath)))
	}
	if pathMap[proposalTimeoutPath] {
		ops = append(ops, NewSetProposalTimeout(v.GetDuration(proposalTimeoutPath)))
	}
	if pathMap[contractVoteEnablePath] {
		ops = append(ops, NewSetContactVoteEnable(v.GetBool(contractVoteEnablePath)))
	}
	if pathMap[contractVoteThresholdPath] {
		ops = append(ops, NewSetContactVoteThreshold(v.GetInt(contractVoteThresholdPath)))
	}
	return NewProposalCreateOperationByConfigOps(ops...), nil
}
