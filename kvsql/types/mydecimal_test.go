package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMyDecimal(t *testing.T) {
	d := new(MyDecimal)
	err := d.FromString([]byte("-0.999999999999999999999999999999"))
	assert.Nil(t, err)
}
