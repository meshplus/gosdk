package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/meshplus/gosdk/common"
)

type routingKey string

const (
	// MQBlock MQBlock
	MQBlock routingKey = "MQBlock"
	// MQLog MQLog
	MQLog routingKey = "MQLog"
	// MQException MQException
	MQException routingKey = "MQException"
	// MQHvm
	//MQHvm routingKey = "MQHvm"
)

// RegisterMeta mq register
type RegisterMeta struct {
	//queue related
	RoutingKeys []routingKey `json:"routingKeys,omitempty"`
	QueueName   string       `json:"queueName,omitempty"`
	//self info
	From      string `json:"from,omitempty"`
	Signature string `json:"signature,omitempty"`
	// block accounts
	IsVerbose bool `json:"isVerbose"`
	// vm log criteria
	FromBlock string           `json:"fromBlock,omitempty"`
	ToBlock   string           `json:"toBlock,omitempty"`
	Addresses []common.Address `json:"addresses,omitempty"`
	Topics    [][]common.Hash  `json:"topics,omitempty"`
	Delay     bool             `json:"delay"`
	// exception criteria
	//Modules        []string `json:"modules,omitempty"`
	//ModulesExclude []string `json:"modules_exclude,omitempty"`
	//SubType        []string `json:"subtypes,omitempty"`
	//SubTypeExclude []string `json:"subtypes_exclude,omitempty"`
	//Code           []int    `json:"error_codes,omitempty"`
	//CodeExclude    []int    `json:"error_codes_exclude,omitempty"`
}

// NewRegisterMeta create a new instance of RegisterMeta
func NewRegisterMeta(from, queueName string, keys ...routingKey) *RegisterMeta {
	//if strings.HasPrefix(from, "0x") {
	//	from = from[2:]
	//}
	return &RegisterMeta{
		From:        from,
		QueueName:   queueName,
		RoutingKeys: keys,
		Topics:      make([][]common.Hash, 0),
	}
}

// Verbose node info is verbose
func (rm *RegisterMeta) Verbose(v bool) *RegisterMeta {
	rm.IsVerbose = v
	return rm
}

// SetFromBlock set from block
func (rm *RegisterMeta) SetFromBlock(from string) *RegisterMeta {
	rm.From = from
	return rm
}

// SetToBlock set to block
func (rm *RegisterMeta) SetToBlock(to string) *RegisterMeta {
	rm.ToBlock = to
	return rm
}

// AddAddress add address
func (rm *RegisterMeta) AddAddress(address ...common.Address) *RegisterMeta {
	rm.Addresses = append(rm.Addresses, address...)
	return rm
}

// SetTopics set topic
func (rm *RegisterMeta) SetTopics(pos int, topics ...common.Hash) *RegisterMeta {
	if len(rm.Topics) >= 4 {
		fmt.Println(fmt.Errorf("you can only set 4 topics at most"))
		logger.Error("you can only set 4 topics at most")
	}
	rm.Topics = append(rm.Topics, topics)
	return rm
}

// SetDelay set delay
func (rm *RegisterMeta) SetDelay(delay bool) *RegisterMeta {
	rm.Delay = delay
	return rm
}

//// AddModules add modules
//func (rm *RegisterMeta) AddModules(modules ...string) *RegisterMeta {
//	rm.Modules = append(rm.Modules, modules...)
//	return rm
//}
//
//// AddModulesExclude add modules exclude
//func (rm *RegisterMeta) AddModulesExclude(modulesExclude ...string) *RegisterMeta {
//	rm.ModulesExclude = append(rm.ModulesExclude, modulesExclude...)
//	return rm
//}
//
//// AddSubType add subtype
//func (rm *RegisterMeta) AddSubType(subtypes ...string) *RegisterMeta {
//	rm.SubType = append(rm.SubType, subtypes...)
//	return rm
//}
//
//// AddSubTypesExclude add subtype exclude
//func (rm *RegisterMeta) AddSubTypesExclude(subtypesExclude ...string) *RegisterMeta {
//	rm.SubTypeExclude = append(rm.SubTypeExclude, subtypesExclude...)
//	return rm
//}
//
//// AddCode add code
//func (rm *RegisterMeta) AddCode(codes ...int) *RegisterMeta {
//	rm.Code = append(rm.Code, codes...)
//	return rm
//}
//
//// AddCodeExclude add code exclude
//func (rm *RegisterMeta) AddCodeExclude(codesExclude ...int) *RegisterMeta {
//	rm.CodeExclude = append(rm.CodeExclude, codesExclude...)
//	return rm
//}

// Sign sign RegisterMeta
func (rm *RegisterMeta) Sign(key interface{}) {
	sig, err := sign(key, concatNeedHash(rm), false, false)
	if err != nil {
		return
	}
	rm.Signature = sig
}

// concatNeedHash need hash string
func concatNeedHash(rm *RegisterMeta) string {
	from := strings.TrimPrefix(strings.ToLower(rm.From), "0x")
	var result bytes.Buffer
	result.WriteString("qname=" + rm.QueueName)
	result.WriteString(":routingKeys=" + arrayToString(rm.RoutingKeys))
	result.WriteString(":from=" + from)
	result.WriteString(":isVerbose=" + strconv.FormatBool(rm.IsVerbose))
	result.WriteString(":fromBlock=" + rm.FromBlock)
	result.WriteString(":toBlock=" + rm.ToBlock)
	result.WriteString(":addresses=" + arrayToString(rm.Addresses))
	result.WriteString(":topics=" + arrayToString(rm.Topics))
	result.WriteString(":delay=" + strconv.FormatBool(rm.Delay))
	//result.WriteString(":modules=" + arrayToString(rm.Modules))
	//result.WriteString(":modulesExclude=" + arrayToString(rm.ModulesExclude))
	//result.WriteString(":subType=" + arrayToString(rm.SubType))
	//result.WriteString(":subTypeExclude=" + arrayToString(rm.SubTypeExclude))
	//result.WriteString(":code=" + arrayToString(rm.Code))
	//result.WriteString(":codeExclude=" + arrayToString(rm.CodeExclude))

	return result.String()
}

// arrayToString hash util
func arrayToString(array interface{}) string {
	var result bytes.Buffer
	switch array := array.(type) {
	case []string:
		return strings.Join(array, ".")
	case []int:
		for i, val := range array {
			result.WriteString(strconv.Itoa(val))
			if i != len(array)-1 {
				result.WriteString(".")
			}
		}
	case []routingKey:
		for i, val := range array {
			result.WriteString(string(val))
			if i != len(array)-1 {
				result.WriteString(".")
			}
		}
	case []common.Address:
		for i, val := range array {
			result.WriteString(val.String())
			if i != len(array)-1 {
				result.WriteString(".")
			}
		}
	case []common.Hash: // not used
		for i, val := range array {
			result.WriteString(val.String())
			if i != len(array)-1 {
				result.WriteString(".")
			}
		}
	case [][]common.Hash:
		for i, array := range array {
			for j, item := range array {
				result.WriteString(item.Hex())
				if j != len(array)-1 {
					result.WriteString(".")
				}
			}
			if i != len(array)-1 {
				result.WriteString(",")
			}
		}
	default:
		logger.Error("not support type")
	}
	return result.String()
}

// Serialize Serialize
func (rm *RegisterMeta) Serialize() interface{} {
	if rm.Signature == "" {
		logger.Warning("this transaction is not Signature")
	}
	data, err := json.Marshal(rm)
	if err != nil {
		return nil
	}
	return data
}

// SerializeToString SerializeToString
func (rm *RegisterMeta) SerializeToString() string {
	return ""
}

// UnRegisterMeta UnRegisterMeta
type UnRegisterMeta struct {
	From         string
	QueueName    string
	ExchangeName string
	Signature    string
}

// NewUnRegisterMeta create new instance
func NewUnRegisterMeta(from, queue, exchange string) *UnRegisterMeta {
	return &UnRegisterMeta{
		From:         from,
		QueueName:    queue,
		ExchangeName: exchange,
	}
}

// Sign sign UnRegisterMeta
func (urm *UnRegisterMeta) Sign(key interface{}) {
	needHash := urm.QueueName + ":" + urm.ExchangeName
	sig, err := sign(key, needHash, false, false)
	if err != nil {
		logger.Error("ecdsa signature error")
		return
	}
	urm.Signature = sig
}
