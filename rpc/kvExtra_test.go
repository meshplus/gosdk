package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKVExtra(t *testing.T) {
	kvStore := NewKVExtra()
	assert.Equal(t, kvStore.data[KVExtraVersionKey], "1.0")
}

func TestAddKV(t *testing.T) {
	kvStore := NewKVExtra()
	err := kvStore.AddKV("foo", "bar")
	assert.Equal(t, err, nil)
	assert.Equal(t, kvStore.data["foo"], "bar")
}

func TestStringify(t *testing.T) {
	stringifyJSON := "{\"__version__\":\"1.0\",\"foo\":\"bar\"}"
	kvStore := NewKVExtra()
	err := kvStore.AddKV("foo", "bar")
	assert.Equal(t, err, nil)
	assert.Equal(t, kvStore.Stringify(), stringifyJSON)
}
