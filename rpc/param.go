package rpc

import (
	"encoding/json"
	"github.com/meshplus/gosdk/common"
)

// QueryTxRange packaged query transaction by block number range
type QueryTxRange struct {
	From     uint64    `json:"from"`
	To       uint64    `json:"to"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

// Serialize serialize to map
func (qtr *QueryTxRange) Serialize() interface{} {
	param := make(map[string]interface{})
	param["from"] = qtr.From
	param["to"] = qtr.To
	param["metadata"] = qtr.Metadata
	return param
}

// SerializeToString serialize to string
func (qtr *QueryTxRange) SerializeToString() string {
	ret, _ := json.Marshal(qtr)
	return string(ret)
}

// mapParam packaged key-value pair param
type mapParam map[string]interface{}

// newMapParam construct the mapParam by init key-value pair
func newMapParam(key string, value interface{}) *mapParam {
	param := make(map[string]interface{})
	param[key] = value
	return (*mapParam)(&param)
}

// Serialize serialize to map
func (mp *mapParam) Serialize() interface{} {
	return mp
}

// SerializeToString serialize to string
func (mp *mapParam) SerializeToString() string {
	bytes, _ := json.Marshal(mp)

	return string(bytes)
}

func (mp *mapParam) addKV(key string, value interface{}) *mapParam {
	(*mp)[key] = value
	return mp
}

/***************** for WebSocket *****************/

// EventFilter event filter
type EventFilter interface {
	Serialize() interface{}
	GetEventType() EventType
}

// BlockEventFilter block filter
type BlockEventFilter struct {
	eventType EventType
	BlockInfo bool
}

// NewBlockEventFilter return a BlockEventFilter
func NewBlockEventFilter() *BlockEventFilter {
	return &BlockEventFilter{
		eventType: BLOCKEVENT,
	}
}

// SetBlockInfo blockInfo setter
func (bf *BlockEventFilter) SetBlockInfo(b bool) {
	bf.BlockInfo = b
}

// Serialize serialize
func (bf *BlockEventFilter) Serialize() interface{} {
	return bf.BlockInfo
}

// GetEventType return event type
func (bf *BlockEventFilter) GetEventType() EventType {
	return bf.eventType
}

// SystemStatusFilter system status filter
type SystemStatusFilter struct {
	eventType         EventType
	Modules           []string `json:"modules,omitempty"`
	ModulesExclude    []string `json:"modules_exclude,omitempty"`
	Subtypes          []string `json:"subtypes,omitempty"`
	SubtypesExclude   []string `json:"subtypes_exclude,omitempty"`
	ErrorCodes        []string `json:"error_codes,omitempty"`
	ErrorCodesExclude []string `json:"error_codes_exclude,omitempty"`
}

// NewSystemStatusFilter init SystemStatusFilter
func NewSystemStatusFilter() *SystemStatusFilter {
	return &SystemStatusFilter{
		eventType: SYSTEMSTATUSEVENT,
	}
}

// AddModules add modules into filter
func (ssf *SystemStatusFilter) AddModules(modules ...string) *SystemStatusFilter {
	ssf.Modules = append(ssf.Modules, modules...)
	return ssf
}

// AddModulesExclude add modules exclude into filter
func (ssf *SystemStatusFilter) AddModulesExclude(modulesExclude ...string) *SystemStatusFilter {
	ssf.ModulesExclude = append(ssf.ModulesExclude, modulesExclude...)
	return ssf
}

// AddSubtypes add subtype into filter
func (ssf *SystemStatusFilter) AddSubtypes(subtypes ...string) *SystemStatusFilter {
	ssf.Subtypes = append(ssf.Subtypes, subtypes...)
	return ssf
}

// AddSubtypesExclude add subtypesExclude into filter
func (ssf *SystemStatusFilter) AddSubtypesExclude(subtypesExclude ...string) *SystemStatusFilter {
	ssf.SubtypesExclude = append(ssf.SubtypesExclude, subtypesExclude...)
	return ssf
}

// AddErrorCode add error code into filter
func (ssf *SystemStatusFilter) AddErrorCode(errorCodes ...string) *SystemStatusFilter {
	ssf.ErrorCodes = append(ssf.ErrorCodes, errorCodes...)
	return ssf
}

// AddErrorCodeExclude add error code exclude into filter
func (ssf *SystemStatusFilter) AddErrorCodeExclude(errorCodesExclude ...string) *SystemStatusFilter {
	ssf.ErrorCodesExclude = append(ssf.ErrorCodesExclude, errorCodesExclude...)
	return ssf
}

// Serialize serialize
func (ssf *SystemStatusFilter) Serialize() interface{} {
	return ssf
}

// GetEventType get event type
func (ssf *SystemStatusFilter) GetEventType() EventType {
	return ssf.eventType
}

// LogsFilter logs filter
type LogsFilter struct {
	eventType EventType
	FromBlock uint64           `json:"fromBlock,omitempty"`
	ToBlock   uint64           `json:"toBlock,omitempty"`
	Addresses []string         `json:"addresses,omitempty"`
	Topics    [4][]common.Hash `json:"topics,omitempty"`
}

// NewLogsFilter init logs filter
func NewLogsFilter() *LogsFilter {
	return &LogsFilter{
		eventType: LOGSEVENT,
	}
}

// SetFromBlock add from block into filter
func (lf *LogsFilter) SetFromBlock(from uint64) *LogsFilter {
	lf.FromBlock = from
	return lf
}

// SetToBlock add to block into filter
func (lf *LogsFilter) SetToBlock(to uint64) *LogsFilter {
	lf.ToBlock = to
	return lf
}

// AddAddress add address into filter
func (lf *LogsFilter) AddAddress(addresses ...string) *LogsFilter {
	lf.Addresses = append(lf.Addresses, addresses...)
	return lf
}

// SetTopic set topic of specific position
func (lf *LogsFilter) SetTopic(pos int, topics ...common.Hash) *LogsFilter {
	lf.Topics[pos] = topics
	return lf
}

// Serialize serialize
func (lf *LogsFilter) Serialize() interface{} {
	return lf
}

// GetEventType get event type
func (lf *LogsFilter) GetEventType() EventType {
	return lf.eventType
}
