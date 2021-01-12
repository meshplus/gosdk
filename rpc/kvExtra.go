package rpc

import (
	"encoding/json"
	"errors"
	"strings"
)

// KVExtra is transaction extra field
type KVExtra struct {
	data map[string]interface{}
}

const (
	// KVExtraVersionValue is private transaction extra field version
	KVExtraVersionValue = "1.0"
	// KVExtraVersionKey is private transaction extra field name
	KVExtraVersionKey = "__version__"
	// KVExtraPrivateKey is private transaction extra payload
	KVExtraPrivateKey = "__privateData__"
	// KVExtraPrivatePayloadKey is private transaction real payload
	KVExtraPrivatePayloadKey = "payload"
	// KVExtraPrivatePubSigKey is private transaction public signature
	KVExtraPrivatePubSigKey = "publicSignature"
)

// NewKVExtra is used to private transaction
func NewKVExtra() *KVExtra {
	defaultData := map[string]interface{}{
		KVExtraVersionKey: KVExtraVersionValue,
	}
	return &KVExtra{
		data: defaultData,
	}
}

// AddKV is used to add key into KVExtra
func (k *KVExtra) AddKV(key string, value interface{}) error {
	if strings.HasPrefix(key, "__") && strings.HasSuffix(key, "__") {
		return errors.New("key with leading and tailing \"__\" is reserved")
	}
	k.add(key, value)
	return nil
}

func (k *KVExtra) add(key string, value interface{}) {
	k.data[key] = value
}

// Stringify is used to convert a KVExtra struct into a string
func (k *KVExtra) Stringify() string {
	v, err := json.Marshal(k.data)
	if err != nil {
		logger.Error("KVExtra stringify error")
	}
	return string(v)
}
