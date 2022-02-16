package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFromFile(t *testing.T) {
	cc, err := New()
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, cc.GetNamespace(), "global")
}
