package pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewConfigWithPath(t *testing.T) {

	s := NewConfigWithPath("../../conf")
	assert.Equal(t, time.Duration(5)*time.Second, s.dailTimeout)
}
