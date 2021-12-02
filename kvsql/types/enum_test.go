package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnum(t *testing.T) {
	e := Enum{
		Name: "hello",
	}
	assert.Equal(t, "hello", e.String())
}
